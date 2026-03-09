```markdown
# Go Calculator Library

A lightweight, high-performance calculator library written in Go with comprehensive error handling and precision control.

## Features

- ✅ Basic arithmetic operations (add, subtract, multiply, divide)
- ✅ Power and square root functions
- ✅ Configurable decimal precision
- ✅ Comprehensive error handling
- ✅ Thread-safe operations
- ✅ Zero external dependencies
- ✅ 95%+ test coverage
- ✅ Go 1.21+ compatible

## Installation

```bash
go get github.com/example/calculator
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/example/calculator"
)

func main() {
    // Create a new calculator with default precision
    calc := calculator.New()
    
    // Basic operations
    sum := calc.Add(5.2, 3.1)
    fmt.Printf("5.2 + 3.1 = %.2f\n", sum) // 8.30
    
    difference := calc.Subtract(10.5, 4.2)
    fmt.Printf("10.5 - 4.2 = %.2f\n", difference) // 6.30
    
    product := calc.Multiply(3.0, 4.5)
    fmt.Printf("3.0 * 4.5 = %.2f\n", product) // 13.50
    
    // Division with error handling
    quotient, err := calc.Divide(10.0, 2.0)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("10.0 / 2.0 = %.2f\n", quotient) // 5.00
    
    // Power and square root
    power := calc.Power(2.0, 3.0)
    fmt.Printf("2.0^3.0 = %.2f\n", power) // 8.00
    
    sqrt, err := calc.Sqrt(16.0)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("sqrt(16.0) = %.2f\n", sqrt) // 4.00
}
```

## Usage Examples

### Calculator with Custom Precision

```go
// Create calculator with 5 decimal places precision
calc, err := calculator.NewWithPrecision(5)
if err != nil {
    log.Fatal(err)
}

result := calc.Divide(22.0, 7.0)
fmt.Printf("22/7 = %.5f\n", result) // 3.14286
```

### Batch Operations

```go
numbers := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
calc := calculator.New()

// Calculate sum
sum := 0.0
for _, n := range numbers {
    sum = calc.Add(sum, n)
}
fmt.Printf("Sum: %.2f\n", sum) // 16.50

// Calculate product
product := 1.0
for _, n := range numbers {
    product = calc.Multiply(product, n)
}
fmt.Printf("Product: %.2f\n", product) // 946.00
```

### Error Handling Examples

```go
calc := calculator.New()

// Division by zero
_, err := calc.Divide(10.0, 0.0)
if err != nil {
    fmt.Println(err) // division by zero
}

// Square root of negative number
_, err = calc.Sqrt(-4.0)
if err != nil {
    fmt.Println(err) // square root of negative number
}

// Invalid precision
_, err = calculator.NewWithPrecision(20)
if err != nil {
    fmt.Println(err) // invalid precision
}
```

## API Reference

### Types

#### `type Calculator`

```go
type Calculator struct {
    // contains filtered or unexported fields
}
```

Calculator represents a calculator instance with configurable precision.

### Functions

#### `func New() *Calculator`

Creates a new Calculator with default precision (2 decimal places).

#### `func NewWithPrecision(precision int) (*Calculator, error)`

Creates a new Calculator with the specified precision (0-10).

### Methods

#### `func (c *Calculator) Add(a, b float64) float64`

Returns the sum of two numbers.

#### `func (c *Calculator) Subtract(a, b float64) float64`

Returns the difference between two numbers.

#### `func (c *Calculator) Multiply(a, b float64) float64`

Returns the product of two numbers.

#### `func (c *Calculator) Divide(a, b float64) (float64, error)`

Returns the quotient of a divided by b. Returns error if b is zero.

#### `func (c *Calculator) Power(base, exponent float64) float64`

Returns base raised to the given exponent.

#### `func (c *Calculator) Sqrt(x float64) (float64, error)`

Returns the square root of x. Returns error if x is negative.

### Errors

```go
var (
    ErrDivisionByZero   = errors.New("division by zero")
    ErrNegativeSqrt     = errors.New("square root of negative number")
    ErrInvalidPrecision = errors.New("invalid precision")
)
```

## Performance

| Operation | Time Complexity | Notes |
|-----------|----------------|-------|
| Add | O(1) | Constant time |
| Subtract | O(1) | Constant time |
| Multiply | O(1) | Constant time |
| Divide | O(1) | Constant time |
| Power | O(log n) | Binary exponentiation |
| Sqrt | O(1) | Uses math.Sqrt |

Benchmarks (ran on Intel i7-1165G7 @ 2.80GHz):

```bash
BenchmarkAdd-8         	1000000000	         0.32 ns/op
BenchmarkDivide-8      	1000000000	         0.48 ns/op
BenchmarkPower-8       	50000000	        25.3 ns/op
BenchmarkSqrt-8        	1000000000	         0.51 ns/op
```

## Testing

Run the test suite:

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. ./...

# Run with race detection
go test -race ./...
```

## Configuration

The calculator supports configuration via environment variables:

```bash
# Set default precision (0-10)
export CALCULATOR_PRECISION=4

# Enable debug logging
export CALCULATOR_DEBUG=true
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/example/calculator/tags).

## Authors

- **Jane Smith** - *Initial work* - [@janesmith](https://github.com/janesmith)
- **John Doe** - *Documentation* - [@johndoe](https://github.com/johndoe)

See also the list of [contributors](https://github.com/example/calculator/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- The Go team for the excellent standard library
- All contributors who have helped shape this project
- Users who have provided valuable feedback

## Support

- 📖 [Documentation](https://pkg.go.dev/github.com/example/calculator)
- 🐛 [Issue Tracker](https://github.com/example/calculator/issues)
- 💬 [Discussion Forum](https://github.com/example/calculator/discussions)

---

## Changelog

### [1.2.0] - 2024-03-15
- Added Power function
- Improved error messages
- Performance optimizations

### [1.1.0] - 2024-02-01
- Added configurable precision
- Added Sqrt function
- Enhanced test coverage

### [1.0.0] - 2024-01-15
- Initial release
- Basic arithmetic operations
- Error handling

---

*Built with ❤️ using Go*
```

## ✅ **What this README provides:**

| Section | Purpose |
|---------|---------|
| **Title & Features** | Quick overview of what the library does |
| **Installation** | Simple go get command |
| **Quick Start** | Complete working example to get started |
| **Usage Examples** | Common use cases with code |
| **API Reference** | Brief function documentation |
| **Performance** | Benchmarks and complexity analysis |
| **Testing** | How to run tests |
| **Contributing** | Link to guidelines |
| **Versioning** | Semantic versioning info |
| **License** | MIT license notice |
| **Changelog** | Version history |

This serves as:
1. **Test data** for validating the readme generator
2. **Reference** for proper Go README formatting
3. **Example** of comprehensive documentation
4. **Template** for future project READMEs