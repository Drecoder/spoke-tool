package model

import (
	"time"

	"example.com/spoke-tool/api/types"
)

// ModelType represents the different SLM models we support
type ModelType string

const (
	// CodeLlamaEncoder - encoder model for code understanding
	CodeLlamaEncoder ModelType = "codellama:7b"

	// CodeLlamaDecoder - decoder model for test generation
	CodeLlamaDecoder ModelType = "codellama:7b"

	// Gemma2B - fast model for simple tasks
	Gemma2B ModelType = "gemma2:2b"
)

// String returns the string representation of the model type
func (m ModelType) String() string {
	return string(m)
}

// SLMResponse represents a response from an SLM model
// This is the output structure from model generation
type SLMResponse struct {
	// Unique identifier for this response
	ID string `json:"id"`

	// ID of the request that generated this response
	RequestID string `json:"request_id"`

	// Model that generated this response
	Model ModelType `json:"model"`

	// Programming language context
	Language types.Language `json:"language"`

	// The generated text
	Response string `json:"response"`

	// Context for continuing the conversation
	Context []int `json:"context,omitempty"`

	// Usage statistics
	TokensUsed int           `json:"tokens_used"`
	Duration   time.Duration `json:"duration_ms"`

	// Metadata
	Timestamp time.Time `json:"timestamp"`

	// Error information (if any)
	Error string `json:"error,omitempty"`

	// Whether generation is complete
	Done bool `json:"done"`
}

// ModelStatus represents the status of a model in Ollama
type ModelStatus struct {
	// Model name
	Name string `json:"name"`

	// Whether the model is available
	Available bool `json:"available"`

	// Model size (human readable)
	Size string `json:"size,omitempty"`

	// When the model was last modified
	Modified time.Time `json:"modified,omitempty"`
}

// GenerationStats represents statistics about a generation request
type GenerationStats struct {
	// Model used
	Model ModelType `json:"model"`

	// Tokens generated
	TokensGenerated int `json:"tokens_generated"`

	// Time taken
	Duration time.Duration `json:"duration_ms"`

	// Tokens per second
	TokensPerSecond float64 `json:"tokens_per_second"`
}

// CodeUnderstanding represents the result of code analysis
// This is a structured representation of code
type CodeUnderstanding struct {
	// Programming language
	Language types.Language `json:"language"`

	// Package or module name
	Package string `json:"package,omitempty"`

	// Functions found in the code
	Functions []FunctionInfo `json:"functions,omitempty"`

	// Types/Structs/Classes found
	Types []TypeInfo `json:"types,omitempty"`

	// Imports/dependencies
	Imports []string `json:"imports,omitempty"`

	// Overall complexity estimate (1-10)
	Complexity int `json:"complexity,omitempty"`

	// Brief summary of the code
	Summary string `json:"summary,omitempty"`
}

// FunctionInfo represents information about a function
type FunctionInfo struct {
	// Function name
	Name string `json:"name"`

	// Full signature
	Signature string `json:"signature"`

	// Parameters
	Parameters []ParameterInfo `json:"parameters,omitempty"`

	// Return type(s)
	Returns []string `json:"returns,omitempty"`

	// Whether this function is exported/public
	IsExported bool `json:"is_exported"`

	// Line numbers in source file
	LineStart int `json:"line_start"`
	LineEnd   int `json:"line_end"`

	// Complexity estimate (1-10)
	Complexity int `json:"complexity,omitempty"`

	// Brief description
	Description string `json:"description,omitempty"`
}

// ParameterInfo represents information about a function parameter
type ParameterInfo struct {
	// Parameter name
	Name string `json:"name"`

	// Parameter type
	Type string `json:"type"`

	// Whether the parameter is optional
	Optional bool `json:"optional,omitempty"`

	// Default value (if any)
	Default string `json:"default,omitempty"`

	// Description of the parameter
	Description string `json:"description,omitempty"`
}

// TypeInfo represents information about a type/struct/class
type TypeInfo struct {
	// Type name
	Name string `json:"name"`

	// Kind of type (struct, interface, class, etc.)
	Kind string `json:"kind"`

	// Fields/Methods
	Fields []FieldInfo `json:"fields,omitempty"`

	// Whether this type is exported/public
	IsExported bool `json:"is_exported"`

	// Line numbers
	LineStart int `json:"line_start"`
	LineEnd   int `json:"line_end"`

	// Description
	Description string `json:"description,omitempty"`
}

// FieldInfo represents information about a struct/class field
type FieldInfo struct {
	// Field name
	Name string `json:"name"`

	// Field type
	Type string `json:"type"`

	// Tags (like Go struct tags)
	Tags map[string]string `json:"tags,omitempty"`

	// Description
	Description string `json:"description,omitempty"`
}

// TestSuggestion represents a generated test
type TestSuggestion struct {
	// Programming language
	Language types.Language `json:"language"`

	// Name of the function being tested
	FunctionName string `json:"function_name"`

	// The generated test code
	TestCode string `json:"test_code"`

	// Where the test file should be saved
	TestFilePath string `json:"test_file_path"`

	// Testing framework used
	Framework string `json:"framework"`

	// Brief description of what the test covers
	Description string `json:"description,omitempty"`

	// Estimated coverage percentage
	EstimatedCoverage float64 `json:"estimated_coverage,omitempty"`
}

// DocSuggestion represents generated documentation
type DocSuggestion struct {
	// Programming language
	Language types.Language `json:"language"`

	// Name of the function being documented
	FunctionName string `json:"function_name"`

	// The generated documentation
	Content string `json:"content"`

	// Documentation format
	Format string `json:"format"`

	// Brief description
	Description string `json:"description,omitempty"`
}

// Example represents a code example extracted from tests
type Example struct {
	// Programming language
	Language types.Language `json:"language"`

	// The example code
	Code string `json:"code"`

	// Description of what the example shows
	Description string `json:"description,omitempty"`

	// Whether this came from a test file
	FromTest bool `json:"from_test"`

	// Source file where example was found
	SourceFile string `json:"source_file,omitempty"`

	// Confidence in this example (0-1)
	Confidence float64 `json:"confidence,omitempty"`
}

// FailureAnalysis represents analysis of a test failure
// This is PURELY EXPLANATORY - no fixes
type FailureAnalysis struct {
	// Test that failed
	TestName string `json:"test_name"`

	// Programming language
	Language types.Language `json:"language"`

	// Error message from the test
	ErrorMessage string `json:"error_message"`

	// Explanation of why the test failed
	Explanation string `json:"explanation"`

	// What was expected
	Expected string `json:"expected,omitempty"`

	// What actually happened
	Actual string `json:"actual,omitempty"`

	// Where in the code the issue occurred
	Location string `json:"location,omitempty"`
}

// BatchRequest represents multiple requests in one batch
type BatchRequest struct {
	// Unique batch ID
	ID string `json:"id"`

	// Individual requests
	Requests []SLMRequest `json:"requests"`

	// Whether to process in parallel
	Parallel bool `json:"parallel"`

	// Maximum concurrency
	MaxConcurrency int `json:"max_concurrency,omitempty"`
}

// BatchResponse represents responses from a batch request
type BatchResponse struct {
	// Batch ID
	BatchID string `json:"batch_id"`

	// Individual responses
	Responses []SLMResponse `json:"responses"`

	// Any errors that occurred
	Errors []string `json:"errors,omitempty"`

	// Statistics
	TotalDuration time.Duration `json:"total_duration_ms"`
	SuccessCount  int           `json:"success_count"`
	FailureCount  int           `json:"failure_count"`
}

// ModelInfo represents detailed information about a model
type ModelInfo struct {
	// Model name
	Name string `json:"name"`

	// Model type/family
	Family string `json:"family,omitempty"`

	// Parameter count (billions)
	Parameters string `json:"parameters,omitempty"`

	// Quantization level
	Quantization string `json:"quantization,omitempty"`

	// File size
	Size string `json:"size"`

	// When the model was pulled
	PulledAt time.Time `json:"pulled_at"`

	// Digest/hash of the model
	Digest string `json:"digest,omitempty"`

	// Whether this is a recommended model
	Recommended bool `json:"recommended,omitempty"`
}

// CompletionChunk represents a streaming response chunk
type CompletionChunk struct {
	// Chunk index
	Index int `json:"index"`

	// Text content
	Content string `json:"content"`

	// Whether this is the final chunk
	Final bool `json:"final"`

	// Cumulative tokens so far
	TokensSoFar int `json:"tokens_so_far"`
}

// ValidateResponse checks if a response is valid
func (r *SLMResponse) ValidateResponse() bool {
	return r.ID != "" && r.Response != ""
}

// IsError checks if the response contains an error
func (r *SLMResponse) IsError() bool {
	return r.Error != ""
}

// Success returns true if generation was successful
func (r *SLMResponse) Success() bool {
	return !r.IsError() && r.Done && r.Response != ""
}

// TokensPerSecond calculates tokens per second
func (r *SLMResponse) TokensPerSecond() float64 {
	if r.Duration == 0 || r.TokensUsed == 0 {
		return 0
	}
	return float64(r.TokensUsed) / r.Duration.Seconds()
}

// GetStats returns generation statistics
func (r *SLMResponse) GetStats() GenerationStats {
	return GenerationStats{
		Model:           r.Model,
		TokensGenerated: r.TokensUsed,
		Duration:        r.Duration,
		TokensPerSecond: r.TokensPerSecond(),
	}
}

// IsComplete checks if a response is complete
func (r *SLMResponse) IsComplete() bool {
	return r.Done && r.Error == ""
}

// HasContext checks if the response has context for continuation
func (r *SLMResponse) HasContext() bool {
	return len(r.Context) > 0
}
