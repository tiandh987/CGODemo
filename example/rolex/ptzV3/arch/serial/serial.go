package serial

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/arch/protocol"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/arch/protocol/pelcod"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	goSerial "go.bug.st/serial"
	"sync"
	"time"
)

type Serial struct {
	mu       sync.Mutex
	port     goSerial.Port
	protocol protocol.InstructRepo
}

func New(comName string, ptz *dsd.PTZ) *Serial {
	log.Debugf("new serial, ptz: %+v, attribute: %+v", ptz, ptz.Attribute)

	if comName == "" || len(comName) == 0 {
		comName = "/dev/ttyS2"
	}

	if ptz == nil {
		ptz = dsd.NewPTZ()
	}

	proto := pelcod.NewPelcoDUseCase(ptz.ConvertAddress())

	// TODO 支持 pelcop 协议
	//if ptz.Protocol == dsd.PELCOP {
	//	proto = pelcop.NewPelcoPUseCase(ptz.ConvertAddress())
	//}

	// Retrieve the port list
	ports, err := goSerial.GetPortsList()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(ports) == 0 {
		log.Error("No serial ports found!")
	}

	for _, port := range ports {
		log.Debugf("Found port: %v\n", port)
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

	return &Serial{
		port:     port,
		protocol: proto,
	}
}

func (s *Serial) Set(ptz *dsd.PTZ) error {
	log.Debugf("set serial mode and protocol, ptz: %+v, attribute: %+v", ptz, ptz.Attribute)

	if ptz.Protocol == dsd.PELCOD {
		s.protocol = pelcod.NewPelcoDUseCase(ptz.ConvertAddress())
	}
	//else if ptz.Protocol == dsd.PELCOP {
	//	s.protocol = pelcop.NewPelcoPUseCase(ptz.ConvertAddress())
	//}

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

func (s *Serial) Send(ct protocol.CommandType, rt protocol.ReplayType, data1, data2 byte) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	instruct := s.protocol.Instruct(ct, data1, data2)
	n, err := s.port.Write(instruct)
	if err != nil {
		log.Errorf("write [%x] to serial failed, err: %s", instruct, err.Error())
		return nil, errors.New("write to serial failed")
	}

	log.Debugf("write instruct: %d %x", n, instruct)

	if rt == protocol.NoneReplay {
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
		//log.Debugf("read data length: %d, data: %x", n, buff[:n])

		if err != nil {
			log.Errorf("read replay data from serial failed, err: %s", err.Error())
			continue
		}

		if err := s.protocol.CheckReplay(rt, buff[:n]); err != nil {
			log.Warnf("replay data is invalid, replay type: %d, replay data: %x", rt, buff[:n])
			continue
		}

		replay = s.protocol.ReplayData(buff[:n])
		break
	}

	return replay, nil
}
