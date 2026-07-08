#!/usr/bin/env bash
set -euo pipefail

base_url="${CARBIDE_DOCS_BASE_URL:-https://carbide.ryangerardwilson.com}"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

headers="$tmp_dir/headers.txt"
body="$tmp_dir/agents.md"

curl -fsS -D "$headers" "$base_url/for/agents" -o "$body"

grep -qi '^content-type: .*text/\(markdown\|plain\)' "$headers"
grep -q '# Carbide for Agents' "$body"
grep -q '## Source Precedence' "$body"
grep -q '## Identify The Current State' "$body"
grep -q '## Laws' "$body"
grep -q '### Law 8. Checked Files Stay Under 1000 Lines' "$body"
grep -q '## Audits' "$body"
grep -q '## Recovery' "$body"
grep -q 'carbide health json' "$body"
grep -q 'carbide deploy prod' "$body"
grep -q 'carbide audit' "$body"
grep -q 'carbide resolve' "$body"
grep -q 'carbide fix' "$body"
grep -q 'carbide audit resolve fix' "$body"
grep -q '.audit/plan.md' "$body"
grep -q 'README.md' "$body"

lines="$(wc -l < "$body" | tr -d ' ')"
if [ "$lines" -lt 120 ]; then
  printf '/for/agents is too short: %s lines\n' "$lines" >&2
  exit 1
fi
