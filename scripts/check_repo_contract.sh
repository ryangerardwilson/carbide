#!/usr/bin/env bash
set -euo pipefail

domain="sealion.ryangerardwilson.com"

required_files=(
  "README.md"
  "install.sh"
  "bin/sealion"
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
  "templates/default/docker-compose.yml"
  "templates/default/sealion.toml"
  "templates/default/src/main.c"
  "templates/default/views/layout.html"
  "templates/default/views/home.html"
  "templates/default/views/register.html"
  "templates/default/views/login.html"
  "templates/default/views/dashboard.html"
  "templates/default/views/not_found.html"
  "templates/default/migrations/001_auth.sql"
)

required_dirs=(
  "src"
  "src/ui"
  "include/sealion"
  "include/sealion/ui"
  "tests/unit"
  "tests/integration"
  "tests/regression"
  "tests/fixtures"
  "examples/hello"
  "infra/compose"
  "infra/schemas"
  "templates/default"
  "templates/default/src"
  "templates/default/views"
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

grep -q "one mandatory app container image" README.md
grep -q "Postgres-only" README.md
grep -q "Separate runtime boundaries" README.md
grep -q "Infrastructure as code" README.md
grep -q "Framework-owned component styling" README.md
grep -q "generated Docker Compose setup" README.md
grep -q "Postgres-backed queues" README.md
grep -q "sealion new" README.md
grep -q "sealion run dev" README.md
grep -q "default_port = 8080" templates/default/sealion.toml
! grep -q 'url = "http://localhost:8080"' templates/default/sealion.toml
grep -q 'PUBLIC_URL: "http://localhost:${SEALION_HTTP_PORT:-8080}"' templates/default/docker-compose.yml
grep -q "develop:" templates/default/docker-compose.yml
grep -q "watch:" templates/default/docker-compose.yml
grep -q "action: rebuild" templates/default/docker-compose.yml
grep -q "path: ./src" templates/default/docker-compose.yml
grep -q "path: ./views" templates/default/docker-compose.yml
grep -q "path: ./Dockerfile" templates/default/docker-compose.yml
grep -q "COPY views ./views" templates/default/Dockerfile
grep -q "{{ title }}" templates/default/views/layout.html
grep -q "{!! content !!}" templates/default/views/layout.html
grep -q "{{ user_email }}" templates/default/views/dashboard.html
grep -q "render_template_text" templates/default/src/main.c
grep -q "respond_view" templates/default/src/main.c
! grep -q "<style>" templates/default/src/main.c
grep -q "listening inside container" templates/default/src/main.c
grep -q "open %s" templates/default/src/main.c
grep -q "compose_supports_watch" bin/sealion
grep -q -- "--watch" bin/sealion

grep -q "$domain" docs/site/index.html
grep -q "Component styling" docs/site/index.html
grep -q "Initial user experience" docs/site/index.html
grep -q "Tailwind-like ergonomics without the Tailwind dependency" docs/site/component-style-system.html
grep -q "Install, create, run, log in" docs/site/initial-user-experience.html
grep -q "CI/CD regression plan" docs/site/ci-cd-regression-tests.html
grep -q "Directory structure" docs/site/repo-structure.html

printf 'repo contract ok\n'
