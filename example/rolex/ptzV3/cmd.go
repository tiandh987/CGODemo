package ptz

import (
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/arch/serial"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/cron"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/cruise"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/idle"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/line"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/powerUp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/preset"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

var ptzBlpIns blp.PTZRepo

func Start() error {

	// 串口设置
	ptz := dsd.NewPTZ()
	if err := config.SetDefault(ptz.ConfigKey(), ptz); err != nil {
		return err
	}
	if err := config.GetConfig(ptz.ConfigKey(), ptz); err != nil {
		return err
	}

	// 预置点
	presets := dsd.NewPresetSlice()
	if err := config.SetDefault(presets.ConfigKey(), presets); err != nil {
		return err
	}
	if err := config.GetConfig(presets.ConfigKey(), &presets); err != nil {
		return err
	}

	// 线扫
	lines := dsd.NewLineSlice()
	if err := config.SetDefault(lines.ConfigKey(), lines); err != nil {
		return err
	}
	if err := config.GetConfig(lines.ConfigKey(), &lines); err != nil {
		return err
	}

	// 巡航
	cruises := dsd.NewCruiseSlice()
	if err := config.SetDefault(cruises.ConfigKey(), cruises); err != nil {
		return err
	}
	if err := config.GetConfig(cruises.ConfigKey(), &cruises); err != nil {
		return err
	}

	////// 巡迹
	////pattern := dsd.Pattern{}
	////cruises := pattern.Default()
	////if err := config.SetDefault(cruise.ConfigKey(), cruises); err != nil {
	////	return err
	////}
	////if err := config.GetConfig(cruise.ConfigKey(), &cruises); err != nil {
	////	return err
	////}

	// 开机动作
	ups := dsd.NewPowerUps()
	if err := config.SetDefault(ups.ConfigKey(), ups); err != nil {
		return err
	}
	if err := config.GetConfig(ups.ConfigKey(), &ups); err != nil {
		return err
	}

	// 空闲动作
	motion := dsd.NewIdleMotion()
	if err := config.SetDefault(motion.ConfigKey(), motion); err != nil {
		return err
	}
	if err := config.GetConfig(motion.ConfigKey(), &motion); err != nil {
		return err
	}

	// 定时任务
	movementSlice := dsd.NewAutoMovementSlice()
	if err := config.SetDefault(movementSlice.ConfigKey(), movementSlice); err != nil {
		return err
	}
	if err := config.GetConfig(movementSlice.ConfigKey(), &movementSlice); err != nil {
		return err
	}

	s := serial.New("", ptz)
	b := basic.New(s)

	p := preset.New(b, presets)

	l := line.New(b, lines)

	c := cruise.New(b, p, cruises)

	up := powerUp.New(ups)

	i := idle.New(motion)

	c2 := cron.New(movementSlice)

	ptzBlpIns = blp.New(b, p, l, c, up, i, c2)

	ptzBlpIns.Manager().Run()

	return nil
}

func Instance() blp.PTZRepo {
	return ptzBlpIns
}
