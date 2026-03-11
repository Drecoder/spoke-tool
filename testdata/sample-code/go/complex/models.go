package models

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================================
// Custom Types
// ============================================================================

// Email represents a validated email address
type Email string

func (e Email) Validate() error {
	if !emailRegex.MatchString(string(e)) {
		return ErrInvalidEmail
	}
	return nil
}

func (e Email) Domain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// PhoneNumber represents a validated phone number
type PhoneNumber string

func (p PhoneNumber) Validate() error {
	if !phoneRegex.MatchString(string(p)) {
		return ErrInvalidPhone
	}
	return nil
}

var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

// Currency represents a monetary amount with currency code
type Currency struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency" validate:"len=3"`
}

func (c Currency) String() string {
	return fmt.Sprintf("%.2f %s", c.Amount, c.Currency)
}

func (c Currency) Add(other Currency) (Currency, error) {
	if c.Currency != other.Currency {
		return Currency{}, ErrCurrencyMismatch
	}
	return Currency{
		Amount:   c.Amount + other.Amount,
		Currency: c.Currency,
	}, nil
}

// JSONB wrapper for PostgreSQL JSONB type
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, &j)
}

// StringArray for PostgreSQL text[] type
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
	default:
		return fmt.Errorf("unsupported type for StringArray: %T", value)
	}
}

// ============================================================================
// Base Models
// ============================================================================

// BaseModel provides common fields for all models
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Auditable adds audit fields
type Auditable struct {
	CreatedBy uint `json:"created_by" gorm:"index"`
	UpdatedBy uint `json:"updated_by" gorm:"index"`
}

// Versioned adds optimistic locking
type Versioned struct {
	Version int `json:"version" gorm:"default:1"`
}

// ============================================================================
// User Models
// ============================================================================

// UserRole enum
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleUser    UserRole = "user"
	RoleGuest   UserRole = "guest"
)

// UserStatus enum
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
	Auditable
	Versioned

	Email        Email       `json:"email" gorm:"uniqueIndex;size:255" validate:"required,email"`
	Phone        PhoneNumber `json:"phone" gorm:"uniqueIndex;size:20" validate:"omitempty"`
	Username     string      `json:"username" gorm:"uniqueIndex;size:50" validate:"required,min=3,max=50,alphanum"`
	PasswordHash string      `json:"-" gorm:"size:255"`
	FirstName    string      `json:"first_name" gorm:"size:100" validate:"required"`
	LastName     string      `json:"last_name" gorm:"size:100" validate:"required"`
	Role         UserRole    `json:"role" gorm:"type:varchar(20);default:'user'" validate:"oneof=admin manager user guest"`
	Status       UserStatus  `json:"status" gorm:"type:varchar(20);default:'pending'" validate:"oneof=active inactive suspended pending"`

	// Profile fields
	AvatarURL string `json:"avatar_url,omitempty" gorm:"size:500"`
	Bio       string `json:"bio,omitempty" gorm:"type:text"`
	Location  string `json:"location,omitempty" gorm:"size:200"`
	Website   string `json:"website,omitempty" gorm:"size:500"`

	// Preferences
	Preferences JSONB `json:"preferences,omitempty" gorm:"type:jsonb"`

	// Metadata
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at,omitempty"`

	// Relationships
	Profile  *UserProfile  `json:"profile,omitempty" gorm:"foreignKey:UserID"`
	Accounts []Account     `json:"accounts,omitempty" gorm:"foreignKey:UserID"`
	Settings *UserSettings `json:"settings,omitempty" gorm:"foreignKey:UserID"`
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) Can(permission string) bool {
	// Check role-based permissions
	switch permission {
	case "users:create", "users:update", "users:delete":
		return u.Role == RoleAdmin || u.Role == RoleManager
	case "users:view":
		return u.Role == RoleAdmin || u.Role == RoleManager || u.Role == RoleUser
	default:
		return u.Role == RoleAdmin
	}
}

// BeforeSave GORM hook
func (u *User) BeforeSave(tx *gorm.DB) error {
	if err := u.Email.Validate(); err != nil {
		return err
	}
	if u.Phone != "" {
		if err := u.Phone.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// UserProfile extends user with profile information
type UserProfile struct {
	BaseModel
	UserID     uint       `json:"user_id" gorm:"uniqueIndex"`
	Company    string     `json:"company,omitempty" gorm:"size:200"`
	JobTitle   string     `json:"job_title,omitempty" gorm:"size:200"`
	Department string     `json:"department,omitempty" gorm:"size:200"`
	Address    string     `json:"address,omitempty" gorm:"type:text"`
	City       string     `json:"city,omitempty" gorm:"size:100"`
	State      string     `json:"state,omitempty" gorm:"size:100"`
	Country    string     `json:"country,omitempty" gorm:"size:100"`
	PostalCode string     `json:"postal_code,omitempty" gorm:"size:20"`
	BirthDate  *time.Time `json:"birth_date,omitempty"`
	Gender     string     `json:"gender,omitempty" gorm:"size:20"`

	// Social links
	LinkedIn string `json:"linkedin,omitempty" gorm:"size:500"`
	Twitter  string `json:"twitter,omitempty" gorm:"size:500"`
	GitHub   string `json:"github,omitempty" gorm:"size:500"`
}

// UserSettings stores user preferences
type UserSettings struct {
	BaseModel
	UserID             uint   `json:"user_id" gorm:"uniqueIndex"`
	Language           string `json:"language" gorm:"default:'en'"`
	Theme              string `json:"theme" gorm:"default:'light'"`
	Timezone           string `json:"timezone" gorm:"default:'UTC'"`
	DateFormat         string `json:"date_format" gorm:"default:'YYYY-MM-DD'"`
	TimeFormat         string `json:"time_format" gorm:"default:'HH:mm:ss'"`
	ItemsPerPage       int    `json:"items_per_page" gorm:"default:25"`
	EmailNotifications bool   `json:"email_notifications" gorm:"default:true"`
	SMSNotifications   bool   `json:"sms_notifications" gorm:"default:false"`
	PushNotifications  bool   `json:"push_notifications" gorm:"default:true"`
}

// ============================================================================
// Account Models
// ============================================================================

// AccountType enum
type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeCredit     AccountType = "credit"
	AccountTypeInvestment AccountType = "investment"
)

// Account represents a financial account
type Account struct {
	BaseModel
	Auditable
	Versioned

	UserID           uint        `json:"user_id" gorm:"index;not null"`
	AccountNumber    string      `json:"account_number" gorm:"uniqueIndex;size:50"`
	AccountType      AccountType `json:"account_type" gorm:"type:varchar(20);not null"`
	Name             string      `json:"name" gorm:"size:200;not null"`
	Description      string      `json:"description,omitempty" gorm:"type:text"`
	Currency         string      `json:"currency" gorm:"size:3;default:'USD'"`
	Balance          float64     `json:"balance" gorm:"default:0"`
	AvailableBalance float64     `json:"available_balance" gorm:"default:0"`
	InterestRate     float64     `json:"interest_rate,omitempty" gorm:"default:0"`
	Status           string      `json:"status" gorm:"type:varchar(20);default:'active'"`
	ClosedAt         *time.Time  `json:"closed_at,omitempty"`

	// Limits
	OverdraftLimit float64 `json:"overdraft_limit,omitempty" gorm:"default:0"`
	DailyLimit     float64 `json:"daily_limit,omitempty" gorm:"default:0"`
	MonthlyLimit   float64 `json:"monthly_limit,omitempty" gorm:"default:0"`

	// Metadata
	Metadata JSONB `json:"metadata,omitempty" gorm:"type:jsonb"`

	// Relationships
	User         *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:AccountID"`
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.Balance += amount
	a.AvailableBalance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if a.AvailableBalance < amount {
		return ErrInsufficientFunds
	}
	a.Balance -= amount
	a.AvailableBalance -= amount
	return nil
}

// ============================================================================
// Transaction Models
// ============================================================================

// TransactionType enum
type TransactionType string

const (
	TxTypeDeposit    TransactionType = "deposit"
	TxTypeWithdrawal TransactionType = "withdrawal"
	TxTypeTransfer   TransactionType = "transfer"
	TxTypePayment    TransactionType = "payment"
	TxTypeFee        TransactionType = "fee"
	TxTypeInterest   TransactionType = "interest"
	TxTypeRefund     TransactionType = "refund"
)

// TransactionStatus enum
type TransactionStatus string

const (
	TxStatusPending   TransactionStatus = "pending"
	TxStatusCompleted TransactionStatus = "completed"
	TxStatusFailed    TransactionStatus = "failed"
	TxStatusCancelled TransactionStatus = "cancelled"
	TxStatusReversed  TransactionStatus = "reversed"
)

// Transaction represents a financial transaction
type Transaction struct {
	BaseModel
	Auditable
	Versioned

	AccountID     uint              `json:"account_id" gorm:"index;not null"`
	TransactionID string            `json:"transaction_id" gorm:"uniqueIndex;size:100"`
	ReferenceID   string            `json:"reference_id,omitempty" gorm:"index;size:100"`
	Type          TransactionType   `json:"type" gorm:"type:varchar(20);not null"`
	Status        TransactionStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Amount        float64           `json:"amount" gorm:"not null"`
	Currency      string            `json:"currency" gorm:"size:3;not null"`
	BalanceBefore float64           `json:"balance_before"`
	BalanceAfter  float64           `json:"balance_after"`
	Description   string            `json:"description,omitempty" gorm:"type:text"`

	// Counterparty
	CounterpartyID      string `json:"counterparty_id,omitempty" gorm:"index;size:100"`
	CounterpartyName    string `json:"counterparty_name,omitempty" gorm:"size:200"`
	CounterpartyAccount string `json:"counterparty_account,omitempty" gorm:"size:100"`

	// Fee breakdown
	Fee       float64 `json:"fee,omitempty" gorm:"default:0"`
	Tax       float64 `json:"tax,omitempty" gorm:"default:0"`
	NetAmount float64 `json:"net_amount"`

	// Timestamps
	EffectiveDate *time.Time `json:"effective_date,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	FailedAt      *time.Time `json:"failed_at,omitempty"`

	// Metadata
	Metadata JSONB `json:"metadata,omitempty" gorm:"type:jsonb"`

	// Relationships
	Account *Account `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	// Generate transaction ID if not provided
	if t.TransactionID == "" {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%d-%d", time.Now().UnixNano(), t.AccountID)))
		t.TransactionID = base64.URLEncoding.EncodeToString(hash[:16])
	}

	// Calculate net amount
	t.NetAmount = t.Amount - t.Fee - t.Tax

	return nil
}

// ============================================================================
// Product Models
// ============================================================================

// ProductCategory represents product categories
type ProductCategory struct {
	BaseModel
	Name        string `json:"name" gorm:"uniqueIndex;size:100;not null"`
	Slug        string `json:"slug" gorm:"uniqueIndex;size:100;not null"`
	Description string `json:"description,omitempty" gorm:"type:text"`
	ParentID    *uint  `json:"parent_id,omitempty" gorm:"index"`
	ImageURL    string `json:"image_url,omitempty" gorm:"size:500"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`

	// Relationships
	Parent   *ProductCategory  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []ProductCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []Product         `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

// Product represents a product in the catalog
type Product struct {
	BaseModel
	Auditable
	Versioned

	SKU         string   `json:"sku" gorm:"uniqueIndex;size:50;not null"`
	Name        string   `json:"name" gorm:"size:200;not null"`
	Slug        string   `json:"slug" gorm:"uniqueIndex;size:200;not null"`
	Description string   `json:"description" gorm:"type:text"`
	ShortDesc   string   `json:"short_description,omitempty" gorm:"size:500"`
	CategoryID  *uint    `json:"category_id,omitempty" gorm:"index"`
	Price       Currency `json:"price" gorm:"embedded;embeddedPrefix:price_"`
	Cost        Currency `json:"cost,omitempty" gorm:"embedded;embeddedPrefix:cost_"`
	Stock       int      `json:"stock" gorm:"default:0"`
	Reserved    int      `json:"reserved" gorm:"default:0"`
	Available   int      `json:"available" gorm:"-:all"`

	// Attributes
	Weight     float64 `json:"weight,omitempty"`
	WeightUnit string  `json:"weight_unit,omitempty" gorm:"size:10"`
	Dimensions string  `json:"dimensions,omitempty" gorm:"size:50"`

	// Media
	Images    StringArray `json:"images,omitempty" gorm:"type:text[]"`
	Thumbnail string      `json:"thumbnail,omitempty" gorm:"size:500"`

	// Status
	Status      string     `json:"status" gorm:"type:varchar(20);default:'draft'"`
	PublishedAt *time.Time `json:"published_at,omitempty"`

	// SEO
	MetaTitle    string      `json:"meta_title,omitempty" gorm:"size:200"`
	MetaDesc     string      `json:"meta_description,omitempty" gorm:"size:500"`
	MetaKeywords StringArray `json:"meta_keywords,omitempty" gorm:"type:text[]"`

	// Metadata
	Attributes JSONB `json:"attributes,omitempty" gorm:"type:jsonb"`

	// Relationships
	Category *ProductCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Variants []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
	Reviews  []Review         `json:"reviews,omitempty" gorm:"foreignKey:ProductID"`
}

func (p *Product) AfterFind(tx *gorm.DB) error {
	p.Available = p.Stock - p.Reserved
	return nil
}

func (p *Product) IsAvailable() bool {
	return p.Status == "published" && p.Available > 0
}

// ProductVariant represents product variations
type ProductVariant struct {
	BaseModel
	ProductID uint        `json:"product_id" gorm:"index;not null"`
	SKU       string      `json:"sku" gorm:"uniqueIndex;size:50;not null"`
	Name      string      `json:"name" gorm:"size:200;not null"`
	Options   JSONB       `json:"options" gorm:"type:jsonb"` // e.g., {"size": "M", "color": "red"}
	Price     Currency    `json:"price" gorm:"embedded;embeddedPrefix:price_"`
	Stock     int         `json:"stock" gorm:"default:0"`
	Images    StringArray `json:"images,omitempty" gorm:"type:text[]"`

	// Relationships
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// ============================================================================
// Order Models
// ============================================================================

// OrderStatus enum
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// Order represents a customer order
type Order struct {
	BaseModel
	Auditable
	Versioned

	OrderNumber string      `json:"order_number" gorm:"uniqueIndex;size:50;not null"`
	UserID      uint        `json:"user_id" gorm:"index;not null"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`

	// Totals
	Subtotal Currency `json:"subtotal" gorm:"embedded;embeddedPrefix:subtotal_"`
	Shipping Currency `json:"shipping" gorm:"embedded;embeddedPrefix:shipping_"`
	Tax      Currency `json:"tax" gorm:"embedded;embeddedPrefix:tax_"`
	Discount Currency `json:"discount" gorm:"embedded;embeddedPrefix:discount_"`
	Total    Currency `json:"total" gorm:"embedded;embeddedPrefix:total_"`

	// Addresses
	ShippingAddress JSONB `json:"shipping_address" gorm:"type:jsonb"`
	BillingAddress  JSONB `json:"billing_address" gorm:"type:jsonb"`

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
	Metadata JSONB `json:"metadata,omitempty" gorm:"type:jsonb"`

	// Relationships
	User     *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items    []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Payments []Payment   `json:"payments,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	BaseModel
	OrderID   uint     `json:"order_id" gorm:"index;not null"`
	ProductID uint     `json:"product_id" gorm:"index;not null"`
	VariantID *uint    `json:"variant_id,omitempty" gorm:"index"`
	SKU       string   `json:"sku" gorm:"size:50;not null"`
	Name      string   `json:"name" gorm:"size:200;not null"`
	Quantity  int      `json:"quantity" gorm:"not null"`
	Price     Currency `json:"price" gorm:"embedded;embeddedPrefix:price_"`
	Total     Currency `json:"total" gorm:"embedded;embeddedPrefix:total_"`

	// Relationships
	Order   *Order          `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Product *Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Variant *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
}

// ============================================================================
// Payment Models
// ============================================================================

// PaymentMethod enum
type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodBank   PaymentMethod = "bank_transfer"
	PaymentMethodCash   PaymentMethod = "cash"
	PaymentMethodCrypto PaymentMethod = "cryptocurrency"
	PaymentMethodWallet PaymentMethod = "digital_wallet"
)

// Payment represents a payment transaction
type Payment struct {
	BaseModel
	Auditable

	PaymentID string        `json:"payment_id" gorm:"uniqueIndex;size:100;not null"`
	OrderID   uint          `json:"order_id" gorm:"index;not null"`
	UserID    uint          `json:"user_id" gorm:"index;not null"`
	Method    PaymentMethod `json:"method" gorm:"type:varchar(50);not null"`
	Amount    Currency      `json:"amount" gorm:"embedded;embeddedPrefix:amount_"`
	Status    string        `json:"status" gorm:"type:varchar(20);default:'pending'"`

	// Payment details
	Provider   string `json:"provider" gorm:"size:100"`
	ProviderID string `json:"provider_id" gorm:"size:100"`
	LastFour   string `json:"last_four,omitempty" gorm:"size:4"`
	CardBrand  string `json:"card_brand,omitempty" gorm:"size:50"`

	// Timestamps
	AuthorizedAt *time.Time `json:"authorized_at,omitempty"`
	SettledAt    *time.Time `json:"settled_at,omitempty"`
	RefundedAt   *time.Time `json:"refunded_at,omitempty"`
	FailedAt     *time.Time `json:"failed_at,omitempty"`

	// Error
	FailureCode    string `json:"failure_code,omitempty" gorm:"size:50"`
	FailureMessage string `json:"failure_message,omitempty" gorm:"type:text"`

	// Metadata
	Metadata JSONB `json:"metadata,omitempty" gorm:"type:jsonb"`

	// Relationships
	Order   *Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Refunds []Refund `json:"refunds,omitempty" gorm:"foreignKey:PaymentID"`
}

// Refund represents a payment refund
type Refund struct {
	BaseModel
	Auditable

	RefundID    string     `json:"refund_id" gorm:"uniqueIndex;size:100;not null"`
	PaymentID   uint       `json:"payment_id" gorm:"index;not null"`
	Amount      Currency   `json:"amount" gorm:"embedded;embeddedPrefix:amount_"`
	Reason      string     `json:"reason,omitempty" gorm:"type:text"`
	Status      string     `json:"status" gorm:"type:varchar(20);default:'pending'"`
	ProviderID  string     `json:"provider_id,omitempty" gorm:"size:100"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Relationships
	Payment *Payment `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
}

// ============================================================================
// Review Models
// ============================================================================

// Review represents a product review
type Review struct {
	BaseModel
	Auditable

	ProductID    uint        `json:"product_id" gorm:"index;not null"`
	UserID       uint        `json:"user_id" gorm:"index;not null"`
	OrderID      *uint       `json:"order_id,omitempty" gorm:"index"`
	Rating       int         `json:"rating" gorm:"not null" validate:"min=1,max=5"`
	Title        string      `json:"title,omitempty" gorm:"size:200"`
	Content      string      `json:"content,omitempty" gorm:"type:text"`
	Pros         StringArray `json:"pros,omitempty" gorm:"type:text[]"`
	Cons         StringArray `json:"cons,omitempty" gorm:"type:text[]"`
	Images       StringArray `json:"images,omitempty" gorm:"type:text[]"`
	Verified     bool        `json:"verified" gorm:"default:false"`
	Status       string      `json:"status" gorm:"type:varchar(20);default:'pending'"`
	HelpfulCount int         `json:"helpful_count" gorm:"default:0"`

	// Relationships
	Product   *Product         `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	User      *User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Order     *Order           `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Responses []ReviewResponse `json:"responses,omitempty" gorm:"foreignKey:ReviewID"`
	Votes     []ReviewVote     `json:"votes,omitempty" gorm:"foreignKey:ReviewID"`
}

// ReviewResponse represents a response to a review
type ReviewResponse struct {
	BaseModel
	Auditable

	ReviewID uint   `json:"review_id" gorm:"index;not null"`
	UserID   uint   `json:"user_id" gorm:"index;not null"`
	Content  string `json:"content" gorm:"type:text;not null"`

	// Relationships
	Review *Review `json:"review,omitempty" gorm:"foreignKey:ReviewID"`
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ReviewVote represents a vote on a review
type ReviewVote struct {
	BaseModel
	ReviewID uint `json:"review_id" gorm:"index;not null"`
	UserID   uint `json:"user_id" gorm:"index;not null"`
	Vote     bool `json:"vote"` // true = helpful, false = not helpful

	// Relationships
	Review *Review `json:"review,omitempty" gorm:"foreignKey:ReviewID"`
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ============================================================================
// Notification Models
// ============================================================================

// NotificationType enum
type NotificationType string

const (
	NotifTypeEmail NotificationType = "email"
	NotifTypeSMS   NotificationType = "sms"
	NotifTypePush  NotificationType = "push"
	NotifTypeInApp NotificationType = "in_app"
)

// Notification represents a user notification
type Notification struct {
	BaseModel

	UserID      uint             `json:"user_id" gorm:"index;not null"`
	Type        NotificationType `json:"type" gorm:"type:varchar(20);not null"`
	Title       string           `json:"title" gorm:"size:200;not null"`
	Content     string           `json:"content" gorm:"type:text"`
	Data        JSONB            `json:"data,omitempty" gorm:"type:jsonb"`
	ReadAt      *time.Time       `json:"read_at,omitempty" gorm:"index"`
	SentAt      *time.Time       `json:"sent_at,omitempty"`
	DeliveredAt *time.Time       `json:"delivered_at,omitempty"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ============================================================================
// Audit Log Models
// ============================================================================

// AuditAction enum
type AuditAction string

const (
	AuditCreate AuditAction = "CREATE"
	AuditUpdate AuditAction = "UPDATE"
	AuditDelete AuditAction = "DELETE"
	AuditView   AuditAction = "VIEW"
	AuditLogin  AuditAction = "LOGIN"
	AuditLogout AuditAction = "LOGOUT"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	BaseModel

	UserID       *uint       `json:"user_id,omitempty" gorm:"index"`
	Action       AuditAction `json:"action" gorm:"type:varchar(20);not null"`
	ResourceType string      `json:"resource_type" gorm:"size:100;not null"`
	ResourceID   string      `json:"resource_id" gorm:"size:100;index"`
	Changes      JSONB       `json:"changes,omitempty" gorm:"type:jsonb"`
	OldValue     JSONB       `json:"old_value,omitempty" gorm:"type:jsonb"`
	NewValue     JSONB       `json:"new_value,omitempty" gorm:"type:jsonb"`
	IPAddress    string      `json:"ip_address" gorm:"size:45"`
	UserAgent    string      `json:"user_agent" gorm:"size:500"`
	RequestID    string      `json:"request_id" gorm:"size:100"`
	Duration     int64       `json:"duration_ms"`
	Status       int         `json:"status_code"`
	Error        string      `json:"error,omitempty" gorm:"type:text"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ============================================================================
// Error Definitions
// ============================================================================

var (
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrInvalidPhone      = errors.New("invalid phone number")
	ErrCurrencyMismatch  = errors.New("currency mismatch")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrProductNotFound   = errors.New("product not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrOrderNotFound     = errors.New("order not found")
)

// ============================================================================
// Validation Helpers
// ============================================================================

var validate = validator.New()

func init() {
	// Register custom validators
	validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		return emailRegex.MatchString(email)
	})

	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		return phoneRegex.MatchString(phone)
	})
}

func (u *User) Validate() error {
	return validate.Struct(u)
}

func (p *Product) Validate() error {
	return validate.Struct(p)
}

// ============================================================================
// DTOs (Data Transfer Objects)
// ============================================================================

// CreateUserRequest DTO
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required,min=3,max=50,alphanum"`
}

// LoginRequest DTO
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse DTO
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// ProductResponse DTO
type ProductResponse struct {
	ID          uint    `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func (p *Product) ToResponse() *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price.Amount,
		Stock:       p.Stock,
	}
}
