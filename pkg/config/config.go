package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config 配置结构
type Config struct {
	// 调度器配置
	Scheduler SchedulerConfig `json:"scheduler"`
	
	// 存储配置
	Storage StorageConfig `json:"storage"`
	
	// 监控配置
	Monitor MonitorConfig `json:"monitor"`
	
	// 日志配置
	Log LogConfig `json:"log"`
}

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
	Workers        int           `json:"workers"`         // 工作协程数量
	MaxQueueSize   int           `json:"max_queue_size"`  // 最大队列大小
	TaskTimeout    time.Duration `json:"task_timeout"`    // 任务超时时间
	RetryInterval  time.Duration `json:"retry_interval"`  // 重试间隔
	MaxRetries     int           `json:"max_retries"`     // 最大重试次数
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type     string `json:"type"`      // 存储类型: memory, file, database
	FilePath  string `json:"file_path"` // 文件存储路径
	Database  string `json:"database"`  // 数据库连接字符串
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	Enabled     bool          `json:"enabled"`      // 是否启用监控
	MetricsPort int           `json:"metrics_port"` // 监控端口
	Interval    time.Duration `json:"interval"`     // 监控间隔
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`       // 日志级别: debug, info, warn, error
	Format     string `json:"format"`      // 日志格式: json, text
	Output     string `json:"output"`      // 输出目标: stdout, file
	FilePath   string `json:"file_path"`   // 日志文件路径
	MaxSize    int    `json:"max_size"`    // 最大文件大小(MB)
	MaxBackups int    `json:"max_backups"` // 最大备份文件数
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Scheduler: SchedulerConfig{
			Workers:        4,
			MaxQueueSize:   1000,
			TaskTimeout:    30 * time.Second,
			RetryInterval:  5 * time.Second,
			MaxRetries:     3,
		},
		Storage: StorageConfig{
			Type:    "memory",
			FilePath: "./data/tasks.json",
		},
		Monitor: MonitorConfig{
			Enabled:     true,
			MetricsPort: 8080,
			Interval:    30 * time.Second,
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "text",
			Output:     "stdout",
			FilePath:   "./logs/task.log",
			MaxSize:    100,
			MaxBackups: 3,
		},
	}
}

// LoadFromFile 从文件加载配置
func LoadFromFile(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() *Config {
	config := DefaultConfig()

	// 调度器配置
	if workers := os.Getenv("TASK_WORKERS"); workers != "" {
		if w, err := strconv.Atoi(workers); err == nil {
			config.Scheduler.Workers = w
		}
	}

	if maxQueueSize := os.Getenv("TASK_MAX_QUEUE_SIZE"); maxQueueSize != "" {
		if mqs, err := strconv.Atoi(maxQueueSize); err == nil {
			config.Scheduler.MaxQueueSize = mqs
		}
	}

	if taskTimeout := os.Getenv("TASK_TIMEOUT"); taskTimeout != "" {
		if tt, err := time.ParseDuration(taskTimeout); err == nil {
			config.Scheduler.TaskTimeout = tt
		}
	}

	if retryInterval := os.Getenv("TASK_RETRY_INTERVAL"); retryInterval != "" {
		if ri, err := time.ParseDuration(retryInterval); err == nil {
			config.Scheduler.RetryInterval = ri
		}
	}

	if maxRetries := os.Getenv("TASK_MAX_RETRIES"); maxRetries != "" {
		if mr, err := strconv.Atoi(maxRetries); err == nil {
			config.Scheduler.MaxRetries = mr
		}
	}

	// 存储配置
	if storageType := os.Getenv("TASK_STORAGE_TYPE"); storageType != "" {
		config.Storage.Type = storageType
	}

	if filePath := os.Getenv("TASK_FILE_PATH"); filePath != "" {
		config.Storage.FilePath = filePath
	}

	if database := os.Getenv("TASK_DATABASE"); database != "" {
		config.Storage.Database = database
	}

	// 监控配置
	if enabled := os.Getenv("TASK_MONITOR_ENABLED"); enabled != "" {
		config.Monitor.Enabled = enabled == "true"
	}

	if metricsPort := os.Getenv("TASK_METRICS_PORT"); metricsPort != "" {
		if mp, err := strconv.Atoi(metricsPort); err == nil {
			config.Monitor.MetricsPort = mp
		}
	}

	// 日志配置
	if level := os.Getenv("TASK_LOG_LEVEL"); level != "" {
		config.Log.Level = level
	}

	if format := os.Getenv("TASK_LOG_FORMAT"); format != "" {
		config.Log.Format = format
	}

	if output := os.Getenv("TASK_LOG_OUTPUT"); output != "" {
		config.Log.Output = output
	}

	if filePath := os.Getenv("TASK_LOG_FILE_PATH"); filePath != "" {
		config.Log.FilePath = filePath
	}

	return config
}

// SaveToFile 保存配置到文件
func (c *Config) SaveToFile(filepath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Scheduler.Workers <= 0 {
		return fmt.Errorf("workers must be greater than 0")
	}

	if c.Scheduler.MaxQueueSize <= 0 {
		return fmt.Errorf("max_queue_size must be greater than 0")
	}

	if c.Scheduler.TaskTimeout <= 0 {
		return fmt.Errorf("task_timeout must be greater than 0")
	}

	if c.Scheduler.MaxRetries < 0 {
		return fmt.Errorf("max_retries must be non-negative")
	}

	return nil
} 