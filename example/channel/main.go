package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 2)

	ch <- "test"

	close(ch)

	time.Sleep(time.Second * 5)

	str := <-ch

	fmt.Printf("str: %s\n", str)

	//ch <- "test2"
	//close(ch)

}
