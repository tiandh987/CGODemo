package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"time"
)

func (b *Blp) startCron() {
	noneTopic, err := b.state.subscribe(b.ctx, stateNoneTopic)
	if err != nil {
		log.Panic(err.Error())
		return
	}

	cronCh := b.cron.Start()

	go func() {
		noneTime := time.Now()
		for {
			select {
			case <-b.ctx.Done():
				b.cron.Stop()
				return
			case msg := <-noneTopic:
				parse, err := time.Parse("2006-01-02 15:04:05", string(msg.Payload))
				if err != nil {
					log.Error(err.Error())
					continue
				}
				noneTime = parse
			case info := <-cronCh:
				if time.Now().Sub(noneTime) < info.AutoHoming {
					log.Infof("the difference between the current time(%s) and noneTime(%s) is less than autoHoming(%d)",
						time.Now().String(), noneTime.String(), info.AutoHoming)
					continue
				}

				ability := None
				switch info.Function {
				case dsd.Preset:
					ability = Preset
				case dsd.Cruise:
					ability = Cruise
				case dsd.Trace:
					ability = Trace
				case dsd.Line:
					ability = LineScan
				case dsd.Region:
					ability = RegionScan
				default:
					continue
				}

				req := Request{
					Trigger: CronTrigger,
					Ability: ability,
					ID:      info.FuncID,
					Speed:   1,
				}

				if err := b.Start(&req); err != nil {
					log.Warnf("cron request: %+v, err: %s", req, err.Error())
				}
			}
		}
	}()
}
