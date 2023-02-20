package dsd

import (
	"testing"
)

func TestNewCommAttribute(t *testing.T) {
	attribute := NewCommAttribute()

	if err := attribute.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewCommAttribute_BaudRate(t *testing.T) {
	attribute := NewCommAttribute()

	attribute.BaudRate = 115200

	if err := attribute.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewCommAttribute_DataBit(t *testing.T) {
	attribute := NewCommAttribute()

	attribute.DataBits = 9

	if err := attribute.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewCommAttribute_StopBit(t *testing.T) {
	attribute := NewCommAttribute()

	attribute.StopBits = 3

	if err := attribute.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewCommAttribute_Parity(t *testing.T) {
	attribute := NewCommAttribute()

	attribute.Parity = 5

	if err := attribute.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewPTZ(t *testing.T) {
	ptz := NewPTZ()

	ptz.Address = 2
	ptz.Protocol = 1
	ptz.Attribute.BaudRate = BaudRate57600
	ptz.Attribute.DataBits = 5
	ptz.Attribute.StopBits = 0
	ptz.Attribute.Parity = 1234

	if err := ptz.Validate(); err != nil {
		t.Errorf(err.Error())
	}
}
