package main

import "fmt"

type Config struct {
	A int
	B int
}

func NewConfig() *Config {
	return &Config{
		A: 1,
		B: 2,
	}
}

func main() {
	a := NewConfig()
	fmt.Printf("a: %+v\n", a)
	b := NewConfig()
	fmt.Printf("b: %+v\n", b)

	a.A = 3
	fmt.Printf("after a: %+v\n", a)
	fmt.Printf("after b: %+v\n", b)
}
