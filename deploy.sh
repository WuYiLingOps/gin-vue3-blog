#!/bin/bash
#
#********************************************************************
# 项目名称：go-vue3-blog
# 文件名称：deploy.sh
# 创建时间：2026-02-15 13:12:45
#
# 系统用户：Administrator
# 作　　者：無以菱
# 联系邮箱：huangjing510@126.com
# 功能描述：博客项目部署管理脚本，支持 build（构建并部署）、start（启动服务）、
#          stop（停止服务）、status（查看状态）四种操作模式
#********************************************************************

# ==================== 颜色定义 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[0;37m'
BOLD='\033[1m'
RESET='\033[0m'

# ==================== 配置定义 ====================
PROJECT_ROOT="/web/gin-vue3-blog"
BACKEND_PORT=8080
GITEE_API_PORT=8081
BACKEND_DIR="blog-backend"
FRONTEND_DIR="blog-frontend"
BACKEND_BIN="blog-backend"
GITEE_API_BIN="gitee-calendar-api"

# ==================== 日志函数 ====================
# 输出信息日志
log_info() {
    echo -e "${GREEN}[INFO]${RESET} $1"
}

# 输出警告日志
log_warn() {
    echo -e "${YELLOW}[WARN]${RESET} $1"
}

# 输出错误日志
log_error() {
    echo -e "${RED}[ERROR]${RESET} $1"
}

# 输出步骤日志
log_step() {
    echo -e "${BLUE}[STEP]${RESET} ${BOLD}$1${RESET}"
}

# 输出成功日志
log_success() {
    echo -e "${GREEN}[SUCCESS]${RESET} ${BOLD}$1${RESET}"
}

# ==================== 工具函数 ====================
# 检查端口是否在监听
# @param port int 端口号
# @return 0 表示端口在监听，1 表示未监听
check_port() {
    local port=$1
    netstat -tunlp 2>/dev/null | grep -q ":${port} "
    return $?
}

# 获取端口对应的进程ID
# @param port int 端口号
# @return string 进程ID，未找到返回空
get_pid_by_port() {
    local port=$1
    local pid=$(netstat -tunlp 2>/dev/null | grep ":${port} " | awk '{print $7}' | cut -d'/' -f1 | head -n1)
    echo "$pid"
}

# 优雅停止进程
# @param pid int 进程ID
# @param name string 服务名称
stop_process() {
    local pid=$1
    local name=$2
    
    if [ -z "$pid" ]; then
        return 1
    fi
    
    log_info "正在停止 $name (PID: $pid)..."
    
    # 先尝试优雅停止
    kill -15 "$pid" 2>/dev/null
    
    # 等待进程退出，最多等待 5 秒
    local count=0
    while [ $count -lt 5 ]; do
        if ! kill -0 "$pid" 2>/dev/null; then
            log_info "$name 已优雅停止"
            return 0
        fi
        sleep 1
        ((count++))
    done
    
    # 强制停止
    log_warn "$name 未响应，强制停止..."
    kill -9 "$pid" 2>/dev/null
    sleep 1
    log_info "$name 已强制停止"
    return 0
}

# 显示使用帮助
show_usage() {
    echo -e "${CYAN}${BOLD}博客项目部署管理脚本${RESET}"
    echo ""
    echo -e "${BOLD}用法:${RESET}"
    echo "  $0 <command>"
    echo ""
    echo -e "${BOLD}可用命令:${RESET}"
    echo -e "  ${GREEN}build${RESET}   - 重新构建并部署项目（编译后端 + 构建前端）"
    echo -e "  ${GREEN}start${RESET}   - 启动服务（仅启动后端服务，不构建前端）"
    echo -e "  ${GREEN}stop${RESET}    - 停止所有服务"
    echo -e "  ${GREEN}status${RESET}  - 查看服务运行状态"
    echo ""
    echo -e "${BOLD}示例:${RESET}"
    echo "  $0 build    # 完整构建并部署"
    echo "  $0 start    # 启动服务"
    echo "  $0 stop     # 停止服务"
    echo "  $0 status   # 查看状态"
    echo ""
    echo -e "${YELLOW}${BOLD}注意:${RESET}"
    echo -e "  使用前请先修改脚本中的 ${CYAN}PROJECT_ROOT${RESET} 变量为实际的项目路径"
    echo -e "  当前配置: ${CYAN}$PROJECT_ROOT${RESET}"
    echo ""
}

# ==================== 功能函数 ====================
# 启动 gitee-calendar-api 服务
start_gitee_api() {
    log_step "检查 gitee-calendar-api 服务（端口 $GITEE_API_PORT）"
    
    if check_port $GITEE_API_PORT; then
        local pid=$(get_pid_by_port $GITEE_API_PORT)
        log_info "gitee-calendar-api 服务已在运行 (PID: $pid)"
        return 0
    fi
    
    log_warn "gitee-calendar-api 服务未运行，正在启动..."
    
    cd "$PROJECT_ROOT" || {
        log_error "无法进入项目目录: $PROJECT_ROOT"
        return 1
    }
    
    if [ ! -f "./$GITEE_API_BIN" ]; then
        log_error "未找到 $GITEE_API_BIN 文件"
        return 1
    fi
    
    nohup ./$GITEE_API_BIN > gitee-calendar-api.log 2>&1 &
    sleep 2
    
    if check_port $GITEE_API_PORT; then
        local pid=$(get_pid_by_port $GITEE_API_PORT)
        log_success "gitee-calendar-api 服务启动成功 (PID: $pid)"
        return 0
    else
        log_error "gitee-calendar-api 服务启动失败，请检查 gitee-calendar-api.log"
        return 1
    fi
}

# 启动后端服务
start_backend() {
    log_step "检查 Go 后端服务（端口 $BACKEND_PORT）"
    
    if check_port $BACKEND_PORT; then
        local pid=$(get_pid_by_port $BACKEND_PORT)
        log_info "Go 后端服务已在运行 (PID: $pid)"
        return 0
    fi
    
    log_warn "Go 后端服务未运行，正在启动..."
    
    cd "$PROJECT_ROOT/$BACKEND_DIR" || {
        log_error "无法进入后端目录: $PROJECT_ROOT/$BACKEND_DIR"
        return 1
    }
    
    if [ ! -f "./$BACKEND_BIN" ]; then
        log_error "未找到 $BACKEND_BIN 文件，请先执行 build 命令编译"
        return 1
    fi
    
    nohup ./$BACKEND_BIN > app.log 2>&1 &
    sleep 2
    
    if check_port $BACKEND_PORT; then
        local pid=$(get_pid_by_port $BACKEND_PORT)
        log_success "Go 后端服务启动成功 (PID: $pid)"
        return 0
    else
        log_error "Go 后端服务启动失败，请检查 app.log"
        return 1
    fi
}

# 停止 gitee-calendar-api 服务
stop_gitee_api() {
    log_step "停止 gitee-calendar-api 服务（端口 $GITEE_API_PORT）"
    
    if ! check_port $GITEE_API_PORT; then
        log_info "gitee-calendar-api 服务未运行"
        return 0
    fi
    
    local pid=$(get_pid_by_port $GITEE_API_PORT)
    stop_process "$pid" "gitee-calendar-api"
    return $?
}

# 停止后端服务
stop_backend() {
    log_step "停止 Go 后端服务（端口 $BACKEND_PORT）"
    
    if ! check_port $BACKEND_PORT; then
        log_info "Go 后端服务未运行"
        return 0
    fi
    
    local pid=$(get_pid_by_port $BACKEND_PORT)
    stop_process "$pid" "Go 后端服务"
    return $?
}

# 编译后端
compile_backend() {
    log_step "编译 Go 后端"
    
    cd "$PROJECT_ROOT/$BACKEND_DIR" || {
        log_error "无法进入后端目录: $PROJECT_ROOT/$BACKEND_DIR"
        return 1
    }
    
    # 删除旧的编译文件
    if [ -f "$BACKEND_BIN" ]; then
        log_info "删除旧的编译文件"
        rm -f "$BACKEND_BIN"
    fi
    
    log_info "开始编译 Go 后端..."
    if go build -o "$BACKEND_BIN" cmd/server/main.go; then
        log_success "Go 后端编译成功"
        return 0
    else
        log_error "Go 后端编译失败"
        return 1
    fi
}

# 构建前端
build_frontend() {
    log_step "构建前端静态资源"
    
    cd "$PROJECT_ROOT/$FRONTEND_DIR" || {
        log_error "无法进入前端目录: $PROJECT_ROOT/$FRONTEND_DIR"
        return 1
    }
    
    log_info "开始构建前端..."
    if pnpm build; then
        log_success "前端构建成功"
        return 0
    else
        log_error "前端构建失败"
        return 1
    fi
}

# ==================== 命令实现 ====================
# build 命令：重新构建并部署项目
cmd_build() {
    log_step "开始重新构建项目..."
    echo ""
    
    # 检查并启动 gitee-calendar-api
    start_gitee_api || exit 1
    echo ""
    
    # 停止后端服务
    if check_port $BACKEND_PORT; then
        local pid=$(get_pid_by_port $BACKEND_PORT)
        log_warn "发现运行中的后端服务（PID: $pid），正在停止..."
        stop_process "$pid" "Go 后端服务"
    fi
    echo ""
    
    # 编译后端
    compile_backend || exit 1
    echo ""
    
    # 启动后端
    start_backend || exit 1
    echo ""
    
    # 构建前端
    build_frontend || exit 1
    echo ""
    
    # 显示状态
    cmd_status
    echo ""
    
    log_success "项目构建并部署完成！"
}

# start 命令：启动服务
cmd_start() {
    log_step "启动服务..."
    echo ""
    
    cd "$PROJECT_ROOT" || {
        log_error "无法进入项目目录: $PROJECT_ROOT"
        exit 1
    }
    
    local has_error=0
    
    # 启动 gitee-calendar-api
    start_gitee_api || has_error=1
    echo ""
    
    # 启动后端
    start_backend || has_error=1
    echo ""
    
    if [ $has_error -eq 0 ]; then
        log_success "所有服务启动完成！"
    else
        log_warn "部分服务启动失败，请检查日志"
    fi
    
    echo ""
    cmd_status
}

# stop 命令：停止服务
cmd_stop() {
    log_step "停止所有服务..."
    echo ""
    
    # 停止后端服务
    stop_backend
    echo ""
    
    # 停止 gitee-calendar-api
    stop_gitee_api
    echo ""
    
    log_success "所有服务已停止"
}

# status 命令：查看服务状态
cmd_status() {
    log_step "服务运行状态"
    echo ""
    echo -e "${BOLD}┌─────────────────────────────────────────────────────────────┐${RESET}"
    echo -e "${BOLD}│                      服务状态概览                           │${RESET}"
    echo -e "${BOLD}├─────────────────────────────────────────────────────────────┤${RESET}"
    
    # 检查后端服务
    if check_port $BACKEND_PORT; then
        local pid=$(get_pid_by_port $BACKEND_PORT)
        printf "${BOLD}│${RESET} %-20s ${GREEN}● 运行中${RESET}    端口: %-6s  PID: %-8s ${BOLD}│${RESET}\n" "Go 后端服务" "$BACKEND_PORT" "$pid"
    else
        printf "${BOLD}│${RESET} %-20s ${RED}○ 已停止${RESET}    端口: %-6s  PID: %-8s ${BOLD}│${RESET}\n" "Go 后端服务" "$BACKEND_PORT" "-"
    fi
    
    # 检查 gitee-calendar-api
    if check_port $GITEE_API_PORT; then
        local pid=$(get_pid_by_port $GITEE_API_PORT)
        printf "${BOLD}│${RESET} %-20s ${GREEN}● 运行中${RESET}    端口: %-6s  PID: %-8s ${BOLD}│${RESET}\n" "gitee-calendar-api" "$GITEE_API_PORT" "$pid"
    else
        printf "${BOLD}│${RESET} %-20s ${RED}○ 已停止${RESET}    端口: %-6s  PID: %-8s ${BOLD}│${RESET}\n" "gitee-calendar-api" "$GITEE_API_PORT" "-"
    fi
    
    echo -e "${BOLD}└─────────────────────────────────────────────────────────────┘${RESET}"
}

# ==================== 主入口 ====================
main() {
    # 检查参数
    if [ $# -eq 0 ]; then
        log_error "缺少命令参数"
        echo ""
        show_usage
        exit 1
    fi
    
    local command=$1
    
    case "$command" in
        build)
            cmd_build
            ;;
        start)
            cmd_start
            ;;
        stop)
            cmd_stop
            ;;
        status)
            cmd_status
            ;;
        -h|--help|help)
            show_usage
            ;;
        *)
            log_error "未知命令: $command"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
