package cli

import (
	"fmt"
	"regexp"
	"strings"
)

func healthDocsProjectShape() healthResult {
	requiredDirs := []string{"web", "api", "db"}
	if missing := missingDirs(requiredDirs); len(missing) > 0 {
		return healthFail("project shape", "missing "+strings.Join(missing, ", "))
	}

	requiredFiles := []string{projectConfigPath, composeFilePath}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("project shape", "missing "+strings.Join(missing, ", "))
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
	extras := rootDirsOutsideContract(allowed)
	if len(extras) > 0 {
		return healthFail("project shape", "non-service root dirs: "+strings.Join(extras, ", "))
	}
	return healthOK("project shape", "docs web api db")
}

func healthDocsConfigContract() healthResult {
	content := readFileString(projectConfigPath)
	required := []string{
		`name = "Carbide Docs"`,
		`slug = "carbide-docs"`,
		`profile = "docs"`,
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
		"[deploy.targets.prod]",
		`script = "./deploy/prod.sh"`,
	}
	if missing := missingNeedles(content, required); len(missing) > 0 {
		return healthFail("config", "missing "+strings.Join(missing, ", "))
	}
	return healthOK("config", "docs deploy script")
}

func healthDocsRuntimeBaselineContract() healthResult {
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
			"github.com/jackc/pgx/v5",
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
	if findings := unsupportedGoDirectiveFindings([]string{"api/go.mod", "db/go.mod"}); len(findings) > 0 {
		return healthFail("runtime baseline", "Go directive drift: "+strings.Join(findings, ", "))
	}
	if findings := packageVersionRangeFindingsFor("web/package.json", []string{"tailwindcss", "@tailwindcss/cli", "typescript", "@types/bun"}); len(findings) > 0 {
		return healthFail("runtime baseline", "package ranges: "+strings.Join(findings, ", "))
	}
	return healthOK("runtime baseline", "docs pinned images")
}

func healthDocsComposeContract() healthResult {
	content := readFileString(composeFilePath)
	services := composeServiceNamesFromFile(composeFilePath)
	for _, service := range []string{"web", "api", "db"} {
		if !containsString(services, service) {
			return healthFail("compose", "missing "+service+" service")
		}
	}
	required := []string{
		"context: ..",
		"dockerfile: app/web/Dockerfile",
		"API_URL: http://api:8080",
		"CARBIDE_HTTP_PORT",
		"service_healthy",
		"./db/migration:/docker-entrypoint-initdb.d:ro",
		"develop:",
		"watch:",
		"path: ./web/src",
		"path: ./web/site",
		"path: ./web/package.json",
		"path: ./web/bun.lock",
		"path: ./web/tsconfig.json",
		"path: ./web/Dockerfile",
		"path: ./api",
		"path: ./db/migration",
	}
	if missing := missingNeedles(content, required); len(missing) > 0 {
		return healthFail("compose", "missing "+strings.Join(missing, ", "))
	}
	return healthOK("compose", "docs web api db")
}

func healthDocsWebContract() healthResult {
	requiredFiles := []string{
		"web/Dockerfile",
		"web/package.json",
		"web/bun.lock",
		"web/tsconfig.json",
		"web/site/index.html",
		"web/site/assets/intro.js",
		"web/site/assets/styles.css",
		"web/site/for/agents/index.md",
		"web/src/build-styles.ts",
		"web/src/server.ts",
		"web/src/styles.css",
		"web/src/styles.d.ts",
		"web/src/lib/cx.ts",
		"web/src/lib/types.ts",
		"web/src/component/l1/index.ts",
		"web/src/component/l1/tokens.ts",
		"web/src/component/l2/DocsChrome.ts",
		"web/src/component/l2/index.ts",
		"web/src/component/l3/DocsSite.ts",
		"web/src/component/l3/index.ts",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("web", "missing "+strings.Join(missing, ", "))
	}
	requiredDirs := []string{
		"web/site",
		"web/site/assets",
		"web/site/for",
		"web/site/for/agents",
		"web/src/component/l1",
		"web/src/component/l2",
		"web/src/component/l3",
		"web/src/lib",
	}
	if missing := missingDirs(requiredDirs); len(missing) > 0 {
		return healthFail("web", "missing "+strings.Join(missing, ", "))
	}
	if anyPathWithExtension("web/src", ".jsx") || anyPathWithExtension("web/src", ".js") || anyPathWithExtension("web/src", ".tsx") {
		return healthFail("web", "docs web source must use TypeScript")
	}
	required := map[string][]string{
		"web/Dockerfile":                     {"COPY app/web/src ./src", "COPY app/web/site ./site", "bun run typecheck", "bun run assets:build", `CMD ["bun", "run", "start"]`},
		"web/package.json":                   {`"build"`, `"assets:build"`, `"docs:styles"`, `"typecheck": "tsc --noEmit"`, `"@tailwindcss/cli":`, `"tailwindcss":`, `"typescript": "6.0.3"`, `"@types/bun": "1.3.14"`, `bun build --target=bun --production --outdir=build/server src/server.ts`},
		"web/tsconfig.json":                  {`"strict": true`, `"types": ["bun-types"]`, `"include": ["src/**/*.ts"]`},
		"web/src/build-styles.ts":            {"tailwindcss", "./src/styles.css", `join(process.cwd(), "site", "assets", "styles.css")`},
		"web/src/styles.css":                 {`@import "tailwindcss";`, `@source "./component/**/*.ts";`, `@source "./lib/**/*.ts";`, `@source "./server.ts";`, `@custom-variant dark`},
		"web/site/index.html":                {"Carbide Documentation", `href="/for/agents"`},
		"web/site/for/agents/index.md":       {"# Carbide for Agents", "https://carbide.ryangerardwilson.com/for/agents"},
		"web/src/server.ts":                  {"join(import.meta.dir, \"..\", \"site\")", "proxy(request", `url.pathname === "/health"`, `url.pathname.startsWith("/api/")`, `requestPath === "/for/agents"`, `"/for/agents/index.md"`, `pathname.endsWith(".html")`, "docsResponseHeaders", "rewriteDocsHtml", "cacheBustHtml", "versionedAssetPath", "createHash", `?v=${hash}`, "cacheControlFor", `return "no-cache"`, `return "no-store"`},
		"web/src/component/l1/tokens.ts":     {"docsClassLayers", "scrollbar", `bg-amber-50`, `dark:text-neutral-50`, `[scrollbar-width:thin]`, "l1:", "l2:", "l3:"},
		"web/src/component/l2/DocsChrome.ts": {"docsScrollbarClass", "docsChromeClassLayers", "docsStaticClassMap", "rewriteDocsClasses", "docsStaticHeaders", `bg-yellow-400`, `text-yellow-300`, `dark:[&_p]:text-neutral-300`, `dark:[&_pre]:text-neutral-50`, `[&_pre]:[scrollbar-width:thin]`},
		"web/src/component/l3/DocsSite.ts":   {"docsSiteClassLayers", "docsWebContract", "rewriteDocsHtml", "docsResponseHeaders", `[scrollbar-width:thin]`},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return healthFail("web", path+" missing "+strings.Join(missing, ", "))
		}
	}
	if findings := docsTailwindInputFindings("web/src/styles.css"); len(findings) > 0 {
		return healthFail("web", "docs Tailwind input contract: "+strings.Join(findings, "; "))
	}
	return healthOK("web", "Bun Tailwind TypeScript docs")
}

const scaffoldTailwindInputHeader = `@import "tailwindcss";

@source "./component/**/*.tsx";
@source "./lib/**/*.ts";
@source "./main.tsx";
@source "./server.ts";
`

const scaffoldTailwindInputExpected = scaffoldTailwindInputHeader + `
@custom-variant dark (&:where([data-theme="dark"], [data-theme="dark"] *));
`

func scaffoldTailwindInputFindings(path string) []string {
	content := readFileString(path)
	var findings []string
	if strings.TrimSpace(content) == "" {
		return []string{"missing Tailwind input"}
	}
	if !strings.HasPrefix(content, scaffoldTailwindInputHeader) {
		findings = append(findings, "must start with Tailwind import/source directives")
	}
	if strings.TrimSpace(content) != strings.TrimSpace(scaffoldTailwindInputExpected) {
		findings = append(findings, "must contain only Tailwind import/source directives and dark variant")
	}
	if lines := fileLineCount(path); lines > 60 {
		findings = append(findings, fmt.Sprintf("too large: %d lines", lines))
	}
	for _, needle := range []string{
		`@custom-variant dark (&:where([data-theme="dark"], [data-theme="dark"] *));`,
	} {
		if !strings.Contains(content, needle) {
			findings = append(findings, "missing "+needle)
		}
	}
	for _, forbidden := range []string{
		"@apply",
		"@layer",
		"@keyframes",
		"@media",
		"@container",
		"@plugin",
		"@config",
		"@theme",
		"--carbide-",
		"<style",
		"::-webkit-scrollbar",
		"scrollbar-color:",
		"scrollbar-width:",
		"html {",
		"body {",
		"font-size:",
		"line-height:",
		"min-width:",
		"margin:",
		"padding:",
	} {
		if strings.Contains(content, forbidden) {
			findings = append(findings, "forbidden "+forbidden)
		}
	}
	if regexp.MustCompile(`(?m)^\s*\.[A-Za-z_-][A-Za-z0-9_-]*(?:[\s,{:.#]|$)`).MatchString(content) {
		findings = append(findings, "custom CSS class selectors belong in Tailwind component classes")
	}
	if regexp.MustCompile(`(?m)^\s*#[A-Za-z_-][A-Za-z0-9_-]*(?:[\s,{:.#]|$)`).MatchString(content) {
		findings = append(findings, "custom ID selectors belong in Tailwind component classes")
	}
	return findings
}

func docsTailwindInputFindings(path string) []string {
	content := readFileString(path)
	var findings []string
	if strings.TrimSpace(content) == "" {
		return []string{"missing Tailwind input"}
	}
	const docsTailwindInputHeader = `@import "tailwindcss";

@source "./component/**/*.ts";
@source "./lib/**/*.ts";
@source "./server.ts";
`
	const docsTailwindInputExpected = docsTailwindInputHeader + `
@custom-variant dark (&:where([data-theme="dark"], [data-theme="dark"] *));
`
	if !strings.HasPrefix(content, docsTailwindInputHeader) {
		findings = append(findings, "must start with Tailwind import/source directives")
	}
	if strings.TrimSpace(content) != strings.TrimSpace(docsTailwindInputExpected) {
		findings = append(findings, "must contain only Tailwind import/source directives and dark variant")
	}
	if lines := fileLineCount(path); lines > 60 {
		findings = append(findings, fmt.Sprintf("too large: %d lines", lines))
	}
	for _, needle := range []string{
		`@custom-variant dark (&:where([data-theme="dark"], [data-theme="dark"] *));`,
	} {
		if !strings.Contains(content, needle) {
			findings = append(findings, "missing "+needle)
		}
	}
	for _, forbidden := range []string{
		"@apply",
		"@layer",
		"@keyframes",
		"@media",
		"@container",
		"@plugin",
		"@config",
		"@theme",
		"--carbide-",
		"<style",
		"::-webkit-scrollbar",
		"scrollbar-color:",
		"scrollbar-width:",
		"html {",
		"body {",
		"font-size:",
		"line-height:",
		"min-width:",
		"margin:",
		"padding:",
	} {
		if strings.Contains(content, forbidden) {
			findings = append(findings, "forbidden "+forbidden)
		}
	}
	if regexp.MustCompile(`(?m)^\s*\.[A-Za-z_-][A-Za-z0-9_-]*(?:[\s,{:.#]|$)`).MatchString(content) {
		findings = append(findings, "custom CSS class selectors belong in Tailwind component classes")
	}
	if regexp.MustCompile(`(?m)^\s*#[A-Za-z_-][A-Za-z0-9_-]*(?:[\s,{:.#]|$)`).MatchString(content) {
		findings = append(findings, "custom ID selectors belong in Tailwind component classes")
	}
	if generated := docsGeneratedTailwindFindings("web/site/assets/styles.css"); len(generated) > 0 {
		findings = append(findings, generated...)
	}
	return findings
}

func docsGeneratedTailwindFindings(path string) []string {
	if !isFile(path) {
		return nil
	}
	content := readFileString(path)
	forbidden := regexp.MustCompile(`\.(docs-layout|docs-sidebar|docs-toc|docs-topbar|docs-content|docs-intro|skip-link|topbar-inner|brand-mark)(?:[\s,{:.#]|$)`)
	if forbidden.MatchString(content) {
		return []string{"generated docs CSS contains custom docs selectors"}
	}
	if strings.Contains(content, "html{font-size:14px}") ||
		strings.Contains(content, "body{min-width:320px") ||
		strings.Contains(content, "body{margin:0;min-width:320px") {
		return []string{"generated docs CSS contains global html/body defaults"}
	}
	return nil
}

func healthDocsAPIContract() healthResult {
	requiredFiles := []string{
		"api/Dockerfile",
		"api/go.mod",
		"api/go.sum",
		"api/main.go",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("api", "missing "+strings.Join(missing, ", "))
	}
	required := map[string][]string{
		"api/go.mod":  {"module carbidedocs/api", "github.com/jackc/pgx/v5"},
		"api/main.go": {"/health", "/api/version", "pgxpool", "database unavailable"},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return healthFail("api", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return healthOK("api", "docs health API")
}

func healthDocsDatabaseContract() healthResult {
	requiredFiles := []string{
		"db/go.mod",
		"db/migration/001_docs.sql",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("database", "missing "+strings.Join(missing, ", "))
	}
	if !fileContains("db/migration/001_docs.sql", "CREATE TABLE IF NOT EXISTS docs_checks") {
		return healthFail("database", "missing docs check migration")
	}
	return healthOK("database", "Postgres docs checks")
}

func healthDocsAgentsContract() healthResult {
	requiredFiles := []string{
		"../../README.md",
	}
	if missing := missingFiles(requiredFiles); len(missing) > 0 {
		return healthFail("agents", "missing "+strings.Join(missing, ", "))
	}
	required := map[string][]string{
		"../../README.md": {
			"/for/agents",
			"source of truth for framework agents",
			"There is no separate internal docs tree under `docs/engineering/`.",
			"## Docs Website",
			"docs/app/web/site/",
			"docs/app/",
			"The docs app does not carry its own `AGENTS.md` or `README.md`.",
			"black and yellow",
			"audits should preserve that",
			"docs/app/deploy/prod.sh",
			"CARBIDE_DOCS_DEPLOY_SSH",
			"CARBIDE_DOCS_POSTGRES_PASSWORD",
			"carbide deploy prod",
			"tests/smoke/docs_for_agents_http.sh",
		},
	}
	for path, needles := range required {
		if missing := missingNeedles(readFileString(path), needles); len(missing) > 0 {
			return healthFail("agents", path+" missing "+strings.Join(missing, ", "))
		}
	}
	return healthOK("agents", "root README docs ops guidance /for/agents")
}

func healthForbiddenRegressions(root string) healthResult {
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
		return healthFail("regressions", strings.Join(hits, ", "))
	}
	return healthOK("regressions", "no legacy markers")
}
