package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/arch/serial"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/cruise"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/idle"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/lineScan"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/powerUp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"sync"
	"time"
)

var (
	mu           sync.RWMutex
	_blpInstance *Blp
)

type Blp struct {
	mu        sync.RWMutex
	state     *ptz.State
	serialCtl control.ControlRepo
	mediaCtl  control.ControlRepo
	basic     *basic.Basic
	preset    *preset.Preset
	line      *lineScan.LineScan
	cruise    *cruise.Cruise
	power     *powerUp.PowerUp
	idle      *idle.Idle
}

//func NewBlp(st *ptz.State, sCtl *serial.Serial, mCtl control.ControlRepo, basic *ptz.Basic, preset *preset.Preset,
//	line *lineScan.LineScan) *Blp {

func NewBlp(st *ptz.State, sCtl *serial.Serial, basic *basic.Basic, preset *preset.Preset,
	line *lineScan.LineScan, cruise *cruise.Cruise, power *powerUp.PowerUp, idle *idle.Idle) *Blp {

	return &Blp{
		state:     st,
		serialCtl: sCtl,
		//mediaCtl:  mCtl,
		basic:  basic,
		preset: preset,
		line:   line,
		cruise: cruise,
		power:  power,
		idle:   idle,
	}
}

func Instance() *Blp {
	mu.Lock()
	ins := _blpInstance
	mu.Unlock()

	return ins
}

func Replace(ins *Blp) {
	mu.Lock()
	_blpInstance = ins
	mu.Unlock()
}

func (b *Blp) Control(trigger ptz.Trigger, function ptz.Function, funcID, cronID int, speed ptz.Speed) error {
	log.Debugf("request param trigger: %d function: %d funcID: %d cronID: %d speed: %d",
		trigger, function, funcID, cronID, speed)

	log.Debugf("current state: %+v\n\n", b.state)

	if err := b.validate(trigger, function, funcID, cronID); err != nil {
		return err
	}

	// 触发者为定时任务时，更新 function、funcID
	if trigger == ptz.Cron {
		//	TODO 更新 function、funcID
	}

	ctl := b.getControl()

	b.mu.Lock()
	defer b.mu.Unlock()

	// 停止当前云台正在运行的动作
	switch b.state.Function() {
	case ptz.Cruise:
		b.cruise.Stop()
	case ptz.Trace:
	case ptz.LineScan:
		b.line.Stop()
	case ptz.RegionScan:
	case ptz.PanMove:
	case ptz.Preset:
		if err := ctl.Stop(); err != nil {
			return err
		}
	case ptz.ManualFunc:
		if err := ctl.Stop(); err != nil {
			return err
		}
	}

	// 等待之前云台动作停止
	time.Sleep(time.Millisecond * 50)

	// 云台空闲动作检测协程
	if trigger != ptz.Idle && function != ptz.None {
		log.Info("pause idle")

		b.idle.Pause()
	}

	// 转动云台
	switch function {
	case ptz.None:
		log.Info("reset idle")

		b.idle.Reset()
	case ptz.Cruise:
		if err := b.cruise.Start(ctl, b.preset, dsd.CruiseID(funcID)); err != nil {
			return err
		}
	case ptz.Trace:
	case ptz.LineScan:
		if err := b.line.Start(ctl, dsd.LineScanID(funcID)); err != nil {
			return err
		}
	case ptz.RegionScan:
	case ptz.PanMove:
	case ptz.Preset:
		if err := b.preset.Start(ctl, dsd.PresetID(funcID)); err != nil {
			return err
		}
	case ptz.ManualFunc:
		if err := b.basic.Start(ctl, funcID, speed); err != nil {
			return err
		}
	}

	if err := b.state.Change(trigger, function, funcID, cronID); err != nil {
		return err
	}

	return nil
}

func (b *Blp) validate(trigger ptz.Trigger, function ptz.Function, funcID, cronID int) error {
	var l ptz.Limit

	// TODO 根据 function、funcID 获取 limit，

	switch function {

	}

	return b.state.Validate(trigger, &l)
}

func (b *Blp) Version() string {
	version, err := b.getControl().Version()
	if err != nil {
		log.Error(err.Error())
		version = ""
	}
	return version
}

func (b *Blp) Model() string {
	model, err := b.getControl().Model()
	if err != nil {
		log.Error(err.Error())
		model = ""
	}
	return model
}

func (b *Blp) Restart() error {
	ctl := b.serialCtl

	// TODO ctl 选择

	return ctl.Restart()
}

func (b *Blp) getControl() control.ControlRepo {
	// TODO 判断使用 serial 还是 media 进行通信

	ctl := b.serialCtl

	return ctl
}

func (b *Blp) State() *dsd.Status {
	return b.state.Convert()
}

// TODO delete
func (b *Blp) Position(pos *dsd.Position) error {
	// TODO 判断使用 serial 还是 media 进行通信

	ctl := b.serialCtl

	return ctl.Goto(pos)
}
