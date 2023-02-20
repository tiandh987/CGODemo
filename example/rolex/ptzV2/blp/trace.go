package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/trace"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (b *Blp) ListTrace() []trace.Record {
	return b.trace.List()
}

func (b *Blp) DefaultTrace() error {
	return nil
}

func (b *Blp) StartRecord(id dsd.TraceID) error {
	pos, _ := b.getControl().Position()
	return b.trace.StartRecord(id, pos)
}

func (b *Blp) StopRecord(id dsd.TraceID) {
	pos, _ := b.getControl().Position()
	b.trace.StopRecord(id, pos)
}

func (b *Blp) Start(id dsd.TraceID) {
	b.trace.Start(b.getControl(), id)
}

func (b *Blp) startTrace() {
	go func() {
		for {
			select {
			case info := <-b.trace.InfoCh():
				if err := b.turn(info.FuncID, info.Speed); err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}()
}

func (b *Blp) quitTrace() {
	b.trace.Quit()
}
