# 直接使用已编译的二进制文件
# 使用前请先在本地编译 Linux 二进制文件：
# Windows PowerShell: 
#   cd blog-backend
#   $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o blog-backend ./cmd/server
# Windows CMD:
#   set GOOS=linux && set GOARCH=amd64 && go build -o blog-backend ./cmd/server
# Linux/Mac: 
#   GOOS=linux GOARCH=amd64 go build -o blog-backend ./cmd/server

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata

# 设置时区为上海
ENV TZ=Asia/Shanghai

WORKDIR /app

# 复制已编译的二进制文件（从 blog-backend 目录的根目录）
COPY blog-backend .

# 复制配置文件
COPY config ./config

# 创建上传目录
RUN mkdir -p ./uploads/avatars

# 给二进制文件执行权限
RUN chmod +x ./blog-backend

# 暴露端口（根据配置文件修改）
EXPOSE 8080

# 运行应用
CMD ["./blog-backend"]



