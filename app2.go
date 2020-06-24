package main

import "fmt"

// 函数 copy 在两个 slice 间复制数据，复制⻓度以 len 小的为准，直接对应位置覆盖
// 注:
// 1.不同类型的切片无法复制
// 2.如果s1的长度大于s2的长度，将s2中对应位置上的值替换s1中对应位置的值
// 3.如果s1的长度小于s2的长度，多余的将不做替换
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
