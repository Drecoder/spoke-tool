package test

import (
	"context"
	"fmt"
	"strings"

	"example.com/spoke-tool/api/types"
	"example.com/spoke-tool/internal/common"
	"example.com/spoke-tool/internal/model"
)

// Interpreter analyzes test failures and explains WHY they happened
// This component ONLY explains - it NEVER suggests fixes
type Interpreter struct {
	config      InterpreterConfig
	modelClient *model.Client
	logger      *common.Logger
}

// InterpreterConfig configures the interpreter
type InterpreterConfig struct {
	// Model to use for analysis (DeepSeek 7B recommended)
	Model model.ModelType

	// Whether to include code context in explanations
	IncludeContext bool

	// Maximum length of explanations
	MaxExplanationLength int
}

// FailureExplanation represents an explanation of a test failure
// This is PURELY EXPLANATORY - no fixes
type FailureExplanation struct {
	// Test that failed
	TestName string `json:"test_name"`

	// Programming language
	Language types.Language `json:"language"`

	// The error message from the test
	ErrorMessage string `json:"error_message"`

	// Human-readable explanation of WHY it failed
	Explanation string `json:"explanation"`

	// What the test expected
	Expected string `json:"expected,omitempty"`

	// What actually happened
	Actual string `json:"actual,omitempty"`

	// Where in the code the issue occurred
	Location string `json:"location,omitempty"`

	// Relevant code snippet
	CodeSnippet string `json:"code_snippet,omitempty"`
}

// NewInterpreter creates a new test failure interpreter
func NewInterpreter(config InterpreterConfig, modelClient *model.Client) *Interpreter {
	if config.Model == "" {
		config.Model = model.DeepSeek7B // Best for reasoning about failures
	}
	if config.MaxExplanationLength == 0 {
		config.MaxExplanationLength = 500
	}

	return &Interpreter{
		config:      config,
		modelClient: modelClient,
		logger:      common.GetLogger().WithField("component", "test-interpreter"),
	}
}

// ExplainFailure explains why a test failed
// This ONLY explains - it NEVER suggests fixes
func (i *Interpreter) ExplainFailure(ctx context.Context, failure *types.TestFailure) (*FailureExplanation, error) {
	i.logger.Info("Explaining test failure", "test", failure.TestName)

	// Build explanation prompt
	prompt := i.buildExplanationPrompt(failure)

	// Get explanation from model
	resp, err := i.modelClient.Generate(ctx, model.SLMRequest{
		Model:       i.config.Model,
		Language:    failure.Language,
		Prompt:      prompt,
		Temperature: 0.3, // Low temperature for consistent explanations
		MaxTokens:   i.config.MaxExplanationLength,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate explanation: %w", err)
	}

	// Parse the explanation
	explanation := &FailureExplanation{
		TestName:     failure.TestName,
		Language:     failure.Language,
		ErrorMessage: failure.ErrorMsg,
		Explanation:  resp.Response,
	}

	// Extract additional details if available
	i.extractDetails(explanation, failure)

	return explanation, nil
}

// ExplainMultipleFailures explains multiple test failures
func (i *Interpreter) ExplainMultipleFailures(ctx context.Context, failures []*types.TestFailure) ([]*FailureExplanation, error) {
	explanations := make([]*FailureExplanation, 0, len(failures))

	for _, failure := range failures {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		explanation, err := i.ExplainFailure(ctx, failure)
		if err != nil {
			i.logger.Warn("Failed to explain failure", "test", failure.TestName, "error", err)
			// Return a basic explanation instead of failing
			explanation = &FailureExplanation{
				TestName:     failure.TestName,
				Language:     failure.Language,
				ErrorMessage: failure.ErrorMsg,
				Explanation:  "Failed to generate detailed explanation. Check the error message above.",
			}
		}
		explanations = append(explanations, explanation)
	}

	return explanations, nil
}

// AnalyzeFailurePattern looks for patterns in failures
// This is PURELY ANALYTICAL - no fixes
func (i *Interpreter) AnalyzeFailurePattern(ctx context.Context, failures []*FailureExplanation) (string, error) {
	if len(failures) == 0 {
		return "", nil
	}

	// Build a summary of failures
	var sb strings.Builder
	sb.WriteString("Analyze these test failures and identify any patterns:\n\n")

	for _, f := range failures {
		sb.WriteString(fmt.Sprintf("Test: %s\n", f.TestName))
		sb.WriteString(fmt.Sprintf("Error: %s\n", f.ErrorMessage))
		sb.WriteString(fmt.Sprintf("Explanation: %s\n", f.Explanation))
		sb.WriteString("---\n")
	}

	sb.WriteString("\nWhat patterns do you notice? Are these failures related?\n")
	sb.WriteString("DO NOT suggest fixes - just identify patterns.\n")

	resp, err := i.modelClient.Generate(ctx, model.SLMRequest{
		Model:       i.config.Model,
		Prompt:      sb.String(),
		Temperature: 0.3,
		MaxTokens:   300,
	})
	if err != nil {
		return "", err
	}

	return resp.Response, nil
}

// GetFailureSuggestion returns a suggestion for what to investigate
// This is NOT a fix - it's a hint about where to look
func (i *Interpreter) GetFailureSuggestion(ctx context.Context, failure *FailureExplanation) (string, error) {
	prompt := fmt.Sprintf(`Based on this test failure, suggest what the developer should investigate.
DO NOT suggest code changes - just point to areas to examine.

Test: %s
Error: %s
Explanation: %s

What part of the code should the developer look at?
What assumptions might be wrong?
What test cases might be missing?

Provide 2-3 specific areas to investigate.`,
		failure.TestName, failure.ErrorMessage, failure.Explanation)

	resp, err := i.modelClient.Generate(ctx, model.SLMRequest{
		Model:       i.config.Model,
		Prompt:      prompt,
		Temperature: 0.4,
		MaxTokens:   200,
	})
	if err != nil {
		return "", err
	}

	return resp.Response, nil
}

// Helper methods

func (i *Interpreter) buildExplanationPrompt(failure *types.TestFailure) string {
	var sb strings.Builder

	sb.WriteString("Explain why this test failed. DO NOT suggest fixes.\n\n")
	sb.WriteString(fmt.Sprintf("Language: %s\n", failure.Language))
	sb.WriteString(fmt.Sprintf("Test: %s\n", failure.TestName))
	sb.WriteString(fmt.Sprintf("Error: %s\n", failure.ErrorMsg))
	sb.WriteString("\n")

	if i.config.IncludeContext {
		sb.WriteString("Test Code:\n")
		sb.WriteString("```\n")
		sb.WriteString(failure.TestCode)
		sb.WriteString("\n```\n\n")

		sb.WriteString("Source Code:\n")
		sb.WriteString("```\n")
		sb.WriteString(failure.SourceCode)
		sb.WriteString("\n```\n\n")
	}

	sb.WriteString("Explain:\n")
	sb.WriteString("- What the test expected vs what actually happened\n")
	sb.WriteString("- Where in the code the mismatch occurred\n")
	sb.WriteString("- Why the actual behavior differs from expected\n")
	sb.WriteString("\n")
	sb.WriteString("Remember: DO NOT suggest code fixes. Only explain what went wrong.\n")

	return sb.String()
}

func (i *Interpreter) extractDetails(explanation *FailureExplanation, failure *types.TestFailure) {
	// Try to extract expected/actual from error message
	lines := strings.Split(failure.ErrorMsg, "\n")
	for _, line := range lines {
		if strings.Contains(line, "expected") && strings.Contains(line, "got") {
			// Try to parse "expected X, got Y" format
			parts := strings.Split(line, ",")
			for _, part := range parts {
				if strings.Contains(part, "expected") {
					explanation.Expected = strings.TrimSpace(strings.TrimPrefix(part, "expected"))
				}
				if strings.Contains(part, "got") {
					explanation.Actual = strings.TrimSpace(strings.TrimPrefix(part, "got"))
				}
			}
		}
	}

	// Try to extract location
	if failure.TestCode != "" && failure.SourceCode != "" {
		// Simple line number extraction (simplified)
		if strings.Contains(failure.ErrorMsg, ":") {
			parts := strings.Split(failure.ErrorMsg, ":")
			if len(parts) > 1 && strings.Contains(parts[0], ".go") {
				explanation.Location = strings.TrimSpace(parts[0])
			}
		}
	}
}

// FormatExplanation formats an explanation for display
func (i *Interpreter) FormatExplanation(explanation *FailureExplanation) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("❌ Test Failed: %s\n", explanation.TestName))
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("Error: %s\n\n", explanation.ErrorMessage))

	sb.WriteString("Explanation:\n")
	sb.WriteString(explanation.Explanation)
	sb.WriteString("\n\n")

	if explanation.Expected != "" && explanation.Actual != "" {
		sb.WriteString("Details:\n")
		sb.WriteString(fmt.Sprintf("  Expected: %s\n", explanation.Expected))
		sb.WriteString(fmt.Sprintf("  Actual:   %s\n", explanation.Actual))
		sb.WriteString("\n")
	}

	if explanation.Location != "" {
		sb.WriteString(fmt.Sprintf("Location: %s\n", explanation.Location))
	}

	return sb.String()
}

// FormatMultipleExplanations formats multiple explanations
func (i *Interpreter) FormatMultipleExplanations(explanations []*FailureExplanation) string {
	if len(explanations) == 0 {
		return "No failures to explain."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 %d Test Failures\n", len(explanations)))
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	for idx, exp := range explanations {
		sb.WriteString(fmt.Sprintf("[%d/%d] ", idx+1, len(explanations)))
		sb.WriteString(i.FormatExplanation(exp))
		if idx < len(explanations)-1 {
			sb.WriteString(strings.Repeat("-", 40))
			sb.WriteString("\n\n")
		}
	}

	return sb.String()
}
