//go:build wireinject
// +build wireinject

package blp

import (
	"github.com/google/wire"
	"rolex/ptz/arch/serial"
	"rolex/ptz/blp/lineScan"
	"rolex/ptz/blp/preset"
	"rolex/ptz/blp/ptz"
	"rolex/ptz/dsd"
)

func New(limit *dsd.Limit, comName string, comCfg *dsd.PTZ, ps []dsd.PresetPoint, ls []dsd.LineScan) *Blp {
	panic(wire.Build(NewBlp, ptz.NewState, serial.New, ptz.NewBasic, preset.New, lineScan.New))
}
