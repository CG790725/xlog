// +build !windows

package xlog

import (
	"os"
	"syscall"
	"time"
)

// CleanupBySize 按大小清理日志（当磁盘空间不足时）
func (r *LogRotator) CleanupBySize(minReservedSpaceGB int, minRetainDays int) error {
	// 获取磁盘剩余空间
	var stat syscall.Statfs_t
	if err := syscall.Statfs(r.logDir, &stat); err != nil {
		return err
	}

	// 计算可用空间（GB）
	freeSpaceGB := int(stat.Bavail*uint64(stat.Bsize) / 1024 / 1024 / 1024)

	// 如果空间充足，不需要清理
	if freeSpaceGB >= minReservedSpaceGB {
		return nil
	}

	// 空间不足，清理旧日志
	minRetainTime := time.Now().AddDate(0, 0, -minRetainDays+1)
	minRetainTime = time.Date(minRetainTime.Year(), minRetainTime.Month(), minRetainTime.Day(), 0, 0, 0, 0, minRetainTime.Location())

	files, err := r.GetLogFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		// 删除超过最小保留天数的日志
		if info.ModTime().Before(minRetainTime) {
			os.Remove(file)
		}
	}

	return nil
}
