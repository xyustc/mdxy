#!/usr/bin/env bash
set -euo pipefail

CUTOVER_DIR=""

log() {
  printf '[rollback] %s\n' "$*"
}

die() {
  printf '[rollback] error: %s\n' "$*" >&2
  exit 1
}

usage() {
  cat <<'EOF'
Usage:
  bash backend-go/scripts/rollback-go-cutover.sh --cutover-dir /path/to/cutover-YYYYmmdd-HHMMSS
EOF
}

find_compose() {
  if docker compose version >/dev/null 2>&1; then
    COMPOSE_CMD=(docker compose)
    return
  fi

  if command -v docker-compose >/dev/null 2>&1; then
    COMPOSE_CMD=(docker-compose)
    return
  fi

  die "未找到 docker compose 或 docker-compose"
}

decode_b64() {
  local value=$1

  if base64 --help 2>&1 | grep -q -- '--decode'; then
    printf '%s' "${value}" | base64 --decode
    return
  fi

  printf '%s' "${value}" | base64 -D
}

load_rollback_var() {
  local key=$1
  local legacy_key=${2:-}
  local line
  local raw

  line=$(grep -E "^${key}=" "${CUTOVER_DIR}/rollback.env" | head -n 1 || true)
  if [ -z "${line}" ] && [ -n "${legacy_key}" ]; then
    line=$(grep -E "^${legacy_key}=" "${CUTOVER_DIR}/rollback.env" | head -n 1 || true)
    [ -n "${line}" ] || die "rollback.env 中缺少 ${key} / ${legacy_key}"

    raw=${line#*=}
    [ -n "${raw}" ] || die "rollback.env 中 ${legacy_key} 为空"
    printf '%s' "${raw}"
    return
  fi

  [ -n "${line}" ] || die "rollback.env 中缺少 ${key}"

  raw=${line#*=}
  [ -n "${raw}" ] || die "rollback.env 中 ${key} 为空"

  decode_b64 "${raw}" 2>/dev/null || die "rollback.env 中 ${key} 解析失败"
}

wait_for_health() {
  local attempt

  [ -n "${HEALTH_URL:-}" ] || return 0
  command -v curl >/dev/null 2>&1 || die "缺少命令: curl"

  for attempt in $(seq 1 30); do
    if curl -fsS "${HEALTH_URL}" >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done

  return 1
}

while [ $# -gt 0 ]; do
  case "$1" in
    --cutover-dir)
      CUTOVER_DIR=${2:-}
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      die "未知参数: $1"
      ;;
  esac
done

[ -n "${CUTOVER_DIR}" ] || {
  usage
  exit 2
}

[ -f "${CUTOVER_DIR}/rollback.env" ] || die "未找到 rollback.env: ${CUTOVER_DIR}/rollback.env"

TARGET_DB=$(load_rollback_var "TARGET_DB_B64" "TARGET_DB")
PRE_CUTOVER_GO_DB=$(load_rollback_var "PRE_CUTOVER_GO_DB_B64" "PRE_CUTOVER_GO_DB")
REPO_ROOT=$(load_rollback_var "REPO_ROOT_B64" "REPO_ROOT")
COMPOSE_SERVICE=$(load_rollback_var "COMPOSE_SERVICE_B64" "COMPOSE_SERVICE")
HEALTH_URL=$(load_rollback_var "HEALTH_URL_B64" "HEALTH_URL")

[ -f "${PRE_CUTOVER_GO_DB}" ] || die "切换前 Go 数据库快照不存在: ${PRE_CUTOVER_GO_DB}"
mkdir -p "$(dirname "${TARGET_DB}")"
cp "${PRE_CUTOVER_GO_DB}" "${TARGET_DB}"
log "已恢复数据库快照: ${TARGET_DB}"

find_compose
(
  cd "${REPO_ROOT}"
  "${COMPOSE_CMD[@]}" up -d --force-recreate "${COMPOSE_SERVICE}"
)

if wait_for_health; then
  log "健康检查通过: ${HEALTH_URL:-<disabled>}"
else
  die "服务重启后健康检查失败，请检查容器日志。"
fi

log "回滚完成。"
