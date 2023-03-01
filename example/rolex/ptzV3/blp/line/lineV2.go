package line

import (
	"context"
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sync"
	"time"
)

type state int

// 线扫状态
const (
	none           state = iota
	leftMargin           // 左边界
	leftResidence        // 左边界停留
	leftToRight          // 左->右
	rightMargin          // 右边界
	rightResidence       // 右边界停留
	rightToLeft          // 右->左
	levelLeft            // 水平旋转-逆时针(左/右边界未设置时)
)

type Line struct {
	mu      sync.RWMutex
	lines   dsd.LineSlice
	running dsd.LineScanID

	basic *basic.Basic

	stateCh chan state
	timer   *time.Timer
	wg      sync.WaitGroup
	cancel  context.CancelFunc
}

func New(b *basic.Basic, s dsd.LineSlice) *Line {
	l := &Line{
		lines:   s,
		running: 0,
		basic:   b,
		stateCh: make(chan state, 1),
	}

	return l
}

func (l *Line) Start(id dsd.LineScanID) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := id.Validate(); err != nil {
		return err
	}

	if l.running != 0 {
		log.Warnf("line scan (%d) is running", l.running)
		return fmt.Errorf("line scan (%d) is running", l.running)
	}

	line := l.lines[id-1]

	if !line.Enable {
		log.Warnf("line scan (%d) is disable", id)
		return fmt.Errorf("line scan (%d) is disable", id)
	}

	l.running = id
	l.lines[id-1].Running = true
	l.stateCh = make(chan state, 1)

	log.Infof("aaaaa")
	l.timer = time.NewTimer(time.Millisecond * 10)
	<-l.timer.C
	log.Infof("bbbbb")

	ctx := context.Background()
	ctx, l.cancel = context.WithCancel(ctx)

	startCh := make(chan struct{})

	l.wg.Add(1)
	go func(id dsd.LineScanID) {
		defer l.wg.Done()

		log.Infof("start line scan id: %d left: %f %ds right: %f %ds", id, line.LeftMargin, line.ResidenceTimeLeft,
			line.RightMargin, line.ResidenceTimeRight)

		var curState state
		for {
			select {
			case <-ctx.Done():
				goto EndLine
			case <-l.timer.C:
				switch curState {
				case leftResidence:
					l.stateCh <- leftToRight
				case rightResidence:
					l.stateCh <- rightToLeft
				}
			case <-startCh:
				log.Infof("startCh receive")
				l.gotoStartPosition(ctx, &line)
				log.Infof("startCh end")
			case curState = <-l.stateCh:
				switch curState {
				case leftMargin:
					l.stateCh <- leftResidence
				case leftResidence:
					l.timer.Reset(time.Second * time.Duration(line.ResidenceTimeLeft))
				case leftToRight:
					if err := l.scan(ctx, basic.DirectionRight, &line); err != nil {
						log.Error(err.Error())
					}
				case rightMargin:
					l.stateCh <- rightResidence
				case rightResidence:
					l.timer.Reset(time.Second * time.Duration(line.ResidenceTimeRight))
				case rightToLeft:
					if err := l.scan(ctx, basic.DirectionLeft, &line); err != nil {
						log.Error(err.Error())
					}
				case levelLeft:
					if err := l.basic.Operation(basic.DirectionLeft, ptz.Speed(line.Speed)); err != nil {
						log.Error(err.Error())
					}
				}
			}
		}
	EndLine:
		log.Infof("end line scan (%d)", id)
		l.basic.Stop()
		l.running = 0
		l.lines[id-1].Running = false
		l.timer.Stop()
		close(l.stateCh)
	}(id)

	if line.LeftMargin == dsd.MarginNoLimit || line.RightMargin == dsd.MarginNoLimit {
		l.stateCh <- levelLeft
	} else {
		startCh <- struct{}{}
	}

	return nil
}

func (l *Line) Stop(id dsd.LineScanID) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := id.Validate(); err != nil {
		return err
	}

	if l.running == 0 {
		return errors.New("line scan is not running")
	}

	l.cancel()
	l.wg.Wait()

	return nil
}

func (l *Line) gotoStartPosition(ctx context.Context, line *dsd.LineScan) {
	pos, err := l.basic.Position()
	if err != nil {
		log.Error(err.Error())
		return
	}
	pos.Pan = line.LeftMargin

	if err := l.basic.Goto(pos); err != nil {
		log.Error(err.Error())
		return
	}

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		defer log.Infof("aaaaaaaaaaaaaaaaaaa")

		timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*30)
		defer cancelFunc()

		ticker := time.NewTicker(time.Millisecond * 100)

		for {
			select {
			case <-timeoutCtx.Done():
				log.Error(timeoutCtx.Err().Error())
				return
			case <-ticker.C:
				curPos, err := l.basic.Position()
				log.Infof("current position: %+v", curPos)
				if err != nil {
					log.Error(err.Error())
					time.Sleep(time.Second)
					continue
				}

				if curPos.Pan >= pos.Pan-2 && curPos.Pan <= pos.Pan+2 &&
					curPos.Tilt >= pos.Tilt-2 && curPos.Tilt <= pos.Tilt+2 &&
					curPos.Zoom >= pos.Zoom-2 && curPos.Zoom <= pos.Zoom+2 {
					log.Infof("xxxxxxxxx stop1")
					l.basic.Stop()
					l.stateCh <- leftMargin
					log.Infof("xxxxxxxxx stop2")
					return
				}
			}
		}
	}()
}

func (l *Line) scan(ctx context.Context, dir basic.Operation, line *dsd.LineScan) error {
	if err := l.basic.Operation(dir, ptz.Speed(line.Speed)); err != nil {
		return err
	}

	pan := line.LeftMargin
	if dir == basic.DirectionRight {
		pan = line.RightMargin
	}

	log.Infof("direction: %d pan: %.2f speed: %d", dir, pan, line.Speed)

	l.wg.Add(1)
	go func() {
		l.wg.Done()
		defer log.Infof("sssssssssssssssssss")

		timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()

		ticker := time.NewTicker(time.Millisecond * 50)
		for {
			select {
			case <-timeoutCtx.Done():
				log.Warn(timeoutCtx.Err().Error())
				return
			case <-ticker.C:
				curPos, err := l.basic.Position()
				log.Debugf("scan current position: %+v", curPos)
				if err != nil {
					log.Error(err.Error())
					time.Sleep(time.Second)
					continue
				}
				if curPos.Pan >= pan-2 && curPos.Pan <= pan+2 {
					l.basic.Stop()

					if dir == basic.DirectionLeft {
						l.stateCh <- leftMargin
					} else if dir == basic.DirectionRight {
						l.stateCh <- rightMargin
					}
					return
				}
			}
		}
	}()

	return nil
}
