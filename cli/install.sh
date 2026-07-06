#!/usr/bin/env bash
set -euo pipefail

repo_url="${CARBIDE_REPO_URL:-https://github.com/ryangerardwilson/carbide.git}"
repo_owner="${CARBIDE_REPO_OWNER:-ryangerardwilson}"
repo_name="${CARBIDE_REPO_NAME:-carbide}"
install_dir="${CARBIDE_HOME:-$HOME/.carbide}"
bin_dir="${CARBIDE_BIN_DIR:-$HOME/.local/bin}"
channel="${CARBIDE_CHANNEL:-release}"
version="${CARBIDE_VERSION:-}"

mkdir -p "$bin_dir"

need() {
  command -v "$1" >/dev/null 2>&1 || {
    printf 'install failed: %s is required\n' "$1" >&2
    exit 1
  }
}

detect_platform() {
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"
  case "$os" in
    linux|darwin) ;;
    *) return 1 ;;
  esac
  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *) return 1 ;;
  esac
  printf '%s_%s\n' "$os" "$arch"
}

latest_release_tag() {
  command -v curl >/dev/null 2>&1 || return 1
  effective="$(
    curl -fsSLI -o /dev/null -w '%{url_effective}' \
      "https://github.com/$repo_owner/$repo_name/releases/latest" 2>/dev/null || true
  )"
  tag="${effective##*/}"
  case "$tag" in
    v*) printf '%s\n' "$tag" ;;
    *) return 1 ;;
  esac
}

release_tag() {
  if [ "$channel" = "main" ]; then
    return 1
  fi
  if [ -n "$version" ]; then
    case "$version" in
      v*) printf '%s\n' "$version" ;;
      *) printf 'v%s\n' "$version" ;;
    esac
    return 0
  fi
  latest_release_tag
}

archive_url_for_ref() {
  ref="$1"
  case "$ref" in
    main) printf 'https://github.com/%s/%s/archive/refs/heads/main.tar.gz\n' "$repo_owner" "$repo_name" ;;
    v*) printf 'https://github.com/%s/%s/archive/refs/tags/%s.tar.gz\n' "$repo_owner" "$repo_name" "$ref" ;;
    *) printf 'https://github.com/%s/%s/archive/refs/heads/%s.tar.gz\n' "$repo_owner" "$repo_name" "$ref" ;;
  esac
}

install_source() {
  ref="$1"
  archive_url="$(archive_url_for_ref "$ref")"
  if command -v git >/dev/null 2>&1; then
    if [ -d "$install_dir/.git" ]; then
      git -C "$install_dir" fetch --quiet --tags origin
      git -C "$install_dir" checkout --quiet "$ref"
      git -C "$install_dir" pull --ff-only --quiet origin "$ref" 2>/dev/null || true
    else
      rm -rf "$install_dir"
      git clone --depth 1 --branch "$ref" "$repo_url" "$install_dir"
    fi
    return
  fi

  need curl
  need tar
  tmp_dir="$(mktemp -d)"
  trap 'rm -rf "$tmp_dir"' EXIT
  curl -fsSL "$archive_url" | tar -xz -C "$tmp_dir" --strip-components=1
  rm -rf "$install_dir"
  mkdir -p "$install_dir"
  cp -R "$tmp_dir/." "$install_dir/"
}

verify_checksum() {
  file="$1"
  checksum_file="$2"
  if [ ! -s "$checksum_file" ]; then
    return 0
  fi
  if command -v sha256sum >/dev/null 2>&1; then
    (cd "$(dirname "$file")" && sha256sum -c "$(basename "$checksum_file")")
    return
  fi
  if command -v shasum >/dev/null 2>&1; then
    (cd "$(dirname "$file")" && shasum -a 256 -c "$(basename "$checksum_file")")
  fi
}

install_release_binary() {
  tag="$1"
  platform="$(detect_platform)" || return 1
  asset="carbide_${platform}.tar.gz"
  base="https://github.com/$repo_owner/$repo_name/releases/download/$tag"
  tmp_dir="$(mktemp -d)"
  archive="$tmp_dir/$asset"
  checksum="$tmp_dir/$asset.sha256"

  command -v curl >/dev/null 2>&1 || return 1
  if ! curl -fsSL "$base/$asset" -o "$archive"; then
    rm -rf "$tmp_dir"
    return 1
  fi
  curl -fsSL "$base/$asset.sha256" -o "$checksum" 2>/dev/null || true
  verify_checksum "$archive" "$checksum"
  tar -xz -C "$tmp_dir" -f "$archive"
  test -x "$tmp_dir/carbide" || {
    rm -rf "$tmp_dir"
    return 1
  }

  build_dir="$install_dir/.cli/bin"
  mkdir -p "$build_dir"
  mv "$tmp_dir/carbide" "$build_dir/carbide"
  chmod +x "$build_dir/carbide"
  rm -rf "$tmp_dir"
  return 0
}

build_from_source() {
  command -v go >/dev/null 2>&1 || {
    printf 'install failed: release binary unavailable and Go is required for source build fallback\n' >&2
    exit 1
  }
  build_dir="$install_dir/.cli/bin"
  mkdir -p "$build_dir"
  commit=""
  if [ -d "$install_dir/.git" ] && command -v git >/dev/null 2>&1; then
    commit="$(git -C "$install_dir" rev-parse --short HEAD 2>/dev/null || true)"
  fi
  tmp_bin="$build_dir/carbide.$$"
  (
    cd "$install_dir/cli"
    go build -ldflags "-X github.com/ryangerardwilson/carbide/cli/internal/cli.commit=$commit" -o "$tmp_bin" ./cmd/carbide
  )
  mv "$tmp_bin" "$build_dir/carbide"
  chmod +x "$build_dir/carbide"
}

ref="main"
tag=""
if tag="$(release_tag 2>/dev/null)"; then
  ref="$tag"
fi

install_source "$ref"

if [ -n "$tag" ] && install_release_binary "$tag"; then
  :
else
  build_from_source
fi

ln -sfn "$install_dir/.cli/bin/carbide" "$bin_dir/carbide"

printf 'installed carbide to %s\n' "$bin_dir/carbide"
case ":$PATH:" in
  *":$bin_dir:"*) ;;
  *) printf 'add %s to PATH if carbide is not found\n' "$bin_dir" ;;
esac
