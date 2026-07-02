#!/usr/bin/env bash
set -euo pipefail

domain="carbide.ryangerardwilson.com"

required_files=(
  ".gitignore"
  "README.md"
  "logo.txt"
  "install.sh"
  "go.mod"
  "bin/carbide"
  "cmd/carbide/main.go"
  "internal/carbide/cli.go"
  "internal/carbide/cli_test.go"
  ".github/workflows/ci.yml"
  ".github/workflows/pages.yml"
  "docs/engineering/CI_CD_REGRESSION_TESTS.md"
  "docs/engineering/COMPONENT_STYLE_SYSTEM.md"
  "docs/engineering/DIRECTORY_STRUCTURE.md"
  "docs/engineering/INITIAL_USER_EXPERIENCE.md"
  "docs/site/CNAME"
  "docs/site/index.html"
  "docs/site/component-style-system.html"
  "docs/site/initial-user-experience.html"
  "docs/site/ci-cd-regression-tests.html"
  "docs/site/repo-structure.html"
  "docs/site/assets/styles.css"
  "scripts/test_cli_scaffold.sh"
  "scripts/test_starter_docker_flow.sh"
  "templates/default/Dockerfile"
  "templates/default/.env.example"
  "templates/default/.gitignore"
  "templates/default/config/env.schema.json"
  "templates/default/doc/runbook/backup-restore.md"
  "templates/default/doc/runbook/deploy.md"
  "templates/default/doc/runbook/env.md"
  "templates/default/docker-compose.yml"
  "templates/default/view/web/Dockerfile"
  "templates/default/view/web/index.html"
  "templates/default/view/web/package.json"
  "templates/default/view/web/bun.lock"
  "templates/default/view/web/src/main.jsx"
  "templates/default/view/web/src/server.jsx"
  "templates/default/view/web/src/styles.css"
  "templates/default/view/web/src/component/utils.js"
  "templates/default/view/web/src/component/l1/Button.jsx"
  "templates/default/view/web/src/component/l1/Field.jsx"
  "templates/default/view/web/src/component/l1/Surface.jsx"
  "templates/default/view/web/src/component/l1/Text.jsx"
  "templates/default/view/web/src/component/l1/index.js"
  "templates/default/view/web/src/component/l1/theme.css"
  "templates/default/view/web/src/component/l1/tokens.js"
  "templates/default/view/web/src/component/l2/Accordion.jsx"
  "templates/default/view/web/src/component/l2/AuthForm.jsx"
  "templates/default/view/web/src/component/l2/Carousel.jsx"
  "templates/default/view/web/src/component/l2/CarouselIntegrations.jsx"
  "templates/default/view/web/src/component/l2/Charts.jsx"
  "templates/default/view/web/src/component/l2/Combobox.jsx"
  "templates/default/view/web/src/component/l2/DateComponents.jsx"
  "templates/default/view/web/src/component/l2/Dropdown.jsx"
  "templates/default/view/web/src/component/l2/EnhancedSelects.jsx"
  "templates/default/view/web/src/component/l2/Layouts.jsx"
  "templates/default/view/web/src/component/l2/Lessons.jsx"
  "templates/default/view/web/src/component/l2/Listbox.jsx"
  "templates/default/view/web/src/component/l2/Modal.jsx"
  "templates/default/view/web/src/component/l2/Notifications.jsx"
  "templates/default/view/web/src/component/l2/Popover.jsx"
  "templates/default/view/web/src/component/l2/RadioGroup.jsx"
  "templates/default/view/web/src/component/l2/Tabs.jsx"
  "templates/default/view/web/src/component/l2/TextEditors.jsx"
  "templates/default/view/web/src/component/l2/Toggle.jsx"
  "templates/default/view/web/src/component/l2/Tooltip.jsx"
  "templates/default/view/web/src/component/l2/index.js"
  "templates/default/view/web/src/component/l3/AuthView.jsx"
  "templates/default/view/web/src/component/l3/ComponentLibraryView.jsx"
  "templates/default/view/web/src/component/l3/DashboardView.jsx"
  "templates/default/view/web/src/component/l3/LoadingView.jsx"
  "templates/default/view/web/src/component/l3/index.js"
  "templates/default/carbide.toml"
  "templates/default/go.mod"
  "templates/default/go.sum"
  "templates/default/src/main.go"
  "templates/default/model/user.go"
  "templates/default/model/session.go"
  "templates/default/controller/auth_controller.go"
  "templates/default/controller/page_controller.go"
  "templates/default/migrations/001_auth.sql"
)

required_dirs=(
  "cmd"
  "cmd/carbide"
  "internal/carbide"
  "src"
  "src/ui"
  "include/carbide"
  "include/carbide/ui"
  "tests/unit"
  "tests/integration"
  "tests/regression"
  "tests/fixtures"
  "examples/hello"
  "infra/compose"
  "infra/schemas"
  "templates/default"
  "templates/default/config"
  "templates/default/doc"
  "templates/default/doc/runbook"
  "templates/default/view"
  "templates/default/view/web"
  "templates/default/view/web/src"
  "templates/default/view/web/src/component"
  "templates/default/view/web/src/component/l1"
  "templates/default/view/web/src/component/l2"
  "templates/default/view/web/src/component/l3"
  "templates/default/src"
  "templates/default/model"
  "templates/default/controller"
  "templates/default/migrations"
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

grep -q "Bun/React/Tailwind frontend container" README.md
grep -q "Postgres-only" README.md
grep -q "Separate runtime boundaries" README.md
grep -q "Infrastructure as code" README.md
grep -q "generated Docker Compose setup" README.md
grep -q "Postgres-backed queues" README.md
grep -q "carbide new" README.md
grep -q "carbide run dev" README.md
grep -q "carbide status" README.md
grep -q "carbide stop dev" README.md
grep -q "carbide follow logs" README.md
grep -q "carbide logs" README.md
grep -q "carbide doctor env" README.md
grep -q "carbide deploy preview" README.md
grep -q "carbide deploy apply" README.md
! grep -q "command_format" bin/carbide
! grep -q "carbide format" bin/carbide
grep -q "module github.com/ryangerardwilson/carbide" go.mod
grep -q "oo_______oo_______oo" logo.txt
grep -q "package main" cmd/carbide/main.go
grep -q "package carbide" internal/carbide/cli.go
grep -q "commandDoctorEnv" internal/carbide/cli.go
grep -q "commandDeployPreview" internal/carbide/cli.go
grep -q "commandDeployApply" internal/carbide/cli.go
! git grep -n -e 'S[e]alion' -e 's[e]alion' -e 'S[E]ALION' -- .
grep -q "composeUpDetached" internal/carbide/cli.go
grep -q "runDevStreams" internal/carbide/cli.go
grep -q -- "--quiet-build" internal/carbide/cli.go
grep -q "Carbide dev" internal/carbide/cli.go
grep -q "Go is required to build the Carbide CLI" install.sh
grep -q ".bin/carbide" install.sh
grep -q "default_port = 8080" templates/default/carbide.toml
grep -q 'schema = "config/env.schema.json"' templates/default/carbide.toml
grep -q "preview_before_apply = true" templates/default/carbide.toml
grep -q ".carbide/" templates/default/.gitignore
grep -q ".env" templates/default/.gitignore
grep -q "config/env.schema.json" templates/default/README.md
grep -q "carbide doctor env" templates/default/README.md
grep -q "carbide deploy preview dev" templates/default/README.md
grep -q "carbide deploy apply dev" templates/default/README.md
grep -q "POSTGRES_PASSWORD" templates/default/.env.example
grep -q '"name": "DATABASE_URL"' templates/default/config/env.schema.json
grep -q '"secret": true' templates/default/config/env.schema.json
grep -q '"browser_exposed": true' templates/default/config/env.schema.json
grep -q '"framework_owned": true' templates/default/config/env.schema.json
grep -q "separate secrets container" templates/default/doc/runbook/env.md
grep -q "preview-before-apply" templates/default/doc/runbook/deploy.md
grep -q "Postgres owns durable application state" templates/default/doc/runbook/backup-restore.md
! grep -q 'url = "http://localhost:8080"' templates/default/carbide.toml
grep -q "frontend:" templates/default/docker-compose.yml
grep -q "backend:" templates/default/docker-compose.yml
grep -q "db:" templates/default/docker-compose.yml
grep -q 'PUBLIC_URL: "http://localhost:${CARBIDE_HTTP_PORT:-8080}"' templates/default/docker-compose.yml
test "$(grep -c 'PUBLIC_URL: "http://localhost:${CARBIDE_HTTP_PORT:-8080}"' templates/default/docker-compose.yml)" -eq 2
grep -q 'PUBLIC_APP_NAME: "${PUBLIC_APP_NAME:-__PROJECT_NAME__}"' templates/default/docker-compose.yml
grep -q 'APP_ENV: "${APP_ENV:-development}"' templates/default/docker-compose.yml
grep -q 'POSTGRES_PASSWORD: "${POSTGRES_PASSWORD:-carbide}"' templates/default/docker-compose.yml
grep -q "develop:" templates/default/docker-compose.yml
grep -q "watch:" templates/default/docker-compose.yml
grep -q "action: rebuild" templates/default/docker-compose.yml
grep -q "context: ./view/web" templates/default/docker-compose.yml
grep -q "path: ./view/web/src" templates/default/docker-compose.yml
grep -q "path: ./view/web/package.json" templates/default/docker-compose.yml
grep -q "path: ./view/web/bun.lock" templates/default/docker-compose.yml
grep -q "path: ./go.mod" templates/default/docker-compose.yml
grep -q "path: ./go.sum" templates/default/docker-compose.yml
grep -q "path: ./src" templates/default/docker-compose.yml
grep -q "path: ./model" templates/default/docker-compose.yml
grep -q "path: ./controller" templates/default/docker-compose.yml
grep -q "path: ./Dockerfile" templates/default/docker-compose.yml
grep -q "FROM golang:" templates/default/Dockerfile
grep -q "go mod download" templates/default/Dockerfile
grep -q "go build" templates/default/Dockerfile
grep -q "COPY model ./model" templates/default/Dockerfile
grep -q "COPY controller ./controller" templates/default/Dockerfile
! grep -q "COPY view ./view" templates/default/Dockerfile
! grep -q "COPY ui_components ./ui_components" templates/default/Dockerfile
! grep -q "gcc" templates/default/Dockerfile
! grep -q "libpq-dev" templates/default/Dockerfile
! test -d templates/default/frontend
! test -f templates/default/view/web/package-lock.json
! test -f templates/default/view/web/vite.config.js
grep -q "oven/bun:1.3.14-debian" templates/default/view/web/Dockerfile
grep -q "bun install --frozen-lockfile" templates/default/view/web/Dockerfile
grep -q '"@tailwindcss/cli": "4.3.2"' templates/default/view/web/package.json
grep -q '"tailwindcss": "4.3.2"' templates/default/view/web/package.json
grep -q '"react": "19.2.7"' templates/default/view/web/package.json
grep -q "Bun.serve" templates/default/view/web/src/server.jsx
grep -q "browser entrypoint" templates/default/view/web/src/server.jsx
grep -q "listening inside container" templates/default/view/web/src/server.jsx
grep -q "proxying /api and /health to backend service" templates/default/view/web/src/server.jsx
! grep -q "Bun frontend listening on http://localhost" templates/default/view/web/src/server.jsx
grep -q '@import "tailwindcss";' templates/default/view/web/src/styles.css
grep -q '@import "./component/l1/theme.css";' templates/default/view/web/src/styles.css
grep -q '/api/${mode}' templates/default/view/web/src/main.jsx
grep -q "./component/l3/index.js" templates/default/view/web/src/main.jsx
grep -q "AuthView" templates/default/view/web/src/main.jsx
grep -q "DashboardView" templates/default/view/web/src/main.jsx
grep -q "LoadingView" templates/default/view/web/src/main.jsx
grep -R -q "Bun + Go + Postgres" templates/default/view/web/src/component
grep -R -q "React + Bun container" templates/default/view/web/src/component
grep -q "export function Button" templates/default/view/web/src/component/l1/Button.jsx
grep -q "export function Field" templates/default/view/web/src/component/l1/Field.jsx
grep -q "export function Panel" templates/default/view/web/src/component/l1/Surface.jsx
grep -q "export const color" templates/default/view/web/src/component/l1/tokens.js
grep -q "export const ui" templates/default/view/web/src/component/l1/tokens.js
grep -q -- "--cb-font-sans" templates/default/view/web/src/component/l1/theme.css
grep -q -- "--cb-color-action" templates/default/view/web/src/component/l1/theme.css
grep -q -- "--cb-color-surface" templates/default/view/web/src/component/l1/theme.css
grep -q "cb-action" templates/default/view/web/src/component/l1/theme.css
grep -q "cb-input" templates/default/view/web/src/component/l1/theme.css
grep -q "ui.action" templates/default/view/web/src/component/l1/Button.jsx
grep -q "ui.input" templates/default/view/web/src/component/l1/Field.jsx
grep -q "export function Dropdown" templates/default/view/web/src/component/l2/Dropdown.jsx
grep -q "export const Menu = Dropdown" templates/default/view/web/src/component/l2/Dropdown.jsx
grep -q "export function Modal" templates/default/view/web/src/component/l2/Modal.jsx
grep -q "export const Dialog = Modal" templates/default/view/web/src/component/l2/Modal.jsx
grep -q "export function Slideover" templates/default/view/web/src/component/l2/Modal.jsx
grep -q "export function Accordion" templates/default/view/web/src/component/l2/Accordion.jsx
grep -q "export const Disclosure = Accordion" templates/default/view/web/src/component/l2/Accordion.jsx
grep -q "export function Carousel" templates/default/view/web/src/component/l2/Carousel.jsx
grep -q "export function Tabs" templates/default/view/web/src/component/l2/Tabs.jsx
grep -q "export function Notifications" templates/default/view/web/src/component/l2/Notifications.jsx
grep -q "export function RadioGroup" templates/default/view/web/src/component/l2/RadioGroup.jsx
grep -q "export function Radio" templates/default/view/web/src/component/l2/RadioGroup.jsx
grep -q "export function Toggle" templates/default/view/web/src/component/l2/Toggle.jsx
grep -q "export const Switch = Toggle" templates/default/view/web/src/component/l2/Toggle.jsx
grep -q "export function Tooltip" templates/default/view/web/src/component/l2/Tooltip.jsx
grep -q "export function Popover" templates/default/view/web/src/component/l2/Popover.jsx
grep -q "export function Listbox" templates/default/view/web/src/component/l2/Listbox.jsx
grep -q "export function Combobox" templates/default/view/web/src/component/l2/Combobox.jsx
grep -q "export function Lessons" templates/default/view/web/src/component/l2/Lessons.jsx
grep -q "export function DashboardLayout" templates/default/view/web/src/component/l2/Layouts.jsx
grep -q "lg:grid-cols-\\[280px_minmax(0,1fr)\\]" templates/default/view/web/src/component/l2/Layouts.jsx
grep -q "aria-label=\"Dashboard\"" templates/default/view/web/src/component/l2/Layouts.jsx
grep -q "aria-current" templates/default/view/web/src/component/l2/Layouts.jsx
grep -q "navItems" templates/default/view/web/src/component/l2/Layouts.jsx
grep -q "export function LandingPageLayout" templates/default/view/web/src/component/l2/Layouts.jsx
grep -q "export function TrixEditor" templates/default/view/web/src/component/l2/TextEditors.jsx
grep -q "export function QuillEditor" templates/default/view/web/src/component/l2/TextEditors.jsx
grep -q "export function SimpleMDEEditor" templates/default/view/web/src/component/l2/TextEditors.jsx
grep -q "export function ChartJsPanel" templates/default/view/web/src/component/l2/Charts.jsx
grep -q "export function ApexChartsPanel" templates/default/view/web/src/component/l2/Charts.jsx
grep -q "export function Select2Select" templates/default/view/web/src/component/l2/EnhancedSelects.jsx
grep -q "export function ChoicesSelect" templates/default/view/web/src/component/l2/EnhancedSelects.jsx
grep -q "export function FlatpickrPicker" templates/default/view/web/src/component/l2/DateComponents.jsx
grep -q "export function DateRangePicker" templates/default/view/web/src/component/l2/DateComponents.jsx
grep -q "export function FullCalendarPanel" templates/default/view/web/src/component/l2/DateComponents.jsx
grep -q "export function GlideCarousel" templates/default/view/web/src/component/l2/CarouselIntegrations.jsx
grep -q "export function SplideCarousel" templates/default/view/web/src/component/l2/CarouselIntegrations.jsx
grep -q "export function AuthView" templates/default/view/web/src/component/l3/AuthView.jsx
grep -q "export function DashboardView" templates/default/view/web/src/component/l3/DashboardView.jsx
grep -q "dashboardNav" templates/default/view/web/src/component/l3/DashboardView.jsx
grep -q "WorkspaceOverview" templates/default/view/web/src/component/l3/DashboardView.jsx
grep -q "onNavItem={setActiveSection}" templates/default/view/web/src/component/l3/DashboardView.jsx
grep -q "github.com/jackc/pgx/v5" templates/default/go.mod
grep -q "package main" templates/default/src/main.go
grep -q "/api/login" templates/default/controller/page_controller.go
grep -q "/api/me" templates/default/controller/page_controller.go
grep -q "handleDashboard" templates/default/controller/page_controller.go
grep -q "CreateUser" templates/default/model/user.go
grep -q "CreateSession" templates/default/model/session.go
! grep -R "admin@carbide.local" templates/default README.md docs >/dev/null
! grep -R "Demo login" templates/default README.md docs >/dev/null
! find templates/default -name '*.c' -o -name '*.h' | grep -q .
! grep -R "seed_admin" templates/default >/dev/null
! grep -R "render_template_text" templates/default >/dev/null
! grep -R "respond_view" templates/default >/dev/null
! find templates/default -path '*/ui_components/*' -print -quit | grep -q .
! grep -R "views/" templates/default README.md docs >/dev/null
grep -q "backend listening on container port" templates/default/src/main.go
grep -q "public API URL is" templates/default/src/main.go
! grep -q "API listening inside backend container" templates/default/src/main.go
! grep -q "frontend proxies API calls" templates/default/src/main.go
grep -q "compose.supports(\"--watch\")" internal/carbide/cli.go
grep -q "newRenderer" internal/carbide/cli.go
grep -q "func (r renderer) Table" internal/carbide/cli.go
grep -q "runDevStreams" internal/carbide/cli.go
grep -q "commandStatus" internal/carbide/cli.go
grep -q "commandStopDev" internal/carbide/cli.go
grep -q "RunServiceProgress" internal/carbide/cli.go
grep -q "RunServiceStopProgress" internal/carbide/cli.go
grep -q "serviceProgressFrameWidth" internal/carbide/cli.go
grep -q "serviceProgressFrame" internal/carbide/cli.go
grep -q "terminalColumns" internal/carbide/cli.go
grep -q "composeServiceStatuses" internal/carbide/cli.go
grep -q "composeServiceSnapshots" internal/carbide/cli.go
grep -q "composePublishedPorts" internal/carbide/cli.go
grep -q "composeInternalPorts" internal/carbide/cli.go
grep -q "streamLogOutput" internal/carbide/cli.go
grep -q "parseComposeLogLine" internal/carbide/cli.go
grep -q "composeLogsArgs" internal/carbide/cli.go
grep -q "openDevLogSink" internal/carbide/cli.go
grep -q "openAppendDevLogSink" internal/carbide/cli.go
grep -q "commandLogs" internal/carbide/cli.go
grep -q "commandFollowLogs" internal/carbide/cli.go
grep -q ".carbide/log/dev.jsonl" internal/carbide/cli.go
grep -q "carbide follow logs" internal/carbide/cli.go
grep -q "carbide status" internal/carbide/cli.go
! grep -q "carbide logs follow" internal/carbide/cli.go
! grep -q 'outputRow{"login"' internal/carbide/cli.go
! grep -q 'outputRow{"mode"' internal/carbide/cli.go

grep -q "$domain" docs/site/index.html
grep -q "Bun frontend" docs/site/index.html
grep -q "Initial user experience" docs/site/index.html
grep -q "Bun frontend, Go API backend, Postgres database" docs/site/component-style-system.html
grep -q "Tailwind is required" docs/site/component-style-system.html
grep -q "carbide follow logs" docs/site/initial-user-experience.html
grep -q "carbide status" docs/site/initial-user-experience.html
grep -q "Install, create, run, register" docs/site/initial-user-experience.html
grep -q "CI/CD regression plan" docs/site/ci-cd-regression-tests.html
grep -q "Directory structure" docs/site/repo-structure.html

printf 'repo contract ok\n'
