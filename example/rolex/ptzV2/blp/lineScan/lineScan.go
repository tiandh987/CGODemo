package lineScan

import (
	"context"
	"errors"
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

	state   lineState
	timer   *time.Timer
	stateCh chan struct{}
	errCh   chan error
	quit    chan struct{}
}

func New(lines []dsd.LineScan) *LineScan {
	return &LineScan{
		lines:   lines,
		state:   none,
		timer:   time.NewTimer(time.Hour),
		stateCh: make(chan struct{}, 1),
		errCh:   make(chan error, 1),
		quit:    make(chan struct{}, 1),
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

	if l.state != none {
		log.Warn("linear scan is running")
		return errors.New("linear scan is running")
	}

	line := l.getLine(id)
	if !line.Enable {
		log.Warnf("linear scan %d is disable", id)
		return nil
	}

	go func() {
		log.Infof("start of linear scan...\nconfig: %+v", line)

		ctx, cancelFunc := context.WithCancel(context.Background())

		l.setLineRunning(id, true)
		if err := l.saveConfig(); err != nil {
			log.Error(err.Error())
			goto EndLinearScan
		}

		// 线扫没有左右边界限制
		if line.LeftMargin == MarginNoLimit || line.RightMargin == MarginNoLimit {
			l.state = noLimit
			l.stateCh <- struct{}{}
		} else {
			l.state = none
			l.stateCh <- struct{}{}
		}

		for {
			select {
			case <-l.quit:
				goto EndLinearScan
			case err := <-l.errCh:
				log.Error(err.Error())
				goto EndLinearScan
			case <-l.timer.C:
				switch l.state {
				case leftResidence:
					l.leftResidence()
				case rightResidence:
					l.rightResidence()
				}
			case <-l.stateCh:
				switch l.state {
				case none:
					l.gotoLeftMargin(ctx, ctl, id)
				case noLimit:
					l.noLimit(ctl, id)
				case leftMargin:
					l.leftMargin(id)
				case leftToRight:
					l.leftToRight(ctx, ctl, id)
				case rightMargin:
					l.rightMargin(id)
				case rightToLeft:
					l.rightToLeft(ctx, ctl, id)
				}
			}
		}
	EndLinearScan:
		log.Infof("end linear scan (%d)", id)

		cancelFunc()
		ctl.Stop()
		l.state = none
		l.setLineRunning(id, false)
		l.saveConfig()
	}()

	return nil
}

func (l *LineScan) Stop() {
	if l.state != none {
		l.quit <- struct{}{}
	}
}

func (l *LineScan) noLimit(ctl control.ControlRepo, id dsd.LineScanID) {
	line := l.getLine(id)

	if err := ctl.Left(ptz.Speed(line.Speed).Convert()); err != nil {
		l.errCh <- err
	}
}

func (l *LineScan) gotoLeftMargin(ctx context.Context, ctl control.ControlRepo, id dsd.LineScanID) {
	line := l.getLine(id)

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

func (l *LineScan) leftMargin(id dsd.LineScanID) {
	log.Debugf("linear scan (%d) left margin", id)

	line := l.getLine(id)
	l.timer.Reset(time.Second * time.Duration(line.ResidenceTimeLeft))
	l.state = leftResidence
}

func (l *LineScan) leftResidence() {
	log.Debug("end of linear scan left residence")

	l.state = leftToRight
	l.stateCh <- struct{}{}
}

func (l *LineScan) leftToRight(ctx context.Context, ctl control.ControlRepo, id dsd.LineScanID) {
	log.Debugf("linear scan (%d) from left to right", id)

	line := l.getLine(id)
	if err := ctl.Right(ptz.Speed(line.Speed).Convert()); err != nil {
		l.errCh <- err
		return
	}

	l.arrivePan(ctx, ctl, line.RightMargin, rightMargin)
}

func (l *LineScan) rightMargin(id dsd.LineScanID) {
	log.Debugf("linear scan (%d) right margin", id)

	line := l.getLine(id)
	l.timer.Reset(time.Second * time.Duration(line.ResidenceTimeRight))
	l.state = rightResidence
}

func (l *LineScan) rightResidence() {
	log.Debug("end of linear scan right residence")

	l.state = rightToLeft
	l.stateCh <- struct{}{}
}

func (l *LineScan) rightToLeft(ctx context.Context, ctl control.ControlRepo, id dsd.LineScanID) {
	log.Debugf("linear scan (%d) from right to left", id)

	line := l.getLine(id)
	if err := ctl.Left(ptz.Speed(line.Speed).Convert()); err != nil {
		l.errCh <- err
		return
	}

	l.arrivePan(ctx, ctl, line.LeftMargin, leftMargin)
}

func (l *LineScan) arrivePan(ctx context.Context, ctl control.ControlRepo, pan float64, state lineState) {

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

					l.state = state
					l.stateCh <- struct{}{}

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
