package dsd

type PowerUps struct {
	Enable       bool `json:"Enable" validate:"boolean"`
	Function     int  `json:"Function" validate:"required,gte=0,lte=5"`
	PresetID     int  `json:"PresetID" validate:"required,gte=1,lte=255"`
	TourID       int  `json:"TourID" validate:"required,gte=1,lte=8"`
	PatternID    int  `json:"PatternID" validate:"required,gte=1,lte=5"`
	LinearScanID int  `json:"LinearScanID" validate:"required,gte=1,lte=5"`
	RegionScanID int  `json:"RegionScanID" validate:"required,gte=1,lte=5"`
}

func NewPowerUps() *PowerUps {
	return &PowerUps{
		Enable:       false,
		Function:     0,
		PresetID:     1,
		TourID:       1,
		PatternID:    1,
		LinearScanID: 1,
		RegionScanID: 1,
	}
}

func (u *PowerUps) ConfigKey() string {
	return "PowerUps"
}
