# 配置文件说明

## 配置文件结构

项目使用基于环境的配置文件管理：

```
blog-backend/config/
├── config.yml       # 环境选择配置（提交到 Git）
├── config-dev.yml   # 开发环境配置（提交到 Git）
└── config-prod.yml  # 生产环境配置（不提交到 Git，需手动创建）
```

## 环境切换

通过修改 `config/config.yml` 文件中的 `env` 字段来切换环境：

### 开发环境（默认）
编辑 `config/config.yml`：
```yaml
env: dev
```

### 生产环境
编辑 `config/config.yml`：
```yaml
env: prod
```

然后直接运行：
```bash
go run ./cmd/server/main.go
# 或
./blog-backend
```

## 配置文件说明

### config-dev.yml（开发环境）
- 已提交到 Git
- 包含开发环境的默认配置
- 日志级别：`debug`
- 服务器模式：`debug`

### config-prod.yml（生产环境）
- **不提交到 Git**（已在 .gitignore 中忽略）
- 需要在生产服务器上手动创建
- 包含生产环境的敏感信息
- 日志级别：`info`
- 服务器模式：`release`

## 生产环境配置模板

创建 `config/config-prod.yml` 文件，内容如下：

```yaml
app:
  port: 8080
server:
  mode: release

db:
  host: localhost
  port: 5432
  user: postgres
  password: your_production_password_here  # 修改为生产环境密码
  dbname: blogdb
  sslmode: disable
jwt:
  secret: "your_production_secret_key_here_please_change_this"  # 修改为强密钥
  expire_hours: 72

log:
  level: info
```

## 重要提示

⚠️ **安全注意事项**：

1. **永远不要**将 `config-prod.yml` 提交到 Git
2. **必须修改**生产环境的数据库密码
3. **必须修改**生产环境的 JWT 密钥（建议使用 32 位以上的随机字符串）
4. 生产环境配置文件应该只存在于生产服务器上
5. 定期更换密钥和密码

## JWT 密钥生成

可以使用以下命令生成安全的随机密钥：

```bash
# Linux/Mac
openssl rand -base64 32

# 或使用 Go
go run -c 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

## 配置项说明

| 配置项 | 说明 | 示例值 |
|--------|------|--------|
| `app.port` | 服务器端口 | `8080` |
| `server.mode` | Gin 运行模式 | `debug` / `release` |
| `db.host` | 数据库主机 | `localhost` |
| `db.port` | 数据库端口 | `5432` |
| `db.user` | 数据库用户 | `postgres` |
| `db.password` | 数据库密码 | `your_password` |
| `db.dbname` | 数据库名称 | `blogdb` |
| `db.sslmode` | SSL 模式 | `disable` / `require` |
| `jwt.secret` | JWT 签名密钥 | 随机字符串 |
| `jwt.expire_hours` | Token 过期时间（小时） | `72` |
| `log.level` | 日志级别 | `debug` / `info` / `warn` / `error` |

## 快速切换环境

只需要修改一行配置即可切换环境：

```bash
# 切换到开发环境
echo "env: dev" > config/config.yml

# 切换到生产环境
echo "env: prod" > config/config.yml
```

## Docker 部署

使用 Docker 部署时，确保 `config.yml` 设置为 `prod`：

```yaml
# config/config.yml
env: prod
```

然后挂载配置文件：

```yaml
# docker-compose.yml
services:
  backend:
    volumes:
      - ./config:/app/config:ro
```

