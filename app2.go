package main

import (
	"fmt"
	"log"
	"time"
)

/* 通过defer记录函数运行时间 */
func main() {
	for i := 0; i < 5; i++ {
		bigSlowOperation()
	}
}

/*
	对于defer trace("bigSlowOperation"):
		表示在结束的时候调用trace, 此时才刚刚初始化trace;

	对于defer trace("bigSlowOperation")():
		表示在结束时调用trace函数返回的函数;
		此时trace为一个闭包, 并且会在进入bigSlowOperation函数时就调用trace, 以初始化start;
		defer真正获得的是trace返回的func()
 */
func bigSlowOperation() {
	//defer trace("bigSlowOperation")
	defer trace("bigSlowOperation")()
	time.Sleep(1 * time.Second)
}

/*
	构成一个闭包;
	进入bigSlowOperation时初始化start, 并且返回一个func供defer调用
 */
func trace(msg string) func() {
	start := time.Now()
	fmt.Printf("enter %s \n\n", msg)

	return func() {
		log.Printf("exit %s (%s)\n", msg, time.Since(start))
	}
}

/*
Output:
	enter bigSlowOperation
	2020/05/24 14:54:44 exit bigSlowOperation (1.0015546s)

	enter bigSlowOperation
	2020/05/24 14:54:45 exit bigSlowOperation (1.0001377s)

	enter bigSlowOperation
	2020/05/24 14:54:46 exit bigSlowOperation (1.0008242s)

	enter bigSlowOperation
	2020/05/24 14:54:47 exit bigSlowOperation (1.0012377s)

	enter bigSlowOperation
	2020/05/24 14:54:48 exit bigSlowOperation (1.0003853s)
*/
