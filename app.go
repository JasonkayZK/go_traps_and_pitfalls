package main

import (
	"fmt"
	"time"
)

func main() {
	intChan := make(chan int, 5)

	go func() {
		intChan <- 1
		intChan <- 2
		intChan <- 3
		intChan <- 4
		intChan <- 5
		close(intChan)
		fmt.Println("channel closed")
	}()

	time.Sleep(3 * time.Second)
	for item := range intChan {
		fmt.Println(item)
	}
}
