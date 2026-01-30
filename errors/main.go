package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// =============================================================================
// BASIC ERROR HANDLING: The (value, error) pattern
// =============================================================================

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func demoBasicErrors() {
	fmt.Println("\n=== BASIC ERROR HANDLING ===")

	// Always check errors!
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("10 / 2 =", result)

	// Error case
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// =============================================================================
// CUSTOM ERROR TYPES: Carry more information
// =============================================================================

// Custom error type (implements error interface)
type ValidationError struct {
	Field   string
	Message string
}

// Implement the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "cannot be negative"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "exceeds maximum"}
	}
	return nil
}

func demoCustomErrors() {
	fmt.Println("\n=== CUSTOM ERROR TYPES ===")

	err := validateAge(-5)
	if err != nil {
		fmt.Println("Error:", err)

		// Type assertion to access custom fields
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("  Field: %s\n", valErr.Field)
			fmt.Printf("  Message: %s\n", valErr.Message)
		}
	}
}

// =============================================================================
// SENTINEL ERRORS: Predefined errors for comparison
// =============================================================================

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
)

func getUser(id int) (string, error) {
	switch id {
	case 1:
		return "Alice", nil
	case 2:
		return "", ErrUnauthorized
	default:
		return "", ErrNotFound
	}
}

func demoSentinelErrors() {
	fmt.Println("\n=== SENTINEL ERRORS ===")

	_, err := getUser(99)

	// Use errors.Is() to compare (handles wrapping)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("User not found - show 404 page")
	} else if errors.Is(err, ErrUnauthorized) {
		fmt.Println("Not authorized - redirect to login")
	}
}

// =============================================================================
// ERROR WRAPPING: Add context while preserving original error
// =============================================================================

func readConfig(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Wrap error with context using %w
		return nil, fmt.Errorf("failed to read config from %s: %w", path, err)
	}
	return data, nil
}

func initializeApp() error {
	_, err := readConfig("/nonexistent/config.json")
	if err != nil {
		// Wrap again with more context
		return fmt.Errorf("app initialization failed: %w", err)
	}
	return nil
}

func demoErrorWrapping() {
	fmt.Println("\n=== ERROR WRAPPING ===")

	err := initializeApp()
	if err != nil {
		// Full error message includes all context
		fmt.Println("Error:", err)

		// Unwrap to check the underlying error
		fmt.Println("\n--- Checking wrapped errors ---")

		// errors.Is checks the entire chain
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("Root cause: file does not exist")
		}

		// Unwrap manually
		fmt.Println("\n--- Unwrapping chain ---")
		for e := err; e != nil; e = errors.Unwrap(e) {
			fmt.Printf("  -> %v\n", e)
		}
	}
}

// =============================================================================
// MULTIPLE ERRORS: errors.Join (Go 1.20+)
// =============================================================================

func validateUser(name string, age int, email string) error {
	var errs []error

	if name == "" {
		errs = append(errs, errors.New("name is required"))
	}
	if age < 0 {
		errs = append(errs, errors.New("age cannot be negative"))
	}
	if email == "" {
		errs = append(errs, errors.New("email is required"))
	}

	return errors.Join(errs...) // Returns nil if errs is empty
}

func demoMultipleErrors() {
	fmt.Println("\n=== MULTIPLE ERRORS (errors.Join) ===")

	err := validateUser("", -5, "")
	if err != nil {
		fmt.Println("Validation errors:")
		fmt.Println(err)
	}
}

// =============================================================================
// PANIC AND RECOVER: Only for truly exceptional cases
// =============================================================================

func riskyOperation() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	fmt.Println("About to panic...")
	panic("something went terribly wrong")
	fmt.Println("This line never executes")
}

func demoPanicRecover() {
	fmt.Println("\n=== PANIC AND RECOVER ===")
	fmt.Println("(See defer-panic module for more details)")

	riskyOperation()
	fmt.Println("Program continues after recovery")
}

// =============================================================================
// ERROR HANDLING PATTERNS
// =============================================================================

func demoPatterns() {
	fmt.Println("\n=== ERROR HANDLING PATTERNS ===")

	// Pattern 1: Early return (most common)
	fmt.Println("\n1. Early return pattern:")
	fmt.Println(`
    func process() error {
        data, err := fetch()
        if err != nil {
            return err  // Return early
        }
        // Continue with happy path
        return nil
    }`)

	// Pattern 2: Error variable
	fmt.Println("\n2. Named error return:")
	fmt.Println(`
    func process() (result int, err error) {
        defer func() {
            if err != nil {
                // Cleanup on error
            }
        }()
        // ...
    }`)

	// Pattern 3: Must functions (panic on error)
	fmt.Println("\n3. Must pattern (for init-time errors):")
	fmt.Println(`
    var template = mustParse("template.html")
    
    func mustParse(name string) *Template {
        t, err := Parse(name)
        if err != nil {
            panic(err)  // Acceptable during init
        }
        return t
    }`)
}

func main() {
	fmt.Println("========================================")
	fmt.Println("GO ERROR HANDLING")
	fmt.Println("========================================")

	demoBasicErrors()
	demoCustomErrors()
	demoSentinelErrors()
	demoErrorWrapping()
	demoMultipleErrors()
	demoPanicRecover()
	demoPatterns()
}
