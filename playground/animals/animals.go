package animals

import "fmt"

type Person struct {
	Name  string
	Age   int
	Email string
	Score []int
}

type Cat struct {
	Name  string
	Breed string
	Age   int
	Color string
}

func (p Person) String() string {
	return fmt.Sprintf("Person: Name: %s, Age: %d, Email: %s", p.GetName(), p.GetAge(), p.GetEmail())
}

func (c Cat) String() string {
	return fmt.Sprintf("Cat: Name: %s, Breed: %s, Age: %d, Color: %s", c.GetName(), c.GetBreed(), c.GetAge(), c.GetColor())
}

func (p Person) GetName() string {
	return p.Name
}

func (c Cat) GetName() string {
	return c.Name
}

func (p Person) GetAge() int {
	return p.Age
}

func (c Cat) GetAge() int {
	return c.Age
}

func (p Person) GetEmail() string {
	return p.Email
}

func (c Cat) GetBreed() string {
	return c.Breed
}

func (c Cat) GetColor() string {
	return c.Color
}