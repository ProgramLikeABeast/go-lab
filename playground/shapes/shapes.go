package shapes

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Square struct {
	Width int
}

type Circle struct {
	Radius int
}

func (s Square) Area() float64 {
	return float64(s.Width * s.Width)
}

func (s Square) Perimeter() float64 {
	return float64(4 * s.Width)
}

func (c Circle) Area() float64 {
	return float64(c.Radius*c.Radius) * math.Pi
}

func (c Circle) Perimeter() float64 {
	return 2 * float64(c.Radius) * math.Pi
}

func (s Square) String() string {
	return fmt.Sprintf("Square: Width: %d, Area: %f, Perimeter: %f", s.Width, s.Area(), s.Perimeter())
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle: Radius: %d, Area: %f, Perimeter: %f", c.Radius, c.Area(), c.Perimeter())
}
