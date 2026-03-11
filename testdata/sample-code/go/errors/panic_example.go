package main

import (
	"errors"
	"fmt"
	"runtime/debug"
	"time"
)

// ============================================================================
// Basic Panic Examples
// ============================================================================

// BasicPanic demonstrates a simple panic
func BasicPanic() {
	fmt.Println("About to panic...")
	panic("something went wrong")
	// This line never executes
	fmt.Println("This will not print")
}

// PanicWithInt demonstrates panicking with different types
func PanicWithInt() {
	panic(42) // Panic with integer
}

// PanicWithError demonstrates panicking with an error
func PanicWithError() {
	err := errors.New("critical error")
	panic(err)
}

// ============================================================================
// Basic Recover Examples
// ============================================================================

// BasicRecover demonstrates simple panic recovery
func BasicRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	fmt.Println("About to panic...")
	panic("something went wrong")
	fmt.Println("This will not print")
}

// RecoverWithType demonstrates recovering and type asserting
func RecoverWithType() {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case string:
				fmt.Printf("Recovered from string panic: %s\n", v)
			case int:
				fmt.Printf("Recovered from int panic: %d\n", v)
			case error:
				fmt.Printf("Recovered from error panic: %v\n", v)
			default:
				fmt.Printf("Recovered from unknown panic type: %v\n", v)
			}
		}
	}()

	// Try different panic types
	panic("urgent message")
	// panic(42)
	// panic(errors.New("critical failure"))
}

// ============================================================================
// Panic in Goroutines
// ============================================================================

// GoroutinePanic demonstrates panic in a goroutine
func GoroutinePanic() {
	go func() {
		panic("goroutine panic")
	}()

	// Give goroutine time to panic
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main function continues") // This may or may not run
}

// GoroutineWithRecover demonstrates recovering panics in goroutines
func GoroutineWithRecover() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Goroutine recovered: %v\n", r)
			}
		}()

		panic("goroutine panic")
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main function continues safely")
}

// ============================================================================
// Nested Panics and Recovers
// ============================================================================

// NestedPanic demonstrates panic in nested function calls
func NestedPanic() {
	defer fmt.Println("Outer defer 1")
	defer fmt.Println("Outer defer 2")

	func() {
		defer fmt.Println("Inner defer 1")
		defer fmt.Println("Inner defer 2")

		fmt.Println("About to panic in inner function")
		panic("inner panic")
	}()
}

// NestedRecover demonstrates recovering at different levels
func NestedRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered at outer level: %v\n", r)
		}
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered at inner level: %v\n", r)
				// Optionally re-panic
				// panic(r)
			}
		}()

		panic("deep panic")
	}()

	fmt.Println("After inner function")
}

// ============================================================================
// Panic with Defer Chain
// ============================================================================

// PanicWithDeferChain demonstrates defer execution order during panic
func PanicWithDeferChain() {
	defer fmt.Println("First defer")
	defer fmt.Println("Second defer")
	defer fmt.Println("Third defer")

	fmt.Println("Before panic")
	panic("panic occurred")
	fmt.Println("After panic") // Never executes
}

// ============================================================================
// Converting Panic to Error
// ============================================================================

// SafeFunction converts panics to errors
func SafeFunction(shouldPanic bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case string:
				err = fmt.Errorf("panic: %s", v)
			case error:
				err = fmt.Errorf("panic: %w", v)
			default:
				err = fmt.Errorf("panic: %v", v)
			}
		}
	}()

	if shouldPanic {
		panic("something went wrong")
	}

	return nil
}

// Must function panics on error (common in initializers)
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// ============================================================================
// Panic in Initialization
// ============================================================================

var (
	// This will panic if the config is invalid
	config = initConfig()
)

func initConfig() map[string]string {
	cfg := make(map[string]string)

	// Simulate loading config
	cfg["host"] = "localhost"
	cfg["port"] = "8080"

	// Validate required fields
	if cfg["host"] == "" {
		panic("config: host is required")
	}

	return cfg
}

// ============================================================================
// Panic for Impossible Cases
// ============================================================================

// Unreachable panics for code that should never be reached
func Unreachable() {
	panic("unreachable code reached")
}

// Assert panics if condition is false (like an assertion)
func Assert(condition bool, message string) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: %s", message))
	}
}

// ============================================================================
// Complex Recovery Patterns
// ============================================================================

// RecoveryWithContext demonstrates recovering with context information
type RecoveryContext struct {
	Function string
	Args     []interface{}
	Panic    interface{}
	Stack    []byte
}

// SafeExecute executes a function with panic recovery and context
func SafeExecute(name string, fn func(), args ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			context := RecoveryContext{
				Function: name,
				Args:     args,
				Panic:    r,
				Stack:    debug.Stack(),
			}

			// Log the context
			fmt.Printf("Panic recovered in %s: %v\n", name, r)
			fmt.Printf("Stack trace:\n%s\n", context.Stack)

			// Convert to error
			err = fmt.Errorf("panic in %s: %v", name, r)
		}
	}()

	fn()
	return nil
}

// ============================================================================
// Resource Cleanup During Panic
// ============================================================================

// Resource represents a resource that needs cleanup
type Resource struct {
	id int
}

func NewResource(id int) *Resource {
	fmt.Printf("Resource %d acquired\n", id)
	return &Resource{id: id}
}

func (r *Resource) Close() {
	fmt.Printf("Resource %d closed\n", r.id)
}

// ResourceCleanup demonstrates proper resource cleanup during panic
func ResourceCleanup() {
	r1 := NewResource(1)
	defer r1.Close()

	r2 := NewResource(2)
	defer r2.Close()

	fmt.Println("Working with resources...")

	// Simulate panic
	panic("something went wrong while using resources")
}

// ============================================================================
// Database Transaction Example
// ============================================================================

// Transaction represents a database transaction
type Transaction struct {
	id     string
	active bool
}

func BeginTransaction() *Transaction {
	tx := &Transaction{
		id:     "tx-123",
		active: true,
	}
	fmt.Printf("Transaction %s started\n", tx.id)
	return tx
}

func (tx *Transaction) Commit() {
	if tx.active {
		fmt.Printf("Transaction %s committed\n", tx.id)
		tx.active = false
	}
}

func (tx *Transaction) Rollback() {
	if tx.active {
		fmt.Printf("Transaction %s rolled back\n", tx.id)
		tx.active = false
	}
}

func (tx *Transaction) Execute(query string) {
	if !tx.active {
		panic("cannot execute on inactive transaction")
	}
	fmt.Printf("Executing: %s\n", query)
}

// TransactionExample demonstrates transaction safety with panic
func TransactionExample() {
	tx := BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic detected: %v\n", r)
			tx.Rollback()
		}
	}()

	tx.Execute("INSERT INTO users (name) VALUES ('Alice')")
	tx.Execute("UPDATE accounts SET balance = balance - 100")

	// Simulate panic
	if true {
		panic("database connection lost")
	}

	tx.Commit()
}

// ============================================================================
// Worker Pool with Panic Recovery
// ============================================================================

// Job represents a unit of work
type Job struct {
	ID   int
	Task func() error
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	jobs    chan Job
	results chan error
}

func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{
		jobs:    make(chan Job, 100),
		results: make(chan error, 100),
	}

	for i := 0; i < size; i++ {
		go pool.worker(i)
	}

	return pool
}

func (p *WorkerPool) worker(id int) {
	for job := range p.jobs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Worker %d panicked on job %d: %v\n", id, job.ID, r)
					p.results <- fmt.Errorf("panic in worker %d: %v", id, r)
				}
			}()

			err := job.Task()
			p.results <- err
		}()
	}
}

func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job
}

func (p *WorkerPool) Wait() []error {
	close(p.jobs)

	var errors []error
	for i := 0; i < cap(p.results); i++ {
		if err := <-p.results; err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// ============================================================================
// Panic in Select Statements
// ============================================================================

// SelectWithPanic demonstrates panic inside select
func SelectWithPanic() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- 2
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in select: %v\n", r)
		}
	}()

	select {
	case <-ch1:
		fmt.Println("Received from ch1")
		panic("panic after receiving from ch1")
	case <-ch2:
		fmt.Println("Received from ch2")
	case <-time.After(50 * time.Millisecond):
		fmt.Println("Timeout")
	}
}

// ============================================================================
// Panic with Custom Types
// ============================================================================

// PanicData represents structured panic information
type PanicData struct {
	Code    int
	Message string
	Time    time.Time
}

func (p PanicData) String() string {
	return fmt.Sprintf("[%d] %s at %s", p.Code, p.Message, p.Time.Format(time.RFC3339))
}

// StructuredPanic demonstrates panicking with custom struct
func StructuredPanic() {
	defer func() {
		if r := recover(); r != nil {
			if data, ok := r.(PanicData); ok {
				fmt.Printf("Structured panic: %s\n", data)
				fmt.Printf("  Code: %d\n", data.Code)
				fmt.Printf("  Message: %s\n", data.Message)
				fmt.Printf("  Time: %v\n", data.Time)
			} else {
				fmt.Printf("Unknown panic type: %v\n", r)
			}
		}
	}()

	panic(PanicData{
		Code:    500,
		Message: "internal server error",
		Time:    time.Now(),
	})
}

// ============================================================================
// Panic Chain
// ============================================================================

// Level1 panics and recovers, then re-panics
func Level1() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Level1 recovered:", r)
			panic("re-panic from level1")
		}
	}()

	Level2()
}

func Level2() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Level2 recovered:", r)
			panic("re-panic from level2")
		}
	}()

	Level3()
}

func Level3() {
	panic("original panic from level3")
}

// ============================================================================
// Main Function
// ============================================================================

func main() {
	fmt.Println("=== Basic Recover ===")
	BasicRecover()

	fmt.Println("\n=== Recover with Type ===")
	RecoverWithType()

	fmt.Println("\n=== Goroutine with Recover ===")
	GoroutineWithRecover()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== Nested Recover ===")
	NestedRecover()

	fmt.Println("\n=== Panic with Defer Chain ===")
	func() {
		defer func() {
			recover()
		}()
		PanicWithDeferChain()
	}()

	fmt.Println("\n=== Converting Panic to Error ===")
	err := SafeFunction(true)
	if err != nil {
		fmt.Printf("Error from safe function: %v\n", err)
	}

	fmt.Println("\n=== Resource Cleanup During Panic ===")
	func() {
		defer func() {
			recover()
		}()
		ResourceCleanup()
	}()

	fmt.Println("\n=== Transaction Example ===")
	func() {
		defer func() {
			recover()
		}()
		TransactionExample()
	}()

	fmt.Println("\n=== Worker Pool with Panic Recovery ===")
	pool := NewWorkerPool(3)

	for i := 0; i < 5; i++ {
		jobID := i
		pool.Submit(Job{
			ID: jobID,
			Task: func() error {
				if jobID == 2 {
					panic(fmt.Sprintf("job %d panicked", jobID))
				}
				fmt.Printf("Job %d completed\n", jobID)
				return nil
			},
		})
	}

	errors := pool.Wait()
	if len(errors) > 0 {
		fmt.Printf("Pool errors: %v\n", errors)
	}

	fmt.Println("\n=== Structured Panic ===")
	StructuredPanic()

	fmt.Println("\n=== Panic Chain ===")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Main recovered final panic:", r)
			}
		}()
		Level1()
	}()

	fmt.Println("\n=== Select with Panic ===")
	SelectWithPanic()

	fmt.Println("\n=== Main function continuing...")
}

// ============================================================================
// Best Practices Summary
// ============================================================================

/*
Panic Best Practices:

1. Use panic for unrecoverable errors:
   - Initialization failures
   - Impossible conditions (should never happen)
   - Programming errors (nil pointer, index out of bounds)

2. Use recover sparingly:
   - In goroutines to prevent crash
   - At API boundaries to convert panics to errors
   - In long-running services to maintain stability

3. Never recover in library code (let caller decide)

4. Always clean up resources with defer

5. Log stack traces when recovering

6. Consider converting panics to errors at package boundaries
*/
