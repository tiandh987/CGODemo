package powerUp

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

// Function 开机功能
type Function int

const (
	None       Function = iota // 自动
	Preset                     // 预置点
	Cruise                     // 巡航
	Trace                      // 巡迹
	LineScan                   // 线性扫描
	RegionScan                 // 区域扫描
)

func (f Function) Validate() error {
	if f < None || f > RegionScan {
		return errors.New("invalid power up function")
	}

	return nil
}

type PowerUp struct {
	enable   bool
	function Function
	funcID   int
}

func New(ups *dsd.PowerUps) *PowerUp {
	up := &PowerUp{}

	if err := up.convert(ups); err != nil {
		log.Panic(err.Error())
	}

	return up
}

func (u *PowerUp) Set(ups *dsd.PowerUps) error {
	if err := u.convert(ups); err != nil {
		return err
	}

	if err := config.SetConfig(ups.ConfigKey(), ups); err != nil {
		return err
	}

	return nil
}

func (u *PowerUp) Get() (*dsd.PowerUps, error) {
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

	return ups, nil
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
		// nothing to do
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
		//if err := dsd.Trace(ups.PresetID).Validate(); err != nil {
		//	return err
		//}
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
