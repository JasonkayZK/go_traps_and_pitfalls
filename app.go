package main

import "fmt"

func main() {
	slice1 := make([]string, 0)

	slice1 = append(slice1, "a")
	slice1 = append(slice1, "b")
	slice1 = append(slice1, "c")

	fmt.Println(slice1)
	fmt.Println(len(slice1))
	fmt.Println(cap(slice1))

	slice2 := make([]string, 3)
	slice2 = append(slice2, "a")
	slice2 = append(slice2, "b")
	slice2 = append(slice2, "c")

	fmt.Println(slice2)
	fmt.Println(len(slice2))
	fmt.Println(cap(slice2))

	slice3 := make([]string, 3, 3)
	slice3 = append(slice3, "a")
	slice3 = append(slice3, "b")
	slice3 = append(slice3, "c")

	fmt.Println(slice3)
	fmt.Println(len(slice3))
	fmt.Println(cap(slice3))

	slice4 := make([]string, 0, 3)
	slice4 = append(slice4, "a")
	slice4 = append(slice4, "b")
	slice4 = append(slice4, "c")

	fmt.Println(slice4)
	fmt.Println(len(slice4))
	fmt.Println(cap(slice4))
}