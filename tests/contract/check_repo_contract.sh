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
  ".github/workflows/pages.yml"
  "docs/engineering/CI_CD_REGRESSION_TESTS.md"
  "docs/engineering/FRONTEND_STARTER_CONTRACT.md"
  "docs/engineering/DIRECTORY_STRUCTURE.md"
  "docs/engineering/INITIAL_USER_EXPERIENCE.md"
  "docs/site/CNAME"
  "docs/site/index.html"
  "docs/site/frontend-starter-contract.html"
  "docs/site/initial-user-experience.html"
  "docs/site/ci-cd-regression-tests.html"
  "docs/site/repo-structure.html"
  "docs/site/assets/styles.css"
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
  "scaffold/web/src/styles.css"
  "scaffold/web/src/lib/cx.js"
  "scaffold/web/src/component/l1/Button.jsx"
  "scaffold/web/src/component/l1/Field.jsx"
  "scaffold/web/src/component/l1/Surface.jsx"
  "scaffold/web/src/component/l1/Text.jsx"
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

grep -qx "$domain" docs/site/CNAME || {
  printf 'docs/site/CNAME must contain only %s\n' "$domain" >&2
  exit 1
}

grep -Eq "Bun/React/Tailwind .*web.* container" README.md
grep -q "Postgres-only" README.md
grep -q "Separate runtime boundaries" README.md
grep -q "Infrastructure as code" README.md
grep -q "generated Docker Compose setup" README.md
grep -q "Postgres-backed queues" README.md
grep -q "carbide new" README.md
grep -q "carbide new My Carbide App" README.md
grep -q 'name = "My Carbide App"' README.md
grep -q 'slug = "my-carbide-app"' README.md
grep -q "carbide run dev" README.md
grep -q "carbide status" README.md
grep -q "carbide stop dev" README.md
grep -q "carbide follow logs" README.md
grep -q "carbide logs" README.md
grep -q "carbide doctor env" README.md
grep -q "carbide deploy preview" README.md
grep -q "carbide deploy apply" README.md
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
test "$(find scaffold -mindepth 1 -maxdepth 1 -type d -printf '%f\n' | sort | tr '\n' ' ')" = "agents.d api db web "
! test -f scaffold/Dockerfile
! test -f scaffold/go.mod
! test -f scaffold/go.sum
grep -q "oo_______oo_______oo" cli/internal/cli/cli.go
grep -q "package main" cli/cmd/carbide/main.go
grep -q "package cli" cli/internal/cli/cli.go
grep -q "commandDoctorEnv" cli/internal/cli/cli.go
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
grep -q "\\[env.variables.DATABASE_URL\\]" scaffold/carbide.toml
grep -q "secret = true" scaffold/carbide.toml
grep -q "browser_exposed = true" scaffold/carbide.toml
grep -q "framework_owned = true" scaffold/carbide.toml
grep -q "preview_before_apply = true" scaffold/carbide.toml
grep -q ".carbide/" scaffold/.gitignore
grep -q ".env" scaffold/.gitignore
grep -q "__PROJECT_NAME__ Agent Context" scaffold/AGENTS.md
grep -q 'agents.d/ENVIRONMENT.md' scaffold/AGENTS.md
grep -q 'agents.d/DEPLOY.md' scaffold/AGENTS.md
grep -q 'agents.d/BACKUP_RESTORE.md' scaffold/AGENTS.md
grep -q 'agents.d/TAILWIND_COMPONENTS.md' scaffold/AGENTS.md
grep -q "carbide.toml" scaffold/README.md
grep -q "carbide doctor env" scaffold/README.md
grep -q "carbide deploy preview dev" scaffold/README.md
grep -q "carbide deploy apply dev" scaffold/README.md
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
grep -q "FROM golang:" scaffold/api/Dockerfile
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
grep -q "oven/bun:1.3.14-debian" scaffold/web/Dockerfile
grep -q "bun install --frozen-lockfile" scaffold/web/Dockerfile
grep -q '"@tailwindcss/cli": "4.3.2"' scaffold/web/package.json
grep -q '"tailwindcss": "4.3.2"' scaffold/web/package.json
grep -q '"react": "19.2.7"' scaffold/web/package.json
grep -q "Bun.serve" scaffold/web/src/server.jsx
grep -q "browser entrypoint" scaffold/web/src/server.jsx
grep -q "listening inside container" scaffold/web/src/server.jsx
grep -q "proxying /api and /health to api service" scaffold/web/src/server.jsx
! grep -q "Bun frontend listening on http://localhost" scaffold/web/src/server.jsx
grep -q '@import "tailwindcss";' scaffold/web/src/styles.css
grep -q "@theme" scaffold/web/src/styles.css
grep -q -- "--color-carbide-action" scaffold/web/src/styles.css
! grep -q "theme.css" scaffold/web/src/styles.css
grep -q '/api/${mode}' scaffold/web/src/main.jsx
grep -q "./component/l3/index.js" scaffold/web/src/main.jsx
grep -q "AuthView" scaffold/web/src/main.jsx
grep -q "DashboardView" scaffold/web/src/main.jsx
grep -q "LoadingView" scaffold/web/src/main.jsx
grep -R -q "Bun + Go + Postgres" scaffold/web/src/component
grep -R -q "React + Bun container" scaffold/web/src/component
grep -q "export function Button" scaffold/web/src/component/l1/Button.jsx
grep -q "export function Field" scaffold/web/src/component/l1/Field.jsx
grep -q "export function Panel" scaffold/web/src/component/l1/Surface.jsx
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
grep -q "lg:grid-cols-\\[280px_minmax(0,1fr)\\]" scaffold/web/src/component/l2/Layouts.jsx
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
! grep -q "carbide logs follow" cli/internal/cli/cli.go
! grep -q 'outputRow{"login"' cli/internal/cli/cli.go
! grep -q 'outputRow{"mode"' cli/internal/cli/cli.go

grep -q "$domain" docs/site/index.html
grep -q "Bun frontend" docs/site/index.html
grep -q "Initial user experience" docs/site/index.html
grep -q "Bun frontend, Go API backend, Postgres database" docs/site/frontend-starter-contract.html
grep -q "Tailwind is required" docs/site/frontend-starter-contract.html
grep -q "carbide follow logs" docs/site/initial-user-experience.html
grep -q "carbide status" docs/site/initial-user-experience.html
grep -q "Install, create, run, register" docs/site/initial-user-experience.html
grep -q "CI/CD regression plan" docs/site/ci-cd-regression-tests.html
grep -q "Directory structure" docs/site/repo-structure.html

printf 'repo contract ok\n'
