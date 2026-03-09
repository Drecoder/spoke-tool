package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// DebugLevel detailed information for debugging
	DebugLevel LogLevel = iota
	// InfoLevel general operational information
	InfoLevel
	// WarnLevel warning information
	WarnLevel
	// ErrorLevel error information
	ErrorLevel
	// FatalLevel fatal error information (process will exit)
	FatalLevel
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color returns ANSI color code for the log level
func (l LogLevel) Color() string {
	switch l {
	case DebugLevel:
		return "\033[36m" // Cyan
	case InfoLevel:
		return "\033[32m" // Green
	case WarnLevel:
		return "\033[33m" // Yellow
	case ErrorLevel:
		return "\033[31m" // Red
	case FatalLevel:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// LogField represents a key-value pair for structured logging
type LogField struct {
	Key   string
	Value interface{}
}

// Fields is a map of key-value pairs for structured logging
type Fields map[string]interface{}

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	Level      LogLevel
	Output     io.Writer
	JSONFormat bool
	Color      bool
	Timestamp  bool
	Caller     bool
	CallerSkip int
	AuditFile  string
}

// Logger represents a structured logger
type Logger struct {
	config    LoggerConfig
	fields    Fields
	mu        sync.RWMutex
	auditFile *os.File
}

// LogEntry represents a single log entry
type LogEntry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Caller    string                 `json:"caller,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Stack     []string               `json:"stack,omitempty"`
	Duration  string                 `json:"duration_ms,omitempty"`
}

// Global default logger
var (
	defaultLogger *Logger
	once          sync.Once
)

// InitLogger initializes the global default logger
func InitLogger(config LoggerConfig) {
	once.Do(func() {
		if config.Output == nil {
			config.Output = os.Stdout
		}
		if config.CallerSkip == 0 {
			config.CallerSkip = 2
		}

		defaultLogger = &Logger{
			config: config,
			fields: make(Fields),
		}

		// Initialize audit file if specified
		if config.AuditFile != "" {
			if err := defaultLogger.initAuditFile(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open audit file: %v\n", err)
			}
		}
	})
}

// GetLogger returns the global default logger
func GetLogger() *Logger {
	if defaultLogger == nil {
		InitLogger(LoggerConfig{
			Level:  InfoLevel,
			Output: os.Stdout,
			Color:  true,
		})
	}
	return defaultLogger
}

// NewLogger creates a new logger instance
func NewLogger(config LoggerConfig) *Logger {
	if config.Output == nil {
		config.Output = os.Stdout
	}
	if config.CallerSkip == 0 {
		config.CallerSkip = 2
	}

	l := &Logger{
		config: config,
		fields: make(Fields),
	}

	if config.AuditFile != "" {
		if err := l.initAuditFile(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open audit file: %v\n", err)
		}
	}

	return l
}

// initAuditFile initializes the audit log file
func (l *Logger) initAuditFile() error {
	// Ensure directory exists
	dir := filepath.Dir(l.config.AuditFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create audit directory: %w", err)
	}

	// Open audit file
	file, err := os.OpenFile(l.config.AuditFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open audit file: %w", err)
	}

	l.auditFile = file
	return nil
}

// WithFields returns a new logger with the given fields
func (l *Logger) WithFields(fields Fields) *Logger {
	newLogger := l.clone()
	newLogger.mu.Lock()
	defer newLogger.mu.Unlock()

	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// WithField returns a new logger with a single field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l.WithFields(Fields{key: value})
}

// WithError returns a new logger with an error field
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}

	fields := Fields{
		"error": err.Error(),
	}

	// Add error code if it's our error type
	var e *Error
	if AsError(err, &e) {
		fields["error_code"] = string(e.Code)
		if len(e.Fields) > 0 {
			fields["error_fields"] = e.Fields
		}
	}

	return l.WithFields(fields)
}

// WithDuration returns a new logger with a duration field
func (l *Logger) WithDuration(d time.Duration) *Logger {
	return l.WithField("duration_ms", d.Milliseconds())
}

// clone creates a copy of the logger
func (l *Logger) clone() *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	newLogger := &Logger{
		config:    l.config,
		fields:    make(Fields),
		auditFile: l.auditFile,
	}

	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// log internal method to write log entries
func (l *Logger) log(level LogLevel, message string, fields Fields) {
	if level < l.config.Level {
		return
	}

	entry := l.buildEntry(level, message, fields)

	// Write to output
	l.writeOutput(entry, level)

	// Write to audit file if it's an audit-worthy event
	if level >= ErrorLevel && l.auditFile != nil {
		l.writeAudit(entry)
	}

	// Handle fatal level
	if level == FatalLevel {
		l.writeAudit(entry)
		os.Exit(1)
	}
}

// buildEntry creates a log entry
func (l *Logger) buildEntry(level LogLevel, message string, fields Fields) *LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	entry := &LogEntry{
		Level:     level.String(),
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	// Add caller info
	if l.config.Caller {
		if _, file, line, ok := runtime.Caller(l.config.CallerSkip); ok {
			entry.Caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
	}

	// Merge fields
	allFields := make(Fields)
	for k, v := range l.fields {
		allFields[k] = v
	}
	for k, v := range fields {
		allFields[k] = v
	}
	if len(allFields) > 0 {
		entry.Fields = allFields
	}

	return entry
}

// writeOutput writes the log entry to the configured output
func (l *Logger) writeOutput(entry *LogEntry, level LogLevel) {
	if l.config.JSONFormat {
		json.NewEncoder(l.config.Output).Encode(entry)
		return
	}

	// Human-readable format
	var sb strings.Builder

	// Timestamp
	if l.config.Timestamp {
		t, _ := time.Parse(time.RFC3339Nano, entry.Timestamp)
		sb.WriteString(t.Format("2006-01-02 15:04:05 "))
	}

	// Level with color
	if l.config.Color {
		sb.WriteString(level.Color())
		sb.WriteString(fmt.Sprintf("[%s]", level.String()))
		sb.WriteString("\033[0m ")
	} else {
		sb.WriteString(fmt.Sprintf("[%s] ", level.String()))
	}

	// Message
	sb.WriteString(entry.Message)

	// Caller
	if entry.Caller != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", entry.Caller))
	}

	// Fields
	if len(entry.Fields) > 0 {
		sb.WriteString(" {")
		first := true
		for k, v := range entry.Fields {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s=%v", k, v))
			first = false
		}
		sb.WriteString("}")
	}

	sb.WriteString("\n")
	fmt.Fprint(l.config.Output, sb.String())
}

// writeAudit writes to the audit file
func (l *Logger) writeAudit(entry *LogEntry) {
	if l.auditFile == nil {
		return
	}

	entry.Timestamp = time.Now().Format(time.RFC3339)
	json.NewEncoder(l.auditFile).Encode(entry)
	l.auditFile.Sync()
}

// Debug logs a debug message
func (l *Logger) Debug(args ...interface{}) {
	l.log(DebugLevel, fmt.Sprint(args...), nil)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...), nil)
}

// Info logs an info message
func (l *Logger) Info(args ...interface{}) {
	l.log(InfoLevel, fmt.Sprint(args...), nil)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...), nil)
}

// Warn logs a warning message
func (l *Logger) Warn(args ...interface{}) {
	l.log(WarnLevel, fmt.Sprint(args...), nil)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...), nil)
}

// Error logs an error message
func (l *Logger) Error(args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprint(args...), nil)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...), nil)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(args ...interface{}) {
	l.log(FatalLevel, fmt.Sprint(args...), nil)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FatalLevel, fmt.Sprintf(format, args...), nil)
}

// LogWithFields logs a message with fields at the specified level
func (l *Logger) LogWithFields(level LogLevel, message string, fields Fields) {
	l.log(level, message, fields)
}

// StartOperation logs the start of an operation and returns a function to log completion
func (l *Logger) StartOperation(name string, fields Fields) func(error) {
	start := time.Now()
	l.Info("Starting "+name, fields)

	return func(err error) {
		duration := time.Since(start)
		if err != nil {
			l.WithError(err).WithDuration(duration).Errorf("%s failed", name)
		} else {
			l.WithDuration(duration).Infof("%s completed", name)
		}
	}
}

// Audit logs an audit event (always written to audit file)
func (l *Logger) Audit(action string, user string, fields Fields) {
	if fields == nil {
		fields = make(Fields)
	}
	fields["action"] = action
	fields["user"] = user

	entry := l.buildEntry(InfoLevel, "AUDIT: "+action, fields)

	if l.auditFile != nil {
		l.writeAudit(entry)
	}

	// Also log to console if debug level
	if l.config.Level <= DebugLevel {
		l.writeOutput(entry, InfoLevel)
	}
}

// Close closes the logger and audit file
func (l *Logger) Close() error {
	if l.auditFile != nil {
		return l.auditFile.Close()
	}
	return nil
}

// Global convenience functions

// Debug logs a debug message using the default logger
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf logs a formatted debug message using the default logger
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info logs an info message using the default logger
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof logs a formatted info message using the default logger
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn logs a warning message using the default logger
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf logs a formatted warning message using the default logger
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error logs an error message using the default logger
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf logs a formatted error message using the default logger
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal logs a fatal message using the default logger
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf logs a formatted fatal message using the default logger
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// WithFields returns a new logger with fields using the default logger
func WithFields(fields Fields) *Logger {
	return GetLogger().WithFields(fields)
}

// WithField returns a new logger with a field using the default logger
func WithField(key string, value interface{}) *Logger {
	return GetLogger().WithField(key, value)
}

// WithError returns a new logger with an error field using the default logger
func WithError(err error) *Logger {
	return GetLogger().WithError(err)
}

// StartOperation logs operation start using the default logger
func StartOperation(name string, fields Fields) func(error) {
	return GetLogger().StartOperation(name, fields)
}

// Audit logs an audit event using the default logger
func Audit(action string, user string, fields Fields) {
	GetLogger().Audit(action, user, fields)
}

// InitFromConfig initializes the logger from a config struct
func InitFromConfig(cfg interface {
	GetLogLevel() string
	GetLogJSON() bool
	GetLogColor() bool
	GetAuditFile() string
}) {
	level := InfoLevel
	switch strings.ToLower(cfg.GetLogLevel()) {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn", "warning":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	}

	InitLogger(LoggerConfig{
		Level:      level,
		Output:     os.Stdout,
		JSONFormat: cfg.GetLogJSON(),
		Color:      cfg.GetLogColor(),
		Timestamp:  true,
		Caller:     level == DebugLevel,
		AuditFile:  cfg.GetAuditFile(),
	})
}
