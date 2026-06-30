#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." >/dev/null 2>&1 && pwd)"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

export SEALION_HOME="$repo_root"

cd "$tmp_dir"
"$repo_root/bin/sealion" new demo

test -f "$tmp_dir/demo/sealion.toml"
test -f "$tmp_dir/demo/docker-compose.yml"
test -f "$tmp_dir/demo/Dockerfile"
test -f "$tmp_dir/demo/src/main.c"
test -f "$tmp_dir/demo/migrations/001_auth.sql"

grep -q 'name = "demo"' "$tmp_dir/demo/sealion.toml"
grep -q 'name: demo' "$tmp_dir/demo/docker-compose.yml"
grep -q 'admin@sealion.local' "$tmp_dir/demo/src/main.c"
! grep -R "__PROJECT_" "$tmp_dir/demo" >/dev/null

mkdir "$tmp_dir/init-app"
cd "$tmp_dir/init-app"
"$repo_root/bin/sealion" init
test -f "$tmp_dir/init-app/sealion.toml"
grep -q 'name = "init-app"' "$tmp_dir/init-app/sealion.toml"

mkdir "$tmp_dir/not-empty"
touch "$tmp_dir/not-empty/file"
cd "$tmp_dir/not-empty"
if "$repo_root/bin/sealion" init >/tmp/sealion-init.out 2>/tmp/sealion-init.err; then
  printf 'sealion init should fail in a non-empty directory\n' >&2
  exit 1
fi
grep -q "requires an empty directory" /tmp/sealion-init.err

printf 'cli scaffold ok\n'

