package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func (a app) commandDeploy(target string) error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}
	if err := ensureDeployTarget(target); err != nil {
		return err
	}
	deploy, found, err := loadDeployTarget(target)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("no checked-in deploy target named %s", target)
	}
	if err := validateDeployTarget(deploy); err != nil {
		return err
	}

	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	scriptPath, err := deployScriptPath(root, deploy)
	if err != nil {
		return err
	}

	r := newRenderer(a.stdout)
	r.Title("Carbide deploy", deploy.Name)
	r.Rows(
		outputRow{"script", deploy.Script},
		outputRow{"root", root},
	)
	if deploy.Description != "" {
		r.Row(outputRow{"about", deploy.Description})
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = root
	cmd.Stdout = a.stdout
	cmd.Stderr = a.stderr
	cmd.Env = append(os.Environ(), deployScriptEnv(root, deploy)...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("deploy target %s failed: %w", deploy.Name, err)
	}

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
			if upgradeBinaryNeedsRebuild(currentHead) {
				if err := buildInstalledBinary(a.home); err != nil {
					return err
				}
				newRenderer(a.stdout).Message(
					"Carbide upgrade",
					"installed CLI",
					outputRow{"status", "refreshed"},
					outputRow{"commit", currentHead},
				)
				return nil
			}
			newRenderer(a.stdout).Message(
				"Carbide upgrade",
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
			"Carbide upgrade",
			"installed CLI",
			outputRow{"status", "upgraded"},
			outputRow{"from", currentHead},
			outputRow{"to", newHead},
		)
		return nil
	}

	installScript := filepath.Join(a.home, "cli/install.sh")
	if !isFile(installScript) {
		return errors.New("cannot find cli/install.sh for this Carbide installation")
	}
	cmd := exec.Command("bash", installScript)
	cmd.Env = append(os.Environ(), "CARBIDE_HOME="+a.home)
	cmd.Stdin = os.Stdin
	cmd.Stdout = a.stdout
	cmd.Stderr = a.stderr
	return cmd.Run()
}

func upgradeBinaryNeedsRebuild(currentHead string) bool {
	currentHead = strings.TrimSpace(currentHead)
	if currentHead == "" {
		return false
	}
	return strings.TrimSpace(commit) != currentHead
}

func (a app) commandRunDev() error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	compose, err := findCompose()
	if err != nil {
		return err
	}

	requestedPort := os.Getenv("CARBIDE_HTTP_PORT")
	port, err := chooseDevPort(requestedPort)
	if err != nil {
		return err
	}

	env := setEnv(os.Environ(), "CARBIDE_HTTP_PORT", strconv.Itoa(port))
	env = setEnv(env, "COMPOSE_MENU", "false")
	env = composeEnv(env)
	watch := compose.supports("--watch")
	logSink, err := openDevLogSink(devLogPath)
	if err != nil {
		return err
	}
	defer logSink.Close()

	r := newRenderer(a.stdout)
	a.printDevHeader(r, port)
	logSink.Write("carbide", "lifecycle", "cli", "starting containers")
	services := composeServices(compose, env)
	if err := r.RunServiceProgress(
		services,
		func() map[string]composeServiceStatus {
			return composeServiceStatuses(compose, env)
		},
		func() error {
			return composeUpDetached(compose, env)
		},
	); err != nil {
		logSink.Write("carbide", "lifecycle", "cli", err.Error())
		return err
	}
	logSink.Write("carbide", "lifecycle", "cli", "ready")
	r.Section("Logs", "live container output")

	return a.runDevStreams(compose, env, watch, logSink)
}

func (a app) commandStopDev() error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	compose, err := findCompose()
	if err != nil {
		return err
	}

	env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
	env = composeEnv(env)
	services := composeServices(compose, env)
	r := newRenderer(a.stdout)
	r.Title("Carbide stop dev", "local stack")

	logSink, _ := openAppendDevLogSink(devLogPath)
	if logSink != nil {
		defer logSink.Close()
		logSink.Write("carbide", "lifecycle", "cli", "stopping containers")
	}

	if err := r.RunServiceStopProgress(
		services,
		func() map[string]composeServiceStatus {
			return composeServiceStatuses(compose, env)
		},
		func() error {
			return composeDown(compose, env)
		},
	); err != nil {
		if logSink != nil {
			logSink.Write("carbide", "lifecycle", "cli", err.Error())
		}
		return err
	}
	if logSink != nil {
		logSink.Write("carbide", "lifecycle", "cli", "stopped containers")
	}
	r.Rows(outputRow{"dev", "stopped"})
	return nil
}

func (a app) commandCleanDev() error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	compose, err := findCompose()
	if err != nil {
		return err
	}

	env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
	env = composeEnv(env)
	services := composeServices(compose, env)
	r := newRenderer(a.stdout)
	r.Title("Carbide clean dev", "local stack")

	snapshots, snapshotErr := composeServiceSnapshots(compose, env)
	if snapshotErr == nil && len(snapshots) == 0 {
		r.Rows(
			outputRow{"dev", "already clean"},
			outputRow{"next", "carbide run dev"},
		)
		return nil
	}

	logSink, _ := openAppendDevLogSink(devLogPath)
	if logSink != nil {
		defer logSink.Close()
		logSink.Write("carbide", "lifecycle", "cli", "cleaning dev state")
	}

	if err := r.RunServiceStopProgress(
		services,
		func() map[string]composeServiceStatus {
			return composeServiceStatuses(compose, env)
		},
		func() error {
			return composeDown(compose, env)
		},
	); err != nil {
		if logSink != nil {
			logSink.Write("carbide", "lifecycle", "cli", err.Error())
		}
		return err
	}
	if logSink != nil {
		logSink.Write("carbide", "lifecycle", "cli", "dev state clean")
	}
	r.Rows(
		outputRow{"dev", "clean"},
		outputRow{"next", "carbide run dev"},
	)
	return nil
}

func (a app) commandStatus() error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	compose, err := findCompose()
	if err != nil {
		return err
	}

	env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
	env = composeEnv(env)
	services := composeServices(compose, env)
	snapshots, err := composeServiceSnapshots(compose, env)
	if err != nil {
		return err
	}

	seen := map[string]bool{}
	rows := make([]tableRow, 0, len(services))
	for _, service := range services {
		snapshot, ok := snapshots[service]
		if !ok {
			rows = append(rows, tableRow{service, "-", "-", "-", "not running"})
			continue
		}
		seen[service] = true
		rows = append(rows, composeStatusRow(snapshot))
	}
	for service, snapshot := range snapshots {
		if !seen[service] {
			rows = append(rows, composeStatusRow(snapshot))
		}
	}

	r := newRenderer(a.stdout)
	r.Title("Carbide status", "local stack")
	r.Table(
		[]string{"service", "container", "ports", "internal", "status"},
		rows,
	)
	return nil
}

func (a app) commandStatusJSON() error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	compose, err := findCompose()
	if err != nil {
		return err
	}

	env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
	env = composeEnv(env)
	services := composeServices(compose, env)
	snapshots, err := composeServiceSnapshots(compose, env)
	if err != nil {
		return err
	}

	report := a.baseJSONReport("status")
	report.OK = true
	report.Services = serviceJSONReports(services, snapshots)
	if urls, err := localURLsFromCompose(compose, env); err == nil {
		report.URLs = &urls
	}
	return a.writeJSON(report)
}

func (a app) commandURLs(asJSON bool) error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	urls, err := localURLs()
	if err != nil {
		return err
	}
	if asJSON {
		report := a.baseJSONReport("urls")
		report.OK = true
		report.URLs = &urls
		return a.writeJSON(report)
	}

	r := newRenderer(a.stdout)
	r.Title("Carbide urls", "local stack")
	r.Rows(
		outputRow{"app", urls.App},
		outputRow{"api", urls.API},
		outputRow{"source", urls.Source},
	)
	return nil
}

func (a app) printDevHeader(r renderer, port int) {
	r.Title("Carbide dev", "local stack")
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
		logSink.Write("carbide", "lifecycle", "cli", "detached from dev logs")
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

	if interrupted {
		r := newRenderer(a.stdout)
		r.Blank()
		r.Message(
			"Carbide dev",
			"detached from logs; containers are still running",
			outputRow{"status", "carbide status"},
			outputRow{"logs", "carbide follow logs"},
			outputRow{"clean", "carbide clean dev"},
		)
		return nil
	}
	if first.err != nil {
		return fmt.Errorf("Docker Compose %s failed: %w", first.name, first.err)
	}
	return nil
}

func serviceJSONReports(services []string, snapshots map[string]composeServiceSnapshot) []serviceJSONReport {
	seen := map[string]bool{}
	reports := make([]serviceJSONReport, 0, len(services)+len(snapshots))
	for _, service := range services {
		snapshot, ok := snapshots[service]
		if !ok {
			reports = append(reports, serviceJSONReport{
				Service:        service,
				Container:      "-",
				PublishedPorts: nil,
				InternalPorts:  nil,
				State:          "not running",
				Status:         "not running",
			})
			continue
		}
		seen[service] = true
		reports = append(reports, serviceJSONReportFromSnapshot(snapshot))
	}
	for service, snapshot := range snapshots {
		if !seen[service] {
			reports = append(reports, serviceJSONReportFromSnapshot(snapshot))
		}
	}
	return reports
}

func serviceJSONReportFromSnapshot(snapshot composeServiceSnapshot) serviceJSONReport {
	return serviceJSONReport{
		Service:        snapshot.Service,
		Container:      statusValue(snapshot.Name),
		PublishedPorts: composePublishedPortValues(snapshot),
		InternalPorts:  composeInternalPortValues(snapshot),
		State:          statusValue(snapshot.State),
		Health:         strings.TrimSpace(snapshot.Health),
		Status:         composeServiceStatusText(snapshot),
	}
}

func localURLs() (urlsJSONReport, error) {
	compose, err := findCompose()
	if err == nil {
		env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
		env = composeEnv(env)
		if urls, err := localURLsFromCompose(compose, env); err == nil {
			return urls, nil
		}
	}
	port := runtimePortFromEnv()
	return urlsJSONReport{
		App:    fmt.Sprintf("http://localhost:%d", port),
		API:    fmt.Sprintf("http://localhost:%d/api", port),
		Source: "CARBIDE_HTTP_PORT or default dev port",
	}, nil
}

func localURLsFromCompose(compose composeCommand, env []string) (urlsJSONReport, error) {
	port := 0
	if snapshots, err := composeServiceSnapshots(compose, env); err == nil {
		if web, ok := snapshots["web"]; ok {
			port = publishedPortFromSnapshot(web)
		}
	}
	source := "running containers"
	if port == 0 {
		port = runtimePortFromEnv()
		source = "CARBIDE_HTTP_PORT or default dev port"
	}
	return urlsJSONReport{
		App:    fmt.Sprintf("http://localhost:%d", port),
		API:    fmt.Sprintf("http://localhost:%d/api", port),
		Source: source,
	}, nil
}

func publishedPortFromSnapshot(snapshot composeServiceSnapshot) int {
	for _, publisher := range snapshot.Publishers {
		if publisher.PublishedPort > 0 && publisher.TargetPort == 8080 {
			return publisher.PublishedPort
		}
	}
	for _, publisher := range snapshot.Publishers {
		if publisher.PublishedPort > 0 {
			return publisher.PublishedPort
		}
	}
	return 0
}

func (a app) baseJSONReport(command string) jsonCommandReport {
	return jsonCommandReport{
		OK:      true,
		Command: command,
		Version: version,
		Commit:  displayCommit(a.home),
		Project: currentProjectJSON(),
	}
}

func (a app) writeJSON(report jsonCommandReport) error {
	encoder := json.NewEncoder(a.stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func currentProjectJSON() *projectJSONReport {
	if !isFile(projectConfigPath) {
		return nil
	}
	metadata, err := readProjectMetadata(projectConfigPath)
	if err != nil {
		return nil
	}
	project := projectJSONReport{
		Name:    strings.TrimSpace(metadata.name),
		Slug:    strings.TrimSpace(metadata.slug),
		Profile: projectProfile(),
	}
	if project.Name == "" && project.Slug == "" && project.Profile == "" {
		return nil
	}
	return &project
}

func envJSONFromReport(report envContractReport) envJSONReport {
	status := "ok"
	if len(report.missingRequired) > 0 || len(report.warnings) > 0 {
		status = "needs attention"
	}
	return envJSONReport{
		FileFound:       report.envFileFound,
		Status:          status,
		MissingRequired: append([]string(nil), report.missingRequired...),
		Warnings:        append([]string(nil), report.warnings...),
		Secrets:         report.secretCount,
		BrowserExposed:  report.browserCount,
		FrameworkOwned:  report.frameworkCount,
	}
}

func envProblemDetails(report envContractReport) []string {
	var details []string
	for _, name := range report.missingRequired {
		details = append(details, "missing required env value: "+name)
	}
	details = append(details, report.warnings...)
	return details
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
