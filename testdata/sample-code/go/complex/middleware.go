package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// ============================================================================
// Context Keys
// ============================================================================

type contextKey string

func (c contextKey) String() string {
	return "middleware." + string(c)
}

var (
	RequestIDKey  = contextKey("request_id")
	UserIDKey     = contextKey("user_id")
	UserEmailKey  = contextKey("user_email")
	UserRoleKey   = contextKey("user_role")
	StartTimeKey  = contextKey("start_time")
	RequestLogger = contextKey("logger")
)

// ============================================================================
// Request/Response Logging
// ============================================================================

// Logger middleware for logging HTTP requests
type Logger struct {
	logger *log.Logger
}

func NewLogger(logger *log.Logger) *Logger {
	if logger == nil {
		logger = log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	}
	return &Logger{logger: logger}
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Add start time to context
		ctx := context.WithValue(r.Context(), StartTimeKey, start)
		r = r.WithContext(ctx)

		// Process request
		next.ServeHTTP(rw, r)

		// Log after request completes
		duration := time.Since(start)

		// Get request ID if present
		reqID, _ := r.Context().Value(RequestIDKey).(string)
		if reqID == "" {
			reqID = "-"
		}

		l.logger.Printf(
			"%s %s %d %s %s %s %v",
			reqID,
			r.Method,
			rw.status,
			r.URL.Path,
			r.RemoteAddr,
			r.UserAgent(),
			duration,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// ============================================================================
// Request ID
// ============================================================================

// RequestID middleware adds a unique ID to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get ID from header or generate new one
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set response header
		w.Header().Set("X-Request-ID", requestID)

		// Add to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDKey).(string)
	return id
}

// ============================================================================
// Recovery / Panic Handling
// ============================================================================

// Recovery middleware handles panics gracefully
type Recovery struct {
	logger *log.Logger
}

func NewRecovery(logger *log.Logger) *Recovery {
	if logger == nil {
		logger = log.New(os.Stdout, "[PANIC] ", log.LstdFlags)
	}
	return &Recovery{logger: logger}
}

func (r *Recovery) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log stack trace
				stack := debug.Stack()
				r.logger.Printf("panic: %v\n%s", err, stack)

				// Return error response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "internal server error",
				})
			}
		}()

		next.ServeHTTP(w, req)
	})
}

// ============================================================================
// CORS
// ============================================================================

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(config *CORSConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowedOrigin := ""
			for _, o := range config.AllowedOrigins {
				if o == "*" || o == origin {
					allowedOrigin = origin
					break
				}
			}

			if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				if len(config.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
				}
			}

			// Handle preflight
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
				w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ============================================================================
// Authentication
// ============================================================================

// Claims represents JWT claims
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret    string
	TokenLookup  string // "header:Authorization", "query:token", "cookie:token"
	ContextKey   string
	Unauthorized http.HandlerFunc
}

// Auth middleware handles JWT authentication
type Auth struct {
	config *AuthConfig
}

func NewAuth(config *AuthConfig) *Auth {
	if config == nil {
		config = &AuthConfig{
			JWTSecret:   os.Getenv("JWT_SECRET"),
			TokenLookup: "header:Authorization",
			ContextKey:  "user",
			Unauthorized: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
			},
		}
	}

	// Default secret for development (should be overridden in production)
	if config.JWTSecret == "" {
		config.JWTSecret = "dev-secret-do-not-use-in-production"
	}

	return &Auth{config: config}
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := a.extractToken(r)
		if err != nil {
			a.config.Unauthorized(w, r)
			return
		}

		claims, err := a.validateToken(token)
		if err != nil {
			a.config.Unauthorized(w, r)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), a.config.ContextKey, claims)
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) extractToken(r *http.Request) (string, error) {
	parts := strings.Split(a.config.TokenLookup, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token lookup format")
	}

	switch parts[0] {
	case "header":
		return a.extractFromHeader(r, parts[1])
	case "query":
		return a.extractFromQuery(r, parts[1])
	case "cookie":
		return a.extractFromCookie(r, parts[1])
	default:
		return "", fmt.Errorf("unsupported token lookup: %s", parts[0])
	}
}

func (a *Auth) extractFromHeader(r *http.Request, header string) (string, error) {
	auth := r.Header.Get(header)
	if auth == "" {
		return "", fmt.Errorf("empty authorization header")
	}

	// Remove 'Bearer ' prefix
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer "), nil
	}

	return auth, nil
}

func (a *Auth) extractFromQuery(r *http.Request, param string) (string, error) {
	token := r.URL.Query().Get(param)
	if token == "" {
		return "", fmt.Errorf("empty token in query")
	}
	return token, nil
}

func (a *Auth) extractFromCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("cookie not found: %w", err)
	}
	return cookie.Value, nil
}

func (a *Auth) validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) int {
	if id, ok := ctx.Value(UserIDKey).(int); ok {
		return id
	}
	return 0
}

// GetUserRole retrieves the user role from context
func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r.Context())

			for _, role := range roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

// ============================================================================
// Basic Auth
// ============================================================================

// BasicAuthConfig holds basic authentication configuration
type BasicAuthConfig struct {
	Realm    string
	Users    map[string]string // username -> password
	Validate func(username, password string) bool
}

// BasicAuth middleware handles HTTP Basic Authentication
func BasicAuth(config *BasicAuthConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = &BasicAuthConfig{
			Realm: "Restricted",
			Users: map[string]string{
				"admin": "password",
			},
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, config.Realm))
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate credentials
			valid := false
			if config.Validate != nil {
				valid = config.Validate(username, password)
			} else {
				expectedPass, exists := config.Users[username]
				valid = exists && subtle.ConstantTimeCompare([]byte(password), []byte(expectedPass)) == 1
			}

			if !valid {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, config.Realm))
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Add username to context
			ctx := context.WithValue(r.Context(), "username", username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ============================================================================
// Rate Limiting
// ============================================================================

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
	keyFunc  func(r *http.Request) string
}

func NewRateLimiter(r rate.Limit, burst int, keyFunc func(r *http.Request) string) *RateLimiter {
	if keyFunc == nil {
		keyFunc = func(r *http.Request) string {
			// Default to IP address
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			return ip
		}
	}

	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
		keyFunc:  keyFunc,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.keyFunc(r)

		limiter := rl.getLimiter(key)

		if !limiter.Allow() {
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.burst))
			w.Header().Set("X-RateLimit-Remaining", "0")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.burst))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", int(limiter.Tokens())))

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Check again in case it was created while we were waiting
		if limiter, exists = rl.limiters[key]; !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[key] = limiter
		}
	}

	return limiter
}

// Cleanup old limiters periodically
func (rl *RateLimiter) Cleanup(maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			rl.mu.Lock()
			for key, limiter := range rl.limiters {
				// If limiter hasn't been used for maxAge, remove it
				// This is a simplification - in practice you'd track last used time
				if limiter.Tokens() == float64(rl.burst) {
					delete(rl.limiters, key)
				}
			}
			rl.mu.Unlock()
		}
	}()
}

// ============================================================================
// Timeout
// ============================================================================

// Timeout middleware adds a timeout to request handling
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				w.WriteHeader(http.StatusGatewayTimeout)
				w.Write([]byte("request timeout"))
			}
		})
	}
}

// ============================================================================
// Compression
// ============================================================================

// Gzip middleware compresses responses
type Gzip struct {
	level int
}

func NewGzip(level int) *Gzip {
	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}
	return &Gzip{level: level}
}

func (g *Gzip) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, g.level)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		gzw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(gzw, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

// ============================================================================
// Security Headers
// ============================================================================

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// Content Type
// ============================================================================

// ContentType middleware ensures proper content type
func ContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}

// JSON is a convenience middleware for JSON responses
func JSON(next http.Handler) http.Handler {
	return ContentType("application/json")(next)
}

// ============================================================================
// Request Size Limiting
// ============================================================================

// MaxBodySize limits the size of request bodies
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// ============================================================================
// IP Whitelist/Blacklist
// ============================================================================

// IPFilter filters requests by IP address
type IPFilter struct {
	whitelist []*net.IPNet
	blacklist []*net.IPNet
	blocked   http.HandlerFunc
}

func NewIPFilter() *IPFilter {
	return &IPFilter{
		blocked: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "access denied", http.StatusForbidden)
		},
	}
}

func (f *IPFilter) WhitelistCIDR(cidr string) error {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	f.whitelist = append(f.whitelist, ipnet)
	return nil
}

func (f *IPFilter) BlacklistCIDR(cidr string) error {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	f.blacklist = append(f.blacklist, ipnet)
	return nil
}

func (f *IPFilter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipStr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			f.blocked(w, r)
			return
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			f.blocked(w, r)
			return
		}

		// Check blacklist first
		for _, cidr := range f.blacklist {
			if cidr.Contains(ip) {
				f.blocked(w, r)
				return
			}
		}

		// If whitelist is not empty, check that
		if len(f.whitelist) > 0 {
			allowed := false
			for _, cidr := range f.whitelist {
				if cidr.Contains(ip) {
					allowed = true
					break
				}
			}
			if !allowed {
				f.blocked(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// Metrics
// ============================================================================

// Metrics middleware collects request metrics
type Metrics struct {
	requests *Counter
	duration *Histogram
	inflight *Gauge
}

func NewMetrics(registry *MetricsRegistry) *Metrics {
	return &Metrics{
		requests: registry.Counter("http_requests_total", "Total HTTP requests"),
		duration: registry.Histogram("http_request_duration_seconds", "HTTP request duration"),
		inflight: registry.Gauge("http_requests_inflight", "In-flight HTTP requests"),
	}
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.inflight.Inc()
		defer m.inflight.Dec()

		m.requests.Inc()

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()

		m.duration.Observe(duration)
	})
}

// Simplified metrics types for demonstration
type Counter struct{}

func (c *Counter) Inc() {}

type Histogram struct{}

func (h *Histogram) Observe(float64) {}

type Gauge struct{}

func (g *Gauge) Inc() {}
func (g *Gauge) Dec() {}

type MetricsRegistry struct{}

func (r *MetricsRegistry) Counter(name, help string) *Counter     { return &Counter{} }
func (r *MetricsRegistry) Histogram(name, help string) *Histogram { return &Histogram{} }
func (r *MetricsRegistry) Gauge(name, help string) *Gauge         { return &Gauge{} }

// ============================================================================
// Chain
// ============================================================================

// Chain helps compose multiple middleware
type Chain struct {
	middlewares []func(http.Handler) http.Handler
}

func NewChain(middlewares ...func(http.Handler) http.Handler) *Chain {
	return &Chain{middlewares: middlewares}
}

func (c *Chain) Then(handler http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}

func (c *Chain) ThenFunc(handlerFunc http.HandlerFunc) http.Handler {
	return c.Then(handlerFunc)
}

func (c *Chain) Append(middlewares ...func(http.Handler) http.Handler) *Chain {
	newChain := make([]func(http.Handler) http.Handler, len(c.middlewares)+len(middlewares))
	copy(newChain, c.middlewares)
	copy(newChain[len(c.middlewares):], middlewares)
	return &Chain{middlewares: newChain}
}

// ============================================================================
// Example Usage
// ============================================================================

func ExampleMiddlewareChain() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create middleware chain
	chain := NewChain(
		RequestID,
		NewLogger(logger).Middleware,
		NewRecovery(logger).Middleware,
		SecurityHeaders,
		JSON,
		NewRateLimiter(rate.Limit(10), 10, nil).Middleware,
		Timeout(30*time.Second),
		CORS(DefaultCORSConfig()),
	)

	// Create handler
	handler := chain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "hello world"}`))
	})

	// Start server
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
