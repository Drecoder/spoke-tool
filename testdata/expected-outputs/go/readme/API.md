```markdown
# Go API Reference

This document provides a comprehensive reference for the Go packages and their APIs.

## 📦 Packages

- [calculator](#package-calculator) - Basic arithmetic operations
- [stringutil](#package-stringutil) - String manipulation utilities
- [fileutil](#package-fileutil) - File system operations
- [validator](#package-validator) - Input validation functions
- [converter](#package-converter) - Type conversion utilities

---

## Package `calculator`

Package calculator provides basic arithmetic operations with error handling.

### Installation

```go
import "github.com/example/calculator"
```

### Constants

```go
const (
    // MaxPrecision defines the maximum decimal precision
    MaxPrecision = 10
    
    // DefaultPrecision is the default precision for calculations
    DefaultPrecision = 2
)
```

### Types

#### type Calculator

```go
type Calculator struct {
    precision int
    // contains filtered or unexported fields
}
```

Calculator represents a calculator instance with configurable precision.

##### func New

```go
func New() *Calculator
```

New creates a new Calculator with default settings.

Example:
```go
calc := calculator.New()
```

##### func NewWithPrecision

```go
func NewWithPrecision(precision int) (*Calculator, error)
```

NewWithPrecision creates a new Calculator with the specified precision.

Parameters:
- `precision` - number of decimal places (must be between 0 and MaxPrecision)

Returns:
- `*Calculator` - new calculator instance
- `error` - error if precision is invalid

Example:
```go
calc, err := calculator.NewWithPrecision(5)
if err != nil {
    log.Fatal(err)
}
```

##### func (c *Calculator) Add

```go
func (c *Calculator) Add(a, b float64) float64
```

Add returns the sum of two numbers.

Parameters:
- `a` - first number
- `b` - second number

Returns:
- `float64` - the sum a + b

Example:
```go
result := calc.Add(5.2, 3.1) // returns 8.3
```

##### func (c *Calculator) Subtract

```go
func (c *Calculator) Subtract(a, b float64) float64
```

Subtract returns the difference between two numbers.

Parameters:
- `a` - first number
- `b` - second number

Returns:
- `float64` - the difference a - b

Example:
```go
result := calc.Subtract(10.5, 4.2) // returns 6.3
```

##### func (c *Calculator) Multiply

```go
func (c *Calculator) Multiply(a, b float64) float64
```

Multiply returns the product of two numbers.

Parameters:
- `a` - first number
- `b` - second number

Returns:
- `float64` - the product a * b

Example:
```go
result := calc.Multiply(3.0, 4.5) // returns 13.5
```

##### func (c *Calculator) Divide

```go
func (c *Calculator) Divide(a, b float64) (float64, error)
```

Divide returns the quotient of a divided by b.

Parameters:
- `a` - dividend
- `b` - divisor (must not be zero)

Returns:
- `float64` - the quotient a / b
- `error` - ErrDivisionByZero if b is zero

Example:
```go
result, err := calc.Divide(10.0, 2.0)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result) // 5.0
```

##### func (c *Calculator) Power

```go
func (c *Calculator) Power(base, exponent float64) float64
```

Power returns base raised to the given exponent.

Parameters:
- `base` - the base number
- `exponent` - the exponent

Returns:
- `float64` - base raised to exponent

Example:
```go
result := calc.Power(2.0, 3.0) // returns 8.0
```

##### func (c *Calculator) Sqrt

```go
func (c *Calculator) Sqrt(x float64) (float64, error)
```

Sqrt returns the square root of x.

Parameters:
- `x` - the number (must be non-negative)

Returns:
- `float64` - square root of x
- `error` - ErrNegativeSqrt if x is negative

Example:
```go
result, err := calc.Sqrt(16.0)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result) // 4.0
```

### Errors

```go
var (
    // ErrDivisionByZero is returned when attempting to divide by zero
    ErrDivisionByZero = errors.New("division by zero")
    
    // ErrNegativeSqrt is returned when attempting sqrt of negative number
    ErrNegativeSqrt = errors.New("square root of negative number")
    
    // ErrInvalidPrecision is returned when precision is out of range
    ErrInvalidPrecision = errors.New("invalid precision")
)
```

---

## Package `stringutil`

Package stringutil provides utilities for string manipulation.

### Installation

```go
import "github.com/example/stringutil"
```

### Functions

#### func Reverse

```go
func Reverse(s string) string
```

Reverse returns its argument string reversed rune-wise left to right.

Parameters:
- `s` - the string to reverse

Returns:
- `string` - the reversed string

Example:
```go
result := stringutil.Reverse("hello") // returns "olleh"
```

#### func ToSnakeCase

```go
func ToSnakeCase(s string) string
```

ToSnakeCase converts a string to snake_case.

Parameters:
- `s` - the input string (e.g., "camelCase" or "PascalCase")

Returns:
- `string` - the snake_case version

Example:
```go
result := stringutil.ToSnakeCase("userName") // returns "user_name"
result := stringutil.ToSnakeCase("UserID")   // returns "user_id"
```

#### func ToCamelCase

```go
func ToCamelCase(s string) string
```

ToCamelCase converts a snake_case string to camelCase.

Parameters:
- `s` - the snake_case string

Returns:
- `string` - the camelCase version

Example:
```go
result := stringutil.ToCamelCase("user_name") // returns "userName"
```

#### func ToPascalCase

```go
func ToPascalCase(s string) string
```

ToPascalCase converts a snake_case string to PascalCase.

Parameters:
- `s` - the snake_case string

Returns:
- `string` - the PascalCase version

Example:
```go
result := stringutil.ToPascalCase("user_name") // returns "UserName"
```

#### func Truncate

```go
func Truncate(s string, maxLen int) string
```

Truncate truncates a string to the specified length, adding "..." if truncated.

Parameters:
- `s` - the string to truncate
- `maxLen` - maximum length (must be >= 3)

Returns:
- `string` - truncated string

Example:
```go
result := stringutil.Truncate("This is a long string", 10) // returns "This is..."
```

#### func ContainsAny

```go
func ContainsAny(s string, substrings []string) bool
```

ContainsAny reports whether any of the substrings are present in s.

Parameters:
- `s` - the string to search
- `substrings` - list of substrings to look for

Returns:
- `bool` - true if any substring is found

Example:
```go
found := stringutil.ContainsAny("hello world", []string{"foo", "world"}) // returns true
```

#### func RemoveDuplicates

```go
func RemoveDuplicates(s string) string
```

RemoveDuplicates removes duplicate lines from a string.

Parameters:
- `s` - multi-line string

Returns:
- `string` - string with duplicate lines removed

Example:
```go
input := "line1\nline2\nline1\nline3"
result := stringutil.RemoveDuplicates(input)
// returns "line1\nline2\nline3"
```

---

## Package `fileutil`

Package fileutil provides utilities for file system operations.

### Installation

```go
import "github.com/example/fileutil"
```

### Types

#### type FileInfo

```go
type FileInfo struct {
    Path      string    `json:"path"`
    Name      string    `json:"name"`
    Size      int64     `json:"size"`
    ModTime   time.Time `json:"mod_time"`
    IsDir     bool      `json:"is_dir"`
    Extension string    `json:"extension"`
    Hash      string    `json:"hash,omitempty"`
}
```

FileInfo represents information about a file.

### Functions

#### func Exists

```go
func Exists(path string) bool
```

Exists reports whether the file or directory exists.

Parameters:
- `path` - file path to check

Returns:
- `bool` - true if the path exists

Example:
```go
if fileutil.Exists("config.yaml") {
    fmt.Println("config file exists")
}
```

#### func IsDir

```go
func IsDir(path string) bool
```

IsDir reports whether the path is a directory.

Parameters:
- `path` - path to check

Returns:
- `bool` - true if path is a directory

Example:
```go
if fileutil.IsDir("./docs") {
    fmt.Println("docs is a directory")
}
```

#### func ReadFile

```go
func ReadFile(path string) (string, error)
```

ReadFile reads the entire file and returns its contents as a string.

Parameters:
- `path` - path to the file

Returns:
- `string` - file contents
- `error` - any error encountered

Example:
```go
content, err := fileutil.ReadFile("data.txt")
if err != nil {
    log.Fatal(err)
}
fmt.Println(content)
```

#### func WriteFile

```go
func WriteFile(path string, data string) error
```

WriteFile writes data to a file, creating directories if needed.

Parameters:
- `path` - path to the file
- `data` - data to write

Returns:
- `error` - any error encountered

Example:
```go
err := fileutil.WriteFile("output.txt", "Hello, World!")
if err != nil {
    log.Fatal(err)
}
```

#### func CopyFile

```go
func CopyFile(src, dst string) error
```

CopyFile copies a file from src to dst.

Parameters:
- `src` - source file path
- `dst` - destination file path

Returns:
- `error` - any error encountered

Example:
```go
err := fileutil.CopyFile("source.txt", "backup/source.txt")
if err != nil {
    log.Fatal(err)
}
```

#### func ListFiles

```go
func ListFiles(dir string, pattern string) ([]string, error)
```

ListFiles returns all files in a directory matching the pattern.

Parameters:
- `dir` - directory to search
- `pattern` - file pattern (e.g., "*.go")

Returns:
- `[]string` - list of matching files
- `error` - any error encountered

Example:
```go
files, err := fileutil.ListFiles("./src", "*.go")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d Go files\n", len(files))
```

#### func GetFileInfo

```go
func GetFileInfo(path string) (*FileInfo, error)
```

GetFileInfo returns detailed information about a file.

Parameters:
- `path` - path to the file

Returns:
- `*FileInfo` - file information
- `error` - any error encountered

Example:
```go
info, err := fileutil.GetFileInfo("config.yaml")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Size: %d bytes\n", info.Size)
fmt.Printf("Modified: %v\n", info.ModTime)
```

---

## Package `validator`

Package validator provides input validation functions.

### Installation

```go
import "github.com/example/validator"
```

### Functions

#### func ValidateEmail

```go
func ValidateEmail(email string) error
```

ValidateEmail checks if an email address is valid.

Parameters:
- `email` - email address to validate

Returns:
- `error` - nil if valid, error describing the issue otherwise

Example:
```go
err := validator.ValidateEmail("user@example.com")
if err != nil {
    fmt.Printf("Invalid email: %v\n", err)
}
```

#### func ValidatePhone

```go
func ValidatePhone(phone string) error
```

ValidatePhone checks if a phone number is valid (US format).

Parameters:
- `phone` - phone number to validate

Returns:
- `error` - nil if valid, error describing the issue otherwise

Example:
```go
err := validator.ValidatePhone("555-123-4567")
if err != nil {
    fmt.Printf("Invalid phone: %v\n", err)
}
```

#### func ValidateURL

```go
func ValidateURL(url string) error
```

ValidateURL checks if a URL is valid.

Parameters:
- `url` - URL to validate

Returns:
- `error` - nil if valid, error describing the issue otherwise

Example:
```go
err := validator.ValidateURL("https://example.com")
if err != nil {
    fmt.Printf("Invalid URL: %v\n", err)
}
```

#### func ValidateRange

```go
func ValidateRange(value, min, max int) error
```

ValidateRange checks if a value is within the specified range.

Parameters:
- `value` - value to check
- `min` - minimum allowed value (inclusive)
- `max` - maximum allowed value (inclusive)

Returns:
- `error` - nil if within range, error otherwise

Example:
```go
err := validator.ValidateRange(42, 1, 100)
if err != nil {
    fmt.Printf("Value out of range: %v\n", err)
}
```

#### func ValidateRequired

```go
func ValidateRequired(value string) error
```

ValidateRequired checks if a string is not empty.

Parameters:
- `value` - string to check

Returns:
- `error` - nil if non-empty, error otherwise

Example:
```go
err := validator.ValidateRequired("")
if err != nil {
    fmt.Println("Field is required") // This will print
}
```

#### func ValidateLength

```go
func ValidateLength(value string, min, max int) error
```

ValidateLength checks if a string length is within bounds.

Parameters:
- `value` - string to check
- `min` - minimum length (inclusive)
- `max` - maximum length (inclusive)

Returns:
- `error` - nil if length is valid, error otherwise

Example:
```go
err := validator.ValidateLength("hello", 3, 10)
if err != nil {
    fmt.Printf("Invalid length: %v\n", err)
}
```

#### func ValidateMatch

```go
func ValidateMatch(value string, pattern *regexp.Regexp) error
```

ValidateMatch checks if a string matches a regex pattern.

Parameters:
- `value` - string to check
- `pattern` - compiled regex pattern

Returns:
- `error` - nil if matches, error otherwise

Example:
```go
pattern := regexp.MustCompile(`^[A-Z][a-z]+$`)
err := validator.ValidateMatch("Hello", pattern)
if err != nil {
    fmt.Printf("Does not match pattern: %v\n", err)
}
```

---

## Package `converter`

Package converter provides type conversion utilities.

### Installation

```go
import "github.com/example/converter"
```

### Functions

#### func ToInt

```go
func ToInt(s string) (int, error)
```

ToInt converts a string to an integer.

Parameters:
- `s` - string to convert

Returns:
- `int` - converted value
- `error` - conversion error

Example:
```go
val, err := converter.ToInt("42")
if err != nil {
    log.Fatal(err)
}
fmt.Println(val) // 42
```

#### func ToFloat

```go
func ToFloat(s string) (float64, error)
```

ToFloat converts a string to a float64.

Parameters:
- `s` - string to convert

Returns:
- `float64` - converted value
- `error` - conversion error

Example:
```go
val, err := converter.ToFloat("3.14")
if err != nil {
    log.Fatal(err)
}
fmt.Println(val) // 3.14
```

#### func ToBool

```go
func ToBool(s string) (bool, error)
```

ToBool converts a string to a boolean.

Parameters:
- `s` - string to convert (accepts "true", "false", "1", "0", "yes", "no")

Returns:
- `bool` - converted value
- `error` - conversion error

Example:
```go
val, err := converter.ToBool("true")
if err != nil {
    log.Fatal(err)
}
fmt.Println(val) // true
```

#### func ToString

```go
func ToString(v interface{}) string
```

ToString converts any value to its string representation.

Parameters:
- `v` - value to convert

Returns:
- `string` - string representation

Example:
```go
s := converter.ToString(42)      // "42"
s := converter.ToString(3.14)    // "3.14"
s := converter.ToString(true)    // "true"
```

#### func ToJSON

```go
func ToJSON(v interface{}) (string, error)
```

ToJSON converts a value to a JSON string.

Parameters:
- `v` - value to convert

Returns:
- `string` - JSON string
- `error` - marshaling error

Example:
```go
data := map[string]interface{}{
    "name": "John",
    "age": 30,
}
json, err := converter.ToJSON(data)
if err != nil {
    log.Fatal(err)
}
fmt.Println(json) // {"age":30,"name":"John"}
```

#### func FromJSON

```go
func FromJSON(jsonStr string, v interface{}) error
```

FromJSON parses a JSON string into the provided value.

Parameters:
- `jsonStr` - JSON string to parse
- `v` - pointer to destination value

Returns:
- `error` - unmarshaling error

Example:
```go
var data map[string]interface{}
err := converter.FromJSON(`{"name":"John","age":30}`, &data)
if err != nil {
    log.Fatal(err)
}
fmt.Println(data["name"]) // John
```

#### func ToBase64

```go
func ToBase64(data []byte) string
```

ToBase64 encodes bytes to a Base64 string.

Parameters:
- `data` - bytes to encode

Returns:
- `string` - Base64 encoded string

Example:
```go
encoded := converter.ToBase64([]byte("hello"))
fmt.Println(encoded) // aGVsbG8=
```

#### func FromBase64

```go
func FromBase64(s string) ([]byte, error)
```

FromBase64 decodes a Base64 string to bytes.

Parameters:
- `s` - Base64 string

Returns:
- `[]byte` - decoded bytes
- `error` - decoding error

Example:
```go
decoded, err := converter.FromBase64("aGVsbG8=")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(decoded)) // hello
```

---

## 📊 Type Index

| Type | Package | Description |
|------|---------|-------------|
| `Calculator` | calculator | Arithmetic calculator |
| `FileInfo` | fileutil | File metadata |
| `ValidationError` | validator | Validation error |

## 🔧 Error Types

| Error | Package | Description |
|-------|---------|-------------|
| `ErrDivisionByZero` | calculator | Division by zero |
| `ErrNegativeSqrt` | calculator | Square root of negative |
| `ErrInvalidPrecision` | calculator | Invalid precision |
| `ErrNotFound` | fileutil | File not found |
| `ErrPermission` | fileutil | Permission denied |

## 📈 Performance Characteristics

| Function | Time Complexity | Space Complexity |
|----------|----------------|------------------|
| `calculator.Add` | O(1) | O(1) |
| `calculator.Divide` | O(1) | O(1) |
| `stringutil.Reverse` | O(n) | O(n) |
| `fileutil.ReadFile` | O(n) | O(n) |
| `validator.ValidateEmail` | O(n) | O(1) |

## 🧪 Examples

See the [examples](examples/) directory for complete runnable examples.

---

*Last Updated: 2024*
```

## ✅ **What this API documentation provides:**

| Section | Description |
|---------|-------------|
| **Package Overview** | Description of each package's purpose |
| **Installation** | Import statements for each package |
| **Constants** | Package-level constants with descriptions |
| **Types** | Type definitions with field descriptions |
| **Functions** | Complete function signatures with parameters and return values |
| **Methods** | Method documentation with receivers |
| **Errors** | Error variables and their meanings |
| **Examples** | Code examples for each function |
| **Type Index** | Quick reference of types by package |
| **Performance** | Time and space complexity notes |

This file serves as:
1. **Test data** for validating the readme generator's API documentation output
2. **Reference** for proper godoc-style formatting
3. **Example** of comprehensive API documentation
4. **Validation** that the generator produces complete API references