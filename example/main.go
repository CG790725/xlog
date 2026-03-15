package main

import (
	"time"

	"github.com/chen1994/xlog"
)

func main() {
	// 创建默认配置
	config := xlog.DefaultConfig()

	// 自定义配置
	config.LogDir = "./logs"
	config.LogBaseName = "myapp"
	config.LogSuffix = "log"
	config.HasDate = true
	config.HasProcessID = false
	config.UTF8Format = true

	// 异步写入配置
	config.WriteInterval = 5 * time.Second
	config.BufferSize = 10000

	// 自动清理配置（可选）
	config.AutoCleanup = true
	config.CleanupCycle = 6 * time.Minute
	config.RetainDays = 3

	// 自动压缩配置（可选）
	config.AutoCompress = true
	config.CompressCycle = 10 * time.Minute
	config.CompressExclude = true

	// 创建日志记录器
	logger, err := xlog.NewXLogger(config)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// 记录不同级别的日志
	logger.Log("这是一条普通日志")

	logger.LogEx(xlog.LevelInfo, "这是一条INFO日志")
	logger.LogEx(xlog.LevelWarn, "这是一条WARN日志")
	logger.LogEx(xlog.LevelError, "这是一条ERROR日志")
	logger.LogEx(xlog.LevelDebug, "这是一条DEBUG日志")

	// 模拟应用程序运行
	for i := 0; i < 10; i++ {
		logger.Log("循环迭代 %d: %s", i, time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(1 * time.Second)
	}

	logger.Log("应用程序结束")
}
