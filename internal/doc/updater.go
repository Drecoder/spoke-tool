package doc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/internal/common"
)

// Updater handles updating README files while preserving manual content
// This component NEVER overwrites manual content - it merges intelligently
type Updater struct {
	config      UpdaterConfig
	fileUtils   *common.FileUtils
	stringUtils *common.StringUtils
	logger      *common.Logger
}

// UpdaterConfig configures the updater
type UpdaterConfig struct {
	// Whether to preserve manual content (ALWAYS true in this design)
	PreserveManual bool

	// Whether to create a backup before updating
	CreateBackup bool

	// Backup directory (default: .readme-backups)
	BackupDir string

	// How to mark generated sections
	GeneratedMarker string

	// Whether to validate after update
	ValidateAfter bool

	// Sections that should NEVER be auto-updated
	ProtectedSections []types.DocSection
}

// UpdateResult represents the result of an update operation
type UpdateResult struct {
	// Path to the README file
	Path string `json:"path"`

	// Whether the file was updated
	Updated bool `json:"updated"`

	// Sections that were added
	AddedSections []string `json:"added_sections,omitempty"`

	// Sections that were updated
	UpdatedSections []string `json:"updated_sections,omitempty"`

	// Sections that were preserved (manual)
	PreservedSections []string `json:"preserved_sections,omitempty"`

	// Path to backup (if created)
	BackupPath string `json:"backup_path,omitempty"`

	// Any validation warnings
	Warnings []string `json:"warnings,omitempty"`
}

// Section represents a parsed section from a README
type ParsedSection struct {
	Type        types.DocSection
	Title       string
	Content     string
	StartLine   int
	EndLine     int
	IsManual    bool // Whether this was manually written
	IsGenerated bool // Whether this was generated
}

// NewUpdater creates a new README updater
func NewUpdater(config UpdaterConfig) *Updater {
	// These should ALWAYS be true in this design
	config.PreserveManual = true

	if config.GeneratedMarker == "" {
		config.GeneratedMarker = "<!-- GENERATED SECTION - DO NOT EDIT MANUALLY -->"
	}
	if config.BackupDir == "" {
		config.BackupDir = ".readme-backups"
	}
	if config.ProtectedSections == nil {
		config.ProtectedSections = []types.DocSection{
			types.DocSectionTitle,       // Title should be manual
			types.DocSectionDescription, // Description should be manual
			types.DocSectionLicense,     // License is usually standard
		}
	}

	return &Updater{
		config:      config,
		fileUtils:   &common.FileUtils{},
		stringUtils: &common.StringUtils{},
		logger:      common.GetLogger().WithField("component", "doc-updater"),
	}
}

// UpdateReadme updates the README file with new generated content
// This NEVER overwrites manual content
func (u *Updater) UpdateReadme(ctx context.Context, path string, generated *FormattedReadme) (*UpdateResult, error) {
	u.logger.Info("Updating README", "path", path)

	result := &UpdateResult{
		Path:              path,
		AddedSections:     []string{},
		UpdatedSections:   []string{},
		PreservedSections: []string{},
		Warnings:          []string{},
	}

	// Check if README exists
	var existingContent string
	var err error

	if u.fileUtils.FileExists(path) {
		// Read existing README
		existingContent, err = u.fileUtils.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read existing README: %w", err)
		}

		// Create backup if configured
		if u.config.CreateBackup {
			backupPath, err := u.createBackup(path, existingContent)
			if err != nil {
				u.logger.Warn("Failed to create backup", "error", err)
			} else {
				result.BackupPath = backupPath
			}
		}
	}

	// Parse existing sections
	var existingSections []*ParsedSection
	if existingContent != "" {
		existingSections = u.parseSections(existingContent)
	}

	// Merge existing and generated content
	newContent, mergeResult := u.mergeContent(existingContent, existingSections, generated)

	// Update result with merge info
	result.AddedSections = mergeResult.Added
	result.UpdatedSections = mergeResult.Updated
	result.PreservedSections = mergeResult.Preserved
	result.Updated = len(mergeResult.Added) > 0 || len(mergeResult.Updated) > 0

	// Validate if configured
	if u.config.ValidateAfter && result.Updated {
		warnings := u.validateReadme(newContent)
		result.Warnings = warnings
	}

	// Write updated content
	if result.Updated {
		if err := u.fileUtils.WriteFile(path, newContent); err != nil {
			return nil, fmt.Errorf("failed to write README: %w", err)
		}
		u.logger.Info("README updated successfully", "path", path)
	} else {
		u.logger.Info("README is up to date, no changes needed", "path", path)
	}

	return result, nil
}

// MergeReadme merges generated content with existing README
// This is the core logic that preserves manual content
func (u *Updater) MergeReadme(existing string, generated *FormattedReadme) (string, error) {
	existingSections := u.parseSections(existing)
	result, _ := u.mergeContent(existing, existingSections, generated)
	return result, nil
}

// parseSections parses a README into sections
func (u *Updater) parseSections(content string) []*ParsedSection {
	var sections []*ParsedSection
	lines := strings.Split(content, "\n")

	var currentSection *ParsedSection
	inGeneratedBlock := false

	for i, line := range lines {
		// Check for generated marker
		if strings.Contains(line, u.config.GeneratedMarker) {
			inGeneratedBlock = true
			continue
		}

		// Check for heading (## Title format)
		if strings.HasPrefix(line, "#") {
			// Save previous section
			if currentSection != nil {
				currentSection.EndLine = i - 1
				sections = append(sections, currentSection)
			}

			// Start new section
			title := strings.TrimSpace(strings.TrimLeft(line, "#"))
			currentSection = &ParsedSection{
				Title:       title,
				StartLine:   i,
				Content:     line + "\n",
				IsManual:    !inGeneratedBlock,
				IsGenerated: inGeneratedBlock,
				Type:        u.detectSectionType(title),
			}
		} else if currentSection != nil {
			// Add line to current section
			currentSection.Content += line + "\n"
		}
	}

	// Add last section
	if currentSection != nil {
		currentSection.EndLine = len(lines) - 1
		sections = append(sections, currentSection)
	}

	return sections
}

// mergeContent merges existing and generated content
func (u *Updater) mergeContent(existing string, existingSections []*ParsedSection, generated *FormattedReadme) (string, *MergeResult) {
	result := &MergeResult{
		Added:     []string{},
		Updated:   []string{},
		Preserved: []string{},
	}

	var sb strings.Builder

	// If no existing content, just use generated
	if existing == "" {
		for _, section := range generated.Sections {
			sb.WriteString(u.formatGeneratedSection(section))
			result.Added = append(result.Added, section.Title)
		}
		return sb.String(), result
	}

	// Map existing sections by title
	existingMap := make(map[string]*ParsedSection)
	for _, s := range existingSections {
		existingMap[s.Title] = s
	}

	// Track which sections we've processed
	processed := make(map[string]bool)

	// Process in order of generated sections
	for _, genSection := range generated.Sections {
		title := strings.TrimSpace(genSection.Title)

		// Check if this section is protected
		if u.isProtected(genSection.Type) {
			if existing, ok := existingMap[title]; ok {
				// Preserve existing manual content
				sb.WriteString(existing.Content)
				result.Preserved = append(result.Preserved, title)
				processed[title] = true
			}
			continue
		}

		// Check if section exists
		if existing, ok := existingMap[title]; ok {
			if existing.IsManual {
				// Section exists and is manual - PRESERVE IT
				sb.WriteString(existing.Content)
				result.Preserved = append(result.Preserved, title)
				u.logger.Debug("Preserved manual section", "section", title)
			} else {
				// Section exists and was generated - UPDATE IT
				sb.WriteString(u.formatGeneratedSection(genSection))
				result.Updated = append(result.Updated, title)
				u.logger.Debug("Updated generated section", "section", title)
			}
			processed[title] = true
		} else {
			// New section - ADD IT
			sb.WriteString(u.formatGeneratedSection(genSection))
			result.Added = append(result.Added, title)
			u.logger.Debug("Added new section", "section", title)
			processed[title] = true
		}
	}

	// Add any remaining existing sections that weren't processed
	// (these are manual sections that don't exist in generated)
	for title, existing := range existingMap {
		if !processed[title] && existing.IsManual {
			sb.WriteString(existing.Content)
			result.Preserved = append(result.Preserved, title)
			u.logger.Debug("Preserved additional manual section", "section", title)
		}
	}

	return sb.String(), result
}

// formatGeneratedSection formats a generated section with markers
func (u *Updater) formatGeneratedSection(section *Section) string {
	var sb strings.Builder

	// Add generated marker
	sb.WriteString(u.config.GeneratedMarker)
	sb.WriteString("\n")

	// Add section title
	sb.WriteString(section.Title)
	sb.WriteString("\n\n")

	// Add section content
	sb.WriteString(section.Content)
	sb.WriteString("\n")

	return sb.String()
}

// createBackup creates a backup of the README file
func (u *Updater) createBackup(path, content string) (string, error) {
	// Create backup directory if it doesn't exist
	backupDir := filepath.Join(filepath.Dir(path), u.config.BackupDir)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("README.md.%s.backup", timestamp)
	backupPath := filepath.Join(backupDir, filename)

	// Write backup
	if err := u.fileUtils.WriteFile(backupPath, content); err != nil {
		return "", err
	}

	return backupPath, nil
}

// validateReadme performs basic validation on the README
func (u *Updater) validateReadme(content string) []string {
	var warnings []string

	// Check for broken links (simple check)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	matches := linkRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 2 {
			link := match[2]
			if strings.HasPrefix(link, "http") {
				// Can't validate external links easily
				continue
			}
			// Check if internal links point to existing files
			// This would need more context
		}
	}

	// Check for empty sections
	sectionRegex := regexp.MustCompile(`## ([^\n]+)\n\n([^#]+)`)
	matches = sectionRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 2 {
			content := strings.TrimSpace(match[2])
			if content == "" {
				warnings = append(warnings, fmt.Sprintf("Empty section: %s", match[1]))
			}
		}
	}

	// Check for placeholder text
	placeholders := []string{"TODO", "FIXME", "XXX", "your project", "your name"}
	for _, ph := range placeholders {
		if strings.Contains(content, ph) {
			warnings = append(warnings, fmt.Sprintf("Contains placeholder: %s", ph))
		}
	}

	return warnings
}

// detectSectionType attempts to detect the section type from title
func (u *Updater) detectSectionType(title string) types.DocSection {
	lower := strings.ToLower(title)

	switch {
	case strings.Contains(lower, "install"):
		return types.DocSectionInstallation
	case strings.Contains(lower, "quick") && strings.Contains(lower, "start"):
		return types.DocSectionQuickStart
	case strings.Contains(lower, "api") || strings.Contains(lower, "reference"):
		return types.DocSectionAPI
	case strings.Contains(lower, "example"):
		return types.DocSectionExamples
	case strings.Contains(lower, "contribut"):
		return types.DocSectionContributing
	case strings.Contains(lower, "license"):
		return types.DocSectionLicense
	case strings.Contains(lower, "description") || lower == "":
		return types.DocSectionDescription
	default:
		return types.DocSectionCustom
	}
}

// isProtected checks if a section type is protected from auto-update
func (u *Updater) isProtected(sectionType types.DocSection) bool {
	for _, p := range u.config.ProtectedSections {
		if p == sectionType {
			return true
		}
	}
	return false
}

// DryRun performs an update without writing to disk
func (u *Updater) DryRun(ctx context.Context, path string, generated *FormattedReadme) (*UpdateResult, error) {
	u.logger.Info("Performing dry run update", "path", path)

	result := &UpdateResult{
		Path:              path,
		AddedSections:     []string{},
		UpdatedSections:   []string{},
		PreservedSections: []string{},
	}

	// Read existing if present
	var existingContent string
	var err error

	if u.fileUtils.FileExists(path) {
		existingContent, err = u.fileUtils.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}

	// Parse existing sections
	var existingSections []*ParsedSection
	if existingContent != "" {
		existingSections = u.parseSections(existingContent)
	}

	// Simulate merge
	_, mergeResult := u.mergeContent(existingContent, existingSections, generated)

	result.AddedSections = mergeResult.Added
	result.UpdatedSections = mergeResult.Updated
	result.PreservedSections = mergeResult.Preserved
	result.Updated = len(mergeResult.Added) > 0 || len(mergeResult.Updated) > 0

	return result, nil
}

// RestoreFromBackup restores a README from the most recent backup
func (u *Updater) RestoreFromBackup(path string) error {
	backupDir := filepath.Join(filepath.Dir(path), u.config.BackupDir)

	// Find most recent backup
	files, err := filepath.Glob(filepath.Join(backupDir, "README.md.*.backup"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no backups found")
	}

	// Sort by name (which includes timestamp) - most recent last
	// Simple approach: take the last one alphabetically
	mostRecent := files[len(files)-1]

	// Read backup
	content, err := u.fileUtils.ReadFile(mostRecent)
	if err != nil {
		return err
	}

	// Write to original location
	if err := u.fileUtils.WriteFile(path, content); err != nil {
		return err
	}

	u.logger.Info("Restored from backup", "backup", mostRecent, "target", path)
	return nil
}

// MergeResult represents the result of a merge operation
type MergeResult struct {
	Added     []string
	Updated   []string
	Preserved []string
}

// GetStats returns statistics about the merge
func (u *Updater) GetStats(result *UpdateResult) map[string]interface{} {
	return map[string]interface{}{
		"updated":            result.Updated,
		"sections_added":     len(result.AddedSections),
		"sections_updated":   len(result.UpdatedSections),
		"sections_preserved": len(result.PreservedSections),
		"warnings":           len(result.Warnings),
		"backup_created":     result.BackupPath != "",
	}
}

// Diff returns the differences between existing and new content
func (u *Updater) Diff(existing string, generated *FormattedReadme) (string, error) {
	existingSections := u.parseSections(existing)
	newContent, _ := u.mergeContent(existing, existingSections, generated)

	// Simple diff - in production you might want to use a proper diff library
	if existing == newContent {
		return "No changes", nil
	}

	var sb strings.Builder
	sb.WriteString("Changes detected:\n\n")

	// This is a very simplified diff
	existingLines := strings.Split(existing, "\n")
	newLines := strings.Split(newContent, "\n")

	if len(existingLines) != len(newLines) {
		sb.WriteString(fmt.Sprintf("- Line count changed: %d → %d\n",
			len(existingLines), len(newLines)))
	}

	return sb.String(), nil
}
