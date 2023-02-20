package trace

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"sync"
	"time"
)

type State int

const (
	none      State = iota
	tracking        // 巡迹中
	recording       // 巡迹记录中
)

type Trace struct {
	mu        sync.RWMutex
	records   []Record
	runState  *runState
	tmpRecord Record

	infoCh chan ScheduleInfo
	stopCh chan struct{}
	quitCh chan struct{}
}

type runState struct {
	recordID dsd.TraceID
	state    State
	index    int
	maxIndex int
	timer    *time.Timer
}

type Record struct {
	ID        dsd.TraceID
	Enable    bool
	Valid     bool
	Infos     []ScheduleInfo
	Start     dsd.Position
	End       dsd.Position
	StartTime time.Time
	EndTime   time.Time
}

type ScheduleInfo struct {
	FuncID    int
	Speed     ptz.Speed
	StartTime time.Time
	duration  time.Duration
}

func (t *Trace) StartRecord(id dsd.TraceID, pos *dsd.Position) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.runState.state != none {
		return errors.New("trace is running")
	}

	t.runState.recordID = id
	t.runState.state = recording
	t.tmpRecord = Record{
		ID:     id,
		Enable: true,
		Valid:  false,
		Infos:  []ScheduleInfo{},
		Start:  *pos,
	}

	return nil
}

func (t *Trace) StopRecord(id dsd.TraceID, pos *dsd.Position) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.runState.state != recording {
		return
	}

	if t.runState.recordID != id {
		return
	}

	t.tmpRecord.End = *pos
	t.tmpRecord.Valid = true

	t.records[id-1] = t.tmpRecord

	return
}

func (t *Trace) Record(info *ScheduleInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.runState.state != recording {
		return
	}

	t.tmpRecord.Infos = append(t.tmpRecord.Infos, *info)
}

func (t *Trace) Start(ctl control.ControlRepo, id dsd.TraceID) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.runState.state != none {
		return
	}

	record := t.records[id-1]

	t.runState.state = tracking
	t.runState.recordID = id
	t.runState.index = 0
	t.runState.maxIndex = len(record.Infos)
	t.runState.timer = time.NewTimer(time.Second)

	for {
		select {
		case <-t.stopCh:
			return
		case <-t.quitCh:
			return
		case <-t.runState.timer.C:
			if t.runState.index == 0 {
				ctl.Goto(&record.Start)
				time.Sleep(record.Infos[t.runState.index].StartTime.Sub(record.StartTime) * time.Second)
			}

			if t.runState.index == t.runState.maxIndex {
				ctl.Goto(&record.End)
				time.Sleep(record.Infos[t.runState.index].StartTime.Sub(record.StartTime) * time.Second)
			}

			t.infoCh <- record.Infos[t.runState.index]
			t.runState.timer.Reset(time.Second * record.Infos[t.runState.index].duration)
			t.runState.index++

			if t.runState.index+1 < t.runState.index {
				ctl.Goto(&record.Start)
				time.Sleep(time.Millisecond * 100)
			}

			info := record.Infos[t.runState.index]
			t.infoCh <- info

			if t.runState.index+1 == t.runState.maxIndex {

			}
		}
	}

}

func (t *Trace) Quit() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.runState.state != none {
		t.quitCh <- struct{}{}
	}
}
