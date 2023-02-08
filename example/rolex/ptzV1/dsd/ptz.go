package dsd

// Position 云台坐标与放大倍数
type Position struct {
	Pan  float64 `json:"Pan" validate:"required,gte=0,lt=360"`  // 水平坐标
	Tile float64 `json:"Tile" validate:"required,gte=0,lte=90"` // 垂直坐标
	Zoom float64 `json:"Zoom" validate:"required,gte=0,lte=20"` // 变倍
}
