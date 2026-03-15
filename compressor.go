package xlog

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogCompressor 日志压缩器
type LogCompressor struct {
	logDir         string
	logBaseName    string
	logSuffix      string
	compressCycle   time.Duration
	excludeCurrent  bool
	lastCompress   time.Time
}

// NewLogCompressor 创建日志压缩器
func NewLogCompressor(logDir, logBaseName, logSuffix string, compressCycle time.Duration, excludeCurrent bool) *LogCompressor {
	return &LogCompressor{
		logDir:        logDir,
		logBaseName:   logBaseName,
		logSuffix:     logSuffix,
		compressCycle:  compressCycle,
		excludeCurrent: excludeCurrent,
		lastCompress:  time.Now(),
	}
}

// Compress 执行日志压缩（如果达到压缩周期）
func (c *LogCompressor) Compress(currentLogPath string) {
	// 检查是否到达压缩周期
	if time.Since(c.lastCompress) < c.compressCycle {
		return
	}

	c.lastCompress = time.Now()
	c.compressLogs(currentLogPath)
}

// compressLogs 压缩日志文件
func (c *LogCompressor) compressLogs(currentLogPath string) {
	// 查找符合条件的日志文件
	pattern := fmt.Sprintf("%s.*.%s", c.logBaseName, c.logSuffix)
	matches, err := filepath.Glob(filepath.Join(c.logDir, pattern))
	if err != nil {
		return
	}

	for _, logFile := range matches {
		// 排除当前日志文件
		if c.excludeCurrent && logFile == currentLogPath {
			continue
		}

		// 检查文件是否已经被压缩
		if strings.HasSuffix(logFile, ".zip") {
			continue
		}

		// 生成压缩文件路径（避免覆盖）
		zipFile := c.getUniqueZipPath(logFile)

		// 执行压缩
		if err := c.compressFile(logFile, zipFile); err != nil {
			// 压缩失败，删除失败的zip文件
			os.Remove(zipFile)
			continue
		}

		// 压缩成功，删除原日志文件
		os.Remove(logFile)
	}
}

// getUniqueZipPath 获取唯一的zip文件路径
func (c *LogCompressor) getUniqueZipPath(logFile string) string {
	zipFile := logFile + ".zip"

	// 如果zip文件已存在，添加序号
	fileNo := 1
	for {
		if _, err := os.Stat(zipFile); os.IsNotExist(err) {
			return zipFile
		}
		zipFile = fmt.Sprintf("%s(%d).zip", logFile, fileNo)
		fileNo++
	}
}

// compressFile 压缩单个文件
func (c *LogCompressor) compressFile(srcFile, dstFile string) error {
	// 打开源文件
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	// 获取源文件信息
	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}

	// 创建目标zip文件
	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 创建zip writer
	zipWriter := zip.NewWriter(dst)
	defer zipWriter.Close()

	// 创建zip文件条目
	header, err := zip.FileInfoHeader(srcInfo)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(srcFile)
	header.Method = zip.Deflate // 使用Deflate压缩

	// 写入文件头
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// 复制文件内容
	_, err = io.Copy(writer, src)
	return err
}

// Decompress 解压日志文件（可选功能）
func (c *LogCompressor) Decompress(zipFile, destDir string) error {
	// 打开zip文件
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 遍历zip文件中的条目
	for _, file := range reader.File {
		// 打开zip中的文件
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// 创建目标文件
		destPath := filepath.Join(destDir, file.Name)
		dest, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dest.Close()

		// 复制文件内容
		_, err = io.Copy(dest, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCompressedFiles 获取所有压缩文件
func (c *LogCompressor) GetCompressedFiles() ([]string, error) {
	pattern := fmt.Sprintf("%s.*.zip", c.logBaseName)
	return filepath.Glob(filepath.Join(c.logDir, pattern))
}
