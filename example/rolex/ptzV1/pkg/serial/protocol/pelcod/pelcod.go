package pelcod

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial/protocol"
)

const (
	SYNC     = iota // 同步位，固定为 0xff
	ADDR            // 地址
	CMD1            // 命令1
	CMD2            // 命令2
	DATA1           // 数据1
	DATA2           // 数据2
	CHECKSUM        // 校验和
)

var _cmdMap = map[protocol.CommandType][]byte{
	protocol.Version: {0x28, 0x08},
	protocol.Model:   {0x28, 0x10},
	protocol.Stop:    {0x00, 0x00},
	protocol.Left:    {0x00, 0x04},
	protocol.Right:   {0x00, 0x02},
}

var _replayMap = map[protocol.ReplayType][]byte{
	protocol.VersionReplay: {0x28, 0x0d},
	protocol.ModelReplay:   {0x28, 0x1c},
}

type pelcoDUseCase struct {
	address byte
}

var _ protocol.InstructRepo = (*pelcoDUseCase)(nil)

func NewPelcoDUseCase(address byte) protocol.InstructRepo {
	return &pelcoDUseCase{
		address: address,
	}
}

func (p *pelcoDUseCase) Instruct(ct protocol.CommandType, data1, data2 byte) []byte {
	instruct := make([]byte, 7)

	log.Debugf("Instruct: %x", instruct)

	instruct[SYNC] = 0xff
	instruct[ADDR] = p.address
	instruct[CMD1] = _cmdMap[ct][0]
	instruct[CMD2] = _cmdMap[ct][1]
	instruct[DATA1] = data1
	instruct[DATA2] = data2
	instruct[CHECKSUM] = (instruct[ADDR] + instruct[CMD1] + instruct[CMD2] + instruct[DATA1] + instruct[DATA2]) % 100

	log.Debugf("Instruct: %x", instruct)

	return instruct
}

func (p *pelcoDUseCase) InstructLen() int {
	return 7
}

func (p *pelcoDUseCase) CheckReplay(rt protocol.ReplayType, replay []byte) error {
	log.Debugf("CheckReplay param, replay type: %d, replay data: %x", rt, replay)

	if _replayMap[rt][0] == replay[CMD1] && _replayMap[rt][1] == replay[CMD2] {
		return nil
	}

	return errors.New("replay data is invalid")
}

func (p *pelcoDUseCase) ReplayData(replay []byte) []byte {
	log.Debugf("ReplayData param, replay data: %x", replay)

	return replay[DATA1:CHECKSUM]
}
