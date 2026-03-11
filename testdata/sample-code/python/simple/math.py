"""
Simple Math Module

Basic mathematical operations demonstrating functions, error handling,
and mathematical concepts in Python.
"""

import math
from typing import List, Union

Number = Union[int, float]


# ============================================================================
# Basic Arithmetic
# ============================================================================

def add(a: Number, b: Number) -> Number:
    """
    Add two numbers together.
    
    Args:
        a: First number
        b: Second number
    
    Returns:
        Sum of a and b
    
    Examples:
        >>> add(5, 3)
        8
        >>> add(2.5, 1.5)
        4.0
    """
    return a + b


def subtract(a: Number, b: Number) -> Number:
    """
    Subtract second number from first.
    
    Args:
        a: First number
        b: Second number
    
    Returns:
        Difference a - b
    """
    return a - b


def multiply(a: Number, b: Number) -> Number:
    """
    Multiply two numbers.
    
    Args:
        a: First number
        b: Second number
    
    Returns:
        Product a * b
    """
    return a * b


def divide(a: Number, b: Number) -> float:
    """
    Divide first number by second.
    
    Args:
        a: Dividend
        b: Divisor
    
    Returns:
        Quotient a / b
    
    Raises:
        ValueError: If b is zero
    """
    if b == 0:
        raise ValueError("Cannot divide by zero")
    return a / b


def power(base: Number, exponent: Number) -> Number:
    """
    Raise base to exponent.
    
    Args:
        base: Base number
        exponent: Exponent
    
    Returns:
        base raised to exponent
    """
    return base ** exponent


def sqrt(x: Number) -> float:
    """
    Calculate square root of a number.
    
    Args:
        x: Number to find square root of
    
    Returns:
        Square root of x
    
    Raises:
        ValueError: If x is negative
    """
    if x < 0:
        raise ValueError("Cannot calculate square root of negative number")
    return math.sqrt(x)


# ============================================================================
# Number Theory
# ============================================================================

def is_even(n: int) -> bool:
    """
    Check if a number is even.
    
    Args:
        n: Integer to check
    
    Returns:
        True if even, False otherwise
    """
    return n % 2 == 0


def is_odd(n: int) -> bool:
    """
    Check if a number is odd.
    
    Args:
        n: Integer to check
    
    Returns:
        True if odd, False otherwise
    """
    return n % 2 != 0


def factorial(n: int) -> int:
    """
    Calculate factorial of a non-negative integer.
    
    Args:
        n: Non-negative integer
    
    Returns:
        n! (n factorial)
    
    Raises:
        ValueError: If n is negative
    """
    if n < 0:
        raise ValueError("Factorial not defined for negative numbers")
    if n <= 1:
        return 1
    
    result = 1
    for i in range(2, n + 1):
        result *= i
    return result


def fibonacci(n: int) -> int:
    """
    Calculate the nth Fibonacci number.
    
    Args:
        n: Position in Fibonacci sequence (0-indexed)
    
    Returns:
        nth Fibonacci number
    
    Raises:
        ValueError: If n is negative
    """
    if n < 0:
        raise ValueError("Fibonacci not defined for negative numbers")
    if n <= 1:
        return n
    
    a, b = 0, 1
    for _ in range(2, n + 1):
        a, b = b, a + b
    return b


def is_prime(n: int) -> bool:
    """
    Check if a number is prime.
    
    Args:
        n: Number to check
    
    Returns:
        True if prime, False otherwise
    """
    if n <= 1:
        return False
    if n <= 3:
        return True
    if n % 2 == 0 or n % 3 == 0:
        return False
    
    i = 5
    while i * i <= n:
        if n % i == 0 or n % (i + 2) == 0:
            return False
        i += 6
    return True


def gcd(a: int, b: int) -> int:
    """
    Calculate greatest common divisor of two integers.
    
    Args:
        a: First integer
        b: Second integer
    
    Returns:
        GCD of a and b
    """
    a, b = abs(a), abs(b)
    while b:
        a, b = b, a % b
    return a


def lcm(a: int, b: int) -> int:
    """
    Calculate least common multiple of two integers.
    
    Args:
        a: First integer
        b: Second integer
    
    Returns:
        LCM of a and b
    """
    if a == 0 or b == 0:
        return 0
    return abs(a * b) // gcd(a, b)


# ============================================================================
# List Operations
# ============================================================================

def min_of(numbers: List[Number]) -> Number:
    """
    Find minimum value in a list.
    
    Args:
        numbers: List of numbers
    
    Returns:
        Minimum value
    
    Raises:
        ValueError: If list is empty
    """
    if not numbers:
        raise ValueError("Cannot find minimum of empty list")
    
    min_val = numbers[0]
    for num in numbers[1:]:
        if num < min_val:
            min_val = num
    return min_val


def max_of(numbers: List[Number]) -> Number:
    """
    Find maximum value in a list.
    
    Args:
        numbers: List of numbers
    
    Returns:
        Maximum value
    
    Raises:
        ValueError: If list is empty
    """
    if not numbers:
        raise ValueError("Cannot find maximum of empty list")
    
    max_val = numbers[0]
    for num in numbers[1:]:
        if num > max_val:
            max_val = num
    return max_val


def sum_of(numbers: List[Number]) -> Number:
    """
    Calculate sum of all numbers in a list.
    
    Args:
        numbers: List of numbers
    
    Returns:
        Sum of all numbers
    """
    total = 0
    for num in numbers:
        total += num
    return total


def average_of(numbers: List[Number]) -> float:
    """
    Calculate average of numbers in a list.
    
    Args:
        numbers: List of numbers
    
    Returns:
        Average value
    
    Raises:
        ValueError: If list is empty
    """
    if not numbers:
        raise ValueError("Cannot calculate average of empty list")
    return sum_of(numbers) / len(numbers)


def median_of(numbers: List[Number]) -> float:
    """
    Calculate median of numbers in a list.
    
    Args:
        numbers: List of numbers
    
    Returns:
        Median value
    
    Raises:
        ValueError: If list is empty
    """
    if not numbers:
        raise ValueError("Cannot calculate median of empty list")
    
    sorted_nums = sorted(numbers)
    n = len(sorted_nums)
    mid = n // 2
    
    if n % 2 == 0:
        return (sorted_nums[mid - 1] + sorted_nums[mid]) / 2
    return float(sorted_nums[mid])


# ============================================================================
# Geometry
# ============================================================================

def area_of_circle(radius: float) -> float:
    """
    Calculate area of a circle.
    
    Args:
        radius: Circle radius
    
    Returns:
        Area of circle
    
    Raises:
        ValueError: If radius is negative
    """
    if radius < 0:
        raise ValueError("Radius cannot be negative")
    return math.pi * radius * radius


def circumference(radius: float) -> float:
    """
    Calculate circumference of a circle.
    
    Args:
        radius: Circle radius
    
    Returns:
        Circumference
    
    Raises:
        ValueError: If radius is negative
    """
    if radius < 0:
        raise ValueError("Radius cannot be negative")
    return 2 * math.pi * radius


def area_of_rectangle(width: float, height: float) -> float:
    """
    Calculate area of a rectangle.
    
    Args:
        width: Rectangle width
        height: Rectangle height
    
    Returns:
        Area of rectangle
    
    Raises:
        ValueError: If width or height is negative
    """
    if width < 0 or height < 0:
        raise ValueError("Width and height must be non-negative")
    return width * height


def perimeter_of_rectangle(width: float, height: float) -> float:
    """
    Calculate perimeter of a rectangle.
    
    Args:
        width: Rectangle width
        height: Rectangle height
    
    Returns:
        Perimeter of rectangle
    """
    return 2 * (width + height)


def area_of_triangle(base: float, height: float) -> float:
    """
    Calculate area of a triangle.
    
    Args:
        base: Triangle base
        height: Triangle height
    
    Returns:
        Area of triangle
    """
    return 0.5 * base * height


def distance(x1: float, y1: float, x2: float, y2: float) -> float:
    """
    Calculate Euclidean distance between two points.
    
    Args:
        x1: First point x coordinate
        y1: First point y coordinate
        x2: Second point x coordinate
        y2: Second point y coordinate
    
    Returns:
        Distance between points
    """
    return math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2)


def slope(x1: float, y1: float, x2: float, y2: float) -> float:
    """
    Calculate slope of line through two points.
    
    Args:
        x1: First point x coordinate
        y1: First point y coordinate
        x2: Second point x coordinate
        y2: Second point y coordinate
    
    Returns:
        Slope of line
    
    Raises:
        ValueError: If line is vertical (x1 == x2)
    """
    if x1 == x2:
        raise ValueError("Vertical line has undefined slope")
    return (y2 - y1) / (x2 - x1)


# ============================================================================
# Conversion Functions
# ============================================================================

def celsius_to_fahrenheit(celsius: float) -> float:
    """
    Convert Celsius to Fahrenheit.
    
    Args:
        celsius: Temperature in Celsius
    
    Returns:
        Temperature in Fahrenheit
    """
    return (celsius * 9/5) + 32


def fahrenheit_to_celsius(fahrenheit: float) -> float:
    """
    Convert Fahrenheit to Celsius.
    
    Args:
        fahrenheit: Temperature in Fahrenheit
    
    Returns:
        Temperature in Celsius
    """
    return (fahrenheit - 32) * 5/9


def km_to_miles(km: float) -> float:
    """
    Convert kilometers to miles.
    
    Args:
        km: Distance in kilometers
    
    Returns:
        Distance in miles
    """
    return km * 0.621371


def miles_to_km(miles: float) -> float:
    """
    Convert miles to kilometers.
    
    Args:
        miles: Distance in miles
    
    Returns:
        Distance in kilometers
    """
    return miles / 0.621371


# ============================================================================
# Utility Functions
# ============================================================================

def clamp(value: Number, min_val: Number, max_val: Number) -> Number:
    """
    Clamp a value between min and max.
    
    Args:
        value: Value to clamp
        min_val: Minimum allowed value
        max_val: Maximum allowed value
    
    Returns:
        Clamped value
    """
    if value < min_val:
        return min_val
    if value > max_val:
        return max_val
    return value


def lerp(a: Number, b: Number, t: float) -> float:
    """
    Linear interpolation between a and b.
    
    Args:
        a: Start value
        b: End value
        t: Interpolation factor (0-1)
    
    Returns:
        Interpolated value
    """
    return a + (b - a) * t


def sign(x: Number) -> int:
    """
    Return sign of a number.
    
    Args:
        x: Number to check
    
    Returns:
        1 if positive, -1 if negative, 0 if zero
    """
    if x > 0:
        return 1
    if x < 0:
        return -1
    return 0


def is_positive(x: Number) -> bool:
    """
    Check if number is positive.
    
    Args:
        x: Number to check
    
    Returns:
        True if positive, False otherwise
    """
    return x > 0


def is_negative(x: Number) -> bool:
    """
    Check if number is negative.
    
    Args:
        x: Number to check
    
    Returns:
        True if negative, False otherwise
    """
    return x < 0


def is_zero(x: Number) -> bool:
    """
    Check if number is zero.
    
    Args:
        x: Number to check
    
    Returns:
        True if zero, False otherwise
    """
    return x == 0