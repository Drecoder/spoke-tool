```markdown
# Edge Cases Test Suite

This directory contains comprehensive edge case tests for validating the test generator's ability to detect and handle boundary conditions, error scenarios, and unusual inputs across multiple programming languages.

## 📋 Overview

The edge cases test suite covers:

- **Zero and empty values** - Testing boundary conditions at the lower limits
- **Numeric boundaries** - MAX_SAFE_INTEGER, MIN_SAFE_INTEGER, overflow, underflow
- **String edge cases** - Unicode, emoji, control characters, surrogate pairs
- **Collection edge cases** - Sparse arrays, holes, frozen objects, circular references
- **Function edge cases** - Recursion depth, closures, argument limits
- **Async edge cases** - Promise resolution, timeouts, error chains
- **Date/Time edge cases** - Invalid dates, epoch boundaries, DST transitions
- **Regular expression edge cases** - Catastrophic backtracking, lookahead/lookbehind
- **Error handling edge cases** - Nested errors, custom errors, stack traces
- **Stream/Buffer edge cases** - Backpressure, zero-length buffers, large streams

## 📁 File Structure

```
edge-cases/
├── README.md                 # This documentation
├── edge_cases_test.go        # Go edge cases test suite
├── edge_cases.test.js        # Node.js/Jest edge cases test suite
├── test_edge_cases.py        # Python/pytest edge cases test suite
└── expected-outputs/         # Expected outputs for validation
    ├── go/
    ├── nodejs/
    └── python/
```

## 🎯 Purpose

These test files serve three main purposes:

1. **Validate the test generator** - Ensure the generator correctly identifies edge cases
2. **Provide examples** - Show proper edge case testing patterns
3. **Verify output correctness** - Compare generated tests against expected outputs

## 🔍 Edge Cases by Language

### Go Edge Cases (`edge_cases_test.go`)

| Category | Test Functions |
|----------|----------------|
| Zero Values | `TestZeroValues`, `TestZeroValues` |
| Boundary Conditions | `TestBoundaryConditions` |
| Overflow/Underflow | `TestOverflow`, `TestIntegerOverflow` |
| Division by Zero | `TestDivisionByZero` |
| Nil Collections | `TestNilSlices`, `TestNilMaps`, `TestNilInterfaces` |
| Index Out of Bounds | `TestIndexOutOfBounds` |
| String Edge Cases | `TestStringEdgeCases` |
| Recursion | `TestRecursionDepth` |
| Channel Edge Cases | `TestChannelEdgeCases` |
| Race Conditions | `TestRaceConditions` |
| Error Handling | `TestErrorWrapping` |
| Type Assertions | `TestTypeAssertions` |

### Node.js Edge Cases (`edge_cases.test.js`)

| Category | Test Suites |
|----------|-------------|
| Zero/Empty Values | `Zero and Empty Values Edge Cases` |
| Numeric Boundaries | `Numeric Boundary Conditions` |
| Integer Overflow | `Integer Overflow and Underflow` |
| String Edge Cases | `String Edge Cases` |
| Array Edge Cases | `Array Edge Cases` |
| Object Edge Cases | `Object Edge Cases` |
| Function Edge Cases | `Function Edge Cases` |
| Promise/Async | `Promise and Async Edge Cases` |
| Date/Time | `Date and Time Edge Cases` |
| Regular Expressions | `Regular Expression Edge Cases` |
| Error Handling | `Error Handling Edge Cases` |
| Event Emitter | `Event Emitter Edge Cases` |
| Streams | `Stream Edge Cases` |
| Buffers | `Buffer Edge Cases` |

### Python Edge Cases (`test_edge_cases.py`)

| Category | Test Classes |
|----------|--------------|
| Zero/Empty | `TestZeroValues` |
| Numeric Boundaries | `TestNumericBoundaries` |
| String Edge Cases | `TestStringEdgeCases` |
| Collection Edge Cases | `TestListEdgeCases`, `TestDictEdgeCases`, `TestSetEdgeCases` |
| Function Edge Cases | `TestFunctionEdgeCases` |
| Async Edge Cases | `TestAsyncEdgeCases` |
| Exception Edge Cases | `TestExceptionEdgeCases` |
| Context Manager | `TestContextManagerEdgeCases` |
| Decorator Edge Cases | `TestDecoratorEdgeCases` |
| Metaclass Edge Cases | `TestMetaclassEdgeCases` |
| Descriptor Edge Cases | `TestDescriptorEdgeCases` |

## 🧪 Running the Tests

### Go Tests

```bash
# Run all edge cases tests
go test -v ./testdata/expected-outputs/edge-cases/edge_cases_test.go

# Run specific test
go test -v -run TestZeroValues ./testdata/expected-outputs/edge-cases/edge_cases_test.go

# Run with race detection
go test -race -v ./testdata/expected-outputs/edge-cases/edge_cases_test.go

# Run with coverage
go test -cover -v ./testdata/expected-outputs/edge-cases/edge_cases_test.go
```

### Node.js/Jest Tests

```bash
# Run all edge cases tests
npx jest testdata/expected-outputs/edge-cases/edge_cases.test.js

# Run with verbose output
npx jest --verbose testdata/expected-outputs/edge-cases/edge_cases.test.js

# Run specific test suite
npx jest -t "Zero and Empty Values Edge Cases" testdata/expected-outputs/edge-cases/edge_cases.test.js

# Run with coverage
npx jest --coverage testdata/expected-outputs/edge-cases/edge_cases.test.js
```

### Python/pytest Tests

```bash
# Run all edge cases tests
pytest testdata/expected-outputs/edge-cases/test_edge_cases.py -v

# Run specific test class
pytest testdata/expected-outputs/edge-cases/test_edge_cases.py::TestZeroValues -v

# Run with coverage
pytest testdata/expected-outputs/edge-cases/test_edge_cases.py --cov=. -v

# Run with verbose output
pytest testdata/expected-outputs/edge-cases/test_edge_cases.py -vv
```

## 📊 Expected Outputs

The `expected-outputs/` directory contains the expected results for each test suite:

```
expected-outputs/
├── go/
│   ├── zero_values.txt           # Expected output for zero value tests
│   ├── boundaries.txt            # Expected output for boundary tests
│   └── overflow.txt              # Expected output for overflow tests
├── nodejs/
│   ├── zero-values.txt           # Expected Jest output
│   ├── numeric-boundaries.txt    # Expected Jest output
│   └── string-edge-cases.txt     # Expected Jest output
└── python/
    ├── test_zero_values.txt      # Expected pytest output
    ├── test_boundaries.txt       # Expected pytest output
    └── test_string_cases.txt     # Expected pytest output
```

## 🔧 Validation Script

Use the validation script to compare generated tests against expected outputs:

```bash
# From the project root
./testdata/expected-outputs/validate.sh

# Or run directly
cd testdata/expected-outputs
./validate.sh
```

## 📝 Adding New Edge Cases

When adding new edge cases, follow these guidelines:

1. **Add to all three languages** - Maintain parity across language implementations
2. **Document the edge case** - Add comments explaining what's being tested
3. **Include expected behavior** - Clearly state what should happen
4. **Add to validation** - Update expected outputs and validation script

Example template for a new edge case:

```go
// TestNewEdgeCase tests [describe the edge case]
func TestNewEdgeCase(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        wantErr bool
        want    interface{}
    }{
        {
            name:    "edge case description",
            input:   edgeValue,
            wantErr: true,
            want:    expectedResult,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := functionUnderTest(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && result != tt.want {
                t.Errorf("got %v, want %v", result, tt.want)
            }
        })
    }
}
```

## 🎯 Edge Case Categories Explained

### Zero and Empty Values
Testing how functions handle the absence of data:
- `0` (integer zero)
- `0.0` (float zero)
- `""` (empty string)
- `[]` (empty slice/array)
- `{}` (empty object/map)
- `nil`/`null`/`None` (null values)

### Numeric Boundaries
Testing limits of numeric types:
- `MAX_SAFE_INTEGER` and beyond
- `MIN_SAFE_INTEGER` and below
- `Number.MAX_VALUE` / `Number.MIN_VALUE`
- `Infinity` and `-Infinity`
- Floating point precision (0.1 + 0.2)

### String Edge Cases
Testing unusual string content:
- Very long strings (memory limits)
- Unicode characters (multi-byte)
- Emoji and surrogate pairs
- Control characters
- Null bytes
- RTL text
- Zero-width joiners

### Collection Edge Cases
Testing unusual collection states:
- Sparse arrays (missing indices)
- Arrays with holes
- Objects with null prototype
- Frozen/sealed objects
- Circular references
- Objects with symbol keys

### Function Edge Cases
Testing function behavior limits:
- Deep recursion (stack overflow)
- Tail recursion optimization
- Many arguments
- Closures with large captured data

### Async Edge Cases
Testing asynchronous edge cases:
- Never-resolving promises
- Promise chains with errors
- Multiple rejections
- Unhandled rejections
- Race conditions

## 📚 Reference

- [Go Testing Documentation](https://pkg.go.dev/testing)
- [Jest Edge Cases](https://jestjs.io/docs/expect)
- [pytest Edge Cases](https://docs.pytest.org/en/stable/example/parametrize.html)
- [MDN Edge Cases](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER)
- [Python Edge Cases](https://docs.python.org/3/library/exceptions.html#bltin-exceptions)

## ✅ Validation Checklist

Use this checklist when validating edge case handling:

- [ ] Zero values handled correctly
- [ ] Empty collections handled correctly
- [ ] Null/nil/None handled correctly
- [ ] Numeric boundaries respected
- [ ] Overflow/underflow detected
- [ ] String encoding issues handled
- [ ] Unicode support correct
- [ ] Recursion limits enforced
- [ ] Async errors caught
- [ ] Race conditions detected
- [ ] Error wrapping preserved
- [ ] Type assertions safe

## 🐛 Common Issues Found by Edge Cases

| Issue | Description | Languages Affected |
|-------|-------------|-------------------|
| Integer Overflow | Not checking for overflow in addition | Go, Node.js, Python |
| Floating Point | Assuming exact decimal representation | All languages |
| Null Pointer | Dereferencing nil pointers | Go, Node.js |
| Index Out of Bounds | Accessing beyond array length | All languages |
| Stack Overflow | Unbounded recursion | All languages |
| Memory Leak | Closures holding large data | Node.js, Python |
| Race Condition | Concurrent access without locks | Go, Node.js |
| Deadlock | Channel/goroutine deadlocks | Go |
| Unhandled Rejection | Promises with no catch | Node.js |
| Exception Safety | Resources not cleaned up | Python |

---

*Last Updated: 2024*