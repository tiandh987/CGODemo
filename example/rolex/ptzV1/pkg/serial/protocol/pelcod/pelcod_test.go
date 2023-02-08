package pelcod

import "testing"

func TestSlice(t *testing.T) {
	instruct := make([]byte, 7)

	t.Logf("ins: %x, len: %d, cap: %d", instruct, len(instruct), cap(instruct))

	instruct[SYNC] = 0xff
	instruct[ADDR] = 1
	instruct[CMD1] = 0x01
	instruct[CMD2] = 0x02
	instruct[DATA1] = 0x03
	instruct[DATA2] = 0x04
	instruct[CHECKSUM] = (instruct[ADDR] + instruct[CMD1] + instruct[CMD2] + instruct[DATA1] + instruct[DATA2]) % 100

	t.Logf("ins: %x, len: %d, cap: %d", instruct, len(instruct), cap(instruct))

}
