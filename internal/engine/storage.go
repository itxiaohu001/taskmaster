package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"task/pkg/task"
)

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
	
	var tasks map[string]task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
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
	
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}
	
	if err := os.WriteFile(f.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
} 