package doc

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/internal/common"
)

// Extractor handles extraction of documentation content from code and tests
type Extractor struct {
	config      ExtractorConfig
	fileUtils   *common.FileUtils
	stringUtils *common.StringUtils
	logger      *common.Logger
}

// ExtractorConfig configures the extractor
type ExtractorConfig struct {
	// Which languages to extract from
	Languages []types.Language

	// Whether to extract from test files
	IncludeTests bool

	// Whether to extract from comments
	IncludeComments bool

	// Maximum examples per function
	MaxExamplesPerFunc int

	// Minimum confidence for extracted content (0-1)
	MinConfidence float64

	// Whether to include edge cases as examples
	IncludeEdgeCases bool

	// Custom extraction patterns per language
	Patterns map[types.Language][]ExtractionPattern
}

// ExtractionPattern defines a regex pattern for extracting content
type ExtractionPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Example     bool
	Description string
}

// ExtractedContent represents content extracted from code
type ExtractedContent struct {
	// The function this content relates to
	Function *types.Function

	// Language of the content
	Language types.Language

	// Source file path
	SourceFile string

	// Extracted examples
	Examples []*ExtractedExample

	// API documentation
	APIDoc string

	// Code comments
	Comments []string

	// Function signature
	Signature string

	// Whether this came from a test file
	FromTest bool

	// Confidence score (0-1)
	Confidence float64
}

// ExtractedExample represents a code example extracted from tests or code
type ExtractedExample struct {
	// The example code
	Code string

	// Language of the example
	Language types.Language

	// Description (if available)
	Description string

	// Whether this example came from a test
	FromTest bool

	// Whether this is an edge case example
	IsEdgeCase bool

	// The test name (if from test)
	TestName string

	// Expected output (if available)
	ExpectedOutput string

	// Confidence score (0-1)
	Confidence float64
}

// NewExtractor creates a new documentation extractor
func NewExtractor(config ExtractorConfig) *Extractor {
	if config.MaxExamplesPerFunc == 0 {
		config.MaxExamplesPerFunc = 3
	}
	if config.MinConfidence == 0 {
		config.MinConfidence = 0.7
	}

	return &Extractor{
		config:      config,
		fileUtils:   &common.FileUtils{},
		stringUtils: &common.StringUtils{},
		logger:      common.GetLogger().WithField("component", "doc-extractor"),
	}
}

// ExtractFromProject extracts documentation from an entire project
func (e *Extractor) ExtractFromProject(ctx context.Context, analysis *types.CodeAnalysis) ([]*ExtractedContent, error) {
	e.logger.Info("Starting documentation extraction", "files", len(analysis.Files))

	var results []*ExtractedContent

	for _, file := range analysis.Files {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Skip if language not in our list
		if !e.shouldProcessLanguage(file.Language) {
			continue
		}

		// Determine if this is a test file
		isTest := e.isTestFile(file.Path, file.Language)

		// Extract based on language
		var content *ExtractedContent
		var err error

		switch file.Language {
		case types.Go:
			content, err = e.extractFromGo(file, isTest)
		case types.NodeJS:
			content, err = e.extractFromNodeJS(file, isTest)
		case types.Python:
			content, err = e.extractFromPython(file, isTest)
		default:
			e.logger.Debug("Skipping unsupported language", "language", file.Language)
			continue
		}

		if err != nil {
			e.logger.Warn("Failed to extract from file", "file", file.Path, "error", err)
			continue
		}

		if content != nil && len(content.Examples) > 0 {
			results = append(results, content)
		}
	}

	e.logger.Info("Extraction complete", "results", len(results))
	return results, nil
}

// ExtractFromFunction extracts documentation from a specific function
func (e *Extractor) ExtractFromFunction(ctx context.Context, fn *types.Function, tests []*types.TestFile) (*ExtractedContent, error) {
	e.logger.Debug("Extracting from function", "function", fn.Name)

	content := &ExtractedContent{
		Function:   fn,
		Language:   fn.Language,
		SourceFile: fn.FilePath,
		Examples:   []*ExtractedExample{},
		Comments:   []string{},
	}

	// Extract signature
	content.Signature = fn.Signature

	// Extract from function comments
	if e.config.IncludeComments {
		comments := e.extractComments(fn.Content, fn.Language)
		content.Comments = comments
	}

	// Extract examples from tests
	if e.config.IncludeTests {
		for _, test := range tests {
			examples, err := e.extractExamplesFromTest(test, fn.Name)
			if err != nil {
				e.logger.Debug("Failed to extract examples from test", "test", test.Path, "error", err)
				continue
			}
			content.Examples = append(content.Examples, examples...)
		}
	}

	// Limit examples
	if len(content.Examples) > e.config.MaxExamplesPerFunc {
		content.Examples = content.Examples[:e.config.MaxExamplesPerFunc]
	}

	// Calculate confidence
	content.Confidence = e.calculateConfidence(content)

	return content, nil
}

// extractFromGo extracts documentation from Go files
func (e *Extractor) extractFromGo(file *types.CodeFile, isTest bool) (*ExtractedContent, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file.Path, file.Content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file: %w", err)
	}

	content := &ExtractedContent{
		Language:   types.Go,
		SourceFile: file.Path,
		FromTest:   isTest,
		Examples:   []*ExtractedExample{},
		Comments:   []string{},
	}

	// Extract package documentation
	if node.Doc != nil {
		content.Comments = append(content.Comments, node.Doc.Text())
	}

	// Walk the AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Extract function
			funcName := x.Name.Name
			content.Function = &types.Function{
				Name:       funcName,
				Language:   types.Go,
				FilePath:   file.Path,
				Signature:  e.formatGoFunc(x),
				IsExported: x.Name.IsExported(),
			}

			// Extract function comment
			if x.Doc != nil {
				content.Comments = append(content.Comments, x.Doc.Text())
			}

			// Extract examples from test functions
			if isTest && strings.HasPrefix(funcName, "Test") {
				example := e.extractGoTestExample(x, file.Content)
				if example != nil {
					content.Examples = append(content.Examples, example)
				}
			}

		case *ast.GenDecl:
			// Extract constants, variables, types
			if x.Doc != nil {
				content.Comments = append(content.Comments, x.Doc.Text())
			}
		}

		return true
	})

	return content, nil
}

// extractFromNodeJS extracts documentation from Node.js files
func (e *Extractor) extractFromNodeJS(file *types.CodeFile, isTest bool) (*ExtractedContent, error) {
	// Simple regex-based extraction for Node.js
	// In production, you might want to use a proper parser like @babel/parser

	content := &ExtractedContent{
		Language:   types.NodeJS,
		SourceFile: file.Path,
		FromTest:   isTest,
		Examples:   []*ExtractedExample{},
		Comments:   []string{},
	}

	// Extract JSDoc comments and function definitions
	jsdocRegex := regexp.MustCompile(`/\*\*([^*]|\*[^/])*\*/`)
	functionRegex := regexp.MustCompile(`(?:function\s+)?(\w+)\s*\([^)]*\)\s*{`)

	// Find all JSDoc comments
	commentMatches := jsdocRegex.FindAllString(file.Content, -1)
	for _, comment := range commentMatches {
		content.Comments = append(content.Comments, comment)

		// Look for examples in comments
		if strings.Contains(comment, "@example") {
			example := e.extractJSDocExample(comment)
			if example != nil {
				content.Examples = append(content.Examples, example)
			}
		}
	}

	// Find function definitions
	functionMatches := functionRegex.FindAllStringSubmatch(file.Content, -1)
	for _, match := range functionMatches {
		if len(match) > 1 {
			funcName := match[1]
			content.Function = &types.Function{
				Name:      funcName,
				Language:  types.NodeJS,
				FilePath:  file.Path,
				Signature: match[0],
			}
		}
	}

	// Extract test examples
	if isTest {
		testRegex := regexp.MustCompile(`(?:test|it)\s*\(\s*['"]([^'"]+)['"]\s*,\s*\([^)]*\)\s*=>\s*{([^}]+)}`)
		testMatches := testRegex.FindAllStringSubmatch(file.Content, -1)
		for _, match := range testMatches {
			if len(match) > 2 {
				example := &ExtractedExample{
					Code:        strings.TrimSpace(match[2]),
					Language:    types.NodeJS,
					Description: match[1],
					FromTest:    true,
					TestName:    match[1],
					Confidence:  0.9,
				}
				content.Examples = append(content.Examples, example)
			}
		}
	}

	return content, nil
}

// extractFromPython extracts documentation from Python files
func (e *Extractor) extractFromPython(file *types.CodeFile, isTest bool) (*ExtractedContent, error) {
	content := &ExtractedContent{
		Language:   types.Python,
		SourceFile: file.Path,
		FromTest:   isTest,
		Examples:   []*ExtractedExample{},
		Comments:   []string{},
	}

	// Extract docstrings
	docstringRegex := regexp.MustCompile(`(?s)"""(.+?)"""|'''(.+?)'''`)
	docstringMatches := docstringRegex.FindAllStringSubmatch(file.Content, -1)
	for _, match := range docstringMatches {
		docstring := match[1]
		if docstring == "" {
			docstring = match[2]
		}
		if docstring != "" {
			content.Comments = append(content.Comments, docstring)
		}
	}

	// Extract function definitions
	functionRegex := regexp.MustCompile(`def\s+(\w+)\s*\([^)]*\):`)
	functionMatches := functionRegex.FindAllStringSubmatch(file.Content, -1)
	for _, match := range functionMatches {
		if len(match) > 1 {
			funcName := match[1]
			content.Function = &types.Function{
				Name:      funcName,
				Language:  types.Python,
				FilePath:  file.Path,
				Signature: match[0],
			}
		}
	}

	// Extract test examples
	if isTest && strings.Contains(file.Path, "test_") {
		testRegex := regexp.MustCompile(`def\s+(test_\w+)\s*\([^)]*\):\s*\n\s+(.+)`)
		testMatches := testRegex.FindAllStringSubmatch(file.Content, -1)
		for _, match := range testMatches {
			if len(match) > 2 {
				example := &ExtractedExample{
					Code:       strings.TrimSpace(match[2]),
					Language:   types.Python,
					FromTest:   true,
					TestName:   match[1],
					Confidence: 0.9,
				}
				content.Examples = append(content.Examples, example)
			}
		}
	}

	return content, nil
}

// extractGoTestExample extracts an example from a Go test function
func (e *Extractor) extractGoTestExample(funcDecl *ast.FuncDecl, content string) *ExtractedExample {
	// Get the function body
	start := funcDecl.Body.Pos() - 1
	end := funcDecl.Body.End() - 1
	if int(start) < len(content) && int(end) < len(content) {
		body := content[start:end]

		// Look for assertions
		lines := strings.Split(body, "\n")
		var exampleLines []string

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			// Look for function calls that might be examples
			if strings.Contains(trimmed, "funcName") || // Replace with actual function name
				strings.Contains(trimmed, "assert") {
				// Clean up test-specific code
				cleaned := strings.ReplaceAll(trimmed, "t.Errorf", "// check")
				cleaned = strings.ReplaceAll(cleaned, "assert.Equal", "// verify")
				exampleLines = append(exampleLines, cleaned)
			}
		}

		if len(exampleLines) > 0 {
			return &ExtractedExample{
				Code:       strings.Join(exampleLines, "\n"),
				Language:   types.Go,
				FromTest:   true,
				TestName:   funcDecl.Name.Name,
				Confidence: 0.8,
			}
		}
	}

	return nil
}

// extractJSDocExample extracts an example from JSDoc
func (e *Extractor) extractJSDocExample(comment string) *ExtractedExample {
	exampleRegex := regexp.MustCompile(`@example\s+([^*]+)`)
	matches := exampleRegex.FindStringSubmatch(comment)
	if len(matches) > 1 {
		return &ExtractedExample{
			Code:        strings.TrimSpace(matches[1]),
			Language:    types.NodeJS,
			Description: "Example from documentation",
			FromTest:    false,
			Confidence:  0.7,
		}
	}
	return nil
}

// extractExamplesFromTest extracts examples from a test file for a specific function
func (e *Extractor) extractExamplesFromTest(test *types.TestFile, functionName string) ([]*ExtractedExample, error) {
	var examples []*ExtractedExample

	// Simple line-by-line extraction for now
	lines := strings.Split(test.Content, "\n")

	for i, line := range lines {
		// Look for lines that call the function
		if strings.Contains(line, functionName) && !strings.Contains(line, "import") {
			// Get context (a few lines around)
			start := i - 2
			if start < 0 {
				start = 0
			}
			end := i + 3
			if end > len(lines) {
				end = len(lines)
			}

			context := strings.Join(lines[start:end], "\n")

			example := &ExtractedExample{
				Code:       context,
				Language:   test.Language,
				FromTest:   true,
				TestName:   filepath.Base(test.Path),
				Confidence: 0.6,
			}

			// Check if it's an edge case
			if strings.Contains(context, "0") ||
				strings.Contains(context, "nil") ||
				strings.Contains(context, "null") ||
				strings.Contains(context, "negative") {
				example.IsEdgeCase = true
			}

			examples = append(examples, example)

			// Limit examples
			if len(examples) >= e.config.MaxExamplesPerFunc {
				break
			}
		}
	}

	return examples, nil
}

// extractComments extracts comments from code
func (e *Extractor) extractComments(content string, lang types.Language) []string {
	var comments []string

	switch lang {
	case types.Go:
		// Extract // comments and /* */ comments
		lineComment := regexp.MustCompile(`//(.+)`)
		blockComment := regexp.MustCompile(`(?s)/\*(.+?)\*/`)

		lineMatches := lineComment.FindAllStringSubmatch(content, -1)
		for _, match := range lineMatches {
			if len(match) > 1 {
				comments = append(comments, strings.TrimSpace(match[1]))
			}
		}

		blockMatches := blockComment.FindAllStringSubmatch(content, -1)
		for _, match := range blockMatches {
			if len(match) > 1 {
				comments = append(comments, strings.TrimSpace(match[1]))
			}
		}

	case types.NodeJS:
		// Extract // comments and /* */ comments
		lineComment := regexp.MustCompile(`//(.+)`)
		blockComment := regexp.MustCompile(`(?s)/\*(.+?)\*/`)

		lineMatches := lineComment.FindAllStringSubmatch(content, -1)
		for _, match := range lineMatches {
			if len(match) > 1 {
				comments = append(comments, strings.TrimSpace(match[1]))
			}
		}

		blockMatches := blockComment.FindAllStringSubmatch(content, -1)
		for _, match := range blockMatches {
			if len(match) > 1 {
				comments = append(comments, strings.TrimSpace(match[1]))
			}
		}

	case types.Python:
		// Extract # comments and """ docstrings """
		lineComment := regexp.MustCompile(`#(.+)`)
		docstring := regexp.MustCompile(`(?s)"""(.+?)"""|'''(.+?)'''`)

		lineMatches := lineComment.FindAllStringSubmatch(content, -1)
		for _, match := range lineMatches {
			if len(match) > 1 {
				comments = append(comments, strings.TrimSpace(match[1]))
			}
		}

		docMatches := docstring.FindAllStringSubmatch(content, -1)
		for _, match := range docMatches {
			doc := match[1]
			if doc == "" {
				doc = match[2]
			}
			if doc != "" {
				comments = append(comments, strings.TrimSpace(doc))
			}
		}
	}

	return common.Slices.Unique(comments)
}

// formatGoFunc formats a Go function declaration
func (e *Extractor) formatGoFunc(f *ast.FuncDecl) string {
	var sb strings.Builder

	sb.WriteString("func ")
	if f.Recv != nil {
		sb.WriteString("(")
		// Add receiver
		sb.WriteString(")")
	}
	sb.WriteString(f.Name.Name)
	sb.WriteString("(")

	// Add parameters
	for i, param := range f.Type.Params.List {
		if i > 0 {
			sb.WriteString(", ")
		}
		for j, name := range param.Names {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(name.Name)
		}
		// Add type
		// This is simplified - you'd want proper type formatting
	}
	sb.WriteString(")")

	// Add return type
	if f.Type.Results != nil {
		sb.WriteString(" ")
		if len(f.Type.Results.List) > 1 {
			sb.WriteString("(")
		}
		// Add return types
		if len(f.Type.Results.List) > 1 {
			sb.WriteString(")")
		}
	}

	return sb.String()
}

// shouldProcessLanguage checks if we should process this language
func (e *Extractor) shouldProcessLanguage(lang types.Language) bool {
	if len(e.config.Languages) == 0 {
		return true // Process all if not specified
	}
	for _, l := range e.config.Languages {
		if l == lang {
			return true
		}
	}
	return false
}

// isTestFile determines if a file is a test file
func (e *Extractor) isTestFile(path string, lang types.Language) bool {
	filename := filepath.Base(path)

	switch lang {
	case types.Go:
		return strings.HasSuffix(filename, "_test.go")
	case types.NodeJS:
		return strings.HasSuffix(filename, ".test.js") ||
			strings.HasSuffix(filename, ".spec.js")
	case types.Python:
		return strings.HasPrefix(filename, "test_") ||
			strings.HasSuffix(filename, "_test.py")
	default:
		return false
	}
}

// calculateConfidence calculates confidence score for extracted content
func (e *Extractor) calculateConfidence(content *ExtractedContent) float64 {
	if content == nil {
		return 0
	}

	var score float64

	// Examples from tests are more reliable
	for _, ex := range content.Examples {
		if ex.FromTest {
			score += 0.3
		}
		if ex.IsEdgeCase {
			score += 0.1
		}
		score += ex.Confidence * 0.2
	}

	// Comments add confidence
	if len(content.Comments) > 0 {
		score += 0.2
	}

	// Having a function signature is good
	if content.Signature != "" {
		score += 0.2
	}

	// Normalize to 0-1
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// ExtractExamples extracts examples from test files
func (e *Extractor) ExtractExamples(ctx context.Context, testFiles []*types.TestFile) ([]*ExtractedExample, error) {
	var examples []*ExtractedExample

	for _, test := range testFiles {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Extract examples based on language
		var fileExamples []*ExtractedExample
		var err error

		switch test.Language {
		case types.Go:
			fileExamples, err = e.extractGoTestExamples(test)
		case types.NodeJS:
			fileExamples, err = e.extractNodeJSTestExamples(test)
		case types.Python:
			fileExamples, err = e.extractPythonTestExamples(test)
		}

		if err != nil {
			e.logger.Debug("Failed to extract examples from test", "file", test.Path, "error", err)
			continue
		}

		examples = append(examples, fileExamples...)
	}

	return examples, nil
}

// extractGoTestExamples extracts examples from Go test files
func (e *Extractor) extractGoTestExamples(test *types.TestFile) ([]*ExtractedExample, error) {
	var examples []*ExtractedExample

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, test.Path, test.Content, 0)
	if err != nil {
		return nil, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if strings.HasPrefix(funcDecl.Name.Name, "Test") {
				example := e.extractGoTestExample(funcDecl, test.Content)
				if example != nil {
					examples = append(examples, example)
				}
			}
		}
		return true
	})

	return examples, nil
}

// extractNodeJSTestExamples extracts examples from Node.js test files
func (e *Extractor) extractNodeJSTestExamples(test *types.TestFile) ([]*ExtractedExample, error) {
	var examples []*ExtractedExample

	// Simple regex-based extraction
	testRegex := regexp.MustCompile(`(?:test|it)\s*\(\s*['"]([^'"]+)['"]\s*,\s*\([^)]*\)\s*=>\s*{([^}]+)}`)
	matches := testRegex.FindAllStringSubmatch(test.Content, -1)

	for _, match := range matches {
		if len(match) > 2 {
			example := &ExtractedExample{
				Code:        strings.TrimSpace(match[2]),
				Language:    types.NodeJS,
				Description: match[1],
				FromTest:    true,
				TestName:    match[1],
				Confidence:  0.9,
			}
			examples = append(examples, example)
		}
	}

	return examples, nil
}

// extractPythonTestExamples extracts examples from Python test files
func (e *Extractor) extractPythonTestExamples(test *types.TestFile) ([]*ExtractedExample, error) {
	var examples []*ExtractedExample

	// Simple regex-based extraction
	testRegex := regexp.MustCompile(`def\s+(test_\w+)\s*\([^)]*\):\s*\n\s+(.+)`)
	matches := testRegex.FindAllStringSubmatch(test.Content, -1)

	for _, match := range matches {
		if len(match) > 2 {
			example := &ExtractedExample{
				Code:       strings.TrimSpace(match[2]),
				Language:   types.Python,
				FromTest:   true,
				TestName:   match[1],
				Confidence: 0.9,
			}
			examples = append(examples, example)
		}
	}

	return examples, nil
}

// CleanExample removes test-specific code from an example
func (e *Extractor) CleanExample(example *ExtractedExample) string {
	code := example.Code

	// Remove assertions
	code = strings.ReplaceAll(code, "t.Errorf", "// check")
	code = strings.ReplaceAll(code, "assert.Equal", "// verify")
	code = strings.ReplaceAll(code, "assert.", "// ")
	code = strings.ReplaceAll(code, "expect(", "// ")
	code = strings.ReplaceAll(code, ".toBe(", "// ")

	// Remove test function wrappers
	code = regexp.MustCompile(`func Test\w+\(t \*testing\.T\) {`).ReplaceAllString(code, "")
	code = regexp.MustCompile(`test\(['"].+['"], async \(\) => {`).ReplaceAllString(code, "")
	code = regexp.MustCompile(`def test_\w+\(self\):`).ReplaceAllString(code, "")

	// Remove trailing brackets
	code = strings.TrimSuffix(code, "}")
	code = strings.TrimSpace(code)

	return code
}
