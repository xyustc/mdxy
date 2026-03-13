#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
BACKEND_DIR=$(cd "${SCRIPT_DIR}/.." && pwd)
REPO_ROOT=$(cd "${BACKEND_DIR}/.." && pwd)

SOURCE_DB=""
TARGET_DB="${BACKEND_DIR}/data/mdxy.db"
BACKUP_ROOT="${BACKEND_DIR}/data/cutover-backups"
COMPOSE_SERVICE="backend"
HEALTH_URL="http://127.0.0.1:8080/health"
SKIP_BUILD=0
SKIP_RESTART=0

log() {
  printf '[cutover] %s\n' "$*"
}

die() {
  printf '[cutover] error: %s\n' "$*" >&2
  exit 1
}

usage() {
  cat <<'EOF'
Usage:
  bash backend-go/scripts/cutover-from-python.sh --source-db /path/to/analytics.db [options]

Options:
  --source-db PATH      Python backend analytics.db path (required)
  --target-db PATH      Go backend mdxy.db path
  --backup-root PATH    Directory used to store cutover artifacts
  --compose-service     Docker Compose service name for the Go backend (default: backend)
  --health-url URL      Health check URL after restart (default: http://127.0.0.1:8080/health)
  --skip-build          Skip docker compose build
  --skip-restart        Only migrate data, do not restart docker service
  -h, --help            Show this help
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

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "缺少命令: $1"
}

encode_b64() {
  printf '%s' "$1" | base64 | tr -d '\n'
}

sqlite_stat() {
  local db_path=$1
  local query=$2
  if command -v sqlite3 >/dev/null 2>&1 && [ -f "$db_path" ]; then
    sqlite3 "$db_path" "$query" 2>/dev/null || true
  fi
}

wait_for_health() {
  local attempt

  [ -n "${HEALTH_URL}" ] || return 0
  require_cmd curl

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
    --source-db)
      SOURCE_DB=${2:-}
      shift 2
      ;;
    --target-db)
      TARGET_DB=${2:-}
      shift 2
      ;;
    --backup-root)
      BACKUP_ROOT=${2:-}
      shift 2
      ;;
    --compose-service)
      COMPOSE_SERVICE=${2:-}
      shift 2
      ;;
    --health-url)
      HEALTH_URL=${2:-}
      shift 2
      ;;
    --skip-build)
      SKIP_BUILD=1
      shift
      ;;
    --skip-restart)
      SKIP_RESTART=1
      shift
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

[ -n "${SOURCE_DB}" ] || {
  usage
  exit 2
}

require_cmd go
require_cmd cp

SOURCE_DB=$(cd "$(dirname "${SOURCE_DB}")" && pwd)/$(basename "${SOURCE_DB}")
TARGET_DB=$(cd "$(dirname "${TARGET_DB}")" && pwd)/$(basename "${TARGET_DB}")
BACKUP_ROOT=$(mkdir -p "${BACKUP_ROOT}" && cd "${BACKUP_ROOT}" && pwd)

[ -f "${SOURCE_DB}" ] || die "源数据库不存在: ${SOURCE_DB}"
mkdir -p "$(dirname "${TARGET_DB}")"

TIMESTAMP=$(date '+%Y%m%d-%H%M%S')
CUTOVER_DIR="${BACKUP_ROOT}/cutover-${TIMESTAMP}"
MIGRATION_BACKUP_DIR="${CUTOVER_DIR}/migration-artifacts"
PYTHON_SOURCE_BACKUP="${CUTOVER_DIR}/python-analytics.db"
PRE_CUTOVER_GO_DB="${CUTOVER_DIR}/go-pre-cutover.db"

mkdir -p "${CUTOVER_DIR}" "${MIGRATION_BACKUP_DIR}"
cp "${SOURCE_DB}" "${PYTHON_SOURCE_BACKUP}"

if [ -f "${TARGET_DB}" ]; then
  cp "${TARGET_DB}" "${PRE_CUTOVER_GO_DB}"
fi

cat > "${CUTOVER_DIR}/rollback.env" <<EOF
CUTOVER_TIMESTAMP=${TIMESTAMP}
REPO_ROOT_B64=$(encode_b64 "${REPO_ROOT}")
BACKEND_DIR_B64=$(encode_b64 "${BACKEND_DIR}")
TARGET_DB_B64=$(encode_b64 "${TARGET_DB}")
PRE_CUTOVER_GO_DB_B64=$(encode_b64 "${PRE_CUTOVER_GO_DB}")
PYTHON_SOURCE_DB_B64=$(encode_b64 "${SOURCE_DB}")
PYTHON_SOURCE_BACKUP_B64=$(encode_b64 "${PYTHON_SOURCE_BACKUP}")
COMPOSE_SERVICE_B64=$(encode_b64 "${COMPOSE_SERVICE}")
HEALTH_URL_B64=$(encode_b64 "${HEALTH_URL}")
EOF
chmod 600 "${CUTOVER_DIR}/rollback.env"

log "Python 源库备份: ${PYTHON_SOURCE_BACKUP}"
if [ -f "${PRE_CUTOVER_GO_DB}" ]; then
  log "Go 目标库备份: ${PRE_CUTOVER_GO_DB}"
else
  log "Go 目标库当前不存在，将创建新库: ${TARGET_DB}"
fi

BEFORE_COUNT=$(sqlite_stat "${SOURCE_DB}" "select count(*) from access_logs;")
if [ -n "${BEFORE_COUNT}" ]; then
  log "Python access_logs 条数: ${BEFORE_COUNT}"
fi

log "执行访问数据迁移..."
(
  cd "${BACKEND_DIR}"
  go run ./cmd/migrate-python-db \
    -source "${SOURCE_DB}" \
    -target "${TARGET_DB}" \
    -backup-dir "${MIGRATION_BACKUP_DIR}"
)

AFTER_COUNT=$(sqlite_stat "${TARGET_DB}" "select count(*) from access_logs;")
if [ -n "${AFTER_COUNT}" ]; then
  log "Go access_logs 条数: ${AFTER_COUNT}"
fi

if [ "${SKIP_RESTART}" -eq 0 ]; then
  find_compose

  if [ "${SKIP_BUILD}" -eq 0 ]; then
    log "构建 Go 后端镜像..."
    (
      cd "${REPO_ROOT}"
      "${COMPOSE_CMD[@]}" build "${COMPOSE_SERVICE}"
    )
  fi

  log "重启 Go 后端服务..."
  (
    cd "${REPO_ROOT}"
    "${COMPOSE_CMD[@]}" up -d --force-recreate "${COMPOSE_SERVICE}"
  )

  if wait_for_health; then
    log "健康检查通过: ${HEALTH_URL}"
  else
    die "健康检查失败，请检查服务日志。可使用 rollback 脚本恢复旧 Go 数据库快照。"
  fi
fi

log "切换完成。"
log "切换产物目录: ${CUTOVER_DIR}"
log "如需恢复切换前的 Go 数据库快照，可执行:"
log "  bash backend-go/scripts/rollback-go-cutover.sh --cutover-dir ${CUTOVER_DIR}"
