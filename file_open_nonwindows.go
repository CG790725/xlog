//go:build !windows

package xlog

import (
	"fmt"
	"os"
)

func (l *XLogger) openLogFile(logPath string) (*os.File, error) {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("获取日志文件信息失败: %v", err)
	}

	if fileInfo.Size() == 0 && !l.config.UTF8Format {
		if _, err := file.Write([]byte{0xFF, 0xFE}); err != nil {
			file.Close()
			return nil, fmt.Errorf("写入UTF-16 BOM失败: %v", err)
		}
	}

	return file, nil
}
