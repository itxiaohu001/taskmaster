package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/itxiaohu001/taskmaster/internal/engine"
	"github.com/itxiaohu001/taskmaster/pkg/config"
	"github.com/itxiaohu001/taskmaster/pkg/manager"
	"github.com/itxiaohu001/taskmaster/pkg/task"
)

func main() {
	// 1. 创建配置
	cfg := config.DefaultConfig()

	// 2. 创建事件处理器
	customEventHandler := CustomEventHandler{}

	// 3. 创建调度器
	scheduler := engine.NewScheduler(cfg.Scheduler.Workers, customEventHandler)

	// 4. 创建任务管理器
	taskManager := manager.NewManager(scheduler)

	// 5. 启动管理器
	if err := taskManager.Start(); err != nil {
		log.Fatalf("启动任务管理器失败: %v", err)
	}
	defer taskManager.Stop()

	// 6. 创建任务
	customTask := &CustomTask{
		BaseTask: task.NewBaseTask("custom1", "自定义任务", 1, 3),
		data:     "自定义数据",
	}

	// 7. 添加任务到管理器
	if err := taskManager.AddTask(customTask); err != nil {
		log.Printf("添加任务失败: %v", err)
	}

	// 8. 阻塞（等待所有 worker 退出）
	taskManager.Wait()
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

type CustomEventHandler struct{}

func (c CustomEventHandler) OnTaskStarted(t task.Task) {
	fmt.Printf("任务开始: %s\n", t.Name())
}

func (c CustomEventHandler) OnTaskCompleted(t task.Task) {
	fmt.Printf("任务完成: %s, 结果: %v\n", t.Name(), t.GetResult())
}

func (c CustomEventHandler) OnTaskFailed(t task.Task) {
	fmt.Printf("任务失败: %s, 错误: %v\n", t.Name(), t.GetError())
}

func (c CustomEventHandler) OnTaskCancelled(t task.Task) {
	fmt.Printf("任务取消: %s\n", t.Name())
}
