package blp

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
)

type Operation int

func (o Operation) ValidateDirection() error {
	if o < DirectionUp || o > DirectionRightDown {
		return errors.New("direction is invalid")
	}

	return nil
}

func (o Operation) ValidateOperation() error {
	if o < FocusFar || o > IrisOpen {
		return errors.New("operation is invalid")
	}

	return nil
}

// 云台操作
const (
	DirectionUp        Operation = iota // 上
	DirectionDown                       // 下
	DirectionLeft                       // 左
	DirectionRight                      // 右
	DirectionLeftUp                     // 左上
	DirectionRightUp                    // 右上
	DirectionLeftDown                   // 左下
	DirectionRightDown                  // 右下
	FocusFar                            // 焦距拉远
	FocusNear                           // 焦距拉进
	ZoomTele                            // 视角变窄
	ZoomWide                            // 视角变宽
	IrisClose                           // 光圈关闭
	IrisOpen                            // 光圈打开
)

func (b *Blp) turn(id int, speed ptz.Speed) error {
	log.Infof("id: %d, speed: %d", id, speed)

	ctl := b.getControl()
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
