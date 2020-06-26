package main

import "fmt"

func deferReturn1() int {
	a := 2
	defer func() {
		a = a + 2
	}()
	return a
}

func add(i *int) {
	*i = *i + 2
}

func deferReturn2() int {
	a := 2
	defer add(&a)
	return a
}

func main() {
	// 控制台输出：2
	fmt.Println(deferReturn1())
	// 控制台输出：2
	fmt.Println(deferReturn2())
}
