-- 创建说说表
CREATE TABLE IF NOT EXISTS moments (
  id BIGSERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  images TEXT,
  user_id BIGINT NOT NULL,
  status SMALLINT NOT NULL DEFAULT 1,
  like_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建普通索引
CREATE INDEX IF NOT EXISTS idx_moments_user_id ON moments(user_id);
CREATE INDEX IF NOT EXISTS idx_moments_status ON moments(status);
CREATE INDEX IF NOT EXISTS idx_moments_created_at ON moments(created_at);

-- 添加全文搜索列（使用 tsvector 类型）
ALTER TABLE moments ADD COLUMN IF NOT EXISTS content_tsv tsvector;

-- 创建 GIN 索引用于全文搜索
CREATE INDEX IF NOT EXISTS idx_moments_content_gin ON moments USING gin(content_tsv);

-- 更新现有数据的 tsvector
UPDATE moments SET content_tsv = to_tsvector('english', content) WHERE content_tsv IS NULL;

-- 添加表注释
COMMENT ON TABLE moments IS '说说表';
COMMENT ON COLUMN moments.content IS '说说内容';
COMMENT ON COLUMN moments.images IS '图片URLs（JSON数组格式）';
COMMENT ON COLUMN moments.user_id IS '用户ID';
COMMENT ON COLUMN moments.status IS '状态：1-公开，0-私密，-1-删除';
COMMENT ON COLUMN moments.like_count IS '点赞数';
COMMENT ON COLUMN moments.content_tsv IS '全文搜索向量';

