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
