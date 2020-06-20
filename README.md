This branch shows some traps and pitfalls in golang-slice practicing.

### slice在Append时的自动扩容机制

在对slice进行append操作时，当slice的元素个数超过当前容量(capacity)的一半，则slice会扩容一倍；

我们可以通过len和cap两个内置函数分别获取slice的元素个数和容量；

但是要注意的是：

**当使用make函数创建slice时指定的初始size大于0时，使用append后，slice中的对象可能不仅仅包括append的对象，还包括多个空对象；**

**甚至使用len求得的元素个数也和调用append加入的元素个数不等(内部含有多个空元素)**

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

所以正确的指定容量初始化一个空的slice应当为`make([]Type, 0, initCap)`;

