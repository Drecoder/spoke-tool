package runners

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"example.com/spoke-tool/api/types"
)

// TestRunner runs tests and collects results
// This runs STANDARD testing frameworks - no custom implementations
type TestRunner struct {
	workDir  string
	timeout  time.Duration
	parallel bool
	verbose  bool
}

// TestConfig configures the test runner
type TestConfig struct {
	// Working directory
	WorkDir string

	// Timeout for test execution
	Timeout time.Duration

	// Whether to run tests in parallel
	Parallel bool

	// Maximum number of parallel test runs
	MaxParallel int

	// Whether to show test output
	Verbose bool

	// Environment variables
	Env map[string]string

	// Test patterns (e.g., "./...", "test_*.py")
	Pattern string
}

// TestResult represents the result of a test run
type TestResult struct {
	// Language of the tests
	Language types.Language `json:"language"`

	// Test suite results
	Suite *types.TestSuite `json:"suite"`

	// Command that was run
	Command string `json:"command"`

	// Raw output
	Output string `json:"output,omitempty"`

	// Error if any
	Error string `json:"error,omitempty"`

	// Timing
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration_ms"`

	// Test files that were run
	TestFiles []string `json:"test_files,omitempty"`
}

// TestSummary represents a summary of test results
type TestSummary struct {
	Total      int                                     `json:"total"`
	Passed     int                                     `json:"passed"`
	Failed     int                                     `json:"failed"`
	Skipped    int                                     `json:"skipped"`
	Duration   time.Duration                           `json:"duration_ms"`
	ByLanguage map[types.Language]*LanguageTestSummary `json:"by_language"`
}

// LanguageTestSummary represents test summary for a language
type LanguageTestSummary struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
	Files   int `json:"files"`
}

// NewTestRunner creates a new test runner
func NewTestRunner(config TestConfig) *TestRunner {
	if config.WorkDir == "" {
		config.WorkDir = "."
	}
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Minute
	}
	if config.MaxParallel == 0 {
		config.MaxParallel = 4
	}

	return &TestRunner{
		workDir:  config.WorkDir,
		timeout:  config.Timeout,
		parallel: config.Parallel,
		verbose:  config.Verbose,
	}
}

// RunGoTests runs Go tests
func (r *TestRunner) RunGoTests(pattern string) (*TestResult, error) {
	result := &TestResult{
		Language:  types.Go,
		StartTime: time.Now(),
		Command:   "go test",
	}

	args := []string{"test"}

	// Add pattern
	if pattern != "" {
		args = append(args, pattern)
	} else {
		args = append(args, "./...")
	}

	// Add verbose flag
	if r.verbose {
		args = append(args, "-v")
	}

	// Run go test
	cmd := exec.Command("go", args...)
	cmd.Dir = r.workDir

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Output = stdout.String() + "\n" + stderr.String()

	// Parse test results
	result.Suite = r.parseGoTestOutput(result.Output)

	if err != nil {
		result.Error = err.Error()
	}

	// Find test files
	result.TestFiles = r.findGoTestFiles()

	return result, nil
}

// RunNodeJSTests runs Node.js/Jest tests
func (r *TestRunner) RunNodeJSTests(pattern string) (*TestResult, error) {
	result := &TestResult{
		Language:  types.NodeJS,
		StartTime: time.Now(),
		Command:   "jest",
	}

	// Check if jest is installed
	if _, err := exec.LookPath("npx"); err != nil {
		return nil, fmt.Errorf("npx not found: %w", err)
	}

	args := []string{"jest"}

	// Add pattern
	if pattern != "" {
		args = append(args, pattern)
	}

	// Add verbose flag
	if r.verbose {
		args = append(args, "--verbose")
	}

	// Add JSON output for parsing
	args = append(args, "--json", "--outputFile=jest-results.json")

	// Run jest
	cmd := exec.Command("npx", args...)
	cmd.Dir = r.workDir

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Output = stdout.String() + "\n" + stderr.String()

	// Parse JSON results
	jsonFile := filepath.Join(r.workDir, "jest-results.json")
	if _, err := os.Stat(jsonFile); err == nil {
		result.Suite = r.parseJestOutput(jsonFile)
		defer os.Remove(jsonFile)
	} else {
		// Fallback to parsing output
		result.Suite = r.parseJestOutputFromText(result.Output)
	}

	if err != nil {
		result.Error = err.Error()
	}

	// Find test files
	result.TestFiles = r.findNodeJSTestFiles()

	return result, nil
}

// RunPythonTests runs Python/pytest tests
func (r *TestRunner) RunPythonTests(pattern string) (*TestResult, error) {
	result := &TestResult{
		Language:  types.Python,
		StartTime: time.Now(),
		Command:   "pytest",
	}

	// Check if pytest is installed
	if _, err := exec.LookPath("pytest"); err != nil {
		return nil, fmt.Errorf("pytest not found: %w", err)
	}

	args := []string{"pytest"}

	// Add pattern
	if pattern != "" {
		args = append(args, pattern)
	}

	// Add verbose flag
	if r.verbose {
		args = append(args, "-v")
	}

	// Add JSON output for parsing
	args = append(args, "--json= pytest-results.json")

	// Run pytest
	cmd := exec.Command("pytest", args...)
	cmd.Dir = r.workDir

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Output = stdout.String() + "\n" + stderr.String()

	// Parse output
	result.Suite = r.parsePytestOutput(result.Output)

	if err != nil {
		result.Error = err.Error()
	}

	// Find test files
	result.TestFiles = r.findPythonTestFiles()

	return result, nil
}

// RunTests automatically detects language and runs appropriate test command
func (r *TestRunner) RunTests(pattern string) (*TestResult, error) {
	// Detect language
	if r.hasGoFiles() {
		return r.RunGoTests(pattern)
	} else if r.hasNodeFiles() {
		return r.RunNodeJSTests(pattern)
	} else if r.hasPythonFiles() {
		return r.RunPythonTests(pattern)
	}

	return nil, fmt.Errorf("unable to detect project language")
}

// RunAllTests runs all tests in the project
func (r *TestRunner) RunAllTests() (map[types.Language][]*TestResult, error) {
	results := make(map[types.Language][]*TestResult)

	// Run Go tests
	if r.hasGoFiles() {
		result, err := r.RunGoTests("")
		if err != nil {
			return nil, err
		}
		results[types.Go] = []*TestResult{result}
	}

	// Run Node.js tests
	if r.hasNodeFiles() {
		result, err := r.RunNodeJSTests("")
		if err != nil {
			return nil, err
		}
		results[types.NodeJS] = []*TestResult{result}
	}

	// Run Python tests
	if r.hasPythonFiles() {
		result, err := r.RunPythonTests("")
		if err != nil {
			return nil, err
		}
		results[types.Python] = []*TestResult{result}
	}

	return results, nil
}

// Parsing methods

func (r *TestRunner) parseGoTestOutput(output string) *types.TestSuite {
	suite := &types.TestSuite{
		File:      "all",
		Framework: "go test",
		Results:   []*types.TestResult{},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	scanner := bufio.NewScanner(strings.NewReader(output))

	var currentTest string
	var failed bool
	var outputLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// Parse test start
		if strings.HasPrefix(line, "=== RUN") {
			currentTest = strings.TrimSpace(strings.TrimPrefix(line, "=== RUN"))
			failed = false
			outputLines = nil
			continue
		}

		// Parse test result
		if strings.HasPrefix(line, "--- PASS") {
			suite.Passed++
			suite.Total++
			suite.Results = append(suite.Results, &types.TestResult{
				Name:   currentTest,
				Status: types.TestStatusPassed,
			})
			continue
		}

		if strings.HasPrefix(line, "--- FAIL") {
			suite.Failed++
			suite.Total++
			failed = true
			result := &types.TestResult{
				Name:   currentTest,
				Status: types.TestStatusFailed,
			}

			// Collect error output
			if len(outputLines) > 0 {
				result.Error = strings.Join(outputLines, "\n")
			}

			suite.Results = append(suite.Results, result)
			continue
		}

		if strings.HasPrefix(line, "--- SKIP") {
			suite.Skipped++
			suite.Total++
			suite.Results = append(suite.Results, &types.TestResult{
				Name:   currentTest,
				Status: types.TestStatusSkipped,
			})
			continue
		}

		// Collect output for failing tests
		if currentTest != "" && failed {
			outputLines = append(outputLines, line)
		}
	}

	return suite
}

func (r *TestRunner) parseJestOutput(jsonPath string) *types.TestSuite {
	suite := &types.TestSuite{
		Framework: "jest",
		Results:   []*types.TestResult{},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return suite
	}

	var jestResult struct {
		NumTotalTests   int `json:"numTotalTests"`
		NumPassedTests  int `json:"numPassedTests"`
		NumFailedTests  int `json:"numFailedTests"`
		NumPendingTests int `json:"numPendingTests"`
		TestResults     []struct {
			Name            string   `json:"name"`
			Status          string   `json:"status"`
			Duration        int      `json:"duration"`
			FailureMessages []string `json:"failureMessages"`
		} `json:"testResults"`
	}

	if err := json.Unmarshal(data, &jestResult); err != nil {
		return suite
	}

	suite.Total = jestResult.NumTotalTests
	suite.Passed = jestResult.NumPassedTests
	suite.Failed = jestResult.NumFailedTests
	suite.Skipped = jestResult.NumPendingTests

	for _, tr := range jestResult.TestResults {
		result := &types.TestResult{
			Name:     tr.Name,
			Duration: fmt.Sprintf("%d", tr.Duration),
		}

		switch tr.Status {
		case "passed":
			result.Status = types.TestStatusPassed
		case "failed":
			result.Status = types.TestStatusFailed
			if len(tr.FailureMessages) > 0 {
				result.Error = tr.FailureMessages[0]
			}
		case "pending":
			result.Status = types.TestStatusSkipped
		}

		suite.Results = append(suite.Results, result)
	}

	return suite
}

func (r *TestRunner) parseJestOutputFromText(output string) *types.TestSuite {
	suite := &types.TestSuite{
		Framework: "jest",
		Results:   []*types.TestResult{},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Simple regex parsing for Jest output
	passRe := regexp.MustCompile(`✓ (.+) \((\d+) ms\)`)
	failRe := regexp.MustCompile(`✕ (.+) \((\d+) ms\)`)

	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()

		if matches := passRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Passed++
			suite.Total++
			suite.Results = append(suite.Results, &types.TestResult{
				Name:     matches[1],
				Status:   types.TestStatusPassed,
				Duration: matches[2],
			})
			continue
		}

		if matches := failRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Failed++
			suite.Total++
			suite.Results = append(suite.Results, &types.TestResult{
				Name:     matches[1],
				Status:   types.TestStatusFailed,
				Duration: matches[2],
			})
			continue
		}
	}

	return suite
}

func (r *TestRunner) parsePytestOutput(output string) *types.TestSuite {
	suite := &types.TestSuite{
		Framework: "pytest",
		Results:   []*types.TestResult{},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Parse test results
	passRe := regexp.MustCompile(`PASSED\s+(\S+)`)
	failRe := regexp.MustCompile(`FAILED\s+(\S+)`)
	skipRe := regexp.MustCompile(`SKIPPED\s+(\S+)`)

	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()

		if matches := passRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Passed++
			suite.Total++
			suite.Results = append(suite.Results, &types.TestResult{
				Name:   matches[1],
				Status: types.TestStatusPassed,
			})
			continue
		}

		if matches := failRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Failed++
			suite.Total++
			suite.Results = append(suite.Results, &types.TestResult{
				Name:   matches[1],
				Status: types.TestStatusFailed,
			})
			continue
		}

		if matches := skipRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Skipped++
			suite.Total++
			suite.Results = append(suite.Results, &types.TestResult{
				Name:   matches[1],
				Status: types.TestStatusSkipped,
			})
			continue
		}
	}

	// Try to get summary line
	summaryRe := regexp.MustCompile(`(\d+) passed, (\d+) failed, (\d+) skipped`)
	if matches := summaryRe.FindStringSubmatch(output); len(matches) > 3 {
		passed, _ := strconv.Atoi(matches[1])
		failed, _ := strconv.Atoi(matches[2])
		skipped, _ := strconv.Atoi(matches[3])

		suite.Passed = passed
		suite.Failed = failed
		suite.Skipped = skipped
		suite.Total = passed + failed + skipped
	}

	return suite
}

// Test file discovery

func (r *TestRunner) findGoTestFiles() []string {
	var files []string

	matches, err := filepath.Glob(filepath.Join(r.workDir, "*_test.go"))
	if err == nil {
		files = append(files, matches...)
	}

	// Recursive search
	filepath.Walk(r.workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})

	return files
}

func (r *TestRunner) findNodeJSTestFiles() []string {
	var files []string

	patterns := []string{"*.test.js", "*.spec.js", "*.test.ts", "*.spec.ts"}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(r.workDir, pattern))
		if err == nil {
			files = append(files, matches...)
		}
	}

	// Check common test directories
	testDirs := []string{"__tests__", "test", "tests"}
	for _, dir := range testDirs {
		testPath := filepath.Join(r.workDir, dir)
		filepath.Walk(testPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && (strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".ts")) {
				files = append(files, path)
			}
			return nil
		})
	}

	return files
}

func (r *TestRunner) findPythonTestFiles() []string {
	var files []string

	patterns := []string{"test_*.py", "*_test.py"}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(r.workDir, pattern))
		if err == nil {
			files = append(files, matches...)
		}
	}

	// Recursive search
	filepath.Walk(r.workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			base := filepath.Base(path)
			if strings.HasPrefix(base, "test_") || strings.HasSuffix(base, "_test.py") {
				files = append(files, path)
			}
		}
		return nil
	})

	return files
}

// Language detection

func (r *TestRunner) hasGoFiles() bool {
	matches, _ := filepath.Glob(filepath.Join(r.workDir, "*.go"))
	return len(matches) > 0
}

func (r *TestRunner) hasNodeFiles() bool {
	// Check for package.json
	if _, err := os.Stat(filepath.Join(r.workDir, "package.json")); err == nil {
		return true
	}

	// Check for JS/TS files
	jsMatches, _ := filepath.Glob(filepath.Join(r.workDir, "*.js"))
	tsMatches, _ := filepath.Glob(filepath.Join(r.workDir, "*.ts"))
	return len(jsMatches) > 0 || len(tsMatches) > 0
}

func (r *TestRunner) hasPythonFiles() bool {
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

// Summary and formatting

// Summarize creates a summary from multiple test results
func (r *TestRunner) Summarize(results map[types.Language][]*TestResult) *TestSummary {
	summary := &TestSummary{
		ByLanguage: make(map[types.Language]*LanguageTestSummary),
	}

	for lang, langResults := range results {
		langSummary := &LanguageTestSummary{
			Files: len(langResults),
		}

		for _, result := range langResults {
			if result.Suite != nil {
				langSummary.Total += result.Suite.Total
				langSummary.Passed += result.Suite.Passed
				langSummary.Failed += result.Suite.Failed
				langSummary.Skipped += result.Suite.Skipped

				summary.Total += result.Suite.Total
				summary.Passed += result.Suite.Passed
				summary.Failed += result.Suite.Failed
				summary.Skipped += result.Suite.Skipped
				summary.Duration += result.Duration
			}
		}

		summary.ByLanguage[lang] = langSummary
	}

	return summary
}

// FormatResult formats a test result for display
func (r *TestRunner) FormatResult(result *TestResult) string {
	if result == nil {
		return "No test results"
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🧪 Test Results (%s)\n", result.Language))
	sb.WriteString(strings.Repeat("=", 60))
	sb.WriteString("\n\n")

	if result.Suite != nil {
		// Summary line
		sb.WriteString(fmt.Sprintf("Total:  %d tests\n", result.Suite.Total))
		sb.WriteString(fmt.Sprintf("Passed: %d ✅\n", result.Suite.Passed))
		sb.WriteString(fmt.Sprintf("Failed: %d ❌\n", result.Suite.Failed))
		sb.WriteString(fmt.Sprintf("Skipped: %d ⏭️\n", result.Suite.Skipped))

		// Pass rate
		if result.Suite.Total > 0 {
			passRate := float64(result.Suite.Passed) / float64(result.Suite.Total) * 100
			sb.WriteString(fmt.Sprintf("Pass rate: %.1f%%\n", passRate))
		}

		sb.WriteString("\n")

		// Failed tests details
		if result.Suite.Failed > 0 {
			sb.WriteString("Failed Tests:\n")
			sb.WriteString("-------------\n")

			for _, tr := range result.Suite.Results {
				if tr.Status == types.TestStatusFailed {
					sb.WriteString(fmt.Sprintf("  ❌ %s\n", tr.Name))
					if tr.Error != "" {
						// Indent error message
						errorLines := strings.Split(tr.Error, "\n")
						for _, line := range errorLines {
							sb.WriteString(fmt.Sprintf("     %s\n", line))
						}
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	sb.WriteString(fmt.Sprintf("\nDuration: %v\n", result.Duration))
	sb.WriteString(fmt.Sprintf("Command: %s\n", result.Command))

	if result.Error != "" {
		sb.WriteString(fmt.Sprintf("\nError: %s\n", result.Error))
	}

	return sb.String()
}

// FormatSummary formats a test summary for display
func (r *TestRunner) FormatSummary(summary *TestSummary) string {
	if summary == nil {
		return "No test summary available"
	}

	var sb strings.Builder

	sb.WriteString("📊 Test Summary\n")
	sb.WriteString(strings.Repeat("=", 60))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("Total Tests: %d\n", summary.Total))
	sb.WriteString(fmt.Sprintf("Passed:      %d ✅\n", summary.Passed))
	sb.WriteString(fmt.Sprintf("Failed:      %d ❌\n", summary.Failed))
	sb.WriteString(fmt.Sprintf("Skipped:     %d ⏭️\n", summary.Skipped))

	if summary.Total > 0 {
		passRate := float64(summary.Passed) / float64(summary.Total) * 100
		sb.WriteString(fmt.Sprintf("Pass Rate:   %.1f%%\n", passRate))
	}

	sb.WriteString(fmt.Sprintf("Duration:    %v\n", summary.Duration))
	sb.WriteString("\n")

	// By language
	if len(summary.ByLanguage) > 0 {
		sb.WriteString("By Language:\n")
		sb.WriteString("------------\n")

		for lang, langSummary := range summary.ByLanguage {
			sb.WriteString(fmt.Sprintf("\n%s:\n", lang))
			sb.WriteString(fmt.Sprintf("  Files:  %d\n", langSummary.Files))
			sb.WriteString(fmt.Sprintf("  Tests:  %d total, %d passed, %d failed, %d skipped\n",
				langSummary.Total, langSummary.Passed, langSummary.Failed, langSummary.Skipped))
		}
	}

	return sb.String()
}

// GetFailedTests returns only the failed tests from a result
func (r *TestRunner) GetFailedTests(result *TestResult) []*types.TestResult {
	if result == nil || result.Suite == nil {
		return nil
	}

	var failed []*types.TestResult
	for _, tr := range result.Suite.Results {
		if tr.Status == types.TestStatusFailed {
			failed = append(failed, tr)
		}
	}
	return failed
}

// WriteJUnitXML writes test results in JUnit XML format (for CI)
func (r *TestRunner) WriteJUnitXML(result *TestResult, outputPath string) error {
	if result == nil || result.Suite == nil {
		return fmt.Errorf("no test results to write")
	}

	// Simple JUnit XML generation
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
    <testsuite name="%s" tests="%d" failures="%d" skipped="%d" time="%.3f">
`,
		result.Language,
		result.Suite.Total,
		result.Suite.Failed,
		result.Suite.Skipped,
		result.Duration.Seconds())

	for _, tr := range result.Suite.Results {
		xml += fmt.Sprintf(`        <testcase name="%s" time="0.0">`, tr.Name)
		if tr.Status == types.TestStatusFailed {
			xml += fmt.Sprintf("\n            <failure message=\"%s\"/>\n        ", tr.Error)
		}
		xml += "</testcase>\n"
	}

	xml += `    </testsuite>
</testsuites>`

	return os.WriteFile(outputPath, []byte(xml), 0644)
}
