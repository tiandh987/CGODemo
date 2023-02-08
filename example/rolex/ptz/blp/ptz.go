package blp

import (
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial/protocol"
)

type ptzRepo interface {
	Version() (string, error)
	Model() (string, error)
}

type ptzUseCase struct {
}

var _ ptzRepo = (*ptzUseCase)(nil)

func NewPtz() ptzRepo {
	return &ptzUseCase{}
}

func (p *ptzUseCase) Version() (string, error) {
	ver, err := serial.Send(protocol.Version, protocol.VersionReplay, 0x00, 0x00)
	if err != nil {
		return "", err
	}

	version := fmt.Sprintf("V%d.%d", ver[0], ver[1])

	log.Debugf("ver: %x, version: %s", ver, version)
	return version, nil
}

func (p *ptzUseCase) Model() (string, error) {

	return "", nil
}
