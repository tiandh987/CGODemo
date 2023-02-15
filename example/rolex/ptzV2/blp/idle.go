package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/idle"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (b *Blp) GetIdle() (*dsd.IdleMotion, error) {
	return b.idle.Get()
}

func (b *Blp) SetIdle(motion *dsd.IdleMotion) error {
	return b.idle.Set(motion)
}

func (b *Blp) DefaultIdle() error {
	return b.idle.Default()
}

func (b *Blp) StartIdle() error {
	b.idle.Start()

	go func() {
		for {
			select {
			case info := <-b.idle.RunCh():
				log.Infof("receive idle action info: %+v", info)

				function := ptz.None

				switch info.Function {
				case idle.Preset:
					function = ptz.Preset
				case idle.Cruise:
					function = ptz.Cruise
				case idle.Trace:
					function = ptz.Trace
				case idle.LineScan:
					function = ptz.LineScan
				case idle.RegionScan:
					function = ptz.RegionScan
				default:
					continue
				}

				if err := b.Control(ptz.Idle, function, info.FuncID, 0, ptz.SpeedOne); err != nil {
					log.Error(err.Error())
					continue
				}
			}
		}
	}()

	return nil
}
