"""
Python Error Examples Package

Comprehensive examples demonstrating error handling patterns,
custom exceptions, and exception hierarchies in Python.
"""

__version__ = "1.0.0"
__author__ = "Spoke Tool Team"

from .custom_exceptions import (
    # Base exceptions
    AppError,
    BusinessError,
    TechnicalError,
    
    # HTTP status based
    BadRequestError,
    UnauthorizedError,
    ForbiddenError,
    NotFoundError,
    ConflictError,
    ValidationError,
    RateLimitError,
    
    # Domain specific - User
    UserNotFoundError,
    UserAlreadyExistsError,
    InvalidCredentialsError,
    AccountLockedError,
    
    # Domain specific - Product
    ProductNotFoundError,
    InsufficientStockError,
    ProductOutOfStockError,
    
    # Domain specific - Order
    OrderNotFoundError,
    OrderCannotBeCancelledError,
    InvalidOrderStatusError,
    
    # Domain specific - Payment
    PaymentError,
    InsufficientFundsError,
    CardDeclinedError,
    
    # Technical errors
    DatabaseError,
    ConnectionError,
    TimeoutError,
    ConfigurationError,
    ThirdPartyServiceError
)

from .error_handling import (
    divide_safe,
    process_user,
    process_order,
    retry_operation,
    validate_user_input,
    safe_file_operation,
    with_error_logging,
    ErrorHandler,
    ErrorAggregator
)

__all__ = [
    # Custom exceptions
    "AppError",
    "BusinessError",
    "TechnicalError",
    "BadRequestError",
    "UnauthorizedError",
    "ForbiddenError",
    "NotFoundError",
    "ConflictError",
    "ValidationError",
    "RateLimitError",
    "UserNotFoundError",
    "UserAlreadyExistsError",
    "InvalidCredentialsError",
    "AccountLockedError",
    "ProductNotFoundError",
    "InsufficientStockError",
    "ProductOutOfStockError",
    "OrderNotFoundError",
    "OrderCannotBeCancelledError",
    "InvalidOrderStatusError",
    "PaymentError",
    "InsufficientFundsError",
    "CardDeclinedError",
    "DatabaseError",
    "ConnectionError",
    "TimeoutError",
    "ConfigurationError",
    "ThirdPartyServiceError",
    
    # Error handling functions
    "divide_safe",
    "process_user",
    "process_order",
    "retry_operation",
    "validate_user_input",
    "safe_file_operation",
    "with_error_logging",
    "ErrorHandler",
    "ErrorAggregator",
]