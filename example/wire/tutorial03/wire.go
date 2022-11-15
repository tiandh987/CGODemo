//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

func InitializeEvent(prase string) (Event, error) {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}, nil
}
