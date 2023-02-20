package trace

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"sort"
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

func New(records []Record) *Trace {
	return &Trace{
		records: records,
		runState: &runState{
			recordID: 0,
			state:    none,
			index:    0,
			maxIndex: 0,
			timer:    time.NewTimer(time.Second),
		},
		tmpRecord: Record{},
		infoCh:    make(chan ScheduleInfo, 1),
		stopCh:    make(chan struct{}, 1),
		quitCh:    make(chan struct{}, 1),
	}
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

func (t *Trace) List() []Record {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.records
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
		ID:        id,
		Enable:    true,
		Valid:     false,
		Infos:     []ScheduleInfo{},
		Start:     *pos,
		StartTime: time.Now(),
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

	t.tmpRecord.Valid = true
	t.tmpRecord.End = *pos
	t.tmpRecord.EndTime = time.Now()

	sort.Slice(t.tmpRecord, func(i, j int) bool {
		return t.tmpRecord.Infos[i].StartTime.Sub(t.tmpRecord.Infos[j].StartTime) < 0
	})

	for i := 0; i < len(t.tmpRecord.Infos); i++ {
		if i == len(t.tmpRecord.Infos)-1 {
			t.tmpRecord.Infos[i].duration = t.tmpRecord.EndTime.Sub(t.tmpRecord.Infos[i].StartTime)
			continue
		}

		t.tmpRecord.Infos[i].duration = t.tmpRecord.Infos[i+1].StartTime.Sub(t.tmpRecord.Infos[i].StartTime)
	}
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
				t.runState.timer.Reset(time.Second * 3)
				t.runState.index = 0
				continue
			}

			t.infoCh <- record.Infos[t.runState.index]
			t.runState.timer.Reset(time.Second * record.Infos[t.runState.index].duration)
			t.runState.index++
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

func (t *Trace) InfoCh() <-chan ScheduleInfo {
	return t.infoCh
}
