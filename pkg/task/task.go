package task

import (
	"context"
	"time"
)

// Status 表示任务状态
type Status int

const (
	StatusPending Status = iota
	StatusRunning
	StatusCompleted
	StatusFailed
	StatusCancelled
	StatusRetrying
)

// String 返回状态字符串表示
func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusCompleted:
		return "completed"
	case StatusFailed:
		return "failed"
	case StatusCancelled:
		return "cancelled"
	case StatusRetrying:
		return "retrying"
	default:
		return "unknown"
	}
}

// Task 定义任务接口
type Task interface {
	// ID 返回任务唯一标识
	ID() string
	
	// Name 返回任务名称
	Name() string
	
	// Execute 执行任务
	Execute(ctx context.Context) error
	
	// GetStatus 获取任务状态
	GetStatus() Status
	
	// SetStatus 设置任务状态
	SetStatus(status Status)
	
	// GetPriority 获取任务优先级
	GetPriority() int
	
	// GetRetryCount 获取重试次数
	GetRetryCount() int
	
	// GetMaxRetries 获取最大重试次数
	GetMaxRetries() int
	
	// GetCreatedAt 获取创建时间
	GetCreatedAt() time.Time
	
	// GetStartedAt 获取开始时间
	GetStartedAt() time.Time
	
	// SetStartedAt 设置开始时间
	SetStartedAt(startedAt time.Time)
	
	// GetCompletedAt 获取完成时间
	GetCompletedAt() time.Time
	
	// SetCompletedAt 设置完成时间
	SetCompletedAt(completedAt time.Time)
	
	// GetError 获取错误信息
	GetError() error
	
	// SetError 设置错误信息
	SetError(err error)
	
	// GetResult 获取任务结果
	GetResult() interface{}
	
	// SetResult 设置任务结果
	SetResult(result interface{})
}

// BaseTask 提供任务的基础实现
type BaseTask struct {
	id          string
	name        string
	status      Status
	priority    int
	retryCount  int
	maxRetries  int
	createdAt   time.Time
	startedAt   time.Time
	completedAt time.Time
	err         error
	result      interface{}
}

// NewBaseTask 创建基础任务
func NewBaseTask(id, name string, priority, maxRetries int) *BaseTask {
	return &BaseTask{
		id:         id,
		name:       name,
		status:     StatusPending,
		priority:   priority,
		maxRetries: maxRetries,
		createdAt:  time.Now(),
	}
}

// ID 返回任务唯一标识
func (t *BaseTask) ID() string {
	return t.id
}

// Name 返回任务名称
func (t *BaseTask) Name() string {
	return t.name
}

// GetStatus 获取任务状态
func (t *BaseTask) GetStatus() Status {
	return t.status
}

// SetStatus 设置任务状态
func (t *BaseTask) SetStatus(status Status) {
	t.status = status
}

// GetPriority 获取任务优先级
func (t *BaseTask) GetPriority() int {
	return t.priority
}

// GetRetryCount 获取重试次数
func (t *BaseTask) GetRetryCount() int {
	return t.retryCount
}

// GetMaxRetries 获取最大重试次数
func (t *BaseTask) GetMaxRetries() int {
	return t.maxRetries
}

// GetCreatedAt 获取创建时间
func (t *BaseTask) GetCreatedAt() time.Time {
	return t.createdAt
}

// GetStartedAt 获取开始时间
func (t *BaseTask) GetStartedAt() time.Time {
	return t.startedAt
}

// GetCompletedAt 获取完成时间
func (t *BaseTask) GetCompletedAt() time.Time {
	return t.completedAt
}

// GetError 获取错误信息
func (t *BaseTask) GetError() error {
	return t.err
}

// SetError 设置错误信息
func (t *BaseTask) SetError(err error) {
	t.err = err
}

// GetResult 获取任务结果
func (t *BaseTask) GetResult() interface{} {
	return t.result
}

// SetResult 设置任务结果
func (t *BaseTask) SetResult(result interface{}) {
	t.result = result
}

// SetStartedAt 设置开始时间
func (t *BaseTask) SetStartedAt(startedAt time.Time) {
	t.startedAt = startedAt
}

// SetCompletedAt 设置完成时间
func (t *BaseTask) SetCompletedAt(completedAt time.Time) {
	t.completedAt = completedAt
}

// SetRetryCount 设置重试次数
func (t *BaseTask) SetRetryCount(retryCount int) {
	t.retryCount = retryCount
}

// Execute 默认执行方法，需要子类重写
func (t *BaseTask) Execute(ctx context.Context) error {
	return nil
} 