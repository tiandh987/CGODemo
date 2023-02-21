package idle

import (
	"errors"
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
}
