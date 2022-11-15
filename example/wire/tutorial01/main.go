package main

import (
	"fmt"
)

// Message is what greeters will use to greet guests.
type Message string

func NewMessage() Message {
	return "Hi there!"
}

// Greeter is the type charged with greeting guests.
type Greeter struct {
	Message Message
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

// Event is a gathering with greeters.
type Event struct {
	Greeter Greeter
}

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	e := InitializeEvent()

	e.Start()
}
