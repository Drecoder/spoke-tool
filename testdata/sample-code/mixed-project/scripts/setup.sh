#!/usr/bin/env bash

# ============================================================================
# Setup Script for Mixed-Language Project
# ============================================================================
#
# This script sets up the entire development environment for the mixed-language
# project, including:
#   - System dependencies
#   - Go environment
#   - Node.js environment
#   - Python environment
#   - Database setup
#   - Docker setup
#   - Git hooks
#   - IDE configuration
#
# ============================================================================

set -e  # Exit on error
set -u  # Exit on undefined variable

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

# Project paths
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GO_SERVER_DIR="${PROJECT_ROOT}/go-server"
WEB_CLIENT_DIR="${PROJECT_ROOT}/web-client"
SCRIPTS_DIR="${PROJECT_ROOT}/scripts"
BACKUP_DIR="${PROJECT_ROOT}/backups"
LOGS_DIR="${PROJECT_ROOT}/logs"
CONFIG_DIR="${PROJECT_ROOT}/config"

# Version requirements
GO_VERSION_MIN="1.21"
NODE_VERSION_MIN="20"
PYTHON_VERSION_MIN="3.11"
DOCKER_VERSION_MIN="24.0"

# Default ports
GO_PORT=8080
NODE_PORT=3000
PYTHON_PORT=8000
POSTGRES_PORT=5432
REDIS_PORT=6379

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

print_step() {
    echo -e "\n${MAGENTA}➡️  $1${NC}"
}

check_command() {
    if command -v "$1" &> /dev/null; then
        print_success "$1 installed"
        return 0
    else
        print_error "$1 not found"
        return 1
    fi
}

version_gte() {
    printf '%s\n%s\n' "$2" "$1" | sort -V -C
}

confirm() {
    read -p "$(echo -e "${YELLOW}$1 [y/N]${NC} ")" -n 1 -r
    echo
    [[ $REPLY =~ ^[Yy]$ ]]
}

# ============================================================================
# System Detection
# ============================================================================

detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        CYGWIN*|MINGW*|MSYS*) echo "windows";;
        *)          echo "unknown";;
    esac
}

OS=$(detect_os)
print_info "Detected OS: ${OS}"

# ============================================================================
# System Dependencies
# ============================================================================

install_system_deps() {
    print_step "Installing system dependencies"
    
    case $OS in
        linux)
            if command -v apt-get &> /dev/null; then
                # Debian/Ubuntu
                sudo apt-get update
                sudo apt-get install -y \
                    build-essential \
                    curl \
                    wget \
                    git \
                    make \
                    unzip \
                    ca-certificates \
                    gnupg \
                    lsb-release \
                    postgresql-client \
                    redis-tools
                
            elif command -v yum &> /dev/null; then
                # CentOS/RHEL
                sudo yum groupinstall -y "Development Tools"
                sudo yum install -y \
                    curl \
                    wget \
                    git \
                    make \
                    unzip \
                    ca-certificates \
                    postgresql \
                    redis
                
            elif command -v pacman &> /dev/null; then
                # Arch Linux
                sudo pacman -Syu --noconfirm \
                    base-devel \
                    curl \
                    wget \
                    git \
                    make \
                    unzip \
                    ca-certificates \
                    postgresql \
                    redis
            else
                print_warning "Unsupported Linux distribution. Please install dependencies manually."
            fi
            ;;
            
        darwin)
            if ! command -v brew &> /dev/null; then
                print_info "Installing Homebrew..."
                /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
            fi
            
            brew update
            brew install \
                curl \
                wget \
                git \
                make \
                unzip \
                postgresql \
                redis
            ;;
            
        windows)
            print_warning "Windows detected. Please install dependencies manually:"
            print_info "  - Git for Windows: https://git-scm.com/download/win"
            print_info "  - PostgreSQL: https://www.postgresql.org/download/windows/"
            print_info "  - Redis: https://github.com/microsoftarchive/redis/releases"
            ;;
    esac
    
    print_success "System dependencies installed"
}

# ============================================================================
# Go Installation
# ============================================================================

install_go() {
    print_step "Setting up Go environment"
    
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_info "Go ${GO_VERSION} already installed"
        
        if version_gte "$GO_VERSION" "$GO_VERSION_MIN"; then
            print_success "Go version ${GO_VERSION} meets minimum requirement (${GO_VERSION_MIN})"
        else
            print_warning "Go version ${GO_VERSION} is below minimum requirement ${GO_VERSION_MIN}"
            if confirm "Install newer version?"; then
                install_go_from_source
            fi
        fi
    else
        print_warning "Go not found"
        if confirm "Install Go?"; then
            install_go_from_source
        fi
    fi
    
    # Set up Go workspace
    if command -v go &> /dev/null; then
        mkdir -p "${HOME}/go/bin"
        mkdir -p "${HOME}/go/pkg"
        mkdir -p "${HOME}/go/src"
        
        # Install Go tools
        print_info "Installing Go tools..."
        go install golang.org/x/tools/cmd/goimports@latest
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        go install github.com/go-delve/delve/cmd/dlv@latest
        go install github.com/cosmtrek/air@latest
        go install github.com/swaggo/swag/cmd/swag@latest
        
        print_success "Go tools installed"
    fi
}

install_go_from_source() {
    case $OS in
        linux)
            wget "https://go.dev/dl/go${GO_VERSION_MIN}.linux-amd64.tar.gz"
            sudo rm -rf /usr/local/go
            sudo tar -C /usr/local -xzf "go${GO_VERSION_MIN}.linux-amd64.tar.gz"
            rm "go${GO_VERSION_MIN}.linux-amd64.tar.gz"
            export PATH=$PATH:/usr/local/go/bin
            echo 'export PATH=$PATH:/usr/local/go/bin' >> "${HOME}/.bashrc"
            ;;
            
        darwin)
            brew install go@${GO_VERSION_MIN}
            ;;
            
        windows)
            print_info "Download Go from: https://go.dev/dl/"
            ;;
    esac
}

# ============================================================================
# Node.js Installation
# ============================================================================

install_node() {
    print_step "Setting up Node.js environment"
    
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node --version | sed 's/v//')
        print_info "Node.js ${NODE_VERSION} already installed"
        
        if version_gte "$NODE_VERSION" "$NODE_VERSION_MIN"; then
            print_success "Node.js version ${NODE_VERSION} meets minimum requirement (${NODE_VERSION_MIN})"
        else
            print_warning "Node.js version ${NODE_VERSION} is below minimum requirement ${NODE_VERSION_MIN}"
            if confirm "Install newer version?"; then
                install_node_from_source
            fi
        fi
    else
        print_warning "Node.js not found"
        if confirm "Install Node.js?"; then
            install_node_from_source
        fi
    fi
    
    # Install global packages
    if command -v npm &> /dev/null; then
        print_info "Installing global npm packages..."
        npm install -g npm@latest
        npm install -g yarn
        npm install -g pnpm
        npm install -g typescript
        npm install -g ts-node
        npm install -g nodemon
        npm install -g pm2
        npm install -g jest
        
        print_success "Global npm packages installed"
    fi
}

install_node_from_source() {
    case $OS in
        linux|darwin)
            curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
            sudo apt-get install -y nodejs
            ;;
            
        windows)
            print_info "Download Node.js from: https://nodejs.org/"
            ;;
    esac
}

# ============================================================================
# Python Installation
# ============================================================================

install_python() {
    print_step "Setting up Python environment"
    
    if command -v python3 &> /dev/null; then
        PYTHON_VERSION=$(python3 --version | awk '{print $2}')
        print_info "Python ${PYTHON_VERSION} already installed"
        
        if version_gte "$PYTHON_VERSION" "$PYTHON_VERSION_MIN"; then
            print_success "Python version ${PYTHON_VERSION} meets minimum requirement (${PYTHON_VERSION_MIN})"
        else
            print_warning "Python version ${PYTHON_VERSION} is below minimum requirement ${PYTHON_VERSION_MIN}"
            if confirm "Install newer version?"; then
                install_python_from_source
            fi
        fi
    else
        print_warning "Python not found"
        if confirm "Install Python?"; then
            install_python_from_source
        fi
    fi
    
    # Install pip if not present
    if ! command -v pip3 &> /dev/null; then
        print_info "Installing pip..."
        curl -sS https://bootstrap.pypa.io/get-pip.py | python3
    fi
    
    # Install global packages
    if command -v pip3 &> /dev/null; then
        print_info "Installing global Python packages..."
        pip3 install --upgrade pip
        pip3 install virtualenv
        pip3 install poetry
        pip3 install pytest
        pip3 install pytest-cov
        pip3 install black
        pip3 install flake8
        pip3 install mypy
        
        print_success "Global Python packages installed"
    fi
}

install_python_from_source() {
    case $OS in
        linux)
            sudo apt-get update
            sudo apt-get install -y python3.11 python3.11-dev python3.11-venv
            ;;
            
        darwin)
            brew install python@3.11
            ;;
            
        windows)
            print_info "Download Python from: https://python.org/"
            ;;
    esac
}

# ============================================================================
# Docker Installation
# ============================================================================

install_docker() {
    print_step "Setting up Docker"
    
    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
        print_info "Docker ${DOCKER_VERSION} already installed"
        
        if version_gte "$DOCKER_VERSION" "$DOCKER_VERSION_MIN"; then
            print_success "Docker version ${DOCKER_VERSION} meets minimum requirement (${DOCKER_VERSION_MIN})"
        else
            print_warning "Docker version ${DOCKER_VERSION} is below minimum requirement ${DOCKER_VERSION_MIN}"
            if confirm "Install newer version?"; then
                install_docker_from_source
            fi
        fi
    else
        print_warning "Docker not found"
        if confirm "Install Docker?"; then
            install_docker_from_source
        fi
    fi
    
    # Install Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        print_info "Installing Docker Compose..."
        sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        print_success "Docker Compose installed"
    fi
}

install_docker_from_source() {
    case $OS in
        linux)
            curl -fsSL https://get.docker.com -o get-docker.sh
            sudo sh get-docker.sh
            sudo usermod -aG docker "${USER}"
            rm get-docker.sh
            ;;
            
        darwin)
            brew install docker docker-compose
            ;;
            
        windows)
            print_info "Download Docker Desktop from: https://docker.com/"
            ;;
    esac
}

# ============================================================================
# Database Setup
# ============================================================================

setup_databases() {
    print_step "Setting up databases"
    
    # Start PostgreSQL
    case $OS in
        linux)
            if command -v systemctl &> /dev/null; then
                sudo systemctl start postgresql
                sudo systemctl enable postgresql
            fi
            ;;
        darwin)
            brew services start postgresql
            ;;
    esac
    
    # Start Redis
    case $OS in
        linux)
            if command -v systemctl &> /dev/null; then
                sudo systemctl start redis-server
                sudo systemctl enable redis-server
            fi
            ;;
        darwin)
            brew services start redis
            ;;
    esac
    
    # Wait for databases to be ready
    print_info "Waiting for databases..."
    sleep 5
    
    # Create databases
    if command -v psql &> /dev/null; then
        print_info "Creating PostgreSQL databases..."
        sudo -u postgres psql -c "CREATE USER myapp WITH PASSWORD 'myapp' SUPERUSER;" 2>/dev/null || true
        sudo -u postgres psql -c "CREATE DATABASE myapp_dev OWNER myapp;" 2>/dev/null || true
        sudo -u postgres psql -c "CREATE DATABASE myapp_test OWNER myapp;" 2>/dev/null || true
        print_success "PostgreSQL databases created"
    fi
    
    if command -v redis-cli &> /dev/null; then
        print_info "Testing Redis..."
        redis-cli ping
        print_success "Redis is running"
    fi
}

# ============================================================================
# Project Setup
# ============================================================================

setup_project() {
    print_step "Setting up project"
    
    # Create directories
    mkdir -p "${BACKUP_DIR}"
    mkdir -p "${LOGS_DIR}"
    mkdir -p "${CONFIG_DIR}"
    
    # Copy example env files
    if [ -f "${GO_SERVER_DIR}/.env.example" ]; then
        cp "${GO_SERVER_DIR}/.env.example" "${GO_SERVER_DIR}/.env"
        print_success "Created Go server .env file"
    fi
    
    if [ -f "${WEB_CLIENT_DIR}/.env.example" ]; then
        cp "${WEB_CLIENT_DIR}/.env.example" "${WEB_CLIENT_DIR}/.env"
        print_success "Created web client .env file"
    fi
    
    if [ -f "${SCRIPTS_DIR}/.env.example" ]; then
        cp "${SCRIPTS_DIR}/.env.example" "${SCRIPTS_DIR}/.env"
        print_success "Created scripts .env file"
    fi
    
    # Make scripts executable
    find "${SCRIPTS_DIR}" -name "*.sh" -exec chmod +x {} \;
    find "${SCRIPTS_DIR}" -name "*.py" -exec chmod +x {} \;
    find "${SCRIPTS_DIR}" -name "*.js" -exec chmod +x {} \;
    
    print_success "Scripts made executable"
}

# ============================================================================
# Go Dependencies
# ============================================================================

setup_go() {
    print_step "Setting up Go dependencies"
    
    if [ -f "${GO_SERVER_DIR}/go.mod" ]; then
        cd "${GO_SERVER_DIR}"
        
        print_info "Downloading Go modules..."
        go mod download
        go mod verify
        
        print_info "Generating Swagger docs..."
        if command -v swag &> /dev/null; then
            swag init
        fi
        
        print_info "Building Go server..."
        go build -o bin/server ./cmd/server
        
        print_success "Go dependencies installed"
    else
        print_warning "No go.mod found in ${GO_SERVER_DIR}"
    fi
}

# ============================================================================
# Node.js Dependencies
# ============================================================================

setup_node() {
    print_step "Setting up Node.js dependencies"
    
    if [ -f "${WEB_CLIENT_DIR}/package.json" ]; then
        cd "${WEB_CLIENT_DIR}"
        
        print_info "Installing npm packages..."
        npm install
        
        print_info "Building web client..."
        npm run build
        
        print_success "Node.js dependencies installed"
    else
        print_warning "No package.json found in ${WEB_CLIENT_DIR}"
    fi
}

# ============================================================================
# Python Dependencies
# ============================================================================

setup_python() {
    print_step "Setting up Python dependencies"
    
    cd "${SCRIPTS_DIR}"
    
    if [ -f "requirements.txt" ]; then
        print_info "Creating virtual environment..."
        python3 -m venv venv
        
        print_info "Activating virtual environment..."
        source venv/bin/activate
        
        print_info "Installing Python packages..."
        pip install --upgrade pip
        pip install -r requirements.txt
        
        print_success "Python dependencies installed"
    else
        print_warning "No requirements.txt found in ${SCRIPTS_DIR}"
    fi
}

# ============================================================================
# Git Hooks
# ============================================================================

setup_git_hooks() {
    print_step "Setting up Git hooks"
    
    HOOKS_DIR="${PROJECT_ROOT}/.git/hooks"
    
    # Pre-commit hook
    cat > "${HOOKS_DIR}/pre-commit" << 'EOF'
#!/bin/bash
echo "Running pre-commit checks..."

# Run tests
make test || exit 1

# Run linters
make lint || exit 1

# Check formatting
make fmt-check || exit 1

echo "✅ Pre-commit checks passed"
EOF
    chmod +x "${HOOKS_DIR}/pre-commit"
    
    # Pre-push hook
    cat > "${HOOKS_DIR}/pre-push" << 'EOF'
#!/bin/bash
echo "Running pre-push checks..."

# Run integration tests
make test-integration || exit 1

echo "✅ Pre-push checks passed"
EOF
    chmod +x "${HOOKS_DIR}/pre-push"
    
    print_success "Git hooks installed"
}

# ============================================================================
# IDE Configuration
# ============================================================================

setup_ide() {
    print_step "Setting up IDE configuration"
    
    # VS Code settings
    VSCODE_DIR="${PROJECT_ROOT}/.vscode"
    mkdir -p "${VSCODE_DIR}"
    
    # Settings
    cat > "${VSCODE_DIR}/settings.json" << EOF
{
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.fixAll.eslint": true
    },
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.vetOnSave": "package",
    "python.formatting.provider": "black",
    "python.linting.enabled": true,
    "python.linting.flake8Enabled": true,
    "python.linting.mypyEnabled": true,
    "files.associations": {
        "*.env*": "dotenv"
    }
}
EOF
    
    # Extensions
    cat > "${VSCODE_DIR}/extensions.json" << EOF
{
    "recommendations": [
        "golang.go",
        "ms-python.python",
        "ms-python.black-formatter",
        "ms-python.flake8",
        "ms-python.mypy",
        "esbenp.prettier-vscode",
        "dbaeumer.vscode-eslint",
        "eamodio.gitlens",
        "ms-azuretools.vscode-docker"
    ]
}
EOF
    
    # Launch configuration
    cat > "${VSCODE_DIR}/launch.json" << EOF
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Go Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "\${workspaceFolder}/go-server/cmd/server"
        },
        {
            "name": "Node Client",
            "type": "node",
            "request": "launch",
            "program": "\${workspaceFolder}/web-client/src/index.js"
        },
        {
            "name": "Python Script",
            "type": "python",
            "request": "launch",
            "program": "\${workspaceFolder}/scripts/main.py"
        }
    ]
}
EOF
    
    print_success "IDE configuration created"
}

# ============================================================================
# Docker Setup
# ============================================================================

setup_docker() {
    print_step "Setting up Docker environment"
    
    if command -v docker-compose &> /dev/null; then
        cd "${PROJECT_ROOT}"
        
        print_info "Building Docker images..."
        docker-compose build
        
        print_info "Starting services..."
        docker-compose up -d
        
        print_success "Docker services started"
    else
        print_warning "Docker Compose not found"
    fi
}

# ============================================================================
# Verification
# ============================================================================

verify_installation() {
    print_step "Verifying installation"
    
    errors=0
    
    # Check Go
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        if version_gte "$GO_VERSION" "$GO_VERSION_MIN"; then
            print_success "Go ${GO_VERSION} ✓"
        else
            print_error "Go ${GO_VERSION} (min ${GO_VERSION_MIN} required)"
            ((errors++))
        fi
    else
        print_error "Go not found"
        ((errors++))
    fi
    
    # Check Node.js
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node --version | sed 's/v//')
        if version_gte "$NODE_VERSION" "$NODE_VERSION_MIN"; then
            print_success "Node.js ${NODE_VERSION} ✓"
        else
            print_error "Node.js ${NODE_VERSION} (min ${NODE_VERSION_MIN} required)"
            ((errors++))
        fi
    else
        print_error "Node.js not found"
        ((errors++))
    fi
    
    # Check Python
    if command -v python3 &> /dev/null; then
        PYTHON_VERSION=$(python3 --version | awk '{print $2}')
        if version_gte "$PYTHON_VERSION" "$PYTHON_VERSION_MIN"; then
            print_success "Python ${PYTHON_VERSION} ✓"
        else
            print_error "Python ${PYTHON_VERSION} (min ${PYTHON_VERSION_MIN} required)"
            ((errors++))
        fi
    else
        print_error "Python not found"
        ((errors++))
    fi
    
    # Check Docker
    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
        if version_gte "$DOCKER_VERSION" "$DOCKER_VERSION_MIN"; then
            print_success "Docker ${DOCKER_VERSION} ✓"
        else
            print_error "Docker ${DOCKER_VERSION} (min ${DOCKER_VERSION_MIN} required)"
            ((errors++))
        fi
    else
        print_error "Docker not found"
        ((errors++))
    fi
    
    # Check PostgreSQL
    if command -v psql &> /dev/null; then
        print_success "PostgreSQL client ✓"
    else
        print_error "PostgreSQL client not found"
        ((errors++))
    fi
    
    # Check Redis
    if command -v redis-cli &> /dev/null; then
        print_success "Redis client ✓"
    else
        print_error "Redis client not found"
        ((errors++))
    fi
    
    if [ $errors -eq 0 ]; then
        print_success "All dependencies verified successfully!"
    else
        print_error "${errors} errors found during verification"
        return 1
    fi
}

# ============================================================================
# Cleanup
# ============================================================================

cleanup() {
    print_step "Cleaning up"
    
    # Remove temporary files
    find "${PROJECT_ROOT}" -name "*.pyc" -delete
    find "${PROJECT_ROOT}" -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
    find "${PROJECT_ROOT}" -name ".DS_Store" -delete
    find "${PROJECT_ROOT}" -name "*.log" -delete
    
    print_success "Cleanup completed"
}

# ============================================================================
# Main Menu
# ============================================================================

show_menu() {
    echo -e "\n${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${WHITE}              Mixed-Language Project Setup Menu${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}\n"
    
    echo "1) Full setup (everything)"
    echo "2) Install system dependencies"
    echo "3) Install Go environment"
    echo "4) Install Node.js environment"
    echo "5) Install Python environment"
    echo "6) Install Docker"
    echo "7) Setup databases"
    echo "8) Setup project files"
    echo "9) Install Go dependencies"
    echo "10) Install Node.js dependencies"
    echo "11) Install Python dependencies"
    echo "12) Setup Git hooks"
    echo "13) Setup IDE configuration"
    echo "14) Setup Docker services"
    echo "15) Verify installation"
    echo "16) Cleanup"
    echo "17) Exit"
    echo
    read -p "$(echo -e "${YELLOW}Enter your choice [1-17]:${NC} ")" choice
    
    case $choice in
        1)
            install_system_deps
            install_go
            install_node
            install_python
            install_docker
            setup_databases
            setup_project
            setup_go
            setup_node
            setup_python
            setup_git_hooks
            setup_ide
            setup_docker
            verify_installation
            ;;
        2) install_system_deps ;;
        3) install_go ;;
        4) install_node ;;
        5) install_python ;;
        6) install_docker ;;
        7) setup_databases ;;
        8) setup_project ;;
        9) setup_go ;;
        10) setup_node ;;
        11) setup_python ;;
        12) setup_git_hooks ;;
        13) setup_ide ;;
        14) setup_docker ;;
        15) verify_installation ;;
        16) cleanup ;;
        17) exit 0 ;;
        *) print_error "Invalid choice" ;;
    esac
    
    if [ "$choice" != "17" ]; then
        echo
        read -p "$(echo -e "${YELLOW}Press Enter to continue...${NC}")"
        show_menu
    fi
}

# ============================================================================
# Main
# ============================================================================

main() {
    print_header "Mixed-Language Project Setup"
    
    # Check if running with appropriate permissions
    if [ "$EUID" -eq 0 ]; then
        print_warning "Running as root"
    fi
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --full)
                install_system_deps
                install_go
                install_node
                install_python
                install_docker
                setup_databases
                setup_project
                setup_go
                setup_node
                setup_python
                setup_git_hooks
                setup_ide
                setup_docker
                verify_installation
                exit 0
                ;;
            --quick)
                setup_project
                setup_go
                setup_node
                setup_python
                exit 0
                ;;
            --verify)
                verify_installation
                exit 0
                ;;
            --clean)
                cleanup
                exit 0
                ;;
            --help)
                echo "Usage: $0 [OPTION]"
                echo "Options:"
                echo "  --full    Perform full setup"
                echo "  --quick   Quick setup (project only)"
                echo "  --verify  Verify installation"
                echo "  --clean   Clean up temporary files"
                echo "  --help    Show this help"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
        shift
    done
    
    # Interactive mode
    show_menu
}

# Run main function
main "$@"