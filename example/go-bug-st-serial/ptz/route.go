package ptz

import (
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.bug.st/serial"
	"log"
	"strconv"
)

// PTSpeed 云台速度
type PTSpeed int

func (s PTSpeed) Convert() byte {
	return _ptSpeedMap[s]
}

const (
	PTSpeedOne PTSpeed = iota + 1
	PTSpeedTwo
	PTSpeedThree
	PTSpeedFour
	PTSpeedFive
	PTSpeedSix
	PTSpeedSeven
	PTSpeedEight
)

var _ptSpeedMap = map[PTSpeed]byte{
	PTSpeedOne:   byte(0x01),
	PTSpeedTwo:   byte(0x09),
	PTSpeedThree: byte(0x12),
	PTSpeedFour:  byte(0x1b),
	PTSpeedFive:  byte(0x24),
	PTSpeedSix:   byte(0x2d),
	PTSpeedSeven: byte(0x36),
	PTSpeedEight: byte(0x3f),
}

func InitRoute(engine *gin.Engine, port serial.Port) {
	restart(engine, port)
	calibration(engine, port)
	model(engine, port)
	version(engine, port)
	control(engine, port)
	preset(engine, port)
	lineScan(engine, port)
	position(engine, port)

	panMove(engine, port)
}

func restart(engine *gin.Engine, port serial.Port) {
	restart := engine.Group("/restart")

	// 云台重启
	restart.POST("", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x0f, 0x00, 0x00).toByteSlice()
		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes(%x)\n", n, bytes)
	})

	// 云台断电重启
	restart.POST("/blackout", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0xf2, 0x21, 0x00, 0x00).toByteSlice()
		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes(%x)\n", n, bytes)
	})
}

// 自动原点校准
func calibration(engine *gin.Engine, port serial.Port) {
	calibration := engine.Group("/calibration")

	// 自动校准
	calibration.POST("", func(context *gin.Context) {
		query := context.Query("mode")

		fmt.Printf("query: %s\n", query)

		mode, _ := strconv.Atoi(query)

		fmt.Printf("mode: %d\n", mode)

		bytes := newPelcodMessage(0x01, 0x20, 0x25, 0x00, byte(mode)).toByteSlice()
		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes(%x)\n", n, bytes)
	})
}

func model(engine *gin.Engine, port serial.Port) {
	model := engine.Group("/model")

	// 查询云台型号
	model.GET("", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x28, 0x10, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("GET model Sent %v bytes (%x)\n", n, bytes)

		read(port)
	})
}

func version(engine *gin.Engine, port serial.Port) {
	version := engine.Group("/version")

	// 查询云台型号
	version.GET("", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x28, 0x08, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)

		read(port)
	})
}

// 云台控制
func control(engine *gin.Engine, port serial.Port) {
	control := engine.Group("/control")
	// 停止
	control.POST("/stop", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x00, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 上
	control.POST("/up", func(context *gin.Context) {
		query := context.Query("speed")

		fmt.Printf("query: %s\n", query)

		intSpeed, _ := strconv.Atoi(query)

		fmt.Printf("intSpeed: %d\n", intSpeed)

		speed := PTSpeed(intSpeed).Convert()

		fmt.Printf("speed: %d\n", speed)

		bytes := newPelcodMessage(0x01, 0x00, 0x08, 0x00, speed).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 下
	control.POST("/down", func(context *gin.Context) {
		query := context.Query("speed")

		speed, _ := strconv.Atoi(query)

		bytes := newPelcodMessage(0x01, 0x00, 0x10, 0x00, PTSpeed(speed).Convert()).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 左
	control.POST("/left", func(context *gin.Context) {
		query := context.Query("speed")

		speed, _ := strconv.Atoi(query)

		bytes := newPelcodMessage(0x01, 0x00, 0x04, PTSpeed(speed).Convert(), 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 右
	control.POST("/right", func(context *gin.Context) {
		query := context.Query("speed")

		speed, _ := strconv.Atoi(query)

		bytes := newPelcodMessage(0x01, 0x00, 0x02, PTSpeed(speed).Convert(), 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 左上
	control.POST("/leftUp", func(context *gin.Context) {
		pQuery := context.Query("pSpeed")
		tQuery := context.Query("tSpeed")

		pSpeed, _ := strconv.Atoi(pQuery)
		tSpeed, _ := strconv.Atoi(tQuery)

		bytes := newPelcodMessage(0x01, 0x00, 0x0c, PTSpeed(pSpeed).Convert(), PTSpeed(tSpeed).Convert()).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 右上
	control.POST("/rightUp", func(context *gin.Context) {
		pQuery := context.Query("pSpeed")
		tQuery := context.Query("tSpeed")

		pSpeed, _ := strconv.Atoi(pQuery)
		tSpeed, _ := strconv.Atoi(tQuery)

		bytes := newPelcodMessage(0x01, 0x00, 0x0a, PTSpeed(pSpeed).Convert(), PTSpeed(tSpeed).Convert()).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 左下
	control.POST("/leftDown", func(context *gin.Context) {
		pQuery := context.Query("pSpeed")
		tQuery := context.Query("tSpeed")

		pSpeed, _ := strconv.Atoi(pQuery)
		tSpeed, _ := strconv.Atoi(tQuery)

		bytes := newPelcodMessage(0x01, 0x00, 0x14, PTSpeed(pSpeed).Convert(), PTSpeed(tSpeed).Convert()).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 右下
	control.POST("/rightDown", func(context *gin.Context) {
		pQuery := context.Query("pSpeed")
		tQuery := context.Query("tSpeed")

		pSpeed, _ := strconv.Atoi(pQuery)
		tSpeed, _ := strconv.Atoi(tQuery)

		bytes := newPelcodMessage(0x01, 0x00, 0x12, PTSpeed(pSpeed).Convert(), PTSpeed(tSpeed).Convert()).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 变倍+
	control.POST("/zoom/tele", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x20, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 变倍-
	control.POST("/zoom/wide", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x40, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 聚焦+
	control.POST("/focus/near", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x01, 0x00, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 聚焦-
	control.POST("/focus/far", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x80, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 光圈+
	control.POST("/iris/open", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x02, 0x00, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 光圈-
	control.POST("/iris/close", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x04, 0x00, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

}

// 预置点
func preset(engine *gin.Engine, port serial.Port) {
	preset := engine.Group("/preset")

	// 设置
	preset.POST("/set", func(context *gin.Context) {
		query := context.Query("id")

		id, _ := strconv.Atoi(query)

		bytes := newPelcodMessage(0x01, 0x00, 0x03, 0x00, byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 删除
	preset.POST("/del", func(context *gin.Context) {
		query := context.Query("id")

		id, _ := strconv.Atoi(query)

		bytes := newPelcodMessage(0x01, 0x00, 0x05, 0x00, byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 调用
	preset.POST("/call", func(context *gin.Context) {
		query := context.Query("id")

		id, _ := strconv.Atoi(query)

		bytes := newPelcodMessage(0x01, 0x00, 0x07, 0x00, byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 速度
	preset.POST("/speed", func(context *gin.Context) {
		query := context.Query("speed")

		speed, _ := strconv.Atoi(query)

		bytes := newPelcodMessage(0x01, 0x00, 0x68, 0x00, byte(speed)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})
}

// 线性扫描
func lineScan(engine *gin.Engine, port serial.Port) {
	line := engine.Group("/linescan")

	// 设置线扫左边界停留时间
	line.POST("/left/stay", func(context *gin.Context) {
		tQuery := context.Query("time")
		time, _ := strconv.Atoi(tQuery)

		idQuery := context.Query("id")
		id, _ := strconv.Atoi(idQuery)

		bytes := newPelcodMessage(0x01, 0x20, 0x03, byte(time), byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 设置线扫右边界停留时间
	line.POST("/right/stay", func(context *gin.Context) {
		tQuery := context.Query("time")
		time, _ := strconv.Atoi(tQuery)

		idQuery := context.Query("id")
		id, _ := strconv.Atoi(idQuery)

		bytes := newPelcodMessage(0x01, 0x20, 0x05, byte(time), byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 设置线扫左边界
	line.POST("/left/bound", func(context *gin.Context) {
		idQuery := context.Query("id")
		id, _ := strconv.Atoi(idQuery)

		bytes := newPelcodMessage(0x01, 0x20, 0x09, 0x00, byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 设置线扫右边界
	line.POST("/right/bound", func(context *gin.Context) {
		idQuery := context.Query("id")
		id, _ := strconv.Atoi(idQuery)

		bytes := newPelcodMessage(0x01, 0x20, 0x0b, 0x00, byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 开始线扫
	line.POST("/start", func(context *gin.Context) {
		idQuery := context.Query("id")
		id, _ := strconv.Atoi(idQuery)

		bytes := newPelcodMessage(0x01, 0x20, 0x0d, 0x00, byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 删除线扫
	line.POST("/del", func(context *gin.Context) {
		idQuery := context.Query("id")
		id, _ := strconv.Atoi(idQuery)

		bytes := newPelcodMessage(0x01, 0x20, 0x0f, 0x00, byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 停止线扫
	line.POST("/stop", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x00, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 线扫速度
	line.POST("/speed", func(context *gin.Context) {
		tSpeed := context.Query("speed")
		speed, _ := strconv.Atoi(tSpeed)

		idQuery := context.Query("id")
		id, _ := strconv.Atoi(idQuery)

		bytes := newPelcodMessage(0x01, 0x20, 0x07, PTSpeed(speed).Convert(), byte(id)).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})
}

// 精准定位
func position(engine *gin.Engine, port serial.Port) {
	position := engine.Group("/position")

	// 查询 Pan 位置
	position.GET("/pan", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x51, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)

		replay := read(port)

		fmt.Printf("pan position: %f\n", float64(int(replay[4])<<8|int(replay[5]))/float64(100))
		fmt.Printf("pan position: %d\n", binary.BigEndian.Uint16(replay[4:6]))
	})

	// 查询 tilt 位置
	position.GET("/tilt", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x53, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)

		replay := read(port)

		fmt.Printf("tilt position: %f\n", float64(int(replay[4])<<8|int(replay[5]))/float64(100))
		fmt.Printf("tilt position: %d\n", binary.BigEndian.Uint16(replay[4:6]))
	})

	// 查询 zoom 位置
	position.GET("/zoom", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x00, 0x55, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)

		replay := read(port)

		fmt.Printf("zoom position: %f\n", float64(int(replay[4])<<8|int(replay[5]))/float64(100))
		fmt.Printf("zoom position: %d\n", binary.BigEndian.Uint16(replay[4:6]))
	})

	// 设置 Pan 位置
	position.POST("/pan", func(context *gin.Context) {
		query := context.Query("position")

		fmt.Printf("query: %s\n", query)

		pos, _ := strconv.Atoi(query)

		fmt.Printf("pos: %d\n", pos)

		var buf = make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(pos))

		fmt.Printf("buf: %x\n", buf)

		bytes := newPelcodMessage(0x01, 0x00, 0x4b, buf[0], buf[1]).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 设置 Tilt 位置
	position.POST("/tilt", func(context *gin.Context) {
		query := context.Query("position")

		fmt.Printf("query: %s\n", query)

		pos, _ := strconv.Atoi(query)

		fmt.Printf("pos: %d\n", pos)

		var buf = make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(pos))

		fmt.Printf("buf: %x\n", buf)

		bytes := newPelcodMessage(0x01, 0x00, 0x4d, buf[0], buf[1]).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 设置 Zoom 位置
	position.POST("/zoom", func(context *gin.Context) {
		query := context.Query("position")

		fmt.Printf("query: %s\n", query)

		pos, _ := strconv.Atoi(query)

		fmt.Printf("pos: %d\n", pos)

		var buf = make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(pos))

		fmt.Printf("buf: %x\n", buf)

		bytes := newPelcodMessage(0x01, 0x00, 0x4f, buf[0], buf[1]).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})
}

// 水平旋转
func panMove(engine *gin.Engine, port serial.Port) {
	panMove := engine.Group("/pan/move")

	// 速度
	panMove.POST("/speed", func(context *gin.Context) {
		query := context.Query("speed")

		fmt.Printf("query: %s\n", query)

		intSpeed, _ := strconv.Atoi(query)

		fmt.Printf("intSpeed: %d\n", intSpeed)

		speed := PTSpeed(intSpeed).Convert()

		fmt.Printf("speed: %d\n", speed)

		bytes := newPelcodMessage(0x01, 0x20, 0x17, speed, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 开始
	panMove.POST("/start", func(context *gin.Context) {
		query := context.Query("direction")

		fmt.Printf("query: %s\n", query)

		direction, _ := strconv.Atoi(query)

		fmt.Printf("direction: %d\n", direction)

		bytes := newPelcodMessage(0x01, 0x20, 0x29, byte(direction), 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})

	// 停止
	panMove.POST("/stop", func(context *gin.Context) {
		bytes := newPelcodMessage(0x01, 0x20, 0x29, 0x00, 0x00).toByteSlice()

		n, err := port.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %d bytes (%x)\n", n, bytes)
	})
}
