package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LoggerConfig{
		Level:  DebugLevel,
		Output: &buf,
		Color:  false,
		Caller: true,
	})

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	if !contains(output, "DEBUG") {
		t.Error("expected DEBUG message")
	}
	if !contains(output, "INFO") {
		t.Error("expected INFO message")
	}
	if !contains(output, "WARN") {
		t.Error("expected WARN message")
	}
	if !contains(output, "ERROR") {
		t.Error("expected ERROR message")
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LoggerConfig{
		Level:  WarnLevel,
		Output: &buf,
		Color:  false,
	})

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	if contains(output, "DEBUG") {
		t.Error("DEBUG should be filtered out")
	}
	if contains(output, "INFO") {
		t.Error("INFO should be filtered out")
	}
	if !contains(output, "WARN") {
		t.Error("expected WARN message")
	}
	if !contains(output, "ERROR") {
		t.Error("expected ERROR message")
	}
}

func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LoggerConfig{
		Level:  InfoLevel,
		Output: &buf,
		Color:  false,
	})

	logger.WithField("key1", "value1").
		WithField("key2", 42).
		Info("test message")

	output := buf.String()
	if !contains(output, "key1=value1") {
		t.Error("expected key1=value1 in output")
	}
	if !contains(output, "key2=42") {
		t.Error("expected key2=42 in output")
	}
}

func TestLoggerJSONFormat(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LoggerConfig{
		Level:      InfoLevel,
		Output:     &buf,
		JSONFormat: true,
		Caller:     true,
	})

	logger.WithFields(Fields{
		"user": "test",
		"id":   123,
	}).Info("json test")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["level"] != "INFO" {
		t.Errorf("expected level INFO, got %v", entry["level"])
	}
	if entry["message"] != "json test" {
		t.Errorf("expected message 'json test', got %v", entry["message"])
	}

	fields, ok := entry["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("expected fields object")
	}
	if fields["user"] != "test" {
		t.Errorf("expected user=test, got %v", fields["user"])
	}
	if fields["id"] != float64(123) { // JSON numbers are float64
		t.Errorf("expected id=123, got %v", fields["id"])
	}
}

func TestLoggerWithError(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LoggerConfig{
		Level:  InfoLevel,
		Output: &buf,
		Color:  false,
	})

	err := errors.New("something went wrong")
	logger.WithError(err).Error("operation failed")

	output := buf.String()
	if !contains(output, "error=something went wrong") {
		t.Error("expected error field in output")
	}
}

func TestLoggerWithDuration(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LoggerConfig{
		Level:  InfoLevel,
		Output: &buf,
		Color:  false,
	})

	logger.WithDuration(1500 * time.Millisecond).Info("slow operation")

	output := buf.String()
	if !contains(output, "duration_ms=1500") {
		t.Error("expected duration field in output")
	}
}

func TestStartOperation(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(LoggerConfig{
		Level:  InfoLevel,
		Output: &buf,
		Color:  false,
	})

	complete := logger.StartOperation("test", Fields{"op": "test"})
	time.Sleep(10 * time.Millisecond)
	complete(nil)

	output := buf.String()
	if !contains(output, "Starting test") {
		t.Error("expected start message")
	}
	if !contains(output, "completed") {
		t.Error("expected completion message")
	}
}

func TestGlobalLogger(t *testing.T) {
	// Reset default logger
	defaultLogger = nil

	var buf bytes.Buffer
	InitLogger(LoggerConfig{
		Level:  InfoLevel,
		Output: &buf,
		Color:  false,
	})

	Info("global test")

	output := buf.String()
	if !contains(output, "global test") {
		t.Error("expected global test message")
	}
}

func TestAuditLogging(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "audit-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	var buf bytes.Buffer
	logger := NewLogger(LoggerConfig{
		Level:     DebugLevel,
		Output:    &buf,
		Color:     false,
		AuditFile: tmpfile.Name(),
	})

	logger.Audit("test_action", "test_user", Fields{
		"resource": "test",
	})

	// Check console output (should only appear in debug)
	consoleOutput := buf.String()
	if !contains(consoleOutput, "AUDIT") {
		t.Error("expected AUDIT in console output for debug level")
	}

	// Check audit file
	auditContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	var auditEntry map[string]interface{}
	if err := json.Unmarshal(auditContent, &auditEntry); err != nil {
		t.Fatalf("failed to parse audit JSON: %v", err)
	}

	if auditEntry["level"] != "INFO" {
		t.Errorf("expected level INFO, got %v", auditEntry["level"])
	}

	fields, ok := auditEntry["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("expected fields object")
	}
	if fields["action"] != "test_action" {
		t.Errorf("expected action=test_action, got %v", fields["action"])
	}
	if fields["user"] != "test_user" {
		t.Errorf("expected user=test_user, got %v", fields["user"])
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
