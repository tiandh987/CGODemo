package dsd

import "testing"

func TestNewLine(t *testing.T) {
	line := NewLine(1)

	//line.ID = 0
	//line.ID = 6
	line.ID = 5

	line.LeftMargin = -2
	line.LeftMargin = -1
	line.LeftMargin = 0
	line.LeftMargin = 359
	line.LeftMargin = 359.99999
	//line.LeftMargin = 360

	line.ResidenceTimeLeft = 0
	line.ResidenceTimeLeft = 1
	//line.ResidenceTimeLeft = 61

	line.Speed = 1
	line.Speed = 0
	line.Speed = 9
	line.Speed = 8

	if err := line.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewLineSlice(t *testing.T) {
	slice := NewLineSlice()

	t.Log(slice.ConfigKey())

	for i, scan := range slice {
		t.Logf("index: %d, scan: %+v\n", i, scan)
	}
}
