package main

import "fmt"

func main() {
	endFlag := make(chan struct{})
	defer close(endFlag)
	go printHello(endFlag)
	<-endFlag
}

func printHello(flag chan struct{}) {
	fmt.Println("ok")
	flag <- struct{}{}
}