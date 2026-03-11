package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"example.com/spoke-tool/api/types"
	"gopkg.in/yaml.v3"
)

// Config represents the complete configuration
type Config struct {
	// Project settings
	ProjectRoot string `json:"project_root" yaml:"project_root"`

	// Model settings
	Models ModelConfig `json:"models" yaml:"models"`

	// Spoke settings
	TestSpoke   TestSpokeConfig   `json:"test_spoke" yaml:"test_spoke"`
	ReadmeSpoke ReadmeSpokeConfig `json:"readme_spoke" yaml:"readme_spoke"`

	// Performance settings
	Squeeze SqueezeConfig `json:"squeeze" yaml:"squeeze"`

	// Audit settings
	Audit AuditConfig `json:"audit" yaml:"audit"`

	// General settings
	LogLevel   string `json:"log_level" yaml:"log_level"`
	LogJSON    bool   `json:"log_json" yaml:"log_json"`
	LogColor   bool   `json:"log_color" yaml:"log_color"`
	configPath string `json:"-" yaml:"-"` // Path to config file
}

// ModelConfig contains model-related settings
type ModelConfig struct {
	Encoder string `json:"encoder" yaml:"encoder"`
	Decoder string `json:"decoder" yaml:"decoder"`
	Fast    string `json:"fast" yaml:"fast"`

	// Model parameters
	Temperature float32 `json:"temperature" yaml:"temperature"`
	MaxTokens   int     `json:"max_tokens" yaml:"max_tokens"`
	Timeout     string  `json:"timeout" yaml:"timeout"`

	// Ollama settings
	OllamaHost string `json:"ollama_host" yaml:"ollama_host"`
}

// TestSpokeConfig contains test generation settings
type TestSpokeConfig struct {
	Enabled           bool              `json:"enabled" yaml:"enabled"`
	AutoRun           bool              `json:"auto_run" yaml:"auto_run"`
	CoverageThreshold float64           `json:"coverage_threshold" yaml:"coverage_threshold"`
	Frameworks        map[string]string `json:"frameworks" yaml:"frameworks"`

	// Generation settings
	MaxTestsPerFunction int               `json:"max_tests_per_function" yaml:"max_tests_per_function"`
	IncludeEdgeCases    bool              `json:"include_edge_cases" yaml:"include_edge_cases"`
	GenerateMocks       bool              `json:"generate_mocks" yaml:"generate_mocks"`
	TestFilePatterns    map[string]string `json:"test_file_patterns" yaml:"test_file_patterns"`

	// Language-specific settings
	Languages map[string]LanguageConfig `json:"languages" yaml:"languages"`
}

// LanguageConfig contains language-specific settings
type LanguageConfig struct {
	Framework     string   `json:"framework" yaml:"framework"`
	TestPattern   string   `json:"test_pattern" yaml:"test_pattern"`
	CoverCommand  string   `json:"cover_command" yaml:"cover_command"`
	TestCommand   string   `json:"test_command" yaml:"test_command"`
	MockFramework string   `json:"mock_framework" yaml:"mock_framework"`
	Extensions    []string `json:"extensions" yaml:"extensions"`
}

// ReadmeSpokeConfig contains documentation generation settings
type ReadmeSpokeConfig struct {
	Enabled    bool     `json:"enabled" yaml:"enabled"`
	AutoUpdate bool     `json:"auto_update" yaml:"auto_update"`
	Sections   []string `json:"sections" yaml:"sections"`

	// Generation settings
	IncludeExamples    bool `json:"include_examples" yaml:"include_examples"`
	MaxExamplesPerFunc int  `json:"max_examples_per_function" yaml:"max_examples_per_function"`
	PreserveManual     bool `json:"preserve_manual" yaml:"preserve_manual"`

	// Template settings
	TemplateFile string `json:"template_file" yaml:"template_file"`
	OutputFile   string `json:"output_file" yaml:"output_file"`

	// Language-specific doc formats
	DocFormats map[string]string `json:"doc_formats" yaml:"doc_formats"`
}

// SqueezeConfig contains performance tuning settings
type SqueezeConfig struct {
	Enabled         bool `json:"enabled" yaml:"enabled"`
	MaxCPUPercent   int  `json:"max_cpu_percent" yaml:"max_cpu_percent"`
	MaxMemoryMB     int  `json:"max_memory_mb" yaml:"max_memory_mb"`
	IdleThresholdMs int  `json:"idle_threshold_ms" yaml:"idle_threshold_ms"`
	MaxConcurrent   int  `json:"max_concurrent" yaml:"max_concurrent"`
	MinConcurrent   int  `json:"min_concurrent" yaml:"min_concurrent"`
}

// AuditConfig contains audit logging settings
type AuditConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Path    string `json:"path" yaml:"path"`
	Retain  int    `json:"retain_days" yaml:"retain_days"`
	JSON    bool   `json:"json" yaml:"json"`
}

// Load loads configuration from a file
func Load(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg types.Config

	// Try YAML first, then JSON based on extension
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	default:
		// Try both
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			if err := json.Unmarshal(data, &cfg); err != nil {
				return nil, fmt.Errorf("failed to parse config (tried YAML and JSON): %w", err)
			}
		}
	}

	// Set defaults for missing values
	cfg = setDefaults(cfg)

	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// LoadWithEnv loads configuration and overrides with environment variables
func LoadWithEnv(path string) (*types.Config, error) {
	cfg, err := Load(path)
	if err != nil {
		return nil, err
	}

	// Override with environment variables
	overrideFromEnv(cfg)

	return cfg, nil
}

// Save saves configuration to a file
func Save(cfg *types.Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *types.Config {
	return &types.Config{
		ProjectRoot: ".",

		Models: struct {
			Encoder string `json:"encoder" yaml:"encoder"`
			Decoder string `json:"decoder" yaml:"decoder"`
			Fast    string `json:"fast" yaml:"fast"`
		}{
			Encoder: "codellama:7b",
			Decoder: "codellama:7b",
			Fast:    "gemma2:2b",
		},

		TestSpoke: struct {
			Enabled           bool                      `json:"enabled" yaml:"enabled"`
			AutoRun           bool                      `json:"auto_run" yaml:"auto_run"`
			CoverageThreshold float64                   `json:"coverage_threshold" yaml:"coverage_threshold"`
			Frameworks        map[types.Language]string `json:"frameworks" yaml:"frameworks"`
		}{
			Enabled:           true,
			AutoRun:           true,
			CoverageThreshold: 80.0,
			Frameworks: map[types.Language]string{
				types.Go:     "testing",
				types.NodeJS: "jest",
				types.Python: "pytest",
			},
		},

		ReadmeSpoke: struct {
			Enabled    bool               `json:"enabled" yaml:"enabled"`
			AutoUpdate bool               `json:"auto_update" yaml:"auto_update"`
			Sections   []types.DocSection `json:"sections" yaml:"sections"`
		}{
			Enabled:    true,
			AutoUpdate: true,
			Sections: []types.DocSection{
				types.DocSectionTitle,
				types.DocSectionInstallation,
				types.DocSectionQuickStart,
				types.DocSectionAPI,
				types.DocSectionExamples,
				types.DocSectionContributing,
				types.DocSectionLicense,
			},
		},

		Squeeze: struct {
			MaxCPUPercent int `json:"max_cpu_percent" yaml:"max_cpu_percent"`
			MaxMemoryMB   int `json:"max_memory_mb" yaml:"max_memory_mb"`
			IdleThreshold int `json:"idle_threshold_ms" yaml:"idle_threshold_ms"`
		}{
			MaxCPUPercent: 80,
			MaxMemoryMB:   4096,
			IdleThreshold: 500,
		},

		Audit: struct {
			Enabled bool   `json:"enabled" yaml:"enabled"`
			Path    string `json:"path" yaml:"path"`
		}{
			Enabled: true,
			Path:    "audit.log",
		},
	}
}

// setDefaults sets default values for missing configuration
func setDefaults(cfg types.Config) types.Config {
	// Project root
	if cfg.ProjectRoot == "" {
		cfg.ProjectRoot = "."
	}

	// Model defaults
	if cfg.Models.Encoder == "" {
		cfg.Models.Encoder = "codellama:7b"
	}
	if cfg.Models.Decoder == "" {
		cfg.Models.Decoder = "codellama:7b"
	}
	if cfg.Models.Fast == "" {
		cfg.Models.Fast = "gemma2:2b"
	}

	// Test spoke defaults
	if cfg.TestSpoke.Frameworks == nil {
		cfg.TestSpoke.Frameworks = map[types.Language]string{
			types.Go:     "testing",
			types.NodeJS: "jest",
			types.Python: "pytest",
		}
	}

	// Readme spoke defaults
	if len(cfg.ReadmeSpoke.Sections) == 0 {
		cfg.ReadmeSpoke.Sections = []types.DocSection{
			types.DocSectionTitle,
			types.DocSectionInstallation,
			types.DocSectionQuickStart,
			types.DocSectionAPI,
			types.DocSectionExamples,
		}
	}

	// Squeeze defaults
	if cfg.Squeeze.MaxCPUPercent == 0 {
		cfg.Squeeze.MaxCPUPercent = 80
	}
	if cfg.Squeeze.MaxMemoryMB == 0 {
		cfg.Squeeze.MaxMemoryMB = 4096
	}
	if cfg.Squeeze.IdleThreshold == 0 {
		cfg.Squeeze.IdleThreshold = 500
	}

	return cfg
}

// validateConfig validates the configuration
func validateConfig(cfg *types.Config) error {
	// Validate models
	if cfg.Models.Encoder == "" {
		return fmt.Errorf("encoder model cannot be empty")
	}
	if cfg.Models.Decoder == "" {
		return fmt.Errorf("decoder model cannot be empty")
	}
	if cfg.Models.Fast == "" {
		return fmt.Errorf("fast model cannot be empty")
	}

	// Validate test spoke
	if cfg.TestSpoke.CoverageThreshold < 0 || cfg.TestSpoke.CoverageThreshold > 100 {
		return fmt.Errorf("coverage threshold must be between 0 and 100")
	}

	// Validate frameworks for supported languages
	supported := map[types.Language]bool{
		types.Go:     true,
		types.NodeJS: true,
		types.Python: true,
	}

	for lang, framework := range cfg.TestSpoke.Frameworks {
		if !supported[lang] {
			return fmt.Errorf("unsupported language: %s", lang)
		}
		if framework == "" {
			return fmt.Errorf("framework cannot be empty for language: %s", lang)
		}
	}

	// Validate squeeze
	if cfg.Squeeze.MaxCPUPercent < 0 || cfg.Squeeze.MaxCPUPercent > 100 {
		return fmt.Errorf("max CPU percent must be between 0 and 100")
	}
	if cfg.Squeeze.MaxMemoryMB < 0 {
		return fmt.Errorf("max memory must be positive")
	}
	if cfg.Squeeze.IdleThreshold < 0 {
		return fmt.Errorf("idle threshold must be positive")
	}

	return nil
}

// overrideFromEnv overrides configuration with environment variables
func overrideFromEnv(cfg *types.Config) {
	// Model overrides
	if val := os.Getenv("SPOKE_ENCODER_MODEL"); val != "" {
		cfg.Models.Encoder = val
	}
	if val := os.Getenv("SPOKE_DECODER_MODEL"); val != "" {
		cfg.Models.Decoder = val
	}
	if val := os.Getenv("SPOKE_FAST_MODEL"); val != "" {
		cfg.Models.Fast = val
	}
	if val := os.Getenv("SPOKE_OLLAMA_HOST"); val != "" {
		// Need to add OllamaHost to ModelConfig
	}

	// Test spoke overrides
	if val := os.Getenv("SPOKE_TEST_ENABLED"); val == "false" {
		cfg.TestSpoke.Enabled = false
	}
	if val := os.Getenv("SPOKE_AUTO_RUN"); val == "false" {
		cfg.TestSpoke.AutoRun = false
	}
	if val := os.Getenv("SPOKE_COVERAGE_THRESHOLD"); val != "" {
		fmt.Sscanf(val, "%f", &cfg.TestSpoke.CoverageThreshold)
	}

	// Readme spoke overrides
	if val := os.Getenv("SPOKE_README_ENABLED"); val == "false" {
		cfg.ReadmeSpoke.Enabled = false
	}
	if val := os.Getenv("SPOKE_AUTO_UPDATE"); val == "false" {
		cfg.ReadmeSpoke.AutoUpdate = false
	}

	// Squeeze overrides
	if val := os.Getenv("SPOKE_MAX_CPU"); val != "" {
		fmt.Sscanf(val, "%d", &cfg.Squeeze.MaxCPUPercent)
	}
	if val := os.Getenv("SPOKE_MAX_MEMORY"); val != "" {
		fmt.Sscanf(val, "%d", &cfg.Squeeze.MaxMemoryMB)
	}

	// Log level
	if val := os.Getenv("SPOKE_LOG_LEVEL"); val != "" {
		// Need to add LogLevel to Config
	}
}

// Merge merges two configurations (second overrides first)
func Merge(base, override *types.Config) *types.Config {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	result := *base

	// Override project root
	if override.ProjectRoot != "" {
		result.ProjectRoot = override.ProjectRoot
	}

	// Override models
	if override.Models.Encoder != "" {
		result.Models.Encoder = override.Models.Encoder
	}
	if override.Models.Decoder != "" {
		result.Models.Decoder = override.Models.Decoder
	}
	if override.Models.Fast != "" {
		result.Models.Fast = override.Models.Fast
	}

	// Override test spoke
	if override.TestSpoke.Enabled != result.TestSpoke.Enabled {
		result.TestSpoke.Enabled = override.TestSpoke.Enabled
	}
	if override.TestSpoke.AutoRun != result.TestSpoke.AutoRun {
		result.TestSpoke.AutoRun = override.TestSpoke.AutoRun
	}
	if override.TestSpoke.CoverageThreshold > 0 {
		result.TestSpoke.CoverageThreshold = override.TestSpoke.CoverageThreshold
	}

	// Merge frameworks
	for k, v := range override.TestSpoke.Frameworks {
		result.TestSpoke.Frameworks[k] = v
	}

	// Override readme spoke
	if override.ReadmeSpoke.Enabled != result.ReadmeSpoke.Enabled {
		result.ReadmeSpoke.Enabled = override.ReadmeSpoke.Enabled
	}
	if override.ReadmeSpoke.AutoUpdate != result.ReadmeSpoke.AutoUpdate {
		result.ReadmeSpoke.AutoUpdate = override.ReadmeSpoke.AutoUpdate
	}
	if len(override.ReadmeSpoke.Sections) > 0 {
		result.ReadmeSpoke.Sections = override.ReadmeSpoke.Sections
	}

	// Override squeeze
	if override.Squeeze.MaxCPUPercent > 0 {
		result.Squeeze.MaxCPUPercent = override.Squeeze.MaxCPUPercent
	}
	if override.Squeeze.MaxMemoryMB > 0 {
		result.Squeeze.MaxMemoryMB = override.Squeeze.MaxMemoryMB
	}
	if override.Squeeze.IdleThreshold > 0 {
		result.Squeeze.IdleThreshold = override.Squeeze.IdleThreshold
	}

	return &result
}

// WriteExample writes an example configuration file
func WriteExample(path string) error {
	// Add comments to the example
	example := `# Spoke Tool Configuration
# 
# This is an example configuration file.
# Copy this to config.yaml and modify as needed.

# Project root directory (relative or absolute)
project_root: "."

# Model settings
models:
  # Encoder model for code understanding
  encoder: "codellama:7b"
  
  # Decoder model for complex tasks (test generation, error analysis)
  decoder: "codellama:7b"
  
  # Fast model for simple tasks (documentation generation)
  fast: "gemma2:2b"
  
  # Ollama host (default: http://localhost:11434)
  ollama_host: "http://localhost:11434"
  
  # Generation parameters
  temperature: 0.7
  max_tokens: 2048
  timeout: "30s"

# Test generation spoke settings
test_spoke:
  enabled: true
  auto_run: true
  coverage_threshold: 80
  
  # Test frameworks by language
  frameworks:
    go: "testing"
    nodejs: "jest"
    python: "pytest"
  
  # Test file patterns by language
  test_file_patterns:
    go: "*_test.go"
    nodejs: "*.test.js"
    python: "test_*.py"
  
  # Generation options
  max_tests_per_function: 10
  include_edge_cases: true
  generate_mocks: true
  
  # Language-specific settings
  languages:
    go:
      framework: "testing"
      test_pattern: "*_test.go"
      cover_command: "go test -cover"
      test_command: "go test"
      extensions: [".go"]
    
    nodejs:
      framework: "jest"
      test_pattern: "*.test.js"
      cover_command: "jest --coverage"
      test_command: "jest"
      extensions: [".js", ".ts"]
    
    python:
      framework: "pytest"
      test_pattern: "test_*.py"
      cover_command: "pytest --cov=."
      test_command: "pytest"
      extensions: [".py"]

# README generation spoke settings
readme_spoke:
  enabled: true
  auto_update: true
  
  # Sections to include (order matters)
  sections:
    - title
    - badges
    - description
    - installation
    - quickstart
    - api
    - examples
    - contributing
    - license
  
  include_examples: true
  max_examples_per_function: 3
  preserve_manual: true
  
  # Template file (optional)
  # template_file: "README.tmpl.md"
  output_file: "README.md"
  
  # Documentation formats by language
  doc_formats:
    go: "godoc"
    nodejs: "jsdoc"
    python: "pydoc"

# Performance tuning (Squeeze mechanism)
squeeze:
  enabled: true
  max_cpu_percent: 80
  max_memory_mb: 4096
  idle_threshold_ms: 500
  max_concurrent: 4
  min_concurrent: 1

# Audit logging
audit:
  enabled: true
  path: "audit.log"
  retain_days: 30
  json: true

# Logging
log_level: "info"  # debug, info, warn, error
log_json: false
log_color: true
`

	return os.WriteFile(path, []byte(example), 0644)
}

// ConfigManager handles configuration operations
type ConfigManager struct {
	config *types.Config
	path   string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(path string) (*ConfigManager, error) {
	cfg, err := Load(path)
	if err != nil {
		// If file doesn't exist, use defaults
		if os.IsNotExist(err) {
			cfg = DefaultConfig()
		} else {
			return nil, err
		}
	}

	return &ConfigManager{
		config: cfg,
		path:   path,
	}, nil
}

// Get returns the current configuration
func (m *ConfigManager) Get() *types.Config {
	return m.config
}

// Update updates the configuration
func (m *ConfigManager) Update(updater func(*types.Config)) error {
	updater(m.config)

	// Validate after update
	if err := validateConfig(m.config); err != nil {
		return err
	}

	// Save to file
	return Save(m.config, m.path)
}

// Reload reloads the configuration from disk
func (m *ConfigManager) Reload() error {
	cfg, err := Load(m.path)
	if err != nil {
		return err
	}
	m.config = cfg
	return nil
}

// GetLogLevel returns the log level as a string
func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

// GetLogJSON returns whether to use JSON logging
func (c *Config) GetLogJSON() bool {
	return c.LogJSON
}

// GetLogColor returns whether to use colored output
func (c *Config) GetLogColor() bool {
	return c.LogColor
}

// GetAuditFile returns the audit file path
func (c *Config) GetAuditFile() string {
	if !c.Audit.Enabled {
		return ""
	}
	return c.Audit.Path
}
