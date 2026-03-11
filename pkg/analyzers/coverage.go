package analyzers

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"example.com/spoke-tool/api/types"
)

// CoverageAnalyzer provides test coverage analysis
// This uses standard coverage tools for each language
type CoverageAnalyzer struct {
	workDir string
}

// CoverageConfig configures the coverage analyzer
type CoverageConfig struct {
	// Working directory
	WorkDir string

	// Coverage output format
	Format string // "text", "html", "xml", "json"

	// Minimum coverage threshold
	Threshold float64

	// Whether to include detailed file coverage
	Detailed bool
}

// CoverageResult represents coverage analysis results
type CoverageResult struct {
	Language   types.Language     `json:"language"`
	Overall    float64            `json:"overall_percent"`
	ByFile     map[string]float64 `json:"by_file,omitempty"`
	ByFunction map[string]float64 `json:"by_function,omitempty"`
	Uncovered  []UncoveredLine    `json:"uncovered,omitempty"`
	Timestamp  string             `json:"timestamp"`
	Command    string             `json:"command"`
	Output     string             `json:"output,omitempty"`
}

// UncoveredLine represents an uncovered line of code
type UncoveredLine struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Content  string `json:"content,omitempty"`
	Function string `json:"function,omitempty"`
}

// GoCoverage represents Go coverage profile
type GoCoverage struct {
	XMLName  xml.Name    `xml:"coverage"`
	Packages []GoPackage `xml:"packages>package"`
	Overall  float64
}

// GoPackage represents a Go package in coverage report
type GoPackage struct {
	Name      string       `xml:"name,attr"`
	Functions []GoFunction `xml:"functions>function"`
	Coverage  float64      `xml:"coverage,attr"`
}

// GoFunction represents a Go function in coverage report
type GoFunction struct {
	Name     string  `xml:"name,attr"`
	Coverage float64 `xml:"coverage,attr"`
}

// NewCoverageAnalyzer creates a new coverage analyzer
func NewCoverageAnalyzer(config CoverageConfig) *CoverageAnalyzer {
	if config.WorkDir == "" {
		config.WorkDir = "."
	}
	if config.Format == "" {
		config.Format = "text"
	}

	return &CoverageAnalyzer{
		workDir: config.WorkDir,
	}
}

// AnalyzeGoCoverage runs go test with coverage and parses results
func (c *CoverageAnalyzer) AnalyzeGoCoverage() (*CoverageResult, error) {
	result := &CoverageResult{
		Language:   types.Go,
		ByFile:     make(map[string]float64),
		ByFunction: make(map[string]float64),
		Uncovered:  []UncoveredLine{},
		Command:    "go test -cover",
	}

	// Create temporary coverage file
	tmpFile, err := os.CreateTemp("", "coverage-*.out")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Run go test with coverage
	cmd := exec.Command("go", "test", "./...", "-coverprofile="+tmpFile.Name())
	cmd.Dir = c.workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Even if tests fail, we might still have coverage
		result.Output = string(output)
	}

	// Parse coverage file
	if _, err := os.Stat(tmpFile.Name()); err == nil {
		if err := c.parseGoCoverageFile(tmpFile.Name(), result); err != nil {
			return nil, err
		}
	}

	// Try to get overall coverage from output
	c.parseGoOverallCoverage(string(output), result)

	return result, nil
}

// AnalyzeNodeJSCoverage runs Jest with coverage and parses results
func (c *CoverageAnalyzer) AnalyzeNodeJSCoverage() (*CoverageResult, error) {
	result := &CoverageResult{
		Language:   types.NodeJS,
		ByFile:     make(map[string]float64),
		ByFunction: make(map[string]float64),
		Uncovered:  []UncoveredLine{},
		Command:    "jest --coverage",
	}

	// Run Jest with coverage
	cmd := exec.Command("npx", "jest", "--coverage", "--json", "--outputFile=coverage.json")
	cmd.Dir = c.workDir
	output, _ := cmd.CombinedOutput()
	result.Output = string(output)

	// Parse Jest coverage output
	coverageFile := filepath.Join(c.workDir, "coverage.json")
	if _, err := os.Stat(coverageFile); err == nil {
		if err := c.parseJestCoverage(coverageFile, result); err != nil {
			return nil, err
		}
		defer os.Remove(coverageFile)
	}

	return result, nil
}

// AnalyzePythonCoverage runs pytest with coverage and parses results
func (c *CoverageAnalyzer) AnalyzePythonCoverage() (*CoverageResult, error) {
	result := &CoverageResult{
		Language:   types.Python,
		ByFile:     make(map[string]float64),
		ByFunction: make(map[string]float64),
		Uncovered:  []UncoveredLine{},
		Command:    "pytest --cov=. --cov-report=xml",
	}

	// Run pytest with coverage
	cmd := exec.Command("pytest", "--cov=.", "--cov-report=xml", "--cov-report=term")
	cmd.Dir = c.workDir
	output, _ := cmd.CombinedOutput()
	result.Output = string(output)

	// Parse coverage XML
	coverageFile := filepath.Join(c.workDir, "coverage.xml")
	if _, err := os.Stat(coverageFile); err == nil {
		if err := c.parsePythonCoverageXML(coverageFile, result); err != nil {
			return nil, err
		}
		defer os.Remove(coverageFile)
	}

	// Also try to get overall from output
	c.parsePythonOverallCoverage(string(output), result)

	return result, nil
}

// AnalyzeCoverage automatically detects language and runs appropriate analyzer
func (c *CoverageAnalyzer) AnalyzeCoverage() (*CoverageResult, error) {
	// Detect language by checking for common files
	if c.hasGoFiles() {
		return c.AnalyzeGoCoverage()
	} else if c.hasNodeFiles() {
		return c.AnalyzeNodeJSCoverage()
	} else if c.hasPythonFiles() {
		return c.AnalyzePythonCoverage()
	}

	return nil, fmt.Errorf("unable to detect project language")
}

// Parse coverage file methods

func (c *CoverageAnalyzer) parseGoCoverageFile(path string, result *CoverageResult) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Skip first line (mode)
	if scanner.Scan() {
		// mode line
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		filePath := parts[0]
		coverageData := parts[1]

		// Parse coverage data (simplified)
		// Format: start.end,count count...
		if strings.Contains(coverageData, ",") {
			c.parseGoCoverageLine(filePath, coverageData, result)
		}
	}

	return scanner.Err()
}

func (c *CoverageAnalyzer) parseGoCoverageLine(filePath string, data string, result *CoverageResult) {
	parts := strings.Split(data, ",")
	if len(parts) < 2 {
		return
	}

	// Get line numbers
	lineRange := strings.Split(parts[0], ".")
	if len(lineRange) < 2 {
		return
	}

	startLine, _ := strconv.Atoi(lineRange[0])
	endLine, _ := strconv.Atoi(lineRange[1])

	// Get count
	count, _ := strconv.Atoi(parts[1])

	// If count is 0, line is uncovered
	if count == 0 {
		for line := startLine; line <= endLine; line++ {
			result.Uncovered = append(result.Uncovered, UncoveredLine{
				File: filePath,
				Line: line,
			})
		}
	}
}

func (c *CoverageAnalyzer) parseGoOverallCoverage(output string, result *CoverageResult) {
	// Look for "coverage: X.X% of statements"
	re := regexp.MustCompile(`coverage:\s*(\d+\.\d+)%`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		overall, _ := strconv.ParseFloat(matches[1], 64)
		result.Overall = overall
	}
}

func (c *CoverageAnalyzer) parseJestCoverage(path string, result *CoverageResult) error {
	// This is a simplified version - Jest coverage JSON is complex
	// In production, you'd want to properly parse the JSON

	// For now, we'll look for summary file
	summaryFile := filepath.Join(filepath.Dir(path), "coverage-summary.json")
	if _, err := os.Stat(summaryFile); err == nil {
		// Parse summary
		defer os.Remove(summaryFile)
	}

	// Try to parse overall from output
	c.parseJestOverallCoverage(result.Output, result)

	return nil
}

func (c *CoverageAnalyzer) parseJestOverallCoverage(output string, result *CoverageResult) {
	// Look for "All files" line with coverage
	re := regexp.MustCompile(`All files\s*\|\s*(\d+\.\d+)`)
	matches := re.FindAllStringSubmatch(output, -1)

	for _, match := range matches {
		if len(match) > 1 {
			overall, _ := strconv.ParseFloat(match[1], 64)
			result.Overall = overall
			break
		}
	}
}

func (c *CoverageAnalyzer) parsePythonCoverageXML(path string, result *CoverageResult) error {
	// This is a simplified version
	// In production, you'd want to properly parse the XML

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Look for line-rate attribute
	re := regexp.MustCompile(`line-rate="(\d+\.?\d*)"`)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			overall, _ := strconv.ParseFloat(matches[1], 64)
			// Convert to percentage
			result.Overall = overall * 100
			break
		}
	}

	return scanner.Err()
}

func (c *CoverageAnalyzer) parsePythonOverallCoverage(output string, result *CoverageResult) {
	// Look for "TOTAL" line with coverage
	re := regexp.MustCompile(`TOTAL\s+\d+\s+\d+\s+(\d+)%`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		overall, _ := strconv.ParseFloat(matches[1], 64)
		result.Overall = overall
	}
}

// Language detection helpers

func (c *CoverageAnalyzer) hasGoFiles() bool {
	matches, _ := filepath.Glob(filepath.Join(c.workDir, "*.go"))
	return len(matches) > 0
}

func (c *CoverageAnalyzer) hasNodeFiles() bool {
	// Check for package.json or JS files
	if _, err := os.Stat(filepath.Join(c.workDir, "package.json")); err == nil {
		return true
	}

	matches, _ := filepath.Glob(filepath.Join(c.workDir, "*.js"))
	matches2, _ := filepath.Glob(filepath.Join(c.workDir, "*.ts"))
	return len(matches) > 0 || len(matches2) > 0
}

func (c *CoverageAnalyzer) hasPythonFiles() bool {
	// Check for setup.py, requirements.txt, or py files
	if _, err := os.Stat(filepath.Join(c.workDir, "setup.py")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(c.workDir, "requirements.txt")); err == nil {
		return true
	}

	matches, _ := filepath.Glob(filepath.Join(c.workDir, "*.py"))
	return len(matches) > 0
}

// CheckCoverageThreshold checks if coverage meets the threshold
func (c *CoverageAnalyzer) CheckCoverageThreshold(result *CoverageResult, threshold float64) (bool, error) {
	if result == nil {
		return false, fmt.Errorf("no coverage result")
	}

	if result.Overall >= threshold {
		return true, nil
	}

	return false, nil
}

// GetUncoveredFiles returns files with coverage below threshold
func (c *CoverageAnalyzer) GetUncoveredFiles(result *CoverageResult, threshold float64) []string {
	var uncovered []string

	for file, cov := range result.ByFile {
		if cov < threshold {
			uncovered = append(uncovered, file)
		}
	}

	return uncovered
}

// GetUncoveredFunctions returns functions with coverage below threshold
func (c *CoverageAnalyzer) GetUncoveredFunctions(result *CoverageResult, threshold float64) []string {
	var uncovered []string

	for fn, cov := range result.ByFunction {
		if cov < threshold {
			uncovered = append(uncovered, fn)
		}
	}

	return uncovered
}

// FormatCoverageReport formats coverage results for display
func (c *CoverageAnalyzer) FormatCoverageReport(result *CoverageResult) string {
	if result == nil {
		return "No coverage data available"
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("📊 Test Coverage Report (%s)\n", result.Language))
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("Overall: %.1f%%\n", result.Overall))

	// Color code based on coverage
	if result.Overall >= 80 {
		sb.WriteString("✅ Excellent coverage\n")
	} else if result.Overall >= 60 {
		sb.WriteString("⚠️  Good coverage, but could improve\n")
	} else {
		sb.WriteString("❌ Low coverage - needs attention\n")
	}

	sb.WriteString("\n")

	// File coverage
	if len(result.ByFile) > 0 {
		sb.WriteString("📁 Coverage by File:\n")
		sb.WriteString("--------------------\n")

		for file, cov := range result.ByFile {
			marker := "✅"
			if cov < 80 {
				marker = "⚠️"
			}
			if cov < 60 {
				marker = "❌"
			}
			sb.WriteString(fmt.Sprintf("%s %s: %.1f%%\n", marker, filepath.Base(file), cov))
		}
		sb.WriteString("\n")
	}

	// Uncovered lines
	if len(result.Uncovered) > 0 {
		sb.WriteString("🔍 Uncovered Lines:\n")
		sb.WriteString("------------------\n")

		// Group by file
		byFile := make(map[string][]UncoveredLine)
		for _, u := range result.Uncovered {
			byFile[u.File] = append(byFile[u.File], u)
		}

		for file, lines := range byFile {
			if len(lines) > 5 {
				sb.WriteString(fmt.Sprintf("  %s: %d lines uncovered\n", filepath.Base(file), len(lines)))
			} else {
				sb.WriteString(fmt.Sprintf("  %s: lines ", filepath.Base(file)))
				for i, line := range lines {
					if i > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(strconv.Itoa(line.Line))
				}
				sb.WriteString("\n")
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Command: %s\n", result.Command))

	return sb.String()
}

// CompareCoverage compares two coverage results
func (c *CoverageAnalyzer) CompareCoverage(old, new *CoverageResult) map[string]float64 {
	deltas := make(map[string]float64)

	if old == nil || new == nil {
		return deltas
	}

	// Overall delta
	deltas["overall"] = new.Overall - old.Overall

	// File deltas
	for file, newCov := range new.ByFile {
		if oldCov, ok := old.ByFile[file]; ok {
			deltas[file] = newCov - oldCov
		}
	}

	return deltas
}

// FormatComparison formats coverage comparison
func (c *CoverageAnalyzer) FormatComparison(deltas map[string]float64) string {
	if len(deltas) == 0 {
		return "No comparison data"
	}

	var sb strings.Builder

	sb.WriteString("📈 Coverage Changes\n")
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	if overall, ok := deltas["overall"]; ok {
		arrow := "➡️"
		if overall > 0 {
			arrow = "⬆️"
		} else if overall < 0 {
			arrow = "⬇️"
		}

		sb.WriteString(fmt.Sprintf("Overall: %s %.1f%%\n", arrow, overall))
		sb.WriteString("\n")
	}

	for file, delta := range deltas {
		if file == "overall" {
			continue
		}

		if delta != 0 {
			arrow := "➡️"
			if delta > 0 {
				arrow = "⬆️"
			} else if delta < 0 {
				arrow = "⬇️"
			}
			sb.WriteString(fmt.Sprintf("%s %s: %.1f%%\n", arrow, filepath.Base(file), delta))
		}
	}

	return sb.String()
}
