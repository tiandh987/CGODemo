package lineScan

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"time"
)

type LineScan struct {
	lines []dsd.LineScan
	state lineState
	quit  chan struct{}
}

func New(lines []dsd.LineScan) *LineScan {
	return &LineScan{
		lines: lines,
		state: none,
		quit:  make(chan struct{}, 1),
	}
}

type lineState int

const (
	none           lineState = iota // None
	leftResidence                   // 左停留
	leftMargin                      // 左边界
	leftToRight                     // 左->右
	rightResidence                  // 右停留
	rightMargin                     // 右边界
	rightToLeft                     // 右->左
)

func (l *LineScan) Start(ctl control.ControlRepo, id int, speed ptz.Speed) error {
	line := l.lines[id]

	pos, err := ctl.Position()
	if err != nil {
		return err
	}

	pos.Pan = line.LeftMargin
	if err := ctl.Goto(pos); err != nil {
		return err
	}

	l.state = leftResidence
	go func() {
		for {
			select {
			case <-l.quit:
				ctl.Stop()
				return
			default:
				switch l.state {
				case leftResidence:
					time.Sleep(time.Second * time.Duration(line.ResidenceTimeLeft))
					l.state = leftMargin
				case leftMargin:
					ctl.Right(speed.Convert())
					l.state = leftToRight
				case leftToRight:
					for {
						pos, _ := ctl.Position()
						if pos.Pan >= line.RightMargin-0.2 || pos.Pan <= line.RightMargin+0.2 {
							break
						}
					}
					ctl.Stop()
					l.state = rightResidence
				case rightResidence:
					time.Sleep(time.Second * time.Duration(line.ResidenceTimeRight))
					l.state = rightMargin
				case rightMargin:
					ctl.Left(speed.Convert())
					l.state = rightToLeft
				case rightToLeft:
					for {
						pos, _ := ctl.Position()
						if pos.Pan >= line.LeftMargin-0.2 || pos.Pan <= line.LeftMargin+0.2 {
							break
						}
					}
					ctl.Stop()
					l.state = leftResidence
				}
			}
		}
	}()

	return nil
}

func (l *LineScan) Stop() {
	l.quit <- struct{}{}
}
