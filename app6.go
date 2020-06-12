package main

import "fmt"

func main() {
	deferTest()
}

func deferTest() {
	fmt.Println("Into deferTest...")
	fmt.Println("Leave deferTest...")
	defer deferCallback()
	return
}

func deferCallback() {
	fmt.Println("Into defer Callback")
}
