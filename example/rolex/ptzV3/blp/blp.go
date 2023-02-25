package blp

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/cron"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/cruise"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/idle"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/line"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/powerUp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sync"
	"time"
)

var _validate = validator.New()

type Blp struct {
	mu    sync.Mutex
	state *state

	ctx        context.Context
	cancelFunc context.CancelFunc
	basic      *basic.Basic
	preset     *preset.Preset
	line       *line.Line
	cruise     *cruise.Cruise
	power      *powerUp.PowerUp
	idle       *idle.Idle
	cron       *cron.Cron
}

var _ PTZRepo = (*Blp)(nil)

type PTZRepo interface {
	Basic() BasicRepo
	Preset() PresetRepo
	Line() LineRepo
	Cruise() CruiseRepo
	Power() PowerRepo
	Idle() IdleRepo
	Cron() CronRepo
	Manager() ManagerRepo
}

type BasicRepo interface {
	Version() string
	Model() string
	Restart()
	Goto(pos *dsd.Position) error
}

type PresetRepo interface {
	List() dsd.PresetSlice
	Update(id dsd.PresetID, name string) error
	Delete(id dsd.PresetID) error
	DeleteAll() error
	Set(id dsd.PresetID, name string) error
}

type LineRepo interface {
	List() dsd.LineSlice
	Default() error
	Set(scan *dsd.LineScan) error
	SetMargin(id dsd.LineScanID, op dsd.LineMarginOp) error
}

type CruiseRepo interface {
	List() dsd.CruiseSlice
	Default() error
	Update(id dsd.CruiseID, name string) error
	Set(cr *dsd.TourPreset) error
	Delete(id dsd.CruiseID) error
}

func (b *Blp) Basic() BasicRepo {
	return b.basic
}

type PowerRepo interface {
	Get() *dsd.PowerUps
	Set(ups *dsd.PowerUps) error
	Default() error
}

func (b *Blp) Preset() PresetRepo {
	return b.preset
}

type IdleRepo interface {
	Get() *dsd.IdleMotion
	Set(motion *dsd.IdleMotion) error
	Default() error
}

type CronRepo interface {
	List() dsd.AutoMovementSlice
	Set(movement *dsd.PtzAutoMovement) error
	Default() error
}

type ManagerRepo interface {
	Run()
	Start(req *Request) error
	Stop(req *Request) error
	Quit()
	State() *dsd.Status
}

func (b *Blp) Line() LineRepo {
	return b.line
}

func (b *Blp) Cruise() CruiseRepo {
	return b.cruise
}

func (b *Blp) Power() PowerRepo {
	return b.power
}

func (b *Blp) Idle() IdleRepo {
	return b.idle
}

func (b *Blp) Cron() CronRepo {
	return b.cron
}

func (b *Blp) Manager() ManagerRepo {
	return b
}

func New(basic *basic.Basic, preset *preset.Preset, line *line.Line, cruise *cruise.Cruise, up *powerUp.PowerUp,
	i *idle.Idle, c *cron.Cron) PTZRepo {

	ctx, cancelFunc := context.WithCancel(context.Background())

	return &Blp{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		state:      newState(),
		basic:      basic,
		preset:     preset,
		line:       line,
		cruise:     cruise,
		power:      up,
		idle:       i,
		cron:       c,
	}
}

func (b *Blp) Run() {
	b.startPowerUp()
	b.startIdle()
	b.startCron()
}

func (b *Blp) Quit() {
	b.cancelFunc()
}

type Request struct {
	Trigger Trigger `json:"Trigger" validate:"gte=0,lte=4"`
	Ability Ability `json:"Ability" validate:"gte=0,lte=7"`
	ID      int     `json:"ID"`
	Speed   int     `json:"Speed" validate:"gte=1,lte=8"`
}

func (r *Request) Validate() error {
	if err := _validate.Struct(r); err != nil {
		return err
	}

	return nil
}

func (b *Blp) Start(req *Request) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err := req.Validate(); err != nil {
		return err
	}

	st := b.state.getInternal()
	// 校验触发者
	if st.trigger.compare(req.Trigger) && st.function != None {
		log.Warnf("trigger current: %d request: %d", st.trigger, req.Trigger)
		return errors.New("trigger priority is low")
	}

	// TODO 限位校验

	reqStop := &Request{
		Trigger: req.Trigger,
		Ability: st.function,
		ID:      st.funcID,
		Speed:   1,
	}

	// 手动触发的云台动作，必须进行手动停止
	if st.trigger == ManualTrigger && req.Trigger == ManualTrigger && st.function != None && st.function != Preset {
		return fmt.Errorf("ptz is running, trigger: manual function: %d, funcID: %d", st.function, st.funcID)
	}

	// 终止当前动作(当触发者为报警联动、定时任务、空闲动作时)
	if st.function != None {
		if err := b.stop(b.ctx, reqStop); err != nil {
			return err
		}
		//time.Sleep(time.Second * 10)
	}

	switch req.Ability {
	case Cruise:
		if err := b.cruise.Start(dsd.CruiseID(req.ID)); err != nil {
			return err
		}
	case Trace:

	case LineScan:
		if err := b.line.Start(b.ctx, dsd.LineScanID(req.ID)); err != nil {
			return err
		}
	case RegionScan:

	case PanMove:

	case Preset:
		if err := b.preset.Start(b.ctx, dsd.PresetID(req.ID)); err != nil {
			return err
		}
	case ManualFunc:
		if err := b.basic.Operation(basic.Operation(req.ID), ptz.Speed(req.Speed)); err != nil {
			return err
		}
	default:
		log.Warnf("invalid ptz ability (%d)", req.Ability)
		return errors.New("ptz ability is invalid")
	}

	// 等待云台启动完成
	time.Sleep(time.Millisecond * 100)

	b.state.update(req)

	return nil
}

func (b *Blp) Stop(req *Request) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	st := b.state.getInternal()
	// 校验触发者
	if st.trigger.compare(req.Trigger) && st.function != None {
		log.Warnf("trigger current: %d request: %d", st.trigger, req.Trigger)
		return errors.New("trigger priority is low")
	}

	return b.stop(b.ctx, req)
}

func (b *Blp) stop(ctx context.Context, req *Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	st := b.state.getInternal()
	if req.Ability != st.function || req.ID != st.funcID {
		return fmt.Errorf("current trigger: %d function: %d funcID: %d, request trigger: %d function: %d funcID: %d",
			st.trigger, st.function, st.funcID, req.Trigger, req.Ability, req.ID)
	}

	switch req.Ability {
	case Cruise:
		if err := b.cruise.Stop(dsd.CruiseID(req.ID)); err != nil {
			return err
		}
	case Trace:

	case LineScan:
		if err := b.line.Stop(dsd.LineScanID(req.ID)); err != nil {
			return err
		}
	case RegionScan:

	case PanMove:

	case Preset:
		if err := b.basic.Stop(); err != nil {
			return err
		}
	case ManualFunc:
		if err := b.basic.Stop(); err != nil {
			return err
		}
	default:
		log.Warnf("invalid ability (%d)", req.Ability)
		return errors.New("invalid ability")
	}

	req.Ability = None
	req.ID = 1
	req.Speed = 1
	b.state.update(req)

	return nil
}

func (b *Blp) State() *dsd.Status {
	return b.state.getExternal()
}
