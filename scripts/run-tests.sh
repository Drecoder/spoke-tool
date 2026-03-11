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
TARGET_DIR=${TARGET_DIR:-""}  # Target directory to test

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
    echo "  --target DIR        Target directory to test (default: current directory)"
}

# Check dependencies
check_deps() {
    local lang=$1
    local missing=false
    
    echo -e "${YELLOW}Checking dependencies for $lang...${NC}"
    
    case $lang in
        go)
            if ! command -v go &> /dev/null; then
                echo -e "${RED}❌ go not found${NC}"; missing=true
            else
                echo -e "${GREEN}✅ go $(go version | awk '{print $3}')${NC}"
            fi
            ;;
        nodejs)
            if ! command -v node &> /dev/null; then
                echo -e "${RED}❌ node not found${NC}"; missing=true
            else
                echo -e "${GREEN}✅ node $(node --version)${NC}"
            fi
            ;;
        python)
            if ! command -v python3 &> /dev/null; then
                echo -e "${RED}❌ python3 not found${NC}"; missing=true
            else
                echo -e "${GREEN}✅ python3 $(python3 --version | cut -d' ' -f2)${NC}"
            fi
            ;;
    esac
    
    if [[ "$missing" == "true" ]]; then
        echo -e "${RED}Missing dependencies for $lang. Exiting.${NC}"
        exit 1
    fi
}

# Run Go tests
run_go_tests() {
    echo -e "${CYAN}Running Go tests...${NC}"
    local start_time=$(date +%s%N)
    local exit_code=0
    local output_file="$OUTPUT_DIR/go-tests.out"
    local coverage_file="$OUTPUT_DIR/go-coverage.out"
    
    mkdir -p "$OUTPUT_DIR"
    local cmd="go test"
    [[ "$VERBOSE" == "true" ]] && cmd="$cmd -v"
    [[ "$RACE" == "true" ]] && cmd="$cmd -race"
    [[ "$COVERAGE" == "true" ]] && cmd="$cmd -coverprofile=$coverage_file -covermode=atomic"
    [[ "$FAIL_FAST" == "true" ]] && cmd="$cmd -failfast"
    
    cmd="$cmd -parallel $PARALLEL -timeout $TIMEOUT $PACKAGES"
    
    echo -e "${YELLOW}Running: $cmd${NC}"
    eval "$cmd" 2>&1 | tee "$output_file"
    exit_code=${PIPESTATUS[0]}
    
    if [[ "$COVERAGE" == "true" && -f "$coverage_file" ]]; then
        go tool cover -html="$coverage_file" -o "$OUTPUT_DIR/go-coverage.html"
    fi
    return $exit_code
}

# Run Node.js tests
run_nodejs_tests() {
    echo -e "${CYAN}Running Node.js tests...${NC}"
    local start_time=$(date +%s%N)
    local exit_code=0
    local output_dir="$OUTPUT_DIR/nodejs"
    mkdir -p "$output_dir"
    
    local runner="npx jest"
    local cmd="$runner --maxWorkers=$PARALLEL"
    [[ "$VERBOSE" == "true" ]] && cmd="$cmd --verbose"
    [[ "$COVERAGE" == "true" ]] && cmd="$cmd --coverage --coverageDirectory=$output_dir/coverage"
    
    echo -e "${YELLOW}Running: $cmd${NC}"
    eval "$cmd" 2>&1 | tee "$output_dir/tests.out"
    exit_code=${PIPESTATUS[0]}
    return $exit_code
}

# Run Python tests
run_python_tests() {
    echo -e "${CYAN}Running Python tests...${NC}"
    local exit_code=0
    local output_dir="$OUTPUT_DIR/python"
    mkdir -p "$output_dir"
    
    local cmd="python3 -m pytest -n $PARALLEL"
    [[ "$VERBOSE" == "true" ]] && cmd="$cmd -v"
    [[ "$COVERAGE" == "true" ]] && cmd="$cmd --cov=. --cov-report=html:$output_dir/htmlcov"
    
    echo -e "${YELLOW}Running: $cmd${NC}"
    eval "$cmd" 2>&1 | tee "$output_dir/tests.out"
    exit_code=${PIPESTATUS[0]}
    return $exit_code
}

generate_summary() {
    local results=("$@")
    echo -e "${YELLOW}Test Summary:${NC}"
    for res in "${results[@]}"; do
        echo -e "$res"
    done
}

cleanup() {
    : # Placeholder for cleanup logic
}

main() {
    local exit_code=0
    local results=()

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help) print_help; exit 0 ;;
            -v|--verbose) VERBOSE=true; shift ;;
            -c|--coverage) COVERAGE=true; shift ;;
            -i|--integration) INTEGRATION=true; shift ;;
            -b|--benchmark) BENCHMARK=true; shift ;;
            -r|--race) RACE=true; shift ;;
            -p|--parallel) PARALLEL="$2"; shift 2 ;;
            -t|--timeout) TIMEOUT="$2"; shift 2 ;;
            -o|--output) OUTPUT_DIR="$2"; shift 2 ;;
            -l|--language) LANGUAGE="$2"; shift 2 ;;
            -f|--fail-fast) FAIL_FAST=true; shift ;;
            --packages) PACKAGES="$2"; shift 2 ;;
            --target) TARGET_DIR="$2"; shift 2 ;;
            *) echo -e "${RED}Unknown option: $1${NC}"; print_help; exit 1 ;;
        esac
    done

    if [[ -n "$TARGET_DIR" ]]; then
        cd "$TARGET_DIR" || exit 1
    fi

    print_banner
    trap cleanup EXIT
    mkdir -p "$OUTPUT_DIR"

    case $LANGUAGE in
        go)
            check_deps "go"
            run_go_tests; exit_code=$?
            results+=("Go: $([[ $exit_code -eq 0 ]] && echo -e "${GREEN}PASS${NC}" || echo -e "${RED}FAIL${NC}")")
            ;;
        nodejs)
            check_deps "nodejs"
            run_nodejs_tests; exit_code=$?
            results+=("NodeJS: $([[ $exit_code -eq 0 ]] && echo -e "${GREEN}PASS${NC}" || echo -e "${RED}FAIL${NC}")")
            ;;
        python)
            check_deps "python"
            run_python_tests; exit_code=$?
            results+=("Python: $([[ $exit_code -eq 0 ]] && echo -e "${GREEN}PASS${NC}" || echo -e "${RED}FAIL${NC}")")
            ;;
        all)
            if [[ -f "go.mod" ]]; then
                check_deps "go"; run_go_tests; r=$?; exit_code=$((exit_code + r))
                results+=("Go: $([[ $r -eq 0 ]] && echo -e "${GREEN}PASS${NC}" || echo -e "${RED}FAIL${NC}")")
            fi
            if [[ -f "package.json" ]]; then
                check_deps "nodejs"; run_nodejs_tests; r=$?; exit_code=$((exit_code + r))
                results+=("NodeJS: $([[ $r -eq 0 ]] && echo -e "${GREEN}PASS${NC}" || echo -e "${RED}FAIL${NC}")")
            fi
            if ls *.py &>/dev/null || [ -d "tests" ]; then
                check_deps "python"; run_python_tests; r=$?; exit_code=$((exit_code + r))
                results+=("Python: $([[ $r -eq 0 ]] && echo -e "${GREEN}PASS${NC}" || echo -e "${RED}FAIL${NC}")")
            fi
            ;;
        *)
            echo -e "${RED}Unsupported language: $LANGUAGE${NC}"; exit 1 ;;
    esac

    generate_summary "${results[@]}"
    exit $exit_code
}

main "$@"