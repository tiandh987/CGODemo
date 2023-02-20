package preset

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sync"
)

type Preset struct {
	mu      sync.RWMutex
	presets dsd.PresetSlice

	basic *basic.Basic
}

func New(basic *basic.Basic, presets dsd.PresetSlice) *Preset {
	return &Preset{
		basic:   basic,
		presets: presets,
	}
}

func (p *Preset) Start(id dsd.PresetID) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	preset := p.presets[id-1]

	if !preset.Enable {
		log.Warnf("preset %d-%s is disable", preset.ID, preset.Name)
		return nil
	}

	if err := p.basic.Goto(&preset.Position); err != nil {
		return err
	}

	return nil
}
