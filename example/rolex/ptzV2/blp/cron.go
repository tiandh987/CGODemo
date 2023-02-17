package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/cron"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (b *Blp) ListCron() []dsd.PtzAutoMovement {
	return b.cron.List()
}

func (b *Blp) DefaultCron() error {
	return b.cron.Default()
}

func (b *Blp) SetCron(movement *dsd.PtzAutoMovement) error {
	return b.cron.Set(movement)
}

func (b *Blp) startCron() {
	b.cron.Start()

	go func() {
		for {
			select {
			case <-b.cron.QuitCh():
				return
			case info := <-b.cron.InfoCh():
				var function ptz.Function
				switch info.Function {
				case cron.Preset:
					function = ptz.Preset
				case cron.Cruise:
					function = ptz.Cruise
				case cron.Trace:
					function = ptz.Trace
				case cron.LineScan:
					function = ptz.LineScan
				case cron.RegionScan:
					function = ptz.RegionScan
				}

				if err := b.Control(ptz.Cron, function, info.FuncID, info.CronID, ptz.SpeedOne); err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}()
}

func (b *Blp) quitCron() {
	b.cron.Quit()
}
