package dsd

import (
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"strings"
)

const (
	MaxCruiseNum = 8
)

type CruiseID int

func (i CruiseID) Validate() error {
	if err := _validate.Var(i, "gte=1,lte=8"); err != nil {
		log.Error(err.Error())
		return errors.New("cruise id is invalid")
	}

	return nil
}

type CruiseName string

func (n CruiseName) Validate() error {
	name := strings.TrimSpace(string(n))
	if err := _validate.Var(name, "min=1,max=64"); err != nil {
		log.Error(err.Error())
		return errors.New("preset name is invalid")
	}

	return nil
}

type TourPreset struct {
	Enable  bool              `json:"Enable" validate:"boolean"`             // 使能
	ID      CruiseID          `json:"ID" validate:"required,gte=1,lte=8"`    // ID
	Name    CruiseName        `json:"Name" validate:"required,min=1,max=64"` // 名字
	Preset  []TourPresetPoint `json:"Preset" validate:"required"`            // 关联预置点
	Running bool              `json:"Runing" validate:"boolean"`             // 运行状态
}

func (p *TourPreset) ConfigKey() string {
	return "CruisePoints"
}

func (p *TourPreset) Default() []TourPreset {
	language := "english"
	if err := config.GetConfig("LocalSettings.Language", &language); err != nil {
		log.Error(err.Error())
	}

	s := make([]TourPreset, MaxCruiseNum)
	for id := 1; id <= MaxCruiseNum; id++ {
		name := fmt.Sprintf("%s%d", "巡航组", id)
		if language != "chinese" {
			name = fmt.Sprintf("%s%d", "Cruise", id)
		}

		preset, _ := NewCruise(CruiseID(id), CruiseName(name))
		s[id-1] = preset
	}

	return s
}

func NewCruise(id CruiseID, name CruiseName) (TourPreset, error) {
	log.Debugf("id: %d, name: %s", id, name)

	if err := id.Validate(); err != nil {
		return TourPreset{}, err
	}

	if err := name.Validate(); err != nil {
		return TourPreset{}, err
	}

	return TourPreset{
		Enable:  false,
		ID:      id,
		Name:    name,
		Preset:  []TourPresetPoint{},
		Running: false,
	}, nil
}

type TourPresetPoint struct {
	ID            PresetID   `json:"ID" validate:"required,gte=1,lte=255"`            // 预设点id
	Name          PresetName `json:"Name" validate:"required,min=1,max=64"`           // 预设点名称
	ResidenceTime int        `json:"ResidenceTime" validate:"required,gte=0,lte=255"` // 停留时间（s）
}
