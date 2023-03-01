package trace

import (
	"context"
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
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
	mu       sync.RWMutex
	records  dsd.RecordSlice
	state    State
	index    int
	maxIndex int
	timer    *time.Timer
	replayCh chan struct{}
	wg       sync.WaitGroup
	cancel   context.CancelFunc

	tmpRecord dsd.Record
	basic     *basic.Basic
}

func New(basic *basic.Basic, records []dsd.Record) *Trace {
	return &Trace{
		records:  records,
		state:    none,
		index:    0,
		maxIndex: 0,
		replayCh: make(chan struct{}, 1),
		tmpRecord: dsd.Record{
			ID:        0,
			Enable:    false,
			Valid:     false,
			Schedules: []dsd.Schedule{},
		},
		basic: basic,
	}
}

func (t *Trace) List() []dsd.Record {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.records
}

func (t *Trace) StartRecord(id dsd.TraceID) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := id.Validate(); err != nil {
		return err
	}

	if t.state == tracking {
		return errors.New("trace is running")
	}

	if t.state == recording {
		return errors.New("trace is recording")
	}

	position, err := t.basic.Position()
	if err != nil {
		return err
	}

	t.state = recording
	t.tmpRecord = dsd.Record{
		ID:            id,
		Enable:        true,
		Valid:         false,
		Schedules:     []dsd.Schedule{},
		StartPosition: *position,
	}

	return nil
}

func (t *Trace) StopRecord(id dsd.TraceID) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != recording {
		return
	}

	if t.tmpRecord.ID != id {
		return
	}

	before := t.records

	t.tmpRecord.Valid = true
	t.records[id-1] = t.tmpRecord
	if err := config.SetConfig(t.records.ConfigKey(), t.records); err != nil {
		t.records = before
		return
	}

	t.state = none

	return
}

func (t *Trace) Record(s *dsd.Schedule) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != recording {
		return
	}

	t.tmpRecord.Schedules = append(t.tmpRecord.Schedules, *s)
}

func (t *Trace) Start(id dsd.TraceID) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if err := id.Validate(); err != nil {
		return err
	}

	if t.state != none {
		return errors.New("trace is running")
	}

	record := t.records[id-1]

	if !record.Enable || !record.Valid || len(record.Schedules) == 0 {
		return errors.New("trace record is invalid")
	}

	t.state = tracking
	t.index = 0
	t.maxIndex = len(record.Schedules)
	t.replayCh = make(chan struct{}, 1)

	log.Infof("aaaaa")
	t.timer = time.NewTimer(time.Millisecond * 10)
	<-t.timer.C
	log.Infof("bbbbb")

	ctx := context.Background()
	ctx, t.cancel = context.WithCancel(ctx)

	startCh := make(chan struct{})
	restartCh := make(chan struct{}, 1)

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()

		log.Infof("start trace id: %d enable: %t valid: %t start position: %+v index: %d maxIndex: %d",
			record.ID, record.Enable, record.Valid, record.StartPosition, t.index, t.maxIndex)
		printTraceInfo(record.Schedules)

		for {
			select {
			case <-ctx.Done():
				goto EndTrace
			case <-startCh:
				log.Infof("startCh receive index: %d", t.index)
				t.gotoStart(ctx, &record.StartPosition)
				log.Infof("startCh end")
			case <-restartCh:
				log.Infof("restartCh receive index: %d", t.index)
				t.gotoStart(ctx, &record.StartPosition)
				log.Infof("restartCh end")

			case <-t.timer.C:
				log.Infof("timer.C receive index: %d", t.index)

				if err := t.basic.Stop(); err != nil {
					log.Error(err.Error())
				}

				t.index++
				if t.index >= t.maxIndex {
					t.index = 0
					restartCh <- struct{}{}
					continue
				}

				log.Infof("timer.C send to replayCh")

				t.replayCh <- struct{}{}
				log.Infof("timer.C send to replayCh end indx: %d", t.index)

			case <-t.replayCh:
				log.Infof("replayCh receive index: %d", t.index)

				schedule := record.Schedules[t.index]
				if err := t.basic.Operation(basic.Operation(schedule.FuncID), ptz.Speed(schedule.Speed)); err != nil {
					log.Error(err.Error())
				}

				log.Infof("start: %s end : %s duration: %d",
					schedule.StartTime, schedule.StopTime, schedule.StopTime.Sub(schedule.StartTime))

				t.timer.Reset(schedule.StopTime.Sub(schedule.StartTime))

				log.Infof("replayCh end")

			}
		}
	EndTrace:
		log.Infof("end trace %d", record.ID)
		t.basic.Stop()
		t.state = none
		t.index = 0
		t.timer.Stop()
		close(t.replayCh)
	}()

	startCh <- struct{}{}

	return nil
}

func (t *Trace) Stop(id dsd.TraceID) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != tracking {
		return errors.New("trace is not tracking")
	}

	t.cancel()
	t.wg.Wait()

	return nil
}

func (t *Trace) IsRecording() bool {
	return t.state == recording
}

func (t *Trace) gotoStart(ctx context.Context, position *dsd.Position) {
	log.Infof("position: %+v", position)

	if err := t.basic.Goto(position); err != nil {
		log.Error(err.Error())
		return
	}

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		defer log.Infof("sssssssssssssssssss")

		timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		ticker := time.NewTicker(time.Millisecond * 100)

		for {
			select {
			case <-timeoutCtx.Done():
				log.Error(timeoutCtx.Err().Error())
				return
			case <-ticker.C:
				curPos, err := t.basic.Position()
				log.Infof("current position: %+v", curPos)
				if err != nil {
					log.Error(err.Error())
					time.Sleep(time.Second)
					continue
				}

				if curPos.Pan >= position.Pan-2 && curPos.Pan <= position.Pan+2 &&
					curPos.Tilt >= position.Tilt-2 && curPos.Tilt <= position.Tilt+2 &&
					curPos.Zoom >= position.Zoom-2 && curPos.Zoom <= position.Zoom+2 {
					log.Infof("xxxxxxxxx stop1")
					t.basic.Stop()
					t.replayCh <- struct{}{}
					log.Infof("xxxxxxxxx stop2")
					return
				}
			}
		}
	}()
}

func printTraceInfo(schedules []dsd.Schedule) {
	for i, schedule := range schedules {
		log.Infof("index: %d funcID: %d speed: %d duration: %d",
			i, schedule.FuncID, schedule.Speed, schedule.StopTime.Sub(schedule.StartTime))
	}
}
