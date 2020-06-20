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

