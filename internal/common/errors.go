package common

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ErrorCode represents a specific error type in the system
type ErrorCode string

const (
	// Configuration errors (1xxx)
	ErrConfigNotFound     ErrorCode = "CONFIG_1001"
	ErrConfigInvalid      ErrorCode = "CONFIG_1002"
	ErrConfigParse        ErrorCode = "CONFIG_1003"
	ErrConfigMissingField ErrorCode = "CONFIG_1004"

	// File system errors (2xxx)
	ErrFileNotFound    ErrorCode = "FS_2001"
	ErrFilePermission  ErrorCode = "FS_2002"
	ErrFileRead        ErrorCode = "FS_2003"
	ErrFileWrite       ErrorCode = "FS_2004"
	ErrDirectoryCreate ErrorCode = "FS_2005"

	// Model/SLM errors (3xxx)
	ErrModelConnection  ErrorCode = "MODEL_3001"
	ErrModelNotFound    ErrorCode = "MODEL_3002"
	ErrModelTimeout     ErrorCode = "MODEL_3003"
	ErrModelResponse    ErrorCode = "MODEL_3004"
	ErrModelPull        ErrorCode = "MODEL_3005"
	ErrModelUnavailable ErrorCode = "MODEL_3006"

	// Analysis errors (4xxx)
	ErrParseError          ErrorCode = "ANALYSIS_4001"
	ErrUnsupportedLanguage ErrorCode = "ANALYSIS_4002"
	ErrNoFunctionsFound    ErrorCode = "ANALYSIS_4003"
	ErrComplexityTooHigh   ErrorCode = "ANALYSIS_4004"

	// Test generation errors (5xxx)
	ErrTestGeneration ErrorCode = "TEST_5001"
	ErrTestRun        ErrorCode = "TEST_5002"
	ErrTestTimeout    ErrorCode = "TEST_5003"
	ErrTestCoverage   ErrorCode = "TEST_5004"
	ErrTestParse      ErrorCode = "TEST_5005"

	// Documentation errors (6xxx)
	ErrDocGeneration ErrorCode = "DOC_6001"
	ErrDocParse      ErrorCode = "DOC_6002"
	ErrDocMerge      ErrorCode = "DOC_6003"
	ErrDocValidation ErrorCode = "DOC_6004"

	// General errors (9xxx)
	ErrInternal        ErrorCode = "GENERAL_9001"
	ErrNotImplemented  ErrorCode = "GENERAL_9002"
	ErrInvalidArgument ErrorCode = "GENERAL_9003"
	ErrTimeout         ErrorCode = "GENERAL_9004"
	ErrCancelled       ErrorCode = "GENERAL_9005"
)

// Error represents a structured error with context
type Error struct {
	Code      ErrorCode      `json:"code"`
	Message   string         `json:"message"`
	Op        string         `json:"operation,omitempty"`
	Err       error          `json:"cause,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
	Stack     []string       `json:"-"`
	Timestamp time.Time      `json:"timestamp"`
	Severity  string         `json:"severity"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements the unwrap interface
func (e *Error) Unwrap() error {
	return e.Err
}

// Is checks if the error matches a target
func (e *Error) Is(target error) bool {
	var t *Error
	if !errors.As(target, &t) {
		return false
	}
	return e.Code == t.Code
}

// WithField adds a field to the error context
func (e *Error) WithField(key string, value any) *Error {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[key] = value
	return e
}

// WithFields adds multiple fields to the error context
func (e *Error) WithFields(fields map[string]any) *Error {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

// WithOp adds operation context
func (e *Error) WithOp(op string) *Error {
	e.Op = op
	return e
}

// WithSeverity sets error severity
func (e *Error) WithSeverity(severity string) *Error {
	e.Severity = severity
	return e
}

// Error constructors

// NewError creates a new base error
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Severity:  "error",
		Stack:     captureStack(2),
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, code ErrorCode, message string) *Error {
	if err == nil {
		return nil
	}

	// If it's already our error type, preserve it
	var e *Error
	if errors.As(err, &e) {
		return &Error{
			Code:      code,
			Message:   message,
			Op:        e.Op,
			Err:       e,
			Fields:    e.Fields,
			Timestamp: time.Now(),
			Severity:  e.Severity,
			Stack:     captureStack(2),
		}
	}

	return &Error{
		Code:      code,
		Message:   message,
		Err:       err,
		Timestamp: time.Now(),
		Severity:  "error",
		Stack:     captureStack(2),
	}
}

// Convenience constructors for common error types

// ConfigError creates a new configuration error
func ConfigError(code ErrorCode, message string, args ...any) *Error {
	return NewError(code, fmt.Sprintf(message, args...)).WithOp("config")
}

// ConfigNotFoundError creates a "config not found" error
func ConfigNotFoundError(path string) *Error {
	return ConfigError(ErrConfigNotFound, "configuration file not found: %s", path).
		WithField("path", path).
		WithSeverity("fatal")
}

// ConfigInvalidError creates an "invalid config" error
func ConfigInvalidError(reason string, fields map[string]any) *Error {
	return ConfigError(ErrConfigInvalid, "invalid configuration: %s", reason).
		WithFields(fields).
		WithSeverity("fatal")
}

// FileSystemError creates a new filesystem error
func FileSystemError(code ErrorCode, message string, path string, err error) *Error {
	return WrapError(err, code, message).
		WithOp("filesystem").
		WithField("path", path)
}

// FileNotFoundError creates a "file not found" error
func FileNotFoundError(path string) *Error {
	return FileSystemError(ErrFileNotFound, "file not found", path, nil).
		WithSeverity("error")
}

// FileReadError creates a "file read" error
func FileReadError(path string, err error) *Error {
	return FileSystemError(ErrFileRead, "failed to read file", path, err).
		WithSeverity("error")
}

// FileWriteError creates a "file write" error
func FileWriteError(path string, err error) *Error {
	return FileSystemError(ErrFileWrite, "failed to write file", path, err).
		WithSeverity("error")
}

// ModelError creates a new model error
func ModelError(code ErrorCode, message string, modelName string, err error) *Error {
	return WrapError(err, code, message).
		WithOp("model").
		WithField("model", modelName)
}

// ModelConnectionError creates a "model connection" error
func ModelConnectionError(host string, err error) *Error {
	return ModelError(ErrModelConnection, "failed to connect to model server", host, err).
		WithField("host", host).
		WithSeverity("fatal")
}

// ModelNotFoundError creates a "model not found" error
func ModelNotFoundError(modelName string) *Error {
	return ModelError(ErrModelNotFound, "model not found", modelName, nil).
		WithSeverity("warning")
}

// ModelTimeoutError creates a "model timeout" error
func ModelTimeoutError(modelName string, timeout time.Duration) *Error {
	return ModelError(ErrModelTimeout, "model request timed out", modelName, nil).
		WithField("timeout", timeout.String()).
		WithSeverity("error")
}

// AnalysisError creates a new analysis error
func AnalysisError(code ErrorCode, message string, language string, file string) *Error {
	return NewError(code, message).
		WithOp("analysis").
		WithField("language", language).
		WithField("file", file)
}

// ParseError creates a "parse error"
func ParseError(language string, file string, err error) *Error {
	return AnalysisError(ErrParseError, "failed to parse code", language, file).
		WithField("error", err.Error()).
		WithSeverity("error")
}

// UnsupportedLanguageError creates an "unsupported language" error
func UnsupportedLanguageError(language string) *Error {
	return AnalysisError(ErrUnsupportedLanguage, "unsupported language", language, "").
		WithSeverity("error")
}

// TestError creates a new test error
func TestError(code ErrorCode, message string, testFile string) *Error {
	return NewError(code, message).
		WithOp("test").
		WithField("test_file", testFile)
}

// TestGenerationError creates a "test generation" error
func TestGenerationError(function string, err error) *Error {
	return WrapError(err, ErrTestGeneration, "failed to generate test").
		WithOp("test").
		WithField("function", function).
		WithSeverity("error")
}

// TestRunError creates a "test run" error
func TestRunError(testFile string, err error) *Error {
	return TestError(ErrTestRun, "test execution failed", testFile).
		WithField("error", err.Error()).
		WithSeverity("warning")
}

// TestFailure represents a test failure (not an error, but a test that failed)
func TestFailure(testName string, message string) *Error {
	return &Error{
		Code:      ErrTestRun,
		Message:   fmt.Sprintf("test failed: %s - %s", testName, message),
		Timestamp: time.Now(),
		Severity:  "test_failure",
		Fields: map[string]any{
			"test_name": testName,
		},
	}
}

// DocError creates a new documentation error
func DocError(code ErrorCode, message string, section string) *Error {
	return NewError(code, message).
		WithOp("documentation").
		WithField("section", section)
}

// DocGenerationError creates a "doc generation" error
func DocGenerationError(function string, err error) *Error {
	return WrapError(err, ErrDocGeneration, "failed to generate documentation").
		WithOp("documentation").
		WithField("function", function).
		WithSeverity("error")
}

// InternalError creates a new internal error
func InternalError(message string, err error) *Error {
	return WrapError(err, ErrInternal, message).
		WithOp("internal").
		WithSeverity("fatal")
}

// NotImplementedError creates a "not implemented" error
func NotImplementedError(feature string) *Error {
	return NewError(ErrNotImplemented, fmt.Sprintf("feature not implemented: %s", feature)).
		WithOp("internal").
		WithField("feature", feature).
		WithSeverity("warning")
}

// InvalidArgumentError creates an "invalid argument" error
func InvalidArgumentError(arg string, reason string) *Error {
	return NewError(ErrInvalidArgument, fmt.Sprintf("invalid argument %s: %s", arg, reason)).
		WithOp("validation").
		WithField("argument", arg).
		WithSeverity("error")
}

// Helper functions

// captureStack captures the current stack trace
func captureStack(skip int) []string {
	var stack []string
	for i := skip; i < skip+10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}
	return stack
}

// IsErrorCode checks if an error has a specific code
func IsErrorCode(err error, code ErrorCode) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// GetErrorCode returns the error code if available
func GetErrorCode(err error) ErrorCode {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

// GetErrorFields returns the error fields if available
func GetErrorFields(err error) map[string]any {
	var e *Error
	if errors.As(err, &e) {
		return e.Fields
	}
	return nil
}

// FormatErrorForDisplay formats an error for user display
func FormatErrorForDisplay(err error) string {
	var e *Error
	if !errors.As(err, &e) {
		return err.Error()
	}

	var sb strings.Builder

	// Add emoji based on severity
	switch e.Severity {
	case "fatal":
		sb.WriteString("💥 ")
	case "error":
		sb.WriteString("❌ ")
	case "warning":
		sb.WriteString("⚠️ ")
	case "test_failure":
		sb.WriteString("🧪 ")
	default:
		sb.WriteString("• ")
	}

	// Add message
	sb.WriteString(e.Message)

	// Add fields if any
	if len(e.Fields) > 0 {
		sb.WriteString(" (")
		first := true
		for k, v := range e.Fields {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s: %v", k, v))
			first = false
		}
		sb.WriteString(")")
	}

	return sb.String()
}

// IsRetryable determines if an error is retryable
func IsRetryable(err error) bool {
	var e *Error
	if !errors.As(err, &e) {
		return false
	}

	switch e.Code {
	case ErrModelConnection, ErrModelTimeout, ErrModelUnavailable:
		return true
	case ErrFilePermission, ErrFileWrite:
		return false
	default:
		return false
	}
}

// ErrorGroup represents a collection of errors
type ErrorGroup struct {
	Errors []error
	Op     string
}

// Add adds an error to the group
func (eg *ErrorGroup) Add(err error) {
	if err != nil {
		eg.Errors = append(eg.Errors, err)
	}
}

// HasErrors returns true if the group has any errors
func (eg *ErrorGroup) HasErrors() bool {
	return len(eg.Errors) > 0
}

// Error implements the error interface
func (eg *ErrorGroup) Error() string {
	if len(eg.Errors) == 0 {
		return ""
	}
	if len(eg.Errors) == 1 {
		return eg.Errors[0].Error()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d errors occurred", len(eg.Errors)))
	for i, err := range eg.Errors {
		sb.WriteString(fmt.Sprintf("\n  %d: %v", i+1, err))
	}
	return sb.String()
}
