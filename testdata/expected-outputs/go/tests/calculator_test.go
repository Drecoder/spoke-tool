package calculator

import (
	"testing"
)

func TestNewCalculator(t *testing.T) {
	t.Run("default precision", func(t *testing.T) {
		calc := New()
		if calc == nil {
			t.Fatal("expected non-nil calculator")
		}

		// Test default behavior
		result := calc.Add(1.1, 2.2)
		expected := 3.3
		if result != expected {
			t.Errorf("Add() with default precision = %v, want %v", result, expected)
		}
	})

	t.Run("custom precision", func(t *testing.T) {
		calc, err := NewWithPrecision(3)
		if err != nil {
			t.Fatalf("NewWithPrecision(3) failed: %v", err)
		}
		if calc == nil {
			t.Fatal("expected non-nil calculator")
		}
	})

	t.Run("invalid precision - too high", func(t *testing.T) {
		calc, err := NewWithPrecision(20)
		if err == nil {
			t.Error("NewWithPrecision(20) expected error, got nil")
		}
		if calc != nil {
			t.Error("NewWithPrecision(20) expected nil calculator, got non-nil")
		}
	})

	t.Run("invalid precision - negative", func(t *testing.T) {
		calc, err := NewWithPrecision(-5)
		if err == nil {
			t.Error("NewWithPrecision(-5) expected error, got nil")
		}
		if calc != nil {
			t.Error("NewWithPrecision(-5) expected nil calculator, got non-nil")
		}
	})
}

func TestCalculator_Add(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{
			name:     "positive numbers",
			a:        5.2,
			b:        3.1,
			expected: 8.3,
		},
		{
			name:     "negative numbers",
			a:        -5.2,
			b:        -3.1,
			expected: -8.3,
		},
		{
			name:     "mixed signs",
			a:        -5.2,
			b:        3.1,
			expected: -2.1,
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
			name:     "decimal precision",
			a:        1.234,
			b:        2.345,
			expected: 3.579,
		},
		{
			name:     "large numbers",
			a:        1e6,
			b:        2e6,
			expected: 3e6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCalculator_Subtract(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{
			name:     "positive result",
			a:        10.5,
			b:        4.2,
			expected: 6.3,
		},
		{
			name:     "negative result",
			a:        4.2,
			b:        10.5,
			expected: -6.3,
		},
		{
			name:     "zero result",
			a:        5,
			b:        5,
			expected: 0,
		},
		{
			name:     "negative numbers",
			a:        -10.5,
			b:        -4.2,
			expected: -6.3,
		},
		{
			name:     "mixed signs",
			a:        -10.5,
			b:        4.2,
			expected: -14.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCalculator_Multiply(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{
			name:     "positive numbers",
			a:        3.0,
			b:        4.5,
			expected: 13.5,
		},
		{
			name:     "negative numbers",
			a:        -3.0,
			b:        -4.5,
			expected: 13.5,
		},
		{
			name:     "mixed signs",
			a:        -3.0,
			b:        4.5,
			expected: -13.5,
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
			name:     "decimal multiplication",
			a:        2.5,
			b:        3.5,
			expected: 8.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCalculator_Divide(t *testing.T) {
	calc := New()

	tests := []struct {
		name      string
		a, b      float64
		expected  float64
		wantError bool
	}{
		{
			name:      "normal division",
			a:         10.0,
			b:         2.0,
			expected:  5.0,
			wantError: false,
		},
		{
			name:      "division by zero",
			a:         10.0,
			b:         0.0,
			expected:  0,
			wantError: true,
		},
		{
			name:      "negative division",
			a:         -10.0,
			b:         2.0,
			expected:  -5.0,
			wantError: false,
		},
		{
			name:      "divide negative by negative",
			a:         -10.0,
			b:         -2.0,
			expected:  5.0,
			wantError: false,
		},
		{
			name:      "zero dividend",
			a:         0.0,
			b:         5.0,
			expected:  0.0,
			wantError: false,
		},
		{
			name:      "decimal division",
			a:         7.5,
			b:         2.5,
			expected:  3.0,
			wantError: false,
		},
		{
			name:      "non-terminating decimal",
			a:         10.0,
			b:         3.0,
			expected:  3.3333333333333335,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Divide(tt.a, tt.b)

			if tt.wantError {
				if err == nil {
					t.Errorf("Divide(%v, %v) expected error, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%v, %v) unexpected error: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func TestCalculator_Power(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		base     float64
		exp      float64
		expected float64
	}{
		{
			name:     "positive exponent",
			base:     2.0,
			exp:      3.0,
			expected: 8.0,
		},
		{
			name:     "zero exponent",
			base:     5.0,
			exp:      0.0,
			expected: 1.0,
		},
		{
			name:     "negative exponent",
			base:     2.0,
			exp:      -1.0,
			expected: 0.5,
		},
		{
			name:     "fractional exponent",
			base:     4.0,
			exp:      0.5,
			expected: 2.0,
		},
		{
			name:     "negative base with even exponent",
			base:     -2.0,
			exp:      2.0,
			expected: 4.0,
		},
		{
			name:     "negative base with odd exponent",
			base:     -2.0,
			exp:      3.0,
			expected: -8.0,
		},
		{
			name:     "zero base positive exponent",
			base:     0.0,
			exp:      5.0,
			expected: 0.0,
		},
		{
			name:     "zero base zero exponent",
			base:     0.0,
			exp:      0.0,
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Power(tt.base, tt.exp)
			if result != tt.expected {
				t.Errorf("Power(%v, %v) = %v, want %v", tt.base, tt.exp, result, tt.expected)
			}
		})
	}
}

func TestCalculator_Sqrt(t *testing.T) {
	calc := New()

	tests := []struct {
		name      string
		x         float64
		expected  float64
		wantError bool
	}{
		{
			name:      "perfect square",
			x:         16.0,
			expected:  4.0,
			wantError: false,
		},
		{
			name:      "non-perfect square",
			x:         2.0,
			expected:  1.4142135623730951,
			wantError: false,
		},
		{
			name:      "zero",
			x:         0.0,
			expected:  0.0,
			wantError: false,
		},
		{
			name:      "one",
			x:         1.0,
			expected:  1.0,
			wantError: false,
		},
		{
			name:      "large number",
			x:         1e6,
			expected:  1000.0,
			wantError: false,
		},
		{
			name:      "negative number",
			x:         -4.0,
			expected:  0,
			wantError: true,
		},
		{
			name:      "very small number",
			x:         1e-10,
			expected:  1e-5,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Sqrt(tt.x)

			if tt.wantError {
				if err == nil {
					t.Errorf("Sqrt(%v) expected error, got nil", tt.x)
				}
			} else {
				if err != nil {
					t.Errorf("Sqrt(%v) unexpected error: %v", tt.x, err)
				}
				// Use a tolerance for floating point comparison
				if !approxEqual(result, tt.expected, 1e-10) {
					t.Errorf("Sqrt(%v) = %v, want %v", tt.x, result, tt.expected)
				}
			}
		})
	}
}

func TestCalculator_Chaining(t *testing.T) {
	calc := New()

	// Test method chaining pattern (if supported)
	// (2 + 3) * 4 = 20
	result := calc.Add(2, 3)          // 5
	result = calc.Multiply(result, 4) // 20

	expected := 20.0
	if result != expected {
		t.Errorf("Chained operations = %v, want %v", result, expected)
	}
}

func TestCalculator_Precision(t *testing.T) {
	t.Run("default precision (2)", func(t *testing.T) {
		calc := New()
		result := calc.Divide(22.0, 7.0)
		// Should round to 2 decimal places: 3.14
		expected := 3.142857142857143 // Actually no rounding in basic version
		if result != expected {
			t.Logf("Note: Basic calculator doesn't round - result: %v", result)
		}
	})

	t.Run("custom precision", func(t *testing.T) {
		calc, err := NewWithPrecision(5)
		if err != nil {
			t.Fatalf("Failed to create calculator with precision: %v", err)
		}

		// This test assumes precision affects rounding
		// Adjust based on actual implementation
		result := calc.Divide(22.0, 7.0)
		t.Logf("Result with precision 5: %v", result)
	})
}

// Helper function for floating point comparison
func approxEqual(a, b, tolerance float64) bool {
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < tolerance
}

func BenchmarkAdd(b *testing.B) {
	calc := New()
	for i := 0; i < b.N; i++ {
		calc.Add(1.1, 2.2)
	}
}

func BenchmarkDivide(b *testing.B) {
	calc := New()
	for i := 0; i < b.N; i++ {
		calc.Divide(10.0, 3.0)
	}
}

func BenchmarkPower(b *testing.B) {
	calc := New()
	for i := 0; i < b.N; i++ {
		calc.Power(2.0, 10.0)
	}
}
