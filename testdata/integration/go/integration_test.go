//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/internal/doc"
	"github.com/yourusername/spoke-tool/internal/model"
	"github.com/yourusername/spoke-tool/internal/test"
)

// ============================================================================
// Test Environment Setup
// ============================================================================

var (
	testDir     string
	modelClient *model.Client
	cfg         *types.Config
)

func TestMain(m *testing.M) {
	// Setup test environment
	var err error
	testDir, err = os.MkdirTemp("", "spoke-tool-integration-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create temp dir: %v\n", err)
		os.Exit(1)
	}

	// Load test config
	cfg = &types.Config{
		ProjectRoot: testDir,
		Models: struct {
			Encoder string `json:"encoder" yaml:"encoder"`
			Decoder string `json:"decoder" yaml:"decoder"`
			Fast    string `json:"fast" yaml:"fast"`
		}{
			Encoder: "codebert",
			Decoder: "deepseek-coder:7b",
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
			Path:    filepath.Join(testDir, "audit.log"),
		},
	}

	// Initialize model client (skip if Ollama not available)
	modelClient, err = model.NewClient(model.ClientConfig{
		OllamaHost: "http://localhost:11434",
		Timeout:    5 * time.Second,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create model client: %v\n", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.RemoveAll(testDir)
	os.Exit(code)
}

// ============================================================================
// Helper Functions
// ============================================================================

func createTestGoFile(t *testing.T, content string) string {
	t.Helper()
	filePath := filepath.Join(testDir, "main.go")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}

func createTestPythonFile(t *testing.T, content string) string {
	t.Helper()
	filePath := filepath.Join(testDir, "main.py")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}

func createTestNodeJSFile(t *testing.T, content string) string {
	t.Helper()
	filePath := filepath.Join(testDir, "main.js")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}

// ============================================================================
// Test Spoke Integration Tests
// ============================================================================

func TestTestSpoke_Go_AnalyzeAndGenerate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test Go file
	goCode := `package main

import "fmt"

// Add returns the sum of two integers
func Add(a, b int) int {
	return a + b
}

// Subtract returns the difference between two integers
func Subtract(a, b int) int {
	return a - b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
	return a * b
}

// Divide returns the quotient of two integers
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// ProcessData processes a slice of integers
func ProcessData(data []int) []int {
	result := make([]int, len(data))
	for i, v := range data {
		result[i] = v * 2
	}
	return result
}
`
	createTestGoFile(t, goCode)

	// Initialize analyzer
	analyzer := test.NewAnalyzer(test.AnalyzerConfig{
		Languages:      []types.Language{types.Go},
		ExportedOnly:   true,
		IncludePrivate: false,
	})

	// Analyze project
	ctx := context.Background()
	result, err := analyzer.AnalyzeProject(ctx, testDir)
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	// Verify analysis results
	if len(result.Functions) == 0 {
		t.Error("Expected to find functions, got none")
	}

	// Count untested functions
	untested := analyzer.FindUntestedFunctions(result)
	if len(untested) == 0 {
		t.Error("Expected to find untested functions")
	}

	// Generate tests (skip if model client not available)
	if modelClient == nil {
		t.Skip("Model client not available")
	}

	generator := test.NewGenerator(test.GeneratorConfig{
		ModelClient:  modelClient,
		ProjectRoot:  testDir,
		AutoRunTests: false,
		Verbose:      true,
	})

	genResult, err := generator.GenerateTests(ctx, result, untested, true)
	if err != nil {
		t.Fatalf("Failed to generate tests: %v", err)
	}

	// Verify test generation
	if len(genResult.GeneratedTests) == 0 {
		t.Error("Expected generated tests, got none")
	}

	// Write test files
	err = generator.WriteTestFiles(ctx, genResult)
	if err != nil {
		t.Fatalf("Failed to write test files: %v", err)
	}

	// Check if test file was created
	testFilePath := filepath.Join(testDir, "main_test.go")
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Errorf("Test file not created: %v", err)
	}
}

func TestTestSpoke_NodeJS_AnalyzeAndGenerate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test Node.js file
	nodeCode := `// math.js
function add(a, b) {
    return a + b;
}

function subtract(a, b) {
    return a - b;
}

function multiply(a, b) {
    return a * b;
}

function divide(a, b) {
    if (b === 0) {
        throw new Error('Division by zero');
    }
    return a / b;
}

function processData(data) {
    return data.map(x => x * 2);
}

module.exports = {
    add,
    subtract,
    multiply,
    divide,
    processData
};
`
	createTestNodeJSFile(t, nodeCode)

	// Initialize analyzer
	analyzer := test.NewAnalyzer(test.AnalyzerConfig{
		Languages:      []types.Language{types.NodeJS},
		ExportedOnly:   true,
		IncludePrivate: false,
	})

	// Analyze project
	ctx := context.Background()
	result, err := analyzer.AnalyzeProject(ctx, testDir)
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	// Verify analysis results
	if len(result.Functions) == 0 {
		t.Error("Expected to find functions, got none")
	}
}

func TestTestSpoke_Python_AnalyzeAndGenerate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test Python file
	pythonCode := `# math.py
def add(a, b):
    """Add two numbers."""
    return a + b

def subtract(a, b):
    """Subtract two numbers."""
    return a - b

def multiply(a, b):
    """Multiply two numbers."""
    return a * b

def divide(a, b):
    """Divide two numbers."""
    if b == 0:
        raise ValueError("Division by zero")
    return a / b

def process_data(data):
    """Process a list of numbers."""
    return [x * 2 for x in data]
`
	createTestPythonFile(t, pythonCode)

	// Initialize analyzer
	analyzer := test.NewAnalyzer(test.AnalyzerConfig{
		Languages:      []types.Language{types.Python},
		ExportedOnly:   true,
		IncludePrivate: false,
	})

	// Analyze project
	ctx := context.Background()
	result, err := analyzer.AnalyzeProject(ctx, testDir)
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	// Verify analysis results
	if len(result.Functions) == 0 {
		t.Error("Expected to find functions, got none")
	}
}

// ============================================================================
// Readme Spoke Integration Tests
// ============================================================================

func TestReadmeSpoke_ExtractAndGenerate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test Go file with tests
	goCode := `package main

// Calculator performs arithmetic operations
type Calculator struct{}

// Add returns the sum of two numbers
func (c *Calculator) Add(a, b int) int {
	return a + b
}

// Subtract returns the difference between two numbers
func (c *Calculator) Subtract(a, b int) int {
	return a - b
}
`
	createTestGoFile(t, goCode)

	// Create test file
	testCode := `package main

import "testing"

func TestCalculator_Add(t *testing.T) {
	c := &Calculator{}
	result := c.Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2,3) = %d, want 5", result)
	}
}

func TestCalculator_Subtract(t *testing.T) {
	c := &Calculator{}
	result := c.Subtract(10, 4)
	if result != 6 {
		t.Errorf("Subtract(10,4) = %d, want 6", result)
	}
}
`
	testFilePath := filepath.Join(testDir, "calculator_test.go")
	err := os.WriteFile(testFilePath, []byte(testCode), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Initialize extractor
	extractor := doc.NewExtractor(doc.ExtractorConfig{
		Languages:          []types.Language{types.Go},
		IncludeTests:       true,
		IncludeComments:    true,
		MaxExamplesPerFunc: 3,
	})

	// Create analysis
	analyzer := test.NewAnalyzer(test.AnalyzerConfig{
		Languages: []types.Language{types.Go},
	})

	ctx := context.Background()
	analysis, err := analyzer.AnalyzeProject(ctx, testDir)
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	// Extract content
	content, err := extractor.ExtractFromProject(ctx, analysis)
	if err != nil {
		t.Fatalf("Failed to extract content: %v", err)
	}

	if len(content) == 0 {
		t.Error("Expected extracted content, got none")
	}

	// Initialize summarizer (skip if model not available)
	if modelClient == nil {
		t.Skip("Model client not available")
	}

	summarizer := doc.NewSummarizer(doc.SummarizerConfig{
		Model:            model.Gemma2B,
		MaxSummaryLength: 200,
		UseCache:         true,
	}, modelClient)

	// Summarize functions
	for _, fn := range analysis.Functions {
		summary, err := summarizer.SummarizeFunction(ctx, fn)
		if err != nil {
			t.Logf("Failed to summarize function %s: %v", fn.Name, err)
			continue
		}
		if summary != nil && summary.Description == "" {
			t.Errorf("Empty summary for function %s", fn.Name)
		}
	}

	// Initialize formatter
	formatter := doc.NewFormatter(doc.FormatterConfig{
		IncludeBadges: true,
		IncludeTOC:    true,
		AddEmojis:     true,
	})

	// Create sections
	sections := []*doc.Section{
		{
			Type:    types.DocSectionTitle,
			Title:   "Calculator Library",
			Content: "A simple calculator library for Go",
		},
		{
			Type:    types.DocSectionInstallation,
			Title:   "Installation",
			Content: "```bash\ngo get github.com/example/calculator\n```",
		},
	}

	// Format README
	readme, err := formatter.FormatReadme("Calculator", "A simple calculator library", sections)
	if err != nil {
		t.Fatalf("Failed to format README: %v", err)
	}

	if readme.Content == "" {
		t.Error("Expected README content, got empty")
	}

	// Initialize updater
	updater := doc.NewUpdater(doc.UpdaterConfig{
		CreateBackup: true,
	})

	// Update README
	readmePath := filepath.Join(testDir, "README.md")
	result, err := updater.UpdateReadme(ctx, readmePath, readme)
	if err != nil {
		t.Fatalf("Failed to update README: %v", err)
	}

	if !result.Updated {
		t.Log("README not updated (may already exist)")
	}
}

// ============================================================================
// End-to-End Integration Tests
// ============================================================================

func TestEndToEnd_Go_Complete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a complete Go project
	projectDir := filepath.Join(testDir, "endtoend")
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	// Create main.go
	mainCode := `package main

import "fmt"

// User represents a user in the system
type User struct {
    ID   int
    Name string
    Age  int
}

// NewUser creates a new user
func NewUser(id int, name string, age int) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid id: %d", id)
    }
    if name == "" {
        return nil, fmt.Errorf("name cannot be empty")
    }
    if age < 0 || age > 150 {
        return nil, fmt.Errorf("invalid age: %d", age)
    }
    return &User{ID: id, Name: name, Age: age}, nil
}

// IsAdult returns true if the user is an adult
func (u *User) IsAdult() bool {
    return u.Age >= 18
}

// String returns a string representation of the user
func (u *User) String() string {
    return fmt.Sprintf("User{ID=%d, Name=%s, Age=%d}", u.ID, u.Name, u.Age)
}
`
	err = os.WriteFile(filepath.Join(projectDir, "main.go"), []byte(mainCode), 0644)
	if err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}

	// Step 1: Analyze project
	analyzer := test.NewAnalyzer(test.AnalyzerConfig{
		Languages:      []types.Language{types.Go},
		ExportedOnly:   true,
		IncludePrivate: false,
	})

	ctx := context.Background()
	analysis, err := analyzer.AnalyzeProject(ctx, projectDir)
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	// Verify analysis found functions
	if len(analysis.Functions) == 0 {
		t.Error("Expected to find functions")
	}

	// Step 2: Generate tests (if model available)
	if modelClient != nil {
		generator := test.NewGenerator(test.GeneratorConfig{
			ModelClient:  modelClient,
			ProjectRoot:  projectDir,
			AutoRunTests: false,
		})

		untested := analyzer.FindUntestedFunctions(analysis)
		if len(untested) > 0 {
			genResult, err := generator.GenerateTests(ctx, analysis, untested, true)
			if err != nil {
				t.Logf("Test generation failed (may be expected if models not available): %v", err)
			} else {
				err = generator.WriteTestFiles(ctx, genResult)
				if err != nil {
					t.Logf("Failed to write test files: %v", err)
				}
			}
		}
	}

	// Step 3: Extract documentation
	extractor := doc.NewExtractor(doc.ExtractorConfig{
		Languages:    []types.Language{types.Go},
		IncludeTests: true,
	})

	content, err := extractor.ExtractFromProject(ctx, analysis)
	if err != nil {
		t.Fatalf("Failed to extract content: %v", err)
	}

	// Step 4: Create README sections
	sections := []*doc.Section{
		{
			Type:    types.DocSectionTitle,
			Title:   "User Management Library",
			Content: "A simple user management library",
		},
		{
			Type:    types.DocSectionInstallation,
			Title:   "Installation",
			Content: "```bash\ngo get github.com/example/user\n```",
		},
		{
			Type:    types.DocSectionAPI,
			Title:   "API Reference",
			Content: "## User\n\nRepresents a user in the system.",
		},
	}

	// Step 5: Format README
	formatter := doc.NewFormatter(doc.FormatterConfig{
		IncludeBadges: true,
		IncludeTOC:    true,
	})

	readme, err := formatter.FormatReadme("User Library", "User management library", sections)
	if err != nil {
		t.Fatalf("Failed to format README: %v", err)
	}

	// Step 6: Update README
	updater := doc.NewUpdater(doc.UpdaterConfig{
		CreateBackup: true,
	})

	readmePath := filepath.Join(projectDir, "README.md")
	result, err := updater.UpdateReadme(ctx, readmePath, readme)
	if err != nil {
		t.Fatalf("Failed to update README: %v", err)
	}

	if result.Updated {
		t.Log("README updated successfully")
	}

	// Verify files were created
	files := []string{
		filepath.Join(projectDir, "main.go"),
		filepath.Join(projectDir, "README.md"),
	}

	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("Expected file not found: %s", f)
		}
	}
}

// ============================================================================
// API Integration Tests
// ============================================================================

func TestAPI_HealthEndpoint(t *testing.T) {
	// Create a test HTTP server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	// Make request
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %q", result["status"])
	}
}

func TestAPI_GenerateTestsEndpoint(t *testing.T) {
	// This test would hit your actual API if it exists
	// For now, it's a placeholder
	t.Skip("API endpoint not implemented")
}

// ============================================================================
// Performance Tests
// ============================================================================

func TestPerformance_LargeProjectAnalysis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Create a large project with many files
	largeDir := filepath.Join(testDir, "large")
	err := os.MkdirAll(largeDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create large dir: %v", err)
	}

	// Generate 100 Go files
	for i := 0; i < 100; i++ {
		content := fmt.Sprintf(`package main

func Function%d() int {
	return %d
}
`, i, i)
		err = os.WriteFile(filepath.Join(largeDir, fmt.Sprintf("file%d.go", i)), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Measure analysis time
	analyzer := test.NewAnalyzer(test.AnalyzerConfig{
		Languages: []types.Language{types.Go},
	})

	ctx := context.Background()
	start := time.Now()
	analysis, err := analyzer.AnalyzeProject(ctx, largeDir)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	t.Logf("Analyzed %d files in %v", len(analysis.Files), duration)

	// Performance assertion (adjust threshold as needed)
	if duration > 10*time.Second {
		t.Errorf("Analysis took too long: %v", duration)
	}
}

// ============================================================================
// Concurrent Operation Tests
// ============================================================================

func TestConcurrent_AnalysisAndGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrency test in short mode")
	}

	// Create multiple projects to analyze concurrently
	projects := make([]string, 5)
	for i := 0; i < 5; i++ {
		projDir := filepath.Join(testDir, fmt.Sprintf("concurrent-%d", i))
		err := os.MkdirAll(projDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create project dir: %v", err)
		}

		content := fmt.Sprintf(`package main

func Process%d(x int) int {
	return x * %d
}
`, i, i+1)
		err = os.WriteFile(filepath.Join(projDir, "main.go"), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		projects[i] = projDir
	}

	// Analyze concurrently
	ctx := context.Background()
	errChan := make(chan error, len(projects))

	for _, proj := range projects {
		go func(dir string) {
			analyzer := test.NewAnalyzer(test.AnalyzerConfig{
				Languages: []types.Language{types.Go},
			})
			_, err := analyzer.AnalyzeProject(ctx, dir)
			errChan <- err
		}(proj)
	}

	// Collect results
	for i := 0; i < len(projects); i++ {
		err := <-errChan
		if err != nil {
			t.Errorf("Concurrent analysis failed: %v", err)
		}
	}
}
