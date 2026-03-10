package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================================
// Custom Types
// ============================================================================

// JSONMap is a custom type for JSON fields
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// StringArray is a custom type for PostgreSQL text[] fields
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}

func (a *StringArray) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		s := string(v)
		if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
			return fmt.Errorf("invalid string array format: %s", s)
		}
		inner := s[1 : len(s)-1]
		if inner == "" {
			*a = StringArray{}
			return nil
		}
		*a = strings.Split(inner, ",")
		return nil
	case []string:
		*a = v
		return nil
	default:
		return fmt.Errorf("unsupported type for StringArray: %T", value)
	}
}

// ============================================================================
// Base Model
// ============================================================================

// BaseModel provides common fields for all models
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// ============================================================================
// User Model
// ============================================================================

// UserRole defines user roles
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleUser    UserRole = "user"
	RoleGuest   UserRole = "guest"
)

// UserStatus defines user status
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
	StatusPending   UserStatus = "pending"
)

// User represents a system user
type User struct {
	BaseModel
	Username     string     `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50,alphanum"`
	Email        string     `json:"email" gorm:"uniqueIndex;size:255;not null" validate:"required,email"`
	PasswordHash string     `json:"-" gorm:"size:255;not null"`
	FirstName    string     `json:"first_name" gorm:"size:100" validate:"required"`
	LastName     string     `json:"last_name" gorm:"size:100" validate:"required"`
	Role         UserRole   `json:"role" gorm:"type:varchar(20);default:'user'" validate:"oneof=admin manager user guest"`
	Status       UserStatus `json:"status" gorm:"type:varchar(20);default:'pending'" validate:"oneof=active inactive suspended pending"`

	// Profile fields
	AvatarURL string `json:"avatar_url,omitempty" gorm:"size:500"`
	Bio       string `json:"bio,omitempty" gorm:"type:text"`
	Phone     string `json:"phone,omitempty" gorm:"size:20"`
	Location  string `json:"location,omitempty" gorm:"size:200"`

	// Preferences
	Preferences JSONMap `json:"preferences,omitempty" gorm:"type:jsonb"`

	// Timestamps
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`

	// Relationships
	Orders      []Order          `json:"orders,omitempty" gorm:"foreignKey:UserID"`
	Products    []Product        `json:"products,omitempty" gorm:"many2many:user_products;"`
	Preferences *UserPreferences `json:"preferences,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook for User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

// AfterCreate hook for User
func (u *User) AfterCreate(tx *gorm.DB) error {
	// Create default preferences
	prefs := &UserPreferences{
		UserID:        u.ID,
		Language:      "en",
		Theme:         "light",
		Timezone:      "UTC",
		Notifications: true,
	}
	return tx.Create(prefs).Error
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

// SetPassword hashes and sets the password
func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.PasswordHash = string(hashed)
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// ============================================================================
// UserPreferences Model
// ============================================================================

// UserPreferences stores user preferences
type UserPreferences struct {
	BaseModel
	UserID        uint   `json:"user_id" gorm:"uniqueIndex;not null"`
	Language      string `json:"language" gorm:"default:'en'"`
	Theme         string `json:"theme" gorm:"default:'light'"`
	Timezone      string `json:"timezone" gorm:"default:'UTC'"`
	DateFormat    string `json:"date_format" gorm:"default:'YYYY-MM-DD'"`
	TimeFormat    string `json:"time_format" gorm:"default:'HH:mm:ss'"`
	ItemsPerPage  int    `json:"items_per_page" gorm:"default:25"`
	Notifications bool   `json:"notifications" gorm:"default:true"`
	EmailUpdates  bool   `json:"email_updates" gorm:"default:true"`
	DarkMode      bool   `json:"dark_mode" gorm:"default:false"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ============================================================================
// Product Model
// ============================================================================

// Product represents a product in the catalog
type Product struct {
	BaseModel
	Name        string  `json:"name" gorm:"size:200;not null" validate:"required"`
	Description string  `json:"description" gorm:"type:text"`
	Price       float64 `json:"price" gorm:"not null" validate:"required,gt=0"`
	Stock       int     `json:"stock" gorm:"default:0" validate:"min=0"`
	SKU         string  `json:"sku" gorm:"uniqueIndex;size:50;not null" validate:"required"`
	CategoryID  *uint   `json:"category_id,omitempty" gorm:"index"`

	// Images
	Images    StringArray `json:"images" gorm:"type:text[]"`
	Thumbnail string      `json:"thumbnail,omitempty" gorm:"size:500"`

	// Attributes
	Weight     float64 `json:"weight,omitempty"`
	WeightUnit string  `json:"weight_unit,omitempty" gorm:"size:10"`
	Dimensions string  `json:"dimensions,omitempty" gorm:"size:50"`

	// Status
	IsPublished bool       `json:"is_published" gorm:"default:false"`
	PublishedAt *time.Time `json:"published_at,omitempty"`

	// Metadata
	Metadata JSONMap `json:"metadata,omitempty" gorm:"type:jsonb"`

	// Relationships
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Variants []Variant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
	Reviews  []Review  `json:"reviews,omitempty" gorm:"foreignKey:ProductID"`
	Orders   []Order   `json:"orders,omitempty" gorm:"many2many:order_items;"`
}

// TableName specifies the table name for Product
func (Product) TableName() string {
	return "products"
}

// BeforeSave hook for Product
func (p *Product) BeforeSave(tx *gorm.DB) error {
	if p.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}

// IsAvailable checks if product is available for purchase
func (p *Product) IsAvailable() bool {
	return p.IsPublished && p.Stock > 0
}

// ReduceStock reduces product stock
func (p *Product) ReduceStock(quantity int) error {
	if p.Stock < quantity {
		return fmt.Errorf("insufficient stock: have %d, need %d", p.Stock, quantity)
	}
	p.Stock -= quantity
	return nil
}

// ============================================================================
// Category Model
// ============================================================================

// Category represents product categories
type Category struct {
	BaseModel
	Name        string `json:"name" gorm:"size:100;not null" validate:"required"`
	Slug        string `json:"slug" gorm:"uniqueIndex;size:100;not null" validate:"required"`
	Description string `json:"description" gorm:"type:text"`
	ParentID    *uint  `json:"parent_id,omitempty" gorm:"index"`
	ImageURL    string `json:"image_url,omitempty" gorm:"size:500"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`

	// Relationships
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []Product  `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}

// ============================================================================
// Variant Model
// ============================================================================

// Variant represents product variations
type Variant struct {
	BaseModel
	ProductID uint        `json:"product_id" gorm:"index;not null"`
	SKU       string      `json:"sku" gorm:"uniqueIndex;size:50;not null"`
	Name      string      `json:"name" gorm:"size:200;not null"`
	Price     float64     `json:"price" gorm:"not null" validate:"gt=0"`
	Stock     int         `json:"stock" gorm:"default:0"`
	Options   JSONMap     `json:"options" gorm:"type:jsonb"` // e.g., {"size": "M", "color": "red"}
	Images    StringArray `json:"images" gorm:"type:text[]"`

	// Relationships
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// TableName specifies the table name for Variant
func (Variant) TableName() string {
	return "variants"
}

// ============================================================================
// Order Model
// ============================================================================

// OrderStatus defines order statuses
type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderPaid       OrderStatus = "paid"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
	OrderRefunded   OrderStatus = "refunded"
)

// Order represents a customer order
type Order struct {
	BaseModel
	OrderNumber string      `json:"order_number" gorm:"uniqueIndex;size:50;not null"`
	UserID      uint        `json:"user_id" gorm:"index;not null"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`

	// Totals
	Subtotal float64 `json:"subtotal" gorm:"not null"`
	Shipping float64 `json:"shipping" gorm:"default:0"`
	Tax      float64 `json:"tax" gorm:"default:0"`
	Discount float64 `json:"discount" gorm:"default:0"`
	Total    float64 `json:"total" gorm:"not null"`

	// Addresses
	ShippingAddress JSONMap `json:"shipping_address" gorm:"type:jsonb"`
	BillingAddress  JSONMap `json:"billing_address" gorm:"type:jsonb"`

	// Payment
	PaymentMethod string     `json:"payment_method" gorm:"size:50"`
	PaymentID     string     `json:"payment_id,omitempty" gorm:"size:100"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`

	// Shipping
	ShippingMethod string     `json:"shipping_method" gorm:"size:50"`
	TrackingNumber string     `json:"tracking_number,omitempty" gorm:"size:100"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`

	// Notes
	CustomerNotes string `json:"customer_notes,omitempty" gorm:"type:text"`
	AdminNotes    string `json:"admin_notes,omitempty" gorm:"type:text"`

	// Metadata
	Metadata JSONMap `json:"metadata,omitempty" gorm:"type:jsonb"`

	// Relationships
	User  *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate hook for Order
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	// Generate order number if not provided
	if o.OrderNumber == "" {
		o.OrderNumber = fmt.Sprintf("ORD-%d-%s", time.Now().Unix(), randomString(6))
	}
	return nil
}

// CalculateTotal calculates the order total
func (o *Order) CalculateTotal() {
	o.Total = o.Subtotal + o.Shipping + o.Tax - o.Discount
}

// ============================================================================
// OrderItem Model
// ============================================================================

// OrderItem represents an item in an order
type OrderItem struct {
	BaseModel
	OrderID   uint    `json:"order_id" gorm:"index;not null"`
	ProductID uint    `json:"product_id" gorm:"index;not null"`
	VariantID *uint   `json:"variant_id,omitempty" gorm:"index"`
	SKU       string  `json:"sku" gorm:"size:50;not null"`
	Name      string  `json:"name" gorm:"size:200;not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Total     float64 `json:"total" gorm:"not null"`

	// Relationships
	Order   *Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Variant *Variant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
}

// TableName specifies the table name for OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeSave hook for OrderItem
func (oi *OrderItem) BeforeSave(tx *gorm.DB) error {
	oi.Total = float64(oi.Quantity) * oi.Price
	return nil
}

// ============================================================================
// Review Model
// ============================================================================

// Review represents a product review
type Review struct {
	BaseModel
	ProductID    uint        `json:"product_id" gorm:"index;not null"`
	UserID       uint        `json:"user_id" gorm:"index;not null"`
	OrderID      *uint       `json:"order_id,omitempty" gorm:"index"`
	Rating       int         `json:"rating" gorm:"not null" validate:"min=1,max=5"`
	Title        string      `json:"title,omitempty" gorm:"size:200"`
	Content      string      `json:"content,omitempty" gorm:"type:text"`
	Pros         StringArray `json:"pros,omitempty" gorm:"type:text[]"`
	Cons         StringArray `json:"cons,omitempty" gorm:"type:text[]"`
	Images       StringArray `json:"images,omitempty" gorm:"type:text[]"`
	IsVerified   bool        `json:"is_verified" gorm:"default:false"`
	IsApproved   bool        `json:"is_approved" gorm:"default:false"`
	HelpfulCount int         `json:"helpful_count" gorm:"default:0"`

	// Relationships
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Review
func (Review) TableName() string {
	return "reviews"
}

// ============================================================================
// Session Model
// ============================================================================

// Session represents a user session
type Session struct {
	BaseModel
	UserID       uint      `json:"user_id" gorm:"index;not null"`
	Token        string    `json:"token" gorm:"uniqueIndex;size:500;not null"`
	IPAddress    string    `json:"ip_address" gorm:"size:45"`
	UserAgent    string    `json:"user_agent" gorm:"size:500"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"index;not null"`
	LastActivity time.Time `json:"last_activity" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Session
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// ============================================================================
// AuditLog Model
// ============================================================================

// AuditAction defines audit action types
type AuditAction string

const (
	ActionCreate AuditAction = "CREATE"
	ActionUpdate AuditAction = "UPDATE"
	ActionDelete AuditAction = "DELETE"
	ActionView   AuditAction = "VIEW"
	ActionLogin  AuditAction = "LOGIN"
	ActionLogout AuditAction = "LOGOUT"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	BaseModel
	UserID       *uint       `json:"user_id,omitempty" gorm:"index"`
	Action       AuditAction `json:"action" gorm:"type:varchar(20);not null"`
	ResourceType string      `json:"resource_type" gorm:"size:100;not null"`
	ResourceID   string      `json:"resource_id" gorm:"size:100;index"`
	Changes      JSONMap     `json:"changes,omitempty" gorm:"type:jsonb"`
	OldValue     JSONMap     `json:"old_value,omitempty" gorm:"type:jsonb"`
	NewValue     JSONMap     `json:"new_value,omitempty" gorm:"type:jsonb"`
	IPAddress    string      `json:"ip_address" gorm:"size:45"`
	UserAgent    string      `json:"user_agent" gorm:"size:500"`
	RequestID    string      `json:"request_id" gorm:"size:100"`
	Duration     int64       `json:"duration_ms"`
	StatusCode   int         `json:"status_code"`
	Error        string      `json:"error,omitempty" gorm:"type:text"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// ============================================================================
// Helper Functions
// ============================================================================

// randomString generates a random string of given length
func randomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(b)
}

// ============================================================================
// Validation Helpers
// ============================================================================

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if email is valid
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// Phone validation regex (simple international format)
var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

// ValidatePhone checks if phone number is valid
func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// ============================================================================
// Model Factories
// ============================================================================

// NewUser creates a new user with default values
func NewUser(username, email, password, firstName, lastName string) (*User, error) {
	user := &User{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      RoleUser,
		Status:    StatusPending,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

// NewProduct creates a new product
func NewProduct(name, sku string, price float64) *Product {
	return &Product{
		Name:        name,
		SKU:         sku,
		Price:       price,
		Stock:       0,
		IsPublished: false,
		Metadata:    make(JSONMap),
	}
}

// NewOrder creates a new order
func NewOrder(userID uint) *Order {
	return &Order{
		UserID:   userID,
		Status:   OrderPending,
		Items:    make([]OrderItem, 0),
		Metadata: make(JSONMap),
	}
}
