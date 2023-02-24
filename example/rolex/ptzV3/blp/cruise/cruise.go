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
	ctx      context.Context
	cancel   context.CancelFunc
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
				go c.jumping(id)
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
				c.residence(id)
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

	go func() {
		log.Infof("1 cruise stop channel %+v", c.stopCh)

		c.initIndex(id)
		c.ctx, c.cancel = context.WithCancel(ctx)
		c.fsmCh = make(chan string, 1)
		c.stopCh = make(chan dsd.CruiseID, 1)
		c.fsmCh <- jumping

		log.Infof("start cruise id: %d index: %d maxIndex: %d", id, c.index, c.maxIndex)

		log.Infof("2 cruise stop channel %+v", c.stopCh)

		for {
			c.cruises[id-1].Running = true

			select {
			case <-ctx.Done():
				log.Warn(ctx.Err().Error())
				goto EndCruise
			case stopId := <-c.stopCh:
				log.Infof("cruise receive stop id: %d", stopId)
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
					if err := c.fsm.Event(c.ctx, jumping, id); err != nil {
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
			default:
				log.Infof("3 cruise stop channel %+v fsm.Current: %s", c.stopCh, c.fsm.Current())
				time.Sleep(time.Second * 5)
			}
		}
	EndCruise:
		log.Infof("end cruise (%d)", id)
		c.cruises[id-1].Running = false
		c.fsm.Event(ctx, none)
		c.cancel()
		//close(c.fsmCh)
	}()

	return nil
}

func (c *Cruise) Stop(ctx context.Context, id dsd.CruiseID) error {
	if err := id.Validate(); err != nil {
		return err
	}

	log.Infof("current: %s, id: %d, running:%t", c.fsm.Current(), id, c.cruises[id-1].Running)

	if (c.fsm.Current() != none && c.cruises[id-1].Running) ||
		(c.fsm.Current() == none && c.cruises[id-1].Running) {
		log.Infof("send %d to cruise stop channel %+v", id, c.stopCh)
		c.stopCh <- id
		close(c.stopCh)

		log.Infof("send %d to cruise stop channel ok", id)

		return nil
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

func (c *Cruise) jumping(id dsd.CruiseID) {
	log.Infof("cruise enter_jumping (%d - %d)", id, c.index)

	p := c.cruises[id-1].Preset[c.index]
	if err := c.preset.Start(c.ctx, p.ID); err != nil {
		log.Error(err.Error())
		c.fsmCh <- skip
		return
	}

	timeoutCtx, cancelFunc := context.WithTimeout(c.ctx, time.Second*60)
	defer cancelFunc()
	if err := c.preset.ReachPreset(timeoutCtx, p.ID); err != nil {
		log.Error(err.Error())
		c.fsmCh <- skip
		return
	}

	c.fsmCh <- residence

	return
}

func (c *Cruise) residence(id dsd.CruiseID) {
	log.Infof("cruise enter_residence (%d - %d)", id, c.index)

	p := c.cruises[id-1].Preset[c.index]
	c.updateIndex()
	c.timer.Reset(time.Second * time.Duration(p.ResidenceTime))

	return
}
