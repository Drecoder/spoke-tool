package test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
		Model:       model.CodeLLamaDecoder, // Use DeepSeek for test generation
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
		// Go tests stay in the same directory with _test.go suffix
		if strings.HasSuffix(fn.FilePath, ".go") {
			return strings.TrimSuffix(fn.FilePath, ".go") + "_test.go"
		}
		return fn.FilePath + "_test.go"

	case types.NodeJS:
		// Get the relative path from project root
		relPath, err := filepath.Rel(g.config.ProjectRoot, fn.FilePath)
		if err != nil {
			g.logger.Debug("Failed to get relative path, using fallback", "error", err)
			return g.fallbackNodeTestPath(fn)
		}

		// Create test path in __tests__ directory mirroring the structure
		testDir := filepath.Join(g.config.ProjectRoot, "__tests__", filepath.Dir(relPath))
		baseName := filepath.Base(fn.FilePath)

		// Change extension to .test.js or .test.ts
		ext := filepath.Ext(baseName)
		testBaseName := strings.TrimSuffix(baseName, ext) + ".test"

		// Determine proper extension for test file
		switch ext {
		case ".js", ".jsx", ".mjs", ".cjs":
			testBaseName += ".js"
		case ".ts", ".tsx":
			testBaseName += ".ts"
		default:
			testBaseName += ".js"
		}

		fullPath := filepath.Join(testDir, testBaseName)
		g.logger.Debug("Node.js test file path", "source", fn.FilePath, "test", fullPath)
		return fullPath

	case types.Python:
		// For Python, tests go in tests/ directory mirroring structure
		relPath, err := filepath.Rel(g.config.ProjectRoot, fn.FilePath)
		if err != nil {
			g.logger.Debug("Failed to get relative path, using fallback", "error", err)
			dir := filepath.Dir(fn.FilePath)
			base := filepath.Base(fn.FilePath)
			return filepath.Join(dir, "test_"+base)
		}

		testDir := filepath.Join(g.config.ProjectRoot, "tests", filepath.Dir(relPath))
		baseName := filepath.Base(fn.FilePath)
		testBaseName := "test_" + baseName

		fullPath := filepath.Join(testDir, testBaseName)
		g.logger.Debug("Python test file path", "source", fn.FilePath, "test", fullPath)
		return fullPath

	default:
		return fn.FilePath + "_test"
	}
}

// fallbackNodeTestPath provides a fallback for Node.js when relative path fails
func (g *Generator) fallbackNodeTestPath(fn *types.Function) string {
	dir := filepath.Dir(fn.FilePath)
	base := filepath.Base(fn.FilePath)
	ext := filepath.Ext(base)

	// Check if we're already in a __tests__ directory
	if strings.Contains(dir, "__tests__") {
		// Just change the extension
		return strings.TrimSuffix(fn.FilePath, ext) + ".test" + ext
	}

	// Put test in __tests__ subdirectory
	testDir := filepath.Join(dir, "__tests__")
	testBase := strings.TrimSuffix(base, ext) + ".test" + ext
	return filepath.Join(testDir, testBase)
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

		// Create directory if it doesn't exist
		dir := filepath.Dir(test.TestFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			g.logger.Error("Failed to create directory", "dir", dir, "error", err)
			result.Errors = append(result.Errors, fmt.Sprintf("mkdir %s: %v", dir, err))
			continue
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
		// Extract package name from file path or use function name
		pkgName := "main"
		dir := filepath.Dir(test.TestFilePath)
		if strings.Contains(dir, "/") {
			parts := strings.Split(dir, "/")
			pkgName = parts[len(parts)-1]
		}
		return fmt.Sprintf(`package %s_test

import (
	"testing"
)`, pkgName)

	case types.NodeJS:
		// For Node.js, we need to import the module correctly
		// Get the relative path from the test file to the source file
		sourceDir := filepath.Dir(strings.Replace(test.TestFilePath, "__tests__", "", 1))
		sourceFile := strings.TrimSuffix(filepath.Base(test.TestFilePath), ".test.js")
		if strings.HasSuffix(test.TestFilePath, ".test.ts") {
			sourceFile = strings.TrimSuffix(filepath.Base(test.TestFilePath), ".test.ts")
		}

		importPath := "./" + filepath.Join("..", filepath.Base(sourceDir), sourceFile)
		return `// Generated tests
const ` + sourceFile + ` = require('` + importPath + `');

describe('` + sourceFile + `', () => {`

	case types.Python:
		baseName := strings.TrimPrefix(filepath.Base(test.TestFilePath), "test_")
		baseName = strings.TrimSuffix(baseName, ".py")
		return `"""Generated tests for ` + baseName + `."""
import pytest
from ..` + baseName + ` import *

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
	g.logger.Info("Starting project analysis", "root", g.config.ProjectRoot)

	// SUPER OBVIOUS DEBUG
	fmt.Println("🚨🚨🚨 ANALYZE PROJECT WAS CALLED! 🚨🚨🚨")
	fmt.Printf("ProjectRoot: %q\n", g.config.ProjectRoot)

	// Walk the directory tree
	var files []types.CodeFile
	var functions []types.Function
	var imports []types.Import

	err := filepath.Walk(g.config.ProjectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories we should ignore
		if info.IsDir() {
			if g.shouldSkipDir(info.Name()) {
				g.logger.Debug("Skipping directory", "dir", info.Name())
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file extension matches supported languages
		lang := g.detectLanguage(path)
		if lang == "" {
			g.logger.Debug("Skipping unsupported file", "file", path)
			return nil
		}

		g.logger.Debug("Processing file", "file", path, "language", lang)

		// Read and parse file
		content, err := os.ReadFile(path)
		if err != nil {
			g.logger.Warn("Failed to read file", "file", path, "error", err)
			return nil
		}

		files = append(files, types.CodeFile{
			Path:     path,
			Language: lang,
			Content:  string(content),
		})

		// Extract functions
		extracted, err := g.extractFunctions(path, string(content), lang)
		if err != nil {
			g.logger.Warn("Failed to extract functions", "file", path, "error", err)
		} else {
			functions = append(functions, extracted...)
			g.logger.Debug("Extracted functions", "file", path, "count", len(extracted))
		}

		return nil
	})

	if err != nil {
		g.logger.Error("Analysis failed", "error", err)
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	g.logger.Info("Analysis complete",
		"files", len(files),
		"functions", len(functions))

	return &types.CodeAnalysis{
		Files:     files,
		Functions: functions,
		Imports:   imports,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// extractFunctions extracts function definitions from source code
func (g *Generator) extractFunctions(path string, content string, lang types.Language) ([]types.Function, error) {
	switch lang {
	case types.Go:
		return g.extractGoFunctions(path, content)
	case types.NodeJS:
		return g.extractNodeFunctions(path, content)
	case types.Python:
		return g.extractPythonFunctions(path, content)
	default:
		return []types.Function{}, nil
	}
}

// extractNodeFunctions extracts functions from Node.js/JavaScript/JSX files
func (g *Generator) extractNodeFunctions(path, content string) ([]types.Function, error) {
	// SUPER OBVIOUS DEBUG
	fmt.Printf("🔍 EXTRACTING NODE FUNCTIONS from: %s\n", path)
	fmt.Printf("   File size: %d bytes\n", len(content))
	fmt.Printf("   First 100 chars: %q\n", content[:min(100, len(content))])

	var functions []types.Function
	lines := strings.Split(content, "\n")
	fmt.Printf("   Lines: %d\n", len(lines))

	// Regular expressions for different function patterns
	funcDeclRegex := regexp.MustCompile(`function\s+(\w+)\s*\([^)]*\)\s*{`)
	arrowFuncRegex := regexp.MustCompile(`(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s*)?\([^)]*\)\s*=>\s*{?`)
	methodRegex := regexp.MustCompile(`(\w+)\s*\([^)]*\)\s*{`)
	exportFuncRegex := regexp.MustCompile(`export\s+(?:default\s+)?function\s+(\w+)\s*\(`)
	exportArrowRegex := regexp.MustCompile(`export\s+(?:default\s+)?(?:const|let|var)\s+(\w+)\s*=`)

	for i, line := range lines {
		lineNum := i + 1

		// Check each pattern
		patterns := []*regexp.Regexp{
			funcDeclRegex,
			arrowFuncRegex,
			exportFuncRegex,
			exportArrowRegex,
			methodRegex,
		}

		for _, pattern := range patterns {
			if matches := pattern.FindStringSubmatch(line); len(matches) > 1 {
				funcName := matches[1]

				// Skip if it's likely a false positive (like common keywords)
				if funcName == "if" || funcName == "for" || funcName == "while" {
					continue
				}

				fmt.Printf("   ✅ Found function at line %d: %s (pattern: %T)\n", lineNum, funcName, pattern)

				// Build signature (simplified)
				signature := line
				if len(line) > 100 {
					signature = line[:100] + "..."
				}

				// Check if function is exported (starts with uppercase or has export keyword)
				isExported := strings.Contains(line, "export") || (len(funcName) > 0 && funcName[0] >= 'A' && funcName[0] <= 'Z')

				// Check if test exists
				hasTest := g.checkNodeTestExists(path, funcName)

				fn := types.Function{
					Name:       funcName,
					Language:   types.NodeJS,
					FilePath:   path,
					Signature:  strings.TrimSpace(signature),
					Content:    line,
					LineStart:  lineNum,
					LineEnd:    lineNum,
					HasTest:    hasTest,
					IsExported: isExported,
				}

				functions = append(functions, fn)
				g.logger.Debug("Found Node.js function",
					"name", funcName,
					"file", path,
					"line", lineNum,
					"exported", isExported,
					"hasTest", hasTest)

				break // Only add once per line
			}
		}
	}

	fmt.Printf("   Total functions found: %d\n", len(functions))
	return functions, nil
}

// extractGoFunctions extracts functions from Go files
func (g *Generator) extractGoFunctions(path, content string) ([]types.Function, error) {
	var functions []types.Function
	lines := strings.Split(content, "\n")

	// Simple regex for Go function declarations
	// Note: This is simplified - real implementation should use go/parser
	funcRegex := regexp.MustCompile(`func\s+(\w+)\s*\([^)]*\)(?:\s*\(?[^{]*\)?)?\s*{`)

	for i, line := range lines {
		lineNum := i + 1

		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			funcName := matches[1]

			// Skip if it's a method receiver or test function
			if strings.Contains(line, "func (") || strings.HasPrefix(funcName, "Test") {
				continue
			}

			// Check if function is exported (starts with uppercase)
			isExported := len(funcName) > 0 && funcName[0] >= 'A' && funcName[0] <= 'Z'

			// Check if test exists
			hasTest := g.checkGoTestExists(path, funcName)

			fn := types.Function{
				Name:       funcName,
				Language:   types.Go,
				FilePath:   path,
				Signature:  strings.TrimSpace(line),
				Content:    line,
				LineStart:  lineNum,
				LineEnd:    lineNum,
				HasTest:    hasTest,
				IsExported: isExported,
			}

			functions = append(functions, fn)
			g.logger.Debug("Found Go function",
				"name", funcName,
				"file", path,
				"line", lineNum,
				"exported", isExported,
				"hasTest", hasTest)
		}
	}

	return functions, nil
}

// extractPythonFunctions extracts functions from Python files
func (g *Generator) extractPythonFunctions(path, content string) ([]types.Function, error) {
	var functions []types.Function
	lines := strings.Split(content, "\n")

	// Regex for Python function definitions
	funcRegex := regexp.MustCompile(`def\s+(\w+)\s*\([^)]*\):`)

	for i, line := range lines {
		lineNum := i + 1

		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			funcName := matches[1]

			// Skip if it's a private method (starts with _)
			isExported := !strings.HasPrefix(funcName, "_")

			// Check if test exists
			hasTest := g.checkPythonTestExists(path, funcName)

			fn := types.Function{
				Name:       funcName,
				Language:   types.Python,
				FilePath:   path,
				Signature:  strings.TrimSpace(line),
				Content:    line,
				LineStart:  lineNum,
				LineEnd:    lineNum,
				HasTest:    hasTest,
				IsExported: isExported,
			}

			functions = append(functions, fn)
			g.logger.Debug("Found Python function",
				"name", funcName,
				"file", path,
				"line", lineNum,
				"exported", isExported,
				"hasTest", hasTest)
		}
	}

	return functions, nil
}

// checkNodeTestExists checks if a test exists for a Node.js function
func (g *Generator) checkNodeTestExists(sourcePath, funcName string) bool {
	// Get the expected test file path using the same logic as getTestFilePath
	// For checking existence, we'll create a temporary function object
	tempFn := &types.Function{
		FilePath: sourcePath,
		Language: types.NodeJS,
		Name:     funcName,
	}
	expectedPath := g.getTestFilePath(tempFn)

	// Check if the file exists
	if _, err := os.Stat(expectedPath); err == nil {
		// File exists, check if it contains the function name (simplified)
		content, err := os.ReadFile(expectedPath)
		if err == nil && strings.Contains(string(content), funcName) {
			return true
		}
	}

	return false
}

// checkGoTestExists checks if a test exists for a Go function
func (g *Generator) checkGoTestExists(sourcePath, funcName string) bool {
	testPath := strings.TrimSuffix(sourcePath, ".go") + "_test.go"
	if _, err := os.Stat(testPath); err != nil {
		return false
	}

	// Check if test file contains a test for this function
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false
	}

	// Look for TestFuncName pattern
	return strings.Contains(string(content), "Test"+funcName)
}

// checkPythonTestExists checks if a test exists for a Python function
func (g *Generator) checkPythonTestExists(sourcePath, funcName string) bool {
	// Get the expected test file path using the same logic as getTestFilePath
	tempFn := &types.Function{
		FilePath: sourcePath,
		Language: types.Python,
		Name:     funcName,
	}
	expectedPath := g.getTestFilePath(tempFn)

	if _, err := os.Stat(expectedPath); err != nil {
		return false
	}

	// Check if test file contains a test for this function
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		return false
	}

	// Look for test_funcName pattern
	return strings.Contains(string(content), "test_"+funcName)
}

// Helper methods
func (g *Generator) shouldSkipDir(name string) bool {
	skipDirs := map[string]bool{
		"node_modules": true,
		".git":         true,
		"vendor":       true,
		"dist":         true,
		"build":        true,
		"__pycache__":  true,
		".venv":        true,
		"env":          true,
		"cache":        true,
		"coverage":     true,
		"test-results": true,
		"__tests__":    false, // Don't skip tests directory
		"tests":        false, // Don't skip tests directory
	}
	return skipDirs[name]
}

func (g *Generator) detectLanguage(path string) types.Language {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return types.Go
	case ".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs":
		return types.NodeJS
	case ".py", ".pyw", ".pyx":
		return types.Python
	default:
		return ""
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
