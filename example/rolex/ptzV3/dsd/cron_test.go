package dsd

import "testing"

func TestNewAutoMovement(t *testing.T) {
	movement := NewAutoMovement(1)

	//movement.ID = 10
	movement.ID = 4
	movement.Enable = true
	if err := movement.Validate(); err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("%+v", movement)
}

func TestNewAutoMovementSlice(t *testing.T) {
	slice := NewAutoMovementSlice()

	t.Logf("config key: %s", slice.ConfigKey())

	for i, movement := range slice {
		t.Logf("index: %d movement: %+v", i, movement)
	}
}
