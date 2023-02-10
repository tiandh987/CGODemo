package lineScan

import (
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (l *LineScan) List() []dsd.LineScan {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.lines
}

func (l *LineScan) Default() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.state != none {
		return errors.New("linear scan is running")
	}

	lines := make([]dsd.LineScan, dsd.MaxLineScanNum)
	for id := 1; id <= dsd.MaxLineScanNum; id++ {
		lineScan, err := dsd.NewLineScan(dsd.LineScanID(id))
		if err != nil {
			return err
		}
		lines[id-1] = lineScan
	}

	if err := l.saveConfig(); err != nil {
		return err
	}

	l.lines = lines

	return nil
}

func (l *LineScan) Set(scan *dsd.LineScan) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.lines[scan.ID-1].Runing {
		return errors.New(fmt.Sprintf("linear scan %d is running", scan.ID))
	}

	before := l.lines[scan.ID]
	scan.Runing = false
	l.lines[scan.ID-1] = *scan

	if err := l.saveConfig(); err != nil {
		l.lines[scan.ID-1] = before
		return err
	}

	return nil
}

type Operation int

const (
	SetLeftMargin Operation = iota
	SetRightMargin
	ClearLeftMargin
	ClearRightMargin
)

func (l *LineScan) SetMargin(ctl control.ControlRepo, id dsd.LineScanID, op Operation) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.state != none && l.lines[id-1].Runing {
		return errors.New(fmt.Sprintf("linear scan %d is running", id))
	}

	pos, err := ctl.Position()
	if err != nil {
		return err
	}

	before := l.lines[id-1]

	switch op {
	case SetLeftMargin:
		l.lines[id-1].LeftMargin = pos.Pan
	case SetRightMargin:
		l.lines[id-1].RightMargin = pos.Pan
	case ClearLeftMargin:
		l.lines[id-1].LeftMargin = MarginNoLimit
	case ClearRightMargin:
		l.lines[id-1].RightMargin = MarginNoLimit
	default:
		log.Warnf("operation(%d) is invalid", op)
		return nil
	}

	if err := l.saveConfig(); err != nil {
		l.lines[id-1] = before
		return err
	}
	return nil
}
