package generators

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/yourusername/spoke-tool/api/types"
)

// DocTemplates provides language-specific documentation templates
// These templates are PURELY FOR GENERATING DOCUMENTATION - no code changes
type DocTemplates struct {
	// Template functions available to all templates
	funcMap template.FuncMap
}

// TemplateData represents data passed to documentation templates
type TemplateData struct {
	// Project information
	ProjectName        string
	ProjectDescription string
	Version            string

	// Language information
	Language     types.Language
	LanguageName string

	// Function being documented
	FunctionName string
	FunctionSig  string
	FunctionCode string

	// Parameters
	Parameters []ParamData

	// Return values
	Returns []ReturnData

	// Examples
	Examples []ExampleData

	// Notes
	Notes []string

	// Edge cases detected
	EdgeCases []string

	// Package information
	PackageName string
	PackagePath string

	// Type information (for structs/classes)
	TypeName string
	TypeKind string
	Fields   []FieldData

	// Custom data
	Extra map[string]interface{}
}

// ParamData represents a function parameter
type ParamData struct {
	Name        string
	Type        string
	Description string
	Optional    bool
	Default     string
}

// ReturnData represents a return value
type ReturnData struct {
	Name        string
	Type        string
	Description string
}

// ExampleData represents a code example
type ExampleData struct {
	Code        string
	Language    string
	Description string
	IsEdgeCase  bool
}

// FieldData represents a struct/class field
type FieldData struct {
	Name        string
	Type        string
	Description string
	Tags        map[string]string
}

// NewDocTemplates creates a new documentation templates instance
func NewDocTemplates() *DocTemplates {
	return &DocTemplates{
		funcMap: template.FuncMap{
			"lower":     strings.ToLower,
			"upper":     strings.ToUpper,
			"title":     strings.Title,
			"trim":      strings.TrimSpace,
			"join":      strings.Join,
			"split":     strings.Split,
			"indent":    indent,
			"codeBlock": formatCodeBlock,
			"table":     formatTable,
			"bullet":    formatBulletList,
			"now":       getCurrentTime,
			"hasPrefix": strings.HasPrefix,
			"hasSuffix": strings.HasSuffix,
			"replace":   strings.ReplaceAll,
		},
	}
}

// GetFunctionTemplate returns the appropriate template for documenting a function
func (d *DocTemplates) GetFunctionTemplate(lang types.Language) string {
	switch lang {
	case types.Go:
		return d.getGoFunctionTemplate()
	case types.NodeJS:
		return d.getNodeJSFunctionTemplate()
	case types.Python:
		return d.getPythonFunctionTemplate()
	default:
		return d.getDefaultFunctionTemplate()
	}
}

// GetTypeTemplate returns the appropriate template for documenting a type
func (d *DocTemplates) GetTypeTemplate(lang types.Language) string {
	switch lang {
	case types.Go:
		return d.getGoTypeTemplate()
	case types.NodeJS:
		return d.getNodeJSTypeTemplate()
	case types.Python:
		return d.getPythonTypeTemplate()
	default:
		return d.getDefaultTypeTemplate()
	}
}

// GetPackageTemplate returns the appropriate template for documenting a package
func (d *DocTemplates) GetPackageTemplate(lang types.Language) string {
	switch lang {
	case types.Go:
		return d.getGoPackageTemplate()
	case types.NodeJS:
		return d.getNodeJSPackageTemplate()
	case types.Python:
		return d.getPythonPackageTemplate()
	default:
		return d.getDefaultPackageTemplate()
	}
}

// GetAPITemplate returns a template for generating API reference
func (d *DocTemplates) GetAPITemplate() string {
	return `## API Reference

{{range .Functions}}
{{template "function" .}}
{{end}}

{{range .Types}}
{{template "type" .}}
{{end}}
`
}

// GetReadmeSectionTemplate returns a template for a README section
func (d *DocTemplates) GetReadmeSectionTemplate() string {
	return `{{if .Title}}## {{.Title}}

{{end}}{{.Content}}
`
}

// Go templates

func (d *DocTemplates) getGoFunctionTemplate() string {
	return `### ` + "`" + `{{.FunctionName}}` + "`" + `

{{if .FunctionSig}}` + "```go" + `
{{.FunctionSig}}
` + "```" + `

{{end}}{{if .Description}}{{.Description}}

{{end}}{{if .Parameters}}**Parameters:**

| Name | Type | Description |
|------|------|-------------|
{{range .Parameters}}| ` + "`" + `{{.Name}}` + "`" + ` | ` + "`" + `{{.Type}}` + "`" + ` | {{.Description}} |
{{end}}

{{end}}{{if .Returns}}**Returns:**

{{range .Returns}}- ` + "`" + `{{.Type}}` + "`" + `{{if .Description}}: {{.Description}}{{end}}
{{end}}

{{end}}{{if .Examples}}**Examples:**

{{range .Examples}}` + "```go" + `
{{.Code}}
` + "```" + `

{{end}}{{end}}{{if .EdgeCases}}**⚠️ Edge Cases:**

{{range .EdgeCases}}- {{.}}
{{end}}

{{end}}{{if .Notes}}**📝 Notes:**

{{range .Notes}}- {{.}}
{{end}}

{{end}}`
}

func (d *DocTemplates) getGoTypeTemplate() string {
	return `### {{.TypeName}} ` + "`" + `{{.TypeKind}}` + "`" + `

{{if .Description}}{{.Description}}

{{end}}{{if .Fields}}**Fields:**

| Name | Type | Description |
|------|------|-------------|
{{range .Fields}}| ` + "`" + `{{.Name}}` + "`" + ` | ` + "`" + `{{.Type}}` + "`" + ` | {{.Description}} |
{{end}}

{{end}}{{if .Methods}}**Methods:**

{{range .Methods}}- ` + "`" + `{{.FunctionSig}}` + "`" + `
{{end}}

{{end}}`
}

func (d *DocTemplates) getGoPackageTemplate() string {
	return `## Package {{.PackageName}}

{{if .Description}}{{.Description}}

{{end}}**Import Path:** ` + "`" + `{{.PackagePath}}` + "`" + `

{{if .Functions}}### Functions

{{range .Functions}}{{template "function" .}}
{{end}}{{end}}{{if .Types}}### Types

{{range .Types}}{{template "type" .}}
{{end}}{{end}}`
}

// Node.js templates

func (d *DocTemplates) getNodeJSFunctionTemplate() string {
	return `### ` + "`" + `{{.FunctionName}}()` + "`" + `

{{if .FunctionSig}}` + "```javascript" + `
{{.FunctionSig}}
` + "```" + `

{{end}}{{if .Description}}{{.Description}}

{{end}}{{if .Parameters}}**Parameters:**

| Name | Type | Description |
|------|------|-------------|
{{range .Parameters}}| ` + "`" + `{{.Name}}` + "`" + ` | ` + "`" + `{{.Type}}` + "`" + ` | {{.Description}}{{if .Optional}} (optional){{end}} |
{{end}}

{{end}}{{if .Returns}}**Returns:**

` + "`" + `{{(index .Returns 0).Type}}` + "`" + `{{if (index .Returns 0).Description}}: {{(index .Returns 0).Description}}{{end}}

{{end}}{{if .Examples}}**Examples:**

{{range .Examples}}` + "```javascript" + `
{{.Code}}
` + "```" + `

{{end}}{{end}}{{if .EdgeCases}}**⚠️ Edge Cases:**

{{range .EdgeCases}}- {{.}}
{{end}}

{{end}}{{if .Notes}}**📝 Notes:**

{{range .Notes}}- {{.}}
{{end}}

{{end}}`
}

func (d *DocTemplates) getNodeJSTypeTemplate() string {
	return `### Class ` + "`" + `{{.TypeName}}` + "`" + `

{{if .Description}}{{.Description}}

{{end}}{{if .Extends}}**Extends:** ` + "`" + `{{.Extends}}` + "`" + `

{{end}}{{if .Implements}}**Implements:** {{range .Implements}}` + "`" + `{{.}}` + "`" + ` {{end}}

{{end}}{{if .Fields}}**Properties:**

| Name | Type | Description |
|------|------|-------------|
{{range .Fields}}| ` + "`" + `{{.Name}}` + "`" + ` | ` + "`" + `{{.Type}}` + "`" + ` | {{.Description}} |
{{end}}

{{end}}{{if .Methods}}**Methods:**

{{range .Methods}}- ` + "`" + `{{.FunctionName}}()` + "`" + `{{if .Description}}: {{.Description}}{{end}}
{{end}}

{{end}}`
}

func (d *DocTemplates) getNodeJSPackageTemplate() string {
	return `## Module ` + "`" + `{{.PackageName}}` + "`" + `

{{if .Description}}{{.Description}}

{{end}}**Installation:** 
` + "```bash" + `
npm install {{.PackageName}}
` + "```" + `

{{if .Functions}}### Functions

{{range .Functions}}{{template "function" .}}
{{end}}{{end}}{{if .Types}}### Classes

{{range .Types}}{{template "type" .}}
{{end}}{{end}}`
}

// Python templates

func (d *DocTemplates) getPythonFunctionTemplate() string {
	return `### ` + "`" + `{{.FunctionName}}()` + "`" + `

{{if .FunctionSig}}` + "```python" + `
{{.FunctionSig}}
` + "```" + `

{{end}}{{if .Description}}{{.Description}}

{{end}}{{if .Parameters}}**Args:**

| Name | Type | Description |
|------|------|-------------|
{{range .Parameters}}| ` + "`" + `{{.Name}}` + "`" + ` | ` + "`" + `{{.Type}}` + "`" + ` | {{.Description}}{{if .Optional}} (default: ` + "`" + `{{.Default}}` + "`" + `){{end}} |
{{end}}

{{end}}{{if .Returns}}**Returns:**

` + "`" + `{{(index .Returns 0).Type}}` + "`" + `{{if (index .Returns 0).Description}}: {{(index .Returns 0).Description}}{{end}}

{{end}}{{if .Examples}}**Examples:**

{{range .Examples}}` + "```python" + `
{{.Code}}
` + "```" + `

{{end}}{{end}}{{if .EdgeCases}}**⚠️ Edge Cases:**

{{range .EdgeCases}}- {{.}}
{{end}}

{{end}}{{if .Notes}}**📝 Notes:**

{{range .Notes}}- {{.}}
{{end}}

{{end}}`
}

func (d *DocTemplates) getPythonTypeTemplate() string {
	return `### Class ` + "`" + `{{.TypeName}}` + "`" + `

{{if .Description}}{{.Description}}

{{end}}{{if .Bases}}**Inherits from:** {{range .Bases}}` + "`" + `{{.}}` + "`" + ` {{end}}

{{end}}{{if .Fields}}**Attributes:**

| Name | Type | Description |
|------|------|-------------|
{{range .Fields}}| ` + "`" + `{{.Name}}` + "`" + ` | ` + "`" + `{{.Type}}` + "`" + ` | {{.Description}} |
{{end}}

{{end}}{{if .Methods}}**Methods:**

{{range .Methods}}- ` + "`" + `{{.FunctionName}}()` + "`" + `{{if .Description}}: {{.Description}}{{end}}
{{end}}

{{end}}`
}

func (d *DocTemplates) getPythonPackageTemplate() string {
	return `## Module ` + "`" + `{{.PackageName}}` + "`" + `

{{if .Description}}{{.Description}}

{{end}}**Installation:** 
` + "```bash" + `
pip install {{.PackageName}}
` + "```" + `

{{if .Functions}}### Functions

{{range .Functions}}{{template "function" .}}
{{end}}{{end}}{{if .Types}}### Classes

{{range .Types}}{{template "type" .}}
{{end}}{{end}}`
}

// Default templates (fallback)

func (d *DocTemplates) getDefaultFunctionTemplate() string {
	return `### {{.FunctionName}}

{{if .Description}}{{.Description}}

{{end}}{{if .Parameters}}**Parameters:**

{{range .Parameters}}- ` + "`" + `{{.Name}}` + "`" + ` ({{.Type}}): {{.Description}}
{{end}}

{{end}}{{if .Returns}}**Returns:**

{{range .Returns}}- ` + "`" + `{{.Type}}` + "`" + `: {{.Description}}
{{end}}

{{end}}{{if .Examples}}**Examples:**

{{range .Examples}}` + "```" + `
{{.Code}}
` + "```" + `

{{end}}{{end}}`
}

func (d *DocTemplates) getDefaultTypeTemplate() string {
	return `### {{.TypeName}} ({{.TypeKind}})

{{if .Description}}{{.Description}}

{{end}}{{if .Fields}}**Fields:**

{{range .Fields}}- ` + "`" + `{{.Name}}` + "`" + ` ({{.Type}}): {{.Description}}
{{end}}

{{end}}`
}

func (d *DocTemplates) getDefaultPackageTemplate() string {
	return `## Package {{.PackageName}}

{{if .Description}}{{.Description}}

{{end}}{{if .Functions}}### Functions

{{range .Functions}}- ` + "`" + `{{.FunctionName}}` + "`" + `
{{end}}

{{end}}{{if .Types}}### Types

{{range .Types}}- ` + "`" + `{{.TypeName}}` + "`" + `
{{end}}

{{end}}`
}

// ExecuteTemplate executes a template with the given data
func (d *DocTemplates) ExecuteTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("doc").Funcs(d.funcMap).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ExecuteNamedTemplate executes a named template with the given data
func (d *DocTemplates) ExecuteNamedTemplate(tmplStr string, name string, data interface{}) (string, error) {
	tmpl, err := template.New(name).Funcs(d.funcMap).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// RegisterTemplate registers a custom template
func (d *DocTemplates) RegisterTemplate(name string, tmplStr string) error {
	_, err := template.New(name).Funcs(d.funcMap).Parse(tmplStr)
	return err
}

// Template helper functions

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(v, "\n", "\n"+pad)
}

func formatCodeBlock(code string, lang string) string {
	return "```" + lang + "\n" + code + "\n```"
}

func formatTable(headers []string, rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	var sb strings.Builder

	// Header
	sb.WriteString("|")
	for _, h := range headers {
		sb.WriteString(" " + h + " |")
	}
	sb.WriteString("\n|")
	for range headers {
		sb.WriteString(" --- |")
	}
	sb.WriteString("\n")

	// Rows
	for _, row := range rows {
		sb.WriteString("|")
		for _, cell := range row {
			sb.WriteString(" " + cell + " |")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func formatBulletList(items []string) string {
	if len(items) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, item := range items {
		sb.WriteString("- " + item + "\n")
	}
	return sb.String()
}

func getCurrentTime() string {
	// This would normally return current time
	// Keeping it simple for templates
	return "2024"
}

// Predefined template collections

// GetCompleteAPITemplate returns a complete API documentation template
func (d *DocTemplates) GetCompleteAPITemplate() string {
	return `# {{.ProjectName}} API Reference

{{if .ProjectDescription}}{{.ProjectDescription}}

{{end}}Version: {{.Version}}

{{template "api" .}}
`
}

// GetMinimalFunctionTemplate returns a minimal function template
func (d *DocTemplates) GetMinimalFunctionTemplate() string {
	return `### ` + "`" + `{{.FunctionName}}` + "`" + `

{{if .Description}}{{.Description}}

{{end}}` + "```" + `{{.Language}}` + "\n" + `{{.FunctionSig}}
` + "```" + `
`
}

// GetExamplesTemplate returns a template for examples section
func (d *DocTemplates) GetExamplesTemplate() string {
	return `## Examples

{{range .Examples}}
### {{if .Description}}{{.Description}}{{else}}Example{{end}}

` + "```" + `{{.Language}}` + "\n" + `{{.Code}}
` + "```" + `

{{end}}`
}

// GetEdgeCasesTemplate returns a template for edge cases section
func (d *DocTemplates) GetEdgeCasesTemplate() string {
	return `## ⚠️ Edge Cases

{{range .EdgeCases}}- {{.}}
{{end}}`
}

// GetNotesTemplate returns a template for notes section
func (d *DocTemplates) GetNotesTemplate() string {
	return `## 📝 Notes

{{range .Notes}}- {{.}}
{{end}}`
}

// GetTableOfContentsTemplate returns a template for table of contents
func (d *DocTemplates) GetTableOfContentsTemplate() string {
	return `## Table of Contents

{{range .Sections}}- [{{.Title}}](#{{.Anchor}})
{{end}}`
}

// GetBadgesTemplate returns a template for badges
func (d *DocTemplates) GetBadgesTemplate() string {
	return `[![Go Version](https://img.shields.io/badge/Go-{{.GoVersion}}-blue)]()
[![Node Version](https://img.shields.io/badge/Node-{{.NodeVersion}}-green)]()
[![Python Version](https://img.shields.io/badge/Python-{{.PythonVersion}}-yellow)]()
[![License](https://img.shields.io/badge/License-{{.License}}-blue)]()
[![Coverage](https://img.shields.io/badge/Coverage-{{.Coverage}}%25-brightgreen)]()
`
}
