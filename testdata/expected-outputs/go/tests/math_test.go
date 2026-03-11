package math

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        2,
			b:        3,
			expected: 5,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -3,
			expected: -5,
		},
		{
			name:     "mixed signs",
			a:        -2,
			b:        3,
			expected: 1,
		},
		{
			name:     "zero values",
			a:        0,
			b:        5,
			expected: 5,
		},
		{
			name:     "both zero",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "large numbers",
			a:        1000000,
			b:        2000000,
			expected: 3000000,
		},
		{
			name:     "max int",
			a:        1<<31 - 1,
			b:        1,
			expected: 1 << 31, // Note: this might overflow in some languages
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive result",
			a:        10,
			b:        4,
			expected: 6,
		},
		{
			name:     "negative result",
			a:        4,
			b:        10,
			expected: -6,
		},
		{
			name:     "zero result",
			a:        5,
			b:        5,
			expected: 0,
		},
		{
			name:     "negative numbers",
			a:        -10,
			b:        -4,
			expected: -6,
		},
		{
			name:     "mixed signs",
			a:        -10,
			b:        4,
			expected: -14,
		},
		{
			name:     "subtract zero",
			a:        7,
			b:        0,
			expected: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        2,
			b:        3,
			expected: 6,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -3,
			expected: 6,
		},
		{
			name:     "mixed signs",
			a:        -2,
			b:        3,
			expected: -6,
		},
		{
			name:     "zero multiplicand",
			a:        0,
			b:        5,
			expected: 0,
		},
		{
			name:     "zero multiplier",
			a:        5,
			b:        0,
			expected: 0,
		},
		{
			name:     "both zero",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "large numbers",
			a:        1000,
			b:        1000,
			expected: 1000000,
		},
		{
			name:     "multiply by one",
			a:        42,
			b:        1,
			expected: 42,
		},
		{
			name:     "multiply by negative one",
			a:        42,
			b:        -1,
			expected: -42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int
		expected  int
		wantError bool
	}{
		{
			name:      "exact division",
			a:         10,
			b:         2,
			expected:  5,
			wantError: false,
		},
		{
			name:      "division with remainder",
			a:         10,
			b:         3,
			expected:  3, // integer division truncates
			wantError: false,
		},
		{
			name:      "division by zero",
			a:         10,
			b:         0,
			expected:  0,
			wantError: true,
		},
		{
			name:      "negative dividend",
			a:         -10,
			b:         2,
			expected:  -5,
			wantError: false,
		},
		{
			name:      "negative divisor",
			a:         10,
			b:         -2,
			expected:  -5,
			wantError: false,
		},
		{
			name:      "both negative",
			a:         -10,
			b:         -2,
			expected:  5,
			wantError: false,
		},
		{
			name:      "zero dividend",
			a:         0,
			b:         5,
			expected:  0,
			wantError: false,
		},
		{
			name:      "divide one",
			a:         1,
			b:         2,
			expected:  0,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)

			if tt.wantError {
				if err == nil {
					t.Errorf("Divide(%d, %d) expected error, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("Divide(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func TestModulo(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int
		expected  int
		wantError bool
	}{
		{
			name:      "positive modulo",
			a:         10,
			b:         3,
			expected:  1,
			wantError: false,
		},
		{
			name:      "exact division",
			a:         10,
			b:         2,
			expected:  0,
			wantError: false,
		},
		{
			name:      "modulo by zero",
			a:         10,
			b:         0,
			expected:  0,
			wantError: true,
		},
		{
			name:      "negative dividend",
			a:         -10,
			b:         3,
			expected:  -1, // Sign follows dividend in Go
			wantError: false,
		},
		{
			name:      "negative divisor",
			a:         10,
			b:         -3,
			expected:  1, // Result sign follows dividend
			wantError: false,
		},
		{
			name:      "both negative",
			a:         -10,
			b:         -3,
			expected:  -1,
			wantError: false,
		},
		{
			name:      "zero dividend",
			a:         0,
			b:         5,
			expected:  0,
			wantError: false,
		},
		{
			name:      "dividend smaller than divisor",
			a:         2,
			b:         5,
			expected:  2,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Modulo(tt.a, tt.b)

			if tt.wantError {
				if err == nil {
					t.Errorf("Modulo(%d, %d) expected error, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Modulo(%d, %d) unexpected error: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("Modulo(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "positive number",
			input:    42,
			expected: 42,
		},
		{
			name:     "negative number",
			input:    -42,
			expected: 42,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "large positive",
			input:    1000000,
			expected: 1000000,
		},
		{
			name:     "large negative",
			input:    -1000000,
			expected: 1000000,
		},
		{
			name:     "min int",
			input:    -1 << 31,
			expected: 1 << 31, // Note: this might overflow in some languages
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Abs(tt.input)
			if result != tt.expected {
				t.Errorf("Abs(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		name     string
		base     int
		exp      int
		expected int
	}{
		{
			name:     "positive exponent",
			base:     2,
			exp:      3,
			expected: 8,
		},
		{
			name:     "zero exponent",
			base:     5,
			exp:      0,
			expected: 1,
		},
		{
			name:     "one exponent",
			base:     5,
			exp:      1,
			expected: 5,
		},
		{
			name:     "negative base even exponent",
			base:     -2,
			exp:      2,
			expected: 4,
		},
		{
			name:     "negative base odd exponent",
			base:     -2,
			exp:      3,
			expected: -8,
		},
		{
			name:     "zero base",
			base:     0,
			exp:      5,
			expected: 0,
		},
		{
			name:     "zero base zero exponent",
			base:     0,
			exp:      0,
			expected: 1, // Mathematical convention
		},
		{
			name:     "large exponent",
			base:     2,
			exp:      10,
			expected: 1024,
		},
		{
			name:     "base one",
			base:     1,
			exp:      100,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pow(tt.base, tt.exp)
			if result != tt.expected {
				t.Errorf("Pow(%d, %d) = %d, want %d", tt.base, tt.exp, result, tt.expected)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "first larger",
			a:        10,
			b:        5,
			expected: 10,
		},
		{
			name:     "second larger",
			a:        3,
			b:        8,
			expected: 8,
		},
		{
			name:     "equal values",
			a:        7,
			b:        7,
			expected: 7,
		},
		{
			name:     "negative numbers",
			a:        -5,
			b:        -10,
			expected: -5,
		},
		{
			name:     "mixed signs",
			a:        -5,
			b:        3,
			expected: 3,
		},
		{
			name:     "with zero",
			a:        0,
			b:        -5,
			expected: 0,
		},
		{
			name:     "large numbers",
			a:        1000000,
			b:        999999,
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Max(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "first smaller",
			a:        3,
			b:        8,
			expected: 3,
		},
		{
			name:     "second smaller",
			a:        10,
			b:        5,
			expected: 5,
		},
		{
			name:     "equal values",
			a:        7,
			b:        7,
			expected: 7,
		},
		{
			name:     "negative numbers",
			a:        -10,
			b:        -5,
			expected: -10,
		},
		{
			name:     "mixed signs",
			a:        -5,
			b:        3,
			expected: -5,
		},
		{
			name:     "with zero",
			a:        0,
			b:        -5,
			expected: -5,
		},
		{
			name:     "large numbers",
			a:        1000000,
			b:        999999,
			expected: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Min(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		min      int
		max      int
		expected int
	}{
		{
			name:     "value within range",
			value:    5,
			min:      1,
			max:      10,
			expected: 5,
		},
		{
			name:     "value below min",
			value:    0,
			min:      1,
			max:      10,
			expected: 1,
		},
		{
			name:     "value above max",
			value:    15,
			min:      1,
			max:      10,
			expected: 10,
		},
		{
			name:     "value equals min",
			value:    1,
			min:      1,
			max:      10,
			expected: 1,
		},
		{
			name:     "value equals max",
			value:    10,
			min:      1,
			max:      10,
			expected: 10,
		},
		{
			name:     "negative values",
			value:    -5,
			min:      -10,
			max:      -1,
			expected: -5,
		},
		{
			name:     "negative below min",
			value:    -15,
			min:      -10,
			max:      -1,
			expected: -10,
		},
		{
			name:     "min greater than max",
			value:    5,
			min:      10,
			max:      1,
			expected: 5, // Should handle gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clamp(tt.value, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.value, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}

func TestIsEven(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{
			name:     "even positive",
			input:    4,
			expected: true,
		},
		{
			name:     "odd positive",
			input:    5,
			expected: false,
		},
		{
			name:     "even negative",
			input:    -4,
			expected: true,
		},
		{
			name:     "odd negative",
			input:    -5,
			expected: false,
		},
		{
			name:     "zero",
			input:    0,
			expected: true,
		},
		{
			name:     "large even",
			input:    1000000,
			expected: true,
		},
		{
			name:     "large odd",
			input:    1000001,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEven(tt.input)
			if result != tt.expected {
				t.Errorf("IsEven(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsOdd(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{
			name:     "odd positive",
			input:    5,
			expected: true,
		},
		{
			name:     "even positive",
			input:    4,
			expected: false,
		},
		{
			name:     "odd negative",
			input:    -5,
			expected: true,
		},
		{
			name:     "even negative",
			input:    -4,
			expected: false,
		},
		{
			name:     "zero",
			input:    0,
			expected: false,
		},
		{
			name:     "large odd",
			input:    1000001,
			expected: true,
		},
		{
			name:     "large even",
			input:    1000000,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsOdd(tt.input)
			if result != tt.expected {
				t.Errorf("IsOdd(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        48,
			b:        18,
			expected: 6,
		},
		{
			name:     "one is zero",
			a:        12,
			b:        0,
			expected: 12,
		},
		{
			name:     "both zero",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "equal numbers",
			a:        15,
			b:        15,
			expected: 15,
		},
		{
			name:     "one is one",
			a:        1,
			b:        100,
			expected: 1,
		},
		{
			name:     "coprime numbers",
			a:        17,
			b:        19,
			expected: 1,
		},
		{
			name:     "negative numbers",
			a:        -48,
			b:        18,
			expected: 6, // GCD is usually positive
		},
		{
			name:     "large numbers",
			a:        1071,
			b:        462,
			expected: 21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GCD(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestLCM(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        12,
			b:        18,
			expected: 36,
		},
		{
			name:     "one is zero",
			a:        12,
			b:        0,
			expected: 0,
		},
		{
			name:     "both zero",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "equal numbers",
			a:        15,
			b:        15,
			expected: 15,
		},
		{
			name:     "coprime numbers",
			a:        17,
			b:        19,
			expected: 323,
		},
		{
			name:     "one divides other",
			a:        6,
			b:        12,
			expected: 12,
		},
		{
			name:     "negative numbers",
			a:        -12,
			b:        18,
			expected: 36, // LCM is usually positive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LCM(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("LCM(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		expected  int
		wantError bool
	}{
		{
			name:      "zero",
			n:         0,
			expected:  1,
			wantError: false,
		},
		{
			name:      "one",
			n:         1,
			expected:  1,
			wantError: false,
		},
		{
			name:      "small number",
			n:         5,
			expected:  120,
			wantError: false,
		},
		{
			name:      "medium number",
			n:         10,
			expected:  3628800,
			wantError: false,
		},
		{
			name:      "negative number",
			n:         -5,
			expected:  0,
			wantError: true,
		},
		{
			name:      "large number (potential overflow)",
			n:         20,
			expected:  2432902008176640000,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Factorial(tt.n)

			if tt.wantError {
				if err == nil {
					t.Errorf("Factorial(%d) expected error, got nil", tt.n)
				}
			} else {
				if err != nil {
					t.Errorf("Factorial(%d) unexpected error: %v", tt.n, err)
				}
				if result != tt.expected {
					t.Errorf("Factorial(%d) = %d, want %d", tt.n, result, tt.expected)
				}
			}
		})
	}
}

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		expected  int
		wantError bool
	}{
		{
			name:      "zero",
			n:         0,
			expected:  0,
			wantError: false,
		},
		{
			name:      "one",
			n:         1,
			expected:  1,
			wantError: false,
		},
		{
			name:      "two",
			n:         2,
			expected:  1,
			wantError: false,
		},
		{
			name:      "three",
			n:         3,
			expected:  2,
			wantError: false,
		},
		{
			name:      "four",
			n:         4,
			expected:  3,
			wantError: false,
		},
		{
			name:      "five",
			n:         5,
			expected:  5,
			wantError: false,
		},
		{
			name:      "six",
			n:         6,
			expected:  8,
			wantError: false,
		},
		{
			name:      "ten",
			n:         10,
			expected:  55,
			wantError: false,
		},
		{
			name:      "negative",
			n:         -5,
			expected:  0,
			wantError: true,
		},
		{
			name:      "large",
			n:         20,
			expected:  6765,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Fibonacci(tt.n)

			if tt.wantError {
				if err == nil {
					t.Errorf("Fibonacci(%d) expected error, got nil", tt.n)
				}
			} else {
				if err != nil {
					t.Errorf("Fibonacci(%d) unexpected error: %v", tt.n, err)
				}
				if result != tt.expected {
					t.Errorf("Fibonacci(%d) = %d, want %d", tt.n, result, tt.expected)
				}
			}
		})
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected bool
	}{
		{
			name:     "zero",
			n:        0,
			expected: false,
		},
		{
			name:     "one",
			n:        1,
			expected: false,
		},
		{
			name:     "two",
			n:        2,
			expected: true,
		},
		{
			name:     "three",
			n:        3,
			expected: true,
		},
		{
			name:     "four",
			n:        4,
			expected: false,
		},
		{
			name:     "five",
			n:        5,
			expected: true,
		},
		{
			name:     "nine",
			n:        9,
			expected: false,
		},
		{
			name:     "seventeen",
			n:        17,
			expected: true,
		},
		{
			name:     "twenty-one",
			n:        21,
			expected: false,
		},
		{
			name:     "negative",
			n:        -7,
			expected: false,
		},
		{
			name:     "large prime",
			n:        997,
			expected: true,
		},
		{
			name:     "large composite",
			n:        999,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrime(tt.n)
			if result != tt.expected {
				t.Errorf("IsPrime(%d) = %v, want %v", tt.n, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

func BenchmarkMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Multiply(100, 200)
	}
}

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Factorial(10)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(997)
	}
}
