package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/itxiaohu001/taskmaster/pkg/task"
)

// TaskData 用于序列化的任务数据结构
type TaskData struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Status      task.Status            `json:"status"`
	Priority    int                    `json:"priority"`
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt time.Time              `json:"completed_at"`
	Error       string                 `json:"error,omitempty"`
	Result      interface{}            `json:"result,omitempty"`
	Type        string                 `json:"type"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// MemoryStorage 内存存储实现
type MemoryStorage struct {
	mu    sync.RWMutex
	tasks map[string]task.Task
}

// NewMemoryStorage 创建内存存储
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[string]task.Task),
	}
}

// Save 保存任务
func (m *MemoryStorage) Save(t task.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks[t.ID()] = t
	return nil
}

// Load 加载任务
func (m *MemoryStorage) Load(taskID string) (task.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, exists := m.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	return t, nil
}

// LoadAll 加载所有任务
func (m *MemoryStorage) LoadAll() ([]task.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]task.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// Delete 删除任务
func (m *MemoryStorage) Delete(taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tasks[taskID]; !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	delete(m.tasks, taskID)
	return nil
}

// FileStorage 文件存储实现
type FileStorage struct {
	filepath string
	mu       sync.RWMutex
}

// NewFileStorage 创建文件存储
func NewFileStorage(filepath string) *FileStorage {
	return &FileStorage{
		filepath: filepath,
	}
}

// Save 保存任务到文件
func (f *FileStorage) Save(t task.Task) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 读取现有任务
	tasks, err := f.loadTasks()
	if err != nil {
		return err
	}

	// 添加或更新任务
	tasks[t.ID()] = t

	// 保存到文件
	return f.saveTasks(tasks)
}

// Load 从文件加载任务
func (f *FileStorage) Load(taskID string) (task.Task, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	tasks, err := f.loadTasks()
	if err != nil {
		return nil, err
	}

	t, exists := tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return t, nil
}

// LoadAll 从文件加载所有任务
func (f *FileStorage) LoadAll() ([]task.Task, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	tasks, err := f.loadTasks()
	if err != nil {
		return nil, err
	}

	result := make([]task.Task, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, t)
	}

	return result, nil
}

// Delete 从文件删除任务
func (f *FileStorage) Delete(taskID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	tasks, err := f.loadTasks()
	if err != nil {
		return err
	}

	if _, exists := tasks[taskID]; !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	delete(tasks, taskID)
	return f.saveTasks(tasks)
}

// loadTasks 从文件加载任务
func (f *FileStorage) loadTasks() (map[string]task.Task, error) {
	// 检查文件是否存在
	if _, err := os.Stat(f.filepath); os.IsNotExist(err) {
		return make(map[string]task.Task), nil
	}

	data, err := os.ReadFile(f.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return make(map[string]task.Task), nil
	}

	var taskDataMap map[string]TaskData
	if err := json.Unmarshal(data, &taskDataMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	// 转换为 Task 接口
	tasks := make(map[string]task.Task)
	for id, taskData := range taskDataMap {
		t, err := f.taskDataToTask(taskData)
		if err != nil {
			return nil, fmt.Errorf("failed to convert task data: %w", err)
		}
		tasks[id] = t
	}

	return tasks, nil
}

// saveTasks 保存任务到文件
func (f *FileStorage) saveTasks(tasks map[string]task.Task) error {
	// 确保目录存在
	dir := f.filepath[:len(f.filepath)-len(f.filepath[len(f.filepath)-1:])]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 转换为可序列化的数据结构
	taskDataMap := make(map[string]TaskData)
	for id, t := range tasks {
		taskData, err := f.taskToTaskData(t)
		if err != nil {
			return fmt.Errorf("failed to convert task: %w", err)
		}
		taskDataMap[id] = taskData
	}

	data, err := json.MarshalIndent(taskDataMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := os.WriteFile(f.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// taskToTaskData 将 Task 转换为 TaskData
func (f *FileStorage) taskToTaskData(t task.Task) (TaskData, error) {
	taskData := TaskData{
		ID:          t.ID(),
		Name:        t.Name(),
		Status:      t.GetStatus(),
		Priority:    t.GetPriority(),
		RetryCount:  t.GetRetryCount(),
		MaxRetries:  t.GetMaxRetries(),
		CreatedAt:   t.GetCreatedAt(),
		StartedAt:   t.GetStartedAt(),
		CompletedAt: t.GetCompletedAt(),
		Result:      t.GetResult(),
	}

	if t.GetError() != nil {
		taskData.Error = t.GetError().Error()
	}

	// 根据任务类型设置 Type 和 Data
	switch task := t.(type) {
	case *task.SimpleTask:
		taskData.Type = "simple"
		taskData.Data = map[string]interface{}{
			"message": task.GetMessage(),
		}
	case *task.HTTPTask:
		taskData.Type = "http"
		taskData.Data = map[string]interface{}{
			"url":     task.GetURL(),
			"method":  task.GetMethod(),
			"timeout": task.GetTimeout().String(),
		}
	case *task.FileProcessTask:
		taskData.Type = "file"
		taskData.Data = map[string]interface{}{
			"file_path": task.GetFilePath(),
			"action":    task.GetAction(),
		}
	case *task.BatchTask:
		taskData.Type = "batch"
		// 批量任务的子任务暂时不序列化，因为可能导致循环引用
		taskData.Data = map[string]interface{}{
			"sub_task_count": len(task.GetTasks()),
		}
	case *task.RetryableTask:
		taskData.Type = "retryable"
		// 重试任务包装的任务暂时不序列化
		taskData.Data = map[string]interface{}{
			"wrapped_task_id": task.GetWrappedTask().ID(),
		}
	default:
		taskData.Type = "unknown"
	}

	return taskData, nil
}

// taskDataToTask 将 TaskData 转换为 Task
func (f *FileStorage) taskDataToTask(taskData TaskData) (task.Task, error) {
	baseTask := task.NewBaseTask(taskData.ID, taskData.Name, taskData.Priority, taskData.MaxRetries)

	// 设置基础属性
	baseTask.SetStatus(taskData.Status)
	baseTask.SetStartedAt(taskData.StartedAt)
	baseTask.SetCompletedAt(taskData.CompletedAt)
	if taskData.Error != "" {
		baseTask.SetError(fmt.Errorf(taskData.Error))
	}
	baseTask.SetResult(taskData.Result)

	// 根据类型创建具体的任务
	switch taskData.Type {
	case "simple":
		message, _ := taskData.Data["message"].(string)
		simpleTask := task.NewSimpleTask(taskData.ID, taskData.Name, message, taskData.Priority, taskData.MaxRetries)
		// 设置状态和结果
		simpleTask.SetStatus(taskData.Status)
		simpleTask.SetStartedAt(taskData.StartedAt)
		simpleTask.SetCompletedAt(taskData.CompletedAt)
		if taskData.Error != "" {
			simpleTask.SetError(fmt.Errorf(taskData.Error))
		}
		simpleTask.SetResult(taskData.Result)
		return simpleTask, nil

	case "http":
		url, _ := taskData.Data["url"].(string)
		method, _ := taskData.Data["method"].(string)
		timeoutStr, _ := taskData.Data["timeout"].(string)
		timeout, _ := time.ParseDuration(timeoutStr)

		httpTask := task.NewHTTPTask(taskData.ID, taskData.Name, url, method, timeout, taskData.Priority, taskData.MaxRetries)
		// 设置状态和结果
		httpTask.SetStatus(taskData.Status)
		httpTask.SetStartedAt(taskData.StartedAt)
		httpTask.SetCompletedAt(taskData.CompletedAt)
		if taskData.Error != "" {
			httpTask.SetError(fmt.Errorf(taskData.Error))
		}
		httpTask.SetResult(taskData.Result)
		return httpTask, nil

	case "file":
		filePath, _ := taskData.Data["file_path"].(string)
		action, _ := taskData.Data["action"].(string)

		fileTask := task.NewFileProcessTask(taskData.ID, taskData.Name, filePath, action, taskData.Priority, taskData.MaxRetries)
		// 设置状态和结果
		fileTask.SetStatus(taskData.Status)
		fileTask.SetStartedAt(taskData.StartedAt)
		fileTask.SetCompletedAt(taskData.CompletedAt)
		if taskData.Error != "" {
			fileTask.SetError(fmt.Errorf(taskData.Error))
		}
		fileTask.SetResult(taskData.Result)
		return fileTask, nil

	case "batch":
		// 批量任务暂时返回基础任务，因为子任务无法序列化
		return baseTask, nil

	case "retryable":
		// 重试任务暂时返回基础任务，因为包装的任务无法序列化
		return baseTask, nil

	default:
		return baseTask, nil
	}
}
