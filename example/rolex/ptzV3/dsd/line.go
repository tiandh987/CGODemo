package dsd

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
)

type LineMarginOp int

const (
	SetLeftMargin LineMarginOp = iota
	SetRightMargin
	ClearLeftMargin
	ClearRightMargin
)

const MarginNoLimit = -1

const (
	MaxLineScanNum = 5
)

type LineScanID int

func (i LineScanID) Validate() error {
	if err := _validate.Var(i, "gte=1,lte=5"); err != nil {
		log.Error(err.Error())
		return errors.New("line scan id is invalid")
	}

	return nil
}

type LineScan struct {
	Enable             bool       `json:"Enable" validate:"boolean"`
	ID                 LineScanID `json:"ID" validate:"required,gte=1,lte=5"`
	LeftMargin         float64    `json:"LeftMargin" validate:"gte=-1,lt=360"`
	ResidenceTimeLeft  int        `json:"ResidenceTimeLeft" validate:"required,gte=1,lte=60"`
	RightMargin        float64    `json:"RightMargin" validate:"gte=-1,lt=360"`
	ResidenceTimeRight int        `json:"ResidenceTimeRight" validate:"required,gte=1,lte=60"`
	Running            bool       `json:"Runing" validate:"boolean"`
	Speed              int        `json:"Speed" validate:"required,gte=1,lte=8"`
}

func NewLine(id LineScanID) LineScan {
	return LineScan{
		Enable:             false,
		ID:                 id,
		LeftMargin:         -1,
		ResidenceTimeLeft:  5,
		RightMargin:        -1,
		ResidenceTimeRight: 5,
		Running:            false,
		Speed:              5,
	}
}

func (l *LineScan) Validate() error {
	if err := _validate.Struct(l); err != nil {
		return err
	}

	return nil
}

type LineSlice []LineScan

func NewLineSlice() LineSlice {
	s := make([]LineScan, MaxLineScanNum)

	for id := 1; id <= MaxLineScanNum; id++ {
		s[id-1] = NewLine(LineScanID(id))
	}

	return s
}

func (l *LineSlice) ConfigKey() string {
	return "LineScans"
}
