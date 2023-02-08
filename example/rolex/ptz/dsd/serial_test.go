package dsd

import (
	"testing"
)

func TestNewPTZ(t *testing.T) {
	ptz := NewPTZ()
	t.Logf("ptz: %+v", ptz)
	t.Logf("ptz.Attribute: %+v", ptz.Attribute)

	//ptz.Address = 1
	//ptz.Attribute.BaudRate = 1234

	if err := ptz.Validate(); err != nil {
		t.Errorf("ptz validate error: %s", err.Error())
	}
}
