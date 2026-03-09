package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	ollama "github.com/ollama/ollama/client"
	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/internal/common"
)

// ModelType represents the different SLM models we support
type ModelType string

const (
	// Encoder models for code understanding
	CodeBERT ModelType = "codebert"

	// Decoder models
	Gemma2B    ModelType = "gemma2:2b"         // Fast, lightweight generation for docs
	DeepSeek7B ModelType = "deepseek-coder:7b" // Complex reasoning for tests
)

// Client handles communication with Ollama SLMs
// This client ONLY handles model communication - no business logic
type Client struct {
	ollamaClient ollama.Client
	modelMap     map[ModelType]string
	timeout      time.Duration
	logger       *common.Logger
}

// ClientConfig holds configuration for the model client
type ClientConfig struct {
	OllamaHost string
	Timeout    time.Duration
	Models     map[ModelType]string
}

// Request represents a request to an SLM
type Request struct {
	ID          string                 `json:"id"`
	Model       ModelType              `json:"model"`
	Language    types.Language         `json:"language"`
	Prompt      string                 `json:"prompt"`
	System      string                 `json:"system,omitempty"`
	Context     []int                  `json:"context,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Temperature float32                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
}

// Response represents a response from an SLM
type Response struct {
	ID         string         `json:"id"`
	RequestID  string         `json:"request_id"`
	Model      ModelType      `json:"model"`
	Language   types.Language `json:"language"`
	Response   string         `json:"response"`
	TokensUsed int            `json:"tokens_used"`
	Duration   time.Duration  `json:"duration_ms"`
	Timestamp  time.Time      `json:"timestamp"`
	Error      string         `json:"error,omitempty"`
	Done       bool           `json:"done"`
}

// ModelStatus represents the status of a model
type ModelStatus struct {
	Name      string    `json:"name"`
	Available bool      `json:"available"`
	Size      string    `json:"size,omitempty"`
	Modified  time.Time `json:"modified,omitempty"`
}

// NewClient creates a new SLM client
func NewClient(config ClientConfig) (*Client, error) {
	if config.OllamaHost == "" {
		config.OllamaHost = "http://localhost:11434"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Create Ollama client
	client, err := ollama.NewClient(config.OllamaHost)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	// Default model mapping
	modelMap := config.Models
	if modelMap == nil {
		modelMap = map[ModelType]string{
			CodeBERT:   "codebert",
			Gemma2B:    "gemma2:2b",
			DeepSeek7B: "deepseek-coder:7b",
		}
	}

	return &Client{
		ollamaClient: client,
		modelMap:     modelMap,
		timeout:      config.Timeout,
		logger:       common.GetLogger().WithField("component", "model-client"),
	}, nil
}

// Generate sends a prompt to the specified model and returns the response
// This is the core method - all other methods build on this
func (c *Client) Generate(ctx context.Context, req Request) (*Response, error) {
	// Set defaults
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}
	if req.Temperature == 0 {
		req.Temperature = 0.7
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 2048
	}
	if req.Language == "" {
		req.Language = types.Go
	}

	// Get model name from map
	modelName, ok := c.modelMap[req.Model]
	if !ok {
		return nil, fmt.Errorf("unknown model type: %s", req.Model)
	}

	c.logger.Debug("Generating with model",
		"model", req.Model,
		"language", req.Language,
		"prompt_length", len(req.Prompt))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Prepare request
	start := time.Now()

	// Build options
	options := map[string]interface{}{
		"temperature": req.Temperature,
		"num_predict": req.MaxTokens,
	}
	// Merge custom options
	for k, v := range req.Options {
		options[k] = v
	}

	// Generate
	resp := &Response{
		ID:        uuid.New().String(),
		RequestID: req.ID,
		Model:     req.Model,
		Language:  req.Language,
		Timestamp: time.Now(),
	}

	var fullResponse strings.Builder
	var tokens int

	err := c.ollamaClient.Generate(ctx, &ollama.GenerateRequest{
		Model:   modelName,
		Prompt:  req.Prompt,
		System:  req.System,
		Context: req.Context,
		Options: options,
	}, func(response ollama.GenerateResponse) error {
		fullResponse.WriteString(response.Response)
		tokens += len(strings.Fields(response.Response))
		resp.Done = response.Done
		return nil
	})

	resp.Duration = time.Since(start)
	resp.TokensUsed = tokens

	if err != nil {
		resp.Error = err.Error()
		c.logger.Error("Generation failed", "error", err, "model", req.Model)
		return resp, fmt.Errorf("generation failed: %w", err)
	}

	resp.Response = fullResponse.String()

	c.logger.Debug("Generation complete",
		"model", req.Model,
		"tokens", tokens,
		"duration", resp.Duration)

	return resp, nil
}

// GenerateWithTemplate uses a template to create the prompt
// Convenience method for common patterns
func (c *Client) GenerateWithTemplate(ctx context.Context, model ModelType, template string, args ...interface{}) (*Response, error) {
	prompt := fmt.Sprintf(template, args...)

	req := Request{
		Model:  model,
		Prompt: prompt,
	}

	return c.Generate(ctx, req)
}

// CheckModels verifies that required models are available
// Returns status map and any missing models
func (c *Client) CheckModels(ctx context.Context) (map[ModelType]bool, error) {
	status := make(map[ModelType]bool)

	// List available models
	models, err := c.ollamaClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	// Create a set of available models
	available := make(map[string]bool)
	for _, m := range models.Models {
		available[m.Name] = true
		c.logger.Debug("Found available model", "model", m.Name)
	}

	// Check each required model
	missing := []string{}
	for modelType, modelName := range c.modelMap {
		status[modelType] = available[modelName]
		if !available[modelName] {
			missing = append(missing, string(modelName))
		}
	}

	if len(missing) > 0 {
		c.logger.Warn("Some models are not available", "missing", missing)
	}

	return status, nil
}

// PullModel pulls a model from Ollama if not already present
// This is an optional convenience method
func (c *Client) PullModel(ctx context.Context, modelType ModelType) error {
	modelName, ok := c.modelMap[modelType]
	if !ok {
		return fmt.Errorf("unknown model type: %s", modelType)
	}

	// Check if already present
	status, err := c.CheckModels(ctx)
	if err != nil {
		return err
	}

	if status[modelType] {
		c.logger.Info("Model already available", "model", modelName)
		return nil // Already present
	}

	c.logger.Info("Pulling model", "model", modelName, "this may take a while")

	// Pull the model
	err = c.ollamaClient.Pull(ctx, modelName, nil)
	if err != nil {
		return fmt.Errorf("failed to pull model %s: %w", modelName, err)
	}

	c.logger.Info("Model pulled successfully", "model", modelName)
	return nil
}

// ListModels returns all available models
func (c *Client) ListModels(ctx context.Context) ([]ModelStatus, error) {
	models, err := c.ollamaClient.List(ctx)
	if err != nil {
		return nil, err
	}

	var status []ModelStatus
	for _, m := range models.Models {
		status = append(status, ModelStatus{
			Name:      m.Name,
			Available: true,
			Size:      formatSize(m.Size),
			Modified:  m.ModifiedAt,
		})
	}

	return status, nil
}

// GetModelInfo returns information about a specific model
func (c *Client) GetModelInfo(ctx context.Context, modelType ModelType) (*ModelStatus, error) {
	models, err := c.ListModels(ctx)
	if err != nil {
		return nil, err
	}

	modelName := c.modelMap[modelType]
	for _, m := range models {
		if m.Name == modelName {
			return &m, nil
		}
	}

	return &ModelStatus{
		Name:      modelName,
		Available: false,
	}, nil
}

// TestConnection tests the connection to Ollama
func (c *Client) TestConnection(ctx context.Context) error {
	_, err := c.ollamaClient.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	c.logger.Info("Successfully connected to Ollama")
	return nil
}

// Close cleans up any resources
func (c *Client) Close() error {
	// Ollama client doesn't need explicit close
	// But we keep this method for interface consistency
	return nil
}

// Helper functions

// formatSize formats byte size to human readable
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// IsRetryableError determines if an error can be retried
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Connection errors are retryable
	if strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "refused") {
		return true
	}

	return false
}

// Prompt templates for common tasks - PURELY DESCRIPTIVE, no fixes
var PromptTemplates = struct {
	// Code understanding (CodeBERT)
	CodeAnalysis string

	// Test generation (DeepSeek)
	TestGeneration string

	// Documentation (Gemma)
	APIDocumentation string

	// Failure analysis (DeepSeek) - EXPLAINS only, no fixes
	FailureAnalysis string
}{
	CodeAnalysis: `Analyze this %s code and describe its structure.
DO NOT suggest improvements - just describe what you see.

Code:
%s

Describe:
- Functions and their purposes
- Parameters and return values
- Dependencies
- Any notable patterns

Keep the description factual and objective.`,

	TestGeneration: `Generate unit tests for this %s function.
DO NOT modify the original code - only create tests.

Function: %s
Code:
%s

Create tests that verify the function's behavior.
Use %s testing framework.
Return ONLY the test code, no explanations.`,

	APIDocumentation: `Write clear documentation for this %s function.
DO NOT suggest changes - just document what it does.

Function: %s
Code:
%s

Include:
- Brief description of purpose
- Parameter explanations
- Return value description
- One simple example

Use %s documentation format.`,

	FailureAnalysis: `Analyze why this test failed.
DO NOT suggest code fixes - just explain what went wrong.

Test: %s
Error: %s
Test Code:
%s
Source Code:
%s

Explain the mismatch between expected and actual behavior.
Focus on what happened, not how to fix it.`,
}
