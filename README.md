This branch shows some traps and pitfalls in golang-init practicing.

### 通过init初始化全局变量

init方法会在包导入之后首先被执行，做一些初始化的工作，例如初始化logger，初始化对象等；

在初始化对象时，要注意：

**防止init创建的是局部变量，从而导致全局变量实际为空(`:=`和`=`的区别)**

例如下面的代码：

```go
var myMap map[string]string

func init() {
	myMap := make(map[string]string)
	myMap["1"] = "test1"
	myMap["2"] = "test2"
}

func main() {
	for k, v := range myMap {
		println(fmt.Sprintf("key: %s, value: %s", k, v))
	}
	// true
	fmt.Println(myMap == nil)
}
```

在实际执行时，myMap是nil；

这是由于，**在init函数中创建的myMap，实际上是一个局部变量！(使用`:=`创建的新的myMap对象覆盖了全局变量myMap)**

所以init函数返回后，全局变量myMap依然是nil；

将init中的初始化`:=`改为`=`则正常运行；

