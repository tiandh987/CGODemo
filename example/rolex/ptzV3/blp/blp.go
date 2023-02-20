package blp

import (
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/preset"
)

type Blp struct {
	basic  *basic.Basic
	preset *preset.Preset
}
