package preset

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"sync"
)

type Preset struct {
	mu      sync.RWMutex
	presets []dsd.PresetPoint
}

func New(ps []dsd.PresetPoint) *Preset {
	return &Preset{
		presets: ps,
	}
}

func (p *Preset) Start(ctl control.ControlRepo, id dsd.PresetID) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if err := id.Validate(); err != nil {
		return err
	}

	preset := p.presets[id-1]

	log.Debugf("param id: %d, preset: %+v, position: %+v", id, preset, preset.Position)

	if !preset.Enable {
		log.Infof("preset %d-%s is disable", preset.ID, preset.Name)
		return nil
	}

	if err := ctl.Goto(preset.Position); err != nil {
		return err
	}

	return nil
}

func (p *Preset) Stop() error {
	return nil
}
