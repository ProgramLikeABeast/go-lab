package main

import (
	"encoding/json"
	"fmt"
)

// =============================================================================
// STRUCT BASICS
// =============================================================================

type Person struct {
	Name    string
	Age     int
	Email   string
	private int // lowercase = unexported (package-private)
}

func demoStructBasics() {
	fmt.Println("\n=== STRUCT BASICS ===")

	// Method 1: Named fields (recommended)
	p1 := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

	// Method 2: Positional (fragile, not recommended)
	p2 := Person{"Bob", 25, "bob@example.com", 0}

	// Method 3: Zero value then assign
	var p3 Person
	p3.Name = "Charlie"
	p3.Age = 35

	fmt.Printf("p1: %+v\n", p1) // %+v shows field names
	fmt.Printf("p2: %+v\n", p2)
	fmt.Printf("p3: %+v\n", p3) // Email is "" (zero value)
}

// =============================================================================
// VALUE vs POINTER RECEIVERS: The most important distinction!
// =============================================================================

type Counter struct {
	value int
}

// Value receiver: operates on a COPY of the struct
func (c Counter) GetValue() int {
	return c.value
}

// Value receiver: this WON'T modify the original!
func (c Counter) IncrementBroken() {
	c.value++ // Modifies the copy, not the original
}

// Pointer receiver: operates on the ORIGINAL struct
func (c *Counter) Increment() {
	c.value++ // Modifies the actual struct
}

// Pointer receiver: can handle nil (useful pattern)
func (c *Counter) SafeIncrement() {
	if c == nil {
		fmt.Println("Cannot increment nil counter")
		return
	}
	c.value++
}

func demoReceivers() {
	fmt.Println("\n=== VALUE vs POINTER RECEIVERS ===")

	c := Counter{value: 10}

	fmt.Println("Initial value:", c.GetValue())

	c.IncrementBroken()
	fmt.Println("After IncrementBroken():", c.GetValue()) // Still 10!

	c.Increment()
	fmt.Println("After Increment():", c.GetValue()) // Now 11

	// Go automatically dereferences: c.Increment() == (&c).Increment()
	// And automatically references: (&c).GetValue() == c.GetValue()

	// Nil receiver handling
	var nilCounter *Counter
	nilCounter.SafeIncrement() // Won't panic
}

// =============================================================================
// STRUCT EMBEDDING: Go's composition (not inheritance!)
// =============================================================================

type Address struct {
	Street  string
	City    string
	Country string
}

func (a Address) FullAddress() string {
	return fmt.Sprintf("%s, %s, %s", a.Street, a.City, a.Country)
}

type Employee struct {
	Person  // Embedded: fields and methods are "promoted"
	Address // Can embed multiple types
	Title   string
	Salary  float64
}

func demoEmbedding() {
	fmt.Println("\n=== STRUCT EMBEDDING ===")

	emp := Employee{
		Person:  Person{Name: "Diana", Age: 28, Email: "diana@corp.com"},
		Address: Address{Street: "123 Main St", City: "Seattle", Country: "USA"},
		Title:   "Engineer",
		Salary:  100000,
	}

	// Access embedded fields directly (promoted)
	fmt.Println("Name:", emp.Name)         // Same as emp.Person.Name
	fmt.Println("City:", emp.City)         // Same as emp.Address.City
	fmt.Println("Title:", emp.Title)

	// Embedded methods are also promoted
	fmt.Println("Address:", emp.FullAddress()) // Same as emp.Address.FullAddress()

	// You can still access the embedded struct explicitly
	fmt.Printf("Person: %+v\n", emp.Person)
}

// =============================================================================
// ANONYMOUS STRUCTS: Quick one-off structures
// =============================================================================

func demoAnonymousStructs() {
	fmt.Println("\n=== ANONYMOUS STRUCTS ===")

	// Useful for one-off data, test fixtures, JSON unmarshaling
	point := struct {
		X, Y int
	}{10, 20}

	fmt.Printf("Point: %+v\n", point)

	// Common pattern: slice of anonymous structs for test cases
	testCases := []struct {
		input    int
		expected int
	}{
		{1, 2},
		{2, 4},
		{5, 10},
	}

	for _, tc := range testCases {
		result := tc.input * 2
		if result == tc.expected {
			fmt.Printf("✓ %d * 2 = %d\n", tc.input, result)
		}
	}
}

// =============================================================================
// STRUCT TAGS: Metadata for reflection (JSON, DB, validation)
// =============================================================================

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-"`                    // Exclude from JSON
	Email     string `json:"email,omitempty"`      // Omit if empty
	CreatedAt string `json:"created_at,omitempty"`
}

func demoStructTags() {
	fmt.Println("\n=== STRUCT TAGS ===")

	user := User{
		ID:       1,
		Username: "alice",
		Password: "secret123", // Won't appear in JSON
		Email:    "",          // Will be omitted (omitempty)
	}

	jsonBytes, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println("JSON output:")
	fmt.Println(string(jsonBytes))
}

// =============================================================================
// CONSTRUCTOR PATTERN: Go doesn't have constructors, use functions
// =============================================================================

type Server struct {
	host string
	port int
	tls  bool
}

// "Constructor" function - convention: NewTypeName
func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
		tls:  false, // sensible default
	}
}

// Functional options pattern (advanced but common)
type ServerOption func(*Server)

func WithTLS() ServerOption {
	return func(s *Server) {
		s.tls = true
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func NewServerWithOptions(host string, opts ...ServerOption) *Server {
	s := &Server{host: host, port: 8080} // defaults
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func demoConstructors() {
	fmt.Println("\n=== CONSTRUCTOR PATTERNS ===")

	// Simple constructor
	s1 := NewServer("localhost", 8080)
	fmt.Printf("Server 1: %+v\n", s1)

	// Functional options pattern
	s2 := NewServerWithOptions("api.example.com", WithTLS(), WithPort(443))
	fmt.Printf("Server 2: %+v\n", s2)
}

// =============================================================================
// COMPARING STRUCTS
// =============================================================================

func demoComparison() {
	fmt.Println("\n=== STRUCT COMPARISON ===")

	type Point struct {
		X, Y int
	}

	p1 := Point{1, 2}
	p2 := Point{1, 2}
	p3 := Point{3, 4}

	fmt.Println("p1 == p2:", p1 == p2) // true (all fields equal)
	fmt.Println("p1 == p3:", p1 == p3) // false

	// Structs with slices/maps/functions are NOT comparable with ==
	// type BadCompare struct {
	//     data []int // This makes the struct non-comparable
	// }
}

func main() {
	fmt.Println("========================================")
	fmt.Println("GO STRUCTS AND METHODS")
	fmt.Println("========================================")

	demoStructBasics()
	demoReceivers()
	demoEmbedding()
	demoAnonymousStructs()
	demoStructTags()
	demoConstructors()
	demoComparison()
}
