package dsd

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
)

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
	Enable             bool       `json:"Enable" validate:"boolean"`                           // 使能
	ID                 LineScanID `json:"ID" validate:"required,gte=1,lte=5"`                  // 线性扫描id
	LeftMargin         float64    `json:"LeftMargin" validate:"gte=-1,lte=360"`                // 设置左边界，-1 清除边界
	ResidenceTimeLeft  int        `json:"ResidenceTimeLeft" validate:"required,gte=1,lte=60"`  // 停留时间（左）秒
	RightMargin        float64    `json:"RightMargin" validate:"gte=-1,lte=360"`               // 设置右边界，-1 清除边界
	ResidenceTimeRight int        `json:"ResidenceTimeRight" validate:"required,gte=1,lte=60"` // 停留时间（右）秒
	Runing             bool       `json:"Runing" validate:"boolean"`                           // 运行状态
	Speed              int        `json:"Speed" validate:"required,gte=1,lte=8"`               // 线扫速度 1-8 1:最慢 8:最快
}

func (l *LineScan) Validate() error {
	if err := _validate.Struct(l); err != nil {
		return err
	}

	return nil
}

func (l *LineScan) Default() []LineScan {
	s := make([]LineScan, MaxLineScanNum)

	for id := 1; id <= MaxLineScanNum; id++ {
		l, _ := NewLineScan(LineScanID(id))
		s[id-1] = l
	}

	return s
}

func (l *LineScan) ConfigKey() string {
	return "LineScans"
}

func NewLineScan(id LineScanID) (LineScan, error) {
	log.Debugf("NewLineScan param id: %d", id)

	if err := id.Validate(); err != nil {
		return LineScan{}, err
	}

	return LineScan{
		Enable:             false,
		ID:                 id,
		LeftMargin:         -1,
		ResidenceTimeLeft:  5,
		RightMargin:        -1,
		ResidenceTimeRight: 5,
		Runing:             false,
		Speed:              5,
	}, nil
}
