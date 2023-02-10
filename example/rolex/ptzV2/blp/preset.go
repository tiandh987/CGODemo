package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (b *Blp) ListPreset() []dsd.PresetPoint {
	return b.preset.List()
}

func (b *Blp) UpdatePreset(id dsd.PresetID, name dsd.PresetName) error {
	return b.preset.Update(id, name)
}

func (b *Blp) DeletePreset(id dsd.PresetID) error {
	return b.preset.Delete(id)
}

func (b *Blp) DeleteAllPreset() error {
	return b.preset.DeleteAll()
}

func (b *Blp) SetPreset(id dsd.PresetID, name dsd.PresetName) error {
	return b.preset.Set(b.getControl(), id, name)
}
