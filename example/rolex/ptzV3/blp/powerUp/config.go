package powerUp

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

func (u *PowerUp) Set(ups *dsd.PowerUps) error {
	if err := u.convert(ups); err != nil {
		return err
	}

	if err := config.SetConfig(ups.ConfigKey(), ups); err != nil {
		return err
	}

	return nil
}

func (u *PowerUp) Get() *dsd.PowerUps {
	ups := dsd.NewPowerUps()
	ups.Enable = u.enable
	ups.Function = int(u.function)

	switch u.function {
	case Preset:
		ups.PresetID = u.funcID
	case Cruise:
		ups.TourID = u.funcID
	case Trace:
		ups.PatternID = u.funcID
	case LineScan:
		ups.LinearScanID = u.funcID
	case RegionScan:
		ups.RegionScanID = u.funcID
	}

	return ups
}

func (u *PowerUp) Enable() bool {
	return u.enable
}

func (u *PowerUp) GetFuncAndId() (Function, int) {
	return u.function, u.funcID
}

func (u *PowerUp) Default() error {
	u.enable = false
	u.function = None
	u.funcID = 1

	ups := dsd.NewPowerUps()
	if err := config.SetConfig(ups.ConfigKey(), ups); err != nil {
		return err
	}

	return nil
}

func (u *PowerUp) convert(ups *dsd.PowerUps) error {
	var funcID int

	switch Function(ups.Function) {
	case None:
		funcID = 1
	case Preset:
		if err := dsd.PresetID(ups.PresetID).Validate(); err != nil {
			return err
		}
		funcID = ups.PresetID
	case Cruise:
		if err := dsd.CruiseID(ups.TourID).Validate(); err != nil {
			return err
		}
		funcID = ups.TourID
	case Trace:
		if err := dsd.TraceID(ups.PatternID).Validate(); err != nil {
			return err
		}
	case LineScan:
		if err := dsd.LineScanID(ups.LinearScanID).Validate(); err != nil {
			return err
		}
		funcID = ups.LinearScanID
	case RegionScan:
	//if err := dsd.PresetID(ups.PresetID).Validate(); err != nil {
	//	return err
	//}
	default:
		return errors.New("invalid power up function")
	}

	u.enable = ups.Enable
	u.function = Function(ups.Function)
	u.funcID = funcID

	return nil
}
