package dsd

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
)

// Protocol 协议
type Protocol int

const (
	PELCOD Protocol = iota
	PELCOP
)

func (p Protocol) Validate() error {
	log.Debugf("Protocol: %d", p)

	if p != PELCOD && p != PELCOP {
		log.Errorf("invalid ptz protocol. Protocol: %d", p)
		return errors.New("invalid ptz protocol")
	}

	return nil
}

// Address 地址
type Address int

const (
	Visible Address = 1 // 可见光
	Thermal Address = 2 // 热成像
)

func (a Address) Validate() error {
	log.Debugf("Address: %d", a)

	if a != Visible && a != Thermal {
		log.Errorf("invalid ptz address. Address: %d", a)
		return errors.New("invalid ptz address")
	}

	return nil
}

// BaudRate 波特率
type BaudRate int

const (
	BaudRate1200   BaudRate = 1200
	BaudRate2400   BaudRate = 2400
	BaudRate4800   BaudRate = 4800
	BaudRate9600   BaudRate = 9600
	BaudRate19200  BaudRate = 19200
	BaudRate38400  BaudRate = 38400
	BaudRate57600  BaudRate = 57600
	BaudRate115200 BaudRate = 115200
)

func (b BaudRate) Validate() error {
	log.Debugf("BaudRate: %d", b)

	if b != BaudRate1200 && b != BaudRate2400 && b != BaudRate4800 && b != BaudRate9600 && b != BaudRate19200 &&
		b != BaudRate38400 && b != BaudRate57600 && b != BaudRate115200 {
		log.Errorf("invalid baud rate. BaudRate: %d", b)
		return errors.New("invalid baud rate")
	}

	return nil
}

// DataBit 数据位
type DataBit int

const (
	FiveDataBit  DataBit = 5
	SixDataBit   DataBit = 6
	SevenDataBit DataBit = 7
	EightDataBit DataBit = 8
)

func (d DataBit) Validate() error {
	log.Debugf("DataBit: %d", d)

	if d != FiveDataBit && d != SixDataBit && d != SevenDataBit && d != EightDataBit {
		log.Errorf("invalid data bit. DataBit: %d", d)
		return errors.New("invalid data bit")
	}

	return nil
}

// StopBit 停止位
type StopBit int

const (
	OneStopBit StopBit = iota
	OnePointFiveStopBits
	TwoStopBits
)

func (s StopBit) Validate() error {
	log.Debugf("StopBit: %d", s)

	if s != OneStopBit && s != OnePointFiveStopBits && s != TwoStopBits {
		log.Errorf("invalid stop bit. StopBit: %d", s)
		return errors.New("invalid stop bit")
	}

	return nil
}

// ParityBit 校验位
type ParityBit int

const (
	NoParity    ParityBit = iota // 无
	OddParity                    // 奇校验
	EvenParity                   // 偶校验
	MarkParity                   // 标志校验
	SpaceParity                  // 空校验
)

func (p ParityBit) Validate() error {
	log.Debugf("ParityBit: %d", p)

	if p != NoParity && p != OddParity && p != EvenParity && p != MarkParity && p != SpaceParity {
		log.Errorf("invalid parity bit. ParityBit: %d", p)
		return errors.New("invalid parity bit")
	}

	return nil
}

// CommAttribute 串口属性
type CommAttribute struct {
	BaudRate BaudRate  `json:"BodeRate"` // 波特率（1200，2400， 4800，9600，19200，38400，57600，115200）
	DataBits DataBit   `json:"DataBits"` // 数据位数（5，6，7，8）
	Parity   ParityBit `json:"Parity"`   // 奇偶校验选项（0：无，1：奇校验，2：偶校验，3：标志校验，4：空校验）
	StopBits StopBit   `json:"StopBits"` // 停止位（0：停止位1；1：停止位1.5；2：停止位2）
}

func NewCommAttribute() *CommAttribute {
	return &CommAttribute{
		BaudRate: BaudRate57600,
		DataBits: EightDataBit,
		Parity:   NoParity,
		StopBits: OneStopBit,
	}
}

func (c *CommAttribute) Validate() error {
	if err := c.BaudRate.Validate(); err != nil {
		return err
	}

	if err := c.DataBits.Validate(); err != nil {
		return err
	}

	if err := c.Parity.Validate(); err != nil {
		return err
	}

	if err := c.StopBits.Validate(); err != nil {
		return err
	}

	return nil
}

// PTZ 串口配置
type PTZ struct {
	Enable    bool           `json:"Enable"`    // 使能
	Protocol  Protocol       `json:"Protocol"`  // 协议类型: 0 PELCOD, 1 PELCOP
	Address   Address        `json:"Address"`   // PELCO协议: 485地址,0-255
	Attribute *CommAttribute `json:"Attribute"` // 串口属性
}

func NewPTZ() *PTZ {
	return &PTZ{
		Enable:    false,
		Protocol:  PELCOD,
		Address:   Visible,
		Attribute: NewCommAttribute(),
	}
}

func (p *PTZ) Validate() error {
	if err := p.Protocol.Validate(); err != nil {
		return err
	}

	if err := p.Address.Validate(); err != nil {
		return err
	}

	if err := p.Attribute.Validate(); err != nil {
		return err
	}

	return nil
}

// ConfigKey 返回用于读取配置文件的 Key
func (p *PTZ) ConfigKey() string {
	return "Serial"
}

func (s *PTZ) ConvertAddress() byte {
	return byte(s.Address)
}
