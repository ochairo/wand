// Package logging provides structured logging capabilities.
package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel defines the severity of a log message
type LogLevel int

const (
	// DEBUG represents debug-level logging.
	DEBUG LogLevel = iota
	// INFO represents info-level logging.
	INFO
	// WARN represents warning-level logging.
	WARN
	// ERROR represents error-level logging.
	ERROR
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging with JSON output and file support
type Logger struct {
	level  LogLevel
	output io.Writer
	file   *os.File
}

// LogEntry represents a single structured log entry
type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

// NewLogger creates a new logger with the given output and log level
func NewLogger(level LogLevel, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}
	return &Logger{
		level:  level,
		output: output,
	}
}

// NewFileLogger creates a logger that writes to a file and stdout
func NewFileLogger(level LogLevel, filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600) //nolint:gosec // G304: filePath from application config
	if err != nil {
		return nil, err
	}

	// Write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)

	return &Logger{
		level:  level,
		output: multiWriter,
		file:   file,
	}, nil
}

// log writes a log entry at the specified level
func (l *Logger) log(level LogLevel, message string, data interface{}) {
	if level < l.level {
		return // Skip messages below configured log level
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level.String(),
		Message:   message,
		Data:      data,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		_, _ = fmt.Fprintf(l.output, "ERROR: Failed to marshal log entry: %v\n", err)
		return
	}

	_, _ = fmt.Fprintf(l.output, "%s\n", string(jsonBytes))
}

// Debug logs a debug-level message
func (l *Logger) Debug(message string, data interface{}) {
	l.log(DEBUG, message, data)
}

// Info logs an info-level message
func (l *Logger) Info(message string, data interface{}) {
	l.log(INFO, message, data)
}

// Warn logs a warning-level message
func (l *Logger) Warn(message string, data interface{}) {
	l.log(WARN, message, data)
}

// Error logs an error-level message
func (l *Logger) Error(message string, data interface{}) {
	l.log(ERROR, message, data)
}

// Close closes the logger and any associated file
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// SetLevel changes the log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() LogLevel {
	return l.level
}
