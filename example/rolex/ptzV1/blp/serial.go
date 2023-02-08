package blp

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/dsd"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial"
)

type SerialRepo interface {
	Get(def bool) (*dsd.PTZ, error)
	Set(ptz *dsd.PTZ) error
}

type serialUseCase struct{}

var _ SerialRepo = (*serialUseCase)(nil)

func NewSerial() SerialRepo {
	return &serialUseCase{}
}

func (s *serialUseCase) Get(def bool) (*dsd.PTZ, error) {
	log.Debugf("param def: %t", def)

	cfg := dsd.NewPTZ()

	if !def {
		if err := config.GetConfig(cfg.ConfigKey(), cfg); err != nil {
			log.Errorf("get %s config failed, err: %s", cfg.ConfigKey(), err.Error())
			return nil, errors.New("get config failed")
		}
	}

	log.Debugf("get %s config, ptz: %+v, attribute: %+v", cfg.ConfigKey(), cfg, cfg.Attribute)

	return cfg, nil
}

func (s *serialUseCase) Set(ptz *dsd.PTZ) error {
	log.Debugf("set serial, ptz: %+v, attribute: %+v", ptz, ptz.Attribute)

	if err := serial.Set(ptz); err != nil {
		return err
	}

	if err := config.SetConfig(ptz.ConfigKey(), ptz); err != nil {
		return err
	}

	return nil
}
