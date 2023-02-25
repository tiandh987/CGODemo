package cruise

import (
	"context"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sync"
	"time"
)

type Cruise struct {
	mu      sync.RWMutex
	cruises dsd.CruiseSlice
	running dsd.CruiseID

	preset *preset.Preset
	basic  *basic.Basic

	indexMu  sync.Mutex
	index    int
	maxIndex int

	jumpCh chan struct{}
	stayCh chan struct{}
	skipCh chan struct{}
	timer  *time.Timer
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func New(basic *basic.Basic, preset *preset.Preset, cruises dsd.CruiseSlice) *Cruise {
	c := &Cruise{
		cruises:  cruises,
		basic:    basic,
		preset:   preset,
		running:  0,
		index:    0,
		maxIndex: 0,
		jumpCh:   make(chan struct{}, 1),
		stayCh:   make(chan struct{}, 1),
		skipCh:   make(chan struct{}, 1),
		timer:    time.NewTimer(time.Second),
	}

	return c
}

func (c *Cruise) Start(id dsd.CruiseID) error {
	if err := id.Validate(); err != nil {
		return err
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

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running > 0 {
		log.Warnf("cruise (%d) is running", c.running)
		return fmt.Errorf("cruise (%d) is running", c.running)
	}

	c.running = id
	c.cruises[id-1].Running = true
	c.initIndex(id)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.timer.Reset(time.Millisecond)
	<-c.timer.C

	c.stayCh = make(chan struct{}, 1)
	c.skipCh = make(chan struct{}, 1)
	c.jumpCh = make(chan struct{}, 1)

	startCh := make(chan struct{})

	c.wg.Add(1)
	go func(ctx context.Context, id dsd.CruiseID) {
		defer c.wg.Done()
		defer log.Info("dfasdfasdfasdf")

		log.Infof("start ctx: %p", ctx)

		log.Infof("start cruise id: %d index: %d maxIndex: %d", id, c.index, c.maxIndex)
		for {
			select {
			case <-ctx.Done():
				log.Infof("start down ctx: %p", ctx)

				log.Warnf("end cruise (%d) %s", id, ctx.Err().Error())
				return
			case <-startCh:
				log.Infof("startCh receive")
				c.jumpCh <- struct{}{}
				log.Infof("startCh send")

			case <-c.timer.C:
				log.Infof("timer receive")
				c.jumpCh <- struct{}{}
				log.Infof("timer send")

			case <-c.jumpCh:
				log.Infof("jumpCh receive")
				c.jumping(ctx, id)
				log.Infof("jumpCh send")

			case <-c.stayCh:
				log.Infof("stayCh receive")
				c.residence(id)
				log.Infof("stayCh send")

			case <-c.skipCh:
				log.Infof("cruise (%d - %d) skip", id, c.index)
				c.updateIndex()
				c.jumpCh <- struct{}{}
				log.Infof("skipCh send")

			}
		}
	}(c.ctx, id)

	log.Infof("start 1")
	startCh <- struct{}{}
	log.Infof("start 2")

	return nil
}

func (c *Cruise) Stop(id dsd.CruiseID) error {
	if err := id.Validate(); err != nil {
		return err
	}

	log.Infof("id: %d running: %d cfg: %t", id, c.running, c.cruises[id-1].Running)

	c.mu.Lock()
	log.Infof("0")

	defer c.mu.Unlock()
	if (c.running > 0 && c.cruises[id-1].Running) || (c.running <= 0 && c.cruises[id-1].Running) {
		log.Infof("1")
		c.cancel()
		log.Infof("2")

		c.wg.Wait()
		log.Infof("3")

		c.timer.Stop()

		c.running = 0
		log.Infof("4")

		c.cruises[id-1].Running = false
		log.Infof("5")

		c.resetIndex()
		log.Infof("6")

		close(c.stayCh)
		close(c.skipCh)
		close(c.jumpCh)

		log.Infof("7")
		return nil
	}

	return fmt.Errorf("cruise (%d) is not running", id)
}

func (c *Cruise) initIndex(id dsd.CruiseID) {
	c.indexMu.Lock()
	defer c.indexMu.Unlock()

	c.index = 0
	c.maxIndex = len(c.cruises[id-1].Preset)
}

func (c *Cruise) updateIndex() {
	c.indexMu.Lock()
	defer c.indexMu.Unlock()

	c.index++

	if c.index >= c.maxIndex {
		c.index = 0
	}
}

func (c *Cruise) resetIndex() {
	c.indexMu.Lock()
	defer c.indexMu.Unlock()

	c.index = 0
	c.maxIndex = 0
}

func (c *Cruise) jumping(ctx context.Context, id dsd.CruiseID) {
	log.Infof("cruise jumping (%d - %d)", id, c.index)

	p := c.cruises[id-1].Preset[c.index]
	if err := c.preset.Start(c.ctx, p.ID); err != nil {
		log.Error(err.Error())
		c.skipCh <- struct{}{}
		return
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer log.Infof("sssssssssssssssssss")

		timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()

		ticker := time.NewTicker(time.Millisecond * 10)
		dst := c.preset.List()[p.ID-1].Position

		log.Infof("position: %+v", dst)

		for {
			select {
			case <-timeoutCtx.Done():

				log.Infof("%p %p", ctx, timeoutCtx)

				log.Warn(timeoutCtx.Err().Error())
				c.skipCh <- struct{}{}
				return
			case <-ticker.C:
				pos, err := c.basic.Position()
				if err != nil {
					log.Warn(ctx.Err().Error())
					c.skipCh <- struct{}{}
					return
				}

				if pos.Pan >= dst.Pan-2 && pos.Pan <= dst.Pan+2 &&
					pos.Tilt >= dst.Tilt-2 && pos.Tilt <= dst.Tilt+2 &&
					pos.Zoom >= dst.Zoom-2 && pos.Zoom <= dst.Zoom+2 {
					c.basic.Stop()
					c.stayCh <- struct{}{}
					return
				}
			}
		}
	}()

	return
}

func (c *Cruise) residence(id dsd.CruiseID) {
	p := c.cruises[id-1].Preset[c.index]
	log.Infof("cruise residence (%d - %d - %ds)", id, c.index, p.ResidenceTime)

	c.timer.Reset(time.Second * time.Duration(p.ResidenceTime))
	c.updateIndex()

	return
}
