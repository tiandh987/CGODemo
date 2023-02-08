package blp

import (
	"context"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial/protocol"
)

type ptzRepo interface {
	Version(ctx context.Context) (string, error)
	Model(ctx context.Context) (string, error)
	Stop(ctx context.Context) error
	Left(ctx context.Context) error
	Right(ctx context.Context) error
}

type ptzUseCase struct {
}

var _ ptzRepo = (*ptzUseCase)(nil)

func NewPtz() ptzRepo {
	return &ptzUseCase{}
}

func (p *ptzUseCase) Version(ctx context.Context) (string, error) {
	ver, err := serial.Send(protocol.Version, protocol.VersionReplay, 0x00, 0x00)
	if err != nil {
		return "", err
	}

	version := fmt.Sprintf("V%d.%d", ver[0], ver[1])

	log.Debugf("ver: %x, version: %s", ver, version)
	return version, nil
}

func (p *ptzUseCase) Model(ctx context.Context) (string, error) {

	return "", nil
}

func (p *ptzUseCase) Stop(ctx context.Context) error {
	if _, err := serial.Send(protocol.Stop, protocol.NoneReplay, 0x00, 0x00); err != nil {
		return err
	}
	return nil
}

func (p *ptzUseCase) Left(ctx context.Context) error {
	if _, err := serial.Send(protocol.Left, protocol.NoneReplay, 0x2f, 0x00); err != nil {
		return err
	}
	return nil
}

func (p *ptzUseCase) Right(ctx context.Context) error {
	if _, err := serial.Send(protocol.Right, protocol.NoneReplay, 0x2f, 0x00); err != nil {
		return err
	}
	return nil
}
