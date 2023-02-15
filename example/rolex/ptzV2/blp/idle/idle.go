package idle

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"time"
)

// Function 开机功能
type Function int

const (
	None       Function = iota // None
	Preset                     // 预置点
	Cruise                     // 巡航
	Trace                      // 巡迹
	LineScan                   // 线性扫描
	RegionScan                 // 区域扫描
)

func (f Function) Validate() error {
	if f < None || f > RegionScan {
		return errors.New("invalid idle function")
	}

	return nil
}

type Idle struct {
	enable   bool
	second   int
	function Function
	funcID   int

	running bool
	timer   *time.Timer
	runCh   chan RunInfo
	pauseCh chan struct{}
	quitCh  chan struct{}
}

func New(motion *dsd.IdleMotion) *Idle {
	idle := &Idle{
		running: false,
		timer:   time.NewTimer(time.Hour),
		runCh:   make(chan RunInfo, 1),
		pauseCh: make(chan struct{}, 1),
		quitCh:  make(chan struct{}, 1),
	}

	if err := idle.convert(motion); err != nil {
		log.Panic(err.Error())
	}

	return idle
}

func (i *Idle) Set(motion *dsd.IdleMotion) error {
	if i.running {
		log.Warnf("idle (%d-%d) is running", i.function, i.funcID)
		return errors.New("idle is running")
	}

	if err := i.convert(motion); err != nil {
		return err
	}

	if err := config.SetConfig(motion.ConfigKey(), motion); err != nil {
		return err
	}

	return nil
}

func (i *Idle) Get() (*dsd.IdleMotion, error) {
	motion := dsd.NewIdleMotion()
	motion.Enable = i.enable
	motion.Function = int(i.function)

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

	motion.Runing = false
	motion.RunningFunction = int(None)
	if i.running {
		motion.Runing = true
		motion.RunningFunction = int(i.function)
	}

	return motion, nil
}

func (i *Idle) Default() error {
	if i.running {
		log.Warnf("idle (%d-%d) is running", i.function, i.funcID)
		return errors.New("idle is running")
	}

	i.enable = false
	i.second = 5
	i.function = None
	i.funcID = 1
	i.running = false

	motion := dsd.NewIdleMotion()
	if err := config.SetConfig(motion.ConfigKey(), motion); err != nil {
		return err
	}

	return nil
}

func (i *Idle) convert(motion *dsd.IdleMotion) error {
	var funcID int

	switch Function(motion.Function) {
	case None:
		// nothing to do
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
		//if err := dsd.Trace(ups.PresetID).Validate(); err != nil {
		//	return err
		//}
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

func (i *Idle) Start() {
	log.Info("start idle detection...")
	go func() {
		for {
			select {
			case <-i.quitCh:
				log.Info("end idle detection...")
				goto EndIdle
			case <-i.pauseCh:
				log.Info("pause idle detection...")
				i.running = false
				i.timer.Stop()
			case <-i.timer.C:
				if i.running || i.function == None {
					continue
				}

				log.Infof("send idle action info (%d-%d)", i.function, i.funcID)
				i.runCh <- RunInfo{
					Function: i.function,
					FuncID:   i.funcID,
				}
				i.running = true
			}
		}
	EndIdle:
		i.running = false
		i.timer.Stop()
		close(i.runCh)
		close(i.pauseCh)
		close(i.quitCh)
		return
	}()
}

func (i *Idle) Reset() {
	i.timer.Reset(time.Second * time.Duration(i.second))
}

func (i *Idle) Pause() {
	if i.running {
		i.pauseCh <- struct{}{}
	}
}

func (i *Idle) Stop() {
	if i.running {
		i.quitCh <- struct{}{}
	}
}

func (i *Idle) RunCh() <-chan RunInfo {
	return i.runCh
}

type RunInfo struct {
	Function Function
	FuncID   int
}
