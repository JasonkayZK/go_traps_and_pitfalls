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

type PointerSlice []*Person

func (ps PointerSlice) String() string {
	str := "["
	for _, p := range ps {
		str = str + `{` + strconv.FormatInt(p.Age, 10) + " " + `[` + strings.Join(p.Hobby, ",") + `]}` + " "
	}
	return strings.TrimSpace(str) + "]"
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

	ps1 := PointerSlice{
		{Age: 18, Hobby: []string{"a", "b", "c"}},
	}
	ps2 := make(PointerSlice, len(ps1))
	copy(ps2, ps1)

	ps1[0].Age = 8
	ps1[0].Hobby[0] = "z"
	fmt.Println(ps1)
	fmt.Println(ps2)
}
