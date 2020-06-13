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

