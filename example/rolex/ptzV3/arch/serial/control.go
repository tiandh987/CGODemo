package serial

import (
	"encoding/binary"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/arch/protocol"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

var _ ptz.AbilityRepo = (*Serial)(nil)

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

func (s *Serial) Up(speed ptz.Speed) error {
	if _, err := s.Send(protocol.Up, protocol.NoneReplay, 0x00, speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Down(speed ptz.Speed) error {
	if _, err := s.Send(protocol.Down, protocol.NoneReplay, 0x00, speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Left(speed ptz.Speed) error {
	if _, err := s.Send(protocol.Left, protocol.NoneReplay, speed.Convert(), 0x00); err != nil {
		return err
	}

	return nil
}

func (s *Serial) Right(speed ptz.Speed) error {
	if _, err := s.Send(protocol.Right, protocol.NoneReplay, speed.Convert(), 0x00); err != nil {
		return err
	}

	return nil
}

func (s *Serial) LeftUp(speed ptz.Speed) error {
	if _, err := s.Send(protocol.LeftUp, protocol.NoneReplay, speed.Convert(), speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *Serial) RightUp(speed ptz.Speed) error {
	if _, err := s.Send(protocol.RightUp, protocol.NoneReplay, speed.Convert(), speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *Serial) LeftDown(speed ptz.Speed) error {
	if _, err := s.Send(protocol.LeftDown, protocol.NoneReplay, speed.Convert(), speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *Serial) RightDown(speed ptz.Speed) error {
	if _, err := s.Send(protocol.RightDown, protocol.NoneReplay, speed.Convert(), speed.Convert()); err != nil {
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
	pan, err := s.Send(protocol.PanGet, protocol.PanReplay, 0x00, 0x00)
	if err != nil {
		return nil, err
	}

	tilt, err := s.Send(protocol.TiltGet, protocol.TiltReplay, 0x00, 0x00)
	if err != nil {
		return nil, err
	}

	zoom, err := s.Send(protocol.ZoomGet, protocol.ZoomReplay, 0x00, 0x00)
	if err != nil {
		return nil, err
	}

	pos, err := s.externalPosition(pan, tilt, zoom)
	if err != nil {
		return nil, err
	}

	return pos, nil
}

func (s *Serial) Goto(pos *dsd.Position) error {
	if err := s.GotoPan(pos); err != nil {
		return err
	}

	if err := s.GotoTilt(pos); err != nil {
		return err
	}

	if err := s.GotoZoom(pos); err != nil {
		return err
	}

	return nil
}

func (s *Serial) GotoPan(pos *dsd.Position) error {
	pan, _, _ := s.internalPosition(pos)

	log.Debugf("pan: %x - %.2f", pan, pos.Pan)

	if _, err := s.Send(protocol.PanSet, protocol.NoneReplay, pan[0], pan[1]); err != nil {
		return err
	}

	return nil
}

func (s *Serial) GotoTilt(pos *dsd.Position) error {
	_, tilt, _ := s.internalPosition(pos)

	log.Debugf("tilt: %x - %.2f", tilt, pos.Tilt)

	if _, err := s.Send(protocol.TiltSet, protocol.NoneReplay, tilt[0], tilt[1]); err != nil {
		return err
	}

	return nil
}

func (s *Serial) GotoZoom(pos *dsd.Position) error {
	_, _, zoom := s.internalPosition(pos)

	log.Debugf("zoom: %x - %.2f", zoom, pos.Zoom)

	if _, err := s.Send(protocol.ZoomSet, protocol.NoneReplay, zoom[0], zoom[1]); err != nil {
		return err
	}

	return nil
}

func (s *Serial) externalPosition(pan, tilt, zoom []byte) (*dsd.Position, error) {
	return &dsd.Position{
		Pan:  float64(int(pan[0])<<8|int(pan[1])) / float64(100),
		Tilt: float64(int(tilt[0])<<8|int(tilt[1])) / float64(100),
		Zoom: float64(int(zoom[0])<<8|int(zoom[1])) / float64(100),
	}, nil
}

func (s *Serial) internalPosition(pos *dsd.Position) (pan, tilt, zoom []byte) {
	pan = make([]byte, 2)
	tilt = make([]byte, 2)
	zoom = make([]byte, 2)

	binary.BigEndian.PutUint16(pan, uint16(pos.Pan*100))
	binary.BigEndian.PutUint16(tilt, uint16(pos.Tilt*100))
	binary.BigEndian.PutUint16(zoom, uint16(pos.Zoom*100))

	return
}
