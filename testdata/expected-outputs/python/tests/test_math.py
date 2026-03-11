"""
Math Operations Tests

Tests for mathematical functions including basic arithmetic,
advanced math, statistics, and property-based tests.
"""

import pytest
import math
import random
from hypothesis import given, strategies as st, assume
from decimal import Decimal, getcontext
from math_utils import (
    # Basic arithmetic
    add,
    subtract,
    multiply,
    divide,
    
    # Advanced math
    power,
    sqrt,
    cbrt,
    nth_root,
    factorial,
    fibonacci,
    
    # Number theory
    is_prime,
    gcd,
    lcm,
    prime_factors,
    is_perfect_square,
    
    # Trigonometry
    sin,
    cos,
    tan,
    arcsin,
    arccos,
    arctan,
    
    # Logarithms
    log,
    log10,
    log2,
    ln,
    
    # Statistics
    mean,
    median,
    mode,
    variance,
    std_dev,
    quartiles,
    percentile,
    
    # Combinatorics
    permutations,
    combinations,
    binomial_coefficient,
    
    # Geometry
    distance,
    midpoint,
    slope,
    circle_area,
    sphere_volume,
    
    # Utility
    clamp,
    lerp,
    map_range,
    round_to,
    sign,
    is_even,
    is_odd,
    is_positive,
    is_negative,
    
    # Series
    arithmetic_series,
    geometric_series,
    fibonacci_series,
    
    # Calculus
    derivative,
    integral,
    limit,
    
    # Complex numbers
    complex_add,
    complex_multiply,
    complex_modulus,
    
    # Matrix operations
    matrix_add,
    matrix_multiply,
    matrix_transpose,
    determinant,
    
    # Constants
    PI,
    E,
    TAU,
    PHI
)

# ============================================================================
# Basic Arithmetic Tests
# ============================================================================

class TestBasicArithmetic:
    """Tests for basic arithmetic operations."""

    def test_add(self):
        """Test addition of numbers."""
        assert add(2, 3) == 5
        assert add(-2, -3) == -5
        assert add(-2, 3) == 1
        assert add(0, 5) == 5
        assert add(0, 0) == 0
        assert add(1.5, 2.5) == 4.0
        assert add(1e10, 2e10) == 3e10

    def test_add_multiple(self):
        """Test addition of multiple numbers."""
        assert add(1, 2, 3, 4, 5) == 15
        assert add(10) == 10
        assert add() == 0

    def test_add_precision(self):
        """Test addition with floating point precision."""
        result = add(0.1, 0.2)
        assert abs(result - 0.3) < 1e-12

    def test_subtract(self):
        """Test subtraction of numbers."""
        assert subtract(10, 4) == 6
        assert subtract(4, 10) == -6
        assert subtract(-10, -4) == -6
        assert subtract(-10, 4) == -14
        assert subtract(0, 5) == -5
        assert subtract(5, 0) == 5
        assert subtract(1.5, 0.5) == 1.0

    def test_subtract_multiple(self):
        """Test subtraction of multiple numbers."""
        assert subtract(100, 20, 30, 10) == 40
        assert subtract(10) == 10

    def test_multiply(self):
        """Test multiplication of numbers."""
        assert multiply(2, 3) == 6
        assert multiply(-2, 3) == -6
        assert multiply(-2, -3) == 6
        assert multiply(0, 5) == 0
        assert multiply(1.5, 2) == 3.0
        assert multiply(1e6, 1e6) == 1e12

    def test_multiply_multiple(self):
        """Test multiplication of multiple numbers."""
        assert multiply(2, 3, 4) == 24
        assert multiply(2, 3, 4, 5) == 120
        assert multiply(10) == 10

    def test_divide(self):
        """Test division of numbers."""
        assert divide(10, 2) == 5.0
        assert divide(9, 3) == 3.0
        assert divide(7, 2) == 3.5
        assert divide(-10, 2) == -5.0
        assert divide(10, -2) == -5.0
        assert divide(-10, -2) == 5.0
        assert divide(0, 5) == 0.0

    def test_divide_by_zero(self):
        """Test division by zero raises error."""
        with pytest.raises(ZeroDivisionError):
            divide(10, 0)

    def test_divide_multiple(self):
        """Test sequential division."""
        assert divide(100, 2, 5) == 10.0
        assert divide(100, 2, 2, 5) == 5.0

# ============================================================================
# Advanced Math Tests
# ============================================================================

class TestAdvancedMath:
    """Tests for advanced mathematical operations."""

    def test_power(self):
        """Test exponentiation."""
        assert power(2, 3) == 8
        assert power(5, 0) == 1
        assert power(2, -1) == 0.5
        assert power(4, 0.5) == 2.0
        assert power(-2, 2) == 4
        assert power(-2, 3) == -8
        assert power(10, 6) == 1_000_000

    @pytest.mark.parametrize("x,expected", [
        (4, 2),
        (9, 3),
        (16, 4),
        (25, 5),
        (2, math.sqrt(2)),
        (0, 0),
        (1, 1)
    ])
    def test_sqrt(self, x, expected):
        """Test square root with parametrized inputs."""
        result = sqrt(x)
        if isinstance(expected, float):
            assert abs(result - expected) < 1e-12
        else:
            assert result == expected

    def test_sqrt_negative(self):
        """Test square root of negative number."""
        with pytest.raises(ValueError, match="Cannot calculate square root of negative number"):
            sqrt(-1)

    @pytest.mark.parametrize("x,expected", [
        (8, 2),
        (27, 3),
        (64, 4),
        (125, 5),
        (2, 2 ** (1/3)),
        (0, 0),
        (1, 1),
        (-8, -2)
    ])
    def test_cbrt(self, x, expected):
        """Test cube root."""
        result = cbrt(x)
        if isinstance(expected, float):
            assert abs(result - expected) < 1e-12
        else:
            assert result == expected

    def test_nth_root(self):
        """Test nth root."""
        assert nth_root(16, 4) == 2
        assert nth_root(27, 3) == 3
        assert nth_root(32, 5) == 2
        assert abs(nth_root(2, 2) - math.sqrt(2)) < 1e-12

    def test_factorial(self):
        """Test factorial."""
        assert factorial(0) == 1
        assert factorial(1) == 1
        assert factorial(5) == 120
        assert factorial(10) == 3_628_800
        
        with pytest.raises(ValueError, match="Factorial not defined for negative numbers"):
            factorial(-1)

    def test_fibonacci(self):
        """Test Fibonacci numbers."""
        fib_sequence = [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]
        for i, expected in enumerate(fib_sequence):
            assert fibonacci(i) == expected
        
        with pytest.raises(ValueError, match="Fibonacci not defined for negative numbers"):
            fibonacci(-1)

# ============================================================================
# Number Theory Tests
# ============================================================================

class TestNumberTheory:
    """Tests for number theory functions."""

    @pytest.mark.parametrize("n,expected", [
        (2, True),
        (3, True),
        (4, False),
        (5, True),
        (9, False),
        (17, True),
        (21, False),
        (97, True),
        (100, False)
    ])
    def test_is_prime(self, n, expected):
        """Test prime number detection."""
        assert is_prime(n) == expected

    def test_gcd(self):
        """Test greatest common divisor."""
        assert gcd(48, 18) == 6
        assert gcd(12, 8) == 4
        assert gcd(17, 19) == 1
        assert gcd(0, 5) == 5
        assert gcd(5, 0) == 5
        assert gcd(0, 0) == 0
        assert gcd(-48, 18) == 6

    def test_lcm(self):
        """Test least common multiple."""
        assert lcm(12, 18) == 36
        assert lcm(8, 12) == 24
        assert lcm(17, 19) == 323
        assert lcm(0, 5) == 0
        assert lcm(5, 0) == 0
        assert lcm(-12, 18) == 36

    def test_prime_factors(self):
        """Test prime factorization."""
        assert prime_factors(12) == [2, 2, 3]
        assert prime_factors(13) == [13]
        assert prime_factors(1) == []
        assert prime_factors(100) == [2, 2, 5, 5]

    def test_is_perfect_square(self):
        """Test perfect square detection."""
        assert is_perfect_square(4) is True
        assert is_perfect_square(9) is True
        assert is_perfect_square(16) is True
        assert is_perfect_square(2) is False
        assert is_perfect_square(0) is True
        assert is_perfect_square(-4) is False

# ============================================================================
# Trigonometry Tests
# ============================================================================

class TestTrigonometry:
    """Tests for trigonometric functions."""

    def test_sin(self):
        """Test sine function."""
        assert abs(sin(0)) < 1e-12
        assert abs(sin(PI/2) - 1) < 1e-12
        assert abs(sin(PI)) < 1e-12
        assert abs(sin(3*PI/2) + 1) < 1e-12
        assert abs(sin(2*PI)) < 1e-12

    def test_cos(self):
        """Test cosine function."""
        assert abs(cos(0) - 1) < 1e-12
        assert abs(cos(PI/2)) < 1e-12
        assert abs(cos(PI) + 1) < 1e-12
        assert abs(cos(3*PI/2)) < 1e-12
        assert abs(cos(2*PI) - 1) < 1e-12

    def test_tan(self):
        """Test tangent function."""
        assert abs(tan(0)) < 1e-12
        assert abs(tan(PI/4) - 1) < 1e-12
        assert abs(tan(PI)) < 1e-12

    def test_arcsin(self):
        """Test arcsine function."""
        assert abs(arcsin(0)) < 1e-12
        assert abs(arcsin(1) - PI/2) < 1e-12
        assert abs(arcsin(-1) + PI/2) < 1e-12
        
        with pytest.raises(ValueError):
            arcsin(2)

    def test_arccos(self):
        """Test arccosine function."""
        assert abs(arccos(1)) < 1e-12
        assert abs(arccos(0) - PI/2) < 1e-12
        assert abs(arccos(-1) - PI) < 1e-12
        
        with pytest.raises(ValueError):
            arccos(2)

# ============================================================================
# Logarithm Tests
# ============================================================================

class TestLogarithms:
    """Tests for logarithmic functions."""

    def test_log(self):
        """Test logarithm with custom base."""
        assert log(1, 10) == 0
        assert log(10, 10) == 1
        assert log(100, 10) == 2
        assert abs(log(8, 2) - 3) < 1e-12
        
        with pytest.raises(ValueError):
            log(0, 10)
        with pytest.raises(ValueError):
            log(-1, 10)
        with pytest.raises(ValueError):
            log(10, 1)

    def test_ln(self):
        """Test natural logarithm."""
        assert ln(1) == 0
        assert abs(ln(E) - 1) < 1e-12
        assert abs(ln(E**2) - 2) < 1e-12

    def test_log10(self):
        """Test base-10 logarithm."""
        assert log10(1) == 0
        assert log10(10) == 1
        assert log10(100) == 2
        assert log10(1000) == 3

    def test_log2(self):
        """Test base-2 logarithm."""
        assert log2(1) == 0
        assert log2(2) == 1
        assert log2(4) == 2
        assert log2(8) == 3
        assert abs(log2(3) - math.log2(3)) < 1e-12

# ============================================================================
# Statistics Tests
# ============================================================================

class TestStatistics:
    """Tests for statistical functions."""

    def test_mean(self):
        """Test arithmetic mean."""
        assert mean([1, 2, 3, 4, 5]) == 3
        assert mean([10, 20, 30]) == 20
        assert mean([-1, 0, 1]) == 0
        assert mean([1.5, 2.5, 3.5]) == 2.5
        
        with pytest.raises(ValueError, match="Cannot calculate mean of empty list"):
            mean([])

    def test_median_odd(self):
        """Test median with odd-length list."""
        assert median([1, 3, 5]) == 3
        assert median([10, 20, 30, 40, 50]) == 30

    def test_median_even(self):
        """Test median with even-length list."""
        assert median([1, 2, 3, 4]) == 2.5
        assert median([10, 20, 30, 40]) == 25

    def test_median_unsorted(self):
        """Test median with unsorted list."""
        assert median([5, 1, 4, 2, 3]) == 3

    def test_mode(self):
        """Test mode."""
        assert mode([1, 2, 2, 3, 4]) == 2
        assert mode([1, 1, 2, 2, 3]) in [1, 2]  # Multiple modes
        assert mode([1, 2, 3, 4]) is None  # No mode

    def test_variance(self):
        """Test variance."""
        data = [1, 2, 3, 4, 5]
        assert variance(data) == 2.0
        assert variance(data, sample=True) == 2.5

    def test_std_dev(self):
        """Test standard deviation."""
        data = [1, 2, 3, 4, 5]
        assert abs(std_dev(data) - math.sqrt(2)) < 1e-12
        assert abs(std_dev(data, sample=True) - math.sqrt(2.5)) < 1e-12

    def test_quartiles(self):
        """Test quartiles."""
        data = [1, 2, 3, 4, 5, 6, 7, 8, 9]
        q1, q2, q3 = quartiles(data)
        assert q1 == 3
        assert q2 == 5
        assert q3 == 7

    def test_percentile(self):
        """Test percentile."""
        data = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
        assert percentile(data, 50) == 5.5
        assert percentile(data, 25) == 3.25
        assert percentile(data, 75) == 7.75

# ============================================================================
# Combinatorics Tests
# ============================================================================

class TestCombinatorics:
    """Tests for combinatorics functions."""

    def test_permutations(self):
        """Test permutations."""
        assert permutations(5, 3) == 60
        assert permutations(5, 0) == 1
        assert permutations(5, 5) == 120
        
        with pytest.raises(ValueError):
            permutations(5, 6)

    def test_combinations(self):
        """Test combinations."""
        assert combinations(5, 3) == 10
        assert combinations(5, 0) == 1
        assert combinations(5, 5) == 1
        
        with pytest.raises(ValueError):
            combinations(5, 6)

    def test_binomial_coefficient(self):
        """Test binomial coefficient."""
        assert binomial_coefficient(5, 3) == 10
        assert binomial_coefficient(10, 5) == 252
        assert binomial_coefficient(0, 0) == 1

# ============================================================================
# Geometry Tests
# ============================================================================

class TestGeometry:
    """Tests for geometry functions."""

    def test_distance(self):
        """Test distance between points."""
        assert distance(0, 0, 3, 4) == 5
        assert distance(1, 1, 4, 5) == 5
        assert distance(-1, -1, 2, 3) == 5

    def test_midpoint(self):
        """Test midpoint calculation."""
        assert midpoint(0, 0, 2, 2) == (1, 1)
        assert midpoint(-1, -1, 1, 1) == (0, 0)
        assert midpoint(1, 2, 3, 4) == (2, 3)

    def test_slope(self):
        """Test slope calculation."""
        assert slope(0, 0, 2, 2) == 1
        assert slope(0, 0, 2, 4) == 2
        assert slope(1, 1, 3, 5) == 2
        
        with pytest.raises(ValueError, match="Vertical line has undefined slope"):
            slope(1, 1, 1, 5)

    def test_circle_area(self):
        """Test circle area."""
        assert abs(circle_area(1) - PI) < 1e-12
        assert circle_area(0) == 0
        assert circle_area(2) == 4 * PI

    def test_sphere_volume(self):
        """Test sphere volume."""
        assert abs(sphere_volume(1) - (4/3) * PI) < 1e-12
        assert sphere_volume(0) == 0
        assert sphere_volume(2) == (32/3) * PI

# ============================================================================
# Utility Function Tests
# ============================================================================

class TestUtilityFunctions:
    """Tests for utility mathematical functions."""

    def test_clamp(self):
        """Test clamping values."""
        assert clamp(5, 1, 10) == 5
        assert clamp(0, 1, 10) == 1
        assert clamp(15, 1, 10) == 10
        assert clamp(-5, -10, -1) == -5

    def test_lerp(self):
        """Test linear interpolation."""
        assert lerp(0, 10, 0) == 0
        assert lerp(0, 10, 0.5) == 5
        assert lerp(0, 10, 1) == 10
        assert lerp(10, 20, 0.3) == 13

    def test_map_range(self):
        """Test mapping values between ranges."""
        assert map_range(5, 0, 10, 0, 100) == 50
        assert map_range(0, -10, 10, 0, 100) == 50
        assert map_range(10, 0, 10, 0, 100) == 100

    def test_round_to(self):
        """Test rounding to specific precision."""
        assert round_to(3.14159, 2) == 3.14
        assert round_to(3.14159, 4) == 3.1416
        assert round_to(123.456, -1) == 120.0

    def test_sign(self):
        """Test sign function."""
        assert sign(5) == 1
        assert sign(-5) == -1
        assert sign(0) == 0
        assert sign(3.14) == 1

    def test_is_even(self):
        """Test even number detection."""
        assert is_even(2) is True
        assert is_even(3) is False
        assert is_even(0) is True
        assert is_even(-2) is True
        assert is_even(-3) is False

    def test_is_odd(self):
        """Test odd number detection."""
        assert is_odd(2) is False
        assert is_odd(3) is True
        assert is_odd(0) is False
        assert is_odd(-2) is False
        assert is_odd(-3) is True

# ============================================================================
# Series Tests
# ============================================================================

class TestSeries:
    """Tests for mathematical series."""

    def test_arithmetic_series(self):
        """Test arithmetic series."""
        assert arithmetic_series(1, 1, 5) == 15
        assert arithmetic_series(1, 2, 5) == 25
        assert arithmetic_series(0, 1, 10) == 45

    def test_geometric_series(self):
        """Test geometric series."""
        assert geometric_series(1, 2, 5) == 31
        assert geometric_series(1, 1, 5) == 5
        assert geometric_series(2, 2, 4) == 30

    def test_fibonacci_series(self):
        """Test Fibonacci series generation."""
        assert fibonacci_series(5) == [0, 1, 1, 2, 3]
        assert fibonacci_series(10) == [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]

# ============================================================================
# Calculus Tests
# ============================================================================

class TestCalculus:
    """Tests for calculus functions."""

    def test_derivative(self):
        """Test numerical derivative."""
        f = lambda x: x**2
        assert abs(derivative(f, 2) - 4) < 1e-6
        
        f = lambda x: math.sin(x)
        assert abs(derivative(f, 0) - 1) < 1e-6

    def test_integral(self):
        """Test numerical integral."""
        f = lambda x: x
        assert abs(integral(f, 0, 1) - 0.5) < 1e-6
        
        f = lambda x: x**2
        assert abs(integral(f, 0, 1) - 1/3) < 1e-6

    def test_limit(self):
        """Test limit calculation."""
        f = lambda x: (x**2 - 1)/(x - 1)
        assert abs(limit(f, 1) - 2) < 1e-6

# ============================================================================
# Complex Number Tests
# ============================================================================

class TestComplexNumbers:
    """Tests for complex number operations."""

    def test_complex_add(self):
        """Test complex addition."""
        assert complex_add(1+2j, 3+4j) == 4+6j
        assert complex_add(1+1j, -1-1j) == 0

    def test_complex_multiply(self):
        """Test complex multiplication."""
        assert complex_multiply(1+2j, 3+4j) == -5+10j
        assert complex_multiply(1+1j, 1-1j) == 2

    def test_complex_modulus(self):
        """Test complex modulus."""
        assert complex_modulus(3+4j) == 5
        assert complex_modulus(1+1j) == math.sqrt(2)

# ============================================================================
# Matrix Tests
# ============================================================================

class TestMatrix:
    """Tests for matrix operations."""

    def test_matrix_add(self):
        """Test matrix addition."""
        A = [[1, 2], [3, 4]]
        B = [[5, 6], [7, 8]]
        result = matrix_add(A, B)
        assert result == [[6, 8], [10, 12]]

    def test_matrix_multiply(self):
        """Test matrix multiplication."""
        A = [[1, 2], [3, 4]]
        B = [[5, 6], [7, 8]]
        result = matrix_multiply(A, B)
        assert result == [[19, 22], [43, 50]]

    def test_matrix_transpose(self):
        """Test matrix transpose."""
        A = [[1, 2, 3], [4, 5, 6]]
        result = matrix_transpose(A)
        assert result == [[1, 4], [2, 5], [3, 6]]

    def test_determinant(self):
        """Test matrix determinant."""
        assert determinant([[1, 2], [3, 4]]) == -2
        assert determinant([[2, 0], [0, 2]]) == 4

# ============================================================================
# Property-Based Tests
# ============================================================================

class TestProperties:
    """Property-based tests using Hypothesis."""

    @given(st.floats(min_value=-1e6, max_value=1e6),
           st.floats(min_value=-1e6, max_value=1e6))
    def test_add_commutative(self, a, b):
        """Test that addition is commutative."""
        assume(not (math.isnan(a) or math.isnan(b)))
        assert add(a, b) == add(b, a)

    @given(st.floats(min_value=-1e6, max_value=1e6),
           st.floats(min_value=-1e6, max_value=1e6),
           st.floats(min_value=-1e6, max_value=1e6))
    def test_add_associative(self, a, b, c):
        """Test that addition is associative."""
        assume(not any(math.isnan(x) for x in [a, b, c]))
        assert add(add(a, b), c) == add(a, add(b, c))

    @given(st.floats(min_value=-1e6, max_value=1e6))
    def test_add_identity(self, a):
        """Test addition identity property."""
        assume(not math.isnan(a))
        assert add(a, 0) == a

    @given(st.floats(min_value=-1e6, max_value=1e6),
           st.floats(min_value=-1e6, max_value=1e6))
    def test_multiply_commutative(self, a, b):
        """Test that multiplication is commutative."""
        assume(not (math.isnan(a) or math.isnan(b)))
        assert multiply(a, b) == multiply(b, a)

    @given(st.integers(min_value=1, max_value=100))
    def test_factorial_positive(self, n):
        """Test factorial properties."""
        result = factorial(n)
        assert result > 0
        assert result == n * factorial(n - 1)

    @given(st.integers(min_value=2, max_value=100))
    def test_is_prime_property(self, n):
        """Test prime number properties."""
        if is_prime(n):
            # Prime numbers are not divisible by any number 2..sqrt(n)
            for i in range(2, int(math.sqrt(n)) + 1):
                assert n % i != 0

    @given(st.lists(st.floats(min_value=-100, max_value=100), min_size=1, max_size=10))
    def test_mean_within_range(self, data):
        """Test that mean is within range of data."""
        assume(not any(math.isnan(x) for x in data))
        m = mean(data)
        assert min(data) <= m <= max(data)

# ============================================================================
# Edge Cases and Error Handling
# ============================================================================

class TestEdgeCases:
    """Tests for edge cases and error handling."""

    def test_infinite_inputs(self):
        """Test handling of infinite inputs."""
        with pytest.raises(ValueError):
            add(float('inf'), 5)
        
        with pytest.raises(ValueError):
            subtract(float('inf'), 5)

    def test_nan_inputs(self):
        """Test handling of NaN inputs."""
        with pytest.raises(ValueError):
            add(float('nan'), 5)

    def test_very_large_numbers(self):
        """Test handling of very large numbers."""
        large = 1e308
        result = add(large, large)
        assert result == float('inf')

    def test_very_small_numbers(self):
        """Test handling of very small numbers."""
        small = 1e-308
        result = multiply(small, small)
        assert result == 0.0

# ============================================================================
# Performance Tests
# ============================================================================

class TestPerformance:
    """Performance tests for mathematical operations."""

    @pytest.mark.benchmark
    def test_factorial_performance(self, benchmark):
        """Benchmark factorial calculation."""
        result = benchmark(factorial, 20)
        assert result == 2_432_902_008_176_640_000

    @pytest.mark.benchmark
    def test_prime_check_performance(self, benchmark):
        """Benchmark prime checking."""
        result = benchmark(is_prime, 99991)
        assert result is True

    @pytest.mark.benchmark
    def test_matrix_multiplication_performance(self, benchmark):
        """Benchmark matrix multiplication."""
        A = [[i + j for j in range(50)] for i in range(50)]
        B = [[i * j for j in range(50)] for i in range(50)]
        result = benchmark(matrix_multiply, A, B)
        assert len(result) == 50

# ============================================================================
# Fixtures
# ============================================================================

@pytest.fixture
def sample_numbers():
    """Provide sample numbers for testing."""
    return [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

@pytest.fixture
def sample_floats():
    """Provide sample floating point numbers."""
    return [1.1, 2.2, 3.3, 4.4, 5.5]

@pytest.fixture
def large_dataset():
    """Provide a large dataset for performance testing."""
    return list(range(10000))

@pytest.fixture
def matrix_2x2():
    """Provide a 2x2 matrix."""
    return [[1, 2], [3, 4]]

@pytest.fixture
def matrix_3x3():
    """Provide a 3x3 matrix."""
    return [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9]
    ]