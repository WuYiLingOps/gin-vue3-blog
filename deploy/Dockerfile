#********************************************************************
# 项目名称：gin-vue3-blog
# 文件名称：Dockerfile
# 创建时间：2026-03-17 21:35:42
#
# 系统用户：root
# 作　　者：無以菱
# 联系邮箱：huangjing510@126.com
# 功能描述：单镜像生产构建，多阶段编译前端与后端 Go 服务，运行时使用 Nginx 统一入口并通过 supervisord 管理多进程
#********************************************************************

############################
# 1) frontend builder
############################
#FROM node:20-bookworm-slim AS frontend-builder
FROM registry.cn-hangzhou.aliyuncs.com/wylhub/node:20-bookworm-slim AS frontend-builder

WORKDIR /src/blog-frontend

# 仅复制依赖描述文件以利用缓存
COPY blog-frontend/package.json ./

# 安装 pnpm 并安装依赖
RUN corepack enable && corepack prepare pnpm@10.26.1 --activate && pnpm install

# 复制源码并构建
COPY blog-frontend/ ./
RUN pnpm build

############################
# 2) blog-backend builder
############################
#FROM golang:1.25-bookworm AS blog-backend-builder
FROM registry.cn-hangzhou.aliyuncs.com/wylhub/golang:1.25-bookworm AS blog-backend-builder

WORKDIR /src/blog-backend

COPY blog-backend/go.mod blog-backend/go.sum ./

# Go 依赖下载代理（构建时可用 --build-arg 覆盖）
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

RUN go mod download

COPY blog-backend/ ./

# 默认尝试构建静态二进制；如项目依赖 CGO，可在构建时传入 CGO_ENABLED=1
ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED}
RUN go build -o /out/blog-backend ./cmd/server

############################
# 3) runtime
############################
#FROM debian:bookworm-slim AS runtime
FROM registry.cn-hangzhou.aliyuncs.com/wylhub/debian:bookworm-slim AS runtime

ENV TZ=Asia/Shanghai

RUN apt-get update \
  && apt-get install -y --no-install-recommends nginx supervisor ca-certificates tzdata curl \
  && rm -f /etc/nginx/sites-enabled/default \
  && rm -rf /var/lib/apt/lists/*

# 路径约定：后端相对路径读取 ./config/config.yml 与 ./uploads
WORKDIR /app

# Go 二进制
COPY --from=blog-backend-builder /out/blog-backend /app/blog-backend

# 后端默认配置（生产请通过挂载 /app/config 覆盖）
COPY blog-backend/config /app/config

# 上传目录（生产请通过挂载 /app/uploads 持久化）
RUN mkdir -p /app/uploads/avatars

# 日志目录（生产请通过挂载 /var/log/gin-vue3-blog 持久化）
RUN mkdir -p /var/log/gin-vue3-blog

# 前端静态
COPY --from=frontend-builder /src/blog-frontend/dist /usr/share/nginx/html

# Nginx 配置与 supervisor 配置
COPY deploy/docker/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY deploy/docker/supervisord.conf /etc/supervisor/supervisord.conf

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD curl -fs http://127.0.0.1:8080/health || exit 1

CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
