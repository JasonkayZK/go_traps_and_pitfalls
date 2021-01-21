package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	x := 1
	// 隐式传参调用
	// 此时传入的是x的“引用值”，即两个x指向的是同一个内存地址，在子routine中修改的值，会改变外部的x！
	go func() {
		fmt.Printf("Implicit invoke: %d\n", x)
		wg.Done()
	}()

	// 直接传参调用
	// 此时为值传递，内部的x不会影响外部的x；
	go func(x int) {
		fmt.Printf("Direct invoke: %d\n", x)
		wg.Done()
	}(x)

	x = 3

	wg.Wait()
}
