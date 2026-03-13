#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env.prod"
BACKEND_DIR="${ROOT_DIR}/backend-go"
FRONTEND_DIR="${ROOT_DIR}/frontend"
SERVICE_NAME="mdxy-backend"

WEB_ROOT="${WEB_ROOT:-/var/www/mdxy}"
DATA_DIR="${DATA_DIR:-/var/lib/mdxy}"
SERVER_PORT="${SERVER_PORT:-8080}"
SERVER_MODE="${SERVER_MODE:-release}"

if [[ -d "${ROOT_DIR}/notes" ]]; then
  NOTES_DIR_DEFAULT="${ROOT_DIR}/notes"
else
  NOTES_DIR_DEFAULT="${ROOT_DIR}/content/notes"
fi
ARTICLES_DIR_DEFAULT="${ROOT_DIR}/content/articles"

DEPLOY_USER="${DEPLOY_USER:-${SUDO_USER:-$(id -un)}}"
DEPLOY_GROUP="${DEPLOY_GROUP:-$(id -gn "${DEPLOY_USER}" 2>/dev/null || id -gn)}"
FRONTEND_NODE_OPTIONS="${FRONTEND_NODE_OPTIONS:---max-old-space-size=1024}"
FRONTEND_SKIP_TYPECHECK="${FRONTEND_SKIP_TYPECHECK:-0}"

log() {
  printf '[deploy] %s\n' "$*"
}

die() {
  printf '[deploy] error: %s\n' "$*" >&2
  exit 1
}

as_root() {
  if [[ "${EUID}" -eq 0 ]]; then
    "$@"
  else
    sudo "$@"
  fi
}

usage() {
  cat <<'EOF'
Usage:
  bash deploy/non-docker-https-deploy.sh [up|restart|status|logs|renew]

Required env (in .env.prod or exported):
  SITE_DOMAIN=example.com
  ACME_EMAIL=admin@example.com
  JWT_SECRET=long-random-secret

Optional env:
  WEB_ROOT=/var/www/mdxy
  DATA_DIR=/var/lib/mdxy
  SERVER_PORT=8080
  SERVER_MODE=release
  CONTENT_NOTES_DIR=/path/to/notes
  CONTENT_ARTICLES_DIR=/path/to/articles
  DEPLOY_USER=ubuntu
  FRONTEND_NODE_OPTIONS=--max-old-space-size=1024
  FRONTEND_SKIP_TYPECHECK=0
EOF
}

normalize_container_path() {
  local key="$1"
  local value="$2"
  local mapped="${value}"

  if [[ "${value}" != /app/* ]]; then
    printf '%s' "${value}"
    return
  fi

  case "${key}" in
    DATABASE_PATH)
      mapped="${DATA_DIR}/mdxy.db"
      ;;
    CONTENT_NOTES_DIR)
      mapped="${NOTES_DIR_DEFAULT}"
      ;;
    CONTENT_ARTICLES_DIR)
      mapped="${ARTICLES_DIR_DEFAULT}"
      ;;
  esac

  if [[ "${mapped}" != "${value}" ]]; then
    log "检测到 ${key} 使用容器路径 ${value}，自动转换为 ${mapped}"
  fi

  printf '%s' "${mapped}"
}

install_runtime_deps() {
  local need_install=0
  command -v nginx >/dev/null 2>&1 || need_install=1
  command -v certbot >/dev/null 2>&1 || need_install=1

  if [[ "${need_install}" -eq 0 ]]; then
    return
  fi

  if command -v apt-get >/dev/null 2>&1; then
    as_root apt-get update -y
    as_root apt-get install -y nginx certbot python3-certbot-nginx
    return
  fi

  if command -v dnf >/dev/null 2>&1; then
    as_root dnf install -y nginx certbot python3-certbot-nginx || \
      as_root dnf install -y nginx certbot certbot-nginx
    return
  fi

  if command -v yum >/dev/null 2>&1; then
    as_root yum install -y epel-release || true
    as_root yum install -y nginx certbot python3-certbot-nginx || \
      as_root yum install -y nginx certbot certbot-nginx
    return
  fi

  die "无法自动安装 nginx/certbot，请手动安装后重试。"
}

require_commands() {
  local missing=()
  command -v go >/dev/null 2>&1 || missing+=("go")
  command -v npm >/dev/null 2>&1 || missing+=("npm")
  command -v systemctl >/dev/null 2>&1 || missing+=("systemctl")
  command -v curl >/dev/null 2>&1 || missing+=("curl")
  if [[ "${#missing[@]}" -gt 0 ]]; then
    die "缺少命令: ${missing[*]}"
  fi
}

load_env() {
  if [[ -f "${ENV_FILE}" ]]; then
    # shellcheck disable=SC1090
    source "${ENV_FILE}"
  fi

  : "${SITE_DOMAIN:?请设置 SITE_DOMAIN（.env.prod 或环境变量）}"
  : "${ACME_EMAIL:?请设置 ACME_EMAIL（.env.prod 或环境变量）}"
  : "${JWT_SECRET:?请设置 JWT_SECRET（.env.prod 或环境变量）}"

  DATABASE_PATH="${DATABASE_PATH:-${DATA_DIR}/mdxy.db}"
  CONTENT_NOTES_DIR="${CONTENT_NOTES_DIR:-${NOTES_DIR_DEFAULT}}"
  CONTENT_ARTICLES_DIR="${CONTENT_ARTICLES_DIR:-${ARTICLES_DIR_DEFAULT}}"
  DATABASE_PATH="$(normalize_container_path DATABASE_PATH "${DATABASE_PATH}")"
  CONTENT_NOTES_DIR="$(normalize_container_path CONTENT_NOTES_DIR "${CONTENT_NOTES_DIR}")"
  CONTENT_ARTICLES_DIR="$(normalize_container_path CONTENT_ARTICLES_DIR "${CONTENT_ARTICLES_DIR}")"

  if [[ ! -d "${CONTENT_NOTES_DIR}" ]]; then
    die "笔记目录不存在: ${CONTENT_NOTES_DIR}"
  fi
}

prepare_dirs() {
  as_root mkdir -p "${WEB_ROOT}" "${DATA_DIR}" "${CONTENT_ARTICLES_DIR}"
  as_root chmod 755 "${WEB_ROOT}" "${CONTENT_ARTICLES_DIR}"
  as_root chmod 750 "${DATA_DIR}" || true
}

build_backend() {
  log "构建后端"
  (
    cd "${BACKEND_DIR}"
    go mod download
    CGO_ENABLED=1 go build -o server ./cmd/server/
  )
}

build_frontend() {
  log "构建前端"
  (
    cd "${FRONTEND_DIR}"
    npm ci

    if [[ "${FRONTEND_SKIP_TYPECHECK}" == "1" ]]; then
      log "按配置跳过 vue-tsc，直接执行 vite build"
      NODE_OPTIONS="${FRONTEND_NODE_OPTIONS}" npx vite build
      return
    fi

    local build_log
    build_log="$(mktemp)"
    trap 'rm -f "${build_log}"' RETURN

    if NODE_OPTIONS="${FRONTEND_NODE_OPTIONS}" npm run build 2>&1 | tee "${build_log}"; then
      return
    fi

    if grep -qiE "killed|out of memory|heap out of memory" "${build_log}"; then
      log "检测到前端构建内存不足，自动降级为 vite build（跳过 vue-tsc）"
      NODE_OPTIONS="${FRONTEND_NODE_OPTIONS}" npx vite build
      return
    fi

    die "前端构建失败（非内存问题）。可设置 FRONTEND_SKIP_TYPECHECK=1 后重试。"
  )
}

install_frontend_assets() {
  log "发布前端静态文件到 ${WEB_ROOT}"
  as_root rm -rf "${WEB_ROOT:?}/"*
  as_root cp -a "${FRONTEND_DIR}/dist/." "${WEB_ROOT}/"
}

write_backend_env() {
  as_root mkdir -p /etc/mdxy
  as_root tee /etc/mdxy/backend.env >/dev/null <<EOF
SERVER_PORT=${SERVER_PORT}
SERVER_MODE=${SERVER_MODE}
DATABASE_PATH=${DATABASE_PATH}
JWT_SECRET=${JWT_SECRET}
CONTENT_NOTES_DIR=${CONTENT_NOTES_DIR}
CONTENT_ARTICLES_DIR=${CONTENT_ARTICLES_DIR}
EOF
  as_root chmod 640 /etc/mdxy/backend.env
}

write_systemd_service() {
  log "写入 systemd 服务"
  as_root tee "/etc/systemd/system/${SERVICE_NAME}.service" >/dev/null <<EOF
[Unit]
Description=MDXY Go Backend
After=network.target

[Service]
Type=simple
User=${DEPLOY_USER}
Group=${DEPLOY_GROUP}
WorkingDirectory=${BACKEND_DIR}
EnvironmentFile=/etc/mdxy/backend.env
ExecStart=${BACKEND_DIR}/server
Restart=on-failure
RestartSec=3
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

  as_root systemctl daemon-reload
  as_root systemctl enable --now "${SERVICE_NAME}"
}

write_nginx_config() {
  log "写入 Nginx 配置"
  local server_names="${SITE_DOMAIN}"
  if [[ "${SITE_DOMAIN}" != www.* ]]; then
    server_names="${server_names} www.${SITE_DOMAIN}"
  fi

  as_root tee /etc/nginx/conf.d/mdxy.conf >/dev/null <<EOF
server {
    listen 80;
    listen [::]:80;
    server_name ${server_names};

    client_max_body_size 10M;

    location / {
        root ${WEB_ROOT};
        index index.html;
        try_files \$uri \$uri/ /index.html;
    }

    location /api {
        proxy_pass http://127.0.0.1:${SERVER_PORT};
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
}
EOF

  as_root nginx -t
  as_root systemctl enable --now nginx
  as_root systemctl restart nginx
}

enable_https() {
  log "申请/更新 HTTPS 证书"
  local certbot_args=(
    --nginx
    --non-interactive
    --agree-tos
    --email "${ACME_EMAIL}"
    --redirect
    -d "${SITE_DOMAIN}"
  )
  if [[ "${SITE_DOMAIN}" != www.* ]]; then
    certbot_args+=(-d "www.${SITE_DOMAIN}")
  fi

  as_root certbot "${certbot_args[@]}"
  as_root systemctl reload nginx
}

verify_services() {
  log "检查后端健康状态"
  curl -fsS "http://127.0.0.1:${SERVER_PORT}/health" >/dev/null
  log "后端健康检查通过"
}

up() {
  require_commands
  install_runtime_deps
  load_env
  prepare_dirs
  build_backend
  build_frontend
  install_frontend_assets
  write_backend_env
  write_systemd_service
  write_nginx_config
  enable_https
  verify_services

  log "部署完成："
  log "  前台: https://${SITE_DOMAIN}"
  log "  后台: https://${SITE_DOMAIN}/admin/login"
  log "请尽快修改默认管理员密码。"
}

restart() {
  as_root systemctl restart "${SERVICE_NAME}"
  as_root systemctl restart nginx
  verify_services
  log "重启完成"
}

status() {
  as_root systemctl --no-pager status "${SERVICE_NAME}" || true
  as_root systemctl --no-pager status nginx || true
}

logs() {
  as_root journalctl -u "${SERVICE_NAME}" -n 200 -f
}

renew() {
  as_root certbot renew
  as_root systemctl reload nginx
  log "证书续期命令执行完成"
}

main() {
  local cmd="${1:-up}"
  case "${cmd}" in
    up) up ;;
    restart) restart ;;
    status) status ;;
    logs) logs ;;
    renew) renew ;;
    -h|--help|help) usage ;;
    *)
      usage
      exit 2
      ;;
  esac
}

main "$@"
