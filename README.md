```markdown
# Spoke Tool

A local AI-powered development assistant that automatically generates and maintains unit tests and README documentation.

## 🎯 Overview

Spoke Tool follows a hub-and-spoke architecture to help developers write better code with less manual effort:

- **Test Spoke** - Generates unit tests for untested functions
- **Readme Spoke** - Creates and updates README documentation from code + tests

All processing happens **locally** using SLMs (Small Language Models) via Ollama. No code ever leaves your machine.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        File System                          │
│                    (Code + Tests + Docs)                    │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      Orchestrator Hub                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Dispatcher  │  │   Monitor    │  │    Queue     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐                         │
│  │    Audit     │  │   Squeeze    │                         │
│  └──────────────┘  └──────────────┘                         │
└─────────────────────────────────────────────────────────────┘
                          │
            ┌─────────────┴─────────────┐
            │                           │
            ▼                           ▼
┌───────────────────────┐    ┌───────────────────────┐
│    Test Spoke         │    │   Readme Spoke        │
│  ┌─────────────────┐  │    │ ┌─────────────────┐   │
│  │   Analyzer      │  │    │ │   Extractor     │   │
│  └─────────────────┘  │    │ └─────────────────┘   │
│  ┌─────────────────┐  │    │ ┌─────────────────┐   │
│  │   Generator     │  │    │ │   Summarizer    │   │
│  └─────────────────┘  │    │ └─────────────────┘   │
│  ┌─────────────────┐  │    │ ┌─────────────────┐   │
│  │    Runner       │  │    │ │   Formatter     │   │
│  └─────────────────┘  │    │ └─────────────────┘   │
│  ┌─────────────────┐  │    │ ┌─────────────────┐   │
│  │  Interpreter    │  │    │ │    Merger       │   │
│  └─────────────────┘  │    │ └─────────────────┘   │
└───────────┬───────────┘    └───────────┬───────────┘
            │                           │
            └─────────────┬─────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       SLM Pool                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  CodeBERT    │  │   Gemma 2B   │  │ DeepSeek 7B  │      │
│  │  (Encoder)   │  │   (Fast)     │  │ (Reasoning)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                    ┌──────────────┐                         │
│                    │   Ollama     │                         │
│                    │    Local     │                         │
│                    └──────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

## ✨ Key Principles

1. **Local First** - All processing happens on your machine
2. **No Auto-Fixes** - Test failures are reported, never automatically fixed
3. **Audit Trail** - All actions logged for compliance
4. **Multi-Language** - Supports Go, Node.js, Python
5. **Resource Aware** - "Squeeze" mechanism prevents overload
6. **Privacy Preserving** - No code leaves your machine

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Ollama installed and running
- Required models pulled:
  ```bash
  ollama pull codebert
  ollama pull gemma2:2b
  ollama pull deepseek-coder:7b
  ```

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/spoke-tool.git
cd spoke-tool

# Build the tools
make build
# or
.\scripts\build.ps1 build

# Install to GOPATH/bin (optional)
make install
# or
.\scripts\install.ps1
```

### Basic Usage

```bash
# Generate tests for your project
testgen -path /path/to/your/project

# Generate/update README
readmegen -path /path/to/your/project

# Watch mode (auto-update on changes)
testgen -watch -path /path/to/your/project
readmegen -watch -path /path/to/your/project
```

## 📦 Components

### Test Spoke (`testgen`)

The test generator analyzes your code and creates unit tests for untested functions.

```bash
testgen [options]

Options:
  -path string        Path to project root (default ".")
  -config string      Path to config file (default "config.yaml")
  -watch              Watch for changes and auto-generate
  -run                Run tests after generation (default true)
  -coverage           Check coverage after tests
  -threshold float    Coverage threshold (default 80.0)
  -lang string        Specific language (go, nodejs, python)
  -verbose            Verbose output
  -version            Show version
```

**Example:**
```bash
# Generate tests for Go project
testgen -path ./my-go-project -verbose

# Generate tests with coverage check
testgen -path ./my-node-project -coverage -threshold 85

# Watch mode for Python project
testgen -path ./my-python-project -watch
```

### Readme Spoke (`readmegen`)

The README generator creates and updates documentation based on your code and tests.

```bash
readmegen [options]

Options:
  -path string        Path to project root (default ".")
  -config string      Path to config file (default "config.yaml")
  -watch              Watch for changes and auto-update
  -force              Force regenerate all sections
  -verbose            Verbose output
  -version            Show version
```

**Example:**
```bash
# Generate README
readmegen -path ./my-project

# Force regenerate all sections
readmegen -path ./my-project -force

# Watch mode
readmegen -path ./my-project -watch
```

## ⚙️ Configuration

Create a `config.yaml` file in your project root:

```yaml
# config.yaml
models:
  encoder: "codebert"
  decoder: "deepseek-coder:7b"
  fast: "gemma2:2b"

test_spoke:
  enabled: true
  auto_run: true
  coverage_threshold: 80
  frameworks:
    go: "testing"
    nodejs: "jest"
    python: "pytest"

readme_spoke:
  enabled: true
  auto_update: true
  sections:
    - title
    - installation
    - quickstart
    - api
    - examples
    - contributing
    - license

squeeze:
  max_cpu_percent: 80
  max_memory_mb: 4096
  idle_threshold_ms: 500

audit:
  enabled: true
  path: "audit.log"
```

## 🔧 Development

### Project Structure

```
spoke-tool/
├── cmd/                    # Entry points
│   ├── testgen/           # Test generator CLI
│   └── readmegen/         # README generator CLI
├── internal/               # Private packages
│   ├── model/             # SLM client
│   ├── test/              # Test generation
│   │   ├── analyzer.go    # Find untested functions
│   │   ├── generator.go   # Generate tests
│   │   ├── runner.go      # Run tests
│   │   └── interpreter.go # Explain failures (NO fixes)
│   ├── doc/               # Doc generation
│   │   ├── extractor.go   # Extract from code/tests
│   │   ├── summarizer.go  # Generate summaries
│   │   ├── formatter.go   # Format markdown
│   │   └── updater.go     # Merge with existing
│   ├── config/            # Configuration
│   └── common/            # Shared utilities
├── pkg/                    # Reusable packages
├── api/                    # Public types
├── scripts/                # Build scripts
├── testdata/               # Test fixtures
└── docs/                   # Documentation
```

### Building

```bash
# Build all tools
make build

# Build specific tool
make build-testgen
make build-readmegen

# Cross-compile for Windows
make build-windows

# Clean artifacts
make clean
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench

# Run integration tests
make test-integration
```

### Adding a New Language

1. Add language to `api/types/code.go`:
   ```go
   const (
       Go     Language = "go"
       NodeJS Language = "nodejs"
       Python Language = "python"
       Rust   Language = "rust"  // New language
   )
   ```

2. Add test framework mapping in `internal/test/generator.go`

3. Add documentation format in `internal/doc/formatter.go`

4. Add extractor methods in `internal/doc/extractor.go`

## 📊 Supported Languages

| Language | Test Framework | Test Pattern | Doc Format |
|----------|---------------|--------------|------------|
| **Go** | `testing` | `*_test.go` | Godoc |
| **Node.js** | Jest | `*.test.js` | JSDoc |
| **Python** | pytest | `test_*.py` | PyDoc |

## 🧪 How It Works

### Test Generation Flow

```
1. Code Change Detected
   ↓
2. Analyze Changed Files
   ↓
3. Identify Functions Without Tests
   ↓
4. Generate Tests (DeepSeek 7B)
   ↓
5. Write Test Files
   ↓
6. Run Tests
   ↓
7. Tests Pass? → Yes → Trigger Readme Update
         ↓ No
    Report Failures with Analysis
         ↓
    [STOP] - Developer Fixes Code/Tests
```

### README Generation Flow

```
1. Code/Tests Stable
   ↓
2. Extract Examples from Tests
   ↓
3. Analyze API Signatures
   ↓
4. Generate Documentation (Gemma 2B)
   ↓
5. Assemble README Sections
   ↓
6. Merge with Existing README
   ↓
7. Write/Update README.md
```

## 🔒 Privacy & Security

- **All processing is local** - No code sent to external APIs
- **SLMs run via Ollama** - Models run on your machine
- **Audit logging** - All actions logged for compliance
- **No auto-fixes** - Test failures require developer action
- **Manual content preserved** - Never overwrites hand-written docs

## ⚡ Performance: The Squeeze Mechanism

The Squeeze mechanism dynamically adjusts concurrency based on system load:

```go
if cpu > 80% || memory > 4GB {
    // Throttle generation
    reduceConcurrency()
} else {
    // Run at full speed
    increaseConcurrency()
}
```

## 📝 Examples

### Generating Tests for Go

```bash
cd my-go-project
testgen -path . -verbose
```

This will:
- Find all untested exported functions
- Generate `*_test.go` files with table-driven tests
- Run the tests and report results

### Generating README for Python

```bash
cd my-python-project
readmegen -path . -force
```

This will:
- Extract examples from `test_*.py` files
- Generate API documentation from docstrings
- Create/update README.md with installation, API reference, and examples

## 🐛 Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| `connection refused` | Ensure Ollama is running: `ollama serve` |
| `model not found` | Pull required models: `ollama pull codebert` |
| `no functions found` | Check language detection or file extensions |
| `tests fail to run` | Verify test framework is installed |

### Debug Mode

```bash
testgen -path . -verbose
readmegen -path . -verbose
```

## 📚 Documentation

- [Architecture Overview](docs/design/architecture.md)
- [Test Spoke Design](docs/design/test-spoke.md)
- [Readme Spoke Design](docs/design/readme-spoke.md)
- [API Reference](docs/api/README.md)

## 🤝 Contributing

Internal use only - see internal documentation for contribution guidelines.

## 📄 License

Private - All rights reserved.

## 👥 Authors

Internal team - see internal documentation.

---

**Built with ❤️ for internal use**
```

## ✅ **What this README provides:**

| Section | Purpose |
|---------|---------|
| **Overview** | What the tool does and its architecture |
| **Key Principles** | Design philosophy (no auto-fixes, local-first) |
| **Quick Start** | Get running in minutes |
| **Components** | Detailed usage of testgen and readmegen |
| **Configuration** | Config file reference |
| **Development** | Building, testing, extending |
| **How It Works** | Flow diagrams for both spokes |
| **Privacy & Security** | Why it's safe to use |
| **Performance** | Squeeze mechanism explanation |
| **Examples** | Real-world usage examples |
| **Troubleshooting** | Common issues and solutions |

## 🎯 **Key Features:**

- ✅ Clear explanation of the **no auto-fixes** principle
- ✅ **Architecture diagram** showing all components
- ✅ **Quick start** for new users
- ✅ **Detailed command reference**
- ✅ **Configuration guide**
- ✅ **Development instructions**
- ✅ **Privacy & security** focus
- ✅ **Troubleshooting** section# spoke-tool
