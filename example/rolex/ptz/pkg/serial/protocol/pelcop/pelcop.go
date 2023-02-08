package pelcop

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial/protocol"
)

type pelcoPUseCase struct {
	address byte
}

var _ protocol.InstructRepo = (*pelcoPUseCase)(nil)

func NewPelcoPUseCase(address byte) protocol.InstructRepo {
	return &pelcoPUseCase{
		address: address,
	}
}

func (p pelcoPUseCase) InstructLen() int {
	//TODO implement me
	panic("implement me")
}

func (p pelcoPUseCase) CheckReplay(rt protocol.ReplayType, replay []byte) error {
	//TODO implement me
	panic("implement me")
}

func (p pelcoPUseCase) ReplayData(replay []byte) []byte {
	//TODO implement me
	panic("implement me")
}

func (p pelcoPUseCase) Instruct(ct protocol.CommandType, data1, data2 byte) []byte {
	//TODO implement me
	panic("implement me")
}
