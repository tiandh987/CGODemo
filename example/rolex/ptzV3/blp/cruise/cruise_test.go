package cruise

import (
	"context"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"testing"
	"time"
)

func TestCruise_Start(t *testing.T) {
	ability := ptz.NewMockAbility()
	b := basic.New(ability)

	presetSlice := dsd.NewPresetSlice()
	presetSlice[0].Enable = true
	presetSlice[0].Position = dsd.Position{
		Pan:  50,
		Tilt: 0,
		Zoom: 1,
	}

	presetSlice[1].Enable = true
	presetSlice[1].Position = dsd.Position{
		Pan:  150,
		Tilt: 0,
		Zoom: 1,
	}

	p := preset.New(b, presetSlice)

	cruiseSlice := dsd.NewCruiseSlice()
	cruiseSlice[0].Enable = true
	cruiseSlice[0].Preset = append(cruiseSlice[0].Preset, dsd.TourPresetPoint{
		ID:            1,
		Name:          "xx",
		ResidenceTime: 5,
	})
	cruiseSlice[0].Preset = append(cruiseSlice[0].Preset, dsd.TourPresetPoint{
		ID:            2,
		Name:          "xxxx",
		ResidenceTime: 3,
	})
	cruiseSlice[0].Preset = append(cruiseSlice[0].Preset, dsd.TourPresetPoint{
		ID:            3,
		Name:          "xxxxccc",
		ResidenceTime: 3,
	})

	cruise := New(p, cruiseSlice)

	//ctx, cancelFunc := context.WithCancel(context.Background())
	ctx := context.Background()
	cruise.Start(ctx, 1)

	t.Logf("xxxxxxxxxxxxxxxxxxxxxxxxx")

	time.Sleep(time.Second * 30)
	cruise.Stop(ctx, 1)

	time.Sleep(time.Second * 30)

	t.Logf("cancelFunc exec")
	//cancelFunc()

	time.Sleep(time.Second * 3600)
}
