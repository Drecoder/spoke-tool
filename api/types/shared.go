package types

import "time"

// EventType represents the type of event in the system
type EventType string

const (
	EventCodeChanged   EventType = "code_changed"
	EventTestsRun      EventType = "tests_run"
	EventTestsPassed   EventType = "tests_passed"
	EventTestsFailed   EventType = "tests_failed"
	EventDocsUpdated   EventType = "docs_updated"
	EventModelRequest  EventType = "model_request"
	EventModelResponse EventType = "model_response"
)

// Event represents an event in the system (for audit logging)
type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Source    string    `json:"source"` // Component that generated event
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user,omitempty"` // For audit trail
}

// Config represents the tool configuration
type Config struct {
	ProjectRoot string `json:"project_root" yaml:"project_root"`
	
	Models struct {
		Encoder   string `json:"encoder" yaml:"encoder"`
		Decoder   string `json:"decoder" yaml:"decoder"`
		Fast      string `json:"fast" yaml:"fast"`
	} `json:"models" yaml:"models"`
	
	TestSpoke struct {
		Enabled          bool    `json:"enabled" yaml:"enabled"`
		AutoRun          bool    `json:"auto_run" yaml:"auto_run"`
		CoverageThreshold float64 `json:"coverage_threshold" yaml:"coverage_threshold"`
		Frameworks       map[Language]string `json:"frameworks" yaml:"frameworks"`
	} `json:"test_spoke" yaml:"test_spoke"`
	
	ReadmeSpoke struct {
		Enabled     bool         `json:"enabled" yaml:"enabled"`
		AutoUpdate  bool         `json:"auto_update" yaml:"auto_update"`
		Sections    []DocSection `json:"sections" yaml:"sections"`
	} `json:"readme_spoke" yaml:"readme_spoke"`
	
	Squeeze struct {
		MaxCPUPercent  int `json:"max_cpu_percent" yaml:"max_cpu_percent"`
		MaxMemoryMB    int `json:"max_memory_mb" yaml:"max_memory_mb"`
		IdleThreshold  int `json:"idle_threshold_ms" yaml:"idle_threshold_ms"`
	} `json:"squeeze" yaml:"squeeze"`
	
	Audit struct {
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Path    string `json:"path" yaml:"path"`
	} `json:"audit" yaml:"audit"`
}

// Error represents a structured error
type Error struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Details   any       `json:"details,omitempty"`
	Component string    `json:"component"`
	Time      time.Time `json:"time"`
}

// Response is a standard API response
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}