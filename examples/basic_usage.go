package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"task/pkg/manager"
	"task/pkg/config"
	"task/pkg/task"
)

func main() {
	// 1. 创建配置
	cfg := config.DefaultConfig()
	cfg.Scheduler.Workers = 2
	cfg.Log.Level = "info"

	// 2. 创建任务管理器
	manager := manager.NewManager(cfg)

	// 3. 启动管理器
	if err := manager.Start(); err != nil {
		log.Fatalf("启动任务管理器失败: %v", err)
	}
	defer manager.Stop()

	// 4. 创建任务
	tasks := []task.Task{
		task.NewSimpleTask("task1", "简单任务1", "Hello World", 1, 3),
		task.NewSimpleTask("task2", "简单任务2", "Hello Go", 2, 3),
		task.NewHTTPTask("task3", "HTTP任务", "https://api.example.com", "GET", 5*time.Second, 1, 3),
		task.NewFileProcessTask("task4", "文件任务", "/tmp/test.txt", "read", 3, 3),
	}

	// 5. 添加任务到管理器
	for _, t := range tasks {
		if err := manager.AddTask(t); err != nil {
			log.Printf("添加任务失败: %v", err)
		}
	}

	// 6. 等待所有任务完成
	manager.Wait()

	fmt.Println("所有任务执行完成!")
}

// 示例：批量任务
func batchTaskExample() {
	cfg := config.DefaultConfig()
	manager := manager.NewManager(cfg)
	
	if err := manager.Start(); err != nil {
		log.Fatalf("启动任务管理器失败: %v", err)
	}
	defer manager.Stop()

	// 创建子任务
	subTasks := []task.Task{
		task.NewSimpleTask("sub1", "子任务1", "子任务1消息", 1, 2),
		task.NewSimpleTask("sub2", "子任务2", "子任务2消息", 1, 2),
		task.NewSimpleTask("sub3", "子任务3", "子任务3消息", 1, 2),
	}

	// 创建批量任务
	batchTask := task.NewBatchTask("batch1", "批量任务", subTasks, 1, 3)
	
	if err := manager.AddTask(batchTask); err != nil {
		log.Printf("添加批量任务失败: %v", err)
	}

	manager.Wait()
}

// 示例：带重试的任务
func retryTaskExample() {
	cfg := config.DefaultConfig()
	manager := manager.NewManager(cfg)
	
	if err := manager.Start(); err != nil {
		log.Fatalf("启动任务管理器失败: %v", err)
	}
	defer manager.Stop()

	// 创建一个可能失败的任务
	failingTask := &task.BaseTask{}
	failingTask.SetStatus(task.StatusPending)
	failingTask.SetError(fmt.Errorf("模拟失败"))

	// 包装为可重试任务
	retryTask := task.NewRetryableTask(failingTask, 3)
	
	if err := manager.AddTask(retryTask); err != nil {
		log.Printf("添加重试任务失败: %v", err)
	}

	manager.Wait()
}

// 示例：自定义任务
func customTaskExample() {
	cfg := config.DefaultConfig()
	manager := manager.NewManager(cfg)
	
	if err := manager.Start(); err != nil {
		log.Fatalf("启动任务管理器失败: %v", err)
	}
	defer manager.Stop()

	// 创建自定义任务
	customTask := &CustomTask{
		BaseTask: task.NewBaseTask("custom1", "自定义任务", 1, 3),
		data:     "自定义数据",
	}

	if err := manager.AddTask(customTask); err != nil {
		log.Printf("添加自定义任务失败: %v", err)
	}

	manager.Wait()
}

// CustomTask 自定义任务示例
type CustomTask struct {
	*task.BaseTask
	data string
}

func (t *CustomTask) Execute(ctx context.Context) error {
	fmt.Printf("执行自定义任务: %s, 数据: %s\n", t.Name(), t.data)
	
	// 模拟处理时间
	time.Sleep(1 * time.Second)
	
	// 设置结果
	t.SetResult(fmt.Sprintf("处理结果: %s", t.data))
	
	return nil
} 