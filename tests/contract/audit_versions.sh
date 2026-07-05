#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." >/dev/null 2>&1 && pwd)"
config="$repo_root/scaffold/carbide.toml"

toml_value() {
  local key="$1"
  sed -n "s/^${key} = \"\\(.*\\)\"$/\\1/p" "$config" | head -n 1
}

print_row() {
  printf '%-18s %-64s %s\n' "$1" "$2" "$3"
}

latest_npm() {
  local package="$1"
  if ! command -v npm >/dev/null 2>&1; then
    printf 'npm unavailable'
    return
  fi
  npm view "$package" version 2>/dev/null || printf 'lookup failed'
}

registry_digest() {
  local ref="$1"
  local tag="${ref%@sha256:*}"
  if ! command -v docker >/dev/null 2>&1; then
    printf 'docker unavailable'
    return
  fi
  docker buildx imagetools inspect "$tag" 2>/dev/null | sed -n 's/^Digest:[[:space:]]*//p' | head -n 1 || printf 'lookup failed'
}

print_docker_row() {
  local name="$1"
  local ref="$2"
  local pinned="${ref##*@}"
  local current
  current="$(registry_digest "$ref")"
  if [ "$current" = "$pinned" ]; then
    print_row "$name" "$ref" "current"
  else
    print_row "$name" "$ref" "registry digest: $current"
  fi
}

printf 'Carbide dependency audit\n'
printf 'reports only; no files are modified\n\n'
printf '%-18s %-64s %s\n' "runtime" "pinned" "latest/check"
printf '%-18s %-64s %s\n' "-------" "------" "------------"

print_row "go module" "$(toml_value go_module)" "builder $(toml_value go_builder_image)"
print_docker_row "go builder" "$(toml_value go_builder_image)"
print_docker_row "api runtime" "$(toml_value api_runtime_image)"
print_docker_row "bun" "$(toml_value bun_image)"
print_docker_row "postgres" "$(toml_value postgres_image)"
print_row "react" "$(toml_value react)" "npm $(latest_npm react)"
print_row "react-dom" "$(toml_value react_dom)" "npm $(latest_npm react-dom)"
print_row "tailwindcss" "$(toml_value tailwindcss)" "npm $(latest_npm tailwindcss)"
print_row "tailwind cli" "$(toml_value tailwind_cli)" "npm $(latest_npm @tailwindcss/cli)"

if command -v go >/dev/null 2>&1; then
  printf '\nGo module update report\n'
  (cd "$repo_root/scaffold/api" && go list -m -u all) || true
  (cd "$repo_root/scaffold/db" && go list -m -u all) || true
else
  printf '\nGo module update report skipped: go unavailable\n'
fi
