package main

import (
	"fmt"
	"github.com/tiandh987/CGODemo/example/archtest/arch/kx"
	"github.com/tiandh987/CGODemo/example/archtest/blp"
)

func main() {
	fmt.Println("hello cgo")

	ar := kx.New()
	//ar := nt.New()

	audioBlp := blp.New(ar)

	play := audioBlp.Play()

	fmt.Printf("play: %d\n", play)
}
