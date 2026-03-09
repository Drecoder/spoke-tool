package doc

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/internal/common"
)

// Formatter handles formatting extracted content into markdown
type Formatter struct {
	config       FormatterConfig
	stringUtils  *common.StringUtils
	templateFunc template.FuncMap
}

// FormatterConfig configures the formatter
type FormatterConfig struct {
	// Whether to include badges
	IncludeBadges bool
	
	// Whether to include a table of contents
	IncludeTOC bool
	
	// Whether to include timestamps
	IncludeTimestamps bool
	
	// Maximum line length for examples
	MaxLineLength int
	
	// Code block style (``` or ~~~)
	CodeBlockStyle string
	
	// Heading style (## or ---)
	HeadingStyle string
	
	// Whether to add emojis to headings
	AddEmojis bool
	
	// Custom templates per section
	Templates map[types.DocSection]string
}

// Section represents a formatted README section
type Section struct {
	Type     types.DocSection
	Title    string
	Content  string
	Order    int
	Children []*Section
}

// FormattedReadme represents a complete formatted README
type FormattedReadme struct {
	Title       string
	Description string
	Badges      []string
	Sections    []*Section
	TOC         string
	Content     string
}

// NewFormatter creates a new documentation formatter
func NewFormatter(config FormatterConfig) *Formatter {
	if config.MaxLineLength == 0 {
		config.MaxLineLength = 80
	}
	if config.CodeBlockStyle == "" {
		config.CodeBlockStyle = "```"
	}
	if config.HeadingStyle == "" {
		config.HeadingStyle = "##"
	}

	return &Formatter{
		config:      config,
		stringUtils: &common.StringUtils{},
		templateFunc: template.FuncMap{
			"lower":    strings.ToLower,
			"upper":    strings.ToUpper,
			"title":    strings.Title,
			"trim":     strings.TrimSpace,
			"now":      common.Times.FormatDuration,
			"code":     e.formatCodeBlock,
			"table":    e.formatTable,
			"badge":    e.formatBadge,
			"link":     e.formatLink,
			"plural":   e.pluralize,
		},
	}
}

// FormatReadme formats a complete README from extracted content
func (f *Formatter) FormatReadme(projectName string, description string, sections []*Section) (*FormattedReadme, error) {
	readme := &FormattedReadme{
		Title:       projectName,
		Description: description,
		Sections:    sections,
	}

	// Add badges
	if f.config.IncludeBadges {
		readme.Badges = f.generateBadges(projectName)
	}

	// Generate table of contents
	if f.config.IncludeTOC {
		readme.TOC = f.generateTOC(sections)
	}

	// Generate full content
	content, err := f.assembleContent(readme)
	if err != nil {
		return nil, err
	}
	readme.Content = content

	return readme, nil
}

// FormatSection formats a single section
func (f *Formatter) FormatSection(sectionType types.DocSection, title string, content interface{}) (*Section, error) {
	var formatted string
	var err error

	// Check for custom template
	if template, ok := f.config.Templates[sectionType]; ok {
		formatted, err = f.executeTemplate(template, content)
	} else {
		formatted, err = f.formatDefaultSection(sectionType, content)
	}

	if err != nil {
		return nil, err
	}

	return &Section{
		Type:    sectionType,
		Title:   f.formatHeading(sectionType, title),
		Content: formatted,
		Order:   f.getSectionOrder(sectionType),
	}, nil
}

// FormatExamples formats a collection of examples
func (f *Formatter) FormatExamples(examples []*ExtractedExample, language types.Language) (string, error) {
	if len(examples) == 0 {
		return "", nil
	}

	var sb strings.Builder

	// Add section header
	sb.WriteString("## Examples\n\n")

	for i, example := range examples {
		// Add example description
		if example.Description != "" {
			sb.WriteString(fmt.Sprintf("### %s\n\n", example.Description))
		}

		// Add the code
		code := f.cleanExampleForDisplay(example)
		sb.WriteString(f.formatCodeBlock(code, string(language)))
		sb.WriteString("\n")

		// Add expected output if available
		if example.ExpectedOutput != "" {
			sb.WriteString("**Expected output:**\n\n")
			sb.WriteString(f.formatCodeBlock(example.ExpectedOutput, ""))
			sb.WriteString("\n")
		}

		// Add separator between examples
		if i < len(examples)-1 {
			sb.WriteString("---\n\n")
		}
	}

	return sb.String(), nil
}

// FormatAPI formats API documentation
func (f *Formatter) FormatAPI(functions []*ExtractedContent) (string, error) {
	if len(functions) == 0 {
		return "", nil
	}

	var sb strings.Builder

	sb.WriteString("## API Reference\n\n")

	for _, fn := range functions {
		if fn.Function == nil {
			continue
		}

		// Function signature
		sb.WriteString(fmt.Sprintf("### `%s`\n\n", fn.Function.Name))

		// Description from comments
		if len(fn.Comments) > 0 {
			sb.WriteString(fn.Comments[0])
			sb.WriteString("\n\n")
		}

		// Parameters table
		if params := f.extractParameters(fn); len(params) > 0 {
			sb.WriteString("**Parameters:**\n\n")
			sb.WriteString(f.formatParameterTable(params))
			sb.WriteString("\n")
		}

		// Return value
		if returns := f.extractReturnValue(fn); returns != "" {
			sb.WriteString(fmt.Sprintf("**Returns:** `%s`\n\n", returns))
		}

		// Examples
		if len(fn.Examples) > 0 {
			sb.WriteString("**Example:**\n\n")
			example := fn.Examples[0] // Take the first example
			sb.WriteString(f.formatCodeBlock(example.Code, string(fn.Language)))
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	return sb.String(), nil
}

// FormatInstallation formats installation instructions
func (f *Formatter) FormatInstallation(projectName string, languages []types.Language) (string, error) {
	var sb strings.Builder

	sb.WriteString("## Installation\n\n")

	hasGo := false
	hasNode := false
	hasPython := false

	for _, lang := range languages {
		switch lang {
		case types.Go:
			hasGo = true
		case types.NodeJS:
			hasNode = true
		case types.Python:
			hasPython = true
		}
	}

	if hasGo {
		sb.WriteString("### Go\n")
		sb.WriteString(f.formatCodeBlock(fmt.Sprintf("go get %s", projectName), "bash"))
		sb.WriteString("\n")
	}

	if hasNode {
		sb.WriteString("### Node.js\n")
		sb.WriteString(f.formatCodeBlock(fmt.Sprintf("npm install %s", projectName), "bash"))
		sb.WriteString("\n")
	}

	if hasPython {
		sb.WriteString("### Python\n")
		sb.WriteString(f.formatCodeBlock(fmt.Sprintf("pip install %s", projectName), "bash"))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// FormatQuickStart formats a quick start guide
func (f *Formatter) FormatQuickStart(examples []*ExtractedExample) (string, error) {
	if len(examples) == 0 {
		return "", nil
	}

	var sb strings.Builder

	sb.WriteString("## Quick Start\n\n")

	// Take the simplest example
	var bestExample *ExtractedExample
	for _, ex := range examples {
		if !ex.IsEdgeCase && len(ex.Code) < 200 { // Prefer simple, short examples
			bestExample = ex
			break
		}
	}

	if bestExample == nil && len(examples) > 0 {
		bestExample = examples[0]
	}

	if bestExample != nil {
		sb.WriteString("Here's a simple example to get you started:\n\n")
		sb.WriteString(f.formatCodeBlock(bestExample.Code, string(bestExample.Language)))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// FormatContributing formats contributing guidelines
func (f *Formatter) FormatContributing() (string, error) {
	return `## Contributing

We love contributions! Here's how you can help:

1. Fork the repository
2. Create a feature branch
3. Write tests for your changes
4. Run the test suite
5. Submit a pull request

### Development Setup

\`\`\`bash
# Clone the repository
git clone https://github.com/yourusername/project.git

# Install dependencies
# [Language-specific instructions]

# Run tests
# [Test command]
\`\`\`

### Code Style

- Follow language-specific style guides
- Write tests for new features
- Update documentation as needed
- Keep pull requests focused

### Reporting Issues

- Check existing issues first
- Include reproduction steps
- Provide example code if possible
- Mention your environment

Thanks for contributing! 🎉
`, nil
}

// FormatLicense formats license information
func (f *Formatter) FormatLicense(license string) (string, error) {
	if license == "" {
		license = "MIT"
	}

	return fmt.Sprintf(`## License

This project is licensed under the %s License - see the [LICENSE](LICENSE) file for details.
`, license), nil
}

// Helper methods

func (f *Formatter) formatHeading(sectionType types.DocSection, title string) string {
	if title == "" {
		switch sectionType {
		case types.DocSectionTitle:
			title = "Project"
		case types.DocSectionInstallation:
			title = "Installation"
		case types.DocSectionQuickStart:
			title = "Quick Start"
		case types.DocSectionAPI:
			title = "API Reference"
		case types.DocSectionExamples:
			title = "Examples"
		case types.DocSectionContributing:
			title = "Contributing"
		case types.DocSectionLicense:
			title = "License"
		default:
			title = string(sectionType)
		}
	}

	if f.config.AddEmojis {
		switch sectionType {
		case types.DocSectionInstallation:
			title = "📦 " + title
		case types.DocSectionQuickStart:
			title = "🚀 " + title
		case types.DocSectionAPI:
			title = "📚 " + title
		case types.DocSectionExamples:
			title = "💡 " + title
		case types.DocSectionContributing:
			title = "🤝 " + title
		case types.DocSectionLicense:
			title = "⚖️ " + title
		}
	}

	return fmt.Sprintf("%s %s", f.config.HeadingStyle, title)
}

func (f *Formatter) formatCodeBlock(code string, language string) string {
	if code == "" {
		return ""
	}

	// Clean up the code
	code = strings.TrimSpace(code)

	// Limit line length
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if len(line) > f.config.MaxLineLength {
			// Don't break code lines, just warn in comment
			lines[i] = line + "  // line too long"
		}
	}
	code = strings.Join(lines, "\n")

	return fmt.Sprintf("%s%s\n%s\n%s\n", 
		f.config.CodeBlockStyle, 
		language, 
		code, 
		f.config.CodeBlockStyle)
}

func (f *Formatter) formatTable(headers []string, rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	var sb strings.Builder

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Header row
	sb.WriteString("|")
	for i, h := range headers {
		sb.WriteString(fmt.Sprintf(" %-*s |", widths[i], h))
	}
	sb.WriteString("\n")

	// Separator row
	sb.WriteString("|")
	for _, w := range widths {
		sb.WriteString(strings.Repeat("-", w+2))
		sb.WriteString("|")
	}
	sb.WriteString("\n")

	// Data rows
	for _, row := range rows {
		sb.WriteString("|")
		for i, cell := range row {
			sb.WriteString(fmt.Sprintf(" %-*s |", widths[i], cell))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (f *Formatter) formatBadge(label, value, color string) string {
	return fmt.Sprintf("![%s](https://img.shields.io/badge/%s-%s-%s)", 
		label, 
		strings.ReplaceAll(label, " ", "_"),
		strings.ReplaceAll(value, " ", "_"),
		color)
}

func (f *Formatter) formatLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func (f *Formatter) pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func (f *Formatter) generateBadges(projectName string) []string {
	badges := []string{
		f.formatBadge("license", "MIT", "green"),
		f.formatBadge("version", "v1.0.0", "blue"),
		f.formatBadge("build", "passing", "brightgreen"),
		f.formatBadge("coverage", "80%", "yellow"),
	}
	return badges
}

func (f *Formatter) generateTOC(sections []*Section) string {
	var sb strings.Builder

	sb.WriteString("## Table of Contents\n\n")

	for _, section := range sections {
		if section.Type == types.DocSectionTitle {
			continue
		}
		link := strings.ToLower(strings.ReplaceAll(section.Title, " ", "-"))
		link = strings.ReplaceAll(link, "##", "")
		link = strings.TrimSpace(link)
		sb.WriteString(fmt.Sprintf("- [%s](#%s)\n", section.Title, link))
	}

	return sb.String()
}

func (f *Formatter) assembleContent(readme *FormattedReadme) (string, error) {
	var sb strings.Builder

	// Title
	sb.WriteString(fmt.Sprintf("# %s\n\n", readme.Title))

	// Description
	if readme.Description != "" {
		sb.WriteString(readme.Description)
		sb.WriteString("\n\n")
	}

	// Badges
	if len(readme.Badges) > 0 {
		for _, badge := range readme.Badges {
			sb.WriteString(badge)
			sb.WriteString(" ")
		}
		sb.WriteString("\n\n")
	}

	// TOC
	if readme.TOC != "" {
		sb.WriteString(readme.TOC)
		sb.WriteString("\n")
	}

	// Sections
	for _, section := range readme.Sections {
		if section.Content != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", section.Title))
			sb.WriteString(section.Content)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

func (f *Formatter) formatDefaultSection(sectionType types.DocSection, content interface{}) (string, error) {
	switch sectionType {
	case types.DocSectionExamples:
		if examples, ok := content.([]*ExtractedExample); ok {
			return f.FormatExamples(examples, "")
		}
	case types.DocSectionAPI:
		if functions, ok := content.([]*ExtractedContent); ok {
			return f.FormatAPI(functions)
		}
	case types.DocSectionInstallation:
		if projectInfo, ok := content.(map[string]interface{}); ok {
			projectName, _ := projectInfo["name"].(string)
			languages, _ := projectInfo["languages"].([]types.Language)
			return f.FormatInstallation(projectName, languages)
		}
	case types.DocSectionQuickStart:
		if examples, ok := content.([]*ExtractedExample); ok {
			return f.FormatQuickStart(examples)
		}
	case types.DocSectionContributing:
		return f.FormatContributing()
	case types.DocSectionLicense:
		if license, ok := content.(string); ok {
			return f.FormatLicense(license)
		}
	}

	return "", fmt.Errorf("unsupported section type: %s", sectionType)
}

func (f *Formatter) executeTemplate(tmpl string, data interface{}) (string, error) {
	t, err := template.New("section").Funcs(f.templateFunc).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (f *Formatter) getSectionOrder(sectionType types.DocSection) int {
	order := map[types.DocSection]int{
		types.DocSectionTitle:        1,
		types.DocSectionInstallation: 2,
		types.DocSectionQuickStart:   3,
		types.DocSectionAPI:          4,
		types.DocSectionExamples:     5,
		types.DocSectionContributing: 6,
		types.DocSectionLicense:      7,
	}
	if val, ok := order[sectionType]; ok {
		return val
	}
	return 99
}

func (f *Formatter) cleanExampleForDisplay(example *ExtractedExample) string {
	code := example.Code

	// Remove test-specific code
	code = strings.ReplaceAll(code, "t.Errorf", "// verify")
	code = strings.ReplaceAll(code, "assert.Equal", "// check")
	code = strings.ReplaceAll(code, "assert.", "// ")
	code = strings.ReplaceAll(code, "expect(", "// ")
	code = strings.ReplaceAll(code, ".toBe(", "// ")

	// Remove test function wrappers
	lines := strings.Split(code, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "func Test") {
			continue
		}
		if strings.HasPrefix(line, "def test_") {
			continue
		}
		if strings.HasPrefix(line, "test(") || strings.HasPrefix(line, "it(") {
			continue
		}
		if line == "}" || line == "});" {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	return strings.Join(cleanLines, "\n")
}

func (f *Formatter) extractParameters(content *ExtractedContent) []map[string]string {
	// This would parse the function signature to extract parameters
	// Simplified version
	var params []map[string]string

	// Look for parameter patterns in comments
	for _, comment := range content.Comments {
		if strings.Contains(comment, "@param") {
			// Parse JSDoc style param
			re := regexp.MustCompile(`@param\s+{([^}]+)}\s+(\w+)\s+-\s+(.+)`)
			matches := re.FindAllStringSubmatch(comment, -1)
			for _, match := range matches {
				if len(match) > 3 {
					params = append(params, map[string]string{
						"name":        match[2],
						"type":        match[1],
						"description": match[3],
					})
				}
			}
		} else if strings.Contains(comment, "Args:") {
			// Parse Google style Python docstring
			lines := strings.Split(comment, "\n")
			inArgs := false
			for _, line := range lines {
				if strings.Contains(line, "Args:") {
					inArgs = true
					continue
				}
				if inArgs {
					if strings.TrimSpace(line) == "" {
						break
					}
					// Parse "name (type): description"
					re := regexp.MustCompile(`\s*(\w+)\s+\(([^)]+)\):\s*(.+)`)
					if matches := re.FindStringSubmatch(line); len(matches) > 3 {
						params = append(params, map[string]string{
							"name":        matches[1],
							"type":        matches[2],
							"description": matches[3],
						})
					}
				}
			}
		}
	}

	return params
}

func (f *Formatter) extractReturnValue(content *ExtractedContent) string {
	// Look for return type in comments
	for _, comment := range content.Comments {
		if strings.Contains(comment, "@returns") {
			re := regexp.MustCompile(`@returns?\s+{([^}]+)}`)
			if matches := re.FindStringSubmatch(comment); len(matches) > 1 {
				return matches[1]
			}
		} else if strings.Contains(comment, "Returns:") {
			re := regexp.MustCompile(`Returns:\s*(.+)`)
			if matches := re.FindStringSubmatch(comment); len(matches) > 1 {
				return matches[1]
			}
		}
	}
	return ""
}

func (f *Formatter) formatParameterTable(params []map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	headers := []string{"Parameter", "Type", "Description"}
	var rows [][]string

	for _, p := range params {
		rows = append(rows, []string{
			p["name"],
			p["type"],
			p["description"],
		})
	}

	return f.formatTable(headers, rows)
}