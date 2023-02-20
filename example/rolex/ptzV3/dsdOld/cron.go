package dsdOld

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/types"
)

const (
	MaxCronNum = 4
)

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
	Function        int                `json:"Function" validate:"required,gte=0,lte=5"`
	PresetID        int                `json:"PresetID" validate:"required,gte=1,lte=255"`
	TourID          int                `json:"TourID" validate:"required,gte=1,lte=8"`
	PatternID       int                `json:"PatternID" validate:"required,gte=1,lte=5"`
	LinearScanID    int                `json:"LinearScanID" validate:"required,gte=1,lte=5"`
	RegionScanID    int                `json:"RegionScanID" validate:"required,gte=1,lte=5"`
	AutoHoming      AutoHoming         `json:"AutoHoming" validate:"required"`
	RunningFunction int                `json:"RunningFunction" validate:"required,gte=0,lte=5"`
	Schedule        types.WeekSchedule `json:"Schedule" validate:"required"`
}

type AutoHoming struct {
	Enable bool `json:"Enable" validate:"boolean"`
	Time   int  `json:"Time" validate:"required,gte=3"`
}

func NewPtzAutoMovement(id CronID) (*PtzAutoMovement, error) {
	schedule := types.WeekSchedule{}
	types.InitWeekSchedule(&schedule)

	if err := id.Validate(); err != nil {
		return nil, err
	}

	return &PtzAutoMovement{
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
	}, nil
}

func (m *PtzAutoMovement) DefaultSlice() []PtzAutoMovement {
	movements := make([]PtzAutoMovement, MaxCronNum)

	for id := 1; id <= MaxCronNum; id++ {
		movement, _ := NewPtzAutoMovement(CronID(id))
		movements[id-1] = *movement
	}
	return movements
}

func (m *PtzAutoMovement) ConfigKey() string {
	return "PtzAutoMovements"
}
