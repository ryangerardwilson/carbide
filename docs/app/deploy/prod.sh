#!/usr/bin/env bash
set -euo pipefail

require_env() {
  local name="$1"
  if [ -z "${!name:-}" ]; then
    printf '%s is required\n' "$name" >&2
    exit 1
  fi
}

require_tool() {
  local name="$1"
  if ! command -v "$name" >/dev/null 2>&1; then
    printf '%s is required\n' "$name" >&2
    exit 1
  fi
}

require_env CARBIDE_DOCS_DEPLOY_SSH
require_env CARBIDE_DOCS_POSTGRES_PASSWORD
require_tool ssh
require_tool rsync

project_root="${CARBIDE_PROJECT_ROOT:-$(pwd)}"
docs_root="$(cd "$project_root/.." && pwd)"

remote_root="${CARBIDE_DOCS_DEPLOY_PATH:-/opt/carbide/docs}"
public_url="${CARBIDE_DOCS_PUBLIC_URL:-https://carbide.ryangerardwilson.com}"
public_app_name="${CARBIDE_DOCS_PUBLIC_APP_NAME:-Carbide Docs}"
http_port="${CARBIDE_DOCS_HTTP_PORT:-18081}"
domain="${CARBIDE_DOCS_DOMAIN:-carbide.ryangerardwilson.com}"
nginx_site="${CARBIDE_DOCS_NGINX_SITE:-carbide}"
manage_nginx="${CARBIDE_DOCS_MANAGE_NGINX:-1}"
ssh_target="${CARBIDE_DOCS_DEPLOY_SSH}"
compose_project_name="${CARBIDE_DOCS_COMPOSE_PROJECT_NAME:-carbide-docs}"
legacy_project_name="${CARBIDE_DOCS_LEGACY_PROJECT_NAME:-app}"

printf 'syncing docs to %s\n' "$ssh_target"

ssh "$ssh_target" "mkdir -p '$remote_root'"

rsync -az --delete \
  --exclude .git \
  --exclude .carbide \
  --exclude app/.env \
  --exclude app/web/node_modules \
  --exclude app/web/public \
  "$docs_root"/ "$ssh_target:$remote_root/"

ssh "$ssh_target" bash -s -- \
  "$remote_root" \
  "$http_port" \
  "$public_url" \
  "$public_app_name" \
  "$CARBIDE_DOCS_POSTGRES_PASSWORD" \
  "$domain" \
  "$nginx_site" \
  "$manage_nginx" \
  "$compose_project_name" \
  "$legacy_project_name" <<'EOF'
set -euo pipefail

remote_root="$1"
http_port="$2"
public_url="$3"
public_app_name="$4"
postgres_password="$5"
domain="$6"
nginx_site="$7"
manage_nginx="$8"
compose_project_name="$9"
legacy_project_name="${10}"
export COMPOSE_PROJECT_NAME="$compose_project_name"
compose_cmd="docker compose --env-file app/.env -f app/docker-compose.yml --project-directory app"

cat > "$remote_root/app/.env" <<ENV
APP_ENV=production
CARBIDE_HTTP_PORT=$http_port
PUBLIC_APP_NAME=$public_app_name
PUBLIC_URL=$public_url
POSTGRES_PASSWORD=$postgres_password
DATABASE_URL=postgres://carbide:$postgres_password@db:5432/carbide?sslmode=disable
ENV

cd "$remote_root"
$compose_cmd config >/dev/null
$compose_cmd up -d --build --remove-orphans

if [ "$legacy_project_name" != "$compose_project_name" ]; then
  COMPOSE_PROJECT_NAME="$legacy_project_name" \
    docker compose --env-file app/.env -f app/docker-compose.yml --project-directory app \
    down --remove-orphans >/dev/null 2>&1 || true
fi

if [ "$manage_nginx" = "1" ]; then
  sudo -n true
  config_path="/etc/nginx/sites-available/$nginx_site"
  enabled_path="/etc/nginx/sites-enabled/$nginx_site"
  cat <<NGINX | sudo tee "$config_path" >/dev/null
server {
    listen 80;
    listen [::]:80;
    server_name $domain;

    location / {
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_pass http://127.0.0.1:$http_port;
    }
}
NGINX
  sudo ln -sfn "$config_path" "$enabled_path"
  sudo nginx -t
  sudo systemctl reload nginx
fi

for _ in $(seq 1 40); do
  if curl -fsS --max-time 5 "http://127.0.0.1:$http_port/health" >/dev/null; then
    exit 0
  fi
  sleep 2
done

curl -fsS --max-time 10 "http://127.0.0.1:$http_port/health" >/dev/null
EOF

printf 'deployed %s\n' "$public_url"
