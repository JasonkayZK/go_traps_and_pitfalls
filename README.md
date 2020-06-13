This branch shows some traps and pitfalls in go-channel practicing.

### 在一个关闭的channel上进行读写

**在一个关闭的channel中进行写操作将会造成panic**，如下：

```go
func main() {
	intChan := make(chan int, 5)
	close(intChan)

	go func() {
		intChan <- 1
	}()

	fmt.Println(<-intChan)
}
// Output：
0
panic: send on closed channel

goroutine 6 [running]:
main.main.func1(0xc00004a060)
	D:/workspace/go_traps_and_pitfalls/app.go:10 +0x3e
created by main.main
	D:/workspace/go_traps_and_pitfalls/app.go:9 +0x76
```

**在一个关闭的channel中进行读操作则可以继续读出**，如下：

```go
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
// Output：
channel closed
1
2
3
4
5
```

