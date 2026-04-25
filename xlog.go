package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 日志级别常量
const (
	LevelError = 0x01
	LevelWarn  = 0x02
	LevelInfo  = 0x04
	LevelDebug = 0x08
)

// XLogger 异步日志记录器（对应C++的XSimpleLogEx）
type XLogger struct {
	config     *Config
	logChan    chan string    // 日志通道
	wg         sync.WaitGroup // 等待组
	stopChan   chan struct{}  // 停止信号
	logPath    string         // 当前日志文件路径
	mu         sync.Mutex     // 互斥锁
	rotator    *LogRotator    // 日志轮转器
	compressor *LogCompressor // 日志压缩器
}

// NewXLogger 创建新的异步日志记录器
func NewXLogger(config *Config) (*XLogger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	logger := &XLogger{
		config:   config,
		logChan:  make(chan string, config.BufferSize),
		stopChan: make(chan struct{}),
	}

	// 确保日志目录存在
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	logger.logPath = logger.getLogPath()

	// 创建日志轮转器
	if config.AutoCleanup {
		logger.rotator = NewLogRotator(
			config.LogDir,
			config.LogBaseName,
			config.LogSuffix,
			config.RetainDays,
			config.CleanupCycle,
		)
	}

	// 创建日志压缩器
	if config.AutoCompress {
		logger.compressor = NewLogCompressor(
			config.LogDir,
			config.LogBaseName,
			config.LogSuffix,
			config.CompressCycle,
			config.CompressExclude,
		)
	}

	// 启动后台写入goroutine
	logger.wg.Add(1)
	go logger.logWriter()

	// 写入启动标识
	logger.Log("** Log Start ***************************")

	return logger, nil
}

// getLogPath 获取日志文件路径
func (l *XLogger) getLogPath() string {
	var path string

	// 拼接目录和基本名称
	path = filepath.Join(l.config.LogDir, l.config.LogBaseName)

	// 拼接日期
	if l.config.HasDate {
		now := time.Now()
		path += fmt.Sprintf(".%04d%02d%02d", now.Year(), now.Month(), now.Day())
	}

	// 拼接进程ID
	if l.config.HasProcessID {
		path += fmt.Sprintf(".%d", os.Getpid())
	}

	// 拼接后缀
	if l.config.LogSuffix != "" {
		path += "." + l.config.LogSuffix
	}

	return path
}

// getCurrentTimeStr 获取当前时间字符串
func (l *XLogger) getCurrentTimeStr() string {
	now := time.Now()
	return fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d %03d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond()/1000000)
}

// logWriter 后台日志写入goroutine
func (l *XLogger) logWriter() {
	defer l.wg.Done()

	ticker := time.NewTicker(l.config.WriteInterval)
	defer ticker.Stop()

	var buffer []string

	for {
		select {
		case log := <-l.logChan:
			buffer = append(buffer, log)

		case <-ticker.C:
			// 定时写入
			if len(buffer) > 0 {
				l.flushLogs(buffer)
				buffer = buffer[:0]
			}

			l.mu.Lock()
			l.logPath = l.getLogPath()
			currentLogPath := l.logPath
			l.mu.Unlock()

			// 执行日志清理
			if l.rotator != nil {
				l.rotator.Rotate()
			}

			// 执行日志压缩
			if l.compressor != nil {
				l.compressor.Compress(currentLogPath)
			}

		case <-l.stopChan:
			// 停止时继续取尽通道中的剩余日志，避免秒退场景丢尾日志
			for {
				select {
				case log := <-l.logChan:
					buffer = append(buffer, log)
				default:
					if len(buffer) > 0 {
						l.flushLogs(buffer)
					}
					return
				}
			}
		}
	}
}

// flushLogs 将日志批量写入文件
func (l *XLogger) flushLogs(logs []string) {
	l.mu.Lock()
	logPath := l.getLogPath()
	l.logPath = logPath
	l.mu.Unlock()

	file, err := l.openLogFile(logPath)
	if err != nil {
		return
	}
	defer file.Close()

	for _, log := range logs {
		if _, err := file.Write(l.encodeLogContent(log)); err != nil {
			return
		}
	}

	file.Sync()
}

// Log 记录日志（INFO级别）
func (l *XLogger) Log(format string, args ...interface{}) {
	l.LogEx(LevelInfo, format, args...)
}

// LogEx 记录指定级别的日志
func (l *XLogger) LogEx(level int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	// 格式化日志：[时间][进程号/线程号][级别]: 内容
	logLine := fmt.Sprintf("[%s][%6d/%6d][%02d]: %s\r\n",
		l.getCurrentTimeStr(),
		os.Getpid(),
		getGoroutineID(),
		level,
		message)

	// 非阻塞发送到通道
	select {
	case l.logChan <- logLine:
	default:
		// 通道满了，丢弃日志（避免阻塞）
	}
}

// Close 关闭日志记录器
func (l *XLogger) Close() {
	close(l.stopChan)
	l.wg.Wait()
}

func (l *XLogger) encodeLogContent(log string) []byte {
	if l.config.UTF8Format {
		return []byte(log)
	}
	return encodeUTF16LE(log)
}

// getGoroutineID 获取当前goroutine ID（简化版本）
func getGoroutineID() int {
	// 注意：Go没有直接获取goroutine ID的标准方法
	// 这里返回0作为简化处理
	// 如果需要真实ID，可以使用第三方库如：github.com/petermattis/goid
	return 0
}
