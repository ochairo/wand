package logging

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestLoggerDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(DEBUG, buf)

	logger.Debug("test debug", map[string]string{"key": "value"})

	var entry LogEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "DEBUG" {
		t.Errorf("Expected level DEBUG, got %s", entry.Level)
	}
	if entry.Message != "test debug" {
		t.Errorf("Expected message 'test debug', got %s", entry.Message)
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(INFO, buf)

	logger.Debug("debug message", nil)
	logger.Info("info message", nil)

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != 1 {
		t.Errorf("Expected 1 log line (DEBUG filtered out), got %d", len(lines))
	}

	var entry LogEntry
	err := json.Unmarshal(lines[0], &entry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}
}

func TestLoggerError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(ERROR, buf)

	logger.Error("critical error", map[string]interface{}{
		"code":    "INSTALL_FAILED",
		"package": "nano",
	})

	var entry LogEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	if err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "ERROR" {
		t.Errorf("Expected level ERROR, got %s", entry.Level)
	}
}

func TestSetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(INFO, buf)

	if logger.GetLevel() != INFO {
		t.Errorf("Expected level INFO, got %v", logger.GetLevel())
	}

	logger.SetLevel(DEBUG)
	if logger.GetLevel() != DEBUG {
		t.Errorf("Expected level DEBUG, got %v", logger.GetLevel())
	}
}
