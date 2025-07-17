package manager

import (
	"fmt"
	"sync"

	"github.com/itxiaohu001/taskmaster/internal/engine"
	"github.com/itxiaohu001/taskmaster/pkg/config"
	"github.com/itxiaohu001/taskmaster/pkg/task"
)

// Manager 任务管理器
type Manager struct {
	config    *config.Config
	scheduler *engine.Scheduler
	mu        sync.RWMutex
	running   bool
}

// NewManager 创建新的任务管理器
func NewManager(scheduler *engine.Scheduler) *Manager {
	return &Manager{
		scheduler: scheduler,
	}
}

// Start 启动任务管理器
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("manager is already running")
	}

	// 启动调度器
	m.scheduler.Start()

	m.running = true

	return nil
}

// Stop 停止任务管理器
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	// 停止调度器
	if m.scheduler != nil {
		m.scheduler.Stop()
	}

	m.running = false

	return nil
}

// AddTask 添加任务
func (m *Manager) AddTask(t task.Task) error {
	if !m.running {
		return fmt.Errorf("manager is not running")
	}

	// 添加到调度器
	if err := m.scheduler.AddTask(t); err != nil {
		return fmt.Errorf("failed to add task to scheduler: %w", err)
	}

	return nil
}

// RemoveTask 移除任务
func (m *Manager) RemoveTask(taskID string) error {
	if !m.running {
		return fmt.Errorf("manager is not running")
	}

	// 从调度器中移除
	if err := m.scheduler.RemoveTask(taskID); err != nil {
		return fmt.Errorf("failed to remove task from scheduler: %w", err)
	}

	return nil
}

// GetTask 获取任务
func (m *Manager) GetTask(taskID string) (task.Task, bool) {
	return m.scheduler.GetTask(taskID)
}

// GetTasks 获取所有任务
func (m *Manager) GetTasks() []task.Task {
	return m.scheduler.GetTasks()
}

// Wait 等待所有任务完成
func (m *Manager) Wait() {
	if m.scheduler != nil {
		m.scheduler.Wait()
	}
}
