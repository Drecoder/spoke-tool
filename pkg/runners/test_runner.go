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
	"strings"
	"time"

	"example.com/spoke-tool/api/types"
)

// TestRunner runs tests and collects results
type TestRunner struct {
	workDir  string
	timeout  time.Duration
	parallel bool
	verbose  bool
}

// TestConfig configures the test runner
type TestConfig struct {
	WorkDir     string
	Timeout     time.Duration
	Parallel    bool
	MaxParallel int
	Verbose     bool
	Env         map[string]string
	Pattern     string
}

// TestResult represents the result of a test run
type TestResult struct {
	Language  types.Language   `json:"language"`
	Suite     *types.TestSuite `json:"suite"`
	Command   string           `json:"command"`
	Output    string           `json:"output,omitempty"`
	Error     string           `json:"error,omitempty"`
	StartTime time.Time        `json:"start_time"`
	EndTime   time.Time        `json:"end_time"`
	Duration  time.Duration    `json:"duration_ms"`
	TestFiles []string         `json:"test_files,omitempty"`
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

type LanguageTestSummary struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
	Files   int `json:"files"`
}

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

// --- Execution Methods ---

func (r *TestRunner) RunGoTests(pattern string) (*TestResult, error) {
	result := &TestResult{
		Language:  types.Go,
		StartTime: time.Now(),
		Command:   "go test",
	}

	args := []string{"test"}
	if pattern != "" {
		args = append(args, pattern)
	} else {
		args = append(args, "./...")
	}

	if r.verbose {
		args = append(args, "-v")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = r.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Output = stdout.String() + "\n" + stderr.String()
	result.Suite = r.parseGoTestOutput(result.Output)

	if err != nil {
		result.Error = err.Error()
	}
	result.TestFiles = r.findGoTestFiles()

	return result, nil
}

func (r *TestRunner) RunNodeJSTests(pattern string) (*TestResult, error) {
	result := &TestResult{
		Language:  types.NodeJS,
		StartTime: time.Now(),
		Command:   "jest",
	}

	if _, err := exec.LookPath("npx"); err != nil {
		return nil, fmt.Errorf("npx not found: %w", err)
	}

	args := []string{"jest"}
	if pattern != "" {
		args = append(args, pattern)
	}
	if r.verbose {
		args = append(args, "--verbose")
	}
	args = append(args, "--json", "--outputFile=jest-results.json")

	cmd := exec.Command("npx", args...)
	cmd.Dir = r.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Output = stdout.String() + "\n" + stderr.String()

	jsonFile := filepath.Join(r.workDir, "jest-results.json")
	if _, err := os.Stat(jsonFile); err == nil {
		result.Suite = r.parseJestOutput(jsonFile)
		defer os.Remove(jsonFile)
	} else {
		result.Suite = r.parseJestOutputFromText(result.Output)
	}

	if err != nil {
		result.Error = err.Error()
	}
	result.TestFiles = r.findNodeJSTestFiles()

	return result, nil
}

func (r *TestRunner) RunPythonTests(pattern string) (*TestResult, error) {
	result := &TestResult{
		Language:  types.Python,
		StartTime: time.Now(),
		Command:   "pytest",
	}

	if _, err := exec.LookPath("pytest"); err != nil {
		return nil, fmt.Errorf("pytest not found: %w", err)
	}

	args := []string{"pytest"}
	if pattern != "" {
		args = append(args, pattern)
	}
	if r.verbose {
		args = append(args, "-v")
	}

	cmd := exec.Command("pytest", args...)
	cmd.Dir = r.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Output = stdout.String() + "\n" + stderr.String()
	result.Suite = r.parsePytestOutput(result.Output)

	if err != nil {
		result.Error = err.Error()
	}
	result.TestFiles = r.findPythonTestFiles()

	return result, nil
}

// --- Parsing Methods ---

func (r *TestRunner) parseGoTestOutput(output string) *types.TestSuite {
	suite := &types.TestSuite{
		File:      "all",
		Framework: "go test",
		Results:   []types.TestResult{}, // FIXED: Value slice
		Timestamp: time.Now().Format(time.RFC3339),
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	var currentTest string
	var outputLines []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "=== RUN") {
			currentTest = strings.TrimSpace(strings.TrimPrefix(line, "=== RUN"))
			outputLines = nil
			continue
		}

		if strings.HasPrefix(line, "--- PASS") {
			suite.Passed++
			suite.Total++
			suite.Results = append(suite.Results, types.TestResult{ // FIXED: Append value
				Name:   currentTest,
				Status: types.TestStatusPassed,
			})
			continue
		}

		if strings.HasPrefix(line, "--- FAIL") {
			suite.Failed++
			suite.Total++
			res := types.TestResult{
				Name:   currentTest,
				Status: types.TestStatusFailed,
			}
			if len(outputLines) > 0 {
				res.Error = strings.Join(outputLines, "\n")
			}
			suite.Results = append(suite.Results, res) // FIXED: Append value
			continue
		}

		if strings.HasPrefix(line, "--- SKIP") {
			suite.Skipped++
			suite.Total++
			suite.Results = append(suite.Results, types.TestResult{ // FIXED: Append value
				Name:   currentTest,
				Status: types.TestStatusSkipped,
			})
			continue
		}
	}
	return suite
}

func (r *TestRunner) parseJestOutput(jsonPath string) *types.TestSuite {
	suite := &types.TestSuite{
		Framework: "jest",
		Results:   []types.TestResult{}, // FIXED
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

	if err := json.Unmarshal(data, &jestResult); err == nil {
		suite.Total = jestResult.NumTotalTests
		suite.Passed = jestResult.NumPassedTests
		suite.Failed = jestResult.NumFailedTests
		suite.Skipped = jestResult.NumPendingTests

		for _, tr := range jestResult.TestResults {
			res := types.TestResult{
				Name:     tr.Name,
				Duration: fmt.Sprintf("%d", tr.Duration),
			}
			switch tr.Status {
			case "passed":
				res.Status = types.TestStatusPassed
			case "failed":
				res.Status = types.TestStatusFailed
				if len(tr.FailureMessages) > 0 {
					res.Error = tr.FailureMessages[0]
				}
			case "pending":
				res.Status = types.TestStatusSkipped
			}
			suite.Results = append(suite.Results, res)
		}
	}
	return suite
}

func (r *TestRunner) parseJestOutputFromText(output string) *types.TestSuite {
	suite := &types.TestSuite{
		Framework: "jest",
		Results:   []types.TestResult{}, // FIXED
		Timestamp: time.Now().Format(time.RFC3339),
	}

	passRe := regexp.MustCompile(`✓ (.+) \((\d+) ms\)`)
	failRe := regexp.MustCompile(`✕ (.+) \((\d+) ms\)`)
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		if matches := passRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Passed++
			suite.Total++
			suite.Results = append(suite.Results, types.TestResult{
				Name:     matches[1],
				Status:   types.TestStatusPassed,
				Duration: matches[2],
			})
		} else if matches := failRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Failed++
			suite.Total++
			suite.Results = append(suite.Results, types.TestResult{
				Name:     matches[1],
				Status:   types.TestStatusFailed,
				Duration: matches[2],
			})
		}
	}
	return suite
}

func (r *TestRunner) parsePytestOutput(output string) *types.TestSuite {
	suite := &types.TestSuite{
		Framework: "pytest",
		Results:   []types.TestResult{}, // FIXED
		Timestamp: time.Now().Format(time.RFC3339),
	}

	passRe := regexp.MustCompile(`PASSED\s+(\S+)`)
	failRe := regexp.MustCompile(`FAILED\s+(\S+)`)
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		if matches := passRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Passed++
			suite.Total++
			suite.Results = append(suite.Results, types.TestResult{
				Name:   matches[1],
				Status: types.TestStatusPassed,
			})
		} else if matches := failRe.FindStringSubmatch(line); len(matches) > 1 {
			suite.Failed++
			suite.Total++
			suite.Results = append(suite.Results, types.TestResult{
				Name:   matches[1],
				Status: types.TestStatusFailed,
			})
		}
	}
	return suite
}

// --- Utility & Formatting ---

func (r *TestRunner) GetFailedTests(result *TestResult) []types.TestResult { // FIXED: Returns values
	if result == nil || result.Suite == nil {
		return nil
	}
	var failed []types.TestResult
	for _, tr := range result.Suite.Results {
		if tr.Status == types.TestStatusFailed {
			failed = append(failed, tr)
		}
	}
	return failed
}

func (r *TestRunner) hasGoFiles() bool {
	matches, _ := filepath.Glob(filepath.Join(r.workDir, "*.go"))
	return len(matches) > 0
}

func (r *TestRunner) hasNodeFiles() bool {
	if _, err := os.Stat(filepath.Join(r.workDir, "package.json")); err == nil {
		return true
	}
	js, _ := filepath.Glob(filepath.Join(r.workDir, "*.js"))
	ts, _ := filepath.Glob(filepath.Join(r.workDir, "*.ts"))
	return len(js) > 0 || len(ts) > 0
}

func (r *TestRunner) hasPythonFiles() bool {
	matches, _ := filepath.Glob(filepath.Join(r.workDir, "*.py"))
	return len(matches) > 0
}

func (r *TestRunner) findGoTestFiles() []string {
	var files []string
	filepath.Walk(r.workDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func (r *TestRunner) findNodeJSTestFiles() []string {
	var files []string
	filepath.Walk(r.workDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && (strings.HasSuffix(path, ".test.js") || strings.HasSuffix(path, ".spec.ts")) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func (r *TestRunner) findPythonTestFiles() []string {
	var files []string
	filepath.Walk(r.workDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && (strings.HasPrefix(filepath.Base(path), "test_") || strings.HasSuffix(path, "_test.py")) {
			files = append(files, path)
		}
		return nil
	})
	return files
}
