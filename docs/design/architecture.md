# Spoke Tool Architecture

## 🎯 Overview

The Spoke Tool is a local AI-powered development assistant that automatically generates and maintains unit tests and README documentation. It follows a hub-and-spoke architecture with a focus on privacy, performance, and developer workflow.

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
└───────────┬───────────┘    └───────────┬───────────┘
            │                           │
            └─────────────┬─────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       SLM Pool                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  codellama:7b    │  │   Gemma 2B   │  │ DeepSeek 7B  │      │
│  │  (Encoder)   │  │   (Fast)     │  │ (Reasoning)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                    ┌──────────────┐                         │
│                    │   Ollama     │                         │
│                    │    Local     │                         │
│                    └──────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

## 🏗️ Core Components

### 1. Orchestrator Hub
The central coordinator that manages all spokes and system resources.

| Component | Responsibility |
|-----------|----------------|
| **Dispatcher** | Routes file change events to appropriate spokes |
| **Monitor** | Tracks CPU, memory, and system load for Squeeze mechanism |
| **Queue** | Manages event processing order and priorities |
| **Audit** | Logs all actions for compliance (SOC2, GDPR) |
| **Squeeze** | Dynamically adjusts concurrency based on system load |

### 2. Test Spoke
Generates and maintains unit tests for code changes.

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Analyze    │ -> │  Generate   │ -> │   Run       │
│  Code Gaps  │    │   Tests     │    │   Tests     │
└─────────────┘    └─────────────┘    └──────┬──────┘
                                             │
                                    ┌────────┴────────┐
                                    ▼                 ▼
                            ┌─────────────┐    ┌─────────────┐
                            │   Pass      │    │   Fail      │
                            │  Update     │    │  Report     │
                            │   Docs      │    │  Analysis   │
                            └─────────────┘    └─────────────┘
```

**Key Features:**
- Multi-language support (Go, Node.js, Python)
- Test gap analysis
- Test generation using SLMs
- Test execution and result collection
- Failure analysis (reporting only, no auto-fixes)
- Coverage tracking

### 3. Readme Spoke
Creates and updates README documentation based on code and tests.

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Extract    │ -> │  Summarize  │ -> │  Generate   │
│  Examples   │    │     API     │    │   README    │
│  from Tests │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                                                    │
                                                    ▼
                                           ┌─────────────┐
                                           │   Write     │
                                           │   to Disk   │
                                           └─────────────┘
```

**Key Features:**
- Example extraction from test files
- API documentation generation
- Multi-language support
- Section-based README assembly
- Preservation of manual content

### 4. SLM Pool
Local language models running via Ollama.

| Model | Purpose | Use Cases |
|-------|---------|-----------|
| **codellama:7b** | Code understanding | Function analysis, test gap detection |
| **Gemma 2B** | Fast generation | Documentation, simple examples |
| **DeepSeek 7B** | Complex reasoning | Test generation, failure analysis |

## 🔄 Data Flow

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
    Developer Fixes Code/Tests
```

### Readme Generation Flow
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
6. Write/Update README.md
```

## 🧠 SLM Interaction Patterns

### Code Understanding (codellama:7b)
```go
// Input: Function code
// Output: Structured analysis
type CodeAnalysis struct {
    Functions   []Function
    Imports     []string
    Complexity  int
    TestGaps    []string
}
```

### Test Generation (DeepSeek 7B)
```go
// Input: Function + Context
// Output: Test code with framework-specific syntax
func GenerateTests(function Function) (string, error) {
    // Language-specific test generation
    // Go: testing package
    // Node: Jest
    // Python: pytest
}
```

### Documentation (Gemma 2B)
```go
// Input: Function + Tests
// Output: Documentation with examples
func GenerateDocs(function Function, tests []Test) (string, error) {
    // Extract working examples from tests
    // Generate API documentation
    // Format in language-standard style
}
```

## 🔒 Privacy & Security

### Local-First Design
- All code analysis happens locally
- No code sent to external servers
- SLMs run locally via Ollama

### Audit Trail
```
Event {
    ID:        uuid,
    Type:      "test_generated",
    Source:    "test-spoke",
    Data:      { function: "Add", file: "math.go" },
    Timestamp: time.Now(),
    User:      "developer@example.com"
}
```

### Data Flow Control
```
Code Changes ──► Local Analysis ──► Local SLM ──► Test Generation
      │                                               │
      └────────── No external API calls ──────────────┘
```

## ⚡ Performance: The Squeeze Mechanism

Dynamically adjusts concurrency based on system load:

```go
type SqueezeConfig struct {
    MaxCPUPercent  int  // Back off when CPU > 80%
    MaxMemoryMB    int  // Back off when memory > 4GB
    IdleThreshold  int  // Resume when idle for 500ms
}

func (s *Squeeze) ShouldThrottle() bool {
    cpu := getCPUUsage()
    mem := getMemoryUsage()
    
    if cpu > s.MaxCPUPercent || mem > s.MaxMemoryMB {
        return true // Throttle generation
    }
    return false
}
```

## 📊 Language Support Matrix

| Feature | Go | Node.js | Python |
|---------|-----|---------|--------|
| Test Framework | `testing` | Jest | pytest |
| Test File Pattern | `*_test.go` | `*.test.js` | `test_*.py` |
| Doc Format | Godoc | JSDoc | PyDoc |
| Coverage Tool | `go test -cover` | `jest --coverage` | `pytest-cov` |
| Mock Support | Interfaces | Jest mocks | `unittest.mock` |

## 🚦 Error Handling Strategy

### Test Failures (Report Only, No Auto-Fix)
```
Test Failure Detected
    ↓
Analyze Failure (DeepSeek 7B)
    ↓
Generate Explanation
    ↓
Report to Developer
    ↓
[STOP] - Developer Action Required
```

### System Errors
```
Error Types:
├── Configuration Error → Exit with code 1
├── Model Error        → Retry with backoff
├── Analysis Error     → Log and skip file
└── Write Error        → Exit with code 2
```

## 📁 Project Structure

```
spoke-tool/
├── cmd/                    # Entry points
│   ├── readmegen/         # README generator CLI
│   └── testgen/           # Test generator CLI
├── internal/               # Private packages
│   ├── model/             # SLM client
│   ├── test/              # Test generation
│   ├── doc/               # Doc generation
│   └── config/            # Configuration
├── api/                    # Public types
│   └── types/             # Shared data structures
├── pkg/                    # Reusable packages
└── docs/                   # Documentation
```

## 🔍 Configuration

```yaml
# config.yaml example
models:
  encoder: "codellama:7b"
  decoder: "codellama:7b"
  fast: "gemma2:2b"

test_spoke:
  enabled: true
  auto_run: true
  coverage_threshold: 80

readme_spoke:
  enabled: true
  auto_update: true
  sections:
    - installation
    - quick-start
    - api
    - examples

squeeze:
  max_cpu_percent: 80
  max_memory_mb: 4096
  idle_threshold_ms: 500
```

## 🧪 Development Workflow

1. **Developer writes code**
2. **Test spoke detects changes**
3. **Tests generated for new functions**
4. **Tests run automatically**
5. **If tests pass → README updates**
6. **If tests fail → Developer gets analysis**

## 🎯 Design Principles

1. **Local First** - All processing happens on developer machine
2. **No Auto-Fixes** - Tests failures require developer action
3. **Audit Trail** - All actions logged for compliance
4. **Multi-Language** - Support primary development languages
5. **Resource Aware** - Squeeze mechanism prevents overload
6. **Privacy Preserving** - No code leaves the machine

## 📈 Future Extensions

- Additional languages (Rust, Java, C#)
- Custom prompt templates
- Plugin system for custom spokes
- IDE integrations (VS Code, JetBrains)
- CI/CD integration

## 📚 References

- [Ollama Documentation](https://ollama.ai/)
- [codellama:7b Paper](https://arxiv.org/abs/2002.08155)
- [DeepSeek Coder](https://github.com/deepseek-ai/DeepSeek-Coder)
- [Gemma Models](https://ai.google.dev/gemma)