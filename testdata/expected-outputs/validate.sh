#!/usr/bin/env bash

# validate.sh - Validate generated outputs against expected outputs
# This script compares generated test files and documentation against
# expected outputs to ensure the generators are working correctly.

set -e

# ============================================================================
# Configuration
# ============================================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
EXPECTED_DIR="$SCRIPT_DIR"
GENERATED_DIR="$PROJECT_ROOT/generated"

# Default options
VERBOSE=false
UPDATE=false
LANGUAGE="all"
TYPE="all"
DIFF_TOOL="diff"
USE_COLORS=true

# Counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# ============================================================================
# Helper Functions
# ============================================================================

print_header() {
    echo -e "\n${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}\n"
}

print_success() {
    if [ "$USE_COLORS" = true ]; then
        echo -e "  ${GREEN}✅ $1${NC}"
    else
        echo "  [PASS] $1"
    fi
}

print_failure() {
    if [ "$USE_COLORS" = true ]; then
        echo -e "  ${RED}❌ $1${NC}"
    else
        echo "  [FAIL] $1"
    fi
}

print_warning() {
    if [ "$USE_COLORS" = true ]; then
        echo -e "  ${YELLOW}⚠️  $1${NC}"
    else
        echo "  [WARN] $1"
    fi
}

print_info() {
    if [ "$USE_COLORS" = true ]; then
        echo -e "  ${CYAN}ℹ️  $1${NC}"
    else
        echo "  [INFO] $1"
    fi
}

print_debug() {
    if [ "$VERBOSE" = true ]; then
        if [ "$USE_COLORS" = true ]; then
            echo -e "  ${MAGENTA}🔍 $1${NC}"
        else
            echo "  [DEBUG] $1"
        fi
    fi
}

print_summary() {
    echo -e "\n${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  Validation Summary${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "  Total:  $TOTAL_TESTS"
    echo -e "  ${GREEN}Passed: $PASSED_TESTS${NC}"
    echo -e "  ${RED}Failed: $FAILED_TESTS${NC}"
    echo -e "  ${YELLOW}Skipped: $SKIPPED_TESTS${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "${GREEN}  ✅ All validations passed!${NC}"
        return 0
    else
        echo -e "${RED}  ❌ Some validations failed!${NC}"
        return 1
    fi
}

# ============================================================================
# Validation Functions
# ============================================================================

validate_file() {
    local expected="$1"
    local generated="$2"
    local description="$3"
    local ignore_pattern="$4"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    print_debug "Validating: $description"
    print_debug "  Expected: $expected"
    print_debug "  Generated: $generated"
    
    # Check if expected file exists
    if [ ! -f "$expected" ]; then
        print_failure "$description - Expected file not found: $expected"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
    
    # Check if generated file exists
    if [ ! -f "$generated" ]; then
        print_failure "$description - Generated file not found: $generated"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
    
    # Compare files
    if [ -n "$ignore_pattern" ]; then
        # Use grep to filter out lines matching ignore pattern
        diff -u \
            <(grep -v "$ignore_pattern" "$expected" 2>/dev/null || true) \
            <(grep -v "$ignore_pattern" "$generated" 2>/dev/null || true) \
            > /tmp/diff_output.txt
    else
        diff -u "$expected" "$generated" > /tmp/diff_output.txt
    fi
    
    if [ $? -eq 0 ]; then
        print_success "$description"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        rm -f /tmp/diff_output.txt
        return 0
    else
        print_failure "$description - Files differ"
        if [ "$VERBOSE" = true ]; then
            echo -e "${YELLOW}Diff output:${NC}"
            cat /tmp/diff_output.txt | head -20
            if [ $(wc -l < /tmp/diff_output.txt) -gt 20 ]; then
                echo -e "${YELLOW}  ... (truncated)${NC}"
            fi
        fi
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

validate_directory() {
    local expected_dir="$1"
    local generated_dir="$2"
    local description="$3"
    
    print_debug "Validating directory: $description"
    
    # Check if expected directory exists
    if [ ! -d "$expected_dir" ]; then
        print_warning "Expected directory not found: $expected_dir"
        SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
        return 1
    fi
    
    # Check if generated directory exists
    if [ ! -d "$generated_dir" ]; then
        print_warning "Generated directory not found: $generated_dir"
        SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
        return 1
    fi
    
    # Find all expected files
    while IFS= read -r expected_file; do
        if [ -f "$expected_file" ]; then
            relative_path="${expected_file#$expected_dir/}"
            generated_file="$generated_dir/$relative_path"
            validate_file "$expected_file" "$generated_file" "$description/$relative_path"
        fi
    done < <(find "$expected_dir" -type f | sort)
}

# ============================================================================
# Language-Specific Validation
# ============================================================================

validate_go() {
    print_header "Validating Go"
    
    # Go tests
    if [[ "$TYPE" == "all" || "$TYPE" == "tests" ]]; then
        validate_directory \
            "$EXPECTED_DIR/go/tests" \
            "$GENERATED_DIR/go/tests" \
            "Go Tests"
    fi
    
    # Go README
    if [[ "$TYPE" == "all" || "$TYPE" == "readme" ]]; then
        validate_directory \
            "$EXPECTED_DIR/go/readme" \
            "$GENERATED_DIR/go/readme" \
            "Go README"
    fi
}

validate_nodejs() {
    print_header "Validating Node.js"
    
    # Node.js tests
    if [[ "$TYPE" == "all" || "$TYPE" == "tests" ]]; then
        validate_directory \
            "$EXPECTED_DIR/nodejs/tests" \
            "$GENERATED_DIR/nodejs/tests" \
            "Node.js Tests"
    fi
    
    # Node.js README
    if [[ "$TYPE" == "all" || "$TYPE" == "readme" ]]; then
        validate_directory \
            "$EXPECTED_DIR/nodejs/readme" \
            "$GENERATED_DIR/nodejs/readme" \
            "Node.js README"
    fi
}

validate_python() {
    print_header "Validating Python"
    
    # Python tests
    if [[ "$TYPE" == "all" || "$TYPE" == "tests" ]]; then
        validate_directory \
            "$EXPECTED_DIR/python/tests" \
            "$GENERATED_DIR/python/tests" \
            "Python Tests"
    fi
    
    # Python README
    if [[ "$TYPE" == "all" || "$TYPE" == "readme" ]]; then
        validate_directory \
            "$EXPECTED_DIR/python/readme" \
            "$GENERATED_DIR/python/readme" \
            "Python README"
    fi
}

validate_edge_cases() {
    print_header "Validating Edge Cases"
    
    # Edge cases tests
    if [[ "$TYPE" == "all" || "$TYPE" == "edge" ]]; then
        validate_file \
            "$EXPECTED_DIR/edge-cases/edge_cases_test.go" \
            "$GENERATED_DIR/edge-cases/edge_cases_test.go" \
            "Edge Cases (Go)"
        
        validate_file \
            "$EXPECTED_DIR/edge-cases/edge_cases.test.js" \
            "$GENERATED_DIR/edge-cases/edge_cases.test.js" \
            "Edge Cases (Node.js)"
        
        validate_file \
            "$EXPECTED_DIR/edge-cases/test_edge_cases.py" \
            "$GENERATED_DIR/edge-cases/test_edge_cases.py" \
            "Edge Cases (Python)"
        
        validate_file \
            "$EXPECTED_DIR/edge-cases/README.md" \
            "$GENERATED_DIR/edge-cases/README.md" \
            "Edge Cases README"
    fi
}

validate_mixed() {
    print_header "Validating Mixed Project"
    
    if [[ "$TYPE" == "all" || "$TYPE" == "mixed" ]]; then
        validate_file \
            "$EXPECTED_DIR/mixed/README.md" \
            "$GENERATED_DIR/mixed/README.md" \
            "Mixed Project README"
        
        validate_file \
            "$EXPECTED_DIR/mixed/docker-compose.yml" \
            "$GENERATED_DIR/mixed/docker-compose.yml" \
            "Mixed Project Docker Compose"
        
        validate_file \
            "$EXPECTED_DIR/mixed/Makefile" \
            "$GENERATED_DIR/mixed/Makefile" \
            "Mixed Project Makefile"
    fi
}

validate_benchmarks() {
    print_header "Validating Benchmarks"
    
    if [[ "$TYPE" == "all" || "$TYPE" == "benchmark" ]]; then
        validate_file \
            "$PROJECT_ROOT/testdata/benchmark/go/benchmark_test.go" \
            "$GENERATED_DIR/benchmark/go/benchmark_test.go" \
            "Go Benchmarks"
        
        validate_file \
            "$PROJECT_ROOT/testdata/benchmark/nodejs/benchmark.test.js" \
            "$GENERATED_DIR/benchmark/nodejs/benchmark.test.js" \
            "Node.js Benchmarks"
        
        validate_file \
            "$PROJECT_ROOT/testdata/benchmark/python/test_benchmark.py" \
            "$GENERATED_DIR/benchmark/python/test_benchmark.py" \
            "Python Benchmarks"
    fi
}

validate_integration() {
    print_header "Validating Integration Tests"
    
    if [[ "$TYPE" == "all" || "$TYPE" == "integration" ]]; then
        validate_file \
            "$PROJECT_ROOT/testdata/integration/go/integration_test.go" \
            "$GENERATED_DIR/integration/go/integration_test.go" \
            "Go Integration Tests"
        
        validate_file \
            "$PROJECT_ROOT/testdata/integration/nodejs/integration.test.js" \
            "$GENERATED_DIR/integration/nodejs/integration.test.js" \
            "Node.js Integration Tests"
        
        validate_file \
            "$PROJECT_ROOT/testdata/integration/python/test_integration.py" \
            "$GENERATED_DIR/integration/python/test_integration.py" \
            "Python Integration Tests"
    fi
}

validate_fuzz() {
    print_header "Validating Fuzz Tests"
    
    if [[ "$TYPE" == "all" || "$TYPE" == "fuzz" ]]; then
        validate_file \
            "$PROJECT_ROOT/testdata/fuzz/go/fuzz_test.go" \
            "$GENERATED_DIR/fuzz/go/fuzz_test.go" \
            "Go Fuzz Tests"
        
        validate_file \
            "$PROJECT_ROOT/testdata/fuzz/go/testdata/fuzz/FuzzAdd" \
            "$GENERATED_DIR/fuzz/go/testdata/fuzz/FuzzAdd" \
            "Go Fuzz Corpus"
        
        validate_file \
            "$PROJECT_ROOT/testdata/fuzz/nodejs/fuzz.test.js" \
            "$GENERATED_DIR/fuzz/nodejs/fuzz.test.js" \
            "Node.js Fuzz Tests"
        
        validate_file \
            "$PROJECT_ROOT/testdata/fuzz/python/test_fuzz.py" \
            "$GENERATED_DIR/fuzz/python/test_fuzz.py" \
            "Python Fuzz Tests"
    fi
}

# ============================================================================
# Main Validation
# ============================================================================

main() {
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  -h, --help           Show this help message"
                echo "  -v, --verbose        Verbose output"
                echo "  -u, --update         Update expected outputs (NOT RECOMMENDED)"
                echo "  -l, --language LANG   Language to validate (go, nodejs, python, all)"
                echo "  -t, --type TYPE       Type to validate (tests, readme, edge, mixed, benchmark, integration, fuzz, all)"
                echo "  -d, --dir DIR         Generated files directory (default: ./generated)"
                echo "  --no-color            Disable colors"
                echo ""
                echo "Examples:"
                echo "  $0                            # Validate all"
                echo "  $0 -l go                       # Validate only Go"
                echo "  $0 -t tests                    # Validate only tests"
                echo "  $0 -l nodejs -t readme         # Validate Node.js README only"
                echo "  $0 -v                          # Verbose output"
                exit 0
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -u|--update)
                UPDATE=true
                print_warning "Update mode enabled - will overwrite expected outputs"
                shift
                ;;
            -l|--language)
                LANGUAGE="$2"
                shift 2
                ;;
            -t|--type)
                TYPE="$2"
                shift 2
                ;;
            -d|--dir)
                GENERATED_DIR="$2"
                shift 2
                ;;
            --no-color)
                USE_COLORS=false
                shift
                ;;
            *)
                echo -e "${RED}Unknown option: $1${NC}"
                exit 1
                ;;
        esac
    done
    
    print_header "Output Validation Tool"
    print_info "Expected directory: $EXPECTED_DIR"
    print_info "Generated directory: $GENERATED_DIR"
    print_info "Language: $LANGUAGE"
    print_info "Type: $TYPE"
    echo ""
    
    # Validate based on language
    case $LANGUAGE in
        go)
            validate_go
            ;;
        nodejs)
            validate_nodejs
            ;;
        python)
            validate_python
            ;;
        all)
            validate_go
            validate_nodejs
            validate_python
            validate_edge_cases
            validate_mixed
            validate_benchmarks
            validate_integration
            validate_fuzz
            ;;
        *)
            echo -e "${RED}Invalid language: $LANGUAGE${NC}"
            exit 1
            ;;
    esac
    
    # Print summary
    print_summary
    
    # Exit with appropriate code
    if [ $FAILED_TESTS -eq 0 ]; then
        exit 0
    else
        exit 1
    fi
}

# Run main function
main "$@"