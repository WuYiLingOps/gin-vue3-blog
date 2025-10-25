-- 为posts表添加全文搜索支持

-- 添加全文搜索列（使用 tsvector 类型）
ALTER TABLE posts ADD COLUMN IF NOT EXISTS search_tsv tsvector;

-- 创建 GIN 索引用于全文搜索（组合标题和内容）
CREATE INDEX IF NOT EXISTS idx_posts_search_gin ON posts USING gin(search_tsv);

-- 更新现有数据的 tsvector（组合标题和内容，标题权重更高）
UPDATE posts 
SET search_tsv = 
    setweight(to_tsvector('english', coalesce(title, '')), 'A') || 
    setweight(to_tsvector('english', coalesce(content, '')), 'B')
WHERE search_tsv IS NULL;

-- 添加列注释
COMMENT ON COLUMN posts.search_tsv IS '全文搜索向量（标题+内容）';

-- 注意：如果需要自动更新search_tsv，可以在应用层处理
-- 或者使用触发器（如果不介意性能影响）
-- 对于博客系统，文章更新不频繁，可以在应用层的Update方法中更新此字段

