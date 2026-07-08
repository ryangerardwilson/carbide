#!/usr/bin/env bash
set -euo pipefail

domain="carbide.ryangerardwilson.com"

repo_search() {
  local pattern="$1"
  shift
  if command -v rg >/dev/null 2>&1; then
    rg -n -- "$pattern" "$@"
    return
  fi
  grep -REn -- "$pattern" "$@"
}

repo_search_files() {
  local pattern="$1"
  shift
  if command -v rg >/dev/null 2>&1; then
    rg -l -- "$pattern" "$@"
    return
  fi
  grep -REl -- "$pattern" "$@"
}

required_files=(
  ".gitignore"
  "README.md"
  "cli/install.sh"
  "cli/go.mod"
  "cli/bin/carbide"
  "cli/cmd/carbide/main.go"
  "cli/internal/cli/types.go"
  "cli/internal/cli/cli.go"
  "cli/internal/cli/audit.go"
  "cli/internal/cli/health.go"
  "cli/internal/cli/dev.go"
  "cli/internal/cli/render.go"
  "cli/internal/cli/logs.go"
  "cli/internal/cli/config.go"
  "cli/internal/cli/project.go"
  "cli/internal/cli/compose.go"
  "cli/internal/cli/system.go"
  "cli/internal/cli/cli_test.go"
  "cli/internal/cli/render_test.go"
  "docs/app/carbide.toml"
  "docs/app/docker-compose.yml"
  "docs/app/web/Dockerfile"
  "docs/app/web/bun.lock"
  "docs/app/web/package.json"
  "docs/app/web/tsconfig.json"
  "docs/app/web/src/build-styles.ts"
  "docs/app/web/src/server.ts"
  "docs/app/web/src/styles.css"
  "docs/app/web/src/styles.d.ts"
  "docs/app/web/src/lib/cx.ts"
  "docs/app/web/src/lib/types.ts"
  "docs/app/web/src/component/l1/index.ts"
  "docs/app/web/src/component/l1/tokens.ts"
  "docs/app/web/src/component/l2/DocsChrome.ts"
  "docs/app/web/src/component/l2/index.ts"
  "docs/app/web/src/component/l3/DocsSite.ts"
  "docs/app/web/src/component/l3/index.ts"
  "docs/app/web/site/index.html"
  "docs/app/web/site/deployment.html"
  "docs/app/web/site/frontend-starter-contract.html"
  "docs/app/web/site/create-your-first-app.html"
  "docs/app/web/site/for/agents/index.md"
  "docs/app/web/site/ci-cd-regression-tests.html"
  "docs/app/web/site/repo-structure.html"
  "docs/app/web/site/version-policy.html"
  "docs/app/web/site/assets/intro.js"
  "docs/app/web/site/assets/styles.css"
  "tests/contract/audit_versions.sh"
  "tests/contract/check_line_limits.sh"
  "tests/contract/check_repo_contract.sh"
  "tests/scaffold/cli_scaffold.sh"
  "tests/smoke/starter_docker_flow.sh"
  "tests/smoke/docs_for_agents_http.sh"
  "scaffold/api/Dockerfile"
  "scaffold/.env.example"
  "scaffold/.gitignore"
  "scaffold/docker-compose.yml"
  "scaffold/web/Dockerfile"
  "scaffold/web/index.html"
  "scaffold/web/package.json"
  "scaffold/web/bun.lock"
  "scaffold/web/tsconfig.json"
  "scaffold/web/src/main.tsx"
  "scaffold/web/src/server.ts"
  "scaffold/web/src/write-index.ts"
  "scaffold/web/src/styles.css"
  "scaffold/web/src/styles.d.ts"
  "scaffold/web/src/lib/cx.ts"
  "scaffold/web/src/lib/types.ts"
  "scaffold/web/src/component/l1/Button.tsx"
  "scaffold/web/src/component/l1/Field.tsx"
  "scaffold/web/src/component/l1/Surface.tsx"
  "scaffold/web/src/component/l1/Text.tsx"
  "scaffold/web/src/component/l1/ThemeToggle.tsx"
  "scaffold/web/src/component/l1/index.ts"
  "scaffold/web/src/component/l1/tokens.ts"
  "scaffold/web/src/component/l2/AuthForm.tsx"
  "scaffold/web/src/component/l2/Layouts.tsx"
  "scaffold/web/src/component/l2/index.ts"
  "scaffold/web/src/component/l3/AuthView.tsx"
  "scaffold/web/src/component/l3/DashboardView.tsx"
  "scaffold/web/src/component/l3/LoadingView.tsx"
  "scaffold/web/src/component/l3/index.ts"
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

bash tests/contract/check_line_limits.sh

for path in "${required_dirs[@]}"; do
  test -d "$path" || {
    printf 'missing required directory: %s\n' "$path" >&2
    exit 1
  }
done

grep -q "Docker-first monorepo framework" README.md
grep -q "changing the Carbide framework repo" README.md
grep -q "changing a generated Carbide app" README.md
grep -q "https://carbide.ryangerardwilson.com/for/agents" README.md
grep -q "source of truth for app agents" README.md
grep -q "source of truth for framework agents" README.md
grep -q "framework-agent entrypoint and routing layer" README.md
grep -q "Its checked-in source is \`docs/app/web/site/for/agents/index.md\`" README.md
grep -q "## Goals" README.md
grep -q "## Non-Goals" README.md
grep -q "## Source Of Truth" README.md
grep -q "Checked-in app-agent contract source:" README.md
grep -q "docs/app/web/site/for/agents/index.md" README.md
grep -q "There is no separate internal docs tree under \`docs/engineering/\`." README.md
grep -q "Human docs page sources:" README.md
grep -q "Executable framework contract:" README.md
grep -q "## Task Router" README.md
grep -q "CLI parsing, health, audit, deploy, logs, upgrade:" README.md
grep -q "Public docs content:" README.md
grep -q "Framework regressions:" README.md
grep -q "## Current App Laws" README.md
grep -q "## Current Starter Taste" README.md
grep -q "Product-owned palettes are allowed." README.md
grep -q "If a framework change alters app-facing commands" README.md
grep -q "docs/app/web/site/for/agents/index.md" README.md
grep -q "## Verification" README.md
grep -q "bash tests/contract/check_line_limits.sh" README.md
grep -q "carbide health framework" README.md
grep -q "## Docs Website" README.md
grep -q "docs/app/web/site/" README.md
grep -q "docs/app/" README.md
grep -q "The docs app does not carry its own \`AGENTS.md\` or \`README.md\`." README.md
grep -q "The public app-agent contract lives at" README.md
grep -q "The checked-in source for that contract is" README.md
grep -q "docs/app/web/site/for/agents/index.md" README.md
grep -q "docs/app/deploy/prod.sh" README.md
grep -q "CARBIDE_DOCS_DEPLOY_SSH" README.md
grep -q "CARBIDE_DOCS_POSTGRES_PASSWORD" README.md
grep -q "carbide deploy prod" README.md
grep -q "tests/smoke/docs_for_agents_http.sh" README.md
! grep -q "Postgres-backed queues" README.md
! grep -q "## Roadmap" README.md
! grep -q "### Phase" README.md
! grep -q "## Core Commands" README.md
! grep -q "## Generated App Layout" README.md
! grep -q "## Runtime And Deploy" README.md
! test -f AGENTS.md
! grep -q "command_format" cli/bin/carbide
! grep -q "carbide format" cli/bin/carbide
grep -q "module github.com/ryangerardwilson/carbide/cli" cli/go.mod
grep -q ".cli/" .gitignore
! test -f go.mod
! test -d src
! test -d examples
! test -d infra
test "$(find . -mindepth 1 -maxdepth 1 -type d ! -name '.git' ! -name '.cli' ! -name '.bin' -printf '%f\n' | sort | tr '\n' ' ')" = "cli docs scaffold tests "
! test -d agents.d
! test -d docs/engineering
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
test "$(find scaffold -mindepth 1 -maxdepth 1 -type d ! -name .carbide -printf '%f\n' | sort | tr '\n' ' ')" = "api db web "
! test -d scaffold/agents.d
! test -d docs/app/agents.d
! test -f scaffold/AGENTS.md
! test -f scaffold/README.md
! test -f scaffold/Dockerfile
! test -f scaffold/go.mod
! test -f scaffold/go.sum
repo_search "oo_______oo_______oo" cli/internal/cli >/dev/null
grep -q "package main" cli/cmd/carbide/main.go
grep -q "package cli" cli/internal/cli/cli.go
repo_search "commandHealth\\(\\)" cli/internal/cli >/dev/null
repo_search "commandHealthEnv" cli/internal/cli >/dev/null
repo_search "commandHealthRuntime" cli/internal/cli >/dev/null
repo_search "commandHealthFramework" cli/internal/cli >/dev/null
repo_search "projectHealthResults" cli/internal/cli >/dev/null
repo_search "commandDeploy\\(" cli/internal/cli >/dev/null
repo_search "commandStatusJSON" cli/internal/cli >/dev/null
repo_search "commandURLs" cli/internal/cli >/dev/null
repo_search "projectDisplayName" cli/internal/cli >/dev/null
! git grep -n -e 'S[e]alion' -e 's[e]alion' -e 'S[E]ALION' -- .
repo_search "composeUpDetached" cli/internal/cli >/dev/null
repo_search "runDevStreams" cli/internal/cli >/dev/null
repo_search "\-\-quiet-build" cli/internal/cli >/dev/null
repo_search "Carbide dev" cli/internal/cli >/dev/null
grep -q "release binary unavailable and Go is required for source build fallback" cli/install.sh
grep -q "CARBIDE_VERSION" cli/install.sh
grep -q "CARBIDE_CHANNEL" cli/install.sh
grep -q "carbide_\${platform}.tar.gz" cli/install.sh
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
! grep -q "^\\[deploy\\]" scaffold/carbide.toml
grep -q ".carbide/" scaffold/.gitignore
grep -q ".env" scaffold/.gitignore
grep -q "web/node_modules/" scaffold/.gitignore
grep -q "web/public/" scaffold/.gitignore
grep -q "web/src/tailwind.css" scaffold/.gitignore
grep -q "POSTGRES_PASSWORD" scaffold/.env.example
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
grep -q "path: ./web/tsconfig.json" scaffold/docker-compose.yml
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
! grep -Eq '"(react|react-dom|tailwindcss|@tailwindcss/cli|typescript|@types/bun|@types/react|@types/react-dom)": "[~^<>*xX]|"(react|react-dom|tailwindcss|@tailwindcss/cli|typescript|@types/bun|@types/react|@types/react-dom)": "latest' scaffold/web/package.json
grep -q "Bun.serve" scaffold/web/src/server.ts
grep -q "browser entrypoint" scaffold/web/src/server.ts
grep -q "listening inside container" scaffold/web/src/server.ts
grep -q "proxying /api and /health to api service" scaffold/web/src/server.ts
grep -q "publicRoot" scaffold/web/src/server.ts
grep -q "Cache-Control" scaffold/web/src/server.ts
grep -q "public, max-age=31536000, immutable" scaffold/web/src/server.ts
grep -q "return 'no-store'" scaffold/web/src/server.ts
grep -q '"assets:build"' scaffold/web/package.json
grep -q '"typecheck": "tsc --noEmit"' scaffold/web/package.json
grep -q '"typescript": "6.0.3"' scaffold/web/package.json
grep -q '"@types/bun": "1.3.14"' scaffold/web/package.json
grep -q '"@types/react": "19.2.17"' scaffold/web/package.json
grep -q '"@types/react-dom": "19.2.3"' scaffold/web/package.json
grep -F -q "assets/[name]-[hash].[ext]" scaffold/web/package.json
grep -q '"strict": true' scaffold/web/tsconfig.json
grep -q '"jsx": "react-jsx"' scaffold/web/tsconfig.json
grep -F -q '"types": ["bun-types"]' scaffold/web/tsconfig.json
grep -q "bun run typecheck" scaffold/web/Dockerfile
grep -q "bun run assets:build" scaffold/web/Dockerfile
grep -q "asset-manifest.json" scaffold/web/src/write-index.ts
grep -F -q '/assets/${scripts[0]}' scaffold/web/src/write-index.ts
! grep -q "Bun frontend listening on http://localhost" scaffold/web/src/server.ts
grep -q '@import "tailwindcss";' scaffold/web/src/styles.css
grep -F -q '@source "./component/**/*.tsx";' scaffold/web/src/styles.css
grep -F -q '@source "./lib/**/*.ts";' scaffold/web/src/styles.css
grep -F -q '@source "./main.tsx";' scaffold/web/src/styles.css
grep -F -q '@source "./server.ts";' scaffold/web/src/styles.css
grep -q "@custom-variant dark" scaffold/web/src/styles.css
grep -q "\\[data-theme=\"dark\"\\]" scaffold/web/src/styles.css
! grep -q "html {" scaffold/web/src/styles.css
! grep -q "body {" scaffold/web/src/styles.css
! grep -q "font-size:" scaffold/web/src/styles.css
! grep -q "line-height:" scaffold/web/src/styles.css
! grep -q "min-width:" scaffold/web/src/styles.css
! grep -q "margin:" scaffold/web/src/styles.css
! grep -q "padding:" scaffold/web/src/styles.css
! grep -q "::-webkit-scrollbar" scaffold/web/src/styles.css
! grep -q "scrollbar-color:" scaffold/web/src/styles.css
! grep -q "scrollbar-width:" scaffold/web/src/styles.css
! grep -q "@theme" scaffold/web/src/styles.css
! grep -q -- "--carbide-" scaffold/web/src/styles.css
! grep -q "carbide-" scaffold/web/src/component/l1/tokens.ts
grep -q "const scrollbar =" scaffold/web/src/component/l1/tokens.ts
grep -F -q "[scrollbar-width:thin]" scaffold/web/src/component/l1/tokens.ts
grep -F -q "dark:[scrollbar-color:rgb(82_82_82)_transparent]" scaffold/web/src/component/l1/tokens.ts
grep -q "bg-white text-neutral-950 dark:bg-black dark:text-neutral-50" scaffold/web/src/component/l1/tokens.ts
! grep -Eq "#0f766e|#115e59|#2dd4bf|#5eead4|#16433c|#0f302c|#16211b|#edf5ef|#ecfdf5|#166534" scaffold/web/src/styles.css
! grep -q "from-carbide-action via-carbide-hero-via" scaffold/web/src/component/l1/tokens.ts
! grep -q "theme.css" scaffold/web/src/styles.css
! grep -Eq '^[[:space:]]*\.[A-Za-z_-]' scaffold/web/src/styles.css
! grep -Eq '^[[:space:]]*#[A-Za-z_-]' scaffold/web/src/styles.css
! grep -Eq '@theme|@apply|@layer|@keyframes|@media|@container|@plugin|@config' scaffold/web/src/styles.css
repo_search "scaffoldTailwindInputFindings" cli/internal/cli >/dev/null
repo_search "scaffold Tailwind input contract" cli/internal/cli >/dev/null
grep -F -q "text-2xl/8 sm:text-3xl/9" scaffold/web/src/component/l1/Text.tsx
grep -F -q "min-h-8 rounded-md border px-2 py-1 text-sm/6" scaffold/web/src/component/l1/Field.tsx
grep -F -q "md: 'min-h-8 px-3 text-xs'" scaffold/web/src/component/l1/Button.tsx
grep -F -q "gap-3 border-l px-4 py-5" scaffold/web/src/component/l2/AuthForm.tsx
grep -F -q "w-full max-w-sm justify-self-center gap-3" scaffold/web/src/component/l2/AuthForm.tsx
grep -F -q "lg:grid-cols-[216px_minmax(0,1fr)]" scaffold/web/src/component/l2/Layouts.tsx
grep -F -q "px-3 py-4 sm:px-5 lg:py-5" scaffold/web/src/component/l2/Layouts.tsx
grep -q "ui.scrollbar" scaffold/web/src/component/l2/Layouts.tsx
! grep -R -E "text-7xl|text-5xl|py-24|lg:py-12|min-h-12 rounded-md border|min-h-10 rounded-md border|lg:grid-cols-\[280px|lg:grid-cols-\[240px|gap-6|p-6|font-extrabold" scaffold/web/src/component >/dev/null
grep -q '/api/${mode}' scaffold/web/src/main.tsx
grep -q "carbide.theme" scaffold/web/src/main.tsx
grep -q "useThemeMode" scaffold/web/src/main.tsx
grep -q "prefers-color-scheme: dark" scaffold/web/index.html
grep -q "dataset.theme" scaffold/web/index.html
grep -q "ThemeToggle" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "data-resolved-theme" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "data-theme-mode" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "SunIcon" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "MoonIcon" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "Switch to light theme" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "Switch to dark theme" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "size-8 rounded-full border" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "onClick={() => onMode?.(nextMode)}" scaffold/web/src/component/l1/ThemeToggle.tsx
! grep -q "<select" scaffold/web/src/component/l1/ThemeToggle.tsx
! grep -q "appearance-none" scaffold/web/src/component/l1/ThemeToggle.tsx
! grep -q "border-x-4 border-t-4 border-x-transparent" scaffold/web/src/component/l1/ThemeToggle.tsx
! grep -q "aria-pressed" scaffold/web/src/component/l1/ThemeToggle.tsx
! grep -q "role=\"group\"" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "./component/l3" scaffold/web/src/main.tsx
grep -q "AuthView" scaffold/web/src/main.tsx
grep -q "DashboardView" scaffold/web/src/main.tsx
grep -q "LoadingView" scaffold/web/src/main.tsx
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
grep -q "export function Button" scaffold/web/src/component/l1/Button.tsx
grep -q "export function Field" scaffold/web/src/component/l1/Field.tsx
grep -q "export function Panel" scaffold/web/src/component/l1/Surface.tsx
grep -q "export function ThemeToggle" scaffold/web/src/component/l1/ThemeToggle.tsx
grep -q "ThemeToggle" scaffold/web/src/component/l1/index.ts
grep -q "export const ui" scaffold/web/src/component/l1/tokens.ts
! test -f scaffold/web/src/component/l1/theme.css
! grep -R "cb-" scaffold/web/src >/dev/null
! grep -R -- "--cb-" scaffold/web/src >/dev/null
grep -q "ui.action" scaffold/web/src/component/l1/Button.tsx
grep -q "ui.input" scaffold/web/src/component/l1/Field.tsx
grep -q "buttonClassLayers" scaffold/web/src/component/l1/Button.tsx
grep -q "fieldClassLayers" scaffold/web/src/component/l1/Field.tsx
grep -q "fieldHintClassLayers" scaffold/web/src/component/l1/Field.tsx
grep -q "fieldErrorClassLayers" scaffold/web/src/component/l1/Field.tsx
grep -q "inputClassLayers" scaffold/web/src/component/l1/Field.tsx
grep -q "panelClassLayers" scaffold/web/src/component/l1/Surface.tsx
grep -q "dividerClassLayers" scaffold/web/src/component/l1/Surface.tsx
grep -q "badgeClassLayers" scaffold/web/src/component/l1/Surface.tsx
grep -q "metricClassLayers" scaffold/web/src/component/l1/Surface.tsx
grep -q "eyebrowClassLayers" scaffold/web/src/component/l1/Text.tsx
grep -q "headingClassLayers" scaffold/web/src/component/l1/Text.tsx
grep -q "mutedClassLayers" scaffold/web/src/component/l1/Text.tsx
grep -q "codeClassLayers" scaffold/web/src/component/l1/Text.tsx
grep -q "formClassLayers" scaffold/web/src/component/l2/AuthForm.tsx
grep -q "formStackClassLayers" scaffold/web/src/component/l2/AuthForm.tsx
grep -q "errorClassLayers" scaffold/web/src/component/l2/AuthForm.tsx
grep -q "modeButtonClassLayers" scaffold/web/src/component/l2/AuthForm.tsx
grep -q "landingClassLayers" scaffold/web/src/component/l2/Layouts.tsx
grep -q "dashboardClassLayers" scaffold/web/src/component/l2/Layouts.tsx
grep -q "screenClassLayers" scaffold/web/src/component/l3/DashboardView.tsx
grep -q "loadingClassLayers" scaffold/web/src/component/l3/LoadingView.tsx
grep -q "ui.focus" scaffold/web/src/component/l1/Field.tsx
grep -q "ui.focus" scaffold/web/src/component/l2/AuthForm.tsx
grep -q "ui.focus" scaffold/web/src/component/l2/Layouts.tsx
! grep -R "text-\\[" scaffold/web/src/component >/dev/null
grep -q "export function DashboardLayout" scaffold/web/src/component/l2/Layouts.tsx
grep -q "lg:grid-cols-\\[216px_minmax(0,1fr)\\]" scaffold/web/src/component/l2/Layouts.tsx
grep -q "aria-label=\"Dashboard\"" scaffold/web/src/component/l2/Layouts.tsx
grep -q "aria-current" scaffold/web/src/component/l2/Layouts.tsx
grep -q "navItems" scaffold/web/src/component/l2/Layouts.tsx
grep -q "export function LandingPageLayout" scaffold/web/src/component/l2/Layouts.tsx
grep -q "export function AuthView" scaffold/web/src/component/l3/AuthView.tsx
grep -q "export function DashboardView" scaffold/web/src/component/l3/DashboardView.tsx
grep -q "dashboardNav" scaffold/web/src/component/l3/DashboardView.tsx
grep -q "WorkspaceOverview" scaffold/web/src/component/l3/DashboardView.tsx
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
repo_search "composeFilePath" cli/internal/cli >/dev/null
repo_search "COMPOSE_FILE" cli/internal/cli >/dev/null
repo_search 'compose.supports\("--watch"\)' cli/internal/cli >/dev/null
repo_search "newRenderer" cli/internal/cli >/dev/null
repo_search "func \\(r renderer\\) Table" cli/internal/cli >/dev/null
repo_search "runDevStreams" cli/internal/cli >/dev/null
repo_search "commandStatus" cli/internal/cli >/dev/null
repo_search "commandStopDev" cli/internal/cli >/dev/null
repo_search "RunServiceProgress" cli/internal/cli >/dev/null
repo_search "RunServiceStopProgress" cli/internal/cli >/dev/null
repo_search "serviceProgressFrameWidth" cli/internal/cli >/dev/null
repo_search "serviceProgressFrame" cli/internal/cli >/dev/null
repo_search "terminalColumns" cli/internal/cli >/dev/null
repo_search "composeServiceStatuses" cli/internal/cli >/dev/null
repo_search "composeServiceSnapshots" cli/internal/cli >/dev/null
repo_search "composePublishedPorts" cli/internal/cli >/dev/null
repo_search "composeInternalPorts" cli/internal/cli >/dev/null
repo_search "streamLogOutput" cli/internal/cli >/dev/null
repo_search "parseComposeLogLine" cli/internal/cli >/dev/null
repo_search "composeLogsArgs" cli/internal/cli >/dev/null
repo_search "openDevLogSink" cli/internal/cli >/dev/null
repo_search "openAppendDevLogSink" cli/internal/cli >/dev/null
repo_search "commandLogs" cli/internal/cli >/dev/null
repo_search "commandFollowLogs" cli/internal/cli >/dev/null
repo_search ".carbide/log/dev.jsonl" cli/internal/cli >/dev/null
repo_search "carbide follow logs" cli/internal/cli >/dev/null
repo_search "carbide status" cli/internal/cli >/dev/null
repo_search "health framework" cli/internal/cli >/dev/null
! repo_search "carbide logs follow" cli/internal/cli >/dev/null
! grep -R -q 'outputRow{"login"' cli/internal/cli
! grep -R -q 'outputRow{"mode"' cli/internal/cli

grep -q "$domain" docs/app/web/site/index.html
grep -q "Bun frontend" docs/app/web/site/index.html
grep -q "Create Your First App" docs/app/web/site/index.html
grep -q 'href="/#start"' docs/app/web/site/index.html
grep -q "<h2>Start</h2>" docs/app/web/site/index.html
grep -q "Paste this into your AI agent" docs/app/web/site/index.html
grep -q "https://carbide.ryangerardwilson.com/for/agents" docs/app/web/site/index.html
grep -q "Treat that Markdown as the source of truth" docs/app/web/site/index.html
! grep -q "Guiding Your Agents to Get Started" docs/app/web/site/index.html
! grep -q "curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/carbide/main/cli/install.sh | bash" docs/app/web/site/index.html
! grep -q "carbide new demo" docs/app/web/site/index.html
! grep -q "The prompt below tells the agent" docs/app/web/site/index.html
! grep -q "Treat the returned Markdown as the source of truth" docs/app/web/site/index.html
! grep -R 'href="/#for-agents"' docs/app/web/site/*.html >/dev/null
grep -q 'href="/version-policy"' docs/app/web/site/index.html
grep -q 'href="/create-your-first-app"' docs/app/web/site/index.html
grep -q 'href="/frontend-starter-contract"' docs/app/web/site/index.html
grep -q 'href="/deployment"' docs/app/web/site/index.html
grep -q 'href="/ci-cd-regression-tests"' docs/app/web/site/index.html
! grep -R -E 'href="[^"]+\.html' docs/app/web/site >/dev/null
grep -q "function isHomePage()" docs/app/web/site/assets/intro.js
grep -q 'pathname === "/"' docs/app/web/site/assets/intro.js
grep -q 'pathname === "/index.html"' docs/app/web/site/assets/intro.js
grep -q "intro=1" docs/app/web/site/assets/intro.js
grep -q "prefers-reduced-motion" docs/app/web/site/assets/intro.js
! grep -q "Carbide docs" docs/app/web/site/assets/intro.js
! grep -q "Carbide Docs" docs/app/web/site/assets/intro.js
! grep -q "chomper" docs/app/web/site/assets/intro.js
! grep -q "pellets" docs/app/web/site/assets/intro.js
! grep -q "clipPath" docs/app/web/site/assets/intro.js
! grep -q "sessionStorage" docs/app/web/site/assets/intro.js
! grep -q "storageKey" docs/app/web/site/assets/intro.js
! grep -q "docs-intro-skip" docs/app/web/site/assets/intro.js
! grep -q ">Skip<" docs/app/web/site/assets/intro.js
! grep -q "skipIntro" docs/app/web/site/assets/intro.js
grep -q "canonicalDocsPath" docs/app/web/src/server.ts
grep -q '"/initial-user-experience": "/create-your-first-app"' docs/app/web/src/server.ts
grep -q 'pathname === "/index.html"' docs/app/web/src/server.ts
grep -q 'requestPath === "/for/agents"' docs/app/web/src/server.ts
grep -q '"/for/agents/index.md"' docs/app/web/src/server.ts
grep -q 'pathname.endsWith(".html")' docs/app/web/src/server.ts
grep -q 'text/markdown; charset=utf-8' docs/app/web/src/server.ts
grep -q "status: 308" docs/app/web/src/server.ts
grep -q 'location: `${pathname}${target.search}`' docs/app/web/src/server.ts
grep -q 'docsResponseHeaders' docs/app/web/src/server.ts
grep -q 'rewriteDocsHtml' docs/app/web/src/server.ts
grep -q 'createHash' docs/app/web/src/server.ts
grep -q 'cacheBustHtml' docs/app/web/src/server.ts
grep -q 'versionedAssetPath' docs/app/web/src/server.ts
grep -F -q '?v=${hash}' docs/app/web/src/server.ts
grep -q 'assets/intro.js' docs/app/web/src/server.ts
grep -q 'assets/styles.css' docs/app/web/src/server.ts
grep -F -q 'output.replaceAll(`"/${assetPath}"`, `"/${versionedPath}"`)' docs/app/web/src/server.ts
grep -q 'join(import.meta.dir, "..", "site")' docs/app/web/src/server.ts
grep -q 'return "no-cache"' docs/app/web/src/server.ts
grep -q 'return "no-store"' docs/app/web/src/server.ts
grep -q '@import "tailwindcss";' docs/app/web/src/styles.css
grep -F -q '@source "./component/**/*.ts";' docs/app/web/src/styles.css
grep -F -q '@source "./lib/**/*.ts";' docs/app/web/src/styles.css
grep -F -q '@source "./server.ts";' docs/app/web/src/styles.css
grep -q "@custom-variant dark" docs/app/web/src/styles.css
! grep -q "html {" docs/app/web/src/styles.css
! grep -q "body {" docs/app/web/src/styles.css
! grep -q "font-size:" docs/app/web/src/styles.css
! grep -q "line-height:" docs/app/web/src/styles.css
! grep -q "min-width:" docs/app/web/src/styles.css
! grep -q "margin:" docs/app/web/src/styles.css
! grep -q "padding:" docs/app/web/src/styles.css
! grep -q "::-webkit-scrollbar" docs/app/web/src/styles.css
! grep -q "scrollbar-color:" docs/app/web/src/styles.css
! grep -q "scrollbar-width:" docs/app/web/src/styles.css
! grep -q "@theme" docs/app/web/src/styles.css
! grep -q -- "--carbide-" docs/app/web/src/styles.css
repo_search "docsTailwindInputFindings" cli/internal/cli >/dev/null
repo_search "docsGeneratedTailwindFindings" cli/internal/cli >/dev/null
grep -q "path: ./web/site" docs/app/docker-compose.yml
grep -q "path: ./api" docs/app/docker-compose.yml
grep -q "path: ./db/migration" docs/app/docker-compose.yml
! grep -q "docs-intro-skip" docs/app/web/src/styles.css
! grep -q "docs-intro" docs/app/web/src/styles.css
grep -q '"assets:build"' docs/app/web/package.json
grep -q '"docs:styles"' docs/app/web/package.json
grep -q '"build"' docs/app/web/package.json
grep -q '"typecheck": "tsc --noEmit"' docs/app/web/package.json
grep -q "tailwindcss" docs/app/web/src/build-styles.ts
grep -q '"@tailwindcss/cli": "4.3.2"' docs/app/web/package.json
grep -q '"tailwindcss": "4.3.2"' docs/app/web/package.json
grep -q '"typescript": "6.0.3"' docs/app/web/package.json
grep -q '"@types/bun": "1.3.14"' docs/app/web/package.json
! grep -q '"react":' docs/app/web/package.json
! grep -q '"react-dom":' docs/app/web/package.json
! grep -q '"@types/react":' docs/app/web/package.json
! grep -q '"@types/react-dom":' docs/app/web/package.json
grep -q '"strict": true' docs/app/web/tsconfig.json
grep -F -q '"types": ["bun-types"]' docs/app/web/tsconfig.json
grep -F -q '"include": ["src/**/*.ts"]' docs/app/web/tsconfig.json
grep -q "bun run typecheck" docs/app/web/Dockerfile
grep -q "bun run assets:build" docs/app/web/Dockerfile
grep -q "COPY app/web/site ./site" docs/app/web/Dockerfile
! test -f docs/app/web/index.html
! test -f docs/app/web/src/main.tsx
! test -f docs/app/web/src/write-index.ts
grep -q "docsClassLayers" docs/app/web/src/component/l1/tokens.ts
grep -q "scrollbar" docs/app/web/src/component/l1/tokens.ts
grep -F -q "[scrollbar-width:thin]" docs/app/web/src/component/l1/tokens.ts
grep -F -q "dark:[scrollbar-color:rgb(250_204_21)_transparent]" docs/app/web/src/component/l1/tokens.ts
grep -q "docsScrollbarClass" docs/app/web/src/component/l2/DocsChrome.ts
grep -q "docsStaticClassMap" docs/app/web/src/component/l2/DocsChrome.ts
grep -q "rewriteDocsClasses" docs/app/web/src/component/l2/DocsChrome.ts
grep -F -q "[&_pre]:[scrollbar-width:thin]" docs/app/web/src/component/l2/DocsChrome.ts
grep -F -q "[&_pre+p]:mt-[18px]" docs/app/web/src/component/l2/DocsChrome.ts
grep -F -q "max-[860px]:mt-[34px]" docs/app/web/src/component/l2/DocsChrome.ts
grep -q "docsChromeClassLayers" docs/app/web/src/component/l2/DocsChrome.ts
grep -q "docsWebContract" docs/app/web/src/component/l3/DocsSite.ts
grep -q "rewriteDocsHtml" docs/app/web/src/component/l3/DocsSite.ts
grep -F -q "[scrollbar-width:thin]" docs/app/web/src/component/l3/DocsSite.ts
repo_search "fileLineCount" cli/internal/cli >/dev/null
! test -f docs/app/AGENTS.md
! test -f docs/app/README.md
grep -q "The docs app does not carry its own \`AGENTS.md\` or \`README.md\`." README.md
grep -q "black and yellow" README.md
grep -q "audits should preserve that" README.md
grep -q "preserve the app's existing" docs/app/web/site/for/agents/index.md
grep -q "A branded black/yellow app may stay black/yellow" docs/app/web/site/for/agents/index.md
grep -q '\[deploy.targets.prod\]' docs/app/carbide.toml
grep -q 'script = "./deploy/prod.sh"' docs/app/carbide.toml
grep -q "CARBIDE_DOCS_POSTGRES_PASSWORD" docs/app/deploy/prod.sh
grep -q "CARBIDE_DOCS_DEPLOY_SSH" docs/app/deploy/prod.sh
grep -q 'compose_project_name="${CARBIDE_DOCS_COMPOSE_PROJECT_NAME:-carbide-docs}"' docs/app/deploy/prod.sh
grep -q 'legacy_project_names="${CARBIDE_DOCS_LEGACY_PROJECT_NAMES:-app 1}"' docs/app/deploy/prod.sh
grep -q 'compose_cmd() {' docs/app/deploy/prod.sh
grep -q -- '--project-directory app' docs/app/deploy/prod.sh
grep -q "bg-amber-50" docs/app/web/src/component/l1/tokens.ts
grep -q "dark:text-neutral-50" docs/app/web/src/component/l1/tokens.ts
grep -q "bg-yellow-400" docs/app/web/src/component/l2/DocsChrome.ts
grep -q "text-yellow-300" docs/app/web/src/component/l2/DocsChrome.ts
grep -F -q "dark:[&_p]:text-neutral-300" docs/app/web/src/component/l2/DocsChrome.ts
grep -F -q "dark:[&_pre]:text-neutral-50" docs/app/web/src/component/l2/DocsChrome.ts
grep -q "Bun frontend, Go API backend, Postgres database" docs/app/web/site/frontend-starter-contract.html
grep -q "Tailwind is required" docs/app/web/site/frontend-starter-contract.html
grep -q "Tailwind Plus and Catalyst" docs/app/web/site/frontend-starter-contract.html
grep -q "Application UI patterns" docs/app/web/site/frontend-starter-contract.html
grep -q "production-ready, fully responsive, accessible, and easy to customize" docs/app/web/site/frontend-starter-contract.html
grep -q "carbide health.*rejects global" docs/app/web/site/frontend-starter-contract.html
grep -q "built-in scrollbar styling" docs/app/web/site/frontend-starter-contract.html
grep -q "custom selectors" docs/app/web/site/frontend-starter-contract.html
grep -q "component styling belongs in Tailwind utility classes" docs/app/web/site/frontend-starter-contract.html
grep -q "web/src/product.css" docs/app/web/site/frontend-starter-contract.html
grep -q "local app docs if they exist" docs/app/web/site/frontend-starter-contract.html
grep -q "ThemeToggle.tsx" docs/app/web/site/frontend-starter-contract.html
grep -q "localStorage" docs/app/web/site/frontend-starter-contract.html
grep -q "matchMedia" docs/app/web/site/frontend-starter-contract.html
grep -q "dataset.themeMode" docs/app/web/site/frontend-starter-contract.html
! repo_search "de-sci|public domain behavior" docs cli/internal/cli >/dev/null
! repo_search "PROJECT\\.md" README.md scaffold docs cli/internal/cli tests >/dev/null
grep -q "scrollbar-width:thin" docs/app/web/site/assets/styles.css
grep -q "scrollbar-color:#d97706 transparent" docs/app/web/site/assets/styles.css
grep -q "scrollbar-color:#facc15 transparent" docs/app/web/site/assets/styles.css
grep -q "carbide follow logs" docs/app/web/site/create-your-first-app.html
grep -q "carbide clean dev" docs/app/web/site/create-your-first-app.html
grep -q "carbide status" docs/app/web/site/create-your-first-app.html
grep -q "carbide audit" docs/app/web/site/create-your-first-app.html
grep -q "carbide resolve" docs/app/web/site/create-your-first-app.html
grep -q "carbide fix" docs/app/web/site/create-your-first-app.html
grep -q "carbide audit resolve fix" docs/app/web/site/create-your-first-app.html
grep -q "carbide health runtime" docs/app/web/site/create-your-first-app.html
grep -q "Troubleshooting" docs/app/web/site/create-your-first-app.html
grep -q "Install, create, run, register" docs/app/web/site/create-your-first-app.html
test ! -f docs/app/web/site/for/agents.html
grep -q "# Carbide for Agents" docs/app/web/site/for/agents/index.md
grep -q "source of truth for AI agents" docs/app/web/site/for/agents/index.md
grep -q "https://raw.githubusercontent.com/ryangerardwilson/carbide/main/docs/app/web/site/for/agents/index.md" docs/app/web/site/for/agents/index.md
grep -q "## Source Precedence" docs/app/web/site/for/agents/index.md
grep -q "framework repo \`README.md\`" docs/app/web/site/for/agents/index.md
grep -q "## Identify The Current State" docs/app/web/site/for/agents/index.md
grep -q "framework-agent entrypoint" docs/app/web/site/for/agents/index.md
grep -q "There is no separate internal \`docs/engineering/\` tree." docs/app/web/site/for/agents/index.md
grep -q "## Prerequisites" docs/app/web/site/for/agents/index.md
grep -q "## Create A New App" docs/app/web/site/for/agents/index.md
grep -q "## Development Loop" docs/app/web/site/for/agents/index.md
grep -q "## Laws" docs/app/web/site/for/agents/index.md
grep -q "### Law 1. One App Repo" docs/app/web/site/for/agents/index.md
grep -q "### Law 8. Checked Files Stay Under 1000 Lines" docs/app/web/site/for/agents/index.md
grep -q "### Law 7. Secrets Are Never Printed" docs/app/web/site/for/agents/index.md
grep -q 'Use `Law 1` through `Law 8`' docs/app/web/site/for/agents/index.md
grep -q "## Ownership Rule" docs/app/web/site/for/agents/index.md
grep -q "## Current Taste" docs/app/web/site/for/agents/index.md
grep -q "### Taste 1. Starter Stack" docs/app/web/site/for/agents/index.md
grep -q "### Taste 6. CLI And Audit Reporting" docs/app/web/site/for/agents/index.md
grep -q 'Use `Taste 1` through `Taste 6`' docs/app/web/site/for/agents/index.md
grep -q "## Frontend Contract" docs/app/web/site/for/agents/index.md
grep -q "web/src/component/l1" docs/app/web/site/for/agents/index.md
grep -q "web/src/component/l2" docs/app/web/site/for/agents/index.md
grep -q "web/src/component/l3" docs/app/web/site/for/agents/index.md
grep -F -q '`l1`: structure and layout' docs/app/web/site/for/agents/index.md
grep -F -q '`l2`: geometry, spacing, borders, radii, and type scale' docs/app/web/site/for/agents/index.md
grep -F -q '`l3`: theme, color, state, motion, and interaction' docs/app/web/site/for/agents/index.md
grep -q "Tailwind Plus / Catalyst style" docs/app/web/site/for/agents/index.md
grep -q "Application UI patterns" docs/app/web/site/for/agents/index.md
grep -q "production-ready, fully responsive, accessible, and easy" docs/app/web/site/for/agents/index.md
grep -q "ship all normal states for interactive components" docs/app/web/site/for/agents/index.md
grep -q "focus-visible, active, disabled, loading, empty, and error" docs/app/web/site/for/agents/index.md
grep -q "ThemeToggle" docs/app/web/site/for/agents/index.md
grep -q "localStorage" docs/app/web/site/for/agents/index.md
grep -q "matchMedia" docs/app/web/site/for/agents/index.md
grep -q "dataset.theme" docs/app/web/site/for/agents/index.md
grep -q "dataset.themeMode" docs/app/web/site/for/agents/index.md
grep -q "## Environment And Secrets" docs/app/web/site/for/agents/index.md
grep -q "## Deployment" docs/app/web/site/for/agents/index.md
grep -q "## Audits" docs/app/web/site/for/agents/index.md
grep -q "## Verification" docs/app/web/site/for/agents/index.md
grep -q "## Recovery" docs/app/web/site/for/agents/index.md
grep -q "## Agent Behavior" docs/app/web/site/for/agents/index.md
grep -q "README.md" docs/app/web/site/for/agents/index.md
grep -q "AGENTS.md" docs/app/web/site/for/agents/index.md
grep -q "carbide.toml" docs/app/web/site/for/agents/index.md
grep -q "docker-compose.yml" docs/app/web/site/for/agents/index.md
grep -q "web/" docs/app/web/site/for/agents/index.md
grep -q "api/" docs/app/web/site/for/agents/index.md
grep -q "db/" docs/app/web/site/for/agents/index.md
grep -q "curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/carbide/main/cli/install.sh | bash" docs/app/web/site/for/agents/index.md
grep -q "carbide new demo" docs/app/web/site/for/agents/index.md
grep -q 'carbide new "My Carbide App"' docs/app/web/site/for/agents/index.md
grep -q "carbide init" docs/app/web/site/for/agents/index.md
grep -q "carbide run dev" docs/app/web/site/for/agents/index.md
grep -q "carbide health" docs/app/web/site/for/agents/index.md
grep -q "carbide status" docs/app/web/site/for/agents/index.md
grep -q "carbide urls json" docs/app/web/site/for/agents/index.md
grep -q "carbide status json" docs/app/web/site/for/agents/index.md
grep -q "carbide health json" docs/app/web/site/for/agents/index.md
grep -q "carbide follow logs" docs/app/web/site/for/agents/index.md
grep -q "carbide clean dev" docs/app/web/site/for/agents/index.md
grep -q "carbide health env" docs/app/web/site/for/agents/index.md
grep -q "carbide health env json" docs/app/web/site/for/agents/index.md
grep -q "carbide health runtime" docs/app/web/site/for/agents/index.md
grep -q "carbide health runtime json" docs/app/web/site/for/agents/index.md
grep -q "carbide stop dev" docs/app/web/site/for/agents/index.md
grep -q "carbide help" docs/app/web/site/for/agents/index.md
grep -q "carbide upgrade" docs/app/web/site/for/agents/index.md
grep -q "carbide deploy prod" docs/app/web/site/for/agents/index.md
grep -q '\[deploy.targets.prod\]' docs/app/web/site/for/agents/index.md
grep -q 'script = "./deploy/prod.sh"' docs/app/web/site/for/agents/index.md
grep -q "carbide audit" docs/app/web/site/for/agents/index.md
grep -q "carbide resolve" docs/app/web/site/for/agents/index.md
grep -q "carbide fix" docs/app/web/site/for/agents/index.md
grep -q "carbide audit resolve fix" docs/app/web/site/for/agents/index.md
grep -q ".audit/report/" docs/app/web/site/for/agents/index.md
grep -q ".audit/plan.md" docs/app/web/site/for/agents/index.md
grep -q ".audit/fix.md" docs/app/web/site/for/agents/index.md
grep -q "Carbide itself should never touch app code" docs/app/web/site/for/agents/index.md
grep -q "Codex edits made intentionally" docs/app/web/site/for/agents/index.md
grep -q "command-shaped JSON output" docs/app/web/site/for/agents/index.md
grep -q "^## Troubleshooting" docs/app/web/site/for/agents/index.md
grep -q 'Carbide does not scaffold `README.md`, `AGENTS.md`, or `agents.d/`' docs/app/web/site/for/agents/index.md
grep -q "Do not add seeded demo credentials" docs/app/web/site/for/agents/index.md
grep -q "Do not print secret values" docs/app/web/site/for/agents/index.md
grep -q "If the current directory already contains" docs/app/web/site/for/agents/index.md
! grep -q "Guiding Your Agents to Get Started" docs/app/web/site/for/agents/index.md
! grep -q "This page is for AI coding agents" docs/app/web/site/for/agents/index.md
! repo_search '<a class="nav-link" href="https://github.com/ryangerardwilson/carbide">Source Repository</a>' docs/app/web/site/*.html >/dev/null
test "$(repo_search_files '<a class="nav-link" href="https://github.com/ryangerardwilson/carbide" target="_blank" rel="noopener noreferrer">Source Repository</a>' docs/app/web/site/*.html | wc -l)" -eq 7
grep -q "Target Contract" docs/app/web/site/deployment.html
grep -q "Script Ownership" docs/app/web/site/deployment.html
grep -q "Script Environment" docs/app/web/site/deployment.html
grep -q "carbide deploy prod" docs/app/web/site/deployment.html
grep -q '\[deploy.targets.prod\]' docs/app/web/site/deployment.html
grep -q 'script = "./deploy/prod.sh"' docs/app/web/site/deployment.html
grep -q "CI/CD regression plan" docs/app/web/site/ci-cd-regression-tests.html
grep -q "carbide health framework" docs/app/web/site/ci-cd-regression-tests.html
grep -q "carbide audit" docs/app/web/site/ci-cd-regression-tests.html
grep -q "carbide resolve" docs/app/web/site/ci-cd-regression-tests.html
grep -q "carbide fix" docs/app/web/site/ci-cd-regression-tests.html
grep -q "carbide audit resolve fix" docs/app/web/site/version-policy.html
grep -q ".audit/plan.md" docs/app/web/site/version-policy.html
grep -q "Directory Structure" docs/app/web/site/repo-structure.html
grep -q "Generated App Layout" docs/app/web/site/repo-structure.html
grep -q 'carbide new "My Carbide App"' docs/app/web/site/repo-structure.html
grep -q "my-carbide-app/" docs/app/web/site/repo-structure.html
grep -q "README.md" docs/app/web/site/repo-structure.html
grep -q "does not scaffold" docs/app/web/site/repo-structure.html
grep -q "web/src/component/l1" docs/app/web/site/repo-structure.html
grep -q "web/src/component/l2" docs/app/web/site/repo-structure.html
grep -q "web/src/component/l3" docs/app/web/site/repo-structure.html
grep -q "scrollbar utility group" docs/app/web/site/repo-structure.html
grep -q "Global .*html.*body.* sizing" docs/app/web/site/repo-structure.html
grep -q "web, api, db" docs/app/web/site/repo-structure.html
grep -q "Generated apps do not include root" docs/app/web/site/repo-structure.html
! grep -q ".github/workflows" docs/app/web/site/repo-structure.html
! grep -q "cli/internal/cli" docs/app/web/site/repo-structure.html
! grep -q "docs/app/" docs/app/web/site/repo-structure.html
grep -q 'class="docs-topbar"' docs/app/web/site/index.html
grep -q 'class="docs-sidebar"' docs/app/web/site/index.html
grep -q 'class="docs-content"' docs/app/web/site/index.html
grep -q 'class="docs-toc"' docs/app/web/site/index.html
grep -q "Search docs" docs/app/web/site/index.html
grep -q "Version v0.2.0" docs/app/web/site/index.html
grep -q 'href="https://github.com/ryangerardwilson/carbide" target="_blank" rel="noopener noreferrer"' docs/app/web/site/index.html
grep -q "Prologue" docs/app/web/site/index.html
grep -q "Getting Started" docs/app/web/site/index.html
grep -q "Architecture" docs/app/web/site/index.html
grep -q "On this page" docs/app/web/site/index.html
! grep -E -q '\.(docs-layout|docs-sidebar|docs-toc|docs-topbar)' docs/app/web/site/assets/styles.css
! grep -q 'html{font-size:14px}' docs/app/web/site/assets/styles.css
! grep -q 'body{min-width:320px' docs/app/web/site/assets/styles.css
! grep -q 'body{margin:0;min-width:320px' docs/app/web/site/assets/styles.css
grep -F -q ".max-\\[860px\\]\\:grid-cols-1" docs/app/web/site/assets/styles.css

for page in docs/app/web/site/*.html; do
  grep -q 'class="docs-topbar"' "$page"
  grep -q 'class="docs-sidebar"' "$page"
  grep -q 'class="docs-content"' "$page"
  grep -q 'class="docs-toc"' "$page"
  grep -q "Search docs" "$page"
  grep -q "Version v0.2.0" "$page"
  grep -q 'href="https://github.com/ryangerardwilson/carbide" target="_blank" rel="noopener noreferrer"' "$page"
  grep -q "On this page" "$page"
done

printf 'repo contract ok\n'
