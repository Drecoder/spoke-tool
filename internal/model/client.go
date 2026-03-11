package model

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"example.com/spoke-tool/api/types"
	"example.com/spoke-tool/internal/common"
	"github.com/google/uuid"
	"github.com/ollama/ollama/api"
)

const (
	CodeLLamaEncoder ModelType = "codellama-encoder" // For encoder role
	CodeLLamaDecoder ModelType = "codellama-decoder" // For decoder role
	Gemma2BChat      ModelType = "gemma2:2b"
)

// SLMRequest represents a request to an SLM
type SLMRequest struct {
	ID          string                 `json:"id"`
	Model       ModelType              `json:"model"`
	Language    types.Language         `json:"language"`
	Prompt      string                 `json:"prompt"`
	System      string                 `json:"system,omitempty"`
	Context     []int                  `json:"context,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Temperature float32                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Response represents a response from an SLM
type Response struct {
	ID         string         `json:"id"`
	RequestID  string         `json:"request_id"`
	Model      ModelType      `json:"model"`
	Language   types.Language `json:"language"`
	Response   string         `json:"response"`
	Error      string         `json:"error,omitempty"`
	Duration   time.Duration  `json:"duration"`
	TokensUsed int            `json:"tokens_used"`
	Done       bool           `json:"done"`
	Timestamp  time.Time      `json:"timestamp"`
}

// Client handles communication with Ollama SLMs
type Client struct {
	client   *api.Client
	modelMap map[ModelType]string
	timeout  time.Duration
	logger   *common.Logger
}

// ClientConfig holds configuration for the model client
type ClientConfig struct {
	OllamaHost string
	Timeout    time.Duration
	Models     map[ModelType]string
}

// NewClient creates a new SLM client
func NewClient(config ClientConfig) (*Client, error) {
	if config.OllamaHost == "" {
		config.OllamaHost = "http://localhost:11434"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Parse URL
	url, err := url.Parse(config.OllamaHost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ollama host: %w", err)
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	// Create Ollama client
	client := api.NewClient(url, httpClient)

	// Default model mapping
	modelMap := config.Models
	if modelMap == nil {
		modelMap = map[ModelType]string{
			CodeLLamaEncoder: "codellama:7b",
			CodeLLamaDecoder: "codellama:7b",
			Gemma2BChat:      "gemma2:2b",
		}
	}

	return &Client{
		client:   client,
		modelMap: modelMap,
		timeout:  config.Timeout,
		logger:   common.GetLogger().WithField("component", "model-client"),
	}, nil
}

// Generate sends a prompt to the specified model and returns the response
func (c *Client) Generate(ctx context.Context, req SLMRequest) (*Response, error) {
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

	// Generate
	resp := &Response{
		ID:        uuid.New().String(),
		RequestID: req.ID,
		Model:     req.Model,
		Language:  req.Language,
		Timestamp: time.Now(),
	}

	// Create generate request
	genReq := &api.GenerateRequest{
		Model:   modelName,
		Prompt:  req.Prompt,
		System:  req.System,
		Options: req.Options,
	}

	// For streaming response
	var fullResponse strings.Builder
	var tokens int

	// Call generate with a callback for streaming
	err := c.client.Generate(ctx, genReq, func(response api.GenerateResponse) error {
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

// CheckModels verifies that required models are available
func (c *Client) CheckModels(ctx context.Context) (map[ModelType]bool, error) {
	status := make(map[ModelType]bool)

	// List available models
	models, err := c.client.List(ctx)
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
	for modelType, modelName := range c.modelMap {
		status[modelType] = available[modelName]
		if !available[modelName] {
			c.logger.Warn("Model not available", "model", modelName)
		}
	}

	return status, nil
}

// ListModels returns all available models
func (c *Client) ListModels(ctx context.Context) ([]api.ListModelResponse, error) {
	resp, err := c.client.List(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Models, nil
}

// Close cleans up any resources
func (c *Client) Close() error {
	return nil
}
