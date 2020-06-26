package main

import "fmt"

func doSomething() {
	status := 0
	
	defer cleanUpByStatus(&status)

	// Job Stuff...

	// Change status, error occurred maybe
	status = 2
}

func cleanUpByStatus(status *int) {
	// Do something by status
	fmt.Println(*status)
}

func main() {
	// 2
	doSomething()
}
