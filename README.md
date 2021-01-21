## Goroutine

This branch shows some traps and pitfalls in goroutine practicing.

### 确保goroutine执行完成后退出

在初学goroutine时，一个go关键字就能实现并发很是神奇，但是还是要关注goroutine的生命周期；

和Java不同的是，go中goroutine的生命周期不会在main函数执行结束后继续存在，例如下面的代码，大概率是不会输出文字的，这是因为：goroutine还没来得及执行完方法，main函数就退出了；

```go
func main() {
	go printHi()
}

func printHi() {
	fmt.Println("hi")
}
```

>   注意这一点和Java有着不同：
>
>   在Java中，线程是由JVM管理的，所以JVM会在最后一个非守护线程退出后才会去主动停止程序；

所以，在go中就需要主动去保证在所有goroutine执行完毕后main函数才退出；

下面提供几种常见的方法：

**① 基于channal信号的同步退出：**

在创建chan时，如果未指定缓冲区的大小，则chan会处于阻塞状态，即：无写入时读阻塞，无读出时写阻塞；

所以可以通过一个空的struct{}类型来实现goroutine同步，如下：

```go
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
```

如果多个goroutine都需要同步，则需要多个endFlag同步；

但这是一种比较粗暴的方法，go中提供了类似于Java中Semaphore的处理方式：sync.WaitGroup

**② 通过sync.WaitGroup同步多个channel**

```go
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
```


### Goroutine显式调用与传参调用

#### **Goroutine调用概述**

对于go关键字创建新goroutine并调用函数的方式有两种：

-   **隐式传参调用**
-   **显式传参调用**

```go
func main() {
	x := 1
	// 隐式传参调用
	// 此时传入的是x的“引用值”，即两个x指向的是同一个内存地址，在子routine中修改的值，会改变外部的x！
	go func() {
		fmt.Println(x)
	}()

	// 直接传参调用
	// 此时为值传递，内部的x不会影响外部的x；
	go func(x int) {
		fmt.Println(x)
	}(x)
}
```

两者的区别在于：

-   当隐式传参时：此时传入的是x的“引用值”，即两个x指向的是同一个内存地址，在子routine中修改的值，会改变外部的x！
-   当显式传参时：此时为值传递，内部的x不会影响外部的x；

另外，需要注意的：<font color="#f00">**显式的传参，在传参时就必须将参数计算好，这一点和defer函数是相同的！**</font>

例如：

app3.go

```go
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
```

上面的函数大概率输出为：

```
Direct invoke: 1
Implicit invoke: 3
```

这是因为，在直接传参调用时，x的值还未被修改(仍然是1)并且已经被确定，而隐式传参调用会根据外部x值的改变而改变；

>   **之所以说是`大概率`是因为，一般情况下，隐式传参调用的goroutine执行速度还是比main中执行至`x=3`语句要慢的，所以，大概率会先执行`x=3`修改x的值，随后才会执行隐式传参调用！**

#### **一道关于 Goroutine 的题**

下面的代码输出什么呢？

app4.go

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan int)
    go fmt.Println(<-ch1)
    ch1 <- 5
    time.Sleep(1 * time.Second)
}
```

以上代码输出什么？（单选）

-   A：5
-   B：不能编译
-   C：运行时死锁

如果你耐心看了上面的讲解，可以很容易知道正确答案是：C；

因为：

在上方创建Goroutine进行调用时，实际上是**显式传参！**

所以，上方的代码其实类似于：

```go
func main() {
    ch1 := make(chan int)
    x := <-ch1
    go fmt.Println(x)
    ch1 <- 5
    time.Sleep(1 * time.Second)
}
```

此时`x := <-ch1`会阻塞main函数，而`ch1 <- 5`也是在main函数中调用的，所以会被阻塞，最终造成死锁！

>   Goroutine题目来源：
>
>   -   [Go语言爱好者周刊：第 78 期 — 这道关于 goroutine 的题](https://mp.weixin.qq.com/s/kma8hvdLVPIkZnKw_MaSKg)

