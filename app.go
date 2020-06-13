package main

import "fmt"

var myMap map[string]string

func init() {
	myMap = make(map[string]string)
	myMap["1"] = "test1"
	myMap["2"] = "test2"
}

func main() {
	for k, v := range myMap {
		println(fmt.Sprintf("key: %s, value: %s", k, v))
	}

	fmt.Println(myMap == nil)
}