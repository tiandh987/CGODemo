package nt

/*
#cgo LDFLAGS: -L ../../cgo/lib/play1 -lirhwp

#include "../../cgo/include/play1/IRAY_PLAY_API.h"
*/
import "C"
import "github.com/tiandh987/CGODemo/example/archtest/blp"

type ntUseCase struct {
}

var _ blp.AudioRepo = (*ntUseCase)(nil)

func (k *ntUseCase) Play() int {
	ret := C.Ir_Play2()

	return int(ret)
}

func New() blp.AudioRepo {
	return &ntUseCase{}
}
