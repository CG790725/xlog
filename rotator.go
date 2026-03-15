package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

// LogRotator 日志轮转和清理器
type LogRotator struct {
	logDir      string
	logBaseName string
	logSuffix   string
	retainDays  int
	cleanupCycle time.Duration
	lastCleanup time.Time
}

// NewLogRotator 创建日志轮转器
func NewLogRotator(logDir, logBaseName, logSuffix string, retainDays int, cleanupCycle time.Duration) *LogRotator {
	return &LogRotator{
		logDir:       logDir,
		logBaseName:  logBaseName,
		logSuffix:    logSuffix,
		retainDays:   retainDays,
		cleanupCycle: cleanupCycle,
		lastCleanup:  time.Now(),
	}
}

// Rotate 执行日志清理（如果达到清理周期）
func (r *LogRotator) Rotate() {
	// 检查是否到达清理周期
	if time.Since(r.lastCleanup) < r.cleanupCycle {
		return
	}

	r.lastCleanup = time.Now()
	r.deleteExcessiveLogs()
}

// deleteExcessiveLogs 删除过期的日志文件
func (r *LogRotator) deleteExcessiveLogs() {
	// 计算清理临界时间点
	cutoffTime := time.Now().AddDate(0, 0, -r.retainDays+1)
	cutoffTime = time.Date(cutoffTime.Year(), cutoffTime.Month(), cutoffTime.Day(), 0, 0, 0, 0, cutoffTime.Location())

	// 查找匹配的日志文件
	pattern := fmt.Sprintf("%s.*", r.logBaseName)
	matches, err := filepath.Glob(filepath.Join(r.logDir, pattern))
	if err != nil {
		return
	}

	// 遍历并删除过期文件
	for _, file := range matches {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		// 检查文件修改时间
		if info.ModTime().Before(cutoffTime) {
			os.Remove(file)
		}
	}
}

// GetLogFiles 获取所有日志文件（按时间排序）
func (r *LogRotator) GetLogFiles() ([]string, error) {
	pattern := fmt.Sprintf("%s.*.%s", r.logBaseName, r.logSuffix)
	matches, err := filepath.Glob(filepath.Join(r.logDir, pattern))
	if err != nil {
		return nil, err
	}

	// 按文件名排序（通常文件名包含日期）
	sort.Strings(matches)

	return matches, nil
}

// ParseLogDate 从日志文件名中解析日期
func (r *LogRotator) ParseLogDate(filename string) (time.Time, error) {
	// 匹配格式：basename.YYYYMMDD.suffix
	re := regexp.MustCompile(`\.(\d{8})\.`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return time.Time{}, fmt.Errorf("无法解析日期")
	}

	dateStr := matches[1]
	return time.Parse("20060102", dateStr)
}

// GetDiskUsage 获取日志目录磁盘使用情况
func (r *LogRotator) GetDiskUsage() (int64, error) {
	var totalSize int64

	files, err := r.GetLogFiles()
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		totalSize += info.Size()
	}

	return totalSize, nil
}
