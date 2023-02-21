package dsd

import (
	"errors"
	"fmt"
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
		return errors.New("curse id is invalid")
	}

	return nil
}

type TourPreset struct {
	Enable  bool              `json:"Enable" validate:"boolean"`
	ID      CruiseID          `json:"ID" validate:"required,gte=1,lte=8"`
	Name    string            `json:"Name" validate:"required,min=1,max=64"`
	Preset  []TourPresetPoint `json:"Preset" `
	Running bool              `json:"Runing" validate:"boolean"`
}

func NewCruise(id CruiseID, name string) TourPreset {
	trimSpaceName := strings.TrimSpace(name)

	return TourPreset{
		Enable:  false,
		ID:      id,
		Name:    trimSpaceName,
		Preset:  []TourPresetPoint{},
		Running: false,
	}
}

func (p *TourPreset) Validate() error {
	if err := _validate.Struct(p); err != nil {
		return err
	}

	for _, point := range p.Preset {
		if err := _validate.Struct(point); err != nil {
			return err
		}
	}

	return nil
}

type TourPresetPoint struct {
	ID            PresetID `json:"ID" validate:"required,gte=1,lte=255"`
	Name          string   `json:"Name" validate:"required,min=1,max=64"`
	ResidenceTime int      `json:"ResidenceTime" validate:"gte=0,lte=255"`
}

func NewTourPresetPoint(id PresetID, name string, time int) TourPresetPoint {
	trimSpaceName := strings.TrimSpace(name)

	return TourPresetPoint{
		ID:            id,
		Name:          trimSpaceName,
		ResidenceTime: time,
	}
}

func (p *TourPresetPoint) Validate() error {
	if err := _validate.Struct(p); err != nil {
		return err
	}

	return nil
}

type CruiseSlice []TourPreset

func (s *CruiseSlice) ConfigKey() string {
	return "CruisePoints"
}

func NewCruiseSlice() CruiseSlice {
	language := "english"
	//if err := config.GetConfig("LocalSettings.Language", &language); err != nil {
	//	log.Error(err.Error())
	//}

	prefix := "Cruise"
	if language == "chinese" {
		prefix = "巡航组"
	}

	s := make([]TourPreset, MaxCruiseNum)
	for id := 1; id <= MaxCruiseNum; id++ {
		name := fmt.Sprintf("%s%d", prefix, id)
		s[id-1] = NewCruise(CruiseID(id), name)
	}

	return s
}
