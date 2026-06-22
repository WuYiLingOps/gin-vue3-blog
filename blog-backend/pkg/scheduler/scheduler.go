/*
 * 项目名称：blog-backend
 * 文件名称：scheduler.go
 * 创建时间：2026-06-22
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：定时任务调度器，基于 robfig/cron 实现定时任务管理
 */
package scheduler

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron *cron.Cron
	jobs map[string]cron.EntryID
	mu   sync.RWMutex
}

// NewScheduler 创建调度器实例
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
		jobs: make(map[string]cron.EntryID),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("定时任务调度器已启动")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("定时任务调度器已停止")
}

// AddJob 添加定时任务
// name: 任务名称（唯一标识）
// spec: cron 表达式
// cmd: 任务执行函数
func (s *Scheduler) AddJob(name, spec string, cmd func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果已存在同名任务，先移除
	if id, exists := s.jobs[name]; exists {
		s.cron.Remove(id)
	}

	entryID, err := s.cron.AddFunc(spec, func() {
		log.Printf("[调度器] 开始执行任务: %s", name)
		cmd()
		log.Printf("[调度器] 任务执行完成: %s", name)
	})
	if err != nil {
		return err
	}

	s.jobs[name] = entryID
	log.Printf("[调度器] 已注册任务: %s (cron: %s)", name, spec)
	return nil
}

// RemoveJob 移除定时任务
func (s *Scheduler) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id, exists := s.jobs[name]; exists {
		s.cron.Remove(id)
		delete(s.jobs, name)
		log.Printf("[调度器] 已移除任务: %s", name)
	}
}
