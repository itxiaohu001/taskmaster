package task

import (
	"context"
	"fmt"
	"time"
)

// SimpleTask 简单任务示例
type SimpleTask struct {
	*BaseTask
	message string
}

// NewSimpleTask 创建简单任务
func NewSimpleTask(id, name, message string, priority, maxRetries int) *SimpleTask {
	return &SimpleTask{
		BaseTask: NewBaseTask(id, name, priority, maxRetries),
		message:  message,
	}
}

// Execute 执行简单任务
func (t *SimpleTask) Execute(ctx context.Context) error {
	// 模拟任务执行
	time.Sleep(1 * time.Second)
	fmt.Printf("执行任务: %s, 消息: %s\n", t.Name(), t.message)
	return nil
}

// HTTPTask HTTP请求任务
type HTTPTask struct {
	*BaseTask
	url     string
	method  string
	timeout time.Duration
}

// NewHTTPTask 创建HTTP任务
func NewHTTPTask(id, name, url, method string, timeout time.Duration, priority, maxRetries int) *HTTPTask {
	return &HTTPTask{
		BaseTask: NewBaseTask(id, name, priority, maxRetries),
		url:      url,
		method:   method,
		timeout:  timeout,
	}
}

// Execute 执行HTTP任务
func (t *HTTPTask) Execute(ctx context.Context) error {
	// 这里可以实现实际的HTTP请求
	// 为了示例，我们只是模拟
	fmt.Printf("执行HTTP任务: %s %s\n", t.method, t.url)
	
	// 模拟网络延迟
	time.Sleep(2 * time.Second)
	
	// 模拟成功
	return nil
}

// FileProcessTask 文件处理任务
type FileProcessTask struct {
	*BaseTask
	filePath string
	action   string // read, write, delete
}

// NewFileProcessTask 创建文件处理任务
func NewFileProcessTask(id, name, filePath, action string, priority, maxRetries int) *FileProcessTask {
	return &FileProcessTask{
		BaseTask: NewBaseTask(id, name, priority, maxRetries),
		filePath: filePath,
		action:   action,
	}
}

// Execute 执行文件处理任务
func (t *FileProcessTask) Execute(ctx context.Context) error {
	fmt.Printf("执行文件任务: %s %s\n", t.action, t.filePath)
	
	// 模拟文件处理
	time.Sleep(500 * time.Millisecond)
	
	// 模拟成功
	return nil
}

// BatchTask 批量任务
type BatchTask struct {
	*BaseTask
	tasks []Task
}

// NewBatchTask 创建批量任务
func NewBatchTask(id, name string, tasks []Task, priority, maxRetries int) *BatchTask {
	return &BatchTask{
		BaseTask: NewBaseTask(id, name, priority, maxRetries),
		tasks:    tasks,
	}
}

// Execute 执行批量任务
func (t *BatchTask) Execute(ctx context.Context) error {
	fmt.Printf("开始执行批量任务: %s, 包含 %d 个子任务\n", t.Name(), len(t.tasks))
	
	for i, task := range t.tasks {
		fmt.Printf("执行子任务 %d: %s\n", i+1, task.Name())
		if err := task.Execute(ctx); err != nil {
			return fmt.Errorf("子任务 %s 执行失败: %w", task.Name(), err)
		}
	}
	
	fmt.Printf("批量任务 %s 执行完成\n", t.Name())
	return nil
}

// RetryableTask 可重试任务包装器
type RetryableTask struct {
	*BaseTask
	task Task
}

// NewRetryableTask 创建可重试任务
func NewRetryableTask(task Task, maxRetries int) *RetryableTask {
	return &RetryableTask{
		BaseTask: NewBaseTask(task.ID(), task.Name(), task.GetPriority(), maxRetries),
		task:     task,
	}
}

// Execute 执行可重试任务
func (t *RetryableTask) Execute(ctx context.Context) error {
	var lastErr error
	
	for i := 0; i <= t.GetMaxRetries(); i++ {
		if err := t.task.Execute(ctx); err != nil {
			lastErr = err
			if i < t.GetMaxRetries() {
				fmt.Printf("任务 %s 执行失败，将在 5s 后重试 (第 %d 次重试)\n", 
					t.Name(), i+1)
				time.Sleep(5 * time.Second)
				continue
			}
		} else {
			return nil
		}
	}
	
	return fmt.Errorf("任务 %s 在 %d 次重试后仍然失败: %w", t.Name(), t.GetMaxRetries(), lastErr)
} 