package main

import "fmt"

func main() {
	level := 51
	fmt.Printf("float32: %+v\n", float32(level))
	fmt.Printf("float32 * 0.72: %+v\n", float32(level)*0.72)
	fmt.Printf("int32(float32 * 0.72): %+v\n", int32(float32(level)*0.72))
}
