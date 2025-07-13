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
	storage   Storage
	monitor   Monitor
	logger    Logger
	mu        sync.RWMutex
	running   bool
}

// Storage 存储接口
type Storage interface {
	Save(t task.Task) error
	Load(taskID string) (task.Task, error)
	LoadAll() ([]task.Task, error)
	Delete(taskID string) error
}

// Monitor 监控接口
type Monitor interface {
	Start() error
	Stop() error
	RecordTaskStarted(t task.Task)
	RecordTaskCompleted(t task.Task)
	RecordTaskFailed(t task.Task)
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// NewManager 创建新的任务管理器
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config: cfg,
	}
}

// Start 启动任务管理器
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("manager is already running")
	}

	// 初始化调度器
	m.scheduler = engine.NewScheduler(m.config.Scheduler.Workers)

	// 初始化存储
	if err := m.initStorage(); err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// 初始化监控
	if err := m.initMonitor(); err != nil {
		return fmt.Errorf("failed to initialize monitor: %w", err)
	}

	// 初始化日志
	if err := m.initLogger(); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 添加事件处理器
	m.scheduler.AddEventHandler(m)

	// 启动调度器
	m.scheduler.Start()

	// 启动监控
	if m.monitor != nil {
		if err := m.monitor.Start(); err != nil {
			return fmt.Errorf("failed to start monitor: %w", err)
		}
	}

	m.running = true
	m.logger.Info("Task manager started successfully")
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

	// 停止监控
	if m.monitor != nil {
		m.monitor.Stop()
	}

	m.running = false
	m.logger.Info("Task manager stopped")
	return nil
}

// AddTask 添加任务
func (m *Manager) AddTask(t task.Task) error {
	if !m.running {
		return fmt.Errorf("manager is not running")
	}

	// 保存任务到存储
	if err := m.storage.Save(t); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// 添加到调度器
	if err := m.scheduler.AddTask(t); err != nil {
		return fmt.Errorf("failed to add task to scheduler: %w", err)
	}

	m.logger.Info("Task added successfully", "task_id", t.ID(), "task_name", t.Name())
	return nil
}

// RemoveTask 移除任务
func (m *Manager) RemoveTask(taskID string) error {
	if !m.running {
		return fmt.Errorf("manager is not running")
	}

	// 从存储中删除
	if err := m.storage.Delete(taskID); err != nil {
		return fmt.Errorf("failed to delete task from storage: %w", err)
	}

	// 从调度器中移除
	if err := m.scheduler.RemoveTask(taskID); err != nil {
		return fmt.Errorf("failed to remove task from scheduler: %w", err)
	}

	m.logger.Info("Task removed successfully", "task_id", taskID)
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

// initStorage 初始化存储
func (m *Manager) initStorage() error {
	switch m.config.Storage.Type {
	case "memory":
		m.storage = engine.NewMemoryStorage()
	case "file":
		m.storage = engine.NewFileStorage(m.config.Storage.FilePath)
	default:
		return fmt.Errorf("unsupported storage type: %s", m.config.Storage.Type)
	}
	return nil
}

// initMonitor 初始化监控
func (m *Manager) initMonitor() error {
	if !m.config.Monitor.Enabled {
		return nil
	}

	m.monitor = engine.NewMetricsMonitor(m.config.Monitor.MetricsPort)
	return nil
}

// initLogger 初始化日志
func (m *Manager) initLogger() error {
	m.logger = engine.NewLogger(m.config.Log)
	return nil
}

// EventHandler 实现

// OnTaskStarted 任务开始事件
func (m *Manager) OnTaskStarted(t task.Task) {
	m.logger.Info("Task started", "task_id", t.ID(), "task_name", t.Name())
	if m.monitor != nil {
		m.monitor.RecordTaskStarted(t)
	}
}

// OnTaskCompleted 任务完成事件
func (m *Manager) OnTaskCompleted(t task.Task) {
	m.logger.Info("Task completed", "task_id", t.ID(), "task_name", t.Name())
	if m.monitor != nil {
		m.monitor.RecordTaskCompleted(t)
	}
}

// OnTaskFailed 任务失败事件
func (m *Manager) OnTaskFailed(t task.Task) {
	m.logger.Error("Task failed", "task_id", t.ID(), "task_name", t.Name(), "error", t.GetError())
	if m.monitor != nil {
		m.monitor.RecordTaskFailed(t)
	}
}

// OnTaskCancelled 任务取消事件
func (m *Manager) OnTaskCancelled(t task.Task) {
	m.logger.Warn("Task cancelled", "task_id", t.ID(), "task_name", t.Name())
}
