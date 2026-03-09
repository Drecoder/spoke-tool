package types

// TestFramework represents the testing framework used
type TestFramework string

const (
	TestFrameworkGoTest  TestFramework = "go test"
	TestFrameworkJest    TestFramework = "jest"
	TestFrameworkMocha   TestFramework = "mocha"
	TestFrameworkPytest  TestFramework = "pytest"
	TestFrameworkUnittest TestFramework = "unittest"
)

// TestFile represents a test file
type TestFile struct {
	Path        string        `json:"path"`
	Language    Language      `json:"language"`
	Content     string        `json:"content"`
	Framework   TestFramework `json:"framework"`
	Functions   []string      `json:"functions_tested"` // Functions this file tests
	Coverage    float64       `json:"coverage,omitempty"`
	LastRun     string        `json:"last_run,omitempty"`
	LastStatus  TestStatus    `json:"last_status,omitempty"`
}

// TestStatus represents the status of a test run
type TestStatus string

const (
	TestStatusPassed  TestStatus = "passed"
	TestStatusFailed  TestStatus = "failed"
	TestStatusSkipped TestStatus = "skipped"
	TestStatusPending TestStatus = "pending"
)

// TestResult represents a single test result
type TestResult struct {
	Name       string     `json:"name"`
	Status     TestStatus `json:"status"`
	Duration   string     `json:"duration_ms"`
	Error      string     `json:"error,omitempty"`
	StackTrace string     `json:"stack_trace,omitempty"`
	Function   string     `json:"function,omitempty"` // Function being tested
}

// TestSuite represents a collection of test results
type TestSuite struct {
	File      string       `json:"file"`
	Framework TestFramework `json:"framework"`
	Results   []TestResult `json:"results"`
	Total     int          `json:"total"`
	Passed    int          `json:"passed"`
	Failed    int          `json:"failed"`
	Skipped   int          `json:"skipped"`
	Duration  string       `json:"duration_ms"`
	Timestamp string       `json:"timestamp"`
}

// TestSuggestion represents a generated test suggestion
type TestSuggestion struct {
	Language     Language      `json:"language"`
	FunctionName string        `json:"function_name"`
	FunctionCode string        `json:"function_code"`
	TestCode     string        `json:"test_code"`
	TestFilePath string        `json:"test_file_path"`
	Framework    TestFramework `json:"framework"`
	Description  string        `json:"description"`
	Confidence   float64       `json:"confidence"` // 0-1
}

// TestCoverage represents code coverage information
type TestCoverage struct {
	Language    Language          `json:"language"`
	Overall     float64           `json:"overall_percent"`
	ByFile      map[string]float64 `json:"by_file"`
	ByFunction  map[string]float64 `json:"by_function"`
	Uncovered   []string          `json:"uncovered_lines"` // Lines without coverage
	Timestamp   string            `json:"timestamp"`
}