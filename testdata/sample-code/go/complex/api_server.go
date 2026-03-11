package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// Domain Models
// ============================================================================

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose in JSON
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product represents a product in the catalog
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Stock       int       `json:"stock"`
	SKU         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Order represents a customer order
type Order struct {
	ID              int           `json:"id"`
	UserID          int           `json:"user_id"`
	Items           []OrderItem   `json:"items"`
	Total           float64       `json:"total"`
	Status          string        `json:"status"`
	ShippingAddress Address       `json:"shipping_address"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Address represents a shipping address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zip_code"`
	Country    string `json:"country"`
}

// ============================================================================
// Database Repository
// ============================================================================

type Repository interface {
	// User operations
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
	
	// Product operations
	GetProduct(ctx context.Context, id int) (*Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id int) error
	ListProducts(ctx context.Context, filter ProductFilter) ([]*Product, error)
	
	// Order operations
	GetOrder(ctx context.Context, id int) (*Order, error)
	CreateOrder(ctx context.Context, order *Order) error
	UpdateOrderStatus(ctx context.Context, id int, status string) error
	ListUserOrders(ctx context.Context, userID int) ([]*Order, error)
	
	// Health check
	Ping(ctx context.Context) error
}

type ProductFilter struct {
	Category string
	MinPrice float64
	MaxPrice float64
	InStock  *bool
	Limit    int
	Offset   int
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
	user := &User{}
	query := `SELECT id, email, password, name, role, active, created_at, updated_at 
	          FROM users WHERE id = $1`
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, password, name, role, active, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	
	err := r.db.QueryRowContext(ctx, query,
		user.Email, user.Password, user.Name, user.Role, user.Active,
		time.Now(), time.Now(),
	).Scan(&user.ID)
	
	if err != nil {
		// Check for duplicate email
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

// Additional repository methods would be implemented here...

func (r *PostgresRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// ============================================================================
// Service Layer
// ============================================================================

type UserService struct {
	repo Repository
}

func NewUserService(repo Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password, name string) (*User, error) {
	// Validate input
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	
	// Check if user already exists
	existing, _ := s.repo.GetUserByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create user
	user := &User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
		Role:     "user",
		Active:   true,
	}
	
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	
	user.Password = "" // Don't return password hash
	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	
	if !user.Active {
		return nil, errors.New("account is disabled")
	}
	
	user.Password = "" // Don't return password hash
	return user, nil
}

type ProductService struct {
	repo Repository
}

func NewProductService(repo Repository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (*Product, error) {
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %d", id)
	}
	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, filter ProductFilter) ([]*Product, error) {
	// Set defaults
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	
	return s.repo.ListProducts(ctx, filter)
}

// ============================================================================
// Authentication Middleware
// ============================================================================

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type AuthMiddleware struct {
	jwtSecret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: []byte(secret)}
}

func (m *AuthMiddleware) GenerateToken(user *User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "api-server",
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.jwtSecret)
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		// Remove 'Bearer ' prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})
		
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		// Add claims to context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("user").(*Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			
			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

// ============================================================================
// Middleware
// ============================================================================

type Middleware struct {
	logger *log.Logger
}

func NewMiddleware(logger *log.Logger) *Middleware {
	return &Middleware{logger: logger}
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create response writer wrapper to capture status code
		rw := &responseWriter{w, http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start)
		m.logger.Printf("%s %s %d %s %v",
			r.Method, r.URL.Path, rw.status, r.RemoteAddr, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (m *Middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Printf("panic: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	limiter := NewRateLimiter(100, time.Minute)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow(r.RemoteAddr) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Simple rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

type visitor struct {
	count     int
	lastSeen  time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	
	// Cleanup old entries
	go func() {
		for {
			time.Sleep(time.Minute)
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > rl.window {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	
	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
		return true
	}
	
	// Reset if window passed
	if time.Since(v.lastSeen) > rl.window {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}
	
	v.lastSeen = time.Now()
	v.count++
	
	return v.count <= rl.limit
}

// ============================================================================
// HTTP Handlers
// ============================================================================

type Handlers struct {
	userService    *UserService
	productService *ProductService
	authMiddleware *AuthMiddleware
	logger         *log.Logger
}

func NewHandlers(
	userService *UserService,
	productService *ProductService,
	authMiddleware *AuthMiddleware,
	logger *log.Logger,
) *Handlers {
	return &Handlers{
		userService:    userService,
		productService: productService,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

// Health check handler
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// User handlers
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	user, err := h.userService.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Generate token
	token, err := h.authMiddleware.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	token, err := h.authMiddleware.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Fetch fresh user data (in a real app, you'd have a user service method for this)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
	})
}

// Product handlers
func (h *Handlers) ListProducts(w http.ResponseWriter, r *http.Request) {
	filter := ProductFilter{
		Category: r.URL.Query().Get("category"),
		Limit:    20,
		Offset:   0,
	}
	
	if limit := r.URL.Query().Get("limit"); limit != "" {
		fmt.Sscanf(limit, "%d", &filter.Limit)
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		fmt.Sscanf(offset, "%d", &filter.Offset)
	}
	
	products, err := h.productService.ListProducts(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to list products", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *Handlers) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	
	product, err := h.productService.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *Handlers) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate product
	if product.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}
	if product.Price <= 0 {
		http.Error(w, "Price must be positive", http.StatusBadRequest)
		return
	}
	
	// In a real app, you'd save to database
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// ============================================================================
// Server
// ============================================================================

type Server struct {
	router       *mux.Router
	handlers     *Handlers
	middleware   *Middleware
	auth         *AuthMiddleware
	repo         Repository
	httpServer   *http.Server
	logger       *log.Logger
}

func NewServer(repo Repository, logger *log.Logger) *Server {
	// Initialize services
	userService := NewUserService(repo)
	productService := NewProductService(repo)
	
	// Initialize middleware
	authMiddleware := NewAuthMiddleware(os.Getenv("JWT_SECRET"))
	middleware := NewMiddleware(logger)
	
	// Initialize handlers
	handlers := NewHandlers(userService, productService, authMiddleware, logger)
	
	s := &Server{
		router:     mux.NewRouter(),
		handlers:   handlers,
		middleware: middleware,
		auth:       authMiddleware,
		repo:       repo,
		logger:     logger,
	}
	
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Public routes
	s.router.HandleFunc("/health", s.handlers.Health).Methods("GET")
	s.router.HandleFunc("/register", s.handlers.Register).Methods("POST")
	s.router.HandleFunc("/login", s.handlers.Login).Methods("POST")
	
	// Public product routes
	s.router.HandleFunc("/products", s.handlers.ListProducts).Methods("GET")
	s.router.HandleFunc("/products/{id}", s.handlers.GetProduct).Methods("GET")
	
	// Protected routes
	api := s.router.PathPrefix("/api").Subrouter()
	api.Use(s.auth.Authenticate)
	api.Use(s.middleware.Logging)
	api.Use(s.middleware.Recover)
	api.Use(s.middleware.CORS)
	api.Use(s.middleware.RateLimit)
	
	api.HandleFunc("/profile", s.handlers.GetProfile).Methods("GET")
	
	// Admin routes
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(s.auth.RequireRole("admin"))
	
	admin.HandleFunc("/products", s.handlers.CreateProduct).Methods("POST")
	// Add more admin routes...
}

func (s *Server) Start(addr string) error {
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	s.logger.Printf("Server starting on %s", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Println("Shutting down server...")
	
	// Close database connection
	if s.repo != nil {
		if dbRepo, ok := s.repo.(*PostgresRepository); ok {
			dbRepo.db.Close()
		}
	}
	
	return s.httpServer.Shutdown(ctx)
}

// ============================================================================
// Main
// ============================================================================

func main() {
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags|log.Lshortfile)
	
	// Load configuration
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	
	dbConnStr := fmt.Sprintf(
		"host=%s port=5432 user=postgres password=postgres dbname=myapp sslmode=disable",
		dbHost,
	)
	
	// Initialize repository
	repo, err := NewPostgresRepository(dbConnStr)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Create server
	server := NewServer(repo, logger)
	
	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		logger.Println("Received shutdown signal")
		cancel()
		
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Printf("Error during shutdown: %v", err)
		}
	}()
	
	// Start server
	if err := server.Start(":8080"); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server failed: %v", err)
	}
	
	logger.Println("Server stopped")
}

// Helper function for strconv.Atoi
import "strconv"