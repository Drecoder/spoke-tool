package main

import (
	"fmt"
	"math"
)

// ============================================================================
// Basic Arithmetic Functions
// ============================================================================

// Add returns the sum of two integers
func Add(a, b int) int {
	return a + b
}

// AddFloat returns the sum of two float64 values
func AddFloat(a, b float64) float64 {
	return a + b
}

// Subtract returns the difference between two integers
func Subtract(a, b int) int {
	return a - b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
	return a * b
}

// Divide returns the quotient of two integers (integer division)
func Divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

// DivideFloat returns the quotient of two float64 values
func DivideFloat(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Mod returns the remainder of a divided by b
func Mod(a, b int) int {
	if b == 0 {
		return 0
	}
	return a % b
}

// ============================================================================
// Power and Root Functions
// ============================================================================

// Power returns a raised to the power of b
func Power(a, b float64) float64 {
	return math.Pow(a, b)
}

// Square returns the square of a number
func Square(x float64) float64 {
	return x * x
}

// Cube returns the cube of a number
func Cube(x float64) float64 {
	return x * x * x
}

// Sqrt returns the square root of x
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

// Cbrt returns the cube root of x
func Cbrt(x float64) float64 {
	return math.Cbrt(x)
}

// Hypot returns the square root of a² + b²
func Hypot(a, b float64) float64 {
	return math.Hypot(a, b)
}

// ============================================================================
// Trigonometric Functions
// ============================================================================

// Sin returns the sine of x (in radians)
func Sin(x float64) float64 {
	return math.Sin(x)
}

// Cos returns the cosine of x (in radians)
func Cos(x float64) float64 {
	return math.Cos(x)
}

// Tan returns the tangent of x (in radians)
func Tan(x float64) float64 {
	return math.Tan(x)
}

// Asin returns the arcsine of x
func Asin(x float64) float64 {
	return math.Asin(x)
}

// Acos returns the arccosine of x
func Acos(x float64) float64 {
	return math.Acos(x)
}

// Atan returns the arctangent of x
func Atan(x float64) float64 {
	return math.Atan(x)
}

// Atan2 returns the arctangent of y/x
func Atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

// ============================================================================
// Logarithmic and Exponential Functions
// ============================================================================

// Exp returns e raised to the power of x
func Exp(x float64) float64 {
	return math.Exp(x)
}

// Log returns the natural logarithm of x
func Log(x float64) float64 {
	return math.Log(x)
}

// Log10 returns the base-10 logarithm of x
func Log10(x float64) float64 {
	return math.Log10(x)
}

// Log2 returns the base-2 logarithm of x
func Log2(x float64) float64 {
	return math.Log2(x)
}

// ============================================================================
// Rounding Functions
// ============================================================================

// Round returns the nearest integer, rounding half away from zero
func Round(x float64) float64 {
	return math.Round(x)
}

// RoundToEven returns the nearest integer, rounding ties to even
func RoundToEven(x float64) float64 {
	return math.RoundToEven(x)
}

// Floor returns the greatest integer value less than or equal to x
func Floor(x float64) float64 {
	return math.Floor(x)
}

// Ceil returns the least integer value greater than or equal to x
func Ceil(x float64) float64 {
	return math.Ceil(x)
}

// Trunc returns the integer value of x
func Trunc(x float64) float64 {
	return math.Trunc(x)
}

// ============================================================================
// Absolute Value and Sign Functions
// ============================================================================

// Abs returns the absolute value of x
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// AbsFloat returns the absolute value of x (float64)
func AbsFloat(x float64) float64 {
	return math.Abs(x)
}

// Sign returns the sign of x (-1, 0, 1)
func Sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

// IsPositive returns true if x is positive
func IsPositive(x int) bool {
	return x > 0
}

// IsNegative returns true if x is negative
func IsNegative(x int) bool {
	return x < 0
}

// IsZero returns true if x is zero
func IsZero(x int) bool {
	return x == 0
}

// ============================================================================
// Min/Max Functions
// ============================================================================

// Min returns the smaller of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinFloat returns the smaller of two float64 values
func MinFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// MaxFloat returns the larger of two float64 values
func MaxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinOf returns the minimum of a slice
func MinOf(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	min := nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
	}
	return min
}

// MaxOf returns the maximum of a slice
func MaxOf(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	max := nums[0]
	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}
	return max
}

// ============================================================================
// Statistical Functions
// ============================================================================

// Sum returns the sum of a slice
func Sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Average returns the average of a slice
func Average(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	return float64(Sum(nums)) / float64(len(nums))
}

// Product returns the product of a slice
func Product(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	result := 1
	for _, n := range nums {
		result *= n
	}
	return result
}

// ============================================================================
// Number Theory Functions
// ============================================================================

// IsEven returns true if n is even
func IsEven(n int) bool {
	return n%2 == 0
}

// IsOdd returns true if n is odd
func IsOdd(n int) bool {
	return n%2 != 0
}

// IsPrime checks if a number is prime
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// Factorial returns n!
func Factorial(n int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// GCD returns the greatest common divisor
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM returns the least common multiple
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return a * b / GCD(a, b)
}

// Fibonacci returns the nth Fibonacci number
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// ============================================================================
// Geometry Functions
// ============================================================================

// AreaOfCircle returns the area of a circle
func AreaOfCircle(radius float64) float64 {
	return math.Pi * radius * radius
}

// Circumference returns the circumference of a circle
func Circumference(radius float64) float64 {
	return 2 * math.Pi * radius
}

// AreaOfRectangle returns the area of a rectangle
func AreaOfRectangle(width, height float64) float64 {
	return width * height
}

// PerimeterOfRectangle returns the perimeter of a rectangle
func PerimeterOfRectangle(width, height float64) float64 {
	return 2 * (width + height)
}

// Distance returns the distance between two points
func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

// Midpoint returns the midpoint between two points
func Midpoint(x1, y1, x2, y2 float64) (float64, float64) {
	return (x1 + x2) / 2, (y1 + y2) / 2
}

// Slope returns the slope between two points
func Slope(x1, y1, x2, y2 float64) float64 {
	if x2-x1 == 0 {
		return math.Inf(1)
	}
	return (y2 - y1) / (x2 - x1)
}

// ============================================================================
// Conversion Functions
// ============================================================================

// DegreesToRadians converts degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// RadiansToDegrees converts radians to degrees
func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// FahrenheitToCelsius converts Fahrenheit to Celsius
func FahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

// CelsiusToFahrenheit converts Celsius to Fahrenheit
func CelsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

// ============================================================================
// Main Function
// ============================================================================

func main() {
	// Basic arithmetic
	fmt.Println("=== Basic Arithmetic ===")
	fmt.Printf("Add(5, 3) = %d\n", Add(5, 3))
	fmt.Printf("Subtract(10, 4) = %d\n", Subtract(10, 4))
	fmt.Printf("Multiply(6, 7) = %d\n", Multiply(6, 7))
	fmt.Printf("Divide(10, 3) = %d\n", Divide(10, 3))
	fmt.Printf("DivideFloat(10, 3) = %f\n", DivideFloat(10, 3))
	fmt.Printf("Mod(10, 3) = %d\n", Mod(10, 3))

	// Powers and roots
	fmt.Println("\n=== Powers and Roots ===")
	fmt.Printf("Power(2, 8) = %f\n", Power(2, 8))
	fmt.Printf("Square(5) = %f\n", Square(5))
	fmt.Printf("Cube(3) = %f\n", Cube(3))
	fmt.Printf("Sqrt(16) = %f\n", Sqrt(16))
	fmt.Printf("Cbrt(27) = %f\n", Cbrt(27))
	fmt.Printf("Hypot(3, 4) = %f\n", Hypot(3, 4))

	// Trigonometry
	fmt.Println("\n=== Trigonometry ===")
	angle := DegreesToRadians(30)
	fmt.Printf("Sin(30°) = %f\n", Sin(angle))
	fmt.Printf("Cos(30°) = %f\n", Cos(angle))
	fmt.Printf("Tan(30°) = %f\n", Tan(angle))

	// Logarithms
	fmt.Println("\n=== Logarithms ===")
	fmt.Printf("Exp(1) = %f\n", Exp(1))
	fmt.Printf("Log(10) = %f\n", Log(10))
	fmt.Printf("Log10(100) = %f\n", Log10(100))
	fmt.Printf("Log2(8) = %f\n", Log2(8))

	// Rounding
	fmt.Println("\n=== Rounding ===")
	fmt.Printf("Round(3.49) = %f\n", Round(3.49))
	fmt.Printf("Round(3.5) = %f\n", Round(3.5))
	fmt.Printf("Floor(3.9) = %f\n", Floor(3.9))
	fmt.Printf("Ceil(3.1) = %f\n", Ceil(3.1))

	// Absolute value
	fmt.Println("\n=== Absolute Value ===")
	fmt.Printf("Abs(-42) = %d\n", Abs(-42))
	fmt.Printf("Sign(-10) = %d\n", Sign(-10))
	fmt.Printf("IsPositive(5) = %v\n", IsPositive(5))
	fmt.Printf("IsNegative(-3) = %v\n", IsNegative(-3))

	// Min/Max
	fmt.Println("\n=== Min/Max ===")
	fmt.Printf("Min(5, 10) = %d\n", Min(5, 10))
	fmt.Printf("Max(5, 10) = %d\n", Max(5, 10))

	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Printf("Min of %v = %d\n", nums, MinOf(nums))
	fmt.Printf("Max of %v = %d\n", nums, MaxOf(nums))

	// Statistics
	fmt.Println("\n=== Statistics ===")
	fmt.Printf("Sum of %v = %d\n", nums, Sum(nums))
	fmt.Printf("Average of %v = %f\n", nums, Average(nums))
	fmt.Printf("Product of %v = %d\n", nums, Product(nums))

	// Number theory
	fmt.Println("\n=== Number Theory ===")
	fmt.Printf("IsEven(42) = %v\n", IsEven(42))
	fmt.Printf("IsOdd(43) = %v\n", IsOdd(43))
	fmt.Printf("IsPrime(17) = %v\n", IsPrime(17))
	fmt.Printf("IsPrime(21) = %v\n", IsPrime(21))
	fmt.Printf("Factorial(5) = %d\n", Factorial(5))
	fmt.Printf("GCD(48, 18) = %d\n", GCD(48, 18))
	fmt.Printf("LCM(12, 18) = %d\n", LCM(12, 18))
	fmt.Printf("Fibonacci(10) = %d\n", Fibonacci(10))

	// Geometry
	fmt.Println("\n=== Geometry ===")
	fmt.Printf("Area of circle (r=5) = %f\n", AreaOfCircle(5))
	fmt.Printf("Distance between (0,0) and (3,4) = %f\n", Distance(0, 0, 3, 4))

	// Conversions
	fmt.Println("\n=== Conversions ===")
	fmt.Printf("30° in radians = %f\n", DegreesToRadians(30))
	fmt.Printf("π/6 in degrees = %f\n", RadiansToDegrees(math.Pi/6))
	fmt.Printf("98.6°F in Celsius = %f\n", FahrenheitToCelsius(98.6))
	fmt.Printf("37°C in Fahrenheit = %f\n", CelsiusToFahrenheit(37))
}
