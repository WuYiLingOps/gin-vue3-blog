package main

import (
	"fmt"
	"log"
	"regexp"

	"blog-backend/config"
	"blog-backend/db"
	"blog-backend/logger"
	"blog-backend/model"
	"blog-backend/util"
)

func main() {
	// 加载配置
	if err := config.LoadConfigByEnv(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 初始化日志
	isDev := config.Cfg.Env == "dev"
	if err := logger.InitLogger(config.Cfg.Log.Level, isDev); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer logger.Sync()

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 获取所有文章
	var posts []model.Post
	if err := db.DB.Find(&posts).Error; err != nil {
		log.Fatalf("获取文章列表失败: %v", err)
	}

	fmt.Printf("找到 %d 篇文章，开始生成slug...\n", len(posts))

	updated := 0
	for _, post := range posts {
		// 检查是否需要更新slug
		// 1. 如果slug为空，需要生成
		// 2. 如果是临时的 "post-XX" 格式，需要重新生成
		// 3. 如果slug看起来不像有效的拼音slug（包含特殊字符、格式不对等），也需要重新生成
		needUpdate := false
		updateReason := ""

		if post.Slug == "" {
			needUpdate = true
			updateReason = "slug为空"
		} else if matched, _ := regexp.MatchString(`^post-\d+$`, post.Slug); matched {
			// 如果是 "post-XX" 格式，需要重新生成
			needUpdate = true
			updateReason = "临时格式 post-XX"
		} else {
			// 检查slug是否看起来像有效的拼音slug
			// 有效的拼音slug应该只包含：小写字母、数字、连字符
			// 且不应该以连字符开头或结尾
			validSlugPattern := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
			if !validSlugPattern.MatchString(post.Slug) {
				needUpdate = true
				updateReason = "slug格式无效（包含特殊字符）"
			} else if len(post.Slug) > 200 {
				// slug太长也不合理
				needUpdate = true
				updateReason = "slug过长"
			}
		}

		if !needUpdate {
			fmt.Printf("跳过文章 ID %d: %s (slug: %s)\n", post.ID, post.Title, post.Slug)
			continue
		}

		// 生成slug
		baseSlug := util.GenerateSlug(post.Title)
		if baseSlug == "" {
			log.Printf("警告: 文章 ID %d 无法生成slug，使用默认值", post.ID)
			baseSlug = fmt.Sprintf("post-%d", post.ID)
		}

		slug := baseSlug
		counter := 1

		// 确保slug唯一
		for {
			var count int64
			db.DB.Model(&model.Post{}).Where("slug = ? AND id != ?", slug, post.ID).Count(&count)
			if count == 0 {
				break
			}
			slug = baseSlug + "-" + fmt.Sprintf("%d", counter)
			counter++
			if counter > 100 {
				log.Printf("警告: 文章 ID %d 的slug生成达到最大重试次数", post.ID)
				slug = fmt.Sprintf("post-%d", post.ID)
				break
			}
		}

		// 更新文章
		oldSlug := post.Slug
		post.Slug = slug
		if err := db.DB.Save(&post).Error; err != nil {
			log.Printf("更新文章 ID %d 失败: %v", post.ID, err)
			continue
		}

		updated++
		if oldSlug == "" {
			fmt.Printf("✓ 文章 ID %d: %s -> %s (%s)\n", post.ID, post.Title, slug, updateReason)
		} else {
			fmt.Printf("✓ 文章 ID %d: %s [%s -> %s] (%s)\n", post.ID, post.Title, oldSlug, slug, updateReason)
		}
	}

	fmt.Printf("\n完成！共更新 %d 篇文章的slug\n", updated)
}
