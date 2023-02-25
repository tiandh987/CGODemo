package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"time"
)

func (b *Blp) startCron() {
	cronCh := b.cron.Start()

	go func() {
		log.Info("start cron action detection process")

		for {
			select {
			case <-b.ctx.Done():
				b.cron.Stop()
				return
			case info := <-cronCh:
				st := b.state.getInternal()
				if st.function == None {
					if st.startTime.Add(time.Second * time.Duration(info.AutoHoming)).After(time.Now()) {
						log.Infof("cron action is triggered but autoHoming time is not up, noneTime: %s autoHoming: %d",
							st.startTime.Format("2006-01-02 15:04:05"), info.AutoHoming)
						continue
					}
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

				if st.trigger == CronTrigger && st.function == ability && st.funcID == info.FuncID {
					log.Debugf("cron(%d %d) is running, skipping", st.function, st.funcID)
					continue
				}

				if st.trigger.compare(CronTrigger) && st.function != None {
					continue
				}

				req := Request{
					Trigger: CronTrigger,
					Ability: ability,
					ID:      info.FuncID,
					Speed:   1,
				}

				log.Infof("trigger cron action, ability: %d, funcID: %d", ability, info.FuncID)

				if err := b.Start(&req); err != nil {
					log.Warnf("cron request: %+v, err: %s", req, err.Error())
				}
			}
		}
	}()
}
