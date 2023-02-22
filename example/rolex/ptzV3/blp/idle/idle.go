package idle

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sync"
)

// Function 空闲功能
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
	mu       sync.RWMutex
	motion   *dsd.IdleMotion
	enable   bool
	second   int
	function Function
	funcID   int
}

func New(motion *dsd.IdleMotion) *Idle {
	i := &Idle{
		motion: motion,
	}

	if err := i.convert(motion); err != nil {
		log.Error(err.Error())

		i.enable = false
		i.second = 0
		i.function = 0
		i.funcID = 0
	}

	return i
}
