package line

import (
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

func (l *Line) List() dsd.LineSlice {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.lines
}

func (l *Line) Default() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fsm.Current() != none {
		log.Warnf("line scan is running")
		return nil
	}

	before := l.lines

	l.lines = dsd.NewLineSlice()
	if err := config.SetConfig(l.lines.ConfigKey(), l.lines); err != nil {
		l.lines = before
		return err
	}

	return nil
}

func (l *Line) Set(scan *dsd.LineScan) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fsm.Current() != none {
		log.Warnf("line scan is running")
		return nil
	}

	before := l.lines[scan.ID-1]

	scan.Running = false
	l.lines[scan.ID-1] = *scan

	if err := config.SetConfig(l.lines.ConfigKey(), l.lines); err != nil {
		l.lines[scan.ID-1] = before
		return err
	}

	return nil
}

func (l *Line) SetMargin(id dsd.LineScanID, op dsd.LineMarginOp) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fsm.Current() != none {
		log.Warnf("line scan is running")
		return nil
	}

	pos, err := l.basic.Position()
	if err != nil {
		return err
	}

	before := l.lines[id-1]

	switch op {
	case dsd.SetLeftMargin:
		l.lines[id-1].LeftMargin = pos.Pan
	case dsd.SetRightMargin:
		l.lines[id-1].RightMargin = pos.Pan
	case dsd.ClearLeftMargin:
		l.lines[id-1].LeftMargin = dsd.MarginNoLimit
	case dsd.ClearRightMargin:
		l.lines[id-1].RightMargin = dsd.MarginNoLimit
	default:
		log.Warnf("operation(%d) is invalid", op)
		return nil
	}

	if err := config.SetConfig(l.lines.ConfigKey(), l.lines); err != nil {
		l.lines[id-1] = before
		return err
	}

	return nil
}
