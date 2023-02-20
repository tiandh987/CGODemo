package protocol

type InstructRepo interface {
	Instruct(ct CommandType, data1, data2 byte) []byte
	InstructLen() int
	CheckReplay(rt ReplayType, replay []byte) error
	ReplayData(replay []byte) []byte
}

type CommandType int

const (
	Version   CommandType = iota // 云台版本号
	Model                        // 云台型号
	Restart                      // 云台重启
	Stop                         // 停止
	Up                           // 上
	Down                         // 下
	Left                         // 左
	Right                        // 右
	LeftUp                       // 左上
	RightUp                      // 右上
	LeftDown                     // 左下
	RightDown                    // 右下
	ZoomAdd                      // 变倍 +
	ZoomSub                      // 变倍 -
	PanSet                       // 设置 Pan 位置
	PanGet                       // 查询 Pan 位置
	TiltSet                      // 设置 Tilt 位置
	TiltGet                      // 查询 Tilt 位置
	ZoomSet                      // 设置 Zoom 位置
	ZoomGet                      // 查询 Zoom 位置
)

type ReplayType int

const (
	NoneReplay    ReplayType = iota
	VersionReplay            // 查询云台型号应答
	ModelReplay              // 查询云台版本号应答
	PanReplay                // 查询 Pan 位置应答
	TiltReplay               // 查询 Tilt 位置应答
	ZoomReplay               // 查询 Zoom 位置应答
)
