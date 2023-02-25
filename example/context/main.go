package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Cruise struct {
	ctx    context.Context
	cancel context.CancelFunc
	stCh   chan struct{}
	nextCh chan struct{}
	wg     sync.WaitGroup
}

func main() {
	cruise := &Cruise{}
	cruise.ctx, cruise.cancel = context.WithCancel(context.Background())

	cruise.stCh = make(chan struct{}, 1)
	cruise.nextCh = make(chan struct{}, 1)

	cruise.wg.Add(1)
	go func() {
		defer cruise.wg.Done()
		defer fmt.Println("xxxxxxxxxxxxxxxx")

		i := 0
		for {
			fmt.Printf("i: %d\n", i)
			i++

			select {
			case <-cruise.ctx.Done():
				fmt.Printf("return first goroutine, %s\n", cruise.ctx.Err())
				return
			case <-cruise.nextCh:
				fmt.Println("nextCh receive")

				cruise.stCh <- struct{}{}
			case <-cruise.stCh:
				fmt.Println("stCh receive")

				cruise.second(cruise.ctx)

				fmt.Println("stCh called receive")
			}
		}
	}()

	cruise.stCh <- struct{}{}

	time.Sleep(time.Second * 14)

	fmt.Println("cancel")
	cruise.cancel()

	cruise.wg.Wait()

	fmt.Println("end")
}

func (c *Cruise) second(ctx context.Context) {
	c.wg.Add(1)

	go func() {
		defer fmt.Println("ssssssssssssssssssss")
		defer c.wg.Done()

		timer := time.NewTimer(time.Second * 3)
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("return second goroutine, %s\n", ctx.Err())
				return
			case <-timer.C:
				fmt.Println("timer return second goroutine")
				c.nextCh <- struct{}{}
				return
			default:
				fmt.Println("second second second")
				time.Sleep(time.Second)
			}
		}
	}()

}
