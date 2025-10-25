-- 博客系统数据库初始化脚本

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

-- 创建标签表
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    color VARCHAR(20),
    post_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

-- 创建文章标签关联表
CREATE TABLE IF NOT EXISTS post_tags (
    post_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

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

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts(title);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_category_id ON posts(category_id);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);

-- 插入默认管理员用户（密码：admin123，需要在首次登录后修改）
-- 注意：密码 hash 需要使用 bcrypt 生成，这里使用 $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi (password)
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
INSERT INTO tags (name, color, post_count, created_at, updated_at)
VALUES 
('Go', '#00ADD8', 0, NOW(), NOW()),
('Vue', '#42b883', 0, NOW(), NOW()),
('TypeScript', '#3178c6', 0, NOW(), NOW()),
('PostgreSQL', '#336791', 0, NOW(), NOW()),
('Docker', '#2496ED', 0, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

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

CREATE INDEX IF NOT EXISTS idx_settings_group ON settings("group");

-- 创建访问量统计表
CREATE TABLE IF NOT EXISTS visit_stats (
    id SERIAL PRIMARY KEY,
    date DATE UNIQUE NOT NULL,
    view_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_visit_stats_date ON visit_stats(date DESC);

-- 创建文章阅读记录表
CREATE TABLE IF NOT EXISTS post_views (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT,
    ip VARCHAR(45) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_post_views_post_id ON post_views(post_id);
CREATE INDEX IF NOT EXISTS idx_post_views_user_id ON post_views(user_id);
CREATE INDEX IF NOT EXISTS idx_post_views_ip ON post_views(ip);
CREATE INDEX IF NOT EXISTS idx_post_views_post_user ON post_views(post_id, user_id);
CREATE INDEX IF NOT EXISTS idx_post_views_post_ip ON post_views(post_id, ip);

-- 插入关于页面的默认配置
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

