"""
Built-in Exceptions Examples Module

Demonstrates usage of Python's built-in exceptions and common patterns.
"""

import sys
import warnings
from typing import Any, Optional


# ============================================================================
# Built-in Exception Examples
# ============================================================================

def demonstrate_value_error():
    """Demonstrate ValueError."""
    try:
        int("not a number")
    except ValueError as e:
        print(f"ValueError caught: {e}")
        raise


def demonstrate_type_error():
    """Demonstrate TypeError."""
    try:
        result = "5" + 5
    except TypeError as e:
        print(f"TypeError caught: {e}")
        raise


def demonstrate_index_error():
    """Demonstrate IndexError."""
    my_list = [1, 2, 3]
    try:
        item = my_list[10]
    except IndexError as e:
        print(f"IndexError caught: {e}")
        raise


def demonstrate_key_error():
    """Demonstrate KeyError."""
    my_dict = {"a": 1, "b": 2}
    try:
        value = my_dict["c"]
    except KeyError as e:
        print(f"KeyError caught: {e}")
        raise


def demonstrate_attribute_error():
    """Demonstrate AttributeError."""
    obj = None
    try:
        obj.some_method()
    except AttributeError as e:
        print(f"AttributeError caught: {e}")
        raise


def demonstrate_zero_division_error():
    """Demonstrate ZeroDivisionError."""
    try:
        result = 10 / 0
    except ZeroDivisionError as e:
        print(f"ZeroDivisionError caught: {e}")
        raise


def demonstrate_file_not_found_error():
    """Demonstrate FileNotFoundError."""
    try:
        with open("nonexistent_file.txt", "r") as f:
            content = f.read()
    except FileNotFoundError as e:
        print(f"FileNotFoundError caught: {e}")
        raise


def demonstrate_permission_error():
    """Demonstrate PermissionError."""
    try:
        with open("/etc/shadow", "r") as f:
            content = f.read()
    except PermissionError as e:
        print(f"PermissionError caught: {e}")
        raise


def demonstrate_import_error():
    """Demonstrate ImportError."""
    try:
        import non_existent_module
    except ImportError as e:
        print(f"ImportError caught: {e}")
        raise


# ============================================================================
# Assertion Examples
# ============================================================================

def divide_positive(a: float, b: float) -> float:
    """
    Divide two numbers with assertions.
    """
    assert a > 0, "Dividend must be positive"
    assert b > 0, "Divisor must be positive"
    assert b != 0, "Divisor cannot be zero"
    
    return a / b


def process_age(age: int) -> str:
    """
    Process age with assertions.
    """
    assert isinstance(age, int), "Age must be an integer"
    assert 0 <= age <= 150, "Age must be between 0 and 150"
    
    if age < 18:
        return "Minor"
    elif age < 65:
        return "Adult"
    else:
        return "Senior"


# ============================================================================
# Warning Examples
# ============================================================================

def deprecated_function():
    """Example of a deprecated function with warning."""
    warnings.warn(
        "deprecated_function is deprecated, use new_function instead",
        DeprecationWarning,
        stacklevel=2
    )
    return "old result"


def experimental_feature():
    """Example of experimental feature with warning."""
    warnings.warn(
        "This feature is experimental and may change",
        UserWarning,
        stacklevel=2
    )
    return "experimental result"


def resource_warning():
    """Example of resource warning."""
    import warnings
    
    # Simulate unclosed file
    f = open("test.txt", "w")
    f.write("test")
    # f.close()  # Deliberately not closing to demonstrate warning
    
    warnings.warn(
        "Unclosed file resource detected",
        ResourceWarning,
        stacklevel=2
    )


# ============================================================================
# Multiple Exception Handling
# ============================================================================

def process_data(data: Any) -> Any:
    """
    Process data with multiple exception handlers.
    """
    try:
        # Multiple operations that could raise different exceptions
        if isinstance(data, str):
            return data.upper()
        elif isinstance(data, (int, float)):
            return data * 2
        elif isinstance(data, list):
            return [x * 2 for x in data]
        else:
            raise TypeError(f"Unsupported type: {type(data)}")
    
    except TypeError as e:
        print(f"Type error: {e}")
        raise
    
    except ValueError as e:
        print(f"Value error: {e}")
        raise
    
    except Exception as e:
        print(f"Unexpected error: {e}")
        raise


def safe_execute(func, *args, **kwargs):
    """
    Safely execute a function, catching specific exceptions.
    """
    try:
        return func(*args, **kwargs)
    except (TypeError, ValueError) as e:
        print(f"Input error: {e}")
        return None
    except (IOError, OSError) as e:
        print(f"IO error: {e}")
        return None
    except Exception as e:
        print(f"Unexpected error: {e}")
        return None


# ============================================================================
# Exception Group Examples (Python 3.11+)
# ============================================================================

def demonstrate_exception_group():
    """
    Demonstrate ExceptionGroup (Python 3.11+).
    """
    try:
        # Simulate multiple errors
        errors = []
        
        try:
            1 / 0
        except ZeroDivisionError as e:
            errors.append(e)
        
        try:
            int("not a number")
        except ValueError as e:
            errors.append(e)
        
        try:
            [1, 2, 3][10]
        except IndexError as e:
            errors.append(e)
        
        if errors:
            if sys.version_info >= (3, 11):
                raise ExceptionGroup("Multiple errors occurred", errors)
            else:
                # Fallback for older Python versions
                raise Exception("Multiple errors occurred") from errors[0]
    
    except Exception as e:
        if sys.version_info >= (3, 11) and hasattr(e, "exceptions"):
            print(f"ExceptionGroup caught with {len(e.exceptions)} exceptions")
            for i, exc in enumerate(e.exceptions):
                print(f"  {i+1}: {type(exc).__name__}: {exc}")
        else:
            print(f"Exception caught: {e}")
        raise


# ============================================================================
# Custom Exception Hierarchy Example
# ============================================================================

class DataProcessingError(Exception):
    """Base error for data processing."""
    pass


class DataValidationError(DataProcessingError):
    """Error during data validation."""
    pass


class DataTransformationError(DataProcessingError):
    """Error during data transformation."""
    pass


class DataStorageError(DataProcessingError):
    """Error during data storage."""
    pass


def process_and_store_data(data: Any) -> Dict[str, Any]:
    """
    Process and store data with custom exception hierarchy.
    """
    try:
        # Validation phase
        if not data:
            raise DataValidationError("Data cannot be empty")
        
        if not isinstance(data, dict):
            raise DataValidationError(f"Expected dict, got {type(data)}")
        
        # Transformation phase
        try:
            transformed = {
                "id": str(data.get("id", "")),
                "name": data.get("name", "").upper(),
                "value": float(data.get("value", 0))
            }
        except (ValueError, TypeError) as e:
            raise DataTransformationError(f"Failed to transform data: {e}") from e
        
        # Storage phase
        try:
            # Simulate storage
            stored = {"status": "stored", "data": transformed}
        except Exception as e:
            raise DataStorageError(f"Failed to store data: {e}") from e
        
        return stored
    
    except DataProcessingError:
        # Handle all data processing errors uniformly
        raise
    except Exception as e:
        # Wrap unexpected errors
        raise DataProcessingError(f"Unexpected error: {e}") from e


# ============================================================================
# Exception Chaining Examples
# ============================================================================

def read_and_parse_file(filepath: str) -> Dict[str, Any]:
    """
    Read and parse JSON file with exception chaining.
    """
    import json
    
    try:
        with open(filepath, 'r') as f:
            return json.load(f)
    except FileNotFoundError as e:
        # Raise a more specific error with the original as cause
        raise ValueError(f"Configuration file {filepath} not found") from e
    except json.JSONDecodeError as e:
        # Chain with context
        raise ValueError(f"Invalid JSON in {filepath}: {e}") from e
    except PermissionError as e:
        raise ValueError(f"Permission denied reading {filepath}") from e


# ============================================================================
# Context Manager with Exception Handling
# ============================================================================

class ExceptionCapture:
    """
    Context manager that captures exceptions for later inspection.
    """
    
    def __init__(self, swallow: bool = False):
        self.swallow = swallow
        self.exception = None
        self.traceback = None
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        if exc_val is not None:
            self.exception = exc_val
            self.traceback = exc_tb
            
            if self.swallow:
                # Swallow the exception
                return True
        return False


# ============================================================================
# Global Exception Handler
# ============================================================================

def global_exception_handler(exc_type, exc_value, exc_traceback):
    """
    Global exception handler for uncaught exceptions.
    """
    if issubclass(exc_type, KeyboardInterrupt):
        # Don't log keyboard interrupt
        sys.__excepthook__(exc_type, exc_value, exc_traceback)
        return
    
    print(f"Uncaught exception: {exc_type.__name__}: {exc_value}")
    import traceback
    traceback.print_tb(exc_traceback)


# Install global exception handler
sys.excepthook = global_exception_handler


# ============================================================================
# Unit Test Helper Functions
# ============================================================================

def assert_raises(
    expected_exception: type,
    callable_obj: callable,
    *args,
    **kwargs
) -> Exception:
    """
    Assert that a callable raises a specific exception.
    Similar to unittest's assertRaises.
    """
    try:
        callable_obj(*args, **kwargs)
    except expected_exception as e:
        return e
    except Exception as e:
        raise AssertionError(
            f"Expected {expected_exception.__name__}, got {type(e).__name__}"
        )
    else:
        raise AssertionError(
            f"Expected {expected_exception.__name__} but no exception raised"
        )


def assert_not_raises(callable_obj: callable, *args, **kwargs) -> Any:
    """
    Assert that a callable does not raise any exception.
    """
    try:
        return callable_obj(*args, **kwargs)
    except Exception as e:
        raise AssertionError(f"Unexpected exception raised: {type(e).__name__}: {e}")


# ============================================================================
# Main demonstration
# ============================================================================

if __name__ == "__main__":
    print("Running exception demonstrations...")
    
    # Demonstrate various exception handlers
    try:
        demonstrate_value_error()
    except ValueError:
        print("Caught ValueError in main")
    
    try:
        process_age(200)
    except AssertionError as e:
        print(f"Caught assertion: {e}")
    
    # Demonstrate exception chaining
    try:
        read_and_parse_file("nonexistent.json")
    except ValueError as e:
        print(f"Chained exception: {e}")
        print(f"Cause: {e.__cause__}")
    
    # Demonstrate custom exception hierarchy
    try:
        result = process_and_store_data({"id": 123, "name": "test", "value": "45.6"})
        print(f"Processing succeeded: {result}")
    except DataProcessingError as e:
        print(f"Processing failed: {e}")