package preset

import (
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

func (p *Preset) List() dsd.PresetSlice {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.presets
}

func (p *Preset) Update(id dsd.PresetID, name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	preset := dsd.NewPreset(id, name)
	if err := preset.Validate(); err != nil {
		return err
	}

	before := p.presets[id-1]

	preset.Enable = before.Enable
	preset.Position = before.Position

	p.presets[id-1] = preset

	if err := config.SetConfig(p.presets.ConfigKey(), p.presets); err != nil {
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
	if err := config.SetConfig(p.presets.ConfigKey(), p.presets); err != nil {
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

	if err := config.SetConfig(p.presets.ConfigKey(), p.presets); err != nil {
		p.presets = before
		return err
	}

	return nil
}

func (p *Preset) Set(id dsd.PresetID, name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	preset := dsd.NewPreset(id, name)
	if err := preset.Validate(); err != nil {
		return err
	}
	preset.Enable = true

	pos, err := p.basic.Position()
	if err != nil {
		return err
	}
	preset.Position = *pos

	before := p.presets[id-1]

	p.presets[id-1] = preset
	if err := config.SetConfig(p.presets.ConfigKey(), p.presets); err != nil {
		p.presets[id-1] = before
		return err
	}

	return nil
}
