package test

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/internal/common"
)

// Analyzer analyzes code to find functions that need tests
// This component ONLY analyzes - it NEVER modifies code
type Analyzer struct {
	config      AnalyzerConfig
	fileUtils   *common.FileUtils
	stringUtils *common.StringUtils
	logger      *common.Logger
}

// AnalyzerConfig configures the analyzer
type AnalyzerConfig struct {
	// Languages to analyze
	Languages []types.Language

	// Whether to include exported functions only
	ExportedOnly bool

	// Whether to include private functions
	IncludePrivate bool

	// Minimum complexity to consider (1-10)
	MinComplexity int

	// File extensions to analyze by language
	Extensions map[types.Language][]string

	// Test file patterns by language
	TestPatterns map[types.Language][]string
}

// AnalysisResult represents the result of code analysis
type AnalysisResult struct {
	// Project root
	ProjectRoot string `json:"project_root"`

	// All files analyzed
	Files []*types.CodeFile `json:"files"`

	// All functions found
	Functions []*types.Function `json:"functions"`

	// Functions that need tests
	UntestedFunctions []*types.Function `json:"untested_functions"`

	// Test files found
	TestFiles []*types.TestFile `json:"test_files"`

	// Coverage information
	Coverage *TestCoverage `json:"coverage,omitempty"`

	// Statistics
	Stats *AnalysisStats `json:"stats"`
}

// AnalysisStats contains statistics about the analysis
type AnalysisStats struct {
	TotalFiles        int                              `json:"total_files"`
	SourceFiles       int                              `json:"source_files"`
	TestFiles         int                              `json:"test_files"`
	TotalFunctions    int                              `json:"total_functions"`
	TestedFunctions   int                              `json:"tested_functions"`
	UntestedFunctions int                              `json:"untested_functions"`
	ByLanguage        map[types.Language]LanguageStats `json:"by_language"`
}

// LanguageStats contains statistics for a specific language
type LanguageStats struct {
	Files             int `json:"files"`
	TestFiles         int `json:"test_files"`
	Functions         int `json:"functions"`
	TestedFunctions   int `json:"tested_functions"`
	UntestedFunctions int `json:"untested_functions"`
}

// TestCoverage represents test coverage information
type TestCoverage struct {
	Overall    float64            `json:"overall"`
	ByFile     map[string]float64 `json:"by_file"`
	ByFunction map[string]float64 `json:"by_function"`
	Uncovered  []string           `json:"uncovered"`
}

// NewAnalyzer creates a new code analyzer
func NewAnalyzer(config AnalyzerConfig) *Analyzer {
	// Set default extensions if not provided
	if config.Extensions == nil {
		config.Extensions = map[types.Language][]string{
			types.Go:     {".go"},
			types.NodeJS: {".js", ".ts", ".jsx", ".tsx"},
			types.Python: {".py"},
		}
	}

	// Set default test patterns if not provided
	if config.TestPatterns == nil {
		config.TestPatterns = map[types.Language][]string{
			types.Go:     {"*_test.go"},
			types.NodeJS: {"*.test.js", "*.spec.js", "*.test.ts", "*.spec.ts"},
			types.Python: {"test_*.py", "*_test.py"},
		}
	}

	return &Analyzer{
		config:      config,
		fileUtils:   &common.FileUtils{},
		stringUtils: &common.StringUtils{},
		logger:      common.GetLogger().WithField("component", "test-analyzer"),
	}
}

// AnalyzeProject analyzes an entire project for testing needs
func (a *Analyzer) AnalyzeProject(ctx context.Context, projectRoot string) (*AnalysisResult, error) {
	a.logger.Info("Analyzing project", "root", projectRoot)

	result := &AnalysisResult{
		ProjectRoot:       projectRoot,
		Files:             []*types.CodeFile{},
		Functions:         []*types.Function{},
		UntestedFunctions: []*types.Function{},
		TestFiles:         []*types.TestFile{},
		Stats: &AnalysisStats{
			ByLanguage: make(map[types.Language]LanguageStats),
		},
	}

	// Walk through all files in the project
	err := a.fileUtils.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip common directories to avoid
			if shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Detect language from file extension
		lang := a.detectLanguage(path)
		if lang == "" {
			return nil // Skip unknown languages
		}

		// Check if we should analyze this language
		if !a.shouldAnalyzeLanguage(lang) {
			return nil
		}

		// Check if this is a test file
		isTest := a.isTestFile(path, lang)

		// Read file content
		content, err := a.fileUtils.ReadFile(path)
		if err != nil {
			a.logger.Warn("Failed to read file", "path", path, "error", err)
			return nil
		}

		// Create code file object
		codeFile := &types.CodeFile{
			Path:     path,
			Language: lang,
			Content:  content,
			Hash:     common.Hashes.MD5(content),
		}

		if isTest {
			// This is a test file
			testFile := &types.TestFile{
				Path:      path,
				Language:  lang,
				Content:   content,
				Framework: a.detectTestFramework(path, content, lang),
			}
			result.TestFiles = append(result.TestFiles, testFile)
			result.Files = append(result.Files, codeFile)

			// Update stats
			stats := result.Stats.ByLanguage[lang]
			stats.TestFiles++
			result.Stats.ByLanguage[lang] = stats
		} else {
			// This is a source file - analyze for functions
			result.Files = append(result.Files, codeFile)

			// Parse functions from this file
			functions, err := a.parseFunctions(ctx, codeFile)
			if err != nil {
				a.logger.Warn("Failed to parse functions", "file", path, "error", err)
				return nil
			}

			result.Functions = append(result.Functions, functions...)

			// Update stats
			stats := result.Stats.ByLanguage[lang]
			stats.Files++
			stats.Functions += len(functions)
			result.Stats.ByLanguage[lang] = stats
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk project: %w", err)
	}

	// Map existing tests to functions
	a.mapTestsToFunctions(result)

	// Find untested functions
	a.findUntestedFunctions(result)

	// Calculate statistics
	a.calculateStats(result)

	a.logger.Info("Analysis complete",
		"files", len(result.Files),
		"functions", len(result.Functions),
		"untested", len(result.UntestedFunctions),
		"test_files", len(result.TestFiles))

	return result, nil
}

// AnalyzeFile analyzes a single file for testing needs
func (a *Analyzer) AnalyzeFile(ctx context.Context, path string) (*types.CodeFile, []*types.Function, error) {
	lang := a.detectLanguage(path)
	if lang == "" {
		return nil, nil, fmt.Errorf("unsupported language: %s", filepath.Ext(path))
	}

	// Read file
	content, err := a.fileUtils.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	codeFile := &types.CodeFile{
		Path:     path,
		Language: lang,
		Content:  content,
		Hash:     common.Hashes.MD5(content),
	}

	// Parse functions
	functions, err := a.parseFunctions(ctx, codeFile)
	if err != nil {
		return nil, nil, err
	}

	return codeFile, functions, nil
}

// FindUntestedFunctions identifies functions that don't have tests
func (a *Analyzer) FindUntestedFunctions(result *AnalysisResult) []*types.Function {
	a.findUntestedFunctions(result)
	return result.UntestedFunctions
}

// parseFunctions extracts functions from a code file based on language
func (a *Analyzer) parseFunctions(ctx context.Context, file *types.CodeFile) ([]*types.Function, error) {
	switch file.Language {
	case types.Go:
		return a.parseGoFunctions(file)
	case types.NodeJS:
		return a.parseNodeJSFunctions(file)
	case types.Python:
		return a.parsePythonFunctions(file)
	default:
		return nil, fmt.Errorf("unsupported language: %s", file.Language)
	}
}

// parseGoFunctions parses Go functions using AST
func (a *Analyzer) parseGoFunctions(file *types.CodeFile) ([]*types.Function, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file.Path, file.Content, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file: %w", err)
	}

	var functions []*types.Function

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			// Skip test functions
			if strings.HasPrefix(fn.Name.Name, "Test") ||
				strings.HasPrefix(fn.Name.Name, "Benchmark") ||
				strings.HasPrefix(fn.Name.Name, "Example") {
				return true
			}

			// Get function content
			start := fset.Position(fn.Pos()).Offset
			end := fset.Position(fn.End()).Offset
			content := file.Content[start:end]

			// Create function object
			function := &types.Function{
				Name:       fn.Name.Name,
				Language:   types.Go,
				FilePath:   file.Path,
				Signature:  a.formatGoSignature(fn),
				Content:    content,
				LineStart:  fset.Position(fn.Pos()).Line,
				LineEnd:    fset.Position(fn.End()).Line,
				Complexity: a.estimateComplexity(fn),
				IsExported: fn.Name.IsExported(),
			}

			// Filter based on config
			if a.shouldIncludeFunction(function) {
				functions = append(functions, function)
			}
		}
		return true
	})

	return functions, nil
}

// parseNodeJSFunctions parses Node.js functions using regex (simplified)
// In production, you might want to use a proper parser like @babel/parser
func (a *Analyzer) parseNodeJSFunctions(file *types.CodeFile) ([]*types.Function, error) {
	var functions []*types.Function

	// Simple regex to find function declarations
	// This is simplified - use a proper parser in production
	funcRegex := regexp.MustCompile(`(?:function\s+)?(\w+)\s*\([^)]*\)\s*{`)
	matches := funcRegex.FindAllStringSubmatchIndex(file.Content, -1)

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		funcName := file.Content[match[2]:match[3]]

		// Skip test functions
		if strings.HasPrefix(funcName, "test") || strings.HasPrefix(funcName, "it") {
			continue
		}

		// Find the end of the function (simplified)
		start := match[0]
		end := a.findBlockEnd(file.Content, start)

		function := &types.Function{
			Name:       funcName,
			Language:   types.NodeJS,
			FilePath:   file.Path,
			Signature:  file.Content[start:match[1]],
			Content:    file.Content[start:end],
			LineStart:  strings.Count(file.Content[:start], "\n") + 1,
			LineEnd:    strings.Count(file.Content[:end], "\n") + 1,
			Complexity: 1, // Simplified
			IsExported: !strings.HasPrefix(funcName, "_"),
		}

		if a.shouldIncludeFunction(function) {
			functions = append(functions, function)
		}
	}

	return functions, nil
}

// parsePythonFunctions parses Python functions using regex (simplified)
// In production, you might want to use the ast module
func (a *Analyzer) parsePythonFunctions(file *types.CodeFile) ([]*types.Function, error) {
	var functions []*types.Function

	// Simple regex to find function definitions
	funcRegex := regexp.MustCompile(`def\s+(\w+)\s*\([^)]*\):`)
	matches := funcRegex.FindAllStringSubmatchIndex(file.Content, -1)

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		funcName := file.Content[match[2]:match[3]]

		// Skip test functions
		if strings.HasPrefix(funcName, "test_") {
			continue
		}

		// Find the end of the function (by dedent)
		start := match[0]
		end := a.findPythonFunctionEnd(file.Content, start)

		function := &types.Function{
			Name:       funcName,
			Language:   types.Python,
			FilePath:   file.Path,
			Signature:  file.Content[start:match[1]],
			Content:    file.Content[start:end],
			LineStart:  strings.Count(file.Content[:start], "\n") + 1,
			LineEnd:    strings.Count(file.Content[:end], "\n") + 1,
			Complexity: 1, // Simplified
			IsExported: !strings.HasPrefix(funcName, "_"),
		}

		if a.shouldIncludeFunction(function) {
			functions = append(functions, function)
		}
	}

	return functions, nil
}

// mapTestsToFunctions matches test files to the functions they test
func (a *Analyzer) mapTestsToFunctions(result *AnalysisResult) {
	for _, testFile := range result.TestFiles {
		// Extract function names from test file
		testedFunctions := a.extractTestedFunctions(testFile)

		// Mark those functions as having tests
		for _, funcName := range testedFunctions {
			for _, fn := range result.Functions {
				if fn.Name == funcName {
					fn.HasTest = true
					fn.TestFile = testFile.Path
					fn.TestName = a.findTestName(testFile, funcName)
					break
				}
			}
		}
	}
}

// findUntestedFunctions identifies functions without tests
func (a *Analyzer) findUntestedFunctions(result *AnalysisResult) {
	for _, fn := range result.Functions {
		if !fn.HasTest {
			result.UntestedFunctions = append(result.UntestedFunctions, fn)
		}
	}

	// Update stats
	result.Stats.UntestedFunctions = len(result.UntestedFunctions)
	result.Stats.TestedFunctions = len(result.Functions) - len(result.UntestedFunctions)
}

// calculateStats calculates final statistics
func (a *Analyzer) calculateStats(result *AnalysisResult) {
	result.Stats.TotalFiles = len(result.Files)
	result.Stats.TotalFunctions = len(result.Functions)

	for lang, stats := range result.Stats.ByLanguage {
		// Count tested functions per language
		tested := 0
		for _, fn := range result.Functions {
			if fn.Language == lang && fn.HasTest {
				tested++
			}
		}
		stats.TestedFunctions = tested
		stats.UntestedFunctions = stats.Functions - tested
		result.Stats.ByLanguage[lang] = stats
	}
}

// Helper functions

func (a *Analyzer) detectLanguage(path string) types.Language {
	ext := filepath.Ext(path)
	for lang, exts := range a.config.Extensions {
		for _, e := range exts {
			if e == ext {
				return lang
			}
		}
	}
	return ""
}

func (a *Analyzer) isTestFile(path string, lang types.Language) bool {
	base := filepath.Base(path)
	patterns := a.config.TestPatterns[lang]

	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, base)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func (a *Analyzer) detectTestFramework(path string, content string, lang types.Language) string {
	switch lang {
	case types.Go:
		return "testing"
	case types.NodeJS:
		if strings.Contains(content, "describe(") || strings.Contains(content, "it(") {
			return "jest"
		}
		if strings.Contains(content, "test(") {
			return "mocha"
		}
		return "unknown"
	case types.Python:
		if strings.Contains(content, "import pytest") {
			return "pytest"
		}
		if strings.Contains(content, "import unittest") {
			return "unittest"
		}
		return "unknown"
	default:
		return "unknown"
	}
}

func (a *Analyzer) shouldAnalyzeLanguage(lang types.Language) bool {
	if len(a.config.Languages) == 0 {
		return true
	}
	for _, l := range a.config.Languages {
		if l == lang {
			return true
		}
	}
	return false
}

func (a *Analyzer) shouldIncludeFunction(fn *types.Function) bool {
	if a.config.ExportedOnly && !fn.IsExported {
		return false
	}
	if !a.config.IncludePrivate && !fn.IsExported {
		return false
	}
	if a.config.MinComplexity > 0 && fn.Complexity < a.config.MinComplexity {
		return false
	}
	return true
}

func (a *Analyzer) formatGoSignature(fn *ast.FuncDecl) string {
	var sb strings.Builder

	sb.WriteString("func ")
	if fn.Recv != nil {
		sb.WriteString("(")
		// Add receiver (simplified)
		sb.WriteString(")")
	}
	sb.WriteString(fn.Name.Name)
	sb.WriteString("(")

	// Add parameters
	for i, param := range fn.Type.Params.List {
		if i > 0 {
			sb.WriteString(", ")
		}
		for j, name := range param.Names {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(name.Name)
		}
		// Add type (simplified)
	}
	sb.WriteString(")")

	return sb.String()
}

func (a *Analyzer) estimateComplexity(fn *ast.FuncDecl) int {
	// Simple cyclomatic complexity estimation
	complexity := 1

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.CaseClause, *ast.CommClause:
			complexity++
		}
		return true
	})

	return complexity
}

func (a *Analyzer) findBlockEnd(content string, start int) int {
	// Simplified block end detection
	// In production, use proper parsing
	braceCount := 0
	inString := false
	escape := false

	for i := start; i < len(content); i++ {
		c := content[i]

		if escape {
			escape = false
			continue
		}

		if c == '\\' && inString {
			escape = true
			continue
		}

		if c == '"' || c == '\'' || c == '`' {
			inString = !inString
			continue
		}

		if !inString {
			if c == '{' {
				braceCount++
			} else if c == '}' {
				braceCount--
				if braceCount == 0 {
					return i + 1
				}
			}
		}
	}

	return len(content)
}

func (a *Analyzer) findPythonFunctionEnd(content string, start int) int {
	// Find the end of a Python function by indentation
	lines := strings.Split(content[start:], "\n")
	if len(lines) == 0 {
		return len(content)
	}

	// Get indentation of first line
	firstLine := lines[0]
	indent := len(firstLine) - len(strings.TrimLeft(firstLine, " "))

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}

		currentIndent := len(line) - len(strings.TrimLeft(line, " "))
		if currentIndent <= indent {
			// Found end of function
			return start + len(strings.Join(lines[:i], "\n"))
		}
	}

	return len(content)
}

func (a *Analyzer) extractTestedFunctions(testFile *types.TestFile) []string {
	var functions []string

	switch testFile.Language {
	case types.Go:
		// Look for TestXxx functions
		re := regexp.MustCompile(`func\s+(Test\w+)\(`)
		matches := re.FindAllStringSubmatch(testFile.Content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				// Remove "Test" prefix to get function name
				name := strings.TrimPrefix(match[1], "Test")
				functions = append(functions, name)
			}
		}

	case types.NodeJS:
		// Look for test('name') or it('name')
		re := regexp.MustCompile(`(?:test|it)\s*\(\s*['"]([^'"]+)['"]`)
		matches := re.FindAllStringSubmatch(testFile.Content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				functions = append(functions, match[1])
			}
		}

	case types.Python:
		// Look for test_ functions
		re := regexp.MustCompile(`def\s+(test_\w+)\(`)
		matches := re.FindAllStringSubmatch(testFile.Content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				// Remove "test_" prefix
				name := strings.TrimPrefix(match[1], "test_")
				functions = append(functions, name)
			}
		}
	}

	return functions
}

func (a *Analyzer) findTestName(testFile *types.TestFile, functionName string) string {
	switch testFile.Language {
	case types.Go:
		return "Test" + functionName
	case types.NodeJS:
		return functionName
	case types.Python:
		return "test_" + functionName
	default:
		return functionName
	}
}

func shouldSkipDir(name string) bool {
	skipDirs := map[string]bool{
		".git":         true,
		".svn":         true,
		".hg":          true,
		".idea":        true,
		".vscode":      true,
		"node_modules": true,
		"vendor":       true,
		"dist":         true,
		"build":        true,
		"__pycache__":  true,
		"venv":         true,
		"env":          true,
	}
	return skipDirs[name]
}
