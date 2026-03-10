
#!/usr/bin/env bash

# create-testdata.sh - Complete testdata setup script
# This script creates all testdata directories and files for the spoke-tool

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# ============================================================================
# Helper Functions
# ============================================================================

print_header() {
    echo -e "\n${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}\n"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${CYAN}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# ============================================================================
# Create Directory Structure
# ============================================================================

print_header "Creating Testdata Directory Structure"

# Expected outputs directories
mkdir -p testdata/expected-outputs/{go,nodejs,python}/{tests,readme}
mkdir -p testdata/expected-outputs/mixed
mkdir -p testdata/expected-outputs/edge-cases

# Sample code directories
mkdir -p testdata/sample-code/go/{simple,complex,errors}
mkdir -p testdata/sample-code/nodejs/{simple,complex,errors}
mkdir -p testdata/sample-code/python/{simple,complex,errors}
mkdir -p testdata/sample-code/mixed-project/{go-server,web-client,scripts}

# Input directories
mkdir -p testdata/input/{go,nodejs,python}

# Benchmark directories
mkdir -p testdata/benchmark/{go,nodejs,python}

# Integration test directories
mkdir -p testdata/integration/{go,nodejs,python}

# Fuzz test directories
mkdir -p testdata/fuzz/go/testdata/fuzz
mkdir -p testdata/fuzz/nodejs
mkdir -p testdata/fuzz/python

print_success "Directory structure created"

# ============================================================================
# Create Expected Outputs - Go
# ============================================================================

print_header "Creating Go Expected Outputs"

# Go test files
touch testdata/expected-outputs/go/tests/math_test.go
touch testdata/expected-outputs/go/tests/calculator_test.go
touch testdata/expected-outputs/go/tests/string_util_test.go
touch testdata/expected-outputs/go/tests/file_util_test.go

# Go README files
touch testdata/expected-outputs/go/readme/README.md
touch testdata/expected-outputs/go/readme/API.md
touch testdata/expected-outputs/go/readme/CONTRIBUTING.md

print_success "Go expected outputs created"

# ============================================================================
# Create Expected Outputs - Node.js
# ============================================================================

print_header "Creating Node.js Expected Outputs"

# Node.js test files
touch testdata/expected-outputs/nodejs/tests/math.test.js
touch testdata/expected-outputs/nodejs/tests/string.test.js
touch testdata/expected-outputs/nodejs/tests/async.test.js
touch testdata/expected-outputs/nodejs/tests/api.test.js

# Node.js README files
touch testdata/expected-outputs/nodejs/readme/README.md
touch testdata/expected-outputs/nodejs/readme/API.md
touch testdata/expected-outputs/nodejs/readme/SETUP.md

print_success "Node.js expected outputs created"

# ============================================================================
# Create Expected Outputs - Python
# ============================================================================

print_header "Creating Python Expected Outputs"

# Python test files
touch testdata/expected-outputs/python/tests/test_math.py
touch testdata/expected-outputs/python/tests/test_string.py
touch testdata/expected-outputs/python/tests/test_file_ops.py
touch testdata/expected-outputs/python/tests/conftest.py

# Python README files
touch testdata/expected-outputs/python/readme/README.md
touch testdata/expected-outputs/python/readme/API.md
touch testdata/expected-outputs/python/readme/INSTALL.md

print_success "Python expected outputs created"

# ============================================================================
# Create Edge Cases
# ============================================================================

print_header "Creating Edge Cases"

touch testdata/expected-outputs/edge-cases/edge_cases_test.go
touch testdata/expected-outputs/edge-cases/edge_cases.test.js
touch testdata/expected-outputs/edge-cases/test_edge_cases.py
touch testdata/expected-outputs/edge-cases/README.md

print_success "Edge cases created"

# ============================================================================
# Create Mixed Project
# ============================================================================

print_header "Creating Mixed Project Outputs"

touch testdata/expected-outputs/mixed/README.md
touch testdata/expected-outputs/mixed/docker-compose.yml
touch testdata/expected-outputs/mixed/Makefile

print_success "Mixed project outputs created"

# ============================================================================
# Create Validation Script
# ============================================================================

print_header "Creating Validation Script"

cat > testdata/expected-outputs/validate.sh << 'EOF'
#!/usr/bin/env bash

# validate.sh - Validate generated outputs against expected outputs

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

TOTAL=0
PASSED=0
FAILED=0
SKIPPED=0

validate_file() {
    local expected=$1
    local generated=$2
    local name=$3
    TOTAL=$((TOTAL + 1))
    
    if [ ! -f "$expected" ]; then
        echo -e "${YELLOW}⚠️  Expected file not found: $expected${NC}"
        SKIPPED=$((SKIPPED + 1))
        return
    fi
    
    if [ ! -f "$generated" ]; then
        echo -e "${RED}❌ Generated file not found: $generated${NC}"
        FAILED=$((FAILED + 1))
        return
    fi
    
    if diff -w "$expected" "$generated" > /dev/null; then
        echo -e "${GREEN}✅ $name matches${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}❌ $name differs${NC}"
        diff -w "$expected" "$generated" | head -20
        FAILED=$((FAILED + 1))
    fi
}

validate_dir() {
    local expected_dir=$1
    local generated_dir=$2
    local name=$3
    
    if [ ! -d "$expected_dir" ]; then
        echo -e "${YELLOW}⚠️  Expected directory not found: $expected_dir${NC}"
        return
    fi
    
    if [ ! -d "$generated_dir" ]; then
        echo -e "${YELLOW}⚠️  Generated directory not found: $generated_dir${NC}"
        return
    fi
    
    for expected_file in "$expected_dir"/*; do
        if [ -f "$expected_file" ]; then
            local filename=$(basename "$expected_file")
            validate_file "$expected_file" "$generated_dir/$filename" "$name/$filename"
        fi
    done
}

echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  Validation Tool${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}\n"

# Validate Go
echo -e "${BLUE}Go Tests:${NC}"
validate_dir "go/tests" "../../../internal/test/testdata" "Go"

echo -e "\n${BLUE}Go README:${NC}"
validate_dir "go/readme" "../../../docs" "Go"

# Validate Node.js
echo -e "\n${BLUE}Node.js Tests:${NC}"
validate_dir "nodejs/tests" "../../../pkg/test/testdata" "Node.js"

# Validate Python
echo -e "\n${BLUE}Python Tests:${NC}"
validate_dir "python/tests" "../../../pkg/test/testdata" "Python"

# Validate Edge Cases
echo -e "\n${BLUE}Edge Cases:${NC}"
validate_file "edge-cases/edge_cases_test.go" "../../../internal/test/testdata/edge_cases_test.go" "Edge Cases (Go)"
validate_file "edge-cases/edge_cases.test.js" "../../../pkg/test/testdata/edge_cases.test.js" "Edge Cases (Node.js)"
validate_file "edge-cases/test_edge_cases.py" "../../../pkg/test/testdata/test_edge_cases.py" "Edge Cases (Python)"

echo -e "\n${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  Summary${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "Total:   $TOTAL"
echo -e "${GREEN}Passed:  $PASSED${NC}"
echo -e "${RED}Failed:  $FAILED${NC}"
echo -e "${YELLOW}Skipped: $SKIPPED${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"

[ $FAILED -eq 0 ] || exit 1
EOF

chmod +x testdata/expected-outputs/validate.sh
print_success "Validation script created"

# ============================================================================
# Create Sample Code - Go Simple
# ============================================================================

print_header "Creating Go Simple Examples"

cat > testdata/sample-code/go/simple/hello.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

func greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
EOF

cat > testdata/sample-code/go/simple/math.go << 'EOF'
package math

func Add(a, b int) int {
    return a + b
}

func Subtract(a, b int) int {
    return a - b
}

func Multiply(a, b int) int {
    return a * b
}

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
EOF

cat > testdata/sample-code/go/simple/strings.go << 'EOF'
package strings

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func ToUpper(s string) string {
    return strings.ToUpper(s)
}
EOF

print_success "Go simple examples created"

# ============================================================================
# Create Sample Code - Node.js Simple
# ============================================================================

print_header "Creating Node.js Simple Examples"

cat > testdata/sample-code/nodejs/simple/math.js << 'EOF'
function add(a, b) {
    return a + b;
}

function subtract(a, b) {
    return a - b;
}

function multiply(a, b) {
    return a * b;
}

function divide(a, b) {
    if (b === 0) {
        throw new Error('Division by zero');
    }
    return a / b;
}

module.exports = { add, subtract, multiply, divide };
EOF

cat > testdata/sample-code/nodejs/simple/strings.js << 'EOF'
function reverse(str) {
    return str.split('').reverse().join('');
}

function toUpper(str) {
    return str.toUpperCase();
}

function toLower(str) {
    return str.toLowerCase();
}

module.exports = { reverse, toUpper, toLower };
EOF

print_success "Node.js simple examples created"

# ============================================================================
# Create Sample Code - Python Simple
# ============================================================================

print_header "Creating Python Simple Examples"

cat > testdata/sample-code/python/simple/math.py << 'EOF'
def add(a: int, b: int) -> int:
    return a + b

def subtract(a: int, b: int) -> int:
    return a - b

def multiply(a: int, b: int) -> int:
    return a * b

def divide(a: int, b: int) -> float:
    if b == 0:
        raise ValueError("Division by zero")
    return a / b
EOF

cat > testdata/sample-code/python/simple/strings.py << 'EOF'
def reverse(s: str) -> str:
    return s[::-1]

def to_upper(s: str) -> str:
    return s.upper()

def to_lower(s: str) -> str:
    return s.lower()
EOF

print_success "Python simple examples created"

# ============================================================================
# Create Benchmark Files
# ============================================================================

print_header "Creating Benchmark Files"

touch testdata/benchmark/go/benchmark_test.go
touch testdata/benchmark/nodejs/benchmark.test.js
touch testdata/benchmark/python/test_benchmark.py

print_success "Benchmark files created"

# ============================================================================
# Create Integration Test Files
# ============================================================================

print_header "Creating Integration Test Files"

touch testdata/integration/go/integration_test.go
touch testdata/integration/nodejs/integration.test.js
touch testdata/integration/python/test_integration.py

print_success "Integration test files created"

# ============================================================================
# Create Fuzz Test Files
# ============================================================================

print_header "Creating Fuzz Test Files"

cat > testdata/fuzz/go/fuzz_config_test.go << 'EOF'
//go:build go1.18
// +build go1.18

package fuzz

import (
    "testing"
    
    "github.com/yourusername/spoke-tool/internal/config"
)

func FuzzConfigParsing(f *testing.F) {
    seeds := []string{
        "",
        "test_spoke:\n  enabled: true",
        "invalid: yaml: [",
    }
    
    for _, seed := range seeds {
        f.Add(seed)
    }
    
    f.Fuzz(func(t *testing.T, data string) {
        cfg, err := config.Parse([]byte(data))
        _ = cfg
        _ = err
    })
}
EOF

touch testdata/fuzz/go/testdata/fuzz/FuzzConfigParsing
touch testdata/fuzz/nodejs/fuzz.test.js
touch testdata/fuzz/python/test_fuzz.py

print_success "Fuzz test files created"

# ============================================================================
# Create Input Files
# ============================================================================

print_header "Creating Input Files"

echo "Sample input for Go" > testdata/input/go/sample_input.txt
echo '{"key": "value"}' > testdata/input/nodejs/sample_input.json
echo "name,age,city\nJohn,30,New York\nJane,25,Boston" > testdata/input/python/sample_input.csv

print_success "Input files created"

# ============================================================================
# Summary
# ============================================================================

print_header "Setup Complete"

echo -e "${GREEN}✅ Testdata structure created successfully!${NC}"
echo ""
echo -e "${CYAN}Directory structure:${NC}"
find testdata -type d -not -path "*/\.*" | sort | sed 's/^testdata/  📂 testdata/'

echo ""
echo -e "${CYAN}File count:${NC} $(find testdata -type f | wc -l) files created"
echo ""
echo -e "${CYAN}Next steps:${NC}"
echo "  1. cd testdata/expected-outputs"
echo "  2. ./validate.sh  (when you have generated outputs to validate)"
echo "  3. Add more sample code as needed"
echo ""
echo -e "${GREEN}Happy testing! 🚀${NC}"