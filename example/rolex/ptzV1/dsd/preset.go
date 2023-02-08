package dsd

type PresetPoint struct {
	Enable   bool     `json:"Enable" validate:"boolean"`             // 使能
	ID       int      `json:"ID" validate:"required,gte=1,lte=255"`  // 预置点id
	Name     string   `json:"Name" validate:"required,min=1,max=64"` // 预置点名称
	Position Position `json:"Position"`                              // 预置点的坐标和放大倍数
}
