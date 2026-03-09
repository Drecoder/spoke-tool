package extractors

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"github.com/yourusername/spoke-tool/api/types"
)

// ExampleExtractor extracts code examples from test files
// This is PURELY EXTRACTIVE - no modifications
type ExampleExtractor struct {
	// No state - pure functions
}

// Example represents a code example extracted from tests
type Example struct {
	// The example code
	Code string `json:"code"`

	// Language of the example
	Language types.Language `json:"language"`

	// Source file where example was found
	SourceFile string `json:"source_file"`

	// Function being demonstrated
	FunctionName string `json:"function_name,omitempty"`

	// Test name that contained this example
	TestName string `json:"test_name,omitempty"`

	// Description of what the example shows
	Description string `json:"description,omitempty"`

	// Input values used in the example
	Inputs map[string]string `json:"inputs,omitempty"`

	// Expected output
	ExpectedOutput string `json:"expected_output,omitempty"`

	// Whether this is an edge case example
	IsEdgeCase bool `json:"is_edge_case"`

	// Line number in source file
	Line int `json:"line"`

	// Confidence score (0-1) that this is a good example
	Confidence float64 `json:"confidence"`
}

// ExtractionOptions configures example extraction
type ExtractionOptions struct {
	// Maximum number of examples to extract per function
	MaxExamplesPerFunction int

	// Minimum confidence threshold
	MinConfidence float64

	// Whether to include edge cases
	IncludeEdgeCases bool

	// Whether to clean test-specific code
	CleanCode bool

	// Languages to extract from
	Languages []types.Language
}

// NewExampleExtractor creates a new example extractor
func NewExampleExtractor() *ExampleExtractor {
	return &ExampleExtractor{}
}

// ExtractFromGoTests extracts examples from Go test files
func (e *ExampleExtractor) ExtractFromGoTests(content string, filePath string, opts *ExtractionOptions) ([]*Example, error) {
	if opts == nil {
		opts = &ExtractionOptions{
			MaxExamplesPerFunction: 3,
			MinConfidence:          0.7,
			IncludeEdgeCases:       true,
			CleanCode:              true,
		}
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, content, 0)
	if err != nil {
		return nil, err
	}

	var examples []*Example

	// Look for test functions
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			// Check if it's a test function
			if strings.HasPrefix(fn.Name.Name, "Test") ||
				strings.HasPrefix(fn.Name.Name, "Example") {

				example := e.extractGoTestExample(fn, content, fset)
				if example != nil && example.Confidence >= opts.MinConfidence {
					examples = append(examples, example)
				}
			}
		}
		return true
	})

	// Limit examples per function
	if opts.MaxExamplesPerFunction > 0 && len(examples) > opts.MaxExamplesPerFunction {
		examples = examples[:opts.MaxExamplesPerFunction]
	}

	return examples, nil
}

// ExtractFromNodeJSTests extracts examples from Node.js/Jest test files
func (e *ExampleExtractor) ExtractFromNodeJSTests(content string, filePath string, opts *ExtractionOptions) ([]*Example, error) {
	if opts == nil {
		opts = &ExtractionOptions{
			MaxExamplesPerFunction: 3,
			MinConfidence:          0.7,
			IncludeEdgeCases:       true,
			CleanCode:              true,
		}
	}

	var examples []*Example

	lines := strings.Split(content, "\n")

	// Look for test blocks
	testRegex := regexp.MustCompile(`(?:test|it)\s*\(\s*['"]([^'"]+)['"]\s*,\s*(?:async\s*)?\(([^)]*)\)\s*=>\s*{`)

	for i, line := range lines {
		lineNum := i + 1

		if matches := testRegex.FindStringSubmatch(line); len(matches) > 1 {
			testName := matches[1]

			// Find the end of the test block
			endLine, blockContent := e.extractBlock(lines, i)

			example := &Example{
				Code:       blockContent,
				Language:   types.NodeJS,
				SourceFile: filePath,
				TestName:   testName,
				Line:       lineNum,
				IsEdgeCase: e.isEdgeCase(testName, blockContent),
				Confidence: 0.8,
			}

			// Try to determine which function is being tested
			example.FunctionName = e.extractFunctionNameFromTest(blockContent)

			// Extract inputs and expected output
			example.Inputs = e.extractInputs(blockContent, types.NodeJS)
			example.ExpectedOutput = e.extractExpectedOutput(blockContent, types.NodeJS)

			// Clean the code if requested
			if opts.CleanCode {
				example.Code = e.cleanNodeJSTestCode(blockContent)
			}

			// Get description from test name
			example.Description = e.cleanTestName(testName)

			if example.Confidence >= opts.MinConfidence {
				examples = append(examples, example)
			}
		}
	}

	// Limit examples
	if opts.MaxExamplesPerFunction > 0 && len(examples) > opts.MaxExamplesPerFunction {
		examples = examples[:opts.MaxExamplesPerFunction]
	}

	return examples, nil
}

// ExtractFromPythonTests extracts examples from Python test files
func (e *ExampleExtractor) ExtractFromPythonTests(content string, filePath string, opts *ExtractionOptions) ([]*Example, error) {
	if opts == nil {
		opts = &ExtractionOptions{
			MaxExamplesPerFunction: 3,
			MinConfidence:          0.7,
			IncludeEdgeCases:       true,
			CleanCode:              true,
		}
	}

	var examples []*Example

	lines := strings.Split(content, "\n")

	// Look for test functions
	testRegex := regexp.MustCompile(`def\s+(test_\w+)\s*\([^)]*\):`)

	for i, line := range lines {
		lineNum := i + 1

		if matches := testRegex.FindStringSubmatch(line); len(matches) > 1 {
			testName := matches[1]

			// Find the end of the test function
			endLine, blockContent := e.extractPythonFunction(lines, i)

			example := &Example{
				Code:       blockContent,
				Language:   types.Python,
				SourceFile: filePath,
				TestName:   testName,
				Line:       lineNum,
				IsEdgeCase: e.isEdgeCase(testName, blockContent),
				Confidence: 0.8,
			}

			// Try to determine which function is being tested
			example.FunctionName = e.extractFunctionNameFromTest(blockContent)

			// Extract inputs and expected output
			example.Inputs = e.extractInputs(blockContent, types.Python)
			example.ExpectedOutput = e.extractExpectedOutput(blockContent, types.Python)

			// Clean the code if requested
			if opts.CleanCode {
				example.Code = e.cleanPythonTestCode(blockContent)
			}

			// Get description from test name
			example.Description = e.cleanTestName(testName)

			if example.Confidence >= opts.MinConfidence {
				examples = append(examples, example)
			}
		}
	}

	// Limit examples
	if opts.MaxExamplesPerFunction > 0 && len(examples) > opts.MaxExamplesPerFunction {
		examples = examples[:opts.MaxExamplesPerFunction]
	}

	return examples, nil
}

// ExtractExamples automatically detects language and extracts examples
func (e *ExampleExtractor) ExtractExamples(content string, filePath string, lang types.Language, opts *ExtractionOptions) ([]*Example, error) {
	switch lang {
	case types.Go:
		return e.ExtractFromGoTests(content, filePath, opts)
	case types.NodeJS:
		return e.ExtractFromNodeJSTests(content, filePath, opts)
	case types.Python:
		return e.ExtractFromPythonTests(content, filePath, opts)
	default:
		return nil, nil
	}
}

// Helper methods for Go extraction

func (e *ExampleExtractor) extractGoTestExample(fn *ast.FuncDecl, content string, fset *token.FileSet) *Example {
	start := fset.Position(fn.Pos()).Offset
	end := fset.Position(fn.End()).Offset

	if start < 0 || end > len(content) {
		return nil
	}

	body := content[start:end]

	example := &Example{
		Code:       body,
		Language:   types.Go,
		TestName:   fn.Name.Name,
		Line:       fset.Position(fn.Pos()).Line,
		Confidence: 0.7,
	}

	// Determine which function is being tested
	if strings.HasPrefix(fn.Name.Name, "Test") {
		example.FunctionName = strings.TrimPrefix(fn.Name.Name, "Test")
	} else if strings.HasPrefix(fn.Name.Name, "Example") {
		example.FunctionName = strings.TrimPrefix(fn.Name.Name, "Example")
	}

	// Check if it's an edge case
	example.IsEdgeCase = e.isEdgeCase(fn.Name.Name, body)

	// Extract inputs and expected output
	example.Inputs = e.extractInputs(body, types.Go)
	example.ExpectedOutput = e.extractExpectedOutput(body, types.Go)

	// Clean the code
	example.Code = e.cleanGoTestCode(body)

	return example
}

// Helper methods for block extraction

func (e *ExampleExtractor) extractBlock(lines []string, startLine int) (int, string) {
	braceCount := 0
	found := false
	var block []string

	for i := startLine; i < len(lines); i++ {
		line := lines[i]
		block = append(block, line)

		// Count braces to find the end
		for _, c := range line {
			if c == '{' {
				braceCount++
				found = true
			} else if c == '}' {
				braceCount--
			}
		}

		if found && braceCount == 0 {
			return i, strings.Join(block, "\n")
		}
	}

	return len(lines) - 1, strings.Join(block, "\n")
}

func (e *ExampleExtractor) extractPythonFunction(lines []string, startLine int) (int, string) {
	// Get indentation of first line
	firstLine := lines[startLine]
	indent := len(firstLine) - len(strings.TrimLeft(firstLine, " "))

	var block []string
	block = append(block, firstLine)

	for i := startLine + 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			block = append(block, line)
			continue
		}

		currentIndent := len(line) - len(strings.TrimLeft(line, " "))
		if currentIndent <= indent && trimmed != "" {
			// End of function
			return i - 1, strings.Join(block, "\n")
		}

		block = append(block, line)
	}

	return len(lines) - 1, strings.Join(block, "\n")
}

// Code cleaning methods

func (e *ExampleExtractor) cleanGoTestCode(code string) string {
	lines := strings.Split(code, "\n")
	var cleaned []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Remove test function declaration
		if strings.HasPrefix(trimmed, "func Test") {
			continue
		}
		if strings.HasPrefix(trimmed, "func Example") {
			continue
		}

		// Remove test-specific assertions
		line = strings.ReplaceAll(line, "t.Errorf", "// verify")
		line = strings.ReplaceAll(line, "t.Fatalf", "// verify")
		line = strings.ReplaceAll(line, "assert.Equal", "// verify")
		line = strings.ReplaceAll(line, "assert.", "// ")

		cleaned = append(cleaned, line)
	}

	// Remove empty lines at start and end
	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

func (e *ExampleExtractor) cleanNodeJSTestCode(code string) string {
	lines := strings.Split(code, "\n")
	var cleaned []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Remove test function declaration
		if strings.HasPrefix(trimmed, "test(") || strings.HasPrefix(trimmed, "it(") {
			continue
		}

		// Remove assertions
		line = strings.ReplaceAll(line, "expect(", "// expect(")
		line = strings.ReplaceAll(line, ".toBe(", " // should be ")
		line = strings.ReplaceAll(line, ".toEqual(", " // should equal ")
		line = strings.ReplaceAll(line, "assert.", "// ")

		cleaned = append(cleaned, line)
	}

	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

func (e *ExampleExtractor) cleanPythonTestCode(code string) string {
	lines := strings.Split(code, "\n")
	var cleaned []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Remove test function declaration
		if strings.HasPrefix(trimmed, "def test_") {
			continue
		}

		// Remove assertions
		line = strings.ReplaceAll(line, "assert ", "# verify: ")
		line = strings.ReplaceAll(line, "self.assertEqual", "# verify")
		line = strings.ReplaceAll(line, "pytest.raises", "# should raise")

		cleaned = append(cleaned, line)
	}

	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

// Extraction helpers

func (e *ExampleExtractor) extractFunctionNameFromTest(code string) string {
	// Look for function calls in the test
	funcRegex := regexp.MustCompile(`(\w+)\s*\(`)
	matches := funcRegex.FindAllStringSubmatch(code, -1)

	// Skip common test functions
	skip := map[string]bool{
		"assert": true, "require": true, "expect": true,
		"t": true, "self": true, "test": true, "it": true,
	}

	for _, match := range matches {
		if len(match) > 1 {
			name := match[1]
			if !skip[name] && !strings.HasPrefix(name, "Test") {
				return name
			}
		}
	}

	return ""
}

func (e *ExampleExtractor) extractInputs(code string, lang types.Language) map[string]string {
	inputs := make(map[string]string)

	switch lang {
	case types.Go:
		// Look for variable assignments
		assignRegex := regexp.MustCompile(`(\w+)\s*:=\s*([^,\n]+)`)
		matches := assignRegex.FindAllStringSubmatch(code, -1)
		for _, match := range matches {
			if len(match) > 2 {
				inputs[match[1]] = strings.TrimSpace(match[2])
			}
		}

	case types.NodeJS:
		// Look for const/let assignments
		assignRegex := regexp.MustCompile(`(?:const|let)\s+(\w+)\s*=\s*([^;\n]+)`)
		matches := assignRegex.FindAllStringSubmatch(code, -1)
		for _, match := range matches {
			if len(match) > 2 {
				inputs[match[1]] = strings.TrimSpace(match[2])
			}
		}

	case types.Python:
		// Look for variable assignments
		assignRegex := regexp.MustCompile(`(\w+)\s*=\s*([^#\n]+)`)
		matches := assignRegex.FindAllStringSubmatch(code, -1)
		for _, match := range matches {
			if len(match) > 2 {
				inputs[match[1]] = strings.TrimSpace(match[2])
			}
		}
	}

	return inputs
}

func (e *ExampleExtractor) extractExpectedOutput(code string, lang types.Language) string {
	switch lang {
	case types.Go:
		// Look for expected values in assertions
		re := regexp.MustCompile(`(?:want|expected)\s*[=:]\s*([^,\n]+)`)
		if matches := re.FindStringSubmatch(code); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}

		// Look for error messages
		re = regexp.MustCompile(`t\.Errorf\([^,]+,\s*([^)]+)\)`)
		if matches := re.FindStringSubmatch(code); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}

	case types.NodeJS:
		// Look for toBe assertions
		re := regexp.MustCompile(`\.toBe\(([^)]+)\)`)
		if matches := re.FindStringSubmatch(code); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}

		// Look for toEqual assertions
		re = regexp.MustCompile(`\.toEqual\(([^)]+)\)`)
		if matches := re.FindStringSubmatch(code); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}

	case types.Python:
		// Look for assert equality
		re := regexp.MustCompile(`assert.*==\s*([^#\n]+)`)
		if matches := re.FindStringSubmatch(code); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

func (e *ExampleExtractor) isEdgeCase(name string, code string) bool {
	lowerName := strings.ToLower(name)
	lowerCode := strings.ToLower(code)

	edgeIndicators := []string{
		"zero", "empty", "nil", "null", "negative",
		"edge", "boundary", "limit", "max", "min",
		"invalid", "error", "exception", "panic",
		"divide by zero", "overflow", "underflow",
	}

	for _, indicator := range edgeIndicators {
		if strings.Contains(lowerName, indicator) {
			return true
		}
		if strings.Contains(lowerCode, indicator) {
			return true
		}
	}

	return false
}

func (e *ExampleExtractor) cleanTestName(name string) string {
	// Remove prefixes
	name = strings.TrimPrefix(name, "Test")
	name = strings.TrimPrefix(name, "test_")
	name = strings.TrimPrefix(name, "Example")

	// Convert camelCase to spaces
	words := regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(name, "$1 $2")

	// Convert underscores to spaces
	words = strings.ReplaceAll(words, "_", " ")

	return strings.ToLower(words)
}

// FormatExample formats an example for display in documentation
func (e *ExampleExtractor) FormatExample(example *Example) string {
	var sb strings.Builder

	if example.Description != "" {
		sb.WriteString(fmt.Sprintf("**%s**\n\n", strings.Title(example.Description)))
	}

	// Show inputs if available
	if len(example.Inputs) > 0 {
		sb.WriteString("Inputs:\n")
		for k, v := range example.Inputs {
			sb.WriteString(fmt.Sprintf("- %s = %s\n", k, v))
		}
		sb.WriteString("\n")
	}

	// Show code
	sb.WriteString("```" + string(example.Language) + "\n")
	sb.WriteString(example.Code)
	sb.WriteString("\n```\n")

	// Show expected output
	if example.ExpectedOutput != "" {
		sb.WriteString("\n**Expected output:**\n")
		sb.WriteString("```\n")
		sb.WriteString(example.ExpectedOutput)
		sb.WriteString("\n```\n")
	}

	if example.IsEdgeCase {
		sb.WriteString("\n⚠️ *This is an edge case example*\n")
	}

	return sb.String()
}

// GetBestExamples returns the highest confidence examples
func (e *ExampleExtractor) GetBestExamples(examples []*Example, count int) []*Example {
	if len(examples) <= count {
		return examples
	}

	// Sort by confidence (descending)
	for i := 0; i < len(examples)-1; i++ {
		for j := i + 1; j < len(examples); j++ {
			if examples[i].Confidence < examples[j].Confidence {
				examples[i], examples[j] = examples[j], examples[i]
			}
		}
	}

	return examples[:count]
}

// ExtractFromTestFiles extracts examples from multiple test files
func (e *ExampleExtractor) ExtractFromTestFiles(files []*types.TestFile, opts *ExtractionOptions) ([]*Example, error) {
	var allExamples []*Example

	for _, file := range files {
		examples, err := e.ExtractExamples(file.Content, file.Path, file.Language, opts)
		if err != nil {
			continue // Skip files with errors
		}
		allExamples = append(allExamples, examples...)
	}

	return allExamples, nil
}

// GroupByFunction groups examples by the function they demonstrate
func (e *ExampleExtractor) GroupByFunction(examples []*Example) map[string][]*Example {
	grouped := make(map[string][]*Example)

	for _, ex := range examples {
		if ex.FunctionName != "" {
			grouped[ex.FunctionName] = append(grouped[ex.FunctionName], ex)
		}
	}

	return grouped
}
