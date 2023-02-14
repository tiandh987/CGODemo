package blp

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/powerUp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (b *Blp) GetPowerUp() (*dsd.PowerUps, error) {
	return b.power.Get()
}

func (b *Blp) SetPowerUp(ups *dsd.PowerUps) error {
	return b.power.Set(ups)
}

func (b *Blp) DefaultPowerUp() error {
	return b.power.Default()
}

func (b *Blp) StartPowerUp() error {
	log.Info("start ptz power up...")
	up, err := b.GetPowerUp()
	if err != nil {
		return err
	}

	if !up.Enable {
		log.Info("power up id disable")
		return nil
	}

	function := ptz.None
	funcID := 0

	switch powerUp.Function(up.Function) {
	case powerUp.None:
		// nothing to do
	case powerUp.Preset:
		function = ptz.Preset
		funcID = up.PresetID
	case powerUp.Cruise:
		function = ptz.Cruise
		funcID = up.TourID
	case powerUp.Trace:
		// TODO
	case powerUp.LineScan:
		function = ptz.LineScan
		funcID = up.LinearScanID
	case powerUp.RegionScan:
		// TODO
	default:
		return errors.New("invalid power up function")
	}

	if err := b.Control(ptz.PowerUp, function, funcID, 0, ptz.SpeedOne); err != nil {
		return err
	}

	return nil
}
