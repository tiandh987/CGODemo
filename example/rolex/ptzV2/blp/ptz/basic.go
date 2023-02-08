package ptz

import "github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"

type Basic struct {
}

func NewBasic() *Basic {
	return &Basic{}
}

func (b *Basic) Start(ctl control.ControlRepo, id int, speed Speed) error {
	switch Operation(id) {
	case DirectionUp:
		return ctl.Up(speed.Convert())
	case DirectionDown:
		return ctl.Down(speed.Convert())
	case DirectionLeft:
		return ctl.Left(speed.Convert())
	case DirectionRight:
		return ctl.Right(speed.Convert())
	case DirectionLeftUp:
		return ctl.LeftUp(speed.Convert())
	case DirectionRightUp:
		return ctl.RightUp(speed.Convert())
	case DirectionLeftDown:
		return ctl.LeftDown(speed.Convert())
	case DirectionRightDown:
		return ctl.RightDown(speed.Convert())
	case FocusFar:

	case FocusNear:

	case ZoomTele:
		return ctl.ZoomAdd()
	case ZoomWide:
		return ctl.ZoomSub()
	case IrisClose:

	case IrisOpen:
	}

	return nil
}

func (b *Basic) Stop() error {
	return nil
}
