package dsd

import (
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"

	"strings"
)

const (
	MaxPresetNum = 255
)

type PresetID int

func (i PresetID) Validate() error {
	if err := _validate.Var(i, "gte=1,lte=255"); err != nil {
		log.Error(err.Error())
		return errors.New("preset id is invalid")
	}

	return nil
}

type PresetName string

func (n PresetName) Validate() error {
	name := strings.TrimSpace(string(n))
	if err := _validate.Var(name, "min=1,max=64"); err != nil {
		log.Error(err.Error())
		return errors.New("preset name is invalid")
	}

	return nil
}

type PresetPoint struct {
	Enable   bool       `json:"Enable" validate:"boolean"`             // 使能
	ID       PresetID   `json:"ID" validate:"required,gte=1,lte=255"`  // 预置点id
	Name     PresetName `json:"Name" validate:"required,min=1,max=64"` // 预置点名称
	Position *Position  `json:"Position"`                              // 预置点的坐标和放大倍数
}

func (p *PresetPoint) ConfigKey() string {
	return "PresetPoints"
}

func (p *PresetPoint) Default() []PresetPoint {
	language := "english"
	if err := config.GetConfig("LocalSettings.Language", &language); err != nil {
		log.Error(err.Error())
	}

	s := make([]PresetPoint, MaxPresetNum)
	for id := 1; id <= MaxPresetNum; id++ {
		name := fmt.Sprintf("%s%d", "预置点", id)
		if language != "chinese" {
			name = fmt.Sprintf("%s%d", "Preset", id)
		}

		preset, _ := NewPreset(PresetID(id), PresetName(name))
		s[id-1] = preset
	}

	return s
}

func NewPreset(id PresetID, name PresetName) (PresetPoint, error) {
	log.Debugf("NewPreset param id: %d, name: %s", id, name)

	if err := id.Validate(); err != nil {
		return PresetPoint{}, err
	}

	if err := name.Validate(); err != nil {
		return PresetPoint{}, err
	}

	return PresetPoint{
		Enable:   false,
		ID:       id,
		Name:     name,
		Position: NewPosition(),
	}, nil
}
