package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	if err := os.MkdirAll(filepath.Join(source, "web", "node_modules", "pkg"), 0755); err != nil {
		t.Fatalf("MkdirAll node_modules returned %v", err)
	}
	if err := os.MkdirAll(filepath.Join(source, "web", "public"), 0755); err != nil {
		t.Fatalf("MkdirAll public returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, ".carbide", "log", "dev.jsonl"), []byte("{}\n"), 0644); err != nil {
		t.Fatalf("WriteFile dev log returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "web", "node_modules", "pkg", "index.js"), []byte("module.exports = {}\n"), 0644); err != nil {
		t.Fatalf("WriteFile node module returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "web", "public", "index.html"), []byte("<html></html>\n"), 0644); err != nil {
		t.Fatalf("WriteFile public returned %v", err)
	}
	if err := os.MkdirAll(filepath.Join(source, "web", "src"), 0755); err != nil {
		t.Fatalf("MkdirAll src returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "web", "src", "tailwind.css"), []byte("generated\n"), 0644); err != nil {
		t.Fatalf("WriteFile tailwind returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, ".env"), []byte("secret=value\n"), 0644); err != nil {
		t.Fatalf("WriteFile env returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "notes.txt"), []byte("__PROJECT_NAME__\n"), 0644); err != nil {
		t.Fatalf("WriteFile notes returned %v", err)
	}

	if err := copyScaffoldPart(source, target, "Demo", "demo"); err != nil {
		t.Fatalf("copyScaffoldPart returned %v", err)
	}
	if isDir(filepath.Join(target, ".carbide")) {
		t.Fatalf("copyScaffoldPart copied runtime .carbide directory")
	}
	if isDir(filepath.Join(target, "web", "node_modules")) {
		t.Fatalf("copyScaffoldPart copied node_modules")
	}
	if isDir(filepath.Join(target, "web", "public")) {
		t.Fatalf("copyScaffoldPart copied public build output")
	}
	if isFile(filepath.Join(target, "web", "src", "tailwind.css")) {
		t.Fatalf("copyScaffoldPart copied generated tailwind output")
	}
	if isFile(filepath.Join(target, ".env")) {
		t.Fatalf("copyScaffoldPart copied local env")
	}
	content, err := os.ReadFile(filepath.Join(target, "notes.txt"))
	if err != nil {
		t.Fatalf("ReadFile notes returned %v", err)
	}
	if got := strings.TrimSpace(string(content)); got != "Demo" {
		t.Fatalf("notes content = %q, want Demo", got)
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
		"Carbide 0.2.0",
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
		"  audit",
		"  clean dev",
		"  deploy prod",
		"  fix",
		"  health",
		"  health json",
		"  health env",
		"  health env json",
		"  health framework",
		"  health framework json",
		"  health runtime",
		"  health runtime json",
		"  help",
		"  init",
		"  logs",
		"  new <project-name>",
		"  resolve",
		"  resolve fix",
		"  status",
		"  status json",
		"  upgrade",
		"  urls",
		"  urls json",
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
		"deploy prod",
		"clean dev",
		"fix",
		"health",
		"health json",
		"health env",
		"health env json",
		"health runtime",
		"health runtime json",
		"health framework",
		"health framework json",
		"audit",
		"resolve",
		"resolve fix",
		"run dev",
		"status json",
		"stop dev",
		"urls",
		"urls json",
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

func TestHealthPrintsAppLaws(t *testing.T) {
	withGeneratedScaffold(t, func(dir string) {
		var out bytes.Buffer
		a := app{stdout: &out}
		if err := a.run([]string{"health"}); err != nil {
			t.Fatalf("health returned %v\n%s", err, out.String())
		}

		got := out.String()
		patterns := []string{
			`(?m)^Carbide health$`,
			`(?m)^app laws$`,
			`(?m)^project shape\s+ok\s+web api db$`,
			`(?m)^config\s+ok\s+carbide\.toml$`,
			`(?m)^deploy targets\s+ok\s+0 checked-in scripts$`,
			`(?m)^env contract\s+ok\s+0 missing, 2 secrets$`,
			`(?m)^compose\s+ok\s+web api db$`,
			`(?m)^line limits\s+ok\s+all checked files <= 1000 lines$`,
			`(?m)^regressions\s+ok\s+no legacy markers$`,
			`(?m)^runtime\s+skip\s+run carbide health runtime$`,
		}
		for _, pattern := range patterns {
			if !regexp.MustCompile(pattern).MatchString(got) {
				t.Fatalf("health output = %q, missing /%s/", got, pattern)
			}
		}
		if strings.Contains(got, "postgres://") {
			t.Fatalf("health output printed secret value: %q", got)
		}
	})
}

func TestHealthRejectsLegacyRootDirectory(t *testing.T) {
	withGeneratedScaffold(t, func(dir string) {
		if err := os.Mkdir("model", 0755); err != nil {
			t.Fatalf("Mkdir model returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		err := a.run([]string{"health"})
		if err == nil {
			t.Fatalf("health should reject legacy root directory")
		}
		if !strings.Contains(out.String(), "legacy root dirs: model") {
			t.Fatalf("health output = %q", out.String())
		}
	})
}

func TestAuditCreatesReportWorkspace(t *testing.T) {
	t.Setenv("CARBIDE_AUDIT_AUTOMATION", "0")
	repoRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatalf("Abs repo root returned %v", err)
	}

	withGeneratedScaffold(t, func(dir string) {
		var out bytes.Buffer
		a := app{home: repoRoot, stdout: &out}
		if err := a.run([]string{"audit"}); err != nil {
			t.Fatalf("audit returned %v\n%s", err, out.String())
		}

		root := filepath.Join(dir, ".audit")
		reportDir := filepath.Join(root, "report")
		latest := filepath.Join(root, "starter-reference")
		for _, path := range []string{
			reportDir,
			filepath.Join(latest, "carbide.toml"),
			filepath.Join(latest, "web", "src", "styles.css"),
			filepath.Join(latest, "web", "src", "main.tsx"),
		} {
			if !isDir(path) && !isFile(path) {
				t.Fatalf("missing audit artifact %s", path)
			}
		}
		if isDir(filepath.Join(latest, "web", "node_modules")) ||
			isDir(filepath.Join(latest, "web", "public")) ||
			isFile(filepath.Join(latest, "web", "src", "tailwind.css")) {
			t.Fatalf("audit scaffold copied generated web output")
		}
		specs := auditSpecs()
		for _, spec := range specs {
			reportPath := filepath.Join(reportDir, spec.fileName)
			if !isFile(reportPath) {
				t.Fatalf("missing report %s", reportPath)
			}
			content, err := os.ReadFile(reportPath)
			if err != nil {
				t.Fatalf("ReadFile report returned %v", err)
			}
			got := string(content)
			for _, want := range []string{
				"status: pending",
				spec.ref,
				spec.title,
				"Run `carbide audit` in an interactive terminal with `codex` installed.",
			} {
				if !strings.Contains(got, want) {
					t.Fatalf("report %s = %q, missing %q", spec.fileName, got, want)
				}
			}
		}

		for _, want := range []string{
			"Carbide audit",
			"workspace prepared",
			".audit",
			"starter-reference",
			"pending files",
			"carbide resolve",
		} {
			if !strings.Contains(out.String(), want) {
				t.Fatalf("audit output = %q, missing %q", out.String(), want)
			}
		}
	})
}

func TestResolveCreatesPendingPlanWithoutCodex(t *testing.T) {
	t.Setenv("CARBIDE_AUDIT_AUTOMATION", "0")
	repoRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatalf("Abs repo root returned %v", err)
	}

	withGeneratedScaffold(t, func(dir string) {
		var out bytes.Buffer
		a := app{home: repoRoot, stdout: &out}
		if err := a.run([]string{"audit"}); err != nil {
			t.Fatalf("audit returned %v\n%s", err, out.String())
		}
		out.Reset()
		if err := a.run([]string{"resolve"}); err != nil {
			t.Fatalf("resolve returned %v\n%s", err, out.String())
		}

		planPath := filepath.Join(dir, ".audit", "plan.md")
		content, err := os.ReadFile(planPath)
		if err != nil {
			t.Fatalf("ReadFile plan returned %v", err)
		}
		got := string(content)
		for _, want := range []string{
			"status: pending",
			"Carbide Resolve Plan",
			"codex is required to resolve pending audit reports automatically",
			"Run `carbide resolve` in an interactive terminal with `codex` installed.",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("plan = %q, missing %q", got, want)
			}
		}
		for _, want := range []string{
			"Carbide resolve",
			"plan stub created",
			".audit/plan.md",
		} {
			if !strings.Contains(out.String(), want) {
				t.Fatalf("resolve output = %q, missing %q", out.String(), want)
			}
		}
	})
}

func TestAuditPromptsUseNewStages(t *testing.T) {
	spec := auditSpecs()[0]
	report := auditReportPrompt(".audit", spec)
	for _, want := range []string{
		"Audit one Carbide contract slice",
		spec.ref,
		spec.title,
		".audit/report/" + spec.fileName,
		"status: complete",
		"## Recommended Changes",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("audit report prompt = %q, missing %q", report, want)
		}
	}

	resolve := resolveCodexPrompt(".audit", []auditClarification{{question: "Keep custom auth copy?", answer: "Yes"}})
	for _, want := range []string{
		"Resolve the Carbide audit reports into one implementation plan.",
		".audit/report",
		".audit/starter-reference",
		"status: ready | needs-clarification",
		"## Implementation Steps",
		"Keep custom auth copy? => Yes",
	} {
		if !strings.Contains(resolve, want) {
			t.Fatalf("resolve prompt = %q, missing %q", resolve, want)
		}
	}

	fix := fixCodexPrompt(".audit")
	for _, want := range []string{
		"Implement the latest Carbide resolve plan.",
		".audit/plan.md",
		"run `carbide health` at the end",
		"run `carbide health runtime` when runtime behavior or containers changed",
		"files changed",
	} {
		if !strings.Contains(fix, want) {
			t.Fatalf("fix prompt = %q, missing %q", fix, want)
		}
	}
}

func TestHealthEnvPrintsContractSummary(t *testing.T) {
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
		if err := a.run([]string{"health", "env"}); err != nil {
			t.Fatalf("run returned %v", err)
		}

		got := out.String()
		for _, want := range []string{
			"Carbide health",
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
				t.Fatalf("health env output = %q, missing %q", got, want)
			}
		}
		if strings.Contains(got, "postgres://") {
			t.Fatalf("health env output printed secret value: %q", got)
		}
	})
}

func TestHealthEnvRejectsMissingRequiredValue(t *testing.T) {
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
		err := a.run([]string{"health", "env"})
		if err == nil {
			t.Fatalf("health env should reject missing required values")
		}
		if !strings.Contains(err.Error(), "missing required value") {
			t.Fatalf("health env error = %v", err)
		}
		if !strings.Contains(out.String(), "missing  API_SECRET") {
			t.Fatalf("health env output = %q", out.String())
		}
	})
}

func TestDeployRequiresCheckedInScript(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "Demo"
slug = "demo"

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
		err := a.run([]string{"deploy", "prod"})
		if err == nil {
			t.Fatalf("deploy should fail without a checked-in target")
		}
		if !strings.Contains(err.Error(), "no checked-in deploy target named prod") {
			t.Fatalf("deploy error = %v", err)
		}
		if out.Len() != 0 {
			t.Fatalf("deploy output should stay empty when target is missing: %q", out.String())
		}
	})
}

func TestUpgradeBinaryNeedsRebuild(t *testing.T) {
	original := commit
	t.Cleanup(func() {
		commit = original
	})

	commit = ""
	if !upgradeBinaryNeedsRebuild("abc123") {
		t.Fatalf("expected empty embedded commit to require rebuild")
	}

	commit = "abc123"
	if upgradeBinaryNeedsRebuild("abc123") {
		t.Fatalf("expected matching embedded commit to skip rebuild")
	}

	commit = "def456"
	if !upgradeBinaryNeedsRebuild("abc123") {
		t.Fatalf("expected mismatched embedded commit to require rebuild")
	}
}

func TestDeployRunsCheckedInScript(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "My Carbide App"
slug = "my-carbide-app"

[deploy.targets.prod]
script = "./deploy/prod.sh"
description = "Ship the current app."
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}
		if err := os.MkdirAll("deploy", 0755); err != nil {
			t.Fatalf("MkdirAll deploy returned %v", err)
		}
		script := `#!/usr/bin/env bash
set -euo pipefail
printf '%s\n%s\n%s\n%s\n' "$CARBIDE_DEPLOY_TARGET" "$CARBIDE_PROJECT_ROOT" "$CARBIDE_PROJECT_NAME" "$CARBIDE_PROJECT_SLUG" > deploy.out
printf 'deploy ok\n'
`
		if err := os.WriteFile(filepath.Join("deploy", "prod.sh"), []byte(script), 0644); err != nil {
			t.Fatalf("WriteFile deploy script returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		err := a.run([]string{"deploy", "prod"})
		if err != nil {
			t.Fatalf("deploy returned %v", err)
		}

		got := out.String()
		for _, want := range []string{
			"Carbide deploy",
			"prod",
			"script  ./deploy/prod.sh",
			"about  Ship the current app.",
			"deploy ok",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("deploy output = %q, missing %q", got, want)
			}
		}

		content, err := os.ReadFile("deploy.out")
		if err != nil {
			t.Fatalf("ReadFile deploy.out returned %v", err)
		}
		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		if len(lines) != 4 {
			t.Fatalf("deploy env lines = %q", lines)
		}
		if lines[0] != "prod" {
			t.Fatalf("deploy target env = %q, want prod", lines[0])
		}
		if lines[2] != "My Carbide App" {
			t.Fatalf("deploy project name env = %q", lines[2])
		}
		if lines[3] != "my-carbide-app" {
			t.Fatalf("deploy project slug env = %q", lines[3])
		}
	})
}

func TestDeployRejectsScriptOutsideProject(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		config := `name = "Demo"
slug = "demo"

[deploy.targets.prod]
script = "../prod.sh"
`
		if err := os.WriteFile("carbide.toml", []byte(config), 0644); err != nil {
			t.Fatalf("WriteFile carbide.toml returned %v", err)
		}

		var out bytes.Buffer
		a := app{stdout: &out}
		err := a.run([]string{"deploy", "prod"})
		if err == nil {
			t.Fatalf("deploy should reject scripts outside the project")
		}
		if !strings.Contains(err.Error(), "script must stay inside the project") {
			t.Fatalf("deploy error = %v", err)
		}
		if out.Len() != 0 {
			t.Fatalf("deploy output should stay empty for invalid script paths: %q", out.String())
		}
	})
}

func TestHealthPrintsDocsProjectContract(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatalf("Abs repo root returned %v", err)
	}
	source, err := filepath.Abs(filepath.Join("..", "..", "..", "docs", "app"))
	if err != nil {
		t.Fatalf("Abs docs app returned %v", err)
	}
	target := filepath.Join(t.TempDir(), "docs-app")
	if err := copyScaffoldPart(source, target, "Carbide Docs", "carbide-docs"); err != nil {
		t.Fatalf("copyScaffoldPart returned %v", err)
	}
	rootReadme, err := os.ReadFile(filepath.Join(repoRoot, "README.md"))
	if err != nil {
		t.Fatalf("ReadFile README.md returned %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "..", "..", "README.md"), rootReadme, 0644); err != nil {
		t.Fatalf("WriteFile README.md returned %v", err)
	}

	withWorkingDir(t, target, func() {
		var out bytes.Buffer
		a := app{stdout: &out}
		if err := a.run([]string{"health"}); err != nil {
			t.Fatalf("health returned %v\n%s", err, out.String())
		}

		got := out.String()
		for _, want := range []string{
			"Carbide health",
			"app laws",
			"project shape     ok",
			"config            ok",
			"runtime baseline  ok",
			"deploy targets    ok      1 checked-in script: prod",
			"env contract      ok",
			"compose           ok      docs web api db",
			"web               ok      Bun Tailwind TypeScript docs",
			"api               ok      docs health API",
			"database          ok      Postgres docs checks",
			"agents            ok      root README docs ops guidance /for/agents",
			"line limits       ok      all checked files <= 1000 lines",
			"runtime           skip    run carbide health runtime",
		} {
			if !strings.Contains(got, want) {
				t.Fatalf("health output = %q, missing %q", got, want)
			}
		}
	})
}

func TestHealthRejectsDocsTailwindInputDrift(t *testing.T) {
	source, err := filepath.Abs(filepath.Join("..", "..", "..", "docs", "app"))
	if err != nil {
		t.Fatalf("Abs docs app returned %v", err)
	}
	target := filepath.Join(t.TempDir(), "docs-app")
	if err := copyScaffoldPart(source, target, "Carbide Docs", "carbide-docs"); err != nil {
		t.Fatalf("copyScaffoldPart returned %v", err)
	}

	stylesPath := filepath.Join(target, "web", "src", "styles.css")
	styles, err := os.ReadFile(stylesPath)
	if err != nil {
		t.Fatalf("ReadFile returned %v", err)
	}
	styles = append(styles, []byte("\n.docs-layout { display: grid; }\n")...)
	if err := os.WriteFile(stylesPath, styles, 0644); err != nil {
		t.Fatalf("WriteFile returned %v", err)
	}

	withWorkingDir(t, target, func() {
		var out bytes.Buffer
		a := app{stdout: &out}
		err := a.run([]string{"health"})
		if err == nil {
			t.Fatalf("health should reject drifted docs Tailwind input")
		}
		if !strings.Contains(out.String(), "docs Tailwind input contract") ||
			!strings.Contains(out.String(), "custom CSS class selectors belong in Tailwind component classes") {
			t.Fatalf("health output = %q", out.String())
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
