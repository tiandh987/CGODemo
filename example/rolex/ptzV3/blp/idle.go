package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/idle"
	"time"
)

func (b *Blp) startIdle() {
	noneTopic, err := b.state.subscribe(b.ctx, stateNoneTopic)
	if err != nil {
		log.Panic(err.Error())
		return
	}

	go func() {
		timer := time.NewTimer(time.Second)
		<-timer.C
		for {
			select {
			case <-b.ctx.Done():
				return
			case <-timer.C:
				function, id := b.idle.GetFuncAndId()

				ability := None
				switch function {
				case idle.Preset:
					ability = Preset
				case idle.Cruise:
					ability = Cruise
				case idle.Trace:
					ability = Trace
				case idle.LineScan:
					ability = LineScan
				case idle.RegionScan:
					ability = RegionScan
				default:
					continue
				}

				req := Request{
					Trigger: IdleTrigger,
					Ability: ability,
					ID:      id,
					Speed:   1,
				}

				if err := b.Start(&req); err != nil {
					log.Warnf("idle request: %+v, err: %s", req, err.Error())
				}
			case <-noneTopic:
				if !b.idle.Enable() {
					continue
				}
				sec := b.idle.GetSecond()
				timer.Reset(time.Second * time.Duration(sec))
			}
		}
	}()
}
