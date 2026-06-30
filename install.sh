#!/usr/bin/env bash
set -euo pipefail

repo_url="${SEALION_REPO_URL:-https://github.com/ryangerardwilson/sealion.git}"
archive_url="${SEALION_ARCHIVE_URL:-https://github.com/ryangerardwilson/sealion/archive/refs/heads/main.tar.gz}"
install_dir="${SEALION_HOME:-$HOME/.sealion}"
bin_dir="${SEALION_BIN_DIR:-$HOME/.local/bin}"

mkdir -p "$bin_dir"

if command -v git >/dev/null 2>&1; then
  if [ -d "$install_dir/.git" ]; then
    git -C "$install_dir" pull --ff-only
  else
    rm -rf "$install_dir"
    git clone --depth 1 "$repo_url" "$install_dir"
  fi
else
  command -v curl >/dev/null 2>&1 || {
    printf 'install failed: git or curl is required\n' >&2
    exit 1
  }
  command -v tar >/dev/null 2>&1 || {
    printf 'install failed: tar is required\n' >&2
    exit 1
  }
  tmp_dir="$(mktemp -d)"
  trap 'rm -rf "$tmp_dir"' EXIT
  curl -fsSL "$archive_url" | tar -xz -C "$tmp_dir" --strip-components=1
  rm -rf "$install_dir"
  mkdir -p "$install_dir"
  cp -R "$tmp_dir/." "$install_dir/"
fi

chmod +x "$install_dir/bin/sealion"
ln -sfn "$install_dir/bin/sealion" "$bin_dir/sealion"

printf 'installed sealion to %s\n' "$bin_dir/sealion"
case ":$PATH:" in
  *":$bin_dir:"*) ;;
  *) printf 'add %s to PATH if sealion is not found\n' "$bin_dir" ;;
esac

