//go:build go1.18
// +build go1.18

package fuzz

import (
	"testing"

	"example.com/spoke-tool/internal/config"
)

// FuzzConfigParsing tests the config parser with various inputs
// to ensure it never panics and handles edge cases gracefully.
func FuzzConfigParsing(f *testing.F) {
	// Seed corpus with various YAML configs
	seeds := []string{
		"",
		"{}",
		"test_spoke:\n  enabled: true",
		"test_spoke:\n  enabled: false",
		"models:\n  encoder: codellama:7b",
		"models:\n  decoder: codellama:7b",
		"models:\n  fast: gemma2:2b",
		"test_spoke:\n  coverage_threshold: 80",
		"test_spoke:\n  auto_run: true",
		"readme_spoke:\n  enabled: true",
		"readme_spoke:\n  auto_update: true",
		"readme_spoke:\n  sections:\n    - installation\n    - quickstart",
		"squeeze:\n  max_cpu_percent: 80",
		"squeeze:\n  max_memory_mb: 4096",
		"squeeze:\n  idle_threshold_ms: 500",
		"audit:\n  enabled: true\n  path: audit.log",
		"project_root: /path/to/project",
		"log_level: debug",
		"log_level: info",
		"log_level: warn",
		"log_level: error",
		"log_json: true",
		"log_color: true",
		"test_spoke:\n  enabled: not-a-bool",
		"test_spoke:\n  coverage_threshold: not-a-number",
		"test_spoke:\n  coverage_threshold: -10",
		"test_spoke:\n  coverage_threshold: 150",
		"test_spoke:\n  frameworks:\n    go: testing\n    nodejs: jest\n    python: pytest",
		"test_spoke:\n  frameworks:\n    go: \n    nodejs: ",
		"test_spoke:\n  test_file_patterns:\n    go: '*_test.go'\n    nodejs: '*.test.js'\n    python: 'test_*.py'",
		"test_spoke:\n  max_tests_per_function: 10",
		"test_spoke:\n  max_tests_per_function: -5",
		"test_spoke:\n  include_edge_cases: true",
		"test_spoke:\n  generate_mocks: true",
		"test_spoke:\n  languages:\n    go:\n      framework: testing\n      test_pattern: '*_test.go'",
		"test_spoke:\n  languages:\n    rust:\n      framework: rust_test",
		"models:\n  temperature: 0.7",
		"models:\n  temperature: 2.5",
		"models:\n  temperature: -1",
		"models:\n  max_tokens: 2048",
		"models:\n  max_tokens: 0",
		"models:\n  max_tokens: -100",
		"models:\n  timeout: 30s",
		"models:\n  timeout: invalid",
		"models:\n  ollama_host: http://localhost:11434",
		"models:\n  ollama_host: ",
		"readme_spoke:\n  sections:\n    - invalid-section",
		"readme_spoke:\n  sections: not-a-list",
		"readme_spoke:\n  include_examples: true",
		"readme_spoke:\n  max_examples_per_function: 3",
		"readme_spoke:\n  preserve_manual: true",
		"readme_spoke:\n  template_file: README.tmpl.md",
		"readme_spoke:\n  output_file: README.md",
		"readme_spoke:\n  doc_formats:\n    go: godoc\n    nodejs: jsdoc\n    python: pydoc",
		"squeeze:\n  enabled: true",
		"squeeze:\n  max_cpu_percent: 0",
		"squeeze:\n  max_cpu_percent: 101",
		"squeeze:\n  max_memory_mb: 0",
		"squeeze:\n  max_memory_mb: -10",
		"squeeze:\n  idle_threshold_ms: 0",
		"squeeze:\n  max_concurrent: 4",
		"squeeze:\n  max_concurrent: 0",
		"squeeze:\n  min_concurrent: 1",
		"squeeze:\n  min_concurrent: 10",
		"audit:\n  enabled: true",
		"audit:\n  path: ",
		"audit:\n  retain_days: 30",
		"audit:\n  retain_days: -5",
		"audit:\n  json: true",
		"---",
		"# comment only",
		"test_spoke:\n  # nested comment\n  enabled: true",
		"\t\ttest_spoke:\n\t\t\tenabled: true",
		"test_spoke: { enabled: true }",
		"test_spoke: { enabled: true, auto_run: false }",
		"[invalid yaml",
		": : :",
		"test_spoke: [list, not, map]",
		"models:\n  - encoder\n  - decoder",
		"test_spoke:\n  enabled: !!str true",
		"test_spoke:\n  coverage_threshold: !!float 80",
		"squeeze:\n  max_cpu_percent: !!int 80",
		"test_spoke:\n  frameworks: !!map\n    go: testing",
		"test_spoke:\n  description: " + string(make([]byte, 10000)),
		"models:\n  encoder: " + string(make([]byte, 1000)),
		"squeeze:\n  max_cpu_percent: " + string(make([]byte, 100)),
		"test_spoke:\n  - " + string(make([]byte, 100)),
		"test_spoke:\n  enabled: true\n" + string(make([]byte, 1000)),
		"%YAML 1.2\n---\ntest_spoke:\n  enabled: true",
		"%TAG ! tag:example.com,2024:\n---\ntest_spoke:\n  enabled: true",
		"test_spoke:\n  enabled: !!binary |\n    c3RyCg==",
		"test_spoke:\n  enabled: !!timestamp 2024-01-01",
		"test_spoke:\n  enabled: !!set {a, b, c}",
		"test_spoke:\n  enabled: !!omap\n    - a: 1\n    - b: 2",
		"<<: *anchor",
		"test_spoke: &anchor\n  enabled: true\nother: *anchor",
		"test_spoke:\n  enabled: !!python/name:sys.stdout",
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F",
		"test_spoke:\n  enabled: \x00\x01\x02",
	}

	// Add seeds to corpus
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data string) {
		// This should NEVER panic - just parse config
		cfg, err := config.Parse([]byte(data))

		// We don't care about the result validity, just that:
		// 1. No panic occurred
		// 2. Return values are valid (not nil pointer dereference)
		// 3. Error may be nil or non-nil, both are acceptable

		if cfg == nil && err == nil {
			// Both nil is suspicious - either should have config or error
			t.Log("Both config and error are nil")
		}

		// Basic sanity checks - if we got a config, access fields to ensure no panic
		if cfg != nil {
			// Access top-level fields (should never panic even if partially parsed)
			_ = cfg.ProjectRoot
			_ = cfg.Models.Encoder
			_ = cfg.Models.Decoder
			_ = cfg.Models.Fast
			_ = cfg.TestSpoke.Enabled
			_ = cfg.TestSpoke.AutoRun
			_ = cfg.TestSpoke.CoverageThreshold
			_ = cfg.ReadmeSpoke.Enabled
			_ = cfg.ReadmeSpoke.AutoUpdate
			_ = cfg.ReadmeSpoke.Sections
			_ = cfg.Squeeze.MaxCPUPercent
			_ = cfg.Squeeze.MaxMemoryMB
			_ = cfg.Squeeze.IdleThreshold
			_ = cfg.Audit.Enabled
			_ = cfg.Audit.Path
			_ = cfg.LogLevel
			_ = cfg.LogJSON
			_ = cfg.LogColor

			// Access map fields safely (should never panic)
			if cfg.TestSpoke.Frameworks != nil {
				for lang, framework := range cfg.TestSpoke.Frameworks {
					_ = lang
					_ = framework
				}
			}
		}
	})
}

// FuzzConfigValidation tests the config validation logic
func FuzzConfigValidation(f *testing.F) {
	seeds := []string{
		"test_spoke:\n  enabled: true",
		"test_spoke:\n  coverage_threshold: 80",
		"test_spoke:\n  coverage_threshold: -10",
		"test_spoke:\n  coverage_threshold: 150",
		"squeeze:\n  max_cpu_percent: 80",
		"squeeze:\n  max_cpu_percent: 101",
		"squeeze:\n  max_memory_mb: -5",
		"models:\n  encoder: ",
		"models:\n  decoder: ",
		"models:\n  fast: ",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data string) {
		cfg, _ := config.Parse([]byte(data))
		if cfg == nil {
			return
		}

		// Validate should never panic
		err := config.Validate(cfg)
		_ = err // We don't care about result, just no panic
	})
}
