package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var version = "0.1.0-dev"
var commit = ""

const devLogPath = ".sealion/log/dev.jsonl"

const helpText = `Sealion
Containerized full-stack apps with React, C, and Postgres.

global actions:
  sealion help
    show this help
  sealion version
    print the installed version
  sealion upgrade
    upgrade the installed CLI from GitHub when a newer commit is available

features:
  install the CLI from GitHub
  # curl -fsSL <github-install-url> | bash
  curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/sealion/main/install.sh | bash

  create a new project directory
  # sealion new <project-name>
  sealion new demo

  initialize the current empty directory
  # sealion init
  mkdir demo && cd demo && sealion init

  run the local development stack
  # sealion run dev
  cd demo && sealion run dev

  inspect structured dev logs
  # sealion logs [service <name>] [containing <text>] [limit <count>] [json]
  sealion logs service backend
  sealion logs containing "/api/login" json
`

type app struct {
	home   string
	stdout io.Writer
	stderr io.Writer
}

type composeCommand struct {
	name    string
	base    []string
	help    string
	logHelp string
}

type renderer struct {
	out    io.Writer
	styled bool
}

type outputRow struct {
	key   string
	value string
}

type runningProcess struct {
	name string
	cmd  *exec.Cmd
}

type processResult struct {
	name string
	err  error
}

type structuredLogEntry struct {
	Time    string `json:"ts"`
	Source  string `json:"source"`
	Stream  string `json:"stream"`
	Service string `json:"service"`
	Message string `json:"message"`
}

type devLogSink struct {
	mu      sync.Mutex
	file    *os.File
	encoder *json.Encoder
}

type logQuery struct {
	service  string
	contains string
	limit    int
	json     bool
}

func main() {
	home, err := resolveHome()
	if err != nil {
		renderError(os.Stderr, err)
		os.Exit(1)
	}

	a := app{
		home:   home,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := a.run(os.Args[1:]); err != nil {
		renderError(os.Stderr, err)
		os.Exit(1)
	}
}

func (a app) run(args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		if len(args) != 1 {
			return errors.New("usage: sealion help")
		}
		a.printHelp()
		return nil
	case "version":
		if len(args) != 1 {
			return errors.New("usage: sealion version")
		}
		return a.commandVersion()
	case "upgrade":
		if len(args) != 1 {
			return errors.New("usage: sealion upgrade")
		}
		return a.commandUpgrade()
	case "new":
		if len(args) != 2 {
			return errors.New("usage: sealion new <project-name>")
		}
		return a.commandNew(args[1])
	case "init":
		if len(args) != 1 {
			return errors.New("usage: sealion init")
		}
		return a.commandInit()
	case "run":
		if len(args) == 2 && args[1] == "dev" {
			return a.commandRunDev()
		}
		return errors.New("usage: sealion run dev")
	case "logs":
		return a.commandLogs(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a app) printHelp() {
	if shouldStyleOutput(a.stdout) {
		fmt.Fprintf(a.stdout, "\033[38;5;245m%s\033[0m", helpText)
		return
	}
	fmt.Fprint(a.stdout, helpText)
}

func (a app) commandVersion() error {
	r := newRenderer(a.stdout)
	if commit != "" {
		r.Title("Sealion", "installed CLI")
		r.Rows(
			outputRow{"version", version},
			outputRow{"commit", commit},
		)
		return nil
	} else if head := gitShortHead(a.home); head != "" {
		r.Title("Sealion", "installed CLI")
		r.Rows(
			outputRow{"version", version},
			outputRow{"commit", head},
		)
		return nil
	}
	r.Title("Sealion", "installed CLI")
	r.Rows(outputRow{"version", version})
	return nil
}

func (a app) commandNew(name string) error {
	if err := ensureProjectName(name); err != nil {
		return err
	}

	target, err := filepath.Abs(filepath.Join(".", name))
	if err != nil {
		return err
	}
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("%s already exists", name)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	slug := projectSlug(name)
	if slug == "" {
		slug = "sealion-app"
	}
	if err := a.copyTemplate(target, name, slug); err != nil {
		return err
	}

	newRenderer(a.stdout).Message(
		"Sealion",
		"project created",
		outputRow{"path", target},
		outputRow{"next", fmt.Sprintf("cd %s", name)},
		outputRow{"", "sealion run dev"},
	)
	return nil
}

func (a app) commandInit() error {
	empty, err := isCurrentDirEmpty()
	if err != nil {
		return err
	}
	if !empty {
		return errors.New("sealion init requires an empty directory")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	name := filepath.Base(pwd)
	if err := ensureProjectName(name); err != nil {
		return err
	}

	slug := projectSlug(name)
	if slug == "" {
		slug = "sealion-app"
	}
	if err := a.copyTemplate(pwd, name, slug); err != nil {
		return err
	}

	newRenderer(a.stdout).Message(
		"Sealion",
		"project initialized",
		outputRow{"path", pwd},
		outputRow{"next", "sealion run dev"},
	)
	return nil
}

func (a app) commandUpgrade() error {
	if isDir(filepath.Join(a.home, ".git")) {
		if _, err := exec.LookPath("git"); err != nil {
			return errors.New("git is required to upgrade this installation")
		}

		status, err := commandOutput(a.home, "git", "status", "--porcelain")
		if err != nil {
			return err
		}
		if strings.TrimSpace(status) != "" {
			return fmt.Errorf("cannot upgrade because %s has local changes", a.home)
		}

		currentHead, err := commandOutput(a.home, "git", "rev-parse", "--short", "HEAD")
		if err != nil {
			return err
		}
		if _, err := commandOutput(a.home, "git", "fetch", "--quiet", "origin", "main"); err != nil {
			return err
		}
		remoteHead, err := commandOutput(a.home, "git", "rev-parse", "--short", "origin/main")
		if err != nil {
			return err
		}
		if currentHead == remoteHead {
			newRenderer(a.stdout).Message(
				"Sealion upgrade",
				"installed CLI",
				outputRow{"status", "up to date"},
				outputRow{"commit", currentHead},
			)
			return nil
		}
		if _, err := commandOutput(a.home, "git", "pull", "--ff-only", "--quiet", "origin", "main"); err != nil {
			return err
		}
		newHead, err := commandOutput(a.home, "git", "rev-parse", "--short", "HEAD")
		if err != nil {
			return err
		}
		if err := buildInstalledBinary(a.home); err != nil {
			return err
		}
		newRenderer(a.stdout).Message(
			"Sealion upgrade",
			"installed CLI",
			outputRow{"status", "upgraded"},
			outputRow{"from", currentHead},
			outputRow{"to", newHead},
		)
		return nil
	}

	installScript := filepath.Join(a.home, "install.sh")
	if !isFile(installScript) {
		return errors.New("cannot find install.sh for this Sealion installation")
	}
	cmd := exec.Command("bash", installScript)
	cmd.Env = append(os.Environ(), "SEALION_HOME="+a.home)
	cmd.Stdin = os.Stdin
	cmd.Stdout = a.stdout
	cmd.Stderr = a.stderr
	return cmd.Run()
}

func (a app) commandRunDev() error {
	if !isFile("sealion.toml") {
		return errors.New("run this inside a Sealion project")
	}

	compose, err := findCompose()
	if err != nil {
		return err
	}

	requestedPort := os.Getenv("SEALION_HTTP_PORT")
	port, err := chooseDevPort(requestedPort)
	if err != nil {
		return err
	}

	env := setEnv(os.Environ(), "SEALION_HTTP_PORT", strconv.Itoa(port))
	env = setEnv(env, "COMPOSE_MENU", "false")
	watch := compose.supports("--watch")
	logSink, err := openDevLogSink(devLogPath)
	if err != nil {
		return err
	}
	defer logSink.Close()

	r := newRenderer(a.stdout)
	a.printDevHeader(r, port)
	logSink.Write("sealion", "lifecycle", "cli", "starting containers")
	if err := composeUpDetached(compose, env); err != nil {
		logSink.Write("sealion", "lifecycle", "cli", err.Error())
		return err
	}
	logSink.Write("sealion", "lifecycle", "cli", "ready")
	r.Section("Logs", "live container output")

	return a.runDevStreams(compose, env, watch, logSink)
}

func (a app) printDevHeader(r renderer, port int) {
	r.Title("Sealion dev", "local stack")
	r.Rows(
		outputRow{"app", fmt.Sprintf("http://localhost:%d", port)},
		outputRow{"api", fmt.Sprintf("http://localhost:%d/api", port)},
	)
}

func (a app) runDevStreams(compose composeCommand, env []string, watch bool, logSink *devLogSink) error {
	var streams sync.WaitGroup
	results := make(chan processResult, 3)
	processes := make([]runningProcess, 0, 2)

	logProcess, err := a.startComposeStream(
		"logs",
		compose,
		env,
		composeLogsArgs(compose),
		func(input io.Reader, r renderer, sink *devLogSink, stream string, wg *sync.WaitGroup) {
			streamLogOutput(input, r, sink, stream, wg)
		},
		logSink,
		&streams,
		results,
	)
	if err != nil {
		return err
	}
	processes = append(processes, logProcess)

	if watch {
		watchProcess, err := a.startComposeStream(
			"watch",
			compose,
			env,
			[]string{"watch", "--no-up", "--quiet"},
			func(input io.Reader, r renderer, sink *devLogSink, stream string, wg *sync.WaitGroup) {
				streamWatchOutput(input, r, sink, stream, wg)
			},
			logSink,
			&streams,
			results,
		)
		if err != nil {
			stopProcesses(processes, syscall.SIGTERM)
			waitForProcesses(len(processes), processes, results, 5*time.Second)
			streams.Wait()
			return err
		}
		processes = append(processes, watchProcess)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	var first processResult
	interrupted := false

	select {
	case sig := <-signals:
		interrupted = true
		logSink.Write("sealion", "lifecycle", "cli", "stopping containers")
		stopProcesses(processes, sig)
	case first = <-results:
		stopProcesses(processes, syscall.SIGTERM)
	}

	alreadyReported := 0
	if !interrupted {
		alreadyReported = 1
	}
	waitForProcesses(len(processes)-alreadyReported, processes, results, 5*time.Second)
	streams.Wait()

	downErr := composeDown(compose, env)
	if interrupted {
		return downErr
	}
	if first.err != nil {
		if downErr != nil {
			return fmt.Errorf("Docker Compose %s failed: %v; cleanup failed: %w", first.name, first.err, downErr)
		}
		return fmt.Errorf("Docker Compose %s failed: %w", first.name, first.err)
	}
	return downErr
}

func (a app) startComposeStream(
	name string,
	compose composeCommand,
	env []string,
	args []string,
	stream func(io.Reader, renderer, *devLogSink, string, *sync.WaitGroup),
	logSink *devLogSink,
	streams *sync.WaitGroup,
	results chan<- processResult,
) (runningProcess, error) {
	cmd := exec.Command(compose.name, compose.args(args...)...)
	cmd.Env = env
	cmd.Stdin = os.Stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return runningProcess{}, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return runningProcess{}, err
	}

	if err := cmd.Start(); err != nil {
		return runningProcess{}, fmt.Errorf("Docker Compose %s failed to start: %w", name, err)
	}

	streams.Add(2)
	go stream(stdout, newRenderer(a.stdout), logSink, "stdout", streams)
	go stream(stderr, newRenderer(a.stderr), logSink, "stderr", streams)

	go func() {
		results <- processResult{name: name, err: cmd.Wait()}
	}()

	return runningProcess{name: name, cmd: cmd}, nil
}

func stopProcesses(processes []runningProcess, sig os.Signal) {
	for _, process := range processes {
		if process.cmd.Process != nil {
			_ = process.cmd.Process.Signal(sig)
		}
	}
}

func waitForProcesses(remaining int, processes []runningProcess, results <-chan processResult, timeout time.Duration) {
	deadline := time.After(timeout)
	for remaining > 0 {
		select {
		case <-results:
			remaining--
		case <-deadline:
			for _, process := range processes {
				if process.cmd.Process != nil {
					_ = process.cmd.Process.Kill()
				}
			}
			for remaining > 0 {
				<-results
				remaining--
			}
		}
	}
}

func newRenderer(out io.Writer) renderer {
	return renderer{out: out, styled: shouldStyleOutput(out)}
}

func renderError(out io.Writer, err error) {
	newRenderer(out).Message(
		"Sealion",
		"command failed",
		outputRow{"error", err.Error()},
		outputRow{"help", "sealion help"},
	)
}

func (r renderer) Message(title string, subtitle string, rows ...outputRow) {
	r.Title(title, subtitle)
	r.Rows(rows...)
}

func (r renderer) Title(title string, subtitle string) {
	if r.styled {
		fmt.Fprintf(r.out, "%s\n", r.paint("1;38;5;81", title))
		if subtitle != "" {
			fmt.Fprintf(r.out, "%s\n", r.paint("2;38;5;245", subtitle))
		}
	} else {
		fmt.Fprintln(r.out, title)
		if subtitle != "" {
			fmt.Fprintln(r.out, subtitle)
		}
	}
	fmt.Fprintln(r.out)
}

func (r renderer) Section(title string, subtitle string) {
	fmt.Fprintln(r.out)
	if r.styled {
		fmt.Fprintf(r.out, "%s\n", r.paint("1;38;5;245", title))
		if subtitle != "" {
			fmt.Fprintf(r.out, "%s\n", r.paint("2;38;5;245", subtitle))
		}
	} else {
		fmt.Fprintln(r.out, title)
		if subtitle != "" {
			fmt.Fprintln(r.out, subtitle)
		}
	}
	fmt.Fprintln(r.out)
}

func (r renderer) Rows(rows ...outputRow) {
	width := rowKeyWidth(rows)
	for _, row := range rows {
		r.writeRow(row, width)
	}
}

func (r renderer) Row(row outputRow) {
	r.writeRow(row, len(row.key))
}

func (r renderer) Blank() {
	fmt.Fprintln(r.out)
}

func (r renderer) writeRow(row outputRow, width int) {
	lines := strings.Split(row.value, "\n")
	if len(lines) > 1 {
		r.writeSingleLine(outputRow{row.key, lines[0]}, width)
		for _, line := range lines[1:] {
			r.writeSingleLine(outputRow{"", line}, width)
		}
		return
	}
	r.writeSingleLine(row, width)
}

func (r renderer) writeSingleLine(row outputRow, width int) {
	if row.key == "" {
		fmt.Fprintf(r.out, "%*s  %s\n", width, "", r.formatValue(row))
		return
	}
	key := r.formatKey(row.key)
	if r.styled {
		fmt.Fprintf(r.out, "%s%s  %s\n", key, strings.Repeat(" ", width-len(row.key)), r.formatValue(row))
		return
	}
	fmt.Fprintf(r.out, "%-*s  %s\n", width, row.key, row.value)
}

func (r renderer) Log(service string, message string) {
	label := service
	if label == "" {
		label = "log"
	}
	width := 9
	if len(label) > width {
		width = len(label)
	}
	if r.styled {
		fmt.Fprintf(r.out, "%s%s  %s\n", r.formatService(label), strings.Repeat(" ", width-len(label)), message)
		return
	}
	fmt.Fprintf(r.out, "%-*s  %s\n", width, label, message)
}

func (r renderer) formatKey(key string) string {
	return r.paint("2;38;5;245", key)
}

func (r renderer) formatService(service string) string {
	switch service {
	case "frontend":
		return r.paint("38;5;81", service)
	case "backend":
		return r.paint("38;5;114", service)
	case "db":
		return r.paint("38;5;222", service)
	case "watch":
		return r.paint("38;5;147", service)
	default:
		return r.paint("38;5;245", service)
	}
}

func (r renderer) formatValue(row outputRow) string {
	if !r.styled {
		return row.value
	}
	switch row.key {
	case "app", "api":
		return r.paint("38;5;81", row.value)
	case "error":
		return r.paint("38;5;203", row.value)
	default:
		return row.value
	}
}

func (r renderer) paint(code string, value string) string {
	if !r.styled {
		return value
	}
	return "\033[" + code + "m" + value + "\033[0m"
}

func rowKeyWidth(rows []outputRow) int {
	width := 0
	for _, row := range rows {
		if len(row.key) > width {
			width = len(row.key)
		}
	}
	return width
}

func streamWatchOutput(input io.Reader, r renderer, logSink *devLogSink, stream string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line == "Watch enabled" {
			continue
		}
		logSink.Write("compose-watch", stream, "watch", line)
		r.Log("watch", line)
	}
}

func streamLogOutput(input io.Reader, r renderer, logSink *devLogSink, stream string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		service, message := parseComposeLogLine(line)
		logSink.Write("compose-log", stream, service, message)
		r.Log(service, message)
	}
}

func openDevLogSink(path string) (*devLogSink, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return &devLogSink{file: file, encoder: json.NewEncoder(file)}, nil
}

func (s *devLogSink) Close() error {
	if s == nil || s.file == nil {
		return nil
	}
	return s.file.Close()
}

func (s *devLogSink) Write(source string, stream string, service string, message string) {
	if s == nil || s.encoder == nil || strings.TrimSpace(message) == "" {
		return
	}
	entry := structuredLogEntry{
		Time:    time.Now().UTC().Format(time.RFC3339Nano),
		Source:  source,
		Stream:  stream,
		Service: service,
		Message: message,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.encoder.Encode(entry)
}

func (a app) commandLogs(args []string) error {
	if !isFile("sealion.toml") {
		return errors.New("run this inside a Sealion project")
	}
	query, err := parseLogQuery(args)
	if err != nil {
		return err
	}
	entries, err := readStructuredLogEntries(devLogPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no dev logs found; run sealion run dev first")
		}
		return err
	}
	entries = filterLogEntries(entries, query)
	entries = limitLogEntries(entries, query.limit)

	if query.json {
		encoder := json.NewEncoder(a.stdout)
		for _, entry := range entries {
			if err := encoder.Encode(entry); err != nil {
				return err
			}
		}
		return nil
	}

	r := newRenderer(a.stdout)
	for _, entry := range entries {
		r.Log(entry.Service, entry.Message)
	}
	return nil
}

func parseLogQuery(args []string) (logQuery, error) {
	query := logQuery{limit: 80}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "service":
			i++
			if i >= len(args) || args[i] == "" {
				return query, errors.New("usage: sealion logs [service <name>] [containing <text>] [limit <count>] [json]")
			}
			query.service = args[i]
		case "containing":
			i++
			if i >= len(args) || args[i] == "" {
				return query, errors.New("usage: sealion logs [service <name>] [containing <text>] [limit <count>] [json]")
			}
			query.contains = args[i]
		case "limit":
			i++
			if i >= len(args) {
				return query, errors.New("usage: sealion logs [service <name>] [containing <text>] [limit <count>] [json]")
			}
			limit, err := strconv.Atoi(args[i])
			if err != nil || limit < 1 {
				return query, errors.New("log limit must be a positive number")
			}
			query.limit = limit
		case "json":
			query.json = true
		default:
			return query, fmt.Errorf("unknown logs option: %s", args[i])
		}
	}
	return query, nil
}

func readStructuredLogEntries(path string) ([]structuredLogEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []structuredLogEntry
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry structuredLogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("invalid structured log line %d: %w", lineNumber, err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func filterLogEntries(entries []structuredLogEntry, query logQuery) []structuredLogEntry {
	filtered := make([]structuredLogEntry, 0, len(entries))
	contains := strings.ToLower(query.contains)
	for _, entry := range entries {
		if query.service != "" && entry.Service != query.service {
			continue
		}
		if contains != "" && !strings.Contains(strings.ToLower(entry.Message), contains) {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func limitLogEntries(entries []structuredLogEntry, limit int) []structuredLogEntry {
	if limit < 1 || len(entries) <= limit {
		return entries
	}
	return entries[len(entries)-limit:]
}

func resolveHome() (string, error) {
	if home := os.Getenv("SEALION_HOME"); home != "" {
		return filepath.Abs(home)
	}

	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(exe)
	switch filepath.Base(dir) {
	case "bin", ".bin":
		return filepath.Dir(dir), nil
	default:
		return filepath.Dir(dir), nil
	}
}

func shouldStyleOutput(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func ensureProjectName(name string) error {
	if name == "" || strings.HasPrefix(name, ".") || strings.ContainsAny(name, `/\`) {
		return errors.New("project name must be a simple directory name")
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9._-]+$`, name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("project name may contain only letters, numbers, dots, underscores, and dashes")
	}
	return nil
}

func projectSlug(input string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(input) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func (a app) copyTemplate(target string, name string, slug string) error {
	template := filepath.Join(a.home, "templates", "default")
	if !isDir(template) {
		return fmt.Errorf("missing template: %s", template)
	}

	return filepath.WalkDir(template, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(template, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return os.MkdirAll(target, 0755)
		}

		dest := filepath.Join(target, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return os.MkdirAll(dest, info.Mode().Perm())
		}
		if entry.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, dest)
		}
		if !entry.Type().IsRegular() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content = bytes.ReplaceAll(content, []byte("__PROJECT_NAME__"), []byte(name))
		content = bytes.ReplaceAll(content, []byte("__PROJECT_SLUG__"), []byte(slug))
		return os.WriteFile(dest, content, info.Mode().Perm())
	})
}

func isCurrentDirEmpty() (bool, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func findCompose() (composeCommand, error) {
	if _, err := commandOutput("", "docker", "compose", "version"); err == nil {
		help, _ := commandOutput("", "docker", "compose", "up", "--help")
		logHelp, _ := commandOutput("", "docker", "compose", "logs", "--help")
		return composeCommand{name: "docker", base: []string{"compose"}, help: help, logHelp: logHelp}, nil
	}
	if _, err := commandOutput("", "docker-compose", "version"); err == nil {
		help, _ := commandOutput("", "docker-compose", "up", "--help")
		logHelp, _ := commandOutput("", "docker-compose", "logs", "--help")
		return composeCommand{name: "docker-compose", help: help, logHelp: logHelp}, nil
	}
	return composeCommand{}, errors.New("Docker Compose is required for sealion run dev")
}

func (c composeCommand) args(extra ...string) []string {
	args := make([]string, 0, len(c.base)+len(extra))
	args = append(args, c.base...)
	args = append(args, extra...)
	return args
}

func (c composeCommand) supports(option string) bool {
	return strings.Contains(c.help, option)
}

func (c composeCommand) logsSupports(option string) bool {
	return strings.Contains(c.logHelp, option)
}

func composeLogsArgs(compose composeCommand) []string {
	args := []string{"logs", "-f", "--tail", "80"}
	if compose.logsSupports("--no-color") {
		args = append(args, "--no-color")
	}
	return args
}

func parseComposeLogLine(line string) (string, string) {
	line = stripANSI(strings.TrimSpace(line))
	parts := strings.SplitN(line, "|", 2)
	if len(parts) != 2 {
		return "log", line
	}
	service := normalizeServiceName(strings.TrimSpace(parts[0]))
	message := strings.TrimSpace(parts[1])
	return service, message
}

func normalizeServiceName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "log"
	}
	parts := strings.Split(name, "-")
	if len(parts) > 1 {
		last := parts[len(parts)-1]
		if _, err := strconv.Atoi(last); err == nil {
			name = strings.Join(parts[:len(parts)-1], "-")
		}
	}
	if idx := strings.LastIndex(name, "-"); idx >= 0 {
		name = name[idx+1:]
	}
	switch name {
	case "frontend", "backend", "db":
		return name
	default:
		return name
	}
}

func stripANSI(value string) string {
	var out strings.Builder
	inEscape := false
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if inEscape {
			if ch >= '@' && ch <= '~' {
				inEscape = false
			}
			continue
		}
		if ch == 0x1b {
			inEscape = true
			continue
		}
		out.WriteByte(ch)
	}
	return out.String()
}

func composeUpDetached(compose composeCommand, env []string) error {
	args := []string{"up", "-d", "--build", "--remove-orphans"}
	if compose.supports("--quiet-build") {
		args = append(args, "--quiet-build")
	}
	if compose.supports("--quiet-pull") {
		args = append(args, "--quiet-pull")
	}
	if compose.supports("--wait") {
		args = append(args, "--wait", "--wait-timeout", "120")
	}
	output, err := runComposeCaptured(compose, env, args...)
	if err != nil {
		return fmt.Errorf("Docker Compose start failed: %w\n%s", err, strings.TrimSpace(output))
	}
	return nil
}

func composeDown(compose composeCommand, env []string) error {
	output, err := runComposeCaptured(compose, env, "down", "--remove-orphans")
	if err != nil {
		return fmt.Errorf("Docker Compose cleanup failed: %w\n%s", err, strings.TrimSpace(output))
	}
	return nil
}

func runComposeCaptured(compose composeCommand, env []string, args ...string) (string, error) {
	cmd := exec.Command(compose.name, compose.args(args...)...)
	cmd.Env = env
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	return output.String(), err
}

func validatePort(value string) (int, error) {
	if value == "" {
		return 0, errors.New("SEALION_HTTP_PORT must be a number from 1 to 65535")
	}
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return 0, errors.New("SEALION_HTTP_PORT must be a number from 1 to 65535")
	}
	return port, nil
}

func chooseDevPort(requested string) (int, error) {
	if requested != "" {
		port, err := validatePort(requested)
		if err != nil {
			return 0, err
		}
		if !portIsAvailable(port) {
			return 0, fmt.Errorf("port %d is already in use; choose another with SEALION_HTTP_PORT=<port> sealion run dev", port)
		}
		return port, nil
	}

	for _, port := range []int{8080, 8081, 8082, 8083, 8084, 8085, 18080, 18081, 18082, 18083, 18084, 18085} {
		if portIsAvailable(port) {
			return port, nil
		}
	}
	return 0, errors.New("no free dev port found; run with SEALION_HTTP_PORT=<port> sealion run dev")
}

func portIsAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return false
	}
	_ = listener.Close()
	return true
}

func buildInstalledBinary(home string) error {
	if _, err := exec.LookPath("go"); err != nil {
		return errors.New("Go is required to build the Sealion CLI")
	}

	outDir := filepath.Join(home, ".bin")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	finalPath := filepath.Join(outDir, "sealion")
	tmpPath := filepath.Join(outDir, fmt.Sprintf(".sealion-%d", os.Getpid()))
	ldflags := "-X main.commit=" + gitShortHead(home)

	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", tmpPath, "./cmd/sealion")
	cmd.Dir = home
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("Go build failed: %w\n%s", err, strings.TrimSpace(output.String()))
	}
	if err := os.Chmod(tmpPath, 0755); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

func commandOutput(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		text := strings.TrimSpace(output.String())
		if text != "" {
			return "", fmt.Errorf("%s %s failed: %w\n%s", name, strings.Join(args, " "), err, text)
		}
		return "", err
	}
	return strings.TrimSpace(output.String()), nil
}

func gitShortHead(dir string) string {
	head, err := commandOutput(dir, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return ""
	}
	return head
}

func setEnv(env []string, key string, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	set := false
	for _, item := range env {
		if strings.HasPrefix(item, prefix) {
			out = append(out, prefix+value)
			set = true
			continue
		}
		out = append(out, item)
	}
	if !set {
		out = append(out, prefix+value)
	}
	return out
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
