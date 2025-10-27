-- =============================================================================
-- 博客系统数据库完整初始化脚本
-- =============================================================================
-- 说明：此脚本包含所有数据库表、索引、默认数据的创建
-- 执行顺序：按照表的依赖关系从基础表到关联表依次创建
-- =============================================================================

-- =============================================================================
-- 1. 用户系统
-- =============================================================================

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    bio VARCHAR(500),
    role VARCHAR(20) DEFAULT 'user',
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 用户表注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.email IS '邮箱';
COMMENT ON COLUMN users.password IS '密码（bcrypt加密）';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.avatar IS '头像URL';
COMMENT ON COLUMN users.bio IS '个人简介';
COMMENT ON COLUMN users.role IS '角色：admin-管理员，user-普通用户';
COMMENT ON COLUMN users.status IS '状态：1-正常，0-禁用';

-- =============================================================================
-- 2. 分类和标签系统
-- =============================================================================

-- 创建分类表
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(200),
    color VARCHAR(20),
    sort INT DEFAULT 0,
    post_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 分类表注释
COMMENT ON TABLE categories IS '分类表';
COMMENT ON COLUMN categories.name IS '分类名称';
COMMENT ON COLUMN categories.description IS '分类描述';
COMMENT ON COLUMN categories.color IS '分类颜色';
COMMENT ON COLUMN categories.sort IS '排序';
COMMENT ON COLUMN categories.post_count IS '文章数量';

-- 创建标签表
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    color VARCHAR(20),
    text_color VARCHAR(20),
    font_size INTEGER,
    post_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 标签表注释
COMMENT ON TABLE tags IS '标签表';
COMMENT ON COLUMN tags.name IS '标签名称';
COMMENT ON COLUMN tags.color IS '标签颜色';
COMMENT ON COLUMN tags.text_color IS '文字颜色';
COMMENT ON COLUMN tags.font_size IS '文字大小(px)';
COMMENT ON COLUMN tags.post_count IS '文章数量';

-- =============================================================================
-- 3. 文章系统
-- =============================================================================

-- 创建文章表
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    summary VARCHAR(500),
    cover VARCHAR(255),
    status INT DEFAULT 1,
    is_top BOOLEAN DEFAULT FALSE,
    view_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    user_id INT NOT NULL,
    category_id INT NOT NULL,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    search_tsv tsvector,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

-- 文章表索引
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts(title);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_category_id ON posts(category_id);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

-- 全文搜索索引（使用 GIN 索引用于全文搜索，组合标题和内容）
CREATE INDEX IF NOT EXISTS idx_posts_search_gin ON posts USING gin(search_tsv);

-- 文章表注释
COMMENT ON TABLE posts IS '文章表';
COMMENT ON COLUMN posts.title IS '文章标题';
COMMENT ON COLUMN posts.content IS '文章内容（Markdown格式）';
COMMENT ON COLUMN posts.summary IS '文章摘要';
COMMENT ON COLUMN posts.cover IS '封面图URL';
COMMENT ON COLUMN posts.status IS '状态：1-已发布，0-草稿，-1-删除';
COMMENT ON COLUMN posts.is_top IS '是否置顶';
COMMENT ON COLUMN posts.view_count IS '浏览量';
COMMENT ON COLUMN posts.like_count IS '点赞数';
COMMENT ON COLUMN posts.user_id IS '作者ID';
COMMENT ON COLUMN posts.category_id IS '分类ID';
COMMENT ON COLUMN posts.published_at IS '发布时间';
COMMENT ON COLUMN posts.search_tsv IS '全文搜索向量（标题+内容）';

-- 创建文章标签关联表
CREATE TABLE IF NOT EXISTS post_tags (
    post_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- 文章标签关联表注释
COMMENT ON TABLE post_tags IS '文章标签关联表';

-- =============================================================================
-- 4. 评论系统
-- =============================================================================

-- 创建评论表
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    parent_id INT,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

-- 评论表索引
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);

-- 评论表注释
COMMENT ON TABLE comments IS '评论表';
COMMENT ON COLUMN comments.content IS '评论内容';
COMMENT ON COLUMN comments.post_id IS '文章ID';
COMMENT ON COLUMN comments.user_id IS '评论用户ID';
COMMENT ON COLUMN comments.parent_id IS '父评论ID（用于回复）';
COMMENT ON COLUMN comments.status IS '状态：1-正常，0-待审核，-1-删除';

-- =============================================================================
-- 5. 说说（动态）系统
-- =============================================================================

-- 创建说说表
CREATE TABLE IF NOT EXISTS moments (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    images TEXT,
    user_id BIGINT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    like_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    content_tsv tsvector
);

-- 说说表索引
CREATE INDEX IF NOT EXISTS idx_moments_user_id ON moments(user_id);
CREATE INDEX IF NOT EXISTS idx_moments_status ON moments(status);
CREATE INDEX IF NOT EXISTS idx_moments_created_at ON moments(created_at);

-- 说说全文搜索索引
CREATE INDEX IF NOT EXISTS idx_moments_content_gin ON moments USING gin(content_tsv);

-- 说说表注释
COMMENT ON TABLE moments IS '说说表';
COMMENT ON COLUMN moments.content IS '说说内容';
COMMENT ON COLUMN moments.images IS '图片URLs（JSON数组格式）';
COMMENT ON COLUMN moments.user_id IS '用户ID';
COMMENT ON COLUMN moments.status IS '状态：1-公开，0-私密，-1-删除';
COMMENT ON COLUMN moments.like_count IS '点赞数';
COMMENT ON COLUMN moments.content_tsv IS '全文搜索向量';

-- =============================================================================
-- 6. 统计系统
-- =============================================================================

-- 创建访问量统计表
CREATE TABLE IF NOT EXISTS visit_stats (
    id SERIAL PRIMARY KEY,
    date DATE UNIQUE NOT NULL,
    view_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 访问量统计表索引
CREATE INDEX IF NOT EXISTS idx_visit_stats_date ON visit_stats(date DESC);

-- 访问量统计表注释
COMMENT ON TABLE visit_stats IS '访问量统计表（按日期）';
COMMENT ON COLUMN visit_stats.date IS '统计日期';
COMMENT ON COLUMN visit_stats.view_count IS '当日访问量';

-- 创建文章阅读记录表
CREATE TABLE IF NOT EXISTS post_views (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT,
    ip VARCHAR(45) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 文章阅读记录表索引
CREATE INDEX IF NOT EXISTS idx_post_views_post_id ON post_views(post_id);
CREATE INDEX IF NOT EXISTS idx_post_views_user_id ON post_views(user_id);
CREATE INDEX IF NOT EXISTS idx_post_views_ip ON post_views(ip);
CREATE INDEX IF NOT EXISTS idx_post_views_post_user ON post_views(post_id, user_id);
CREATE INDEX IF NOT EXISTS idx_post_views_post_ip ON post_views(post_id, ip);

-- 文章阅读记录表注释
COMMENT ON TABLE post_views IS '文章阅读记录表（用于去重统计）';
COMMENT ON COLUMN post_views.post_id IS '文章ID';
COMMENT ON COLUMN post_views.user_id IS '用户ID（匿名用户为NULL）';
COMMENT ON COLUMN post_views.ip IS '访客IP地址';
COMMENT ON COLUMN post_views.created_at IS '阅读时间';

-- =============================================================================
-- 7. 系统配置
-- =============================================================================

-- 创建系统配置表
CREATE TABLE IF NOT EXISTS settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,
    value TEXT,
    type VARCHAR(20) DEFAULT 'text',
    "group" VARCHAR(50),
    label VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 系统配置表索引
CREATE INDEX IF NOT EXISTS idx_settings_group ON settings("group");

-- 系统配置表注释
COMMENT ON TABLE settings IS '系统配置表';
COMMENT ON COLUMN settings.key IS '配置键（唯一）';
COMMENT ON COLUMN settings.value IS '配置值';
COMMENT ON COLUMN settings.type IS '配置类型：text-文本，json-JSON，image-图片';
COMMENT ON COLUMN settings."group" IS '配置分组：site-网站，about-关于';
COMMENT ON COLUMN settings.label IS '配置标签（显示名称）';

-- =============================================================================
-- 8. 点赞系统
-- =============================================================================

-- 创建文章点赞记录表
CREATE TABLE IF NOT EXISTS post_likes (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    user_id INTEGER,
    ip VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_post_likes_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT fk_post_likes_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 文章点赞记录表索引
CREATE INDEX IF NOT EXISTS idx_post_likes_post_id ON post_likes(post_id);
CREATE INDEX IF NOT EXISTS idx_post_likes_user_id ON post_likes(user_id);
CREATE INDEX IF NOT EXISTS idx_post_likes_ip ON post_likes(ip);

-- 创建唯一约束，防止同一用户/IP重复点赞
CREATE UNIQUE INDEX IF NOT EXISTS idx_post_likes_unique_user ON post_likes(post_id, user_id) WHERE user_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_post_likes_unique_ip ON post_likes(post_id, ip) WHERE user_id IS NULL;

-- 文章点赞记录表注释
COMMENT ON TABLE post_likes IS '文章点赞记录表';
COMMENT ON COLUMN post_likes.post_id IS '文章ID';
COMMENT ON COLUMN post_likes.user_id IS '用户ID（已登录用户）';
COMMENT ON COLUMN post_likes.ip IS 'IP地址（匿名用户）';

-- 创建说说点赞记录表
CREATE TABLE IF NOT EXISTS moment_likes (
    id SERIAL PRIMARY KEY,
    moment_id BIGINT NOT NULL,
    user_id BIGINT,
    ip VARCHAR(45) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(moment_id, user_id),
    UNIQUE(moment_id, ip)
);

-- 说说点赞记录表索引
CREATE INDEX IF NOT EXISTS idx_moment_likes_moment_id ON moment_likes(moment_id);
CREATE INDEX IF NOT EXISTS idx_moment_likes_user_id ON moment_likes(user_id);
CREATE INDEX IF NOT EXISTS idx_moment_likes_ip ON moment_likes(ip);

-- 说说点赞记录表注释
COMMENT ON TABLE moment_likes IS '说说点赞记录表';
COMMENT ON COLUMN moment_likes.moment_id IS '说说ID';
COMMENT ON COLUMN moment_likes.user_id IS '用户ID（匿名用户为NULL）';
COMMENT ON COLUMN moment_likes.ip IS '用户IP地址';
COMMENT ON COLUMN moment_likes.created_at IS '点赞时间';

-- =============================================================================
-- 9. IP 黑名单系统
-- =============================================================================

-- 创建IP黑名单表
CREATE TABLE IF NOT EXISTS ip_blacklist (
    id SERIAL PRIMARY KEY,
    ip VARCHAR(45) UNIQUE NOT NULL,
    reason VARCHAR(255),
    ban_type SMALLINT DEFAULT 1,
    expire_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- IP黑名单表索引
CREATE INDEX IF NOT EXISTS idx_ip_blacklist_ip ON ip_blacklist(ip);
CREATE INDEX IF NOT EXISTS idx_ip_blacklist_expire_at ON ip_blacklist(expire_at);

-- IP黑名单表注释
COMMENT ON TABLE ip_blacklist IS 'IP黑名单表';
COMMENT ON COLUMN ip_blacklist.ip IS 'IP地址';
COMMENT ON COLUMN ip_blacklist.reason IS '封禁原因';
COMMENT ON COLUMN ip_blacklist.ban_type IS '封禁类型：1-自动封禁，2-手动封禁';
COMMENT ON COLUMN ip_blacklist.expire_at IS '过期时间，NULL表示永久封禁';

-- =============================================================================
-- 10. 初始化默认数据
-- =============================================================================

-- 插入默认管理员用户
-- 用户名：admin
-- 密码：password （实际使用时请修改）
-- 密码 hash: $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
INSERT INTO users (username, email, password, nickname, avatar, bio, role, status, created_at, updated_at)
VALUES 
('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '管理员', '', '博客管理员', 'admin', 1, NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- 插入默认分类
INSERT INTO categories (name, description, color, sort, post_count, created_at, updated_at)
VALUES 
('技术', '技术文章', '#2196F3', 1, 0, NOW(), NOW()),
('生活', '生活随笔', '#4CAF50', 2, 0, NOW(), NOW()),
('思考', '思考感悟', '#FF9800', 3, 0, NOW(), NOW()),
('教程', '教程文档', '#9C27B0', 4, 0, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- 插入默认标签
INSERT INTO tags (name, color, text_color, font_size, post_count, created_at, updated_at)
VALUES 
('Go', '#00ADD8', NULL, NULL, 0, NOW(), NOW()),
('Vue', '#42b883', NULL, NULL, 0, NOW(), NOW()),
('TypeScript', '#3178c6', NULL, NULL, 0, NOW(), NOW()),
('PostgreSQL', '#336791', NULL, NULL, 0, NOW(), NOW()),
('Docker', '#2496ED', NULL, NULL, 0, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- 插入网站配置（关于页面）
INSERT INTO settings (key, value, type, "group", label, created_at, updated_at)
VALUES 
('about_title', '👋 你好', 'text', 'about', '关于标题', NOW(), NOW()),
('about_intro', '欢迎来到我的个人博客！这里记录了我在技术学习旅程中的点点滴滴。', 'text', 'about', '个人简介', NOW(), NOW()),
('about_avatar', '', 'image', 'about', '个人头像', NOW(), NOW()),
('about_skills', '["Vue 3","Go","TypeScript","PostgreSQL","Docker"]', 'json', 'about', '技术栈', NOW(), NOW()),
('about_email', 'your-email@example.com', 'text', 'about', '联系邮箱', NOW(), NOW()),
('about_github', 'github.com/yourname', 'text', 'about', 'GitHub', NOW(), NOW()),
('about_site_intro', '本站基于 Vue 3 + Go 构建，采用前后端分离架构。使用 Naive UI 组件库，支持 Markdown 写作。

如果你觉得这个博客不错，欢迎 Star 或 Fork 源码！', 'text', 'about', '关于本站', NOW(), NOW()),
('site_name', '我的博客', 'text', 'site', '网站名称', NOW(), NOW()),
('site_icp', '', 'text', 'site', 'ICP备案号', NOW(), NOW()),
('site_police', '', 'text', 'site', '公安备案号', NOW(), NOW())
ON CONFLICT (key) DO NOTHING;

-- 插入最近7天的访问统计初始记录
INSERT INTO visit_stats (date, view_count, created_at, updated_at)
SELECT 
    CURRENT_DATE - INTERVAL '1 day' * generate_series(0, 6) AS date,
    0 AS view_count,
    NOW() AS created_at,
    NOW() AS updated_at
ON CONFLICT (date) DO NOTHING;

-- =============================================================================
-- 11. 更新现有数据的全文搜索向量
-- =============================================================================

-- 更新文章的全文搜索向量（组合标题和内容，标题权重更高）
UPDATE posts 
SET search_tsv = 
    setweight(to_tsvector('english', coalesce(title, '')), 'A') || 
    setweight(to_tsvector('english', coalesce(content, '')), 'B')
WHERE search_tsv IS NULL;

-- 更新说说的全文搜索向量
UPDATE moments 
SET content_tsv = to_tsvector('english', content) 
WHERE content_tsv IS NULL;

-- =============================================================================
-- 初始化完成
-- =============================================================================
-- 说明：
-- 1. 默认管理员账号：admin / password（请首次登录后修改）
-- 2. 全文搜索使用 PostgreSQL 的 tsvector 和 GIN 索引
-- 3. 应用层更新文章/说说时，需要同时更新 search_tsv/content_tsv 字段
-- 4. 文章阅读记录用于去重统计，避免同一用户/IP重复计数
-- =============================================================================
