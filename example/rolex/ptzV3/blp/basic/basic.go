package basic

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

type Basic struct {
	ar ptz.AbilityRepo
}

func New(repo ptz.AbilityRepo) *Basic {
	return &Basic{ar: repo}
}

func (b *Basic) Operation(id int, speed ptz.Speed) error {
	log.Infof("id: %d, speed: %d", id, speed)

	switch Operation(id) {
	case DirectionUp:
		return b.ar.Up(speed)
	case DirectionDown:
		return b.ar.Down(speed)
	case DirectionLeft:
		return b.ar.Left(speed)
	case DirectionRight:
		return b.ar.Right(speed)
	case DirectionLeftUp:
		return b.ar.LeftUp(speed)
	case DirectionRightUp:
		return b.ar.RightUp(speed)
	case DirectionLeftDown:
		return b.ar.LeftDown(speed)
	case DirectionRightDown:
		return b.ar.RightDown(speed)
	case FocusFar:

	case FocusNear:

	case ZoomTele:
		return b.ar.ZoomAdd()
	case ZoomWide:
		return b.ar.ZoomSub()
	case IrisClose:

	case IrisOpen:

	default:
		log.Warnf("invalid ptz operation (%d)", id)
	}

	return nil
}

func (b *Basic) Version() string {
	version, err := b.ar.Version()
	if err != nil {
		log.Errorf(err.Error())
		version = ""
	}
	return version
}

func (b *Basic) Model() string {
	model, err := b.ar.Model()
	if err != nil {
		log.Errorf(err.Error())
		model = ""
	}
	return model
}

func (b *Basic) Restart() {
	if err := b.ar.Restart(); err != nil {
		log.Errorf(err.Error())
	}
}

func (b *Basic) Stop() {
	if err := b.ar.Stop(); err != nil {
		log.Errorf(err.Error())
	}
}

func (b *Basic) Position() (*dsd.Position, error) {
	return b.ar.Position()
}

func (b *Basic) Goto(pos *dsd.Position) error {
	return b.ar.Goto(pos)
}
