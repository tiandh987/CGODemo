package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/cruise"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/idle"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/line"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/powerUp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

type Blp struct {
	basic  *basic.Basic
	preset *preset.Preset
	line   *line.Line
	cruise *cruise.Cruise
	power  *powerUp.PowerUp
	idle   *idle.Idle
}

var _ PTZRepo = (*Blp)(nil)

type PTZRepo interface {
	Preset() PresetRepo
	Line() LineRepo
	Cruise() CruiseRepo
	Power() PowerRepo
	Idle() IdleRepo
}

type PresetRepo interface {
	List() dsd.PresetSlice
	Update(id dsd.PresetID, name string) error
	Delete(id dsd.PresetID) error
	DeleteAll() error
	Set(id dsd.PresetID, name string) error
}

func (b *Blp) Preset() PresetRepo {
	return b.preset
}

type LineRepo interface {
	List() dsd.LineSlice
	Default() error
	Set(scan *dsd.LineScan) error
	SetMargin(id dsd.LineScanID, op dsd.LineMarginOp) error
}

func (b *Blp) Line() LineRepo {
	return b.line
}

type CruiseRepo interface {
	List() dsd.CruiseSlice
	Default() error
	Update(id dsd.CruiseID, name string) error
	Set(cr *dsd.TourPreset) error
	Delete(id dsd.CruiseID) error
}

func (b *Blp) Cruise() CruiseRepo {
	return b.cruise
}

type PowerRepo interface {
	Get() *dsd.PowerUps
	Set(ups *dsd.PowerUps) error
	Default() error
}

func (b *Blp) Power() PowerRepo {
	return b.power
}

type IdleRepo interface {
	Get() *dsd.IdleMotion
	Set(motion *dsd.IdleMotion) error
	Default() error
}

func (b *Blp) Idle() IdleRepo {
	return b.idle
}

func New(basic *basic.Basic, preset *preset.Preset, line *line.Line, cruise *cruise.Cruise, up *powerUp.PowerUp,
	i *idle.Idle) PTZRepo {
	return &Blp{
		basic:  basic,
		preset: preset,
		line:   line,
		cruise: cruise,
		power:  up,
		idle:   i,
	}
}
