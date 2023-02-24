package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	pkgVersion "github.com/tiandh987/CGODemo/example/rolex/pkg/version"
	ptzV3 "github.com/tiandh987/CGODemo/example/rolex/ptzV3"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/basic"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"net/http"
	"strconv"
	"time"
)

// loginResponse
type Response struct {
	Code      int64       `json:"Code"`      // 200
	Data      interface{} `json:"Data"`      // 返回数据
	Detail    string      `json:"Detail"`    // 详细错误信息
	Message   string      `json:"Message"`   // msg: success
	Translate string      `json:"Translate"` // 操作成功
}

// GOOS=linux GOARCH=arm go build -a && mv rolex rolex_test && scp rolex_test tiandenghao@172.17.132.250:/home/tiandenghao/nfsroot/rolex_nb
func main() {
	pflag.Parse()
	pkgVersion.PrintAndExitIfRequested()

	// 应用程序启动文件初始化
	appOpt := &config.Option{
		Path: "/mnt/custom/tian/rolex_nb/config",
		Name: "rolex",
		Type: "yaml",
	}
	appCfg := config.New(appOpt)

	// 日志包初始化
	logOpt := log.NewOptions()
	if err := appCfg.UnmarshalKey("log", logOpt); err != nil {
		if err := appCfg.Set("log", logOpt); err != nil {
			panic(err)
		}
	}
	if errs := logOpt.Validate(); errs != nil {
		panic(errs)
	}

	log.Init(logOpt)
	defer log.Flush()

	engine := gin.Default()

	if err := ptzV3.Start(); err != nil {
		panic(err)
	}

	blpInstance := ptzV3.Instance()

	tokenGroup := engine.Group("/v1/token")
	tokenGroup.POST("", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      "token",
			Detail:    "",
			Message:   "",
			Translate: "",
		})
	})

	ptzGroup := engine.Group("/v1/ptz")

	// 查询云台版本
	ptzGroup.GET("/ptzversion", func(c *gin.Context) {
		version := blpInstance.Basic().Version()

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      version,
			Detail:    "",
			Message:   "",
			Translate: "",
		})
	})

	// 查询云台型号
	ptzGroup.GET("/ptzmodel", func(c *gin.Context) {
		model := blpInstance.Basic().Model()

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      model,
			Detail:    "",
			Message:   "",
			Translate: "",
		})
	})

	// 查询云台状态
	ptzGroup.GET("/status", func(c *gin.Context) {
		state := blpInstance.Manager().State()

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      state,
			Detail:    "",
			Message:   "",
			Translate: "",
		})
	})

	// 云台转动-单击
	ptzGroup.PUT("/turn", func(c *gin.Context) {
		direction := c.Query("direction")
		dirNum, err := strconv.Atoi(direction)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		speed := c.Query("speed")
		speedNum, err := strconv.Atoi(speed)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := basic.Operation(dirNum).ValidateDirection(); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := ptz.Speed(speedNum).Validate(); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		req := blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.ManualFunc,
			ID:      dirNum,
			Speed:   speedNum,
		}

		log.Infof("req: %+v", req)

		if err := blpInstance.Manager().Start(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		time.Sleep(time.Millisecond * 200)
		//time.Sleep(time.Second * 5)

		req2 := blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.ManualFunc,
			ID:      dirNum,
			Speed:   speedNum,
		}

		log.Infof("req: %+v", req)

		if err := blpInstance.Manager().Stop(&req2); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      nil,
			Detail:    "",
			Message:   "turn ok",
			Translate: "",
		})
		return
	})

	// 云台转动-长按
	ptzGroup.PUT("/moveContinuously", func(c *gin.Context) {
		direction := c.Query("direction")
		dirNum, err := strconv.Atoi(direction)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		speed := c.Query("speed")
		speedNum, err := strconv.Atoi(speed)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := basic.Operation(dirNum).ValidateDirection(); err != nil {
			log.Error(err.Error())

			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := ptz.Speed(speedNum).Validate(); err != nil {
			log.Error(err.Error())

			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		req := blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.ManualFunc,
			ID:      dirNum,
			Speed:   speedNum,
		}

		if err := blpInstance.Manager().Start(&req); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "moveContinuously ok")
		return
	})

	// 云台变倍
	ptzGroup.PUT("/operation", func(c *gin.Context) {
		method := c.Query("method")
		methodNum, err := strconv.Atoi(method)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := basic.Operation(methodNum).ValidateOperation(); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		req := blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.ManualFunc,
			ID:      methodNum,
			Speed:   1,
		}

		if err := blpInstance.Manager().Start(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "operation ok")
		return
	})

	// 云台停止转动
	ptzGroup.PUT("/stop", func(c *gin.Context) {
		req := blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.ManualFunc,
			ID:      blpInstance.Manager().State().FunctionID,
			Speed:   1,
		}

		if err := blpInstance.Manager().Stop(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "stop ok")
		return
	})

	// 精确定位
	ptzGroup.PUT("/ptzPosition", func(c *gin.Context) {
		var pos dsd.Position

		if c.ShouldBind(&pos) != nil {
			c.JSON(401, gin.H{"status": "bind error"})
		}

		if err := blpInstance.Basic().Goto(&pos); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "goto position ok")
		return
	})

	presetRouter(engine)
	lineRouter(engine)
	cruiseRouter(engine)
	powerRouter(engine)
	idleRouter(engine)
	cronRouter(engine)

	engine.Run(":8089")

}

func presetRouter(engine *gin.Engine) {
	blpInstance := ptzV3.Instance()

	presetGroup := engine.Group("/v1/ptz/preset")

	// 获取所有预置点
	presetGroup.GET("/getpresets", func(c *gin.Context) {
		list := blpInstance.Preset().List()

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      list,
			Detail:    "",
			Message:   "",
			Translate: "",
		})
		return
	})

	// 转至预置点
	presetGroup.PUT("/gotopreset", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		req := &blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.Preset,
			ID:      idNum,
			Speed:   1,
		}

		if err := blpInstance.Manager().Start(req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "goto preset ok")
		return
	})

	// 修改预置点名称
	presetGroup.PUT("/modifypreset", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		name := c.Query("name")

		if err := blpInstance.Preset().Update(dsd.PresetID(idNum), name); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "modify preset ok")
		return
	})

	// 删除预置点
	presetGroup.DELETE("/removepreset", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := blpInstance.Preset().Delete(dsd.PresetID(idNum)); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "delete preset ok")
		return
	})

	// 删除全部预置点
	presetGroup.DELETE("/removepresets", func(c *gin.Context) {
		if err := blpInstance.Preset().DeleteAll(); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "delete preset ok")
		return
	})

	// 设置预置点
	presetGroup.POST("/setpreset", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		name := c.Query("name")

		if err := blpInstance.Preset().Set(dsd.PresetID(idNum), name); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "set preset ok")
		return
	})
}

func lineRouter(engine *gin.Engine) {
	blpInstance := ptzV3.Instance()

	lineGroup := engine.Group("/v1/ptz/linearscan")

	// 获取线性扫描配置
	lineGroup.GET("", func(c *gin.Context) {

		list := blpInstance.Line().List()

		c.JSON(200, Response{
			Code:      200,
			Data:      list,
			Detail:    "",
			Message:   "get linear scan success",
			Translate: "",
		})
		return
	})

	// 获取线性扫描配置
	lineGroup.PUT("", func(c *gin.Context) {

		if err := blpInstance.Line().Default(); err != nil {
			log.Error(err.Error())
			return
		}

		c.JSON(200, Response{
			Code:      200,
			Data:      "",
			Detail:    "",
			Message:   "default linear scan success",
			Translate: "",
		})
		return
	})

	// 设置线性扫描参数
	lineGroup.POST("", func(c *gin.Context) {

		var line dsd.LineScan

		if c.ShouldBind(&line) != nil {
			c.JSON(401, gin.H{"status": "bind error"})
			return
		}

		log.Infof("line: %+v", line)

		if err := line.Validate(); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := blpInstance.Line().Set(&line); err != nil {
			log.Error(err.Error())
			return
		}

		c.JSON(200, Response{
			Code:      200,
			Data:      "",
			Detail:    "",
			Message:   "set linear scan success",
			Translate: "",
		})
		return
	})

	// 开始线扫
	lineGroup.POST("/start", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		req := &blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.LineScan,
			ID:      idNum,
			Speed:   1,
		}

		if err := blpInstance.Manager().Start(req); err != nil {
			log.Error(err.Error())
			c.JSON(400, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(200, "success")
		return
	})

	// 停止线扫
	lineGroup.POST("/stop", func(c *gin.Context) {

		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		req := &blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.LineScan,
			ID:      idNum,
			Speed:   1,
		}

		if err := blpInstance.Manager().Stop(req); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(200, "success")
		return
	})

	// 设置线性扫描左右边界
	lineGroup.POST("/limit", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		limit := c.Query("limit")
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		clear := c.Query("clear")
		clearBool, err := strconv.ParseBool(clear)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		var op dsd.LineMarginOp

		if clearBool && limitNum == 1 {
			op = dsd.ClearLeftMargin
		} else if clearBool && limitNum == 2 {
			op = dsd.ClearRightMargin
		} else if !clearBool && limitNum == 1 {
			op = dsd.SetLeftMargin
		} else if !clearBool && limitNum == 2 {
			op = dsd.SetRightMargin
		} else {
			err := fmt.Errorf("param is invalid. limit: %d, clear: %t", limitNum, clearBool)
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := blpInstance.Line().SetMargin(dsd.LineScanID(idNum), op); err != nil {
			c.JSON(200, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})

		}

		c.JSON(200, "success")
		return
	})
}

func cruiseRouter(engine *gin.Engine) {
	blpInstance := ptzV3.Instance()

	cruiseGroup := engine.Group("/v1/ptz/tour")

	// 获取巡航组配置
	cruiseGroup.GET("/gettours", func(c *gin.Context) {

		cruises := blpInstance.Cruise().List()

		c.JSON(200, Response{
			Code:      200,
			Data:      cruises,
			Detail:    "",
			Message:   "get cruise success",
			Translate: "",
		})
		return
	})

	// 巡航组恢复默认配置
	cruiseGroup.PUT("", func(c *gin.Context) {

		if err := blpInstance.Cruise().Default(); err != nil {
			c.JSON(400, Response{
				Code:      400,
				Data:      "",
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
		}

		c.JSON(200, Response{
			Code:      200,
			Data:      "",
			Detail:    "",
			Message:   "get cruise success",
			Translate: "",
		})
		return
	})

	// 修改巡航组名称
	cruiseGroup.PUT("/modifytour", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		name := c.Query("name")

		if err := blpInstance.Cruise().Update(dsd.CruiseID(idNum), name); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "modify cruise ok")
		return
	})

	// 清除巡航线路
	cruiseGroup.DELETE("/removetour", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := blpInstance.Cruise().Delete(dsd.CruiseID(idNum)); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, "delete cruise ok")
		return
	})

	// 设置巡航线路
	cruiseGroup.POST("/settour", func(c *gin.Context) {
		var cruise dsd.TourPreset

		if c.ShouldBind(&cruise) != nil {
			c.JSON(401, gin.H{"status": "bind error"})
			return
		}

		log.Infof("cruise: %+v", cruise)

		if err := cruise.Validate(); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		if err := blpInstance.Cruise().Set(&cruise); err != nil {
			log.Error(err.Error())
			return
		}

		c.JSON(200, Response{
			Code:      200,
			Data:      "",
			Detail:    "",
			Message:   "set cruise success",
			Translate: "",
		})
		return
	})

	// 开始
	cruiseGroup.PUT("/starttour", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		req := &blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.Cruise,
			ID:      idNum,
			Speed:   1,
		}

		if err := blpInstance.Manager().Start(req); err != nil {
			log.Error(err.Error())
			c.JSON(400, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		c.JSON(200, "success")
		return
	})

	// 停止
	cruiseGroup.PUT("/stoptour", func(c *gin.Context) {
		id := c.Query("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		req := &blp.Request{
			Trigger: blp.ManualTrigger,
			Ability: blp.Cruise,
			ID:      idNum,
			Speed:   1,
		}

		if err := blpInstance.Manager().Stop(req); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
			return
		}

		c.JSON(200, "success")
		return
	})
}

func powerRouter(engine *gin.Engine) {
	blpInstance := ptzV3.Instance()

	powerGroup := engine.Group("/v1/ptz/powerUpAction")

	// 获取开机动作
	powerGroup.GET("", func(c *gin.Context) {

		up := blpInstance.Power().Get()

		c.JSON(200, Response{
			Code:      200,
			Data:      up,
			Detail:    "",
			Message:   "get power up success",
			Translate: "",
		})
		return
	})

	// 设置开机动作
	powerGroup.POST("", func(c *gin.Context) {

		var ups dsd.PowerUps

		if c.ShouldBind(&ups) != nil {
			c.JSON(401, gin.H{"status": "bind error"})
		}

		if err := blpInstance.Power().Set(&ups); err != nil {
			c.JSON(200, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
		}

		c.JSON(200, Response{
			Code:      200,
			Data:      nil,
			Detail:    "",
			Message:   "set power up success",
			Translate: "",
		})
		return
	})

	// 开机动作恢复默认配置
	powerGroup.PUT("/defaultconfig", func(c *gin.Context) {

		if err := blpInstance.Power().Default(); err != nil {
			c.JSON(200, Response{
				Code:      400,
				Data:      nil,
				Detail:    err.Error(),
				Message:   "",
				Translate: "",
			})
		}

		c.JSON(200, Response{
			Code:      200,
			Data:      nil,
			Detail:    "",
			Message:   "default power up success",
			Translate: "",
		})
		return
	})
}

func idleRouter(engine *gin.Engine) {
	blpInstance := ptzV3.Instance()

	idleGroup := engine.Group("/v1/ptz/idlemotion")

	// 获取空闲动作配置
	idleGroup.GET("", func(c *gin.Context) {
		idle := blpInstance.Idle().Get()

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      idle,
			Detail:    "",
			Message:   "",
			Translate: "",
		})
		return
	})

	// 空闲动作恢复默认配置
	idleGroup.PUT("", func(c *gin.Context) {
		err := blpInstance.Idle().Default()
		if err != nil {
			c.JSON(200, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
		}

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      "",
			Detail:    "",
			Message:   "default idle success",
			Translate: "",
		})
		return
	})

	// 空闲动作开始和停止
	idleGroup.POST("", func(c *gin.Context) {
		var motion dsd.IdleMotion

		err := c.ShouldBind(&motion)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		err = blpInstance.Idle().Set(&motion)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      "",
			Detail:    "",
			Message:   "set idle success",
			Translate: "",
		})
		return
	})
}

func cronRouter(engine *gin.Engine) {
	blpInstance := ptzV3.Instance()

	cronGroup := engine.Group("/v1/ptz/autoMovement")

	// 获取定时任务
	cronGroup.GET("", func(c *gin.Context) {
		cron := blpInstance.Cron().List()

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      cron,
			Detail:    "",
			Message:   "",
			Translate: "",
		})
		return
	})

	// 云台定时任务恢复默认配置
	cronGroup.PUT("", func(c *gin.Context) {
		err := blpInstance.Cron().Default()
		if err != nil {
			c.JSON(200, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
		}

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      "",
			Detail:    "",
			Message:   "default cron success",
			Translate: "",
		})
		return
	})

	// 设置定时任务
	cronGroup.POST("", func(c *gin.Context) {
		var movement dsd.PtzAutoMovement

		err := c.ShouldBind(&movement)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		err = blpInstance.Cron().Set(&movement)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Code:      400,
				Data:      nil,
				Detail:    "",
				Message:   err.Error(),
				Translate: "",
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      "",
			Detail:    "",
			Message:   "set cron success",
			Translate: "",
		})
		return
	})
}
