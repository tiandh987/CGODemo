package idle

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

func (i *Idle) Set(motion *dsd.IdleMotion) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	idle := &Idle{}
	if err := idle.convert(motion); err != nil {
		return err
	}

	if err := config.SetConfig(motion.ConfigKey(), motion); err != nil {
		return err
	}

	i.motion = motion
	i.enable = idle.enable
	i.second = idle.second
	i.function = idle.function
	i.funcID = idle.funcID

	return nil
}

func (i *Idle) Get() *dsd.IdleMotion {
	i.mu.RLock()
	defer i.mu.RUnlock()

	motion := dsd.NewIdleMotion()
	motion.Enable = i.enable
	motion.Function = int(i.function)
	motion.Second = i.second

	switch i.function {
	case Preset:
		motion.PresetID = i.funcID
	case Cruise:
		motion.TourID = i.funcID
	case Trace:
		motion.PatternID = i.funcID
	case LineScan:
		motion.LinearScanID = i.funcID
	case RegionScan:
		motion.RegionScanID = i.funcID
	}

	return motion
}

func (i *Idle) GetSecond() int {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.second
}

func (i *Idle) Enable() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.enable
}

func (i *Idle) GetFuncAndId() (Function, int) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.function, i.funcID
}

func (i *Idle) Default() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	motion := dsd.NewIdleMotion()
	if err := config.SetConfig(motion.ConfigKey(), motion); err != nil {
		return err
	}

	i.enable = false
	i.second = 5
	i.function = None
	i.funcID = 1

	return nil
}

func (i *Idle) convert(motion *dsd.IdleMotion) error {
	var funcID int

	switch Function(motion.Function) {
	case None:
		funcID = 1
	case Preset:
		if err := dsd.PresetID(motion.PresetID).Validate(); err != nil {
			return err
		}
		funcID = motion.PresetID
	case Cruise:
		if err := dsd.CruiseID(motion.TourID).Validate(); err != nil {
			return err
		}
		funcID = motion.TourID
	case Trace:
		if err := dsd.TraceID(motion.PatternID).Validate(); err != nil {
			return err
		}
	case LineScan:
		if err := dsd.LineScanID(motion.LinearScanID).Validate(); err != nil {
			return err
		}
		funcID = motion.LinearScanID
	case RegionScan:
		//if err := dsd.PresetID(ups.PresetID).Validate(); err != nil {
		//	return err
		//}
	default:
		return errors.New("invalid idle function")
	}

	i.enable = motion.Enable
	i.second = motion.Second
	i.function = Function(motion.Function)
	i.funcID = funcID

	return nil
}
