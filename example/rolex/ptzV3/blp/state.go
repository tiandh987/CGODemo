package blp

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"time"
)

// Ability 云台支持能力
type Ability int

const (
	None       Ability = iota //
	Cruise                    // 巡航
	Trace                     // 巡迹
	LineScan                  // 线性扫描
	RegionScan                // 区域扫描
	PanMove                   // 水平旋转
	Preset                    // 预置点
	ManualFunc                // 手动转动云台
)

// Trigger 云台运动状态变更的触发者
type Trigger int

const (
	ManualTrigger  Trigger = iota // 手动
	AlarmTrigger                  // 报警联动
	CronTrigger                   // 定时任务
	PowerUpTrigger                // 开机动作
	IdleTrigger                   // 空闲动作
)

// Compare 触发者优先级对比，数字越小，优先级越高
// 返回 ture，表示 t 的优先级高于 trigger
// 返回 false，表示 t 的优先级低于 trigger
func (t Trigger) compare(trigger Trigger) bool {
	return t < trigger
}

const (
	stateNoneTopic = "stateNoneTopic" // 发布云台 None 状态
)

type state struct {
	trigger   Trigger   // 触发者
	function  Ability   // 当前运行功能
	funcID    int       // 当前运行功能ID
	startTime time.Time // 开始时间，单位：毫秒

	pubSub *gochannel.GoChannel
}

func newState() *state {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{
			Persistent:                     false,
			BlockPublishUntilSubscriberAck: false,
		},
		watermill.NewStdLogger(false, false),
	)

	return &state{
		trigger:   IdleTrigger,
		function:  None,
		funcID:    0,
		startTime: time.Now(),
		pubSub:    pubSub,
	}
}

func (s *state) string() string {
	str := fmt.Sprintf("trigger: %d function: %d funcID: %d startTime: %s",
		s.trigger, s.function, s.funcID, s.startTime.Format("2006-01-02 15:04:05"))

	return str
}

func (s *state) update(req *Request) {
	s.trigger = req.Trigger
	s.function = req.Ability
	s.funcID = req.ID
	s.startTime = time.Now()

	if req.Ability == None {
		msg := message.NewMessage(watermill.NewUUID(), message.Payload(s.startTime.Format("2006-01-02 15:04:05")))
		if err := s.pubSub.Publish(stateNoneTopic, msg); err != nil {
			log.Error(err.Error())
		}
	}
}

func (s *state) subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	subscribe, err := s.pubSub.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}

	return subscribe, nil
}

func (s *state) getExternal() *dsd.Status {
	status := dsd.NewStatus()
	status.StartTime = s.startTime.UnixMilli()
	status.Trigger = int(s.trigger)

	if s.function == None {
		return &status
	}

	status.Moving = true
	status.Function = int(s.function)
	status.FunctionID = s.funcID

	return &status
}

func (s *state) getInternal() state {
	st := state{
		trigger:   s.trigger,
		function:  s.function,
		funcID:    s.funcID,
		startTime: s.startTime,
	}

	return st
}
