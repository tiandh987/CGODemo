package basic

import (
	"context"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"time"
)

type Basic struct {
	ar ptz.AbilityRepo
}

func New(repo ptz.AbilityRepo) *Basic {
	return &Basic{ar: repo}
}

func (b *Basic) Operation(id Operation, speed ptz.Speed) error {
	log.Debugf("id: %d, speed: %d", id, speed)

	switch id {
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

func (b *Basic) Stop() error {
	if err := b.ar.Stop(); err != nil {
		log.Errorf(err.Error())
		return err
	}

	return nil
}

func (b *Basic) Position() (*dsd.Position, error) {
	return b.ar.Position()
}

func (b *Basic) Goto(pos *dsd.Position) error {
	return b.ar.Goto(pos)
}

func (b *Basic) ReachPosition(ctx context.Context, dst *dsd.Position) error {
	ticker := time.NewTicker(time.Millisecond * 10)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			pos, err := b.Position()
			if err != nil {
				return err
			}

			if pos.Pan >= dst.Pan-2 && pos.Pan <= dst.Pan+2 &&
				pos.Tilt >= dst.Tilt-2 && pos.Tilt <= dst.Tilt+2 &&
				pos.Zoom >= dst.Zoom-2 && pos.Zoom <= dst.Zoom+2 {
				b.Stop()
				return nil
			}
		}
	}
}
