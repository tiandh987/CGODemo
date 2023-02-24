package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/powerUp"
)

func (b *Blp) startPowerUp() {
	if !b.power.Enable() {
		log.Info("power up is disable")

		return
	}

	function, id := b.power.GetFuncAndId()
	log.Infof("start ptz power up (%d %d)", function, id)

	ability := None

	switch function {
	case powerUp.Preset:
		ability = Preset
	case powerUp.Cruise:
		ability = Cruise
	case powerUp.Trace:
		ability = Trace
	case powerUp.LineScan:
		ability = LineScan
	case powerUp.RegionScan:
		ability = RegionScan
	default:
		log.Warnf("invalid power up function (%d)", function)
		return
	}

	req := Request{
		Trigger: PowerUpTrigger,
		Ability: ability,
		ID:      id,
		Speed:   1,
	}

	if err := b.Start(&req); err != nil {
		log.Warnf("power up request: %+v, err: %s", req, err.Error())
	}

	return
}
