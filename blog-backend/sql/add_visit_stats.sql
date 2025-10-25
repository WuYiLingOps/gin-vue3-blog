-- 添加访问量统计表的迁移脚本
-- 如果你的数据库已经初始化，运行此脚本以添加 visit_stats 表

-- 创建访问量统计表
CREATE TABLE IF NOT EXISTS visit_stats (
    id SERIAL PRIMARY KEY,
    date DATE UNIQUE NOT NULL,
    view_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_visit_stats_date ON visit_stats(date DESC);

-- 插入最近7天的初始记录（访问量为0）
INSERT INTO visit_stats (date, view_count, created_at, updated_at)
SELECT 
    CURRENT_DATE - INTERVAL '1 day' * generate_series(0, 6) AS date,
    0 AS view_count,
    NOW() AS created_at,
    NOW() AS updated_at
ON CONFLICT (date) DO NOTHING;

-- 查看已插入的记录
SELECT * FROM visit_stats ORDER BY date DESC;

