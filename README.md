# Go 任务处理库

一个功能强大、易于使用的 Go 任务处理库，支持任务调度、优先级队列、重试机制、监控和日志记录。

## 特性

- 🚀 **高性能调度**: 基于优先级队列的任务调度
- 🔄 **自动重试**: 支持任务失败自动重试
- 📊 **实时监控**: 任务执行状态监控和指标收集
- 💾 **多种存储**: 支持内存和文件存储
- 📝 **详细日志**: 分级日志记录和文件轮转
- ⚙️ **灵活配置**: 支持配置文件和环境变量
- 🔧 **易于扩展**: 插件化架构，支持自定义任务类型

## 快速开始

### 安装

```bash
go get github.com/your-username/task
```

### 基本使用

```go
package main

import (
    "log"
    "task/pkg/manager"
    "task/pkg/config"
    "task/pkg/task"
)

func main() {
    // 1. 创建配置
    cfg := config.DefaultConfig()
    cfg.Scheduler.Workers = 4

    // 2. 创建任务管理器
    manager := manager.NewManager(cfg)

    // 3. 启动管理器
    if err := manager.Start(); err != nil {
        log.Fatalf("启动失败: %v", err)
    }
    defer manager.Stop()

    // 4. 创建任务
    simpleTask := task.NewSimpleTask("task1", "简单任务", "Hello World", 1, 3)

    // 5. 添加任务
    if err := manager.AddTask(simpleTask); err != nil {
        log.Printf("添加任务失败: %v", err)
    }

    // 6. 等待完成
    manager.Wait()
}
```

## 任务类型

### 简单任务

```go
task := task.NewSimpleTask("id", "name", "message", priority, maxRetries)
```

### HTTP 任务

```go
task := task.NewHTTPTask("id", "name", "https://api.example.com", "GET", 5*time.Second, priority, maxRetries)
```

### 文件处理任务

```go
task := task.NewFileProcessTask("id", "name", "/path/to/file", "read", priority, maxRetries)
```

### 批量任务

```go
subTasks := []task.Task{
    task.NewSimpleTask("sub1", "子任务1", "消息1", 1, 2),
    task.NewSimpleTask("sub2", "子任务2", "消息2", 1, 2),
}
batchTask := task.NewBatchTask("batch1", "批量任务", subTasks, 1, 3)
```

### 自定义任务

```go
type CustomTask struct {
    *task.BaseTask
    data string
}

func (t *CustomTask) Execute(ctx context.Context) error {
    // 实现你的任务逻辑
    fmt.Printf("处理数据: %s\n", t.data)
    return nil
}

task := &CustomTask{
    BaseTask: task.NewBaseTask("custom1", "自定义任务", 1, 3),
    data:     "自定义数据",
}
```

## 配置

### 配置文件

创建 `config.json`:

```json
{
  "scheduler": {
    "workers": 4,
    "max_queue_size": 1000,
    "task_timeout": "30s",
    "retry_interval": "5s",
    "max_retries": 3
  },
  "storage": {
    "type": "file",
    "file_path": "./data/tasks.json"
  },
  "monitor": {
    "enabled": true,
    "metrics_port": 8080,
    "interval": "30s"
  },
  "log": {
    "level": "info",
    "format": "text",
    "output": "file",
    "file_path": "./logs/task.log",
    "max_size": 100,
    "max_backups": 3
  }
}
```

### 环境变量

```bash
export TASK_WORKERS=4
export TASK_MAX_QUEUE_SIZE=1000
export TASK_TIMEOUT=30s
export TASK_STORAGE_TYPE=file
export TASK_LOG_LEVEL=info
```

## API 参考

### Task 接口

```go
type Task interface {
    ID() string
    Name() string
    Execute(ctx context.Context) error
    GetStatus() Status
    SetStatus(status Status)
    GetPriority() int
    GetRetryCount() int
    GetMaxRetries() int
    GetCreatedAt() time.Time
    GetStartedAt() time.Time
    GetCompletedAt() time.Time
    GetError() error
    SetError(err error)
    GetResult() interface{}
    SetResult(result interface{})
}
```

### Manager 接口

```go
type Manager interface {
    Start() error
    Stop() error
    AddTask(t Task) error
    RemoveTask(taskID string) error
    GetTask(taskID string) (Task, bool)
    GetTasks() []Task
    Wait()
}
```

## 任务状态

- `Pending`: 等待执行
- `Running`: 正在执行
- `Completed`: 执行完成
- `Failed`: 执行失败
- `Cancelled`: 已取消
- `Retrying`: 重试中

## 监控

库内置监控功能，可以通过 HTTP 端点查看任务执行指标：

```bash
curl http://localhost:8080/metrics
```

## 日志

支持多种日志级别和输出格式：

- **级别**: debug, info, warn, error
- **格式**: text, json
- **输出**: stdout, file

## 存储

支持多种存储后端：

- **内存存储**: 适用于临时任务
- **文件存储**: 适用于持久化任务

## 最佳实践

### 1. 任务设计

```go
// 好的做法：任务应该是幂等的
func (t *MyTask) Execute(ctx context.Context) error {
    // 检查任务是否已经执行过
    if t.GetStatus() == task.StatusCompleted {
        return nil
    }
    
    // 执行任务逻辑
    return nil
}
```

### 2. 错误处理

```go
func (t *MyTask) Execute(ctx context.Context) error {
    defer func() {
        if r := recover(); r != nil {
            t.SetError(fmt.Errorf("panic: %v", r))
        }
    }()
    
    // 任务逻辑
    return nil
}
```

### 3. 资源管理

```go
func (t *MyTask) Execute(ctx context.Context) error {
    // 设置超时
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    // 任务逻辑
    return nil
}
```

## 示例

查看 `examples/` 目录中的完整示例：

- `basic_usage.go`: 基本使用示例
- `batch_tasks.go`: 批量任务示例
- `custom_tasks.go`: 自定义任务示例

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License 