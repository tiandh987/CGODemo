package dsd

type IdleMotion struct {
	Enable          bool `json:"Enable" validate:"boolean"`
	Second          int  `json:"Second" validate:"required,gte=5,lte=720"`
	Function        int  `json:"Function" validate:"gte=0,lte=5"`
	PresetID        int  `json:"PresetID" validate:"required,gte=1,lte=255"`
	TourID          int  `json:"TourID" validate:"required,gte=1,lte=8"`
	PatternID       int  `json:"PatternID" validate:"required,gte=1,lte=5"`
	LinearScanID    int  `json:"LinearScanID" validate:"required,gte=1,lte=5"`
	RegionScanID    int  `json:"RegionScanID" validate:"required,gte=1,lte=5"`
	Running         bool `json:"Runing" validate:"boolean"`
	RunningFunction int  `json:"RunningFunction" validate:"required,gte=0,lte=5"`
}

func NewIdleMotion() *IdleMotion {
	return &IdleMotion{
		Enable:          false,
		Second:          5,
		Function:        0,
		PresetID:        1,
		TourID:          1,
		PatternID:       1,
		LinearScanID:    1,
		RegionScanID:    1,
		Running:         false,
		RunningFunction: 0,
	}
}

func (m *IdleMotion) ConfigKey() string {
	return "IdleMotion"
}

func (m *IdleMotion) Validate() error {
	if err := _validate.Struct(m); err != nil {
		return err
	}

	return nil
}
