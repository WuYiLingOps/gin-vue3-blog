package model

import (
	"time"
)

// User 用户模型
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
type Post struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null;size:200;index"`
	Content     string     `json:"content" gorm:"type:text"`
	Summary     string     `json:"summary" gorm:"size:500"`
	Cover       string     `json:"cover" gorm:"size:255"`
	Status      int        `json:"status" gorm:"default:1;index"` // 1:发布 0:草稿 -1:删除
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
}

// Category 分类模型
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
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Color     string    `json:"color" gorm:"size:20"`
	PostCount int       `json:"post_count" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Comment 评论模型
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	PostID    uint      `json:"post_id" gorm:"index"`
	UserID    uint      `json:"user_id" gorm:"index"`
	ParentID  *uint     `json:"parent_id" gorm:"index"` // 父评论ID，用于回复
	Status    int       `json:"status" gorm:"default:1"` // 1:正常 0:隐藏
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联关系
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Post     Post      `json:"post" gorm:"foreignKey:PostID"`
	Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Comment `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// Setting 系统配置模型
type Setting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"uniqueIndex;not null;size:100"`
	Value     string    `json:"value" gorm:"type:text"`
	Type      string    `json:"type" gorm:"size:20"` // text, json, image
	Group     string    `json:"group" gorm:"size:50;index"` // about, site, etc.
	Label     string    `json:"label" gorm:"size:100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VisitStat 访问量统计模型
type VisitStat struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Date      time.Time `json:"date" gorm:"uniqueIndex;not null"` // 日期（只保留年月日）
	ViewCount int       `json:"view_count" gorm:"default:0"`      // 当天浏览量
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostView 文章阅读记录模型
type PostView struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"index;not null"`
	UserID    *uint     `json:"user_id" gorm:"index"` // 已登录用户ID，可为空
	IP        string    `json:"ip" gorm:"size:45;index"` // 访客IP地址
	CreatedAt time.Time `json:"created_at"`
}

// Moment 说说模型
type Moment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	Images    string    `json:"images" gorm:"type:text"` // JSON数组格式存储图片URLs
	UserID    uint      `json:"user_id" gorm:"index"`
	Status    int       `json:"status" gorm:"default:1;index"` // 1:公开 0:私密 -1:删除
	LikeCount int       `json:"like_count" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联关系
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (Post) TableName() string {
	return "posts"
}

func (Category) TableName() string {
	return "categories"
}

func (Tag) TableName() string {
	return "tags"
}

func (Comment) TableName() string {
	return "comments"
}

func (Setting) TableName() string {
	return "settings"
}

func (VisitStat) TableName() string {
	return "visit_stats"
}

func (PostView) TableName() string {
	return "post_views"
}

func (Moment) TableName() string {
	return "moments"
}
