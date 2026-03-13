#!/bin/bash
# 开发环境服务管理脚本
# 用法: ./dev.sh {start|stop|restart|status|logs}
# 可选环境变量:
#   FRONTEND_HOST=0.0.0.0   # 默认监听所有网卡，便于局域网访问
#   FRONTEND_PORT=5173

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$PROJECT_DIR/backend-go"
FRONTEND_DIR="$PROJECT_DIR/frontend"
BACKEND_PID_FILE="$PROJECT_DIR/.backend.pid"
FRONTEND_PID_FILE="$PROJECT_DIR/.frontend.pid"
BACKEND_LOG="$PROJECT_DIR/.backend.log"
FRONTEND_LOG="$PROJECT_DIR/.frontend.log"

BACKEND_PORT="${BACKEND_PORT:-8080}"
FRONTEND_PORT="${FRONTEND_PORT:-5173}"
FRONTEND_HOST="${FRONTEND_HOST:-0.0.0.0}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# ── 工具函数 ──

check_deps() {
    local missing=()
    command -v go  >/dev/null 2>&1 || missing+=("go")
    command -v node >/dev/null 2>&1 || missing+=("node")
    command -v npx  >/dev/null 2>&1 || missing+=("npx")
    if [ ${#missing[@]} -gt 0 ]; then
        echo -e "${RED}缺少依赖: ${missing[*]}${NC}"
        exit 1
    fi
}

get_lan_ip() {
    local ip_addr=""

    if [[ "$(uname -s)" == "Darwin" ]]; then
        ip_addr=$(ipconfig getifaddr en0 2>/dev/null || true)
        [ -z "$ip_addr" ] && ip_addr=$(ipconfig getifaddr en1 2>/dev/null || true)
    else
        if command -v ip >/dev/null 2>&1; then
            ip_addr=$(ip -4 route get 1.1.1.1 2>/dev/null | awk '{for(i=1;i<=NF;i++) if($i=="src"){print $(i+1); exit}}')
        fi
        [ -z "$ip_addr" ] && ip_addr=$(hostname -I 2>/dev/null | awk '{print $1}')
    fi

    if [ -z "$ip_addr" ] && command -v ifconfig >/dev/null 2>&1; then
        ip_addr=$(ifconfig 2>/dev/null | awk '/inet / && $2 != "127.0.0.1" {print $2; exit}')
    fi

    echo "$ip_addr"
}

print_frontend_access_urls() {
    echo -e "前端本机地址: ${CYAN}http://localhost:${FRONTEND_PORT}${NC}"
    if [[ "$FRONTEND_HOST" == "0.0.0.0" || "$FRONTEND_HOST" == "::" ]]; then
        local lan_ip
        lan_ip=$(get_lan_ip)
        if [ -n "$lan_ip" ]; then
            echo -e "前端局域网地址: ${CYAN}http://${lan_ip}:${FRONTEND_PORT}${NC}"
            echo -e "${YELLOW}提示: 局域网 HTTP 下，浏览器可能禁用摄像头 API（如 Stark Shapes 手势控制）。${NC}"
            echo -e "${YELLOW}如需手势控制，请使用 HTTPS 或本机 localhost 访问。${NC}"
        else
            echo -e "${YELLOW}未自动识别局域网 IP，请使用本机 IP + 端口 ${FRONTEND_PORT}${NC}"
        fi
    else
        echo -e "前端绑定地址: ${CYAN}http://${FRONTEND_HOST}:${FRONTEND_PORT}${NC}"
    fi
}

check_port() {
    local port="$1"
    local name="$2"
    local pid
    pid=$(lsof -ti :"$port" 2>/dev/null || true)
    if [ -n "$pid" ]; then
        local cmd
        cmd=$(ps -p "$pid" -o comm= 2>/dev/null || echo "unknown")
        echo -e "${YELLOW}端口 $port 已被占用 (PID: $pid, 进程: $cmd)${NC}"
        echo -n "是否强制终止该进程以启动${name}? [y/N] "
        read -r answer
        if [[ "$answer" =~ ^[Yy]$ ]]; then
            kill "$pid" 2>/dev/null
            sleep 1
            if lsof -ti :"$port" >/dev/null 2>&1; then
                kill -9 "$pid" 2>/dev/null
                sleep 0.5
            fi
            echo -e "端口 $port 已释放"
        else
            echo -e "${RED}跳过启动${name}${NC}"
            return 1
        fi
    fi
    return 0
}

is_running() {
    local pid_file="$1"
    if [ -f "$pid_file" ]; then
        local pid
        pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            return 0
        fi
        rm -f "$pid_file"
    fi
    return 1
}

cleanup_stale_pid() {
    local pid_file="$1"
    local name="$2"
    if [ -f "$pid_file" ]; then
        local pid
        pid=$(cat "$pid_file")
        if ! kill -0 "$pid" 2>/dev/null; then
            echo -e "${YELLOW}清理${name}残留 PID 文件${NC}"
            rm -f "$pid_file"
        fi
    fi
}

# ── 启动 ──

start_backend() {
    cleanup_stale_pid "$BACKEND_PID_FILE" "后端"
    if is_running "$BACKEND_PID_FILE"; then
        echo -e "${YELLOW}后端已在运行 (PID: $(cat "$BACKEND_PID_FILE"))${NC}"
        return 0
    fi
    check_port "$BACKEND_PORT" "后端" || return 1

    echo -n "编译后端..."
    cd "$BACKEND_DIR" || exit 1
    if ! go build -o server ./cmd/server/ 2>"$BACKEND_LOG"; then
        echo -e " ${RED}编译失败${NC}"
        echo -e "错误日志: $BACKEND_LOG"
        tail -5 "$BACKEND_LOG" 2>/dev/null
        return 1
    fi
    echo -e " ${GREEN}OK${NC}"

    echo -n "启动后端..."
    nohup ./server >> "$BACKEND_LOG" 2>&1 &
    echo $! > "$BACKEND_PID_FILE"

    # 等待健康检查通过
    local retries=0
    while [ $retries -lt 15 ]; do
        if curl -sf "http://localhost:${BACKEND_PORT}/health" >/dev/null 2>&1; then
            echo -e " ${GREEN}OK${NC} (PID: $(cat "$BACKEND_PID_FILE"), 端口: $BACKEND_PORT)"
            return 0
        fi
        if ! is_running "$BACKEND_PID_FILE"; then
            echo -e " ${RED}进程异常退出${NC}"
            echo "最近日志:"
            tail -10 "$BACKEND_LOG" 2>/dev/null
            return 1
        fi
        sleep 0.5
        retries=$((retries + 1))
    done
    echo -e " ${RED}启动超时${NC} (健康检查未通过)"
    tail -10 "$BACKEND_LOG" 2>/dev/null
    return 1
}

start_frontend() {
    cleanup_stale_pid "$FRONTEND_PID_FILE" "前端"
    if is_running "$FRONTEND_PID_FILE"; then
        echo -e "${YELLOW}前端已在运行 (PID: $(cat "$FRONTEND_PID_FILE"))${NC}"
        return 0
    fi
    check_port "$FRONTEND_PORT" "前端" || return 1

    # 检查 node_modules
    if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
        echo -n "安装前端依赖..."
        cd "$FRONTEND_DIR" || exit 1
        if ! npm install >> "$FRONTEND_LOG" 2>&1; then
            echo -e " ${RED}npm install 失败${NC}"
            tail -10 "$FRONTEND_LOG" 2>/dev/null
            return 1
        fi
        echo -e " ${GREEN}OK${NC}"
    fi

    echo -n "启动前端..."
    cd "$FRONTEND_DIR" || exit 1
    nohup npx vite --host "$FRONTEND_HOST" --port "$FRONTEND_PORT" --strictPort >> "$FRONTEND_LOG" 2>&1 &
    echo $! > "$FRONTEND_PID_FILE"

    # 等待 Vite 就绪
    local retries=0
    while [ $retries -lt 20 ]; do
        if curl -sf "http://localhost:${FRONTEND_PORT}" >/dev/null 2>&1; then
            echo -e " ${GREEN}OK${NC} (PID: $(cat "$FRONTEND_PID_FILE"), 端口: $FRONTEND_PORT)"
            print_frontend_access_urls
            return 0
        fi
        if ! is_running "$FRONTEND_PID_FILE"; then
            echo -e " ${RED}进程异常退出${NC}"
            tail -10 "$FRONTEND_LOG" 2>/dev/null
            return 1
        fi
        sleep 0.5
        retries=$((retries + 1))
    done
    echo -e " ${RED}启动超时${NC}"
    tail -10 "$FRONTEND_LOG" 2>/dev/null
    return 1
}

# ── 停止 ──

stop_service() {
    local name="$1"
    local pid_file="$2"
    local port="$3"
    if is_running "$pid_file"; then
        local pid
        pid=$(cat "$pid_file")
        echo -n "停止${name} (PID: $pid)..."
        kill "$pid" 2>/dev/null
        local count=0
        while kill -0 "$pid" 2>/dev/null && [ $count -lt 10 ]; do
            sleep 0.5
            count=$((count + 1))
        done
        if kill -0 "$pid" 2>/dev/null; then
            kill -9 "$pid" 2>/dev/null
            sleep 0.5
        fi
        rm -f "$pid_file"
        echo -e " ${GREEN}OK${NC}"
    else
        echo -e "${name}未在运行"
    fi
    # 清理可能残留的端口占用（子进程）
    local remaining
    remaining=$(lsof -ti :"$port" 2>/dev/null || true)
    if [ -n "$remaining" ]; then
        echo -e "${YELLOW}清理端口 $port 残留进程: $remaining${NC}"
        echo "$remaining" | xargs kill 2>/dev/null || true
        sleep 0.5
    fi
}

# ── 状态 ──

show_status() {
    echo -e "${CYAN}── 服务状态 ──${NC}"
    if is_running "$BACKEND_PID_FILE"; then
        echo -e "后端:  ${GREEN}运行中${NC} (PID: $(cat "$BACKEND_PID_FILE"), 端口: $BACKEND_PORT)"
    else
        echo -e "后端:  ${RED}未运行${NC}"
    fi
    if is_running "$FRONTEND_PID_FILE"; then
        echo -e "前端:  ${GREEN}运行中${NC} (PID: $(cat "$FRONTEND_PID_FILE"), 端口: $FRONTEND_PORT)"
        print_frontend_access_urls
    else
        echo -e "前端:  ${RED}未运行${NC}"
    fi
}

# ── 日志 ──

show_logs() {
    local target="${2:-all}"
    case "$target" in
        backend|b)
            echo -e "${CYAN}── 后端日志 ──${NC}"
            tail -50 "$BACKEND_LOG" 2>/dev/null || echo "无日志"
            ;;
        frontend|f)
            echo -e "${CYAN}── 前端日志 ──${NC}"
            tail -50 "$FRONTEND_LOG" 2>/dev/null || echo "无日志"
            ;;
        *)
            echo -e "${CYAN}── 后端日志 (最近 20 行) ──${NC}"
            tail -20 "$BACKEND_LOG" 2>/dev/null || echo "无日志"
            echo ""
            echo -e "${CYAN}── 前端日志 (最近 20 行) ──${NC}"
            tail -20 "$FRONTEND_LOG" 2>/dev/null || echo "无日志"
            ;;
    esac
}

# ── 主入口 ──

check_deps

case "${1:-}" in
    start)
        start_backend
        start_frontend
        echo ""
        show_status
        ;;
    stop)
        stop_service "前端" "$FRONTEND_PID_FILE" "$FRONTEND_PORT"
        stop_service "后端" "$BACKEND_PID_FILE" "$BACKEND_PORT"
        ;;
    restart)
        stop_service "前端" "$FRONTEND_PID_FILE" "$FRONTEND_PORT"
        stop_service "后端" "$BACKEND_PID_FILE" "$BACKEND_PORT"
        sleep 1
        start_backend
        start_frontend
        echo ""
        show_status
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs "$@"
        ;;
    *)
        echo "用法: $0 {start|stop|restart|status|logs [backend|frontend]}"
        exit 1
        ;;
esac
