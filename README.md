This branch shows some pits and falls in golang-method practicing.

### 方法Receiver的类型的意义

在Golang中定义的方法，在定义方法的Receiver时有两种情况：Receiver是指针还是非指针类型，例如：

```go
func (t T) Function() {/* func body */}
func (t *T) Function() {/* func body */}
```

而不同的Receiver类型，有着不同的行为；

通俗来讲，当Receiver是非指针类型时，对应的方法就是值传递；而当Receiver是指针类型时，对应的方法就是引用传递（其实传递的是指针的值）；

对于熟悉C/C++编程的同学来说应该很好理解；

而使用过Java的同学都知道，在Java中，八种基本的数据类型，如int，double在方法中都是值传递，而其他对象类型（包括Integer这种类对象）都是引用传递；

对于值传递来说，方法的调用对于方法调用来说是`无副作用`的，也就是说方法的调用会将整个对象复制一份，然后在方法调用时改变的是这个复制的值；

对于引用传递来说，方法的调用是会直接改变这个对象的属性的；

在app.go中的例子展示了这一点：

```go
package main

import "fmt"

type Pointer struct {
	x int64
	y int64
}

func (p Pointer) modifiedWithNonPointer() {
	p.x = 10000
	p.y = 10000
}

func (p *Pointer) modifiedWithPointer() {
	p.x = 10000
	p.y = 10000
}

/* 注意receiver是指针类型还是非指针类型 */
func main() {
	p1 := Pointer{x: 1, y: 1}
	p2 := Pointer{x: 2, y: 2}

	p1.modifiedWithNonPointer()
	fmt.Printf("modified with non-pointer receiver: %v\n", p1)

	p2.modifiedWithPointer()
	fmt.Printf("modified with pointer receiver: %v\n", p2)
}

/*
Output:
	modified with non-pointer receiver: {1 1}
	modified with pointer receiver: {10000 10000}
*/
```

对于值传递：

由于在每次调用时，会先将整个方法调用对象复制一份，当对象较大时，方法的执行效率会严重降低（而指针是固定大小的指向内存的一个数字）；并且值传递会浪费更多的内存空间；

但是值传递可以**很简单的保证线程安全**，对于同一个对象调用在两个线程中的操作并不会修改原对象的内容；这也是函数式编程的好处之一！

对于引用传递：

每一次在调用时，传递的是对象指针的值，而不会因为对象的大小而降低方法的执行效率；但是由于引用传递会导致多个方法或者单个方法在多线程调用时存在多个线程（goroutine）同时修改同一个对象，而造成并发问题，从而加大并发编程难度；

>   Golang相比于Java一刀切的引用传递，给了Coder更多的选择空间；但是到底使用哪种传递方式还需要根据特定的场景来自行判断；

### 方法Receiver的类型与调用方式

对于指针和非指针类型的Receiver，在进行调用时有着不同的形式，如下：

对于引用传递的方法：

```go
func (t *T) function() {/* body */}
```

下面的几种调用形式都是可以的：

```go
t := T{}
t.function()
(&t).function()

t := &T{}
t.function()
(*t).function()
```

而对于值传递的方法：

```go
func (t T) function() {/* body */}
```

下面的两种调用方式也都是可以的：

```go
t := T{}
t.function()
(&t).function()

t := &T{}
t.function()
(*t).function()
```

Golang会在编译阶段自动帮你加上对应的指针/取值符号；

但是在某些情况下，这种优化会导致一些困惑；例如在app2.go中展示的例子：

```go
package main

import (
	"bytes"
	"fmt"
)

type Rectangle struct {
	x, y int64
}

func (r *Rectangle) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	_, _ = fmt.Fprintf(&buf, "x: %d, y: %d", r.x, r.y)
	buf.WriteByte('}')

	return buf.String()
}

func main() {
	rectangle := Rectangle{x: 1, y: 1}
	fmt.Println(&rectangle)
	fmt.Println(rectangle.String())
	fmt.Println(rectangle)
}

/*
Output:
	{x: 1, y: 1}
	{x: 1, y: 1}
	{1 1}
*/
```

Println正确调用了`&rectangle`，`rectangle.String()`中的String方法，原因是：

对于`*Retangle`类型的指针，的确存在了绑定在这个指针类型上的String方法，所以会主动去调用这个方法；而`rectangle.String()`在编译期会被优化为`(&rectangle).String()`，所以也可以正常输出；

但是对于Rectangle类型，由于并没有绑定在其上的String方法，所以会输出对象的原始值；

>   并且Golang不支持任何意义上的重载！
>
>   所以也不要妄想定义一个值传递的String了！

