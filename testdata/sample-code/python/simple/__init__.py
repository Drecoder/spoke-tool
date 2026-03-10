"""
Python Simple Examples Package

Basic Python examples demonstrating fundamental concepts
for testing code generation tools.
"""

__version__ = "1.0.0"
__author__ = "Spoke Tool Team"

from .hello import greet, main
from .math import (
    add, subtract, multiply, divide,
    power, sqrt, factorial, fibonacci,
    is_even, is_odd, min_of, max_of,
    sum_of, average_of
)
from .strings import (
    reverse, to_upper, to_lower, capitalize,
    is_palindrome, count_vowels, count_consonants,
    remove_whitespace, truncate
)

__all__ = [
    # Hello module
    "greet",
    "main",
    
    # Math module
    "add",
    "subtract",
    "multiply",
    "divide",
    "power",
    "sqrt",
    "factorial",
    "fibonacci",
    "is_even",
    "is_odd",
    "min_of",
    "max_of",
    "sum_of",
    "average_of",
    
    # Strings module
    "reverse",
    "to_upper",
    "to_lower",
    "capitalize",
    "is_palindrome",
    "count_vowels",
    "count_consonants",
    "remove_whitespace",
    "truncate",
]