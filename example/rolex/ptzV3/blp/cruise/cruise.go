package cruise

import (
	"context"
	"errors"
	"fmt"
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
	skip      = "skip"      // 跳过预置点
)

type Cruise struct {
	mu      sync.RWMutex
	cruises dsd.CruiseSlice

	fsm      *fsm.FSM
	fsmCh    chan string
	index    int
	maxIndex int
	timer    *time.Timer
	preset   *preset.Preset
	stopCh   chan dsd.CruiseID
}

func New(preset *preset.Preset, cruises dsd.CruiseSlice) *Cruise {
	c := &Cruise{
		cruises: cruises,
		preset:  preset,
		timer:   time.NewTimer(time.Second),
		fsmCh:   make(chan string, 1),
		stopCh:  make(chan dsd.CruiseID, 1),
	}

	c.fsm = fsm.NewFSM(
		none,
		fsm.Events{
			{Name: none, Src: []string{jumping, skip, residence}, Dst: none},
			{Name: jumping, Src: []string{none, skip, residence}, Dst: jumping},
			{Name: residence, Src: []string{jumping}, Dst: residence},
			{Name: skip, Src: []string{jumping}, Dst: skip},
		},
		fsm.Callbacks{
			"enter_none": func(ctx context.Context, event *fsm.Event) {
				return
			},
			"enter_jumping": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.CruiseID)
				c.jumping(ctx, id)
				return
			},
			"enter_skip": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.CruiseID)
				log.Infof("cruise (%d - %d) enter_skip", id, c.index)
				c.updateIndex()
				c.fsmCh <- jumping
				return
			},
			"enter_residence": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.CruiseID)
				c.residence(ctx, id)
				return
			},
		})

	return c
}

func (c *Cruise) Start(ctx context.Context, id dsd.CruiseID) error {
	if err := id.Validate(); err != nil {
		return err
	}

	if c.fsm.Current() != none {
		log.Warnf("cruise is running")
		return errors.New("cruise is running")
	}

	cruise := c.cruises[id-1]
	if !cruise.Enable {
		log.Warnf("cruise (%d - %s) is disable", id, cruise.Name)
		return fmt.Errorf("cruise (%d - %s) is disable", id, cruise.Name)
	}

	if len(cruise.Preset) == 0 {
		log.Warnf("cruise (%d - %s) is empty", id, cruise.Name)
		return fmt.Errorf("cruise (%d - %s) is empty", id, cruise.Name)
	}

	go func(id dsd.CruiseID) {
		c.initIndex(id)
		c.cruises[id-1].Running = true
		c.fsmCh = make(chan string, 1)
		c.stopCh = make(chan dsd.CruiseID, 1)
		c.timer.Reset(time.Second * 3)

		log.Infof("start cruise id: %d index: %d maxIndex: %d", id, c.index, c.maxIndex)

		for {
			select {
			case <-ctx.Done():
				log.Warn(ctx.Err().Error())
				goto EndCruise
			case stopId := <-c.stopCh:
				if stopId != id {
					log.Warnf("current id (%d), request id (%d)", id, stopId)
					continue
				}
				goto EndCruise
			case <-c.timer.C:
				c.fsmCh <- jumping
			case state := <-c.fsmCh:
				switch state {
				case jumping:
					if err := c.fsm.Event(ctx, jumping, id); err != nil {
						log.Error(err.Error())
						goto EndCruise
					}
				case residence:
					if err := c.fsm.Event(ctx, residence, id); err != nil {
						log.Error(err.Error())
						goto EndCruise
					}
				case skip:
					if err := c.fsm.Event(ctx, skip, id); err != nil {
						log.Error(err.Error())
						goto EndCruise
					}
				}
			}
		}
	EndCruise:
		log.Infof("end cruise (%d)", id)
		c.cruises[id-1].Running = false
		c.fsm.Event(ctx, none)
		close(c.fsmCh)
		close(c.stopCh)
	}(id)

	return nil
}

func (c *Cruise) Stop(ctx context.Context, id dsd.CruiseID) error {
	if err := id.Validate(); err != nil {
		return err
	}

	log.Infof("current: %s, id: %d, running:%t", c.fsm.Current(), id, c.cruises[id-1].Running)

	if c.fsm.Current() != none && c.cruises[id-1].Running {
		c.stopCh <- id
		return nil
	}

	if c.fsm.Current() == none && c.cruises[id-1].Running {
		c.cruises[id-1].Running = false
	}

	return fmt.Errorf("cruise (%d) is not running", id)
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

func (c *Cruise) jumping(ctx context.Context, id dsd.CruiseID) {
	log.Infof("cruise enter_jumping (%d - %d)", id, c.index)

	preset := c.cruises[id-1].Preset[c.index]
	if err := c.preset.Start(ctx, preset.ID); err != nil {
		log.Error(err.Error())
		c.fsmCh <- skip
		return
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc()
	if err := c.preset.ReachPreset(timeoutCtx, preset.ID); err != nil {
		log.Error(err.Error())
		c.fsmCh <- skip
		return
	}

	c.fsmCh <- residence

	return
}

func (c *Cruise) residence(ctx context.Context, id dsd.CruiseID) {
	log.Infof("cruise enter_residence (%d - %d)", id, c.index)

	preset := c.cruises[id-1].Preset[c.index]
	c.updateIndex()
	c.timer.Reset(time.Second * time.Duration(preset.ResidenceTime))

	return
}
