package blp

import (
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/lineScan"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func (b *Blp) ListLine() []dsd.LineScan {
	return b.line.List()
}

func (b *Blp) DefaultLine() error {
	return b.line.Default()
}

func (b *Blp) SetLine(scan *dsd.LineScan) error {
	return b.line.Set(scan)
}

// SetLineMargin
// limit  左右边界  1-左边界   2-右边界
// clear  清除边界  true-清除  false-设置
func (b *Blp) SetLineMargin(id dsd.LineScanID, limit int, clear bool) error {
	var op lineScan.Operation

	if clear && limit == 1 {
		op = lineScan.ClearLeftMargin
	} else if clear && limit == 2 {
		op = lineScan.ClearRightMargin
	} else if !clear && limit == 1 {
		op = lineScan.SetLeftMargin
	} else if !clear && limit == 2 {
		op = lineScan.SetRightMargin
	} else {
		return errors.New(fmt.Sprintf("param is invalid. limit: %d, clear: %t", limit, clear))
	}

	return b.line.SetMargin(b.getControl(), id, op)
}
