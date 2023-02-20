package dsd

import (
	"fmt"
	"strings"
)

const (
	MaxPresetNum = 255
)

type PresetID int

func (i PresetID) Validate() error {
	if err := _validate.Var(i, "gte=1,lte=255"); err != nil {
		return err
	}

	return nil
}

type PresetPoint struct {
	Enable   bool     `json:"Enable" validate:"boolean"`
	ID       PresetID `json:"ID" validate:"required,gte=1,lte=255"`
	Name     string   `json:"Name" validate:"required,min=1,max=64"`
	Position Position `json:"Position"`
}

func NewPreset(id PresetID, name string) PresetPoint {
	trimSpaceName := strings.TrimSpace(name)

	return PresetPoint{
		Enable:   false,
		ID:       id,
		Name:     trimSpaceName,
		Position: NewPosition(),
	}
}

func (p *PresetPoint) Validate() error {
	if err := _validate.Struct(p); err != nil {
		return err
	}
	return nil
}

type PresetSlice []PresetPoint

func (p *PresetSlice) ConfigKey() string {
	return "PresetPoints"
}

func NewPresetSlice() PresetSlice {
	language := "english"
	//if err := config.GetConfig("LocalSettings.Language", &language); err != nil {
	//	log.Error(err.Error())
	//}

	prefix := "Preset"
	if language == "chinese" {
		prefix = "预置点"
	}

	s := make([]PresetPoint, MaxPresetNum)
	for id := 1; id <= MaxPresetNum; id++ {
		name := fmt.Sprintf("%s%d", prefix, id)
		s[id-1] = NewPreset(PresetID(id), name)
	}

	return s
}
