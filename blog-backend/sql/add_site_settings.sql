-- 添加网站设置的迁移脚本

INSERT INTO settings (key, value, type, "group", label, created_at, updated_at)
VALUES 
('site_name', '我的博客', 'text', 'site', '网站名称', NOW(), NOW()),
('site_icp', '', 'text', 'site', 'ICP备案号', NOW(), NOW()),
('site_police', '', 'text', 'site', '公安备案号', NOW(), NOW())
ON CONFLICT (key) DO NOTHING;

