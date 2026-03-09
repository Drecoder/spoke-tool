package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/spoke-tool/api/types"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if cfg.Models.Encoder != "codebert" {
		t.Errorf("expected encoder 'codebert', got %s", cfg.Models.Encoder)
	}

	if cfg.TestSpoke.CoverageThreshold != 80.0 {
		t.Errorf("expected coverage threshold 80, got %f", cfg.TestSpoke.CoverageThreshold)
	}

	if len(cfg.ReadmeSpoke.Sections) == 0 {
		t.Error("expected at least one README section")
	}
}

func TestLoadSave(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	configPath := filepath.Join(tmpdir, "config.yaml")

	// Create a test config
	original := DefaultConfig()
	original.ProjectRoot = "/test/project"
	original.Models.Encoder = "test-encoder"

	// Save it
	err = Save(original, configPath)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load it back
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Compare
	if loaded.ProjectRoot != original.ProjectRoot {
		t.Errorf("ProjectRoot: expected %s, got %s", original.ProjectRoot, loaded.ProjectRoot)
	}
	if loaded.Models.Encoder != original.Models.Encoder {
		t.Errorf("Encoder: expected %s, got %s", original.Models.Encoder, loaded.Models.Encoder)
	}
}

func TestValidateConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Valid config should not error
	err := validateConfig(cfg)
	if err != nil {
		t.Errorf("validateConfig failed on valid config: %v", err)
	}

	// Invalid coverage threshold
	cfg.TestSpoke.CoverageThreshold = 150
	err = validateConfig(cfg)
	if err == nil {
		t.Error("validateConfig should error on invalid coverage threshold")
	}
	cfg.TestSpoke.CoverageThreshold = 80 // Reset

	// Invalid CPU percent
	cfg.Squeeze.MaxCPUPercent = 120
	err = validateConfig(cfg)
	if err == nil {
		t.Error("validateConfig should error on invalid CPU percent")
	}
}

func TestMerge(t *testing.T) {
	base := DefaultConfig()
	base.ProjectRoot = "/base"
	base.Models.Encoder = "base-encoder"
	base.TestSpoke.Frameworks = map[types.Language]string{
		types.Go: "base-testing",
	}

	override := &types.Config{
		ProjectRoot: "/override",
		Models: struct {
			Encoder string `json:"encoder" yaml:"encoder"`
			Decoder string `json:"decoder" yaml:"decoder"`
			Fast    string `json:"fast" yaml:"fast"`
		}{
			Encoder: "override-encoder",
		},
		TestSpoke: struct {
			Enabled           bool                      `json:"enabled" yaml:"enabled"`
			AutoRun           bool                      `json:"auto_run" yaml:"auto_run"`
			CoverageThreshold float64                   `json:"coverage_threshold" yaml:"coverage_threshold"`
			Frameworks        map[types.Language]string `json:"frameworks" yaml:"frameworks"`
		}{
			Frameworks: map[types.Language]string{
				types.Go: "override-testing",
			},
		},
	}

	merged := Merge(base, override)

	if merged.ProjectRoot != "/override" {
		t.Errorf("expected project root /override, got %s", merged.ProjectRoot)
	}
	if merged.Models.Encoder != "override-encoder" {
		t.Errorf("expected encoder override-encoder, got %s", merged.Models.Encoder)
	}
	if merged.TestSpoke.Frameworks[types.Go] != "override-testing" {
		t.Errorf("expected framework override-testing, got %s", merged.TestSpoke.Frameworks[types.Go])
	}
}

func TestConfigManager(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	configPath := filepath.Join(tmpdir, "config.yaml")

	// Create manager with non-existent file (should use defaults)
	manager, err := NewConfigManager(configPath)
	if err != nil {
		t.Fatalf("NewConfigManager failed: %v", err)
	}

	if manager.Get() == nil {
		t.Fatal("manager.Get() returned nil")
	}

	// Update config
	err = manager.Update(func(cfg *types.Config) {
		cfg.ProjectRoot = "/updated"
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if manager.Get().ProjectRoot != "/updated" {
		t.Errorf("expected project root /updated, got %s", manager.Get().ProjectRoot)
	}

	// Reload from disk
	err = manager.Reload()
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	if manager.Get().ProjectRoot != "/updated" {
		t.Errorf("after reload, expected /updated, got %s", manager.Get().ProjectRoot)
	}
}

func TestWriteExample(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	examplePath := filepath.Join(tmpdir, "example.yaml")

	err = WriteExample(examplePath)
	if err != nil {
		t.Fatalf("WriteExample failed: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(examplePath); err != nil {
		t.Errorf("example file not created: %v", err)
	}
}
