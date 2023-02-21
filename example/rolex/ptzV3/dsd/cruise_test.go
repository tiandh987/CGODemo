package dsd

import "testing"

func TestNewCruise(t *testing.T) {
	cruise := NewCruise(1, "aaa")

	//cruise.ID = 10
	//cruise.ID = 0
	//cruise.ID = -1
	cruise.ID = 5

	point := NewTourPresetPoint(255, "aaa", 10)

	if err := point.Validate(); err != nil {
		t.Error(err.Error())
		return
	}
	cruise.Preset = append(cruise.Preset, point)

	if err := cruise.Validate(); err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("%+v", cruise)
}

func TestNewCruiseSlice(t *testing.T) {
	slice := NewCruiseSlice()

	for i, preset := range slice {
		t.Logf("index: %d, preset: %+v", i, preset)
	}

	t.Log(slice.ConfigKey())
}
