package serial

import (
	"encoding/binary"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/arch/protocol"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
)

var _ ptz.AbilityRepo = (*SerialV2)(nil)

func (s *SerialV2) Version() (string, error) {
	ver, err := s.Send(protocol.Version, protocol.VersionReplay, 0x00, 0x00)
	if err != nil {
		return "", err
	}

	version := fmt.Sprintf("V%d.%d", ver[0], ver[1])

	log.Debugf("ver: %x, version: %s", ver, version)
	return version, nil
}

func (s *SerialV2) Model() (string, error) {
	mod, err := s.Send(protocol.Model, protocol.ModelReplay, 0x00, 0x00)
	if err != nil {
		return "", err
	}

	model := fmt.Sprintf("V%d.%d", mod[0], mod[1])

	log.Debugf("mod: %x, model: %s", mod, model)
	return model, nil
}

func (s *SerialV2) Restart() error {
	if _, err := s.Send(protocol.Restart, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) Stop() error {
	if _, err := s.Send(protocol.Stop, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) Up(speed ptz.Speed) error {
	if _, err := s.Send(protocol.Up, protocol.NoneReplay, 0x00, speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) Down(speed ptz.Speed) error {
	if _, err := s.Send(protocol.Down, protocol.NoneReplay, 0x00, speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) Left(speed ptz.Speed) error {
	if _, err := s.Send(protocol.Left, protocol.NoneReplay, speed.Convert(), 0x00); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) Right(speed ptz.Speed) error {
	if _, err := s.Send(protocol.Right, protocol.NoneReplay, speed.Convert(), 0x00); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) LeftUp(speed ptz.Speed) error {
	if _, err := s.Send(protocol.LeftUp, protocol.NoneReplay, speed.Convert(), speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) RightUp(speed ptz.Speed) error {
	if _, err := s.Send(protocol.RightUp, protocol.NoneReplay, speed.Convert(), speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) LeftDown(speed ptz.Speed) error {
	if _, err := s.Send(protocol.LeftDown, protocol.NoneReplay, speed.Convert(), speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) RightDown(speed ptz.Speed) error {
	if _, err := s.Send(protocol.RightDown, protocol.NoneReplay, speed.Convert(), speed.Convert()); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) ZoomAdd() error {
	if _, err := s.Send(protocol.ZoomAdd, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) ZoomSub() error {
	if _, err := s.Send(protocol.ZoomSub, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) Position() (*dsd.Position, error) {
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

func (s *SerialV2) Goto(pos *dsd.Position) error {
	log.Infof("xxxxxxxxxx Goto %+v", pos)

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

func (s *SerialV2) GotoPan(pos *dsd.Position) error {
	pan, _, _ := s.internalPosition(pos)

	log.Debugf("pan: %x - %.2f", pan, pos.Pan)

	if _, err := s.Send(protocol.PanSet, protocol.NoneReplay, pan[0], pan[1]); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) GotoTilt(pos *dsd.Position) error {
	_, tilt, _ := s.internalPosition(pos)

	log.Debugf("tilt: %x - %.2f", tilt, pos.Tilt)

	if _, err := s.Send(protocol.TiltSet, protocol.NoneReplay, tilt[0], tilt[1]); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) GotoZoom(pos *dsd.Position) error {
	_, _, zoom := s.internalPosition(pos)

	log.Debugf("zoom: %x - %.2f", zoom, pos.Zoom)

	if _, err := s.Send(protocol.ZoomSet, protocol.NoneReplay, zoom[0], zoom[1]); err != nil {
		return err
	}

	return nil
}

func (s *SerialV2) externalPosition(pan, tilt, zoom []byte) (*dsd.Position, error) {
	return &dsd.Position{
		Pan:  float64(int(pan[0])<<8|int(pan[1])) / float64(100),
		Tilt: float64(int(tilt[0])<<8|int(tilt[1])) / float64(100),
		Zoom: float64(int(zoom[0])<<8|int(zoom[1])) / float64(100),
	}, nil
}

func (s *SerialV2) internalPosition(pos *dsd.Position) (pan, tilt, zoom []byte) {
	pan = make([]byte, 2)
	tilt = make([]byte, 2)
	zoom = make([]byte, 2)

	binary.BigEndian.PutUint16(pan, uint16(pos.Pan*100))
	binary.BigEndian.PutUint16(tilt, uint16(pos.Tilt*100))
	binary.BigEndian.PutUint16(zoom, uint16(pos.Zoom*100))

	return
}
