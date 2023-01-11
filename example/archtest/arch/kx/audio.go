package kx

/*
#cgo LDFLAGS: -L ../../cgo/lib/play1 -lirhwp

#include "../../cgo/include/play1/IRAY_PLAY_API.h"
*/
import "C"
import "github.com/tiandh987/CGODemo/example/archtest/blp"

type kxUseCase struct {
}

var _ blp.AudioRepo = (*kxUseCase)(nil)

func (k *kxUseCase) Play() int {
	ret := C.Ir_Play1()

	return int(ret)
}

func New() blp.AudioRepo {
	return &kxUseCase{}
}
