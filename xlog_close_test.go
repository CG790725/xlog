package xlog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCloseFlushesPendingChannelLogs(t *testing.T) {
	tempDir := t.TempDir()

	cfg := DefaultConfig()
	cfg.LogDir = tempDir
	cfg.LogBaseName = "closeflush"
	cfg.LogSuffix = "log"
	cfg.HasDate = false
	cfg.HasProcessID = false
	cfg.UTF8Format = true
	cfg.WriteInterval = time.Hour
	cfg.BufferSize = 16
	cfg.AutoCleanup = false
	cfg.AutoCompress = false

	logger, err := NewXLogger(cfg)
	if err != nil {
		t.Fatalf("NewXLogger() error = %v", err)
	}

	logger.Log("first pending log")
	logger.LogEx(LevelError, "second pending log")
	logger.Close()

	data, err := os.ReadFile(filepath.Join(tempDir, "closeflush.log"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "first pending log") {
		t.Fatalf("log file missing first pending log: %s", content)
	}
	if !strings.Contains(content, "second pending log") {
		t.Fatalf("log file missing second pending log: %s", content)
	}
}
