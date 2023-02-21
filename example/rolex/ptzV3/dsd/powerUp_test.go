package dsd

import "testing"

func TestNewPowerUps(t *testing.T) {
	ups := NewPowerUps()

	if err := ups.Validate(); err != nil {
		t.Errorf(err.Error())
		return
	}

	t.Logf("%+v", ups)
}
