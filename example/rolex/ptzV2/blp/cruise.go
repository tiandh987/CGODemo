package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (b *Blp) ListCruise() []dsd.TourPreset {
	return b.cruise.List()
}

func (b *Blp) DefaultCruise() error {
	return b.cruise.Default()
}

func (b *Blp) UpdateCruise(id dsd.CruiseID, name dsd.CruiseName) error {
	return b.cruise.Update(id, name)
}

func (b *Blp) DeleteCruise(id dsd.CruiseID) error {
	return b.cruise.Delete(id)
}

func (b *Blp) SetCruise(cruise *dsd.TourPreset) error {
	return b.cruise.Set(cruise)
}
