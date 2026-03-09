cat > Makefile << 'EOF'
.PHONY: build test clean install watch-readme watch-test

VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)"

# Build all tools
build:
	go build $(LDFLAGS) -o bin/readmegen.exe cmd/readmegen/main.go
	go build $(LDFLAGS) -o bin/testgen.exe cmd/testgen/main.go

# Build specific tools
build-readmegen:
	go build $(LDFLAGS) -o bin/readmegen.exe cmd/readmegen/main.go

build-testgen:
	go build $(LDFLAGS) -o bin/testgen.exe cmd/testgen/main.go

# Run tests
test:
	go test ./...

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/readmegen
	go install $(LDFLAGS) ./cmd/testgen

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Watch modes
watch-readme:
	go run cmd/readmegen/main.go -watch -path .

watch-test:
	go run cmd/testgen/main.go -watch -path .

# Run once
run-readme:
	go run cmd/readmegen/main.go $(ARGS)

run-test:
	go run cmd/testgen/main.go $(ARGS)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build           Build all tools"
	@echo "  build-readmegen Build only readmegen"
	@echo "  build-testgen   Build only testgen"
	@echo "  test            Run tests"
	@echo "  install         Install to GOPATH/bin"
	@echo "  clean           Clean build artifacts"
	@echo "  watch-readme    Run readmegen in watch mode"
	@echo "  watch-test      Run testgen in watch mode"
	@echo "  run-readme      Run readmegen once (use ARGS for flags)"
	@echo "  run-test        Run testgen once (use ARGS for flags)"
	@echo ""
	@echo "Examples:"
	@echo "  make run-test ARGS='-path ../myproject -verbose'"
	@echo "  make run-readme ARGS='-path ../myproject -force'"
	@echo "  make build"
EOF