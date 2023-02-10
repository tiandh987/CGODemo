package preset

import (
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (p *Preset) List() []dsd.PresetPoint {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.presets
}

func (p *Preset) Update(id dsd.PresetID, name dsd.PresetName) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	before := p.presets[id-1]
	preset, err := dsd.NewPreset(id, name)
	if err != nil {
		return err
	}
	preset.Enable = before.Enable
	preset.Position = before.Position

	p.presets[id-1] = preset

	log.Debugf("before[%p]: %+v\npreset[%p]: %+v\npresets[%d][%p]: %+v\n",
		&before, before, &preset, preset, id-1, &(p.presets[id-1]), p.presets[id-1])

	if err := config.SetConfig(preset.ConfigKey(), p.presets); err != nil {
		p.presets[id-1] = before
		return err
	}
	return nil
}

func (p *Preset) Delete(id dsd.PresetID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	before := p.presets[id-1]
	p.presets[id-1].Enable = false
	if err := config.SetConfig(p.presets[id-1].ConfigKey(), p.presets); err != nil {
		p.presets[id-1] = before
		return err
	}
	return nil
}

func (p *Preset) DeleteAll() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	before := p.presets

	for i, _ := range p.presets {
		p.presets[i].Enable = false
	}
	if err := config.SetConfig(p.presets[0].ConfigKey(), p.presets); err != nil {
		p.presets = before
		return err
	}
	return nil
}

func (p *Preset) Set(ctl control.ControlRepo, id dsd.PresetID, name dsd.PresetName) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	preset, err := dsd.NewPreset(id, name)
	if err != nil {
		return err
	}
	preset.Enable = true

	pos, err := ctl.Position()
	if err != nil {
		return err
	}
	preset.Position = pos

	before := p.presets[id-1]
	p.presets[id-1] = preset
	if err := config.SetConfig(preset.ConfigKey(), p.presets); err != nil {
		p.presets[id-1] = before
		return err
	}
	return nil
}
