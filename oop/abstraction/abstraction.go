// Package abstraction demonstrates how Go achieves abstraction through interfaces.
//
// KEY CONCEPT: Abstraction = hiding implementation details, exposing only what's needed.
//
// In Java:
//   abstract class Database { abstract void save(String data); }
//   interface PaymentProcessor { void processPayment(double amount); }
//
// In Go:
//   type Database interface { Save(data string) error }
//   type PaymentProcessor interface { ProcessPayment(amount float64) error }
//
// Go has NO abstract classes. Use interfaces for abstraction.
package abstraction

import (
	"errors"
	"fmt"
	"strings"
)

// =============================================================================
// ABSTRACTION THROUGH INTERFACES
// Interfaces define WHAT something does, not HOW it does it.
// =============================================================================

// Database abstracts away the details of data storage.
// Users of this interface don't need to know if it's MySQL, MongoDB, or a file.
type Database interface {
	Save(key string, data string) error
	Load(key string) (string, error)
	Delete(key string) error
}

// --- Concrete Implementations (hidden details) ---

// MemoryDatabase stores data in memory.
// The implementation details are hidden from users of the Database interface.
type MemoryDatabase struct {
	data map[string]string
}

// NewMemoryDatabase creates a new in-memory database.
func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		data: make(map[string]string),
	}
}

func (m *MemoryDatabase) Save(key string, data string) error {
	m.data[key] = data
	fmt.Printf("[MemoryDB] Saved '%s'\n", key)
	return nil
}

func (m *MemoryDatabase) Load(key string) (string, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", errors.New("key not found")
}

func (m *MemoryDatabase) Delete(key string) error {
	delete(m.data, key)
	fmt.Printf("[MemoryDB] Deleted '%s'\n", key)
	return nil
}

// FileDatabase simulates file-based storage.
type FileDatabase struct {
	basePath string
	files    map[string]string // Simulated file system
}

// NewFileDatabase creates a new file-based database.
func NewFileDatabase(basePath string) *FileDatabase {
	return &FileDatabase{
		basePath: basePath,
		files:    make(map[string]string),
	}
}

func (f *FileDatabase) Save(key string, data string) error {
	path := f.basePath + "/" + key + ".txt"
	f.files[path] = data
	fmt.Printf("[FileDB] Wrote to file: %s\n", path)
	return nil
}

func (f *FileDatabase) Load(key string) (string, error) {
	path := f.basePath + "/" + key + ".txt"
	if val, ok := f.files[path]; ok {
		return val, nil
	}
	return "", errors.New("file not found")
}

func (f *FileDatabase) Delete(key string) error {
	path := f.basePath + "/" + key + ".txt"
	delete(f.files, path)
	fmt.Printf("[FileDB] Deleted file: %s\n", path)
	return nil
}

// =============================================================================
// USING ABSTRACTION: Code depends on interface, not implementation
// =============================================================================

// UserService uses the Database interface.
// It doesn't care HOW data is stored - that's abstracted away.
type UserService struct {
	db Database // Depends on interface, not concrete type
}

// NewUserService creates a service with any Database implementation.
func NewUserService(db Database) *UserService {
	return &UserService{db: db}
}

// CreateUser saves a user - implementation details are hidden.
func (s *UserService) CreateUser(id, name, email string) error {
	userData := fmt.Sprintf("name=%s,email=%s", name, email)
	return s.db.Save("user:"+id, userData)
}

// GetUser loads a user - doesn't know if it's from memory, file, or cloud.
func (s *UserService) GetUser(id string) (string, error) {
	return s.db.Load("user:" + id)
}

// =============================================================================
// ANOTHER EXAMPLE: Payment Processing Abstraction
// =============================================================================

// PaymentProcessor abstracts payment handling.
type PaymentProcessor interface {
	ProcessPayment(amount float64) error
	Refund(transactionID string, amount float64) error
	GetName() string
}

// CreditCardProcessor handles credit card payments.
type CreditCardProcessor struct {
	cardNumber string
}

func NewCreditCardProcessor(cardNumber string) *CreditCardProcessor {
	// Mask card number for security
	masked := "****-****-****-" + cardNumber[len(cardNumber)-4:]
	return &CreditCardProcessor{cardNumber: masked}
}

func (c *CreditCardProcessor) ProcessPayment(amount float64) error {
	fmt.Printf("[CreditCard %s] Processing $%.2f...\n", c.cardNumber, amount)
	fmt.Println("[CreditCard] Contacting bank... Approved!")
	return nil
}

func (c *CreditCardProcessor) Refund(txID string, amount float64) error {
	fmt.Printf("[CreditCard] Refunding $%.2f for transaction %s\n", amount, txID)
	return nil
}

func (c *CreditCardProcessor) GetName() string {
	return "Credit Card"
}

// PayPalProcessor handles PayPal payments.
type PayPalProcessor struct {
	email string
}

func NewPayPalProcessor(email string) *PayPalProcessor {
	return &PayPalProcessor{email: email}
}

func (p *PayPalProcessor) ProcessPayment(amount float64) error {
	fmt.Printf("[PayPal %s] Processing $%.2f...\n", p.email, amount)
	fmt.Println("[PayPal] Redirecting to PayPal... Payment confirmed!")
	return nil
}

func (p *PayPalProcessor) Refund(txID string, amount float64) error {
	fmt.Printf("[PayPal] Refunding $%.2f for transaction %s\n", amount, txID)
	return nil
}

func (p *PayPalProcessor) GetName() string {
	return "PayPal"
}

// CryptoProcessor handles cryptocurrency payments.
type CryptoProcessor struct {
	walletAddress string
}

func NewCryptoProcessor(walletAddress string) *CryptoProcessor {
	return &CryptoProcessor{walletAddress: walletAddress}
}

func (c *CryptoProcessor) ProcessPayment(amount float64) error {
	fmt.Printf("[Crypto] Sending $%.2f to wallet %s...\n", amount, c.walletAddress[:8]+"...")
	fmt.Println("[Crypto] Waiting for blockchain confirmation... Confirmed!")
	return nil
}

func (c *CryptoProcessor) Refund(txID string, amount float64) error {
	fmt.Printf("[Crypto] Initiating refund of $%.2f for %s\n", amount, txID)
	return nil
}

func (c *CryptoProcessor) GetName() string {
	return "Cryptocurrency"
}

// =============================================================================
// CHECKOUT SERVICE: Uses abstraction to support any payment method
// =============================================================================

// CheckoutService processes orders using any payment processor.
type CheckoutService struct {
	processor PaymentProcessor
}

// NewCheckoutService creates a checkout with the given payment processor.
func NewCheckoutService(p PaymentProcessor) *CheckoutService {
	return &CheckoutService{processor: p}
}

// Checkout processes a purchase - implementation details are abstracted.
func (c *CheckoutService) Checkout(items []string, total float64) error {
	fmt.Printf("\n--- Checkout ---\n")
	fmt.Printf("Items: %s\n", strings.Join(items, ", "))
	fmt.Printf("Total: $%.2f\n", total)
	fmt.Printf("Payment Method: %s\n", c.processor.GetName())
	fmt.Println()
	return c.processor.ProcessPayment(total)
}

// =============================================================================
// Demo function
// =============================================================================

// Demo demonstrates abstraction in Go.
func Demo() {
	fmt.Println("=== ABSTRACTION (Hiding Implementation Details) ===")

	// --- Database Abstraction Demo ---
	fmt.Println("\n--- Database Abstraction ---")
	fmt.Println("Same UserService code works with different database implementations:\n")

	// Use memory database
	fmt.Println("Using MemoryDatabase:")
	memDB := NewMemoryDatabase()
	userService1 := NewUserService(memDB)
	userService1.CreateUser("001", "Alice", "alice@example.com")

	// Use file database
	fmt.Println("\nUsing FileDatabase:")
	fileDB := NewFileDatabase("/data/users")
	userService2 := NewUserService(fileDB)
	userService2.CreateUser("002", "Bob", "bob@example.com")

	fmt.Println("\nThe UserService doesn't know or care which database is used!")

	// --- Payment Abstraction Demo ---
	fmt.Println("\n--- Payment Processor Abstraction ---")
	fmt.Println("Same CheckoutService works with different payment methods:")

	items := []string{"Laptop", "Mouse", "Keyboard"}

	// Checkout with Credit Card
	ccProcessor := NewCreditCardProcessor("1234567890123456")
	checkout1 := NewCheckoutService(ccProcessor)
	checkout1.Checkout(items, 1299.99)

	// Checkout with PayPal
	paypalProcessor := NewPayPalProcessor("user@example.com")
	checkout2 := NewCheckoutService(paypalProcessor)
	checkout2.Checkout(items, 1299.99)

	// Checkout with Crypto
	cryptoProcessor := NewCryptoProcessor("0x1234567890abcdef")
	checkout3 := NewCheckoutService(cryptoProcessor)
	checkout3.Checkout(items, 1299.99)

	// --- Key Points ---
	fmt.Println("\n--- Key Points about Go Abstraction ---")
	fmt.Println("1. Use interfaces to define WHAT, hide HOW")
	fmt.Println("2. Go has NO abstract classes - use interfaces instead")
	fmt.Println("3. Code depends on interfaces (abstractions), not concrete types")
	fmt.Println("4. Easy to swap implementations without changing dependent code")
	fmt.Println("5. Makes testing easier - can inject mock implementations")
	fmt.Println("6. 'Accept interfaces, return structs' - common Go idiom")
}
