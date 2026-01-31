// Package encapsulation demonstrates Go's approach to data hiding and access control.
//
// KEY CONCEPT: Go uses CAPITALIZATION to control visibility.
//
// In Java:
//   private String name;      // Only accessible within the class
//   public String getName()   // Accessible from anywhere
//
// In Go:
//   name string               // Unexported (lowercase) - only accessible within the package
//   Name string               // Exported (uppercase) - accessible from other packages
//   getName() string          // Unexported method
//   GetName() string          // Exported method
package encapsulation

import (
	"errors"
	"fmt"
)

// =============================================================================
// EXPORTED vs UNEXPORTED
// =============================================================================

// BankAccount demonstrates encapsulation with private (unexported) fields
// and public (exported) methods.
type BankAccount struct {
	// Exported fields (uppercase) - accessible from other packages
	AccountNumber string
	HolderName    string

	// Unexported fields (lowercase) - only accessible within this package
	balance    float64 // Private: can't be directly accessed from outside
	pin        string  // Private: sensitive data
	isLocked   bool    // Private: internal state
	loginAttempts int  // Private: internal counter
}

// =============================================================================
// CONSTRUCTOR FUNCTION
// Since fields are private, we provide a constructor to create valid objects.
// This is Go's equivalent of a Java constructor.
// =============================================================================

// NewBankAccount creates a new account with validated initial values.
// This is a "constructor" pattern - the recommended way to create objects
// when you need to enforce invariants.
func NewBankAccount(accountNumber, holderName, pin string, initialDeposit float64) (*BankAccount, error) {
	// Validate inputs (encapsulation lets us enforce business rules)
	if len(pin) != 4 {
		return nil, errors.New("PIN must be exactly 4 digits")
	}
	if initialDeposit < 0 {
		return nil, errors.New("initial deposit cannot be negative")
	}
	if accountNumber == "" || holderName == "" {
		return nil, errors.New("account number and holder name are required")
	}

	return &BankAccount{
		AccountNumber: accountNumber,
		HolderName:    holderName,
		balance:       initialDeposit,
		pin:           pin,
		isLocked:      false,
		loginAttempts: 0,
	}, nil
}

// =============================================================================
// GETTER METHODS (Exported)
// Provide controlled read access to private fields.
// =============================================================================

// GetBalance returns the current balance.
// This is a getter - provides read-only access to the private balance field.
func (a *BankAccount) GetBalance() float64 {
	return a.balance
}

// IsLocked returns whether the account is locked.
func (a *BankAccount) IsLocked() bool {
	return a.isLocked
}

// =============================================================================
// BUSINESS LOGIC METHODS (Exported)
// These methods encapsulate complex operations and enforce business rules.
// =============================================================================

// Deposit adds money to the account.
// Encapsulation ensures we can validate and control how balance changes.
func (a *BankAccount) Deposit(amount float64) error {
	if a.isLocked {
		return errors.New("account is locked")
	}
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}

	a.balance += amount
	return nil
}

// Withdraw removes money from the account with PIN verification.
// This shows encapsulation protecting sensitive operations.
func (a *BankAccount) Withdraw(amount float64, pin string) error {
	if a.isLocked {
		return errors.New("account is locked")
	}

	// Verify PIN (internal validation)
	if !a.verifyPIN(pin) {
		a.loginAttempts++
		if a.loginAttempts >= 3 {
			a.isLocked = true
			return errors.New("account locked due to too many failed attempts")
		}
		return fmt.Errorf("invalid PIN (%d attempts remaining)", 3-a.loginAttempts)
	}

	// Reset attempts on successful verification
	a.loginAttempts = 0

	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if amount > a.balance {
		return errors.New("insufficient funds")
	}

	a.balance -= amount
	return nil
}

// ChangePIN allows changing the PIN with old PIN verification.
func (a *BankAccount) ChangePIN(oldPIN, newPIN string) error {
	if a.isLocked {
		return errors.New("account is locked")
	}
	if !a.verifyPIN(oldPIN) {
		return errors.New("current PIN is incorrect")
	}
	if len(newPIN) != 4 {
		return errors.New("new PIN must be exactly 4 digits")
	}

	a.pin = newPIN
	return nil
}

// Unlock unlocks a locked account (would require admin verification in real system).
func (a *BankAccount) Unlock(adminOverride bool) error {
	if !adminOverride {
		return errors.New("admin override required to unlock account")
	}
	a.isLocked = false
	a.loginAttempts = 0
	return nil
}

// =============================================================================
// PRIVATE (UNEXPORTED) HELPER METHODS
// These can only be called within the encapsulation package.
// =============================================================================

// verifyPIN is a private method - only accessible within this package.
// External code cannot bypass PIN verification.
func (a *BankAccount) verifyPIN(pin string) bool {
	return a.pin == pin
}

// =============================================================================
// ANOTHER EXAMPLE: User with private validation
// =============================================================================

// User demonstrates encapsulation with validation.
type User struct {
	// Exported
	Username string

	// Unexported - must use methods to access/modify
	email    string
	password string
	isActive bool
}

// NewUser creates a validated User.
func NewUser(username, email, password string) (*User, error) {
	if len(username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	// In real code, you'd validate email format too

	return &User{
		Username: username,
		email:    email,
		password: hashPassword(password), // Store hashed, not plain text
		isActive: true,
	}, nil
}

// GetEmail returns the user's email.
func (u *User) GetEmail() string {
	return u.email
}

// SetEmail updates email with validation.
func (u *User) SetEmail(newEmail string) error {
	// Add email validation here
	if newEmail == "" {
		return errors.New("email cannot be empty")
	}
	u.email = newEmail
	return nil
}

// CheckPassword verifies a password without exposing the stored hash.
func (u *User) CheckPassword(password string) bool {
	return u.password == hashPassword(password)
}

// Private helper function (unexported).
func hashPassword(password string) string {
	// In real code, use bcrypt or similar
	return "hashed_" + password
}

// =============================================================================
// Demo function
// =============================================================================

// Demo demonstrates encapsulation in Go.
func Demo() {
	fmt.Println("=== ENCAPSULATION (Data Hiding & Access Control) ===")

	// --- BankAccount Demo ---
	fmt.Println("\n--- BankAccount with Private Balance ---")

	// Create account using constructor (enforces validation)
	account, err := NewBankAccount("ACC001", "John Doe", "1234", 1000.0)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}

	fmt.Printf("Account: %s, Holder: %s\n", account.AccountNumber, account.HolderName)
	fmt.Printf("Balance: $%.2f\n", account.GetBalance())

	// Can't do: account.balance = 1000000 (unexported field)
	// Must use methods that enforce business rules

	// Deposit (doesn't need PIN)
	account.Deposit(500)
	fmt.Printf("After depositing $500: $%.2f\n", account.GetBalance())

	// Withdraw (needs PIN)
	err = account.Withdraw(200, "1234")
	if err != nil {
		fmt.Println("Withdraw error:", err)
	} else {
		fmt.Printf("After withdrawing $200: $%.2f\n", account.GetBalance())
	}

	// Try wrong PIN
	fmt.Println("\nTrying wrong PIN 3 times:")
	for i := 0; i < 3; i++ {
		err = account.Withdraw(100, "0000")
		fmt.Printf("  Attempt %d: %v\n", i+1, err)
	}
	fmt.Printf("Account locked: %v\n", account.IsLocked())

	// --- User Demo ---
	fmt.Println("\n--- User with Private Password ---")

	user, _ := NewUser("johndoe", "john@example.com", "securepass123")
	fmt.Printf("Username: %s, Email: %s\n", user.Username, user.GetEmail())

	// Can't access: user.password (unexported)
	// Can only verify password through method
	fmt.Printf("Password 'wrongpass' correct: %v\n", user.CheckPassword("wrongpass"))
	fmt.Printf("Password 'securepass123' correct: %v\n", user.CheckPassword("securepass123"))

	// --- Key Points ---
	fmt.Println("\n--- Key Points about Go Encapsulation ---")
	fmt.Println("1. Lowercase = unexported (private to package)")
	fmt.Println("2. Uppercase = exported (public, accessible from other packages)")
	fmt.Println("3. Use constructor functions (NewXxx) to enforce creation rules")
	fmt.Println("4. Use getter/setter methods to control access to private fields")
	fmt.Println("5. Encapsulation scope is PACKAGE level, not type level")
	fmt.Println("   (All code in the same package can access unexported fields)")
}
