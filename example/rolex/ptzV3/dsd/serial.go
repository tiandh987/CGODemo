package dsd

// Protocol 云台协议类型
type Protocol int

const (
	PELCOD Protocol = iota
	PELCOP
)

// Address 地址
type Address int

const (
	Visible Address = 1 // 可见光
	Thermal Address = 2 // 热成像
)

// BaudRate 串口波特率
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

// DataBit 数据位
type DataBit int

const (
	FiveDataBit  DataBit = 5
	SixDataBit   DataBit = 6
	SevenDataBit DataBit = 7
	EightDataBit DataBit = 8
)

// StopBit 停止位
type StopBit int

const (
	OneStopBit StopBit = iota
	OnePointFiveStopBits
	TwoStopBits
)

// ParityBit 校验位
type ParityBit int

const (
	NoParity    ParityBit = iota // 无
	OddParity                    // 奇校验
	EvenParity                   // 偶校验
	MarkParity                   // 标志校验
	SpaceParity                  // 空校验
)

// CommAttribute 串口属性
type CommAttribute struct {
	BaudRate BaudRate  `json:"BodeRate" validate:"oneof=1200 2400 4800 9600 19200 38400 57600 115200"`
	DataBits DataBit   `json:"DataBits" validate:"oneof=5 6 7 8"`
	Parity   ParityBit `json:"Parity" validate:"oneof=0 1 2 3 4"`
	StopBits StopBit   `json:"StopBits" validate:"oneof=0 1 2"`
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
	if err := _validate.Struct(c); err != nil {
		return err
	}

	return nil
}

// PTZ 串口配置
type PTZ struct {
	Enable    bool           `json:"Enable"`
	Protocol  Protocol       `json:"Protocol" validate:"oneof=0 1"`
	Address   Address        `json:"Address" validate:"oneof=1 2"`
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
	if err := _validate.Struct(p); err != nil {
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
