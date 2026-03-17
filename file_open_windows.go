//go:build windows

package xlog

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/windows"
)

func (l *XLogger) openLogFile(logPath string) (*os.File, error) {
	pathPtr, err := windows.UTF16PtrFromString(logPath)
	if err != nil {
		return nil, fmt.Errorf("转换日志路径失败: %v", err)
	}

	handle, err := windows.CreateFile(
		pathPtr,
		windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_ALWAYS,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}

	file := os.NewFile(uintptr(handle), logPath)
	if file == nil {
		windows.CloseHandle(handle)
		return nil, fmt.Errorf("创建日志文件对象失败")
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

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		file.Close()
		return nil, fmt.Errorf("定位日志文件末尾失败: %v", err)
	}

	return file, nil
}
