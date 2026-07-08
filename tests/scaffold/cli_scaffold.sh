#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." >/dev/null 2>&1 && pwd)"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

export CARBIDE_HOME="$repo_root"

"$repo_root/cli/bin/carbide" > "$tmp_dir/no-args.out"
grep -q "________________________oo_______oo_______oo_________" "$tmp_dir/no-args.out"
grep -q "Carbide 0.2.0" "$tmp_dir/no-args.out"
grep -q "Usage:" "$tmp_dir/no-args.out"
grep -q "Commands:" "$tmp_dir/no-args.out"
grep -q "carbide <command> \\[arguments\\]" "$tmp_dir/no-args.out"
grep -q "new <project-name>" "$tmp_dir/no-args.out"
grep -q "init" "$tmp_dir/no-args.out"
grep -q "help" "$tmp_dir/no-args.out"
! grep -q "Options:" "$tmp_dir/no-args.out"
! grep -q "Available commands:" "$tmp_dir/no-args.out"
! grep -q "run dev" "$tmp_dir/no-args.out"
! grep -q "status" "$tmp_dir/no-args.out"
! grep -q "stop dev" "$tmp_dir/no-args.out"
! grep -q "follow logs" "$tmp_dir/no-args.out"
! grep -q "upgrade" "$tmp_dir/no-args.out"
! grep -q "version" "$tmp_dir/no-args.out"
! grep -q "features:" "$tmp_dir/no-args.out"
! grep -q "raw.githubusercontent.com/ryangerardwilson/carbide" "$tmp_dir/no-args.out"

"$repo_root/cli/bin/carbide" help > "$tmp_dir/help.out"
awk 'length($0) > 79 { print "help line exceeds 79 chars: " $0; exit 1 }' "$tmp_dir/help.out"
grep -q "^Usage:$" "$tmp_dir/help.out"
grep -q "^  carbide <command> \\[arguments\\]$" "$tmp_dir/help.out"
grep -q "^Available commands:$" "$tmp_dir/help.out"
grep -q "^  clean dev " "$tmp_dir/help.out"
grep -q "^  deploy apply prod " "$tmp_dir/help.out"
grep -q "^  deploy check prod " "$tmp_dir/help.out"
grep -q "^  deploy check prod json " "$tmp_dir/help.out"
grep -q "^  deploy preview prod " "$tmp_dir/help.out"
grep -q "^  deploy preview prod json " "$tmp_dir/help.out"
grep -q "^  health " "$tmp_dir/help.out"
grep -q "^  health json " "$tmp_dir/help.out"
grep -q "^  health env " "$tmp_dir/help.out"
grep -q "^  health env json " "$tmp_dir/help.out"
grep -q "^  health framework " "$tmp_dir/help.out"
grep -q "^  health framework json " "$tmp_dir/help.out"
grep -q "^  health runtime " "$tmp_dir/help.out"
grep -q "^  health runtime json " "$tmp_dir/help.out"
grep -q "^  help " "$tmp_dir/help.out"
grep -q "^  init " "$tmp_dir/help.out"
grep -q "^  logs " "$tmp_dir/help.out"
grep -q "^  new <project-name> " "$tmp_dir/help.out"
grep -q "^  audit " "$tmp_dir/help.out"
grep -q "^  status " "$tmp_dir/help.out"
grep -q "^  status json " "$tmp_dir/help.out"
grep -q "^  upgrade " "$tmp_dir/help.out"
grep -q "^  urls " "$tmp_dir/help.out"
grep -q "^  urls json " "$tmp_dir/help.out"
grep -q "^  version " "$tmp_dir/help.out"
grep -q "^follow$" "$tmp_dir/help.out"
grep -q "^  follow logs " "$tmp_dir/help.out"
grep -q "^  follow logs service api " "$tmp_dir/help.out"
grep -q "^logs$" "$tmp_dir/help.out"
grep -q "^  logs containing \"/api/login\" json " "$tmp_dir/help.out"
grep -q "^run$" "$tmp_dir/help.out"
grep -q "^  run dev " "$tmp_dir/help.out"
grep -q "^stop$" "$tmp_dir/help.out"
grep -q "^  stop dev " "$tmp_dir/help.out"
! grep -q "^area" "$tmp_dir/help.out"
! grep -q "^command  .*purpose" "$tmp_dir/help.out"
! grep -q "carbide help" "$tmp_dir/help.out"
! grep -q "carbide run dev" "$tmp_dir/help.out"
! grep -q "^Carbide$" "$tmp_dir/help.out"
! grep -q "Containerized full-stack apps with React, Go, and Postgres." "$tmp_dir/help.out"
! grep -q "_____________________________________________________" "$tmp_dir/help.out"
! grep -q "________________________oo_______oo_______oo_________" "$tmp_dir/help.out"
! grep -q "install the CLI" "$tmp_dir/help.out"
! grep -q "<github-install-url>" "$tmp_dir/help.out"
! grep -q "curl -fsSL" "$tmp_dir/help.out"
! grep -q "raw.githubusercontent.com/ryangerardwilson/carbide" "$tmp_dir/help.out"
! grep -q "features:" "$tmp_dir/help.out"
! grep -q "global actions:" "$tmp_dir/help.out"
! grep -q "carbide logs follow" "$tmp_dir/help.out"
! grep -q "carbide format" "$tmp_dir/help.out"

if "$repo_root/cli/bin/carbide" format >/tmp/carbide-format.out 2>/tmp/carbide-format.err; then
  printf 'carbide format should not exist\n' >&2
  exit 1
fi
grep -q "unknown command: format" /tmp/carbide-format.err

cd "$tmp_dir"
"$repo_root/cli/bin/carbide" new demo

test -f "$tmp_dir/demo/carbide.toml"
! test -f "$tmp_dir/demo/AGENTS.md"
! test -f "$tmp_dir/demo/README.md"
! test -f "$tmp_dir/demo/PROJECT.md"
test -f "$tmp_dir/demo/.env.example"
test -f "$tmp_dir/demo/.gitignore"
! test -d "$tmp_dir/demo/config"
! test -d "$tmp_dir/demo/view"
! test -d "$tmp_dir/demo/agents.d"
test -f "$tmp_dir/demo/docker-compose.yml"
test -f "$tmp_dir/demo/api/Dockerfile"
test -f "$tmp_dir/demo/web/Dockerfile"
test -f "$tmp_dir/demo/web/index.html"
test -f "$tmp_dir/demo/web/package.json"
test -f "$tmp_dir/demo/web/bun.lock"
test -f "$tmp_dir/demo/web/tsconfig.json"
test -f "$tmp_dir/demo/web/src/main.tsx"
test -f "$tmp_dir/demo/web/src/server.ts"
test -f "$tmp_dir/demo/web/src/write-index.ts"
test -f "$tmp_dir/demo/web/src/styles.css"
test -f "$tmp_dir/demo/web/src/styles.d.ts"
test -f "$tmp_dir/demo/web/src/lib/cx.ts"
test -f "$tmp_dir/demo/web/src/lib/types.ts"
! test -d "$tmp_dir/demo/web/node_modules"
! test -d "$tmp_dir/demo/web/public"
! test -f "$tmp_dir/demo/web/src/tailwind.css"
test -d "$tmp_dir/demo/web/src/component/l1"
test -d "$tmp_dir/demo/web/src/component/l2"
test -d "$tmp_dir/demo/web/src/component/l3"
test -f "$tmp_dir/demo/web/src/component/l1/Button.tsx"
test -f "$tmp_dir/demo/web/src/component/l1/Field.tsx"
test -f "$tmp_dir/demo/web/src/component/l1/Surface.tsx"
test -f "$tmp_dir/demo/web/src/component/l1/Text.tsx"
test -f "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
test -f "$tmp_dir/demo/web/src/component/l1/index.ts"
test -f "$tmp_dir/demo/web/src/component/l1/tokens.ts"
test -f "$tmp_dir/demo/web/src/component/l2/AuthForm.tsx"
test -f "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
test -f "$tmp_dir/demo/web/src/component/l2/index.ts"
test -f "$tmp_dir/demo/web/src/component/l3/AuthView.tsx"
test -f "$tmp_dir/demo/web/src/component/l3/DashboardView.tsx"
test -f "$tmp_dir/demo/web/src/component/l3/LoadingView.tsx"
test -f "$tmp_dir/demo/web/src/component/l3/index.ts"
! test -f "$tmp_dir/demo/web/package-lock.json"
! test -f "$tmp_dir/demo/web/vite.config.js"
test -f "$tmp_dir/demo/api/go.mod"
test -f "$tmp_dir/demo/api/go.sum"
test -f "$tmp_dir/demo/db/go.mod"
test -f "$tmp_dir/demo/db/go.sum"
test -f "$tmp_dir/demo/api/main.go"
test -f "$tmp_dir/demo/api/auth.go"
test -f "$tmp_dir/demo/api/routes.go"
test -f "$tmp_dir/demo/db/user.go"
test -f "$tmp_dir/demo/db/session.go"
test -f "$tmp_dir/demo/db/migration/001_auth.sql"
! test -d "$tmp_dir/demo/src"
! test -d "$tmp_dir/demo/model"
! test -d "$tmp_dir/demo/controller"
! test -d "$tmp_dir/demo/migrations"
! test -d "$tmp_dir/demo/infra"
! test -d "$tmp_dir/demo/frontend"
! test -d "$tmp_dir/demo/doc"
test "$(find "$tmp_dir/demo" -mindepth 1 -maxdepth 1 -type d -printf '%f\n' | sort | tr '\n' ' ')" = "api db web "
! test -f "$tmp_dir/demo/Dockerfile"
! test -f "$tmp_dir/demo/go.mod"
! test -f "$tmp_dir/demo/go.sum"
! test -f "$tmp_dir/demo/api/app.h"
! test -f "$tmp_dir/demo/api/main.c"
! test -f "$tmp_dir/demo/db/user.c"
! test -f "$tmp_dir/demo/db/session.c"

grep -q 'name = "Demo"' "$tmp_dir/demo/carbide.toml"
grep -q 'slug = "demo"' "$tmp_dir/demo/carbide.toml"
grep -q "default_port = 8080" "$tmp_dir/demo/carbide.toml"
grep -q "contract_version = 1" "$tmp_dir/demo/carbide.toml"
grep -q "\\[env.variables.DATABASE_URL\\]" "$tmp_dir/demo/carbide.toml"
grep -q "secret = true" "$tmp_dir/demo/carbide.toml"
grep -q "browser_exposed = true" "$tmp_dir/demo/carbide.toml"
grep -q "framework_owned = true" "$tmp_dir/demo/carbide.toml"
grep -q "preview_before_apply = true" "$tmp_dir/demo/carbide.toml"
! grep -q 'url = "http://localhost:8080"' "$tmp_dir/demo/carbide.toml"
grep -q 'name: demo' "$tmp_dir/demo/docker-compose.yml"
grep -q ".carbide/" "$tmp_dir/demo/.gitignore"
grep -q ".env" "$tmp_dir/demo/.gitignore"
grep -q "web/node_modules/" "$tmp_dir/demo/.gitignore"
grep -q "web/public/" "$tmp_dir/demo/.gitignore"
grep -q "web/src/tailwind.css" "$tmp_dir/demo/.gitignore"
grep -q "POSTGRES_PASSWORD" "$tmp_dir/demo/.env.example"
grep -q "web:" "$tmp_dir/demo/docker-compose.yml"
grep -q "api:" "$tmp_dir/demo/docker-compose.yml"
grep -q "db:" "$tmp_dir/demo/docker-compose.yml"
! grep -q "backend:" "$tmp_dir/demo/docker-compose.yml"
! grep -q "database:" "$tmp_dir/demo/docker-compose.yml"
grep -q "API_URL: http://api:8080" "$tmp_dir/demo/docker-compose.yml"
grep -q "@db:5432/carbide" "$tmp_dir/demo/docker-compose.yml"
grep -q 'PUBLIC_URL: "http://localhost:${CARBIDE_HTTP_PORT:-8080}"' "$tmp_dir/demo/docker-compose.yml"
test "$(grep -c 'PUBLIC_URL: "http://localhost:${CARBIDE_HTTP_PORT:-8080}"' "$tmp_dir/demo/docker-compose.yml")" -eq 2
grep -q 'PUBLIC_APP_NAME: "${PUBLIC_APP_NAME:-Demo}"' "$tmp_dir/demo/docker-compose.yml"
grep -q 'APP_ENV: "${APP_ENV:-development}"' "$tmp_dir/demo/docker-compose.yml"
grep -q 'POSTGRES_PASSWORD: "${POSTGRES_PASSWORD:-carbide}"' "$tmp_dir/demo/docker-compose.yml"
grep -q "develop:" "$tmp_dir/demo/docker-compose.yml"
grep -q "watch:" "$tmp_dir/demo/docker-compose.yml"
grep -q "action: rebuild" "$tmp_dir/demo/docker-compose.yml"
grep -q "context: ./web" "$tmp_dir/demo/docker-compose.yml"
grep -q "context: ." "$tmp_dir/demo/docker-compose.yml"
grep -q "dockerfile: api/Dockerfile" "$tmp_dir/demo/docker-compose.yml"
grep -q "path: ./web/src" "$tmp_dir/demo/docker-compose.yml"
grep -q "path: ./web/package.json" "$tmp_dir/demo/docker-compose.yml"
grep -q "path: ./web/bun.lock" "$tmp_dir/demo/docker-compose.yml"
grep -q "path: ./web/tsconfig.json" "$tmp_dir/demo/docker-compose.yml"
grep -q "path: ./api" "$tmp_dir/demo/docker-compose.yml"
grep -q "path: ./db" "$tmp_dir/demo/docker-compose.yml"
grep -q "path: ./api/Dockerfile" "$tmp_dir/demo/docker-compose.yml"
! grep -q "path: ./go.mod" "$tmp_dir/demo/docker-compose.yml"
! grep -q "path: ./go.sum" "$tmp_dir/demo/docker-compose.yml"
! grep -R 'admin@carbide.local' "$tmp_dir/demo" >/dev/null
! grep -R 'Demo login' "$tmp_dir/demo" >/dev/null
grep -q "module carbideapp/api" "$tmp_dir/demo/api/go.mod"
grep -q "carbideapp/db" "$tmp_dir/demo/api/go.mod"
grep -q "replace carbideapp/db => ../db" "$tmp_dir/demo/api/go.mod"
grep -q "module carbideapp/db" "$tmp_dir/demo/db/go.mod"
grep -q "github.com/jackc/pgx/v5" "$tmp_dir/demo/db/go.mod"
grep -q "package main" "$tmp_dir/demo/api/main.go"
grep -q "/api/login" "$tmp_dir/demo/api/routes.go"
grep -q "/api/me" "$tmp_dir/demo/api/routes.go"
grep -q "handleDashboard" "$tmp_dir/demo/api/routes.go"
grep -q "api listening on container port" "$tmp_dir/demo/api/main.go"
grep -q "public API URL is" "$tmp_dir/demo/api/main.go"
! grep -q "API listening inside api container" "$tmp_dir/demo/api/main.go"
! find "$tmp_dir/demo" -name '*.c' -o -name '*.h' | grep -q .
! grep -R "render_template_text" "$tmp_dir/demo" >/dev/null
! grep -R "respond_view" "$tmp_dir/demo" >/dev/null
grep -q "oven/bun:1.3.14-debian@sha256:9dba1a1b43ce28c9d7931bfc4eb00feb63b0114720a0277a8f939ae4dfc9db6f" "$tmp_dir/demo/web/Dockerfile"
grep -q "bun install --frozen-lockfile" "$tmp_dir/demo/web/Dockerfile"
grep -q '"@tailwindcss/cli": "4.3.2"' "$tmp_dir/demo/web/package.json"
grep -q '"tailwindcss": "4.3.2"' "$tmp_dir/demo/web/package.json"
grep -q '"react": "19.2.7"' "$tmp_dir/demo/web/package.json"
grep -q "go 1.25.0" "$tmp_dir/demo/api/go.mod"
grep -q "go 1.25.0" "$tmp_dir/demo/db/go.mod"
grep -q "postgres:17-alpine@sha256:dc17045ccfd343b49600570ea734b9c4991cf1c3f3302e67df51e3b402dd55c4" "$tmp_dir/demo/docker-compose.yml"
grep -q "FROM golang:1.26-bookworm@sha256:b305420a68d0f229d91eb3b3ed9e519fcf2cf5461da4bef997bf927e8c0bfd2b" "$tmp_dir/demo/api/Dockerfile"
grep -q "FROM debian:trixie-slim@sha256:28de0877c2189802884ccd20f15ee41c203573bd87bb6b883f5f46362d24c5c2" "$tmp_dir/demo/api/Dockerfile"
! grep -q "go 1.23.0" "$tmp_dir/demo/api/go.mod"
! grep -q "postgres:16-alpine" "$tmp_dir/demo/docker-compose.yml"
grep -q "Bun.serve" "$tmp_dir/demo/web/src/server.ts"
grep -q "browser entrypoint" "$tmp_dir/demo/web/src/server.ts"
grep -q "listening inside container" "$tmp_dir/demo/web/src/server.ts"
grep -q "proxying /api and /health to api service" "$tmp_dir/demo/web/src/server.ts"
grep -q "publicRoot" "$tmp_dir/demo/web/src/server.ts"
grep -q "Cache-Control" "$tmp_dir/demo/web/src/server.ts"
grep -q "public, max-age=31536000, immutable" "$tmp_dir/demo/web/src/server.ts"
grep -q "return 'no-store'" "$tmp_dir/demo/web/src/server.ts"
grep -q '"assets:build"' "$tmp_dir/demo/web/package.json"
grep -q '"typecheck": "tsc --noEmit"' "$tmp_dir/demo/web/package.json"
grep -q '"typescript": "6.0.3"' "$tmp_dir/demo/web/package.json"
grep -q '"@types/bun": "1.3.14"' "$tmp_dir/demo/web/package.json"
grep -q '"@types/react": "19.2.17"' "$tmp_dir/demo/web/package.json"
grep -q '"@types/react-dom": "19.2.3"' "$tmp_dir/demo/web/package.json"
grep -F -q "assets/[name]-[hash].[ext]" "$tmp_dir/demo/web/package.json"
grep -q '"strict": true' "$tmp_dir/demo/web/tsconfig.json"
grep -q '"jsx": "react-jsx"' "$tmp_dir/demo/web/tsconfig.json"
grep -F -q '"types": ["bun-types"]' "$tmp_dir/demo/web/tsconfig.json"
grep -q "bun run typecheck" "$tmp_dir/demo/web/Dockerfile"
grep -q "bun run assets:build" "$tmp_dir/demo/web/Dockerfile"
grep -q "asset-manifest.json" "$tmp_dir/demo/web/src/write-index.ts"
grep -F -q '/assets/${scripts[0]}' "$tmp_dir/demo/web/src/write-index.ts"
! grep -q "Bun frontend listening on http://localhost" "$tmp_dir/demo/web/src/server.ts"
grep -q '@import "tailwindcss";' "$tmp_dir/demo/web/src/styles.css"
grep -F -q '@source "./component/**/*.tsx";' "$tmp_dir/demo/web/src/styles.css"
grep -F -q '@source "./lib/**/*.ts";' "$tmp_dir/demo/web/src/styles.css"
grep -F -q '@source "./main.tsx";' "$tmp_dir/demo/web/src/styles.css"
grep -F -q '@source "./server.ts";' "$tmp_dir/demo/web/src/styles.css"
grep -q "@custom-variant dark" "$tmp_dir/demo/web/src/styles.css"
grep -q "\\[data-theme=\"dark\"\\]" "$tmp_dir/demo/web/src/styles.css"
! grep -q "html {" "$tmp_dir/demo/web/src/styles.css"
! grep -q "body {" "$tmp_dir/demo/web/src/styles.css"
! grep -q "font-size:" "$tmp_dir/demo/web/src/styles.css"
! grep -q "line-height:" "$tmp_dir/demo/web/src/styles.css"
! grep -q "min-width:" "$tmp_dir/demo/web/src/styles.css"
! grep -q "margin:" "$tmp_dir/demo/web/src/styles.css"
! grep -q "padding:" "$tmp_dir/demo/web/src/styles.css"
! grep -q "::-webkit-scrollbar" "$tmp_dir/demo/web/src/styles.css"
! grep -q "scrollbar-color:" "$tmp_dir/demo/web/src/styles.css"
! grep -q "scrollbar-width:" "$tmp_dir/demo/web/src/styles.css"
! grep -q "@theme" "$tmp_dir/demo/web/src/styles.css"
! grep -q -- "--carbide-" "$tmp_dir/demo/web/src/styles.css"
! grep -q "carbide-" "$tmp_dir/demo/web/src/component/l1/tokens.ts"
grep -q "const scrollbar =" "$tmp_dir/demo/web/src/component/l1/tokens.ts"
grep -F -q "[scrollbar-width:thin]" "$tmp_dir/demo/web/src/component/l1/tokens.ts"
grep -F -q "dark:[scrollbar-color:rgb(82_82_82)_transparent]" "$tmp_dir/demo/web/src/component/l1/tokens.ts"
grep -q "bg-white text-neutral-950 dark:bg-black dark:text-neutral-50" "$tmp_dir/demo/web/src/component/l1/tokens.ts"
! grep -Eq "#0f766e|#115e59|#2dd4bf|#5eead4|#16433c|#0f302c|#16211b|#edf5ef|#ecfdf5|#166534" "$tmp_dir/demo/web/src/styles.css"
! grep -q "from-carbide-action via-carbide-hero-via" "$tmp_dir/demo/web/src/component/l1/tokens.ts"
! grep -q "theme.css" "$tmp_dir/demo/web/src/styles.css"
! grep -Eq '^[[:space:]]*\.[A-Za-z_-]' "$tmp_dir/demo/web/src/styles.css"
! grep -Eq '^[[:space:]]*#[A-Za-z_-]' "$tmp_dir/demo/web/src/styles.css"
! grep -Eq '@theme|@apply|@layer|@keyframes|@media|@container|@plugin|@config' "$tmp_dir/demo/web/src/styles.css"
grep -F -q "text-2xl/8 sm:text-3xl/9" "$tmp_dir/demo/web/src/component/l1/Text.tsx"
grep -F -q "min-h-8 rounded-md border px-2 py-1 text-sm/6" "$tmp_dir/demo/web/src/component/l1/Field.tsx"
grep -F -q "md: 'min-h-8 px-3 text-xs'" "$tmp_dir/demo/web/src/component/l1/Button.tsx"
grep -F -q "gap-3 border-l px-4 py-5" "$tmp_dir/demo/web/src/component/l2/AuthForm.tsx"
grep -F -q "w-full max-w-sm justify-self-center gap-3" "$tmp_dir/demo/web/src/component/l2/AuthForm.tsx"
grep -F -q "lg:grid-cols-[216px_minmax(0,1fr)]" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -F -q "px-3 py-4 sm:px-5 lg:py-5" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "ui.scrollbar" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
! grep -R -E "text-7xl|text-5xl|py-24|lg:py-12|min-h-12 rounded-md border|min-h-10 rounded-md border|lg:grid-cols-\[280px|lg:grid-cols-\[240px|gap-6|p-6|font-extrabold" "$tmp_dir/demo/web/src/component" >/dev/null
grep -q '/api/${mode}' "$tmp_dir/demo/web/src/main.tsx"
grep -q "carbide.theme" "$tmp_dir/demo/web/src/main.tsx"
grep -q "useThemeMode" "$tmp_dir/demo/web/src/main.tsx"
grep -q "prefers-color-scheme: dark" "$tmp_dir/demo/web/index.html"
grep -q "dataset.theme" "$tmp_dir/demo/web/index.html"
grep -F -q "[scrollbar-width:thin]" "$tmp_dir/demo/web/index.html"
grep -F -q "dark:[scrollbar-color:rgb(82_82_82)_transparent]" "$tmp_dir/demo/web/index.html"
grep -q "ThemeToggle" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "data-resolved-theme" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "data-theme-mode" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "SunIcon" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "MoonIcon" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "Switch to light theme" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "Switch to dark theme" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "size-8 rounded-full border" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "onClick={() => onMode?.(nextMode)}" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
! grep -q "<select" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
! grep -q "appearance-none" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
! grep -q "border-x-4 border-t-4 border-x-transparent" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
! grep -q "aria-pressed" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
! grep -q "role=\"group\"" "$tmp_dir/demo/web/src/component/l1/ThemeToggle.tsx"
grep -q "./component/l3" "$tmp_dir/demo/web/src/main.tsx"
grep -q "AuthView" "$tmp_dir/demo/web/src/main.tsx"
grep -q "DashboardView" "$tmp_dir/demo/web/src/main.tsx"
grep -R -q "Lorem ipsum dolor sit amet" "$tmp_dir/demo/web/src/component"
grep -R -q "Consectetur adipiscing elit" "$tmp_dir/demo/web/src/component"
grep -R -q "Sed do eiusmod tempor" "$tmp_dir/demo/web/src/component"
! grep -R -q "Create the owner account" "$tmp_dir/demo/web/src/component"
! grep -R -q "Use your account email" "$tmp_dir/demo/web/src/component"
! grep -R -q "Your session is active" "$tmp_dir/demo/web/src/component"
! grep -R -q "Bun + Go + Postgres" "$tmp_dir/demo/web/src/component"
! grep -R -q "React + Bun container" "$tmp_dir/demo/web/src/component"
! grep -R -q "React and Tailwind" "$tmp_dir/demo/web/src/component"
! grep -R -q "Go owns" "$tmp_dir/demo/web/src/component"
! grep -R -q "Postgres" "$tmp_dir/demo/web/src/component"
grep -q "export const ui" "$tmp_dir/demo/web/src/component/l1/tokens.ts"
grep -q "ThemeToggle" "$tmp_dir/demo/web/src/component/l1/index.ts"
! test -f "$tmp_dir/demo/web/src/component/l1/theme.css"
! grep -R "cb-" "$tmp_dir/demo/web/src" >/dev/null
! grep -R -- "--cb-" "$tmp_dir/demo/web/src" >/dev/null
grep -q "ui.action" "$tmp_dir/demo/web/src/component/l1/Button.tsx"
grep -q "ui.input" "$tmp_dir/demo/web/src/component/l1/Field.tsx"
grep -q "buttonClassLayers" "$tmp_dir/demo/web/src/component/l1/Button.tsx"
grep -q "fieldClassLayers" "$tmp_dir/demo/web/src/component/l1/Field.tsx"
grep -q "fieldHintClassLayers" "$tmp_dir/demo/web/src/component/l1/Field.tsx"
grep -q "fieldErrorClassLayers" "$tmp_dir/demo/web/src/component/l1/Field.tsx"
grep -q "inputClassLayers" "$tmp_dir/demo/web/src/component/l1/Field.tsx"
grep -q "panelClassLayers" "$tmp_dir/demo/web/src/component/l1/Surface.tsx"
grep -q "dividerClassLayers" "$tmp_dir/demo/web/src/component/l1/Surface.tsx"
grep -q "badgeClassLayers" "$tmp_dir/demo/web/src/component/l1/Surface.tsx"
grep -q "metricClassLayers" "$tmp_dir/demo/web/src/component/l1/Surface.tsx"
grep -q "eyebrowClassLayers" "$tmp_dir/demo/web/src/component/l1/Text.tsx"
grep -q "headingClassLayers" "$tmp_dir/demo/web/src/component/l1/Text.tsx"
grep -q "mutedClassLayers" "$tmp_dir/demo/web/src/component/l1/Text.tsx"
grep -q "codeClassLayers" "$tmp_dir/demo/web/src/component/l1/Text.tsx"
grep -q "formClassLayers" "$tmp_dir/demo/web/src/component/l2/AuthForm.tsx"
grep -q "formStackClassLayers" "$tmp_dir/demo/web/src/component/l2/AuthForm.tsx"
grep -q "errorClassLayers" "$tmp_dir/demo/web/src/component/l2/AuthForm.tsx"
grep -q "modeButtonClassLayers" "$tmp_dir/demo/web/src/component/l2/AuthForm.tsx"
grep -q "landingClassLayers" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "dashboardClassLayers" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "screenClassLayers" "$tmp_dir/demo/web/src/component/l3/DashboardView.tsx"
grep -q "loadingClassLayers" "$tmp_dir/demo/web/src/component/l3/LoadingView.tsx"
grep -q "ui.focus" "$tmp_dir/demo/web/src/component/l1/Field.tsx"
grep -q "ui.focus" "$tmp_dir/demo/web/src/component/l2/AuthForm.tsx"
grep -q "ui.focus" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
! grep -R "text-\\[" "$tmp_dir/demo/web/src/component" >/dev/null
grep -q "export function DashboardLayout" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "lg:grid-cols-\\[216px_minmax(0,1fr)\\]" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "aria-label=\"Dashboard\"" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "aria-current" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "navItems" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "export function LandingPageLayout" "$tmp_dir/demo/web/src/component/l2/Layouts.tsx"
grep -q "dashboardNav" "$tmp_dir/demo/web/src/component/l3/DashboardView.tsx"
grep -q "WorkspaceOverview" "$tmp_dir/demo/web/src/component/l3/DashboardView.tsx"
! grep -R "ComponentLibraryView" "$tmp_dir/demo/web/src/component" >/dev/null
! find "$tmp_dir/demo" -path '*/ui_components/*' -print -quit | grep -q .
test -d "$tmp_dir/demo/web/src/component/l1"
test -d "$tmp_dir/demo/web/src/component/l2"
test -d "$tmp_dir/demo/web/src/component/l3"
test "$(find "$tmp_dir/demo/web/src/component" -mindepth 1 -maxdepth 1 -type d -printf '%f\n' | sort | tr '\n' ' ')" = "l1 l2 l3 "
! find "$tmp_dir/demo/web/src/component" -mindepth 1 -maxdepth 1 -type f -print -quit | grep -q .
! test -d "$tmp_dir/demo/web/src/component/ui"
! test -d "$tmp_dir/demo/web/src/component/screen"
! grep -R "views/" "$tmp_dir/demo" >/dev/null
! grep -R "__PROJECT_" "$tmp_dir/demo" >/dev/null

cd "$tmp_dir/demo"
CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" health env > "$tmp_dir/health-env.out"
grep -q "Carbide health" "$tmp_dir/health-env.out"
grep -q "environment contract" "$tmp_dir/health-env.out"
grep -Eq "^status[[:space:]]+ok" "$tmp_dir/health-env.out"
grep -Eq "^required[[:space:]]+0 missing" "$tmp_dir/health-env.out"
grep -Eq "^secrets[[:space:]]+2 declared" "$tmp_dir/health-env.out"
! grep -q "postgres://carbide:carbide" "$tmp_dir/health-env.out"

CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" health > "$tmp_dir/health.out"
grep -q "Carbide health" "$tmp_dir/health.out"
grep -q "app laws" "$tmp_dir/health.out"
grep -Eq "^project shape[[:space:]]+ok[[:space:]]+web api db" "$tmp_dir/health.out"
grep -Eq "^config[[:space:]]+ok[[:space:]]+carbide.toml" "$tmp_dir/health.out"
grep -Eq "^env contract[[:space:]]+ok[[:space:]]+0 missing, 2 secrets" "$tmp_dir/health.out"
grep -Eq "^compose[[:space:]]+ok[[:space:]]+web api db" "$tmp_dir/health.out"
grep -Eq "^regressions[[:space:]]+ok[[:space:]]+no legacy markers" "$tmp_dir/health.out"
grep -Eq "^runtime[[:space:]]+skip[[:space:]]+run carbide health runtime" "$tmp_dir/health.out"
! grep -q "postgres://carbide:carbide" "$tmp_dir/health.out"

CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" health json > "$tmp_dir/health.json"
grep -q '"command": "health"' "$tmp_dir/health.json"
grep -q '"ok": true' "$tmp_dir/health.json"
grep -q '"check": "config"' "$tmp_dir/health.json"
! grep -q '"check": "agents"' "$tmp_dir/health.json"

CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" health env json > "$tmp_dir/health-env.json"
grep -q '"command": "health env"' "$tmp_dir/health-env.json"
grep -q '"status": "ok"' "$tmp_dir/health-env.json"

CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" urls > "$tmp_dir/urls.out"
grep -q "Carbide urls" "$tmp_dir/urls.out"
grep -q "http://localhost:8080" "$tmp_dir/urls.out"

CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" urls json > "$tmp_dir/urls.json"
grep -q '"command": "urls"' "$tmp_dir/urls.json"
grep -q '"app": "http://localhost:8080"' "$tmp_dir/urls.json"

CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" deploy check prod > "$tmp_dir/deploy-check.out"
grep -q "Carbide deploy" "$tmp_dir/deploy-check.out"
grep -q "check prod" "$tmp_dir/deploy-check.out"
grep -Eq "^state[[:space:]]+missing-target" "$tmp_dir/deploy-check.out"
grep -Eq "^apply[[:space:]]+no" "$tmp_dir/deploy-check.out"

CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" deploy check prod json > "$tmp_dir/deploy-check.json"
grep -q '"command": "deploy check"' "$tmp_dir/deploy-check.json"
grep -q '"classification": "missing-target"' "$tmp_dir/deploy-check.json"

CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" deploy preview prod > "$tmp_dir/deploy-preview.out"
grep -q "Carbide deploy" "$tmp_dir/deploy-preview.out"
grep -q "preview prod" "$tmp_dir/deploy-preview.out"
grep -Eq "^state[[:space:]]+missing-target" "$tmp_dir/deploy-preview.out"
grep -Eq "^mutates[[:space:]]+no" "$tmp_dir/deploy-preview.out"
grep -q "add a checked-in deploy target to carbide.toml" "$tmp_dir/deploy-preview.out"
grep -q "single-VM ssh-compose targets can apply" "$tmp_dir/deploy-preview.out"

if CARBIDE_HOME="$repo_root" "$repo_root/cli/bin/carbide" deploy apply prod > "$tmp_dir/deploy-apply.out" 2> "$tmp_dir/deploy-apply.err"; then
  printf 'carbide deploy apply should be guarded until a target exists\n' >&2
  exit 1
fi
grep -q "status.*disabled" "$tmp_dir/deploy-apply.out"
grep -q "no checked-in deploy target exists" "$tmp_dir/deploy-apply.out"
grep -q "disabled until a checked-in deploy target exists" "$tmp_dir/deploy-apply.err"

mkdir "$tmp_dir/init-app"
cd "$tmp_dir/init-app"
"$repo_root/cli/bin/carbide" init
test -f "$tmp_dir/init-app/carbide.toml"
grep -q 'name = "Init App"' "$tmp_dir/init-app/carbide.toml"
grep -q 'slug = "init-app"' "$tmp_dir/init-app/carbide.toml"

cd "$tmp_dir"
"$repo_root/cli/bin/carbide" new My Carbide App
test -f "$tmp_dir/my-carbide-app/carbide.toml"
grep -q 'name = "My Carbide App"' "$tmp_dir/my-carbide-app/carbide.toml"
grep -q 'slug = "my-carbide-app"' "$tmp_dir/my-carbide-app/carbide.toml"
grep -q 'PUBLIC_APP_NAME: "${PUBLIC_APP_NAME:-My Carbide App}"' "$tmp_dir/my-carbide-app/docker-compose.yml"

mkdir "$tmp_dir/not-empty"
touch "$tmp_dir/not-empty/file"
cd "$tmp_dir/not-empty"
if "$repo_root/cli/bin/carbide" init >/tmp/carbide-init.out 2>/tmp/carbide-init.err; then
  printf 'carbide init should fail in a non-empty directory\n' >&2
  exit 1
fi
grep -q "requires an empty directory" /tmp/carbide-init.err

if command -v python3 >/dev/null 2>&1; then
  fake_bin="$tmp_dir/fake-bin"
  port_file="$tmp_dir/selected-port"
  args_file="$tmp_dir/docker-args"
  mkdir "$fake_bin"
  cat > "$fake_bin/docker" <<'SH'
#!/usr/bin/env bash
set -euo pipefail

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "version" ]; then
  printf 'Docker Compose fake\n'
  exit 0
fi

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "up" ] && [ "${3:-}" = "--help" ]; then
  printf 'Usage: docker compose up [OPTIONS]\n'
  printf '      --quiet-build    Suppress the build output\n'
  printf '      --quiet-pull     Pull without printing progress information\n'
  printf '      --wait           Wait for services to be running|healthy\n'
  printf '      --wait-timeout int\n'
  printf '      --watch    Watch source code and rebuild/refresh containers when files are updated.\n'
  exit 0
fi

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "logs" ] && [ "${3:-}" = "--help" ]; then
  printf 'Usage: docker compose logs [OPTIONS]\n'
  printf '      --no-color    Produce monochrome output\n'
  printf '      --tail string Number of lines to show from the end of the logs\n'
  exit 0
fi

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "config" ] && [ "${3:-}" = "--services" ]; then
  printf 'web\napi\ndb\n'
  exit 0
fi

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "ps" ] && [ "${3:-}" = "--format" ] && [ "${4:-}" = "json" ]; then
  status_port="8082"
  if [ -n "${FAKE_DOCKER_PORT_FILE:-}" ] && [ -s "$FAKE_DOCKER_PORT_FILE" ]; then
    status_port="$(cat "$FAKE_DOCKER_PORT_FILE")"
  fi
  printf '{"Service":"web","Name":"demo-web-1","State":"running","Health":"healthy","Publishers":[{"URL":"0.0.0.0","TargetPort":8080,"PublishedPort":%s,"Protocol":"tcp"},{"URL":"::","TargetPort":8080,"PublishedPort":%s,"Protocol":"tcp"}]}\n' "$status_port" "$status_port"
  printf '{"Service":"api","Name":"demo-api-1","State":"running","Health":"healthy","Publishers":[{"URL":"","TargetPort":8080,"PublishedPort":0,"Protocol":"tcp"}]}\n'
  printf '{"Service":"db","Name":"demo-db-1","State":"running","Health":"healthy","Publishers":[{"URL":"","TargetPort":5432,"PublishedPort":0,"Protocol":"tcp"}]}\n'
  exit 0
fi

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "up" ]; then
  printf '%s\n' "${CARBIDE_HTTP_PORT:-}" > "$FAKE_DOCKER_PORT_FILE"
  printf '%s\n' "$*" > "$FAKE_DOCKER_ARGS_FILE"
  exit 0
fi

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "logs" ]; then
  printf '%s\n' "$*" >> "$FAKE_DOCKER_ARGS_FILE"
  printf 'api-1  | GET /health\n'
  printf 'web-1 | listening on :8080\n'
  if [ "${FAKE_DOCKER_STREAM_LONG:-}" = "1" ]; then
    trap 'exit 0' INT TERM
    while true; do sleep 1; done
  fi
  sleep 0.2
  exit 0
fi

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "watch" ]; then
  printf '%s\n' "$*" >> "$FAKE_DOCKER_ARGS_FILE"
  printf 'Watch enabled\n'
  printf 'rebuilding api\n'
  if [ "${FAKE_DOCKER_STREAM_LONG:-}" = "1" ]; then
    trap 'exit 0' INT TERM
    while true; do sleep 1; done
  fi
  sleep 0.2
  exit 0
fi

if [ "${1:-}" = "compose" ] && [ "${2:-}" = "down" ]; then
  printf '%s\n' "$*" >> "$FAKE_DOCKER_ARGS_FILE"
  exit 0
fi

printf 'unexpected fake docker command: %s\n' "$*" >&2
exit 1
SH
  chmod +x "$fake_bin/docker"

  python3 - <<'PY' &
import socket
import sys
import time

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
try:
    sock.bind(("0.0.0.0", 8080))
    sock.listen(1)
except OSError:
    sys.exit(0)
time.sleep(60)
PY
  listener_pid="$!"
  sleep 0.5

  cd "$tmp_dir/demo"
  PATH="$fake_bin:$PATH" FAKE_DOCKER_PORT_FILE="$port_file" FAKE_DOCKER_ARGS_FILE="$args_file" "$repo_root/cli/bin/carbide" run dev > "$tmp_dir/run-dev.out"
  grep -q "Carbide dev" "$tmp_dir/run-dev.out"
  grep -Eq "^app[[:space:]]+http://localhost:" "$tmp_dir/run-dev.out"
  grep -Eq "^api[[:space:]]+http://localhost:" "$tmp_dir/run-dev.out"
  ! grep -Eq "^port[[:space:]]+" "$tmp_dir/run-dev.out"
  ! grep -Eq "^login[[:space:]]+" "$tmp_dir/run-dev.out"
  ! grep -Eq "^mode[[:space:]]+" "$tmp_dir/run-dev.out"
  ! grep -Eq "^status[[:space:]]+" "$tmp_dir/run-dev.out"
  ! grep -Eq "^containers[[:space:]]+" "$tmp_dir/run-dev.out"
  ! grep -Eq "^logs[[:space:]]+" "$tmp_dir/run-dev.out"
  ! grep -Eq "^stop[[:space:]]+" "$tmp_dir/run-dev.out"
  ! grep -Eq "^watch[[:space:]]+enabled" "$tmp_dir/run-dev.out"
  ! grep -q "busy, using" "$tmp_dir/run-dev.out"
  ! grep -q "^Watch enabled$" "$tmp_dir/run-dev.out"
  grep -Eq "^[0-9]{2}:[0-9]{2}:[0-9]{2}[[:space:]]+api[[:space:]]+GET /health" "$tmp_dir/run-dev.out"
  grep -Eq "^[0-9]{2}:[0-9]{2}:[0-9]{2}[[:space:]]+web[[:space:]]+listening on :8080" "$tmp_dir/run-dev.out"
  grep -Eq "^[0-9]{2}:[0-9]{2}:[0-9]{2}[[:space:]]+watch[[:space:]]+rebuilding api" "$tmp_dir/run-dev.out"
  test -f "$tmp_dir/demo/.carbide/log/dev.jsonl"
  grep -q '"service":"api"' "$tmp_dir/demo/.carbide/log/dev.jsonl"
  grep -q '"message":"GET /health"' "$tmp_dir/demo/.carbide/log/dev.jsonl"
  PATH="$fake_bin:$PATH" "$repo_root/cli/bin/carbide" logs service api > "$tmp_dir/logs-api.out"
  grep -Eq "^[0-9]{2}:[0-9]{2}:[0-9]{2}[[:space:]]+api[[:space:]]+GET /health" "$tmp_dir/logs-api.out"
  PATH="$fake_bin:$PATH" "$repo_root/cli/bin/carbide" logs json containing listening > "$tmp_dir/logs-json.out"
  grep -q '"service":"web"' "$tmp_dir/logs-json.out"
  grep -q '"message":"listening on :8080"' "$tmp_dir/logs-json.out"
  PATH="$fake_bin:$PATH" FAKE_DOCKER_PORT_FILE="$port_file" "$repo_root/cli/bin/carbide" status > "$tmp_dir/status.out"
  grep -q "Carbide status" "$tmp_dir/status.out"
  grep -Eq "^service[[:space:]]+container[[:space:]]+ports[[:space:]]+internal[[:space:]]+status" "$tmp_dir/status.out"
  grep -Eq "^web[[:space:]]+demo-web-1[[:space:]]+localhost:[0-9]+[[:space:]]+8080/tcp[[:space:]]+running \\(healthy\\)" "$tmp_dir/status.out"
  grep -Eq "^api[[:space:]]+demo-api-1[[:space:]]+-[[:space:]]+8080/tcp[[:space:]]+running \\(healthy\\)" "$tmp_dir/status.out"
  grep -Eq "^db[[:space:]]+demo-db-1[[:space:]]+-[[:space:]]+5432/tcp[[:space:]]+running \\(healthy\\)" "$tmp_dir/status.out"
  PATH="$fake_bin:$PATH" FAKE_DOCKER_PORT_FILE="$port_file" "$repo_root/cli/bin/carbide" status json > "$tmp_dir/status.json"
  grep -q '"command": "status"' "$tmp_dir/status.json"
  grep -q '"service": "web"' "$tmp_dir/status.json"
  grep -F -q '"published_ports": [' "$tmp_dir/status.json"
  grep -F -q '"internal_ports": [' "$tmp_dir/status.json"
  grep -q '"status": "running (healthy)"' "$tmp_dir/status.json"
  grep -q -- "--quiet-build" "$args_file"
  grep -q -- "--quiet-pull" "$args_file"
  grep -q "compose logs -f --tail 80 --no-color" "$args_file"
  grep -q "compose watch --no-up --quiet" "$args_file"
  ! grep -q "compose down" "$args_file"
  PATH="$fake_bin:$PATH" FAKE_DOCKER_ARGS_FILE="$args_file" "$repo_root/cli/bin/carbide" stop dev > "$tmp_dir/stop-dev.out"
  grep -q "Carbide stop dev" "$tmp_dir/stop-dev.out"
  grep -Eq "^dev[[:space:]]+stopped" "$tmp_dir/stop-dev.out"
  grep -q "compose down --remove-orphans" "$args_file"
  : > "$args_file"
  PATH="$fake_bin:$PATH" FAKE_DOCKER_ARGS_FILE="$args_file" "$repo_root/cli/bin/carbide" clean dev > "$tmp_dir/clean-dev.out"
  grep -q "Carbide clean dev" "$tmp_dir/clean-dev.out"
  grep -Eq "^dev[[:space:]]+clean" "$tmp_dir/clean-dev.out"
  grep -Eq "^next[[:space:]]+carbide run dev" "$tmp_dir/clean-dev.out"
  grep -q "compose down --remove-orphans" "$args_file"
  PATH="$fake_bin:$PATH" FAKE_DOCKER_ARGS_FILE="$args_file" "$repo_root/cli/bin/carbide" follow logs service api > "$tmp_dir/logs-follow.out"
  grep -Eq "^[0-9]{2}:[0-9]{2}:[0-9]{2}[[:space:]]+api[[:space:]]+GET /health" "$tmp_dir/logs-follow.out"
  ! grep -q "web" "$tmp_dir/logs-follow.out"
  selected_port="$(cat "$port_file")"
  if [ "$selected_port" = "8080" ]; then
    printf 'carbide run dev should not select occupied port 8080\n' >&2
    exit 1
  fi

  if PATH="$fake_bin:$PATH" FAKE_DOCKER_PORT_FILE="$port_file" FAKE_DOCKER_ARGS_FILE="$args_file" CARBIDE_HTTP_PORT=8080 "$repo_root/cli/bin/carbide" run dev > "$tmp_dir/explicit-port.out" 2> "$tmp_dir/explicit-port.err"; then
    printf 'explicit occupied CARBIDE_HTTP_PORT should fail before compose starts\n' >&2
    exit 1
  fi
  grep -q "port 8080 is already in use" "$tmp_dir/explicit-port.err"

  : > "$args_file"
  PATH="$fake_bin:$PATH" FAKE_DOCKER_STREAM_LONG=1 FAKE_DOCKER_PORT_FILE="$port_file" FAKE_DOCKER_ARGS_FILE="$args_file" "$repo_root/cli/bin/carbide" run dev > "$tmp_dir/run-dev-detach.out" &
  run_dev_pid="$!"
  for _ in $(seq 1 50); do
    if grep -q "GET /health" "$tmp_dir/run-dev-detach.out" 2>/dev/null; then
      break
    fi
    sleep 0.1
  done
  kill -INT "$run_dev_pid" >/dev/null 2>&1 || true
  for _ in $(seq 1 30); do
    if ! kill -0 "$run_dev_pid" >/dev/null 2>&1; then
      break
    fi
    sleep 0.1
  done
  if kill -0 "$run_dev_pid" >/dev/null 2>&1; then
    kill -TERM "$run_dev_pid" >/dev/null 2>&1 || true
  fi
  wait "$run_dev_pid"
  grep -q "detached from logs; containers are still running" "$tmp_dir/run-dev-detach.out"
  grep -Eq "^status[[:space:]]+carbide status" "$tmp_dir/run-dev-detach.out"
  grep -Eq "^logs[[:space:]]+carbide follow logs" "$tmp_dir/run-dev-detach.out"
  grep -Eq "^clean[[:space:]]+carbide clean dev" "$tmp_dir/run-dev-detach.out"
  ! grep -q "compose down" "$args_file"

  kill "$listener_pid" >/dev/null 2>&1 || true
fi

remote_repo="$tmp_dir/carbide-origin.git"
installed_repo="$tmp_dir/installed-carbide"
upgrade_work="$tmp_dir/upgrade-work"

git init --bare "$remote_repo" >/dev/null
git init "$installed_repo" >/dev/null
cp "$repo_root/.gitignore" "$installed_repo/.gitignore"
cp -R "$repo_root/cli" "$installed_repo/cli"
git -C "$installed_repo" add .gitignore cli
git -C "$installed_repo" -c user.name="Carbide Test" -c user.email="test@carbide.local" commit -m "Initial install" >/dev/null
git -C "$installed_repo" branch -M main
git -C "$installed_repo" remote add origin "$remote_repo"
git -C "$installed_repo" push -u origin main >/dev/null
git --git-dir="$remote_repo" symbolic-ref HEAD refs/heads/main

CARBIDE_HOME="$installed_repo" "$repo_root/cli/bin/carbide" upgrade > "$tmp_dir/upgrade-current.out"
grep -q "Carbide upgrade" "$tmp_dir/upgrade-current.out"
grep -Eq "^status[[:space:]]+up to date" "$tmp_dir/upgrade-current.out"

git clone --branch main "$remote_repo" "$upgrade_work" >/dev/null
printf '# changed\n' >> "$upgrade_work/README.md"
git -C "$upgrade_work" add README.md
git -C "$upgrade_work" -c user.name="Carbide Test" -c user.email="test@carbide.local" commit -m "Remote update" >/dev/null
git -C "$upgrade_work" push >/dev/null

CARBIDE_HOME="$installed_repo" "$repo_root/cli/bin/carbide" upgrade > "$tmp_dir/upgrade-new.out"
grep -q "Carbide upgrade" "$tmp_dir/upgrade-new.out"
grep -Eq "^status[[:space:]]+upgraded" "$tmp_dir/upgrade-new.out"
test -x "$installed_repo/.cli/bin/carbide"

printf 'cli scaffold ok\n'
