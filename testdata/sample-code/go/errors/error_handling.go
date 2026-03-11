package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// Basic Error Handling
// ============================================================================

// BasicErrorHandling demonstrates simple error checking
func BasicErrorHandling(filename string) {
	// Common pattern: check error immediately
	file, err := os.Open(filename)
	if err != nil {
		// Handle error
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()
	
	// Process file...
	fmt.Println("File opened successfully")
}

// MultipleErrorChecks demonstrates handling multiple error points
func MultipleErrorChecks(filename string) error {
	// First operation
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}
	
	// Second operation
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}
	
	// Third operation
	if value, ok := config["threshold"].(float64); ok {
		if value < 0 {
			return errors.New("threshold cannot be negative")
		}
	}
	
	return nil
}

// ============================================================================
// Error Wrapping and Unwrapping
// ============================================================================

// WrapError demonstrates error wrapping with context
func WrapError() error {
	err := errors.New("original error")
	
	// Add context at each level
	err = fmt.Errorf("reading config: %w", err)
	err = fmt.Errorf("initializing app: %w", err)
	
	return err
}

// UnwrapError demonstrates error unwrapping
func UnwrapError() {
	err := WrapError()
	
	// Print the full error chain
	fmt.Printf("Error: %v\n", err)
	
	// Unwrap step by step
	for err != nil {
		fmt.Printf("Level: %v\n", err)
		err = errors.Unwrap(err)
	}
}

// IsError demonstrates errors.Is for sentinel errors
func IsError() {
	var ErrNotFound = errors.New("not found")
	var ErrPermission = errors.New("permission denied")
	
	err := fmt.Errorf("accessing file: %w", ErrNotFound)
	
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Error is a not-found error")
	}
	
	if errors.Is(err, ErrPermission) {
		fmt.Println("Error is a permission error")
	}
}

// AsError demonstrates errors.As for custom error types
func AsError() {
	type TimeoutError struct {
		Timeout time.Duration
		Message string
	}
	
	func (e TimeoutError) Error() string {
		return fmt.Sprintf("timeout after %v: %s", e.Timeout, e.Message)
	}
	
	err := TimeoutError{Timeout: 5 * time.Second, Message: "connection timed out"}
	wrapped := fmt.Errorf("operation failed: %w", err)
	
	var timeoutErr TimeoutError
	if errors.As(wrapped, &timeoutErr) {
		fmt.Printf("Timeout: %v, Message: %s\n", timeoutErr.Timeout, timeoutErr.Message)
	}
}

// ============================================================================
// Sentinel Errors
// ============================================================================

// Define sentinel errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrInvalidInput = errors.New("invalid input")
	ErrTimeout      = errors.New("operation timed out")
	ErrClosed       = errors.New("connection closed")
)

// SentinelErrorExample demonstrates using sentinel errors
func SentinelErrorExample(id int) error {
	if id <= 0 {
		return ErrInvalidInput
	}
	
	// Simulate not found
	if id == 42 {
		return fmt.Errorf("user %d: %w", id, ErrNotFound)
	}
	
	return nil
}

func HandleSentinelErrors() {
	err := SentinelErrorExample(42)
	
	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Println("Resource not found, creating it...")
	case errors.Is(err, ErrInvalidInput):
		fmt.Println("Invalid input, please check your data")
	case err != nil:
		fmt.Printf("Unexpected error: %v\n", err)
	default:
		fmt.Println("Success!")
	}
}

// ============================================================================
// Custom Error Types
// ============================================================================

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s=%v: %s", e.Field, e.Value, e.Message)
}

// NotFoundError represents a resource not found
type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %v", e.Resource, e.ID)
}

// BusinessError represents a business rule violation
type BusinessError struct {
	Rule    string
	Message string
	Details map[string]interface{}
}

func (e BusinessError) Error() string {
	return fmt.Sprintf("business rule violation [%s]: %s", e.Rule, e.Message)
}

// CustomErrorExample demonstrates using custom error types
func CustomErrorExample() error {
	// Validation error
	if err := validateAge(-5); err != nil {
		return err
	}
	
	// Not found error
	if err := findUser(999); err != nil {
		return err
	}
	
	// Business error
	if err := processOrder(-100); err != nil {
		return err
	}
	
	return nil
}

func validateAge(age int) error {
	if age < 0 {
		return ValidationError{
			Field:   "age",
			Value:   age,
			Message: "age cannot be negative",
		}
	}
	if age > 150 {
		return ValidationError{
			Field:   "age",
			Value:   age,
			Message: "age exceeds maximum",
		}
	}
	return nil
}

func findUser(id int) error {
	// Simulate not found
	if id > 100 {
		return NotFoundError{
			Resource: "user",
			ID:       id,
		}
	}
	return nil
}

func processOrder(amount float64) error {
	if amount < 0 {
		return BusinessError{
			Rule:    "minimum_order",
			Message: "order amount cannot be negative",
			Details: map[string]interface{}{
				"amount":      amount,
				"minimum_allowed": 0,
			},
		}
	}
	return nil
}

// HandleCustomErrors demonstrates type switching on custom errors
func HandleCustomErrors() {
	err := CustomErrorExample()
	
	switch e := err.(type) {
	case nil:
		fmt.Println("Success!")
		
	case ValidationError:
		fmt.Printf("Validation error on field %s: %s\n", e.Field, e.Message)
		
	case NotFoundError:
		fmt.Printf("Not found: %s with ID %v\n", e.Resource, e.ID)
		
	case BusinessError:
		fmt.Printf("Business error: %s (rule: %s)\n", e.Message, e.Rule)
		if len(e.Details) > 0 {
			fmt.Printf("Details: %v\n", e.Details)
		}
		
	default:
		fmt.Printf("Unknown error: %v\n", err)
	}
}

// ============================================================================
// Error Handling with defer
// ============================================================================

// ErrorHandlingWithDefer demonstrates using defer for error handling
func ErrorHandlingWithDefer() (err error) {
	defer func() {
		if err != nil {
			// Log or wrap the error
			fmt.Printf("Function failed with error: %v\n", err)
		}
	}()
	
	// Simulate work that might fail
	if time.Now().Unix()%2 == 0 {
		return errors.New("random failure")
	}
	
	return nil
}

// ============================================================================
// Multiple Return Values
// ============================================================================

// MultiReturnError demonstrates functions with multiple return values
func MultiReturnError() (result int, err error) {
	value, err := strconv.Atoi("42")
	if err != nil {
		return 0, fmt.Errorf("conversion failed: %w", err)
	}
	
	return value * 2, nil
}

// ============================================================================
// Error Aggregation
// ============================================================================

// ErrorList aggregates multiple errors
type ErrorList struct {
	errors []error
}

func (e *ErrorList) Add(err error) {
	if err != nil {
		e.errors = append(e.errors, err)
	}
}

func (e *ErrorList) Error() string {
	if len(e.errors) == 0 {
		return "no errors"
	}
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d errors occurred", len(e.errors)))
	for i, err := range e.errors {
		sb.WriteString(fmt.Sprintf("\n  %d: %v", i+1, err))
	}
	return sb.String()
}

func (e *ErrorList) HasErrors() bool {
	return len(e.errors) > 0
}

// ValidateAll demonstrates error aggregation
func ValidateAll(inputs []string) error {
	errs := &ErrorList{}
	
	for i, input := range inputs {
		if input == "" {
			errs.Add(fmt.Errorf("input %d is empty", i))
		}
		if len(input) > 10 {
			errs.Add(fmt.Errorf("input %d too long: %d > 10", i, len(input)))
		}
	}
	
	if errs.HasErrors() {
		return errs
	}
	return nil
}

// ============================================================================
// Error Handling Patterns
// ============================================================================

// CheckError pattern - simple error check
func CheckError(err error) bool {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return true
	}
	return false
}

// Must pattern - panic on error (use sparingly)
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// IgnoreError pattern - explicitly ignore error (use with caution)
func IgnoreError() {
	// Read a file that might not exist
	data, _ := os.ReadFile("optional.txt")
	if len(data) > 0 {
		fmt.Println("File had content")
	}
}

// ============================================================================
// Retry Pattern
// ============================================================================

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Backoff     float64
}

// Retry executes a function with retries
func Retry(fn func() error, config RetryConfig) error {
	var err error
	delay := config.Delay
	
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		
		if attempt < config.MaxAttempts {
			fmt.Printf("Attempt %d failed: %v. Retrying in %v...\n", attempt, err, delay)
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * config.Backoff)
		}
	}
	
	return fmt.Errorf("all %d attempts failed: %w", config.MaxAttempts, err)
}

// RetryableErrorExample demonstrates the retry pattern
func RetryableErrorExample() error {
	attempt := 0
	return Retry(func() error {
		attempt++
		if attempt < 3 {
			return fmt.Errorf("attempt %d failed", attempt)
		}
		return nil
	}, RetryConfig{
		MaxAttempts: 5,
		Delay:       100 * time.Millisecond,
		Backoff:     2.0,
	})
}

// ============================================================================
// Timeout Pattern
// ============================================================================

// TimeoutError represents a timeout
type TimeoutError struct {
	Operation string
	Timeout   time.Duration
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("%s timed out after %v", e.Operation, e.Timeout)
}

// WithTimeout executes a function with timeout
func WithTimeout(timeout time.Duration, fn func() error) error {
	done := make(chan error)
	
	go func() {
		done <- fn()
	}()
	
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return TimeoutError{
			Operation: "function execution",
			Timeout:   timeout,
		}
	}
}

// ============================================================================
// Context-Based Error Handling
// ============================================================================

// ContextError adds context to errors
type ContextError struct {
	Context string
	Err     error
}

func (e ContextError) Error() string {
	return fmt.Sprintf("%s: %v", e.Context, e.Err)
}

func (e ContextError) Unwrap() error {
	return e.Err
}

// WithContext adds context to an error
func WithContext(err error, context string) error {
	if err == nil {
		return nil
	}
	return ContextError{
		Context: context,
		Err:     err,
	}
}

// ============================================================================
// Error Handling in Goroutines
// ============================================================================

// GoroutineError demonstrates error handling in goroutines
func GoroutineError() error {
	errCh := make(chan error, 3)
	
	// Launch goroutines
	for i := 0; i < 3; i++ {
		go func(id int) {
			if id == 1 {
				errCh <- fmt.Errorf("goroutine %d failed", id)
				return
			}
			errCh <- nil
		}(i)
	}
	
	// Collect errors
	var errs []error
	for i := 0; i < 3; i++ {
		if err := <-errCh; err != nil {
			errs = append(errs, err)
		}
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("%d goroutines failed: %v", len(errs), errs)
	}
	return nil
}

// ============================================================================
// Error Handling with defer and panic/recover
// ============================================================================

// SafeExecute executes a function safely, converting panics to errors
func SafeExecute(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	
	fn()
	return nil
}

// ============================================================================
// File Processing Examples
// ============================================================================

// ProcessFile demonstrates comprehensive error handling
func ProcessFile(filename string) error {
	// Check if file exists
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist: %w", filename, ErrNotFound)
		}
		return fmt.Errorf("cannot stat file: %w", err)
	}
	
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: %w", ErrUnauthorized)
		}
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()
	
	// Read file
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}
	
	// Parse JSON
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return ValidationError{
			Field:   "file content",
			Value:   string(data[:min(50, len(data))]),
			Message: "invalid JSON: " + err.Error(),
		}
	}
	
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ============================================================================
// Network Error Handling
// ============================================================================

// NetworkOperation demonstrates network error handling
func NetworkOperation(address string) error {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return TimeoutError{
				Operation: "connection",
				Timeout:   5 * time.Second,
			}
		}
		return fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()
	
	// Set deadline
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("set deadline failed: %w", err)
	}
	
	// Write data
	if _, err := conn.Write([]byte("PING\n")); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return TimeoutError{
				Operation: "write",
				Timeout:   10 * time.Second,
			}
		}
		return fmt.Errorf("write failed: %w", err)
	}
	
	return nil
}

// ============================================================================
// Main Function
// ============================================================================

func main() {
	fmt.Println("=== Basic Error Handling ===")
	BasicErrorHandling("nonexistent.txt")
	
	fmt.Println("\n=== Error Wrapping ===")
	err := WrapError()
	fmt.Printf("Wrapped error: %v\n", err)
	
	fmt.Println("\n=== Error Unwrapping ===")
	UnwrapError()
	
	fmt.Println("\n=== Sentinel Errors ===")
	HandleSentinelErrors()
	
	fmt.Println("\n=== Custom Error Types ===")
	HandleCustomErrors()
	
	fmt.Println("\n=== Error Aggregation ===")
	err = ValidateAll([]string{"valid", "", "too long input", ""})
	if err != nil {
		fmt.Printf("Validation errors:\n%v\n", err)
	}
	
	fmt.Println("\n=== Retry Pattern ===")
	if err := RetryableErrorExample(); err != nil {
		fmt.Printf("Retry failed: %v\n", err)
	} else {
		fmt.Println("Retry succeeded!")
	}
	
	fmt.Println("\n=== Timeout Pattern ===")
	err = WithTimeout(100*time.Millisecond, func() error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	if err != nil {
		fmt.Printf("Timeout error: %v\n", err)
	}
	
	fmt.Println("\n=== Goroutine Error Handling ===")
	if err := GoroutineError(); err != nil {
		fmt.Printf("Goroutine errors: %v\n", err)
	}
	
	fmt.Println("\n=== File Processing ===")
	if err := ProcessFile("config.json"); err != nil {
		fmt.Printf("File processing error: %v\n", err)
		
		// Check error types
		if errors.Is(err, ErrNotFound) {
			fmt.Println("  -> File not found, using defaults")
		}
		
		var valErr ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("  -> Validation error on %s\n", valErr.Field)
		}
	}
	
	fmt.Println("\n=== Network Operation ===")
	if err := NetworkOperation("localhost:9999"); err != nil {
		fmt.Printf("Network error: %v\n", err)
		
		var timeoutErr TimeoutError
		if errors.As(err, &timeoutErr) {
			fmt.Printf("  -> Operation timed out after %v\n", timeoutErr.Timeout)
		}
	}
}