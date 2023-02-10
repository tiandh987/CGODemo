package ptz

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"

	"time"
)

// Trigger 云台运动状态变更的触发者
type Trigger int

// Compare 触发者优先级对比，数字越小，优先级越高
// 返回 ture，表示 t 的优先级高于 trigger
// 返回 false，表示 t 的优先级低于 trigger
func (t Trigger) Compare(trigger Trigger) bool {
	return t <= trigger
}

const (
	Manual  Trigger = iota // 手动
	Alarm                  // 报警联动
	Cron                   // 定时任务
	PowerUp                // 开机动作
	Idle                   // 空闲动作
)

// Function 云台当前运行功能
type Function int

const (
	None       Function = iota // 未运行任务
	Cruise                     // 巡航
	Trace                      // 巡迹
	LineScan                   // 线性扫描
	RegionScan                 // 区域扫描
	PanMove                    // 水平旋转
	Preset                     // 预置点
	ManualFunc                 // 手动转动云台
)

type Limit struct {
	level         bool    // 水平限位使能
	leftBoundary  float64 // 左边界位置
	rightBoundary float64 // 右边界位置

	verticalEnable bool    // 垂直限位使能
	upBoundary     float64 // 上边界位置
	downBoundary   float64 // 下边界位置
}

func newLimit(l *dsd.Limit) *Limit {
	if l == nil {
		return &Limit{
			level:          false,
			leftBoundary:   0,
			rightBoundary:  0,
			verticalEnable: false,
			upBoundary:     0,
			downBoundary:   0,
		}
	}

	return &Limit{
		level:          l.LevelEnable,
		leftBoundary:   l.LeftBoundary,
		rightBoundary:  l.RightBoundary,
		verticalEnable: l.VerticalEnable,
		upBoundary:     l.UpBoundary,
		downBoundary:   l.DownBoundary,
	}
}

type State struct {
	trigger    Trigger  // 触发者
	function   Function // 当前运行功能
	functionID int      // 当前运行功能ID
	cronID     int      // 定时任务ID
	startTime  int64    // 开始时间，单位：毫秒
	limit      *Limit   // 云台限位
	version    string   // 云台版本
	model      string   // 云台型号
}

func NewState(limit *dsd.Limit) *State {
	return &State{
		trigger:    Idle,
		function:   None,
		functionID: 0,
		cronID:     0,
		startTime:  0,
		limit:      newLimit(limit),
		version:    "",
		model:      "",
	}
}

func (s *State) Validate(trigger Trigger, l *Limit) error {
	// 云台处于停止状态
	if s.function == None {
		return nil
	}

	// 云台处于运动状态，判断触发者优先级
	if !s.trigger.Compare(trigger) {
		return errors.New("trigger is low priority")
	}

	// TODO 限位校验

	//if l.upBoundary < s.limit.upBoundary || l.upBoundary > s.limit.downBoundary {
	//	return false
	//}
	//
	//if l.downBoundary < s.limit.upBoundary || l.downBoundary > s.limit.downBoundary {
	//	return false
	//}
	//
	//if l.leftBoundary < s.limit.leftBoundary || l.leftBoundary > s.limit.rightBoundary {
	//	return false
	//}
	//
	//if l.downBoundary < s.limit.upBoundary || l.downBoundary > s.limit.downBoundary {
	//	return false
	//}

	return nil
}

func (s *State) Change(trigger Trigger, function Function, funcID int, cronID int) error {
	log.Debugf("request param trigger: %d function: %d funcID: %d cronID: %d",
		trigger, function, funcID, cronID)

	s.trigger = trigger
	s.function = function
	s.functionID = funcID
	s.cronID = cronID
	s.startTime = time.Now().UnixMilli()

	// TODO pub-sub 模式，向其它模块报告云台当前状态，以便于测温、智能分析及时做出策略的变更

	return nil
}

func (s *State) Version() string {
	return s.version
}

func (s *State) Model() string {
	return s.model
}

func (s *State) Function() Function {
	return s.function
}
