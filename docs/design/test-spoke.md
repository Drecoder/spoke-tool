# Test Spoke Design

## 🎯 Purpose

The Test Spoke automatically generates and maintains unit tests for code changes across multiple languages. It analyzes code, identifies untested functions, generates appropriate tests, runs them, and reports results—all without modifying source code.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Test Spoke                               │
│                                                               │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐   │
│  │   Analyzer   │───▶│  Generator   │───▶│    Runner    │   │
│  │              │    │              │    │              │   │
│  │ • AST Parser │    │ • Templates  │    │ • Execute    │   │
│  │ • Coverage   │    │ • SLM Calls  │    │ • Collect    │   │
│  │ • Dep Graph  │    │ • Mocks      │    │ • Report     │   │
│  └──────────────┘    └──────────────┘    └──────┬───────┘   │
│                                                  │           │
│                                          ┌───────▼───────┐   │
│                                          │   Interpreter │   │
│                                          │               │   │
│                                          │ • Parse Error │   │
│                                          │ • Explain Why │   │
│                                          │ • NO FIXES    │   │
│                                          └───────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 🔍 Component Details

### 1. Analyzer

The analyzer examines code to find what needs testing.

**Responsibilities:**
- Parse code into AST (Abstract Syntax Tree)
- Identify functions, methods, and classes
- Map dependencies between files
- Check existing test coverage
- Flag untested public APIs

**Language-Specific Parsers:**

| Language | Parser | Test Detection |
|----------|--------|----------------|
| Go | `go/ast` | `*_test.go` files |
| Node.js | `@babel/parser` | `*.test.js`, `*.spec.js` |
| Python | `ast` module | `test_*.py` files |

**Output:**
```go
type AnalysisResult struct {
    Language     Language
    Files        []CodeFile
    Functions    []Function
    TestCoverage map[string]bool  // function name -> has test
    Dependencies map[string][]string
    Complexity   map[string]int
}
```

### 2. Generator

The generator creates test files using SLMs.

**Responsibilities:**
- Select appropriate model (DeepSeek 7B for complex, Gemma 2B for simple)
- Build language-specific prompts
- Generate test code
- Create mock objects
- Handle edge cases

**Test Generation Strategy:**

```
Function Input
    ↓
Analyze Signature
    ↓
Identify Test Cases:
    ├── Happy Path
    ├── Edge Cases
    ├── Error Conditions
    └── Table-Driven (where appropriate)
    ↓
Generate Test Code
    ↓
Validate Syntax
    ↓
Output Test File
```

**Language-Specific Templates:**

#### Go Template
```go
func Test{{.FunctionName}}(t *testing.T) {
    tests := []struct {
        name string
        {{.InputParams}}
        expected {{.ReturnType}}
        expectError bool
    }{
        // Generated test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := {{.FunctionName}}(tt.args)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

#### Jest Template
```javascript
describe('{{.FunctionName}}', () => {
    // Generated test cases
    test.each([
        // Test cases
    ])('%s', (_, input, expected) => {
        expect({{.FunctionName}}(input)).toBe(expected);
    });
});
```

#### Pytest Template
```python
@pytest.mark.parametrize("input,expected", [
    # Generated test cases
])
def test_{{.function_name}}(input, expected):
    assert {{.function_name}}(input) == expected
```

### 3. Runner

The runner executes tests and collects results.

**Responsibilities:**
- Run language-specific test commands
- Capture output and exit codes
- Parse test results
- Measure coverage
- Format results for reporting

**Test Commands:**

| Language | Command | Coverage Flag |
|----------|---------|---------------|
| Go | `go test ./...` | `-cover` |
| Node.js | `npm test` or `jest` | `--coverage` |
| Python | `pytest` | `--cov=.` |

**Output:**
```go
type TestResult struct {
    Passed     bool
    Failed     int
    Total      int
    Coverage   float64
    Failures   []TestFailure
    Duration   time.Duration
}
```

### 4. Interpreter

The interpreter analyzes test failures and explains WHY they failed—without suggesting fixes.

**Responsibilities:**
- Parse error messages
- Correlate failures with source code
- Generate human-readable explanations
- **NEVER** suggest code changes

**Failure Analysis Flow:**

```
Test Failure
    ↓
Parse Error: "expected 3, got 4"
    ↓
Locate Source: Add(1, 2) returns 4
    ↓
Analyze: Function adds extra 1
    ↓
Explain: "The function returns a+b+1 but test expects a+b. The calculation is off by +1."
    ↓
Report to Developer (NO FIX SUGGESTED)
```

**Example Explanations:**

```go
// Bad (DO NOT DO THIS):
"Change the function to return a+b instead of a+b+1"

// Good (DO THIS):
"The function returns a+b+1 but the test expects a+b. The calculation includes an extra +1."
```

## 🔄 Workflow States

```
                    ┌─────────────┐
                    │ Code Change │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │  Analyze    │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │ Any Gaps?   │───No───┐
                    └──────┬──────┘        │
                          Yes               │
                           ▼                 │
                    ┌─────────────┐        │
                    │  Generate   │        │
                    │    Tests    │        │
                    └──────┬──────┘        │
                           ▼                 │
                    ┌─────────────┐        │
                    │    Run      │        │
                    │    Tests    │        │
                    └──────┬──────┘        │
                           ▼                 │
                    ┌─────────────┐        │
                    │ All Pass?   │──No──┐ │
                    └──────┬──────┘      │ │
                          Yes             │ │
                           ▼               │ │
                    ┌─────────────┐      │ │
                    │  Update     │      │ │
                    │   README    │      │ │
                    └─────────────┘      │ │
                           ▼               │ │
                    ┌─────────────┐      │ │
                    │    Done     │      │ │
                    └─────────────┘      │ │
                                         ▼ ▼
                               ┌─────────────────┐
                               │   Analyze       │
                               │   Failures      │
                               └────────┬────────┘
                                        ▼
                               ┌─────────────────┐
                               │   Report to     │
                               │   Developer     │
                               └─────────────────┘
```

## 📝 SLM Prompts

### Test Generation Prompt
```
You are a {language} testing expert using {framework}.

Generate comprehensive unit tests for this function:

Function: {functionName}
Code:
{code}
Dependencies: {dependencies}

Requirements:
- Use {framework} syntax
- Cover happy path, edge cases, errors
- Use table-driven tests where appropriate
- Include mocks for dependencies
- Return ONLY the test code, no explanations

Test File Path: {testFilePath}
```

### Failure Analysis Prompt
```
Explain why this test failed. DO NOT suggest code fixes.

Language: {language}
Test Name: {testName}
Error: {errorMessage}

Test Code:
{testCode}

Source Code:
{sourceCode}

Provide a clear explanation of WHY the test failed.
Focus on the mismatch between expected and actual behavior.
DO NOT suggest changes to the code.
```

## 🧪 Example Walkthrough

### Input Function (Go)
```go
func CalculateDiscount(price float64, isMember bool) float64 {
    if isMember {
        return price * 0.9 // 10% discount
    }
    return price
}
```

### Generated Test
```go
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name     string
        price    float64
        isMember bool
        expected float64
    }{
        {"non-member full price", 100.0, false, 100.0},
        {"member 10% discount", 100.0, true, 90.0},
        {"zero price non-member", 0.0, false, 0.0},
        {"zero price member", 0.0, true, 0.0},
        {"negative price", -50.0, true, -45.0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateDiscount(tt.price, tt.isMember)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### If Test Fails
```
Error: got -45.0, want -45.0? (they match, so this passes)
But if there was an error, analysis would explain why.
```

## 📊 Configuration Options

```yaml
test_spoke:
  enabled: true
  auto_run: true                     # Run tests after generation
  coverage_threshold: 80              # Minimum coverage %
  languages:
    go:
      framework: "testing"
      test_pattern: "*_test.go"
      cover_command: "go test -cover"
    nodejs:
      framework: "jest"
      test_pattern: "*.test.js"
      cover_command: "jest --coverage"
    python:
      framework: "pytest"
      test_pattern: "test_*.py"
      cover_command: "pytest --cov=."
  
  generation:
    model: "codellama:7b"        # For complex tests
    fallback_model: "gemma2:2b"       # For simple tests
    max_tests_per_function: 10
    include_edge_cases: true
    generate_mocks: true
  
  analysis:
    model: "codellama:7b"        # For failure analysis
    detailed_reports: true
  
  reporting:
    format: "console"                 # console, json, html
    show_failures_only: false
    show_coverage: true
```

## 🚦 Error Handling

### Test Failures
```
Test Failed
    ↓
[INTERPRETER] Parse error message
    ↓
[INTERPRETER] Analyze source context
    ↓
[INTERPRETER] Generate explanation
    ↓
[REPORT] Show developer:
    - What test failed
    - Expected vs actual
    - Why it failed
    ↓
[STOP] - Developer must fix
```

### System Errors
```
Error Type          → Action
─────────────────────────────────────────
Config Error        → Exit with code 1
Model Unavailable   → Retry with backoff
Parse Error         → Skip file, log warning
Write Permission    → Exit with code 2
```

## 🔒 Security Considerations

1. **No Code Modification**: Never changes source code
2. **Local Analysis**: All processing on developer machine
3. **Audit Trail**: All test generations logged
4. **Safe Execution**: Tests run in isolated environment
5. **No Auto-Fixes**: Prevents hallucinated "fixes"

## 📈 Performance Optimization

### Squeeze Integration
```go
func (t *TestSpoke) GenerateTests(functions []Function) {
    if squeeze.ShouldThrottle() {
        // Reduce batch size
        functions = prioritizeFunctions(functions)
    }
    
    // Generate in parallel with limits
    semaphore := make(chan struct{}, maxConcurrent)
    for _, fn := range functions {
        semaphore <- struct{}{}
        go func(f Function) {
            defer func() { <-semaphore }()
            t.generateForFunction(f)
        }(fn)
    }
}
```

## 🎯 Design Principles Applied

1. **No Auto-Fixes** - Tests failures always require developer action
2. **Explain Don't Fix** - Analysis explains why, never suggests code changes
3. **Local First** - All processing on developer machine
4. **Language Agnostic** - Support for multiple languages
5. **Resource Aware** - Respects system load via Squeeze
6. **Privacy Preserving** - No code leaves the machine