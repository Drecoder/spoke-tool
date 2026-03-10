"""
Error Handling Patterns Module

Demonstrates various error handling patterns and best practices in Python.
"""

import functools
import logging
import time
from contextlib import contextmanager
from typing import Any, Callable, Dict, List, Optional, Type, Union

from .custom_exceptions import (
    AppError,
    ValidationError,
    NotFoundError,
    DatabaseError,
    ConnectionError,
    TimeoutError,
    PaymentError,
    InsufficientStockError
)

logger = logging.getLogger(__name__)


# ============================================================================
# Basic Try/Except Patterns
# ============================================================================

def divide_safe(a: float, b: float) -> float:
    """
    Safe division with error handling.
    Demonstrates basic try/except with specific exception types.
    """
    try:
        result = a / b
        return result
    except ZeroDivisionError as e:
        logger.error(f"Division by zero: {e}")
        raise ValueError("Cannot divide by zero") from e
    except TypeError as e:
        logger.error(f"Invalid type for division: {e}")
        raise TypeError("Both arguments must be numbers") from e
    except Exception as e:
        logger.error(f"Unexpected error during division: {e}")
        raise


def read_config_safe(file_path: str) -> Dict[str, Any]:
    """
    Safely read configuration file with multiple error handlers.
    Demonstrates handling different exception types.
    """
    import json
    
    try:
        with open(file_path, 'r') as f:
            config = json.load(f)
            return config
    except FileNotFoundError as e:
        logger.warning(f"Config file not found: {e}")
        return {}  # Return default config
    except json.JSONDecodeError as e:
        logger.error(f"Invalid JSON in config file: {e}")
        raise ValueError("Config file contains invalid JSON") from e
    except PermissionError as e:
        logger.error(f"Permission denied reading config: {e}")
        raise PermissionError("Cannot read config file - permission denied") from e
    except Exception as e:
        logger.error(f"Unexpected error reading config: {e}")
        raise


# ============================================================================
# Try/Except/Else/Finally Patterns
# ============================================================================

def process_user(user_id: int) -> Dict[str, Any]:
    """
    Process user data with comprehensive error handling.
    Demonstrates try/except/else/finally pattern.
    """
    db_connection = None
    try:
        # Attempt to get database connection
        db_connection = get_database_connection()
        
        # Fetch user
        user = fetch_user_from_db(db_connection, user_id)
        
        if not user:
            raise NotFoundError("User", user_id)
        
        # Process user data
        processed_data = {
            "id": user["id"],
            "name": user["name"].upper(),
            "email": user["email"].lower(),
            "processed_at": time.time()
        }
        
    except NotFoundError as e:
        logger.info(f"User not found: {user_id}")
        raise  # Re-raise the exception
    
    except DatabaseError as e:
        logger.error(f"Database error while fetching user {user_id}: {e}")
        raise  # Re-raise
    
    except Exception as e:
        logger.error(f"Unexpected error processing user {user_id}: {e}")
        raise RuntimeError(f"Failed to process user {user_id}") from e
    
    else:
        # This block runs only if no exception occurred
        logger.info(f"Successfully processed user {user_id}")
        return processed_data
    
    finally:
        # This block always runs, regardless of exceptions
        if db_connection:
            close_database_connection(db_connection)
            logger.debug("Database connection closed")


# ============================================================================
# Context Managers for Resource Management
# ============================================================================

@contextmanager
def database_transaction():
    """
    Context manager for database transactions.
    Automatically commits on success, rolls back on error.
    """
    connection = get_database_connection()
    try:
        logger.debug("Starting database transaction")
        yield connection
        connection.commit()
        logger.debug("Transaction committed")
    except Exception as e:
        connection.rollback()
        logger.error(f"Transaction rolled back due to error: {e}")
        raise
    finally:
        close_database_connection(connection)


@contextmanager
def managed_file(file_path: str, mode: str = 'r'):
    """
    Context manager for file operations with automatic cleanup.
    """
    file_obj = None
    try:
        file_obj = open(file_path, mode)
        logger.debug(f"Opened file: {file_path}")
        yield file_obj
    except IOError as e:
        logger.error(f"File operation failed: {e}")
        raise
    finally:
        if file_obj:
            file_obj.close()
            logger.debug(f"Closed file: {file_path}")


def safe_file_operation(file_path: str, operation: Callable) -> Any:
    """
    Perform safe file operation using context manager.
    """
    with managed_file(file_path, 'r') as f:
        return operation(f)


# ============================================================================
# Decorators for Error Handling
# ============================================================================

def retry(
    max_attempts: int = 3,
    delay: float = 1.0,
    backoff: float = 2.0,
    exceptions: tuple = (Exception,)
):
    """
    Decorator to retry a function on failure with exponential backoff.
    """
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            current_delay = delay
            last_exception = None
            
            for attempt in range(max_attempts):
                try:
                    return func(*args, **kwargs)
                except exceptions as e:
                    last_exception = e
                    if attempt < max_attempts - 1:
                        logger.warning(
                            f"Attempt {attempt + 1}/{max_attempts} failed: {e}. "
                            f"Retrying in {current_delay}s"
                        )
                        time.sleep(current_delay)
                        current_delay *= backoff
            
            # All attempts failed
            logger.error(f"All {max_attempts} attempts failed for {func.__name__}")
            raise last_exception
        
        return wrapper
    return decorator


def with_error_logging(func: Callable) -> Callable:
    """
    Decorator to log function entry, exit, and errors.
    """
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        func_name = func.__name__
        logger.debug(f"Entering {func_name}")
        
        try:
            result = func(*args, **kwargs)
            logger.debug(f"Exiting {func_name} successfully")
            return result
        except Exception as e:
            logger.error(f"Error in {func_name}: {e}", exc_info=True)
            raise
    
    return wrapper


def handle_errors(
    error_map: Dict[Type[Exception], Union[Type[Exception], Callable]]
):
    """
    Decorator to map specific exceptions to custom handlers.
    """
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            try:
                return func(*args, **kwargs)
            except tuple(error_map.keys()) as e:
                handler = error_map[type(e)]
                
                if callable(handler):
                    return handler(e)
                elif issubclass(handler, Exception):
                    raise handler(str(e)) from e
                else:
                    raise
        return wrapper
    return decorator


# ============================================================================
# Error Handling with Callbacks
# ============================================================================

def with_error_callback(
    func: Callable,
    on_success: Optional[Callable] = None,
    on_error: Optional[Callable] = None,
    on_finally: Optional[Callable] = None
) -> Any:
    """
    Execute function with success/error/finally callbacks.
    """
    try:
        result = func()
        if on_success:
            on_success(result)
        return result
    except Exception as e:
        if on_error:
            on_error(e)
        raise
    finally:
        if on_finally:
            on_finally()


# ============================================================================
# Graceful Degradation
# ============================================================================

def get_user_data_with_fallback(
    user_id: int,
    fallback_data: Optional[Dict] = None
) -> Dict:
    """
    Get user data with graceful degradation.
    Returns fallback data if primary source fails.
    """
    try:
        # Try primary data source (database)
        user = fetch_user_from_db(user_id)
        if user:
            return user
        
        # If not found in DB, try cache
        user = fetch_user_from_cache(user_id)
        if user:
            logger.info(f"User {user_id} retrieved from cache")
            return user
        
        # If all else fails, use fallback
        logger.warning(f"User {user_id} not found, using fallback data")
        return fallback_data or {"id": user_id, "name": "Unknown User"}
    
    except DatabaseError as e:
        logger.error(f"Database error fetching user {user_id}: {e}")
        # Try cache as fallback
        try:
            user = fetch_user_from_cache(user_id)
            if user:
                return user
        except Exception as cache_error:
            logger.error(f"Cache error as well: {cache_error}")
        
        # Ultimate fallback
        return fallback_data or {"id": user_id, "name": "Unknown User"}
    
    except Exception as e:
        logger.error(f"Unexpected error fetching user {user_id}: {e}")
        return fallback_data or {"id": user_id, "name": "Unknown User"}


# ============================================================================
# Circuit Breaker Pattern
# ============================================================================

class CircuitBreaker:
    """
    Circuit breaker pattern implementation.
    Prevents repeated calls to failing services.
    """
    
    CLOSED = "CLOSED"  # Normal operation
    OPEN = "OPEN"      # Failing, don't attempt calls
    HALF_OPEN = "HALF_OPEN"  # Testing if service recovered
    
    def __init__(
        self,
        failure_threshold: int = 5,
        recovery_timeout: int = 60,
        name: str = "default"
    ):
        self.failure_threshold = failure_threshold
        self.recovery_timeout = recovery_timeout
        self.name = name
        
        self.state = self.CLOSED
        self.failure_count = 0
        self.last_failure_time = None
        self.success_count = 0
    
    def __call__(self, func: Callable) -> Callable:
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            if self.state == self.OPEN:
                if self._should_attempt_recovery():
                    self.state = self.HALF_OPEN
                    logger.info(f"Circuit breaker {self.name} half-open, testing recovery")
                else:
                    raise Exception(f"Circuit breaker {self.name} is OPEN")
            
            try:
                result = func(*args, **kwargs)
                self._on_success()
                return result
            except Exception as e:
                self._on_failure()
                raise
        
        return wrapper
    
    def _should_attempt_recovery(self) -> bool:
        if not self.last_failure_time:
            return True
        
        elapsed = time.time() - self.last_failure_time
        return elapsed >= self.recovery_timeout
    
    def _on_success(self):
        if self.state == self.HALF_OPEN:
            self.success_count += 1
            if self.success_count >= 2:  # Two successes in half-open state
                self.state = self.CLOSED
                self.failure_count = 0
                self.success_count = 0
                logger.info(f"Circuit breaker {self.name} closed - service recovered")
    
    def _on_failure(self):
        self.failure_count += 1
        self.last_failure_time = time.time()
        
        if self.state == self.HALF_OPEN:
            self.state = self.OPEN
            logger.warning(f"Circuit breaker {self.name} open - recovery attempt failed")
        elif self.failure_count >= self.failure_threshold:
            self.state = self.OPEN
            logger.warning(f"Circuit breaker {self.name} open - threshold reached")


# ============================================================================
# Error Aggregator
# ============================================================================

class ErrorAggregator:
    """
    Collects and aggregates multiple errors.
    Useful for batch operations where multiple failures can occur.
    """
    
    def __init__(self, raise_on_first: bool = False):
        self.errors: List[Dict[str, Any]] = []
        self.warnings: List[Dict[str, Any]] = []
        self.raise_on_first = raise_on_first
    
    def execute(self, operation: Callable, *args, **kwargs) -> Optional[Any]:
        """
        Execute operation and collect any errors.
        """
        try:
            return operation(*args, **kwargs)
        except Exception as e:
            if self.raise_on_first:
                raise
            
            self.add_error(e, operation.__name__, args, kwargs)
            return None
    
    def add_error(
        self,
        error: Exception,
        operation: str = None,
        args: tuple = None,
        kwargs: dict = None
    ):
        """Add an error to the aggregator."""
        self.errors.append({
            "error": str(error),
            "type": error.__class__.__name__,
            "operation": operation,
            "args": args,
            "kwargs": kwargs,
            "timestamp": time.time()
        })
    
    def add_warning(self, message: str, context: dict = None):
        """Add a warning to the aggregator."""
        self.warnings.append({
            "message": message,
            "context": context or {},
            "timestamp": time.time()
        })
    
    def has_errors(self) -> bool:
        """Check if any errors were collected."""
        return len(self.errors) > 0
    
    def has_warnings(self) -> bool:
        """Check if any warnings were collected."""
        return len(self.warnings) > 0
    
    def get_summary(self) -> Dict[str, Any]:
        """Get summary of collected errors and warnings."""
        return {
            "error_count": len(self.errors),
            "warning_count": len(self.warnings),
            "errors": self.errors,
            "warnings": self.warnings
        }
    
    def clear(self):
        """Clear all collected errors and warnings."""
        self.errors.clear()
        self.warnings.clear()


# ============================================================================
# Exception Chaining and Context
# ============================================================================

def process_order(order_id: int, user_id: int) -> Dict[str, Any]:
    """
    Process order with proper exception chaining.
    """
    try:
        # Validate user
        user = fetch_user_from_db(user_id)
        if not user:
            raise NotFoundError("User", user_id)
        
        # Validate order
        order = fetch_order_from_db(order_id)
        if not order:
            raise NotFoundError("Order", order_id)
        
        # Process payment
        try:
            payment_result = process_payment(order["total"], user["payment_method"])
        except PaymentError as e:
            # Chain with context
            raise PaymentError(
                f"Payment failed for order {order_id}",
                details={"order_id": order_id, "user_id": user_id}
            ) from e
        
        return {
            "order": order,
            "user": user,
            "payment": payment_result,
            "status": "completed"
        }
        
    except (NotFoundError, PaymentError):
        # Re-raise domain exceptions
        raise
    except DatabaseError as e:
        # Convert database error to domain error with context
        raise DatabaseError(
            f"Failed to process order {order_id}",
            operation="process_order",
            original_error=e
        )
    except Exception as e:
        # Catch-all with context
        raise RuntimeError(f"Unexpected error processing order {order_id}") from e


# ============================================================================
# Validation with Error Collection
# ============================================================================

def validate_user_input(data: Dict[str, Any]) -> Dict[str, Any]:
    """
    Validate user input with comprehensive error collection.
    """
    errors = []
    
    # Validate email
    email = data.get("email")
    if not email:
        errors.append({"field": "email", "message": "Email is required"})
    elif "@" not in email:
        errors.append({"field": "email", "message": "Invalid email format"})
    
    # Validate age
    age = data.get("age")
    if age is not None:
        try:
            age = int(age)
            if age < 0 or age > 150:
                errors.append({"field": "age", "message": "Age must be between 0 and 150"})
        except (ValueError, TypeError):
            errors.append({"field": "age", "message": "Age must be a number"})
    
    # Validate name
    name = data.get("name")
    if not name:
        errors.append({"field": "name", "message": "Name is required"})
    elif len(name) < 2:
        errors.append({"field": "name", "message": "Name must be at least 2 characters"})
    
    # Validate phone (optional)
    phone = data.get("phone")
    if phone:
        import re
        if not re.match(r"^\+?[\d\s-]{10,}$", phone):
            errors.append({"field": "phone", "message": "Invalid phone format"})
    
    if errors:
        raise ValidationError(
            message="User input validation failed",
            errors=errors
        )
    
    return data


# ============================================================================
# Mock functions for demonstration
# ============================================================================

def get_database_connection():
    """Mock database connection."""
    return {"connected": True}


def close_database_connection(conn):
    """Mock close connection."""
    pass


def fetch_user_from_db(user_id: int) -> Optional[Dict]:
    """Mock fetching user from database."""
    if user_id == 1:
        return {"id": 1, "name": "Alice", "email": "alice@example.com"}
    return None


def fetch_user_from_cache(user_id: int) -> Optional[Dict]:
    """Mock fetching user from cache."""
    if user_id == 2:
        return {"id": 2, "name": "Bob", "email": "bob@example.com"}
    return None


def fetch_order_from_db(order_id: int) -> Optional[Dict]:
    """Mock fetching order from database."""
    if order_id == 100:
        return {"id": 100, "total": 150.00, "status": "pending"}
    return None


def process_payment(amount: float, method: str) -> Dict:
    """Mock payment processing."""
    if amount > 1000:
        raise PaymentError("Amount exceeds limit")
    return {"status": "success", "transaction_id": "txn_123"}