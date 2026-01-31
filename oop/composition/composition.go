// Package composition demonstrates Go's alternative to inheritance.
//
// KEY CONCEPT: Go does NOT have inheritance like Java.
// Instead, Go uses COMPOSITION and EMBEDDING to achieve code reuse.
//
// In Java:
//   class Dog extends Animal { ... }  // Dog inherits from Animal
//
// In Go:
//   type Dog struct { Animal }        // Dog embeds Animal (composition)
package composition

import "fmt"

// =============================================================================
// EMBEDDING: Go's way to "inherit" fields and methods
// =============================================================================

// Animal is a base type with common fields and methods.
// In Java, this would be a parent class.
type Animal struct {
	Name string
	Age  int
}

// Speak is a method on Animal.
// Embedded types will "inherit" this method automatically.
func (a Animal) Speak() string {
	return fmt.Sprintf("%s makes a sound", a.Name)
}

// Info returns basic animal information.
func (a Animal) Info() string {
	return fmt.Sprintf("Name: %s, Age: %d", a.Name, a.Age)
}

// -----------------------------------------------------------------------------

// Dog embeds Animal - it gets all Animal's fields and methods automatically.
// This is similar to "class Dog extends Animal" in Java.
type Dog struct {
	Animal        // Embedded (anonymous) field - Dog "inherits" Animal
	Breed  string // Dog's own field
}

// Speak overrides the Animal's Speak method.
// This is similar to method overriding in Java.
func (d Dog) Speak() string {
	return fmt.Sprintf("%s says: Woof!", d.Name) // Can access Animal's Name directly
}

// Fetch is a Dog-specific method.
func (d Dog) Fetch() string {
	return fmt.Sprintf("%s fetches the ball!", d.Name)
}

// -----------------------------------------------------------------------------

// Cat also embeds Animal.
type Cat struct {
	Animal        // Embedded Animal
	Indoor bool   // Cat's own field
}

// Speak overrides the Animal's Speak method.
func (c Cat) Speak() string {
	return fmt.Sprintf("%s says: Meow!", c.Name)
}

// Purr is a Cat-specific method.
func (c Cat) Purr() string {
	return fmt.Sprintf("%s purrs contentedly", c.Name)
}

// =============================================================================
// COMPOSITION (has-a relationship): An alternative to embedding
// =============================================================================

// Engine is a component that can be composed into other types.
type Engine struct {
	Horsepower int
	Type       string // "V6", "V8", "Electric", etc.
}

// Start starts the engine.
func (e Engine) Start() string {
	return fmt.Sprintf("%s engine with %d HP started", e.Type, e.Horsepower)
}

// Car HAS-A Engine (composition), rather than IS-A Engine.
// The Engine field is named, so you access it via car.Engine.Start()
type Car struct {
	Brand  string
	Model  string
	Engine Engine // Named field - composition (has-a)
}

// Start delegates to the engine.
func (c Car) Start() string {
	return fmt.Sprintf("%s %s: %s", c.Brand, c.Model, c.Engine.Start())
}

// =============================================================================
// Demo function to show composition in action
// =============================================================================

// Demo demonstrates composition and embedding concepts.
func Demo() {
	fmt.Println("=== COMPOSITION (Go's Alternative to Inheritance) ===")

	// --- Embedding Demo ---
	fmt.Println("\n--- Embedding (similar to inheritance) ---")

	// Create a Dog - note how we initialize the embedded Animal
	dog := Dog{
		Animal: Animal{Name: "Buddy", Age: 3},
		Breed:  "Golden Retriever",
	}

	// Dog "inherits" Info() from Animal
	fmt.Println("Dog Info:", dog.Info())
	// Dog overrides Speak()
	fmt.Println("Dog Speak:", dog.Speak())
	// Dog has its own method
	fmt.Println("Dog Fetch:", dog.Fetch())
	// Can access embedded fields directly
	fmt.Printf("Dog Name: %s, Breed: %s\n", dog.Name, dog.Breed)

	fmt.Println()

	// Create a Cat
	cat := Cat{
		Animal: Animal{Name: "Whiskers", Age: 2},
		Indoor: true,
	}

	fmt.Println("Cat Info:", cat.Info()) // "Inherited" from Animal
	fmt.Println("Cat Speak:", cat.Speak()) // Overridden
	fmt.Println("Cat Purr:", cat.Purr())   // Cat-specific

	// --- Composition Demo ---
	fmt.Println("\n--- Composition (has-a relationship) ---")

	car := Car{
		Brand: "Tesla",
		Model: "Model S",
		Engine: Engine{
			Horsepower: 670,
			Type:       "Electric",
		},
	}

	fmt.Println(car.Start())
	// Access composed object directly
	fmt.Printf("Engine details: %s, %d HP\n", car.Engine.Type, car.Engine.Horsepower)

	// --- Key Differences from Java ---
	fmt.Println("\n--- Key Differences from Java Inheritance ---")
	fmt.Println("1. No 'extends' keyword - use embedding instead")
	fmt.Println("2. No 'super' keyword - access embedded type by name (e.g., dog.Animal.Speak())")
	fmt.Println("3. Methods are 'inherited' but can be overridden")
	fmt.Println("4. Embedded fields are accessed directly (dog.Name, not dog.Animal.Name)")

	// Show accessing the "parent" method
	fmt.Println("\nCalling 'parent' method explicitly:")
	fmt.Println("dog.Animal.Speak():", dog.Animal.Speak()) // Calls Animal's Speak, not Dog's
}
