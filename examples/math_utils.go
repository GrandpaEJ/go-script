package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Math Utils Module ===")
	PI := 3.14159265359
	a := 5
	b := 3
	sum := (a + b)
	fmt.Println("Add 5 + 3 =", sum)
	x := 4
	y := 6
	product := (x * y)
	fmt.Println("Multiply 4 * 6 =", product)
	radius := 3
	area := ((PI * radius) * radius)
	fmt.Println("Circle area (r=3) =", area)
	fmt.Println("=== Module Demo Complete ===")
}
