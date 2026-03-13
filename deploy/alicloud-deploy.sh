#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env.prod"
COMPOSE_FILE_BASE="${ROOT_DIR}/docker-compose.yml"
COMPOSE_FILE_PROD="${ROOT_DIR}/docker-compose.prod.yml"

log() {
  printf '[deploy] %s\n' "$*"
}

die() {
  printf '[deploy] error: %s\n' "$*" >&2
  exit 1
}

usage() {
  cat <<'EOF'
Usage:
  bash deploy/alicloud-deploy.sh [up|restart|down|logs|status|config]

Commands:
  up       Build and start production stack (default)
  restart  Recreate and restart production stack
  down     Stop production stack
  logs     Follow logs of caddy/frontend/backend
  status   Show compose service status
  config   Validate merged compose config
EOF
}

detect_compose() {
  if docker compose version >/dev/null 2>&1; then
    COMPOSE_CMD=(docker compose)
    return
  fi
  if command -v docker-compose >/dev/null 2>&1; then
    COMPOSE_CMD=(docker-compose)
    return
  fi
  die "未检测到 docker compose（或 docker-compose）"
}

compose() {
  "${COMPOSE_CMD[@]}" \
    --env-file "${ENV_FILE}" \
    -f "${COMPOSE_FILE_BASE}" \
    -f "${COMPOSE_FILE_PROD}" \
    "$@"
}

require_env() {
  [ -f "${ENV_FILE}" ] || die "缺少 ${ENV_FILE}。请先执行: cp .env.prod.example .env.prod 并填写变量"
  # shellcheck disable=SC1090
  source "${ENV_FILE}"

  [ -n "${SITE_DOMAIN:-}" ] || die ".env.prod 缺少 SITE_DOMAIN"
  [ -n "${ACME_EMAIL:-}" ] || die ".env.prod 缺少 ACME_EMAIL"
  [ -n "${JWT_SECRET:-}" ] || die ".env.prod 缺少 JWT_SECRET"
}

main() {
  local cmd="${1:-up}"
  detect_compose
  require_env

  case "${cmd}" in
    up)
      log "启动生产环境（Caddy + frontend + backend）"
      compose up -d --build
      compose ps
      log "已启动。请确认域名 ${SITE_DOMAIN} 已解析并开放 80/443 安全组端口。"
      ;;
    restart)
      log "重启生产环境"
      compose up -d --build --force-recreate
      compose ps
      ;;
    down)
      log "停止生产环境"
      compose down
      ;;
    logs)
      compose logs -f caddy frontend backend
      ;;
    status)
      compose ps
      ;;
    config)
      compose config >/dev/null
      log "Compose 配置校验通过"
      ;;
    -h|--help|help)
      usage
      ;;
    *)
      usage
      exit 2
      ;;
  esac
}

main "$@"
