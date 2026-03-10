package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yourusername/spoke-tool/testdata/sample-code/mixed-project/go-server/config"
	"github.com/yourusername/spoke-tool/testdata/sample-code/mixed-project/go-server/handlers"
	"github.com/yourusername/spoke-tool/testdata/sample-code/mixed-project/go-server/middleware"
	"github.com/yourusername/spoke-tool/testdata/sample-code/mixed-project/go-server/models"
)

// ============================================================================
// Application Configuration
// ============================================================================

// Config holds all application configuration
type Config struct {
	Environment  string        `env:"ENVIRONMENT" default:"development"`
	Port         int           `env:"PORT" default:"8080"`
	Host         string        `env:"HOST" default:"0.0.0.0"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" default:"15s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"15s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" default:"60s"`

	// Database
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DBUser     string `env:"DB_USER" default:"postgres"`
	DBPassword string `env:"DB_PASSWORD" default:"postgres"`
	DBName     string `env:"DB_NAME" default:"go_server_db"`
	DBSSLMode  string `env:"DB_SSL_MODE" default:"disable"`

	// Redis
	RedisHost     string `env:"REDIS_HOST" default:"localhost"`
	RedisPort     int    `env:"REDIS_PORT" default:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" default:""`
	RedisDB       int    `env:"REDIS_DB" default:"0"`

	// JWT
	JWTSecret     string `env:"JWT_SECRET" default:"your-secret-key-change-in-production"`
	JWTExpiration int    `env:"JWT_EXPIRATION" default:"24"` // hours

	// Rate Limiting
	RateLimit      int `env:"RATE_LIMIT" default:"100"`
	RateLimitBurst int `env:"RATE_LIMIT_BURST" default:"50"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" default:"info"`
	LogJSON  bool   `env:"LOG_JSON" default:"false"`

	// CORS
	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" default:"*"`

	// Metrics
	EnableMetrics bool   `env:"ENABLE_METRICS" default:"true"`
	MetricsPath   string `env:"METRICS_PATH" default:"/metrics"`

	// Profiling
	EnableProfiling bool `env:"ENABLE_PROFILING" default:"false"`
}

// ============================================================================
// Application Server
// ============================================================================

// Server represents the HTTP server
type Server struct {
	config     *Config
	router     *gin.Engine
	db         *gorm.DB
	logger     *zap.Logger
	httpServer *http.Server
	handlers   *Handlers
}

// Handlers aggregates all handler instances
type Handlers struct {
	User    *handlers.UserHandler
	Auth    *handlers.AuthHandler
	Product *handlers.ProductHandler
	Health  *handlers.HealthHandler
}

// NewServer creates a new server instance
func NewServer(config *Config, logger *zap.Logger) (*Server, error) {
	// Set Gin mode based on environment
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()

	// Initialize database
	db, err := initDatabase(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize handlers
	handlers := &Handlers{
		User:    handlers.NewUserHandler(db, logger),
		Auth:    handlers.NewAuthHandler(db, logger, config.JWTSecret),
		Product: handlers.NewProductHandler(db, logger),
		Health:  handlers.NewHealthHandler(db, logger),
	}

	server := &Server{
		config:   config,
		router:   router,
		db:       db,
		logger:   logger,
		handlers: handlers,
	}

	// Setup middleware and routes
	server.setupMiddleware()
	server.setupRoutes()

	return server, nil
}

// ============================================================================
// Database Initialization
// ============================================================================

func initDatabase(config *Config, logger *zap.Logger) (*gorm.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword,
		config.DBName, config.DBSSLMode,
	)

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate schemas
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	logger.Info("Database connected successfully",
		zap.String("host", config.DBHost),
		zap.String("database", config.DBName),
	)

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
}

// ============================================================================
// Middleware Setup
// ============================================================================

func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Logger middleware
	s.router.Use(middleware.Logger(s.logger))

	// Request ID middleware
	s.router.Use(middleware.RequestID())

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = s.config.CORSAllowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	corsConfig.AllowHeaders = []string{
		"Origin", "Content-Type", "Accept", "Authorization",
		"X-Request-ID", "X-API-Key",
	}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour
	s.router.Use(cors.New(corsConfig))

	// Rate limiting middleware
	s.router.Use(middleware.RateLimiter(s.config.RateLimit, s.config.RateLimitBurst))

	// Timeout middleware
	s.router.Use(middleware.Timeout(30 * time.Second))

	// Security headers middleware
	s.router.Use(middleware.SecurityHeaders())

	// Metrics middleware
	if s.config.EnableMetrics {
		s.router.Use(middleware.Metrics())
	}

	// Profiling
	if s.config.EnableProfiling && s.config.Environment != "production" {
		pprof.Register(s.router)
		s.logger.Info("Profiling enabled")
	}
}

// ============================================================================
// Route Setup
// ============================================================================

func (s *Server) setupRoutes() {
	// Health check endpoints (no auth)
	s.router.GET("/health", s.handlers.Health.Health)
	s.router.GET("/ready", s.handlers.Health.Readiness)

	// API v1 group
	v1 := s.router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/public")
		{
			public.GET("/products", s.handlers.Product.ListProducts)
			public.GET("/products/:id", s.handlers.Product.GetProduct)
		}

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.handlers.User.CreateUser)
			auth.POST("/login", s.handlers.Auth.Login)
			auth.POST("/refresh", s.handlers.Auth.RefreshToken)
			auth.POST("/logout", s.handlers.Auth.Logout)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.Auth(s.config.JWTSecret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/", s.handlers.User.ListUsers)
				users.GET("/:id", s.handlers.User.GetUser)
				users.PUT("/:id", s.handlers.User.UpdateUser)
				users.DELETE("/:id", s.handlers.User.DeleteUser)
				users.GET("/profile", s.handlers.User.GetUser) // Get current user
			}

			// Product routes (admin only)
			products := protected.Group("/products")
			products.Use(middleware.RequireRole("admin"))
			{
				products.POST("/", s.handlers.Product.CreateProduct)
				products.PUT("/:id", s.handlers.Product.UpdateProduct)
				products.DELETE("/:id", s.handlers.Product.DeleteProduct)
			}
		}
	}

	// Metrics endpoint
	if s.config.EnableMetrics {
		s.router.GET(s.config.MetricsPath, middleware.PrometheusHandler())
	}

	// 404 handler
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
		})
	})
}

// ============================================================================
// Server Lifecycle
// ============================================================================

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	s.logger.Info("Starting server",
		zap.String("address", addr),
		zap.String("environment", s.config.Environment),
	)

	// Start server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")

	// Close database connection
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				s.logger.Error("Failed to close database", zap.Error(err))
			}
		}
	}

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.Info("Server stopped gracefully")
	return nil
}

// ============================================================================
// Main Function
// ============================================================================

func main() {
	// Parse command line flags
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := &Config{}
	if configPath != "" {
		if err := config.LoadFromFile(configPath, cfg); err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	} else {
		if err := config.LoadFromEnv(cfg); err != nil {
			log.Fatalf("Failed to load config from environment: %v", err)
		}
	}

	// Initialize logger
	logger, err := initLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create server
	server, err := NewServer(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to create server", zap.Error(err))
	}

	// Start server
	if err := server.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// ============================================================================
// Logger Initialization
// ============================================================================

func initLogger(cfg *Config) (*zap.Logger, error) {
	var level zapcore.Level
	switch cfg.LogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.LogJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

// ============================================================================
// Health Check Handler (additional)
// ============================================================================

// HealthCheckHandler provides detailed health information
func (s *Server) HealthCheckHandler(c *gin.Context) {
	status := "healthy"
	checks := make(map[string]interface{})

	// Check database
	if err := s.db.Exec("SELECT 1").Error; err != nil {
		status = "degraded"
		checks["database"] = map[string]interface{}{
			"status": "down",
			"error":  err.Error(),
		}
	} else {
		sqlDB, _ := s.db.DB()
		stats := sqlDB.Stats()
		checks["database"] = map[string]interface{}{
			"status":           "up",
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
			"max_open":         stats.MaxOpenConnections,
		}
	}

	// Check Redis (if configured)
	if s.config.RedisHost != "" {
		// Add Redis check here if you have Redis client
		checks["redis"] = map[string]interface{}{
			"status": "not_configured",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    status,
		"timestamp": time.Now(),
		"uptime":    time.Since(serverStartTime).String(),
		"version":   "1.0.0",
		"checks":    checks,
	})
}

// Track server start time
var serverStartTime = time.Now()
