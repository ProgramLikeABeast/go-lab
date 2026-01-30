package main

import (
	"fmt"
	"example/playground/animals"
	"example/playground/shapes"
)

func main() {
	person := animals.Person{Name: "John", Email: "john@example.com", Score: []int{1, 2, 3, 4, 5}}
	fmt.Println(person.String())

	cat := animals.Cat{Name: "Whiskers", Breed: "Siamese", Age: 3, Color: "White"}
	fmt.Println(cat.String())

	square := shapes.Square{Width: 10}
	fmt.Println(square.String())

	circle := shapes.Circle{Radius: 10}
	fmt.Println(circle.String())
}