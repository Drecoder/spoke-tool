package doc

import (
	"context"
	"fmt"
	"time"

	"example.com/spoke-tool/api/types"
	"example.com/spoke-tool/internal/common"
	"example.com/spoke-tool/internal/model"
)

// Generator handles README generation using extracted content and SLMs
type Generator struct {
	config      GeneratorConfig
	modelClient *model.Client
	fileUtils   *common.FileUtils
	logger      *common.Logger
	extractor   *Extractor
	summarizer  *Summarizer
	formatter   *Formatter
	updater     *Updater
}

// GeneratorConfig configures the README generator
type GeneratorConfig struct {
	ModelClient *model.Client
	ProjectRoot string
	Sections    []types.DocSection
	Verbose     bool
}

// NewGenerator creates a new README generator
func NewGenerator(config GeneratorConfig) *Generator {
	return &Generator{
		config:      config,
		modelClient: config.ModelClient,
		fileUtils:   &common.FileUtils{},
		logger:      common.GetLogger().WithField("component", "doc-generator"),
		extractor: NewExtractor(ExtractorConfig{
			IncludeTests:       true,
			IncludeComments:    true,
			MaxExamplesPerFunc: 3,
		}),
		summarizer: NewSummarizer(SummarizerConfig{
			Model:            model.Gemma2B,
			MaxSummaryLength: 200,
			DetectEdgeCases:  true,
			UseCache:         true,
			CacheTTL:         24 * time.Hour,
		}, config.ModelClient),
		formatter: NewFormatter(FormatterConfig{
			IncludeBadges: true,
			IncludeTOC:    true,
			AddEmojis:     true,
			MaxLineLength: 80,
		}),
		updater: NewUpdater(UpdaterConfig{
			CreateBackup:  true,
			ValidateAfter: true,
		}),
	}
}

// AnalyzeProject analyzes the project to extract documentation content
func (g *Generator) AnalyzeProject(ctx context.Context) (*types.CodeAnalysis, error) {
	g.logger.Info("Analyzing project", "root", g.config.ProjectRoot)

	// This would actually analyze the project
	// For now, return a placeholder
	return &types.CodeAnalysis{
		Files:     []types.CodeFile{},
		Functions: []types.Function{},
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// NeedsUpdate checks if the README needs to be updated
func (g *Generator) NeedsUpdate(ctx context.Context, analysis *types.CodeAnalysis) (bool, error) {
	// Check if README exists
	readmePath := g.config.ProjectRoot + "/README.md"
	if !g.fileUtils.FileExists(readmePath) {
		return true, nil
	}

	// In a real implementation, would check timestamps and content changes
	return true, nil
}

// GenerateSections generates README sections from the analysis
func (g *Generator) GenerateSections(ctx context.Context, analysis *types.CodeAnalysis) ([]*Section, error) {
	g.logger.Info("Generating README sections")

	var sections []*Section

	// Extract content from the project
	content, err := g.extractor.ExtractFromProject(ctx, analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content: %w", err)
	}

	// Generate sections based on config
	for _, sectionType := range g.config.Sections {
		var section *Section
		var err error

		switch sectionType {
		case types.DocSectionTitle:
			section, err = g.formatter.FormatSection(sectionType, "Project", nil)
		case types.DocSectionDescription:
			section, err = g.generateDescriptionSection(ctx, analysis)
		case types.DocSectionInstallation:
			section, err = g.generateInstallationSection(ctx, analysis)
		case types.DocSectionQuickStart:
			section, err = g.generateQuickStartSection(ctx, content)
		case types.DocSectionAPI:
			section, err = g.generateAPISection(ctx, content)
		case types.DocSectionExamples:
			section, err = g.generateExamplesSection(ctx, content)
		case types.DocSectionContributing:
			section, err = g.formatter.FormatSection(sectionType, "Contributing", nil)
		case types.DocSectionLicense:
			section, err = g.formatter.FormatSection(sectionType, "License", "MIT")
		default:
			continue
		}

		if err != nil {
			g.logger.Warn("Failed to generate section", "section", sectionType, "error", err)
			continue
		}

		if section != nil {
			sections = append(sections, section)
		}
	}

	return sections, nil
}

// generateDescriptionSection creates a project description
func (g *Generator) generateDescriptionSection(ctx context.Context, analysis *types.CodeAnalysis) (*Section, error) {
	description := "A powerful tool built with Go, Node.js, and Python."

	// In a real implementation, would generate this from the code
	return g.formatter.FormatSection(types.DocSectionDescription, "Description", description)
}

// generateInstallationSection creates installation instructions
func (g *Generator) generateInstallationSection(ctx context.Context, analysis *types.CodeAnalysis) (*Section, error) {
	// Detect languages present
	var languages []types.Language
	langMap := make(map[types.Language]bool)

	for _, fn := range analysis.Functions {
		if !langMap[fn.Language] {
			langMap[fn.Language] = true
			languages = append(languages, fn.Language)
		}
	}

	content, err := g.formatter.FormatInstallation("spoke-tool", languages)
	if err != nil {
		return nil, err
	}

	return &Section{
		Type:    types.DocSectionInstallation,
		Title:   "Installation",
		Content: content,
	}, nil
}

// generateQuickStartSection creates a quick start guide
func (g *Generator) generateQuickStartSection(ctx context.Context, content []*ExtractedContent) (*Section, error) {
	// Collect examples from extracted content
	var examples []*ExtractedExample
	for _, c := range content {
		examples = append(examples, c.Examples...)
	}

	quickStart, err := g.formatter.FormatQuickStart(examples)
	if err != nil {
		return nil, err
	}

	return &Section{
		Type:    types.DocSectionQuickStart,
		Title:   "Quick Start",
		Content: quickStart,
	}, nil
}

// generateAPISection creates API documentation
func (g *Generator) generateAPISection(ctx context.Context, content []*ExtractedContent) (*Section, error) {
	api, err := g.formatter.FormatAPI(content)
	if err != nil {
		return nil, err
	}

	return &Section{
		Type:    types.DocSectionAPI,
		Title:   "API Reference",
		Content: api,
	}, nil
}

// generateExamplesSection creates examples
func (g *Generator) generateExamplesSection(ctx context.Context, content []*ExtractedContent) (*Section, error) {
	// Collect all examples
	var allExamples []*ExtractedExample
	for _, c := range content {
		allExamples = append(allExamples, c.Examples...)
	}

	// Limit to a reasonable number
	if len(allExamples) > 5 {
		allExamples = allExamples[:5]
	}

	examples, err := g.formatter.FormatExamples(allExamples, types.Go)
	if err != nil {
		return nil, err
	}

	return &Section{
		Type:    types.DocSectionExamples,
		Title:   "Examples",
		Content: examples,
	}, nil
}

// AssembleReadme assembles the final README from sections
func (g *Generator) AssembleReadme(ctx context.Context, analysis *types.CodeAnalysis, sections []*Section) (*FormattedReadme, error) {
	projectName := "Spoke Tool"
	if analysis != nil && len(analysis.Files) > 0 {
		// Try to extract project name from files
	}

	return g.formatter.FormatReadme(projectName, "A local AI-powered development tool", sections)
}
