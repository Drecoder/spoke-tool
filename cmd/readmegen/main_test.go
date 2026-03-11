package main

import (
	"path/filepath"
	"testing"

	"example.com/spoke-tool/api/types"
	"example.com/spoke-tool/internal/model"
)

func TestLoadConfig(t *testing.T) {
	// Create temp dir for test
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Test with non-existent config (should return defaults)
	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("loadConfig with missing file failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.ProjectRoot != "" {
		t.Errorf("expected empty project root, got %s", cfg.ProjectRoot)
	}
}

func loadConfig(configPath string) (*types.Config, error) {
	// TODO: Implement loadConfig function
	return &types.Config{}, nil
}

func hasChanged(old, new *types.CodeAnalysis) bool {
	if len(old.Files) != len(new.Files) || len(old.Functions) != len(new.Functions) {
		return true
	}
	return old.Timestamp != new.Timestamp
}

func getModelName(modelType model.ModelType, cfg *types.Config) string {
	switch modelType {
	case model.CodeLLamaEncoder:
		return cfg.Models.Encoder
	case model.CodeLLamaDecoder:
		return cfg.Models.Decoder
	case model.Gemma2B:
		return cfg.Models.Fast
	default:
		return ""
	}
}

func TestGetModelName(t *testing.T) {
	cfg := &types.Config{
		Models: struct {
			Encoder string `json:"encoder" yaml:"encoder"`
			Decoder string `json:"decoder" yaml:"decoder"`
			Fast    string `json:"fast" yaml:"fast"`
		}{
			Encoder: "codellama:7b",
			Decoder: "codellama:7b",
			Fast:    "gemma2:2b",
		},
	}

	tests := []struct {
		modelType model.ModelType
		expected  string
	}{
		{model.CodeLLamaEncoder, "codellama:7b"},
		{model.CodeLLamaDecoder, "codellama:7b"},
		{model.Gemma2B, "gemma2:2b"},
	}

	for _, tt := range tests {
		t.Run(string(tt.modelType), func(t *testing.T) {
			got := getModelName(tt.modelType, cfg)
			if got != tt.expected {
				t.Errorf("getModelName() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestHasChanged(t *testing.T) {
	old := &types.CodeAnalysis{
		Files:     []types.CodeFile{{Path: "test.go"}},
		Functions: []types.Function{{Name: "Test"}},
		Timestamp: "2024-01-01T00:00:00Z",
	}

	new := &types.CodeAnalysis{
		Files:     []types.CodeFile{{Path: "test.go"}},
		Functions: []types.Function{{Name: "Test"}},
		Timestamp: "2024-01-01T00:00:01Z",
	}

	// Same content, different timestamp - should return true
	if !hasChanged(old, new) {
		t.Error("hasChanged() should return true for different timestamp")
	}

	// Different number of files
	new.Files = append(new.Files, types.CodeFile{Path: "test2.go"})
	if !hasChanged(old, new) {
		t.Error("hasChanged() should return true for different file count")
	}
}
