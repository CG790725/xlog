package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chen1994/xlog"
)

// 测试辅助函数
func setupTestDir(t *testing.T) string {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("xlog_test_%d", time.Now().UnixNano()))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	return dir
}

func cleanupTestDir(dir string) {
	os.RemoveAll(dir)
}

func countFiles(dir, pattern string) (int, error) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return 0, err
	}
	return len(matches), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ==================== 基础功能测试 ====================

func TestBasicLogging(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "test"
	config.LogSuffix = "log"
	config.HasDate = false
	config.HasProcessID = false
	config.WriteInterval = 100 * time.Millisecond
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, err := xlog.NewXLogger(config)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}

	// 记录一些日志
	logger.Log("测试日志1")
	logger.Log("测试日志2: %s", "参数")
	logger.LogEx(xlog.LevelInfo, "INFO日志")
	logger.LogEx(xlog.LevelWarn, "WARN日志")
	logger.LogEx(xlog.LevelError, "ERROR日志")
	logger.LogEx(xlog.LevelDebug, "DEBUG日志")

	// 等待日志写入
	time.Sleep(500 * time.Millisecond)

	logger.Close()

	// 检查文件是否创建
	logPath := filepath.Join(dir, "test.log")
	if !fileExists(logPath) {
		t.Errorf("日志文件未创建: %s", logPath)
	}

	// 读取文件内容检查
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	contentStr := string(content)
	expectedStrings := []string{
		"** Log Start",
		"测试日志1",
		"测试日志2",
		"INFO日志",
		"WARN日志",
		"ERROR日志",
		"DEBUG日志",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("日志内容缺失: %s", expected)
		}
	}

	t.Logf("✓ 基础日志功能测试通过")
}

// ==================== 日志级别测试 ====================

func TestLogLevels(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "levels"
	config.HasDate = false
	config.WriteInterval = 100 * time.Millisecond
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)

	// 测试所有日志级别
	levels := []int{
		xlog.LevelError,
		xlog.LevelWarn,
		xlog.LevelInfo,
		xlog.LevelDebug,
	}

	for _, level := range levels {
		logger.LogEx(level, "级别%d的日志", level)
	}

	time.Sleep(200 * time.Millisecond)
	logger.Close()

	// 验证日志文件
	content, _ := os.ReadFile(filepath.Join(dir, "levels.log"))
	contentStr := string(content)

	for _, level := range levels {
		expected := fmt.Sprintf("级别%d的日志", level)
		if !strings.Contains(contentStr, expected) {
			t.Errorf("日志级别 %d 未正确记录", level)
		}
	}

	t.Logf("✓ 日志级别测试通过")
}

// ==================== 文件命名测试 ====================

func TestFileNaming(t *testing.T) {
	testCases := []struct {
		name          string
		hasDate       bool
		hasProcessID  bool
		expectDate    bool
		expectPID     bool
	}{
		{"无日期无进程ID", false, false, false, false},
		{"只有日期", true, false, true, false},
		{"只有进程ID", false, true, false, true},
		{"日期和进程ID", true, true, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := setupTestDir(t)
			defer cleanupTestDir(dir)

			config := xlog.DefaultConfig()
			config.LogDir = dir
			config.LogBaseName = "naming"
			config.HasDate = tc.hasDate
			config.HasProcessID = tc.hasProcessID
			config.WriteInterval = 100 * time.Millisecond
			config.AutoCleanup = false
			config.AutoCompress = false

			logger, _ := xlog.NewXLogger(config)
			logger.Log("测试")
			time.Sleep(200 * time.Millisecond)
			logger.Close()

			// 检查文件名
			files, _ := filepath.Glob(filepath.Join(dir, "naming.*"))
			if len(files) == 0 {
				t.Fatalf("没有找到日志文件")
			}

			filename := filepath.Base(files[0])

			// 检查日期
			if tc.expectDate {
				datePattern := time.Now().Format("20060102")
				if !strings.Contains(filename, datePattern) {
					t.Errorf("文件名应该包含日期: %s", filename)
				}
			}

			// 检查进程ID
			if tc.expectPID {
				pid := fmt.Sprintf("%d", os.Getpid())
				if !strings.Contains(filename, pid) {
					t.Errorf("文件名应该包含进程ID: %s", filename)
				}
			}
		})
	}

	t.Logf("✓ 文件命名测试通过")
}

// ==================== 异步写入测试 ====================

func TestAsyncWrite(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "async"
	config.HasDate = false
	config.WriteInterval = 200 * time.Millisecond
	config.BufferSize = 1000
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)

	// 快速记录大量日志
	start := time.Now()
	logCount := 100
	for i := 0; i < logCount; i++ {
		logger.Log("异步日志 #%d", i)
	}
	elapsed := time.Since(start)

	// 异步写入应该非常快（不应该阻塞）
	if elapsed > 100*time.Millisecond {
		t.Errorf("异步日志记录太慢: %v", elapsed)
	}

	// 等待后台写入完成
	time.Sleep(500 * time.Millisecond)
	logger.Close()

	// 验证所有日志都写入了
	content, _ := os.ReadFile(filepath.Join(dir, "async.log"))
	lines := strings.Count(string(content), "\n")

	if lines < logCount {
		t.Errorf("日志丢失: 期望至少 %d 行，实际 %d 行", logCount, lines)
	}

	t.Logf("✓ 异步写入测试通过 (记录 %d 条日志耗时 %v)", logCount, elapsed)
}

// ==================== 并发安全测试 ====================

func TestConcurrency(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "concurrent"
	config.HasDate = false
	config.WriteInterval = 100 * time.Millisecond
	config.BufferSize = 10000
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)

	// 启动多个goroutine并发写日志
	var wg sync.WaitGroup
	goroutines := 10
	logsPerGoroutine := 50

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				logger.Log("Goroutine-%d 日志-%d", id, j)
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(500 * time.Millisecond)
	logger.Close()

	// 验证日志完整性
	content, _ := os.ReadFile(filepath.Join(dir, "concurrent.log"))

	// 统计每个goroutine的日志数量
	for i := 0; i < goroutines; i++ {
		pattern := fmt.Sprintf("Goroutine-%d", i)
		count := strings.Count(string(content), pattern)
		if count != logsPerGoroutine {
			t.Errorf("Goroutine %d 日志数量不正确: 期望 %d, 实际 %d", i, logsPerGoroutine, count)
		}
	}

	t.Logf("✓ 并发安全测试通过 (%d 个goroutine × %d 条日志)", goroutines, logsPerGoroutine)
}

// ==================== 自动清理测试 ====================

func TestAutoCleanup(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	// 创建不同时间的旧日志文件
	now := time.Now()

	// 创建4天前的文件（应该被删除）
	oldFile1 := filepath.Join(dir, fmt.Sprintf("cleanup.%s.log", now.AddDate(0, 0, -4).Format("20060102")))
	os.WriteFile(oldFile1, []byte("old log 4 days ago"), 0644)
	os.Chtimes(oldFile1, now.AddDate(0, 0, -4), now.AddDate(0, 0, -4))

	// 创建3天前的文件（应该被删除）
	oldFile2 := filepath.Join(dir, fmt.Sprintf("cleanup.%s.log", now.AddDate(0, 0, -3).Format("20060102")))
	os.WriteFile(oldFile2, []byte("old log 3 days ago"), 0644)
	os.Chtimes(oldFile2, now.AddDate(0, 0, -3), now.AddDate(0, 0, -3))

	// 创建1天前的文件（应该保留）
	recentFile := filepath.Join(dir, fmt.Sprintf("cleanup.%s.log", now.AddDate(0, 0, -1).Format("20060102")))
	os.WriteFile(recentFile, []byte("recent log"), 0644)
	os.Chtimes(recentFile, now.AddDate(0, 0, -1), now.AddDate(0, 0, -1))

	// 创建日志记录器，保留2天
	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "cleanup"
	config.HasDate = true
	config.WriteInterval = 100 * time.Millisecond
	config.AutoCleanup = true
	config.CleanupCycle = 100 * time.Millisecond
	config.RetainDays = 2 // 只保留2天
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)
	logger.Log("当前日志")

	// 等待清理执行
	time.Sleep(300 * time.Millisecond)
	logger.Close()

	// 验证旧文件被删除
	if fileExists(oldFile1) {
		t.Errorf("4天前的文件应该被删除: %s", oldFile1)
	}
	if fileExists(oldFile2) {
		t.Errorf("3天前的文件应该被删除: %s", oldFile2)
	}

	// 验证最近文件保留
	if !fileExists(recentFile) {
		t.Errorf("1天前的文件应该保留: %s", recentFile)
	}

	t.Logf("✓ 自动清理测试通过")
}

// ==================== 自动压缩测试 ====================

func TestAutoCompression(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	// 创建一些旧日志文件
	oldFiles := []string{
		"compress.20260312.log", // 3天前
		"compress.20260313.log", // 2天前
	}

	for _, filename := range oldFiles {
		path := filepath.Join(dir, filename)
		content := strings.Repeat("这是旧日志内容，需要被压缩。", 100)
		os.WriteFile(path, []byte(content), 0644)

		// 设置修改时间
		dateStr := filename[10:18]
		date, _ := time.Parse("20060102", dateStr)
		os.Chtimes(path, date, date)
	}

	// 创建日志记录器
	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "compress"
	config.HasDate = true
	config.WriteInterval = 100 * time.Millisecond
	config.AutoCleanup = false
	config.AutoCompress = true
	config.CompressCycle = 100 * time.Millisecond
	config.CompressExclude = true // 排除当前日志

	logger, _ := xlog.NewXLogger(config)
	logger.Log("当前日志")

	// 等待压缩执行
	time.Sleep(300 * time.Millisecond)
	logger.Close()

	// 检查zip文件是否创建
	zipFiles, _ := filepath.Glob(filepath.Join(dir, "*.zip"))
	if len(zipFiles) == 0 {
		t.Errorf("没有创建zip压缩文件")
	} else {
		t.Logf("创建了 %d 个zip文件", len(zipFiles))
		for _, f := range zipFiles {
			t.Logf("  - %s", filepath.Base(f))
		}
	}

	// 检查原始文件是否被删除
	for _, filename := range oldFiles {
		path := filepath.Join(dir, filename)
		if fileExists(path) {
			t.Errorf("原始日志文件未被删除: %s", filename)
		}
	}

	t.Logf("✓ 自动压缩测试通过")
}

// ==================== UTF8编码测试 ====================

func TestUTF8Encoding(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	testCases := []struct {
		name       string
		utf8Format bool
	}{
		{"UTF-8格式", true},
		{"UTF-16格式", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := xlog.DefaultConfig()
			config.LogDir = dir
			config.LogBaseName = fmt.Sprintf("utf8_%v", tc.utf8Format)
			config.HasDate = false
			config.UTF8Format = tc.utf8Format
			config.WriteInterval = 100 * time.Millisecond
			config.AutoCleanup = false
			config.AutoCompress = false

			logger, _ := xlog.NewXLogger(config)

			// 记录中文和特殊字符
			logger.Log("中文测试：你好世界")
			logger.Log("特殊字符：!@#$%%^&*()")

			time.Sleep(200 * time.Millisecond)
			logger.Close()

			// 验证文件存在
			logPath := filepath.Join(dir, fmt.Sprintf("utf8_%v.log", tc.utf8Format))
			if !fileExists(logPath) {
				t.Errorf("日志文件未创建")
			}

			// 读取并验证内容
			content, _ := os.ReadFile(logPath)
			if !strings.Contains(string(content), "你好世界") {
				t.Errorf("UTF编码测试失败")
			}
		})
	}

	t.Logf("✓ UTF编码测试通过")
}

// ==================== 缓冲区满测试 ====================

func TestBufferFull(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "buffer"
	config.HasDate = false
	config.WriteInterval = 1 * time.Second // 长间隔，让缓冲区填满
	config.BufferSize = 100                 // 小缓冲区
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)

	// 快速记录超过缓冲区大小的日志
	successCount := 0
	for i := 0; i < 1000; i++ {
		logger.Log("缓冲区测试日志 #%d", i)
		successCount++
	}

	// 等待写入
	time.Sleep(1500 * time.Millisecond)
	logger.Close()

	// 验证至少有部分日志被写入
	content, _ := os.ReadFile(filepath.Join(dir, "buffer.log"))
	if len(content) == 0 {
		t.Errorf("没有日志被写入")
	}

	t.Logf("✓ 缓冲区满测试通过 (成功记录 %d 条)", successCount)
}

// ==================== 日志轮转测试 ====================

func TestLogFileRotation(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "rotation"
	config.HasDate = true
	config.WriteInterval = 100 * time.Millisecond
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)

	// 模拟跨天（通过多次检查文件）
	logger.Log("第一天日志")

	// 检查是否创建了当天的日志文件
	today := time.Now().Format("20060102")
	expectedFile := filepath.Join(dir, fmt.Sprintf("rotation.%s.log", today))

	time.Sleep(200 * time.Millisecond)

	if !fileExists(expectedFile) {
		t.Errorf("日志文件未创建: %s", expectedFile)
	}

	logger.Close()

	t.Logf("✓ 日志轮转测试通过")
}

// ==================== Close测试 ====================

func TestLoggerClose(t *testing.T) {
	dir := setupTestDir(t)
	defer cleanupTestDir(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.LogBaseName = "close"
	config.HasDate = false
	config.WriteInterval = 50 * time.Millisecond
	config.BufferSize = 1000
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)

	// 记录多条日志
	logger.Log("关闭前的日志1")
	logger.Log("关闭前的日志2")
	logger.Log("关闭前的日志3")

	// 给一点时间让日志进入buffer
	time.Sleep(10 * time.Millisecond)

	// Close应该等待所有日志写入完成
	logger.Close()

	// 验证日志已写入
	logPath := filepath.Join(dir, "close.log")
	if !fileExists(logPath) {
		t.Fatalf("日志文件未创建")
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "关闭前的日志1") {
		t.Errorf("Close()未正确刷新缓冲区，缺少日志1")
	}
	if !strings.Contains(contentStr, "关闭前的日志2") {
		t.Errorf("Close()未正确刷新缓冲区，缺少日志2")
	}
	if !strings.Contains(contentStr, "关闭前的日志3") {
		t.Errorf("Close()未正确刷新缓冲区，缺少日志3")
	}

	t.Logf("✓ Close测试通过")
}

// ==================== 默认配置测试 ====================

func TestDefaultConfig(t *testing.T) {
	config := xlog.DefaultConfig()

	// 验证默认值
	tests := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{"LogSuffix", config.LogSuffix, "log"},
		{"HasDate", config.HasDate, true},
		{"HasProcessID", config.HasProcessID, false},
		{"UTF8Format", config.UTF8Format, true},
		{"WriteInterval", config.WriteInterval, 5 * time.Second},
		{"BufferSize", config.BufferSize, 10000},
		{"AutoCleanup", config.AutoCleanup, false},
		{"RetainDays", config.RetainDays, 3},
		{"AutoCompress", config.AutoCompress, false},
		{"CompressExclude", config.CompressExclude, true},
	}

	for _, tt := range tests {
		if tt.actual != tt.expected {
			t.Errorf("默认配置 %s 错误: 期望 %v, 实际 %v", tt.name, tt.expected, tt.actual)
		}
	}

	t.Logf("✓ 默认配置测试通过")
}

// ==================== 性能测试 ====================

func BenchmarkLogging(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("xlog_bench_%d", time.Now().UnixNano()))
	os.MkdirAll(dir, 0755)
	defer os.RemoveAll(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.HasDate = false
	config.WriteInterval = 100 * time.Millisecond
	config.BufferSize = 100000
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Log("性能测试日志 #%d", i)
	}
}

func BenchmarkConcurrentLogging(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("xlog_bench_%d", time.Now().UnixNano()))
	os.MkdirAll(dir, 0755)
	defer os.RemoveAll(dir)

	config := xlog.DefaultConfig()
	config.LogDir = dir
	config.HasDate = false
	config.WriteInterval = 100 * time.Millisecond
	config.BufferSize = 100000
	config.AutoCleanup = false
	config.AutoCompress = false

	logger, _ := xlog.NewXLogger(config)
	defer logger.Close()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Log("并发性能测试 #%d", i)
			i++
		}
	})
}

// ==================== 主函数 ====================

func main() {
	// 运行所有测试
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{"TestBasicLogging", TestBasicLogging},
			{"TestLogLevels", TestLogLevels},
			{"TestFileNaming", TestFileNaming},
			{"TestAsyncWrite", TestAsyncWrite},
			{"TestConcurrency", TestConcurrency},
			{"TestAutoCleanup", TestAutoCleanup},
			{"TestAutoCompression", TestAutoCompression},
			{"TestUTF8Encoding", TestUTF8Encoding},
			{"TestBufferFull", TestBufferFull},
			{"TestLogFileRotation", TestLogFileRotation},
			{"TestLoggerClose", TestLoggerClose},
			{"TestDefaultConfig", TestDefaultConfig},
		},
		[]testing.InternalBenchmark{
			{"BenchmarkLogging", BenchmarkLogging},
			{"BenchmarkConcurrentLogging", BenchmarkConcurrentLogging},
		},
		nil,
	)
}
