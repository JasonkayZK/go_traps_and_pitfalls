package main

import "fmt"

/* defer modify the return value */
func main() {
	// 12
	fmt.Println(func(x int) (result int) {
		defer func() {result += x}()
		return x * 2
	}(4))
}
