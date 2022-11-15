package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Printf("main start...\n")

	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	defer cancelFunc()

	err := A(timeoutCtx)
	if err != nil {
		fmt.Printf("main call A err: %s\n", err.Error())

	}

	fmt.Printf("main end\n")
}

func A(ctx context.Context) error {
	fmt.Printf("A start...\n")

	ch := make(chan bool, 1)

	go func() {
		ok := B()
		ch <- ok
	}()

	select {
	case <-ctx.Done():
		fmt.Printf("A Done, err: %s\n", ctx.Err().Error())
		return ctx.Err()
	case ok := <-ch:
		if !ok {
			return fmt.Errorf("call b error")
		}
	}

	fmt.Printf("A end\n")

	return nil
}

func B() bool {
	fmt.Printf("B start...\n")

	for i := 0; i < 30; i++ {
		if i > 10 {
			return false
		}
		fmt.Printf("B content: %d\n", i)
		time.Sleep(time.Second)
	}

	fmt.Printf("B end\n")

	return true
}
