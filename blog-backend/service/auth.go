/*
 * 项目名称：blog-backend
 * 文件名称：auth.go
 * 创建时间：2026-01-31 16:34:35
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：认证业务逻辑层，提供用户注册、登录、密码重置、邮箱修改等认证相关业务处理
 */
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"blog-backend/config"
	"blog-backend/constant"
	"blog-backend/db"
	"blog-backend/model"
	"blog-backend/repository"
	"blog-backend/util"

	"gorm.io/gorm"
)

// AuthService 认证业务逻辑层结构体
type AuthService struct {
	userRepo        *repository.UserRepository
	resetTokenRepo  *repository.PasswordResetRepository
	emailChangeRepo *repository.EmailChangeRepository
	settingRepo     *repository.SettingRepository
}

// NewAuthService 创建认证业务逻辑层实例
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:        repository.NewUserRepository(),
		resetTokenRepo:  repository.NewPasswordResetRepository(),
		emailChangeRepo: repository.NewEmailChangeRepository(),
		settingRepo:     repository.NewSettingRepository(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code" binding:"required,len=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaID string `json:"captcha_id" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest, ip string) (*model.User, error) {
	// 检查注册功能是否被禁用
	if isRegisterDisabled, err := s.isRegisterDisabled(); err == nil && isRegisterDisabled {
		return nil, errors.New("用户注册功能已关闭")
	}

	// 验证邮箱验证码
	resetToken, err := s.resetTokenRepo.GetValidToken(req.Email, req.Code)
	if err != nil {
		return nil, errors.New("验证码无效或已过期")
	}

	if resetToken.IsUsed {
		return nil, errors.New("验证码已被使用")
	}

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
		Role:     constant.RoleUser,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 标记验证码为已使用
	resetToken.IsUsed = true
	s.resetTokenRepo.Update(resetToken)

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest, ip string) (*LoginResponse, error) {
	// 验证验证码
	if err := util.VerifyCaptcha(req.CaptchaID, req.Captcha, ip); err != nil {
		return nil, err
	}

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

	// 如果更新的是系统拥有者（super_admin），清理博主信息缓存，确保前台个人名片和关于我页面立即生效
	if user.Role == constant.RoleSuperAdmin {
		ctx := context.Background()
		_ = db.RDB.Del(ctx, "blog:author_profile").Err()
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

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ForgotPassword 发送重置密码邮件
func (s *AuthService) ForgotPassword(req *ForgotPasswordRequest, ip string) error {
	// 检查邮箱是否存在
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		// 为了安全，不告诉用户邮箱是否存在，但仍然返回成功
		return nil
	}

	// 检查发送频率限制（1分钟内只能发送一次）
	if recent, err := s.resetTokenRepo.GetRecentByEmail(req.Email, 1*60*1000000000); err == nil && recent != nil {
		remainingTime := 60 - int(time.Since(recent.CreatedAt).Seconds())
		if remainingTime > 0 {
			return errors.New(fmt.Sprintf("验证码发送过于频繁，请%d秒后再试", remainingTime))
		}
	}

	// 生成验证码和令牌
	code := util.GenerateVerificationCode()
	token := util.GenerateRandomString(32)

	// 保存到数据库
	userID := user.ID
	resetToken := &model.PasswordResetToken{
		UserID:   &userID,
		Email:    req.Email,
		Token:    token,
		Code:     code,
		ExpireAt: util.GetTimeAfterMinutes(15), // 15分钟有效期
		IsUsed:   false,
	}

	if err := s.resetTokenRepo.Create(resetToken); err != nil {
		return errors.New("系统错误，请稍后重试")
	}

	// 获取网站名称
	siteName := s.getSiteName()

	// 获取邮箱配置
	emailConfig := util.EmailConfig{
		Host:     config.Cfg.Email.Host,
		Port:     config.Cfg.Email.Port,
		Username: config.Cfg.Email.Username,
		Password: config.Cfg.Email.Password,
		FromName: config.Cfg.Email.FromName,
		SiteName: siteName,
	}

	// 异步发送邮件，避免阻塞请求
	go func(config util.EmailConfig, email string, verificationCode string) {
		if err := util.SendResetPasswordEmail(config, email, verificationCode); err != nil {
			// 记录错误日志，但不影响主流程
			fmt.Printf("发送密码重置邮件失败: %v\n", err)
		}
	}(emailConfig, req.Email, code)

	return nil
}

// ResetPassword 重置密码
func (s *AuthService) ResetPassword(req *ResetPasswordRequest) error {
	// 查找有效的重置令牌
	resetToken, err := s.resetTokenRepo.GetValidToken(req.Email, req.Code)
	if err != nil {
		return errors.New("验证码无效或已过期")
	}

	if resetToken.IsUsed {
		return errors.New("验证码已被使用")
	}

	// 验证密码强度
	if !util.ValidatePassword(req.NewPassword) {
		return errors.New("新密码长度至少为6个字符")
	}

	// 获取用户（密码重置时UserID不为空）
	if resetToken.UserID == nil {
		return errors.New("无效的重置令牌")
	}

	user, err := s.userRepo.GetByID(*resetToken.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 加密新密码
	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("密码重置失败")
	}

	// 标记令牌为已使用
	resetToken.IsUsed = true
	s.resetTokenRepo.Update(resetToken)

	return nil
}

// UpdateEmailRequest 修改邮箱请求
type UpdateEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

// UpdateEmail 修改邮箱
func (s *AuthService) UpdateEmail(userID uint, req *UpdateEmailRequest) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证新邮箱格式
	if !util.ValidateEmail(req.NewEmail) {
		return errors.New("邮箱格式不正确")
	}

	// 检查新邮箱是否与当前邮箱相同
	if user.Email == req.NewEmail {
		return errors.New("新邮箱与当前邮箱相同")
	}

	// 检查新邮箱是否已被使用
	if _, err := s.userRepo.GetByEmail(req.NewEmail); err == nil {
		return errors.New("该邮箱已被其他用户使用")
	}

	// 检查一年内的修改次数
	count, err := s.emailChangeRepo.CountByUserIDInYear(userID)
	if err != nil {
		return errors.New("系统错误，请稍后重试")
	}
	if count >= 2 {
		return errors.New("一年内只能修改两次邮箱，您已达到上限")
	}

	// 记录旧邮箱
	oldEmail := user.Email

	// 更新邮箱
	user.Email = req.NewEmail
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("邮箱修改失败")
	}

	// 记录修改历史
	record := &model.EmailChangeRecord{
		UserID:   userID,
		OldEmail: oldEmail,
		NewEmail: req.NewEmail,
	}
	if err := s.emailChangeRepo.Create(record); err != nil {
		// 记录失败不影响主流程，只记录日志
		println("邮箱修改记录创建失败:", err.Error())
	}

	return nil
}

// GetEmailChangeCount 获取用户一年内的邮箱修改次数
func (s *AuthService) GetEmailChangeCount(userID uint) (int64, error) {
	return s.emailChangeRepo.CountByUserIDInYear(userID)
}

// SendRegisterCodeRequest 发送注册验证码请求
type SendRegisterCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// SendRegisterCode 发送注册验证码
func (s *AuthService) SendRegisterCode(req *SendRegisterCodeRequest, ip string) error {
	// 检查注册功能是否被禁用
	if isRegisterDisabled, err := s.isRegisterDisabled(); err == nil && isRegisterDisabled {
		return errors.New("用户注册功能已关闭")
	}

	// 验证邮箱格式
	if !util.ValidateEmail(req.Email) {
		return errors.New("邮箱格式不正确")
	}

	// 检查邮箱是否已被注册
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return errors.New("该邮箱已被注册")
	}

	// 检查发送频率限制（1分钟内只能发送一次）
	if recent, err := s.resetTokenRepo.GetRecentByEmail(req.Email, 1*60*1000000000); err == nil && recent != nil {
		remainingTime := 60 - int(time.Since(recent.CreatedAt).Seconds())
		if remainingTime > 0 {
			return errors.New(fmt.Sprintf("验证码发送过于频繁，请%d秒后再试", remainingTime))
		}
	}

	// 生成验证码和令牌
	code := util.GenerateVerificationCode()
	token := util.GenerateRandomString(32)

	// 保存到数据库（注册验证码不需要关联用户ID，使用nil）
	resetToken := &model.PasswordResetToken{
		UserID:   nil, // 注册时还没有用户ID
		Email:    req.Email,
		Token:    token,
		Code:     code,
		ExpireAt: util.GetTimeAfterMinutes(15), // 15分钟有效期
		IsUsed:   false,
	}

	if err := s.resetTokenRepo.Create(resetToken); err != nil {
		return errors.New("系统错误，请稍后重试")
	}

	// 获取网站名称
	siteName := s.getSiteName()

	// 获取邮箱配置
	emailConfig := util.EmailConfig{
		Host:     config.Cfg.Email.Host,
		Port:     config.Cfg.Email.Port,
		Username: config.Cfg.Email.Username,
		Password: config.Cfg.Email.Password,
		FromName: config.Cfg.Email.FromName,
		SiteName: siteName,
	}

	// 异步发送邮件，避免阻塞请求
	go func(config util.EmailConfig, email string, verificationCode string) {
		if err := util.SendRegisterVerificationEmail(config, email, verificationCode); err != nil {
			// 记录错误日志，但不影响主流程
			fmt.Printf("发送注册验证码邮件失败: %v\n", err)
		}
	}(emailConfig, req.Email, code)

	return nil
}

// getSiteName 获取网站名称
func (s *AuthService) getSiteName() string {
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

// isRegisterDisabled 检查注册功能是否被禁用
func (s *AuthService) isRegisterDisabled() (bool, error) {
	setting, err := s.settingRepo.GetByKey("disable_register")
	if err != nil {
		// 如果配置不存在，默认允许注册（返回false）
		return false, nil
	}
	// "1" 表示禁用，"0" 或空表示允许
	return setting.Value == "1", nil
}
