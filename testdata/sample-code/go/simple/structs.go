package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// ============================================================================
// Basic Struct
// ============================================================================

// Person represents a basic person struct
type Person struct {
	Name string
	Age  int
}

// Basic struct usage
func basicStruct() {
	// Different ways to create structs
	p1 := Person{"Alice", 30}
	p2 := Person{Name: "Bob", Age: 25}
	p3 := Person{Name: "Charlie"} // Age defaults to 0

	fmt.Printf("p1: %+v\n", p1)
	fmt.Printf("p2: %+v\n", p2)
	fmt.Printf("p3: %+v\n", p3)
}

// ============================================================================
// Struct with Methods
// ============================================================================

// Employee represents an employee with methods
type Employee struct {
	ID       int
	Name     string
	Salary   float64
	HireDate time.Time
}

// Value receiver method
func (e Employee) YearsOfService() float64 {
	return time.Since(e.HireDate).Hours() / 24 / 365
}

// Pointer receiver method (can modify the struct)
func (e *Employee) GiveRaise(percent float64) {
	e.Salary = e.Salary * (1 + percent/100)
}

// Stringer implementation
func (e Employee) String() string {
	return fmt.Sprintf("Employee #%d: %s ($%.2f)", e.ID, e.Name, e.Salary)
}

// ============================================================================
// Struct with Tags
// ============================================================================

// User represents a user with JSON tags
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"-"` // This field will be ignored in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user with default values
func NewUser(username, email string) *User {
	now := time.Now()
	return &User{
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ============================================================================
// Nested Structs
// ============================================================================

// Address represents a physical address
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

// Contact represents contact information
type Contact struct {
	Email     string
	Phone     string
	Preferred string // "email" or "phone"
}

// Customer represents a customer with nested structs
type Customer struct {
	ID        int
	Name      string
	Address   Address
	Contact   Contact
	CreatedAt time.Time
}

// FullAddress returns the complete address as a string
func (c Customer) FullAddress() string {
	return fmt.Sprintf("%s, %s, %s %s, %s",
		c.Address.Street,
		c.Address.City,
		c.Address.State,
		c.Address.ZipCode,
		c.Address.Country,
	)
}

// ContactInfo returns the preferred contact method
func (c Customer) ContactInfo() string {
	if c.Contact.Preferred == "email" {
		return c.Contact.Email
	}
	return c.Contact.Phone
}

// ============================================================================
// Embedded Structs (Composition)
// ============================================================================

// Animal is a base struct
type Animal struct {
	Name   string
	Age    int
	Weight float64
}

// Sounder interface for animals that make sounds
type Sounder interface {
	MakeSound() string
}

// Dog "inherits" from Animal through embedding
type Dog struct {
	Animal
	Breed string
}

func (d Dog) MakeSound() string {
	return "Woof!"
}

// Cat "inherits" from Animal through embedding
type Cat struct {
	Animal
	Indoor bool
}

func (c Cat) MakeSound() string {
	return "Meow!"
}

// Bird "inherits" from Animal through embedding
type Bird struct {
	Animal
	Wingspan float64
}

func (b Bird) MakeSound() string {
	return "Tweet!"
}

// ============================================================================
// Struct with Private Fields
// ============================================================================

// BankAccount has private fields with accessor methods
type BankAccount struct {
	accountNumber string
	balance       float64
	owner         string
}

// NewBankAccount creates a new account
func NewBankAccount(owner string) *BankAccount {
	return &BankAccount{
		accountNumber: generateAccountNumber(),
		balance:       0,
		owner:         owner,
	}
}

// generateAccountNumber is a helper function
func generateAccountNumber() string {
	return fmt.Sprintf("ACC-%d", time.Now().UnixNano())
}

// Deposit adds money to the account
func (b *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}
	b.balance += amount
	return nil
}

// Withdraw removes money from the account
func (b *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive")
	}
	if b.balance < amount {
		return fmt.Errorf("insufficient funds")
	}
	b.balance -= amount
	return nil
}

// GetBalance returns the current balance
func (b *BankAccount) GetBalance() float64 {
	return b.balance
}

// GetAccountNumber returns the account number
func (b *BankAccount) GetAccountNumber() string {
	return b.accountNumber
}

// ============================================================================
// Struct with Anonymous Fields
// ============================================================================

// Product has anonymous fields
type Product struct {
	string  // name
	float64 // price
	int     // stock
}

// NewProduct creates a new product
func NewProduct(name string, price float64, stock int) Product {
	return Product{name, price, stock}
}

// Name returns the product name
func (p Product) Name() string {
	return p.string
}

// Price returns the product price
func (p Product) Price() float64 {
	return p.float64
}

// Stock returns the product stock
func (p Product) Stock() int {
	return p.int
}

// ============================================================================
// Struct with Slice and Map Fields
// ============================================================================

// Team represents a team with members
type Team struct {
	Name     string
	Members  []string
	Tags     map[string]string
	Metadata map[string]interface{}
}

// AddMember adds a member to the team
func (t *Team) AddMember(member string) {
	t.Members = append(t.Members, member)
}

// AddTag adds a tag to the team
func (t *Team) AddTag(key, value string) {
	if t.Tags == nil {
		t.Tags = make(map[string]string)
	}
	t.Tags[key] = value
}

// ============================================================================
// Struct with Time Fields
// ============================================================================

// Event represents an event with time fields
type Event struct {
	ID          int
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IsOngoing checks if the event is currently happening
func (e Event) IsOngoing() bool {
	now := time.Now()
	return now.After(e.StartTime) && now.Before(e.EndTime)
}

// Duration returns the event duration
func (e Event) Duration() time.Duration {
	return e.EndTime.Sub(e.StartTime)
}

// ============================================================================
// JSON Serialization Examples
// ============================================================================

// jsonExample demonstrates JSON marshaling/unmarshaling
func jsonExample() {
	user := User{
		ID:        1,
		Username:  "alice",
		Email:     "alice@example.com",
		Password:  "secret123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
		return
	}
	fmt.Printf("User JSON:\n%s\n", jsonData)

	// Unmarshal from JSON
	var decodedUser User
	err = json.Unmarshal(jsonData, &decodedUser)
	if err != nil {
		fmt.Printf("JSON unmarshal error: %v\n", err)
		return
	}
	fmt.Printf("Decoded user: %+v\n", decodedUser)
}

// ============================================================================
// Constructor Patterns
// ============================================================================

// Config represents application configuration
type Config struct {
	Host     string
	Port     int
	Timeout  time.Duration
	MaxConns int
	Debug    bool
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     8080,
		Timeout:  30 * time.Second,
		MaxConns: 100,
		Debug:    false,
	}
}

// WithHost sets the host
func (c *Config) WithHost(host string) *Config {
	c.Host = host
	return c
}

// WithPort sets the port
func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}

// WithDebug sets debug mode
func (c *Config) WithDebug(debug bool) *Config {
	c.Debug = debug
	return c
}

// ============================================================================
// Comparison and Equality
// ============================================================================

// Point represents a 2D point
type Point struct {
	X, Y int
}

// Equals checks if two points are equal
func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

// Distance calculates the distance to another point
func (p Point) Distance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return float64(dx*dx + dy*dy)
}

// ============================================================================
// Main Function
// ============================================================================

func main() {
	fmt.Println("=== Basic Struct ===")
	basicStruct()

	fmt.Println("\n=== Struct with Methods ===")
	emp := Employee{
		ID:       1,
		Name:     "Alice",
		Salary:   50000,
		HireDate: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	fmt.Printf("Employee: %s\n", emp.String())
	fmt.Printf("Years of service: %.1f\n", emp.YearsOfService())
	emp.GiveRaise(10)
	fmt.Printf("After raise: %s\n", emp.String())

	fmt.Println("\n=== Nested Structs ===")
	customer := Customer{
		ID:   1,
		Name: "Bob Smith",
		Address: Address{
			Street:  "123 Main St",
			City:    "Springfield",
			State:   "IL",
			ZipCode: "62701",
			Country: "USA",
		},
		Contact: Contact{
			Email:     "bob@example.com",
			Phone:     "555-123-4567",
			Preferred: "email",
		},
		CreatedAt: time.Now(),
	}
	fmt.Printf("Customer: %+v\n", customer)
	fmt.Printf("Full address: %s\n", customer.FullAddress())
	fmt.Printf("Contact info: %s\n", customer.ContactInfo())

	fmt.Println("\n=== Embedded Structs (Composition) ===")
	dog := Dog{
		Animal: Animal{Name: "Rex", Age: 3, Weight: 25.5},
		Breed:  "German Shepherd",
	}
	cat := Cat{
		Animal: Animal{Name: "Whiskers", Age: 2, Weight: 4.5},
		Indoor: true,
	}
	bird := Bird{
		Animal:   Animal{Name: "Tweety", Age: 1, Weight: 0.1},
		Wingspan: 15.5,
	}

	animals := []Sounder{dog, cat, bird}
	for _, animal := range animals {
		fmt.Printf("%s says: %s\n", animal.(interface{ Name() string }).(interface{ Name() string }).Name(), animal.MakeSound())
		// Note: In real code, you'd use type assertion or interface with Name() method
	}

	// Better way: using a custom interface
	type NamedSounder interface {
		Sounder
		Name() string
	}

	// Add Name method to Animal for this example
	// (in real code, you'd define the method on the embedded type)
	dog.Animal.Name = "Rex"
	cat.Animal.Name = "Whiskers"
	bird.Animal.Name = "Tweety"

	fmt.Println("\n=== Private Fields ===")
	account := NewBankAccount("Alice")
	account.Deposit(1000)
	account.Withdraw(250)
	fmt.Printf("Account %s balance: $%.2f\n", account.GetAccountNumber(), account.GetBalance())

	fmt.Println("\n=== Anonymous Fields ===")
	product := NewProduct("Laptop", 999.99, 10)
	fmt.Printf("Product: %s, Price: $%.2f, Stock: %d\n", product.Name(), product.Price(), product.Stock())

	fmt.Println("\n=== Slice and Map Fields ===")
	team := Team{Name: "Developers"}
	team.AddMember("Alice")
	team.AddMember("Bob")
	team.AddTag("language", "Go")
	team.AddTag("framework", "Gin")
	fmt.Printf("Team %s: members=%v, tags=%v\n", team.Name, team.Members, team.Tags)

	fmt.Println("\n=== Time Fields ===")
	event := Event{
		ID:          1,
		Title:       "Meeting",
		Description: "Project sync",
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(1 * time.Hour),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	fmt.Printf("Event ongoing: %v\n", event.IsOngoing())
	fmt.Printf("Duration: %v\n", event.Duration())

	fmt.Println("\n=== JSON Serialization ===")
	jsonExample()

	fmt.Println("\n=== Constructor Pattern ===")
	config := DefaultConfig().
		WithHost("example.com").
		WithPort(9090).
		WithDebug(true)
	fmt.Printf("Config: %+v\n", config)

	fmt.Println("\n=== Comparison ===")
	p1 := Point{3, 4}
	p2 := Point{3, 4}
	p3 := Point{5, 12}
	fmt.Printf("p1 equals p2: %v\n", p1.Equals(p2))
	fmt.Printf("p1 equals p3: %v\n", p1.Equals(p3))
}
