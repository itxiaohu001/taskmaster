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
}

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
	Workers       int           `json:"workers"`        // 工作协程数量
	TaskTimeout   time.Duration `json:"task_timeout"`   // 任务超时时间
	RetryInterval time.Duration `json:"retry_interval"` // 重试间隔
	MaxRetries    int           `json:"max_retries"`    // 最大重试次数
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type     string `json:"type"`      // 存储类型: memory, file, database
	FilePath string `json:"file_path"` // 文件存储路径
	Database string `json:"database"`  // 数据库连接字符串
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Scheduler: SchedulerConfig{
			Workers:       4,
			TaskTimeout:   30 * time.Second,
			RetryInterval: 5 * time.Second,
			MaxRetries:    3,
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

	if c.Scheduler.TaskTimeout <= 0 {
		return fmt.Errorf("task_timeout must be greater than 0")
	}

	if c.Scheduler.MaxRetries < 0 {
		return fmt.Errorf("max_retries must be non-negative")
	}

	return nil
}
