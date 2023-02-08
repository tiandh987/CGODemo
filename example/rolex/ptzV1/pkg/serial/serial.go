package serial

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/dsd"
	protocol2 "github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial/protocol"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial/protocol/pelcod"
	"github.com/tiandh987/CGODemo/example/rolex/ptz/pkg/serial/protocol/pelcop"
	goSerial "go.bug.st/serial"
	"sync"
	"time"
)

var (
	_instanceMu sync.Mutex
	_serialIns  = new("/dev/ttyS2", dsd.NewPTZ())
)

type serial struct {
	mu       sync.Mutex
	port     goSerial.Port
	protocol protocol2.InstructRepo
}

func new(comName string, ptz *dsd.PTZ) *serial {
	log.Debugf("new serial, ptz: %+v, attribute: %+v", ptz, ptz.Attribute)

	if comName == "" || len(comName) == 0 {
		comName = "/dev/ttyS2"
	}

	if ptz == nil {
		ptz = dsd.NewPTZ()
	}

	proto := pelcod.NewPelcoDUseCase(ptz.ConvertAddress())
	if ptz.Protocol == dsd.PELCOP {
		proto = pelcop.NewPelcoPUseCase(ptz.ConvertAddress())
	}

	// Retrieve the port list
	ports, err := goSerial.GetPortsList()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	for _, port := range ports {
		log.Infof("Found port: %v\n", port)
	}

	mode := &goSerial.Mode{
		BaudRate: int(ptz.Attribute.BaudRate),
		DataBits: int(ptz.Attribute.DataBits),
		Parity:   goSerial.Parity(ptz.Attribute.Parity),
		StopBits: goSerial.StopBits(ptz.Attribute.StopBits),
	}

	log.Debugf("comName: %s, mode: %+v", comName, mode)

	port, err := goSerial.Open(comName, mode)
	if err != nil {
		log.Errorf("open serial failed, err: %s", err.Error())
		return nil
	}
	port.SetReadTimeout(time.Second * 5)

	return &serial{
		port:     port,
		protocol: proto,
	}
}

func (s *serial) Set(ptz *dsd.PTZ) error {
	log.Debugf("set serial mode and protocol, ptz: %+v, attribute: %+v", ptz, ptz.Attribute)

	if ptz.Protocol == dsd.PELCOD {
		s.protocol = pelcod.NewPelcoDUseCase(ptz.ConvertAddress())
	} else if ptz.Protocol == dsd.PELCOP {
		s.protocol = pelcop.NewPelcoPUseCase(ptz.ConvertAddress())
	}

	mode := goSerial.Mode{
		BaudRate: int(ptz.Attribute.BaudRate),
		DataBits: int(ptz.Attribute.DataBits),
		Parity:   goSerial.Parity(ptz.Attribute.Parity),
		StopBits: goSerial.StopBits(ptz.Attribute.StopBits),
	}

	if err := s.port.SetMode(&mode); err != nil {
		log.Errorf("set serial mode failed, err: %s", err.Error())
		return errors.New("set serial mode failed")
	}

	return nil
}

func (s *serial) Send(ct protocol2.CommandType, rt protocol2.ReplayType, data1, data2 byte) ([]byte, error) {
	log.Debugf("serial write param, CommandType: %d, ReplayType: %d, data1: %x, data2: %x",
		ct, rt, data1, data2)

	instruct := s.protocol.Instruct(ct, data1, data2)

	log.Debugf("instruct: %x", instruct)

	n, err := s.port.Write(instruct)
	if err != nil {
		log.Errorf("write [%x] to serial failed, err: %s", instruct, err.Error())
		return nil, errors.New("write to serial failed")
	}
	log.Debugf("write serial number: %d", n)

	if rt == protocol2.NoneReplay {
		return nil, nil
	}

	var replay []byte
	buff := make([]byte, 1024)

	for {
		n, err := s.port.Read(buff)
		if n != s.protocol.InstructLen() {
			log.Warnf("read replay data length is invalid. len: %d, data: %x", n, buff[:n])
			continue
		}

		log.Debugf("read data length: %d, data: %x", n, buff[:n])

		if err != nil {
			log.Errorf("read replay data from serial failed, err: %s", err.Error())
			continue
		}

		if err := s.protocol.CheckReplay(rt, buff[:n]); err != nil {
			log.Warnf("replay data is invalid, replay type: %d, replay data: %x", rt, buff[:n])
			continue
		}

		replay = s.protocol.ReplayData(buff[:n])

		log.Debugf("get replay data: %x", replay)
		break
	}

	log.Debugf("query success, replay: %x", replay)

	return replay, nil
}

func Init(ptz *dsd.PTZ) {
	_instanceMu.Lock()
	_serialIns = new("", ptz)
	_instanceMu.Unlock()
}

func Set(ptz *dsd.PTZ) error {
	return _serialIns.Set(ptz)
}

func Send(ct protocol2.CommandType, rt protocol2.ReplayType, data1, data2 byte) ([]byte, error) {
	return _serialIns.Send(ct, rt, data1, data2)
}
