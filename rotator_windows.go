// +build windows

package xlog

import (
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// CleanupBySize 按大小清理日志（当磁盘空间不足时）
func (r *LogRotator) CleanupBySize(minReservedSpaceGB int, minRetainDays int) error {
	// 获取磁盘剩余空间
	var freeBytesAvailable uint64
	var totalNumberOfBytes uint64
	var totalNumberOfFreeBytes uint64

	// 转换路径为Windows指针
	pathPtr, err := windows.UTF16PtrFromString(r.logDir)
	if err != nil {
		return err
	}

	// 调用Windows API获取磁盘空间
	err = windows.GetDiskFreeSpaceEx(
		pathPtr,
		&freeBytesAvailable,
		&totalNumberOfBytes,
		&totalNumberOfFreeBytes,
	)
	if err != nil {
		return err
	}

	// 计算可用空间（GB）
	freeSpaceGB := int(freeBytesAvailable / 1024 / 1024 / 1024)

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
