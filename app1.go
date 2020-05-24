package main

import "fmt"

/* A simple demo for showing how multiple defer works in a single function */
func main() {
	defer func() {
		fmt.Println("First defer declared!")
	}()

	defer func() {
		fmt.Println("Second defer declared!")
	}()

	func() {
		fmt.Println("A function declared in main!")
	}()
}

/*
Output:
	A function declared in main!
	Second defer declared!
	First defer declared!
*/