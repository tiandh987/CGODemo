package dsd

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"time"
)

const (
	MaxTraceNum = 5
)

type TraceID int

func (i TraceID) Validate() error {
	if err := _validate.Var(i, "gte=1,lte=5"); err != nil {
		log.Error(err.Error())
		return errors.New("trace id is invalid")
	}

	return nil
}

type Pattern struct {
	ID            int       `json:"ID" validate:"required,gte=1,lte=5"` // ID
	Enable        bool      `json:"Enable" validate:"boolean"`          // 使能
	Check         int       `json:"Check"`                              // 校验 点击开始记录和停止记录都加2，结果为4才能开始
	Commands      []Command `json:"Commands"`                           // 用户操作记录
	StartPosition Position  `json:"StartPosition"`                      // 巡迹PTZ坐标起始位置
	EndPosition   Position  `json:"EndPosition"`                        // 巡迹PTZ坐标结束位置
	Running       bool      `json:"Runing"`                             // 运行状态
}

type Command struct {
	CommInfo  *CommInfo `json:"CommInfo"`
	StartTime int       `json:"StartTime"` // 执行时间戳:单位毫秒
}

type CommInfo struct {
	CommType int      `json:"CommType"`
	Data     []string `json:"Data"`
}

// =====

type Record struct {
	ID            TraceID
	Enable        bool
	Valid         bool
	Schedules     []Schedule
	StartPosition Position
}

type Schedule struct {
	FuncID    int
	Speed     int
	StartTime time.Time
	StopTime  time.Time
}

type RecordSlice []Record

func (s *RecordSlice) ConfigKey() string {
	return "Traces"
}

func NewRecordSlice() RecordSlice {
	s := make([]Record, MaxTraceNum)

	for id := 1; id <= MaxTraceNum; id++ {
		s[id-1] = Record{
			ID:            TraceID(id),
			Enable:        true,
			Valid:         false,
			Schedules:     []Schedule{},
			StartPosition: NewPosition(),
		}
	}

	return s
}
