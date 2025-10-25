package service

import (
	"errors"

	"blog-backend/model"
	"blog-backend/repository"
	"blog-backend/util"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	// 验证用户名格式
	if !util.ValidateUsername(req.Username) {
		return nil, errors.New("用户名格式不正确（3-20个字符，只能包含字母、数字、下划线）")
	}

	// 验证邮箱格式
	if !util.ValidateEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 验证密码强度
	if !util.ValidatePassword(req.Password) {
		return nil, errors.New("密码长度至少为6个字符")
	}

	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 加密密码
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Username,
		Role:     "user",
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, errors.New("登录失败")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if !util.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 Token
	token, err := util.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("Token 生成失败")
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetProfile 获取用户信息
func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("获取用户信息失败")
	}
	return user, nil
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// UpdateProfile 更新用户信息
func (s *AuthService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新用户信息失败")
	}

	return user, nil
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdatePassword 修改密码
func (s *AuthService) UpdatePassword(userID uint, req *UpdatePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !util.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 验证新密码
	if !util.ValidatePassword(req.NewPassword) {
		return errors.New("新密码长度至少为6个字符")
	}

	// 加密新密码
	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("密码修改失败")
	}

	return nil
}
