package main

import (
	"example/greetings"
	"fmt"
	"log"
	"rsc.io/quote"
)

func main() {
	fmt.Println("Hello, World!")
	// from standard library
	fmt.Println(quote.Go())
	// from local different module
	message, err := greetings.Hello("Gladys")
	fmt.Println(message)

	println("//////////////////////////////////////////////////////////////")//////////////////////////////////////////////////////////////

	names := []string{"Gladys", "Samantha", "Darrin"}
	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}
	for name, message := range messages {
		fmt.Println(name, message)
	}

    println("//////////////////////////////////////////////////////////////")//////////////////////////////////////////////////////////////

	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	message, err = greetings.Hello("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)
}