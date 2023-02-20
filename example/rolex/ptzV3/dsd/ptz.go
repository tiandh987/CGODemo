package dsd

// Position 云台坐标与放大倍数
type Position struct {
	Pan  float64 `json:"Pan" validate:"gte=0,lt=360"`  // 水平坐标
	Tilt float64 `json:"Tilt" validate:"gte=0,lte=90"` // 垂直坐标
	Zoom float64 `json:"Zoom" validate:"gte=1,lte=20"` // 变倍
}

func NewPosition() Position {
	return Position{
		Pan:  0,
		Tilt: 0,
		Zoom: 1,
	}
}

func (p *Position) Validate() error {
	if err := _validate.Struct(p); err != nil {
		return err
	}

	return nil
}

type Status struct {
	Moving       bool  `json:"Moving" validate:"boolean"`
	Trigger      int   `json:"Trigger" validate:"required,gte=0,lte=4"`
	Function     int   `json:"Function" validate:"required,gte=0,lte=5"`
	FunctionID   int   `json:"FunctionID" validate:"required"`
	TimingTaskID int   `json:"TimingTaskID" validate:"required,gte=1,lte=4"`
	StartTime    int64 `json:"StartTime" validate:"required"`
}
