# 访问量统计从 PV 改为 UV 的迁移说明

## 改动概述

将博客系统的访问量统计从 **PV（页面访问量）** 改为 **UV（独立访客量）**。

- **PV (Page View)**: 统计所有页面访问次数，同一用户多次访问会重复计数
- **UV (Unique Visitor)**: 统计独立访客数，基于 IP 地址去重，同一 IP 在同一天只计数一次

## 主要修改内容

### 1. 数据库层面

#### 新增表：`visit_records`
用于记录每个 IP 的访问记录，实现去重统计。

```sql
CREATE TABLE visit_records (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    ip VARCHAR(45) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(date, ip)  -- 唯一索引确保同一 IP 每天只有一条记录
);
```

#### 更新表：`visit_stats`
- 字段 `view_count` 的含义从"当日访问量(PV)"改为"当日独立访客数(UV)"

### 2. 代码层面

#### 修改文件清单：
1. **`model/models.go`**
   - 新增 `VisitRecord` 模型
   - 更新 `VisitStat` 字段注释

2. **`repository/visit_stat.go`**
   - 将 `IncrementTodayViewCount()` 改为 `RecordVisit(ip string)`
   - 实现基于 IP 的去重统计逻辑
   - 使用事务确保数据一致性

3. **`middleware/visit_stat.go`**
   - 调用 `util.GetClientIP()` 获取访客真实 IP
   - 将 IP 传递给 `RecordVisit()` 方法

4. **`sql/init.sql`**
   - 添加 `visit_records` 表的创建语句
   - 更新相关注释

5. **`sql/migrations/add_visit_records.sql`**
   - 新增数据库迁移脚本

## 迁移步骤

### 对于新部署

直接使用更新后的 `sql/init.sql` 初始化数据库即可。

### 对于已有数据库

执行迁移脚本：

```bash
psql -U your_username -d your_database -f sql/migrations/add_visit_records.sql
```

或者手动执行以下 SQL：

```sql
-- 创建访问记录表
CREATE TABLE IF NOT EXISTS visit_records (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    ip VARCHAR(45) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_visit_records_date ON visit_records(date);
CREATE INDEX IF NOT EXISTS idx_visit_records_ip ON visit_records(ip);
CREATE UNIQUE INDEX IF NOT EXISTS idx_visit_records_date_ip ON visit_records(date, ip);

-- 更新注释
COMMENT ON TABLE visit_records IS '访问记录表（用于 UV 统计去重）';
COMMENT ON COLUMN visit_records.date IS '访问日期';
COMMENT ON COLUMN visit_records.ip IS '访客IP地址';
COMMENT ON COLUMN visit_records.created_at IS '首次访问时间';
COMMENT ON COLUMN visit_stats.view_count IS '当日独立访客数（UV）';
```

## 注意事项

1. **历史数据**: 现有的 `visit_stats` 表中的数据是 PV 统计，改动后新数据将是 UV 统计。建议：
   - 清空 `visit_stats` 表重新统计，或
   - 保留历史数据并记录切换时间点

2. **IP 获取**: 系统使用 `util.GetClientIP()` 获取访客 IP，支持以下场景：
   - 直接访问：使用 `RemoteAddr`
   - 反向代理：从 `X-Forwarded-For` 或 `X-Real-IP` 获取

3. **性能考虑**: 
   - 访问记录异步处理，不影响响应速度
   - 使用唯一索引防止重复插入
   - 使用事务确保数据一致性

4. **统计范围**: 继续只统计：
   - GET 请求
   - 200-299 状态码
   - 非静态资源（图片、CSS、JS 等）
   - 非 API 请求

## 验证方法

1. 启动后端服务
2. 访问博客首页
3. 检查数据库：
```sql
-- 查看今天的 UV 统计
SELECT * FROM visit_stats WHERE date = CURRENT_DATE;

-- 查看今天的访问记录
SELECT * FROM visit_records WHERE date = CURRENT_DATE;
```
4. 同一 IP 多次访问，UV 应该不变

## 回滚方案

如需回滚到 PV 统计，可以：
1. 恢复 `repository/visit_stat.go` 中的 `IncrementTodayViewCount()` 方法
2. 恢复 `middleware/visit_stat.go` 中的调用
3. 删除 `visit_records` 表（可选）

---

**修改日期**: 2024-10-27  
**版本**: v1.0

