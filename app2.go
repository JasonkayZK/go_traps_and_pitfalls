package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	fmt.Println("main start")

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sleepTime := rand.Intn(5)
			fmt.Println(fmt.Sprintf("work %d: %d seconds", i, sleepTime))
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}(i)
	}

	wg.Wait()

	fmt.Println("main end")
}

