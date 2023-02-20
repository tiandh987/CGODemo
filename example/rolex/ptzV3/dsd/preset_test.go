package dsd

import "testing"

func TestNewPreset(t *testing.T) {

	//preset := NewPreset(999, "aaa")
	//preset := NewPreset(1, "01234567890123456789012345678901234567890123456789012345678912345")
	//preset := NewPreset(1, "                                                                 ")

	preset := NewPreset(1, "aaa")
	if err := preset.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewPresetSlice(t *testing.T) {
	slice := NewPresetSlice()

	for i, point := range slice {
		t.Logf("index: %d - point: %+v\n", i, point)
	}

	t.Log(slice.ConfigKey())
}
