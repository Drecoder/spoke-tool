"""
Python Benchmark Tests

This file contains benchmarks for various operations to measure performance
and validate the test generator's benchmark support.
"""

import pytest
import time
import json
import hashlib
import re
from collections import defaultdict, Counter
from functools import reduce
import math
import random
import string
import os
import tempfile

# ============================================================================
# String Operation Benchmarks
# ============================================================================

class TestStringBenchmarks:
    """Benchmarks for string operations."""
    
    @pytest.fixture
    def sample_string(self):
        return "The quick brown fox jumps over the lazy dog"
    
    @pytest.fixture
    def string_parts(self):
        return ["hello", "world", "this", "is", "a", "test", "string"]
    
    @pytest.fixture
    def long_string(self):
        return "The quick brown fox jumps over the lazy dog. " * 100
    
    def test_string_concatenation_plus(self, benchmark, string_parts):
        """Benchmark string concatenation with + operator."""
        def concat_with_plus():
            result = ""
            for part in string_parts:
                result += part
            return result
        
        result = benchmark(concat_with_plus)
        assert result == "helloworldthisisateststring"
    
    def test_string_concatenation_join(self, benchmark, string_parts):
        """Benchmark string concatenation with join."""
        def concat_with_join():
            return "".join(string_parts)
        
        result = benchmark(concat_with_join)
        assert result == "helloworldthisisateststring"
    
    def test_string_format_percent(self, benchmark):
        """Benchmark % formatting."""
        name = "John"
        age = 30
        
        def percent_format():
            return "Name: %s, Age: %d" % (name, age)
        
        result = benchmark(percent_format)
        assert result == "Name: John, Age: 30"
    
    def test_string_format_format(self, benchmark):
        """Benchmark str.format()."""
        name = "John"
        age = 30
        
        def format_method():
            return "Name: {}, Age: {}".format(name, age)
        
        result = benchmark(format_method)
        assert result == "Name: John, Age: 30"
    
    def test_string_format_fstring(self, benchmark):
        """Benchmark f-strings."""
        name = "John"
        age = 30
        
        def fstring():
            return f"Name: {name}, Age: {age}"
        
        result = benchmark(fstring)
        assert result == "Name: John, Age: 30"
    
    def test_string_slicing(self, benchmark, sample_string):
        """Benchmark string slicing."""
        def slice_string():
            return sample_string[10:20]
        
        result = benchmark(slice_string)
        assert result == "brown fox"
    
    def test_string_split(self, benchmark, sample_string):
        """Benchmark string split."""
        def split_string():
            return sample_string.split()
        
        result = benchmark(split_string)
        assert len(result) == 8
    
    def test_string_replace(self, benchmark, long_string):
        """Benchmark string replace."""
        def replace_string():
            return long_string.replace("fox", "cat")
        
        result = benchmark(replace_string)
        assert "cat" in result
    
    def test_string_regex_search(self, benchmark, long_string):
        """Benchmark regex search."""
        pattern = re.compile(r'\b\w{5}\b')  # 5-letter words
        
        def regex_search():
            return pattern.findall(long_string)
        
        result = benchmark(regex_search)
        assert len(result) > 0

# ============================================================================
# List Operation Benchmarks
# ============================================================================

class TestListBenchmarks:
    """Benchmarks for list operations."""
    
    @pytest.fixture(params=[10, 100, 1000])
    def list_size(self, request):
        return request.param
    
    @pytest.fixture
    def small_list(self):
        return list(range(100))
    
    @pytest.fixture
    def medium_list(self):
        return list(range(1000))
    
    @pytest.fixture
    def large_list(self):
        return list(range(10000))
    
    def test_list_comprehension(self, benchmark, list_size):
        """Benchmark list comprehension."""
        def list_comp():
            return [x * 2 for x in range(list_size)]
        
        result = benchmark(list_comp)
        assert len(result) == list_size
    
    def test_list_append(self, benchmark, list_size):
        """Benchmark list append."""
        def list_append():
            result = []
            for i in range(list_size):
                result.append(i * 2)
            return result
        
        result = benchmark(list_append)
        assert len(result) == list_size
    
    def test_list_map(self, benchmark, small_list):
        """Benchmark map function."""
        def list_map():
            return list(map(lambda x: x * 2, small_list))
        
        result = benchmark(list_map)
        assert result == [x * 2 for x in small_list]
    
    def test_list_filter(self, benchmark, medium_list):
        """Benchmark filter function."""
        def list_filter():
            return list(filter(lambda x: x % 2 == 0, medium_list))
        
        result = benchmark(list_filter)
        assert all(x % 2 == 0 for x in result)
    
    def test_list_reduce(self, benchmark, small_list):
        """Benchmark reduce function."""
        def list_reduce():
            return reduce(lambda x, y: x + y, small_list)
        
        result = benchmark(list_reduce)
        assert result == sum(small_list)
    
    def test_list_sort(self, benchmark, large_list):
        """Benchmark list sort."""
        def list_sort():
            return sorted(large_list, reverse=True)
        
        result = benchmark(list_sort)
        assert result[0] == 9999
    
    def test_list_reverse(self, benchmark, medium_list):
        """Benchmark list reversal."""
        def list_reverse():
            return medium_list[::-1]
        
        result = benchmark(list_reverse)
        assert result[0] == 999
    
    def test_list_index(self, benchmark, large_list):
        """Benchmark list indexing."""
        def list_index():
            return large_list[5000]
        
        result = benchmark(list_index)
        assert result == 5000

# ============================================================================
# Dictionary Operation Benchmarks
# ============================================================================

class TestDictBenchmarks:
    """Benchmarks for dictionary operations."""
    
    @pytest.fixture
    def small_dict(self):
        return {str(i): i for i in range(100)}
    
    @pytest.fixture
    def large_dict(self):
        return {str(i): i for i in range(10000)}
    
    def test_dict_creation_literal(self, benchmark):
        """Benchmark dict creation with literal."""
        def create_dict():
            return {"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
        
        result = benchmark(create_dict)
        assert len(result) == 5
    
    def test_dict_creation_comprehension(self, benchmark):
        """Benchmark dict comprehension."""
        def dict_comprehension():
            return {i: i * 2 for i in range(100)}
        
        result = benchmark(dict_comprehension)
        assert len(result) == 100
    
    def test_dict_creation_zip(self, benchmark):
        """Benchmark dict creation with zip."""
        keys = range(100)
        values = [x * 2 for x in range(100)]
        
        def dict_zip():
            return dict(zip(keys, values))
        
        result = benchmark(dict_zip)
        assert len(result) == 100
    
    def test_dict_get(self, benchmark, large_dict):
        """Benchmark dict get operation."""
        def dict_get():
            return large_dict.get("5000", None)
        
        result = benchmark(dict_get)
        assert result == 5000
    
    def test_dict_set(self, benchmark):
        """Benchmark dict set operation."""
        d = {}
        
        def dict_set():
            for i in range(100):
                d[str(i)] = i
        
        benchmark(dict_set)
        assert len(d) == 100
    
    def test_dict_iteration_keys(self, benchmark, large_dict):
        """Benchmark dict keys iteration."""
        def iterate_keys():
            return [k for k in large_dict.keys()]
        
        result = benchmark(iterate_keys)
        assert len(result) == 10000
    
    def test_dict_iteration_values(self, benchmark, large_dict):
        """Benchmark dict values iteration."""
        def iterate_values():
            return [v for v in large_dict.values()]
        
        result = benchmark(iterate_values)
        assert len(result) == 10000
    
    def test_dict_iteration_items(self, benchmark, large_dict):
        """Benchmark dict items iteration."""
        def iterate_items():
            return [(k, v) for k, v in large_dict.items()]
        
        result = benchmark(iterate_items)
        assert len(result) == 10000
    
    def test_dict_comparison(self, benchmark):
        """Benchmark dict comparison."""
        d1 = {str(i): i for i in range(1000)}
        d2 = {str(i): i for i in range(1000)}
        
        def dict_equal():
            return d1 == d2
        
        assert benchmark(dict_equal) is True

# ============================================================================
# Set Operation Benchmarks
# ============================================================================

class TestSetBenchmarks:
    """Benchmarks for set operations."""
    
    @pytest.fixture
    def set_a(self):
        return set(range(0, 5000, 2))  # evens
    
    @pytest.fixture
    def set_b(self):
        return set(range(0, 5000, 3))  # multiples of 3
    
    def test_set_union(self, benchmark, set_a, set_b):
        """Benchmark set union."""
        def set_union():
            return set_a | set_b
        
        result = benchmark(set_union)
        assert len(result) > len(set_a)
    
    def test_set_intersection(self, benchmark, set_a, set_b):
        """Benchmark set intersection."""
        def set_intersection():
            return set_a & set_b
        
        result = benchmark(set_intersection)
        assert len(result) < len(set_a)
    
    def test_set_difference(self, benchmark, set_a, set_b):
        """Benchmark set difference."""
        def set_difference():
            return set_a - set_b
        
        result = benchmark(set_difference)
        assert len(result) < len(set_a)
    
    def test_set_membership(self, benchmark, set_a):
        """Benchmark set membership test."""
        def set_membership():
            return 2500 in set_a
        
        assert benchmark(set_membership) is True

# ============================================================================
# Function Call Benchmarks
# ============================================================================

class TestFunctionCallBenchmarks:
    """Benchmarks for function call overhead."""
    
    def regular_function(self, x):
        return x * 2
    
    @staticmethod
    def static_method(x):
        return x * 2
    
    def lambda_function(self):
        return lambda x: x * 2
    
    def test_regular_function(self, benchmark):
        """Benchmark regular function call."""
        def call_function():
            return self.regular_function(42)
        
        result = benchmark(call_function)
        assert result == 84
    
    def test_static_method(self, benchmark):
        """Benchmark static method call."""
        def call_static():
            return self.static_method(42)
        
        result = benchmark(call_static)
        assert result == 84
    
    def test_lambda_call(self, benchmark):
        """Benchmark lambda function call."""
        f = lambda x: x * 2
        
        def call_lambda():
            return f(42)
        
        result = benchmark(call_lambda)
        assert result == 84
    
    def test_builtin_function(self, benchmark):
        """Benchmark built-in function call."""
        def call_builtin():
            return len([1, 2, 3, 4, 5])
        
        result = benchmark(call_builtin)
        assert result == 5

# ============================================================================
# Math Operation Benchmarks
# ============================================================================

class TestMathBenchmarks:
    """Benchmarks for mathematical operations."""
    
    def test_factorial_math(self, benchmark):
        """Benchmark math.factorial."""
        def math_factorial():
            return math.factorial(20)
        
        result = benchmark(math_factorial)
        assert result == 2432902008176640000
    
    def test_factorial_loop(self, benchmark):
        """Benchmark factorial with loop."""
        def loop_factorial():
            result = 1
            for i in range(1, 21):
                result *= i
            return result
        
        result = benchmark(loop_factorial)
        assert result == 2432902008176640000
    
    def test_power_operator(self, benchmark):
        """Benchmark ** operator."""
        def power_op():
            return 2 ** 20
        
        result = benchmark(power_op)
        assert result == 1048576
    
    def test_power_function(self, benchmark):
        """Benchmark pow function."""
        def power_func():
            return pow(2, 20)
        
        result = benchmark(power_func)
        assert result == 1048576
    
    def test_sqrt_math(self, benchmark):
        """Benchmark math.sqrt."""
        def sqrt_math():
            return math.sqrt(123456789)
        
        result = benchmark(sqrt_math)
        assert abs(result - 11111.111) < 0.1
    
    def test_sqrt_operator(self, benchmark):
        """Benchmark ** 0.5 operator."""
        def sqrt_op():
            return 123456789 ** 0.5
        
        result = benchmark(sqrt_op)
        assert abs(result - 11111.111) < 0.1

# ============================================================================
# Random Number Benchmarks
# ============================================================================

class TestRandomBenchmarks:
    """Benchmarks for random number generation."""
    
    def test_random_random(self, benchmark):
        """Benchmark random.random()."""
        def rand_random():
            return random.random()
        
        result = benchmark(rand_random)
        assert 0 <= result <= 1
    
    def test_random_randint(self, benchmark):
        """Benchmark random.randint()."""
        def rand_randint():
            return random.randint(1, 1000000)
        
        result = benchmark(rand_randint)
        assert 1 <= result <= 1000000
    
    def test_random_choice(self, benchmark):
        """Benchmark random.choice()."""
        sequence = list(range(1000))
        
        def rand_choice():
            return random.choice(sequence)
        
        result = benchmark(rand_choice)
        assert 0 <= result <= 999
    
    def test_random_shuffle(self, benchmark):
        """Benchmark random.shuffle()."""
        sequence = list(range(1000))
        
        def rand_shuffle():
            random.shuffle(sequence[:])  # Shuffle a copy
        
        benchmark(rand_shuffle)

# ============================================================================
# JSON Operation Benchmarks
# ============================================================================

class TestJSONBenchmarks:
    """Benchmarks for JSON operations."""
    
    @pytest.fixture
    def test_data(self):
        return {
            "id": 42,
            "name": "Test Object",
            "tags": ["python", "benchmark", "json"],
            "active": True,
            "score": 3.14159,
            "metadata": {
                "created": "2024-01-01",
                "version": 1,
                "source": "benchmark"
            },
            "nested": {
                "array": [1, 2, 3, 4, 5],
                "object": {"a": 1, "b": 2}
            }
        }
    
    def test_json_dumps(self, benchmark, test_data):
        """Benchmark json.dumps()."""
        def json_serialize():
            return json.dumps(test_data)
        
        result = benchmark(json_serialize)
        assert "Test Object" in result
    
    def test_json_loads(self, benchmark, test_data):
        """Benchmark json.loads()."""
        json_str = json.dumps(test_data)
        
        def json_deserialize():
            return json.loads(json_str)
        
        result = benchmark(json_deserialize)
        assert result["name"] == "Test Object"
    
    def test_json_roundtrip(self, benchmark, test_data):
        """Benchmark JSON round trip."""
        def json_roundtrip():
            return json.loads(json.dumps(test_data))
        
        result = benchmark(json_roundtrip)
        assert result["name"] == "Test Object"

# ============================================================================
# Hash Function Benchmarks
# ============================================================================

class TestHashBenchmarks:
    """Benchmarks for hash functions."""
    
    @pytest.fixture(params=[10, 100, 1000, 10000])
    def data_size(self, request):
        return request.param
    
    @pytest.fixture
    def test_data(self, data_size):
        return os.urandom(data_size)
    
    def test_md5(self, benchmark, test_data):
        """Benchmark MD5 hashing."""
        def md5_hash():
            return hashlib.md5(test_data).hexdigest()
        
        result = benchmark(md5_hash)
        assert len(result) == 32
    
    def test_sha1(self, benchmark, test_data):
        """Benchmark SHA1 hashing."""
        def sha1_hash():
            return hashlib.sha1(test_data).hexdigest()
        
        result = benchmark(sha1_hash)
        assert len(result) == 40
    
    def test_sha256(self, benchmark, test_data):
        """Benchmark SHA256 hashing."""
        def sha256_hash():
            return hashlib.sha256(test_data).hexdigest()
        
        result = benchmark(sha256_hash)
        assert len(result) == 64
    
    def test_blake2b(self, benchmark, test_data):
        """Benchmark BLAKE2b hashing."""
        def blake2b_hash():
            return hashlib.blake2b(test_data).hexdigest()
        
        result = benchmark(blake2b_hash)
        assert len(result) == 128

# ============================================================================
# Sorting Benchmarks
# ============================================================================

class TestSortingBenchmarks:
    """Benchmarks for sorting algorithms."""
    
    @pytest.fixture(params=[100, 1000, 10000])
    def array_size(self, request):
        return request.param
    
    @pytest.fixture
    def random_array(self, array_size):
        return [random.randint(1, 10000) for _ in range(array_size)]
    
    @pytest.fixture
    def sorted_array(self, array_size):
        return list(range(array_size))
    
    @pytest.fixture
    def reverse_array(self, array_size):
        return list(range(array_size, 0, -1))
    
    def test_sort_random(self, benchmark, random_array):
        """Benchmark sorting random array."""
        def sort_random():
            return sorted(random_array)
        
        result = benchmark(sort_random)
        assert len(result) == len(random_array)
        assert result == sorted(random_array)
    
    def test_sort_sorted(self, benchmark, sorted_array):
        """Benchmark sorting already sorted array."""
        def sort_sorted():
            return sorted(sorted_array)
        
        result = benchmark(sort_sorted)
        assert result == sorted_array
    
    def test_sort_reverse(self, benchmark, reverse_array):
        """Benchmark sorting reverse-sorted array."""
        def sort_reverse():
            return sorted(reverse_array)
        
        result = benchmark(sort_reverse)
        assert result == list(range(1, len(reverse_array) + 1))
    
    def test_sort_key(self, benchmark, random_array):
        """Benchmark sorting with key function."""
        def sort_with_key():
            return sorted(random_array, key=lambda x: -x)
        
        result = benchmark(sort_with_key)
        assert result == sorted(random_array, reverse=True)

# ============================================================================
# Collection Counter Benchmarks
# ============================================================================

class TestCounterBenchmarks:
    """Benchmarks for collections.Counter."""
    
    @pytest.fixture
    def word_list(self):
        words = []
        for _ in range(10000):
            words.append(random.choice(string.ascii_lowercase))
        return words
    
    def test_counter_manual(self, benchmark, word_list):
        """Benchmark manual counting with dict."""
        def manual_count():
            counts = {}
            for word in word_list:
                counts[word] = counts.get(word, 0) + 1
            return counts
        
        result = benchmark(manual_count)
        assert sum(result.values()) == len(word_list)
    
    def test_counter_defaultdict(self, benchmark, word_list):
        """Benchmark counting with defaultdict."""
        def defaultdict_count():
            counts = defaultdict(int)
            for word in word_list:
                counts[word] += 1
            return dict(counts)
        
        result = benchmark(defaultdict_count)
        assert sum(result.values()) == len(word_list)
    
    def test_counter_class(self, benchmark, word_list):
        """Benchmark collections.Counter."""
        def counter_count():
            return Counter(word_list)
        
        result = benchmark(counter_count)
        assert sum(result.values()) == len(word_list)

# ============================================================================
# File I/O Benchmarks
# ============================================================================

class TestFileIOBenchmarks:
    """Benchmarks for file I/O operations."""
    
    @pytest.fixture
    def temp_file(self):
        fd, path = tempfile.mkstemp()
        os.close(fd)
        yield path
        os.unlink(path)
    
    @pytest.fixture
    def test_data(self):
        return "Hello, World!\n" * 10000
    
    def test_file_write(self, benchmark, temp_file, test_data):
        """Benchmark file write."""
        def file_write():
            with open(temp_file, 'w') as f:
                f.write(test_data)
        
        benchmark(file_write)
        
        # Verify
        with open(temp_file, 'r') as f:
            content = f.read()
        assert len(content) == len(test_data)
    
    def test_file_read(self, benchmark, temp_file, test_data):
        """Benchmark file read."""
        # Write test data first
        with open(temp_file, 'w') as f:
            f.write(test_data)
        
        def file_read():
            with open(temp_file, 'r') as f:
                return f.read()
        
        result = benchmark(file_read)
        assert len(result) == len(test_data)
    
    def test_file_readlines(self, benchmark, temp_file, test_data):
        """Benchmark file readlines."""
        with open(temp_file, 'w') as f:
            f.write(test_data)
        
        def file_readlines():
            with open(temp_file, 'r') as f:
                return f.readlines()
        
        result = benchmark(file_readlines)
        assert len(result) == 10000

# ============================================================================
# Exception Handling Benchmarks
# ============================================================================

class TestExceptionBenchmarks:
    """Benchmarks for exception handling overhead."""
    
    def test_try_except_no_exception(self, benchmark):
        """Benchmark try/except with no exception."""
        def try_except():
            try:
                result = 42
            except Exception:
                result = 0
            return result
        
        result = benchmark(try_except)
        assert result == 42
    
    def test_try_except_with_exception(self, benchmark):
        """Benchmark try/except with exception."""
        def try_except_raise():
            try:
                raise ValueError("test")
            except ValueError:
                return 42
        
        result = benchmark(try_except_raise)
        assert result == 42
    
    def test_if_else(self, benchmark):
        """Benchmark if/else for comparison."""
        def if_else():
            condition = False
            if condition:
                result = 0
            else:
                result = 42
            return result
        
        result = benchmark(if_else)
        assert result == 42

# ============================================================================
# List vs Generator Benchmarks
# ============================================================================

class TestListVsGeneratorBenchmarks:
    """Benchmarks comparing lists and generators."""
    
    @pytest.fixture(params=[100, 1000, 10000])
    def size(self, request):
        return request.param
    
    def test_list_comprehension_memory(self, benchmark, size):
        """Benchmark list comprehension memory usage."""
        def list_comp():
            return [x * 2 for x in range(size)]
        
        result = benchmark(list_comp)
        assert len(result) == size
    
    def test_generator_expression(self, benchmark, size):
        """Benchmark generator expression."""
        def generator_expr():
            return sum(x * 2 for x in range(size))
        
        result = benchmark(generator_expr)
        expected = sum(x * 2 for x in range(size))
        assert result == expected
    
    def test_map_object(self, benchmark, size):
        """Benchmark map object."""
        def map_object():
            return sum(map(lambda x: x * 2, range(size)))
        
        result = benchmark(map_object)
        expected = sum(x * 2 for x in range(size))
        assert result == expected

# ============================================================================
# Decorator Benchmarks
# ============================================================================

class TestDecoratorBenchmarks:
    """Benchmarks for decorator overhead."""
    
    def simple_decorator(self, func):
        def wrapper(*args, **kwargs):
            return func(*args, **kwargs)
        return wrapper
    
    @simple_decorator
    def decorated_function(self, x):
        return x * 2
    
    def undecorated_function(self, x):
        return x * 2
    
    def test_undecorated(self, benchmark):
        """Benchmark undecorated function."""
        def call_undecorated():
            return self.undecorated_function(21)
        
        result = benchmark(call_undecorated)
        assert result == 42
    
    def test_decorated(self, benchmark):
        """Benchmark decorated function."""
        def call_decorated():
            return self.decorated_function(21)
        
        result = benchmark(call_decorated)
        assert result == 42

# ============================================================================
# Parameterized Benchmarks
# ============================================================================

@pytest.mark.parametrize("n", [10, 100, 1000, 10000])
def test_fibonacci_recursive(benchmark, n):
    """Benchmark recursive fibonacci (parameterized)."""
    def fib(n):
        if n <= 1:
            return n
        return fib(n-1) + fib(n-2)
    
    # Only run for small n to avoid excessive time
    if n > 30:
        pytest.skip("Too slow for large n")
    
    result = benchmark(lambda: fib(min(n, 30)))
    assert result >= 0

@pytest.mark.parametrize("n", [10, 100, 1000, 10000])
def test_fibonacci_iterative(benchmark, n):
    """Benchmark iterative fibonacci (parameterized)."""
    def fib_iterative(n):
        a, b = 0, 1
        for _ in range(n):
            a, b = b, a + b
        return a
    
    result = benchmark(lambda: fib_iterative(n))
    expected = fib_iterative(n)
    assert result == expected

# ============================================================================
# Main execution
# ============================================================================

if __name__ == "__main__":
    pytest.main([__file__, "-v", "--benchmark-only"])