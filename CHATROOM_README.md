# 聊天室功能说明

## 功能概述

为博客系统添加了一个实时聊天室功能，支持：

- ✅ 实时消息发送和接收（基于WebSocket）
- ✅ 消息历史记录存储
- ✅ 在线用户列表显示
- ✅ 支持匿名用户和登录用户
- ✅ 管理员可以管理聊天消息
- ✅ 管理员可以发送系统广播消息

## 技术实现

### 后端（Go + Gin）

1. **数据模型** (`blog-backend/model/models.go`)
   - 添加了 `ChatMessage` 模型
   - 支持用户ID（可选）、用户名、头像、IP、消息内容等字段

2. **WebSocket Hub** (`blog-backend/service/chat.go`)
   - 实现了 `Hub` 结构体管理所有WebSocket连接
   - 支持消息广播、用户加入/离开通知
   - 自动发送历史消息和在线用户列表
   - 心跳检测和自动重连机制

3. **Repository层** (`blog-backend/repository/chat.go`)
   - 消息的增删查改
   - 分页查询历史消息
   - 获取最近消息列表

4. **Handler层** (`blog-backend/handler/chat.go`)
   - WebSocket连接处理
   - HTTP API接口（获取消息列表、在线信息等）
   - 管理员接口（删除消息、发送系统广播）

5. **路由配置** (`blog-backend/router/router.go`)
   - `/api/chat/ws` - WebSocket连接端点
   - `/api/chat/messages` - 获取消息列表
   - `/api/chat/online` - 获取在线信息
   - `/api/admin/chat/*` - 管理员接口

### 前端（Vue 3 + TypeScript）

1. **API接口** (`blog-frontend/src/api/chat.ts`)
   - 封装了所有聊天室相关的HTTP API
   - 定义了消息、在线用户等TypeScript类型

2. **WebSocket管理** (`blog-frontend/src/utils/websocket.ts`)
   - `ChatWebSocket` 类封装了WebSocket连接管理
   - 支持自动重连、心跳检测
   - 事件驱动的消息处理机制

3. **聊天室页面** (`blog-frontend/src/pages/blog/Chat.vue`)
   - 实时消息显示（支持滚动到最新消息）
   - 消息输入和发送
   - 在线用户数量显示
   - 支持匿名用户设置昵称
   - 登录用户自动使用账户信息

4. **管理后台** (`blog-frontend/src/pages/admin/ChatManage.vue`)
   - 查看所有聊天消息
   - 删除不当消息
   - 发送系统广播通知
   - 查看在线用户统计

## 数据库迁移

运行以下SQL脚本创建聊天消息表：

```bash
psql -U your_username -d your_database -f blog-backend/sql/add_chat_messages.sql
```

或者直接运行应用，GORM会自动创建表（如果配置了AutoMigrate）。

## 使用方法

### 启动后端

```bash
cd blog-backend
go run cmd/server/main.go
```

### 启动前端

```bash
cd blog-frontend
npm install
npm run dev
```

### 访问聊天室

- **前台用户**: 访问 `http://localhost:5173/chat`
  - 未登录用户：需要输入昵称才能进入
  - 已登录用户：自动使用账户信息

- **管理后台**: 访问 `http://localhost:5173/admin/chat`
  - 需要管理员权限
  - 可以查看、删除消息
  - 可以发送系统广播

## API接口说明

### 公开接口

- `GET /api/chat/ws?username=xxx&avatar=xxx` - WebSocket连接
- `GET /api/chat/messages?page=1&page_size=50` - 获取消息列表
- `GET /api/chat/online` - 获取在线信息

### 管理员接口

- `GET /api/admin/chat/messages?page=1&page_size=20` - 管理员获取消息列表
- `DELETE /api/admin/chat/messages/:id` - 删除消息
- `POST /api/admin/chat/broadcast` - 发送系统广播
  ```json
  {
    "content": "系统广播内容"
  }
  ```

## WebSocket消息格式

### 客户端发送

```json
{
  "type": "message",
  "content": "消息内容"
}
```

### 服务端推送

```json
{
  "type": "message|history|user_join|user_leave|user_list|system",
  "data": { ... },
  "timestamp": 1234567890
}
```

消息类型说明：
- `message` - 新消息
- `history` - 历史消息列表
- `user_join` - 用户加入
- `user_leave` - 用户离开
- `user_list` - 在线用户列表
- `system` - 系统消息

## 配置说明

### 后端配置

WebSocket相关配置在 `blog-backend/handler/chat.go` 中：
- `ReadBufferSize` - 读缓冲区大小（默认1024）
- `WriteBufferSize` - 写缓冲区大小（默认1024）
- `CheckOrigin` - 跨域检查（生产环境建议配置具体域名）

### 前端配置

WebSocket URL配置在 `blog-frontend/src/utils/websocket.ts` 中：
- 自动根据当前协议选择 `ws://` 或 `wss://`
- 自动使用环境变量 `VITE_API_BASE_URL` 或当前host

## 注意事项

1. **生产环境**：
   - 需要配置HTTPS和WSS（WebSocket Secure）
   - 修改CORS配置，限制允许的域名
   - 考虑添加消息内容过滤和敏感词检测

2. **性能优化**：
   - 当前在线用户数较多时，考虑使用Redis存储在线状态
   - 历史消息可以考虑添加缓存
   - 可以限制单个用户的消息发送频率

3. **安全性**：
   - 已实现IP记录，可配合IP黑名单功能
   - 建议添加消息长度限制
   - 可以添加用户发言频率限制

4. **扩展功能**（可选）：
   - 私聊功能
   - 消息撤回
   - 表情包支持
   - 图片/文件发送
   - @提及功能
   - 聊天室分组/频道

## 故障排查

1. **连接失败**：
   - 检查后端服务是否正常运行
   - 检查WebSocket URL是否正确
   - 检查防火墙设置

2. **消息发送失败**：
   - 检查WebSocket连接状态
   - 查看浏览器控制台错误信息
   - 检查后端日志

3. **消息不实时**：
   - 检查心跳机制是否正常
   - 检查网络连接稳定性
   - 查看WebSocket连接是否频繁断开重连

## 开发者信息

- 开发时间：2025-11
- 框架版本：Go 1.25, Vue 3, TypeScript
- 依赖包：gorilla/websocket, naive-ui

