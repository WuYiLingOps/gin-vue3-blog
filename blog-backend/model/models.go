/*
 * 项目名称：blog-backend
 * 文件名称：models.go
 * 创建时间：2026-01-31 16:26:19
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：数据模型定义文件，包含所有数据库表的模型结构体和关联关系
 */
package model

import (
	"time"
)

// User 用户模型
// 功能说明：存储系统用户信息，包括管理员和普通用户
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password  string    `json:"-" gorm:"not null"`
	Nickname  string    `json:"nickname" gorm:"size:50"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	Bio       string    `json:"bio" gorm:"size:500"`
	Role      string    `json:"role" gorm:"default:user;size:20"` // admin, user
	Status    int       `json:"status" gorm:"default:1"`          // 1:正常 0:禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Post 文章模型
// 功能说明：存储博客文章信息，支持草稿、发布、删除等状态，支持公开和私密可见性
type Post struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null;size:200;index"`
	Slug        string     `json:"slug" gorm:"size:255;uniqueIndex"` // URL友好的标识符（拼音）
	Content     string     `json:"content" gorm:"type:text"`
	Summary     string     `json:"summary" gorm:"size:500"`
	Cover       string     `json:"cover" gorm:"size:255"`
	Status      int        `json:"status" gorm:"default:1;index"`     // 1:发布 0:草稿 -1:删除
	Visibility  int        `json:"visibility" gorm:"default:1;index"` // 1:公开 0:私密
	IsTop       bool       `json:"is_top" gorm:"default:false"`
	ViewCount   int        `json:"view_count" gorm:"default:0"`
	LikeCount   int        `json:"like_count" gorm:"default:0"`
	UserID      uint       `json:"user_id" gorm:"index"`
	CategoryID  uint       `json:"category_id" gorm:"index"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联关系
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Category Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Tags     []Tag     `json:"tags" gorm:"many2many:post_tags;"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	Liked    bool      `json:"liked" gorm:"-"` // 当前用户是否点赞（不存储到数据库）
}

// Category 分类模型
// 功能说明：存储文章分类信息，用于对文章进行分类管理
type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Description string    `json:"description" gorm:"size:200"`
	Color       string    `json:"color" gorm:"size:20"`
	Sort        int       `json:"sort" gorm:"default:0"`
	PostCount   int       `json:"post_count" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag 标签模型
// 功能说明：存储文章标签信息，支持自定义颜色和样式，用于文章标签云展示
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Color     string    `json:"color" gorm:"size:20"`
	TextColor *string   `json:"text_color" gorm:"size:20"` // 文字颜色，可选
	FontSize  *int      `json:"font_size"`                 // 文字大小，可选
	PostCount int       `json:"post_count" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Comment 评论模型
// 功能说明：存储评论信息，支持文章评论和友链评论，支持评论回复（父子关系）
type Comment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Content     string    `json:"content" gorm:"not null;type:text"`
	CommentType string    `json:"comment_type" gorm:"default:post;size:20;index"` // 评论类型：post-文章评论，friendlink-友链评论
	PostID      *uint     `json:"post_id" gorm:"index"`                           // 文章ID（文章评论时使用，友链评论时为NULL）
	TargetID    *uint     `json:"target_id" gorm:"index"`                         // 目标ID（通用目标ID，根据comment_type不同含义不同）
	UserID      uint      `json:"user_id" gorm:"index"`
	ParentID    *uint     `json:"parent_id" gorm:"index"`  // 父评论ID，用于回复
	Status      int       `json:"status" gorm:"default:1"` // 1:正常 0:隐藏
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Post     *Post     `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Comment `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// Setting 系统配置模型
// 功能说明：存储系统配置信息，采用键值对形式，支持分组管理（网站配置、上传配置等）
type Setting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"uniqueIndex;not null;size:100"`
	Value     string    `json:"value" gorm:"type:text"`
	Type      string    `json:"type" gorm:"size:20"`        // text, json, image
	Group     string    `json:"group" gorm:"size:50;index"` // about, site, etc.
	Label     string    `json:"label" gorm:"size:100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostView 文章阅读记录模型
// 功能说明：记录文章访问记录，用于统计文章阅读量和访客信息
type PostView struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"index;not null"`
	UserID    *uint     `json:"user_id" gorm:"index"`    // 已登录用户ID，可为空
	IP        string    `json:"ip" gorm:"size:45;index"` // 访客IP地址
	CreatedAt time.Time `json:"created_at"`
}

// Moment 说说模型
// 功能说明：存储说说（动态）信息，支持公开和私密状态，支持图片上传
type Moment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	Images    string    `json:"images" gorm:"type:text"` // JSON数组格式存储图片URLs
	UserID    uint      `json:"user_id" gorm:"index"`
	Status    int       `json:"status" gorm:"index"` // 1:公开 0:私密 -1:删除（默认值在业务层处理）
	LikeCount int       `json:"like_count" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联关系
	User  User `json:"user" gorm:"foreignKey:UserID"`
	Liked bool `json:"liked" gorm:"-"` // 当前用户是否点赞（不存储到数据库）
}

// MomentLike 说说点赞记录模型
// 功能说明：记录说说的点赞信息，支持登录用户和匿名用户点赞
type MomentLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MomentID  uint      `json:"moment_id" gorm:"index;not null"`
	UserID    *uint     `json:"user_id" gorm:"index"`    // 已登录用户ID，可为空
	IP        string    `json:"ip" gorm:"size:45;index"` // 访客IP地址
	CreatedAt time.Time `json:"created_at"`
}

// PostLike 文章点赞记录模型
// 功能说明：记录文章的点赞信息，支持登录用户和匿名用户点赞
type PostLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"index;not null"`
	UserID    *uint     `json:"user_id" gorm:"index"`    // 已登录用户ID，可为空
	IP        string    `json:"ip" gorm:"size:45;index"` // 访客IP地址
	CreatedAt time.Time `json:"created_at"`
}

// IPBlacklist IP黑名单模型
// 功能说明：存储被封禁的IP地址信息，支持自动封禁和手动封禁，支持过期时间设置
type IPBlacklist struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	IP        string     `json:"ip" gorm:"column:ip;uniqueIndex;not null;size:45"` // 显式指定列名，确保正确映射
	Reason    string     `json:"reason" gorm:"size:255"`
	BanType   int        `json:"ban_type" gorm:"default:1"` // 1:自动封禁 2:手动封禁
	ExpireAt  *time.Time `json:"expire_at"`                 // 过期时间，NULL表示永久封禁
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// IPWhitelist IP白名单模型
// 功能说明：存储白名单IP地址信息，支持CIDR格式，支持过期时间设置
type IPWhitelist struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	IP        string     `json:"ip" gorm:"column:ip;uniqueIndex;not null;size:45"` // 支持 CIDR 格式，显式指定列名
	Reason    string     `json:"reason" gorm:"size:255"`                           // 添加原因
	ExpireAt  *time.Time `json:"expire_at"`                                        // 过期时间，NULL表示永久有效
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// PasswordResetToken 密码重置令牌模型
// 功能说明：存储密码重置和注册验证的令牌信息，包含验证码和过期时间
type PasswordResetToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id" gorm:"index"` // 注册时为NULL，密码重置时为实际用户ID
	Email     string    `json:"email" gorm:"size:100;index;not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;size:100;not null"`
	Code      string    `json:"code" gorm:"size:6;not null"` // 6位验证码
	ExpireAt  time.Time `json:"expire_at" gorm:"not null"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

// EmailChangeRecord 邮箱修改记录模型
// 功能说明：记录用户邮箱修改历史，用于追踪和审计
type EmailChangeRecord struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	OldEmail  string    `json:"old_email" gorm:"size:100;not null"`
	NewEmail  string    `json:"new_email" gorm:"size:100;not null"`
	ChangedAt time.Time `json:"changed_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// ChatMessage 聊天消息模型
// 功能说明：存储聊天室消息和系统公告信息，支持匿名用户和登录用户
type ChatMessage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Content     string    `json:"content" gorm:"not null;type:text"`
	UserID      *uint     `json:"user_id" gorm:"index"`                       // 登录用户ID，可为空（匿名用户）
	Username    string    `json:"username" gorm:"size:50;not null"`           // 用户名（登录用户为真实用户名，匿名用户为临时昵称）
	Avatar      string    `json:"avatar" gorm:"size:255"`                     // 头像URL
	IP          string    `json:"ip" gorm:"size:45"`                          // IP地址
	Priority    int       `json:"priority" gorm:"default:0"`                  // 优先级：0-普通，1-置顶
	Target      string    `json:"target" gorm:"size:20;default:announcement"` // 投递目标：announcement / chat / both
	IsBroadcast bool      `json:"is_broadcast" gorm:"default:false;index"`    // 是否为系统广播
	Status      int       `json:"status" gorm:"default:1;index"`              // 1:正常 0:删除
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定User模型的数据库表名
func (User) TableName() string {
	return "users"
}

// TableName 指定Post模型的数据库表名
func (Post) TableName() string {
	return "posts"
}

// TableName 指定Category模型的数据库表名
func (Category) TableName() string {
	return "categories"
}

// TableName 指定Tag模型的数据库表名
func (Tag) TableName() string {
	return "tags"
}

// TableName 指定Comment模型的数据库表名
func (Comment) TableName() string {
	return "comments"
}

// TableName 指定Setting模型的数据库表名
func (Setting) TableName() string {
	return "settings"
}

// TableName 指定PostView模型的数据库表名
func (PostView) TableName() string {
	return "post_views"
}

// TableName 指定Moment模型的数据库表名
func (Moment) TableName() string {
	return "moments"
}

// TableName 指定MomentLike模型的数据库表名
func (MomentLike) TableName() string {
	return "moment_likes"
}

// TableName 指定IPBlacklist模型的数据库表名
func (IPBlacklist) TableName() string {
	return "ip_blacklist"
}

// TableName 指定IPWhitelist模型的数据库表名
func (IPWhitelist) TableName() string {
	return "ip_whitelist"
}

// TableName 指定PasswordResetToken模型的数据库表名
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// TableName 指定EmailChangeRecord模型的数据库表名
func (EmailChangeRecord) TableName() string {
	return "email_change_records"
}

// TableName 指定ChatMessage模型的数据库表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// FriendLinkCategory 友链分类模型
// 功能说明：存储友链分类信息，用于对友链进行分类管理
type FriendLinkCategory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:50"`
	Description string    `json:"description" gorm:"size:200"`
	SortOrder   int       `json:"sort_order" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	FriendLinks []FriendLink `json:"friend_links,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName 指定FriendLinkCategory模型的数据库表名
func (FriendLinkCategory) TableName() string {
	return "friend_link_categories"
}

// FriendLink 友链模型
// 功能说明：存储友情链接信息，包含网站名称、URL、图标、描述等信息
type FriendLink struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:100"`
	URL         string    `json:"url" gorm:"not null;size:255"`
	Icon        string    `json:"icon" gorm:"size:255"`
	Description string    `json:"description" gorm:"type:text"`
	Screenshot  string    `json:"screenshot" gorm:"size:255"`
	AtomURL     string    `json:"atom_url" gorm:"size:255"`
	CategoryID  uint      `json:"category_id" gorm:"not null;index"` // 分类ID（必选）
	SortOrder   int       `json:"sort_order" gorm:"default:0;index"`
	Status      int       `json:"status" gorm:"default:1;index"` // 1:启用 0:禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Category FriendLinkCategory `json:"category" gorm:"foreignKey:CategoryID"`
}

// TableName 指定FriendLink模型的数据库表名
func (FriendLink) TableName() string {
	return "friend_links"
}

// Album 相册模型
// 功能说明：存储相册照片信息，用于展示个人相册
type Album struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ImageURL    string    `json:"image_url" gorm:"not null;size:500"`
	Title       string    `json:"title" gorm:"size:200"`
	Description string    `json:"description" gorm:"size:500"`
	SortOrder   int       `json:"sort_order" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定Album模型的数据库表名
func (Album) TableName() string {
	return "albums"
}
