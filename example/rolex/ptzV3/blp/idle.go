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
		log.Info("start idle action detection process")

		timer := time.NewTimer(time.Second)
		<-timer.C
		for {
			select {
			case <-b.ctx.Done():
				return
			case <-timer.C:
				function, id := b.idle.GetFuncAndId()

				if !b.idle.Enable() {
					log.Info("idle action triggered but disabled")
					continue
				}

				state := b.state.getInternal()
				if state.function != None {
					log.Info("idle action triggered but current state is not none")
					continue
				}

				if state.startTime.Add(time.Second * time.Duration(b.idle.GetSecond())).After(time.Now()) {
					log.Info("idle action is triggered but time is not up")
					continue
				}

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

				log.Infof("trigger idle action, none time: %s ability: %d, funcID: %d",
					state.startTime.Format("2006-01-02 15:04:05"), ability, id)

				if err := b.Start(&req); err != nil {
					log.Warnf("idle request: %+v, err: %s", req, err.Error())
				}
			case msg := <-noneTopic:
				msg.Ack()

				if !b.idle.Enable() {
					log.Info("idle action triggered but disabled")
					continue
				}

				sec := b.idle.GetSecond()
				timer.Reset(time.Second * time.Duration(sec))
			}
		}
	}()
}
