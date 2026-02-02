#!/bin/bash
#
#********************************************************************
#Author:           YiLing Wu (hj)
#email:            huangjing510@126.com
#Date:             2026-01-10
#FileName:         deploy.sh
#URL:              https://script.huangjingblog.cn
#Description:      重新构建并部署博客项目（go后端 + 前端）,使用之前先更新项目目录PROJECT_ROOT
#Copyright (C):    2026 All rights reserved
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

# ==================== 函数定义 ====================
# 带颜色的日志输出
log_info() {
    echo -e "${GREEN}[INFO]${RESET} $1"
}
log_warn() {
    echo -e "${YELLOW}[WARN]${RESET} $1"
}
log_error() {
    echo -e "${RED}[ERROR]${RESET} $1"
}
log_step() {
    echo -e "${BLUE}[STEP]${RESET} ${BOLD}$1${RESET}"
}

# 检查端口是否在监听
check_port() {
    local port=$1
    netstat -tunlp 2>/dev/null | grep -q ":${port} "
    return $?
}

# 获取端口对应的进程ID
get_pid_by_port() {
    local port=$1
    local pid=$(netstat -tunlp 2>/dev/null | grep ":${port} " | awk '{print $7}' | cut -d'/' -f1 | head -n1)
    echo "$pid"
}

# 项目根目录
PROJECT_ROOT="/web/gin-vue3-blog"
cd "$PROJECT_ROOT" || {
    log_error "无法进入项目目录: $PROJECT_ROOT"
    exit 1
}

log_step "开始重新构建项目..."

# ==================== 步骤1: 检查并启动 gitee-calendar-api 服务 ====================
log_step "检查 gitee-calendar-api 服务（端口 8081）"
if check_port 8081; then
    log_info "gitee-calendar-api 服务已在运行"
else
    log_warn "gitee-calendar-api 服务未运行，正在启动..."
    if [ -f "./gitee-calendar-api" ]; then
        nohup ./gitee-calendar-api > gitee-calendar-api.log 2>&1 &
        sleep 2
        if check_port 8081; then
            log_info "gitee-calendar-api 服务启动成功"
        else
            log_error "gitee-calendar-api 服务启动失败"
            exit 1
        fi
    else
        log_error "未找到 gitee-calendar-api 文件"
        exit 1
    fi
fi

# ==================== 步骤2: 停止并重新编译启动 go 后端服务 ====================
log_step "检查 go 后端服务（端口 8080）"
if check_port 8080; then
    PID=$(get_pid_by_port 8080)
    if [ -n "$PID" ]; then
        log_warn "发现运行中的后端服务（PID: $PID），正在停止..."
        kill -9 "$PID" 2>/dev/null
        sleep 1
        log_info "后端服务已停止"
    else
        log_warn "端口 8080 被占用但无法获取进程ID"
    fi
else
    log_info "后端服务未运行"
fi

# ==================== 步骤3: 重新编译 go 后端 ====================
log_step "重新编译 go 后端"
cd blog-backend || {
    log_error "无法进入 blog-backend 目录"
    exit 1
}

# 删除旧的编译文件
if [ -f "blog-backend" ]; then
    log_info "删除旧的编译文件"
    rm -f blog-backend
fi

# 重新编译
log_info "开始编译 go 后端..."
if go build -o blog-backend cmd/server/main.go; then
    log_info "go 后端编译成功"
else
    log_error "go 后端编译失败"
    exit 1
fi

# ==================== 步骤4: 启动 go 后端服务 ====================
log_step "启动 go 后端服务"
nohup ./blog-backend > app.log 2>&1 &
sleep 2

if check_port 8080; then
    log_info "go 后端服务启动成功"
else
    log_error "go 后端服务启动失败，请检查 app.log"
    exit 1
fi

# ==================== 步骤5: 重新构建前端 ====================
log_step "重新构建前端静态资源"
cd "$PROJECT_ROOT/blog-frontend" || {
    log_error "无法进入 blog-frontend 目录"
    exit 1
}

log_info "开始构建前端..."
if pnpm build; then
    log_info "前端构建成功"
else
    log_error "前端构建失败"
    exit 1
fi

# ==================== 步骤6: 检查服务状态 ====================
log_step "检查服务端口状态"
echo ""
log_info "检查端口 8080 (go 后端服务)..."
if check_port 8080; then
    PID=$(get_pid_by_port 8080)
    log_info "✓ 端口 8080 正常运行 (PID: $PID)"
else
    log_error "✗ 端口 8080 未运行"
fi

echo ""
log_info "检查端口 8081 (gitee-calendar-api 服务)..."
if check_port 8081; then
    PID=$(get_pid_by_port 8081)
    log_info "✓ 端口 8081 正常运行 (PID: $PID)"
else
    log_error "✗ 端口 8081 未运行"
fi

echo ""
log_step "重新构建项目完成！"
log_info "所有服务已重新部署并运行"

