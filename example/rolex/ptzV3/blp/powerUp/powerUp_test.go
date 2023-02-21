package powerUp

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"testing"
)

func TestPowerUp(t *testing.T) {
	ups := dsd.NewPowerUps()
	ups.Function = 1
	ups.Function = 2
	ups.Enable = true
	if err := ups.Validate(); err != nil {
		t.Errorf(err.Error())
		return
	}

	up := New(ups)

	t.Logf("%+v", up)
}
