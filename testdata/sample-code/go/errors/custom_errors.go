package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// ============================================================================
// Base Error Types
// ============================================================================

// ErrorCode represents a specific error type
type ErrorCode string

const (
	// Validation errors (1xxx)
	ErrValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrRequiredField    ErrorCode = "REQUIRED_FIELD"
	ErrInvalidFormat    ErrorCode = "INVALID_FORMAT"
	ErrInvalidValue     ErrorCode = "INVALID_VALUE"
	ErrOutOfRange       ErrorCode = "OUT_OF_RANGE"

	// Authentication errors (2xxx)
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrTokenInvalid       ErrorCode = "TOKEN_INVALID"
	ErrInsufficientScope  ErrorCode = "INSUFFICIENT_SCOPE"

	// Resource errors (3xxx)
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrAlreadyExists  ErrorCode = "ALREADY_EXISTS"
	ErrConflict       ErrorCode = "CONFLICT"
	ErrResourceLocked ErrorCode = "RESOURCE_LOCKED"

	// Business logic errors (4xxx)
	ErrInsufficientFunds ErrorCode = "INSUFFICIENT_FUNDS"
	ErrLimitExceeded     ErrorCode = "LIMIT_EXCEEDED"
	ErrInvalidOperation  ErrorCode = "INVALID_OPERATION"
	ErrDuplicateEntry    ErrorCode = "DUPLICATE_ENTRY"

	// System errors (5xxx)
	ErrInternal           ErrorCode = "INTERNAL_ERROR"
	ErrDatabase           ErrorCode = "DATABASE_ERROR"
	ErrNetwork            ErrorCode = "NETWORK_ERROR"
	ErrTimeout            ErrorCode = "TIMEOUT"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// Severity level for errors
type Severity string

const (
	SeverityDebug    Severity = "DEBUG"
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityCritical Severity = "CRITICAL"
)

// ============================================================================
// DomainError - Base custom error
// ============================================================================

// DomainError represents a domain-specific error with context
type DomainError struct {
	Code      ErrorCode   `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Severity  Severity    `json:"severity"`
	Timestamp time.Time   `json:"timestamp"`
	Stack     []string    `json:"-"`
	cause     error
}

func NewDomainError(code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:      code,
		Message:   message,
		Severity:  SeverityError,
		Timestamp: time.Now(),
		Stack:     captureStack(2),
	}
}

func (e *DomainError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.cause
}

func (e *DomainError) WithCause(err error) *DomainError {
	e.cause = err
	return e
}

func (e *DomainError) WithDetails(details interface{}) *DomainError {
	e.Details = details
	return e
}

func (e *DomainError) WithSeverity(severity Severity) *DomainError {
	e.Severity = severity
	return e
}

// Is implements error matching for DomainError
func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// ============================================================================
// Validation Errors
// ============================================================================

// ValidationError represents validation failures
type ValidationError struct {
	*DomainError
	Field string      `json:"field"`
	Value interface{} `json:"value,omitempty"`
	Rule  string      `json:"rule,omitempty"`
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		DomainError: NewDomainError(ErrValidationFailed, message),
		Field:       field,
	}
}

func NewRequiredFieldError(field string) *ValidationError {
	return &ValidationError{
		DomainError: NewDomainError(ErrRequiredField, fmt.Sprintf("%s is required", field)),
		Field:       field,
		Rule:        "required",
	}
}

func NewInvalidFormatError(field, format string) *ValidationError {
	return &ValidationError{
		DomainError: NewDomainError(ErrInvalidFormat, fmt.Sprintf("%s must be in %s format", field, format)),
		Field:       field,
		Rule:        "format",
		Value:       format,
	}
}

func NewOutOfRangeError(field string, min, max interface{}) *ValidationError {
	return &ValidationError{
		DomainError: NewDomainError(ErrOutOfRange, fmt.Sprintf("%s must be between %v and %v", field, min, max)),
		Field:       field,
		Rule:        "range",
		Details: map[string]interface{}{
			"min": min,
			"max": max,
		},
	}
}

// ValidationErrors aggregates multiple validation errors
type ValidationErrors struct {
	Errors []*ValidationError `json:"errors"`
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]*ValidationError, 0),
	}
}

func (v *ValidationErrors) Add(err *ValidationError) {
	v.Errors = append(v.Errors, err)
}

func (v *ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}

	var sb strings.Builder
	sb.WriteString("validation failed: ")
	for i, err := range v.Errors {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Message)
	}
	return sb.String()
}

func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// ============================================================================
// Authentication Errors
// ============================================================================

// AuthError represents authentication/authorization errors
type AuthError struct {
	*DomainError
	UserID   string `json:"user_id,omitempty"`
	Resource string `json:"resource,omitempty"`
	Required string `json:"required,omitempty"`
}

func NewUnauthorizedError(message string) *AuthError {
	return &AuthError{
		DomainError: NewDomainError(ErrUnauthorized, message),
	}
}

func NewInvalidCredentialsError() *AuthError {
	return &AuthError{
		DomainError: NewDomainError(ErrInvalidCredentials, "invalid email or password"),
	}
}

func NewTokenExpiredError() *AuthError {
	return &AuthError{
		DomainError: NewDomainError(ErrTokenExpired, "token has expired"),
	}
}

func NewInsufficientScopeError(required, actual string) *AuthError {
	return &AuthError{
		DomainError: NewDomainError(ErrInsufficientScope, "insufficient permissions"),
		Required:    required,
		Details: map[string]string{
			"required": required,
			"actual":   actual,
		},
	}
}

// ============================================================================
// Resource Errors
// ============================================================================

// NotFoundError represents resource not found
type NotFoundError struct {
	*DomainError
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
}

func NewNotFoundError(resourceType, id string) *NotFoundError {
	return &NotFoundError{
		DomainError:  NewDomainError(ErrNotFound, fmt.Sprintf("%s not found: %s", resourceType, id)),
		ResourceType: resourceType,
		ResourceID:   id,
	}
}

func (e *NotFoundError) IsNotFound() bool {
	return true
}

// ConflictError represents resource conflict
type ConflictError struct {
	*DomainError
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
}

func NewConflictError(resourceType, id, message string) *ConflictError {
	return &ConflictError{
		DomainError:  NewDomainError(ErrConflict, message),
		ResourceType: resourceType,
		ResourceID:   id,
	}
}

func NewAlreadyExistsError(resourceType, field, value string) *ConflictError {
	return &ConflictError{
		DomainError:  NewDomainError(ErrAlreadyExists, fmt.Sprintf("%s with %s '%s' already exists", resourceType, field, value)),
		ResourceType: resourceType,
		Details: map[string]string{
			"field": field,
			"value": value,
		},
	}
}

// ============================================================================
// Business Logic Errors
// ============================================================================

// BusinessError represents business rule violations
type BusinessError struct {
	*DomainError
	Rule    string                 `json:"rule"`
	Context map[string]interface{} `json:"context"`
}

func NewBusinessError(rule, message string) *BusinessError {
	return &BusinessError{
		DomainError: NewDomainError(ErrInvalidOperation, message),
		Rule:        rule,
		Context:     make(map[string]interface{}),
	}
}

func (e *BusinessError) WithContext(key string, value interface{}) *BusinessError {
	e.Context[key] = value
	return e
}

// InsufficientFundsError for financial operations
type InsufficientFundsError struct {
	*BusinessError
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
	Required  float64 `json:"required"`
	Shortfall float64 `json:"shortfall"`
}

func NewInsufficientFundsError(accountID string, balance, required float64) *InsufficientFundsError {
	shortfall := required - balance
	return &InsufficientFundsError{
		BusinessError: NewBusinessError("insufficient_funds", "insufficient funds for transaction"),
		AccountID:     accountID,
		Balance:       balance,
		Required:      required,
		Shortfall:     shortfall,
	}
}

// LimitExceededError for rate limits and quotas
type LimitExceededError struct {
	*BusinessError
	LimitType string    `json:"limit_type"`
	Limit     int64     `json:"limit"`
	Current   int64     `json:"current"`
	ResetAt   time.Time `json:"reset_at"`
}

func NewLimitExceededError(limitType string, limit, current int64, resetAt time.Time) *LimitExceededError {
	return &LimitExceededError{
		BusinessError: NewBusinessError("limit_exceeded", fmt.Sprintf("%s limit exceeded", limitType)),
		LimitType:     limitType,
		Limit:         limit,
		Current:       current,
		ResetAt:       resetAt,
	}
}

// ============================================================================
// System Errors
// ============================================================================

// SystemError represents internal system errors
type SystemError struct {
	*DomainError
	Component string `json:"component"`
	Operation string `json:"operation"`
}

func NewInternalError(message string) *SystemError {
	return &SystemError{
		DomainError: NewDomainError(ErrInternal, message),
	}
}

func NewDatabaseError(operation string, err error) *SystemError {
	return &SystemError{
		DomainError: NewDomainError(ErrDatabase, fmt.Sprintf("database error during %s", operation)).WithCause(err),
		Component:   "database",
		Operation:   operation,
	}
}

func NewNetworkError(operation string, err error) *SystemError {
	return &SystemError{
		DomainError: NewDomainError(ErrNetwork, fmt.Sprintf("network error during %s", operation)).WithCause(err),
		Component:   "network",
		Operation:   operation,
	}
}

func NewTimeoutError(operation string, timeout time.Duration) *SystemError {
	return &SystemError{
		DomainError: NewDomainError(ErrTimeout, fmt.Sprintf("timeout after %v during %s", timeout, operation)),
		Component:   "timeout",
		Operation:   operation,
	}
}

// ============================================================================
// HTTP Error Mapping
// ============================================================================

// HTTPError represents an error with HTTP status code
type HTTPError struct {
	error
	StatusCode int `json:"-"`
}

func (e HTTPError) Error() string {
	return e.error.Error()
}

func (e HTTPError) Status() int {
	return e.StatusCode
}

// MapErrorToHTTP converts domain errors to HTTP responses
func MapErrorToHTTP(err error) (int, interface{}) {
	var statusCode int
	var response map[string]interface{}

	switch e := err.(type) {
	case *ValidationError:
		statusCode = http.StatusBadRequest
		response = map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
			"field":   e.Field,
		}

	case *ValidationErrors:
		statusCode = http.StatusBadRequest
		response = map[string]interface{}{
			"code":    ErrValidationFailed,
			"message": e.Error(),
			"errors":  e.Errors,
		}

	case *AuthError:
		switch e.Code {
		case ErrUnauthorized, ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
		case ErrInsufficientScope:
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusUnauthorized
		}
		response = map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		}

	case *NotFoundError:
		statusCode = http.StatusNotFound
		response = map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		}

	case *ConflictError:
		statusCode = http.StatusConflict
		response = map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		}

	case *BusinessError:
		statusCode = http.StatusUnprocessableEntity
		response = map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
			"rule":    e.Rule,
		}

	case *InsufficientFundsError:
		statusCode = http.StatusPaymentRequired
		response = map[string]interface{}{
			"code":      e.Code,
			"message":   e.Message,
			"balance":   e.Balance,
			"required":  e.Required,
			"shortfall": e.Shortfall,
		}

	case *LimitExceededError:
		statusCode = http.StatusTooManyRequests
		response = map[string]interface{}{
			"code":     e.Code,
			"message":  e.Message,
			"limit":    e.Limit,
			"current":  e.Current,
			"reset_at": e.ResetAt,
		}

	case *SystemError:
		statusCode = http.StatusInternalServerError
		response = map[string]interface{}{
			"code":    e.Code,
			"message": "an internal error occurred",
		}

	default:
		statusCode = http.StatusInternalServerError
		response = map[string]interface{}{
			"code":    ErrInternal,
			"message": "an unexpected error occurred",
		}
	}

	// Add request ID if available (would come from context)
	if requestID := getRequestID(); requestID != "" {
		response["request_id"] = requestID
	}

	return statusCode, response
}

// Helper function (would be implemented in actual code)
func getRequestID() string {
	return "" // Placeholder
}

// ============================================================================
// Error Wrapping Utilities
// ============================================================================

// WrapError wraps an error with additional context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorWithCode wraps an error with a domain error
func WrapErrorWithCode(err error, code ErrorCode, message string) *DomainError {
	if err == nil {
		return nil
	}

	domainErr := NewDomainError(code, message)
	domainErr.cause = err
	return domainErr
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	var nf *NotFoundError
	return errors.As(err, &nf)
}

// IsValidation checks if an error is a validation error
func IsValidation(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// IsAuth checks if an error is an authentication error
func IsAuth(err error) bool {
	var ae *AuthError
	return errors.As(err, &ae)
}

// IsConflict checks if an error is a conflict error
func IsConflict(err error) bool {
	var ce *ConflictError
	return errors.As(err, &ce)
}

// ============================================================================
// Stack Trace Capture
// ============================================================================

func captureStack(skip int) []string {
	var stack []string
	for i := skip; i < skip+10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}
	return stack
}

// ============================================================================
// Error Handling Middleware
// ============================================================================

// ErrorHandler middleware for HTTP handlers
type ErrorHandler struct {
	logger interface {
		Error(msg string, fields ...interface{})
	}
}

func NewErrorHandler(logger interface {
	Error(msg string, fields ...interface{})
}) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (h *ErrorHandler) Handle(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			h.handleError(w, r, err)
		}
	}
}

func (h *ErrorHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	statusCode, response := MapErrorToHTTP(err)

	// Log error
	h.logger.Error("request failed",
		"error", err.Error(),
		"status", statusCode,
		"path", r.URL.Path,
		"method", r.Method,
	)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// Retryable Error Detection
// ============================================================================

// IsRetryable determines if an error can be retried
func IsRetryable(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Code {
		case ErrTimeout, ErrNetwork, ErrServiceUnavailable:
			return true
		case ErrInternal:
			// Some internal errors might be retryable
			return true
		}
	}

	// Network errors are retryable
	var netErr interface {
		Timeout() bool
		Temporary() bool
	}
	if errors.As(err, &netErr) && (netErr.Timeout() || netErr.Temporary()) {
		return true
	}

	return false
}

// ============================================================================
// Error Aggregation
// ============================================================================

// AggregateError collects multiple errors
type AggregateError struct {
	Errors []error `json:"errors"`
}

func NewAggregateError() *AggregateError {
	return &AggregateError{
		Errors: make([]error, 0),
	}
}

func (a *AggregateError) Add(err error) {
	if err != nil {
		a.Errors = append(a.Errors, err)
	}
}

func (a *AggregateError) Error() string {
	if len(a.Errors) == 0 {
		return "no errors"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d errors occurred", len(a.Errors)))
	for i, err := range a.Errors {
		sb.WriteString(fmt.Sprintf("\n  %d: %v", i+1, err))
	}
	return sb.String()
}

func (a *AggregateError) HasErrors() bool {
	return len(a.Errors) > 0
}

// ============================================================================
// Error Factories
// ============================================================================

// Common error factories
var (
	ErrEmailRequired    = NewRequiredFieldError("email")
	ErrPasswordRequired = NewRequiredFieldError("password")
	ErrNameRequired     = NewRequiredFieldError("name")

	ErrInvalidEmail = NewInvalidFormatError("email", "email")
	ErrInvalidPhone = NewInvalidFormatError("phone", "E.164")
	ErrInvalidURL   = NewInvalidFormatError("url", "URL")

	ErrUserNotFound    = NewNotFoundError("user", "unknown")
	ErrProductNotFound = NewNotFoundError("product", "unknown")
	ErrOrderNotFound   = NewNotFoundError("order", "unknown")

	ErrEmailAlreadyExists = NewAlreadyExistsError("user", "email", "")
)

// ============================================================================
// Error Examples
// ============================================================================

func ExampleUsage() error {
	// Basic domain error
	err := NewDomainError(ErrValidationFailed, "invalid input")

	// With cause
	cause := errors.New("underlying error")
	err = NewDomainError(ErrDatabase, "query failed").WithCause(cause)

	// Validation error
	err = NewValidationError("email", "invalid email format")

	// Multiple validation errors
	ve := NewValidationErrors()
	ve.Add(NewRequiredFieldError("name"))
	ve.Add(NewInvalidFormatError("email", "email"))

	// Business error with context
	err = NewBusinessError("max_amount", "amount exceeds limit").
		WithContext("max_allowed", 1000).
		WithContext("requested", 1500)

	// Check error types
	switch e := err.(type) {
	case *NotFoundError:
		fmt.Printf("resource not found: %s/%s", e.ResourceType, e.ResourceID)
	case *ValidationError:
		fmt.Printf("validation failed for field %s: %s", e.Field, e.Message)
	case *AuthError:
		fmt.Printf("auth error: %s", e.Message)
	}

	// Check with errors.Is
	if errors.Is(err, ErrUserNotFound) {
		// Handle not found
	}

	return err
}

// ============================================================================
// JSON Marshaling
// ============================================================================

func (e *DomainError) MarshalJSON() ([]byte, error) {
	type Alias DomainError
	return json.Marshal(&struct {
		*Alias
		Error string `json:"error"`
	}{
		Alias: (*Alias)(e),
		Error: e.Error(),
	})
}

// ============================================================================
// Error Stringer
// ============================================================================

func (c ErrorCode) String() string {
	return string(c)
}

func (s Severity) String() string {
	return string(s)
}
