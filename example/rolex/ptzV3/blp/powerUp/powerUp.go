package powerUp

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
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
		log.Error(err.Error())
		up.enable = false
		up.function = None
		up.funcID = 1
	}

	return up
}
