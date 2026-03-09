package runners

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yourusername/spoke-tool/api/types"
)

// CoverageRunner runs coverage tools and collects results
// This runs STANDARD coverage tools - no custom implementations
type CoverageRunner struct {
	workDir string
	timeout time.Duration
	verbose bool
}

// CoverageConfig configures the coverage runner
type CoverageConfig struct {
	// Working directory
	WorkDir string

	// Timeout for coverage run
	Timeout time.Duration

	// Coverage threshold (for reporting)
	Threshold float64

	// Output format (text, html, json, xml)
	Format string

	// Whether to include detailed file coverage
	Detailed bool

	// Whether to show uncovered lines
	ShowUncovered bool

	// Verbose output
	Verbose bool
}

// CoverageResult represents the result of a coverage run
type CoverageResult struct {
	// Language of the project
	Language types.Language `json:"language"`

	// Overall coverage percentage
	Overall float64 `json:"overall_percent"`

	// Coverage by file
	ByFile map[string]float64 `json:"by_file,omitempty"`

	// Coverage by function
	ByFunction map[string]float64 `json:"by_function,omitempty"`

	// Uncovered lines
	Uncovered []UncoveredLine `json:"uncovered,omitempty"`

	// Command that was run
	Command string `json:"command"`

	// Raw output
	Output string `json:"output,omitempty"`

	// Error if any
	Error string `json:"error,omitempty"`

	// Timestamp
	Timestamp time.Time `json:"timestamp"`

	// Duration
	Duration time.Duration `json:"duration_ms"`
}

// UncoveredLine represents an uncovered line of code
type UncoveredLine struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Content  string `json:"content,omitempty"`
	Function string `json:"function,omitempty"`
}

// GoCoverageProfile represents a Go coverage profile
type GoCoverageProfile struct {
	Mode   string
	Blocks []GoCoverageBlock
}

// GoCoverageBlock represents a block in a Go coverage profile
type GoCoverageBlock struct {
	FileName  string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	NumStmt   int
	Count     int
}

// CoberturaCoverage represents Cobertura XML format (used by Python)
type CoberturaCoverage struct {
	XMLName    xml.Name `xml:"coverage"`
	LineRate   float64  `xml:"line-rate,attr"`
	BranchRate float64  `xml:"branch-rate,attr"`
	Version    string   `xml:"version,attr"`
	Timestamp  int64    `xml:"timestamp,attr"`
	Packages   []struct {
		Name     string  `xml:"name,attr"`
		LineRate float64 `xml:"line-rate,attr"`
		Classes  []struct {
			Name     string  `xml:"name,attr"`
			Filename string  `xml:"filename,attr"`
			LineRate float64 `xml:"line-rate,attr"`
			Lines    []struct {
				Number int `xml:"number,attr"`
				Hits   int `xml:"hits,attr"`
			} `xml:"lines>line"`
		} `xml:"classes>class"`
	} `xml:"packages>package"`
}

// JestCoverageSummary represents Jest coverage summary
type JestCoverageSummary struct {
	Total struct {
		Lines struct {
			Pct float64 `json:"pct"`
		} `json:"lines"`
	} `json:"total"`
}

// NewCoverageRunner creates a new coverage runner
func NewCoverageRunner(config CoverageConfig) *CoverageRunner {
	if config.WorkDir == "" {
		config.WorkDir = "."
	}
	if config.Timeout == 0 {
		config.Timeout = 2 * time.Minute
	}
	if config.Format == "" {
		config.Format = "text"
	}

	return &CoverageRunner{
		workDir: config.WorkDir,
		timeout: config.Timeout,
		verbose: config.Verbose,
	}
}

// RunGoCoverage runs go test with coverage
func (r *CoverageRunner) RunGoCoverage() (*CoverageResult, error) {
	result := &CoverageResult{
		Language:   types.Go,
		ByFile:     make(map[string]float64),
		ByFunction: make(map[string]float64),
		Uncovered:  []UncoveredLine{},
		Timestamp:  time.Now(),
		Command:    "go test -coverprofile=coverage.out",
	}

	start := time.Now()

	// Create temporary coverage file
	tmpFile, err := os.CreateTemp(r.workDir, "coverage-*.out")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpName)

	// Run go test with coverage
	cmd := exec.Command("go", "test", "./...", "-coverprofile="+tmpName)
	cmd.Dir = r.workDir
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(start)

	// Parse coverage profile if it exists
	if _, err := os.Stat(tmpName); err == nil {
		if err := r.parseGoCoverageProfile(tmpName, result); err != nil {
			if r.verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to parse coverage: %v\n", err)
			}
		}
	}

	// Try to get overall from output
	r.parseGoOverallCoverage(result.Output, result)

	if err != nil {
		result.Error = err.Error()
	}

	return result, nil
}

// RunNodeJSCoverage runs Jest with coverage
func (r *CoverageRunner) RunNodeJSCoverage() (*CoverageResult, error) {
	result := &CoverageResult{
		Language:   types.NodeJS,
		ByFile:     make(map[string]float64),
		ByFunction: make(map[string]float64),
		Timestamp:  time.Now(),
		Command:    "jest --coverage",
	}

	start := time.Now()

	// Check if jest is installed
	if _, err := exec.LookPath("npx"); err != nil {
		return nil, fmt.Errorf("npx not found: %w", err)
	}

	// Run jest with coverage
	cmd := exec.Command("npx", "jest", "--coverage", "--json", "--outputFile=coverage-results.json")
	cmd.Dir = r.workDir
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(start)

	// Parse coverage summary
	summaryFile := filepath.Join(r.workDir, "coverage-results.json")
	if _, err := os.Stat(summaryFile); err == nil {
		if err := r.parseJestCoverage(summaryFile, result); err != nil {
			if r.verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to parse Jest coverage: %v\n", err)
			}
		}
		defer os.Remove(summaryFile)
	}

	// Also look for coverage summary in standard location
	lcovFile := filepath.Join(r.workDir, "coverage", "lcov.info")
	if _, err := os.Stat(lcovFile); err == nil {
		if err := r.parseLcovFile(lcovFile, result); err != nil && r.verbose {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse lcov: %v\n", err)
		}
	}

	// Try to get overall from output
	r.parseJestOverallCoverage(result.Output, result)

	if err != nil {
		result.Error = err.Error()
	}

	return result, nil
}

// RunPythonCoverage runs pytest with coverage
func (r *CoverageRunner) RunPythonCoverage() (*CoverageResult, error) {
	result := &CoverageResult{
		Language:   types.Python,
		ByFile:     make(map[string]float64),
		ByFunction: make(map[string]float64),
		Uncovered:  []UncoveredLine{},
		Timestamp:  time.Now(),
		Command:    "pytest --cov=. --cov-report=xml",
	}

	start := time.Now()

	// Check if pytest is installed
	if _, err := exec.LookPath("pytest"); err != nil {
		return nil, fmt.Errorf("pytest not found: %w", err)
	}

	// Run pytest with coverage
	cmd := exec.Command("pytest", "--cov=.", "--cov-report=xml", "--cov-report=term")
	cmd.Dir = r.workDir
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(start)

	// Parse coverage XML
	xmlFile := filepath.Join(r.workDir, "coverage.xml")
	if _, err := os.Stat(xmlFile); err == nil {
		if err := r.parsePythonCoverageXML(xmlFile, result); err != nil {
			if r.verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to parse coverage XML: %v\n", err)
			}
		}
		defer os.Remove(xmlFile)
	}

	// Try to get overall from output
	r.parsePythonOverallCoverage(result.Output, result)

	if err != nil {
		result.Error = err.Error()
	}

	return result, nil
}

// RunCoverage automatically detects language and runs appropriate coverage tool
func (r *CoverageRunner) RunCoverage() (*CoverageResult, error) {
	// Detect language
	if r.hasGoFiles() {
		return r.RunGoCoverage()
	} else if r.hasNodeFiles() {
		return r.RunNodeJSCoverage()
	} else if r.hasPythonFiles() {
		return r.RunPythonCoverage()
	}

	return nil, fmt.Errorf("unable to detect project language")
}

// RunCoverageWithThreshold runs coverage and checks against threshold
func (r *CoverageRunner) RunCoverageWithThreshold(threshold float64) (*CoverageResult, bool, error) {
	result, err := r.RunCoverage()
	if err != nil {
		return result, false, err
	}

	passed := result.Overall >= threshold
	return result, passed, nil
}

// Parsing methods

func (r *CoverageRunner) parseGoCoverageProfile(path string, result *CoverageResult) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Parse mode line
	if scanner.Scan() {
		// mode line, ignore
	}

	fileStats := make(map[string]struct {
		covered   int
		total     int
		uncovered []int
	})

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		fileName := parts[0]
		counts := parts[1]

		// Parse "start.end,count count..."
		countParts := strings.Split(counts, ",")
		if len(countParts) < 2 {
			continue
		}

		// Get line range
		lineRange := strings.Split(countParts[0], ".")
		if len(lineRange) < 2 {
			continue
		}

		startLine, _ := strconv.Atoi(lineRange[0])
		endLine, _ := strconv.Atoi(lineRange[1])

		// Get count
		count, _ := strconv.Atoi(countParts[1])

		// Update file stats
		stats := fileStats[fileName]
		stats.total += endLine - startLine + 1
		if count > 0 {
			stats.covered += endLine - startLine + 1
		} else {
			for line := startLine; line <= endLine; line++ {
				stats.uncovered = append(stats.uncovered, line)
			}
		}
		fileStats[fileName] = stats
	}

	// Calculate percentages
	var totalCovered, totalLines int
	for file, stats := range fileStats {
		if stats.total > 0 {
			percentage := float64(stats.covered) / float64(stats.total) * 100
			result.ByFile[file] = percentage

			totalCovered += stats.covered
			totalLines += stats.total

			// Add uncovered lines
			for _, line := range stats.uncovered {
				result.Uncovered = append(result.Uncovered, UncoveredLine{
					File: file,
					Line: line,
				})
			}
		}
	}

	if totalLines > 0 {
		result.Overall = float64(totalCovered) / float64(totalLines) * 100
	}

	return scanner.Err()
}

func (r *CoverageRunner) parseGoOverallCoverage(output string, result *CoverageResult) {
	// Look for "coverage: X.X% of statements"
	re := regexp.MustCompile(`coverage:\s*(\d+\.\d+)%`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		overall, _ := strconv.ParseFloat(matches[1], 64)
		result.Overall = overall
	}
}

func (r *CoverageRunner) parseJestCoverage(path string, result *CoverageResult) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var summary JestCoverageSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		return err
	}

	result.Overall = summary.Total.Lines.Pct
	return nil
}

func (r *CoverageRunner) parseJestOverallCoverage(output string, result *CoverageResult) {
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

func (r *CoverageRunner) parseLcovFile(path string, result *CoverageResult) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var currentFile string
	var linesFound, linesHit int

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "SF:") {
			// Start of a new file
			if currentFile != "" && linesFound > 0 {
				percentage := float64(linesHit) / float64(linesFound) * 100
				result.ByFile[currentFile] = percentage
			}
			currentFile = strings.TrimPrefix(line, "SF:")
			linesFound = 0
			linesHit = 0
		} else if strings.HasPrefix(line, "LF:") {
			// Lines found
			lf := strings.TrimPrefix(line, "LF:")
			linesFound, _ = strconv.Atoi(lf)
		} else if strings.HasPrefix(line, "LH:") {
			// Lines hit
			lh := strings.TrimPrefix(line, "LH:")
			linesHit, _ = strconv.Atoi(lh)
		}
	}

	// Last file
	if currentFile != "" && linesFound > 0 {
		percentage := float64(linesHit) / float64(linesFound) * 100
		result.ByFile[currentFile] = percentage
	}

	return scanner.Err()
}

func (r *CoverageRunner) parsePythonCoverageXML(path string, result *CoverageResult) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var coverage CoberturaCoverage
	if err := xml.Unmarshal(data, &coverage); err != nil {
		return err
	}

	// Overall coverage
	result.Overall = coverage.LineRate * 100

	// File coverage
	for _, pkg := range coverage.Packages {
		for _, class := range pkg.Classes {
			result.ByFile[class.Filename] = class.LineRate * 100

			// Parse uncovered lines
			if len(class.Lines) > 0 {
				for _, line := range class.Lines {
					if line.Hits == 0 {
						result.Uncovered = append(result.Uncovered, UncoveredLine{
							File: class.Filename,
							Line: line.Number,
						})
					}
				}
			}
		}
	}

	return nil
}

func (r *CoverageRunner) parsePythonOverallCoverage(output string, result *CoverageResult) {
	// Look for "TOTAL" line with coverage
	re := regexp.MustCompile(`TOTAL\s+\d+\s+\d+\s+(\d+)%`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		overall, _ := strconv.ParseFloat(matches[1], 64)
		result.Overall = overall
	}
}

// Language detection

func (r *CoverageRunner) hasGoFiles() bool {
	matches, _ := filepath.Glob(filepath.Join(r.workDir, "*.go"))
	return len(matches) > 0
}

func (r *CoverageRunner) hasNodeFiles() bool {
	// Check for package.json
	if _, err := os.Stat(filepath.Join(r.workDir, "package.json")); err == nil {
		return true
	}

	// Check for JS/TS files
	jsMatches, _ := filepath.Glob(filepath.Join(r.workDir, "*.js"))
	tsMatches, _ := filepath.Glob(filepath.Join(r.workDir, "*.ts"))
	return len(jsMatches) > 0 || len(tsMatches) > 0
}

func (r *CoverageRunner) hasPythonFiles() bool {
	// Check for setup.py or requirements.txt
	if _, err := os.Stat(filepath.Join(r.workDir, "setup.py")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(r.workDir, "requirements.txt")); err == nil {
		return true
	}

	// Check for Python files
	matches, _ := filepath.Glob(filepath.Join(r.workDir, "*.py"))
	return len(matches) > 0
}

// Reporting methods

// FormatResult formats the coverage result for display
func (r *CoverageRunner) FormatResult(result *CoverageResult) string {
	if result == nil {
		return "No coverage data available"
	}

	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("📊 Test Coverage Report (%s)\n", result.Language))
	sb.WriteString(strings.Repeat("=", 60))
	sb.WriteString("\n\n")

	// Overall
	sb.WriteString(fmt.Sprintf("Overall Coverage: %.1f%%", result.Overall))

	// Color indicator
	if result.Overall >= 80 {
		sb.WriteString(" ✅\n")
	} else if result.Overall >= 60 {
		sb.WriteString(" ⚠️\n")
	} else {
		sb.WriteString(" ❌\n")
	}

	sb.WriteString("\n")

	// File coverage
	if len(result.ByFile) > 0 {
		sb.WriteString("Coverage by File:\n")
		sb.WriteString("-----------------\n")

		// Sort files by coverage (lowest first)
		type fileCov struct {
			name string
			cov  float64
		}
		var files []fileCov
		for name, cov := range result.ByFile {
			files = append(files, fileCov{name, cov})
		}

		// Simple bubble sort
		for i := 0; i < len(files)-1; i++ {
			for j := i + 1; j < len(files); j++ {
				if files[i].cov > files[j].cov {
					files[i], files[j] = files[j], files[i]
				}
			}
		}

		for _, f := range files {
			marker := "✅"
			if f.cov < 80 {
				marker = "⚠️"
			}
			if f.cov < 60 {
				marker = "❌"
			}
			sb.WriteString(fmt.Sprintf("%s %s: %.1f%%\n", marker, filepath.Base(f.name), f.cov))
		}
		sb.WriteString("\n")
	}

	// Uncovered lines
	if len(result.Uncovered) > 0 && len(result.Uncovered) <= 20 {
		sb.WriteString("Uncovered Lines:\n")
		sb.WriteString("----------------\n")

		// Group by file
		byFile := make(map[string][]int)
		for _, u := range result.Uncovered {
			byFile[u.File] = append(byFile[u.File], u.Line)
		}

		for file, lines := range byFile {
			sb.WriteString(fmt.Sprintf("  %s: ", filepath.Base(file)))
			for i, line := range lines {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(strconv.Itoa(line))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	} else if len(result.Uncovered) > 20 {
		sb.WriteString(fmt.Sprintf("📋 %d lines uncovered across %d files\n\n",
			len(result.Uncovered), len(result.ByFile)))
	}

	// Command and timing
	sb.WriteString(fmt.Sprintf("Command: %s\n", result.Command))
	sb.WriteString(fmt.Sprintf("Duration: %v\n", result.Duration))

	if result.Error != "" {
		sb.WriteString(fmt.Sprintf("\nError: %s\n", result.Error))
	}

	return sb.String()
}

// CheckThreshold checks if coverage meets threshold
func (r *CoverageRunner) CheckThreshold(result *CoverageResult, threshold float64) (bool, map[string]float64) {
	if result == nil {
		return false, nil
	}

	failing := make(map[string]float64)

	// Check overall
	if result.Overall < threshold {
		failing["overall"] = result.Overall
	}

	// Check files
	for file, cov := range result.ByFile {
		if cov < threshold {
			failing[file] = cov
		}
	}

	return len(failing) == 0, failing
}

// GenerateHTMLReport generates an HTML coverage report
func (r *CoverageRunner) GenerateHTMLReport(result *CoverageResult, outputPath string) error {
	// This would generate a nice HTML report
	// For now, just write a simple HTML file
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Coverage Report</title>
    <style>
        body { font-family: sans-serif; margin: 40px; }
        .header { background: #f0f0f0; padding: 20px; border-radius: 5px; }
        .overall { font-size: 24px; margin: 20px 0; }
        .good { color: green; }
        .warning { color: orange; }
        .bad { color: red; }
        table { border-collapse: collapse; width: 100%%; }
        th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Coverage Report</h1>
        <p>Generated: %s</p>
        <p>Language: %s</p>
    </div>
    
    <div class="overall %s">
        Overall Coverage: %.1f%%
    </div>
    
    <h2>Coverage by File</h2>
    <table>
        <tr>
            <th>File</th>
            <th>Coverage</th>
        </tr>
        %s
    </table>
</body>
</html>`,
		time.Now().Format(time.RFC3339),
		result.Language,
		getCoverageClass(result.Overall),
		result.Overall,
		generateFileRows(result),
	)

	return os.WriteFile(outputPath, []byte(html), 0644)
}

func getCoverageClass(cov float64) string {
	if cov >= 80 {
		return "good"
	}
	if cov >= 60 {
		return "warning"
	}
	return "bad"
}

func generateFileRows(result *CoverageResult) string {
	var rows strings.Builder
	for file, cov := range result.ByFile {
		class := getCoverageClass(cov)
		rows.WriteString(fmt.Sprintf(
			"<tr><td>%s</td><td class=\"%s\">%.1f%%</td></tr>\n",
			filepath.Base(file), class, cov))
	}
	return rows.String()
}

// GetSummary returns a brief summary
func (r *CoverageRunner) GetSummary(result *CoverageResult) string {
	if result == nil {
		return "No coverage data"
	}

	return fmt.Sprintf("%s coverage: %.1f%% (%d files)",
		result.Language,
		result.Overall,
		len(result.ByFile))
}
