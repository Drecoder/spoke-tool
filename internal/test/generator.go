package test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"example.com/spoke-tool/api/types"
	"example.com/spoke-tool/cmd/shared"
	"example.com/spoke-tool/internal/common"
	"example.com/spoke-tool/internal/model"
)

// Generator handles test generation using SLMs
type Generator struct {
	config      GeneratorConfig
	modelClient *model.Client
	fileUtils   *common.FileUtils
	stringUtils *common.StringUtils
	logger      *common.Logger
}

// GeneratorConfig configures the test generator
type GeneratorConfig struct {
	ModelClient       *model.Client
	ProjectRoot       string
	AutoRunTests      bool
	CheckCoverage     bool
	CoverageThreshold float64
	TargetLanguage    types.Language
	Verbose           bool
	MaxTestsPerFunc   int
	IncludeEdgeCases  bool
	GenerateMocks     bool
}

// GenerationResult represents the result of test generation
type GenerationResult struct {
	Success         bool                `json:"success"`
	Message         string              `json:"message"`
	ExitCode        shared.ExitCode     `json:"exit_code"`
	GeneratedTests  []*GeneratedTest    `json:"generated_tests"`
	TestsGenerated  int                 `json:"tests_generated"`
	TestResults     *types.TestSuite    `json:"test_results,omitempty"`
	Coverage        *types.TestCoverage `json:"coverage,omitempty"`
	FunctionsTested int                 `json:"functions_tested"`
	ModelsQueried   int                 `json:"models_queried"`
	Errors          []string            `json:"errors,omitempty"`
}

// GeneratedTest represents a single generated test
type GeneratedTest struct {
	FunctionName string         `json:"function_name"`
	Language     types.Language `json:"language"`
	TestCode     string         `json:"test_code"`
	TestFilePath string         `json:"test_file_path"`
	Framework    string         `json:"framework"`
	Confidence   float64        `json:"confidence"`
}

// NewGenerator creates a new test generator
func NewGenerator(config GeneratorConfig) *Generator {
	if config.MaxTestsPerFunc == 0 {
		config.MaxTestsPerFunc = 5
	}

	return &Generator{
		config:      config,
		modelClient: config.ModelClient,
		fileUtils:   &common.FileUtils{},
		stringUtils: &common.StringUtils{},
		logger:      common.GetLogger().WithField("component", "test-generator"),
	}
}

// GenerateTests generates tests for the specified functions
func (g *Generator) GenerateTests(ctx context.Context, analysis *types.CodeAnalysis, functions []*types.Function, force bool) (*GenerationResult, error) {
	g.logger.Info("Generating tests", "functions", len(functions))

	result := &GenerationResult{
		Success:        true,
		GeneratedTests: []*GeneratedTest{},
		Errors:         []string{},
	}

	// Group functions by language for batch processing
	byLanguage := make(map[types.Language][]*types.Function)
	for _, fn := range functions {
		byLanguage[fn.Language] = append(byLanguage[fn.Language], fn)
	}

	for lang, funcs := range byLanguage {
		// Skip if we're targeting a specific language
		if g.config.TargetLanguage != "" && g.config.TargetLanguage != lang {
			continue
		}

		g.logger.Info("Generating tests for language", "language", lang, "count", len(funcs))

		// Get dependencies for this language
		deps := g.getDependencies(analysis, lang)

		// Generate tests for each function
		for _, fn := range funcs {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			// Skip if function already has tests and we're not forcing
			if fn.HasTest && !force {
				g.logger.Debug("Skipping function with existing tests", "function", fn.Name)
				continue
			}

			// Generate test for this function
			test, err := g.generateForFunction(ctx, fn, deps)
			if err != nil {
				g.logger.Error("Failed to generate test", "function", fn.Name, "error", err)
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", fn.Name, err))
				continue
			}

			if test != nil {
				result.GeneratedTests = append(result.GeneratedTests, test)
				result.FunctionsTested++
				result.ModelsQueried++
			}
		}
	}

	if len(result.Errors) > 0 {
		result.Success = false
		result.Message = fmt.Sprintf("Generated %d tests with %d errors",
			len(result.GeneratedTests), len(result.Errors))
	} else {
		result.Message = fmt.Sprintf("Successfully generated %d tests", len(result.GeneratedTests))
	}

	g.logger.Info("Test generation complete",
		"generated", len(result.GeneratedTests),
		"errors", len(result.Errors))

	return result, nil
}

// generateForFunction generates a test for a single function
func (g *Generator) generateForFunction(ctx context.Context, fn *types.Function, deps string) (*GeneratedTest, error) {
	g.logger.Debug("Generating test for function", "function", fn.Name)

	// Build the prompt based on language
	prompt := g.buildTestPrompt(fn, deps)

	// Generate test code using the model
	resp, err := g.modelClient.Generate(ctx, model.SLMRequest{
		Model:       model.DeepSeek7B, // Use DeepSeek for test generation
		Language:    fn.Language,
		Prompt:      prompt,
		Temperature: 0.3, // Lower temperature for more consistent tests
		MaxTokens:   2048,
	})
	if err != nil {
		return nil, fmt.Errorf("model generation failed: %w", err)
	}

	// Parse the generated test code
	testCode := g.cleanGeneratedTest(resp.Response, fn.Language)
	if testCode == "" {
		return nil, fmt.Errorf("generated empty test code")
	}

	// Determine test file path
	testFilePath := g.getTestFilePath(fn)

	// Determine test framework
	framework := g.getTestFramework(fn.Language)

	// Calculate confidence (simplified)
	confidence := g.calculateConfidence(resp)

	return &GeneratedTest{
		FunctionName: fn.Name,
		Language:     fn.Language,
		TestCode:     testCode,
		TestFilePath: testFilePath,
		Framework:    framework,
		Confidence:   confidence,
	}, nil
}

// buildTestPrompt creates a language-specific prompt for test generation
func (g *Generator) buildTestPrompt(fn *types.Function, deps string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Generate unit tests for the following %s function.\n", fn.Language))
	sb.WriteString("DO NOT modify the original code - only create tests.\n\n")
	sb.WriteString(fmt.Sprintf("Function: %s\n", fn.Name))
	sb.WriteString(fmt.Sprintf("Signature: %s\n", fn.Signature))
	sb.WriteString("\nCode:\n")
	sb.WriteString(fn.Content)
	sb.WriteString("\n")

	if deps != "" {
		sb.WriteString(fmt.Sprintf("\nDependencies: %s\n", deps))
	}

	// Language-specific instructions
	switch fn.Language {
	case types.Go:
		sb.WriteString("\nCreate tests using the testing package that verify:\n")
		sb.WriteString("- Happy path (normal inputs)\n")
		sb.WriteString("- Edge cases (zero values, boundaries)\n")
		sb.WriteString("- Error conditions\n")
		sb.WriteString("- Table-driven tests where appropriate\n")
		sb.WriteString("\nReturn ONLY the test code, no explanations.")

	case types.NodeJS:
		sb.WriteString("\nCreate tests using Jest that verify:\n")
		sb.WriteString("- Happy path (normal inputs)\n")
		sb.WriteString("- Edge cases (null, undefined, boundaries)\n")
		sb.WriteString("- Error conditions\n")
		sb.WriteString("- Async behavior if applicable\n")
		sb.WriteString("- Mocks for dependencies\n")
		sb.WriteString("\nReturn ONLY the test code, no explanations.")

	case types.Python:
		sb.WriteString("\nCreate tests using pytest that verify:\n")
		sb.WriteString("- Happy path (normal inputs)\n")
		sb.WriteString("- Edge cases (None, empty, boundaries)\n")
		sb.WriteString("- Error conditions (exceptions)\n")
		sb.WriteString("- Fixtures for setup\n")
		sb.WriteString("- Mocks for dependencies\n")
		sb.WriteString("\nReturn ONLY the test code, no explanations.")
	}

	return sb.String()
}

// cleanGeneratedTest removes any explanatory text and ensures just the test code
func (g *Generator) cleanGeneratedTest(response string, lang types.Language) string {
	lines := strings.Split(response, "\n")
	var cleanLines []string
	inCode := false

	for _, line := range lines {
		// Look for code block markers
		if strings.Contains(line, "```") {
			inCode = !inCode
			continue
		}

		// If we're in a code block, keep the line
		if inCode {
			cleanLines = append(cleanLines, line)
			continue
		}

		// If not in a code block, try to detect test code
		if g.looksLikeTestCode(line, lang) {
			cleanLines = append(cleanLines, line)
		}
	}

	// If we didn't find any code blocks, use the whole response
	if len(cleanLines) == 0 {
		return strings.TrimSpace(response)
	}

	return strings.Join(cleanLines, "\n")
}

// looksLikeTestCode heuristically determines if a line looks like test code
func (g *Generator) looksLikeTestCode(line string, lang types.Language) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}

	switch lang {
	case types.Go:
		return strings.HasPrefix(line, "func Test") ||
			strings.Contains(line, "t *testing.T") ||
			strings.Contains(line, "t.Run(") ||
			strings.Contains(line, "t.Errorf") ||
			strings.Contains(line, "assert.")

	case types.NodeJS:
		return strings.HasPrefix(line, "test(") ||
			strings.HasPrefix(line, "it(") ||
			strings.HasPrefix(line, "describe(") ||
			strings.Contains(line, "expect(") ||
			strings.Contains(line, ".toBe(") ||
			strings.Contains(line, ".toEqual(")

	case types.Python:
		return strings.HasPrefix(line, "def test_") ||
			strings.Contains(line, "assert ") ||
			strings.Contains(line, "pytest.raises") ||
			strings.Contains(line, "@pytest.mark")

	default:
		return false
	}
}

// getTestFilePath determines where to write the test file
func (g *Generator) getTestFilePath(fn *types.Function) string {
	switch fn.Language {
	case types.Go:
		// For Go, test files go in the same package with _test.go suffix
		if strings.HasSuffix(fn.FilePath, ".go") {
			return strings.TrimSuffix(fn.FilePath, ".go") + "_test.go"
		}
		return fn.FilePath + "_test.go"

	case types.NodeJS:
		// For Node.js, test files go alongside with .test.js suffix
		if strings.HasSuffix(fn.FilePath, ".js") {
			return strings.TrimSuffix(fn.FilePath, ".js") + ".test.js"
		}
		if strings.HasSuffix(fn.FilePath, ".ts") {
			return strings.TrimSuffix(fn.FilePath, ".ts") + ".test.ts"
		}
		return fn.FilePath + ".test.js"

	case types.Python:
		// For Python, test files go in same dir with test_ prefix
		dir := fn.FilePath[:strings.LastIndex(fn.FilePath, "/")]
		base := fn.FilePath[strings.LastIndex(fn.FilePath, "/")+1:]
		return dir + "/test_" + base

	default:
		return fn.FilePath + "_test"
	}
}

// getTestFramework returns the framework name for the language
func (g *Generator) getTestFramework(lang types.Language) string {
	switch lang {
	case types.Go:
		return "go test"
	case types.NodeJS:
		return "jest"
	case types.Python:
		return "pytest"
	default:
		return "unknown"
	}
}

// getDependencies extracts dependencies for the language
func (g *Generator) getDependencies(analysis *types.CodeAnalysis, lang types.Language) string {
	var deps []string

	for _, imp := range analysis.Imports {
		// Filter by language if needed
		deps = append(deps, imp.Path)
	}

	return strings.Join(deps, ", ")
}

// calculateConfidence estimates how good the generated test is
func (g *Generator) calculateConfidence(resp *model.Response) float64 {
	// Start with base confidence from model
	confidence := 0.7

	// Adjust based on response length (very short responses are suspect)
	if len(resp.Response) < 50 {
		confidence -= 0.2
	}

	// Check for indicators of good tests
	response := strings.ToLower(resp.Response)
	if strings.Contains(response, "test") {
		confidence += 0.1
	}
	if strings.Contains(response, "assert") || strings.Contains(response, "expect") {
		confidence += 0.1
	}
	if strings.Contains(response, "error") {
		confidence += 0.1
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// WriteTestFiles writes the generated tests to disk
func (g *Generator) WriteTestFiles(ctx context.Context, result *GenerationResult) error {
	g.logger.Info("Writing test files", "count", len(result.GeneratedTests))

	for _, test := range result.GeneratedTests {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Check if test file already exists
		exists := g.fileUtils.FileExists(test.TestFilePath)

		if exists {
			// Append to existing test file
			if err := g.appendToTestFile(test); err != nil {
				g.logger.Error("Failed to append to test file", "file", test.TestFilePath, "error", err)
				result.Errors = append(result.Errors, fmt.Sprintf("append %s: %v", test.TestFilePath, err))
				continue
			}
			g.logger.Debug("Appended to test file", "file", test.TestFilePath)
		} else {
			// Create new test file
			if err := g.createTestFile(test); err != nil {
				g.logger.Error("Failed to create test file", "file", test.TestFilePath, "error", err)
				result.Errors = append(result.Errors, fmt.Sprintf("create %s: %v", test.TestFilePath, err))
				continue
			}
			g.logger.Debug("Created test file", "file", test.TestFilePath)
		}
	}

	return nil
}

// createTestFile creates a new test file
func (g *Generator) createTestFile(test *GeneratedTest) error {
	// Add package/import headers based on language
	header := g.getTestFileHeader(test)
	content := header + "\n\n" + test.TestCode

	return g.fileUtils.WriteFile(test.TestFilePath, content)
}

// appendToTestFile appends tests to an existing file
func (g *Generator) appendToTestFile(test *GeneratedTest) error {
	// Read existing content
	content, err := g.fileUtils.ReadFile(test.TestFilePath)
	if err != nil {
		return err
	}

	// Add a separator and the new test
	newContent := content + "\n\n" + test.TestCode

	return g.fileUtils.WriteFile(test.TestFilePath, newContent)
}

// getTestFileHeader returns the appropriate file header based on language
func (g *Generator) getTestFileHeader(test *GeneratedTest) string {
	switch test.Language {
	case types.Go:
		return fmt.Sprintf(`package %s_test

import (
	"testing"
)`, test.FunctionName)

	case types.NodeJS:
		return `// Generated tests
const { ` + test.FunctionName + ` } = require('./` + test.FunctionName + `');

describe('` + test.FunctionName + `', () => {`

	case types.Python:
		return `"""Generated tests for ` + test.FunctionName + `."""
import pytest
from .` + test.FunctionName + ` import ` + test.FunctionName + `
`

	default:
		return ""
	}
}

// Config returns the generator configuration
func (g *Generator) Config() GeneratorConfig {
	return g.config
}

// FindUntestedFunctions identifies functions without tests
func (g *Generator) FindUntestedFunctions(analysis *types.CodeAnalysis) []*types.Function {
	var untested []*types.Function
	for _, fn := range analysis.Functions {
		if !fn.HasTest {
			untested = append(untested, &fn)
		}
	}
	return untested
}

// AnalyzeProject analyzes the project structure
func (g *Generator) AnalyzeProject(ctx context.Context) (*types.CodeAnalysis, error) {
	// Placeholder - you'll implement actual analysis later
	return &types.CodeAnalysis{
		Files:     []types.CodeFile{},
		Functions: []types.Function{},
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// CheckCoverage checks test coverage
func (g *Generator) CheckCoverage(ctx context.Context) (*types.TestCoverage, error) {
	return &types.TestCoverage{
		Overall: 0.0,
		ByFile:  make(map[string]float64),
	}, nil
}

// AnalyzeFailure analyzes a test failure
func (g *Generator) AnalyzeFailure(ctx context.Context, test *types.TestResult) (string, error) {
	return "Test failure analysis placeholder", nil
}

// RunTests runs the generated tests
func (g *Generator) RunTests(ctx context.Context, result *GenerationResult) (*types.TestSuite, error) {
	return &types.TestSuite{
		Total:   0,
		Passed:  0,
		Failed:  0,
		Results: []types.TestResult{},
	}, nil
}
