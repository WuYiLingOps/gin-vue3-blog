.PHONY: help run build test clean docker-build docker-up docker-down migrate

# 默认目标
help:
	@echo "可用命令："
	@echo "  make run          - 运行应用"
	@echo "  make build        - 编译应用"
	@echo "  make test         - 运行测试"
	@echo "  make clean        - 清理构建文件"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make docker-up    - 启动 Docker 容器"
	@echo "  make docker-down  - 停止 Docker 容器"
	@echo "  make migrate      - 运行数据库迁移"

# 运行应用
run:
	go run cmd/server/main.go

# 编译应用
build:
	go build -o bin/blog-backend cmd/server/main.go

# 运行测试
test:
	go test -v ./...

# 清理构建文件
clean:
	rm -rf bin/
	go clean

# 构建 Docker 镜像
docker-build:
	docker build -t blog-backend:latest .

# 启动 Docker 容器
docker-up:
	docker-compose up -d

# 停止 Docker 容器
docker-down:
	docker-compose down

# 查看日志
docker-logs:
	docker-compose logs -f backend

# 数据库迁移
migrate:
	@echo "执行数据库迁移..."
	psql -U postgres -d blogdb -f sql/init.sql
	@echo "数据库迁移完成"

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 安装依赖
deps:
	go mod download
	go mod tidy

# 生成默认管理员用户
init-admin:
	@echo "生成默认管理员账户..."
	@echo "用户名: admin"
	@echo "密码: admin123"
	@echo "请在生产环境中修改密码"

