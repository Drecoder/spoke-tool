#!/usr/bin/env bash

# run-tests.sh - Run tests for spoke-tool
# This script runs unit tests, integration tests, and coverage reports

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Default values
VERBOSE=${VERBOSE:-false}
COVERAGE=${COVERAGE:-false}
INTEGRATION=${INTEGRATION:-false}
BENCHMARK=${BENCHMARK:-false}
RACE=${RACE:-false}
PARALLEL=${PARALLEL:-4}
TIMEOUT=${TIMEOUT:-"5m"}
OUTPUT_DIR=${OUTPUT_DIR:-"test-results"}
PACKAGES=${PACKAGES:-"./..."}
LANGUAGE=${LANGUAGE:-"go"}  # go, nodejs, python, all
FAIL_FAST=${FAIL_FAST:-false}

# Print banner
print_banner() {
    echo -e "${BLUE}"
    echo '  _____           _        _    _           _        '
    echo ' |_   _|         | |      | |  | |         | |       '
    echo '   | | ___  _ __ | |_ __ _| |__| |_ __  ___| |_ __ _ '
    echo '   | |/ _ \| '"'"'_ \| __/ _` |  __  | '"'"'_ \/ __| __/ _` |'
    echo '   | | (_) | | | | || (_| | |  | | |_) \__ \ || (_| |'
    echo '   \_/\___/|_| |_|\__\__,_|_|  |_| .__/|___/\__\__,_|'
    echo '                                  | |                  '
    echo '                                  |_|                  '
    echo -e "${NC}"
    echo -e "${CYAN}Test Runner${NC}"
    echo ""
}

# Print help
print_help() {
    echo -e "${YELLOW}Usage:${NC} ./scripts/run-tests.sh [options]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -v, --verbose       Verbose output"
    echo "  -c, --coverage      Run with coverage"
    echo "  -i, --integration   Run integration tests"
    echo "  -b, --benchmark     Run benchmarks"
    echo "  -r, --race          Enable race detection (Go only)"
    echo "  -p, --parallel N    Set parallel test count (default: 4)"
    echo "  -t, --timeout T     Set test timeout (default: 5m)"
    echo "  -o, --output DIR    Output directory for results (default: test-results)"
    echo "  -l, --language LANG Language to test (go, nodejs, python, all)"
    echo "  -f, --fail-fast     Stop on first failure"
    echo "  --packages PKGS     Packages to test (default: ./...)"
    echo ""
    echo "Examples:"
    echo "  ./scripts/run-tests.sh                    # Run all Go tests"
    echo "  ./scripts/run-tests.sh -v -c              # Verbose with coverage"
    echo "  ./scripts/run-tests.sh -l all             # Test all languages"
    echo "  ./scripts/run-tests.sh -i -b              # Integration tests and benchmarks"
    echo "  ./scripts/run-tests.sh -l nodejs -c       # Node.js tests with coverage"
}

# Check dependencies
check_deps() {
    local lang=$1
    local missing=false
    
    echo -e "${YELLOW}Checking dependencies for $lang...${NC}"
    
    case $lang in
        go)
            if ! command -v go &> /dev/null; then
                echo -e "${RED}❌ go not found${NC}"
                missing=true
            else
                echo -e "${GREEN}✅ go $(go version | awk '{print $3}')${NC}"
            fi
            
            if [[ "$COVERAGE" == "true" ]]; then
                if ! command -v go &> /dev/null; then
                    echo -e "${RED}❌ go tool cover required for coverage${NC}"
                    missing=true
                fi
            fi
            ;;
            
        nodejs)
            if ! command -v node &> /dev/null; then
                echo -e "${RED}❌ node not found${NC}"
                missing=true
            else
                echo -e "${GREEN}✅ node $(node --version)${NC}"
            fi
            
            if ! command -v npm &> /dev/null; then
                echo -e "${RED}❌ npm not found${NC}"
                missing=true
            else
                echo -e "${GREEN}✅ npm $(npm --version)${NC}"
            fi
            
            # Check for test frameworks
            if [[ -f "package.json" ]]; then
                if grep -q '"jest"' package.json; then
                    echo -e "${GREEN}✅ jest found${NC}"
                else
                    echo -e "${YELLOW}⚠️  jest not found in package.json${NC}"
                fi
            fi
            ;;
            
        python)
            if ! command -v python3 &> /dev/null; then
                echo -e "${RED}❌ python3 not found${NC}"
                missing=true
            else
                echo -e "${GREEN}✅ python3 $(python3 --version | cut -d' ' -f2)${NC}"
            fi
            
            # Check for pytest
            if ! python3 -c "import pytest" 2>/dev/null; then
                echo -e "${YELLOW}⚠️  pytest not found${NC}"
                echo -e "${YELLOW}  Run: pip install pytest pytest-cov${NC}"
            else
                echo -e "${GREEN}✅ pytest found${NC}"
            fi
            ;;
    esac
    
    if [[ "$missing" == "true" ]]; then
        echo -e "${RED}Missing dependencies. Please install them and try again.${NC}"
        exit 1
    fi
    
    echo ""
}

# Run Go tests
run_go_tests() {
    echo -e "${CYAN}Running Go tests...${NC}"
    
    local start_time=$(date +%s%N)
    local exit_code=0
    local output_file="$OUTPUT_DIR/go-tests.out"
    local coverage_file="$OUTPUT_DIR/go-coverage.out"
    local json_file="$OUTPUT_DIR/go-tests.json"
    local junit_file="$OUTPUT_DIR/go-tests.xml"
    
    mkdir -p "$OUTPUT_DIR"
    
    # Build test command
    local cmd="go test"
    
    # Add flags
    [[ "$VERBOSE" == "true" ]] && cmd="$cmd -v"
    [[ "$RACE" == "true" ]] && cmd="$cmd -race"
    [[ "$COVERAGE" == "true" ]] && cmd="$cmd -coverprofile=$coverage_file -covermode=atomic"
    [[ "$INTEGRATION" == "true" ]] && cmd="$cmd -tags=integration"
    [[ "$BENCHMARK" == "true" ]] && cmd="$cmd -bench=. -benchmem"
    [[ "$FAIL_FAST" == "true" ]] && cmd="$cmd -failfast"
    
    cmd="$cmd -parallel $PARALLEL -timeout $TIMEOUT"
    
    # Add packages
    if [[ "$BENCHMARK" == "true" ]]; then
        cmd="$cmd -run=^$ $PACKAGES"  # Run only benchmarks
    else
        cmd="$cmd $PACKAGES"
    fi
    
    echo -e "${YELLOW}Running: $cmd${NC}"
    echo ""
    
    # Run tests with tee to capture output
    if [[ "$VERBOSE" == "true" ]]; then
        eval "$cmd" 2>&1 | tee "$output_file"
        exit_code=${PIPESTATUS[0]}
    else
        eval "$cmd" &> "$output_file"
        exit_code=$?
        cat "$output_file"
    fi
    
    local end_time=$(date +%s%N)
    local duration=$(( ($end_time - $start_time) / 1000000 ))
    
    echo ""
    
    # Generate coverage report
    if [[ "$COVERAGE" == "true" && -f "$coverage_file" ]]; then
        echo -e "${YELLOW}Coverage summary:${NC}"
        go tool cover -func="$coverage_file" | grep -E "^total:" | awk '{print $3}'
        
        # Generate HTML coverage report
        go tool cover -html="$coverage_file" -o "$OUTPUT_DIR/go-coverage.html"
        echo -e "${GREEN}✅ Coverage HTML report: $OUTPUT_DIR/go-coverage.html${NC}"
    fi
    
    # Generate JSON output
    if command -v go-test2json &> /dev/null; then
        cat "$output_file" | go-test2json > "$json_file"
        echo -e "${GREEN}✅ JSON output: $json_file${NC}"
    fi
    
    # Generate JUnit XML (if go-junit-report is installed)
    if command -v go-junit-report &> /dev/null; then
        cat "$output_file" | go-junit-report > "$junit_file"
        echo -e "${GREEN}✅ JUnit XML: $junit_file${NC}"
    fi
    
    echo -e "${GREEN}✅ Go tests completed in ${duration}ms${NC}"
    
    return $exit_code
}

# Run Node.js tests
run_nodejs_tests() {
    echo -e "${CYAN}Running Node.js tests...${NC}"
    
    local start_time=$(date +%s%N)
    local exit_code=0
    local output_dir="$OUTPUT_DIR/nodejs"
    local output_file="$output_dir/tests.out"
    local json_file="$output_dir/jest-results.json"
    local junit_file="$output_dir/junit.xml"
    local coverage_dir="$output_dir/coverage"
    
    mkdir -p "$output_dir" "$coverage_dir"
    
    # Determine test runner
    local runner="npm test"
    if [[ -f "node_modules/.bin/jest" ]]; then
        runner="npx jest"
    elif [[ -f "node_modules/.bin/mocha" ]]; then
        runner="npx mocha"
    fi
    
    # Build test command
    local cmd="$runner"
    
    # Add flags
    [[ "$VERBOSE" == "true" ]] && cmd="$cmd --verbose"
    [[ "$COVERAGE" == "true" ]] && cmd="$cmd --coverage --coverageDirectory=$coverage_dir"
    [[ "$INTEGRATION" == "true" ]] && cmd="$cmd --testMatch='**/*.integration.js'"
    [[ "$BENCHMARK" == "true" ]] && cmd="$cmd --testMatch='**/*.bench.js'"
    
    # Add Jest specific flags
    if [[ "$runner" == *"jest"* ]]; then
        cmd="$cmd --json --outputFile=$json_file"
        [[ "$FAIL_FAST" == "true" ]] && cmd="$cmd --bail"
        cmd="$cmd --maxWorkers=$PARALLEL"
    fi
    
    echo -e "${YELLOW}Running: $cmd${NC}"
    echo ""
    
    # Run tests
    if [[ "$VERBOSE" == "true" ]]; then
        eval "$cmd" 2>&1 | tee "$output_file"
        exit_code=${PIPESTATUS[0]}
    else
        eval "$cmd" &> "$output_file"
        exit_code=$?
        cat "$output_file"
    fi
    
    local end_time=$(date +%s%N)
    local duration=$(( ($end_time - $start_time) / 1000000 ))
    
    echo ""
    
    # Generate JUnit XML from Jest JSON
    if [[ -f "$json_file" && "$runner" == *"jest"* ]]; then
        # Convert Jest JSON to JUnit XML
        if command -v npx &> /dev/null; then
            npx jest-junit --output="$junit_file" --outputDirectory="$output_dir" 2>/dev/null
        fi
    fi
    
    # Show coverage summary
    if [[ "$COVERAGE" == "true" && -f "$coverage_dir/lcov-report/index.html" ]]; then
        echo -e "${GREEN}✅ Coverage HTML report: $coverage_dir/lcov-report/index.html${NC}"
        
        # Show coverage summary
        if [[ -f "$coverage_dir/coverage-final.json" ]]; then
            echo -e "${YELLOW}Coverage summary:${NC}"
            grep -o '"pct":[^,]*' "$coverage_dir/coverage-final.json" | head -1
        fi
    fi
    
    echo -e "${GREEN}✅ Node.js tests completed in ${duration}ms${NC}"
    
    return $exit_code
}

# Run Python tests
run_python_tests() {
    echo -e "${CYAN}Running Python tests...${NC}"
    
    local start_time=$(date +%s%N)
    local exit_code=0
    local output_dir="$OUTPUT_DIR/python"
    local output_file="$output_dir/tests.out"
    local xml_file="$output_dir/pytest.xml"
    local coverage_dir="$output_dir/coverage"
    local html_dir="$output_dir/htmlcov"
    
    mkdir -p "$output_dir" "$coverage_dir" "$html_dir"
    
    # Build test command
    local cmd="python3 -m pytest"
    
    # Add flags
    [[ "$VERBOSE" == "true" ]] && cmd="$cmd -v"
    [[ "$COVERAGE" == "true" ]] && cmd="$cmd --cov=. --cov-report=html:$html_dir --cov-report=xml:$coverage_dir/coverage.xml"
    [[ "$INTEGRATION" == "true" ]] && cmd="$cmd -m integration"
    [[ "$BENCHMARK" == "true" ]] && cmd="$cmd --benchmark-only"
    [[ "$FAIL_FAST" == "true" ]] && cmd="$cmd -x"
    
    cmd="$cmd -n $PARALLEL --timeout=$(echo $TIMEOUT | sed 's/m//g') --junitxml=$xml_file"
    
    echo -e "${YELLOW}Running: $cmd${NC}"
    echo ""
    
    # Run tests
    if [[ "$VERBOSE" == "true" ]]; then
        eval "$cmd" 2>&1 | tee "$output_file"
        exit_code=${PIPESTATUS[0]}
    else
        eval "$cmd" &> "$output_file"
        exit_code=$?
        cat "$output_file"
    fi
    
    local end_time=$(date +%s%N)
    local duration=$(( ($end_time - $start_time) / 1000000 ))
    
    echo ""
    
    # Show coverage summary
    if [[ "$COVERAGE" == "true" && -f "$html_dir/index.html" ]]; then
        echo -e "${GREEN}✅ Coverage HTML report: $html_dir/index.html${NC}"
        
        # Extract coverage percentage
        if [[ -f "$coverage_dir/coverage.xml" ]]; then
            local coverage=$(grep -o 'line-rate="[^"]*"' "$coverage_dir/coverage.xml" | head -1 | cut -d'"' -f2)
            coverage=$(echo "$coverage * 100" | bc)
            echo -e "${YELLOW}Coverage: ${coverage}%${NC}"
        fi
    fi
    
    echo -e "${GREEN}✅ Python tests completed in ${duration}ms${NC}"
    
    return $exit_code
}

# Generate summary report
generate_summary() {
    local results=("$@")
    local summary_file="$OUTPUT_DIR/summary.md"
    
    echo -e "${YELLOW}Generating test summary...${NC}"
    
    cat > "$summary_file" << EOF
# Test Summary Report

Generated: $(date)

## Overview

| Language | Status | Duration | Coverage |
|----------|--------|----------|----------|
EOF
    
    for result in "${results[@]}"; do
        echo "$result" >> "$summary_file"
    done
    
    cat >> "$summary_file" << EOF

## Details

### Go Tests
- Package: $PACKAGES
- Parallel: $PARALLEL
- Timeout: $TIMEOUT

### Node.js Tests
- Test files: $(find . -name "*.test.js" -o -name "*.spec.js" 2>/dev/null | wc -l)
- Coverage: $([[ "$COVERAGE" == "true" ]] && echo "Yes" || echo "No")

### Python Tests
- Test files: $(find . -name "test_*.py" -o -name "*_test.py" 2>/dev/null | wc -l)
- Integration: $([[ "$INTEGRATION" == "true" ]] && echo "Yes" || echo "No")

## Command
\`\`\`bash
$0 $@
\`\`\`
EOF
    
    echo -e "${GREEN}✅ Summary report: $summary_file${NC}"
}

# Clean up
cleanup() {
    echo -e "${YELLOW}Cleaning up...${NC}"
    # Remove temporary files if needed
}

# Main function
main() {
    local exit_code=0
    local results=()
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                print_help
                exit 0
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -c|--coverage)
                COVERAGE=true
                shift
                ;;
            -i|--integration)
                INTEGRATION=true
                shift
                ;;
            -b|--benchmark)
                BENCHMARK=true
                shift
                ;;
            -r|--race)
                RACE=true
                shift
                ;;
            -p|--parallel)
                PARALLEL="$2"
                shift 2
                ;;
            -t|--timeout)
                TIMEOUT="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -l|--language)
                LANGUAGE="$2"
                shift 2
                ;;
            -f|--fail-fast)
                FAIL_FAST=true
                shift
                ;;
            --packages)
                PACKAGES="$2"
                shift 2
                ;;
            *)
                echo -e "${RED}Unknown option: $1${NC}"
                print_help
                exit 1
                ;;
        esac
    done
    
    # Print banner
    print_banner
    
    # Set trap for cleanup
    trap cleanup EXIT
    
    # Create output directory
    mkdir -p "$OUTPUT_DIR"
    
    # Save start time
    global_start=$(date +%s%N)
    
    # Run tests based on language
    case $LANGUAGE in
        go)
            check_deps "go"
            run_go_tests
            exit_code=$?
            results+=("| Go | $([[ $exit_code -eq 0 ]] && echo "✅ Passed" || echo "❌ Failed") | - | $([[ "$COVERAGE" == "true" ]] && echo "Yes" || echo "No") |")
            ;;
            
        nodejs)
            check_deps "nodejs"
            run_nodejs_tests
            exit_code=$?
            results+=("| Node.js | $([[ $exit_code -eq 0 ]] && echo "✅ Passed" || echo "❌ Failed") | - | $([[ "$COVERAGE" == "true" ]] && echo "Yes" || echo "No") |")
            ;;
            
        python)
            check_deps "python"
            run_python_tests
            exit_code=$?
            results+=("| Python | $([[ $exit_code -eq 0 ]] && echo "✅ Passed" || echo "❌ Failed") | - | $([[ "$COVERAGE" == "true" ]] && echo "Yes" || echo "No") |")
            ;;
            
        all)
            # Run Go tests
            check_deps "go"
            run_go_tests
            go_code=$?
            results+=("| Go | $([[ $go_code -eq 0 ]] && echo "✅ Passed" || echo "❌ Failed") | - | $([[ "$COVERAGE" == "true" ]] && echo "Yes" || echo "No") |")
            
            # Run Node.js tests
            if [[ -f "package.json" ]] || ls *.js &>/dev/null; then
                echo ""
                check_deps "nodejs"
                run_nodejs_tests
                node_code=$?
                results+=("| Node.js | $([[ $node_code -eq 0 ]] && echo "✅ Passed" || echo "❌ Failed") | - | $([[ "$COVERAGE" == "true" ]] && echo "Yes" || echo "No") |")
                exit_code=$((exit_code + node_code))
            fi
            
            # Run Python tests
            if ls *.py &>/dev/null; then
                echo ""
                check_deps "python"
                run_python_tests
                py_code=$?
                results+=("| Python | $([[ $py_code -eq 0 ]] && echo "✅ Passed" || echo "❌ Failed") | - | $([[ "$COVERAGE" == "true" ]] && echo "Yes" || echo "No") |")
                exit_code=$((exit_code + py_code))
            fi
            ;;
            
        *)
            echo -e "${RED}Unsupported language: $LANGUAGE${NC}"
            exit 1
            ;;
    esac
    
    # Calculate total duration
    global_end=$(date +%s%N)
    global_duration=$(( ($global_end - $global_start) / 1000000 ))
    
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════════${NC}"
    echo ""
    
    # Generate summary
    generate_summary "${results[@]}"
    
    # Final status
    if [[ $exit_code -eq 0 ]]; then
        echo -e "${GREEN}✅ All tests passed! (${global_duration}ms)${NC}"
    else
        echo -e "${RED}❌ Some tests failed. Check the output above.${NC}"
    fi
    
    echo -e "${CYAN}Results saved to: $OUTPUT_DIR/${NC}"
    
    exit $exit_code
}

# Run main function
main "$@"