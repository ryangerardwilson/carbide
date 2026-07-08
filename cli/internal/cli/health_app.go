package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func (a app) projectHealthResults() []healthResult {
	if !isFile(projectConfigPath) {
		return []healthResult{healthFail("project", "missing carbide.toml")}
	}

	if projectProfile() == "docs" {
		return []healthResult{
			healthDocsProjectShape(),
			healthDocsConfigContract(),
			healthDocsRuntimeBaselineContract(),
			healthDeployTargetsContract(),
			healthEnvContract(),
			healthDocsComposeContract(),
			healthDocsWebContract(),
			healthDocsAPIContract(),
			healthDocsDatabaseContract(),
			healthDocsAgentsContract(),
			healthLineLimits("."),
			healthForbiddenRegressions("."),
		}
	}

	return []healthResult{
		healthProjectShape(),
		healthConfigContract(),
		healthDeployTargetsContract(),
		healthEnvContract(),
		healthComposeContract(),
		healthLineLimits("."),
		healthForbiddenRegressions("."),
	}
}

func healthProjectShape() healthResult {
	requiredDirs := []string{"web", "api", "db"}
	if missing := missingDirs(requiredDirs); len(missing) > 0 {
		return healthFail("project shape", "missing "+strings.Join(missing, ", "))
	}

	requiredFiles := []string{projectConfigPath, composeFilePath}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("project shape", "missing "+strings.Join(missing, ", "))
	}

	forbidden := []string{"src", "model", "controller", "view", "views", "frontend", "templates", "include", "infra", "doc"}
	if found := existingDirs(forbidden); len(found) > 0 {
		return healthFail("project shape", "legacy root dirs: "+strings.Join(found, ", "))
	}
	if isFile("go.mod") || isFile("go.sum") || isFile("Dockerfile") {
		return healthFail("project shape", "root Go/Docker files are not allowed")
	}

	services := composeServiceNamesFromFile(composeFilePath)
	allowed := map[string]bool{}
	for _, service := range services {
		allowed[service] = true
	}
	allowed["deploy"] = true
	if len(services) == 0 {
		for _, service := range defaultComposeServices() {
			allowed[service] = true
		}
		allowed["deploy"] = true
	}
	extras := rootDirsOutsideContract(allowed)
	if len(extras) > 0 {
		return healthFail("project shape", "non-service root dirs: "+strings.Join(extras, ", "))
	}
	return healthOK("project shape", "web api db")
}

func healthConfigContract() healthResult {
	content, err := os.ReadFile(projectConfigPath)
	if err != nil {
		return healthFail("config", err.Error())
	}
	text := string(content)
	required := []string{
		"name = ",
		"slug = ",
		"[dev]",
		"default_port = 8080",
		`database = "postgres"`,
		"[env]",
		"contract_version = 1",
	}
	if missing := missingNeedles(text, required); len(missing) > 0 {
		return healthFail("config", "missing "+strings.Join(missing, ", "))
	}
	return healthOK("config", "carbide.toml")
}

func healthDeployTargetsContract() healthResult {
	targets, err := loadDeployTargets()
	if err != nil {
		return healthFail("deploy targets", err.Error())
	}
	if len(targets) == 0 {
		return healthOK("deploy targets", "0 checked-in scripts")
	}

	names := make([]string, 0, len(targets))
	for _, target := range targets {
		if err := validateDeployTarget(target); err != nil {
			return healthFail("deploy targets", err.Error())
		}
		names = append(names, target.Name)
	}
	sort.Strings(names)
	if len(names) == 1 {
		return healthOK("deploy targets", "1 checked-in script: "+names[0])
	}
	return healthOK("deploy targets", fmt.Sprintf("%d checked-in scripts: %s", len(names), strings.Join(names, ", ")))
}

func healthRuntimeBaselineContract() healthResult {
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
			fmt.Sprintf(`"tailwindcss": "%s"`, baselineTailwindVersion),
			fmt.Sprintf(`"@tailwindcss/cli": "%s"`, baselineTailwindVersion),
			`"typescript": "6.0.3"`,
			`"@types/bun": "1.3.14"`,
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
			return healthFail("runtime baseline", path+" missing "+strings.Join(missing, ", "))
		}
	}
	if findings := floatingDockerReferences([]string{"api/Dockerfile", "web/Dockerfile", composeFilePath}); len(findings) > 0 {
		return healthFail("runtime baseline", "floating Docker refs: "+strings.Join(findings, ", "))
	}
	if findings := packageVersionRangeFindings("web/package.json"); len(findings) > 0 {
		return healthFail("runtime baseline", "package ranges: "+strings.Join(findings, ", "))
	}
	if findings := unsupportedGoDirectiveFindings([]string{"api/go.mod", "db/go.mod"}); len(findings) > 0 {
		return healthFail("runtime baseline", "Go directive drift: "+strings.Join(findings, ", "))
	}
	return healthOK("runtime baseline", "Go 1.25 React 19.2 Tailwind 4.3 Bun 1.3 Postgres 17")
}

func healthEnvContract() healthResult {
	report, err := inspectEnvContract()
	if err != nil {
		return healthFail("env contract", err.Error())
	}
	detail := fmt.Sprintf("%d missing, %d secrets", len(report.missingRequired), report.secretCount)
	if len(report.missingRequired) > 0 {
		return healthFail("env contract", detail)
	}
	if len(report.warnings) > 0 {
		return healthFail("env contract", strings.Join(report.warnings, "; "))
	}
	return healthOK("env contract", detail)
}

func healthComposeContract() healthResult {
	content, err := os.ReadFile(composeFilePath)
	if err != nil {
		return healthFail("compose", "missing docker-compose.yml")
	}
	text := string(content)
	services := composeServiceNamesFromFile(composeFilePath)
	for _, service := range []string{"web", "api", "db"} {
		if !containsString(services, service) {
			return healthFail("compose", "missing "+service+" service")
		}
	}
	if containsString(services, "backend") || containsString(services, "database") {
		return healthFail("compose", "legacy service names present")
	}
	required := []string{
		"API_URL: http://api:8080",
		"@db:5432/carbide",
		`PUBLIC_URL: "http://localhost:${CARBIDE_HTTP_PORT:-8080}"`,
		"develop:",
		"watch:",
		"action: rebuild",
		"path: ./web/src",
		"path: ./web/tsconfig.json",
		"path: ./api",
		"path: ./db",
	}
	if missing := missingNeedles(text, required); len(missing) > 0 {
		return healthFail("compose", "missing "+strings.Join(missing, ", "))
	}
	return healthOK("compose", "web api db")
}

func healthLineLimits(root string) healthResult {
	violations, err := lawLineLimitViolations(root, maxLawFileLines)
	if err != nil {
		return healthFail("line limits", err.Error())
	}
	if len(violations) > 0 {
		return healthFail("line limits", strings.Join(violations, ", "))
	}
	return healthOK("line limits", fmt.Sprintf("all checked files <= %d lines", maxLawFileLines))
}

func healthFrontendContract() healthResult {
	requiredFiles := []string{
		"web/Dockerfile",
		"web/package.json",
		"web/bun.lock",
		"web/tsconfig.json",
		"web/index.html",
		"web/src/main.tsx",
		"web/src/server.ts",
		"web/src/write-index.ts",
		"web/src/styles.css",
		"web/src/styles.d.ts",
		"web/src/lib/cx.ts",
		"web/src/lib/types.ts",
		"web/src/component/l1/Button.tsx",
		"web/src/component/l1/Field.tsx",
		"web/src/component/l1/Surface.tsx",
		"web/src/component/l1/Text.tsx",
		"web/src/component/l1/ThemeToggle.tsx",
		"web/src/component/l1/tokens.ts",
		"web/src/component/l2/AuthForm.tsx",
		"web/src/component/l2/Layouts.tsx",
		"web/src/component/l3/AuthView.tsx",
		"web/src/component/l3/DashboardView.tsx",
		"web/src/component/l3/LoadingView.tsx",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("frontend", "missing "+strings.Join(missing, ", "))
	}
	requiredDirs := []string{"web/src/component/l1", "web/src/component/l2", "web/src/component/l3"}
	if missing := missingDirs(requiredDirs); len(missing) > 0 {
		return healthFail("frontend", "missing "+strings.Join(missing, ", "))
	}
	forbiddenFiles := []string{"web/package-lock.json", "web/vite.config.js", "web/src/component/l1/theme.css"}
	if found := existingFiles(forbiddenFiles); len(found) > 0 {
		return healthFail("frontend", "forbidden "+strings.Join(found, ", "))
	}
	if anyPathWithExtension("web/src", ".jsx") || anyPathWithExtension("web/src", ".js") || anyPathWithExtension("web/src", ".mjs") {
		return healthFail("frontend", "frontend source must use TypeScript")
	}
	if fileContains("web/src/styles.css", "theme.css") ||
		treeContains("web/src", "cb-") ||
		treeContains("web/src", "--cb-") {
		return healthFail("frontend", "parallel CSS theme detected")
	}
	if findings := scaffoldTailwindInputFindings("web/src/styles.css"); len(findings) > 0 {
		return healthFail("frontend", "scaffold Tailwind input contract: "+strings.Join(findings, "; "))
	}
	if treeContains("web/src", "carbide-") || treeContains("web/src", "--carbide-") {
		return healthFail("frontend", "generated carbide styling hooks detected")
	}
	if lines := fileLineCount("web/src/styles.css"); lines > 60 {
		return healthFail("frontend", fmt.Sprintf("Tailwind input too large: %d lines", lines))
	}
	if fileContains("web/src/styles.css", "#0f766e") ||
		fileContains("web/src/styles.css", "#115e59") ||
		fileContains("web/src/styles.css", "#2dd4bf") ||
		fileContains("web/src/styles.css", "#5eead4") ||
		fileContains("web/src/styles.css", "#16433c") ||
		fileContains("web/src/styles.css", "#0f302c") ||
		fileContains("web/src/component/l1/tokens.ts", "from-carbide-action via-carbide-hero-via") {
		return healthFail("frontend", "green scaffold palette detected")
	}
	if fileContains("web/src/component/l2/Layouts.tsx", "text-7xl") ||
		fileContains("web/src/component/l2/Layouts.tsx", "text-5xl") ||
		fileContains("web/src/component/l2/Layouts.tsx", "py-24") ||
		fileContains("web/src/component/l2/Layouts.tsx", "lg:py-12") ||
		fileContains("web/src/component/l2/Layouts.tsx", "lg:grid-cols-[280px") ||
		fileContains("web/src/component/l2/Layouts.tsx", "lg:grid-cols-[240px") ||
		fileContains("web/src/component/l3/DashboardView.tsx", "gap-6") ||
		fileContains("web/src/component/l3/DashboardView.tsx", "p-6") ||
		fileContains("web/src/component/l1/Field.tsx", "min-h-12 rounded-md border") ||
		fileContains("web/src/component/l1/Field.tsx", "min-h-10 rounded-md border") ||
		treeContains("web/src/component", "font-extrabold") {
		return healthFail("frontend", "oversized scaffold density detected")
	}
	if fileContains("web/src/component/l1/ThemeToggle.tsx", "aria-pressed") ||
		fileContains("web/src/component/l1/ThemeToggle.tsx", `role="group"`) ||
		fileContains("web/src/component/l1/ThemeToggle.tsx", `<select`) ||
		fileContains("web/src/component/l1/ThemeToggle.tsx", `appearance-none`) {
		return healthFail("frontend", "non-icon theme toggle detected")
	}
	if !fileContains("web/package.json", `"react":`) ||
		!fileContains("web/package.json", `"tailwindcss":`) ||
		!fileContains("web/package.json", `"@tailwindcss/cli":`) ||
		!fileContains("web/package.json", `"typescript": "6.0.3"`) ||
		!fileContains("web/package.json", `"@types/bun": "1.3.14"`) ||
		!fileContains("web/package.json", `"@types/react": "19.2.17"`) ||
		!fileContains("web/package.json", `"@types/react-dom": "19.2.3"`) ||
		!fileContains("web/package.json", `"typecheck": "tsc --noEmit"`) ||
		!fileContains("web/package.json", `"assets:build":`) ||
		!fileContains("web/package.json", `--entry-naming='assets/[name]-[hash].[ext]'`) ||
		!fileContains("web/tsconfig.json", `"strict": true`) ||
		!fileContains("web/tsconfig.json", `"jsx": "react-jsx"`) ||
		!fileContains("web/tsconfig.json", `"types": ["bun-types"]`) ||
		!fileContains("web/Dockerfile", `bun run typecheck`) ||
		!fileContains("web/Dockerfile", `bun run assets:build`) ||
		!fileContains("web/src/server.ts", `publicRoot`) ||
		!fileContains("web/src/server.ts", `Cache-Control`) ||
		!fileContains("web/src/server.ts", `public, max-age=31536000, immutable`) ||
		!fileContains("web/src/server.ts", `return 'no-store'`) ||
		!fileContains("web/src/write-index.ts", `asset-manifest.json`) ||
		!fileContains("web/src/write-index.ts", `/assets/${scripts[0]}`) ||
		!fileContains("web/src/styles.css", `@import "tailwindcss";`) ||
		!fileContains("web/src/styles.css", `@source "./component/**/*.tsx";`) ||
		!fileContains("web/src/styles.css", `@source "./lib/**/*.ts";`) ||
		!fileContains("web/src/styles.css", `@source "./main.tsx";`) ||
		!fileContains("web/src/styles.css", `@source "./server.ts";`) ||
		!fileContains("web/src/styles.css", `@custom-variant dark`) ||
		!fileContains("web/index.html", `[scrollbar-width:thin]`) ||
		!fileContains("web/index.html", `dark:[scrollbar-color:rgb(82_82_82)_transparent]`) ||
		!fileContains("web/index.html", `prefers-color-scheme: dark`) ||
		!fileContains("web/src/main.tsx", `carbide.theme`) ||
		!fileContains("web/src/component/l1/ThemeToggle.tsx", `SunIcon`) ||
		!fileContains("web/src/component/l1/ThemeToggle.tsx", `MoonIcon`) ||
		!fileContains("web/src/component/l1/ThemeToggle.tsx", `Switch to light theme`) ||
		!fileContains("web/src/component/l1/ThemeToggle.tsx", `Switch to dark theme`) ||
		!fileContains("web/src/component/l1/ThemeToggle.tsx", `size-8 rounded-full border`) ||
		!fileContains("web/src/component/l1/ThemeToggle.tsx", `data-theme-mode`) ||
		!fileContains("web/src/component/l1/tokens.ts", `bg-white text-neutral-950 dark:bg-black dark:text-neutral-50`) ||
		!fileContains("web/src/component/l1/tokens.ts", `const scrollbar =`) ||
		!fileContains("web/src/component/l1/tokens.ts", `[scrollbar-width:thin]`) ||
		!fileContains("web/src/component/l1/tokens.ts", `dark:[scrollbar-color:rgb(82_82_82)_transparent]`) ||
		!fileContains("web/src/component/l1/Text.tsx", `text-2xl/8 sm:text-3xl/9`) ||
		!fileContains("web/src/component/l1/Field.tsx", `min-h-8 rounded-md border px-2 py-1 text-sm/6`) ||
		!fileContains("web/src/component/l1/Button.tsx", `md: 'min-h-8 px-3 text-xs'`) ||
		!fileContains("web/src/component/l2/AuthForm.tsx", `gap-3 border-l px-4 py-5`) ||
		!fileContains("web/src/component/l2/AuthForm.tsx", `w-full max-w-sm justify-self-center gap-3`) ||
		!fileContains("web/src/component/l2/Layouts.tsx", `lg:grid-cols-[216px_minmax(0,1fr)]`) ||
		!fileContains("web/src/component/l2/Layouts.tsx", `px-3 py-4 sm:px-5 lg:py-5`) ||
		!fileContains("web/src/component/l2/Layouts.tsx", `ui.scrollbar`) ||
		!fileContains("web/src/main.tsx", "./component/l3") {
		return healthFail("frontend", "React/Bun/Tailwind contract drifted")
	}
	return healthOK("frontend", "Bun React Tailwind TypeScript")
}

func healthAPIContract() healthResult {
	requiredFiles := []string{
		"api/Dockerfile",
		"api/go.mod",
		"api/go.sum",
		"api/main.go",
		"api/auth.go",
		"api/routes.go",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("api", "missing "+strings.Join(missing, ", "))
	}
	if fileContains("api/Dockerfile", "gcc") || fileContains("api/Dockerfile", "libpq-dev") || anyPathWithExtension("api", ".c", ".h") {
		return healthFail("api", "legacy C backend artifacts present")
	}
	required := map[string][]string{
		"api/go.mod":     {"module carbideapp/api", "carbideapp/db", "replace carbideapp/db => ../db"},
		"api/Dockerfile": {"FROM golang:", "go mod download", "COPY api ./api", "COPY db ./db"},
		"api/routes.go":  {"/api/register", "/api/login", "/api/me", "handleDashboard"},
		"api/main.go":    {"api listening on container port", "public API URL is"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return healthFail("api", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return healthOK("api", "Go HTTP API")
}

func healthDatabaseContract() healthResult {
	requiredFiles := []string{
		"db/go.mod",
		"db/go.sum",
		"db/user.go",
		"db/session.go",
		"db/migration/001_auth.sql",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("database", "missing "+strings.Join(missing, ", "))
	}
	required := map[string][]string{
		"db/go.mod":                 {"module carbideapp/db", "github.com/jackc/pgx/v5"},
		"db/user.go":                {"CreateUser", "VerifyUser", "pgxpool"},
		"db/session.go":             {"CreateSession", "CurrentUser", "DestroySession"},
		"db/migration/001_auth.sql": {"CREATE TABLE IF NOT EXISTS users", "CREATE TABLE IF NOT EXISTS sessions"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return healthFail("database", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return healthOK("database", "Postgres users sessions")
}
