#!/usr/bin/env bash
set -euo pipefail

domain="carbide.ryangerardwilson.com"

required_files=(
  ".gitignore"
  "README.md"
  "cli/install.sh"
  "cli/go.mod"
  "cli/bin/carbide"
  "cli/cmd/carbide/main.go"
  "cli/internal/cli/cli.go"
  "cli/internal/cli/cli_test.go"
  ".github/workflows/ci.yml"
  ".github/workflows/dependency-audit.yml"
  "docs/engineering/CI_CD_REGRESSION_TESTS.md"
  "docs/engineering/DEPLOYMENT.md"
  "docs/engineering/FRONTEND_STARTER_CONTRACT.md"
  "docs/engineering/DIRECTORY_STRUCTURE.md"
  "docs/engineering/CREATE_YOUR_FIRST_APP.md"
  "docs/engineering/VERSION_POLICY.md"
  "docs/app/carbide.toml"
  "docs/app/docker-compose.yml"
  "docs/app/agents.d/TAILWIND_COMPONENTS.md"
  "docs/app/web/bun.lock"
  "docs/app/web/src/build-styles.js"
  "docs/app/web/src/server.jsx"
  "docs/app/web/src/styles.css"
  "docs/app/web/src/lib/cx.js"
  "docs/app/web/src/component/l1/Text.jsx"
  "docs/app/web/src/component/l1/Surface.jsx"
  "docs/app/web/src/component/l1/index.js"
  "docs/app/web/src/component/l1/tokens.js"
  "docs/app/web/src/component/l2/DocsChrome.jsx"
  "docs/app/web/src/component/l2/index.js"
  "docs/app/web/src/component/l3/DocsSite.jsx"
  "docs/app/web/src/component/l3/index.js"
  "docs/site/index.html"
  "docs/site/deployment.html"
  "docs/site/frontend-starter-contract.html"
  "docs/site/create-your-first-app.html"
  "docs/site/ci-cd-regression-tests.html"
  "docs/site/repo-structure.html"
  "docs/site/version-policy.html"
  "docs/site/assets/styles.css"
  "tests/contract/audit_versions.sh"
  "tests/contract/check_repo_contract.sh"
  "tests/scaffold/cli_scaffold.sh"
  "tests/smoke/starter_docker_flow.sh"
  "scaffold/api/Dockerfile"
  "scaffold/AGENTS.md"
  "scaffold/README.md"
  "scaffold/.env.example"
  "scaffold/.gitignore"
  "scaffold/agents.d/BACKUP_RESTORE.md"
  "scaffold/agents.d/DEPLOY.md"
  "scaffold/agents.d/ENVIRONMENT.md"
  "scaffold/agents.d/TAILWIND_COMPONENTS.md"
  "scaffold/docker-compose.yml"
  "scaffold/web/Dockerfile"
  "scaffold/web/index.html"
  "scaffold/web/package.json"
  "scaffold/web/bun.lock"
  "scaffold/web/src/main.jsx"
  "scaffold/web/src/server.jsx"
  "scaffold/web/src/write-index.mjs"
  "scaffold/web/src/styles.css"
  "scaffold/web/src/lib/cx.js"
  "scaffold/web/src/component/l1/Button.jsx"
  "scaffold/web/src/component/l1/Field.jsx"
  "scaffold/web/src/component/l1/Surface.jsx"
  "scaffold/web/src/component/l1/Text.jsx"
  "scaffold/web/src/component/l1/ThemeToggle.jsx"
  "scaffold/web/src/component/l1/index.js"
  "scaffold/web/src/component/l1/tokens.js"
  "scaffold/web/src/component/l2/AuthForm.jsx"
  "scaffold/web/src/component/l2/Layouts.jsx"
  "scaffold/web/src/component/l2/index.js"
  "scaffold/web/src/component/l3/AuthView.jsx"
  "scaffold/web/src/component/l3/DashboardView.jsx"
  "scaffold/web/src/component/l3/LoadingView.jsx"
  "scaffold/web/src/component/l3/index.js"
  "scaffold/carbide.toml"
  "scaffold/api/go.mod"
  "scaffold/api/go.sum"
  "scaffold/db/go.mod"
  "scaffold/db/go.sum"
  "scaffold/api/main.go"
  "scaffold/api/auth.go"
  "scaffold/api/routes.go"
  "scaffold/db/user.go"
  "scaffold/db/session.go"
  "scaffold/db/migration/001_auth.sql"
)

required_dirs=(
  "cli"
  "cli/bin"
  "cli/cmd"
  "cli/cmd/carbide"
  "cli/internal"
  "cli/internal/cli"
  "tests/contract"
  "tests/scaffold"
  "tests/smoke"
  "scaffold"
  "scaffold/agents.d"
  "scaffold/web"
  "scaffold/web/src"
  "scaffold/web/src/component"
  "scaffold/web/src/lib"
  "scaffold/web/src/component/l1"
  "scaffold/web/src/component/l2"
  "scaffold/web/src/component/l3"
  "scaffold/api"
  "scaffold/db"
  "scaffold/db/migration"
)

for path in "${required_files[@]}"; do
  test -f "$path" || {
    printf 'missing required file: %s\n' "$path" >&2
    exit 1
  }
done

for path in "${required_dirs[@]}"; do
  test -d "$path" || {
    printf 'missing required directory: %s\n' "$path" >&2
    exit 1
  }
done

grep -Eq "Bun/React/Tailwind .*web.* container" README.md
grep -q "Postgres-only" README.md
grep -q "Separate runtime boundaries" README.md
grep -q "Infrastructure as code" README.md
grep -q "generated Docker Compose setup" README.md
grep -q "Postgres-backed queues" README.md
grep -q "carbide new" README.md
grep -F -q 'carbide new "My Carbide App"' README.md
grep -q 'name = "My Carbide App"' README.md
grep -q 'slug = "my-carbide-app"' README.md
grep -q "carbide run dev" README.md
grep -q "carbide status" README.md
grep -q "carbide stop dev" README.md
grep -q "carbide follow logs" README.md
grep -q "carbide logs" README.md
grep -q "carbide doctor" README.md
grep -q "carbide doctor env" README.md
grep -q "carbide doctor runtime" README.md
grep -q "carbide doctor framework" README.md
grep -q "carbide deploy preview" README.md
grep -q "carbide deploy apply" README.md
grep -q "VERSION_POLICY.md" README.md
grep -q "explicit Carbide baselines" README.md
grep -q "floating Docker images" README.md
! grep -q "command_format" cli/bin/carbide
! grep -q "carbide format" cli/bin/carbide
grep -q "module github.com/ryangerardwilson/carbide/cli" cli/go.mod
grep -q ".cli/" .gitignore
! test -f go.mod
! test -d src
! test -d examples
! test -d infra
test "$(find . -mindepth 1 -maxdepth 1 -type d ! -name '.git' ! -name '.github' ! -name '.cli' ! -name '.bin' -printf '%f\n' | sort | tr '\n' ' ')" = "cli docs scaffold tests "
! test -d include
! test -d templates
! test -d scaffold/config
! test -d scaffold/view
! test -d scaffold/src
! test -d scaffold/model
! test -d scaffold/controller
! test -d scaffold/migrations
! test -d scaffold/infra
! test -d scaffold/frontend
! test -d scaffold/doc
test "$(find scaffold -mindepth 1 -maxdepth 1 -type d ! -name .carbide -printf '%f\n' | sort | tr '\n' ' ')" = "agents.d api db web "
! test -f scaffold/Dockerfile
! test -f scaffold/go.mod
! test -f scaffold/go.sum
grep -q "oo_______oo_______oo" cli/internal/cli/cli.go
grep -q "package main" cli/cmd/carbide/main.go
grep -q "package cli" cli/internal/cli/cli.go
grep -q "commandDoctor()" cli/internal/cli/cli.go
grep -q "commandDoctorEnv" cli/internal/cli/cli.go
grep -q "commandDoctorRuntime" cli/internal/cli/cli.go
grep -q "commandDoctorFramework" cli/internal/cli/cli.go
grep -q "projectDoctorResults" cli/internal/cli/cli.go
grep -q "commandDeployPreview" cli/internal/cli/cli.go
grep -q "commandDeployApply" cli/internal/cli/cli.go
grep -q "projectDisplayName" cli/internal/cli/cli.go
! git grep -n -e 'S[e]alion' -e 's[e]alion' -e 'S[E]ALION' -- .
grep -q "composeUpDetached" cli/internal/cli/cli.go
grep -q "runDevStreams" cli/internal/cli/cli.go
grep -q -- "--quiet-build" cli/internal/cli/cli.go
grep -q "Carbide dev" cli/internal/cli/cli.go
grep -q "Go is required to build the Carbide CLI" cli/install.sh
grep -q ".cli/bin/carbide" cli/install.sh
grep -q "default_port = 8080" scaffold/carbide.toml
grep -q "contract_version = 1" scaffold/carbide.toml
grep -q "\\[runtime\\]" scaffold/carbide.toml
grep -q 'policy = "explicit-baseline"' scaffold/carbide.toml
grep -q 'go_module = "1.25.0"' scaffold/carbide.toml
grep -q 'go_builder_image = "golang:1.26-bookworm@sha256:b305420a68d0f229d91eb3b3ed9e519fcf2cf5461da4bef997bf927e8c0bfd2b"' scaffold/carbide.toml
grep -q 'api_runtime_image = "debian:trixie-slim@sha256:28de0877c2189802884ccd20f15ee41c203573bd87bb6b883f5f46362d24c5c2"' scaffold/carbide.toml
grep -q 'bun_image = "oven/bun:1.3.14-debian@sha256:9dba1a1b43ce28c9d7931bfc4eb00feb63b0114720a0277a8f939ae4dfc9db6f"' scaffold/carbide.toml
grep -q 'postgres_image = "postgres:17-alpine@sha256:dc17045ccfd343b49600570ea734b9c4991cf1c3f3302e67df51e3b402dd55c4"' scaffold/carbide.toml
grep -q 'react = "19.2.7"' scaffold/carbide.toml
grep -q 'tailwindcss = "4.3.2"' scaffold/carbide.toml
grep -q "\\[env.variables.DATABASE_URL\\]" scaffold/carbide.toml
grep -q "secret = true" scaffold/carbide.toml
grep -q "browser_exposed = true" scaffold/carbide.toml
grep -q "framework_owned = true" scaffold/carbide.toml
grep -q "preview_before_apply = true" scaffold/carbide.toml
grep -q ".carbide/" scaffold/.gitignore
grep -q ".env" scaffold/.gitignore
grep -q "web/node_modules/" scaffold/.gitignore
grep -q "web/public/" scaffold/.gitignore
grep -q "web/src/tailwind.css" scaffold/.gitignore
grep -q "__PROJECT_NAME__ Agent Context" scaffold/AGENTS.md
grep -q 'agents.d/ENVIRONMENT.md' scaffold/AGENTS.md
grep -q 'agents.d/DEPLOY.md' scaffold/AGENTS.md
grep -q 'agents.d/BACKUP_RESTORE.md' scaffold/AGENTS.md
grep -q 'agents.d/TAILWIND_COMPONENTS.md' scaffold/AGENTS.md
grep -q "carbide.toml" scaffold/README.md
grep -q "carbide doctor" scaffold/README.md
grep -q "carbide doctor env" scaffold/README.md
grep -q "carbide doctor runtime" scaffold/README.md
grep -q "explicit runtime baseline" scaffold/README.md
grep -q "Postgres major-version baseline change" scaffold/README.md
grep -q "carbide deploy preview prod" scaffold/README.md
grep -q "carbide deploy apply prod" scaffold/README.md
grep -q "POSTGRES_PASSWORD" scaffold/.env.example
grep -q "separate secrets container" scaffold/agents.d/ENVIRONMENT.md
grep -q "preview-before-apply" scaffold/agents.d/DEPLOY.md
grep -q "Postgres owns durable application state" scaffold/agents.d/BACKUP_RESTORE.md
grep -q "Tailwind Component Organization" scaffold/agents.d/TAILWIND_COMPONENTS.md
grep -q "Use L1/L2/L3 in two related ways" scaffold/agents.d/TAILWIND_COMPONENTS.md
grep -q "component/l1/" scaffold/agents.d/TAILWIND_COMPONENTS.md
grep -q "component/l2/" scaffold/agents.d/TAILWIND_COMPONENTS.md
grep -q "component/l3/" scaffold/agents.d/TAILWIND_COMPONENTS.md
grep -q "web/src/lib/cx.js" scaffold/agents.d/TAILWIND_COMPONENTS.md
grep -q "visible focus styling" scaffold/agents.d/TAILWIND_COMPONENTS.md
grep -q "Do not add a parallel .*theme.css" scaffold/agents.d/TAILWIND_COMPONENTS.md
! grep -q 'url = "http://localhost:8080"' scaffold/carbide.toml
grep -q "web:" scaffold/docker-compose.yml
grep -q "api:" scaffold/docker-compose.yml
grep -q "db:" scaffold/docker-compose.yml
! grep -q "backend:" scaffold/docker-compose.yml
! grep -q "database:" scaffold/docker-compose.yml
grep -q "postgres:17-alpine@sha256:dc17045ccfd343b49600570ea734b9c4991cf1c3f3302e67df51e3b402dd55c4" scaffold/docker-compose.yml
! grep -q "postgres:16-alpine" scaffold/docker-compose.yml
grep -q "API_URL: http://api:8080" scaffold/docker-compose.yml
grep -q "@db:5432/carbide" scaffold/docker-compose.yml
grep -q 'PUBLIC_URL: "http://localhost:${CARBIDE_HTTP_PORT:-8080}"' scaffold/docker-compose.yml
test "$(grep -c 'PUBLIC_URL: "http://localhost:${CARBIDE_HTTP_PORT:-8080}"' scaffold/docker-compose.yml)" -eq 2
grep -q 'PUBLIC_APP_NAME: "${PUBLIC_APP_NAME:-__PROJECT_NAME__}"' scaffold/docker-compose.yml
grep -q 'APP_ENV: "${APP_ENV:-development}"' scaffold/docker-compose.yml
grep -q 'POSTGRES_PASSWORD: "${POSTGRES_PASSWORD:-carbide}"' scaffold/docker-compose.yml
grep -q "develop:" scaffold/docker-compose.yml
grep -q "watch:" scaffold/docker-compose.yml
grep -q "action: rebuild" scaffold/docker-compose.yml
grep -q "context: ./web" scaffold/docker-compose.yml
grep -q "context: ." scaffold/docker-compose.yml
grep -q "dockerfile: api/Dockerfile" scaffold/docker-compose.yml
grep -q "path: ./web/src" scaffold/docker-compose.yml
grep -q "path: ./web/package.json" scaffold/docker-compose.yml
grep -q "path: ./web/bun.lock" scaffold/docker-compose.yml
grep -q "path: ./api" scaffold/docker-compose.yml
grep -q "path: ./db" scaffold/docker-compose.yml
grep -q "path: ./api/Dockerfile" scaffold/docker-compose.yml
! grep -q "path: ./go.mod" scaffold/docker-compose.yml
! grep -q "path: ./go.sum" scaffold/docker-compose.yml
grep -q "FROM golang:1.26-bookworm@sha256:b305420a68d0f229d91eb3b3ed9e519fcf2cf5461da4bef997bf927e8c0bfd2b" scaffold/api/Dockerfile
grep -q "FROM debian:trixie-slim@sha256:28de0877c2189802884ccd20f15ee41c203573bd87bb6b883f5f46362d24c5c2" scaffold/api/Dockerfile
! grep -q "debian:bookworm-slim" scaffold/api/Dockerfile
grep -q "go mod download" scaffold/api/Dockerfile
grep -q "go build" scaffold/api/Dockerfile
grep -q "COPY api/go.mod api/go.sum ./api/" scaffold/api/Dockerfile
grep -q "COPY db/go.mod db/go.sum ./db/" scaffold/api/Dockerfile
grep -q "COPY api ./api" scaffold/api/Dockerfile
grep -q "COPY db ./db" scaffold/api/Dockerfile
! grep -q "COPY view ./view" scaffold/api/Dockerfile
! grep -q "COPY model ./model" scaffold/api/Dockerfile
! grep -q "COPY controller ./controller" scaffold/api/Dockerfile
! grep -q "COPY ui_components ./ui_components" scaffold/api/Dockerfile
! grep -q "gcc" scaffold/api/Dockerfile
! grep -q "libpq-dev" scaffold/api/Dockerfile
! test -f scaffold/web/package-lock.json
! test -f scaffold/web/vite.config.js
grep -q "oven/bun:1.3.14-debian@sha256:9dba1a1b43ce28c9d7931bfc4eb00feb63b0114720a0277a8f939ae4dfc9db6f" scaffold/web/Dockerfile
grep -q "bun install --frozen-lockfile" scaffold/web/Dockerfile
grep -q '"@tailwindcss/cli": "4.3.2"' scaffold/web/package.json
grep -q '"tailwindcss": "4.3.2"' scaffold/web/package.json
grep -q '"react": "19.2.7"' scaffold/web/package.json
grep -q "go 1.25.0" scaffold/api/go.mod
grep -q "go 1.25.0" scaffold/db/go.mod
! grep -q "go 1.23.0" scaffold/api/go.mod
! grep -q "go 1.23.0" scaffold/db/go.mod
! grep -Eq '"(react|react-dom|tailwindcss|@tailwindcss/cli)": "[~^<>*xX]|"(react|react-dom|tailwindcss|@tailwindcss/cli)": "latest' scaffold/web/package.json
grep -q "Bun.serve" scaffold/web/src/server.jsx
grep -q "browser entrypoint" scaffold/web/src/server.jsx
grep -q "listening inside container" scaffold/web/src/server.jsx
grep -q "proxying /api and /health to api service" scaffold/web/src/server.jsx
grep -q "publicRoot" scaffold/web/src/server.jsx
grep -q "Cache-Control" scaffold/web/src/server.jsx
grep -q "public, max-age=31536000, immutable" scaffold/web/src/server.jsx
grep -q "return 'no-store'" scaffold/web/src/server.jsx
grep -q '"assets:build"' scaffold/web/package.json
grep -F -q "assets/[name]-[hash].[ext]" scaffold/web/package.json
grep -q "bun run assets:build" scaffold/web/Dockerfile
grep -q "asset-manifest.json" scaffold/web/src/write-index.mjs
grep -F -q '/assets/${scripts[0]}' scaffold/web/src/write-index.mjs
! grep -q "Bun frontend listening on http://localhost" scaffold/web/src/server.jsx
grep -q '@import "tailwindcss";' scaffold/web/src/styles.css
grep -q "@theme" scaffold/web/src/styles.css
grep -q -- "--color-carbide-action" scaffold/web/src/styles.css
grep -q "\\[data-theme=\"dark\"\\]" scaffold/web/src/styles.css
grep -q "color-scheme: dark" scaffold/web/src/styles.css
grep -q "font-size: 14px" scaffold/web/src/styles.css
grep -q "line-height: 1.4" scaffold/web/src/styles.css
grep -q "var(--carbide-page)" scaffold/web/src/styles.css
grep -q -- "--carbide-page: #ffffff" scaffold/web/src/styles.css
grep -q -- "--carbide-page: #000000" scaffold/web/src/styles.css
grep -q "bg-carbide-hero text-carbide-hero-text" scaffold/web/src/component/l1/tokens.js
! grep -Eq "#0f766e|#115e59|#2dd4bf|#5eead4|#16433c|#0f302c|#16211b|#edf5ef|#ecfdf5|#166534" scaffold/web/src/styles.css
! grep -q "from-carbide-action via-carbide-hero-via" scaffold/web/src/component/l1/tokens.js
! grep -q "theme.css" scaffold/web/src/styles.css
grep -F -q "text-2xl/8 sm:text-3xl/9" scaffold/web/src/component/l1/Text.jsx
grep -F -q "min-h-8 rounded-md border px-2 py-1 text-sm/6" scaffold/web/src/component/l1/Field.jsx
grep -F -q "md: 'min-h-8 px-3 text-xs'" scaffold/web/src/component/l1/Button.jsx
grep -F -q "gap-3 border-l px-4 py-5" scaffold/web/src/component/l2/AuthForm.jsx
grep -F -q "w-full max-w-sm justify-self-center gap-3" scaffold/web/src/component/l2/AuthForm.jsx
grep -F -q "lg:grid-cols-[216px_minmax(0,1fr)]" scaffold/web/src/component/l2/Layouts.jsx
grep -F -q "px-3 py-4 sm:px-5 lg:py-5" scaffold/web/src/component/l2/Layouts.jsx
! grep -R -E "text-7xl|text-5xl|py-24|lg:py-12|min-h-12 rounded-md border|min-h-10 rounded-md border|lg:grid-cols-\[280px|lg:grid-cols-\[240px|gap-6|p-6|font-extrabold" scaffold/web/src/component >/dev/null
grep -q '/api/${mode}' scaffold/web/src/main.jsx
grep -q "carbide.theme" scaffold/web/src/main.jsx
grep -q "useThemeMode" scaffold/web/src/main.jsx
grep -q "prefers-color-scheme: dark" scaffold/web/index.html
grep -q "dataset.theme" scaffold/web/index.html
grep -q "ThemeToggle" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "data-resolved-theme" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "data-theme-mode" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "SunIcon" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "MoonIcon" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "Switch to light theme" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "Switch to dark theme" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "size-8 rounded-full border" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "onClick={() => onMode?.(nextMode)}" scaffold/web/src/component/l1/ThemeToggle.jsx
! grep -q "<select" scaffold/web/src/component/l1/ThemeToggle.jsx
! grep -q "appearance-none" scaffold/web/src/component/l1/ThemeToggle.jsx
! grep -q "border-x-4 border-t-4 border-x-transparent" scaffold/web/src/component/l1/ThemeToggle.jsx
! grep -q "aria-pressed" scaffold/web/src/component/l1/ThemeToggle.jsx
! grep -q "role=\"group\"" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "./component/l3/index.js" scaffold/web/src/main.jsx
grep -q "AuthView" scaffold/web/src/main.jsx
grep -q "DashboardView" scaffold/web/src/main.jsx
grep -q "LoadingView" scaffold/web/src/main.jsx
grep -R -q "Lorem ipsum dolor sit amet" scaffold/web/src/component
grep -R -q "Consectetur adipiscing elit" scaffold/web/src/component
grep -R -q "Sed do eiusmod tempor" scaffold/web/src/component
! grep -R -q "Create the owner account" scaffold/web/src/component
! grep -R -q "Use your account email" scaffold/web/src/component
! grep -R -q "Your session is active" scaffold/web/src/component
! grep -R -q "Bun + Go + Postgres" scaffold/web/src/component
! grep -R -q "React + Bun container" scaffold/web/src/component
! grep -R -q "React and Tailwind" scaffold/web/src/component
! grep -R -q "Go owns" scaffold/web/src/component
! grep -R -q "Postgres" scaffold/web/src/component
grep -q "export function Button" scaffold/web/src/component/l1/Button.jsx
grep -q "export function Field" scaffold/web/src/component/l1/Field.jsx
grep -q "export function Panel" scaffold/web/src/component/l1/Surface.jsx
grep -q "export function ThemeToggle" scaffold/web/src/component/l1/ThemeToggle.jsx
grep -q "ThemeToggle" scaffold/web/src/component/l1/index.js
grep -q "export const ui" scaffold/web/src/component/l1/tokens.js
! test -f scaffold/web/src/component/l1/theme.css
! grep -R "cb-" scaffold/web/src >/dev/null
! grep -R -- "--cb-" scaffold/web/src >/dev/null
grep -q "ui.action" scaffold/web/src/component/l1/Button.jsx
grep -q "ui.input" scaffold/web/src/component/l1/Field.jsx
grep -q "buttonClassLayers" scaffold/web/src/component/l1/Button.jsx
grep -q "fieldClassLayers" scaffold/web/src/component/l1/Field.jsx
grep -q "fieldHintClassLayers" scaffold/web/src/component/l1/Field.jsx
grep -q "fieldErrorClassLayers" scaffold/web/src/component/l1/Field.jsx
grep -q "inputClassLayers" scaffold/web/src/component/l1/Field.jsx
grep -q "panelClassLayers" scaffold/web/src/component/l1/Surface.jsx
grep -q "dividerClassLayers" scaffold/web/src/component/l1/Surface.jsx
grep -q "badgeClassLayers" scaffold/web/src/component/l1/Surface.jsx
grep -q "metricClassLayers" scaffold/web/src/component/l1/Surface.jsx
grep -q "eyebrowClassLayers" scaffold/web/src/component/l1/Text.jsx
grep -q "headingClassLayers" scaffold/web/src/component/l1/Text.jsx
grep -q "mutedClassLayers" scaffold/web/src/component/l1/Text.jsx
grep -q "codeClassLayers" scaffold/web/src/component/l1/Text.jsx
grep -q "formClassLayers" scaffold/web/src/component/l2/AuthForm.jsx
grep -q "formStackClassLayers" scaffold/web/src/component/l2/AuthForm.jsx
grep -q "errorClassLayers" scaffold/web/src/component/l2/AuthForm.jsx
grep -q "modeButtonClassLayers" scaffold/web/src/component/l2/AuthForm.jsx
grep -q "landingClassLayers" scaffold/web/src/component/l2/Layouts.jsx
grep -q "dashboardClassLayers" scaffold/web/src/component/l2/Layouts.jsx
grep -q "screenClassLayers" scaffold/web/src/component/l3/DashboardView.jsx
grep -q "loadingClassLayers" scaffold/web/src/component/l3/LoadingView.jsx
grep -q "ui.focus" scaffold/web/src/component/l1/Field.jsx
grep -q "ui.focus" scaffold/web/src/component/l2/AuthForm.jsx
grep -q "ui.focus" scaffold/web/src/component/l2/Layouts.jsx
! grep -R "text-\\[" scaffold/web/src/component >/dev/null
grep -q "export function DashboardLayout" scaffold/web/src/component/l2/Layouts.jsx
grep -q "lg:grid-cols-\\[216px_minmax(0,1fr)\\]" scaffold/web/src/component/l2/Layouts.jsx
grep -q "aria-label=\"Dashboard\"" scaffold/web/src/component/l2/Layouts.jsx
grep -q "aria-current" scaffold/web/src/component/l2/Layouts.jsx
grep -q "navItems" scaffold/web/src/component/l2/Layouts.jsx
grep -q "export function LandingPageLayout" scaffold/web/src/component/l2/Layouts.jsx
grep -q "export function AuthView" scaffold/web/src/component/l3/AuthView.jsx
grep -q "export function DashboardView" scaffold/web/src/component/l3/DashboardView.jsx
grep -q "dashboardNav" scaffold/web/src/component/l3/DashboardView.jsx
grep -q "WorkspaceOverview" scaffold/web/src/component/l3/DashboardView.jsx
! grep -R "ComponentLibraryView" scaffold/web/src/component >/dev/null
grep -q "module carbideapp/api" scaffold/api/go.mod
grep -q "carbideapp/db" scaffold/api/go.mod
grep -q "replace carbideapp/db => ../db" scaffold/api/go.mod
grep -q "module carbideapp/db" scaffold/db/go.mod
grep -q "github.com/jackc/pgx/v5" scaffold/db/go.mod
grep -q "package main" scaffold/api/main.go
grep -q "/api/login" scaffold/api/routes.go
grep -q "/api/me" scaffold/api/routes.go
grep -q "handleDashboard" scaffold/api/routes.go
grep -q "CreateUser" scaffold/db/user.go
grep -q "CreateSession" scaffold/db/session.go
! grep -R "admin@carbide.local" scaffold README.md docs >/dev/null
! grep -R "Demo login" scaffold README.md docs >/dev/null
! find scaffold -name '*.c' -o -name '*.h' | grep -q .
! grep -R "seed_admin" scaffold >/dev/null
! grep -R "render_template_text" scaffold >/dev/null
! grep -R "respond_view" scaffold >/dev/null
! find scaffold -path '*/ui_components/*' -print -quit | grep -q .
test -d scaffold/web/src/component/l1
test -d scaffold/web/src/component/l2
test -d scaffold/web/src/component/l3
test "$(find scaffold/web/src/component -mindepth 1 -maxdepth 1 -type d -printf '%f\n' | sort | tr '\n' ' ')" = "l1 l2 l3 "
! find scaffold/web/src/component -mindepth 1 -maxdepth 1 -type f -print -quit | grep -q .
! test -d scaffold/web/src/component/ui
! test -d scaffold/web/src/component/screen
! grep -R "views/" scaffold README.md docs >/dev/null
grep -q "api listening on container port" scaffold/api/main.go
grep -q "public API URL is" scaffold/api/main.go
! grep -q "API listening inside api container" scaffold/api/main.go
! grep -q "frontend proxies API calls" scaffold/api/main.go
grep -q "composeFilePath" cli/internal/cli/cli.go
grep -q "COMPOSE_FILE" cli/internal/cli/cli.go
grep -q "compose.supports(\"--watch\")" cli/internal/cli/cli.go
grep -q "newRenderer" cli/internal/cli/cli.go
grep -q "func (r renderer) Table" cli/internal/cli/cli.go
grep -q "runDevStreams" cli/internal/cli/cli.go
grep -q "commandStatus" cli/internal/cli/cli.go
grep -q "commandStopDev" cli/internal/cli/cli.go
grep -q "RunServiceProgress" cli/internal/cli/cli.go
grep -q "RunServiceStopProgress" cli/internal/cli/cli.go
grep -q "serviceProgressFrameWidth" cli/internal/cli/cli.go
grep -q "serviceProgressFrame" cli/internal/cli/cli.go
grep -q "terminalColumns" cli/internal/cli/cli.go
grep -q "composeServiceStatuses" cli/internal/cli/cli.go
grep -q "composeServiceSnapshots" cli/internal/cli/cli.go
grep -q "composePublishedPorts" cli/internal/cli/cli.go
grep -q "composeInternalPorts" cli/internal/cli/cli.go
grep -q "streamLogOutput" cli/internal/cli/cli.go
grep -q "parseComposeLogLine" cli/internal/cli/cli.go
grep -q "composeLogsArgs" cli/internal/cli/cli.go
grep -q "openDevLogSink" cli/internal/cli/cli.go
grep -q "openAppendDevLogSink" cli/internal/cli/cli.go
grep -q "commandLogs" cli/internal/cli/cli.go
grep -q "commandFollowLogs" cli/internal/cli/cli.go
grep -q ".carbide/log/dev.jsonl" cli/internal/cli/cli.go
grep -q "carbide follow logs" cli/internal/cli/cli.go
grep -q "carbide status" cli/internal/cli/cli.go
grep -q "doctor framework" cli/internal/cli/cli.go
! grep -q "carbide logs follow" cli/internal/cli/cli.go
! grep -q 'outputRow{"login"' cli/internal/cli/cli.go
! grep -q 'outputRow{"mode"' cli/internal/cli/cli.go

grep -q "$domain" docs/site/index.html
grep -q "Bun frontend" docs/site/index.html
grep -q "Create Your First App" docs/site/index.html
grep -q 'href="/version-policy"' docs/site/index.html
grep -q 'href="/create-your-first-app"' docs/site/index.html
grep -q 'href="/frontend-starter-contract"' docs/site/index.html
grep -q 'href="/deployment"' docs/site/index.html
grep -q 'href="/ci-cd-regression-tests"' docs/site/index.html
! grep -R -E 'href="[^"]+\.html' docs/site >/dev/null
grep -q "canonicalDocsPath" docs/app/web/src/server.jsx
grep -q '"/initial-user-experience": "/create-your-first-app"' docs/app/web/src/server.jsx
grep -q 'pathname === "/index.html"' docs/app/web/src/server.jsx
grep -q 'pathname.endsWith(".html")' docs/app/web/src/server.jsx
grep -q "status: 308" docs/app/web/src/server.jsx
grep -q 'location: `${pathname}${target.search}`' docs/app/web/src/server.jsx
grep -q 'docsResponseHeaders' docs/app/web/src/server.jsx
grep -q '@import "tailwindcss";' docs/app/web/src/styles.css
grep -F -q '@source "./component/**/*.jsx";' docs/app/web/src/styles.css
grep -q '"tailwind:build"' docs/app/web/package.json
grep -q "tailwindcss" docs/app/web/src/build-styles.js
grep -q '"@tailwindcss/cli": "4.3.2"' docs/app/web/package.json
grep -q '"tailwindcss": "4.3.2"' docs/app/web/package.json
grep -q '"react": "19.2.7"' docs/app/web/package.json
grep -q '"react-dom": "19.2.7"' docs/app/web/package.json
grep -q "bun run tailwind:build" docs/app/web/Dockerfile
grep -q "docsClassLayers" docs/app/web/src/component/l1/tokens.js
grep -q "docsChromeClassLayers" docs/app/web/src/component/l2/DocsChrome.jsx
grep -q "docsWebContract" docs/app/web/src/component/l3/DocsSite.jsx
grep -q "component/l1" docs/app/agents.d/TAILWIND_COMPONENTS.md
grep -q "component/l2" docs/app/agents.d/TAILWIND_COMPONENTS.md
grep -q "component/l3" docs/app/agents.d/TAILWIND_COMPONENTS.md
grep -q "Bun frontend, Go API backend, Postgres database" docs/site/frontend-starter-contract.html
grep -q "Tailwind is required" docs/site/frontend-starter-contract.html
grep -q "carbide follow logs" docs/site/create-your-first-app.html
grep -q "carbide status" docs/site/create-your-first-app.html
grep -q "carbide doctor runtime" docs/site/create-your-first-app.html
grep -q "Install, create, run, register" docs/site/create-your-first-app.html
grep -q "Single VM" docs/site/deployment.html
grep -q "Multiple VMs" docs/site/deployment.html
grep -q 'type = "ssh-compose"' docs/site/deployment.html
grep -q 'type = "ssh-compose-environment"' docs/site/deployment.html
grep -q "clustered orchestration is implemented" docs/site/deployment.html
grep -q "CI/CD regression plan" docs/site/ci-cd-regression-tests.html
grep -q "carbide doctor framework" docs/site/ci-cd-regression-tests.html
grep -q "Directory structure" docs/site/repo-structure.html
grep -q 'class="docs-topbar"' docs/site/index.html
grep -q 'class="docs-sidebar"' docs/site/index.html
grep -q 'class="docs-content"' docs/site/index.html
grep -q 'class="docs-toc"' docs/site/index.html
grep -q "Search docs" docs/site/index.html
grep -q "Version v0.1" docs/site/index.html
grep -q "Prologue" docs/site/index.html
grep -q "Getting Started" docs/site/index.html
grep -q "Architecture" docs/site/index.html
grep -q "On this page" docs/site/index.html
grep -q ".docs-layout" docs/site/assets/styles.css
grep -q ".docs-sidebar" docs/site/assets/styles.css
grep -q ".docs-toc" docs/site/assets/styles.css
grep -q ".docs-topbar" docs/site/assets/styles.css
grep -E -q '@media \(max-width: ?860px\)' docs/site/assets/styles.css

for page in docs/site/*.html; do
  grep -q 'class="docs-topbar"' "$page"
  grep -q 'class="docs-sidebar"' "$page"
  grep -q 'class="docs-content"' "$page"
  grep -q 'class="docs-toc"' "$page"
  grep -q "Search docs" "$page"
  grep -q "Version v0.1" "$page"
  grep -q "On this page" "$page"
done

printf 'repo contract ok\n'
