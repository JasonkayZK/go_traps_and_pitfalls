package main

import (
	"fmt"
	"time"
)

func main() {
	// timestamp
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().In(time.UTC).Unix())
	fmt.Println(time.Unix(time.Now().Unix(), 0).Unix())

	// Time to other location
	fmt.Println(time.Unix(time.Now().Unix(), 0))
	fmt.Println(time.Unix(time.Now().Unix(), 0).UTC())
}