This branch shows how multiple defer works in the same function;

### 一个函数中的多个defer

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

### defer与闭包

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

### defer修改返回值

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

但是下面的两段代码却都没有修改返回值：

```go
func deferReturn1() int {
	a := 2
	defer func() {
		a = a + 2
	}()
	return a
}
// 控制台输出：2
```

```go
func add(i *int) {
	*i = *i + 2
}

func deferReturn2() int {
	a := 2
	defer add(&a)
	return a
}
// 控制台输出：2
```

原因是因为：

`return a`代码经过编译后，会被拆分为：

   1. 返回值 = a
   2. 调用 defer 函数
   3. return

所以：对于**未声明返回变量名的返回值**会：先放入一个匿名的返回值容器中，然后调用defer函数，最后再将匿名容器中的值返回；

则：通过defer修改返回值的方法为：

<font color="#f00">**通过修改返回值声明的值(显式返回值容器)来修改返回值，如app3中的代码；**</font>

### 循环中的defer

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

### defer与文件关闭

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

### Return前的Defer

通常情况下，Defer和Java等语言中finally关键字所做的事情类似：在正常或出现异常情况下都关闭资源等；

但是在go中defer是一个函数调用语句，这就意味着，如果执行不到defer语句，则在当前函数正常或异常return时，并不会执行defer函数，例如：

```go
package main

import "fmt"

func main() {
	deferTest()
}

func deferTest() {
	fmt.Println("Into deferTest...")
	return

	defer deferCallback()
	fmt.Println("Leave deferTest...")
}

func deferCallback() {
	fmt.Println("Into defer Callback")
}

// Output：
Into deferTest...
```

而如果将defer函数放在函数开头调用，则defer函数会在return时执行，例如：

```go
package main

import "fmt"

func main() {
	deferTest()
}

func deferTest() {
	defer deferCallback()
	fmt.Println("Into deferTest...")
	fmt.Println("Leave deferTest...")

	return
}

func deferCallback() {
	fmt.Println("Into defer Callback")
}
// Output：
Into deferTest...
Leave deferTest...
Into defer Callback
```

>   **注意：**
>
>   **defer是一条语句，在程序实际运行时才能被决定是否被执行**
>
>   **而try-catch-finally是一个语句块，在Java中是通过在class文件中创建的异常表来处理的，所以在finally声明的语句，无论如何都会被执行；**

所以在return前声明defer是没有意义的，如：

```go
package main

import "fmt"

func main() {
	deferTest()
}

func deferTest() {
	fmt.Println("Into deferTest...")
	fmt.Println("Leave deferTest...")
//	defer deferCallback()
    deferCallback()
	return
}

func deferCallback() {
	fmt.Println("Into defer Callback")
}
```

### defer函数传参

#### 使用场景

有时希望在当前函数执行后根据执行的结果执行一些操作；

为了保证函数在执行到中途时直接执行return返回，我们需要在程序开始时就声明defer函数，以确保在任何情况下都可以执行defer函数，如下：

```go
func doSomething() {
	status := 0
	
	defer cleanUpByStatus(status)

	// Job Stuff...

	// Change status, error occurred maybe
	status = 2
}

func cleanUpByStatus(status int) {
	// Do something by status
	fmt.Println(status)
}

func main() {
    // 0
	doSomething()
}
```

但是由于defer函数在**声明时就已经计算出了传函的值**，所以上面在声明了defer函数之后修改了status的之后，并**没有修改传入defer函数的值**；

一个简单的方法是：将defer声明的变量类型改为指针类型，这样在defer声明时传入的其实是变量的指针(地址值)，此后在修改变量的值时，地址下的值也会变化，例如：

```go
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
```

下面的说明展示了不同情况下传入defer的变量类型的场景；

#### 说明

① 非引用传参给`defer`调用的函数，且为非闭包函数，值`不会`受后面的改变影响

```go
func defer1() {
	a := 3  // a 作为演示的参数
	defer fmt.Println(a) // 非引用传参，非闭包函数中，a 的值 不会 受后面的改变影响
	a = a + 2
}
// 控制台输出 3
```

****

② 传递引用给`defer`调用的函数，即使不使用闭包函数，值也`会`受后面的改变影响

```go
func myPrintln(point *int)  {
	fmt.Println(*point) // 输出引用所指向的值
}
func defer2() {
	a := 3
	// &a 是 a 的引用。内存中的形式： 0x .... ---> 3
	defer myPrintln(&a) // 传递引用给函数，即使不使用闭包函数，值 会 受后面的改变影响
	a = a + 2
}
// 控制台输出 5
```

****

③ `defer`调用闭包函数，且内调用外部非传参进来的变量，值`会`受后面的改变影响

```go
// 闭包函数内，事实是该值的引用
func defer3() {
	a := 3
	defer func() {
        // 闭包函数内调用外部非传参进来的变量，事实是该值的引用，值 会 受后面的改变影响
		fmt.Println(a) 
	}()
	a = a + 2  // 3 + 2 = 5
}
// 控制台输出： 5
```

```go
// defer4 会抛出数组越界错误。
func defer4() {
	a := []int{1,2,3}
	for i:=0;i<len(a);i++ {
		// 同 defer3 的闭包形式。因为 i 是外部变量，没用通过传参的形式调用。在闭包内，是引用。
		// 值 会 受 ++ 改变影响。导致最终 i 是3， a[3] 越界
		defer func() {
			fmt.Println(a[i])
		}()
	}
}
// 结果：数组越界错误
```

****

④ `defer`调用闭包函数，若内部使用了传参参数的值。使用的是 值

```go
func defer5() {
	a := []int{1,2,3}
	for i:=0; i<len(a); i++ {
		// 闭包函数内部使用传参参数的值。内部的值为传参的值。同时引用是不同的
		defer func(index int) {
		        // index 有一个新地址指向它
			fmt.Println(a[index]) // index == i
		}(i)
		// 后进先出，3 2 1
	}
}
// 控制台输出： 
//     3
//     2
//     1
```

****

⑤ `defer`所调用的非闭包函数，参数如果是函数，会按顺序先执行（函数参数）

```go
func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}
func defer6()  {
	a := 1
	b := 2
	// calc 充当了函数中的函数参数。即使在 defer 的函数中，它作为函数参数，定义的时候也会首先调用函数进行求值
	// 按照正常的顺序，calc("10", a, b) 首先被调用求值。calc("122", a, b) 排第二被调用
	defer calc("1", a, calc("10", a, b))
	defer calc("12",a, calc("122", a, b))
}
// 控制台输出：
/**
10 1 2 3   // 第一个函数参数
122 1 2 3  // 第二个函数参数
12 1 3 4   // 倒数第一个 calc
1 1 3 4    // 倒数第二个 calc
*/
```
