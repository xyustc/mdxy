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
FRONTEND_NPM_CI="${FRONTEND_NPM_CI:-auto}"
ENABLE_WWW="${ENABLE_WWW:-0}"
CACHE_ENABLED="${CACHE_ENABLED:-1}"
CACHE_DIR="${CACHE_DIR:-${ROOT_DIR}/.deploy-cache}"
FORCE_REBUILD="${FORCE_REBUILD:-0}"
STOP_NGINX_ON_STOP="${STOP_NGINX_ON_STOP:-0}"

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
  bash deploy/non-docker-https-deploy.sh [up|build|start|stop|restart|status|health|logs|renew]

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
  FRONTEND_NPM_CI=auto   # auto|always|never
  ENABLE_WWW=0          # 1 时申请并支持 www 子域名
  CACHE_ENABLED=1       # 1 时启用构建缓存
  CACHE_DIR=.deploy-cache
  FORCE_REBUILD=0       # 1 时强制重建后端和前端
  STOP_NGINX_ON_STOP=0  # 1 时 stop 命令会一并停止 nginx
EOF
}

is_true() {
  case "${1:-}" in
    1|true|TRUE|yes|YES|on|ON)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

cache_file_path() {
  local key="$1"
  printf '%s/%s' "${CACHE_DIR}" "${key}"
}

hash_stdin() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum | awk '{print $1}'
    return
  fi
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 | awk '{print $1}'
    return
  fi
  cksum | awk '{print $1}'
}

fingerprint_file_set() {
  if [[ "$#" -eq 0 ]]; then
    printf 'empty'
    return
  fi

  {
    local file rel
    for file in "$@"; do
      [[ -f "${file}" ]] || continue
      rel="${file#${ROOT_DIR}/}"
      printf 'FILE:%s\n' "${rel}"
      cat "${file}"
      printf '\n'
    done
  } | hash_stdin
}

ensure_cache_dir() {
  if ! is_true "${CACHE_ENABLED}"; then
    return
  fi
  mkdir -p "${CACHE_DIR}"
}

read_cache_value() {
  if ! is_true "${CACHE_ENABLED}"; then
    return
  fi
  local key="$1"
  local path
  path="$(cache_file_path "${key}")"
  if [[ -f "${path}" ]]; then
    cat "${path}"
  fi
}

write_cache_value() {
  if ! is_true "${CACHE_ENABLED}"; then
    return
  fi
  local key="$1"
  local value="$2"
  ensure_cache_dir
  printf '%s' "${value}" > "$(cache_file_path "${key}")"
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

install_rpm_nginx_certbot() {
  local pkg_manager="$1"
  if as_root "${pkg_manager}" install -y nginx certbot python3-certbot-nginx; then
    return
  fi
  as_root "${pkg_manager}" install -y nginx certbot certbot-nginx
}

frontend_lock_fingerprint() {
  local -a files=()
  [[ -f "${FRONTEND_DIR}/package.json" ]] && files+=("${FRONTEND_DIR}/package.json")
  [[ -f "${FRONTEND_DIR}/package-lock.json" ]] && files+=("${FRONTEND_DIR}/package-lock.json")
  fingerprint_file_set "${files[@]}"
}

backend_build_fingerprint() {
  local -a files=()
  mapfile -t files < <(
    find "${BACKEND_DIR}" -type f \
      \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' -o -name '*.yaml' \) \
      ! -path "${BACKEND_DIR}/data/*" \
      ! -path "${BACKEND_DIR}/bin/*" \
      ! -name 'server' \
      | LC_ALL=C sort
  )
  fingerprint_file_set "${files[@]}"
}

frontend_build_fingerprint() {
  local -a files=()
  mapfile -t files < <(
    {
      find "${FRONTEND_DIR}/src" -type f 2>/dev/null
      find "${FRONTEND_DIR}" -maxdepth 1 -type f \
        \( -name 'package.json' -o -name 'package-lock.json' -o -name 'vite.config.ts' -o -name 'index.html' -o -name 'tsconfig*.json' \)
    } | LC_ALL=C sort
  )
  fingerprint_file_set "${files[@]}"
}

should_run_frontend_install() {
  local lock_fp cached_lock_fp
  case "${FRONTEND_NPM_CI}" in
    auto)
      if [[ ! -d node_modules ]]; then
        return 0
      fi
      if ! is_true "${CACHE_ENABLED}" || is_true "${FORCE_REBUILD}"; then
        return 0
      fi
      lock_fp="$(frontend_lock_fingerprint)"
      cached_lock_fp="$(read_cache_value frontend-npm-lock.hash)"
      if [[ -z "${cached_lock_fp}" || "${lock_fp}" != "${cached_lock_fp}" ]]; then
        return 0
      fi
      return 1
      ;;
    always|1|true)
      return 0
      ;;
    never|0|false)
      return 1
      ;;
    *)
      die "FRONTEND_NPM_CI 仅支持 auto|always|never（或 1/0/true/false）"
      ;;
  esac
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
    install_rpm_nginx_certbot dnf
    return
  fi

  if command -v yum >/dev/null 2>&1; then
    as_root yum install -y epel-release || true
    install_rpm_nginx_certbot yum
    return
  fi

  die "无法自动安装 nginx/certbot，请手动安装后重试。"
}

require_commands() {
  local missing=()
  local cmd
  for cmd in "$@"; do
    command -v "${cmd}" >/dev/null 2>&1 || missing+=("${cmd}")
  done
  if [[ "${#missing[@]}" -gt 0 ]]; then
    die "缺少命令: ${missing[*]}"
  fi
}

require_build_commands() {
  require_commands go npm
}

require_systemd_commands() {
  require_commands systemctl
}

require_health_commands() {
  require_commands curl
}

require_certbot_commands() {
  require_commands certbot
}

validate_switches() {
  case "${ENABLE_WWW}" in
    0|1|true|false|TRUE|FALSE|yes|no|YES|NO|on|off|ON|OFF)
      ;;
    *)
      die "ENABLE_WWW 仅支持 0/1/true/false/yes/no/on/off"
      ;;
  esac

  case "${CACHE_ENABLED}" in
    0|1|true|false|TRUE|FALSE|yes|no|YES|NO|on|off|ON|OFF)
      ;;
    *)
      die "CACHE_ENABLED 仅支持 0/1/true/false/yes/no/on/off"
      ;;
  esac

  case "${FORCE_REBUILD}" in
    0|1|true|false|TRUE|FALSE|yes|no|YES|NO|on|off|ON|OFF)
      ;;
    *)
      die "FORCE_REBUILD 仅支持 0/1/true/false/yes/no/on/off"
      ;;
  esac

  case "${STOP_NGINX_ON_STOP}" in
    0|1|true|false|TRUE|FALSE|yes|no|YES|NO|on|off|ON|OFF)
      ;;
    *)
      die "STOP_NGINX_ON_STOP 仅支持 0/1/true/false/yes/no/on/off"
      ;;
  esac
}

normalize_runtime_paths() {
  if [[ "${CACHE_DIR}" != /* ]]; then
    CACHE_DIR="${ROOT_DIR}/${CACHE_DIR}"
  fi
}

load_env() {
  if [[ -f "${ENV_FILE}" ]]; then
    # shellcheck disable=SC1090
    source "${ENV_FILE}"
  fi
  validate_switches
  normalize_runtime_paths
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

load_optional_env() {
  if [[ -f "${ENV_FILE}" ]]; then
    # shellcheck disable=SC1090
    source "${ENV_FILE}"
  fi
  validate_switches
  normalize_runtime_paths
}

prepare_dirs() {
  as_root mkdir -p "${WEB_ROOT}" "${DATA_DIR}" "${CONTENT_ARTICLES_DIR}"
  as_root chmod 755 "${WEB_ROOT}" "${CONTENT_ARTICLES_DIR}"
  as_root chmod 750 "${DATA_DIR}" || true
}

build_backend() {
  log "构建后端"
  local source_fp cached_fp
  if is_true "${CACHE_ENABLED}" && ! is_true "${FORCE_REBUILD}"; then
    source_fp="$(backend_build_fingerprint)"
    cached_fp="$(read_cache_value backend-build.hash)"
    if [[ -f "${BACKEND_DIR}/server" && -n "${cached_fp}" && "${source_fp}" == "${cached_fp}" ]]; then
      log "后端源码未变化，跳过构建（缓存命中）"
      return
    fi
  fi

  (
    cd "${BACKEND_DIR}"
    go mod download
    CGO_ENABLED=1 go build -o server ./cmd/server/
  )

  if is_true "${CACHE_ENABLED}"; then
    if [[ -z "${source_fp:-}" ]]; then
      source_fp="$(backend_build_fingerprint)"
    fi
    write_cache_value backend-build.hash "${source_fp}"
  fi
}

build_frontend() {
  log "构建前端"
  local source_fp cached_fp lock_fp
  (
    cd "${FRONTEND_DIR}"
    if should_run_frontend_install; then
      npm ci
      if is_true "${CACHE_ENABLED}"; then
        lock_fp="$(frontend_lock_fingerprint)"
        write_cache_value frontend-npm-lock.hash "${lock_fp}"
      fi
    else
      log "跳过 npm ci（FRONTEND_NPM_CI=${FRONTEND_NPM_CI}）"
    fi

    if is_true "${CACHE_ENABLED}" && ! is_true "${FORCE_REBUILD}"; then
      source_fp="$(frontend_build_fingerprint)"
      cached_fp="$(read_cache_value frontend-build.hash)"
      if [[ -f dist/index.html && -n "${cached_fp}" && "${source_fp}" == "${cached_fp}" ]]; then
        log "前端源码未变化，跳过构建（缓存命中）"
        return
      fi
    fi

    if [[ "${FRONTEND_SKIP_TYPECHECK}" == "1" ]]; then
      log "按配置跳过 vue-tsc，直接执行 vite build"
      if ! NODE_OPTIONS="${FRONTEND_NODE_OPTIONS}" npx vite build; then
        die "vite build 失败。若是内存不足，请增加 swap 或在本地/CI 构建后上传 dist。"
      fi
    else
      local build_log build_rc
      build_log="$(mktemp)"
      trap 'rm -f "${build_log}"' RETURN

      set +e
      NODE_OPTIONS="${FRONTEND_NODE_OPTIONS}" npm run build 2>&1 | tee "${build_log}"
      build_rc=${PIPESTATUS[0]}
      set -e

      if [[ "${build_rc}" -ne 0 ]]; then
        if [[ "${build_rc}" -eq 137 || "${build_rc}" -eq 143 ]] || grep -qiE "killed|out of memory|heap out of memory" "${build_log}"; then
          log "检测到前端构建内存不足，自动降级为 vite build（跳过 vue-tsc）"
          if ! NODE_OPTIONS="${FRONTEND_NODE_OPTIONS}" npx vite build; then
            die "前端构建内存不足，且 vite build 仍失败。建议增加 swap 或在本地/CI 构建后上传 dist。"
          fi
        else
          die "前端构建失败（非内存问题）。可设置 FRONTEND_SKIP_TYPECHECK=1 后重试。"
        fi
      fi
    fi

    if is_true "${CACHE_ENABLED}"; then
      if [[ -z "${source_fp}" ]]; then
        source_fp="$(frontend_build_fingerprint)"
      fi
      write_cache_value frontend-build.hash "${source_fp}"
    fi
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
  if is_true "${ENABLE_WWW}" && [[ "${SITE_DOMAIN}" != www.* ]]; then
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
  as_root systemctl enable nginx
  if as_root systemctl is-active --quiet nginx; then
    as_root systemctl reload nginx
  else
    as_root systemctl start nginx
  fi
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
  if is_true "${ENABLE_WWW}" && [[ "${SITE_DOMAIN}" != www.* ]]; then
    certbot_args+=(-d "www.${SITE_DOMAIN}")
  fi

  as_root certbot "${certbot_args[@]}"
  as_root systemctl reload nginx
}

verify_services() {
  log "检查后端健康状态"
  require_health_commands
  curl -fsS "http://127.0.0.1:${SERVER_PORT}/health" >/dev/null
  log "后端健康检查通过"
}

build() {
  load_optional_env
  require_build_commands
  ensure_cache_dir
  build_backend
  build_frontend
  log "构建完成"
}

up() {
  load_env
  require_build_commands
  require_systemd_commands
  require_health_commands
  install_runtime_deps
  ensure_cache_dir
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

start() {
  load_optional_env
  require_systemd_commands
  as_root systemctl start "${SERVICE_NAME}"
  as_root systemctl start nginx || true
  verify_services
  log "服务启动完成"
}

stop() {
  load_optional_env
  require_systemd_commands
  as_root systemctl stop "${SERVICE_NAME}" || true
  if is_true "${STOP_NGINX_ON_STOP}"; then
    as_root systemctl stop nginx || true
  fi
  log "服务停止完成（backend）"
}

restart() {
  load_optional_env
  require_systemd_commands
  as_root systemctl restart "${SERVICE_NAME}"
  if as_root systemctl is-active --quiet nginx; then
    as_root systemctl reload nginx
  else
    as_root systemctl start nginx
  fi
  verify_services
  log "重启完成"
}

status() {
  load_optional_env
  require_systemd_commands
  as_root systemctl --no-pager status "${SERVICE_NAME}" || true
  as_root systemctl --no-pager status nginx || true
}

health() {
  load_optional_env
  verify_services
}

logs() {
  load_optional_env
  require_systemd_commands
  as_root journalctl -u "${SERVICE_NAME}" -n 200 -f
}

renew() {
  load_optional_env
  require_certbot_commands
  require_systemd_commands
  as_root certbot renew
  as_root systemctl reload nginx
  log "证书续期命令执行完成"
}

main() {
  local cmd="${1:-up}"
  case "${cmd}" in
    up) up ;;
    build) build ;;
    start) start ;;
    stop) stop ;;
    restart) restart ;;
    status) status ;;
    health) health ;;
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
