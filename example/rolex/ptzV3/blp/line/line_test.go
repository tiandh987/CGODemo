package line

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"testing"
	"time"
)

func TestLine_Start(t *testing.T) {
	ability := ptz.NewMockAbility()
	b := basic.New(ability)

	slice := dsd.NewLineSlice()
	slice[0].Enable = true
	slice[0].LeftMargin = 50
	slice[0].ResidenceTimeLeft = 5
	slice[0].RightMargin = 150
	slice[0].ResidenceTimeRight = 3
	slice[0].Speed = 6

	line := New(b, slice)

	//ctx := context.Background()

	line.Start(1)

	t.Logf("xxxxxxxxxxxxxxxxxxxxxxxxx")

	time.Sleep(time.Second * 60)
	line.Stop(1)
}
