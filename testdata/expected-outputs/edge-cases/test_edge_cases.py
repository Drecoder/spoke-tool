"""
Edge Cases Test Suite for Python

This file contains tests for various edge cases, boundary conditions,
and error scenarios to validate the test generator's edge case handling.
"""

import pytest
import sys
import os
import math
import json
import time
import asyncio
import gc
import weakref
from collections import defaultdict, deque
from dataclasses import dataclass
from contextlib import contextmanager
from unittest.mock import Mock, patch
import threading
import queue
import signal
import tempfile
import re

# ============================================================================
# Zero and Empty Values Edge Cases
# ============================================================================

class TestZeroValues:
    """Test edge cases with zero and empty values."""

    def test_zero_integer(self):
        """Test handling of zero integer."""
        result = process_integer(0)
        assert result == "zero"

    def test_zero_float(self):
        """Test handling of zero float."""
        result = process_float(0.0)
        assert result == "zero"

    def test_negative_zero(self):
        """Test handling of negative zero."""
        result = process_float(-0.0)
        assert result == "zero"

    def test_empty_string(self):
        """Test handling of empty string."""
        result = process_string("")
        assert result == "empty"

    def test_none_value(self):
        """Test handling of None."""
        with pytest.raises(ValueError, match="Value cannot be None"):
            process_value(None)

    def test_empty_list(self):
        """Test handling of empty list."""
        result = process_list([])
        assert result == "empty"

    def test_empty_tuple(self):
        """Test handling of empty tuple."""
        result = process_tuple(())
        assert result == "empty"

    def test_empty_dict(self):
        """Test handling of empty dict."""
        result = process_dict({})
        assert result == "empty"

    def test_empty_set(self):
        """Test handling of empty set."""
        result = process_set(set())
        assert result == "empty"

    def test_empty_frozenset(self):
        """Test handling of empty frozenset."""
        result = process_frozenset(frozenset())
        assert result == "empty"

    def test_empty_deque(self):
        """Test handling of empty deque."""
        result = process_deque(deque())
        assert result == "empty"

    def test_empty_bytes(self):
        """Test handling of empty bytes."""
        result = process_bytes(b'')
        assert result == "empty"

    def test_empty_bytearray(self):
        """Test handling of empty bytearray."""
        result = process_bytearray(bytearray())
        assert result == "empty"

    def test_zero_length_range(self):
        """Test handling of zero-length range."""
        result = process_range(range(0))
        assert result == "empty"


# ============================================================================
# Numeric Boundary Conditions
# ============================================================================

class TestNumericBoundaries:
    """Test edge cases with numeric boundaries."""

    def test_max_int(self):
        """Test handling of maximum integer."""
        result = process_large_number(sys.maxsize)
        assert result == "max"

    def test_min_int(self):
        """Test handling of minimum integer."""
        result = process_large_number(-sys.maxsize - 1)
        assert result == "min"

    def test_max_float(self):
        """Test handling of maximum float."""
        result = process_float(sys.float_info.max)
        assert result == "max"

    def test_min_float(self):
        """Test handling of minimum float."""
        result = process_float(sys.float_info.min)
        assert result == "min"

    def test_infinity(self):
        """Test handling of infinity."""
        with pytest.raises(ValueError, match="Cannot process infinity"):
            process_float(float('inf'))

    def test_negative_infinity(self):
        """Test handling of negative infinity."""
        with pytest.raises(ValueError, match="Cannot process infinity"):
            process_float(float('-inf'))

    def test_nan(self):
        """Test handling of NaN."""
        with pytest.raises(ValueError, match="Cannot process NaN"):
            process_float(float('nan'))

    def test_epsilon(self):
        """Test handling of epsilon."""
        result = process_float(sys.float_info.epsilon)
        assert result == "epsilon"

    def test_floating_point_precision(self):
        """Test floating point precision issues."""
        # Famous floating point precision issue
        assert 0.1 + 0.2 != 0.3
        result = process_float(0.1 + 0.2)
        assert result == "approx-0.3"

    def test_subnormal_numbers(self):
        """Test handling of subnormal numbers."""
        subnormal = sys.float_info.min * sys.float_info.epsilon
        result = process_float(subnormal)
        assert result == "subnormal"

    def test_overflow_error(self):
        """Test overflow error."""
        with pytest.raises(OverflowError):
            _ = 10.0 ** 1000000


# ============================================================================
# Integer Overflow and Underflow
# ============================================================================

class TestIntegerOverflow:
    """Test edge cases with integer overflow."""

    def test_addition_overflow(self):
        """Test integer addition overflow."""
        result = add_with_overflow(sys.maxsize, 1)
        # Python handles big ints automatically, no overflow
        assert result > sys.maxsize

    def test_multiplication_overflow(self):
        """Test integer multiplication overflow."""
        result = sys.maxsize * sys.maxsize
        assert result > sys.maxsize

    def test_bitwise_shift_overflow(self):
        """Test bitwise shift overflow."""
        result = 1 << 1000  # Python handles arbitrarily large ints
        assert result.bit_length() == 1001

    def test_recursive_addition(self):
        """Test recursive addition with large numbers."""
        def add_recursive(a, b):
            if b == 0:
                return a
            return add_recursive(a + 1, b - 1)
        
        with pytest.raises(RecursionError):
            add_recursive(0, 10000)


# ============================================================================
# String Edge Cases
# ============================================================================

class TestStringEdgeCases:
    """Test edge cases with strings."""

    def test_very_long_string(self):
        """Test handling of very long string."""
        long_string = "a" * 1_000_000
        result = process_string(long_string)
        assert result == "very long"

    def test_string_with_null(self):
        """Test handling of string with null character."""
        s = "hello\x00world"
        result = process_string(s)
        assert result == "contains null"

    def test_unicode_string(self):
        """Test handling of Unicode string."""
        s = "Hello, 世界! 🚀 🌍"
        result = process_string(s)
        assert result == "unicode"

    def test_emoji_modifiers(self):
        """Test handling of emoji modifiers."""
        s = "👨‍👩‍👧‍👦"  # Family emoji (multiple code points)
        result = process_string(s)
        assert len(s) >= 1  # Length may be >1 due to combining chars
        assert result == "has emoji"

    def test_rtl_string(self):
        """Test handling of RTL string."""
        s = "مرحبا بالعالم"  # Arabic
        result = process_string(s)
        assert result == "rtl"

    def test_control_characters(self):
        """Test handling of control characters."""
        s = "Hello\x1b[31mWorld\x1b[0m"
        result = process_string(s)
        assert result == "contains control"

    def test_surrogate_pairs(self):
        """Test handling of surrogate pairs."""
        s = "𠜎𠜱𠝹𠱓"  # Characters outside BMP
        result = process_string(s)
        assert len(s) == 4  # Each character is a surrogate pair
        assert result == "has surrogate"

    def test_zero_width_joiner(self):
        """Test handling of zero-width joiner."""
        s = "वि" + '\u200D' + "कास"
        result = process_string(s)
        assert result == "has zwj"

    def test_bidirectional_text(self):
        """Test handling of bidirectional text."""
        s = "Hello (مرحبا) World"
        result = process_string(s)
        assert result == "bidi"

    def test_invalid_utf8(self):
        """Test handling of invalid UTF-8."""
        invalid_bytes = b'\xff\xfe\xfd'
        with pytest.raises(UnicodeDecodeError):
            invalid_bytes.decode('utf-8')


# ============================================================================
# List Edge Cases
# ============================================================================

class TestListEdgeCases:
    """Test edge cases with lists."""

    def test_very_large_list(self):
        """Test handling of very large list."""
        large_list = list(range(1_000_000))
        result = process_list(large_list)
        assert result == "very large"

    def test_list_with_none(self):
        """Test handling of list with None."""
        lst = [1, None, 2, None, 3]
        result = process_list(lst)
        assert result == "has none"

    def test_nested_lists(self):
        """Test handling of deeply nested lists."""
        lst = []
        current = lst
        for _ in range(1000):
            current.append([])
            current = current[0]
        
        result = process_list(lst)
        assert result == "deeply nested"

    def test_list_with_circular_reference(self):
        """Test handling of list with circular reference."""
        lst = [1, 2, 3]
        lst.append(lst)
        
        with pytest.raises(RecursionError):
            process_list(lst)  # Should detect circular reference

    def test_list_comprehension_memory(self):
        """Test list comprehension memory usage."""
        # This might use a lot of memory
        large_list = [x for x in range(10_000_000)]
        assert len(large_list) == 10_000_000

    def test_slice_with_step(self):
        """Test list slicing with step."""
        lst = list(range(10))
        assert lst[::2] == [0, 2, 4, 6, 8]
        assert lst[::-1] == [9, 8, 7, 6, 5, 4, 3, 2, 1, 0]

    def test_out_of_range_index(self):
        """Test accessing out-of-range index."""
        lst = [1, 2, 3]
        with pytest.raises(IndexError):
            _ = lst[10]

    def test_negative_index(self):
        """Test negative indexing."""
        lst = [1, 2, 3, 4, 5]
        assert lst[-1] == 5
        assert lst[-5] == 1
        with pytest.raises(IndexError):
            _ = lst[-6]


# ============================================================================
# Dictionary Edge Cases
# ============================================================================

class TestDictEdgeCases:
    """Test edge cases with dictionaries."""

    def test_very_large_dict(self):
        """Test handling of very large dict."""
        large_dict = {i: i * 2 for i in range(100_000)}
        result = process_dict(large_dict)
        assert result == "very large"

    def test_dict_with_none_keys(self):
        """Test dict with None as key."""
        d = {None: "value"}
        result = process_dict(d)
        assert result == "has none key"

    def test_dict_with_mixed_types(self):
        """Test dict with mixed key types."""
        d = {
            1: "int",
            "key": "string",
            (1, 2): "tuple",
            None: "none",
            True: "bool"
        }
        result = process_dict(d)
        assert result == "mixed types"

    def test_dict_with_custom_objects(self):
        """Test dict with custom objects as keys."""
        class Custom:
            def __init__(self, value):
                self.value = value
            
            def __hash__(self):
                return hash(self.value)
        
        obj = Custom(42)
        d = {obj: "value"}
        result = process_dict(d)
        assert result == "has custom keys"

    def test_dict_with_missing_key(self):
        """Test accessing missing key."""
        d = {"a": 1, "b": 2}
        with pytest.raises(KeyError):
            _ = d["c"]
        
        # Safe access
        assert d.get("c") is None
        assert d.get("c", "default") == "default"

    def test_dict_comparison(self):
        """Test dict comparison."""
        d1 = {"a": 1, "b": 2}
        d2 = {"b": 2, "a": 1}
        d3 = {"a": 1, "b": 3}
        
        assert d1 == d2  # Order doesn't matter
        assert d1 != d3

    def test_dict_views(self):
        """Test dict views."""
        d = {"a": 1, "b": 2}
        keys = d.keys()
        values = d.values()
        items = d.items()
        
        assert len(keys) == 2
        assert len(values) == 2
        assert len(items) == 2
        
        # Views reflect changes
        d["c"] = 3
        assert len(keys) == 3


# ============================================================================
# Set Edge Cases
# ============================================================================

class TestSetEdgeCases:
    """Test edge cases with sets."""

    def test_very_large_set(self):
        """Test handling of very large set."""
        large_set = set(range(100_000))
        result = process_set(large_set)
        assert result == "very large"

    def test_set_with_mixed_types(self):
        """Test set with mixed types."""
        s = {1, "string", (1, 2), None, True}
        result = process_set(s)
        assert result == "mixed types"

    def test_set_operations(self):
        """Test set operations."""
        s1 = {1, 2, 3, 4, 5}
        s2 = {4, 5, 6, 7, 8}
        
        assert s1 | s2 == {1, 2, 3, 4, 5, 6, 7, 8}  # Union
        assert s1 & s2 == {4, 5}  # Intersection
        assert s1 - s2 == {1, 2, 3}  # Difference
        assert s1 ^ s2 == {1, 2, 3, 6, 7, 8}  # Symmetric difference

    def test_set_membership(self):
        """Test set membership."""
        s = set(range(1000))
        assert 500 in s
        assert 1000 not in s

    def test_frozen_set(self):
        """Test frozen set."""
        fs = frozenset([1, 2, 3])
        with pytest.raises(AttributeError):
            fs.add(4)  # Frozen set is immutable


# ============================================================================
# Function Edge Cases
# ============================================================================

class TestFunctionEdgeCases:
    """Test edge cases with functions."""

    def test_recursion_depth(self):
        """Test deep recursion."""
        def recurse(n):
            if n <= 0:
                return 0
            return 1 + recurse(n - 1)
        
        with pytest.raises(RecursionError):
            recurse(10000)

    def test_tail_recursion(self):
        """Test tail recursion (not optimized in Python)."""
        def tail_recurse(n, acc=0):
            if n <= 0:
                return acc
            return tail_recurse(n - 1, acc + 1)
        
        with pytest.raises(RecursionError):
            tail_recurse(10000)

    def test_many_arguments(self):
        """Test function with many arguments."""
        def many_args(*args):
            return len(args)
        
        result = many_args(*range(1000))
        assert result == 1000

    def test_kwargs_many(self):
        """Test function with many keyword arguments."""
        def many_kwargs(**kwargs):
            return len(kwargs)
        
        kwargs = {f"key{i}": i for i in range(1000)}
        result = many_kwargs(**kwargs)
        assert result == 1000

    def test_closure_memory(self):
        """Test closure holding large data."""
        def make_closure():
            large_data = list(range(1_000_000))
            return lambda: len(large_data)
        
        fn = make_closure()
        result = fn()
        assert result == 1_000_000

    def test_decorator_stack(self):
        """Test multiple decorators."""
        def decorator(func):
            def wrapper(*args, **kwargs):
                return func(*args, **kwargs) + 1
            return wrapper
        
        @decorator
        @decorator
        @decorator
        def func(x):
            return x
        
        assert func(0) == 3


# ============================================================================
# Generator Edge Cases
# ============================================================================

class TestGeneratorEdgeCases:
    """Test edge cases with generators."""

    def test_infinite_generator(self):
        """Test infinite generator."""
        def infinite():
            i = 0
            while True:
                yield i
                i += 1
        
        gen = infinite()
        for i, val in enumerate(gen):
            if i >= 1000:
                break
            assert val == i

    def test_generator_exhaustion(self):
        """Test generator exhaustion."""
        def count_up_to(n):
            for i in range(n):
                yield i
        
        gen = count_up_to(5)
        assert list(gen) == [0, 1, 2, 3, 4]
        assert list(gen) == []  # Generator exhausted

    def test_generator_with_send(self):
        """Test generator with send()."""
        def accumulator():
            total = 0
            while True:
                value = yield total
                if value is not None:
                    total += value
        
        gen = accumulator()
        next(gen)  # Start generator
        assert gen.send(10) == 10
        assert gen.send(5) == 15
        assert gen.send(3) == 18

    def test_generator_throw(self):
        """Test generator with throw()."""
        def gen_func():
            try:
                yield 1
                yield 2
            except ValueError:
                yield "error handled"
        
        gen = gen_func()
        assert next(gen) == 1
        assert gen.throw(ValueError) == "error handled"

    def test_generator_close(self):
        """Test generator close()."""
        closed = False
        
        def gen_func():
            nonlocal closed
            try:
                yield 1
                yield 2
            finally:
                closed = True
        
        gen = gen_func()
        assert next(gen) == 1
        gen.close()
        assert closed is True


# ============================================================================
# Async/Await Edge Cases
# ============================================================================

@pytest.mark.asyncio
class TestAsyncEdgeCases:
    """Test edge cases with async/await."""

    async def test_never_resolving(self):
        """Test coroutine that never resolves."""
        async def never():
            await asyncio.Event().wait()
        
        with pytest.raises(asyncio.TimeoutError):
            await asyncio.wait_for(never(), timeout=0.1)

    async def test_many_tasks(self):
        """Test many concurrent tasks."""
        async def task(i):
            await asyncio.sleep(0.001)
            return i
        
        tasks = [task(i) for i in range(1000)]
        results = await asyncio.gather(*tasks)
        assert len(results) == 1000
        assert results[500] == 500

    async def test_task_cancellation(self):
        """Test task cancellation."""
        async def slow_task():
            try:
                await asyncio.sleep(10)
            except asyncio.CancelledError:
                return "cancelled"
        
        task = asyncio.create_task(slow_task())
        await asyncio.sleep(0.1)
        task.cancel()
        result = await task
        assert result == "cancelled"

    async def test_gather_with_exceptions(self):
        """Test gather with exceptions."""
        async def task_ok(i):
            return i
        
        async def task_error():
            raise ValueError("error")
        
        tasks = [task_ok(1), task_error(), task_ok(3)]
        
        with pytest.raises(ValueError):
            await asyncio.gather(*tasks)

    async def test_gather_return_exceptions(self):
        """Test gather with return_exceptions=True."""
        async def task_ok(i):
            return i
        
        async def task_error():
            raise ValueError("error")
        
        tasks = [task_ok(1), task_error(), task_ok(3)]
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        assert results[0] == 1
        assert isinstance(results[1], ValueError)
        assert results[2] == 3


# ============================================================================
# Exception Edge Cases
# ============================================================================

class TestExceptionEdgeCases:
    """Test edge cases with exceptions."""

    def test_nested_exceptions(self):
        """Test nested exception handling."""
        try:
            try:
                raise ValueError("inner")
            except ValueError as inner:
                raise RuntimeError("outer") from inner
        except RuntimeError as outer:
            assert outer.__cause__ is not None
            assert outer.__cause__.args[0] == "inner"

    def test_custom_exception(self):
        """Test custom exception."""
        class CustomError(Exception):
            def __init__(self, message, code):
                super().__init__(message)
                self.code = code
                self.timestamp = time.time()
        
        with pytest.raises(CustomError) as excinfo:
            raise CustomError("test error", 500)
        
        assert excinfo.value.code == 500
        assert excinfo.value.timestamp is not None

    def test_exception_chain(self):
        """Test exception chaining."""
        try:
            try:
                raise ValueError("first")
            except ValueError:
                raise RuntimeError("second")
        except RuntimeError as e:
            assert e.__context__ is not None
            assert e.__context__.args[0] == "first"

    def test_finally_always_runs(self):
        """Test that finally always runs."""
        result = []
        
        try:
            result.append("try")
            raise ValueError("test")
        except ValueError:
            result.append("except")
        finally:
            result.append("finally")
        
        assert result == ["try", "except", "finally"]

    def test_finally_with_return(self):
        """Test finally with return."""
        def func():
            try:
                return "try"
            finally:
                return "finally"  # This overrides the try return
        
        assert func() == "finally"

    def test_exception_in_finally(self):
        """Test exception in finally block."""
        def func():
            try:
                return "try"
            finally:
                raise ValueError("finally error")
        
        with pytest.raises(ValueError):
            func()


# ============================================================================
# Context Manager Edge Cases
# ============================================================================

class TestContextManagerEdgeCases:
    """Test edge cases with context managers."""

    def test_basic_context_manager(self):
        """Test basic context manager."""
        class ManagedResource:
            def __init__(self):
                self.entered = False
                self.exited = False
            
            def __enter__(self):
                self.entered = True
                return self
            
            def __exit__(self, *args):
                self.exited = True
                return False
        
        with ManagedResource() as resource:
            assert resource.entered is True
            assert resource.exited is False
        
        assert resource.exited is True

    def test_context_manager_exception(self):
        """Test context manager with exception."""
        class Resource:
            def __enter__(self):
                return self
            
            def __exit__(self, exc_type, exc_val, exc_tb):
                return False  # Don't suppress exception
        
        with pytest.raises(ValueError):
            with Resource():
                raise ValueError("test")

    def test_context_manager_suppress(self):
        """Test context manager suppressing exception."""
        class Suppressor:
            def __enter__(self):
                return self
            
            def __exit__(self, exc_type, exc_val, exc_tb):
                return True  # Suppress exception
        
        with Suppressor():
            raise ValueError("This won't propagate")

    def test_nested_context_managers(self):
        """Test nested context managers."""
        order = []
        
        class CM:
            def __init__(self, name):
                self.name = name
            
            def __enter__(self):
                order.append(f"enter {self.name}")
                return self
            
            def __exit__(self, *args):
                order.append(f"exit {self.name}")
                return False
        
        with CM("a") as a, CM("b") as b:
            order.append("inside")
        
        assert order == ["enter a", "enter b", "inside", "exit b", "exit a"]

    @contextmanager
    def generator_cm(self):
        """Context manager implemented as generator."""
        yield "resource"
    
    def test_generator_context_manager(self):
        """Test generator-based context manager."""
        with self.generator_cm() as resource:
            assert resource == "resource"


# ============================================================================
# Decorator Edge Cases
# ============================================================================

class TestDecoratorEdgeCases:
    """Test edge cases with decorators."""

    def test_simple_decorator(self):
        """Test simple decorator."""
        def double(func):
            def wrapper(*args, **kwargs):
                return func(*args, **kwargs) * 2
            return wrapper
        
        @double
        def add(a, b):
            return a + b
        
        assert add(2, 3) == 10

    def test_decorator_with_args(self):
        """Test decorator with arguments."""
        def multiply(multiplier):
            def decorator(func):
                def wrapper(*args, **kwargs):
                    return func(*args, **kwargs) * multiplier
                return wrapper
            return decorator
        
        @multiply(3)
        def add(a, b):
            return a + b
        
        assert add(2, 3) == 15

    def test_multiple_decorators(self):
        """Test multiple decorators."""
        def add_one(func):
            def wrapper(*args, **kwargs):
                return func(*args, **kwargs) + 1
            return wrapper
        
        def double(func):
            def wrapper(*args, **kwargs):
                return func(*args, **kwargs) * 2
            return wrapper
        
        @double
        @add_one
        def add(a, b):
            return a + b
        
        assert add(2, 3) == 12  # (2+3+1)*2

    def test_class_decorator(self):
        """Test class decorator."""
        def add_method(cls):
            cls.new_method = lambda self: "added"
            return cls
        
        @add_method
        class MyClass:
            pass
        
        obj = MyClass()
        assert obj.new_method() == "added"


# ============================================================================
# Metaclass Edge Cases
# ============================================================================

class TestMetaclassEdgeCases:
    """Test edge cases with metaclasses."""

    def test_simple_metaclass(self):
        """Test simple metaclass."""
        class Meta(type):
            def __new__(cls, name, bases, dct):
                dct['added'] = 'by metaclass'
                return super().__new__(cls, name, bases, dct)
        
        class MyClass(metaclass=Meta):
            pass
        
        assert MyClass.added == 'by metaclass'

    def test_metaclass_with_methods(self):
        """Test metaclass with methods."""
        class Meta(type):
            def __call__(cls, *args, **kwargs):
                instance = super().__call__(*args, **kwargs)
                instance.created = True
                return instance
        
        class MyClass(metaclass=Meta):
            pass
        
        obj = MyClass()
        assert obj.created is True

    def test_metaclass_inheritance(self):
        """Test metaclass inheritance."""
        class MetaA(type):
            def __new__(cls, name, bases, dct):
                dct['from_a'] = True
                return super().__new__(cls, name, bases, dct)
        
        class MetaB(MetaA):
            def __new__(cls, name, bases, dct):
                dct['from_b'] = True
                return super().__new__(cls, name, bases, dct)
        
        class MyClass(metaclass=MetaB):
            pass
        
        assert MyClass.from_a is True
        assert MyClass.from_b is True


# ============================================================================
# Descriptor Edge Cases
# ============================================================================

class TestDescriptorEdgeCases:
    """Test edge cases with descriptors."""

    def test_data_descriptor(self):
        """Test data descriptor."""
        class ValidatedAttribute:
            def __init__(self, validator):
                self.validator = validator
                self.data = {}
            
            def __get__(self, obj, objtype=None):
                if obj is None:
                    return self
                return self.data.get(id(obj), None)
            
            def __set__(self, obj, value):
                if not self.validator(value):
                    raise ValueError(f"Invalid value: {value}")
                self.data[id(obj)] = value
        
        class Person:
            age = ValidatedAttribute(lambda x: 0 <= x <= 150)
            
            def __init__(self, age):
                self.age = age
        
        p = Person(30)
        assert p.age == 30
        
        with pytest.raises(ValueError):
            Person(200)

    def test_non_data_descriptor(self):
        """Test non-data descriptor."""
        class cached_property:
            def __init__(self, func):
                self.func = func
                self.cache = {}
            
            def __get__(self, obj, objtype=None):
                if obj is None:
                    return self
                if id(obj) not in self.cache:
                    self.cache[id(obj)] = self.func(obj)
                return self.cache[id(obj)]
        
        class Data:
            @cached_property
            def expensive(self):
                return sum(range(1000000))
        
        d = Data()
        result1 = d.expensive
        result2 = d.expensive
        assert result1 == result2


# ============================================================================
# Helper Functions
# ============================================================================

def process_integer(n):
    """Process integer value."""
    if n == 0:
        return "zero"
    if n == sys.maxsize:
        return "max"
    if n == -sys.maxsize - 1:
        return "min"
    return n

def process_float(f):
    """Process float value."""
    if f == 0.0:
        return "zero"
    if math.isinf(f):
        raise ValueError("Cannot process infinity")
    if math.isnan(f):
        raise ValueError("Cannot process NaN")
    if f == sys.float_info.max:
        return "max"
    if f == sys.float_info.min:
        return "min"
    if f == sys.float_info.epsilon:
        return "epsilon"
    if f < sys.float_info.min:
        return "subnormal"
    if abs(f - 0.3) < 0.000001:
        return "approx-0.3"
    return f

def process_string(s):
    """Process string value."""
    if s == "":
        return "empty"
    if len(s) > 1000000:
        return "very long"
    if '\x00' in s:
        return "contains null"
    if any(ord(c) > 0xFFFF for c in s):
        return "has surrogate"
    if any('\u200D' in s):
        return "has zwj"
    if any('\u0590' <= c <= '\u05FF' for c in s):
        return "rtl"
    if any('\u1F600' <= c <= '\u1F6FF' for c in s):
        return "has emoji"
    if '\x1b' in s:
        return "contains control"
    if any(0x0600 <= ord(c) <= 0x06FF for c in s):
        return "bidi"
    if any(ord(c) > 0x7F for c in s):
        return "unicode"
    return "normal"

def process_value(v):
    """Process any value."""
    if v is None:
        raise ValueError("Value cannot be None")
    return v

def process_list(lst):
    """Process list value."""
    if not lst:
        return "empty"
    if len(lst) > 1000000:
        return "very large"
    if None in lst:
        return "has none"
    
    # Check for circular reference
    try:
        repr(lst)
        return "normal"
    except RecursionError:
        return "circular"

def process_tuple(t):
    """Process tuple value."""
    return "empty" if not t else "has items"

def process_dict(d):
    """Process dict value."""
    if not d:
        return "empty"
    if len(d) > 100000:
        return "very large"
    if None in d:
        return "has none key"
    
    # Check for mixed types
    types = {type(k) for k in d.keys()}
    if len(types) > 1:
        return "mixed types"
    
    # Check for custom objects
    if any(not isinstance(k, (int, str, tuple, type(None), bool)) for k in d):
        return "has custom keys"
    
    return "normal"

def process_set(s):
    """Process set value."""
    if not s:
        return "empty"
    if len(s) > 100000:
        return "very large"
    
    types = {type(e) for e in s}
    if len(types) > 1:
        return "mixed types"
    
    return "has items"

def process_frozenset(fs):
    """Process frozenset value."""
    return "empty" if not fs else "has items"

def process_deque(d):
    """Process deque value."""
    return "empty" if not d else "has items"

def process_bytes(b):
    """Process bytes value."""
    return "empty" if not b else "has data"

def process_bytearray(ba):
    """Process bytearray value."""
    return "empty" if not ba else "has data"

def process_range(r):
    """Process range value."""
    return "empty" if len(r) == 0 else "has items"

def process_large_number(n):
    """Process large number."""
    if n == sys.maxsize:
        return "max"
    if n == -sys.maxsize - 1:
        return "min"
    return "normal"

def add_with_overflow(a, b):
    """Add two numbers (Python handles big ints automatically)."""
    return a + b