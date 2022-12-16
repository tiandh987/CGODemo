package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tiandh987/CGODemo/example/go-bug-st-serial/ptz"
	"go.bug.st/serial"
	"log"
)

func main() {

	// Retrieve the port list
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	// Print the list of detected ports
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	// Open the first serial port detected at 9600bps N81
	mode := &serial.Mode{
		BaudRate: 57600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open("/dev/ttyS2", mode)
	if err != nil {
		log.Fatal(err)
	}

	// 启动 gin 服务器
	engine := gin.Default()

	ptz.InitRoute(engine, port)

	engine.Run(":8089")
}
