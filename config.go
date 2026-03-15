package xlog

import (
	"os"
	"path/filepath"
	"time"
)

// Config 日志配置
type Config struct {
	// 基础配置
	LogDir       string // 日志目录
	LogBaseName  string // 日志基本名称
	LogSuffix    string // 日志后缀（默认.log）
	HasDate      bool   // 文件名是否包含日期
	HasProcessID bool   // 文件名是否包含进程ID
	UTF8Format   bool   // 是否UTF8格式（默认true）

	// 异步写入配置
	WriteInterval time.Duration // 写入间隔（默认5秒）
	BufferSize    int           // 缓冲区大小（默认10000）

	// 自动清理配置
	AutoCleanup   bool          // 是否启用自动清理
	CleanupCycle   time.Duration // 清理周期（默认6分钟）
	RetainDays     int           // 保留天数（默认3天）

	// 自动压缩配置
	AutoCompress    bool          // 是否启用自动压缩
	CompressCycle    time.Duration // 压缩周期（默认10分钟）
	CompressExclude  bool          // 是否排除当天日志（默认true）
}

// DefaultConfig 创建默认配置
func DefaultConfig() *Config {
	exePath, _ := os.Executable()
	exeName := filepath.Base(exePath)
	exeDir := filepath.Dir(exePath)

	return &Config{
		LogDir:        exeDir,
		LogBaseName:   exeName,
		LogSuffix:     "log",
		HasDate:       true,
		HasProcessID:  false,
		UTF8Format:    true,
		WriteInterval: 5 * time.Second,
		BufferSize:    10000,
		AutoCleanup:   false,
		CleanupCycle:  6 * time.Minute,
		RetainDays:    3,
		AutoCompress:  false,
		CompressCycle: 10 * time.Minute,
		CompressExclude: true,
	}
}
