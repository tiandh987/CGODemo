package dsd

type Limit struct {
	LevelEnable   bool    `json:"LevelEnable" validate:"boolean"`        // 水平限位使能
	CheckLeft     int     `json:"CheckLeft" validate:"gte=0,lte=1"`      // 水平左校验位 点左限位为check=1
	CheckRight    int     `json:"CheckRight" validate:"gte=0,lte=1"`     // 水平右校验位 点右限位为check=1
	LeftBoundary  float64 `json:"LeftBoundary" validate:"gte=0,lt=360"`  // 左边界位置
	RightBoundary float64 `json:"RightBoundary" validate:"gte=0,lt=360"` // 右边界位置

	VerticalEnable bool    `json:"VerticalEnable" validate:"boolean"`    // 垂直限位使能
	CheckUp        int     `json:"CheckUp" validate:"gte=0,lte=1"`       // 水平上校验位 点左限位为check=1
	CheckDown      int     `json:"CheckDown" validate:"gte=0,lte=1"`     // 水平下校验位 点右限位为check=1
	DownBoundary   float64 `json:"DownBoundary" validate:"gte=0,lte=90"` // 下边界位置
	UpBoundary     float64 `json:"UpBoundary" validate:"gte=0,lte=90"`   // 上边界位置
}

func NewLimit() *Limit {
	return &Limit{
		LevelEnable:    false,
		CheckLeft:      0,
		CheckRight:     0,
		LeftBoundary:   0,
		RightBoundary:  0,
		VerticalEnable: false,
		CheckUp:        0,
		CheckDown:      0,
		DownBoundary:   0,
		UpBoundary:     0,
	}
}

func (l *Limit) ConfigKey() string {
	return "PTZLimit"
}
