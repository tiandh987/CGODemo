package basic

import "errors"

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
