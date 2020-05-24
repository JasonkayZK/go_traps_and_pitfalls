This branch shows how multiple defer works in the same function;

### app1

在app1.go中展示了，当一个函数中含有多个被声明为defer的函数时，在当前函数结束（正常通过return/异常通过panic）时，**defer的调用顺序为defer声明的反方向；**

```go
package main

import "fmt"

/* A simple demo for showing how multiple defer works in a single function */
func main() {
	defer func() {
		fmt.Println("First defer declared!")
	}()

	defer func() {
		fmt.Println("Second defer declared!")
	}()

	func() {
		fmt.Println("A function declared in main!")
	}()
}

/*
Output:
	A function declared in main!
	Second defer declared!
	First defer declared!
*/
```

>注意：defer反向调用顺序的陷阱

### app2

app2.go展示了在defer中获取一个闭包函数，使用该闭包函数记录函数的运行时间：

```go
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
```

>   注意：在defer一个闭包时，不要忘记最后的括号

### app3

app3展示了在defer中也可以修改返回值：

```go
package main

import "fmt"

/* defer modify the return value */
func main() {
	// 12
	fmt.Println(func(x int) (result int) {
		defer func() {result += x}()
		return x * 2
	}(4))
}
```

### app4

app4展示了在一个循环中defer函数的坑；

defer声明的方法会在循环结束后被反方向重新调用，而并非在每一次循环中调用，这就导致了在一个循环中对文件等资源的处理关闭不及时，导致资源枯竭；

一个比较好的解决方法是在循环中调用另外一个函数，并在另一个函数中使用defer；

```go
package main

import "os"

/* Fault: */
/* defer in circulation, run out of file descriptors */
/*
func main() {
	filenames := os.Args[1:]
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		// all file will close by the end of circulation!
		defer f.Close()
		// process of each file
	}
}
*/

/* Solution 1: extract file operation to another function */
func main() {
	filenames := os.Args[1:]
	for _, filename := range filenames {
		if err := func(filename string) error {
			f, err := os.Open(filename)
			if err != nil {
				return err
			}
			defer f.Close()
			// process f
			// ...
			return nil
		}(filename); err != nil {
			panic(err)
		}
	}
}
```

### app5

在示例app5中，通过os.Create打开了一个文件进行写入，但是在关闭时，没有对f.Close采用defer机制，原因是：这可能会产生一些微妙的错误；

对于许多文件系统，尤其是NFS，写入文件时发生的错误会被延迟到关闭文件时才反馈；如果没有检查文件关闭时的反馈信息，可能会导致数据的丢失，并且操作被误以为写入成功！

所以当io.Copy和f.Close都失败了，在处理时更倾向于将io.Copy的错误信息反馈给调用者；

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

/* defer in Copy file */
func main() {
	filename, n, _ := fetch(os.Args[1])
	fmt.Printf("download file: %s, size: %d bit\n", filename, n)
}

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	// url末尾作为文件名
	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}

	// 文件大小
	n, err = io.Copy(f, resp.Body)

	// 不使用defer关闭文件
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}

	return local, n ,err
}
```





