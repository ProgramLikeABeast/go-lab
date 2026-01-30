package main

import (
	"fmt"
	"math"
)

// =============================================================================
// IMPLICIT INTERFACES: No "implements" keyword needed!
// If a type has all the methods, it implements the interface automatically.
// =============================================================================

// Define an interface
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle "implements" Shape by having Area() and Perimeter() methods
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle also implements Shape - no explicit declaration needed
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// This function accepts ANY Shape
func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func demoImplicitInterfaces() {
	fmt.Println("\n=== IMPLICIT INTERFACES ===")

	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}

	fmt.Print("Rectangle: ")
	printShapeInfo(rect) // Rectangle satisfies Shape interface

	fmt.Print("Circle: ")
	printShapeInfo(circle) // Circle satisfies Shape interface

	// You can put different types in a slice if they share an interface
	shapes := []Shape{rect, circle}
	fmt.Println("\nAll shapes:")
	for _, s := range shapes {
		printShapeInfo(s)
	}
}

// =============================================================================
// EMPTY INTERFACE: any / interface{} accepts ANY type
// =============================================================================

func printAnything(val any) { // "any" is alias for "interface{}"
	fmt.Printf("Value: %v, Type: %T\n", val, val)
}

func demoEmptyInterface() {
	fmt.Println("\n=== EMPTY INTERFACE (any) ===")

	printAnything(42)
	printAnything("hello")
	printAnything(3.14)
	printAnything([]int{1, 2, 3})
	printAnything(Rectangle{10, 5})
}

// =============================================================================
// TYPE ASSERTIONS: Extract concrete type from interface
// =============================================================================

func demoTypeAssertions() {
	fmt.Println("\n=== TYPE ASSERTIONS ===")

	var i any = "hello"

	// Type assertion: interface.(ConcreteType)
	s := i.(string)
	fmt.Println("String value:", s)

	// This would panic if i is not an int:
	// n := i.(int) // panic: interface conversion

	// Safe type assertion with ok pattern
	n, ok := i.(int)
	if ok {
		fmt.Println("Int value:", n)
	} else {
		fmt.Println("Not an int!")
	}

	// Type switch: elegant way to handle multiple types
	fmt.Println("\n=== TYPE SWITCH ===")
	values := []any{42, "hello", 3.14, true, Rectangle{10, 5}}

	for _, v := range values {
		switch val := v.(type) {
		case int:
			fmt.Printf("Integer: %d (doubled: %d)\n", val, val*2)
		case string:
			fmt.Printf("String: %q (length: %d)\n", val, len(val))
		case float64:
			fmt.Printf("Float: %.2f\n", val)
		case bool:
			fmt.Printf("Boolean: %t\n", val)
		case Shape:
			fmt.Printf("Shape with area: %.2f\n", val.Area())
		default:
			fmt.Printf("Unknown type: %T\n", val)
		}
	}
}

// =============================================================================
// INTERFACE COMPOSITION: Build interfaces from other interfaces
// =============================================================================

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

// ReadWriter is composed of Reader and Writer
type ReadWriter interface {
	Reader
	Writer
}

// =============================================================================
// STRINGER INTERFACE: Go's equivalent of toString()
// =============================================================================

type Person struct {
	Name string
	Age  int
}

// Implement fmt.Stringer interface
func (p Person) String() string {
	return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

func demoStringer() {
	fmt.Println("\n=== STRINGER INTERFACE ===")

	p := Person{Name: "Alice", Age: 30}

	// fmt.Println automatically calls String() if available
	fmt.Println(p) // Output: Alice (30 years old)
}

// =============================================================================
// NIL INTERFACE VALUES: A subtle gotcha
// =============================================================================

func demoNilInterfaces() {
	fmt.Println("\n=== NIL INTERFACE GOTCHA ===")

	var s *Rectangle = nil  // nil pointer of type *Rectangle
	var shape Shape = s     // Interface holds (type=*Rectangle, value=nil)

	// The interface is NOT nil! It has a type.
	if shape == nil {
		fmt.Println("shape is nil")
	} else {
		fmt.Println("shape is NOT nil (has type *Rectangle with nil value)")
		fmt.Printf("Type: %T, Value: %v\n", shape, shape)
	}

	// Truly nil interface
	var shape2 Shape // zero value is nil
	if shape2 == nil {
		fmt.Println("shape2 is truly nil (no type, no value)")
	}
}

func main() {
	fmt.Println("========================================")
	fmt.Println("GO INTERFACES")
	fmt.Println("========================================")

	demoImplicitInterfaces()
	demoEmptyInterface()
	demoTypeAssertions()
	demoStringer()
	demoNilInterfaces()
}
