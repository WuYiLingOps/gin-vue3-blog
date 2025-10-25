-- 添加文章阅读记录表的迁移脚本
-- 用于已有数据库升级

-- 创建文章阅读记录表
CREATE TABLE IF NOT EXISTS post_views (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT,
    ip VARCHAR(45) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_post_views_post_id ON post_views(post_id);
CREATE INDEX IF NOT EXISTS idx_post_views_user_id ON post_views(user_id);
CREATE INDEX IF NOT EXISTS idx_post_views_ip ON post_views(ip);

-- 组合索引用于快速查询是否已阅读
CREATE INDEX IF NOT EXISTS idx_post_views_post_user ON post_views(post_id, user_id);
CREATE INDEX IF NOT EXISTS idx_post_views_post_ip ON post_views(post_id, ip);

-- 说明：
-- 1. post_id: 文章ID，关联 posts 表
-- 2. user_id: 用户ID，关联 users 表（可为 NULL，表示匿名用户）
-- 3. ip: 访客IP地址，用于匿名用户去重
-- 4. created_at: 阅读时间

-- 使用示例：
-- 1. 检查登录用户是否已阅读：
--    SELECT COUNT(*) FROM post_views WHERE post_id = ? AND user_id = ?;
-- 
-- 2. 检查匿名用户是否已阅读：
--    SELECT COUNT(*) FROM post_views WHERE post_id = ? AND ip = ? AND user_id IS NULL;
-- 
-- 3. 记录阅读：
--    INSERT INTO post_views (post_id, user_id, ip) VALUES (?, ?, ?);

