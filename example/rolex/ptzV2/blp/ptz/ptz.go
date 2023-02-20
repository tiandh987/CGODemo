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
