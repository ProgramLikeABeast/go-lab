// Package polymorphism demonstrates how Go achieves polymorphism through interfaces.
//
// KEY CONCEPT: Polymorphism = "many forms"
// The same method call behaves differently depending on the actual type.
//
// In Java:
//   Animal animal = new Dog();  // Animal reference, Dog object
//   animal.speak();             // Calls Dog's speak()
//
// In Go:
//   var animal Speaker = Dog{}  // Speaker interface, Dog value
//   animal.Speak()              // Calls Dog's Speak()
package polymorphism

import (
	"fmt"
	"math"
)

// =============================================================================
// INTERFACE-BASED POLYMORPHISM
// =============================================================================

// Speaker is an interface that defines a behavior.
// Any type with a Speak() method automatically implements this interface.
type Speaker interface {
	Speak() string
}

// Dog implements Speaker (implicitly - no "implements" keyword needed).
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return fmt.Sprintf("%s says: Woof! Woof!", d.Name)
}

// Cat implements Speaker.
type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return fmt.Sprintf("%s says: Meow!", c.Name)
}

// Cow implements Speaker.
type Cow struct {
	Name string
}

func (c Cow) Speak() string {
	return fmt.Sprintf("%s says: Moo!", c.Name)
}

// Robot also implements Speaker - interfaces work with ANY type!
type Robot struct {
	Model string
}

func (r Robot) Speak() string {
	return fmt.Sprintf("Robot %s says: BEEP BOOP!", r.Model)
}

// =============================================================================
// POLYMORPHIC FUNCTIONS
// =============================================================================

// MakeSpeak accepts ANY type that implements Speaker.
// This is polymorphism: one function, many types.
func MakeSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

// Chorus makes multiple speakers speak together.
func Chorus(speakers []Speaker) {
	fmt.Println("--- All together now! ---")
	for _, s := range speakers {
		fmt.Println(s.Speak())
	}
}

// =============================================================================
// ANOTHER EXAMPLE: Shape polymorphism
// =============================================================================

// Shape interface - any type with Area() and Perimeter() is a Shape.
type Shape interface {
	Area() float64
	Perimeter() float64
	Name() string
}

// Rectangle implements Shape.
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }
func (r Rectangle) Name() string       { return "Rectangle" }

// Circle implements Shape.
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func (c Circle) Name() string       { return "Circle" }

// Triangle implements Shape.
type Triangle struct {
	A, B, C float64 // Side lengths
}

func (t Triangle) Area() float64 {
	// Heron's formula
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }
func (t Triangle) Name() string       { return "Triangle" }

// PrintShapeInfo works with ANY Shape - that's polymorphism!
func PrintShapeInfo(s Shape) {
	fmt.Printf("%s: Area = %.2f, Perimeter = %.2f\n",
		s.Name(), s.Area(), s.Perimeter())
}

// TotalArea calculates the total area of multiple shapes.
// Different types, same interface, same function.
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// =============================================================================
// Demo function
// =============================================================================

// Demo demonstrates polymorphism in Go.
func Demo() {
	fmt.Println("=== POLYMORPHISM (One Interface, Many Types) ===")

	// --- Speaker Polymorphism ---
	fmt.Println("\n--- Speaker Interface Polymorphism ---")

	// Create different types that all implement Speaker
	dog := Dog{Name: "Buddy"}
	cat := Cat{Name: "Whiskers"}
	cow := Cow{Name: "Bessie"}
	robot := Robot{Model: "R2D2"}

	// Same function (MakeSpeak) works with different types
	fmt.Println("Calling MakeSpeak() with different types:")
	MakeSpeak(dog)   // Calls Dog.Speak()
	MakeSpeak(cat)   // Calls Cat.Speak()
	MakeSpeak(cow)   // Calls Cow.Speak()
	MakeSpeak(robot) // Calls Robot.Speak()

	// Store different types in same slice (of interface type)
	fmt.Println("\nStoring different types in []Speaker slice:")
	speakers := []Speaker{dog, cat, cow, robot}
	Chorus(speakers)

	// --- Shape Polymorphism ---
	fmt.Println("\n--- Shape Interface Polymorphism ---")

	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	triangle := Triangle{A: 3, B: 4, C: 5}

	// Same function works with different shape types
	fmt.Println("Calling PrintShapeInfo() with different shapes:")
	PrintShapeInfo(rect)
	PrintShapeInfo(circle)
	PrintShapeInfo(triangle)

	// Calculate total area of mixed shapes
	shapes := []Shape{rect, circle, triangle}
	fmt.Printf("\nTotal area of all shapes: %.2f\n", TotalArea(shapes))

	// --- Key Points ---
	fmt.Println("\n--- Key Points about Go Polymorphism ---")
	fmt.Println("1. Interfaces define BEHAVIOR, not data")
	fmt.Println("2. Types implement interfaces IMPLICITLY (no 'implements' keyword)")
	fmt.Println("3. A function accepting an interface can work with ANY implementing type")
	fmt.Println("4. []Interface can hold values of different concrete types")
	fmt.Println("5. The correct method is called based on the actual type (dynamic dispatch)")
}
