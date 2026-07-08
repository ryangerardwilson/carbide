package cli

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"sync"
)

var version = "0.2.0"
var commit = ""

const devLogPath = ".carbide/log/dev.jsonl"
const auditRootDir = ".audit"
const auditReportDirName = "report"
const auditStarterDirName = "starter-reference"
const auditPlanFileName = "plan.md"
const projectConfigPath = "carbide.toml"
const composeFilePath = "docker-compose.yml"
const legacyComposeFilePath = "compose.yml"
const defaultTerminalWidth = 80
const helpOutputMaxWidth = 79
const progressStateColumnWidth = 8
const minimumProgressFrameWidth = 4
const maxLawFileLines = 1000

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

type auditSpec struct {
	ref         string
	title       string
	fileName    string
	description string
}

type auditStageResult struct {
	root         string
	reportDir    string
	starterDir   string
	reportCount  int
	reportsReady bool
}

type resolveStageResult struct {
	planPath string
	ready    bool
}

type auditClarification struct {
	question string
	answer   string
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
	Name        string
	Script      string
	Description string
}

type jsonCommandReport struct {
	OK       bool                `json:"ok"`
	Command  string              `json:"command"`
	Version  string              `json:"version"`
	Commit   string              `json:"commit,omitempty"`
	Project  *projectJSONReport  `json:"project,omitempty"`
	URLs     *urlsJSONReport     `json:"urls,omitempty"`
	Env      *envJSONReport      `json:"env,omitempty"`
	Checks   []healthJSONCheck   `json:"checks,omitempty"`
	Services []serviceJSONReport `json:"services,omitempty"`
	Errors   []string            `json:"errors,omitempty"`
	Next     []string            `json:"next,omitempty"`
}

type projectJSONReport struct {
	Name    string `json:"name,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Profile string `json:"profile,omitempty"`
}

type urlsJSONReport struct {
	App    string `json:"app"`
	API    string `json:"api"`
	Source string `json:"source"`
}

type envJSONReport struct {
	FileFound       bool     `json:"file_found"`
	Status          string   `json:"status"`
	MissingRequired []string `json:"missing_required"`
	Warnings        []string `json:"warnings"`
	Secrets         int      `json:"secrets"`
	BrowserExposed  int      `json:"browser_exposed"`
	FrameworkOwned  int      `json:"framework_owned"`
}

type healthJSONCheck struct {
	Check  string `json:"check"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

type serviceJSONReport struct {
	Service        string   `json:"service"`
	Container      string   `json:"container"`
	PublishedPorts []string `json:"published_ports"`
	InternalPorts  []string `json:"internal_ports"`
	State          string   `json:"state"`
	Health         string   `json:"health,omitempty"`
	Status         string   `json:"status"`
}
