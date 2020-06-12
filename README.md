This branch shows some traps and pitfalls in golang-time practicing.

### 注意time包中的Unix方法

在time包中存在两个Unix方法：
```go
// 计算当前时间的时间戳
func (t Time) Unix() int64 {
    return t.unixSec()
}

// 根据时间戳求当前时间
func Unix(sec int64, nsec int64) Time {
	if nsec < 0 || nsec >= 1e9 {
		n := nsec / 1e9
		sec += n
		nsec -= n * 1e9
		if nsec < 0 {
			nsec += 1e9
			sec--
		}
	}
	return unixTime(sec, int32(nsec))
}
```

其中Time中的Unix方法用于求取当前time变量的时间戳(UTC+0)；

而Unix函数用于通过时间戳创建一个Time类型的时间对象，但是此对象是基于Local时间的，即：比如在中国，则求值之后会转换到UTC+8；

例如，下面的代码：

app.go

```GO
package main

import (
	"fmt"
	"time"
)

func main() {
	// timestamp
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().In(time.UTC).Unix())
	fmt.Println(time.Unix(time.Now().Unix(), 0).Unix())

	// Time to other location
	fmt.Println(time.Unix(time.Now().Unix(), 0))
	fmt.Println(time.Unix(time.Now().Unix(), 0).UTC())
}

// Output：
1591881833
1591881833
1591881833
2020-06-11 21:23:53 +0800 CST
2020-06-11 13:23:53 +0000 UTC
```

不管在什么地域，最后Unix()都是返回的同一个值；

而Unix(time.Now, 0)求出的结果是个地域相关的；

所以在处理国际化时：

**使用time.Unix()创建的时间最好先调用UTC转化为标准时间，然后再根据时区进行处理**

一个比较好的实践如下:

app2.go

```go
package main

import (
	"fmt"
	"time"
)

const (
	timeLayout = `2006-01-02 15:04:00`
)

// timezone timestamp offset
var timezoneConst = map[string]time.Duration{
	"UTC-12 IDLW":    -12 * 3600,
	"UTC-11 MIT":     -11 * 3600,
	"UTC-10 HST":     -10 * 3600,
	"UTC-9:30 MSIT":  -9.5 * 3600,
	"UTC-9 AKST":     -9 * 3600,
	"UTC-8 PST":      -8 * 3600,
	"UTC-7 MST":      -7 * 3600,
	"UTC-6 CST":      -6 * 3600,
	"UTC-5 EST":      -5 * 3600,
	"UTC-4 AST":      -4 * 3600,
	"UTC-3:30 NST":   -3.5 * 3600,
	"UTC-3 SAT":      -3 * 3600,
	"UTC-2":          -2 * 3600,
	"UTC-1 CVT":      -1 * 3600,
	"UTC":            0,
	"UTC+1 CET":      3600,
	"UTC+2 EET":      2 * 3600,
	"UTC+3 MSK":      3 * 3600,
	"UTC+3:30 IRT":   3.5 * 3600,
	"UTC+4 META":     4 * 3600,
	"UTC+4:30 AFT":   4.5 * 3600,
	"UTC+5 METB":     5 * 3600,
	"UTC+5:30 IDT":   5.5 * 3600,
	"UTC+6 BHT":      6 * 3600,
	"UTC+6:30 MRT":   6.5 * 3600,
	"UTC+7 IST":      7 * 3600,
	"UTC+8 EAT":      8 * 3600,
	"UTC+9 FET":      9 * 3600,
	"UTC+9:30 ACST":  9.5 * 3600,
	"UTC+10 AEST":    10 * 3600,
	"UTC+10:30 FAST": 10.5 * 3600,
	"UTC+11 VIT":     11 * 3600,
	"UTC+11:30 NFT":  11.5 * 3600,
	"UTC+12 PSTB":    12 * 3600,
	"UTC+12:45 CIT":  12.75 * 3600,
	"UTC+13 PSTC":    13 * 3600,
	"UTC+14 PSTD":    14 * 3600,
}

func main() {
	fmt.Println(time.Now().UTC().Add(timezoneConst["UTC+14 PSTD"] * time.Second))
	fmt.Println(time.Now().UTC().Add(timezoneConst["UTC+8 EAT"] * time.Second))
}
```
