// Package logger provides logging utilities for the TrackTaste application.
// Log format: [LEVEL] YYYY-MM-DD HH:mm:ss [Feature] Message
package logger

import (
	"fmt"
	"time"
)

// Level represents log severity levels.
type Level string

// Log level constants in order of decreasing severity.
const (
	// LevelFatal indicates a fatal error that causes application termination.
	LevelFatal Level = "FATAL"
	// LevelError indicates an error that should be investigated.
	LevelError Level = "ERROR"
	// LevelWarning indicates a warning that may require attention.
	LevelWarning Level = "WARNING"
	// LevelInfo indicates general informational messages.
	LevelInfo Level = "INFO"
	// LevelDebug indicates debug-level messages for development.
	LevelDebug Level = "DEBUG"
)

// Log outputs a log message with the specified level and feature name.
// Format: [LEVEL] YYYY-MM-DD HH:mm:ss [feature] message
func Log(level Level, feature string, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s [%s] %s\n", level, timestamp, feature, message)
}

// Fatal outputs a fatal log message and should be followed by application termination.
func Fatal(feature string, message string) {
	Log(LevelFatal, feature, message)
}

// Error outputs an error log message for errors that need investigation.
func Error(feature string, message string) {
	Log(LevelError, feature, message)
}

// Warning outputs a warning log message for conditions that may require attention.
func Warning(feature string, message string) {
	Log(LevelWarning, feature, message)
}

// Info outputs an info log message for general operational information.
func Info(feature string, message string) {
	Log(LevelInfo, feature, message)
}

// Debug outputs a debug log message for development and troubleshooting.
func Debug(feature string, message string) {
	Log(LevelDebug, feature, message)
}
