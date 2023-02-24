package line

import (
	"context"
	"errors"
	"fmt"
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

	fsm    *fsm.FSM
	fsmCh  chan string
	basic  *basic.Basic
	timer  *time.Timer
	stopCh chan dsd.LineScanID
	ctx    context.Context
	cancel context.CancelFunc
}

func New(b *basic.Basic, s dsd.LineSlice) *Line {
	l := &Line{
		basic:  b,
		lines:  s,
		fsmCh:  make(chan string, 1),
		timer:  time.NewTimer(time.Second),
		stopCh: make(chan dsd.LineScanID, 1),
	}

	l.fsm = fsm.NewFSM(
		none,
		fsm.Events{
			{Name: none, Src: []string{leftMargin, leftResidence, leftToRight, rightMargin, rightResidence,
				rightToLeft, levelLeft}, Dst: none},
			{Name: leftMargin, Src: []string{none, rightToLeft}, Dst: leftMargin},
			{Name: leftResidence, Src: []string{leftMargin}, Dst: leftResidence},
			{Name: leftToRight, Src: []string{leftResidence}, Dst: leftToRight},
			{Name: rightMargin, Src: []string{leftToRight}, Dst: rightMargin},
			{Name: rightResidence, Src: []string{rightMargin}, Dst: rightResidence},
			{Name: rightToLeft, Src: []string{rightResidence}, Dst: rightToLeft},

			{Name: levelLeft, Src: []string{none}, Dst: levelLeft},
		},
		fsm.Callbacks{
			"enter_none": func(ctx context.Context, event *fsm.Event) {
				l.basic.Stop()
				return
			},

			"enter_leftMargin": func(ctx context.Context, event *fsm.Event) {
				l.fsmCh <- leftResidence
			},
			"enter_leftResidence": func(ctx context.Context, event *fsm.Event) {
				line := event.Args[0].(dsd.LineScan)
				if err := l.leftResidence(line); err != nil {
					log.Error(err.Error())
					return
				}
			},
			"enter_leftToRight": func(ctx context.Context, event *fsm.Event) {
				line := event.Args[0].(dsd.LineScan)
				if err := l.leftToRight(line); err != nil {
					log.Error(err.Error())
					return
				}
				//l.fsmCh <- rightMargin
			},
			"enter_rightMargin": func(ctx context.Context, event *fsm.Event) {
				l.fsmCh <- rightResidence
			},
			"enter_rightResidence": func(ctx context.Context, event *fsm.Event) {
				line := event.Args[0].(dsd.LineScan)
				if err := l.rightResidence(line); err != nil {
					log.Error(err.Error())
					return
				}
			},
			"enter_rightToLeft": func(ctx context.Context, event *fsm.Event) {
				line := event.Args[0].(dsd.LineScan)
				if err := l.rightToLeft(line); err != nil {
					log.Error(err.Error())
					return
				}
				//l.fsmCh <- leftMargin
			},
			"enter_levelLeft": func(ctx context.Context, event *fsm.Event) {
				line := event.Args[0].(dsd.LineScan)
				log.Infof("enter_levelLeft (%d)", line.ID)

				if err := l.basic.Operation(basic.DirectionLeft, ptz.Speed(line.Speed)); err != nil {
					log.Error(err.Error())
					return
				}
			},
		},
	)

	return l
}

func (l *Line) Start(ctx context.Context, id dsd.LineScanID) error {
	if err := id.Validate(); err != nil {
		return err
	}

	if l.fsm.Current() != none {
		log.Warnf("line scan is running")
		return errors.New("line scan is running")
	}

	if !l.lines[id-1].Enable {
		log.Warnf("line scan (%d) is disable", id)
		return fmt.Errorf("line scan (%d) is disable", id)
	}

	go func(id dsd.LineScanID) {
		line := l.lines[id-1]
		l.ctx, l.cancel = context.WithCancel(context.Background())
		l.fsmCh = make(chan string, 1)
		l.stopCh = make(chan dsd.LineScanID, 1)

		log.Infof("start line scan id: %d left: %f %ds right: %f %ds", id, line.LeftMargin, line.ResidenceTimeLeft,
			line.RightMargin, line.ResidenceTimeRight)

		if line.LeftMargin == dsd.MarginNoLimit || line.RightMargin == dsd.MarginNoLimit {
			l.fsmCh <- levelLeft
		} else {
			if err := l.gotoStartPosition(line); err != nil {
				log.Error(err.Error())
				goto EndLine
			}
			l.fsmCh <- leftMargin
		}

		for {
			l.lines[id-1].Running = true

			select {
			case <-ctx.Done():
				log.Warn(ctx.Err().Error())
				goto EndLine
			case stopId := <-l.stopCh:
				log.Infof("line receive stop id: %d", stopId)
				if stopId != id {
					log.Warnf("current id (%d), request id (%d)", id, stopId)
					continue
				}
				goto EndLine
			case <-l.timer.C:
				switch l.fsm.Current() {
				case leftResidence:
					l.fsmCh <- leftToRight
				case rightResidence:
					l.fsmCh <- rightToLeft
				}
			case state := <-l.fsmCh:
				switch state {
				case leftMargin:
					if err := l.fsm.Event(ctx, leftMargin, line); err != nil {
						log.Error(err.Error())
						goto EndLine
					}
				case leftResidence:
					if err := l.fsm.Event(ctx, leftResidence, line); err != nil {
						log.Error(err.Error())
						goto EndLine
					}
				case leftToRight:
					if err := l.fsm.Event(ctx, leftToRight, line); err != nil {
						log.Error(err.Error())
						goto EndLine
					}
				case rightMargin:
					if err := l.fsm.Event(ctx, rightMargin, line); err != nil {
						log.Error(err.Error())
						goto EndLine
					}
				case rightResidence:
					if err := l.fsm.Event(ctx, rightResidence, line); err != nil {
						log.Error(err.Error())
						goto EndLine
					}
				case rightToLeft:
					if err := l.fsm.Event(ctx, rightToLeft, line); err != nil {
						log.Error(err.Error())
						goto EndLine
					}
				case levelLeft:
					if err := l.fsm.Event(ctx, levelLeft, line); err != nil {
						log.Error(err.Error())
						goto EndLine
					}
				}
			}
		}
	EndLine:
		log.Infof("end line scan (%d)", id)
		l.lines[id-1].Running = false
		l.fsm.Event(ctx, none)
		close(l.fsmCh)
		l.cancel()
	}(id)

	return nil
}

func (l *Line) Stop(id dsd.LineScanID) error {
	if err := id.Validate(); err != nil {
		return err
	}

	log.Infof("current: %s, id: %d, running:%t", l.fsm.Current(), id, l.lines[id-1].Running)

	if (l.fsm.Current() != none && l.lines[id-1].Running) ||
		(l.fsm.Current() == none && l.lines[id-1].Running) {
		l.stopCh <- id
		close(l.stopCh)
		return nil
	}

	return fmt.Errorf("line scan (%d) is not running", id)
}

func (l *Line) gotoStartPosition(line dsd.LineScan) error {
	pos, err := l.basic.Position()
	if err != nil {
		return err
	}
	pos.Pan = line.LeftMargin

	if err := l.basic.Goto(pos); err != nil {
		return err
	}

	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	if err := l.basic.ReachPosition(timeoutCtx, pos); err != nil {
		return err
	}

	return nil
}

func (l *Line) leftResidence(line dsd.LineScan) error {
	log.Infof("line scan enter_leftResidence (%d %ds)", line.ID, line.ResidenceTimeLeft)

	l.timer.Reset(time.Second * time.Duration(line.ResidenceTimeLeft))

	return nil
}

func (l *Line) leftToRight(line dsd.LineScan) error {
	log.Infof("line scan enter_leftToRight (%d)", line.ID)

	if err := l.basic.Operation(basic.DirectionRight, ptz.Speed(line.Speed)); err != nil {
		return err
	}

	l.arrivePan(line.RightMargin)

	return nil
}

func (l *Line) rightResidence(line dsd.LineScan) error {
	log.Infof("line scan enter_rightResidence (%d)", line.ID)

	l.timer.Reset(time.Second * time.Duration(line.ResidenceTimeRight))

	return nil
}

func (l *Line) rightToLeft(line dsd.LineScan) error {
	log.Infof("line scan enter_rightToLeft (%d)", line.ID)

	if err := l.basic.Operation(basic.DirectionLeft, ptz.Speed(line.Speed)); err != nil {
		return err
	}

	l.arrivePan(line.LeftMargin)

	return nil
}

func (l *Line) arrivePan(pan float64) {
	go func() {
		timeoutCtx, cancelFunc := context.WithTimeout(l.ctx, time.Second*60)
		defer cancelFunc()

		ticker := time.NewTicker(time.Millisecond * 10)
		for {
			select {
			case <-timeoutCtx.Done():
				log.Warn(timeoutCtx.Err().Error())
				return
			case <-ticker.C:
				pos, err := l.basic.Position()
				if err != nil {
					log.Error(err.Error())
					return
				}
				if pos.Pan >= pan-2 && pos.Pan <= pan+2 {
					if err := l.basic.Stop(); err != nil {
						log.Error(err.Error())
						return
					}

					switch l.fsm.Current() {
					case leftToRight:
						l.fsmCh <- rightMargin
					case rightToLeft:
						l.fsmCh <- leftMargin
					}

					return
				}
			}
		}
	}()
}
