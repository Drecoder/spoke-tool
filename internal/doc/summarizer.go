package doc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/internal/common"
	"github.com/yourusername/spoke-tool/internal/model"
)

// Summarizer handles generating human-readable summaries from code
// NOTE: This component ONLY summarizes - it never suggests or changes code
type Summarizer struct {
	config      SummarizerConfig
	modelClient *model.Client
	stringUtils *common.StringUtils
	logger      *common.Logger
	cache       *SummaryCache
}

// SummarizerConfig configures the summarizer
type SummarizerConfig struct {
	// Model to use for summarization (usually Gemma 2B for speed)
	Model model.ModelType

	// Maximum length of summaries
	MaxSummaryLength int

	// Whether to detect and highlight edge cases
	DetectEdgeCases bool

	// Whether to use caching (recommended)
	UseCache bool

	// Cache duration
	CacheTTL time.Duration
}

// SummaryCache caches generated summaries to avoid redundant API calls
type SummaryCache struct {
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// CacheEntry represents a cached summary
type CacheEntry struct {
	Summary   string
	Timestamp time.Time
	Function  string
	Language  types.Language
}

// FunctionSummary represents a human-readable summary of a function
// This is PURELY DESCRIPTIVE - no suggestions or code changes
type FunctionSummary struct {
	// The function being summarized
	FunctionName string         `json:"function_name"`
	Language     types.Language `json:"language"`

	// What the function does (descriptive only)
	Description string `json:"description"`

	// Parameters (descriptive only)
	Parameters []ParameterInfo `json:"parameters,omitempty"`

	// Return value (descriptive only)
	Returns string `json:"returns,omitempty"`

	// Edge Cases - DETECTED but NOT fixed (informational only)
	EdgeCases []string `json:"edge_cases,omitempty"`

	// Important behavioral notes (purely informational)
	Notes []string `json:"notes,omitempty"`
}

// ParameterInfo describes a function parameter
type ParameterInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"` // What it's used for
}

// APISummary represents a high-level summary of an API
type APISummary struct {
	// Package/Module name
	Name string `json:"name"`

	// Overall description
	Description string `json:"description"`

	// Functions in this API
	Functions []*FunctionSummary `json:"functions"`

	// Types defined (names only, descriptive)
	Types []string `json:"types,omitempty"`

	// Common edge cases across the API
	CommonEdgeCases []string `json:"common_edge_cases,omitempty"`
}

// NewSummarizer creates a new documentation summarizer
// This component is PURELY DESCRIPTIVE - it NEVER suggests or changes code
func NewSummarizer(config SummarizerConfig, modelClient *model.Client) *Summarizer {
	if config.MaxSummaryLength == 0 {
		config.MaxSummaryLength = 200
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 24 * time.Hour
	}
	if config.Model == "" {
		config.Model = model.Gemma2B // Fast model for summaries
	}

	return &Summarizer{
		config:      config,
		modelClient: modelClient,
		stringUtils: &common.StringUtils{},
		logger:      common.GetLogger().WithField("component", "doc-summarizer"),
		cache: &SummaryCache{
			entries: make(map[string]*CacheEntry),
			ttl:     config.CacheTTL,
		},
	}
}

// SummarizeFunction generates a human-readable description of what a function does
// This includes detection of edge cases - but NEVER suggests fixes
func (s *Summarizer) SummarizeFunction(ctx context.Context, fn *types.Function) (*FunctionSummary, error) {
	s.logger.Debug("Generating summary for function", "function", fn.Name, "language", fn.Language)

	// Check cache
	if s.config.UseCache {
		if cached := s.getCachedSummary(fn.Name, fn.Language); cached != nil {
			s.logger.Debug("Using cached summary", "function", fn.Name)
			return cached, nil
		}
	}

	// Build prompt - PURELY DESCRIPTIVE, DETECTS edge cases but DOES NOT suggest fixes
	prompt := s.buildDescriptionPrompt(fn)

	// Generate summary
	resp, err := s.modelClient.Generate(ctx, model.ModelRequest{
		Model:       s.config.Model,
		Language:    fn.Language,
		Prompt:      prompt,
		Temperature: 0.3, // Lower temperature for consistent descriptions
		MaxTokens:   400,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	// Parse response - ONLY extract descriptive information
	summary, err := s.parseDescription(resp.Response, fn)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if s.config.UseCache {
		s.cacheSummary(fn.Name, fn.Language, summary)
	}

	return summary, nil
}

// SummarizeAPI generates a high-level description of an API
// This includes aggregated edge cases - but NEVER suggests fixes
func (s *Summarizer) SummarizeAPI(ctx context.Context, functions []*types.Function) (*APISummary, error) {
	s.logger.Info("Generating API summary", "functions", len(functions))

	if len(functions) == 0 {
		return &APISummary{}, nil
	}

	// Build package/module name from context
	name := s.extractPackageName(functions)

	// Generate overall description
	description, err := s.generateAPIDescription(ctx, functions)
	if err != nil {
		s.logger.Warn("Failed to generate API description", "error", err)
	}

	// Summarize each function
	var functionSummaries []*FunctionSummary
	var allEdgeCases []string
	edgeCaseCount := make(map[string]int)

	for _, fn := range functions {
		summary, err := s.SummarizeFunction(ctx, fn)
		if err != nil {
			s.logger.Warn("Failed to summarize function", "function", fn.Name, "error", err)
			continue
		}
		functionSummaries = append(functionSummaries, summary)

		// Collect edge cases for API-level view
		for _, ec := range summary.EdgeCases {
			edgeCaseCount[ec]++
			if edgeCaseCount[ec] == 1 {
				allEdgeCases = append(allEdgeCases, ec)
			}
		}
	}

	// Extract type names (just names, no implementation)
	types := s.extractTypeNames(functions)

	// Identify common edge cases (appear in multiple functions)
	var commonEdgeCases []string
	for ec, count := range edgeCaseCount {
		if count > 1 {
			commonEdgeCases = append(commonEdgeCases, fmt.Sprintf("%s (affects %d functions)", ec, count))
		}
	}

	return &APISummary{
		Name:            name,
		Description:     description,
		Functions:       functionSummaries,
		Types:           types,
		CommonEdgeCases: commonEdgeCases,
	}, nil
}

// Helper methods - ALL PURELY DESCRIPTIVE

func (s *Summarizer) buildDescriptionPrompt(fn *types.Function) string {
	var sb strings.Builder

	sb.WriteString("Describe what this function does. DO NOT suggest improvements or changes.\n\n")
	sb.WriteString(fmt.Sprintf("Function Name: %s\n", fn.Name))
	sb.WriteString(fmt.Sprintf("Language: %s\n", fn.Language))
	sb.WriteString("\nCode:\n")
	sb.WriteString(fn.Content)
	sb.WriteString("\n\n")
	sb.WriteString("Provide a clear, factual description including:\n")
	sb.WriteString("- What the function does (purpose)\n")
	sb.WriteString("- What each parameter is used for\n")
	sb.WriteString("- What it returns\n")

	if s.config.DetectEdgeCases {
		sb.WriteString("- ANY EDGE CASES you observe in the code (e.g., zero values, empty inputs, boundary conditions, error cases)\n")
		sb.WriteString("  List these as observations, NOT as suggestions to fix them\n")
	}

	sb.WriteString("- Any other important behavioral notes\n")
	sb.WriteString("\n")
	sb.WriteString("IMPORTANT: Only describe what the code does. If you see edge cases, note them as observations.\n")
	sb.WriteString("DO NOT say 'should', 'could', 'might', 'consider', or suggest changes of any kind.")

	return sb.String()
}

func (s *Summarizer) parseDescription(response string, fn *types.Function) (*FunctionSummary, error) {
	summary := &FunctionSummary{
		FunctionName: fn.Name,
		Language:     fn.Language,
		Parameters:   []ParameterInfo{},
		EdgeCases:    []string{},
		Notes:        []string{},
	}

	// Split into paragraphs
	paragraphs := strings.Split(response, "\n\n")
	if len(paragraphs) > 0 {
		summary.Description = strings.TrimSpace(paragraphs[0])
	}

	// Look for parameter descriptions, edge cases, and notes
	for _, p := range paragraphs {
		lines := strings.Split(p, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Look for parameter descriptions (just descriptive, no suggestions)
			if strings.Contains(strings.ToLower(line), "parameter") ||
				strings.Contains(strings.ToLower(line), "param") ||
				strings.Contains(strings.ToLower(line), "argument") ||
				strings.Contains(line, ":") && len(strings.Split(line, ":")) == 2 {

				s.parseParameterLine(line, summary)
				continue
			}

			// Look for return value description
			if strings.Contains(strings.ToLower(line), "return") {
				cleanLine := strings.TrimPrefix(line, "-")
				cleanLine = strings.TrimPrefix(cleanLine, "•")
				cleanLine = strings.TrimSpace(cleanLine)
				summary.Returns = cleanLine
				continue
			}

			// Look for edge cases (DETECT but DON'T suggest fixes)
			if s.config.DetectEdgeCases && s.isEdgeCase(line) {
				cleanLine := strings.TrimPrefix(line, "-")
				cleanLine = strings.TrimPrefix(cleanLine, "•")
				cleanLine = strings.TrimSpace(cleanLine)

				// Add as observation, NOT suggestion
				if !strings.Contains(strings.ToLower(cleanLine), "should") &&
					!strings.Contains(strings.ToLower(cleanLine), "could") &&
					!strings.Contains(strings.ToLower(cleanLine), "might") &&
					!strings.Contains(strings.ToLower(cleanLine), "consider") &&
					!strings.Contains(strings.ToLower(cleanLine), "suggest") &&
					!strings.Contains(strings.ToLower(cleanLine), "fix") &&
					!strings.Contains(strings.ToLower(cleanLine), "improve") {

					// Format as observation, not suggestion
					if !strings.HasPrefix(cleanLine, "Edge case:") &&
						!strings.HasPrefix(cleanLine, "Note:") {
						cleanLine = "Edge case: " + cleanLine
					}
					summary.EdgeCases = append(summary.EdgeCases, cleanLine)
				}
				continue
			}

			// Look for important notes (purely informational)
			if strings.Contains(strings.ToLower(line), "note") ||
				strings.Contains(strings.ToLower(line), "important") ||
				strings.Contains(strings.ToLower(line), "behavior") ||
				strings.Contains(strings.ToLower(line), "caution") ||
				strings.Contains(strings.ToLower(line), "warning") {

				cleanLine := strings.TrimPrefix(line, "-")
				cleanLine = strings.TrimPrefix(cleanLine, "•")
				cleanLine = strings.TrimSpace(cleanLine)

				// Only add if it's descriptive, not suggestive
				if !strings.Contains(strings.ToLower(cleanLine), "should") &&
					!strings.Contains(strings.ToLower(cleanLine), "could") &&
					!strings.Contains(strings.ToLower(cleanLine), "might") &&
					!strings.Contains(strings.ToLower(cleanLine), "consider") &&
					!strings.Contains(strings.ToLower(cleanLine), "suggest") {
					summary.Notes = append(summary.Notes, cleanLine)
				}
			}
		}
	}

	// Truncate description if too long
	if len(summary.Description) > s.config.MaxSummaryLength {
		summary.Description = summary.Description[:s.config.MaxSummaryLength-3] + "..."
	}

	return summary, nil
}

// parseParameterLine extracts parameter information from a line
func (s *Summarizer) parseParameterLine(line string, summary *FunctionSummary) {
	// Remove common prefixes
	cleanLine := strings.TrimPrefix(line, "-")
	cleanLine = strings.TrimPrefix(cleanLine, "•")
	cleanLine = strings.TrimSpace(cleanLine)

	// Try different formats
	parts := strings.SplitN(cleanLine, ":", 2)
	if len(parts) == 2 {
		paramName := strings.TrimSpace(parts[0])
		paramDesc := strings.TrimSpace(parts[1])

		// Extract type if present in parentheses
		if strings.Contains(paramName, "(") {
			// Parse "name (type)" format
			nameParts := strings.SplitN(paramName, "(", 2)
			if len(nameParts) == 2 {
				paramName = strings.TrimSpace(nameParts[0])
				paramType := strings.TrimSuffix(strings.TrimSpace(nameParts[1]), ")")

				summary.Parameters = append(summary.Parameters, ParameterInfo{
					Name:        paramName,
					Type:        paramType,
					Description: paramDesc,
				})
				return
			}
		}

		// Try "name - type - description" format
		dashParts := strings.Split(paramName, "-")
		if len(dashParts) >= 2 {
			paramName = strings.TrimSpace(dashParts[0])
			paramType := strings.TrimSpace(dashParts[1])
			summary.Parameters = append(summary.Parameters, ParameterInfo{
				Name:        paramName,
				Type:        paramType,
				Description: paramDesc,
			})
			return
		}

		// Simple format
		summary.Parameters = append(summary.Parameters, ParameterInfo{
			Name:        paramName,
			Description: paramDesc,
		})
	}
}

// isEdgeCase detects if a line describes an edge case
func (s *Summarizer) isEdgeCase(line string) bool {
	lowerLine := strings.ToLower(line)

	edgeCaseIndicators := []string{
		"edge case",
		"boundary",
		"zero",
		"empty",
		"nil",
		"null",
		"negative",
		"invalid",
		"error case",
		"exception",
		"limit",
		"maximum",
		"minimum",
		"overflow",
		"underflow",
		"corner case",
		"special case",
		"when empty",
		"when zero",
		"when nil",
		"when null",
		"when negative",
		"if empty",
		"if zero",
		"if nil",
		"handles empty",
		"handles zero",
		"handles nil",
	}

	for _, indicator := range edgeCaseIndicators {
		if strings.Contains(lowerLine, indicator) {
			return true
		}
	}

	return false
}

func (s *Summarizer) generateAPIDescription(ctx context.Context, functions []*types.Function) (string, error) {
	// Build a list of function names
	var names []string
	for _, fn := range functions {
		names = append(names, fn.Name)
	}

	prompt := fmt.Sprintf(`Based on these function names, describe what this API does in one sentence.
DO NOT suggest improvements - just describe the purpose.

Functions: %s

Provide a brief, one-sentence description.`,
		strings.Join(names, ", "))

	resp, err := s.modelClient.Generate(ctx, model.ModelRequest{
		Model:       s.config.Model,
		Prompt:      prompt,
		Temperature: 0.3,
		MaxTokens:   100,
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp.Response), nil
}

func (s *Summarizer) extractPackageName(functions []*types.Function) string {
	if len(functions) == 0 {
		return "Package"
	}

	// Extract from file paths (just naming, no changes)
	paths := make(map[string]bool)
	for _, fn := range functions {
		parts := strings.Split(fn.FilePath, "/")
		if len(parts) > 1 {
			paths[parts[len(parts)-2]] = true
		}
	}

	if len(paths) == 1 {
		for p := range paths {
			return p
		}
	}

	return "Package"
}

func (s *Summarizer) extractTypeNames(functions []*types.Function) []string {
	// Extract type names for documentation - PURELY DESCRIPTIVE
	var types []string
	seen := make(map[string]bool)

	for _, fn := range functions {
		// Look for type names in signature (just names, no implementation)
		words := strings.Fields(fn.Signature)
		for _, w := range words {
			if strings.Contains(w, "struct") ||
				strings.Contains(w, "interface") ||
				strings.Contains(w, "type") {
				// Extract just the name
				if strings.Contains(w, ".") {
					parts := strings.Split(w, ".")
					w = parts[len(parts)-1]
				}
				if !seen[w] {
					seen[w] = true
					types = append(types, w)
				}
			}
		}
	}

	return types
}

func (s *Summarizer) getCachedSummary(functionName string, language types.Language) *FunctionSummary {
	if !s.config.UseCache {
		return nil
	}

	key := fmt.Sprintf("%s:%s", language, functionName)
	if entry, ok := s.cache.entries[key]; ok {
		if time.Since(entry.Timestamp) < s.cache.ttl {
			return &FunctionSummary{
				FunctionName: entry.Function,
				Language:     entry.Language,
				Description:  entry.Summary,
			}
		}
		// Expired
		delete(s.cache.entries, key)
	}
	return nil
}

func (s *Summarizer) cacheSummary(functionName string, language types.Language, summary *FunctionSummary) {
	if !s.config.UseCache || summary == nil {
		return
	}

	key := fmt.Sprintf("%s:%s", language, functionName)
	s.cache.entries[key] = &CacheEntry{
		Summary:   summary.Description,
		Timestamp: time.Now(),
		Function:  functionName,
		Language:  language,
	}
}

// FormatSummary formats a function summary for display in documentation
// This is PURELY FORMATTING - no content generation
func (s *Summarizer) FormatSummary(summary *FunctionSummary) string {
	if summary == nil {
		return ""
	}

	var sb strings.Builder

	// Function name and description
	sb.WriteString(fmt.Sprintf("### `%s`\n\n", summary.FunctionName))

	if summary.Description != "" {
		sb.WriteString(summary.Description)
		sb.WriteString("\n\n")
	}

	// Parameters table (if any)
	if len(summary.Parameters) > 0 {
		sb.WriteString("**Parameters:**\n\n")
		sb.WriteString("| Name | Type | Description |\n")
		sb.WriteString("|------|------|-------------|\n")

		for _, p := range summary.Parameters {
			typeStr := p.Type
			if typeStr == "" {
				typeStr = "-"
			}
			sb.WriteString(fmt.Sprintf("| `%s` | `%s` | %s |\n",
				p.Name, typeStr, p.Description))
		}
		sb.WriteString("\n")
	}

	// Return value (if any)
	if summary.Returns != "" {
		sb.WriteString(fmt.Sprintf("**Returns:** %s\n\n", summary.Returns))
	}

	// Edge Cases - DETECTED, not fixed (informational only)
	if len(summary.EdgeCases) > 0 {
		sb.WriteString("**⚠️ Edge Cases to be aware of:**\n")
		for _, ec := range summary.EdgeCases {
			sb.WriteString(fmt.Sprintf("- %s\n", ec))
		}
		sb.WriteString("\n")
	}

	// Important notes (if any)
	if len(summary.Notes) > 0 {
		sb.WriteString("**📝 Notes:**\n")
		for _, note := range summary.Notes {
			sb.WriteString(fmt.Sprintf("- %s\n", note))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatAPISummary formats an API summary for display
// This is PURELY FORMATTING - no content generation
func (s *Summarizer) FormatAPISummary(summary *APISummary) string {
	if summary == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", summary.Name))

	if summary.Description != "" {
		sb.WriteString(summary.Description)
		sb.WriteString("\n\n")
	}

	if len(summary.Types) > 0 {
		sb.WriteString("**Types:**\n")
		for _, t := range summary.Types {
			sb.WriteString(fmt.Sprintf("- `%s`\n", t))
		}
		sb.WriteString("\n")
	}

	// Common edge cases across the API
	if len(summary.CommonEdgeCases) > 0 {
		sb.WriteString("**⚠️ Common Edge Cases Across the API:**\n")
		for _, ec := range summary.CommonEdgeCases {
			sb.WriteString(fmt.Sprintf("- %s\n", ec))
		}
		sb.WriteString("\n")
	}

	if len(summary.Functions) > 0 {
		sb.WriteString("## Functions\n\n")
		for _, fn := range summary.Functions {
			sb.WriteString(s.FormatSummary(fn))
			sb.WriteString("---\n\n")
		}
	}

	return sb.String()
}
