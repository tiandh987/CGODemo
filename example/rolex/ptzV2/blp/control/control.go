package control

import "github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"

type ControlRepo interface {
	Version() (string, error)
	Model() (string, error)
	Restart() error
	Stop() error
	Up(speed byte) error
	Down(speed byte) error
	Left(speed byte) error
	Right(speed byte) error
	LeftUp(speed byte) error
	RightUp(speed byte) error
	LeftDown(speed byte) error
	RightDown(speed byte) error
	ZoomAdd() error
	ZoomSub() error
	Position() (*dsd.Position, error)
	Goto(*dsd.Position) error
}
