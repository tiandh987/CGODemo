package lineScan

import (
	"context"
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"sync"
	"time"
)

type LineScan struct {
	mu    sync.RWMutex
	lines []dsd.LineScan

	state *state
	timer *time.Timer
	errCh chan error
	quit  chan struct{}
}

type state struct {
	lst lineState
	ch  chan struct{}
}

func newState() *state {
	return &state{
		lst: none,
		ch:  make(chan struct{}, 1),
	}
}

func (s *state) reset() {
	s.lst = none
	s.ch = make(chan struct{}, 1)
}

func (s *state) update(lst lineState) {
	s.lst = lst
	if lst != leftResidence && lst != rightResidence {
		s.ch <- struct{}{}
	}
}

func New(lines []dsd.LineScan) *LineScan {
	return &LineScan{
		lines: lines,
		state: newState(),
		timer: time.NewTimer(time.Hour),
		errCh: make(chan error, 1),
		quit:  make(chan struct{}, 1),
	}
}

type lineState int

const (
	none lineState = iota // None
	noLimit
	leftResidence  // 左停留
	leftMargin     // 左边界
	leftToRight    // 左->右
	rightResidence // 右停留
	rightMargin    // 右边界
	rightToLeft    // 右->左
)

const MarginNoLimit = -1

func (l *LineScan) Start(ctl control.ControlRepo, id dsd.LineScanID) error {
	if err := id.Validate(); err != nil {
		return err
	}

	if l.state.lst != none {
		log.Warn("linear scan is running")
		return errors.New("linear scan is running")
	}

	line := l.getLine(id)
	if !line.Enable {
		log.Warnf("linear scan %d is disable", id)
		return errors.New(fmt.Sprintf("linear scan %d is disable", id))
	}

	go func() {
		log.Infof("start of linear scan...\nconfig: %+v", line)

		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()

		l.setLineRunning(id, true)
		if err := l.saveConfig(); err != nil {
			log.Error(err.Error())
			goto EndLinearScan
		}

		l.state.reset()
		// 线扫没有左右边界限制
		if line.LeftMargin == MarginNoLimit || line.RightMargin == MarginNoLimit {
			l.state.update(noLimit)
		} else {
			l.state.update(none)
		}

		for {
			select {
			case <-l.quit:
				goto EndLinearScan
			case err := <-l.errCh:
				log.Error(err.Error())
				goto EndLinearScan
			case <-l.timer.C:
				switch l.state.lst {
				case leftResidence:
					l.leftResidence()
				case rightResidence:
					l.rightResidence()
				}
			case <-l.state.ch:
				switch l.state.lst {
				case none:
					l.gotoLeftMargin(ctx, ctl, &line)
				case noLimit:
					l.noLimit(ctl, &line)
				case leftMargin:
					l.leftMargin(&line)
				case leftToRight:
					l.leftToRight(ctx, ctl, &line)
				case rightMargin:
					l.rightMargin(&line)
				case rightToLeft:
					l.rightToLeft(ctx, ctl, &line)
				}
			}
		}
	EndLinearScan:
		log.Infof("end linear scan (%d)", id)

		ctl.Stop()
		l.state.reset()
		l.setLineRunning(id, false)
		l.saveConfig()
	}()

	return nil
}

func (l *LineScan) Stop() {
	if l.state.lst != none {
		l.quit <- struct{}{}
	}
}

func (l *LineScan) noLimit(ctl control.ControlRepo, line *dsd.LineScan) {
	if err := ctl.Left(ptz.Speed(line.Speed).Convert()); err != nil {
		l.errCh <- err
	}
}

func (l *LineScan) gotoLeftMargin(ctx context.Context, ctl control.ControlRepo, line *dsd.LineScan) {
	pos, err := ctl.Position()
	if err != nil {
		l.errCh <- err
		return
	}

	pos.Pan = line.LeftMargin
	if err := ctl.Goto(pos); err != nil {
		l.errCh <- err
		return
	}

	l.arrivePan(ctx, ctl, line.LeftMargin, leftMargin)
}

func (l *LineScan) leftMargin(line *dsd.LineScan) {
	log.Debugf("linear scan (%d) left margin", line.ID)

	l.timer.Reset(time.Second * time.Duration(line.ResidenceTimeLeft))
	l.state.update(leftResidence)
}

func (l *LineScan) leftResidence() {
	log.Debug("end of linear scan left residence")

	l.state.update(leftToRight)
}

func (l *LineScan) leftToRight(ctx context.Context, ctl control.ControlRepo, line *dsd.LineScan) {
	log.Debugf("linear scan (%d) from left to right", line.ID)

	if err := ctl.Right(ptz.Speed(line.Speed).Convert()); err != nil {
		l.errCh <- err
		return
	}

	l.arrivePan(ctx, ctl, line.RightMargin, rightMargin)
}

func (l *LineScan) rightMargin(line *dsd.LineScan) {
	log.Debugf("linear scan (%d) right margin", line.ID)

	l.timer.Reset(time.Second * time.Duration(line.ResidenceTimeRight))
	l.state.update(rightResidence)
}

func (l *LineScan) rightResidence() {
	log.Debug("end of linear scan right residence")

	l.state.update(rightToLeft)
}

func (l *LineScan) rightToLeft(ctx context.Context, ctl control.ControlRepo, line *dsd.LineScan) {
	log.Debugf("linear scan (%d) from right to left", line.ID)

	if err := ctl.Left(ptz.Speed(line.Speed).Convert()); err != nil {
		l.errCh <- err
		return
	}

	l.arrivePan(ctx, ctl, line.LeftMargin, leftMargin)
}

func (l *LineScan) arrivePan(ctx context.Context, ctl control.ControlRepo, pan float64, lst lineState) {
	go func() {
		timer := time.NewTimer(time.Second * 300)

		for {
			select {
			case <-ctx.Done():
				log.Info(ctx.Err().Error())
				return
			case <-timer.C:
				l.errCh <- errors.New("timeout waiting for jump to the expect pan")
				return
			default:
				pos, err := ctl.Position()
				if err != nil {
					l.errCh <- err
					return
				}
				if pos.Pan >= pan-2 && pos.Pan <= pan+2 {
					if err := ctl.Stop(); err != nil {
						l.errCh <- err
						return
					}

					l.state.update(lst)
					return
				}
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()
}

func (l *LineScan) saveConfig() error {
	if err := config.SetConfig(l.lines[0].ConfigKey(), l.lines); err != nil {
		return err
	}

	return nil
}

func (l *LineScan) getLine(id dsd.LineScanID) dsd.LineScan {
	l.mu.RLock()
	line := l.lines[id-1]
	l.mu.RUnlock()

	return line
}

func (l *LineScan) setLineRunning(id dsd.LineScanID, running bool) {
	l.mu.Lock()
	l.lines[id-1].Runing = running
	l.mu.Unlock()
}
