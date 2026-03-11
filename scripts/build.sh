#!/usr/bin/env bash

# build.sh - Build script for spoke-tool
# This script builds both readmegen and testgen binaries

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Version info
VERSION=${VERSION:-"dev"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.Date=$DATE"

# Directories
BIN_DIR="bin"
CMD_DIR="cmd"

# Create bin directory if it doesn't exist
mkdir -p "$BIN_DIR"

# Print banner
print_banner() {
    echo -e "${BLUE}"
    echo '   _____               _        _______          _ '
    echo '  / ____|             | |      |__   __|        | |'
    echo ' | (___  _ __   ___  | | _____   | | ___   ___ | |'
    echo '  \___ \| '"'"'_ \ / _ \ | |/ / _ \  | |/ _ \ / _ \| |'
    echo '  ____) | |_) | (_) |   <  __/  | | (_) | (_) | |'
    echo ' |_____/| .__/ \___/|_|\_\___|  |_|\___/ \___/|_|'
    echo '        | |                                      '
    echo '        |_|                                      '
    echo -e "${NC}"
    echo -e "${CYAN}Version: $VERSION (commit: $COMMIT, built: $DATE)${NC}"
    echo ""
}

# Print help
print_help() {
    echo -e "${YELLOW}Usage:${NC} ./scripts/build.sh [options] [targets]"
    echo ""
    echo "Options:"
    echo "  -h, --help      Show this help message"
    echo "  -v, --verbose   Verbose output"
    echo "  -r, --release   Build release binaries (stripped, optimized)"
    echo "  -o, --output    Output directory (default: bin/)"
    echo "  --os            Target OS (linux, windows, darwin) - for cross-compilation"
    echo "  --arch          Target architecture (amd64, arm64, 386)"
    echo ""
    echo "Targets:"
    echo "  all             Build all tools (default)"
    echo "  readmegen       Build only readmegen"
    echo "  testgen         Build only testgen"
    echo "  clean           Clean build artifacts"
    echo ""
    echo "Examples:"
    echo "  ./scripts/build.sh                # Build all tools"
    echo "  ./scripts/build.sh readmegen       # Build only readmegen"
    echo "  ./scripts/build.sh -r all          # Build release binaries"
    echo "  ./scripts/build.sh --os windows --arch amd64 testgen  # Cross-compile testgen for Windows"
}

# Build a specific tool
build_tool() {
    local tool=$1
    local output_name=$2
    
    echo -e "${YELLOW}Building ${tool}...${NC}"
    
    # Set output path
    local output_path="$BIN_DIR/$output_name"
    
    # Add .exe extension for Windows
    if [[ "$GOOS" == "windows" ]]; then
        output_path="$output_path.exe"
    fi
    
    # Build command
    local cmd="go build"
    
    # Add release flags if requested
    if [[ "$RELEASE" == "true" ]]; then
        echo -e "${PURPLE}Release build - stripping debug info${NC}"
        LDFLAGS="$LDFLAGS -s -w"
    fi
    
    cmd="$cmd -ldflags \"$LDFLAGS\" -o $output_path $CMD_DIR/$tool/main.go"
    
    # Add verbosity
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${CYAN}Running: $cmd${NC}"
        eval "$cmd -v"
    else
        eval "$cmd" 2>&1 | grep -v '^#' || exit 1
    fi
    
    # Check if build succeeded
    if [[ $? -eq 0 ]]; then
        echo -e "${GREEN}✅ Built $output_path${NC}"
        
        # Show file info
        if [[ -f "$output_path" ]]; then
            local size=$(du -h "$output_path" | cut -f1)
            echo -e "   Size: $size"
            
            # Make executable
            chmod +x "$output_path"
        fi
    else
        echo -e "${RED}❌ Failed to build $tool${NC}"
        exit 1
    fi
}

# Clean build artifacts
clean() {
    echo -e "${YELLOW}Cleaning build artifacts...${NC}"
    rm -rf "$BIN_DIR"
    go clean
    echo -e "${GREEN}✅ Clean complete${NC}"
}

# Check dependencies
check_deps() {
    echo -e "${YELLOW}Checking dependencies...${NC}"
    
    # Check Go version
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed${NC}"
        exit 1
    fi
    
    local go_version=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go $go_version${NC}"
    
    # Check git (optional, for version info)
    if ! command -v git &> /dev/null; then
        echo -e "${YELLOW}⚠️  git not found - version info will be limited${NC}"
    fi
    
    # Run go mod tidy
    echo -e "${YELLOW}Tidying dependencies...${NC}"
    go mod tidy
    
    echo ""
}

# Main build function
main() {
    # Default values
    VERBOSE="false"
    RELEASE="false"
    TARGET="all"
    GOOS=${GOOS:-$(go env GOOS)}
    GOARCH=${GOARCH:-$(go env GOARCH)}
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                print_help
                exit 0
                ;;
            -v|--verbose)
                VERBOSE="true"
                shift
                ;;
            -r|--release)
                RELEASE="true"
                shift
                ;;
            -o|--output)
                BIN_DIR="$2"
                shift 2
                ;;
            --os)
                GOOS="$2"
                shift 2
                ;;
            --arch)
                GOARCH="$2"
                shift 2
                ;;
            clean)
                clean
                exit 0
                ;;
            readmegen|testgen|all)
                TARGET="$1"
                shift
                ;;
            *)
                echo -e "${RED}Unknown option: $1${NC}"
                print_help
                exit 1
                ;;
        esac
    done
    
    # Export build environment
    export GOOS
    export GOARCH
    export CGO_ENABLED=0
    
    # Print banner
    print_banner
    
    # Show build environment
    echo -e "${CYAN}Build Environment:${NC}"
    echo -e "  OS:   $GOOS"
    echo -e "  Arch: $GOARCH"
    echo -e "  Dir:  $BIN_DIR"
    echo ""
    
    # Check dependencies
    check_deps
    
    # Build targets
    case $TARGET in
        all)
            echo -e "${CYAN}Building all tools${NC}"
            build_tool "readmegen" "readmegen"
            build_tool "testgen" "testgen"
            ;;
        readmegen)
            build_tool "readmegen" "readmegen"
            ;;
        testgen)
            build_tool "testgen" "testgen"
            ;;
    esac
    
    echo ""
    echo -e "${GREEN}✨ Build complete!${NC}"
    
    # Show binaries
    echo -e "${CYAN}Binaries in $BIN_DIR:${NC}"
    ls -lh "$BIN_DIR" 2>/dev/null || echo "No binaries found"
}

# Run main function
main "$@"
