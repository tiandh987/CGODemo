package cruise

import (
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

func (c *Cruise) List() dsd.CruiseSlice {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cruises
}

func (c *Cruise) Default() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running > 0 {
		log.Warnf("cruise(%d) is running", c.running)
		return fmt.Errorf("cruise(%d) is running", c.running)
	}

	before := c.cruises

	c.cruises = dsd.NewCruiseSlice()
	if err := config.SetConfig(c.cruises.ConfigKey(), c.cruises); err != nil {
		c.cruises = before
		return err
	}

	return nil
}

func (c *Cruise) Update(id dsd.CruiseID, name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cruise := dsd.NewCruise(id, name)
	if err := cruise.Validate(); err != nil {
		return err
	}

	if c.cruises[id-1].Running {
		log.Warnf("cruise (%d) is running", id)
		return errors.New("cruise is running")
	}

	before := c.cruises[id-1]

	cruise.Enable = before.Enable
	cruise.Preset = before.Preset
	cruise.Running = before.Running

	c.cruises[id-1] = cruise
	if err := config.SetConfig(c.cruises.ConfigKey(), c.cruises); err != nil {
		c.cruises[id-1] = before
		return err
	}

	return nil
}

func (c *Cruise) Set(cr *dsd.TourPreset) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cruises[cr.ID-1].Running {
		log.Warnf("cruise (%d) is running", cr.ID)
		return errors.New("cruise is running")
	}

	before := c.cruises[cr.ID-1]

	cr.Running = false
	if cr.Preset == nil {
		cr.Preset = []dsd.TourPresetPoint{}
	}
	c.cruises[cr.ID-1] = *cr

	if err := config.SetConfig(c.cruises.ConfigKey(), c.cruises); err != nil {
		c.cruises[cr.ID-1] = before
		return err
	}

	return nil
}

func (c *Cruise) Delete(id dsd.CruiseID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cruises[id-1].Running {
		log.Warnf("cruise (%d) is running", id)
		return errors.New("cruise is running")
	}

	before := c.cruises[id-1]

	c.cruises[id-1].Preset = []dsd.TourPresetPoint{}
	c.cruises[id-1].Running = false
	c.cruises[id-1].Enable = false

	if err := config.SetConfig(c.cruises.ConfigKey(), c.cruises); err != nil {
		c.cruises[id-1] = before
		return err
	}

	return nil
}
