package service

import (
	"fmt"
	"time"

	"blog-backend/repository"
)

type CleanupService struct {
	resetTokenRepo *repository.PasswordResetRepository
}

func NewCleanupService() *CleanupService {
	return &CleanupService{
		resetTokenRepo: repository.NewPasswordResetRepository(),
	}
}

// StartCleanupTasks 启动定期清理任务
func (s *CleanupService) StartCleanupTasks() {
	// 每小时清理一次过期的密码重置令牌
	go s.cleanupExpiredTokensPeriodically(1 * time.Hour)
}

// cleanupExpiredTokensPeriodically 定期清理过期令牌
func (s *CleanupService) cleanupExpiredTokensPeriodically(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 立即执行一次清理
	s.cleanupExpiredTokens()

	// 定期执行
	for range ticker.C {
		s.cleanupExpiredTokens()
	}
}

// cleanupExpiredTokens 清理过期令牌
func (s *CleanupService) cleanupExpiredTokens() {
	if err := s.resetTokenRepo.DeleteExpired(); err != nil {
		fmt.Printf("清理过期令牌失败: %v\n", err)
	} else {
		fmt.Printf("过期令牌清理完成: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	}
}

