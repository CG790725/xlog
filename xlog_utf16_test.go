package xlog

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestXLoggerWritesUTF16LEWithBOMWhenUTF8Disabled(t *testing.T) {
	tempDir := t.TempDir()

	cfg := DefaultConfig()
	cfg.LogDir = tempDir
	cfg.LogBaseName = "utf16-test"
	cfg.LogSuffix = "log"
	cfg.HasDate = false
	cfg.HasProcessID = false
	cfg.UTF8Format = false
	cfg.WriteInterval = 20 * time.Millisecond
	cfg.BufferSize = 16
	cfg.AutoCleanup = false
	cfg.AutoCompress = false

	logger, err := NewXLogger(cfg)
	if err != nil {
		t.Fatalf("NewXLogger failed: %v", err)
	}

	testMessage := "UTF16 检查 ABC 123"
	logger.Log(testMessage)

	time.Sleep(80 * time.Millisecond)
	logger.Close()

	logPath := filepath.Join(tempDir, "utf16-test.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if len(data) < 2 {
		t.Fatalf("log file too short: %d", len(data))
	}
	if !bytes.Equal(data[:2], []byte{0xFF, 0xFE}) {
		t.Fatalf("expected UTF-16LE BOM, got %v", data[:2])
	}

	expectedStart := encodeUTF16LE("** Log Start ***************************\r\n")
	if !bytes.Contains(data[2:], expectedStart) {
		t.Fatalf("expected UTF-16LE encoded start marker")
	}

	if !bytes.Contains(data[2:], encodeUTF16LE(testMessage)) {
		t.Fatalf("expected UTF-16LE encoded test message")
	}

	if bytes.Contains(data, []byte(testMessage)) {
		t.Fatalf("found raw UTF-8 bytes in UTF-16 log output")
	}
}

func TestEncodeUTF16LEProducesLittleEndianBytes(t *testing.T) {
	encoded := encodeUTF16LE("A中")
	expected := []byte{
		0x41, 0x00,
		0x2D, 0x4E,
	}
	if !bytes.Equal(encoded, expected) {
		t.Fatalf("unexpected UTF-16LE encoding: got %v want %v", encoded, expected)
	}

	if strings.Contains(string(encoded), "A中") {
		t.Fatalf("encoded bytes should not be directly readable as UTF-8 text")
	}
}
