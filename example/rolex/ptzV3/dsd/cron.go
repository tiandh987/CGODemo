package dsd

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/types"
)

const (
	MaxCronNum = 4
)

// CronFunction 定时功能
type CronFunction int

const (
	None   CronFunction = iota // None
	Preset                     // 预置点
	Cruise                     // 巡航
	Trace                      // 巡迹
	Line                       // 线性扫描
	Region                     // 区域扫描
)

func (f CronFunction) Validate() error {
	if f < None || f > Region {
		return errors.New("invalid cron function")
	}

	return nil
}

type CronID int

func (i CronID) Validate() error {
	if err := _validate.Var(i, "gte=1,lte=4"); err != nil {
		log.Error(err.Error())
		return errors.New("cron id is invalid")
	}

	return nil
}

type PtzAutoMovement struct {
	ID              CronID             `json:"ID" validate:"required,gte=1,lte=4"`
	Enable          bool               `json:"Enable" validate:"boolean"`
	Function        int                `json:"Function" validate:"gte=0,lte=5"`
	PresetID        int                `json:"PresetID" validate:"required,gte=1,lte=255"`
	TourID          int                `json:"TourID" validate:"required,gte=1,lte=8"`
	PatternID       int                `json:"PatternID" validate:"required,gte=1,lte=5"`
	LinearScanID    int                `json:"LinearScanID" validate:"required,gte=1,lte=5"`
	RegionScanID    int                `json:"RegionScanID" validate:"required,gte=1,lte=5"`
	AutoHoming      AutoHoming         `json:"AutoHoming" validate:"required"`
	RunningFunction int                `json:"RunningFunction" validate:"gte=0,lte=5"`
	Schedule        types.WeekSchedule `json:"Schedule" validate:"required"`
}

type AutoHoming struct {
	Enable bool `json:"Enable" validate:"boolean"`
	Time   int  `json:"Time" validate:"required,gte=3"`
}

func NewAutoMovement(id CronID) PtzAutoMovement {
	schedule := types.WeekSchedule{}
	types.InitWeekSchedule(&schedule)

	return PtzAutoMovement{
		ID:           id,
		Enable:       false,
		Function:     0,
		PresetID:     1,
		TourID:       1,
		PatternID:    1,
		LinearScanID: 1,
		RegionScanID: 1,
		AutoHoming: AutoHoming{
			Enable: true,
			Time:   3,
		},
		RunningFunction: 0,
		Schedule:        schedule,
	}
}

func (a *PtzAutoMovement) Validate() error {
	if err := _validate.Struct(a); err != nil {
		return err
	}
	return nil
}

func (a *PtzAutoMovement) GetFuncID() (int, error) {
	funcID := 0

	switch CronFunction(a.Function) {
	case Preset:
		funcID = a.PresetID
	case Cruise:
		funcID = a.TourID
	case Trace:
		funcID = a.PatternID
	case Line:
		funcID = a.LinearScanID
	case Region:
		funcID = a.RegionScanID
	default:
		return funcID, errors.New("invalid cron function")
	}

	return funcID, nil
}

type AutoMovementSlice []PtzAutoMovement

func NewAutoMovementSlice() AutoMovementSlice {
	movements := make([]PtzAutoMovement, MaxCronNum)

	for id := 1; id <= MaxCronNum; id++ {
		movement := NewAutoMovement(CronID(id))
		movements[id-1] = movement
	}
	return movements
}

func (s *AutoMovementSlice) ConfigKey() string {
	return "PtzAutoMovements"
}
