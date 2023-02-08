package protocol

type InstructRepo interface {
	Instruct(ct CommandType, data1, data2 byte) []byte
	InstructLen() int
	CheckReplay(rt ReplayType, replay []byte) error
	ReplayData(replay []byte) []byte
}

type CommandType int

const (
	Version CommandType = iota
	Model
)

type ReplayType int

const (
	NoneReplay ReplayType = iota
	VersionReplay
	ModelReplay
)
