package main

import "fmt"

type Pointer struct {
	x int64
	y int64
}

func (p Pointer) modifiedWithNonPointer() {
	p.x = 10000
	p.y = 10000
}

func (p *Pointer) modifiedWithPointer() {
	p.x = 10000
	p.y = 10000
}

/* 注意receiver是指针类型还是非指针类型 */
func main() {
	p1 := Pointer{x: 1, y: 1}
	p2 := Pointer{x: 2, y: 2}

	p1.modifiedWithNonPointer()
	fmt.Printf("modified with non-pointer receiver: %v\n", p1)

	p2.modifiedWithPointer()
	fmt.Printf("modified with pointer receiver: %v\n", p2)
}

/*
Output:
	modified with non-pointer receiver: {1 1}
	modified with pointer receiver: {10000 10000}
*/