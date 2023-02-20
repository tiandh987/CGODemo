package main

import (
	"context"
	"fmt"
	"github.com/looplab/fsm"
)

func main() {
	newFSM := fsm.NewFSM(
		"green",
		fsm.Events{
			{Name: "warn", Src: []string{"green"}, Dst: "yellow"},
			{Name: "panic", Src: []string{"yellow"}, Dst: "red"},
			{Name: "panic", Src: []string{"green"}, Dst: "red"},
			{Name: "calm", Src: []string{"red"}, Dst: "yellow"},
			{Name: "clear", Src: []string{"yellow"}, Dst: "green"},
		},
		fsm.Callbacks{
			"before_warn": func(_ context.Context, e *fsm.Event) {
				fmt.Printf("before_warn, args: %+v\n", e.Args[0])
			},
			"before_event": func(_ context.Context, e *fsm.Event) {
				fmt.Printf("before_event, args: %+v\n", e.Args[0])
			},
			"leave_green": func(_ context.Context, e *fsm.Event) {
				fmt.Printf("leave_green, args: %+v\n", e.Args[0])
			},
			"leave_state": func(_ context.Context, e *fsm.Event) {
				fmt.Printf("leave_state, args: %+v\n", e.Args[0])
			},
			"enter_yellow": func(_ context.Context, e *fsm.Event) {
				fmt.Printf("enter_yellow, args: %+v\n", e.Args[0])
			},
			"enter_state": func(_ context.Context, e *fsm.Event) {
				fmt.Printf("enter_state, args: %+v\n", e.Args[0])
			},
			"after_warn": func(_ context.Context, e *fsm.Event) {
				fmt.Printf("after_warn, args: %+v\n", e.Args[0])
			},
			"after_event": func(_ context.Context, e *fsm.Event) {
				fmt.Printf("after_event, args: %+v\n", e.Args[0])
			},
		},
	)

	fmt.Println(newFSM.Current())

	b := aaa{
		num:  15,
		name: "aaaname",
	}

	err := newFSM.Event(context.Background(), "warn", b)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(newFSM.Current())

}

type aaa struct {
	num  int
	name string
}
