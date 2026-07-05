package cli

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var version = "0.1.0-dev"
var commit = ""

const devLogPath = ".carbide/log/dev.jsonl"
const projectConfigPath = "carbide.toml"
const composeFilePath = "docker-compose.yml"
const legacyComposeFilePath = "compose.yml"
const defaultTerminalWidth = 80
const progressStateColumnWidth = 8
const minimumProgressFrameWidth = 4

const runtimeContractVersion = 1
const baselineGoModuleVersion = "1.25.0"
const baselineGoBuilderImage = "golang:1.26-bookworm@sha256:b305420a68d0f229d91eb3b3ed9e519fcf2cf5461da4bef997bf927e8c0bfd2b"
const baselineAPIRuntimeImage = "debian:trixie-slim@sha256:28de0877c2189802884ccd20f15ee41c203573bd87bb6b883f5f46362d24c5c2"
const baselineBunImage = "oven/bun:1.3.14-debian@sha256:9dba1a1b43ce28c9d7931bfc4eb00feb63b0114720a0277a8f939ae4dfc9db6f"
const baselinePostgresImage = "postgres:17-alpine@sha256:dc17045ccfd343b49600570ea734b9c4991cf1c3f3302e67df51e3b402dd55c4"
const baselineReactVersion = "19.2.7"
const baselineTailwindVersion = "4.3.2"

const defaultLogoText = `_____________________________________________________
________________________oo_______oo_______oo_________
_ooooo___ooooo__oo_ooo__oooooo________oooooo__ooooo__
oo___oo_oo___oo_ooo___o_oo___oo__oo__oo___oo_oo____o_
oo______oo___oo_oo______oo___oo__oo__oo___oo_ooooooo_
oo______oo___oo_oo______oo___oo__oo__oo___oo_oo______
_ooooo___oooo_o_oo______oooooo__oooo__oooooo__ooooo__
_____________________________________________________
`

const commandListText = `Carbide %s

Usage:
  carbide <command> [arguments]

Commands:
  new <project-name>   Create a new Carbide project
  init                 Initialize the current empty directory
  help                 Show detailed help
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
	out         io.Writer
	styled      bool
	interactive bool
	termWidth   int
}

type outputRow struct {
	key   string
	value string
}

type tableRow []string

type helpCommandSection struct {
	name string
	rows []outputRow
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

type composeServiceStatus struct {
	service string
	state   string
	health  string
}

type composeServicePort struct {
	URL           string `json:"URL"`
	TargetPort    int    `json:"TargetPort"`
	PublishedPort int    `json:"PublishedPort"`
	Protocol      string `json:"Protocol"`
}

type composeServiceSnapshot struct {
	Service    string               `json:"Service"`
	Name       string               `json:"Name"`
	State      string               `json:"State"`
	Health     string               `json:"Health"`
	Status     string               `json:"Status"`
	Ports      string               `json:"Ports"`
	Publishers []composeServicePort `json:"Publishers"`
}

type envSchema struct {
	Version   int           `json:"version"`
	Variables []envVariable `json:"variables"`
}

type envVariable struct {
	Name           string `json:"name"`
	Service        string `json:"service"`
	Required       bool   `json:"required"`
	Secret         bool   `json:"secret"`
	BrowserExposed bool   `json:"browser_exposed"`
	FrameworkOwned bool   `json:"framework_owned"`
	LocalDefault   string `json:"local_default"`
	Description    string `json:"description"`
}

type envContractReport struct {
	schema          envSchema
	envFileFound    bool
	missingRequired []string
	warnings        []string
	secretCount     int
	browserCount    int
	frameworkCount  int
}

type deployTarget struct {
	Name             string
	Type             string
	Host             string
	Hosts            map[string]deployHost
	Roles            []deployRole
	Domain           string
	RemotePath       string
	SourcePath       string
	ComposeFile      string
	ProjectDirectory string
	PublicPort       int
	HealthPath       string
	Nginx            bool
	NginxSite        string
	Strategy         string
}

type deployHost struct {
	Name        string
	SSH         string
	Address     string
	Description string
}

type deployRole struct {
	Name             string
	Hosts            []string
	RemotePath       string
	ComposeFile      string
	ProjectDirectory string
	PublicPort       int
	HealthPath       string
	Nginx            bool
	Primary          string
	Migration        string
}

func SetCommit(value string) {
	if value != "" {
		commit = value
	}
}

func Main() {
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
		a.printCommandList()
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		if len(args) != 1 {
			return errors.New("usage: carbide help")
		}
		a.printHelp()
		return nil
	case "version":
		if len(args) != 1 {
			return errors.New("usage: carbide version")
		}
		return a.commandVersion()
	case "upgrade":
		if len(args) != 1 {
			return errors.New("usage: carbide upgrade")
		}
		return a.commandUpgrade()
	case "new":
		if len(args) < 2 {
			return errors.New("usage: carbide new <project-name>")
		}
		return a.commandNew(strings.Join(args[1:], " "))
	case "init":
		if len(args) != 1 {
			return errors.New("usage: carbide init")
		}
		return a.commandInit()
	case "doctor":
		if len(args) == 1 {
			return a.commandDoctor()
		}
		if len(args) == 2 && args[1] == "env" {
			return a.commandDoctorEnv()
		}
		if len(args) == 2 && args[1] == "runtime" {
			return a.commandDoctorRuntime()
		}
		if len(args) == 2 && args[1] == "framework" {
			return a.commandDoctorFramework()
		}
		return errors.New("usage: carbide doctor [env|runtime|framework]")
	case "deploy":
		if len(args) == 3 && args[1] == "preview" {
			return a.commandDeployPreview(args[2])
		}
		if len(args) == 3 && args[1] == "apply" {
			return a.commandDeployApply(args[2])
		}
		return errors.New("usage: carbide deploy preview prod | carbide deploy apply prod")
	case "run":
		if len(args) == 2 && args[1] == "dev" {
			return a.commandRunDev()
		}
		return errors.New("usage: carbide run dev")
	case "status":
		if len(args) == 1 {
			return a.commandStatus()
		}
		return errors.New("usage: carbide status")
	case "stop":
		if len(args) == 2 && args[1] == "dev" {
			return a.commandStopDev()
		}
		return errors.New("usage: carbide stop dev")
	case "follow":
		if len(args) >= 2 && args[1] == "logs" {
			return a.commandFollowLogs(args[2:])
		}
		return errors.New("usage: carbide follow logs [service <name>] [containing <text>]")
	case "logs":
		return a.commandLogs(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a app) printCommandList() {
	r := newRenderer(a.stdout)
	logo := carbideLogo()
	if r.interactive {
		r.AnimateLogo(logo)
	} else {
		r.Logo(logo)
	}
	text := fmt.Sprintf(commandListText, version)
	if r.styled {
		fmt.Fprint(a.stdout, r.paint("38;5;245", text))
		return
	}
	fmt.Fprint(a.stdout, text)
}

func (a app) printHelp() {
	r := newRenderer(a.stdout)
	r.CommandList([]helpCommandSection{
		{
			rows: []outputRow{
				{"deploy apply prod", "apply production deploy"},
				{"deploy preview prod", "preview production deploy"},
				{"doctor", "check project contract"},
				{"doctor env", "validate env contract"},
				{"doctor framework", "run framework regressions"},
				{"doctor runtime", "run Docker runtime checks"},
				{"help", "show this help"},
				{"init", "init current directory"},
				{"logs", "query saved logs"},
				{"new <project-name>", "create project directory"},
				{"status", "show containers and ports"},
				{"upgrade", "upgrade CLI from GitHub"},
				{"version", "print installed version"},
			},
		},
		{
			name: "follow",
			rows: []outputRow{
				{"follow logs", "stream live logs"},
				{"follow logs service api", "stream one service"},
			},
		},
		{
			name: "logs",
			rows: []outputRow{
				{"logs containing \"/api/login\" json", "query logs as JSON"},
			},
		},
		{
			name: "run",
			rows: []outputRow{
				{"run dev", "start Docker dev stack"},
			},
		},
		{
			name: "stop",
			rows: []outputRow{
				{"stop dev", "stop dev containers"},
			},
		},
	})
}

func (a app) commandVersion() error {
	r := newRenderer(a.stdout)
	if commit != "" {
		r.Title("Carbide", "installed CLI")
		r.Rows(
			outputRow{"version", version},
			outputRow{"commit", commit},
		)
		return nil
	} else if head := gitShortHead(a.home); head != "" {
		r.Title("Carbide", "installed CLI")
		r.Rows(
			outputRow{"version", version},
			outputRow{"commit", head},
		)
		return nil
	}
	r.Title("Carbide", "installed CLI")
	r.Rows(outputRow{"version", version})
	return nil
}

func (a app) commandNew(name string) error {
	if err := ensureProjectName(name); err != nil {
		return err
	}

	slug := projectSlug(name)
	if slug == "" {
		slug = "carbide-app"
	}
	displayName := projectDisplayName(name)

	target, err := filepath.Abs(filepath.Join(".", slug))
	if err != nil {
		return err
	}
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("%s already exists", slug)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := a.copyScaffold(target, displayName, slug); err != nil {
		return err
	}

	newRenderer(a.stdout).Message(
		"Carbide",
		"project created",
		outputRow{"path", target},
		outputRow{"next", fmt.Sprintf("cd %s", slug)},
		outputRow{"", "carbide run dev"},
	)
	return nil
}

func (a app) commandInit() error {
	empty, err := isCurrentDirEmpty()
	if err != nil {
		return err
	}
	if !empty {
		return errors.New("carbide init requires an empty directory")
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
		slug = "carbide-app"
	}
	displayName := projectDisplayName(name)
	if err := a.copyScaffold(pwd, displayName, slug); err != nil {
		return err
	}

	newRenderer(a.stdout).Message(
		"Carbide",
		"project initialized",
		outputRow{"path", pwd},
		outputRow{"next", "carbide run dev"},
	)
	return nil
}

type doctorResult struct {
	check  string
	status string
	detail string
}

func (a app) commandDoctor() error {
	results := a.projectDoctorResults()
	results = append(results, doctorResult{"runtime", "skip", "run carbide doctor runtime"})
	return a.renderDoctorResults("project contract", results)
}

func (a app) commandDoctorEnv() error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}

	report, err := inspectEnvContract()
	if err != nil {
		return err
	}

	status := "ok"
	if len(report.missingRequired) > 0 || len(report.warnings) > 0 {
		status = "needs attention"
	}
	envFile := ".env not found; local defaults active"
	if report.envFileFound {
		envFile = ".env found"
	}

	r := newRenderer(a.stdout)
	r.Title("Carbide doctor", "environment contract")
	r.Rows(
		outputRow{"contract", projectConfigPath},
		outputRow{"env", envFile},
		outputRow{"status", status},
		outputRow{"required", fmt.Sprintf("%d missing", len(report.missingRequired))},
		outputRow{"secrets", fmt.Sprintf("%d declared", report.secretCount)},
		outputRow{"browser", fmt.Sprintf("%d exposed", report.browserCount)},
		outputRow{"framework", fmt.Sprintf("%d owned", report.frameworkCount)},
	)
	for _, name := range report.missingRequired {
		r.Row(outputRow{"missing", name})
	}
	for _, warning := range report.warnings {
		r.Row(outputRow{"warning", warning})
	}
	if len(report.missingRequired) > 0 {
		return fmt.Errorf("environment contract has %d missing required value(s)", len(report.missingRequired))
	}
	return nil
}

func (a app) commandDoctorRuntime() error {
	results := a.projectDoctorResults()
	if doctorFailures(results) > 0 {
		results = append(results, doctorResult{"runtime", "skip", "fix project contract first"})
		return a.renderDoctorResults("runtime contract", results)
	}

	runtimeResults := a.runtimeDoctorResults()
	results = append(results, runtimeResults...)
	return a.renderDoctorResults("runtime contract", results)
}

func (a app) commandDoctorFramework() error {
	results := a.frameworkDoctorResults()
	return a.renderDoctorResults("framework regressions", results)
}

func (a app) renderDoctorResults(subtitle string, results []doctorResult) error {
	rows := make([]tableRow, 0, len(results))
	for _, result := range results {
		rows = append(rows, tableRow{result.check, result.status, result.detail})
	}

	r := newRenderer(a.stdout)
	r.Title("Carbide doctor", subtitle)
	r.Table([]string{"check", "status", "detail"}, rows)

	failures := doctorFailures(results)
	if failures > 0 {
		return fmt.Errorf("doctor found %d failing check(s)", failures)
	}
	return nil
}

func doctorFailures(results []doctorResult) int {
	count := 0
	for _, result := range results {
		if result.status == "fail" {
			count++
		}
	}
	return count
}

func doctorOK(check string, detail string) doctorResult {
	return doctorResult{check: check, status: "ok", detail: detail}
}

func doctorFail(check string, detail string) doctorResult {
	return doctorResult{check: check, status: "fail", detail: detail}
}

func doctorWarn(check string, detail string) doctorResult {
	return doctorResult{check: check, status: "warn", detail: detail}
}

func doctorSkip(check string, detail string) doctorResult {
	return doctorResult{check: check, status: "skip", detail: detail}
}

func (a app) projectDoctorResults() []doctorResult {
	if !isFile(projectConfigPath) {
		return []doctorResult{doctorFail("project", "missing carbide.toml")}
	}

	if projectProfile() == "docs" {
		return []doctorResult{
			doctorDocsProjectShape(),
			doctorDocsConfigContract(),
			doctorDocsRuntimeBaselineContract(),
			doctorEnvContract(),
			doctorDocsComposeContract(),
			doctorDocsWebContract(),
			doctorDocsAPIContract(),
			doctorDocsDatabaseContract(),
			doctorDocsAgentsContract(),
			doctorForbiddenRegressions("."),
		}
	}

	return []doctorResult{
		doctorProjectShape(),
		doctorConfigContract(),
		doctorRuntimeBaselineContract(),
		doctorEnvContract(),
		doctorComposeContract(),
		doctorFrontendContract(),
		doctorAPIContract(),
		doctorDatabaseContract(),
		doctorAgentsContract(),
		doctorForbiddenRegressions("."),
	}
}

func doctorProjectShape() doctorResult {
	requiredDirs := []string{"web", "api", "db", "agents.d"}
	if missing := missingDirs(requiredDirs); len(missing) > 0 {
		return doctorFail("project shape", "missing "+strings.Join(missing, ", "))
	}

	requiredFiles := []string{projectConfigPath, composeFilePath, "AGENTS.md"}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("project shape", "missing "+strings.Join(missing, ", "))
	}

	forbidden := []string{"src", "model", "controller", "view", "views", "frontend", "templates", "include", "infra", "doc"}
	if found := existingDirs(forbidden); len(found) > 0 {
		return doctorFail("project shape", "legacy root dirs: "+strings.Join(found, ", "))
	}
	if isFile("go.mod") || isFile("go.sum") || isFile("Dockerfile") {
		return doctorFail("project shape", "root Go/Docker files are not allowed")
	}

	services := composeServiceNamesFromFile(composeFilePath)
	allowed := map[string]bool{"agents.d": true}
	for _, service := range services {
		allowed[service] = true
	}
	if len(services) == 0 {
		for _, service := range defaultComposeServices() {
			allowed[service] = true
		}
	}
	extras := rootDirsOutsideContract(allowed)
	if len(extras) > 0 {
		return doctorFail("project shape", "non-service root dirs: "+strings.Join(extras, ", "))
	}
	return doctorOK("project shape", "web api db agents.d")
}

func doctorConfigContract() doctorResult {
	content, err := os.ReadFile(projectConfigPath)
	if err != nil {
		return doctorFail("config", err.Error())
	}
	text := string(content)
	required := []string{
		"name = ",
		"slug = ",
		"carbide_version = ",
		"[dev]",
		"default_port = 8080",
		`database = "postgres"`,
		"[runtime]",
		fmt.Sprintf("contract_version = %d", runtimeContractVersion),
		`policy = "explicit-baseline"`,
		fmt.Sprintf(`go_module = "%s"`, baselineGoModuleVersion),
		fmt.Sprintf(`go_builder_image = "%s"`, baselineGoBuilderImage),
		fmt.Sprintf(`api_runtime_image = "%s"`, baselineAPIRuntimeImage),
		fmt.Sprintf(`bun_image = "%s"`, baselineBunImage),
		fmt.Sprintf(`postgres_image = "%s"`, baselinePostgresImage),
		fmt.Sprintf(`react = "%s"`, baselineReactVersion),
		fmt.Sprintf(`react_dom = "%s"`, baselineReactVersion),
		fmt.Sprintf(`tailwindcss = "%s"`, baselineTailwindVersion),
		fmt.Sprintf(`tailwind_cli = "%s"`, baselineTailwindVersion),
		"[env]",
		"contract_version = 1",
		"[deploy]",
		"preview_before_apply = true",
	}
	if missing := missingNeedles(text, required); len(missing) > 0 {
		return doctorFail("config", "missing "+strings.Join(missing, ", "))
	}
	return doctorOK("config", "carbide.toml")
}

func doctorRuntimeBaselineContract() doctorResult {
	required := map[string][]string{
		projectConfigPath: {
			fmt.Sprintf("contract_version = %d", runtimeContractVersion),
			`policy = "explicit-baseline"`,
			fmt.Sprintf(`go_module = "%s"`, baselineGoModuleVersion),
			fmt.Sprintf(`go_builder_image = "%s"`, baselineGoBuilderImage),
			fmt.Sprintf(`api_runtime_image = "%s"`, baselineAPIRuntimeImage),
			fmt.Sprintf(`bun_image = "%s"`, baselineBunImage),
			fmt.Sprintf(`postgres_image = "%s"`, baselinePostgresImage),
			fmt.Sprintf(`react = "%s"`, baselineReactVersion),
			fmt.Sprintf(`react_dom = "%s"`, baselineReactVersion),
			fmt.Sprintf(`tailwindcss = "%s"`, baselineTailwindVersion),
			fmt.Sprintf(`tailwind_cli = "%s"`, baselineTailwindVersion),
		},
		"api/Dockerfile": {
			"FROM " + baselineGoBuilderImage,
			"FROM " + baselineAPIRuntimeImage,
		},
		"web/Dockerfile": {
			"FROM " + baselineBunImage,
		},
		"web/package.json": {
			fmt.Sprintf(`"react": "%s"`, baselineReactVersion),
			fmt.Sprintf(`"react-dom": "%s"`, baselineReactVersion),
			fmt.Sprintf(`"tailwindcss": "%s"`, baselineTailwindVersion),
			fmt.Sprintf(`"@tailwindcss/cli": "%s"`, baselineTailwindVersion),
		},
		composeFilePath: {
			"image: " + baselinePostgresImage,
		},
		"api/go.mod": {
			"go " + baselineGoModuleVersion,
		},
		"db/go.mod": {
			"go " + baselineGoModuleVersion,
		},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return doctorFail("runtime baseline", path+" missing "+strings.Join(missing, ", "))
		}
	}
	if findings := floatingDockerReferences([]string{"api/Dockerfile", "web/Dockerfile", composeFilePath}); len(findings) > 0 {
		return doctorFail("runtime baseline", "floating Docker refs: "+strings.Join(findings, ", "))
	}
	if findings := packageVersionRangeFindings("web/package.json"); len(findings) > 0 {
		return doctorFail("runtime baseline", "package ranges: "+strings.Join(findings, ", "))
	}
	if findings := unsupportedGoDirectiveFindings([]string{"api/go.mod", "db/go.mod"}); len(findings) > 0 {
		return doctorFail("runtime baseline", "Go directive drift: "+strings.Join(findings, ", "))
	}
	return doctorOK("runtime baseline", "Go 1.25 React 19.2 Tailwind 4.3 Bun 1.3 Postgres 17")
}

func doctorEnvContract() doctorResult {
	report, err := inspectEnvContract()
	if err != nil {
		return doctorFail("env contract", err.Error())
	}
	detail := fmt.Sprintf("%d missing, %d secrets", len(report.missingRequired), report.secretCount)
	if len(report.missingRequired) > 0 {
		return doctorFail("env contract", detail)
	}
	if len(report.warnings) > 0 {
		return doctorFail("env contract", strings.Join(report.warnings, "; "))
	}
	return doctorOK("env contract", detail)
}

func doctorComposeContract() doctorResult {
	content, err := os.ReadFile(composeFilePath)
	if err != nil {
		return doctorFail("compose", "missing docker-compose.yml")
	}
	text := string(content)
	services := composeServiceNamesFromFile(composeFilePath)
	for _, service := range []string{"web", "api", "db"} {
		if !containsString(services, service) {
			return doctorFail("compose", "missing "+service+" service")
		}
	}
	if containsString(services, "backend") || containsString(services, "database") {
		return doctorFail("compose", "legacy service names present")
	}
	required := []string{
		"API_URL: http://api:8080",
		"@db:5432/carbide",
		`PUBLIC_URL: "http://localhost:${CARBIDE_HTTP_PORT:-8080}"`,
		"develop:",
		"watch:",
		"action: rebuild",
		"path: ./web/src",
		"path: ./api",
		"path: ./db",
	}
	if missing := missingNeedles(text, required); len(missing) > 0 {
		return doctorFail("compose", "missing "+strings.Join(missing, ", "))
	}
	return doctorOK("compose", "web api db")
}

func doctorFrontendContract() doctorResult {
	requiredFiles := []string{
		"web/Dockerfile",
		"web/package.json",
		"web/bun.lock",
		"web/index.html",
		"web/src/main.jsx",
		"web/src/server.jsx",
		"web/src/write-index.mjs",
		"web/src/styles.css",
		"web/src/lib/cx.js",
		"web/src/component/l1/Button.jsx",
		"web/src/component/l1/Field.jsx",
		"web/src/component/l1/Surface.jsx",
		"web/src/component/l1/Text.jsx",
		"web/src/component/l1/ThemeToggle.jsx",
		"web/src/component/l1/tokens.js",
		"web/src/component/l2/AuthForm.jsx",
		"web/src/component/l2/Layouts.jsx",
		"web/src/component/l3/AuthView.jsx",
		"web/src/component/l3/DashboardView.jsx",
		"web/src/component/l3/LoadingView.jsx",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("frontend", "missing "+strings.Join(missing, ", "))
	}
	requiredDirs := []string{"web/src/component/l1", "web/src/component/l2", "web/src/component/l3"}
	if missing := missingDirs(requiredDirs); len(missing) > 0 {
		return doctorFail("frontend", "missing "+strings.Join(missing, ", "))
	}
	forbiddenFiles := []string{"web/package-lock.json", "web/vite.config.js", "web/src/component/l1/theme.css"}
	if found := existingFiles(forbiddenFiles); len(found) > 0 {
		return doctorFail("frontend", "forbidden "+strings.Join(found, ", "))
	}
	if fileContains("web/src/styles.css", "theme.css") || treeContains("web/src", "cb-") || treeContains("web/src", "--cb-") {
		return doctorFail("frontend", "parallel CSS theme detected")
	}
	if fileContains("web/src/styles.css", "#0f766e") ||
		fileContains("web/src/styles.css", "#115e59") ||
		fileContains("web/src/styles.css", "#2dd4bf") ||
		fileContains("web/src/styles.css", "#5eead4") ||
		fileContains("web/src/styles.css", "#16433c") ||
		fileContains("web/src/styles.css", "#0f302c") ||
		fileContains("web/src/component/l1/tokens.js", "from-carbide-action via-carbide-hero-via") {
		return doctorFail("frontend", "green scaffold palette detected")
	}
	if fileContains("web/src/component/l2/Layouts.jsx", "text-7xl") ||
		fileContains("web/src/component/l2/Layouts.jsx", "text-5xl") ||
		fileContains("web/src/component/l2/Layouts.jsx", "py-24") ||
		fileContains("web/src/component/l2/Layouts.jsx", "lg:py-12") ||
		fileContains("web/src/component/l2/Layouts.jsx", "lg:grid-cols-[280px") ||
		fileContains("web/src/component/l2/Layouts.jsx", "lg:grid-cols-[240px") ||
		fileContains("web/src/component/l3/DashboardView.jsx", "gap-6") ||
		fileContains("web/src/component/l3/DashboardView.jsx", "p-6") ||
		fileContains("web/src/component/l1/Field.jsx", "min-h-12 rounded-md border") ||
		fileContains("web/src/component/l1/Field.jsx", "min-h-10 rounded-md border") ||
		treeContains("web/src/component", "font-extrabold") {
		return doctorFail("frontend", "oversized scaffold density detected")
	}
	if fileContains("web/src/component/l1/ThemeToggle.jsx", "aria-pressed") ||
		fileContains("web/src/component/l1/ThemeToggle.jsx", `role="group"`) ||
		fileContains("web/src/component/l1/ThemeToggle.jsx", `<select`) ||
		fileContains("web/src/component/l1/ThemeToggle.jsx", `appearance-none`) {
		return doctorFail("frontend", "non-icon theme toggle detected")
	}
	if !fileContains("web/package.json", `"react":`) ||
		!fileContains("web/package.json", `"tailwindcss":`) ||
		!fileContains("web/package.json", `"@tailwindcss/cli":`) ||
		!fileContains("web/package.json", `"assets:build":`) ||
		!fileContains("web/package.json", `--entry-naming='assets/[name]-[hash].[ext]'`) ||
		!fileContains("web/Dockerfile", `bun run assets:build`) ||
		!fileContains("web/src/server.jsx", `publicRoot`) ||
		!fileContains("web/src/server.jsx", `Cache-Control`) ||
		!fileContains("web/src/server.jsx", `public, max-age=31536000, immutable`) ||
		!fileContains("web/src/server.jsx", `return 'no-store'`) ||
		!fileContains("web/src/write-index.mjs", `asset-manifest.json`) ||
		!fileContains("web/src/write-index.mjs", `/assets/${scripts[0]}`) ||
		!fileContains("web/src/styles.css", `@import "tailwindcss";`) ||
		!fileContains("web/src/styles.css", `[data-theme="dark"]`) ||
		!fileContains("web/src/styles.css", `font-size: 14px`) ||
		!fileContains("web/src/styles.css", `line-height: 1.4`) ||
		!fileContains("web/src/styles.css", `--carbide-page: #ffffff`) ||
		!fileContains("web/src/styles.css", `--carbide-page: #000000`) ||
		!fileContains("web/index.html", `prefers-color-scheme: dark`) ||
		!fileContains("web/src/main.jsx", `carbide.theme`) ||
		!fileContains("web/src/component/l1/ThemeToggle.jsx", `SunIcon`) ||
		!fileContains("web/src/component/l1/ThemeToggle.jsx", `MoonIcon`) ||
		!fileContains("web/src/component/l1/ThemeToggle.jsx", `Switch to light theme`) ||
		!fileContains("web/src/component/l1/ThemeToggle.jsx", `Switch to dark theme`) ||
		!fileContains("web/src/component/l1/ThemeToggle.jsx", `size-8 rounded-full border`) ||
		!fileContains("web/src/component/l1/ThemeToggle.jsx", `data-theme-mode`) ||
		!fileContains("web/src/component/l1/tokens.js", `bg-carbide-hero text-carbide-hero-text`) ||
		!fileContains("web/src/component/l1/Text.jsx", `text-2xl/8 sm:text-3xl/9`) ||
		!fileContains("web/src/component/l1/Field.jsx", `min-h-8 rounded-md border px-2 py-1 text-sm/6`) ||
		!fileContains("web/src/component/l1/Button.jsx", `md: 'min-h-8 px-3 text-xs'`) ||
		!fileContains("web/src/component/l2/AuthForm.jsx", `gap-3 border-l px-4 py-5`) ||
		!fileContains("web/src/component/l2/AuthForm.jsx", `w-full max-w-sm justify-self-center gap-3`) ||
		!fileContains("web/src/component/l2/Layouts.jsx", `lg:grid-cols-[216px_minmax(0,1fr)]`) ||
		!fileContains("web/src/component/l2/Layouts.jsx", `px-3 py-4 sm:px-5 lg:py-5`) ||
		!fileContains("web/src/main.jsx", "./component/l3/index.js") {
		return doctorFail("frontend", "React/Bun/Tailwind contract drifted")
	}
	return doctorOK("frontend", "Bun React Tailwind")
}

func doctorAPIContract() doctorResult {
	requiredFiles := []string{
		"api/Dockerfile",
		"api/go.mod",
		"api/go.sum",
		"api/main.go",
		"api/auth.go",
		"api/routes.go",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("api", "missing "+strings.Join(missing, ", "))
	}
	if fileContains("api/Dockerfile", "gcc") || fileContains("api/Dockerfile", "libpq-dev") || anyPathWithExtension("api", ".c", ".h") {
		return doctorFail("api", "legacy C backend artifacts present")
	}
	required := map[string][]string{
		"api/go.mod":     {"module carbideapp/api", "carbideapp/db", "replace carbideapp/db => ../db"},
		"api/Dockerfile": {"FROM golang:", "go mod download", "COPY api ./api", "COPY db ./db"},
		"api/routes.go":  {"/api/register", "/api/login", "/api/me", "handleDashboard"},
		"api/main.go":    {"api listening on container port", "public API URL is"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return doctorFail("api", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return doctorOK("api", "Go HTTP API")
}

func doctorDatabaseContract() doctorResult {
	requiredFiles := []string{
		"db/go.mod",
		"db/go.sum",
		"db/user.go",
		"db/session.go",
		"db/migration/001_auth.sql",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("database", "missing "+strings.Join(missing, ", "))
	}
	required := map[string][]string{
		"db/go.mod":                 {"module carbideapp/db", "github.com/jackc/pgx/v5"},
		"db/user.go":                {"CreateUser", "VerifyUser", "pgxpool"},
		"db/session.go":             {"CreateSession", "CurrentUser", "DestroySession"},
		"db/migration/001_auth.sql": {"CREATE TABLE IF NOT EXISTS users", "CREATE TABLE IF NOT EXISTS sessions"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return doctorFail("database", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return doctorOK("database", "Postgres users sessions")
}

func doctorAgentsContract() doctorResult {
	requiredFiles := []string{
		"AGENTS.md",
		"agents.d/ENVIRONMENT.md",
		"agents.d/DEPLOY.md",
		"agents.d/BACKUP_RESTORE.md",
		"agents.d/TAILWIND_COMPONENTS.md",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("agents", "missing "+strings.Join(missing, ", "))
	}
	required := map[string][]string{
		"AGENTS.md":                       {"carbide run dev", "carbide status", "carbide doctor"},
		"agents.d/ENVIRONMENT.md":         {"separate secrets container", "carbide doctor"},
		"agents.d/DEPLOY.md":              {"preview-before-apply", "carbide deploy preview"},
		"agents.d/BACKUP_RESTORE.md":      {"Postgres owns durable application state"},
		"agents.d/TAILWIND_COMPONENTS.md": {"Tailwind Component Organization", "component/l1/", "component/l2/", "component/l3/"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return doctorFail("agents", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return doctorOK("agents", "AGENTS.md agents.d")
}

func doctorDocsProjectShape() doctorResult {
	requiredDirs := []string{"web", "api", "db", "agents.d"}
	if missing := missingDirs(requiredDirs); len(missing) > 0 {
		return doctorFail("project shape", "missing "+strings.Join(missing, ", "))
	}

	requiredFiles := []string{projectConfigPath, composeFilePath, "AGENTS.md"}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("project shape", "missing "+strings.Join(missing, ", "))
	}
	if isFile("go.mod") || isFile("go.sum") || isFile("Dockerfile") {
		return doctorFail("project shape", "root Go/Docker files are not allowed")
	}

	services := composeServiceNamesFromFile(composeFilePath)
	allowed := map[string]bool{"agents.d": true}
	for _, service := range services {
		allowed[service] = true
	}
	extras := rootDirsOutsideContract(allowed)
	if len(extras) > 0 {
		return doctorFail("project shape", "non-service root dirs: "+strings.Join(extras, ", "))
	}
	return doctorOK("project shape", "docs web api db")
}

func doctorDocsConfigContract() doctorResult {
	content := readFileString(projectConfigPath)
	required := []string{
		`name = "Carbide Docs"`,
		`slug = "carbide-docs"`,
		`profile = "docs"`,
		"carbide_version = ",
		"[dev]",
		"default_port = 8080",
		`database = "postgres"`,
		"[runtime]",
		fmt.Sprintf("contract_version = %d", runtimeContractVersion),
		`policy = "explicit-baseline"`,
		fmt.Sprintf(`go_module = "%s"`, baselineGoModuleVersion),
		fmt.Sprintf(`go_builder_image = "%s"`, baselineGoBuilderImage),
		fmt.Sprintf(`api_runtime_image = "%s"`, baselineAPIRuntimeImage),
		fmt.Sprintf(`bun_image = "%s"`, baselineBunImage),
		fmt.Sprintf(`postgres_image = "%s"`, baselinePostgresImage),
		"[env]",
		"contract_version = 1",
		"[deploy]",
		"preview_before_apply = true",
		"[deploy.hosts.de-sci]",
		`ssh = "de-sci"`,
		"[deploy.targets.de-sci]",
		`type = "ssh-compose"`,
		`host = "de-sci"`,
		`remote_path = "/opt/carbide/docs"`,
		"[deploy.targets.de-sci-environment]",
		`type = "ssh-compose-environment"`,
		"[deploy.targets.de-sci-environment.roles.web]",
		"[deploy.targets.de-sci-environment.roles.api]",
		"[deploy.targets.de-sci-environment.roles.db]",
		`migration = "once"`,
	}
	if missing := missingNeedles(content, required); len(missing) > 0 {
		return doctorFail("config", "missing "+strings.Join(missing, ", "))
	}
	return doctorOK("config", "docs deploy target")
}

func doctorDocsRuntimeBaselineContract() doctorResult {
	required := map[string][]string{
		projectConfigPath: {
			fmt.Sprintf("contract_version = %d", runtimeContractVersion),
			`policy = "explicit-baseline"`,
			fmt.Sprintf(`go_builder_image = "%s"`, baselineGoBuilderImage),
			fmt.Sprintf(`api_runtime_image = "%s"`, baselineAPIRuntimeImage),
			fmt.Sprintf(`bun_image = "%s"`, baselineBunImage),
			fmt.Sprintf(`postgres_image = "%s"`, baselinePostgresImage),
		},
		"api/Dockerfile": {
			"FROM " + baselineGoBuilderImage,
			"FROM " + baselineAPIRuntimeImage,
		},
		"web/Dockerfile": {
			"FROM " + baselineBunImage,
		},
		"web/package.json": {
			fmt.Sprintf(`"react": "%s"`, baselineReactVersion),
			fmt.Sprintf(`"react-dom": "%s"`, baselineReactVersion),
			fmt.Sprintf(`"tailwindcss": "%s"`, baselineTailwindVersion),
			fmt.Sprintf(`"@tailwindcss/cli": "%s"`, baselineTailwindVersion),
		},
		composeFilePath: {
			"image: " + baselinePostgresImage,
		},
		"api/go.mod": {
			"go " + baselineGoModuleVersion,
			"github.com/jackc/pgx/v5",
		},
		"db/go.mod": {
			"go " + baselineGoModuleVersion,
		},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return doctorFail("runtime baseline", path+" missing "+strings.Join(missing, ", "))
		}
	}
	if findings := floatingDockerReferences([]string{"api/Dockerfile", "web/Dockerfile", composeFilePath}); len(findings) > 0 {
		return doctorFail("runtime baseline", "floating Docker refs: "+strings.Join(findings, ", "))
	}
	if findings := unsupportedGoDirectiveFindings([]string{"api/go.mod", "db/go.mod"}); len(findings) > 0 {
		return doctorFail("runtime baseline", "Go directive drift: "+strings.Join(findings, ", "))
	}
	if findings := packageVersionRangeFindings("web/package.json"); len(findings) > 0 {
		return doctorFail("runtime baseline", "package ranges: "+strings.Join(findings, ", "))
	}
	return doctorOK("runtime baseline", "docs pinned images")
}

func doctorDocsComposeContract() doctorResult {
	content := readFileString(composeFilePath)
	services := composeServiceNamesFromFile(composeFilePath)
	for _, service := range []string{"web", "api", "db"} {
		if !containsString(services, service) {
			return doctorFail("compose", "missing "+service+" service")
		}
	}
	required := []string{
		"context: ..",
		"dockerfile: app/web/Dockerfile",
		"API_URL: http://api:8080",
		"CARBIDE_HTTP_PORT",
		"service_healthy",
		"./db/migration:/docker-entrypoint-initdb.d:ro",
	}
	if missing := missingNeedles(content, required); len(missing) > 0 {
		return doctorFail("compose", "missing "+strings.Join(missing, ", "))
	}
	return doctorOK("compose", "docs web api db")
}

func doctorDocsWebContract() doctorResult {
	requiredFiles := []string{
		"web/Dockerfile",
		"web/package.json",
		"web/bun.lock",
		"web/src/build-styles.js",
		"web/src/server.jsx",
		"web/src/styles.css",
		"web/src/lib/cx.js",
		"web/src/component/l1/Text.jsx",
		"web/src/component/l1/Surface.jsx",
		"web/src/component/l1/index.js",
		"web/src/component/l1/tokens.js",
		"web/src/component/l2/DocsChrome.jsx",
		"web/src/component/l2/index.js",
		"web/src/component/l3/DocsSite.jsx",
		"web/src/component/l3/index.js",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("web", "missing "+strings.Join(missing, ", "))
	}
	requiredDirs := []string{"web/src/component/l1", "web/src/component/l2", "web/src/component/l3", "web/src/lib"}
	if missing := missingDirs(requiredDirs); len(missing) > 0 {
		return doctorFail("web", "missing "+strings.Join(missing, ", "))
	}
	required := map[string][]string{
		"web/Dockerfile":                      {"COPY app/web/src ./src", "bun run tailwind:build", `CMD ["bun", "run", "start"]`, "COPY site ./site"},
		"web/package.json":                    {`"tailwind:build"`, `"@tailwindcss/cli":`, `"react":`, `"react-dom":`, `"tailwindcss":`},
		"web/src/build-styles.js":             {"tailwindcss", "./src/styles.css", "styles.css"},
		"web/src/styles.css":                  {`@import "tailwindcss";`, `@source "./component/**/*.jsx";`},
		"web/src/server.jsx":                  {"serveStatic", "proxy(request", `url.pathname === "/health"`, `url.pathname.startsWith("/api/")`, `./component/l3/index.js`, "docsResponseHeaders", "cacheBustHtml", "versionedAssetPath", "createHash", `?v=${hash}`},
		"web/src/component/l1/tokens.js":      {"docsClassLayers", "l1:", "l2:", "l3:"},
		"web/src/component/l2/DocsChrome.jsx": {"docsChromeClassLayers", "docsStaticHeaders"},
		"web/src/component/l3/DocsSite.jsx":   {"docsSiteClassLayers", "docsWebContract", "docsResponseHeaders"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return doctorFail("web", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return doctorOK("web", "Bun React Tailwind docs")
}

func doctorDocsAPIContract() doctorResult {
	requiredFiles := []string{
		"api/Dockerfile",
		"api/go.mod",
		"api/go.sum",
		"api/main.go",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("api", "missing "+strings.Join(missing, ", "))
	}
	required := map[string][]string{
		"api/go.mod":  {"module carbidedocs/api", "github.com/jackc/pgx/v5"},
		"api/main.go": {"/health", "/api/version", "pgxpool", "database unavailable"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return doctorFail("api", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return doctorOK("api", "docs health API")
}

func doctorDocsDatabaseContract() doctorResult {
	requiredFiles := []string{
		"db/go.mod",
		"db/migration/001_docs.sql",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("database", "missing "+strings.Join(missing, ", "))
	}
	if !fileContains("db/migration/001_docs.sql", "CREATE TABLE IF NOT EXISTS deploy_checks") {
		return doctorFail("database", "missing deploy check migration")
	}
	return doctorOK("database", "Postgres deploy checks")
}

func doctorDocsAgentsContract() doctorResult {
	requiredFiles := []string{
		"AGENTS.md",
		"agents.d/ENVIRONMENT.md",
		"agents.d/DEPLOY.md",
		"agents.d/BACKUP_RESTORE.md",
		"agents.d/TAILWIND_COMPONENTS.md",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return doctorFail("agents", "missing "+strings.Join(missing, ", "))
	}
	required := map[string][]string{
		"AGENTS.md":                       {"carbide doctor", "carbide deploy preview de-sci"},
		"agents.d/ENVIRONMENT.md":         {"remote `.env`", "POSTGRES_PASSWORD"},
		"agents.d/DEPLOY.md":              {"preview-before-apply", "ssh-compose"},
		"agents.d/BACKUP_RESTORE.md":      {"Postgres", "carbide_docs_pgdata"},
		"agents.d/TAILWIND_COMPONENTS.md": {"Tailwind", "component/l1", "component/l2", "component/l3"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return doctorFail("agents", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return doctorOK("agents", "docs agents.d")
}

func doctorForbiddenRegressions(root string) doctorResult {
	forbidden := []string{
		"Sea" + "lion",
		"sea" + "lion",
		"admin@carbide.local",
		"Demo login",
		"seed_admin",
		"render_template_text",
		"respond_view",
	}
	if hits := treeContainsAny(root, forbidden); len(hits) > 0 {
		return doctorFail("regressions", strings.Join(hits, ", "))
	}
	return doctorOK("regressions", "no legacy markers")
}

func (a app) runtimeDoctorResults() []doctorResult {
	if !isFile(projectConfigPath) {
		return []doctorResult{doctorFail("runtime", "run this inside a Carbide project")}
	}
	profile := projectProfile()

	compose, err := findCompose()
	if err != nil {
		return []doctorResult{doctorFail("runtime", err.Error())}
	}

	env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
	env = composeEnv(env)
	alreadyRunning := composeHasRunningServices(compose, env)
	port := 0
	if alreadyRunning {
		port = publishedWebPort(compose, env)
		if port == 0 {
			port = runtimePortFromEnv()
		}
	} else {
		selected, err := chooseDevPort(os.Getenv("CARBIDE_HTTP_PORT"))
		if err != nil {
			return []doctorResult{doctorFail("runtime", err.Error())}
		}
		port = selected
		env = setEnv(env, "CARBIDE_HTTP_PORT", strconv.Itoa(port))
	}

	results := []doctorResult{}
	if _, err := runComposeCaptured(compose, env, "config"); err != nil {
		return append(results, doctorFail("compose config", err.Error()))
	}
	results = append(results, doctorOK("compose config", "valid"))

	startedByDoctor := !alreadyRunning
	if startedByDoctor {
		if err := composeUpDetached(compose, env); err != nil {
			results = append(results, doctorFail("stack start", err.Error()))
			return results
		}
		results = append(results, doctorOK("stack start", fmt.Sprintf("localhost:%d", port)))
	} else {
		results = append(results, doctorOK("stack start", "already running"))
	}

	cleanupNeeded := startedByDoctor
	if cleanupNeeded {
		defer func() {
			if cleanupNeeded {
				_ = composeDown(compose, env)
			}
		}()
	}

	client := &http.Client{Timeout: 10 * time.Second}
	if err := waitForHTTP(client, fmt.Sprintf("http://localhost:%d/health", port), 60*time.Second); err != nil {
		results = append(results, doctorFail("health", err.Error()))
		return results
	}
	results = append(results, doctorOK("health", "/health"))

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if profile == "docs" {
		if err := httpGetContains(client, baseURL+"/api/version", `"name":"Carbide Docs"`); err != nil {
			results = append(results, doctorFail("version api", err.Error()))
			return results
		}
		results = append(results, doctorOK("version api", "/api/version"))

		if startedByDoctor {
			if err := composeDown(compose, env); err != nil {
				results = append(results, doctorWarn("cleanup", err.Error()))
			} else {
				results = append(results, doctorOK("cleanup", "stopped doctor stack"))
			}
			cleanupNeeded = false
		}
		return results
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		results = append(results, doctorFail("auth flow", err.Error()))
		return results
	}
	client.Jar = jar
	if err := httpGetContains(client, baseURL+"/api/me", `"authenticated":false`); err != nil {
		results = append(results, doctorFail("anonymous", err.Error()))
		return results
	}
	results = append(results, doctorOK("anonymous", "/api/me"))

	email := fmt.Sprintf("doctor-%d@carbide.local", time.Now().UnixNano())
	if err := httpPostFormContains(client, baseURL+"/api/register", url.Values{"email": {email}, "password": {"password"}}, `"ok":true`); err != nil {
		results = append(results, doctorFail("register", err.Error()))
		return results
	}
	results = append(results, doctorOK("register", "first-user flow"))

	if err := httpGetContains(client, baseURL+"/api/dashboard", email); err != nil {
		results = append(results, doctorFail("dashboard api", err.Error()))
		return results
	}
	if err := httpGetContains(client, baseURL+"/dashboard", `<div id="root"></div>`); err != nil {
		results = append(results, doctorFail("dashboard web", err.Error()))
		return results
	}
	results = append(results, doctorOK("dashboard", "api and web shell"))

	if err := httpPostFormContains(client, baseURL+"/api/logout", nil, `"ok":true`); err != nil {
		results = append(results, doctorFail("logout", err.Error()))
		return results
	}
	if err := httpGetContains(client, baseURL+"/api/me", `"authenticated":false`); err != nil {
		results = append(results, doctorFail("logout", err.Error()))
		return results
	}
	results = append(results, doctorOK("logout", "session cleared"))

	if startedByDoctor {
		if err := composeDown(compose, env); err != nil {
			results = append(results, doctorWarn("cleanup", err.Error()))
		} else {
			results = append(results, doctorOK("cleanup", "stopped doctor stack"))
		}
		cleanupNeeded = false
	}
	return results
}

func (a app) frameworkDoctorResults() []doctorResult {
	if !isFile(filepath.Join(a.home, "cli", "go.mod")) || !isDir(filepath.Join(a.home, "tests")) {
		return []doctorResult{doctorFail("framework", "run from a Carbide source checkout")}
	}
	env, cleanup, err := frameworkDoctorCommandEnv(a.home)
	if err != nil {
		return []doctorResult{doctorFail("framework", err.Error())}
	}
	defer cleanup()

	type frameworkCheck struct {
		name string
		run  func() error
	}
	checks := []frameworkCheck{
		{
			name: "shell syntax",
			run: func() error {
				_, err := commandOutputEnv(
					a.home,
					env,
					"bash",
					"-n",
					"tests/contract/audit_versions.sh",
					"tests/contract/check_repo_contract.sh",
					"tests/scaffold/cli_scaffold.sh",
					"tests/smoke/starter_docker_flow.sh",
					"cli/bin/carbide",
					"cli/install.sh",
				)
				return err
			},
		},
		{name: "Go CLI tests", run: func() error { return runFrameworkGoTests(a.home) }},
		{name: "repo contract", run: func() error {
			_, err := commandOutputEnv(a.home, env, "bash", "tests/contract/check_repo_contract.sh")
			return err
		}},
		{name: "CLI scaffold", run: func() error {
			_, err := commandOutputEnv(a.home, env, "bash", "tests/scaffold/cli_scaffold.sh")
			return err
		}},
		{name: "Docker smoke", run: func() error {
			_, err := commandOutputEnv(a.home, env, "bash", "tests/smoke/starter_docker_flow.sh")
			return err
		}},
	}

	results := make([]doctorResult, 0, len(checks))
	for _, check := range checks {
		if err := check.run(); err != nil {
			results = append(results, doctorFail(check.name, firstLine(err.Error())))
			continue
		}
		results = append(results, doctorOK(check.name, "passed"))
	}
	return results
}

func missingFiles(paths []string) []string {
	var missing []string
	for _, path := range paths {
		if !isFile(path) {
			missing = append(missing, path)
		}
	}
	return missing
}

func missingDirs(paths []string) []string {
	var missing []string
	for _, path := range paths {
		if !isDir(path) {
			missing = append(missing, path)
		}
	}
	return missing
}

func existingFiles(paths []string) []string {
	var found []string
	for _, path := range paths {
		if isFile(path) {
			found = append(found, path)
		}
	}
	return found
}

func existingDirs(paths []string) []string {
	var found []string
	for _, path := range paths {
		if isDir(path) {
			found = append(found, path)
		}
	}
	return found
}

func missingNeedles(content string, needles []string) []string {
	var missing []string
	for _, needle := range needles {
		if !strings.Contains(content, needle) {
			missing = append(missing, needle)
		}
	}
	return missing
}

func readFileString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func fileContains(path string, needle string) bool {
	return strings.Contains(readFileString(path), needle)
}

func floatingDockerReferences(paths []string) []string {
	var findings []string
	for _, path := range paths {
		content := readFileString(path)
		for _, ref := range dockerImageRefs(path, content) {
			if isFloatingImageRef(ref) {
				findings = append(findings, path+" "+ref)
			}
		}
	}
	return findings
}

func dockerImageRefs(path string, content string) []string {
	var refs []string
	if strings.HasSuffix(path, "Dockerfile") {
		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 2 && strings.EqualFold(fields[0], "FROM") {
				refs = append(refs, strings.Trim(fields[1], `"'`))
			}
		}
		return refs
	}

	imageLine := regexp.MustCompile(`(?m)^\s*image:\s*["']?([^"'\s]+)`)
	for _, match := range imageLine.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			refs = append(refs, strings.Trim(match[1], `"'`))
		}
	}
	return refs
}

func isFloatingImageRef(ref string) bool {
	if ref == "" || strings.Contains(ref, "${") {
		return false
	}
	return !strings.Contains(ref, "@sha256:") || strings.Contains(ref, ":latest")
}

func packageVersionRangeFindings(path string) []string {
	content := readFileString(path)
	var findings []string
	for _, name := range []string{"react", "react-dom", "tailwindcss", "@tailwindcss/cli"} {
		version := packageVersion(content, name)
		if version == "" {
			findings = append(findings, name+" missing")
			continue
		}
		if isSemverRange(version) {
			findings = append(findings, name+" "+version)
		}
	}
	return findings
}

func packageVersion(content string, name string) string {
	pattern := regexp.MustCompile(`"` + regexp.QuoteMeta(name) + `"\s*:\s*"([^"]+)"`)
	match := pattern.FindStringSubmatch(content)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func isSemverRange(version string) bool {
	trimmed := strings.TrimSpace(version)
	if trimmed == "" {
		return true
	}
	rangeMarkers := []string{"^", "~", ">", "<", "*", "x", "X", "latest", "||", " - "}
	for _, marker := range rangeMarkers {
		if strings.Contains(trimmed, marker) {
			return true
		}
	}
	return false
}

func unsupportedGoDirectiveFindings(paths []string) []string {
	var findings []string
	for _, path := range paths {
		version := goDirectiveVersion(path)
		if version != baselineGoModuleVersion {
			if version == "" {
				version = "missing"
			}
			findings = append(findings, path+" "+version)
		}
	}
	return findings
}

func goDirectiveVersion(path string) string {
	scanner := bufio.NewScanner(strings.NewReader(readFileString(path)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[0] == "go" {
			return fields[1]
		}
	}
	return ""
}

func containsString(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

func composeServiceNamesFromFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var services []string
	seen := map[string]bool{}
	inServices := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !strings.HasPrefix(raw, " ") && strings.HasSuffix(trimmed, ":") {
			inServices = trimmed == "services:"
			continue
		}
		if !inServices {
			continue
		}
		if strings.HasPrefix(raw, "  ") && !strings.HasPrefix(raw, "    ") && strings.HasSuffix(trimmed, ":") {
			service := strings.TrimSuffix(trimmed, ":")
			if service != "" && !seen[service] {
				seen[service] = true
				services = append(services, service)
			}
		}
	}
	return services
}

func rootDirsOutsideContract(allowed map[string]bool) []string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return []string{err.Error()}
	}
	var extras []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if !allowed[name] {
			extras = append(extras, name)
		}
	}
	return extras
}

func treeContains(root string, needle string) bool {
	return len(treeContainsAny(root, []string{needle})) > 0
}

func treeContainsAny(root string, needles []string) []string {
	found := map[string]bool{}
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		name := entry.Name()
		if entry.IsDir() {
			switch name {
			case ".git", ".carbide", ".cli", ".bin", "node_modules", "vendor":
				return filepath.SkipDir
			default:
				return nil
			}
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(data)
		for _, needle := range needles {
			if strings.Contains(text, needle) {
				found[needle] = true
			}
		}
		return nil
	})

	var hits []string
	for _, needle := range needles {
		if found[needle] {
			hits = append(hits, needle)
		}
	}
	return hits
}

func anyPathWithExtension(root string, extensions ...string) bool {
	found := false
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || found {
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		for _, extension := range extensions {
			if strings.HasSuffix(path, extension) {
				found = true
				return nil
			}
		}
		return nil
	})
	return found
}

func composeHasRunningServices(compose composeCommand, env []string) bool {
	snapshots, err := composeServiceSnapshots(compose, env)
	if err != nil {
		return false
	}
	for _, snapshot := range snapshots {
		if strings.EqualFold(strings.TrimSpace(snapshot.State), "running") {
			return true
		}
	}
	return false
}

func publishedWebPort(compose composeCommand, env []string) int {
	snapshots, err := composeServiceSnapshots(compose, env)
	if err != nil {
		return 0
	}
	web, ok := snapshots["web"]
	if !ok {
		return 0
	}
	for _, publisher := range web.Publishers {
		if publisher.PublishedPort > 0 && publisher.TargetPort == 8080 {
			return publisher.PublishedPort
		}
	}
	for _, publisher := range web.Publishers {
		if publisher.PublishedPort > 0 {
			return publisher.PublishedPort
		}
	}
	return 0
}

func runtimePortFromEnv() int {
	if value := os.Getenv("CARBIDE_HTTP_PORT"); value != "" {
		if port, err := validatePort(value); err == nil {
			return port
		}
	}
	return 8080
}

func waitForHTTP(client *http.Client, endpoint string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(endpoint)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
			lastErr = fmt.Errorf("%s returned %s", endpoint, resp.Status)
		} else {
			lastErr = err
		}
		time.Sleep(time.Second)
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("%s did not become ready", endpoint)
	}
	return lastErr
}

func httpGetContains(client *http.Client, endpoint string, needle string) error {
	resp, err := client.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s returned %s", endpoint, resp.Status)
	}
	if !strings.Contains(string(body), needle) {
		return fmt.Errorf("%s did not contain %q", endpoint, needle)
	}
	return nil
}

func httpPostFormContains(client *http.Client, endpoint string, values url.Values, needle string) error {
	if values == nil {
		values = url.Values{}
	}
	resp, err := client.PostForm(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s returned %s", endpoint, resp.Status)
	}
	if !strings.Contains(string(body), needle) {
		return fmt.Errorf("%s did not contain %q", endpoint, needle)
	}
	return nil
}

func runFrameworkGoTests(home string) error {
	if _, err := exec.LookPath("go"); err == nil {
		_, err := commandOutput(filepath.Join(home, "cli"), "go", "test", "./...")
		return err
	}
	if _, err := exec.LookPath("docker"); err != nil {
		return errors.New("Go is not installed and Docker is unavailable for Go test fallback")
	}
	args := []string{
		"run",
		"--rm",
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-e", "HOME=/tmp",
		"-e", "GOCACHE=/tmp/gocache",
		"-e", "GOMODCACHE=/tmp/gomodcache",
		"-v", home + ":" + home,
		"-w", filepath.Join(home, "cli"),
		"golang:1.25-bookworm",
		"bash",
		"-lc",
		"export PATH=/usr/local/go/bin:$PATH; go test ./...",
	}
	_, err := commandOutput("", "docker", args...)
	return err
}

func frameworkDoctorCommandEnv(home string) ([]string, func(), error) {
	env := setEnv(os.Environ(), "CARBIDE_HOME", home)
	if _, err := exec.LookPath("go"); err == nil {
		return env, func() {}, nil
	}
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return nil, func() {}, errors.New("Go is not installed and Docker is unavailable for framework checks")
	}

	dir, err := os.MkdirTemp("", "carbide-doctor-go-")
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	wrapper := fmt.Sprintf(`#!/usr/bin/env bash
set -euo pipefail
root="${CARBIDE_HOME:-$(pwd -P)}"
workdir="$(pwd -P)"
exec %q run --rm \
  --user "$(id -u):$(id -g)" \
  -e HOME=/tmp \
  -e GOCACHE=/tmp/gocache \
  -e GOMODCACHE=/tmp/gomodcache \
  -v "$root:$root" \
  -w "$workdir" \
  golang:1.25-bookworm \
  /usr/local/go/bin/go "$@"
`, dockerPath)
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(wrapper), 0755); err != nil {
		cleanup()
		return nil, func() {}, err
	}
	env = setEnv(env, "PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return env, cleanup, nil
}

func firstLine(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "failed"
	}
	if before, _, ok := strings.Cut(value, "\n"); ok {
		return before
	}
	return value
}

func (a app) commandDeployPreview(target string) error {
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

	report, err := inspectEnvContract()
	if err != nil {
		return err
	}

	envStatus := "ok"
	if len(report.missingRequired) > 0 || len(report.warnings) > 0 {
		envStatus = "needs attention"
	}

	r := newRenderer(a.stdout)
	r.Title("Carbide deploy", fmt.Sprintf("preview %s", target))
	if found {
		if err := validateDeployTarget(deploy); err != nil {
			return err
		}
		source, err := filepath.Abs(deploy.SourcePath)
		if err != nil {
			return err
		}
		if deploy.Type == "ssh-compose-environment" {
			r.Rows(
				outputRow{"target", deploy.Name},
				outputRow{"type", deploy.Type},
				outputRow{"domain", deploy.Domain},
				outputRow{"source", source},
				outputRow{"mutates", "no"},
				outputRow{"env", envStatus},
				outputRow{"hosts", strings.Join(deployHostRows(deploy), "\n")},
				outputRow{"roles", strings.Join(deployRoleRows(deploy), "\n")},
				outputRow{"apply", "disabled until clustered orchestration is implemented"},
			)
			return nil
		}
		r.Rows(
			outputRow{"target", deploy.Name},
			outputRow{"type", deploy.Type},
			outputRow{"host", deploy.Host},
			outputRow{"domain", deploy.Domain},
			outputRow{"source", source},
			outputRow{"remote", deploy.RemotePath},
			outputRow{"compose", deploy.ComposeFile},
			outputRow{"port", strconv.Itoa(deploy.PublicPort)},
			outputRow{"mutates", "no"},
			outputRow{"env", envStatus},
			outputRow{"apply", fmt.Sprintf("carbide deploy apply %s", deploy.Name)},
		)
		return nil
	}

	r.Rows(
		outputRow{"target", target},
		outputRow{"mutates", "no"},
		outputRow{"env", envStatus},
		outputRow{"plan", "validate env contract\nuse checked-in deploy target when one exists\nrefuse apply until target is implemented"},
	)
	return nil
}

func (a app) commandDeployApply(target string) error {
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

	r := newRenderer(a.stdout)
	r.Title("Carbide deploy", fmt.Sprintf("apply %s", target))
	if found {
		if err := validateDeployTarget(deploy); err != nil {
			return err
		}
		if deploy.Type == "ssh-compose-environment" {
			r.Rows(
				outputRow{"target", deploy.Name},
				outputRow{"type", deploy.Type},
				outputRow{"status", "guarded"},
				outputRow{"reason", "clustered apply needs explicit orchestration"},
				outputRow{"preview", fmt.Sprintf("carbide deploy preview %s", deploy.Name)},
			)
			return fmt.Errorf("deploy apply %s is guarded until clustered orchestration is implemented", target)
		}
		report, err := inspectEnvContract()
		if err != nil {
			return err
		}
		if len(report.missingRequired) > 0 {
			return fmt.Errorf("environment contract has %d missing required value(s)", len(report.missingRequired))
		}
		r.Rows(
			outputRow{"target", deploy.Name},
			outputRow{"type", deploy.Type},
			outputRow{"host", deploy.Host},
			outputRow{"remote", deploy.RemotePath},
		)
		if err := a.applySSHComposeDeploy(deploy, r); err != nil {
			return err
		}
		r.Row(outputRow{"status", "deployed"})
		return nil
	}

	r.Rows(
		outputRow{"target", target},
		outputRow{"status", "disabled"},
		outputRow{"reason", "no deploy target is implemented yet"},
		outputRow{"preview", fmt.Sprintf("carbide deploy preview %s", target)},
	)
	return fmt.Errorf("deploy apply %s is disabled until a deploy target exists", target)
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
		r.Rows(
			outputRow{"logs", "detached"},
			outputRow{"dev", "running"},
			outputRow{"follow", "carbide follow logs"},
			outputRow{"stop", "carbide stop dev"},
		)
		return nil
	}
	if first.err != nil {
		return fmt.Errorf("Docker Compose %s failed: %w", first.name, first.err)
	}
	return nil
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
	interactive := isTerminalOutput(out)
	termWidth := 0
	if interactive {
		termWidth = terminalColumns(out)
		if termWidth == 0 {
			termWidth = terminalColumnsFromEnv()
		}
	}
	return renderer{
		out:         out,
		interactive: interactive,
		styled:      interactive && os.Getenv("NO_COLOR") == "",
		termWidth:   termWidth,
	}
}

func renderError(out io.Writer, err error) {
	newRenderer(out).Message(
		"Carbide",
		"command failed",
		outputRow{"error", err.Error()},
		outputRow{"help", "carbide help"},
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

func (r renderer) Table(headers []string, rows []tableRow) {
	widths := make([]int, len(headers))
	for index, header := range headers {
		widths[index] = len(header)
	}
	for _, row := range rows {
		for index := range headers {
			value := ""
			if index < len(row) {
				value = row[index]
			}
			if len(value) > widths[index] {
				widths[index] = len(value)
			}
		}
	}

	writeCells := func(cells []string, header bool) {
		for index := range headers {
			value := ""
			if index < len(cells) {
				value = cells[index]
			}
			if index > 0 {
				fmt.Fprint(r.out, "  ")
			}
			padded := value
			if index < len(headers)-1 {
				padded += strings.Repeat(" ", widths[index]-len(value))
			}
			if header {
				padded = r.paint("2;38;5;245", padded)
			}
			fmt.Fprint(r.out, padded)
		}
		fmt.Fprintln(r.out)
	}

	writeCells(headers, true)
	for _, row := range rows {
		writeCells([]string(row), false)
	}
}

func (r renderer) CommandList(sections []helpCommandSection) {
	fmt.Fprintln(r.out, r.formatHelpHeading("Usage:"))
	fmt.Fprintln(r.out, "  carbide <command> [arguments]")
	fmt.Fprintln(r.out)
	fmt.Fprintln(r.out, r.formatHelpHeading("Available commands:"))

	width := helpCommandWidth(sections)
	for _, section := range sections {
		if section.name != "" {
			fmt.Fprintln(r.out, r.formatHelpGroup(section.name))
		}
		for _, row := range section.rows {
			r.writeHelpCommand(row, width)
		}
	}
}

func (r renderer) writeHelpCommand(row outputRow, width int) {
	if r.styled {
		key := r.paint("38;5;245", row.key)
		fmt.Fprintf(r.out, "  %s%s  %s\n", key, strings.Repeat(" ", width-len(row.key)), row.value)
		return
	}
	fmt.Fprintf(r.out, "  %-*s  %s\n", width, row.key, row.value)
}

func (r renderer) formatHelpHeading(value string) string {
	return r.paint("1;38;5;245", value)
}

func (r renderer) formatHelpGroup(value string) string {
	return r.paint("1;38;5;245", value)
}

func (r renderer) Row(row outputRow) {
	r.writeRow(row, len(row.key))
}

func (r renderer) Blank() {
	fmt.Fprintln(r.out)
}

func (r renderer) Logo(logo string) {
	for index, line := range logoLines(logo) {
		fmt.Fprintln(r.out, r.formatLogoLine(index, line))
	}
	fmt.Fprintln(r.out)
}

func (r renderer) AnimateLogo(logo string) {
	lines := logoLines(logo)
	if len(lines) == 0 {
		return
	}

	width := maxLineWidth(lines)
	chompFrames := width + (len(lines)-1)*2 + 1
	for frame := 0; frame <= chompFrames; frame++ {
		if frame > 0 {
			fmt.Fprintf(r.out, "\033[%dA", len(lines))
		}
		for index, line := range lines {
			position := frame - index*2
			fmt.Fprintf(r.out, "\r\033[K%s\n", r.formatLogoPacmanLine(line, position, frame+index))
		}
		if frame < chompFrames {
			time.Sleep(9 * time.Millisecond)
		}
	}
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
	r.LogAt(time.Now(), service, message)
}

func (r renderer) LogEntry(entry structuredLogEntry) {
	r.LogAt(entryTimestamp(entry), entry.Service, entry.Message)
}

func (r renderer) LogAt(timestamp time.Time, service string, message string) {
	label := service
	if label == "" {
		label = "log"
	}
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	stamp := timestamp.Local().Format("15:04:05")
	width := 9
	if len(label) > width {
		width = len(label)
	}
	if r.styled {
		fmt.Fprintf(
			r.out,
			"%s  %s%s  %s\n",
			r.paint("2;38;5;245", stamp),
			r.formatService(label),
			strings.Repeat(" ", width-len(label)),
			message,
		)
		return
	}
	fmt.Fprintf(r.out, "%s  %-*s  %s\n", stamp, width, label, message)
}

func (r renderer) RunServiceProgress(
	services []string,
	poll func() map[string]composeServiceStatus,
	work func() error,
) error {
	return r.runServiceProgress("start", services, poll, work)
}

func (r renderer) RunServiceStopProgress(
	services []string,
	poll func() map[string]composeServiceStatus,
	work func() error,
) error {
	return r.runServiceProgress("stop", services, poll, work)
}

func (r renderer) runServiceProgress(
	mode string,
	services []string,
	poll func() map[string]composeServiceStatus,
	work func() error,
) error {
	if !r.interactive {
		return work()
	}

	done := make(chan error, 1)
	go func() {
		done <- work()
	}()

	if len(services) == 0 {
		services = []string{"containers"}
	}
	statuses := map[string]composeServiceStatus{}
	step := 0
	r.writeServiceProgress(mode, services, statuses, step)
	ticker := time.NewTicker(140 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			if err != nil {
				for _, service := range services {
					statuses[service] = composeServiceStatus{service: service, state: "failed"}
				}
				r.rewriteServiceProgress(mode, services, statuses, step)
				return err
			}
			for _, service := range services {
				statuses[service] = composeServiceStatus{service: service, state: progressDoneState(mode), health: "healthy"}
			}
			r.rewriteServiceProgress(mode, services, statuses, step)
			return nil
		case <-ticker.C:
			step++
			if current := poll(); len(current) > 0 {
				statuses = current
			}
			r.rewriteServiceProgress(mode, services, statuses, step)
		}
	}
}

func (r renderer) rewriteServiceProgress(mode string, services []string, statuses map[string]composeServiceStatus, step int) {
	fmt.Fprintf(r.out, "\033[%dA", len(services))
	r.writeServiceProgress(mode, services, statuses, step)
}

func (r renderer) writeServiceProgress(mode string, services []string, statuses map[string]composeServiceStatus, step int) {
	width := rowTextWidth(services)
	frameWidth := r.serviceProgressFrameWidth(width)
	for index, service := range services {
		status := statuses[service]
		state := progressState(mode, status)
		frame := serviceProgressFrame(frameWidth, step+index, state)
		stateLabel := padRight(state, progressStateColumnWidth)
		fmt.Fprintf(
			r.out,
			"\r\033[K%s%s  %s %s\n",
			r.formatService(service),
			strings.Repeat(" ", width-len(service)),
			r.paint(serviceProgressColor(state), frame),
			r.paint("2;38;5;245", stateLabel),
		)
	}
}

func (r renderer) serviceProgressFrameWidth(serviceWidth int) int {
	termWidth := r.currentTerminalWidth()
	frameWidth := termWidth - serviceWidth - progressStateColumnWidth - 5
	if frameWidth < minimumProgressFrameWidth {
		return minimumProgressFrameWidth
	}
	return frameWidth
}

func (r renderer) currentTerminalWidth() int {
	if r.interactive {
		if width := terminalColumns(r.out); width > 0 {
			return width
		}
	}
	if r.termWidth > 0 {
		return r.termWidth
	}
	return defaultTerminalWidth
}

func progressDoneState(mode string) string {
	if mode == "stop" {
		return "stopped"
	}
	return "running"
}

func progressState(mode string, status composeServiceStatus) string {
	if mode == "stop" {
		return serviceStopProgressState(status)
	}
	return serviceProgressState(status)
}

func serviceProgressFrame(width int, step int, state string) string {
	if width < 1 {
		width = 1
	}
	switch state {
	case "ready":
		return "[" + strings.Repeat("#", width) + "]"
	case "stopped":
		return "[" + strings.Repeat(" ", width) + "]"
	case "failed":
		return "[" + strings.Repeat("!", width) + "]"
	case "stopping":
		return reversePacmanFrame(width, step)
	default:
		return pacmanFrame(width, step)
	}
}

func pacmanFrame(width int, step int) string {
	position := step % width
	if position < 0 {
		position = 0
	}
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < width; i++ {
		switch {
		case i == position:
			b.WriteByte(pacmanMouth(step, "right"))
		case i < position:
			b.WriteByte('-')
		case isCandyPosition(i - position):
			b.WriteByte('o')
		default:
			b.WriteByte(' ')
		}
	}
	b.WriteByte(']')
	return b.String()
}

func reversePacmanFrame(width int, step int) string {
	position := width - 1 - (step % width)
	if position < 0 {
		position = 0
	}
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < width; i++ {
		switch {
		case i == position:
			b.WriteByte(pacmanMouth(step, "left"))
		case i > position:
			b.WriteByte('-')
		case isCandyPosition(position - i):
			b.WriteByte('o')
		default:
			b.WriteByte(' ')
		}
	}
	b.WriteByte(']')
	return b.String()
}

func isCandyPosition(distance int) bool {
	return distance > 1 && (distance-2)%3 == 0
}

func pacmanMouth(step int, direction string) byte {
	open := step%2 == 0
	if direction == "left" {
		if open {
			return 'D'
		}
		return 'd'
	}
	if open {
		return 'C'
	}
	return 'c'
}

func serviceProgressState(status composeServiceStatus) string {
	state := strings.ToLower(status.state)
	health := strings.ToLower(status.health)
	if state == "failed" || state == "exited" || state == "dead" {
		return "failed"
	}
	if state == "running" && (health == "" || health == "healthy") {
		return "ready"
	}
	if health == "healthy" {
		return "ready"
	}
	return "starting"
}

func serviceStopProgressState(status composeServiceStatus) string {
	state := strings.ToLower(status.state)
	if state == "failed" || state == "dead" {
		return "failed"
	}
	if state == "stopped" {
		return "stopped"
	}
	return "stopping"
}

func serviceProgressColor(state string) string {
	switch state {
	case "ready":
		return "38;5;114"
	case "stopped":
		return "2;38;5;245"
	case "failed":
		return "38;5;203"
	default:
		return "38;5;81"
	}
}

func rowTextWidth(values []string) int {
	width := 0
	for _, value := range values {
		if len(value) > width {
			width = len(value)
		}
	}
	return width
}

func helpCommandWidth(sections []helpCommandSection) int {
	width := 0
	for _, section := range sections {
		for _, row := range section.rows {
			if len(row.key) > width {
				width = len(row.key)
			}
		}
	}
	return width
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

func clamp(value int, minValue int, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func (r renderer) formatKey(key string) string {
	return r.paint("2;38;5;245", key)
}

func (r renderer) formatService(service string) string {
	switch service {
	case "web":
		return r.paint("38;5;81", service)
	case "api":
		return r.paint("38;5;114", service)
	case "db":
		return r.paint("38;5;222", service)
	case "watch":
		return r.paint("38;5;147", service)
	default:
		return r.paint("38;5;245", service)
	}
}

func (r renderer) formatLogoLine(_ int, line string) string {
	if !r.styled {
		return line
	}

	var out strings.Builder
	activeColor := ""
	for i := 0; i < len(line); i++ {
		color := logoGlyphColor(line[i])
		if color != activeColor {
			if activeColor != "" {
				out.WriteString("\033[0m")
			}
			if color != "" {
				out.WriteString("\033[")
				out.WriteString(color)
				out.WriteString("m")
			}
			activeColor = color
		}
		out.WriteByte(line[i])
	}
	if activeColor != "" {
		out.WriteString("\033[0m")
	}
	return out.String()
}

func (r renderer) formatLogoPacmanLine(line string, position int, step int) string {
	width := len(line)
	if position >= width {
		return r.formatLogoLine(0, line)
	}

	var out strings.Builder
	if position > 0 {
		out.WriteString(r.formatLogoLine(0, visiblePrefix(line, position)))
	}
	start := clamp(position, 0, width)
	for column := start; column < width; column++ {
		switch {
		case column == position:
			out.WriteString(r.formatLogoChomper(pacmanMouth(step, "right")))
		case isCandyPosition(column - position):
			out.WriteString(r.formatLogoPellet())
		default:
			out.WriteByte(' ')
		}
	}
	return out.String()
}

func (r renderer) formatLogoChomper(ch byte) string {
	return r.paint("1;38;5;226", string(ch))
}

func (r renderer) formatLogoPellet() string {
	return r.paint("2;38;5;220", "o")
}

func logoGlyphColor(ch byte) string {
	switch ch {
	case '_':
		return "2;38;5;245"
	case 'o', 'O', '0':
		return "38;5;220"
	default:
		return ""
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

func logoLines(logo string) []string {
	logo = strings.TrimRight(logo, "\n")
	if strings.TrimSpace(logo) == "" {
		return nil
	}
	return strings.Split(logo, "\n")
}

func maxLineWidth(lines []string) int {
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	return width
}

func visiblePrefix(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if width >= len(value) {
		return value
	}
	return value[:width]
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
		entry := newStructuredLogEntry("compose-watch", stream, "watch", line)
		logSink.WriteEntry(entry)
		r.LogEntry(entry)
	}
}

func streamLogOutput(input io.Reader, r renderer, logSink *devLogSink, stream string, wg *sync.WaitGroup) {
	streamLogOutputWithQuery(input, r, logSink, stream, logQuery{}, wg)
}

func streamLogOutputWithQuery(input io.Reader, r renderer, logSink *devLogSink, stream string, query logQuery, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		service, message := parseComposeLogLine(line)
		entry := newStructuredLogEntry("compose-log", stream, service, message)
		logSink.WriteEntry(entry)
		if logEntryMatchesQuery(entry, query) {
			r.LogEntry(entry)
		}
	}
}

func newStructuredLogEntry(source string, stream string, service string, message string) structuredLogEntry {
	return structuredLogEntry{
		Time:    time.Now().UTC().Format(time.RFC3339Nano),
		Source:  source,
		Stream:  stream,
		Service: service,
		Message: message,
	}
}

func openDevLogSink(path string) (*devLogSink, error) {
	return openDevLogSinkWithFlags(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
}

func openAppendDevLogSink(path string) (*devLogSink, error) {
	return openDevLogSinkWithFlags(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND)
}

func openDevLogSinkWithFlags(path string, flags int) (*devLogSink, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, flags, 0644)
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
	s.WriteEntry(newStructuredLogEntry(source, stream, service, message))
}

func (s *devLogSink) WriteEntry(entry structuredLogEntry) {
	if s == nil || s.encoder == nil || strings.TrimSpace(entry.Message) == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.encoder.Encode(entry)
}

func (a app) commandLogs(args []string) error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}
	query, err := parseLogQuery(args)
	if err != nil {
		return err
	}
	entries, err := readStructuredLogEntries(devLogPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no dev logs found; run carbide run dev first")
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
		r.LogEntry(entry)
	}
	return nil
}

func (a app) commandFollowLogs(args []string) error {
	if !isFile("carbide.toml") {
		return errors.New("run this inside a Carbide project")
	}
	query, err := parseLogQuery(args)
	if err != nil {
		return err
	}
	if query.json {
		return errors.New("carbide follow logs does not support json")
	}
	if query.limit != 80 {
		return errors.New("carbide follow logs does not support limit")
	}

	compose, err := findCompose()
	if err != nil {
		return err
	}
	env := setEnv(os.Environ(), "COMPOSE_MENU", "false")
	env = composeEnv(env)
	logSink, err := openAppendDevLogSink(devLogPath)
	if err != nil {
		return err
	}
	defer logSink.Close()

	var streams sync.WaitGroup
	results := make(chan processResult, 1)
	process, err := a.startComposeStream(
		"logs",
		compose,
		env,
		composeLogsArgs(compose),
		func(input io.Reader, r renderer, sink *devLogSink, stream string, wg *sync.WaitGroup) {
			streamLogOutputWithQuery(input, r, sink, stream, query, wg)
		},
		logSink,
		&streams,
		results,
	)
	if err != nil {
		return err
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
		stopProcesses([]runningProcess{process}, sig)
	case first = <-results:
	}

	if interrupted {
		waitForProcesses(1, []runningProcess{process}, results, 5*time.Second)
		streams.Wait()
		r := newRenderer(a.stdout)
		r.Blank()
		r.Rows(outputRow{"logs", "detached"})
		return nil
	}
	streams.Wait()
	if first.err != nil {
		return fmt.Errorf("Docker Compose %s failed: %w", first.name, first.err)
	}
	return nil
}

func inspectEnvContract() (envContractReport, error) {
	schema, err := readEnvSchema(projectConfigPath)
	if err != nil {
		return envContractReport{}, err
	}
	dotenv, envFileFound, err := readDotenv(".env")
	if err != nil {
		return envContractReport{}, err
	}

	report := envContractReport{
		schema:       schema,
		envFileFound: envFileFound,
	}
	seen := map[string]bool{}
	for _, variable := range schema.Variables {
		name := strings.TrimSpace(variable.Name)
		if name == "" {
			report.warnings = append(report.warnings, "contract contains an unnamed variable")
			continue
		}
		if seen[name] {
			report.warnings = append(report.warnings, fmt.Sprintf("%s is declared more than once", name))
			continue
		}
		seen[name] = true
		if variable.Secret {
			report.secretCount++
		}
		if variable.BrowserExposed {
			report.browserCount++
		}
		if variable.FrameworkOwned {
			report.frameworkCount++
		}
		if variable.Secret && variable.BrowserExposed {
			report.warnings = append(report.warnings, fmt.Sprintf("%s is secret and browser-exposed", name))
		}
		if variable.Required && envValue(name, variable.LocalDefault, dotenv) == "" {
			report.missingRequired = append(report.missingRequired, name)
		}
	}
	return report, nil
}

func readEnvSchema(path string) (envSchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return envSchema{}, fmt.Errorf("missing %s", path)
		}
		return envSchema{}, err
	}
	schema, err := parseEnvContractTOML(string(data))
	if err != nil {
		return envSchema{}, fmt.Errorf("invalid %s: %w", path, err)
	}
	if schema.Version == 0 {
		return envSchema{}, fmt.Errorf("%s is missing env.contract_version", path)
	}
	if len(schema.Variables) == 0 {
		return envSchema{}, fmt.Errorf("%s declares no env variables", path)
	}
	return schema, nil
}

func parseEnvContractTOML(content string) (envSchema, error) {
	var schema envSchema
	var variables []envVariable
	var current *envVariable
	seen := map[string]bool{}
	section := ""

	flushVariable := func() error {
		if current == nil {
			return nil
		}
		current.Name = strings.TrimSpace(current.Name)
		if current.Name == "" {
			return errors.New("env variable table has an empty name")
		}
		if seen[current.Name] {
			return fmt.Errorf("%s is declared more than once", current.Name)
		}
		seen[current.Name] = true
		variables = append(variables, *current)
		current = nil
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(stripTomlComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if err := flushVariable(); err != nil {
				return envSchema{}, err
			}
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			if strings.HasPrefix(section, "env.variables.") {
				name := strings.TrimSpace(strings.TrimPrefix(section, "env.variables."))
				name = strings.Trim(name, `"`)
				current = &envVariable{Name: name}
			}
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch {
		case section == "env" && key == "contract_version":
			version, err := strconv.Atoi(value)
			if err != nil {
				return envSchema{}, fmt.Errorf("line %d has invalid env.contract_version", lineNumber)
			}
			schema.Version = version
		case strings.HasPrefix(section, "env.variables."):
			if current == nil {
				return envSchema{}, fmt.Errorf("line %d is outside an env variable table", lineNumber)
			}
			if err := assignEnvVariableField(current, key, value, lineNumber); err != nil {
				return envSchema{}, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return envSchema{}, err
	}
	if err := flushVariable(); err != nil {
		return envSchema{}, err
	}
	schema.Variables = variables
	return schema, nil
}

func assignEnvVariableField(variable *envVariable, key string, value string, lineNumber int) error {
	switch key {
	case "service":
		variable.Service = parseTomlString(value)
	case "required":
		variable.Required = parseTomlBool(value)
	case "secret":
		variable.Secret = parseTomlBool(value)
	case "browser_exposed":
		variable.BrowserExposed = parseTomlBool(value)
	case "framework_owned":
		variable.FrameworkOwned = parseTomlBool(value)
	case "local_default":
		variable.LocalDefault = parseTomlString(value)
	case "description":
		variable.Description = parseTomlString(value)
	default:
		return fmt.Errorf("line %d has unknown env variable field %q", lineNumber, key)
	}
	return nil
}

func parseTomlBool(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "true")
}

func parseTomlString(value string) string {
	return unquoteEnvValue(strings.TrimSpace(value))
}

func parseTomlStringArray(value string) ([]string, error) {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "[") || !strings.HasSuffix(value, "]") {
		return nil, errors.New("expected TOML string array")
	}
	body := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"))
	if body == "" {
		return []string{}, nil
	}

	var values []string
	for _, item := range strings.Split(body, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			return nil, errors.New("empty array item")
		}
		if !(strings.HasPrefix(item, `"`) && strings.HasSuffix(item, `"`)) &&
			!(strings.HasPrefix(item, `'`) && strings.HasSuffix(item, `'`)) {
			return nil, fmt.Errorf("%s is not quoted", item)
		}
		values = append(values, parseTomlString(item))
	}
	return values, nil
}

func stripTomlComment(line string) string {
	inString := false
	escaped := false
	for i, r := range line {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' && inString {
			escaped = true
			continue
		}
		if r == '"' {
			inString = !inString
			continue
		}
		if r == '#' && !inString {
			return line[:i]
		}
	}
	return line
}

func readDotenv(path string) (map[string]string, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, false, nil
		}
		return nil, false, err
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		values[key] = unquoteEnvValue(strings.TrimSpace(value))
	}
	if err := scanner.Err(); err != nil {
		return nil, false, err
	}
	return values, true, nil
}

func envValue(name string, localDefault string, dotenv map[string]string) string {
	if value, ok := dotenv[name]; ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return strings.TrimSpace(localDefault)
}

func unquoteEnvValue(value string) string {
	if len(value) < 2 {
		return value
	}
	if (value[0] == '"' && value[len(value)-1] == '"') ||
		(value[0] == '\'' && value[len(value)-1] == '\'') {
		return value[1 : len(value)-1]
	}
	return value
}

func ensureDeployTarget(target string) error {
	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(target) {
		return errors.New("deploy target must use lowercase letters, numbers, and dashes")
	}
	return nil
}

func loadDeployTarget(name string) (deployTarget, bool, error) {
	content, err := os.ReadFile(projectConfigPath)
	if err != nil {
		return deployTarget{}, false, err
	}

	target := deployTarget{
		Name:             name,
		Hosts:            map[string]deployHost{},
		SourcePath:       ".",
		ComposeFile:      composeFilePath,
		ProjectDirectory: ".",
		HealthPath:       "/health",
	}
	section := ""
	found := false
	roleIndexes := map[string]int{}
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(stripTomlComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if strings.HasPrefix(section, "deploy.hosts.") {
			hostName := strings.TrimSpace(strings.TrimPrefix(section, "deploy.hosts."))
			if hostName == "" || strings.Contains(hostName, ".") {
				return deployTarget{}, false, fmt.Errorf("line %d has invalid deploy host section", lineNumber)
			}
			host := target.Hosts[hostName]
			host.Name = hostName
			if err := assignDeployHostField(&host, key, value, lineNumber); err != nil {
				return deployTarget{}, false, err
			}
			target.Hosts[hostName] = host
			continue
		}

		targetSection := "deploy.targets." + name
		if section == targetSection {
			found = true
			if err := assignDeployTargetField(&target, key, value, lineNumber); err != nil {
				return deployTarget{}, false, err
			}
			continue
		}

		rolePrefix := targetSection + ".roles."
		if strings.HasPrefix(section, rolePrefix) {
			found = true
			roleName := strings.TrimSpace(strings.TrimPrefix(section, rolePrefix))
			if roleName == "" || strings.Contains(roleName, ".") {
				return deployTarget{}, false, fmt.Errorf("line %d has invalid deploy role section", lineNumber)
			}
			index, ok := roleIndexes[roleName]
			if !ok {
				index = len(target.Roles)
				roleIndexes[roleName] = index
				target.Roles = append(target.Roles, deployRole{Name: roleName})
			}
			if err := assignDeployRoleField(&target.Roles[index], key, value, lineNumber); err != nil {
				return deployTarget{}, false, err
			}
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		return deployTarget{}, false, err
	}
	return target, found, nil
}

func assignDeployTargetField(target *deployTarget, key string, value string, lineNumber int) error {
	switch key {
	case "type":
		target.Type = parseTomlString(value)
	case "host":
		target.Host = parseTomlString(value)
	case "domain":
		target.Domain = parseTomlString(value)
	case "remote_path":
		target.RemotePath = parseTomlString(value)
	case "source_path":
		target.SourcePath = parseTomlString(value)
	case "compose_file":
		target.ComposeFile = parseTomlString(value)
	case "project_directory":
		target.ProjectDirectory = parseTomlString(value)
	case "public_port":
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("line %d has invalid public_port", lineNumber)
		}
		target.PublicPort = port
	case "health_path":
		target.HealthPath = parseTomlString(value)
	case "nginx":
		target.Nginx = parseTomlBool(value)
	case "nginx_site":
		target.NginxSite = parseTomlString(value)
	case "strategy":
		target.Strategy = parseTomlString(value)
	default:
		return fmt.Errorf("line %d has unknown deploy target field %q", lineNumber, key)
	}
	return nil
}

func assignDeployHostField(host *deployHost, key string, value string, lineNumber int) error {
	switch key {
	case "ssh":
		host.SSH = parseTomlString(value)
	case "address":
		host.Address = parseTomlString(value)
	case "description":
		host.Description = parseTomlString(value)
	default:
		return fmt.Errorf("line %d has unknown deploy host field %q", lineNumber, key)
	}
	return nil
}

func assignDeployRoleField(role *deployRole, key string, value string, lineNumber int) error {
	switch key {
	case "hosts":
		hosts, err := parseTomlStringArray(value)
		if err != nil {
			return fmt.Errorf("line %d has invalid hosts: %w", lineNumber, err)
		}
		role.Hosts = hosts
	case "remote_path":
		role.RemotePath = parseTomlString(value)
	case "compose_file":
		role.ComposeFile = parseTomlString(value)
	case "project_directory":
		role.ProjectDirectory = parseTomlString(value)
	case "public_port":
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("line %d has invalid public_port", lineNumber)
		}
		role.PublicPort = port
	case "health_path":
		role.HealthPath = parseTomlString(value)
	case "nginx":
		role.Nginx = parseTomlBool(value)
	case "primary":
		role.Primary = parseTomlString(value)
	case "migration":
		role.Migration = parseTomlString(value)
	default:
		return fmt.Errorf("line %d has unknown deploy role field %q", lineNumber, key)
	}
	return nil
}

func validateDeployTarget(target deployTarget) error {
	switch target.Type {
	case "ssh-compose":
		return validateSSHComposeTarget(target)
	case "ssh-compose-environment":
		return validateSSHComposeEnvironmentTarget(target)
	default:
		return fmt.Errorf("deploy target %s has unsupported type %q", target.Name, target.Type)
	}
}

func validateSSHComposeTarget(target deployTarget) error {
	if target.Host == "" {
		return fmt.Errorf("deploy target %s is missing host", target.Name)
	}
	if target.RemotePath == "" || !strings.HasPrefix(target.RemotePath, "/") {
		return fmt.Errorf("deploy target %s remote_path must be absolute", target.Name)
	}
	if strings.ContainsAny(target.RemotePath, "\r\n") {
		return fmt.Errorf("deploy target %s remote_path is invalid", target.Name)
	}
	if target.PublicPort < 1 || target.PublicPort > 65535 {
		return fmt.Errorf("deploy target %s public_port must be between 1 and 65535", target.Name)
	}
	if target.HealthPath == "" || !strings.HasPrefix(target.HealthPath, "/") {
		return fmt.Errorf("deploy target %s health_path must start with /", target.Name)
	}
	if target.Domain != "" && !regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`).MatchString(target.Domain) {
		return fmt.Errorf("deploy target %s has invalid domain", target.Name)
	}
	if target.Nginx && target.Domain == "" {
		return fmt.Errorf("deploy target %s requires domain when nginx is enabled", target.Name)
	}
	if target.NginxSite != "" && !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(target.NginxSite) {
		return fmt.Errorf("deploy target %s nginx_site must use lowercase letters, numbers, and dashes", target.Name)
	}
	return nil
}

func validateSSHComposeEnvironmentTarget(target deployTarget) error {
	if len(target.Hosts) == 0 {
		return fmt.Errorf("deploy target %s requires [deploy.hosts.*] entries", target.Name)
	}
	for name, host := range target.Hosts {
		if host.SSH == "" {
			return fmt.Errorf("deploy host %s is missing ssh", name)
		}
		if strings.ContainsAny(host.SSH, "\r\n") {
			return fmt.Errorf("deploy host %s ssh is invalid", name)
		}
	}
	if len(target.Roles) == 0 {
		return fmt.Errorf("deploy target %s requires role tables", target.Name)
	}

	roleNames := map[string]bool{}
	for _, role := range target.Roles {
		roleNames[role.Name] = true
		if len(role.Hosts) == 0 {
			return fmt.Errorf("deploy role %s is missing hosts", role.Name)
		}
		for _, hostName := range role.Hosts {
			if _, ok := target.Hosts[hostName]; !ok {
				return fmt.Errorf("deploy role %s references unknown host %s", role.Name, hostName)
			}
		}
		if role.RemotePath != "" && !strings.HasPrefix(role.RemotePath, "/") {
			return fmt.Errorf("deploy role %s remote_path must be absolute", role.Name)
		}
		if role.PublicPort < 0 || role.PublicPort > 65535 {
			return fmt.Errorf("deploy role %s public_port must be between 0 and 65535", role.Name)
		}
		if role.HealthPath != "" && !strings.HasPrefix(role.HealthPath, "/") {
			return fmt.Errorf("deploy role %s health_path must start with /", role.Name)
		}
		if role.Nginx && target.Domain == "" {
			return fmt.Errorf("deploy role %s requires target domain when nginx is enabled", role.Name)
		}
	}
	for _, required := range []string{"web", "api", "db"} {
		if !roleNames[required] {
			return fmt.Errorf("deploy target %s is missing %s role", target.Name, required)
		}
	}
	db, ok := deployRoleByName(target.Roles, "db")
	if !ok {
		return fmt.Errorf("deploy target %s is missing db role", target.Name)
	}
	if db.Primary == "" {
		return fmt.Errorf("deploy role db requires primary host")
	}
	if !containsString(db.Hosts, db.Primary) {
		return fmt.Errorf("deploy role db primary %s is not in db hosts", db.Primary)
	}
	if db.Migration != "once" {
		return fmt.Errorf("deploy role db migration must be once")
	}
	if target.Domain != "" && !regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`).MatchString(target.Domain) {
		return fmt.Errorf("deploy target %s has invalid domain", target.Name)
	}
	return nil
}

func deployRoleByName(roles []deployRole, name string) (deployRole, bool) {
	for _, role := range roles {
		if role.Name == name {
			return role, true
		}
	}
	return deployRole{}, false
}

func deployHostRows(target deployTarget) []string {
	names := make([]string, 0, len(target.Hosts))
	for name := range target.Hosts {
		names = append(names, name)
	}
	sort.Strings(names)

	rows := make([]string, 0, len(names))
	for _, name := range names {
		host := target.Hosts[name]
		value := host.SSH
		if host.Address != "" {
			value += " " + host.Address
		}
		rows = append(rows, fmt.Sprintf("%s -> %s", name, value))
	}
	return rows
}

func deployRoleRows(target deployTarget) []string {
	roles := append([]deployRole(nil), target.Roles...)
	sort.Slice(roles, func(i int, j int) bool {
		return roles[i].Name < roles[j].Name
	})

	rows := make([]string, 0, len(roles))
	for _, role := range roles {
		parts := []string{fmt.Sprintf("%s: %s", role.Name, strings.Join(role.Hosts, ","))}
		if role.PublicPort > 0 {
			parts = append(parts, fmt.Sprintf("port %d", role.PublicPort))
		}
		if role.Nginx {
			parts = append(parts, "nginx")
		}
		if role.Primary != "" {
			parts = append(parts, "primary "+role.Primary)
		}
		if role.Migration != "" {
			parts = append(parts, "migrate "+role.Migration)
		}
		rows = append(rows, strings.Join(parts, " "))
	}
	return rows
}

func (a app) applySSHComposeDeploy(target deployTarget, r renderer) error {
	if _, err := exec.LookPath("ssh"); err != nil {
		return errors.New("ssh is required for ssh-compose deploys")
	}
	if _, err := exec.LookPath("rsync"); err != nil {
		return errors.New("rsync is required for ssh-compose deploys")
	}

	source, err := filepath.Abs(target.SourcePath)
	if err != nil {
		return err
	}
	if !isDir(source) {
		return fmt.Errorf("deploy source does not exist: %s", source)
	}

	r.Row(outputRow{"prepare", "remote directory"})
	if _, err := runSSHScript(target.Host, fmt.Sprintf(
		"set -euo pipefail\nsudo mkdir -p %s\nsudo chown \"$(id -u):$(id -g)\" %s\n",
		shellQuote(target.RemotePath),
		shellQuote(target.RemotePath),
	)); err != nil {
		return err
	}

	r.Row(outputRow{"sync", source})
	sourceWithSlash := source + string(os.PathSeparator)
	remoteDest := target.Host + ":" + target.RemotePath + "/"
	if _, err := commandOutput("", "rsync",
		"-az",
		"--delete",
		"--exclude", ".git",
		"--exclude", ".carbide",
		"--exclude", ".env",
		"--exclude", "app/.env",
		sourceWithSlash,
		remoteDest,
	); err != nil {
		return err
	}

	r.Row(outputRow{"env", "ensure remote .env"})
	envContent, err := deployEnvContent(target)
	if err != nil {
		return err
	}
	if _, err := runSSHScriptInput(target.Host, fmt.Sprintf(
		"set -euo pipefail\ncd %s\nif [ ! -f .env ]; then umask 077; cat > .env; else cat >/dev/null; fi\n",
		shellQuote(target.RemotePath),
	), envContent); err != nil {
		return err
	}

	compose := remoteComposeCommand(target)
	r.Row(outputRow{"compose", "config"})
	if _, err := runSSHScript(target.Host, fmt.Sprintf(
		"set -euo pipefail\ncd %s\n%s config >/dev/null\n",
		shellQuote(target.RemotePath),
		compose,
	)); err != nil {
		return err
	}

	r.Row(outputRow{"compose", "up -d --build"})
	if _, err := runSSHScript(target.Host, fmt.Sprintf(
		"set -euo pipefail\ncd %s\n%s up -d --build --remove-orphans\n",
		shellQuote(target.RemotePath),
		compose,
	)); err != nil {
		return err
	}

	if target.Nginx {
		r.Row(outputRow{"nginx", target.Domain})
		if err := installNginxSite(target); err != nil {
			return err
		}
	}

	healthURL := fmt.Sprintf("http://127.0.0.1:%d%s", target.PublicPort, target.HealthPath)
	r.Row(outputRow{"health", healthURL})
	if _, err := runSSHScript(target.Host, fmt.Sprintf(
		"set -euo pipefail\nfor i in $(seq 1 40); do curl -fsS --max-time 5 %s >/dev/null && exit 0; sleep 2; done\ncurl -fsS --max-time 10 %s >/dev/null\n",
		shellQuote(healthURL),
		shellQuote(healthURL),
	)); err != nil {
		return err
	}
	return nil
}

func deployEnvContent(target deployTarget) (string, error) {
	password, err := randomHex(24)
	if err != nil {
		return "", err
	}
	composeName, appName := deployProjectNames()
	publicURL := fmt.Sprintf("http://127.0.0.1:%d", target.PublicPort)
	if target.Domain != "" {
		scheme := "http"
		if target.Nginx {
			scheme = "https"
		}
		publicURL = scheme + "://" + target.Domain
	}
	return fmt.Sprintf(`COMPOSE_PROJECT_NAME=%s
APP_ENV=production
CARBIDE_HTTP_PORT=%d
PUBLIC_APP_NAME=%s
PUBLIC_URL=%s
POSTGRES_PASSWORD=%s
DATABASE_URL=postgres://carbide:%s@db:5432/carbide?sslmode=disable
`, composeName, target.PublicPort, appName, publicURL, password, password), nil
}

func deployProjectNames() (string, string) {
	metadata, err := readProjectMetadata(projectConfigPath)
	if err != nil {
		return "carbide-app", "Carbide App"
	}

	composeName := normalizeComposeProjectName(metadata.slug)
	if composeName == "" {
		composeName = normalizeComposeProjectName(metadata.name)
	}
	if composeName == "" {
		composeName = "carbide-app"
	}

	appName := strings.TrimSpace(metadata.name)
	if appName == "" || isTemplatePlaceholder(appName) || strings.ContainsAny(appName, "\r\n") {
		appName = projectDisplayName(composeName)
	}
	return composeName, appName
}

func remoteComposeCommand(target deployTarget) string {
	return fmt.Sprintf(
		"docker compose --env-file .env -f %s --project-directory %s",
		shellQuote(target.ComposeFile),
		shellQuote(target.ProjectDirectory),
	)
}

func installNginxSite(target deployTarget) error {
	name := target.NginxSite
	if name == "" {
		name = "carbide-" + target.Name
	}
	available := "/etc/nginx/sites-available/" + name
	enabled := "/etc/nginx/sites-enabled/" + name
	certPath := "/etc/letsencrypt/live/" + target.Domain + "/fullchain.pem"
	updateProxyPass := fmt.Sprintf(
		`s#proxy_pass http://127\.0\.0\.1:[0-9]+;#proxy_pass http://127.0.0.1:%d;#g`,
		target.PublicPort,
	)
	config := fmt.Sprintf(`server {
    listen 80;
    listen [::]:80;
    server_name %s;

    location / {
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://127.0.0.1:%d;
    }
}
`, target.Domain, target.PublicPort)

	script := fmt.Sprintf(`set -euo pipefail
sudo -n true
if [ -f %s ] && sudo grep -q %s %s; then
  cat >/dev/null
  sudo perl -0pi -e %s %s
else
  tmp="$(mktemp)"
  cat > "$tmp"
  sudo mv "$tmp" %s
  sudo ln -sf %s %s
fi
sudo nginx -t
sudo systemctl reload nginx
`,
		shellQuote(available),
		shellQuote(certPath),
		shellQuote(available),
		shellQuote(updateProxyPass),
		shellQuote(available),
		shellQuote(available),
		shellQuote(available),
		shellQuote(enabled),
	)
	_, err := runSSHScriptInput(target.Host, script, config)
	return err
}

func runSSHScript(host string, script string) (string, error) {
	return commandOutput("", "ssh", host, "bash", "-lc", script)
}

func runSSHScriptInput(host string, script string, input string) (string, error) {
	return commandOutputInput("", input, "ssh", host, "bash", "-lc", script)
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func randomHex(bytesCount int) (string, error) {
	buf := make([]byte, bytesCount)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func entryTimestamp(entry structuredLogEntry) time.Time {
	timestamp, err := time.Parse(time.RFC3339Nano, entry.Time)
	if err != nil {
		return time.Now()
	}
	return timestamp
}

func parseLogQuery(args []string) (logQuery, error) {
	query := logQuery{limit: 80}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "service":
			i++
			if i >= len(args) || args[i] == "" {
				return query, errors.New("usage: carbide logs [service <name>] [containing <text>] [limit <count>] [json]")
			}
			query.service = args[i]
		case "containing":
			i++
			if i >= len(args) || args[i] == "" {
				return query, errors.New("usage: carbide logs [service <name>] [containing <text>] [limit <count>] [json]")
			}
			query.contains = args[i]
		case "limit":
			i++
			if i >= len(args) {
				return query, errors.New("usage: carbide logs [service <name>] [containing <text>] [limit <count>] [json]")
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
	for _, entry := range entries {
		if logEntryMatchesQuery(entry, query) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func logEntryMatchesQuery(entry structuredLogEntry, query logQuery) bool {
	if query.service != "" && entry.Service != query.service {
		return false
	}
	if query.contains != "" && !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(query.contains)) {
		return false
	}
	return true
}

func limitLogEntries(entries []structuredLogEntry, limit int) []structuredLogEntry {
	if limit < 1 || len(entries) <= limit {
		return entries
	}
	return entries[len(entries)-limit:]
}

func carbideLogo() string {
	return defaultLogoText
}

func resolveHome() (string, error) {
	if home := os.Getenv("CARBIDE_HOME"); home != "" {
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
	case "bin":
		parent := filepath.Dir(dir)
		switch filepath.Base(parent) {
		case ".cli", "cli":
			return filepath.Dir(parent), nil
		default:
			return parent, nil
		}
	case ".bin":
		return filepath.Dir(dir), nil
	default:
		return filepath.Dir(dir), nil
	}
}

func shouldStyleOutput(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return isTerminalOutput(w)
}

func isTerminalOutput(w io.Writer) bool {
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

type terminalWindowSize struct {
	rows    uint16
	columns uint16
	xpixels uint16
	ypixels uint16
}

func terminalColumns(w io.Writer) int {
	file, ok := w.(*os.File)
	if !ok {
		return 0
	}
	var size terminalWindowSize
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		file.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&size)),
	)
	if errno != 0 || size.columns == 0 {
		return 0
	}
	return int(size.columns)
}

func terminalColumnsFromEnv() int {
	columns, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err != nil || columns <= 0 {
		return 0
	}
	return columns
}

func ensureProjectName(name string) error {
	if name == "" || strings.HasPrefix(name, ".") || strings.ContainsAny(name, `/\`) {
		return errors.New("project name must be a simple directory name")
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9._ -]+$`, name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("project name may contain only letters, numbers, spaces, dots, underscores, and dashes")
	}
	return nil
}

func projectDisplayName(input string) string {
	words := projectWords(input)
	if len(words) == 0 {
		return "My Carbide App"
	}

	for i, word := range words {
		words[i] = titleWord(word)
	}
	return strings.Join(words, " ")
}

func projectWords(input string) []string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ' ' || r == '-' || r == '_' || r == '.'
	})
	words := make([]string, 0, len(parts))
	for _, part := range parts {
		word := strings.TrimSpace(part)
		if word != "" {
			words = append(words, word)
		}
	}
	return words
}

func titleWord(word string) string {
	lower := strings.ToLower(word)
	runes := []rune(lower)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
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

type projectMetadata struct {
	name string
	slug string
}

func projectProfile() string {
	data, err := os.ReadFile(projectConfigPath)
	if err != nil {
		return "app"
	}

	section := ""
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(stripTomlComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			continue
		}
		if section != "" {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) != "profile" {
			continue
		}
		if profile := strings.TrimSpace(parseTomlString(value)); profile != "" {
			return profile
		}
	}
	return "app"
}

func readProjectMetadata(path string) (projectMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return projectMetadata{}, err
	}

	var metadata projectMetadata
	section := ""
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(stripTomlComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			continue
		}
		if section != "" {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "name":
			metadata.name = parseTomlString(value)
		case "slug":
			metadata.slug = parseTomlString(value)
		}
	}
	if err := scanner.Err(); err != nil {
		return projectMetadata{}, err
	}
	return metadata, nil
}

func composeProjectName() string {
	metadata, err := readProjectMetadata(projectConfigPath)
	if err == nil {
		if slug := normalizeComposeProjectName(metadata.slug); slug != "" {
			return slug
		}
	}

	pwd, err := os.Getwd()
	if err == nil {
		if slug := normalizeComposeProjectName(filepath.Base(pwd)); slug != "" {
			return slug
		}
	}
	return "carbide-app"
}

func normalizeComposeProjectName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || isTemplatePlaceholder(value) {
		return ""
	}
	return projectSlug(value)
}

func isTemplatePlaceholder(value string) bool {
	return strings.Contains(value, "__PROJECT_")
}

func (a app) copyScaffold(target string, name string, slug string) error {
	source := filepath.Join(a.home, "scaffold")
	if !isDir(source) {
		return fmt.Errorf("missing scaffold source: %s", source)
	}
	return copyScaffoldPart(source, target, name, slug)
}

func copyScaffoldPart(source string, target string, name string, slug string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return os.MkdirAll(target, 0755)
		}
		if entry.IsDir() && rel == ".carbide" {
			return filepath.SkipDir
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
	return composeCommand{}, errors.New("Docker Compose is required for carbide run dev")
}

func composeEnv(env []string) []string {
	if isFile(composeFilePath) {
		env = setEnv(env, "COMPOSE_FILE", composeFilePath)
	} else if isFile(legacyComposeFilePath) {
		env = setEnv(env, "COMPOSE_FILE", legacyComposeFilePath)
	}
	if isFile(projectConfigPath) {
		env = setEnv(env, "COMPOSE_PROJECT_NAME", composeProjectName())
	}
	return env
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

func composeServices(compose composeCommand, env []string) []string {
	output, err := runComposeCaptured(compose, env, "config", "--services")
	if err != nil {
		return defaultComposeServices()
	}
	seen := map[string]bool{}
	var services []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		service := strings.TrimSpace(scanner.Text())
		if service == "" || seen[service] {
			continue
		}
		seen[service] = true
		services = append(services, service)
	}
	if len(services) == 0 {
		return defaultComposeServices()
	}
	return services
}

func defaultComposeServices() []string {
	return []string{"web", "api", "db"}
}

func composeServiceStatuses(compose composeCommand, env []string) map[string]composeServiceStatus {
	snapshots, err := composeServiceSnapshots(compose, env)
	if err != nil {
		return nil
	}
	return composeStatusesFromSnapshots(snapshots)
}

func composeServiceSnapshots(compose composeCommand, env []string) (map[string]composeServiceSnapshot, error) {
	output, err := runComposeCaptured(compose, env, "ps", "--format", "json")
	if err != nil {
		return nil, err
	}
	return parseComposeServiceSnapshots(output)
}

func parseComposeServiceStatuses(output string) (map[string]composeServiceStatus, error) {
	snapshots, err := parseComposeServiceSnapshots(output)
	if err != nil {
		return nil, err
	}
	return composeStatusesFromSnapshots(snapshots), nil
}

func parseComposeServiceSnapshots(output string) (map[string]composeServiceSnapshot, error) {
	parseRecords := func(records []composeServiceSnapshot) map[string]composeServiceSnapshot {
		snapshots := map[string]composeServiceSnapshot{}
		for _, record := range records {
			service := strings.TrimSpace(record.Service)
			if service == "" {
				continue
			}
			record.Service = service
			record.Name = strings.TrimSpace(record.Name)
			record.State = strings.TrimSpace(record.State)
			record.Health = strings.TrimSpace(record.Health)
			record.Status = strings.TrimSpace(record.Status)
			record.Ports = strings.TrimSpace(record.Ports)
			snapshots[service] = record
		}
		return snapshots
	}

	var records []composeServiceSnapshot
	if err := json.Unmarshal([]byte(output), &records); err == nil {
		return parseRecords(records), nil
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var record composeServiceSnapshot
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return parseRecords(records), nil
}

func composeStatusesFromSnapshots(snapshots map[string]composeServiceSnapshot) map[string]composeServiceStatus {
	statuses := map[string]composeServiceStatus{}
	for service, snapshot := range snapshots {
		statuses[service] = composeServiceStatus{
			service: service,
			state:   snapshot.State,
			health:  snapshot.Health,
		}
	}
	return statuses
}

func composeStatusRow(snapshot composeServiceSnapshot) tableRow {
	return tableRow{
		snapshot.Service,
		statusValue(snapshot.Name),
		composePublishedPorts(snapshot),
		composeInternalPorts(snapshot),
		composeServiceStatusText(snapshot),
	}
}

func statusValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

func composePublishedPorts(snapshot composeServiceSnapshot) string {
	seen := map[string]bool{}
	var ports []string
	for _, publisher := range snapshot.Publishers {
		if publisher.PublishedPort <= 0 {
			continue
		}
		host := strings.TrimSpace(publisher.URL)
		if host == "" || host == "0.0.0.0" || host == "::" {
			host = "localhost"
		}
		value := fmt.Sprintf("%s:%d", host, publisher.PublishedPort)
		if !seen[value] {
			seen[value] = true
			ports = append(ports, value)
		}
	}
	if len(ports) == 0 {
		return "-"
	}
	return strings.Join(ports, ", ")
}

func composeInternalPorts(snapshot composeServiceSnapshot) string {
	seen := map[string]bool{}
	var ports []string
	for _, publisher := range snapshot.Publishers {
		if publisher.TargetPort <= 0 {
			continue
		}
		protocol := strings.TrimSpace(publisher.Protocol)
		if protocol == "" {
			protocol = "tcp"
		}
		value := fmt.Sprintf("%d/%s", publisher.TargetPort, protocol)
		if !seen[value] {
			seen[value] = true
			ports = append(ports, value)
		}
	}
	if len(ports) == 0 && strings.TrimSpace(snapshot.Ports) != "" {
		return strings.TrimSpace(snapshot.Ports)
	}
	if len(ports) == 0 {
		return "-"
	}
	return strings.Join(ports, ", ")
}

func composeServiceStatusText(snapshot composeServiceSnapshot) string {
	state := strings.TrimSpace(strings.ToLower(snapshot.State))
	health := strings.TrimSpace(strings.ToLower(snapshot.Health))
	if state == "" {
		return statusValue(snapshot.Status)
	}
	if health != "" {
		return fmt.Sprintf("%s (%s)", state, health)
	}
	return state
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
	case "web", "api", "db":
		return name
	default:
		return name
	}
}

func stripANSI(value string) string {
	var out strings.Builder
	inEscape := false
	inCSI := false
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if inEscape {
			if !inCSI && ch == '[' {
				inCSI = true
				continue
			}
			if inCSI {
				if ch >= '@' && ch <= '~' {
					inEscape = false
					inCSI = false
				}
				continue
			}
			if ch >= 0x30 && ch <= 0x7e {
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
		return 0, errors.New("CARBIDE_HTTP_PORT must be a number from 1 to 65535")
	}
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return 0, errors.New("CARBIDE_HTTP_PORT must be a number from 1 to 65535")
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
			return 0, fmt.Errorf("port %d is already in use; choose another with CARBIDE_HTTP_PORT=<port> carbide run dev", port)
		}
		return port, nil
	}

	for _, port := range []int{8080, 8081, 8082, 8083, 8084, 8085, 18080, 18081, 18082, 18083, 18084, 18085} {
		if portIsAvailable(port) {
			return port, nil
		}
	}
	return 0, errors.New("no free dev port found; run with CARBIDE_HTTP_PORT=<port> carbide run dev")
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
		return errors.New("Go is required to build the Carbide CLI")
	}

	outDir := filepath.Join(home, ".cli", "bin")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	finalPath := filepath.Join(outDir, "carbide")
	tmpPath := filepath.Join(outDir, fmt.Sprintf(".carbide-%d", os.Getpid()))
	ldflags := "-X github.com/ryangerardwilson/carbide/cli/internal/cli.commit=" + gitShortHead(home)

	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", tmpPath, "./cmd/carbide")
	cmd.Dir = filepath.Join(home, "cli")
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
	return commandOutputEnv(dir, nil, name, args...)
}

func commandOutputEnv(dir string, env []string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if env != nil {
		cmd.Env = env
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

func commandOutputInput(dir string, input string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdin = strings.NewReader(input)
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
