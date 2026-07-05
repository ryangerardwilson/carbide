package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestProjectSlug(t *testing.T) {
	tests := map[string]string{
		"Demo":           "demo",
		"my_app.test":    "my-app-test",
		"  Weird Name  ": "weird-name",
		"already--clean": "already-clean",
		"My Carbide App": "my-carbide-app",
		"my carbide app": "my-carbide-app",
		"My_Carbide.App": "my-carbide-app",
		"___":            "",
	}

	for input, want := range tests {
		if got := projectSlug(input); got != want {
			t.Fatalf("projectSlug(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestProjectDisplayName(t *testing.T) {
	tests := map[string]string{
		"demo":           "Demo",
		"my-carbide-app": "My Carbide App",
		"My Carbide App": "My Carbide App",
		"my_carbide.app": "My Carbide App",
		"___":            "My Carbide App",
	}

	for input, want := range tests {
		if got := projectDisplayName(input); got != want {
			t.Fatalf("projectDisplayName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestEnsureProjectName(t *testing.T) {
	valid := []string{"demo", "demo_app", "demo-app", "demo.app", "Demo1", "My Carbide App"}
	for _, name := range valid {
		if err := ensureProjectName(name); err != nil {
			t.Fatalf("ensureProjectName(%q) returned %v", name, err)
		}
	}

	invalid := []string{"", ".hidden", "nested/app", "bad*name"}
	for _, name := range invalid {
		if err := ensureProjectName(name); err == nil {
			t.Fatalf("ensureProjectName(%q) should fail", name)
		}
	}
}

func TestValidatePort(t *testing.T) {
	for _, value := range []string{"", "0", "65536", "abc"} {
		if _, err := validatePort(value); err == nil {
			t.Fatalf("validatePort(%q) should fail", value)
		}
	}

	got, err := validatePort("8080")
	if err != nil {
		t.Fatalf("validatePort returned %v", err)
	}
	if got != 8080 {
		t.Fatalf("validatePort returned %d, want 8080", got)
	}
}

func TestComposeEnvUsesProjectSlug(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		if err := os.WriteFile("carbide.toml", []byte(`name = "Demo App"
slug = "demo-app"
`), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}
		if err := os.WriteFile("docker-compose.yml", []byte("services: {}\n"), 0644); err != nil {
			t.Fatalf("WriteFile docker-compose.yml returned %v", err)
		}

		env := composeEnv(nil)
		if got := envSliceValue(env, "COMPOSE_PROJECT_NAME"); got != "demo-app" {
			t.Fatalf("COMPOSE_PROJECT_NAME = %q, want demo-app", got)
		}
		if got := envSliceValue(env, "COMPOSE_FILE"); got != "docker-compose.yml" {
			t.Fatalf("COMPOSE_FILE = %q, want docker-compose.yml", got)
		}
	})
}

func TestComposeEnvFallsBackForTemplateSlug(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "Raw Scaffold")
	if err := os.Mkdir(project, 0755); err != nil {
		t.Fatalf("Mkdir returned %v", err)
	}

	withWorkingDir(t, project, func() {
		if err := os.WriteFile("carbide.toml", []byte(`name = "__PROJECT_NAME__"
slug = "__PROJECT_SLUG__"
`), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}
		if err := os.WriteFile("docker-compose.yml", []byte("services: {}\n"), 0644); err != nil {
			t.Fatalf("WriteFile docker-compose.yml returned %v", err)
		}

		env := composeEnv(nil)
		if got := envSliceValue(env, "COMPOSE_PROJECT_NAME"); got != "raw-scaffold" {
			t.Fatalf("COMPOSE_PROJECT_NAME = %q, want raw-scaffold", got)
		}
	})
}

func TestCopyScaffoldSkipsRuntimeDirectory(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source")
	target := filepath.Join(t.TempDir(), "target")
	if err := os.MkdirAll(filepath.Join(source, ".carbide", "log"), 0755); err != nil {
		t.Fatalf("MkdirAll .carbide returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, ".carbide", "log", "dev.jsonl"), []byte("{}\n"), 0644); err != nil {
		t.Fatalf("WriteFile dev log returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "README.md"), []byte("__PROJECT_NAME__\n"), 0644); err != nil {
		t.Fatalf("WriteFile README returned %v", err)
	}

	if err := copyScaffoldPart(source, target, "Demo", "demo"); err != nil {
		t.Fatalf("copyScaffoldPart returned %v", err)
	}
	if isDir(filepath.Join(target, ".carbide")) {
		t.Fatalf("copyScaffoldPart copied runtime .carbide directory")
	}
	content, err := os.ReadFile(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatalf("ReadFile README returned %v", err)
	}
	if got := strings.TrimSpace(string(content)); got != "Demo" {
		t.Fatalf("README content = %q, want Demo", got)
	}
}

func TestBareCommandPrintsCommandList(t *testing.T) {
	var out bytes.Buffer
	a := app{stdout: &out}

	if err := a.run(nil); err != nil {
		t.Fatalf("run returned %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"_____________________________________________________",
		"________________________oo_______oo_______oo_________",
		"Carbide 0.1.0-dev",
		"Usage:",
		"carbide <command> [arguments]",
		"Commands:",
		"new <project-name>",
		"init",
		"help",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("bare command output = %q, missing %q", got, want)
		}
	}
	for _, unwanted := range []string{
		"Options:",
		"Available commands:",
		"run dev",
		"status",
		"stop dev",
		"follow logs",
		"upgrade",
		"version",
		"features:",
		"raw.githubusercontent.com/ryangerardwilson/carbide",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("bare command output = %q, should not contain %q", got, unwanted)
		}
	}
}

func TestHelpPrintsRuntimeReference(t *testing.T) {
	var out bytes.Buffer
	a := app{stdout: &out}

	if err := a.run([]string{"help"}); err != nil {
		t.Fatalf("run returned %v", err)
	}

	got := out.String()
	if !strings.HasPrefix(got, "Usage:\n") {
		t.Fatalf("help output = %q, should start with Usage section", got)
	}
	for _, line := range strings.Split(strings.TrimRight(got, "\n"), "\n") {
		if width := len(stripANSI(line)); width > 79 {
			t.Fatalf("help line is %d columns, want <= 79: %q", width, line)
		}
	}
	assertOutputOrder(t, got, []string{
		"Usage:",
		"  carbide <command> [arguments]",
		"Available commands:",
		"  deploy apply prod",
		"  deploy preview prod",
		"  doctor",
		"  doctor env",
		"  doctor framework",
		"  doctor runtime",
		"  help",
		"  init",
		"  logs",
		"  new <project-name>",
		"  status",
		"  upgrade",
		"  version",
		"follow\n",
		"  follow logs",
		"  follow logs service api",
		"logs\n",
		"  logs containing \"/api/login\" json",
		"run\n",
		"  run dev",
		"stop\n",
		"  stop dev",
	})
	for _, want := range []string{
		"Usage:",
		"Available commands:",
		"new <project-name>",
		"deploy preview prod",
		"deploy apply prod",
		"doctor",
		"doctor env",
		"doctor runtime",
		"doctor framework",
		"run dev",
		"stop dev",
		"follow logs",
		"logs containing \"/api/login\" json",
		"upgrade",
		"version",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("help output = %q, missing %q", got, want)
		}
	}
	for _, unwanted := range []string{
		"area",
		"command  ",
		"purpose",
		"carbide help",
		"carbide run dev",
		"Carbide\n",
		"Containerized full-stack apps with React, Go, and Postgres.",
		"_____________________________________________________",
		"________________________oo_______oo_______oo_________",
		"install the CLI",
		"<github-install-url>",
		"curl -fsSL",
		"raw.githubusercontent.com/ryangerardwilson/carbide",
		"features:",
		"global actions:",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("help output = %q, should not contain %q", got, unwanted)
		}
	}
}

func TestDoctorPrintsProjectContract(t *testing.T) {
	withGeneratedScaffold(t, func(dir string) {
		var out bytes.Buffer
		a := app{stdout: &out}
		if err := a.run([]string{"doctor"}); err != nil {
			t.Fatalf("doctor returned %v\n%s", err, out.String())
		}

		got := out.String()
		for _, want := range []string{
			"Carbide doctor",
			"project contract",
			"project shape     ok",
			"runtime baseline  ok",
			"env contract      ok      0 missing, 2 secrets",
			"compose           ok      web api db",
			"frontend          ok      Bun React Tailwind",
			"api               ok      Go HTTP API",
			"database          ok      Postgres users sessions",
			"agents            ok      AGENTS.md agents.d",
			"regressions       ok      no legacy markers",
			"runtime           skip    run carbide doctor runtime",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("doctor output = %q, missing %q", got, want)
			}
		}
		if strings.Contains(got, "postgres://") {
			t.Fatalf("doctor output printed secret value: %q", got)
		}
	})
}

func TestDoctorRejectsLegacyRootDirectory(t *testing.T) {
	withGeneratedScaffold(t, func(dir string) {
		if err := os.Mkdir("model", 0755); err != nil {
			t.Fatalf("Mkdir model returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		err := a.run([]string{"doctor"})
		if err == nil {
			t.Fatalf("doctor should reject legacy root directory")
		}
		if !strings.Contains(out.String(), "legacy root dirs: model") {
			t.Fatalf("doctor output = %q", out.String())
		}
	})
}

func TestDoctorEnvPrintsContractSummary(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "demo"

[env]
contract_version = 1

[env.variables.DATABASE_URL]
service = "api"
required = true
secret = true
browser_exposed = false
framework_owned = true
local_default = "postgres://carbide:carbide@db:5432/carbide"

[env.variables.PUBLIC_APP_NAME]
service = "web"
required = false
secret = false
browser_exposed = true
framework_owned = false
local_default = "Demo"
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		if err := a.run([]string{"doctor", "env"}); err != nil {
			t.Fatalf("run returned %v", err)
		}

		got := out.String()
		for _, want := range []string{
			"Carbide doctor",
			"environment contract",
			"contract   carbide.toml",
			"env        .env not found; local defaults active",
			"status     ok",
			"required   0 missing",
			"secrets    1 declared",
			"browser    1 exposed",
			"framework  1 owned",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("doctor env output = %q, missing %q", got, want)
			}
		}
		if strings.Contains(got, "postgres://") {
			t.Fatalf("doctor env output printed secret value: %q", got)
		}
	})
}

func TestDoctorEnvRejectsMissingRequiredValue(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "demo"

[env]
contract_version = 1

[env.variables.API_SECRET]
required = true
secret = true
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		err := a.run([]string{"doctor", "env"})
		if err == nil {
			t.Fatalf("doctor env should reject missing required values")
		}
		if !strings.Contains(err.Error(), "missing required value") {
			t.Fatalf("doctor env error = %v", err)
		}
		if !strings.Contains(out.String(), "missing  API_SECRET") {
			t.Fatalf("doctor env output = %q", out.String())
		}
	})
}

func TestDeployPreviewAndApplyAreGuarded(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "demo"

[env]
contract_version = 1

[env.variables.APP_ENV]
required = true
local_default = "development"
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		if err := a.run([]string{"deploy", "preview", "prod"}); err != nil {
			t.Fatalf("preview returned %v", err)
		}
		got := out.String()
		for _, want := range []string{
			"Carbide deploy",
			"preview prod",
			"target   prod",
			"mutates  no",
			"plan     validate env contract",
			"refuse apply until target is implemented",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("deploy preview output = %q, missing %q", got, want)
			}
		}

		out.Reset()
		err := a.run([]string{"deploy", "apply", "prod"})
		if err == nil {
			t.Fatalf("deploy apply should be disabled")
		}
		if !strings.Contains(err.Error(), "disabled until a deploy target exists") {
			t.Fatalf("deploy apply error = %v", err)
		}
		if !strings.Contains(out.String(), "status   disabled") {
			t.Fatalf("deploy apply output = %q", out.String())
		}
	})
}

func TestDeployPreviewReadsSSHComposeTarget(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "Carbide Docs"
slug = "carbide-docs"

[env]
contract_version = 1

[env.variables.APP_ENV]
required = false
local_default = "production"

[deploy.targets.de-sci]
type = "ssh-compose"
host = "de-sci"
domain = "carbide.ryangerardwilson.com"
remote_path = "/opt/carbide/docs"
source_path = ".."
compose_file = "app/docker-compose.yml"
project_directory = "app"
public_port = 18081
health_path = "/health"
nginx = true
nginx_site = "carbide"
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		if err := a.run([]string{"deploy", "preview", "de-sci"}); err != nil {
			t.Fatalf("preview returned %v", err)
		}
		got := out.String()
		for _, want := range []string{
			"Carbide deploy",
			"preview de-sci",
			"target   de-sci",
			"type     ssh-compose",
			"host     de-sci",
			"domain   carbide.ryangerardwilson.com",
			"remote   /opt/carbide/docs",
			"compose  app/docker-compose.yml",
			"port     18081",
			"mutates  no",
			"apply    carbide deploy apply de-sci",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("deploy preview output = %q, missing %q", got, want)
			}
		}
		if strings.Contains(got, "refuse apply until target is implemented") {
			t.Fatalf("deploy preview output still shows disabled stub: %q", got)
		}
	})
}

func TestDeployEnvContentUsesProjectMetadata(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "My Carbide App"
slug = "my-carbide-app"
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}

		got, err := deployEnvContent(deployTarget{
			Name:       "prod",
			Domain:     "app.example.com",
			PublicPort: 18080,
			Nginx:      true,
		})
		if err != nil {
			t.Fatalf("deployEnvContent returned %v", err)
		}
		for _, want := range []string{
			"COMPOSE_PROJECT_NAME=my-carbide-app",
			"PUBLIC_APP_NAME=My Carbide App",
			"PUBLIC_URL=https://app.example.com",
			"CARBIDE_HTTP_PORT=18080",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("deploy env = %q, missing %q", got, want)
			}
		}
		if strings.Contains(got, "carbide-docs") || strings.Contains(got, "Carbide Docs") {
			t.Fatalf("deploy env still contains docs app defaults: %q", got)
		}
	})
}

func TestDeployPreviewReadsSSHComposeEnvironmentTarget(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "Carbide Docs"
slug = "carbide-docs"

[env]
contract_version = 1

[env.variables.APP_ENV]
required = false
local_default = "production"

[deploy.hosts.web-1]
ssh = "web-1"

[deploy.hosts.api-1]
ssh = "api-1"

[deploy.hosts.db-1]
ssh = "db-1"

[deploy.targets.prod]
type = "ssh-compose-environment"
domain = "carbide.example.com"
remote_path = "/opt/carbide/app"
source_path = "."
compose_file = "docker-compose.yml"
project_directory = "."
health_path = "/health"
strategy = "preview-only"

[deploy.targets.prod.roles.web]
hosts = ["web-1"]
public_port = 8080
nginx = true

[deploy.targets.prod.roles.api]
hosts = ["api-1"]

[deploy.targets.prod.roles.db]
hosts = ["db-1"]
primary = "db-1"
migration = "once"
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		if err := a.run([]string{"deploy", "preview", "prod"}); err != nil {
			t.Fatalf("preview returned %v", err)
		}
		got := out.String()
		for _, want := range []string{
			"Carbide deploy",
			"preview prod",
			"target   prod",
			"type     ssh-compose-environment",
			"domain   carbide.example.com",
			"mutates  no",
			"hosts    api-1 -> api-1",
			"         db-1 -> db-1",
			"         web-1 -> web-1",
			"roles    api: api-1",
			"         db: db-1 primary db-1 migrate once",
			"         web: web-1 port 8080 nginx",
			"apply    disabled until clustered orchestration is implemented",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("deploy preview output = %q, missing %q", got, want)
			}
		}
	})
}

func TestDeployApplyGuardsSSHComposeEnvironmentTarget(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "Carbide Docs"
slug = "carbide-docs"

[env]
contract_version = 1

[env.variables.APP_ENV]
required = false
local_default = "production"

[deploy.hosts.de-sci]
ssh = "de-sci"

[deploy.targets.prod]
type = "ssh-compose-environment"
domain = "carbide.example.com"

[deploy.targets.prod.roles.web]
hosts = ["de-sci"]
public_port = 8080
nginx = true

[deploy.targets.prod.roles.api]
hosts = ["de-sci"]

[deploy.targets.prod.roles.db]
hosts = ["de-sci"]
primary = "de-sci"
migration = "once"
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		err := a.run([]string{"deploy", "apply", "prod"})
		if err == nil {
			t.Fatalf("deploy apply should guard environment targets")
		}
		if !strings.Contains(err.Error(), "clustered orchestration is implemented") {
			t.Fatalf("deploy apply error = %v", err)
		}
		got := out.String()
		for _, want := range []string{
			"target   prod",
			"type     ssh-compose-environment",
			"status   guarded",
			"reason   clustered apply needs explicit orchestration",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("deploy apply output = %q, missing %q", got, want)
			}
		}
	})
}

func TestDoctorPrintsDocsProjectContract(t *testing.T) {
	source, err := filepath.Abs(filepath.Join("..", "..", "..", "docs", "app"))
	if err != nil {
		t.Fatalf("Abs docs app returned %v", err)
	}
	target := filepath.Join(t.TempDir(), "docs-app")
	if err := copyScaffoldPart(source, target, "Carbide Docs", "carbide-docs"); err != nil {
		t.Fatalf("copyScaffoldPart returned %v", err)
	}

	withWorkingDir(t, target, func() {
		var out bytes.Buffer
		a := app{stdout: &out}
		if err := a.run([]string{"doctor"}); err != nil {
			t.Fatalf("doctor returned %v\n%s", err, out.String())
		}

		got := out.String()
		for _, want := range []string{
			"Carbide doctor",
			"project contract",
			"project shape     ok",
			"config            ok",
			"runtime baseline  ok",
			"env contract      ok",
			"compose           ok      docs web api db",
			"web               ok      Bun React Tailwind docs",
			"api               ok      docs health API",
			"database          ok      Postgres deploy checks",
			"agents            ok      docs agents.d",
			"runtime           skip    run carbide doctor runtime",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("doctor output = %q, missing %q", got, want)
			}
		}
	})
}

func withWorkingDir(t *testing.T, dir string, work func()) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir returned %v", err)
	}
	defer func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore Chdir returned %v", err)
		}
	}()
	work()
}

func withGeneratedScaffold(t *testing.T, work func(string)) {
	t.Helper()
	source, err := filepath.Abs(filepath.Join("..", "..", "..", "scaffold"))
	if err != nil {
		t.Fatalf("Abs scaffold returned %v", err)
	}
	if !isDir(source) {
		t.Fatalf("missing scaffold source at %s", source)
	}
	target := filepath.Join(t.TempDir(), "demo")
	if err := copyScaffoldPart(source, target, "Demo", "demo"); err != nil {
		t.Fatalf("copyScaffoldPart returned %v", err)
	}
	withWorkingDir(t, target, func() {
		work(target)
	})
}

func envSliceValue(env []string, key string) string {
	prefix := key + "="
	for _, item := range env {
		if strings.HasPrefix(item, prefix) {
			return strings.TrimPrefix(item, prefix)
		}
	}
	return ""
}

func assertOutputOrder(t *testing.T, output string, values []string) {
	t.Helper()
	offset := 0
	for _, value := range values {
		index := strings.Index(output[offset:], value)
		if index < 0 {
			t.Fatalf("output missing %q after byte %d:\n%s", value, offset, output)
		}
		offset += index + len(value)
	}
}

func TestRendererStyledLogoUsesGlyphColors(t *testing.T) {
	r := renderer{styled: true}

	got := r.formatLogoLine(0, "_o0Ox")
	want := "\033[2;38;5;245m_\033[0m\033[38;5;220mo0O\033[0mx"
	if got != want {
		t.Fatalf("styled logo line = %q, want %q", got, want)
	}
	if plain := stripANSI(got); plain != "_o0Ox" {
		t.Fatalf("styled logo line strips to %q, want %q", plain, "_o0Ox")
	}
}

func TestRendererPlainOutput(t *testing.T) {
	var out bytes.Buffer
	newRenderer(&out).Message(
		"Carbide",
		"project created",
		outputRow{"path", "/tmp/demo"},
		outputRow{"next", "cd demo"},
		outputRow{"", "carbide run dev"},
	)

	want := "Carbide\nproject created\n\npath  /tmp/demo\nnext  cd demo\n      carbide run dev\n"
	if out.String() != want {
		t.Fatalf("renderer output = %q, want %q", out.String(), want)
	}
}

func TestRendererIndentsMultilineValues(t *testing.T) {
	var out bytes.Buffer
	newRenderer(&out).Rows(outputRow{"error", "first line\nsecond line"})

	want := "error  first line\n       second line\n"
	if out.String() != want {
		t.Fatalf("renderer output = %q, want %q", out.String(), want)
	}
}

func TestRendererTable(t *testing.T) {
	var out bytes.Buffer
	newRenderer(&out).Table(
		[]string{"service", "container", "ports"},
		[]tableRow{
			{"web", "demo-web-1", "localhost:8082"},
			{"api", "demo-api-1", "-"},
		},
	)

	got := out.String()
	for _, want := range []string{
		"service  container   ports",
		"web      demo-web-1  localhost:8082",
		"api      demo-api-1  -",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("table output = %q, missing %q", got, want)
		}
	}
}

func TestServiceProgressFrame(t *testing.T) {
	tests := []struct {
		state string
		step  int
		want  string
	}{
		{"starting", 0, "[C o  o  o ]"},
		{"starting", 1, "[-c o  o  o]"},
		{"stopping", 0, "[ o  o  o D]"},
		{"stopping", 1, "[o  o  o d-]"},
		{"ready", 0, "[##########]"},
		{"stopped", 0, "[          ]"},
		{"failed", 0, "[!!!!!!!!!!]"},
	}
	for _, test := range tests {
		if got := serviceProgressFrame(10, test.step, test.state); got != test.want {
			t.Fatalf("serviceProgressFrame(10, %d, %q) = %q, want %q", test.step, test.state, got, test.want)
		}
	}
}

func TestRendererLogoPacmanLine(t *testing.T) {
	r := renderer{}
	tests := []struct {
		position int
		step     int
		want     string
	}{
		{-1, 0, " o  "},
		{1, 0, "_C o"},
		{1, 1, "_c o"},
		{4, 0, "_o0_"},
	}
	for _, test := range tests {
		got := r.formatLogoPacmanLine("_o0_", test.position, test.step)
		if got != test.want {
			t.Fatalf("formatLogoPacmanLine(position=%d, step=%d) = %q, want %q", test.position, test.step, got, test.want)
		}
	}
}

func TestRendererStyledLogoPacmanLine(t *testing.T) {
	r := renderer{styled: true}

	got := r.formatLogoPacmanLine("_o0_", 1, 0)
	for _, want := range []string{
		"\033[2;38;5;245m_\033[0m",
		"\033[1;38;5;226mC\033[0m",
		"\033[2;38;5;220mo\033[0m",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("styled pacman logo line = %q, missing %q", got, want)
		}
	}
	if plain := stripANSI(got); plain != "_C o" {
		t.Fatalf("styled pacman logo line strips to %q, want %q", plain, "_C o")
	}
}

func TestServiceProgressState(t *testing.T) {
	tests := []struct {
		status composeServiceStatus
		want   string
	}{
		{composeServiceStatus{state: "running"}, "ready"},
		{composeServiceStatus{state: "running", health: "healthy"}, "ready"},
		{composeServiceStatus{state: "running", health: "starting"}, "starting"},
		{composeServiceStatus{state: "exited"}, "failed"},
		{composeServiceStatus{}, "starting"},
	}
	for _, test := range tests {
		if got := serviceProgressState(test.status); got != test.want {
			t.Fatalf("serviceProgressState(%#v) = %q, want %q", test.status, got, test.want)
		}
	}
}

func TestServiceStopProgressState(t *testing.T) {
	tests := []struct {
		status composeServiceStatus
		want   string
	}{
		{composeServiceStatus{state: "running"}, "stopping"},
		{composeServiceStatus{state: "exited"}, "stopping"},
		{composeServiceStatus{state: "stopped"}, "stopped"},
		{composeServiceStatus{state: "failed"}, "failed"},
		{composeServiceStatus{}, "stopping"},
	}
	for _, test := range tests {
		if got := serviceStopProgressState(test.status); got != test.want {
			t.Fatalf("serviceStopProgressState(%#v) = %q, want %q", test.status, got, test.want)
		}
	}
}

func TestServiceProgressRunsWithoutColor(t *testing.T) {
	var out bytes.Buffer
	r := renderer{out: &out, interactive: true, termWidth: 52}

	err := r.RunServiceProgress(
		[]string{"api"},
		func() map[string]composeServiceStatus {
			return nil
		},
		func() error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("RunServiceProgress returned %v", err)
	}

	got := out.String()
	lines := visibleTerminalLines(got)
	starting := terminalLineContaining(t, lines, "starting")
	ready := terminalLineContaining(t, lines, "ready")
	wantFrameWidth := 52 - len("api") - progressStateColumnWidth - 5

	wantStarting := "api  " + serviceProgressFrame(wantFrameWidth, 0, "starting") + " " + padRight("starting", progressStateColumnWidth)
	if starting != wantStarting {
		t.Fatalf("starting progress line = %q, want %q", starting, wantStarting)
	}
	if len(starting) != 52 {
		t.Fatalf("starting progress line width = %d, want 52: %q", len(starting), starting)
	}

	wantReady := "api  " + serviceProgressFrame(wantFrameWidth, 0, "ready") + " " + padRight("ready", progressStateColumnWidth)
	if ready != wantReady {
		t.Fatalf("ready progress line = %q, want %q", ready, wantReady)
	}
	if len(ready) != 52 {
		t.Fatalf("ready progress line width = %d, want 52: %q", len(ready), ready)
	}
}

func TestServiceStopProgressRunsWithoutColor(t *testing.T) {
	var out bytes.Buffer
	r := renderer{out: &out, interactive: true, termWidth: 52}

	err := r.RunServiceStopProgress(
		[]string{"api"},
		func() map[string]composeServiceStatus {
			return nil
		},
		func() error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("RunServiceStopProgress returned %v", err)
	}

	got := out.String()
	lines := visibleTerminalLines(got)
	stopping := terminalLineContaining(t, lines, "stopping")
	stopped := terminalLineContaining(t, lines, "stopped")
	wantFrameWidth := 52 - len("api") - progressStateColumnWidth - 5

	wantStopping := "api  " + serviceProgressFrame(wantFrameWidth, 0, "stopping") + " " + padRight("stopping", progressStateColumnWidth)
	if stopping != wantStopping {
		t.Fatalf("stopping progress line = %q, want %q", stopping, wantStopping)
	}
	if len(stopping) != 52 {
		t.Fatalf("stopping progress line width = %d, want 52: %q", len(stopping), stopping)
	}

	wantStopped := "api  " + serviceProgressFrame(wantFrameWidth, 0, "stopped") + " " + padRight("stopped", progressStateColumnWidth)
	if stopped != wantStopped {
		t.Fatalf("stopped progress line = %q, want %q", stopped, wantStopped)
	}
	if len(stopped) != 52 {
		t.Fatalf("stopped progress line width = %d, want 52: %q", len(stopped), stopped)
	}
}

func TestServiceProgressFrameWidthUsesTerminalWidth(t *testing.T) {
	r := renderer{termWidth: 80}
	if got := r.serviceProgressFrameWidth(len("web")); got != 64 {
		t.Fatalf("serviceProgressFrameWidth = %d, want 64", got)
	}

	r = renderer{termWidth: 20}
	if got := r.serviceProgressFrameWidth(len("web")); got != minimumProgressFrameWidth {
		t.Fatalf("narrow serviceProgressFrameWidth = %d, want %d", got, minimumProgressFrameWidth)
	}
}

func visibleTerminalLines(output string) []string {
	plain := stripANSI(output)
	var lines []string
	for _, line := range strings.Split(plain, "\n") {
		if line == "" {
			continue
		}
		lines = append(lines, strings.TrimLeft(line, "\r"))
	}
	return lines
}

func terminalLineContaining(t *testing.T, lines []string, value string) string {
	t.Helper()
	for _, line := range lines {
		if strings.Contains(line, value) {
			return line
		}
	}
	t.Fatalf("terminal output missing %q in %#v", value, lines)
	return ""
}

func TestStreamWatchOutputFiltersNoise(t *testing.T) {
	var out bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	streamWatchOutput(strings.NewReader("Watch enabled\n\nrebuilt api\n"), newRenderer(&out), nil, "stdout", &wg)
	wg.Wait()

	got := out.String()
	if !strings.Contains(got, "watch      rebuilt api\n") {
		t.Fatalf("watch output = %q", got)
	}
	if len(strings.Fields(got)[0]) != len("15:04:05") {
		t.Fatalf("watch output missing timestamp: %q", got)
	}
}

func TestStreamLogOutputParsesComposeServices(t *testing.T) {
	var out bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	streamLogOutput(
		strings.NewReader("api-1  | GET /health\nweb-1 | listening\ndemo-db-1 | ready\n"),
		newRenderer(&out),
		nil,
		"stdout",
		&wg,
	)
	wg.Wait()

	got := out.String()
	for _, want := range []string{
		"api        GET /health",
		"web        listening",
		"db         ready",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("log output = %q, missing %q", got, want)
		}
	}
}

func TestStreamLogOutputWritesStructuredJSON(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "dev.jsonl")
	sink, err := openDevLogSink(logPath)
	if err != nil {
		t.Fatalf("openDevLogSink returned %v", err)
	}

	var out bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	streamLogOutput(strings.NewReader("api-1 | GET /health\n"), newRenderer(&out), sink, "stdout", &wg)
	wg.Wait()
	if err := sink.Close(); err != nil {
		t.Fatalf("Close returned %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile returned %v", err)
	}
	text := string(data)
	for _, want := range []string{
		`"source":"compose-log"`,
		`"stream":"stdout"`,
		`"service":"api"`,
		`"message":"GET /health"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("structured log %q missing %s", text, want)
		}
	}
}

func TestParseLogQuery(t *testing.T) {
	query, err := parseLogQuery([]string{"service", "api", "containing", "health", "limit", "5", "json"})
	if err != nil {
		t.Fatalf("parseLogQuery returned %v", err)
	}
	if query.service != "api" || query.contains != "health" || query.limit != 5 || !query.json {
		t.Fatalf("query = %#v", query)
	}

	if _, err = parseLogQuery([]string{"follow", "service", "web"}); err == nil {
		t.Fatalf("parseLogQuery should reject follow as a logs option")
	}
}

func TestParseComposeServiceStatuses(t *testing.T) {
	statuses, err := parseComposeServiceStatuses(`[
{"Service":"web","Name":"demo-web-1","State":"running","Health":"","Publishers":[{"URL":"0.0.0.0","TargetPort":8080,"PublishedPort":8082,"Protocol":"tcp"},{"URL":"::","TargetPort":8080,"PublishedPort":8082,"Protocol":"tcp"}]},
{"Service":"api","Name":"demo-api-1","State":"running","Health":"healthy","Publishers":[{"URL":"","TargetPort":8080,"PublishedPort":0,"Protocol":"tcp"}]},
{"Service":"db","Name":"demo-db-1","State":"created","Health":"starting","Publishers":[{"URL":"","TargetPort":5432,"PublishedPort":0,"Protocol":"tcp"}]}
]`)
	if err != nil {
		t.Fatalf("parseComposeServiceStatuses returned %v", err)
	}
	if serviceProgressState(statuses["web"]) != "ready" {
		t.Fatalf("web status = %#v", statuses["web"])
	}
	if serviceProgressState(statuses["api"]) != "ready" {
		t.Fatalf("api status = %#v", statuses["api"])
	}
	if serviceProgressState(statuses["db"]) != "starting" {
		t.Fatalf("db status = %#v", statuses["db"])
	}

	snapshots, err := parseComposeServiceSnapshots(`[
{"Service":"web","Name":"demo-web-1","State":"running","Health":"","Publishers":[{"URL":"0.0.0.0","TargetPort":8080,"PublishedPort":8082,"Protocol":"tcp"},{"URL":"::","TargetPort":8080,"PublishedPort":8082,"Protocol":"tcp"}]},
{"Service":"api","Name":"demo-api-1","State":"running","Health":"healthy","Publishers":[{"URL":"","TargetPort":8080,"PublishedPort":0,"Protocol":"tcp"}]}
]`)
	if err != nil {
		t.Fatalf("parseComposeServiceSnapshots returned %v", err)
	}
	if got := composePublishedPorts(snapshots["web"]); got != "localhost:8082" {
		t.Fatalf("web published ports = %q", got)
	}
	if got := composeInternalPorts(snapshots["api"]); got != "8080/tcp" {
		t.Fatalf("api internal ports = %q", got)
	}
	if got := composeServiceStatusText(snapshots["api"]); got != "running (healthy)" {
		t.Fatalf("api status text = %q", got)
	}

	statuses, err = parseComposeServiceStatuses(`{"Service":"api","State":"exited","Health":""}`)
	if err != nil {
		t.Fatalf("parseComposeServiceStatuses line mode returned %v", err)
	}
	if serviceProgressState(statuses["api"]) != "failed" {
		t.Fatalf("api status = %#v", statuses["api"])
	}
}

func TestFilterAndLimitLogEntries(t *testing.T) {
	entries := []structuredLogEntry{
		{Service: "web", Message: "listening"},
		{Service: "api", Message: "GET /health"},
		{Service: "api", Message: "POST /api/login"},
	}

	filtered := filterLogEntries(entries, logQuery{service: "api", contains: "api"})
	if len(filtered) != 1 || filtered[0].Message != "POST /api/login" {
		t.Fatalf("filtered = %#v", filtered)
	}

	limited := limitLogEntries(entries, 2)
	if len(limited) != 2 || limited[0].Message != "GET /health" {
		t.Fatalf("limited = %#v", limited)
	}
}
