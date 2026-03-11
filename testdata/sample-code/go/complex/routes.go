package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"example.com/spoke-tool/internal/auth"
	"example.com/spoke-tool/internal/handlers"
	"example.com/spoke-tool/internal/middleware"
	"example.com/spoke-tool/internal/services"
)

// ============================================================================
// Router Configuration
// ============================================================================

// RouterConfig holds configuration for the router
type RouterConfig struct {
	Environment     string        `json:"environment"`
	AllowedOrigins  []string      `json:"allowed_origins"`
	RateLimit       int           `json:"rate_limit"`
	RateLimitWindow time.Duration `json:"rate_limit_window"`
	Timeout         time.Duration `json:"timeout"`
	MaxBodySize     int64         `json:"max_body_size"`
	EnableMetrics   bool          `json:"enable_metrics"`
	EnableProfiling bool          `json:"enable_profiling"`
	EnableSwagger   bool          `json:"enable_swagger"`
}

// DefaultRouterConfig returns a default configuration
func DefaultRouterConfig() *RouterConfig {
	return &RouterConfig{
		Environment:     "development",
		AllowedOrigins:  []string{"*"},
		RateLimit:       100,
		RateLimitWindow: time.Minute,
		Timeout:         30 * time.Second,
		MaxBodySize:     10 * 1024 * 1024, // 10MB
		EnableMetrics:   true,
		EnableProfiling: false,
		EnableSwagger:   true,
	}
}

// ============================================================================
// Main Router
// ============================================================================

// Router wraps chi router with additional functionality
type Router struct {
	*chi.Mux
	config     *RouterConfig
	logger     *zap.Logger
	services   *services.ServiceContainer
	middleware *middleware.MiddlewareManager
	auth       *auth.JWTAuth
}

// NewRouter creates a new router with all middleware and routes
func NewRouter(
	config *RouterConfig,
	logger *zap.Logger,
	services *services.ServiceContainer,
	middleware *middleware.MiddlewareManager,
	auth *auth.JWTAuth,
) *Router {
	r := &Router{
		Mux:        chi.NewRouter(),
		config:     config,
		logger:     logger,
		services:   services,
		middleware: middleware,
		auth:       auth,
	}

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

// setupMiddleware configures global middleware
func (r *Router) setupMiddleware() {
	// Core middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Compression
	r.Use(middleware.Compress(5))

	// Timeout
	r.Use(middleware.Timeout(r.config.Timeout))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   r.config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting
	r.Use(httprate.Limit(
		r.config.RateLimit,
		r.config.RateLimitWindow,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	// Request size limit
	r.Use(middleware.MaxBodySize(r.config.MaxBodySize))

	// Health check endpoint (bypasses auth)
	r.Get("/health", r.healthCheck)
	r.Get("/ready", r.readinessCheck)

	// Metrics
	if r.config.EnableMetrics {
		r.Mount("/metrics", promhttp.Handler())
	}

	// Debug profiling
	if r.config.EnableProfiling && r.config.Environment == "development" {
		r.Mount("/debug", middleware.Profiler())
	}

	// Swagger documentation
	if r.config.EnableSwagger {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
	}
}

// setupRoutes configures all API routes
func (r *Router) setupRoutes() {
	// API version group
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Use(r.middleware.RateLimit("public", 60, time.Minute))
			r.Mount("/auth", r.authRoutes())
			r.Mount("/public", r.publicRoutes())
		})

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(r.auth.Verifier())
			r.Use(r.auth.Authenticator)

			// User routes
			r.Mount("/users", r.userRoutes())

			// Account routes
			r.Mount("/accounts", r.accountRoutes())

			// Transaction routes
			r.Mount("/transactions", r.transactionRoutes())

			// Product routes
			r.Mount("/products", r.productRoutes())

			// Order routes
			r.Mount("/orders", r.orderRoutes())

			// Notification routes
			r.Mount("/notifications", r.notificationRoutes())
		})

		// Admin routes (require admin role)
		r.Group(func(r chi.Router) {
			r.Use(r.auth.Verifier())
			r.Use(r.auth.Authenticator)
			r.Use(r.middleware.RequireRole("admin"))

			r.Mount("/admin", r.adminRoutes())
		})

		// Webhook routes (no auth, signature verification)
		r.Mount("/webhooks", r.webhookRoutes())
	})
}

// ============================================================================
// Health Routes
// ============================================================================

func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	render.JSON(w, req, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (r *Router) readinessCheck(w http.ResponseWriter, req *http.Request) {
	// Check database connection
	if err := r.services.DB.Ping(); err != nil {
		render.Status(req, http.StatusServiceUnavailable)
		render.JSON(w, req, map[string]string{
			"status": "not ready",
			"error":  "database connection failed",
		})
		return
	}

	render.JSON(w, req, map[string]string{
		"status": "ready",
	})
}

// ============================================================================
// Auth Routes
// ============================================================================

func (r *Router) authRoutes() chi.Router {
	router := chi.NewRouter()
	authHandler := handlers.NewAuthHandler(r.services.Auth, r.logger)

	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)
	router.Post("/logout", authHandler.Logout)
	router.Post("/refresh", authHandler.RefreshToken)
	router.Post("/verify-email", authHandler.VerifyEmail)
	router.Post("/resend-verification", authHandler.ResendVerification)
	router.Post("/forgot-password", authHandler.ForgotPassword)
	router.Post("/reset-password", authHandler.ResetPassword)
	router.Post("/change-password", authHandler.ChangePassword)

	// OAuth routes
	router.Get("/{provider}", authHandler.OAuthRedirect)
	router.Get("/{provider}/callback", authHandler.OAuthCallback)

	return router
}

// ============================================================================
// Public Routes
// ============================================================================

func (r *Router) publicRoutes() chi.Router {
	router := chi.NewRouter()
	publicHandler := handlers.NewPublicHandler(r.services.Public, r.logger)

	// Public product catalog
	router.Get("/products", publicHandler.ListProducts)
	router.Get("/products/{id}", publicHandler.GetProduct)
	router.Get("/categories", publicHandler.ListCategories)
	router.Get("/categories/{slug}", publicHandler.GetCategory)

	// Public content
	router.Get("/blog", publicHandler.ListBlogPosts)
	router.Get("/blog/{slug}", publicHandler.GetBlogPost)
	router.Get("/faq", publicHandler.ListFAQ)

	// Contact
	router.Post("/contact", publicHandler.Contact)

	return router
}

// ============================================================================
// User Routes
// ============================================================================

func (r *Router) userRoutes() chi.Router {
	router := chi.NewRouter()
	userHandler := handlers.NewUserHandler(r.services.User, r.logger)

	router.Get("/profile", userHandler.GetProfile)
	router.Put("/profile", userHandler.UpdateProfile)
	router.Patch("/profile", userHandler.PartialUpdateProfile)
	router.Delete("/profile", userHandler.DeleteProfile)

	router.Get("/settings", userHandler.GetSettings)
	router.Put("/settings", userHandler.UpdateSettings)

	router.Get("/preferences", userHandler.GetPreferences)
	router.Put("/preferences", userHandler.UpdatePreferences)

	router.Post("/avatar", userHandler.UploadAvatar)
	router.Delete("/avatar", userHandler.DeleteAvatar)

	// User's accounts
	router.Get("/accounts", userHandler.ListUserAccounts)
	router.Post("/accounts", userHandler.CreateAccount)

	// User's transactions
	router.Get("/transactions", userHandler.ListUserTransactions)
	router.Get("/transactions/{id}", userHandler.GetTransaction)

	// User's orders
	router.Get("/orders", userHandler.ListUserOrders)
	router.Get("/orders/{id}", userHandler.GetOrder)

	// User's notifications
	router.Get("/notifications", userHandler.ListNotifications)
	router.Put("/notifications/{id}/read", userHandler.MarkNotificationRead)
	router.Put("/notifications/read-all", userHandler.MarkAllNotificationsRead)

	return router
}

// ============================================================================
// Account Routes
// ============================================================================

func (r *Router) accountRoutes() chi.Router {
	router := chi.NewRouter()
	accountHandler := handlers.NewAccountHandler(r.services.Account, r.logger)

	router.Get("/", accountHandler.ListAccounts)
	router.Post("/", accountHandler.CreateAccount)
	router.Get("/{id}", accountHandler.GetAccount)
	router.Put("/{id}", accountHandler.UpdateAccount)
	router.Delete("/{id}", accountHandler.CloseAccount)

	router.Get("/{id}/balance", accountHandler.GetBalance)
	router.Get("/{id}/statement", accountHandler.GetStatement)
	router.Get("/{id}/transactions", accountHandler.ListTransactions)

	router.Post("/{id}/deposit", accountHandler.Deposit)
	router.Post("/{id}/withdraw", accountHandler.Withdraw)
	router.Post("/{id}/transfer", accountHandler.Transfer)

	return router
}

// ============================================================================
// Transaction Routes
// ============================================================================

func (r *Router) transactionRoutes() chi.Router {
	router := chi.NewRouter()
	txHandler := handlers.NewTransactionHandler(r.services.Transaction, r.logger)

	router.Get("/", txHandler.ListTransactions)
	router.Get("/{id}", txHandler.GetTransaction)
	router.Post("/{id}/cancel", txHandler.CancelTransaction)
	router.Post("/{id}/reverse", txHandler.ReverseTransaction)

	router.Get("/pending", txHandler.ListPendingTransactions)
	router.Get("/failed", txHandler.ListFailedTransactions)

	router.Get("/export", txHandler.ExportTransactions)

	return router
}

// ============================================================================
// Product Routes
// ============================================================================

func (r *Router) productRoutes() chi.Router {
	router := chi.NewRouter()
	productHandler := handlers.NewProductHandler(r.services.Product, r.logger)

	// Categories
	router.Route("/categories", func(r chi.Router) {
		r.Get("/", productHandler.ListCategories)
		r.Post("/", productHandler.CreateCategory)
		r.Get("/{id}", productHandler.GetCategory)
		r.Put("/{id}", productHandler.UpdateCategory)
		r.Delete("/{id}", productHandler.DeleteCategory)
	})

	// Products
	router.Route("/", func(r chi.Router) {
		r.Get("/", productHandler.ListProducts)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Patch("/{id}/stock", productHandler.UpdateStock)

		// Product variants
		r.Route("/{id}/variants", func(r chi.Router) {
			r.Get("/", productHandler.ListVariants)
			r.Post("/", productHandler.CreateVariant)
			r.Get("/{variantId}", productHandler.GetVariant)
			r.Put("/{variantId}", productHandler.UpdateVariant)
			r.Delete("/{variantId}", productHandler.DeleteVariant)
		})

		// Product reviews
		r.Route("/{id}/reviews", func(r chi.Router) {
			r.Get("/", productHandler.ListReviews)
			r.Post("/", productHandler.CreateReview)
			r.Get("/{reviewId}", productHandler.GetReview)
			r.Put("/{reviewId}", productHandler.UpdateReview)
			r.Delete("/{reviewId}", productHandler.DeleteReview)
			r.Post("/{reviewId}/helpful", productHandler.MarkHelpful)
		})

		// Product images
		r.Post("/{id}/images", productHandler.UploadImages)
		r.Delete("/{id}/images/{imageId}", productHandler.DeleteImage)
		r.Put("/{id}/images/{imageId}/primary", productHandler.SetPrimaryImage)
	})

	return router
}

// ============================================================================
// Order Routes
// ============================================================================

func (r *Router) orderRoutes() chi.Router {
	router := chi.NewRouter()
	orderHandler := handlers.NewOrderHandler(r.services.Order, r.logger)

	router.Get("/", orderHandler.ListOrders)
	router.Post("/", orderHandler.CreateOrder)
	router.Get("/{id}", orderHandler.GetOrder)
	router.Put("/{id}", orderHandler.UpdateOrder)
	router.Delete("/{id}", orderHandler.CancelOrder)

	router.Post("/{id}/checkout", orderHandler.Checkout)
	router.Post("/{id}/pay", orderHandler.ProcessPayment)
	router.Post("/{id}/refund", orderHandler.RefundOrder)

	router.Get("/{id}/invoice", orderHandler.GetInvoice)
	router.Get("/{id}/tracking", orderHandler.GetTracking)

	// Order items
	router.Route("/{id}/items", func(r chi.Router) {
		r.Get("/", orderHandler.ListOrderItems)
		r.Post("/", orderHandler.AddOrderItem)
		r.Put("/{itemId}", orderHandler.UpdateOrderItem)
		r.Delete("/{itemId}", orderHandler.RemoveOrderItem)
	})

	// Order status updates (admin only)
	router.Route("/{id}/status", func(r chi.Router) {
		r.Use(r.middleware.RequireRole("admin", "manager"))
		r.Put("/", orderHandler.UpdateOrderStatus)
		r.Post("/ship", orderHandler.MarkAsShipped)
		r.Post("/deliver", orderHandler.MarkAsDelivered)
	})

	return router
}

// ============================================================================
// Admin Routes
// ============================================================================

func (r *Router) adminRoutes() chi.Router {
	router := chi.NewRouter()
	adminHandler := handlers.NewAdminHandler(r.services.Admin, r.logger)

	// User management
	router.Route("/users", func(r chi.Router) {
		r.Get("/", adminHandler.ListUsers)
		r.Get("/{id}", adminHandler.GetUser)
		r.Put("/{id}", adminHandler.UpdateUser)
		r.Delete("/{id}", adminHandler.DeleteUser)
		r.Post("/{id}/suspend", adminHandler.SuspendUser)
		r.Post("/{id}/activate", adminHandler.ActivateUser)
		r.Put("/{id}/role", adminHandler.ChangeUserRole)
	})

	// System settings
	router.Route("/settings", func(r chi.Router) {
		r.Get("/", adminHandler.GetSettings)
		r.Put("/", adminHandler.UpdateSettings)
		r.Get("/{key}", adminHandler.GetSetting)
		r.Put("/{key}", adminHandler.UpdateSetting)
	})

	// Audit logs
	router.Get("/audit-logs", adminHandler.ListAuditLogs)
	router.Get("/audit-logs/{id}", adminHandler.GetAuditLog)
	router.Get("/audit-logs/export", adminHandler.ExportAuditLogs)

	// Metrics and monitoring
	router.Get("/metrics/system", adminHandler.GetSystemMetrics)
	router.Get("/metrics/business", adminHandler.GetBusinessMetrics)
	router.Get("/metrics/performance", adminHandler.GetPerformanceMetrics)

	// Background jobs
	router.Get("/jobs", adminHandler.ListJobs)
	router.Get("/jobs/{id}", adminHandler.GetJob)
	router.Post("/jobs/{id}/retry", adminHandler.RetryJob)
	router.Post("/jobs/{id}/cancel", adminHandler.CancelJob)

	// Feature flags
	router.Get("/features", adminHandler.ListFeatures)
	router.Post("/features", adminHandler.CreateFeature)
	router.Put("/features/{key}", adminHandler.UpdateFeature)
	router.Delete("/features/{key}", adminHandler.DeleteFeature)

	// Announcements
	router.Get("/announcements", adminHandler.ListAnnouncements)
	router.Post("/announcements", adminHandler.CreateAnnouncement)
	router.Put("/announcements/{id}", adminHandler.UpdateAnnouncement)
	router.Delete("/announcements/{id}", adminHandler.DeleteAnnouncement)

	return router
}

// ============================================================================
// Webhook Routes
// ============================================================================

func (r *Router) webhookRoutes() chi.Router {
	router := chi.NewRouter()
	webhookHandler := handlers.NewWebhookHandler(r.services.Webhook, r.logger)

	// Payment webhooks
	router.Post("/stripe", webhookHandler.StripeWebhook)
	router.Post("/paypal", webhookHandler.PayPalWebhook)
	router.Post("/braintree", webhookHandler.BraintreeWebhook)

	// Communication webhooks
	router.Post("/sendgrid", webhookHandler.SendGridWebhook)
	router.Post("/twilio", webhookHandler.TwilioWebhook)
	router.Post("/slack", webhookHandler.SlackWebhook)

	// External service webhooks
	router.Post("/github", webhookHandler.GitHubWebhook)
	router.Post("/shopify", webhookHandler.ShopifyWebhook)

	return router
}

// ============================================================================
// Static Routes
// ============================================================================

func (r *Router) staticRoutes() chi.Router {
	router := chi.NewRouter()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/*", fileServer)

	// Serve uploaded files
	router.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	return router
}

// ============================================================================
// Route Groups Example
// ============================================================================

// APIVersion represents an API version with its routes
type APIVersion struct {
	Version string
	Routes  func(r chi.Router)
}

// VersionedRouter supports multiple API versions
type VersionedRouter struct {
	*chi.Mux
	versions []APIVersion
	logger   *zap.Logger
}

func NewVersionedRouter(logger *zap.Logger) *VersionedRouter {
	return &VersionedRouter{
		Mux:      chi.NewRouter(),
		versions: make([]APIVersion, 0),
		logger:   logger,
	}
}

func (r *VersionedRouter) RegisterVersion(version string, routes func(r chi.Router)) {
	r.versions = append(r.versions, APIVersion{
		Version: version,
		Routes:  routes,
	})

	// Mount versioned routes
	r.Route("/"+version, func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("X-API-Version", version)
				next.ServeHTTP(w, req)
			})
		})
		routes(r)
	})
}

// ============================================================================
// Route Documentation Examples (Swagger)
// ============================================================================

// HealthCheck godoc
// @Summary      Health check endpoint
// @Description  Returns service health status
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (r *Router) healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	// Implementation in handlers
}

// ListUsers godoc
// @Summary      List users
// @Description  Returns paginated list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page     query     int  false  "Page number"  default(1)
// @Param        limit    query     int  false  "Items per page"  default(20)
// @Param        sort     query     string  false  "Sort field"
// @Param        order    query     string  false  "Sort order (asc/desc)"
// @Success      200  {array}   models.User
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      403  {object}  handlers.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/users [get]
func listUsersHandler(w http.ResponseWriter, req *http.Request) {
	// Implementation in handlers
}

// ============================================================================
// Route Testing Helpers
// ============================================================================

// TestRouter creates a router for testing
func TestRouter() *chi.Mux {
	r := chi.NewRouter()

	// Test middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	// Test routes
	r.Get("/test", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("test ok"))
	})

	r.Get("/panic", func(w http.ResponseWriter, req *http.Request) {
		panic("test panic")
	})

	return r
}

// ============================================================================
// Route Registration Example
// ============================================================================

func main() {
	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// RESTy routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/", listUsers)   // GET /users
		r.Post("/", createUser) // POST /users
		r.Put("/", deleteUsers) // PUT /users

		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", getUser)       // GET /users/123
			r.Put("/", updateUser)    // PUT /users/123
			r.Delete("/", deleteUser) // DELETE /users/123

			r.Get("/profile", getUserProfile)    // GET /users/123/profile
			r.Put("/profile", updateUserProfile) // PUT /users/123/profile
		})
	})

	// Mount sub-routers
	r.Mount("/api/v1", apiV1Router())
	r.Mount("/admin", adminRouter())

	// Start server
	http.ListenAndServe(":8080", r)
}

// Placeholder handlers for example
func listUsers(w http.ResponseWriter, r *http.Request)         {}
func createUser(w http.ResponseWriter, r *http.Request)        {}
func deleteUsers(w http.ResponseWriter, r *http.Request)       {}
func getUser(w http.ResponseWriter, r *http.Request)           {}
func updateUser(w http.ResponseWriter, r *http.Request)        {}
func deleteUser(w http.ResponseWriter, r *http.Request)        {}
func getUserProfile(w http.ResponseWriter, r *http.Request)    {}
func updateUserProfile(w http.ResponseWriter, r *http.Request) {}
func apiV1Router() http.Handler                                { return nil }
func adminRouter() http.Handler                                { return nil }
