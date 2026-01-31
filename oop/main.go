// OOP in Go - Demonstrating the Four Pillars of Object-Oriented Programming
//
// Go is NOT a traditional OOP language like Java, but it supports OOP concepts
// in its own way. This module demonstrates how Go achieves each pillar.
//
// THE FOUR PILLARS:
// 1. Inheritance   → Go uses COMPOSITION and EMBEDDING instead
// 2. Polymorphism  → Go uses INTERFACES for polymorphic behavior
// 3. Encapsulation → Go uses EXPORTED/UNEXPORTED (uppercase/lowercase)
// 4. Abstraction   → Go uses INTERFACES to hide implementation details
//
// Run this program to see all demos: go run .
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example/oop/abstraction"
	"example/oop/composition"
	"example/oop/encapsulation"
	"example/oop/polymorphism"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     OOP IN GO - THE FOUR PILLARS (Go's Approach)             ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  1. Composition   - Go's alternative to inheritance          ║")
	fmt.Println("║  2. Polymorphism  - Many forms through interfaces            ║")
	fmt.Println("║  3. Encapsulation - Data hiding with exported/unexported     ║")
	fmt.Println("║  4. Abstraction   - Hiding details with interfaces           ║")
	fmt.Println("║  5. Run ALL demos                                            ║")
	fmt.Println("║  0. Exit                                                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nEnter your choice (0-5): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			composition.Demo()
		case "2":
			polymorphism.Demo()
		case "3":
			encapsulation.Demo()
		case "4":
			abstraction.Demo()
		case "5":
			fmt.Println("\n" + strings.Repeat("=", 70))
			composition.Demo()
			fmt.Println("\n" + strings.Repeat("=", 70))
			polymorphism.Demo()
			fmt.Println("\n" + strings.Repeat("=", 70))
			encapsulation.Demo()
			fmt.Println("\n" + strings.Repeat("=", 70))
			abstraction.Demo()
			fmt.Println("\n" + strings.Repeat("=", 70))
			fmt.Println("All demos completed!")
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please enter 0-5.")
		}
	}
}
