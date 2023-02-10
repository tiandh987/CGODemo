package main

import "fmt"

func main() {
	s := make([]int, 10)

	for i := 0; i < 10; i++ {
		s[i] = i
	}

	fmt.Printf("s: %+v\n\n", s)

	tmp := &s[1]
	tmp1 := s[1]
	fmt.Printf("s: %p\ntmp: %p\ntmp1: %p\n\n", s, tmp, &tmp1)
	fmt.Printf("s: %d\ntmp: %d\ntmp1: %d\n\n", s, *tmp, tmp1)

	*tmp = 10
	fmt.Printf("s: %+v\n\n", s)

}
