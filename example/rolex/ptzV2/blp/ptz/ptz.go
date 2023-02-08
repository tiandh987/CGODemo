package ptz

import (
	"errors"
)

// Speed 云台速度
type Speed int

func (s Speed) Convert() byte {
	return _ptSpeedMap[s]
}

func (s Speed) Validate() error {
	if s < SpeedOne || s > SpeedEight {
		return errors.New("speed is invalid")
	}

	return nil
}

const (
	SpeedOne Speed = iota + 1
	SpeedTwo
	SpeedThree
	SpeedFour
	SpeedFive
	SpeedSix
	SpeedSeven
	SpeedEight
)

var _ptSpeedMap = map[Speed]byte{
	SpeedOne:   byte(0x01),
	SpeedTwo:   byte(0x09),
	SpeedThree: byte(0x12),
	SpeedFour:  byte(0x1b),
	SpeedFive:  byte(0x24),
	SpeedSix:   byte(0x2d),
	SpeedSeven: byte(0x36),
	SpeedEight: byte(0x3f),
}

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
