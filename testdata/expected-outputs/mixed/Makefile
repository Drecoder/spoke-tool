# ============================================================================
# Mixed Language Project Makefile
# Supports: Go, Node.js, Python
# ============================================================================

.PHONY: help setup build clean test lint run dev docker prod deploy

# ============================================================================
# Configuration
# ============================================================================

SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go configuration
GO_VERSION ?= 1.21
GO_PACKAGES ?= ./...
GO_FLAGS ?= -v
GO_LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"
GO_TEST_FLAGS ?= -race -cover

# Node.js configuration
NODE_VERSION ?= 20
NPM_FLAGS ?= --silent
NODE_ENV ?= development

# Python configuration
PYTHON_VERSION ?= 3.11
PYTHON_FLAGS ?= -v
PIP_FLAGS ?= --quiet
PYTHON_ENV ?= development

# Docker configuration
DOCKER_REGISTRY ?= docker.io
DOCKER_NAMESPACE ?= myorg
DOCKER_TAG ?= $(VERSION)
DOCKER_BUILDKIT ?= 1
COMPOSE_PROJECT_NAME ?= myapp

# Paths
GO_SRC = ./go-server
NODE_SRC = ./web-client
PYTHON_SRC = ./scripts
BIN_DIR = ./bin
COVERAGE_DIR = ./coverage
TEST_RESULTS_DIR = ./test-results
DIST_DIR = ./dist

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
CYAN := \033[0;36m
NC := \033[0m # No Color

# ============================================================================
# Help
# ============================================================================

help: ## Show this help message
	@printf "$(CYAN)Usage: make [target]$(NC)\n\n"
	@printf "$(YELLOW)Targets:$(NC)\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@printf "\n$(YELLOW)Services:$(NC)\n"
	@printf "  $(BLUE)go-api$(NC)       - Go backend service\n"
	@printf "  $(BLUE)node-frontend$(NC) - Node.js frontend service\n"
	@printf "  $(BLUE)python-processor$(NC) - Python data processor\n"

# ============================================================================
# Setup
# ============================================================================

setup: setup-go setup-node setup-python ## Install all dependencies
	@printf "$(GREEN)✓ All dependencies installed$(NC)\n"

setup-go: ## Install Go dependencies
	@printf "$(YELLOW)Installing Go dependencies...$(NC)\n"
	cd $(GO_SRC) && go mod download
	cd $(GO_SRC) && go mod verify
	@printf "$(GREEN)✓ Go dependencies installed$(NC)\n"

setup-node: ## Install Node.js dependencies
	@printf "$(YELLOW)Installing Node.js dependencies...$(NC)\n"
	cd $(NODE_SRC) && npm install $(NPM_FLAGS)
	@printf "$(GREEN)✓ Node.js dependencies installed$(NC)\n"

setup-python: ## Install Python dependencies
	@printf "$(YELLOW)Installing Python dependencies...$(NC)\n"
	cd $(PYTHON_SRC) && pip install $(PIP_FLAGS) -r requirements.txt
	cd $(PYTHON_SRC) && pip install $(PIP_FLAGS) -r requirements-dev.txt
	@printf "$(GREEN)✓ Python dependencies installed$(NC)\n"

# ============================================================================
# Build
# ============================================================================

build: build-go build-node build-python ## Build all services

build-go: ## Build Go service
	@printf "$(YELLOW)Building Go service...$(NC)\n"
	mkdir -p $(BIN_DIR)
	cd $(GO_SRC) && go build $(GO_FLAGS) $(GO_LDFLAGS) -o ../$(BIN_DIR)/go-api ./cmd/api
	@printf "$(GREEN)✓ Go service built: $(BIN_DIR)/go-api$(NC)\n"

build-node: ## Build Node.js service
	@printf "$(YELLOW)Building Node.js service...$(NC)\n"
	mkdir -p $(DIST_DIR)/node
	cd $(NODE_SRC) && npm run build
	cp -r $(NODE_SRC)/dist/* $(DIST_DIR)/node/
	@printf "$(GREEN)✓ Node.js service built: $(DIST_DIR)/node$(NC)\n"

build-python: ## Build Python service (package)
	@printf "$(YELLOW)Building Python service...$(NC)\n"
	mkdir -p $(DIST_DIR)/python
	cd $(PYTHON_SRC) && python setup.py sdist -d ../$(DIST_DIR)/python
	@printf "$(GREEN)✓ Python service built: $(DIST_DIR)/python$(NC)\n"

build-all: build ## Alias for build

# ============================================================================
# Clean
# ============================================================================

clean: clean-go clean-node clean-python clean-docker ## Clean all build artifacts

clean-go: ## Clean Go artifacts
	@printf "$(YELLOW)Cleaning Go artifacts...$(NC)\n"
	cd $(GO_SRC) && go clean -cache -testcache -modcache
	rm -f $(BIN_DIR)/go-api
	rm -rf $(COVERAGE_DIR)/go
	@printf "$(GREEN)✓ Go artifacts cleaned$(NC)\n"

clean-node: ## Clean Node.js artifacts
	@printf "$(YELLOW)Cleaning Node.js artifacts...$(NC)\n"
	cd $(NODE_SRC) && rm -rf node_modules dist
	rm -rf $(DIST_DIR)/node
	rm -rf $(COVERAGE_DIR)/node
	@printf "$(GREEN)✓ Node.js artifacts cleaned$(NC)\n"

clean-python: ## Clean Python artifacts
	@printf "$(YELLOW)Cleaning Python artifacts...$(NC)\n"
	cd $(PYTHON_SRC) && rm -rf __pycache__ *.pyc .pytest_cache .coverage
	cd $(PYTHON_SRC) && rm -rf build dist *.egg-info
	rm -rf $(DIST_DIR)/python
	rm -rf $(COVERAGE_DIR)/python
	@printf "$(GREEN)✓ Python artifacts cleaned$(NC)\n"

clean-docker: ## Clean Docker artifacts
	@printf "$(YELLOW)Cleaning Docker artifacts...$(NC)\n"
	docker system prune -f
	@printf "$(GREEN)✓ Docker artifacts cleaned$(NC)\n"

clean-all: clean ## Alias for clean

# ============================================================================
# Test
# ============================================================================

test: test-go test-node test-python ## Run all tests

test-go: ## Run Go tests
	@printf "$(YELLOW)Running Go tests...$(NC)\n"
	mkdir -p $(COVERAGE_DIR)/go $(TEST_RESULTS_DIR)
	cd $(GO_SRC) && go test $(GO_TEST_FLAGS) -coverprofile=../$(COVERAGE_DIR)/go/coverage.out $(GO_PACKAGES) | \
		tee $(TEST_RESULTS_DIR)/go-test.out
	cd $(GO_SRC) && go tool cover -html=../$(COVERAGE_DIR)/go/coverage.out -o ../$(COVERAGE_DIR)/go/coverage.html
	@printf "$(GREEN)✓ Go tests completed$(NC)\n"
	@printf "  Coverage report: $(COVERAGE_DIR)/go/coverage.html\n"

test-node: ## Run Node.js tests
	@printf "$(YELLOW)Running Node.js tests...$(NC)\n"
	mkdir -p $(COVERAGE_DIR)/node $(TEST_RESULTS_DIR)
	cd $(NODE_SRC) && npm test -- --coverage --coverageDirectory=../$(COVERAGE_DIR)/node | \
		tee $(TEST_RESULTS_DIR)/node-test.out
	@printf "$(GREEN)✓ Node.js tests completed$(NC)\n"
	@printf "  Coverage report: $(COVERAGE_DIR)/node/lcov-report/index.html\n"

test-python: ## Run Python tests
	@printf "$(YELLOW)Running Python tests...$(NC)\n"
	mkdir -p $(COVERAGE_DIR)/python $(TEST_RESULTS_DIR)
	cd $(PYTHON_SRC) && pytest --cov=. --cov-report=html:../$(COVERAGE_DIR)/python \
		--cov-report=xml:../$(COVERAGE_DIR)/python/coverage.xml | \
		tee $(TEST_RESULTS_DIR)/python-test.out
	@printf "$(GREEN)✓ Python tests completed$(NC)\n"
	@printf "  Coverage report: $(COVERAGE_DIR)/python/index.html\n"

test-integration: ## Run integration tests
	@printf "$(YELLOW)Running integration tests...$(NC)\n"
	mkdir -p $(TEST_RESULTS_DIR)/integration
	docker-compose up -d --wait
	cd tests && go test -tags=integration -v ./... | tee ../$(TEST_RESULTS_DIR)/integration/integration.out
	docker-compose down
	@printf "$(GREEN)✓ Integration tests completed$(NC)\n"

test-benchmark: ## Run benchmarks
	@printf "$(YELLOW)Running benchmarks...$(NC)\n"
	mkdir -p $(TEST_RESULTS_DIR)/benchmark
	cd $(GO_SRC) && go test -bench=. -benchmem -run=^$ ./... | tee ../$(TEST_RESULTS_DIR)/benchmark/go-bench.txt
	cd $(NODE_SRC) && npm run benchmark | tee ../$(TEST_RESULTS_DIR)/benchmark/node-bench.txt
	cd $(PYTHON_SRC) && pytest --benchmark-only --benchmark-json=../$(TEST_RESULTS_DIR)/benchmark/python-bench.json
	@printf "$(GREEN)✓ Benchmarks completed$(NC)\n"

test-all: test test-integration test-benchmark ## Run all tests

# ============================================================================
# Lint
# ============================================================================

lint: lint-go lint-node lint-python ## Run all linters

lint-go: ## Run Go linters
	@printf "$(YELLOW)Running Go linters...$(NC)\n"
	cd $(GO_SRC) && golangci-lint run ./...
	cd $(GO_SRC) && go vet ./...
	@printf "$(GREEN)✓ Go linting passed$(NC)\n"

lint-node: ## Run Node.js linters
	@printf "$(YELLOW)Running Node.js linters...$(NC)\n"
	cd $(NODE_SRC) && npm run lint
	cd $(NODE_SRC) && npm run prettier -- --check
	@printf "$(GREEN)✓ Node.js linting passed$(NC)\n"

lint-python: ## Run Python linters
	@printf "$(YELLOW)Running Python linters...$(NC)\n"
	cd $(PYTHON_SRC) && flake8 .
	cd $(PYTHON_SRC) && black --check .
	cd $(PYTHON_SRC) && mypy .
	@printf "$(GREEN)✓ Python linting passed$(NC)\n"

lint-fix: ## Fix linting issues automatically
	@printf "$(YELLOW)Fixing linting issues...$(NC)\n"
	cd $(GO_SRC) && go fmt ./...
	cd $(NODE_SRC) && npm run prettier -- --write
	cd $(PYTHON_SRC) && black .
	@printf "$(GREEN)✓ Linting fixes applied$(NC)\n"

# ============================================================================
# Development
# ============================================================================

dev: dev-go dev-node dev-python ## Run all services in development mode

dev-go: ## Run Go service in development mode
	@printf "$(YELLOW)Starting Go service in development mode...$(NC)\n"
	cd $(GO_SRC) && go run $(GO_FLAGS) ./cmd/api

dev-node: ## Run Node.js service in development mode
	@printf "$(YELLOW)Starting Node.js service in development mode...$(NC)\n"
	cd $(NODE_SRC) && npm run dev

dev-python: ## Run Python service in development mode
	@printf "$(YELLOW)Starting Python service in development mode...$(NC)\n"
	cd $(PYTHON_SRC) && python -m uvicorn main:app --reload

dev-docker: ## Run all services with Docker Compose
	@printf "$(YELLOW)Starting Docker Compose development environment...$(NC)\n"
	export COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME)-dev && \
	docker-compose up --build

# ============================================================================
# Docker
# ============================================================================

docker-build: docker-build-go docker-build-node docker-build-python ## Build all Docker images

docker-build-go: ## Build Go Docker image
	@printf "$(YELLOW)Building Go Docker image...$(NC)\n"
	docker build -t $(DOCKER_NAMESPACE)/go-api:$(DOCKER_TAG) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-f $(GO_SRC)/Dockerfile $(GO_SRC)
	@printf "$(GREEN)✓ Go Docker image built: $(DOCKER_NAMESPACE)/go-api:$(DOCKER_TAG)$(NC)\n"

docker-build-node: ## Build Node.js Docker image
	@printf "$(YELLOW)Building Node.js Docker image...$(NC)\n"
	docker build -t $(DOCKER_NAMESPACE)/node-frontend:$(DOCKER_TAG) \
		--build-arg NODE_VERSION=$(NODE_VERSION) \
		-f $(NODE_SRC)/Dockerfile $(NODE_SRC)
	@printf "$(GREEN)✓ Node.js Docker image built: $(DOCKER_NAMESPACE)/node-frontend:$(DOCKER_TAG)$(NC)\n"

docker-build-python: ## Build Python Docker image
	@printf "$(YELLOW)Building Python Docker image...$(NC)\n"
	docker build -t $(DOCKER_NAMESPACE)/python-processor:$(DOCKER_TAG) \
		--build-arg PYTHON_VERSION=$(PYTHON_VERSION) \
		-f $(PYTHON_SRC)/Dockerfile $(PYTHON_SRC)
	@printf "$(GREEN)✓ Python Docker image built: $(DOCKER_NAMESPACE)/python-processor:$(DOCKER_TAG)$(NC)\n"

docker-push: docker-push-go docker-push-node docker-push-python ## Push all Docker images

docker-push-go: ## Push Go Docker image
	@printf "$(YELLOW)Pushing Go Docker image...$(NC)\n"
	docker push $(DOCKER_NAMESPACE)/go-api:$(DOCKER_TAG)
	@printf "$(GREEN)✓ Go Docker image pushed$(NC)\n"

docker-push-node: ## Push Node.js Docker image
	@printf "$(YELLOW)Pushing Node.js Docker image...$(NC)\n"
	docker push $(DOCKER_NAMESPACE)/node-frontend:$(DOCKER_TAG)
	@printf "$(GREEN)✓ Node.js Docker image pushed$(NC)\n"

docker-push-python: ## Push Python Docker image
	@printf "$(YELLOW)Pushing Python Docker image...$(NC)\n"
	docker push $(DOCKER_NAMESPACE)/python-processor:$(DOCKER_TAG)
	@printf "$(GREEN)✓ Python Docker image pushed$(NC)\n"

docker-up: ## Start all services with Docker Compose
	@printf "$(YELLOW)Starting Docker Compose environment...$(NC)\n"
	docker-compose up -d
	@printf "$(GREEN)✓ Docker Compose environment started$(NC)\n"

docker-down: ## Stop all services with Docker Compose
	@printf "$(YELLOW)Stopping Docker Compose environment...$(NC)\n"
	docker-compose down
	@printf "$(GREEN)✓ Docker Compose environment stopped$(NC)\n"

docker-logs: ## View Docker Compose logs
	docker-compose logs -f

docker-ps: ## List Docker Compose services
	docker-compose ps

# ============================================================================
# Production
# ============================================================================

prod: docker-build docker-push ## Build and push for production

deploy: ## Deploy to production (example with kubectl)
	@printf "$(YELLOW)Deploying to production...$(NC)\n"
	kubectl apply -f k8s/
	kubectl rollout status deployment/go-api
	kubectl rollout status deployment/node-frontend
	kubectl rollout status deployment/python-processor
	@printf "$(GREEN)✓ Deployment completed$(NC)\n"

# ============================================================================
# Version
# ============================================================================

version: ## Show version information
	@printf "$(CYAN)Version:$(NC) $(VERSION)\n"
	@printf "$(CYAN)Commit:$(NC) $(COMMIT)\n"
	@printf "$(CYAN)Build Date:$(NC) $(BUILD_DATE)\n"
	@printf "$(CYAN)Go Version:$(NC) $(GO_VERSION)\n"
	@printf "$(CYAN)Node Version:$(NC) $(NODE_VERSION)\n"
	@printf "$(CYAN)Python Version:$(NC) $(PYTHON_VERSION)\n"

# ============================================================================
# Utility
# ============================================================================

info: ## Display project information
	@printf "$(CYAN)Project Information$(NC)\n"
	@printf "  $(YELLOW)Project Root:$(NC) $(PWD)\n"
	@printf "  $(YELLOW)Go Service:$(NC) $(GO_SRC)\n"
	@printf "  $(YELLOW)Node Service:$(NC) $(NODE_SRC)\n"
	@printf "  $(YELLOW)Python Service:$(NC) $(PYTHON_SRC)\n"
	@printf "  $(YELLOW)Binary Directory:$(NC) $(BIN_DIR)\n"
	@printf "  $(YELLOW)Distribution Directory:$(NC) $(DIST_DIR)\n"
	@printf "  $(YELLOW)Coverage Directory:$(NC) $(COVERAGE_DIR)\n"
	@printf "  $(YELLOW)Test Results Directory:$(NC) $(TEST_RESULTS_DIR)\n"

tree: ## Show directory tree
	@tree -I 'node_modules|__pycache__|*.pyc|.git|dist'

backup: ## Create backup of important files
	@printf "$(YELLOW)Creating backup...$(NC)\n"
	tar -czf backup-$(shell date +%Y%m%d-%H%M%S).tar.gz \
		$(GO_SRC) $(NODE_SRC) $(PYTHON_SRC) \
		docker-compose.yml Makefile README.md
	@printf "$(GREEN)✓ Backup created$(NC)\n"

# ============================================================================
# Quality
# ============================================================================

quality: lint test build ## Run quality checks (lint, test, build)

ci: quality docker-build ## CI pipeline target

# ============================================================================
# Default target
# ============================================================================

.DEFAULT_GOAL := help