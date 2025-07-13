package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/itxiaohu001/taskmaster/pkg/task"
)

// Scheduler 任务调度器
type Scheduler struct {
	mu            sync.RWMutex
	tasks         map[string]task.Task
	queue         *PriorityQueue
	workers       int
	running       bool
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	eventHandlers []EventHandler
}

// EventHandler 事件处理器接口
type EventHandler interface {
	OnTaskStarted(t task.Task)
	OnTaskCompleted(t task.Task)
	OnTaskFailed(t task.Task)
	OnTaskCancelled(t task.Task)
}

// NewScheduler 创建新的调度器
func NewScheduler(workers int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:   make(map[string]task.Task),
		queue:   NewPriorityQueue(),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// AddTask 添加任务到调度器
func (s *Scheduler) AddTask(t task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[t.ID()]; exists {
		return fmt.Errorf("task with ID %s already exists", t.ID())
	}

	s.tasks[t.ID()] = t
	s.queue.PushTask(t)
	return nil
}

// RemoveTask 从调度器移除任务
func (s *Scheduler) RemoveTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	delete(s.tasks, taskID)
	// 从队列中移除任务
	s.queue.Remove(taskID)
	return nil
}

// GetTask 获取任务
func (s *Scheduler) GetTask(taskID string) (task.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, exists := s.tasks[taskID]
	return t, exists
}

// GetTasks 获取所有任务
func (s *Scheduler) GetTasks() []task.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]task.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	s.cancel()
	s.wg.Wait()
}

// Wait 等待所有任务完成
func (s *Scheduler) Wait() {
	s.wg.Wait()
}

// AddEventHandler 添加事件处理器
func (s *Scheduler) AddEventHandler(handler EventHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventHandlers = append(s.eventHandlers, handler)
}

// worker 工作协程
func (s *Scheduler) worker(id int) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			task := s.queue.PopTask()
			if task == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			s.executeTask(task)
		}
	}
}

// executeTask 执行任务
func (s *Scheduler) executeTask(t task.Task) {
	// 设置任务状态为运行中
	t.SetStatus(task.StatusRunning)
	t.SetStartedAt(time.Now())

	// 触发任务开始事件
	s.triggerTaskStarted(t)

	// 执行任务
	err := t.Execute(s.ctx)

	// 处理执行结果
	if err != nil {
		t.SetError(err)
		t.SetStatus(task.StatusFailed)
		s.triggerTaskFailed(t)
	} else {
		t.SetStatus(task.StatusCompleted)
		t.SetCompletedAt(time.Now())
		s.triggerTaskCompleted(t)
	}
}

// triggerTaskStarted 触发任务开始事件
func (s *Scheduler) triggerTaskStarted(t task.Task) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, handler := range s.eventHandlers {
		handler.OnTaskStarted(t)
	}
}

// triggerTaskCompleted 触发任务完成事件
func (s *Scheduler) triggerTaskCompleted(t task.Task) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, handler := range s.eventHandlers {
		handler.OnTaskCompleted(t)
	}
}

// triggerTaskFailed 触发任务失败事件
func (s *Scheduler) triggerTaskFailed(t task.Task) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, handler := range s.eventHandlers {
		handler.OnTaskFailed(t)
	}
}

// triggerTaskCancelled 触发任务取消事件
func (s *Scheduler) triggerTaskCancelled(t task.Task) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, handler := range s.eventHandlers {
		handler.OnTaskCancelled(t)
	}
}
