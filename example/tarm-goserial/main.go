package main

import (
	"fmt"
	"github.com/tarm/goserial"
	"log"
	"time"
)

func main() {
	cfg := &serial.Config{
		Name:        "/dev/ttyS2",
		Baud:        57600,
		ReadTimeout: 3,
	}

	iorwc, err := serial.OpenPort(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer iorwc.Close()

	stop := []byte("\xFF\x01\x00\x00\x00\x00\x01")
	left := []byte("\xFF\x01\x00\x04\x20\x00\x25")
	leftUp := []byte("\xFF\x01\x00\x0c\x20\x20\x4d")

	// 向左转动
	num, err := iorwc.Write(left)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 6)

	// 停止
	num, err = iorwc.Write(stop)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 向左上转动
	num, err = iorwc.Write(leftUp)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 3)

	// 停止
	num, err = iorwc.Write(stop)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(num)

	version := []byte("\xff\x01\x28\x10\x00\x00\x39")

	_, err = iorwc.Write(version)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 1024)
	lens, err := iorwc.Read(buf)
	if err != nil {
		log.Println(err)
	}

	revData := buf[:lens]
	fmt.Println("revData: %s -- %X\n", revData, revData)
}
