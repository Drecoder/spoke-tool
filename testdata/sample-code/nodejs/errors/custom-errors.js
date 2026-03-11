/**
 * Custom Error Classes
 * Demonstrates different patterns for creating and using custom error types
 */

// ============================================================================
// Base Custom Error
// ============================================================================

/**
 * Base application error class
 * Extends the native Error with additional properties
 */
class AppError extends Error {
    constructor(message, statusCode = 500, code = 'INTERNAL_ERROR') {
        super(message);
        this.name = this.constructor.name;
        this.statusCode = statusCode;
        this.code = code;
        this.isOperational = true;
        this.timestamp = new Date().toISOString();
        
        // Capture stack trace
        Error.captureStackTrace(this, this.constructor);
    }
    
    toJSON() {
        return {
            name: this.name,
            message: this.message,
            code: this.code,
            statusCode: this.statusCode,
            timestamp: this.timestamp,
            stack: this.stack
        };
    }
}

// ============================================================================
// HTTP Status Based Errors
// ============================================================================

/**
 * 400 Bad Request Error
 */
class BadRequestError extends AppError {
    constructor(message = 'Bad request', code = 'BAD_REQUEST') {
        super(message, 400, code);
    }
}

/**
 * 401 Unauthorized Error
 */
class UnauthorizedError extends AppError {
    constructor(message = 'Unauthorized', code = 'UNAUTHORIZED') {
        super(message, 401, code);
    }
}

/**
 * 403 Forbidden Error
 */
class ForbiddenError extends AppError {
    constructor(message = 'Forbidden', code = 'FORBIDDEN') {
        super(message, 403, code);
    }
}

/**
 * 404 Not Found Error
 */
class NotFoundError extends AppError {
    constructor(resource = 'Resource', id = null) {
        const message = id 
            ? `${resource} with id ${id} not found`
            : `${resource} not found`;
        super(message, 404, 'NOT_FOUND');
        this.resource = resource;
        this.resourceId = id;
    }
}

/**
 * 409 Conflict Error
 */
class ConflictError extends AppError {
    constructor(message = 'Conflict occurred', code = 'CONFLICT') {
        super(message, 409, code);
    }
}

/**
 * 422 Unprocessable Entity Error
 */
class ValidationError extends AppError {
    constructor(message = 'Validation failed', errors = []) {
        super(message, 422, 'VALIDATION_ERROR');
        this.errors = errors;
    }
    
    addError(field, message) {
        this.errors.push({ field, message });
    }
}

/**
 * 429 Too Many Requests Error
 */
class RateLimitError extends AppError {
    constructor(message = 'Too many requests', retryAfter = 60) {
        super(message, 429, 'RATE_LIMIT_EXCEEDED');
        this.retryAfter = retryAfter;
    }
}

// ============================================================================
// Domain-Specific Errors
// ============================================================================

/**
 * User-related errors
 */
class UserNotFoundError extends NotFoundError {
    constructor(userId) {
        super('User', userId);
        this.code = 'USER_NOT_FOUND';
    }
}

class UserAlreadyExistsError extends ConflictError {
    constructor(email) {
        super(`User with email ${email} already exists`, 'USER_ALREADY_EXISTS');
        this.email = email;
    }
}

class InvalidCredentialsError extends UnauthorizedError {
    constructor() {
        super('Invalid email or password', 'INVALID_CREDENTIALS');
    }
}

class AccountLockedError extends ForbiddenError {
    constructor(userId, lockedUntil) {
        super(`Account locked until ${lockedUntil}`, 'ACCOUNT_LOCKED');
        this.userId = userId;
        this.lockedUntil = lockedUntil;
    }
}

/**
 * Product-related errors
 */
class ProductNotFoundError extends NotFoundError {
    constructor(productId) {
        super('Product', productId);
        this.code = 'PRODUCT_NOT_FOUND';
    }
}

class InsufficientStockError extends BadRequestError {
    constructor(productId, requested, available) {
        super(
            `Insufficient stock for product ${productId}. Requested: ${requested}, Available: ${available}`,
            'INSUFFICIENT_STOCK'
        );
        this.productId = productId;
        this.requested = requested;
        this.available = available;
    }
}

class ProductOutOfStockError extends BadRequestError {
    constructor(productId) {
        super(`Product ${productId} is out of stock`, 'OUT_OF_STOCK');
        this.productId = productId;
    }
}

/**
 * Order-related errors
 */
class OrderNotFoundError extends NotFoundError {
    constructor(orderId) {
        super('Order', orderId);
        this.code = 'ORDER_NOT_FOUND';
    }
}

class OrderCannotBeCancelledError extends BadRequestError {
    constructor(orderId, status) {
        super(
            `Order ${orderId} cannot be cancelled in status: ${status}`,
            'ORDER_CANNOT_BE_CANCELLED'
        );
        this.orderId = orderId;
        this.status = status;
    }
}

class InvalidOrderStatusError extends BadRequestError {
    constructor(orderId, currentStatus, newStatus) {
        super(
            `Cannot transition order ${orderId} from ${currentStatus} to ${newStatus}`,
            'INVALID_ORDER_STATUS'
        );
        this.orderId = orderId;
        this.currentStatus = currentStatus;
        this.newStatus = newStatus;
    }
}

/**
 * Payment-related errors
 */
class PaymentFailedError extends AppError {
    constructor(message = 'Payment failed', paymentMethod = null) {
        super(message, 402, 'PAYMENT_FAILED');
        this.paymentMethod = paymentMethod;
    }
}

class InsufficientFundsError extends PaymentFailedError {
    constructor(required, available) {
        super(`Insufficient funds. Required: ${required}, Available: ${available}`);
        this.required = required;
        this.available = available;
        this.code = 'INSUFFICIENT_FUNDS';
    }
}

class CardDeclinedError extends PaymentFailedError {
        super('Card was declined', 'CREDIT_CARD');
        this.code = 'CARD_DECLINED';
    }
}

class InvalidPaymentMethodError extends BadRequestError {
    constructor(method) {
        super(`Invalid payment method: ${method}`, 'INVALID_PAYMENT_METHOD');
        this.method = method;
    }
}

// ============================================================================
// Database Errors
// ============================================================================

class DatabaseError extends AppError {
    constructor(message = 'Database error occurred', originalError = null) {
        super(message, 500, 'DATABASE_ERROR');
        this.originalError = originalError?.message;
    }
}

class ConnectionError extends DatabaseError {
    constructor(host, port) {
        super(`Failed to connect to database at ${host}:${port}`);
        this.code = 'DB_CONNECTION_ERROR';
        this.host = host;
        this.port = port;
    }
}

class QueryError extends DatabaseError {
    constructor(query, originalError) {
        super(`Query execution failed: ${originalError?.message}`);
        this.code = 'DB_QUERY_ERROR';
        this.query = query;
    }
}

class DuplicateKeyError extends DatabaseError {
    constructor(collection, key, value) {
        super(`Duplicate key error in ${collection}: ${key} = ${value}`);
        this.code = 'DUPLICATE_KEY';
        this.collection = collection;
        this.key = key;
        this.value = value;
    }
}

// ============================================================================
// External Service Errors
// ============================================================================

class ExternalServiceError extends AppError {
    constructor(service, message = 'External service error', statusCode = 503) {
        super(`${service}: ${message}`, statusCode, 'EXTERNAL_SERVICE_ERROR');
        this.service = service;
    }
}

class ServiceTimeoutError extends ExternalServiceError {
    constructor(service, timeout) {
        super(service, `Request timed out after ${timeout}ms`, 504);
        this.code = 'SERVICE_TIMEOUT';
        this.timeout = timeout;
    }
}

class ServiceUnavailableError extends ExternalServiceError {
    constructor(service) {
        super(service, 'Service is unavailable', 503);
        this.code = 'SERVICE_UNAVAILABLE';
    }
}

// ============================================================================
// File System Errors
// ============================================================================

class FileSystemError extends AppError {
    constructor(message = 'File system error', path = null) {
        super(message, 500, 'FILE_SYSTEM_ERROR');
        this.path = path;
    }
}

class FileNotFoundError extends FileSystemError {
    constructor(path) {
        super(`File not found: ${path}`, path);
        this.code = 'FILE_NOT_FOUND';
        this.statusCode = 404;
    }
}

class PermissionDeniedError extends FileSystemError {
    constructor(path) {
        super(`Permission denied: ${path}`, path);
        this.code = 'PERMISSION_DENIED';
        this.statusCode = 403;
    }
}

class FileExistsError extends FileSystemError {
    constructor(path) {
        super(`File already exists: ${path}`, path);
        this.code = 'FILE_EXISTS';
        this.statusCode = 409;
    }
}

// ============================================================================
// Configuration Errors
// ============================================================================

class ConfigurationError extends AppError {
    constructor(message = 'Configuration error') {
        super(message, 500, 'CONFIG_ERROR');
    }
}

class MissingConfigError extends ConfigurationError {
    constructor(key) {
        super(`Missing required configuration: ${key}`);
        this.code = 'MISSING_CONFIG';
        this.key = key;
    }
}

class InvalidConfigError extends ConfigurationError {
        super(`Invalid configuration value for: ${key}`);
        this.code = 'INVALID_CONFIG';
        this.key = key;
        this.value = value;
    }
}

// ============================================================================
// Export all error classes
// ============================================================================

module.exports = {
    // Base
    AppError,
    
    // HTTP status based
    BadRequestError,
    UnauthorizedError,
    ForbiddenError,
    NotFoundError,
    ConflictError,
    ValidationError,
    RateLimitError,
    
    // Domain specific - User
    UserNotFoundError,
    UserAlreadyExistsError,
    InvalidCredentialsError,
    AccountLockedError,
    
    // Domain specific - Product
    ProductNotFoundError,
    InsufficientStockError,
    ProductOutOfStockError,
    
    // Domain specific - Order
    OrderNotFoundError,
    OrderCannotBeCancelledError,
    InvalidOrderStatusError,
    
    // Domain specific - Payment
    PaymentFailedError,
    InsufficientFundsError,
    CardDeclinedError,
    InvalidPaymentMethodError,
    
    // Database
    DatabaseError,
    ConnectionError,
    QueryError,
    DuplicateKeyError,
    
    // External services
    ExternalServiceError,
    ServiceTimeoutError,
    ServiceUnavailableError,
    
    // File system
    FileSystemError,
    FileNotFoundError,
    PermissionDeniedError,
    FileExistsError,
    
    // Configuration
    ConfigurationError,
    MissingConfigError,
    InvalidConfigError
};