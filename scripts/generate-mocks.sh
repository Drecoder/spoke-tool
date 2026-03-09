#!/usr/bin/env bash

# generate-mocks.sh - Generate mock implementations for testing
# This script uses standard mocking tools for each language

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Default values
OUTPUT_DIR=${OUTPUT_DIR:-"mocks"}
LANGUAGE=${LANGUAGE:-"auto"}
VERBOSE=${VERBOSE:-false}
RECURSIVE=${RECURSIVE:-false}
FORCE=${FORCE:-false}

# Print banner
print_banner() {
    echo -e "${BLUE}"
    echo '   _____                           _        '
    echo '  / ____|                         | |       '
    echo ' | |  __  ___ _ __   ___ _ __ __ _| | _____ '
    echo ' | | |_ |/ _ \ '"'"'_ \ / _ \ '"'"'__/ _` | |/ / _ \\'
    echo ' | |__| |  __/ | | |  __/ | | (_| |   <  __/'
    echo '  \_____|\___|_| |_|\___|_|  \__,_|_|\_\___|'
    echo -e "${NC}"
    echo -e "${CYAN}Mock Generator${NC}"
    echo ""
}

# Print help
print_help() {
    echo -e "${YELLOW}Usage:${NC} ./scripts/generate-mocks.sh [options] [targets]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -v, --verbose       Verbose output"
    echo "  -l, --language      Language (go, nodejs, python, auto)"
    echo "  -o, --output        Output directory (default: mocks/)"
    echo "  -r, --recursive     Generate mocks recursively"
    echo "  -f, --force         Force overwrite existing files"
    echo "  -p, --package       Package name (for Go)"
    echo "  -i, --interfaces    Specific interfaces to mock (comma-separated)"
    echo ""
    echo "Examples:"
    echo "  ./scripts/generate-mocks.sh                    # Auto-detect and generate mocks"
    echo "  ./scripts/generate-mocks.sh -l go              # Generate Go mocks"
    echo "  ./scripts/generate-mocks.sh -l nodejs -r       # Generate Node.js mocks recursively"
    echo "  ./scripts/generate-mocks.sh -i Reader,Writer   # Mock specific interfaces"
}

# Check dependencies
check_deps() {
    local lang=$1
    local missing=false
    
    echo -e "${YELLOW}Checking dependencies for $lang...${NC}"
    
    case $lang in
        go)
            if ! command -v mockgen &> /dev/null; then
                echo -e "${RED}❌ mockgen not found${NC}"
                echo -e "${YELLOW}Installing mockgen...${NC}"
                go install go.uber.org/mock/mockgen@latest
                if [[ $? -ne 0 ]]; then
                    echo -e "${RED}Failed to install mockgen${NC}"
                    missing=true
                fi
            else
                echo -e "${GREEN}✅ mockgen $(mockgen --version)${NC}"
            fi
            ;;
            
        nodejs)
            if ! command -v jest &> /dev/null; then
                echo -e "${YELLOW}⚠️  jest not found globally, checking local...${NC}"
                if [[ ! -f "node_modules/.bin/jest" ]]; then
                    echo -e "${RED}❌ jest not found${NC}"
                    echo -e "${YELLOW}Run: npm install --save-dev jest @types/jest${NC}"
                    missing=true
                else
                    echo -e "${GREEN}✅ jest (local)${NC}"
                fi
            else
                echo -e "${GREEN}✅ jest${NC}"
            fi
            
            # Check for ts-jest if TypeScript
            if [[ -f "tsconfig.json" ]]; then
                if [[ ! -f "node_modules/.bin/ts-jest" ]]; then
                    echo -e "${YELLOW}⚠️  ts-jest not found for TypeScript${NC}"
                fi
            fi
            ;;
            
        python)
            if ! command -v pytest &> /dev/null; then
                echo -e "${YELLOW}⚠️  pytest not found globally, checking venv...${NC}"
            fi
            
            # Check for pytest-mock
            python -c "import pytest_mock" 2>/dev/null
            if [[ $? -ne 0 ]]; then
                echo -e "${YELLOW}⚠️  pytest-mock not installed${NC}"
                echo -e "${YELLOW}Run: pip install pytest-mock${NC}"
            else
                echo -e "${GREEN}✅ pytest-mock${NC}"
            fi
            ;;
    esac
    
    if [[ "$missing" == "true" ]]; then
        echo -e "${RED}Missing dependencies. Please install them and try again.${NC}"
        exit 1
    fi
    
    echo ""
}

# Detect language from project files
detect_language() {
    if [[ -f "go.mod" ]] || ls *.go &>/dev/null; then
        echo "go"
    elif [[ -f "package.json" ]] || ls *.js &>/dev/null || ls *.ts &>/dev/null; then
        echo "nodejs"
    elif [[ -f "setup.py" ]] || [[ -f "requirements.txt" ]] || ls *.py &>/dev/null; then
        echo "python"
    else
        echo "unknown"
    fi
}

# Find interfaces in Go files
find_go_interfaces() {
    local dir=${1:-"."}
    local interfaces=()
    
    while IFS= read -r file; do
        if [[ -f "$file" ]]; then
            # Look for interface declarations
            while IFS= read -r line; do
                if [[ $line =~ ^type[[:space:]]+([a-zA-Z0-9_]+)[[:space:]]+interface ]]; then
                    interfaces+=("${BASH_REMATCH[1]}")
                fi
            done < "$file"
        fi
    done < <(find "$dir" -name "*.go" -not -name "*_test.go" 2>/dev/null)
    
    printf '%s\n' "${interfaces[@]}"
}

# Generate Go mocks
generate_go_mocks() {
    local pkg=${1:-"."}
    local output=${2:-"mocks"}
    local specific_interfaces=($3)
    
    echo -e "${CYAN}Generating Go mocks...${NC}"
    
    # Create output directory
    mkdir -p "$output"
    
    # Find all interfaces
    if [[ ${#specific_interfaces[@]} -eq 0 ]]; then
        echo -e "${YELLOW}Finding interfaces in $pkg...${NC}"
        mapfile -t interfaces < <(find_go_interfaces "$pkg")
    else
        interfaces=("${specific_interfaces[@]}")
    fi
    
    if [[ ${#interfaces[@]} -eq 0 ]]; then
        echo -e "${YELLOW}No interfaces found${NC}"
        return 0
    fi
    
    echo -e "${GREEN}Found ${#interfaces[@]} interfaces${NC}"
    
    # Generate mocks
    for iface in "${interfaces[@]}"; do
        echo -e "  Generating mock for $iface..."
        
        local mock_file="$output/mock_${iface}.go"
        
        # Check if file exists and not forcing
        if [[ -f "$mock_file" && "$FORCE" != "true" ]]; then
            echo -e "  ${YELLOW}⚠️  $mock_file exists (use -f to overwrite)${NC}"
            continue
        fi
        
        # Generate mock
        if [[ "$VERBOSE" == "true" ]]; then
            mockgen -package "$output" -destination "$mock_file" "$pkg" "$iface"
        else
            mockgen -package "$output" -destination "$mock_file" "$pkg" "$iface" 2>/dev/null
        fi
        
        if [[ $? -eq 0 ]]; then
            echo -e "  ${GREEN}✅ Generated $mock_file${NC}"
        else
            echo -e "  ${RED}❌ Failed to generate mock for $iface${NC}"
        fi
    done
}

# Generate Node.js mocks
generate_nodejs_mocks() {
    local output=${1:-"mocks"}
    local specific_interfaces=($2)
    
    echo -e "${CYAN}Generating Node.js mocks...${NC}"
    
    mkdir -p "$output"
    
    # Create Jest mock file
    local mock_file="$output/__mocks__"
    mkdir -p "$mock_file"
    
    # Generate manual mocks
    cat > "$mock_file/jest.setup.js" << 'EOF'
// Auto-generated Jest mocks
jest.mock('fs', () => ({
    readFileSync: jest.fn(),
    writeFileSync: jest.fn(),
    existsSync: jest.fn(),
}));

jest.mock('path', () => ({
    join: jest.fn((...args) => args.join('/')),
    resolve: jest.fn((...args) => args.join('/')),
    dirname: jest.fn((p) => p.split('/').slice(0, -1).join('/')),
    basename: jest.fn((p) => p.split('/').pop()),
}));

// Helper to create mock functions
global.createMock = (implementation) => {
    const mock = jest.fn(implementation);
    mock.mockImplementation = jest.fn().mockImplementation(implementation);
    mock.mockReturnValue = jest.fn().mockReturnValue;
    mock.mockResolvedValue = jest.fn().mockResolvedValue;
    mock.mockRejectedValue = jest.fn().mockRejectedValue;
    return mock;
};
EOF
    
    echo -e "${GREEN}✅ Generated Jest setup in $mock_file/jest.setup.js${NC}"
    
    # Generate TypeScript types if needed
    if [[ -f "tsconfig.json" ]]; then
        cat > "$mock_file/types.d.ts" << 'EOF'
// Auto-generated mock types
declare module 'fs' {
    export const readFileSync: jest.Mock;
    export const writeFileSync: jest.Mock;
    export const existsSync: jest.Mock;
}

declare module 'path' {
    export const join: jest.Mock;
    export const resolve: jest.Mock;
    export const dirname: jest.Mock;
    export const basename: jest.Mock;
}

declare global {
    function createMock<T extends (...args: any[]) => any>(implementation?: T): jest.Mock<ReturnType<T>, Parameters<T>>;
}
EOF
        echo -e "${GREEN}✅ Generated TypeScript types${NC}"
    fi
    
    # Update package.json if needed
    if [[ -f "package.json" ]]; then
        if ! grep -q '"jest"' package.json; then
            echo -e "${YELLOW}Adding Jest configuration to package.json...${NC}"
            # This would need jq for proper JSON manipulation
            # For now, just show instructions
            echo -e "${YELLOW}Add to package.json:${NC}"
            cat << 'EOF'
{
  "jest": {
    "setupFilesAfterEnv": ["<rootDir>/mocks/__mocks__/jest.setup.js"],
    "moduleNameMapper": {
      "^fs$": "<rootDir>/mocks/__mocks__/fs.js"
    }
  }
}
EOF
        fi
    fi
}

# Generate Python mocks
generate_python_mocks() {
    local output=${1:-"mocks"}
    local specific_interfaces=($2)
    
    echo -e "${CYAN}Generating Python mocks...${NC}"
    
    mkdir -p "$output"
    
    # Create conftest.py with fixtures
    local conftest="$output/conftest.py"
    
    cat > "$conftest" << 'EOF'
"""Auto-generated pytest fixtures and mocks."""
import pytest
from unittest.mock import Mock, patch, MagicMock

@pytest.fixture
def mock_file_system():
    """Mock file system operations."""
    with patch('builtins.open', create=True) as mock_open:
        mock_open.return_value.__enter__.return_value.read.return_value = ""
        mock_open.return_value.__enter__.return_value.write.return_value = None
        yield mock_open

@pytest.fixture
def mock_os():
    """Mock os module."""
    with patch('os.path') as mock_path, \
         patch('os.makedirs') as mock_makedirs, \
         patch('os.listdir') as mock_listdir:
        mock_path.exists.return_value = True
        mock_path.join.side_effect = lambda *args: '/'.join(args)
        mock_makedirs.return_value = None
        mock_listdir.return_value = []
        yield {
            'path': mock_path,
            'makedirs': mock_makedirs,
            'listdir': mock_listdir
        }

@pytest.fixture
def mock_logger():
    """Mock logger."""
    return MagicMock(
        debug=Mock(),
        info=Mock(),
        warning=Mock(),
        error=Mock(),
        critical=Mock()
    )

# Helper to create mock objects
def create_mock(spec=None, **kwargs):
    """Create a mock object with optional spec."""
    return Mock(spec=spec, **kwargs)

# Helper to patch multiple objects
def patch_multiple(targets):
    """Context manager to patch multiple targets."""
    patches = [patch(target) for target in targets]
    for p in patches:
        p.start()
    yield
    for p in patches:
        p.stop()
EOF
    
    echo -e "${GREEN}✅ Generated pytest fixtures in $conftest${NC}"
    
    # Create __init__.py to make it a package
    touch "$output/__init__.py"
    
    # Generate specific mocks if requested
    if [[ ${#specific_interfaces[@]} -gt 0 ]]; then
        local mock_file="$output/mocks.py"
        cat > "$mock_file" << 'EOF'
"""Auto-generated mock classes."""
from unittest.mock import Mock, MagicMock

EOF
        
        for iface in "${specific_interfaces[@]}"; do
            cat >> "$mock_file" << EOF

class Mock${iface}:
    """Mock implementation of ${iface}."""
    
    def __init__(self, **kwargs):
        ${iface}_methods = [m for m in dir(${iface}) if not m.startswith('_')]
        for method in ${iface}_methods:
            setattr(self, method, Mock(**kwargs))
EOF
            echo -e "${GREEN}✅ Added mock for $iface${NC}"
        done
    fi
}

# Main function
main() {
    local language="$LANGUAGE"
    local output="$OUTPUT_DIR"
    local recursive="$RECURSIVE"
    local force="$FORCE"
    local interfaces=()
    
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
            -l|--language)
                LANGUAGE="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -r|--recursive)
                RECURSIVE=true
                shift
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            -p|--package)
                PACKAGE="$2"
                shift 2
                ;;
            -i|--interfaces)
                IFS=',' read -ra interfaces <<< "$2"
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
    
    # Detect language if auto
    if [[ "$LANGUAGE" == "auto" ]]; then
        LANGUAGE=$(detect_language)
        echo -e "${CYAN}Detected language: $LANGUAGE${NC}"
    fi
    
    # Check if language is supported
    if [[ "$LANGUAGE" != "go" && "$LANGUAGE" != "nodejs" && "$LANGUAGE" != "python" ]]; then
        echo -e "${RED}Unsupported language: $LANGUAGE${NC}"
        exit 1
    fi
    
    # Check dependencies
    check_deps "$LANGUAGE"
    
    # Generate mocks based on language
    case $LANGUAGE in
        go)
            local pkg=${PACKAGE:-"."}
            if [[ "$RECURSIVE" == "true" ]]; then
                # Find all packages
                while IFS= read -r p; do
                    if [[ -d "$p" ]] && ls "$p"/*.go &>/dev/null; then
                        echo -e "${CYAN}Processing package: $p${NC}"
                        generate_go_mocks "$p" "$OUTPUT_DIR/$p" "${interfaces[*]}"
                    fi
                done < <(find . -type d -not -path "*/\.*" -not -path "*/vendor/*" -not -path "*/node_modules/*")
            else
                generate_go_mocks "$pkg" "$OUTPUT_DIR" "${interfaces[*]}"
            fi
            ;;
            
        nodejs)
            generate_nodejs_mocks "$OUTPUT_DIR" "${interfaces[*]}"
            ;;
            
        python)
            generate_python_mocks "$OUTPUT_DIR" "${interfaces[*]}"
            ;;
    esac
    
    echo ""
    echo -e "${GREEN}✅ Mock generation complete!${NC}"
    
    # Show output location
    echo -e "${CYAN}Mocks generated in: $OUTPUT_DIR${NC}"
    
    # Show next steps
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    case $LANGUAGE in
        go)
            echo "  import ("
            echo "    \"yourproject/$OUTPUT_DIR\""
            echo "  )"
            ;;
        nodejs)
            echo "  import mock from '../$OUTPUT_DIR/__mocks__';"
            echo "  // Add to jest.config.js: setupFilesAfterEnv: ['<rootDir>/$OUTPUT_DIR/__mocks__/jest.setup.js']"
            ;;
        python)
            echo "  # Add to conftest.py or import fixtures"
            echo "  from $OUTPUT_DIR.conftest import *"
            ;;
    esac
}

# Run main function
main "$@"