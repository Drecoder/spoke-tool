package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================================
// Request/Response Types
// ============================================================================

// APIResponse standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *APIMeta    `json:"meta,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// APIMeta contains pagination metadata
type APIMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ============================================================================
// User Handlers
// ============================================================================

// UserHandler handles user-related requests
type UserHandler struct {
	db       *gorm.DB
	logger   *zap.Logger
	validate *validator.Validate
}

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

// UserResponse represents user data response
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *gorm.DB, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		db:       db,
		logger:   logger,
		validate: validator.New(),
	}
}

// CreateUser handles user creation
// @Summary Create a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation request"
// @Success 201 {object} UserResponse
// @Failure 400 {object} APIError
// @Failure 409 {object} APIError
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: err.Error(),
		})
		return
	}

	// Check if user exists
	var existingUser User
	if err := h.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, APIError{
			Code:    "USER_EXISTS",
			Message: "User with this email or username already exists",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		h.logger.Error("Database error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to check user existence",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to process password",
		})
		return
	}

	// Create user
	user := User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Active:       true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create user",
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// GetUser handles retrieving a user by ID
// @Summary Get user by ID
// @Description Retrieve user details
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404 {object} APIError
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_ID",
			Message: "Invalid user ID",
		})
		return
	}

	var user User
	if err := h.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIError{
				Code:    "USER_NOT_FOUND",
				Message: fmt.Sprintf("User with ID %d not found", id),
			})
			return
		}
		h.logger.Error("Database error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to retrieve user",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// ListUsersRequest represents list users request parameters
type ListUsersRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PerPage  int    `form:"per_page,default=20" binding:"min=1,max=100"`
	Search   string `form:"search"`
	SortBy   string `form:"sort_by" binding:"oneof=id username email created_at"`
	SortDesc bool   `form:"sort_desc"`
}

// ListUsers handles listing users with pagination
// @Summary List users
// @Description Retrieve paginated list of users
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort field" Enums(id,username,email,created_at)
// @Param sort_desc query bool false "Sort descending"
// @Success 200 {object} []UserResponse
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_PARAMETERS",
			Message: "Invalid query parameters",
		})
		return
	}

	// Build query
	query := h.db.Model(&User{})

	// Apply search
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ? OR full_name ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		h.logger.Error("Failed to count users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to count users",
		})
		return
	}

	// Apply sorting
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortDesc {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("id DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PerPage
	var users []User
	if err := query.Offset(offset).Limit(req.PerPage).Find(&users).Error; err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list users",
		})
		return
	}

	// Build response
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	// Calculate total pages
	totalPages := (int(total) + req.PerPage - 1) / req.PerPage

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    responses,
		Meta: &APIMeta{
			Page:       req.Page,
			PerPage:    req.PerPage,
			Total:      int(total),
			TotalPages: totalPages,
		},
	})
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email" validate:"omitempty,email"`
}

// UpdateUser handles updating a user
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Update request"
// @Success 200 {object} UserResponse
// @Failure 400 {object} APIError
// @Failure 404 {object} APIError
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_ID",
			Message: "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	// Find user
	var user User
	if err := h.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIError{
				Code:    "USER_NOT_FOUND",
				Message: fmt.Sprintf("User with ID %d not found", id),
			})
			return
		}
		h.logger.Error("Database error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to retrieve user",
		})
		return
	}

	// Update fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		// Check if email is already taken
		var existingUser User
		if err := h.db.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, APIError{
				Code:    "EMAIL_EXISTS",
				Message: "Email already in use",
			})
			return
		}
		user.Email = req.Email
	}

	user.UpdatedAt = time.Now()

	if err := h.db.Save(&user).Error; err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Description Delete a user
// @Tags users
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} APIError
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_ID",
			Message: "Invalid user ID",
		})
		return
	}

	result := h.db.Delete(&User{}, id)
	if result.Error != nil {
		h.logger.Error("Failed to delete user", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to delete user",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, APIError{
			Code:    "USER_NOT_FOUND",
			Message: fmt.Sprintf("User with ID %d not found", id),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ============================================================================
// Auth Handlers
// ============================================================================

// AuthHandler handles authentication requests
type AuthHandler struct {
	db     *gorm.DB
	logger *zap.Logger
	jwtKey []byte
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int64        `json:"expires_in"`
	User      UserResponse `json:"user"`
}

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB, logger *zap.Logger, jwtKey string) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
		jwtKey: []byte(jwtKey),
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} APIError
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	// Find user
	var user User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Invalid email or password",
			})
			return
		}
		h.logger.Error("Database error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Authentication failed",
		})
		return
	}

	// Check if user is active
	if !user.Active {
		c.JSON(http.StatusUnauthorized, APIError{
			Code:    "ACCOUNT_INACTIVE",
			Message: "Account is inactive",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, APIError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-server",
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtKey)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to generate token",
		})
		return
	}

	// Update last login
	h.db.Model(&user).Update("last_login_at", time.Now())

	c.JSON(http.StatusOK, LoginResponse{
		Token:     tokenString,
		ExpiresIn: expirationTime.Unix(),
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// RefreshToken handles token refresh
// @Summary Refresh JWT token
// @Description Get a new JWT token using valid token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} LoginResponse
// @Failure 401 {object} APIError
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIError{
			Code:    "UNAUTHORIZED",
			Message: "Unauthorized",
		})
		return
	}

	// Find user
	var user User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, APIError{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	// Generate new token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-server",
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtKey)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:     tokenString,
		ExpiresIn: expirationTime.Unix(),
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// Logout handles user logout
// @Summary User logout
// @Description Invalidate user session
// @Tags auth
// @Security BearerAuth
// @Success 200
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real implementation, you might blacklist the token
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ============================================================================
// Health Handler
// ============================================================================

// HealthHandler handles health checks
type HealthHandler struct {
	db        *gorm.DB
	logger    *zap.Logger
	startTime time.Time
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		db:        db,
		logger:    logger,
		startTime: time.Now(),
	}
}

// Health handles health check requests
// @Summary Health check
// @Description Check service health
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
		Version:   "1.0.0",
		Services:  make(map[string]string),
	}

	// Check database
	if err := h.db.Exec("SELECT 1").Error; err != nil {
		response.Status = "degraded"
		response.Services["database"] = "down"
		h.logger.Error("Database health check failed", zap.Error(err))
	} else {
		response.Services["database"] = "up"
	}

	c.JSON(http.StatusOK, response)
}

// Readiness handles readiness checks
// @Summary Readiness check
// @Description Check if service is ready to accept traffic
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @ServiceUnavailable 503
// @Router /ready [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	// Check database connectivity
	if err := h.db.Exec("SELECT 1").Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"reason":  "database connection failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now(),
	})
}

// ============================================================================
// Product Handlers (Example)
// ============================================================================

// ProductHandler handles product-related requests
type ProductHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// Product represents a product model
type Product struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	Stock       int       `json:"stock" binding:"min=0"`
	SKU         string    `json:"sku" binding:"required" gorm:"uniqueIndex"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewProductHandler creates a new product handler
func NewProductHandler(db *gorm.DB, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		db:     db,
		logger: logger,
	}
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid product data",
			Details: err.Error(),
		})
		return
	}

	// Check if SKU already exists
	var existingProduct Product
	if err := h.db.Where("sku = ?", product.SKU).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusConflict, APIError{
			Code:    "SKU_EXISTS",
			Message: "Product with this SKU already exists",
		})
		return
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if err := h.db.Create(&product).Error; err != nil {
		h.logger.Error("Failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct handles retrieving a product
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_ID",
			Message: "Invalid product ID",
		})
		return
	}

	var product Product
	if err := h.db.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIError{
				Code:    "PRODUCT_NOT_FOUND",
				Message: fmt.Sprintf("Product with ID %d not found", id),
			})
			return
		}
		h.logger.Error("Database error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to retrieve product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListProducts handles listing products with filters
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var products []Product
	query := h.db.Model(&Product{})

	// Apply filters
	if category := c.Query("category"); category != "" {
		query = query.Where("category_id = ?", category)
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	if inStock := c.Query("in_stock"); inStock == "true" {
		query = query.Where("stock > 0")
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	offset := (page - 1) * perPage

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(perPage).Find(&products).Error; err != nil {
		h.logger.Error("Failed to list products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list products",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    products,
		Meta: &APIMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      int(total),
			TotalPages: (int(total) + perPage - 1) / perPage,
		},
	})
}

// UpdateProduct handles product updates
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_ID",
			Message: "Invalid product ID",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_REQUEST",
			Message: "Invalid update data",
		})
		return
	}

	// Remove fields that shouldn't be updated
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "sku") // SKU shouldn't be updated

	updates["updated_at"] = time.Now()

	result := h.db.Model(&Product{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		// Check for duplicate SKU error
		if strings.Contains(result.Error.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, APIError{
				Code:    "SKU_EXISTS",
				Message: "SKU already in use",
			})
			return
		}
		h.logger.Error("Failed to update product", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update product",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, APIError{
			Code:    "PRODUCT_NOT_FOUND",
			Message: fmt.Sprintf("Product with ID %d not found", id),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct handles product deletion
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Code:    "INVALID_ID",
			Message: "Invalid product ID",
		})
		return
	}

	result := h.db.Delete(&Product{}, id)
	if result.Error != nil {
		h.logger.Error("Failed to delete product", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to delete product",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, APIError{
			Code:    "PRODUCT_NOT_FOUND",
			Message: fmt.Sprintf("Product with ID %d not found", id),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
