#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="${ROOT_DIR}/backend-go"
ENV_FILE="${ROOT_DIR}/.env.prod"
DATA_DIR="${DATA_DIR:-/var/lib/mdxy}"

USERNAME="admin"
DB_PATH=""
NEW_PASSWORD=""

log() {
  printf '[admin-passwd] %s\n' "$*"
}

die() {
  printf '[admin-passwd] error: %s\n' "$*" >&2
  exit 1
}

usage() {
  cat <<'EOF'
Usage:
  bash deploy/reset-admin-password.sh [options]

Options:
  -u, --username <name>   管理员用户名（默认：admin）
  -d, --db-path <path>    SQLite 数据库路径（默认从 .env.prod / DATABASE_PATH 推断）
  -p, --password <pass>   新密码（不推荐明文传参，建议交互输入）
      --env-file <path>   指定环境变量文件（默认：.env.prod）
  -h, --help              查看帮助

Examples:
  # 交互式输入新密码（推荐）
  bash deploy/reset-admin-password.sh

  # 指定用户和数据库
  bash deploy/reset-admin-password.sh --username admin --db-path /var/lib/mdxy/mdxy.db
EOF
}

parse_args() {
  local arg
  while [[ $# -gt 0 ]]; do
    arg="$1"
    shift
    case "${arg}" in
      -u|--username)
        [[ $# -gt 0 ]] || die "--username 需要参数"
        USERNAME="$1"
        shift
        ;;
      -d|--db-path)
        [[ $# -gt 0 ]] || die "--db-path 需要参数"
        DB_PATH="$1"
        shift
        ;;
      -p|--password)
        [[ $# -gt 0 ]] || die "--password 需要参数"
        NEW_PASSWORD="$1"
        shift
        ;;
      --env-file)
        [[ $# -gt 0 ]] || die "--env-file 需要参数"
        ENV_FILE="$1"
        shift
        ;;
      -h|--help|help)
        usage
        exit 0
        ;;
      *)
        die "未知参数: ${arg}"
        ;;
    esac
  done
}

load_env_if_exists() {
  if [[ -f "${ENV_FILE}" ]]; then
    # shellcheck disable=SC1090
    source "${ENV_FILE}"
  fi
}

resolve_db_path() {
  if [[ -z "${DB_PATH}" ]]; then
    DB_PATH="${DATABASE_PATH:-${DATA_DIR}/mdxy.db}"
  fi

  if [[ "${DB_PATH}" == /app/* ]]; then
    log "检测到容器路径 ${DB_PATH}，自动映射为 ${DATA_DIR}/mdxy.db"
    DB_PATH="${DATA_DIR}/mdxy.db"
  fi

  if [[ "${DB_PATH}" != /* ]]; then
    DB_PATH="${BACKEND_DIR}/${DB_PATH}"
  fi
}

read_password_if_needed() {
  if [[ -n "${NEW_PASSWORD}" ]]; then
    return
  fi

  local p1 p2
  read -r -s -p "请输入新的管理员密码: " p1
  printf '\n'
  read -r -s -p "请再次输入新的管理员密码: " p2
  printf '\n'

  [[ -n "${p1}" ]] || die "密码不能为空"
  [[ "${p1}" == "${p2}" ]] || die "两次输入的密码不一致"
  NEW_PASSWORD="${p1}"
}

require_commands() {
  command -v sqlite3 >/dev/null 2>&1 || die "未找到 sqlite3，请先安装"
  command -v go >/dev/null 2>&1 || die "未找到 go，请先安装 Go 1.21+"
}

validate_inputs() {
  [[ -d "${BACKEND_DIR}" ]] || die "后端目录不存在: ${BACKEND_DIR}"
  [[ -f "${DB_PATH}" ]] || die "数据库文件不存在: ${DB_PATH}"
  [[ "${USERNAME}" =~ ^[A-Za-z0-9._-]+$ ]] || die "用户名包含非法字符"
  [[ -n "${NEW_PASSWORD}" ]] || die "密码不能为空"
}

admin_exists() {
  local count
  count="$(sqlite3 "${DB_PATH}" "SELECT COUNT(1) FROM admins WHERE username='${USERNAME}';")"
  [[ "${count}" != "0" ]] || die "未找到管理员用户: ${USERNAME}"
}

generate_bcrypt_hash() {
  local tmp_base tmp_go
  tmp_base="$(mktemp "${BACKEND_DIR}/tmp-bcrypt-XXXXXX")" || die "创建临时文件失败"
  tmp_go="${tmp_base}.go"
  mv "${tmp_base}" "${tmp_go}" || die "临时文件重命名失败"
  trap "rm -f '${tmp_go}'" EXIT

  cat > "${tmp_go}" <<'EOF'
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password, err := os.ReadFile("/dev/stdin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(password) == 0 {
		fmt.Fprintln(os.Stderr, "empty password")
		os.Exit(2)
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(string(hash))
}
EOF

  (cd "${BACKEND_DIR}" && printf '%s' "${NEW_PASSWORD}" | go run "${tmp_go}")
}

update_password() {
  local password_hash
  password_hash="$(generate_bcrypt_hash)"
  [[ -n "${password_hash}" ]] || die "生成密码哈希失败"

  sqlite3 "${DB_PATH}" "UPDATE admins SET password_hash='${password_hash}' WHERE username='${USERNAME}';"

  local updated
  updated="$(sqlite3 "${DB_PATH}" "SELECT COUNT(1) FROM admins WHERE username='${USERNAME}' AND password_hash='${password_hash}';")"
  [[ "${updated}" != "0" ]] || die "密码更新失败，请检查数据库权限"
}

main() {
  parse_args "$@"
  load_env_if_exists
  resolve_db_path
  read_password_if_needed
  require_commands
  validate_inputs
  admin_exists
  update_password

  log "管理员密码已更新（username=${USERNAME}）"
  log "数据库: ${DB_PATH}"
}

main "$@"
