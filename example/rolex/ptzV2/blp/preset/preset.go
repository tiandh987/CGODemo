package preset

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

//type presetRepo interface {
//	List(def bool) (dsd.PresetPointSlice, error)
//	Goto(id dsd.PresetID) error
//	GotoOk(id dsd.PresetID) (bool, error)
//	Update(point *dsd.PresetPoint) error
//	Delete(id dsd.PresetID) error
//	DeleteAll() error
//	Set(point *dsd.PresetPoint) error
//	Speed(speed Speed) error
//	PresetID() dsd.PresetID
//}
//
//type presetUseCase struct{}
//
//var _ presetRepo = (*presetUseCase)(nil)
//
//func NewPreset() presetRepo {
//	return &presetUseCase{}
//}
//
//func (p *presetUseCase) List(def bool) (dsd.PresetPointSlice, error) {
//	ps := dsd.NewPresetSlice()
//
//	if !def {
//		if err := config.GetConfig(ps.ConfigKey(), &ps); err != nil {
//			return nil, err
//		}
//	}
//
//	return ps, nil
//}
//
//func (p *presetUseCase) Goto(id dsd.PresetID) error {
//	ps := dsd.NewPresetSlice()
//	if err := config.GetConfig(ps.ConfigKey(), &ps); err != nil {
//		return err
//	}
//
//	if !ps[id-1].Enable {
//		return errors.New("preset is disable")
//	}
//
//	if _, err := serial.Send(protocol.PresetCall, protocol.NoneReplay, 0x00, byte(id)); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *presetUseCase) GotoOk(id dsd.PresetID) (bool, error) {
//	//TODO implement me
//	return false, nil
//}
//
//func (p *presetUseCase) Update(point *dsd.PresetPoint) error {
//	ps := dsd.NewPresetSlice()
//	if err := config.GetConfig(ps.ConfigKey(), &ps); err != nil {
//		return err
//	}
//	ps[point.ID-1].Name = point.Name
//
//	if err := config.SetConfig(ps.ConfigKey(), ps); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *presetUseCase) Delete(id dsd.PresetID) error {
//	ps := dsd.NewPresetSlice()
//	if err := config.GetConfig(ps.ConfigKey(), &ps); err != nil {
//		return err
//	}
//
//	preset, err := dsd.NewPreset(ps[id-1].ID, ps[id-1].Name)
//	if err != nil {
//		return err
//	}
//	ps[id-1] = preset
//
//	if err := config.SetConfig(ps.ConfigKey(), ps); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *presetUseCase) DeleteAll() error {
//	ps := dsd.NewPresetSlice()
//	if err := config.SetConfig(ps.ConfigKey(), &ps); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *presetUseCase) Set(point *dsd.PresetPoint) error {
//	if _, err := serial.Send(protocol.PresetSet, protocol.NoneReplay, 0x00, byte(point.ID)); err != nil {
//		return err
//	}
//
//	position, err := _blpInstance.Ptz.Position()
//	if err != nil {
//		return err
//	}
//	point.Enable = true
//	point.Position = position
//
//	ps := dsd.NewPresetSlice()
//	if err := config.GetConfig(ps.ConfigKey(), &ps); err != nil {
//		return err
//	}
//	ps[point.ID-1] = point
//
//	if err := config.SetConfig(ps.ConfigKey(), ps); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *presetUseCase) Speed(speed Speed) error {
//	if _, err := serial.Send(protocol.PresetSpeed, protocol.NoneReplay, 0x00, speed.Convert()); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *presetUseCase) PresetID() dsd.PresetID {
//	//TODO implement me
//	return 0
//}

type Preset struct {
	presets []dsd.PresetPoint
}

func New(ps []dsd.PresetPoint) *Preset {
	return &Preset{
		presets: ps,
	}
}

func (p *Preset) Start(ctl control.ControlRepo, id int, speed ptz.Speed) error {
	if err := ctl.Goto(p.presets[id].Position); err != nil {
		return err
	}

	return nil
}

func (p *Preset) Stop() error {
	return nil
}
