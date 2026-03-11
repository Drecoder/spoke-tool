package edgecases

import (
	"errors"
	"fmt"
	"testing"
)

// ============================================================================
// Basic Edge Cases
// ============================================================================

func TestZeroValues(t *testing.T) {
	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "zero integer input",
			fn: func() error {
				return ProcessInt(0)
			},
		},
		{
			name: "empty string",
			fn: func() error {
				return ProcessString("")
			},
		},
		{
			name: "zero float",
			fn: func() error {
				return ProcessFloat(0.0)
			},
		},
		{
			name: "nil pointer",
			fn: func() error {
				var p *int
				return ProcessPointer(p)
			},
		},
		{
			name: "zero length slice",
			fn: func() error {
				return ProcessSlice([]int{})
			},
		},
		{
			name: "zero length map",
			fn: func() error {
				return ProcessMap(map[string]int{})
			},
		},
		{
			name: "zero channel",
			fn: func() error {
				return ProcessChannel(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

// ============================================================================
// Boundary Conditions
// ============================================================================

func TestBoundaryConditions(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{
			name:    "minimum int32",
			input:   -2147483648,
			wantErr: false,
		},
		{
			name:    "maximum int32",
			input:   2147483647,
			wantErr: false,
		},
		{
			name:    "minimum int64",
			input:   -9223372036854775808,
			wantErr: false,
		},
		{
			name:    "maximum int64",
			input:   9223372036854775807,
			wantErr: false,
		},
		{
			name:    "just below threshold",
			input:   99,
			wantErr: false,
		},
		{
			name:    "at threshold",
			input:   100,
			wantErr: false,
		},
		{
			name:    "just above threshold",
			input:   101,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateThreshold(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateThreshold(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// Overflow and Underflow
// ============================================================================

func TestOverflow(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		wantErr bool
	}{
		{
			name:    "max int + 1",
			a:       9223372036854775807,
			b:       1,
			wantErr: true,
		},
		{
			name:    "min int - 1",
			a:       -9223372036854775808,
			b:       -1,
			wantErr: true,
		},
		{
			name:    "large addition within bounds",
			a:       500000000,
			b:       500000000,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SafeAdd(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeAdd(%d, %d) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// Division by Zero
// ============================================================================

func TestDivisionByZero(t *testing.T) {
	tests := []struct {
		name      string
		a, b      float64
		wantPanic bool
	}{
		{
			name:      "integer division by zero",
			a:         10,
			b:         0,
			wantPanic: true,
		},
		{
			name:      "float division by zero",
			a:         3.14,
			b:         0.0,
			wantPanic: true,
		},
		{
			name:      "negative division by zero",
			a:         -10,
			b:         0,
			wantPanic: true,
		},
		{
			name:      "zero divided by zero",
			a:         0,
			b:         0,
			wantPanic: true,
		},
		{
			name:      "normal division",
			a:         10,
			b:         2,
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("Divide(%f, %f) panic = %v, wantPanic %v", tt.a, tt.b, r, tt.wantPanic)
				}
			}()
			Divide(tt.a, tt.b)
		})
	}
}

// ============================================================================
// Nil and Empty Collections
// ============================================================================

func TestNilSlices(t *testing.T) {
	var nilSlice []int
	emptySlice := []int{}

	tests := []struct {
		name    string
		slice   []int
		wantErr bool
	}{
		{
			name:    "nil slice",
			slice:   nilSlice,
			wantErr: true,
		},
		{
			name:    "empty slice",
			slice:   emptySlice,
			wantErr: true,
		},
		{
			name:    "non-empty slice",
			slice:   []int{1, 2, 3},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessSlice(tt.slice)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessSlice(%v) error = %v, wantErr %v", tt.slice, err, tt.wantErr)
			}
		})
	}
}

func TestNilMaps(t *testing.T) {
	var nilMap map[string]int
	emptyMap := make(map[string]int)

	tests := []struct {
		name    string
		m       map[string]int
		wantErr bool
	}{
		{
			name:    "nil map",
			m:       nilMap,
			wantErr: true,
		},
		{
			name:    "empty map",
			m:       emptyMap,
			wantErr: true,
		},
		{
			name:    "non-empty map",
			m:       map[string]int{"one": 1, "two": 2},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessMap(tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessMap(%v) error = %v, wantErr %v", tt.m, err, tt.wantErr)
			}
		})
	}
}

func TestNilInterfaces(t *testing.T) {
	var nilInterface interface{}
	var nilReader interface{ Read([]byte) (int, error) } = nil

	tests := []struct {
		name    string
		iface   interface{}
		wantErr bool
	}{
		{
			name:    "nil interface",
			iface:   nilInterface,
			wantErr: true,
		},
		{
			name:    "nil reader",
			iface:   nilReader,
			wantErr: true,
		},
		{
			name:    "non-nil interface",
			iface:   struct{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessInterface(tt.iface)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessInterface(%v) error = %v, wantErr %v", tt.iface, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// Index Out of Bounds
// ============================================================================

func TestIndexOutOfBounds(t *testing.T) {
	slice := []int{10, 20, 30}

	tests := []struct {
		name      string
		index     int
		wantPanic bool
	}{
		{
			name:      "negative index",
			index:     -1,
			wantPanic: true,
		},
		{
			name:      "index zero",
			index:     0,
			wantPanic: false,
		},
		{
			name:      "last index",
			index:     len(slice) - 1,
			wantPanic: false,
		},
		{
			name:      "index equal to length",
			index:     len(slice),
			wantPanic: true,
		},
		{
			name:      "index greater than length",
			index:     len(slice) + 5,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("GetElement(%d) panic = %v, wantPanic %v", tt.index, r, tt.wantPanic)
				}
			}()
			GetElement(slice, tt.index)
		})
	}
}

// ============================================================================
// String Edge Cases
// ============================================================================

func TestStringEdgeCases(t *testing.T) {
	longString := string(make([]byte, 1000000)) // 1MB string

	tests := []struct {
		name    string
		s       string
		wantErr bool
	}{
		{
			name:    "empty string",
			s:       "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			s:       "   \t\n  ",
			wantErr: false,
		},
		{
			name:    "very long string",
			s:       longString,
			wantErr: false,
		},
		{
			name:    "unicode string",
			s:       "Hello, 世界! 🚀",
			wantErr: false,
		},
		{
			name:    "string with null bytes",
			s:       "hello\x00world",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessString(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessString(%q) error = %v, wantErr %v", tt.s, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// Recursion Edge Cases
// ============================================================================

func TestRecursionDepth(t *testing.T) {
	tests := []struct {
		name      string
		depth     int
		wantPanic bool
	}{
		{
			name:      "zero depth",
			depth:     0,
			wantPanic: false,
		},
		{
			name:      "shallow recursion",
			depth:     10,
			wantPanic: false,
		},
		{
			name:      "deep recursion",
			depth:     10000,
			wantPanic: true, // May cause stack overflow
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("RecursiveFunction(%d) panic = %v, wantPanic %v", tt.depth, r, tt.wantPanic)
				}
			}()
			RecursiveFunction(tt.depth)
		})
	}
}

// ============================================================================
// Channel Edge Cases
// ============================================================================

func TestChannelEdgeCases(t *testing.T) {
	ch := make(chan int)

	tests := []struct {
		name      string
		op        func()
		wantPanic bool
	}{
		{
			name: "send on nil channel",
			op: func() {
				var nilCh chan int
				nilCh <- 42
			},
			wantPanic: true,
		},
		{
			name: "receive from nil channel",
			op: func() {
				var nilCh chan int
				<-nilCh
			},
			wantPanic: true,
		},
		{
			name: "close nil channel",
			op: func() {
				var nilCh chan int
				close(nilCh)
			},
			wantPanic: true,
		},
		{
			name: "close already closed channel",
			op: func() {
				close(ch)
				close(ch)
			},
			wantPanic: true,
		},
		{
			name: "send on closed channel",
			op: func() {
				close(ch)
				ch <- 42
			},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("%s panic = %v, wantPanic %v", tt.name, r, tt.wantPanic)
				}
			}()
			tt.op()
		})
	}
}

// ============================================================================
// Race Conditions (for testing race detector)
// ============================================================================

func TestRaceConditions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race condition tests in short mode")
	}

	t.Run("concurrent map access", func(t *testing.T) {
		m := make(map[int]int)
		done := make(chan bool)

		// Writer
		go func() {
			for i := 0; i < 100; i++ {
				m[i] = i
			}
			done <- true
		}()

		// Reader
		go func() {
			for i := 0; i < 100; i++ {
				_ = m[i]
			}
			done <- true
		}()

		<-done
		<-done
	})

	t.Run("concurrent counter", func(t *testing.T) {
		var counter int
		done := make(chan bool)

		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 1000; j++ {
					counter++
				}
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// ============================================================================
// Error Handling Edge Cases
// ============================================================================

func TestErrorWrapping(t *testing.T) {
	err := errors.New("root cause")

	tests := []struct {
		name       string
		err        error
		wantString string
	}{
		{
			name:       "nil error",
			err:        nil,
			wantString: "",
		},
		{
			name:       "simple error",
			err:        err,
			wantString: "root cause",
		},
		{
			name:       "wrapped error",
			err:        WrapError(err, "context"),
			wantString: "context: root cause",
		},
		{
			name:       "doubly wrapped error",
			err:        WrapError(WrapError(err, "level1"), "level2"),
			wantString: "level2: level1: root cause",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				if tt.err.Error() != tt.wantString {
					t.Errorf("Error = %q, want %q", tt.err.Error(), tt.wantString)
				}
			}
		})
	}
}

// ============================================================================
// Type Assertion Edge Cases
// ============================================================================

func TestTypeAssertions(t *testing.T) {
	var i interface{} = "hello"

	tests := []struct {
		name      string
		value     interface{}
		target    interface{}
		wantPanic bool
	}{
		{
			name:      "valid type assertion",
			value:     i,
			target:    "string",
			wantPanic: false,
		},
		{
			name:      "invalid type assertion",
			value:     i,
			target:    42,
			wantPanic: true,
		},
		{
			name:      "nil interface",
			value:     nil,
			target:    "string",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("TypeAssertion(%v) panic = %v, wantPanic %v", tt.value, r, tt.wantPanic)
				}
			}()

			switch tt.target.(type) {
			case string:
				_ = tt.value.(string)
			case int:
				_ = tt.value.(int)
			}
		})
	}
}

// ============================================================================
// Helper Functions (these would be in your actual code)
// ============================================================================

func ProcessInt(x int) error {
	if x == 0 {
		return errors.New("zero value not allowed")
	}
	return nil
}

func ProcessString(s string) error {
	if s == "" {
		return errors.New("empty string not allowed")
	}
	return nil
}

func ProcessFloat(f float64) error {
	if f == 0.0 {
		return errors.New("zero value not allowed")
	}
	return nil
}

func ProcessPointer(p *int) error {
	if p == nil {
		return errors.New("nil pointer not allowed")
	}
	return nil
}

func ProcessSlice(s []int) error {
	if s == nil {
		return errors.New("nil slice not allowed")
	}
	if len(s) == 0 {
		return errors.New("empty slice not allowed")
	}
	return nil
}

func ProcessMap(m map[string]int) error {
	if m == nil {
		return errors.New("nil map not allowed")
	}
	if len(m) == 0 {
		return errors.New("empty map not allowed")
	}
	return nil
}

func ProcessChannel(ch chan int) error {
	if ch == nil {
		return errors.New("nil channel not allowed")
	}
	return nil
}

func ProcessInterface(i interface{}) error {
	if i == nil {
		return errors.New("nil interface not allowed")
	}
	return nil
}

func ValidateThreshold(x int) error {
	if x > 100 {
		return errors.New("value exceeds threshold")
	}
	return nil
}

func SafeAdd(a, b int) (int, error) {
	if a > 0 && b > 0 && a > (1<<63-1)-b {
		return 0, errors.New("integer overflow")
	}
	if a < 0 && b < 0 && a < (-1<<63)-b {
		return 0, errors.New("integer underflow")
	}
	return a + b, nil
}

func Divide(a, b float64) float64 {
	if b == 0 {
		panic("division by zero")
	}
	return a / b
}

func GetElement(slice []int, index int) int {
	if index < 0 || index >= len(slice) {
		panic("index out of bounds")
	}
	return slice[index]
}

func RecursiveFunction(n int) int {
	if n <= 0 {
		return 0
	}
	return 1 + RecursiveFunction(n-1)
}

func WrapError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
