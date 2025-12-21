package service

import (
	"errors"
	"fmt"
	"strings"

	"blog-backend/config"
	"blog-backend/model"
	"blog-backend/repository"
	"blog-backend/util"

	"gorm.io/gorm"
)

type CommentService struct {
	repo        *repository.CommentRepository
	postRepo    *repository.PostRepository
	momentRepo  *repository.MomentRepository
	userRepo    *repository.UserRepository
	settingRepo *repository.SettingRepository
}

func NewCommentService() *CommentService {
	return &CommentService{
		repo:        repository.NewCommentRepository(),
		postRepo:    repository.NewPostRepository(),
		momentRepo:  repository.NewMomentRepository(),
		userRepo:    repository.NewUserRepository(),
		settingRepo: repository.NewSettingRepository(),
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content     string `json:"content" binding:"required"`
	CommentType string `json:"comment_type"` // 评论类型：post-文章评论，friendlink-友链评论，moment-说说评论
	PostID      *uint  `json:"post_id"`      // 文章ID（文章评论时使用）
	TargetID    *uint  `json:"target_id"`    // 目标ID（友链评论时使用，可以为0或友链ID）
	ParentID    *uint  `json:"parent_id"`    // 父评论ID（用于回复）
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content"`
}

// Create 创建评论
func (s *CommentService) Create(userID uint, req *CreateCommentRequest) (*model.Comment, error) {
	// 确定评论类型（默认为post，向后兼容）
	commentType := req.CommentType
	if commentType == "" {
		commentType = "post"
	}

	// 根据评论类型进行验证
	if commentType == "post" {
		// 文章评论：必须提供post_id，且文章必须存在
		if req.PostID == nil || *req.PostID == 0 {
			return nil, errors.New("文章ID不能为空")
		}
		if _, err := s.postRepo.GetByID(*req.PostID); err != nil {
			return nil, errors.New("文章不存在")
		}
	} else if commentType == "friendlink" {
		// 友链评论：不需要post_id，使用target_id（可以为0表示友链页面）
		if req.TargetID == nil {
			targetID := uint(0)
			req.TargetID = &targetID
		}
	} else if commentType == "moment" {
		// 说说评论：必须提供target_id（说说ID）
		if req.TargetID == nil || *req.TargetID == 0 {
			return nil, errors.New("说说ID不能为空")
		}
		// 注意：这里不验证说说是否存在，因为可能没有 moment repository
		// 如果需要验证，可以添加 moment repository
	} else {
		return nil, errors.New("不支持的评论类型")
	}

	// 如果是回复评论，检查父评论是否存在
	if req.ParentID != nil {
		if _, err := s.repo.GetByID(*req.ParentID); err != nil {
			return nil, errors.New("父评论不存在")
		}
	}

	comment := &model.Comment{
		Content:     req.Content,
		CommentType: commentType,
		PostID:      req.PostID,
		TargetID:    req.TargetID,
		UserID:      userID,
		ParentID:    req.ParentID,
		Status:      1,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, errors.New("评论创建失败")
	}

	// 获取完整的评论信息（包含用户信息）
	createdComment, err := s.repo.GetByID(comment.ID)
	if err != nil {
		return nil, errors.New("获取评论失败")
	}

	return createdComment, nil
}

// CreateWithContext 创建评论（带上下文，用于获取请求信息）
func (s *CommentService) CreateWithContext(userID uint, req *CreateCommentRequest, siteURL string) (*model.Comment, error) {
	comment, err := s.Create(userID, req)
	if err != nil {
		return nil, err
	}

	// 异步发送通知邮件（不阻塞请求）
	go s.sendCommentNotifications(comment, userID, siteURL)

	return comment, nil
}

// sendCommentNotifications 发送评论通知邮件
func (s *CommentService) sendCommentNotifications(comment *model.Comment, commenterID uint, requestSiteURL string) {
	// 检查邮件配置是否完整
	if config.Cfg.Email.Host == "" || config.Cfg.Email.Username == "" {
		return // 邮件未配置，不发送通知
	}

	// 获取网站名称
	siteName := s.getSiteName()

	emailConfig := util.EmailConfig{
		Host:     config.Cfg.Email.Host,
		Port:     config.Cfg.Email.Port,
		Username: config.Cfg.Email.Username,
		Password: config.Cfg.Email.Password,
		FromName: config.Cfg.Email.FromName,
		SiteName: siteName,
	}

	// 获取评论者信息
	commenter, err := s.userRepo.GetByID(commenterID)
	if err != nil {
		return
	}
	commenterName := commenter.Nickname
	if commenterName == "" {
		commenterName = commenter.Username
	}

	// 获取网站URL（用于构建文章链接）
	// 优先使用请求中的URL，其次使用数据库配置，最后使用默认值
	siteURL := s.getSiteURL(requestSiteURL)

	// 处理文章评论的通知
	if comment.CommentType == "post" && comment.PostID != nil {
		// 获取文章信息
		post, err := s.postRepo.GetByID(*comment.PostID)
		if err != nil {
			return
		}

		// 构建文章URL
		postURL := fmt.Sprintf("%s/post/%d", siteURL, post.ID)
		commentPreview := s.stripMarkdown(comment.Content)

		// 通知管理员（如果启用了管理员通知）
		// 需要单独配置开关，开启后所有管理员都会收到通知
		// 注意：由于普通用户没有权限写文章，文章作者只能是管理员，因此统一通过管理员通知处理
		if s.shouldNotifyAdmin() {
			admins, err := s.userRepo.GetAdmins()
			if err == nil {
				for _, admin := range admins {
					// 不通知评论者本人（如果评论者是管理员）
					if admin.ID != commenterID && admin.Email != "" {
						if err := util.SendAdminCommentNotificationEmail(
							emailConfig,
							admin.Email,
							commenterName,
							post.Title,
							commentPreview,
							postURL,
						); err != nil {
							// 记录错误，但不影响主流程（邮件发送是异步的）
							fmt.Printf("发送评论通知邮件给管理员失败: %v\n", err)
						}
					}
				}
			}
		}
	} else if comment.CommentType == "moment" && comment.TargetID != nil && *comment.TargetID > 0 {
		// 处理说说评论的通知
		// 获取说说信息
		moment, err := s.momentRepo.GetByID(*comment.TargetID)
		if err != nil {
			return
		}

		// 构建说说页面URL
		momentURL := fmt.Sprintf("%s/moments", siteURL)
		commentPreview := s.stripMarkdown(comment.Content)

		// 截取说说内容作为标题（最多50个字符）
		momentTitle := moment.Content
		if len([]rune(momentTitle)) > 50 {
			momentTitle = string([]rune(momentTitle)[:50]) + "..."
		}
		if momentTitle == "" {
			momentTitle = "说说"
		}

		// 通知管理员（如果启用了管理员通知）
		if s.shouldNotifyAdmin() {
			admins, err := s.userRepo.GetAdmins()
			if err == nil {
				for _, admin := range admins {
					// 不通知评论者本人（如果评论者是管理员）
					if admin.ID != commenterID && admin.Email != "" {
						if err := util.SendAdminCommentNotificationEmail(
							emailConfig,
							admin.Email,
							commenterName,
							"说说："+momentTitle,
							commentPreview,
							momentURL,
						); err != nil {
							// 记录错误，但不影响主流程（邮件发送是异步的）
							fmt.Printf("发送说说评论通知邮件给管理员失败: %v\n", err)
						}
					}
				}
			}
		}
	}
}

// shouldNotifyAdmin 检查是否应该通知管理员
func (s *CommentService) shouldNotifyAdmin() bool {
	setting, err := s.settingRepo.GetByKey("notify_admin_on_comment")
	if err != nil {
		// 如果配置不存在，默认不通知
		return false
	}
	// 配置值为 "1" 或 "true" 时通知
	return setting.Value == "1" || setting.Value == "true"
}

// getSiteURL 获取网站URL
// 优先级：我的友链信息中的url > 数据库配置的site_url > 默认值
// 注意：不使用请求头中的 Host，因为那是后端地址，不是前端地址
func (s *CommentService) getSiteURL(requestSiteURL string) string {
	// 如果请求URL不为空，检查是否是后端地址（包含常见后端端口）
	// 如果是后端地址，则忽略，使用数据库配置
	if requestSiteURL != "" {
		// 检查是否包含后端常见端口（8080, 8081, 3001等）
		// 如果包含，说明是后端地址，应该忽略
		if !containsBackendPort(requestSiteURL) {
			return requestSiteURL
		}
		// 如果是后端地址，继续使用数据库配置
	}

	// 1. 优先从"我的友链信息"中获取url
	friendLinkInfoSettings, err := s.settingRepo.GetByGroup("friendlink_info")
	if err == nil {
		// 查找url字段
		for _, setting := range friendLinkInfoSettings {
			if setting.Key == "url" && setting.Value != "" {
				// 如果友链信息中的url是默认值，则使用默认值http://localhost:3000
				if setting.Value == "https://xxxxx.cn/" {
					return "http://localhost:3000"
				}
				// 否则使用友链信息中的url
				return setting.Value
			}
		}
	}

	// 2. 其次从数据库配置的site_url读取
	setting, err := s.settingRepo.GetByKey("site_url")
	if err == nil && setting.Value != "" {
		return setting.Value
	}

	// 3. 返回默认值（开发环境默认值）
	return "http://localhost:3000"
}

// containsBackendPort 检查URL是否包含后端常见端口
func containsBackendPort(url string) bool {
	backendPorts := []string{":8080", ":8081", ":3001", ":8000", ":8001"}
	for _, port := range backendPorts {
		if strings.Contains(url, port) {
			return true
		}
	}
	return false
}

// stripMarkdown 去除Markdown格式，只保留纯文本（简单实现）
func (s *CommentService) stripMarkdown(content string) string {
	// 简单的Markdown去除：移除代码块、链接等
	// 这里只做基本处理，更复杂的可以用专门的库
	result := content
	// 简单的正则替换（移除常见的Markdown标记）
	// 注意：这是一个简化版本，对于复杂的Markdown可能需要专门的库
	// 移除代码块 ```code```
	// 移除行内代码 `code`
	// 移除链接 [text](url) -> text
	// 移除图片 ![alt](url) -> alt
	// 移除粗体 **text** -> text
	// 移除斜体 *text* -> text
	// 这里我们只做基本的清理，保留原内容的前200个字符作为预览
	if len([]rune(result)) > 200 {
		result = string([]rune(result)[:200]) + "..."
	}
	return result
}

// GetByID 获取评论详情
func (s *CommentService) GetByID(id uint) (*model.Comment, error) {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		return nil, errors.New("获取评论失败")
	}
	return comment, nil
}

// Update 更新评论
func (s *CommentService) Update(id, userID uint, role string, req *UpdateCommentRequest) (*model.Comment, error) {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("评论不存在")
	}

	// 权限检查：只有作者和管理员可以修改
	if comment.UserID != userID && role != "admin" {
		return nil, errors.New("无权限修改此评论")
	}

	if req.Content != "" {
		comment.Content = req.Content
	}

	if err := s.repo.Update(comment); err != nil {
		return nil, errors.New("评论更新失败")
	}

	return comment, nil
}

// Delete 删除评论
func (s *CommentService) Delete(id, userID uint, role string) error {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 权限检查：只有作者和管理员可以删除
	if comment.UserID != userID && role != "admin" {
		return errors.New("无权限删除此评论")
	}

	return s.repo.Delete(id)
}

// GetByPostID 获取文章的评论列表（向后兼容）
func (s *CommentService) GetByPostID(postID uint) ([]model.Comment, error) {
	return s.repo.GetByPostID(postID)
}

// GetByTypeAndTarget 根据评论类型和目标ID获取评论列表
func (s *CommentService) GetByTypeAndTarget(commentType string, targetID uint) ([]model.Comment, error) {
	return s.repo.GetByTypeAndTarget(commentType, targetID)
}

// List 获取评论列表（管理后台）
func (s *CommentService) List(page, pageSize int) ([]model.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.List(page, pageSize)
}

// UpdateStatus 更新评论状态
func (s *CommentService) UpdateStatus(id uint, status int) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.New("评论不存在")
	}

	return s.repo.UpdateStatus(id, status)
}

// getSiteName 获取网站名称
func (s *CommentService) getSiteName() string {
	settings, err := s.settingRepo.GetByGroup("site")
	if err != nil {
		return ""
	}

	for _, setting := range settings {
		if setting.Key == "site_name" {
			return setting.Value
		}
	}

	return ""
}
