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
legacy_project_names="${CARBIDE_DOCS_LEGACY_PROJECT_NAMES:-app 1}"

printf 'syncing docs to %s\n' "$ssh_target"

ssh "$ssh_target" "mkdir -p '$remote_root'"

rsync -az --delete \
  --exclude .git \
  --exclude .carbide \
  --exclude app/.env \
  --exclude app/web/node_modules \
  --exclude app/web/public \
  "$docs_root"/ "$ssh_target:$remote_root/"

printf -v remote_root_q '%q' "$remote_root"
printf -v http_port_q '%q' "$http_port"
printf -v public_url_q '%q' "$public_url"
printf -v public_app_name_q '%q' "$public_app_name"
printf -v postgres_password_q '%q' "$CARBIDE_DOCS_POSTGRES_PASSWORD"
printf -v domain_q '%q' "$domain"
printf -v nginx_site_q '%q' "$nginx_site"
printf -v manage_nginx_q '%q' "$manage_nginx"
printf -v compose_project_name_q '%q' "$compose_project_name"
printf -v legacy_project_names_q '%q' "$legacy_project_names"

remote_script="$(cat <<'EOF'
set -euo pipefail

compose_cmd() {
  docker compose \
    -p "$compose_project_name" \
    --env-file app/.env \
    -f app/docker-compose.yml \
    --project-directory app \
    "$@"
}

legacy_compose_down() {
  local project_name="$1"
  docker compose \
    -p "$project_name" \
    --env-file app/.env \
    -f app/docker-compose.yml \
    --project-directory app \
    down --remove-orphans
}

cat > "$remote_root/app/.env" <<ENV
APP_ENV=production
CARBIDE_HTTP_PORT=$http_port
PUBLIC_APP_NAME=$public_app_name
PUBLIC_URL=$public_url
POSTGRES_PASSWORD=$postgres_password
DATABASE_URL=postgres://carbide:$postgres_password@db:5432/carbide?sslmode=disable
ENV

cd "$remote_root"
compose_cmd config >/dev/null
compose_cmd up -d --build --remove-orphans

for legacy_project_name in $legacy_project_names; do
  if [ -n "$legacy_project_name" ] && [ "$legacy_project_name" != "$compose_project_name" ]; then
    legacy_compose_down "$legacy_project_name" >/dev/null 2>&1 || true
  fi
done

if [ "$manage_nginx" = "1" ]; then
  sudo -n true
  config_path="/etc/nginx/sites-available/$nginx_site"
  enabled_path="/etc/nginx/sites-enabled/$nginx_site"
  cert_dir="/etc/letsencrypt/live/$domain"
  fullchain_path="$cert_dir/fullchain.pem"
  privkey_path="$cert_dir/privkey.pem"
  options_path="/etc/letsencrypt/options-ssl-nginx.conf"
  dhparams_path="/etc/letsencrypt/ssl-dhparams.pem"
  if sudo test -f "$fullchain_path" && sudo test -f "$privkey_path" && sudo test -f "$options_path" && sudo test -f "$dhparams_path"; then
    cat <<NGINX | sudo tee "$config_path" >/dev/null
server {
    listen 80;
    listen [::]:80;
    server_name $domain;

    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name $domain;

    ssl_certificate $fullchain_path;
    ssl_certificate_key $privkey_path;
    include $options_path;
    ssl_dhparam $dhparams_path;

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
  else
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
  fi
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
)"

ssh "$ssh_target" "bash -s" <<EOF
remote_root=$remote_root_q
http_port=$http_port_q
public_url=$public_url_q
public_app_name=$public_app_name_q
postgres_password=$postgres_password_q
domain=$domain_q
nginx_site=$nginx_site_q
manage_nginx=$manage_nginx_q
compose_project_name=$compose_project_name_q
legacy_project_names=$legacy_project_names_q
$remote_script
EOF

printf 'deployed %s\n' "$public_url"
