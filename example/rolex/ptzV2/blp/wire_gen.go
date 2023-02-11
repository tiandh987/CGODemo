// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/arch/serial"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/cruise"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/lineScan"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

// Injectors from wire.go:

func New(limit *dsd.Limit, comName string, comCfg *dsd.PTZ, ps []dsd.PresetPoint, ls []dsd.LineScan, cs []dsd.TourPreset) *Blp {
	state := ptz.NewState(limit)
	serialSerial := serial.New(comName, comCfg)
	basic := basic.New()
	presetPreset := preset.New(ps)
	lineScanLineScan := lineScan.New(ls)
	cruiseCruise := cruise.New(cs)
	blp := NewBlp(state, serialSerial, basic, presetPreset, lineScanLineScan, cruiseCruise)
	return blp
}
