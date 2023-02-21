package pelcod

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/arch/protocol"
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
	protocol.Version:   {0x28, 0x08},
	protocol.Model:     {0x28, 0x10},
	protocol.Restart:   {0x00, 0x0f},
	protocol.Stop:      {0x00, 0x00},
	protocol.Up:        {0x00, 0x08},
	protocol.Down:      {0x00, 0x10},
	protocol.Left:      {0x00, 0x04},
	protocol.Right:     {0x00, 0x02},
	protocol.LeftUp:    {0x00, 0x0c},
	protocol.RightUp:   {0x00, 0x0a},
	protocol.LeftDown:  {0x00, 0x14},
	protocol.RightDown: {0x00, 0x12},
	protocol.ZoomAdd:   {0x00, 0x20},
	protocol.ZoomSub:   {0x00, 0x40},
	protocol.PanSet:    {0x00, 0x4b},
	protocol.PanGet:    {0x00, 0x51},
	protocol.TiltSet:   {0x00, 0x4d},
	protocol.TiltGet:   {0x00, 0x53},
	protocol.ZoomSet:   {0x00, 0x4f},
	protocol.ZoomGet:   {0x00, 0x55},
}

var _replayCmdMap = map[protocol.ReplayType][]byte{
	protocol.VersionReplay: {0x28, 0x0d},
	protocol.ModelReplay:   {0x28, 0x1c},
	protocol.PanReplay:     {0x00, 0x59},
	protocol.TiltReplay:    {0x00, 0x5b},
	protocol.ZoomReplay:    {0x00, 0x5d},
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

	instruct[SYNC] = 0xff
	instruct[ADDR] = p.address
	instruct[CMD1] = _cmdMap[ct][0]
	instruct[CMD2] = _cmdMap[ct][1]
	instruct[DATA1] = data1
	instruct[DATA2] = data2
	instruct[CHECKSUM] = instruct[ADDR] + instruct[CMD1] + instruct[CMD2] + instruct[DATA1] + instruct[DATA2]

	//log.Debugf("command type: %d, data1: %x, data2: %x, instruct: %x", ct, data1, data2, instruct)

	return instruct
}

func (p *pelcoDUseCase) InstructLen() int {
	return 7
}

func (p *pelcoDUseCase) CheckReplay(rt protocol.ReplayType, replay []byte) error {
	if _replayCmdMap[rt][0] == replay[CMD1] && _replayCmdMap[rt][1] == replay[CMD2] {
		return nil
	}

	return errors.New("replay data is invalid")
}

func (p *pelcoDUseCase) ReplayData(replay []byte) []byte {
	return replay[DATA1:CHECKSUM]
}
