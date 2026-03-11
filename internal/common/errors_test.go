package common

import (
	"errors"
	"testing"
)

func TestErrorCreation(t *testing.T) {
	err := NewError(ErrConfigNotFound, "test error")
	if err.Code != ErrConfigNotFound {
		t.Errorf("expected code %s, got %s", ErrConfigNotFound, err.Code)
	}
	if err.Message != "test error" {
		t.Errorf("expected message 'test error', got '%s'", err.Message)
	}
}

func TestErrorWrapping(t *testing.T) {
	original := errors.New("original error")
	wrapped := WrapError(original, ErrFileRead, "failed to read")

	if !errors.Is(wrapped, original) {
		t.Error("wrapped error should be Is the original")
	}

	var e *Error
	if !errors.As(wrapped, &e) {
		t.Error("wrapped error should be As *Error")
	}

	if e.Code != ErrFileRead {
		t.Errorf("expected code %s, got %s", ErrFileRead, e.Code)
	}
}

func TestErrorWithFields(t *testing.T) {
	err := NewError(ErrConfigInvalid, "invalid config").
		WithField("file", "config.yaml").
		WithField("line", 42)

	if err.Fields["file"] != "config.yaml" {
		t.Errorf("expected field file=config.yaml, got %v", err.Fields["file"])
	}
	if err.Fields["line"] != 42 {
		t.Errorf("expected field line=42, got %v", err.Fields["line"])
	}
}

func TestConfigNotFoundError(t *testing.T) {
	err := ConfigNotFoundError("test.yaml")

	if err.Code != ErrConfigNotFound {
		t.Errorf("expected code %s, got %s", ErrConfigNotFound, err.Code)
	}

	if err.Fields["path"] != "test.yaml" {
		t.Errorf("expected path field, got %v", err.Fields)
	}
}

func TestIsErrorCode(t *testing.T) {
	err := FileNotFoundError("test.txt")

	if !IsErrorCode(err, ErrFileNotFound) {
		t.Error("IsErrorCode should return true for matching code")
	}

	if IsErrorCode(err, ErrFileRead) {
		t.Error("IsErrorCode should return false for non-matching code")
	}
}

func TestFormatErrorForDisplay(t *testing.T) {
	err := FileNotFoundError("config.yaml")
	formatted := FormatErrorForDisplay(err)

	if formatted == "" {
		t.Error("expected non-empty formatted error")
	}
	t.Logf("Formatted: %s", formatted)
}

func TestIsRetryable(t *testing.T) {
	retryable := ModelTimeoutError("gemma", 0)
	if !IsRetryable(retryable) {
		t.Error("ModelTimeoutError should be retryable")
	}

	notRetryable := FileNotFoundError("test.txt")
	if IsRetryable(notRetryable) {
		t.Error("FileNotFoundError should not be retryable")
	}
}

func TestErrorGroup(t *testing.T) {
	group := &ErrorGroup{Op: "test"}

	group.Add(FileNotFoundError("a.txt"))
	group.Add(FileNotFoundError("b.txt"))

	if !group.HasErrors() {
		t.Error("group should have errors")
	}

	if len(group.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(group.Errors))
	}

	errStr := group.Error()
	if errStr == "" {
		t.Error("expected non-empty error string")
	}
	t.Logf("Group error: %s", errStr)
}
