// Package logging provides a centralized logging system for Askeladden Discord bot.
// This package standardizes log formatting, levels, and output across the application.
package logging

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the log level
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

// Logger represents the centralized logger for the application
type Logger struct {
	level     LogLevel
	component string
	stdLogger *log.Logger
}

// Global logger instance
var globalLogger *Logger

// Initialize sets up the global logger with the specified minimum log level
func Initialize(level LogLevel) {
	globalLogger = &Logger{
		level:     level,
		component: "APP",
		stdLogger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// SetLevel updates the minimum log level
func SetLevel(level LogLevel) {
	if globalLogger != nil {
		globalLogger.level = level
	}
}

// GetLogger returns a logger instance for a specific component
func GetLogger(component string) *Logger {
	if globalLogger == nil {
		Initialize(INFO) // Default to INFO level
	}

	return &Logger{
		level:     globalLogger.level, // Reference to global level
		component: strings.ToUpper(component),
		stdLogger: globalLogger.stdLogger,
	}
}

// shouldLog determines if a message should be logged based on the current level
func (l *Logger) shouldLog(level LogLevel) bool {
	// Always check against the global logger's current level
	if globalLogger != nil {
		return level >= globalLogger.level
	}
	return level >= l.level
}

// formatMessage creates a consistently formatted log message
func (l *Logger) formatMessage(level LogLevel, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	return fmt.Sprintf("[%s] [%s] %s", level, l.component, message)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.shouldLog(DEBUG) {
		l.stdLogger.Println(l.formatMessage(DEBUG, format, args...))
	}
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.shouldLog(INFO) {
		l.stdLogger.Println(l.formatMessage(INFO, format, args...))
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.shouldLog(WARN) {
		l.stdLogger.Println(l.formatMessage(WARN, format, args...))
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.shouldLog(ERROR) {
		l.stdLogger.Println(l.formatMessage(ERROR, format, args...))
	}
}

// Fatal logs an error message and exits the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	message := l.formatMessage(ERROR, format, args...)
	l.stdLogger.Fatalf("%s", message)
}

// Global convenience functions that use the default logger

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	GetLogger("APP").Debug(format, args...)
}

// Info logs an informational message using the default logger
func Info(format string, args ...interface{}) {
	GetLogger("APP").Info(format, args...)
}

// Warn logs a warning message using the default logger
func Warn(format string, args ...interface{}) {
	GetLogger("APP").Warn(format, args...)
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	GetLogger("APP").Error(format, args...)
}

// Fatal logs an error message and exits the program using the default logger
func Fatal(format string, args ...interface{}) {
	GetLogger("APP").Fatal(format, args...)
}

// ParseLogLevel converts a string to a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}
