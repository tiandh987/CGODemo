package basic

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
)

type Basic struct {
}

func New() *Basic {
	return &Basic{}
}

func (b *Basic) Start(ctl control.ControlRepo, id int, speed ptz.Speed) error {
	log.Infof("id: %d, speed: %d", id, speed)

	switch ptz.Operation(id) {
	case ptz.DirectionUp:
		return ctl.Up(speed.Convert())
	case ptz.DirectionDown:
		return ctl.Down(speed.Convert())
	case ptz.DirectionLeft:
		return ctl.Left(speed.Convert())
	case ptz.DirectionRight:
		return ctl.Right(speed.Convert())
	case ptz.DirectionLeftUp:
		return ctl.LeftUp(speed.Convert())
	case ptz.DirectionRightUp:
		return ctl.RightUp(speed.Convert())
	case ptz.DirectionLeftDown:
		return ctl.LeftDown(speed.Convert())
	case ptz.DirectionRightDown:
		return ctl.RightDown(speed.Convert())
	case ptz.FocusFar:

	case ptz.FocusNear:

	case ptz.ZoomTele:
		return ctl.ZoomAdd()
	case ptz.ZoomWide:
		return ctl.ZoomSub()
	case ptz.IrisClose:

	case ptz.IrisOpen:
	}

	return nil
}

func (b *Basic) Stop() error {
	return nil
}
