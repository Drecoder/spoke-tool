package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// ============================================================================
// Basic Defer Examples
// ============================================================================

// BasicDefer demonstrates simple defer usage for cleanup
func BasicDefer() {
	// Defer statements are executed in LIFO order (last-in-first-out)
	defer fmt.Println("First defer - actually prints third")
	defer fmt.Println("Second defer - actually prints second")
	defer fmt.Println("Third defer - actually prints first")

	fmt.Println("Function body - prints first")
	// Output:
	// Function body - prints first
	// Third defer - actually prints first
	// Second defer - actually prints second
	// First defer - actually prints third
}

// DeferWithArguments demonstrates that defer captures arguments at the time of defer
func DeferWithArguments() {
	x := 10
	defer fmt.Println("Deferred x =", x) // Captures x=10

	x = 20
	fmt.Println("Current x =", x)
	// Output:
	// Current x = 20
	// Deferred x = 10
}

// DeferWithClosure demonstrates using a closure to capture current values
func DeferWithClosure() {
	x := 10
	defer func() {
		fmt.Println("Closure x =", x) // Captures current x at execution time
	}()

	x = 20
	fmt.Println("Current x =", x)
	// Output:
	// Current x = 20
	// Closure x = 20
}

// ============================================================================
// File Operations
// ============================================================================

// ReadFile demonstrates defer for closing files
func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close() // Ensures file is closed when function exits

	// Read file contents
	data := make([]byte, 1024)
	n, err := file.Read(data)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data[:n], nil
}

// WriteFile demonstrates multiple defers for cleanup
func WriteFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close() // Ensure file is closed

	// Write data
	if _, err := file.Write(data); err != nil {
		// file.Close() will still be called via defer
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Sync to disk
	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

// CopyFile demonstrates using defer for multiple resources
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	return nil
}

// ============================================================================
// Database Operations
// ============================================================================

// DatabaseQuery demonstrates defer for closing rows and transaction handling
func DatabaseQuery(db *sql.DB, query string) error {
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close() // Ensure rows are closed

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows iteration error: %w", err)
	}

	return nil
}

// TransactionWithDefer demonstrates transaction handling with defer
func TransactionWithDefer(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback - if commit succeeds, this will be a no-op
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Perform multiple operations
	if _, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "Alice"); err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}

	if _, err := tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE user_id = ?", 1); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

// ============================================================================
// Mutex Unlocking
// ============================================================================

// Counter demonstrates defer for mutex unlocking
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // Ensure mutex is unlocked even if function panics

	c.value++
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.value
}

// ComplexMutexExample demonstrates multiple mutexes with defer
type BankAccount struct {
	balance int
	mu      sync.RWMutex
}

func (b *BankAccount) Deposit(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.balance += amount
}

func (b *BankAccount) Withdraw(amount int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.balance < amount {
		return errors.New("insufficient funds")
	}

	b.balance -= amount
	return nil
}

func (b *BankAccount) Balance() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.balance
}

// ============================================================================
// Panic Recovery
// ============================================================================

// SafeDivision demonstrates defer with recover
func SafeDivision(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
			err = fmt.Errorf("panic recovered: %v", r)
			result = 0
		}
	}()

	result = a / b
	return result, nil
}

// ProcessWithRecovery demonstrates recovery in goroutines
func ProcessWithRecovery(items []int) {
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Goroutine panicked for value %d: %v", val, r)
					debug.PrintStack()
				}
			}()

			// Simulate processing that might panic
			if val == 0 {
				panic("cannot process zero")
			}
			fmt.Printf("Processing %d\n", 100/val)
		}(item)
	}

	wg.Wait()
}

// ============================================================================
// Performance Measurement
// ============================================================================

// Timer demonstrates defer for measuring function execution time
func Timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("Function took %v\n", time.Since(start))
	}
}

func ExpensiveOperation() {
	defer Timer()() // Defer the execution of the returned function

	// Simulate expensive work
	time.Sleep(2 * time.Second)
	fmt.Println("Expensive operation completed")
}

// TrackTime is a reusable timing decorator
func TrackTime(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func ProcessData() {
	defer TrackTime("ProcessData")()

	// Simulate work
	time.Sleep(500 * time.Millisecond)
}

// ============================================================================
// Resource Cleanup
// ============================================================================

// Resource represents a resource that needs cleanup
type Resource struct {
	id int
}

func NewResource(id int) *Resource {
	fmt.Printf("Resource %d acquired\n", id)
	return &Resource{id: id}
}

func (r *Resource) Close() error {
	fmt.Printf("Resource %d released\n", r.id)
	return nil
}

func (r *Resource) DoWork() error {
	if r.id == 0 {
		return errors.New("invalid resource")
	}
	fmt.Printf("Working with resource %d\n", r.id)
	return nil
}

// MultiResourceExample demonstrates cleanup of multiple resources
func MultiResourceExample() error {
	r1 := NewResource(1)
	defer r1.Close()

	r2 := NewResource(2)
	defer r2.Close()

	r3 := NewResource(3)
	defer r3.Close()

	// Use resources
	if err := r1.DoWork(); err != nil {
		return fmt.Errorf("r1 work failed: %w", err)
	}

	if err := r2.DoWork(); err != nil {
		return fmt.Errorf("r2 work failed: %w", err)
	}

	return nil
}

// ============================================================================
// Logging with Defer
// ============================================================================

// LogFunctionEntry demonstrates defer for logging function entry/exit
func LogFunctionEntry() func() {
	fmt.Println("Function entered")
	return func() {
		fmt.Println("Function exited")
	}
}

func ProcessWithLogging() {
	defer LogFunctionEntry()()

	fmt.Println("Processing...")
	time.Sleep(100 * time.Millisecond)
}

// LogWithError demonstrates defer for logging function result
func LogWithError() func(error) {
	return func(err error) {
		if err != nil {
			fmt.Printf("Function failed: %v\n", err)
		} else {
			fmt.Println("Function succeeded")
		}
	}
}

func OperationThatMayFail(succeed bool) (err error) {
	defer LogWithError()(err)

	if !succeed {
		return errors.New("operation failed")
	}

	fmt.Println("Operation succeeded")
	return nil
}

// ============================================================================
// Conditional Defer
// ============================================================================

// ConditionalDefer demonstrates conditional execution in defer
func ConditionalDefer(shouldLog bool) {
	if shouldLog {
		defer fmt.Println("This will be logged")
	}

	defer fmt.Println("This will always be logged")

	fmt.Println("Doing work...")
}

// ============================================================================
// Nested Functions with Defer
// ============================================================================

func OuterFunction() {
	defer fmt.Println("Outer defer")

	fmt.Println("Outer function")

	func() {
		defer fmt.Println("Inner defer")
		fmt.Println("Inner function")
	}()

	// Output:
	// Outer function
	// Inner function
	// Inner defer
	// Outer defer
}

// ============================================================================
// Error Handling Patterns
// ============================================================================

// DeferWithNamedReturn demonstrates using defer to modify return values
func DeferWithNamedReturn() (result int) {
	defer func() {
		result *= 2
	}()

	result = 10
	return // Returns 20
}

// DeferWithError demonstrates using defer to wrap errors
func DeferWithError() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("operation failed: %w", err)
		}
	}()

	// Some operation that might fail
	if time.Now().Unix()%2 == 0 {
		return errors.New("random error")
	}

	return nil
}

// ============================================================================
// Channel Cleanup
// ============================================================================

// ChannelProducer demonstrates defer for closing channels
func ChannelProducer() <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch) // Ensure channel is closed when goroutine exits
		defer fmt.Println("Channel closed")

		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()

	return ch
}

// ============================================================================
// Complex Real-World Examples
// ============================================================================

// HTTPHandler demonstrates defer in HTTP handlers
func HTTPHandler() {
	// Simulated HTTP handler
	func() {
		// Log request
		defer fmt.Println("Request completed")

		// Track time
		defer TrackTime("HTTPHandler")()

		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Handler panicked: %v\n", r)
				// Return 500 error in real handler
			}
		}()

		fmt.Println("Processing request...")
		// Simulate work
		time.Sleep(100 * time.Millisecond)

		// Simulate panic for demonstration
		// panic("something went wrong")
	}()
}

// FileProcessor demonstrates complex file processing with multiple defers
func FileProcessor(inputFile, outputFile string) error {
	// Open input file
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}
	defer func() {
		input.Close()
		fmt.Println("Input file closed")
	}()

	// Create output file
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer func() {
		output.Close()
		fmt.Println("Output file closed")
	}()

	// Create backup file
	backup, err := os.Create(outputFile + ".bak")
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	defer func() {
		backup.Close()
		// Optionally remove backup if everything succeeded
		if err == nil {
			os.Remove(backup.Name())
		}
		fmt.Println("Backup file closed")
	}()

	// Process data
	buf := make([]byte, 1024)
	for {
		n, readErr := input.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return fmt.Errorf("read error: %w", readErr)
		}
		if n == 0 {
			break
		}

		// Write to both output and backup
		if _, err := output.Write(buf[:n]); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		if _, err := backup.Write(buf[:n]); err != nil {
			return fmt.Errorf("backup write error: %w", err)
		}
	}

	return nil
}

// ============================================================================
// Testing Examples
// ============================================================================

// TestHelper demonstrates defer in test helpers
func TestHelper(t interface {
	Helper()
	Log(args ...interface{})
}) func() {
	t.Helper()
	t.Log("Setup test")

	return func() {
		t.Helper()
		t.Log("Teardown test")
	}
}

func RunTest(t interface {
	Helper()
	Log(args ...interface{})
}) {
	defer TestHelper(t)()

	// Test logic here
	t.Log("Running test...")
}

// ============================================================================
// Main Function for Demonstrations
// ============================================================================

func main() {
	fmt.Println("=== Basic Defer ===")
	BasicDefer()

	fmt.Println("\n=== Defer With Arguments ===")
	DeferWithArguments()

	fmt.Println("\n=== Defer With Closure ===")
	DeferWithClosure()

	fmt.Println("\n=== Performance Measurement ===")
	ExpensiveOperation()
	ProcessData()

	fmt.Println("\n=== Panic Recovery ===")
	result, err := SafeDivision(10, 0)
	fmt.Printf("SafeDivision(10,0) = %d, err=%v\n", result, err)

	fmt.Println("\n=== Mutex Example ===")
	counter := &Counter{}
	counter.Increment()
	counter.Increment()
	fmt.Printf("Counter value: %d\n", counter.GetValue())

	fmt.Println("\n=== Resource Cleanup ===")
	MultiResourceExample()

	fmt.Println("\n=== Named Return Values ===")
	fmt.Printf("DeferWithNamedReturn() = %d\n", DeferWithNamedReturn())

	fmt.Println("\n=== Channel Cleanup ===")
	ch := ChannelProducer()
	for val := range ch {
		fmt.Printf("Received: %d\n", val)
	}

	fmt.Println("\n=== HTTP Handler ===")
	HTTPHandler()

	fmt.Println("\n=== Conditional Defer ===")
	ConditionalDefer(true)
	ConditionalDefer(false)

	fmt.Println("\n=== Nested Functions ===")
	OuterFunction()

	fmt.Println("\n=== Goroutine Recovery ===")
	ProcessWithRecovery([]int{1, 2, 0, 3, 4})

	// Give goroutines time to complete
	time.Sleep(100 * time.Millisecond)
}
