package model

import (
	"fmt"
	"strings"

	"example.com/spoke-tool/api/types"
)

// PromptTemplates holds all prompt templates for different tasks and languages
// This is PURELY TEMPLATES - no logic, no suggestions, no fixes
type PromptTemplates struct {
	// Code understanding templates (codellama:7b)
	CodeUnderstanding map[types.Language]string

	// Test generation templates (DeepSeek 7B)
	TestGeneration map[types.Language]string

	// Documentation templates (Gemma 2B)
	Documentation map[types.Language]string

	// Failure analysis templates (DeepSeek 7B) - EXPLAINS only
	FailureAnalysis map[types.Language]string

	// Example extraction templates (codellama:7b/Gemma)
	ExampleExtraction map[types.Language]string
}

// NewPromptTemplates creates a new set of prompt templates
func NewPromptTemplates() *PromptTemplates {
	return &PromptTemplates{
		CodeUnderstanding: getCodeUnderstandingTemplates(),
		TestGeneration:    getTestGenerationTemplates(),
		Documentation:     getDocumentationTemplates(),
		FailureAnalysis:   getFailureAnalysisTemplates(),
		ExampleExtraction: getExampleExtractionTemplates(),
	}
}

// GetCodeUnderstandingPrompt returns a prompt for code understanding
func (p *PromptTemplates) GetCodeUnderstandingPrompt(lang types.Language, code string) string {
	tmpl, ok := p.CodeUnderstanding[lang]
	if !ok {
		tmpl = p.CodeUnderstanding[types.Go] // Fallback to Go
	}
	return fmt.Sprintf(tmpl, code)
}

// GetTestGenerationPrompt returns a prompt for test generation
func (p *PromptTemplates) GetTestGenerationPrompt(lang types.Language, functionName, code, deps string) string {
	tmpl, ok := p.TestGeneration[lang]
	if !ok {
		tmpl = p.TestGeneration[types.Go] // Fallback to Go
	}

	framework := getTestFramework(lang)
	return fmt.Sprintf(tmpl, lang, functionName, code, deps, framework)
}

// GetDocumentationPrompt returns a prompt for documentation generation
func (p *PromptTemplates) GetDocumentationPrompt(lang types.Language, functionName, code string) string {
	tmpl, ok := p.Documentation[lang]
	if !ok {
		tmpl = p.Documentation[types.Go] // Fallback to Go
	}

	docFormat := getDocFormat(lang)
	return fmt.Sprintf(tmpl, lang, functionName, code, docFormat)
}

// GetFailureAnalysisPrompt returns a prompt for test failure analysis
// This ONLY asks for explanation - NO fixes
func (p *PromptTemplates) GetFailureAnalysisPrompt(lang types.Language, testName, errorMsg, testCode, sourceCode string) string {
	tmpl, ok := p.FailureAnalysis[lang]
	if !ok {
		tmpl = p.FailureAnalysis[types.Go] // Fallback to Go
	}

	return fmt.Sprintf(tmpl, testName, errorMsg, testCode, sourceCode, lang)
}

// GetExampleExtractionPrompt returns a prompt for extracting examples from tests
func (p *PromptTemplates) GetExampleExtractionPrompt(lang types.Language, testCode string) string {
	tmpl, ok := p.ExampleExtraction[lang]
	if !ok {
		tmpl = p.ExampleExtraction[types.Go] // Fallback to Go
	}

	return fmt.Sprintf(tmpl, lang, testCode)
}

// Template definitions - PURE STRINGS, no logic

func getCodeUnderstandingTemplates() map[types.Language]string {
	return map[types.Language]string{
		types.Go: `Analyze this Go code and describe its structure.
DO NOT suggest improvements - just describe what you see.

Code:
%s

Describe:
- Package name
- Functions and their purposes
- Parameters and return values
- Types and structs
- Dependencies
- Any notable patterns

Keep the description factual and objective.`,

		types.NodeJS: `Analyze this Node.js code and describe its structure.
DO NOT suggest improvements - just describe what you see.

Code:
%s

Describe:
- Module/exports
- Functions and their purposes
- Parameters and return values
- Classes and methods
- Dependencies (require/import)
- Async patterns if any
- Any notable patterns

Keep the description factual and objective.`,

		types.Python: `Analyze this Python code and describe its structure.
DO NOT suggest improvements - just describe what you see.

Code:
%s

Describe:
- Module name
- Functions and their purposes
- Parameters and return values
- Classes and methods
- Imports
- Decorators if any
- Any notable patterns

Keep the description factual and objective.`,
	}
}

func getTestGenerationTemplates() map[types.Language]string {
	return map[types.Language]string{
		types.Go: `Generate unit tests for this Go function.
DO NOT modify the original code - only create tests.

Language: %s
Function: %s
Code:
%s
Dependencies: %s

Create tests using the %s framework that verify:
- Happy path (normal inputs)
- Edge cases (zero values, boundaries)
- Error conditions
- Table-driven tests where appropriate

Return ONLY the test code, no explanations.`,

		types.NodeJS: `Generate unit tests for this Node.js function.
DO NOT modify the original code - only create tests.

Language: %s
Function: %s
Code:
%s
Dependencies: %s

Create tests using %s that verify:
- Happy path (normal inputs)
- Edge cases (null, undefined, boundaries)
- Error conditions
- Async behavior if applicable
- Mocks for dependencies

Return ONLY the test code, no explanations.`,

		types.Python: `Generate unit tests for this Python function.
DO NOT modify the original code - only create tests.

Language: %s
Function: %s
Code:
%s
Dependencies: %s

Create tests using %s that verify:
- Happy path (normal inputs)
- Edge cases (None, empty, boundaries)
- Error conditions (exceptions)
- Fixtures for setup
- Mocks for dependencies

Return ONLY the test code, no explanations.`,
	}
}

func getDocumentationTemplates() map[types.Language]string {
	return map[types.Language]string{
		types.Go: `Write clear documentation for this Go function.
DO NOT suggest changes - just document what it does.

Language: %s
Function: %s
Code:
%s

Write a godoc-style comment including:
- Brief description of what the function does
- Parameter explanations
- Return value description
- Any panics or errors
- One simple example

Use %s format.
Return ONLY the documentation comment.`,

		types.NodeJS: `Write clear JSDoc documentation for this Node.js function.
DO NOT suggest changes - just document what it does.

Language: %s
Function: %s
Code:
%s

Write JSDoc including:
- @description of what the function does
- @param tags with types and descriptions
- @returns tag with type and description
- @throws if applicable
- @example of usage

Use %s format.
Return ONLY the JSDoc comment.`,

		types.Python: `Write clear docstring documentation for this Python function.
DO NOT suggest changes - just document what it does.

Language: %s
Function: %s
Code:
%s

Write a Google-style docstring including:
- Description of what the function does
- Args: with types and descriptions
- Returns: with type and description
- Raises: if applicable
- Example: usage example

Use %s format.
Return ONLY the docstring.`,
	}
}

func getFailureAnalysisTemplates() map[types.Language]string {
	return map[types.Language]string{
		types.Go: `Analyze why this Go test failed.
DO NOT suggest code fixes - just explain what went wrong.

Test: %s
Error: %s
Test Code:
%s
Source Code:
%s

Explain:
- What the test expected vs what actually happened
- Where in the code the mismatch occurred
- Why the actual behavior differs from expected

Focus on explaining the failure, NOT how to fix it.`,

		types.NodeJS: `Analyze why this Node.js/Jest test failed.
DO NOT suggest code fixes - just explain what went wrong.

Test: %s
Error: %s
Test Code:
%s
Source Code:
%s

Explain:
- What the test expected vs what actually happened
- Where in the code the mismatch occurred
- Whether async issues might be involved
- Why the actual behavior differs from expected

Focus on explaining the failure, NOT how to fix it.`,

		types.Python: `Analyze why this Python/pytest test failed.
DO NOT suggest code fixes - just explain what went wrong.

Test: %s
Error: %s
Test Code:
%s
Source Code:
%s

Explain:
- What the test expected vs what actually happened
- Where in the code the mismatch occurred
- Whether exception handling is involved
- Why the actual behavior differs from expected

Focus on explaining the failure, NOT how to fix it.`,
	}
}

func getExampleExtractionTemplates() map[types.Language]string {
	return map[types.Language]string{
		types.Go: `Extract a clean usage example from this Go test.
DO NOT modify the example - just extract the relevant parts.

Language: %s
Test Code:
%s

Extract the function calls and assertions that demonstrate usage.
Remove test-specific boilerplate (t.Run, t.Errorf, etc.).
Format as a clean, runnable example.
Return ONLY the example code.`,

		types.NodeJS: `Extract a clean usage example from this Node.js/Jest test.
DO NOT modify the example - just extract the relevant parts.

Language: %s
Test Code:
%s

Extract the function calls and expectations that demonstrate usage.
Remove test-specific boilerplate (describe, it, test, etc.).
Format as a clean example.
Return ONLY the example code.`,

		types.Python: `Extract a clean usage example from this Python/pytest test.
DO NOT modify the example - just extract the relevant parts.

Language: %s
Test Code:
%s

Extract the function calls and assertions that demonstrate usage.
Remove test-specific boilerplate (def test_, assert, etc.).
Format as a clean example.
Return ONLY the example code.`,
	}
}

// Helper functions to get framework and format names

func getTestFramework(lang types.Language) string {
	switch lang {
	case types.Go:
		return "testing"
	case types.NodeJS:
		return "Jest"
	case types.Python:
		return "pytest"
	default:
		return "testing"
	}
}

func getDocFormat(lang types.Language) string {
	switch lang {
	case types.Go:
		return "godoc"
	case types.NodeJS:
		return "JSDoc"
	case types.Python:
		return "Google-style docstring"
	default:
		return "documentation"
	}
}

// PromptBuilder helps build prompts with proper formatting
type PromptBuilder struct {
	sb strings.Builder
}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

// AddLine adds a line to the prompt
func (p *PromptBuilder) AddLine(line string) *PromptBuilder {
	p.sb.WriteString(line)
	p.sb.WriteString("\n")
	return p
}

// Addf adds a formatted line to the prompt
func (p *PromptBuilder) Addf(format string, args ...interface{}) *PromptBuilder {
	p.sb.WriteString(fmt.Sprintf(format, args...))
	p.sb.WriteString("\n")
	return p
}

// AddSection adds a section with a header
func (p *PromptBuilder) AddSection(header string, content string) *PromptBuilder {
	p.sb.WriteString(header)
	p.sb.WriteString(":\n")
	p.sb.WriteString(content)
	p.sb.WriteString("\n\n")
	return p
}

// AddCode adds a code block
func (p *PromptBuilder) AddCode(language string, code string) *PromptBuilder {
	p.sb.WriteString("```")
	p.sb.WriteString(language)
	p.sb.WriteString("\n")
	p.sb.WriteString(code)
	p.sb.WriteString("\n```\n\n")
	return p
}

// AddInstruction adds an instruction line
func (p *PromptBuilder) AddInstruction(instruction string) *PromptBuilder {
	p.sb.WriteString("- ")
	p.sb.WriteString(instruction)
	p.sb.WriteString("\n")
	return p
}

// String returns the built prompt
func (p *PromptBuilder) String() string {
	return p.sb.String()
}

// Reset clears the builder
func (p *PromptBuilder) Reset() {
	p.sb.Reset()
}

// Predefined instruction sets

// GetStandardTestInstructions returns standard test generation instructions
func GetStandardTestInstructions() []string {
	return []string{
		"Happy path (normal inputs)",
		"Edge cases (zero values, boundaries)",
		"Error conditions",
		"Table-driven tests where appropriate",
		"Use the language's testing framework",
		"Return ONLY the test code, no explanations",
	}
}

// GetStandardDocInstructions returns standard documentation instructions
func GetStandardDocInstructions() []string {
	return []string{
		"Brief description of what the function does",
		"Parameter explanations",
		"Return value description",
		"One simple example",
		"Use the language's doc format",
		"Return ONLY the documentation",
	}
}

// GetStandardAnalysisInstructions returns standard failure analysis instructions
func GetStandardAnalysisInstructions() []string {
	return []string{
		"What the test expected vs what actually happened",
		"Where in the code the mismatch occurred",
		"Why the actual behavior differs from expected",
		"DO NOT suggest fixes",
		"Focus on explaining, not fixing",
	}
}
