package dsdOld

import "testing"

func TestNewPreset(t *testing.T) {
	preset1, err := NewPreset(1, "aaaaaaaaaa")
	if err != nil {
		t.Error(err)
	}
	t.Logf("preset1: %+v, position: %+v", preset1, preset1.Position)

	preset100, err := NewPreset(100, "0123456789012345678901234567890123456789012345678901234567891234")
	if err != nil {
		t.Error(err)
	}
	t.Logf("preset100: %+v, position: %+v", preset100, preset100.Position)

	preset255, err := NewPreset(255, "this is 255")
	if err != nil {
		t.Error(err)
	}
	t.Logf("preset255: %+v, position: %+v", preset255, preset100.Position)

	//preset256, err := NewPreset(256, "this is 256")
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Logf("preset256: %+v, position: %+v", preset256, preset256.Position)
}
