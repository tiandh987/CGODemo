package cruise

import (
	"context"
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"sync"
	"time"
)

type state struct {
	cst      cruiseState
	index    int
	maxIndex int
	ch       chan struct{}
}

func newState() *state {
	return &state{
		cst:      none,
		index:    0,
		maxIndex: 0,
		ch:       make(chan struct{}, 1),
	}
}

func (s *state) reset(max int) {
	s.cst = none
	s.index = 0
	s.maxIndex = max
	s.ch = make(chan struct{}, 1)
}

func (s *state) updateCst(cst cruiseState) {
	s.cst = cst
	if cst != residence {
		s.ch <- struct{}{}
	}
}

func (s *state) updateIndex() {
	if s.index >= s.maxIndex {
		s.index = 0
	} else {
		s.index++
	}
}

type Cruise struct {
	mu      sync.RWMutex
	cruises []dsd.TourPreset

	timer *time.Timer
	state *state
	errCh chan error
	quit  chan struct{}
}

func New(cruises []dsd.TourPreset) *Cruise {
	return &Cruise{
		cruises: cruises,
		state:   newState(),
		timer:   time.NewTimer(time.Hour),
		errCh:   make(chan error, 1),
		quit:    make(chan struct{}, 1),
	}
}

type cruiseState int

const (
	none cruiseState = iota // None
	gotoPreset
	atPreset
	residence
)

func (c *Cruise) Start(ctl control.ControlRepo, preset *preset.Preset, id dsd.CruiseID) error {
	if err := id.Validate(); err != nil {
		return err
	}

	if c.state.cst != none {
		log.Warn("cruise is running")
		return errors.New("cruise is running")
	}

	cruise := c.getCruise(id)
	if !cruise.Enable {
		log.Warnf("cruise %d is disable", id)
		return errors.New(fmt.Sprintf("cruise %d is disable", id))
	}

	if len(cruise.Preset) == 0 {
		log.Warnf("cruise %d preset is empty", id)
		return errors.New(fmt.Sprintf("cruise %d preset is empty", id))
	}

	go func() {
		log.Infof("start of cruise...\nconfig: %+v", cruise)

		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()

		// 初始化状态
		c.setCruiseRunning(id, true)
		if err := c.saveConfig(); err != nil {
			log.Error(err.Error())
			goto EndCruise
		}

		c.state.reset(len(cruise.Preset) - 1)
		c.state.updateCst(none)
		log.Debugf("cruise init state: %+v", cruise)

		for {
			select {
			case <-c.quit:
				goto EndCruise
			case err := <-c.errCh:
				log.Error(err.Error())
				goto EndCruise
			case <-c.timer.C:
				switch c.state.cst {
				case residence:
					c.state.updateIndex()
					c.state.updateCst(gotoPreset)
				}
			case <-c.state.ch:
				switch c.state.cst {
				case none:
					c.state.updateCst(gotoPreset)
				case gotoPreset:
					c.gotoPreset(ctx, ctl, preset, &cruise)
				case atPreset:
					c.timer.Reset(time.Second * time.Duration(cruise.Preset[c.state.index].ResidenceTime))
					c.state.updateCst(residence)
				}
			}
		}
	EndCruise:
		log.Infof("end cruise (%d)", id)

		ctl.Stop()
		c.state.reset(0)
		c.setCruiseRunning(id, false)
		c.saveConfig()
	}()

	return nil
}

func (c *Cruise) Stop() {
	if c.state.cst != none {
		c.quit <- struct{}{}
	}
}

func (c *Cruise) gotoPreset(ctx context.Context, ctl control.ControlRepo, preset *preset.Preset, cruise *dsd.TourPreset) {
	if err := preset.Start(ctl, cruise.Preset[c.state.index].ID); err != nil {
		c.errCh <- err
	}

	presetPosition, err := preset.GetPosition(cruise.Preset[c.state.index].ID)
	if err != nil {
		c.errCh <- err
	}
	log.Debugf("expect preset id: %d position: %+v", cruise.Preset[c.state.index].ID, presetPosition)

	go func() {
		//timer := time.NewTimer(time.Second * 30)
		timer := time.NewTimer(time.Second * 10)

		for {
			select {
			case <-ctx.Done():
				log.Info(ctx.Err().Error())
				return
			case <-timer.C:
				c.errCh <- errors.New("timeout waiting for jump to the expect preset")
				return
			default:
				pos, err := ctl.Position()
				if err != nil {
					c.errCh <- err
					return
				}

				if pos.Pan >= presetPosition.Pan-2 && pos.Pan <= presetPosition.Pan+2 &&
					pos.Tilt >= presetPosition.Tilt-2 && pos.Tilt <= presetPosition.Tilt+2 &&
					pos.Zoom >= presetPosition.Zoom-2 && pos.Zoom <= presetPosition.Zoom+2 {
					c.state.updateCst(atPreset)
					return
				}
				//time.Sleep(time.Millisecond * 10)
				time.Sleep(time.Second)
			}
		}
	}()
}

func (c *Cruise) getCruise(id dsd.CruiseID) dsd.TourPreset {
	c.mu.RLock()
	cruise := c.cruises[id-1]
	c.mu.RUnlock()

	return cruise
}

func (c *Cruise) setCruiseRunning(id dsd.CruiseID, running bool) {
	c.mu.Lock()
	c.cruises[id-1].Running = running
	c.mu.Unlock()
}

func (c *Cruise) saveConfig() error {
	if err := config.SetConfig(c.cruises[0].ConfigKey(), c.cruises); err != nil {
		return err
	}

	return nil
}
