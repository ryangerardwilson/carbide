#!/usr/bin/env bash
set -euo pipefail

limit=1000

is_checked_file() {
  local path="$1"
  case "$path" in
    *.go|*.ts|*.tsx|*.js|*.jsx|*.sh|*.css|*.sql|*.toml|*.yml|*.yaml|*Dockerfile)
      ;;
    *)
      return 1
      ;;
  esac

  case "$path" in
    */bun.lock|*/package-lock.json|*/pnpm-lock.yaml|*/yarn.lock)
      return 1
      ;;
  esac

  return 0
}

violations=()

while IFS= read -r path; do
  if ! is_checked_file "$path"; then
    continue
  fi

  line_count=$(wc -l < "$path")
  if [ "$line_count" -gt "$limit" ]; then
    violations+=("$(printf '%5d  %s' "$line_count" "$path")")
  fi
done < <(git ls-files)

if [ "${#violations[@]}" -gt 0 ]; then
  printf 'files over %d lines:\n' "$limit" >&2
  printf '%s\n' "${violations[@]}" | sort -nr >&2
  exit 1
fi

printf 'file line limits ok\n'
