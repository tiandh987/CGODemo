package dsd

import "testing"

func TestNewPosition(t *testing.T) {
	position := NewPosition()
	if err := position.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewPositionPan(t *testing.T) {
	position := NewPosition()

	position.Pan = 15.99
	if err := position.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewPositionPanLess(t *testing.T) {
	position := NewPosition()

	position.Pan = -4
	if err := position.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewPositionPanMore(t *testing.T) {
	position := NewPosition()

	//position.Pan = 359.999
	position.Pan = 360
	if err := position.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewPositionTilt(t *testing.T) {
	position := NewPosition()

	//position.Tilt = -1
	//position.Tilt = 0
	//position.Tilt = 45
	//position.Tilt = 90
	position.Tilt = 91

	if err := position.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewPositionZoom(t *testing.T) {
	position := NewPosition()

	//position.Zoom = -1
	//position.Zoom = 0
	//position.Zoom = 1
	//position.Zoom = 15
	//position.Zoom = 20
	position.Zoom = 21

	if err := position.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}
