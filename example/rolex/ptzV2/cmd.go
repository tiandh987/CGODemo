package ptz

import (
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

func Start() error {
	// 限位
	limit := dsd.NewLimit()
	if err := config.SetDefault(limit.ConfigKey(), limit); err != nil {
		return err
	}
	if err := config.GetConfig(limit.ConfigKey(), limit); err != nil {
		return err
	}

	// 串口设置
	ptz := dsd.NewPTZ()
	if err := config.SetDefault(ptz.ConfigKey(), ptz); err != nil {
		return err
	}
	if err := config.GetConfig(ptz.ConfigKey(), ptz); err != nil {
		return err
	}

	// 预置点
	preset := dsd.PresetPoint{}
	presets := preset.Default()
	if err := config.SetDefault(preset.ConfigKey(), presets); err != nil {
		return err
	}
	if err := config.GetConfig(preset.ConfigKey(), &presets); err != nil {
		return err
	}

	// 线扫
	line := dsd.LineScan{}
	lines := line.Default()
	if err := config.SetDefault(line.ConfigKey(), lines); err != nil {
		return err
	}
	if err := config.GetConfig(line.ConfigKey(), &lines); err != nil {
		return err
	}

	// 巡迹
	cruise := dsd.TourPreset{}
	cruises := cruise.Default()
	if err := config.SetDefault(cruise.ConfigKey(), cruises); err != nil {
		return err
	}
	if err := config.GetConfig(cruise.ConfigKey(), &cruises); err != nil {
		return err
	}

	blpInstance := blp.New(limit, "", ptz, presets, lines, cruises)
	blp.Replace(blpInstance)

	return nil
}
