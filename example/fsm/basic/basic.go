package main

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

// FSM 是 Go 语言开发的有限状态机

func main() {
	fsm := fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{},
	)

	fmt.Println(fsm.Current())

	err := fsm.Event(context.Background(), "open")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm.Current())

	err = fsm.Event(context.Background(), "close")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm.Current())
}
