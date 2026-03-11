"""
Custom Exception Classes Module

Defines a comprehensive hierarchy of custom exceptions for the application.
"""

import logging
from typing import Optional, Any, Dict, List

logger = logging.getLogger(__name__)


# ============================================================================
# Base Exception Classes
# ============================================================================

class AppError(Exception):
    """
    Base application exception class.
    All custom exceptions should inherit from this.
    """
    
    def __init__(
        self,
        message: str = "An application error occurred",
        code: str = "APP_ERROR",
        status_code: int = 500,
        details: Optional[Dict[str, Any]] = None,
        cause: Optional[Exception] = None
    ):
        self.message = message
        self.code = code
        self.status_code = status_code
        self.details = details or {}
        self.cause = cause
        self.timestamp = __import__("datetime").datetime.utcnow()
        
        # Build error message
        full_message = f"[{self.code}] {self.message}"
        if self.details:
            full_message += f" - Details: {self.details}"
        if self.cause:
            full_message += f" - Caused by: {self.cause}"
        
        super().__init__(full_message)
        
        # Log error creation
        logger.debug(f"Created exception: {self.__class__.__name__} - {self.code}")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert exception to dictionary for API responses."""
        return {
            "error": {
                "code": self.code,
                "message": self.message,
                "status_code": self.status_code,
                "details": self.details,
                "timestamp": self.timestamp.isoformat()
            }
        }
    
    def __str__(self) -> str:
        return f"{self.__class__.__name__}: [{self.code}] {self.message}"


class BusinessError(AppError):
    """
    Base class for business logic errors.
    These are expected errors that occur during normal operation.
    """
    
    def __init__(
        self,
        message: str = "Business rule violation",
        code: str = "BUSINESS_ERROR",
        status_code: int = 422,
        **kwargs
    ):
        super().__init__(message, code, status_code, **kwargs)


class TechnicalError(AppError):
    """
    Base class for technical/system errors.
    These are unexpected errors that indicate system problems.
    """
    
    def __init__(
        self,
        message: str = "Technical error occurred",
        code: str = "TECHNICAL_ERROR",
        status_code: int = 500,
        **kwargs
    ):
        super().__init__(message, code, status_code, **kwargs)


# ============================================================================
# HTTP Status Based Exceptions
# ============================================================================

class BadRequestError(BusinessError):
    """400 Bad Request error."""
    
    def __init__(
        self,
        message: str = "Bad request",
        code: str = "BAD_REQUEST",
        details: Optional[Dict[str, Any]] = None,
        **kwargs
    ):
        super().__init__(message, code, 400, details=details, **kwargs)


class UnauthorizedError(BusinessError):
    """401 Unauthorized error."""
    
    def __init__(
        self,
        message: str = "Authentication required",
        code: str = "UNAUTHORIZED",
        **kwargs
    ):
        super().__init__(message, code, 401, **kwargs)


class ForbiddenError(BusinessError):
    """403 Forbidden error."""
    
    def __init__(
        self,
        message: str = "Insufficient permissions",
        code: str = "FORBIDDEN",
        **kwargs
    ):
        super().__init__(message, code, 403, **kwargs)


class NotFoundError(BusinessError):
    """404 Not Found error."""
    
    def __init__(
        self,
        resource: str = "Resource",
        identifier: Any = None,
        message: Optional[str] = None,
        code: str = "NOT_FOUND",
        **kwargs
    ):
        if message is None:
            message = f"{resource} not found"
            if identifier:
                message += f": {identifier}"
        
        details = kwargs.get("details", {})
        details.update({
            "resource": resource,
            "identifier": str(identifier) if identifier else None
        })
        
        super().__init__(message, code, 404, details=details, **kwargs)


class ConflictError(BusinessError):
    """409 Conflict error."""
    
    def __init__(
        self,
        message: str = "Resource conflict",
        code: str = "CONFLICT",
        **kwargs
    ):
        super().__init__(message, code, 409, **kwargs)


class ValidationError(BadRequestError):
    """422 Validation error."""
    
    def __init__(
        self,
        message: str = "Validation failed",
        errors: Optional[List[Dict[str, Any]]] = None,
        code: str = "VALIDATION_ERROR",
        **kwargs
    ):
        details = kwargs.get("details", {})
        details["validation_errors"] = errors or []
        super().__init__(message, code, details=details, **kwargs)
    
    def add_error(self, field: str, message: str, value: Any = None):
        """Add a validation error."""
        if "validation_errors" not in self.details:
            self.details["validation_errors"] = []
        
        self.details["validation_errors"].append({
            "field": field,
            "message": message,
            "value": str(value) if value else None
        })


class RateLimitError(BusinessError):
    """429 Too Many Requests error."""
    
    def __init__(
        self,
        message: str = "Rate limit exceeded",
        retry_after: int = 60,
        code: str = "RATE_LIMIT_EXCEEDED",
        **kwargs
    ):
        details = kwargs.get("details", {})
        details["retry_after"] = retry_after
        
        super().__init__(message, code, 429, details=details, **kwargs)


# ============================================================================
# Domain-Specific Exceptions - User
# ============================================================================

class UserNotFoundError(NotFoundError):
    """User not found error."""
    
    def __init__(
        self,
        identifier: Any = None,
        message: Optional[str] = None,
        **kwargs
    ):
        super().__init__("User", identifier, message, code="USER_NOT_FOUND", **kwargs)


class UserAlreadyExistsError(ConflictError):
    """User already exists error."""
    
    def __init__(
        self,
        email: str = None,
        username: str = None,
        **kwargs
    ):
        if email:
            message = f"User with email {email} already exists"
            details = {"email": email}
        elif username:
            message = f"User with username {username} already exists"
            details = {"username": username}
        else:
            message = "User already exists"
            details = {}
        
        super().__init__(message, "USER_ALREADY_EXISTS", details=details, **kwargs)


class InvalidCredentialsError(UnauthorizedError):
    """Invalid login credentials error."""
    
    def __init__(self, **kwargs):
        super().__init__(
            message="Invalid email or password",
            code="INVALID_CREDENTIALS",
            **kwargs
        )


class AccountLockedError(ForbiddenError):
    """Account locked error."""
    
    def __init__(
        self,
        user_id: int = None,
        locked_until: str = None,
        **kwargs
    ):
        message = "Account is locked"
        if locked_until:
            message += f" until {locked_until}"
        
        details = {"user_id": user_id, "locked_until": locked_until}
        
        super().__init__(
            message=message,
            code="ACCOUNT_LOCKED",
            details=details,
            **kwargs
        )


# ============================================================================
# Domain-Specific Exceptions - Product
# ============================================================================

class ProductNotFoundError(NotFoundError):
    """Product not found error."""
    
    def __init__(
        self,
        identifier: Any = None,
        message: Optional[str] = None,
        **kwargs
    ):
        super().__init__("Product", identifier, message, code="PRODUCT_NOT_FOUND", **kwargs)


class InsufficientStockError(BadRequestError):
    """Insufficient stock error."""
    
    def __init__(
        self,
        product_id: int,
        requested: int,
        available: int,
        **kwargs
    ):
        message = f"Insufficient stock for product {product_id}. Requested: {requested}, Available: {available}"
        details = {
            "product_id": product_id,
            "requested": requested,
            "available": available
        }
        
        super().__init__(
            message=message,
            code="INSUFFICIENT_STOCK",
            details=details,
            **kwargs
        )


class ProductOutOfStockError(BadRequestError):
    """Product out of stock error."""
    
    def __init__(
        self,
        product_id: int,
        **kwargs
    ):
        message = f"Product {product_id} is out of stock"
        details = {"product_id": product_id}
        
        super().__init__(
            message=message,
            code="OUT_OF_STOCK",
            details=details,
            **kwargs
        )


# ============================================================================
# Domain-Specific Exceptions - Order
# ============================================================================

class OrderNotFoundError(NotFoundError):
    """Order not found error."""
    
    def __init__(
        self,
        identifier: Any = None,
        message: Optional[str] = None,
        **kwargs
    ):
        super().__init__("Order", identifier, message, code="ORDER_NOT_FOUND", **kwargs)


class OrderCannotBeCancelledError(BadRequestError):
    """Order cannot be cancelled error."""
    
    def __init__(
        self,
        order_id: int,
        status: str,
        **kwargs
    ):
        message = f"Order {order_id} cannot be cancelled in status: {status}"
        details = {"order_id": order_id, "status": status}
        
        super().__init__(
            message=message,
            code="ORDER_CANNOT_BE_CANCELLED",
            details=details,
            **kwargs
        )


class InvalidOrderStatusError(BadRequestError):
    """Invalid order status transition error."""
    
    def __init__(
        self,
        order_id: int,
        current_status: str,
        new_status: str,
        **kwargs
    ):
        message = f"Cannot transition order {order_id} from {current_status} to {new_status}"
        details = {
            "order_id": order_id,
            "current_status": current_status,
            "new_status": new_status
        }
        
        super().__init__(
            message=message,
            code="INVALID_ORDER_STATUS",
            details=details,
            **kwargs
        )


# ============================================================================
# Domain-Specific Exceptions - Payment
# ============================================================================

class PaymentError(BusinessError):
    """Payment processing error."""
    
    def __init__(
        self,
        message: str = "Payment failed",
        payment_method: str = None,
        transaction_id: str = None,
        **kwargs
    ):
        details = kwargs.get("details", {})
        if payment_method:
            details["payment_method"] = payment_method
        if transaction_id:
            details["transaction_id"] = transaction_id
        
        super().__init__(message, "PAYMENT_ERROR", 402, details=details, **kwargs)


class InsufficientFundsError(PaymentError):
    """Insufficient funds error."""
    
    def __init__(
        self,
        required: float,
        available: float,
        **kwargs
    ):
        message = f"Insufficient funds. Required: {required}, Available: {available}"
        details = {"required": required, "available": available}
        
        super().__init__(
            message=message,
            code="INSUFFICIENT_FUNDS",
            details=details,
            **kwargs
        )


class CardDeclinedError(PaymentError):
    """Card declined error."""
    
    def __init__(
        self,
        reason: str = None,
        **kwargs
    ):
        message = "Card was declined"
        if reason:
            message += f": {reason}"
        
        super().__init__(
            message=message,
            code="CARD_DECLINED",
            **kwargs
        )


# ============================================================================
# Technical Exceptions
# ============================================================================

class DatabaseError(TechnicalError):
    """Database error."""
    
    def __init__(
        self,
        message: str = "Database error occurred",
        operation: str = None,
        query: str = None,
        original_error: Exception = None,
        **kwargs
    ):
        details = kwargs.get("details", {})
        if operation:
            details["operation"] = operation
        if query:
            details["query"] = query
        
        super().__init__(
            message=message,
            code="DATABASE_ERROR",
            details=details,
            cause=original_error,
            **kwargs
        )


class ConnectionError(TechnicalError):
    """Connection error."""
    
    def __init__(
        self,
        service: str = None,
        host: str = None,
        port: int = None,
        original_error: Exception = None,
        **kwargs
    ):
        message = "Connection failed"
        if service:
            message += f" to {service}"
        if host and port:
            message += f" at {host}:{port}"
        
        details = kwargs.get("details", {})
        if service:
            details["service"] = service
        if host:
            details["host"] = host
        if port:
            details["port"] = port
        
        super().__init__(
            message=message,
            code="CONNECTION_ERROR",
            details=details,
            cause=original_error,
            **kwargs
        )


class TimeoutError(TechnicalError):
    """Timeout error."""
    
    def __init__(
        self,
        operation: str = None,
        timeout_seconds: int = None,
        **kwargs
    ):
        message = "Operation timed out"
        if operation:
            message += f": {operation}"
        
        details = kwargs.get("details", {})
        if operation:
            details["operation"] = operation
        if timeout_seconds:
            details["timeout_seconds"] = timeout_seconds
        
        super().__init__(message, "TIMEOUT_ERROR", 408, details=details, **kwargs)


class ConfigurationError(TechnicalError):
    """Configuration error."""
    
    def __init__(
        self,
        message: str = "Configuration error",
        setting: str = None,
        value: Any = None,
        **kwargs
    ):
        details = kwargs.get("details", {})
        if setting:
            details["setting"] = setting
        if value is not None:
            details["value"] = str(value)
        
        super().__init__(message, "CONFIGURATION_ERROR", 500, details=details, **kwargs)


class ThirdPartyServiceError(TechnicalError):
    """Third-party service error."""
    
    def __init__(
        self,
        service: str,
        message: str = "External service error",
        response_code: int = None,
        response_body: str = None,
        **kwargs
    ):
        full_message = f"{service}: {message}"
        details = kwargs.get("details", {})
        details["service"] = service
        if response_code:
            details["response_code"] = response_code
        if response_body:
            details["response_body"] = response_body[:200]  # Truncate long responses
        
        super().__init__(
            message=full_message,
            code="THIRD_PARTY_ERROR",
            details=details,
            **kwargs
        )