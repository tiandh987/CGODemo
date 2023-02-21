package cruise

import (
	"context"
	"github.com/looplab/fsm"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sync"
	"time"
)

// 巡航状态
const (
	none      = "none"
	jumping   = "jumping"   // 跳转中
	residence = "residence" // 停留
)

type Cruise struct {
	mu      sync.RWMutex
	cruises dsd.CruiseSlice

	fsm      *fsm.FSM
	preset   *preset.Preset
	index    int
	maxIndex int
}

func New(preset *preset.Preset, cruises dsd.CruiseSlice) *Cruise {
	c := &Cruise{
		cruises: cruises,
		preset:  preset,
	}

	c.fsm = fsm.NewFSM(
		none,
		fsm.Events{
			{Name: none, Src: []string{jumping, residence}, Dst: none},
			{Name: jumping, Src: []string{none, residence}, Dst: jumping},
			{Name: residence, Src: []string{jumping}, Dst: residence},
		},
		fsm.Callbacks{
			"enter_none": func(ctx context.Context, event *fsm.Event) {
				log.Info("cruise enter_none")
				return
			},
			"enter_jumping": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.CruiseID)
				if err := c.jumping(ctx, id); err != nil {
					log.Error(err.Error())
					c.fsm.Event(ctx, none)
					return
				}
			},

			"enter_residence": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.CruiseID)
				if err := c.residence(ctx, id); err != nil {
					log.Error(err.Error())
					c.fsm.Event(ctx, none)
					return
				}
			},
		})

	return c
}

func (c *Cruise) Start(ctx context.Context, id dsd.CruiseID) {
	if c.fsm.Current() != none {
		log.Warnf("cruise is running")
		return
	}

	cruise := c.cruises[id-1]
	if !cruise.Enable {
		log.Warnf("cruise (%d - %s) is disable", id, cruise.Name)
		return
	}

	c.initIndex(id)

	event := jumping
	if !c.fsm.Can(event) {
		log.Warnf("line scan can not convert to %s", event)
		return
	}

	log.Infof("event: %s id: %d index : %d maxIndex: %d", event, id, c.index, c.maxIndex)

	go c.fsm.Event(ctx, event, id)

	c.cruises[id-1].Running = true
}

func (c *Cruise) Stop(ctx context.Context, id dsd.CruiseID) {
	if c.fsm.Current() == none || !c.cruises[id-1].Running {
		return
	}

	c.fsm.Event(ctx, none)
	c.cruises[id-1].Running = false
}

func (c *Cruise) initIndex(id dsd.CruiseID) {
	c.index = 0
	c.maxIndex = len(c.cruises[id-1].Preset)
}

func (c *Cruise) updateIndex() {
	c.index++

	if c.index >= c.maxIndex {
		c.index = 0
	}
}

func (c *Cruise) resetIndex() {
	c.index = 0
	c.maxIndex = 0
}

func (c *Cruise) jumping(ctx context.Context, id dsd.CruiseID) error {
	log.Infof("cruise enter_jumping (%d - %d)", id, c.index)

	preset := c.cruises[id-1].Preset[c.index]
	if err := c.preset.Start(ctx, preset.ID); err != nil {
		return err
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc()

	if err := c.preset.ReachPreset(timeoutCtx, preset.ID); err != nil {
		return err
	}

	if err := c.fsm.Event(ctx, residence, id); err != nil {
		return err
	}

	return nil
}

func (c *Cruise) residence(ctx context.Context, id dsd.CruiseID) error {
	log.Infof("cruise enter_residence (%d - %d)", id, c.index)

	preset := c.cruises[id-1].Preset[c.index]

	time.Sleep(time.Second * time.Duration(preset.ResidenceTime))
	c.updateIndex()

	if err := c.fsm.Event(ctx, jumping, id); err != nil {
		return err
	}

	return nil
}
