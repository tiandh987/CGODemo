package ptz

import (
	"fmt"
	"go.bug.st/serial"
	"log"
)

// pelco d byte names
const (
	SYNC     = 0
	ADDR     = 1
	COMMAND1 = 2
	COMMAND2 = 3
	DATA1    = 4
	DATA2    = 5
	CHECKSUM = 6
)

type pelcodMessage [7]byte

func newPelcodMessage(addr, cmd1, cmd2, data1, data2 byte) *pelcodMessage {
	msg := &pelcodMessage{}
	msg.sync()
	msg.addr(addr)
	msg.cmd1(cmd1)
	msg.cmd2(cmd2)
	msg.data1(data1)
	msg.data2(data2)
	msg.checksum()

	return msg
}

func (p *pelcodMessage) sync() {
	p[SYNC] = byte('\xff')
}

func (p *pelcodMessage) addr(addr byte) {
	p[ADDR] = addr
}

func (p *pelcodMessage) cmd1(cmd byte) {
	p[COMMAND1] = cmd
}

func (p *pelcodMessage) cmd2(cmd byte) {
	p[COMMAND2] = cmd
}

func (p *pelcodMessage) data1(data byte) {
	p[DATA1] = data
}

func (p *pelcodMessage) data2(data byte) {
	p[DATA2] = data
}

func (p *pelcodMessage) checksum() {
	sum := p[ADDR] + p[COMMAND1] + p[COMMAND2] + p[DATA1] + p[DATA2]
	mod := sum % 100
	p[CHECKSUM] = mod
}

func (p *pelcodMessage) toByteSlice() []byte {
	var res []byte
	res = append(res, p[SYNC], p[ADDR], p[COMMAND1], p[COMMAND2], p[DATA1], p[DATA2], p[CHECKSUM])
	return res
}

func read(port serial.Port) []byte {
	buff := make([]byte, 1024)

	n, err := port.Read(buff)
	if err != nil {
		log.Fatal(err)
	}

	if n == 0 {
		fmt.Println("receive value length is zero")
	}

	fmt.Printf("receive value: %x\n", buff[:n])

	return buff[:n]
}
