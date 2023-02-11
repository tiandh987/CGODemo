package cruise

import (
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (c *Cruise) List() []dsd.TourPreset {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cruises
}

func (c *Cruise) Default() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state.cst != none {
		log.Warn("cruise is running")
		return errors.New("cruise is running")
	}

	before := c.cruises

	var cruise dsd.TourPreset
	c.cruises = cruise.Default()

	if err := c.saveConfig(); err != nil {
		c.cruises = before
		return err
	}

	return nil
}

func (c *Cruise) Update(id dsd.CruiseID, name dsd.CruiseName) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cruise, err := dsd.NewCruise(id, name)
	if err != nil {
		return err
	}

	if c.state.cst != none && c.cruises[id-1].Running == true {
		log.Warnf(fmt.Sprintf("cruise (%d) is running", id))
		return errors.New(fmt.Sprintf("cruise (%d) is running", id))
	}

	before := c.cruises[id-1]

	cruise.Enable = before.Enable
	cruise.Preset = before.Preset
	cruise.Running = before.Running

	c.cruises[id-1] = cruise
	if err := c.saveConfig(); err != nil {
		c.cruises[id-1] = before
		return err
	}
	return nil
}

func (c *Cruise) Set(cr *dsd.TourPreset) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state.cst != none && c.cruises[cr.ID-1].Running == true {
		return errors.New(fmt.Sprintf("cruise (%d) is running", cr.ID))
	}

	before := c.cruises[cr.ID-1]

	cr.Running = false
	if cr.Preset == nil {
		cr.Preset = []dsd.TourPresetPoint{}
	}
	c.cruises[cr.ID-1] = *cr

	if err := c.saveConfig(); err != nil {
		c.cruises[cr.ID-1] = before
		return err
	}

	return nil
}

func (c *Cruise) Delete(id dsd.CruiseID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state.cst != none && c.cruises[id-1].Running == true {
		log.Warn(fmt.Sprintf("cruise (%d) is running", id))
		return errors.New(fmt.Sprintf("cruise (%d) is running", id))
	}

	before := c.cruises[id-1]

	c.cruises[id-1].Preset = []dsd.TourPresetPoint{}
	c.cruises[id-1].Running = false

	if err := c.saveConfig(); err != nil {
		c.cruises[id-1] = before
		return err
	}

	return nil
}
