package line

import (
	"context"
	"github.com/looplab/fsm"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sync"
	"time"
)

// 线扫状态
const (
	none           = "none"
	leftMargin     = "leftMargin"     // 左边界
	leftResidence  = "leftResidence"  // 左边界停留
	leftToRight    = "leftToRight"    // 左->右
	rightMargin    = "rightMargin"    // 右边界
	rightResidence = "rightResidence" // 右边界停留
	rightToLeft    = "rightToLeft"    // 右->左
	levelLeft      = "levelLeft"      // 水平旋转-逆时针(左/右边界未设置时)
)

type Line struct {
	mu    sync.RWMutex
	lines dsd.LineSlice

	fsm     *fsm.FSM
	eventCh chan string
	basic   *basic.Basic
}

func New(b *basic.Basic, s dsd.LineSlice) *Line {
	l := &Line{
		basic:   b,
		lines:   s,
		eventCh: make(chan string, 1),
	}

	l.fsm = fsm.NewFSM(
		none,
		fsm.Events{
			{Name: none, Src: []string{leftMargin, leftResidence, leftToRight, rightMargin, rightResidence,
				rightToLeft}, Dst: none},
			{Name: leftMargin, Src: []string{none, rightToLeft}, Dst: leftMargin},
			{Name: leftResidence, Src: []string{none, leftMargin}, Dst: leftResidence},
			{Name: leftToRight, Src: []string{leftResidence}, Dst: leftToRight},
			{Name: rightMargin, Src: []string{leftToRight}, Dst: rightMargin},
			{Name: rightResidence, Src: []string{rightMargin}, Dst: rightResidence},
			{Name: rightToLeft, Src: []string{rightResidence}, Dst: rightToLeft},

			{Name: levelLeft, Src: []string{none}, Dst: levelLeft},
		},
		fsm.Callbacks{
			"enter_none": func(ctx context.Context, event *fsm.Event) {
				log.Info("line scan enter_none")
				if l.fsm.Current() != "none" {
					l.basic.Stop()
				}
				return
			},
			"enter_leftMargin": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.LineScanID)

				if err := l.leftMargin(ctx, id); err != nil {
					log.Error(err.Error())
					l.fsm.Event(ctx, none)
					return
				}
			},
			"enter_leftResidence": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.LineScanID)

				if err := l.leftResidence(ctx, id); err != nil {
					log.Error(err.Error())
					l.fsm.Event(ctx, none)
					return
				}

			},
			"enter_leftToRight": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.LineScanID)

				if err := l.leftToRight(ctx, id); err != nil {
					log.Error(err.Error())
					l.fsm.Event(ctx, none)
					return
				}
			},
			"enter_rightMargin": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.LineScanID)

				log.Infof("line scan enter_rightMargin (%d)", id)

				if err := l.fsm.Event(ctx, rightResidence, id); err != nil {
					log.Error(err.Error())
					l.fsm.Event(ctx, none)
					return
				}
			},
			"enter_rightResidence": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.LineScanID)

				if err := l.rightResidence(ctx, id); err != nil {
					log.Error(err.Error())
					l.fsm.Event(ctx, none)
					return
				}
			},
			"enter_rightToLeft": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.LineScanID)

				if err := l.rightToLeft(ctx, id); err != nil {
					log.Error(err.Error())
					l.fsm.Event(ctx, none)
					return
				}
			},
			"enter_levelLeft": func(ctx context.Context, event *fsm.Event) {
				id := event.Args[0].(dsd.LineScanID)

				if err := l.basic.Operation(basic.DirectionLeft, ptz.Speed(l.lines[id-1].Speed)); err != nil {
					log.Error(err.Error())
					l.fsm.Event(ctx, none)
					return
				}
			},
		},
	)

	return l
}

func (l *Line) Start(ctx context.Context, id dsd.LineScanID) {
	if l.fsm.Current() != none {
		log.Warnf("line scan is running")
		return
	}

	if !l.lines[id-1].Enable {
		log.Warnf("line scan (%d) is disable")
		return
	}

	event := leftMargin
	if l.lines[id-1].LeftMargin == dsd.MarginNoLimit || l.lines[id-1].RightMargin == dsd.MarginNoLimit {
		event = levelLeft
	}

	if !l.fsm.Can(event) {
		log.Warnf("line scan can not convert to %s", event)
		return
	}

	go l.fsm.Event(ctx, event, id)

	l.lines[id-1].Running = true
}

func (l *Line) Stop(ctx context.Context, id dsd.LineScanID) {
	if l.fsm.Current() == none || !l.lines[id-1].Running {
		return
	}

	l.fsm.Event(ctx, none)
	l.lines[id-1].Running = false
}

func (l *Line) leftMargin(ctx context.Context, id dsd.LineScanID) error {
	log.Infof("line scan enter_leftMargin (%d)", id)

	line := l.lines[id-1]
	pos, err := l.basic.Position()
	if err != nil {
		return err
	}
	pos.Pan = line.LeftMargin

	if err := l.basic.Goto(pos); err != nil {
		return err
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc()

	if err := l.basic.ReachPosition(timeoutCtx, pos); err != nil {
		return err
	}

	if err := l.fsm.Event(ctx, leftResidence, id); err != nil {
		return err
	}

	return nil
}

func (l *Line) leftResidence(ctx context.Context, id dsd.LineScanID) error {
	log.Infof("line scan enter_leftResidence (%d)", id)

	time.Sleep(time.Duration(l.lines[id-1].ResidenceTimeLeft) * time.Second)

	if err := l.fsm.Event(ctx, leftToRight, id); err != nil {
		return err
	}

	return nil
}

func (l *Line) leftToRight(ctx context.Context, id dsd.LineScanID) error {
	log.Infof("line scan enter_leftToRight (%d)", id)

	line := l.lines[id-1]
	pos, err := l.basic.Position()
	if err != nil {
		return err
	}
	pos.Pan = line.RightMargin

	if err := l.basic.Operation(basic.DirectionRight, ptz.Speed(line.Speed)); err != nil {
		return err
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc()
	if err := l.basic.ReachPosition(timeoutCtx, pos); err != nil {
		return err
	}

	if err = l.fsm.Event(ctx, rightMargin, id); err != nil {
		return err
	}

	return nil
}

func (l *Line) rightResidence(ctx context.Context, id dsd.LineScanID) error {
	log.Infof("line scan enter_rightResidence (%d)", id)

	line := l.lines[id-1]
	time.Sleep(time.Duration(line.ResidenceTimeRight) * time.Second)

	if err := l.fsm.Event(ctx, rightToLeft, id); err != nil {
		return err
	}

	return nil
}

func (l *Line) rightToLeft(ctx context.Context, id dsd.LineScanID) error {
	log.Infof("line scan enter_rightToLeft (%d)", id)

	line := l.lines[id-1]
	pos, err := l.basic.Position()
	if err != nil {
		return err
	}
	pos.Pan = line.RightMargin

	if err := l.basic.Operation(basic.DirectionLeft, ptz.Speed(line.Speed)); err != nil {
		return err
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc()
	if err := l.basic.ReachPosition(timeoutCtx, pos); err != nil {
		return err
	}

	if err := l.fsm.Event(ctx, leftMargin, id); err != nil {
		return err
	}

	return nil
}
