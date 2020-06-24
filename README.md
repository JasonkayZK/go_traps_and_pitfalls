This branch shows some traps and pitfalls in golang-slice practicing.

### slice初始化

在对slice进行初始化操作时，要注意在使用make函数创建时传入的len数量将会在slice中添加len个空的元素！

以下面的代码为例：

```go
func main() {
	slice1 := make([]string, 0)

	slice1 = append(slice1, "a")
	slice1 = append(slice1, "b")
	slice1 = append(slice1, "c")

	fmt.Println(slice1)
	fmt.Println(len(slice1))
	fmt.Println(cap(slice1))

	slice2 := make([]string, 3)
	slice2 = append(slice2, "a")
	slice2 = append(slice2, "b")
	slice2 = append(slice2, "c")

	fmt.Println(slice2)
	fmt.Println(len(slice2))
	fmt.Println(cap(slice2))

	slice3 := make([]string, 3, 3)
	slice3 = append(slice3, "a")
	slice3 = append(slice3, "b")
	slice3 = append(slice3, "c")

	fmt.Println(slice3)
	fmt.Println(len(slice3))
	fmt.Println(cap(slice3))

	slice4 := make([]string, 0, 3)
	slice4 = append(slice4, "a")
	slice4 = append(slice4, "b")
	slice4 = append(slice4, "c")

	fmt.Println(slice4)
	fmt.Println(len(slice4))
	fmt.Println(cap(slice4))
}
// Output：
[a b c]
3
4
[   a b c]
6
6
[   a b c]
6
6
[a b c]
3
3
```

原因在于当在make(T type, len int, cap len)中指定初始化大小len时，会同时创建len个空的元素进入slice中；

所以**正确的指定容量初始化一个空的slice应当为`make([]Type, 0, initCap)`;**

### 使用copy函数复制

golang原生提供了copy函数用于数组、切片等数据结构进行数据复制；

#### copy复制规则

copy函数在两个slice间**复制数据**，**复制⻓度以len小的为准，直接对应索引位置覆盖**

即copy函数满足以下规则：

1.  **不同类型的切片无法复制**
2.  如果s1的长度大于s2的长度，将s2中**对应位置上的值替换**s1中对应位置的值
3.  如果s1的长度小于s2的长度，**多余的将不做替换**

以下面的代码为例：

```go
func main() {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5}
	s3 := []int{6, 7, 8, 9}
	copy(s1, s2)
	// [4 5 3]
	fmt.Println(s1)

	// [6 7]
	copy(s2, s3)
	fmt.Println(s2)
}
```

都无法满足完全将一个切片复制到另一个切片；

如有这种需求可以采用类似于下面的代码：

```go
s1 := []int{1, 2, 3}
s2 := make([]int, len(s1))
copy(s2, s1)
```

#### copy函数的浅复制

go中的内置函数实现的是浅复制，就是说对于所有的数据进行的是值的复制(指针类型复制的是地址值)；

所以：

-   **对于非指针数据类型(int64而非*int64)：将直接拷贝其值；**
-   **对于指针类型(切片、数组、map以及其他指针等)：拷贝的是地址信息；**

则对于下面的例子：copy之后修改值类型将会使两个切片的值不一致，而修改指针类型后两个切片的数据都发生的改变：

```go
package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Person struct {
	Age   int64
	Hobby []string
}

func main() {
	s1 := []Person{
		{Age: 18, Hobby: []string{"a", "b", "c"}},
	}
	s2 := make([]Person, len(s1))
	copy(s2, s1)

	s1[0].Age = 8
	s1[0].Hobby[0] = "z"
	fmt.Println(s1)
	fmt.Println(s2)
}
// Output：
[{8 [z b c]}]
[{18 [z b c]}]
```

#### 为什么使用copy函数复制

既然copy函数有这么多的坑，为什么不直接使用for循环进行逐个复制？

由于copy是builtin函数，所以使用内置的copy函数进行复制时的效率会高于使用for循环逐个复制的效率，尤其是当复制的元素较多时；

