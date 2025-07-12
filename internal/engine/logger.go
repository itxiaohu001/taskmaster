package engine

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"task/pkg/config"
	"task/pkg/task"
)

// Logger 日志实现
type Logger struct {
	config config.LogConfig
	logger *log.Logger
	file   *os.File
}

// NewLogger 创建新的日志器
func NewLogger(cfg config.LogConfig) *Logger {
	l := &Logger{
		config: cfg,
	}

	// 设置输出
	switch cfg.Output {
	case "file":
		if err := l.setupFileOutput(); err != nil {
			// 如果文件输出失败，回退到标准输出
			l.logger = log.New(os.Stdout, "", log.LstdFlags)
		}
	default:
		l.logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	return l
}

// setupFileOutput 设置文件输出
func (l *Logger) setupFileOutput() error {
	// 确保目录存在
	dir := filepath.Dir(l.config.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 打开日志文件
	file, err := os.OpenFile(l.config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.file = file
	l.logger = log.New(file, "", log.LstdFlags)
	return nil
}

// formatMessage 格式化消息和键值对参数
func (l *Logger) formatMessage(msg string, args ...interface{}) string {
	if len(args) == 0 {
		return msg
	}

	// 检查是否是键值对格式
	if len(args)%2 == 0 {
		var pairs []string
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				pairs = append(pairs, fmt.Sprintf("%v=%v", args[i], args[i+1]))
			}
		}
		if len(pairs) > 0 {
			return msg + " " + strings.Join(pairs, " ")
		}
	}

	// 如果不是键值对，使用标准格式化
	return fmt.Sprintf(msg, args...)
}

// Debug 输出调试日志
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.shouldLog("debug") {
		formattedMsg := l.formatMessage(msg, args...)
		l.logger.Printf("[DEBUG] %s", formattedMsg)
	}
}

// Info 输出信息日志
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.shouldLog("info") {
		formattedMsg := l.formatMessage(msg, args...)
		l.logger.Printf("[INFO] %s", formattedMsg)
	}
}

// Warn 输出警告日志
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.shouldLog("warn") {
		formattedMsg := l.formatMessage(msg, args...)
		l.logger.Printf("[WARN] %s", formattedMsg)
	}
}

// Error 输出错误日志
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.shouldLog("error") {
		formattedMsg := l.formatMessage(msg, args...)
		l.logger.Printf("[ERROR] %s", formattedMsg)
	}
}

// shouldLog 判断是否应该输出日志
func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	configLevel := levels[l.config.Level]
	currentLevel := levels[level]

	return currentLevel >= configLevel
}

// Close 关闭日志器
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// MetricsMonitor 监控实现
type MetricsMonitor struct {
	port int
	// 这里可以添加更多的监控指标
	startedTasks   int64
	completedTasks int64
	failedTasks    int64
}

// NewMetricsMonitor 创建监控器
func NewMetricsMonitor(port int) *MetricsMonitor {
	return &MetricsMonitor{
		port: port,
	}
}

// Start 启动监控
func (m *MetricsMonitor) Start() error {
	// 这里可以实现HTTP服务器来暴露监控指标
	// 为了简化，这里只是记录启动
	return nil
}

// Stop 停止监控
func (m *MetricsMonitor) Stop() error {
	// 停止监控服务
	return nil
}

// RecordTaskStarted 记录任务开始
func (m *MetricsMonitor) RecordTaskStarted(t task.Task) {
	// 这里可以记录更详细的指标
}

// RecordTaskCompleted 记录任务完成
func (m *MetricsMonitor) RecordTaskCompleted(t task.Task) {
	// 这里可以记录更详细的指标
}

// RecordTaskFailed 记录任务失败
func (m *MetricsMonitor) RecordTaskFailed(t task.Task) {
	// 这里可以记录更详细的指标
} 