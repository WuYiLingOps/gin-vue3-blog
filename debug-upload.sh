#!/bin/bash
# 图片上传问题排查脚本

echo "========== 博客图片上传问题排查 =========="
echo ""

echo "1. 检查 Docker 容器状态..."
docker ps | grep blog

echo ""
echo "2. 检查后端容器内的上传目录..."
docker exec blog-backend ls -la /app/uploads/avatars/ 2>/dev/null || echo "❌ 无法访问容器或目录不存在"

echo ""
echo "3. 检查 Docker volume..."
docker volume ls | grep uploads

echo ""
echo "4. 测试后端是否能直接访问静态文件..."
echo "   访问: http://localhost:8080/uploads/avatars/"
docker exec blog-backend ls /app/uploads/avatars/ | head -1 | xargs -I {} echo "   测试文件: {}"
docker exec blog-backend ls /app/uploads/avatars/ | head -1 | xargs -I {} curl -I http://localhost:8080/uploads/avatars/{}

echo ""
echo "5. 检查后端日志..."
docker logs blog-backend --tail 20

echo ""
echo "========== 排查完成 =========="
echo ""
echo "💡 如果上传目录为空，请尝试以下解决方案："
echo "   1. 修改 docker-compose.yml，将 volume 改为本地目录："
echo "      volumes:"
echo "        - ./uploads:/app/uploads"
echo ""
echo "   2. 然后重启容器："
echo "      docker compose down && docker compose up -d"

