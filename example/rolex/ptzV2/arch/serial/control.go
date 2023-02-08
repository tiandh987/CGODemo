package serial

import (
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/arch/protocol"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/control"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
)

var _ control.ControlRepo = (*Serial)(nil)

func (s *Serial) Version() (string, error) {
	ver, err := s.Send(protocol.Version, protocol.VersionReplay, 0x00, 0x00)
	if err != nil {
		return "", err
	}

	version := fmt.Sprintf("V%d.%d", ver[0], ver[1])

	log.Debugf("ver: %x, version: %s", ver, version)
	return version, nil
}

func (s *Serial) Model() (string, error) {
	mod, err := s.Send(protocol.Model, protocol.ModelReplay, 0x00, 0x00)
	if err != nil {
		return "", err
	}

	model := fmt.Sprintf("V%d.%d", mod[0], mod[1])

	log.Debugf("mod: %x, model: %s", mod, model)
	return model, nil
}

func (s *Serial) Restart() error {
	if _, err := s.Send(protocol.Restart, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Stop() error {
	if _, err := s.Send(protocol.Stop, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Up(speed byte) error {
	if _, err := s.Send(protocol.Up, protocol.NoneReplay, 0x00, speed); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Down(speed byte) error {
	if _, err := s.Send(protocol.Down, protocol.NoneReplay, 0x00, speed); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Left(speed byte) error {
	if _, err := s.Send(protocol.Left, protocol.NoneReplay, speed, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Right(speed byte) error {
	if _, err := s.Send(protocol.Right, protocol.NoneReplay, speed, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *Serial) LeftUp(speed byte) error {
	if _, err := s.Send(protocol.LeftUp, protocol.NoneReplay, speed, speed); err != nil {
		return err
	}

	return nil
}

func (s *Serial) RightUp(speed byte) error {
	if _, err := s.Send(protocol.RightUp, protocol.NoneReplay, speed, speed); err != nil {
		return err
	}

	return nil
}

func (s *Serial) LeftDown(speed byte) error {
	if _, err := s.Send(protocol.LeftDown, protocol.NoneReplay, speed, speed); err != nil {
		return err
	}

	return nil
}

func (s *Serial) RightDown(speed byte) error {
	if _, err := s.Send(protocol.RightDown, protocol.NoneReplay, speed, speed); err != nil {
		return err
	}

	return nil
}

func (s *Serial) ZoomAdd() error {
	if _, err := s.Send(protocol.ZoomAdd, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *Serial) ZoomSub() error {
	if _, err := s.Send(protocol.ZoomSub, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Position() (*dsd.Position, error) {
	//pan, err := s.Send(protocol.PanGet, protocol.PanReplay, 0x00, 0x00)
	//if err != nil {
	//	return nil, err
	//}
	//
	//tilt, err := s.Send(protocol.TiltGet, protocol.TiltReplay, 0x00, 0x00)
	//if err != nil {
	//	return nil, err
	//}
	//
	//zoom, err := s.Send(protocol.ZoomGet, protocol.ZoomReplay, 0x00, 0x00)
	//if err != nil {
	//	return nil, err
	//}
	//
	//pos := control.Position{
	//	Pan:  0,
	//	Tilt: 0,
	//	Zoom: 0,
	//}
	//external := pos.External(pan, tilt, zoom)
	//
	//log.Debugf("pan: %x - %d, tilt: %x - %d, zoom: %x - %d, external: %+v",
	//	pan, pan, tilt, tilt, zoom, zoom, external)
	//
	return nil, nil
}

func (s *Serial) Goto(position *dsd.Position) error {
	//pos := &ptzPosition{}
	//pos.Internal(position)
	//
	//if _, err := p.serial.Send(protocol.PanSet, protocol.NoneReplay, pos.pan[0], pos.pan[1]); err != nil {
	//	return err
	//}
	//
	//if _, err := p.serial.Send(protocol.TiltSet, protocol.NoneReplay, pos.tilt[0], pos.tilt[1]); err != nil {
	//	return err
	//}
	//
	//if _, err := p.serial.Send(protocol.ZoomSet, protocol.NoneReplay, pos.zoom[0], pos.zoom[1]); err != nil {
	//	return err
	//}
	//
	return nil
}
