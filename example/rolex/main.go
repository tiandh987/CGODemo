package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	pkgVersion "github.com/tiandh987/CGODemo/example/rolex/pkg/version"
	ptzV2 "github.com/tiandh987/CGODemo/example/rolex/ptzV2"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/blp/ptz"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
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

	if err := ptzV2.Start(); err != nil {
		panic(err)
	}

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
		version := blp.Instance().Version()

		c.JSON(http.StatusOK, version)
	})

	// 查询云台型号
	ptzGroup.GET("/ptzmodel", func(c *gin.Context) {
		model := blp.Instance().Model()

		c.JSON(http.StatusOK, model)
	})

	// 查询云台状态
	ptzGroup.GET("/status", func(c *gin.Context) {
		state := blp.Instance().State()

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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		speed := c.Query("speed")
		speedNum, err := strconv.Atoi(speed)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := ptz.Operation(dirNum).ValidateDirection(); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := ptz.Speed(speedNum).Validate(); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().Control(ptz.Manual, ptz.ManualFunc, dirNum, 0, ptz.Speed(speedNum)); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		time.Sleep(time.Millisecond * 200)
		//time.Sleep(time.Second * 5)

		if err := blp.Instance().Control(ptz.Manual, ptz.None, 0, 0, 0); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, "turn ok")
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
				Message:   "",
				Translate: "",
			})
			return
		}

		speed := c.Query("speed")
		speedNum, err := strconv.Atoi(speed)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := ptz.Operation(dirNum).ValidateDirection(); err != nil {
			log.Error(err.Error())

			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := ptz.Speed(speedNum).Validate(); err != nil {
			log.Error(err.Error())

			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().Control(ptz.Manual, ptz.ManualFunc, dirNum, 0, ptz.Speed(speedNum)); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, err)
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := ptz.Operation(methodNum).ValidateOperation(); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().Control(ptz.Manual, ptz.ManualFunc, methodNum, 0, 0); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, "operation ok")
		return
	})

	// 云台停止转动
	ptzGroup.PUT("/stop", func(c *gin.Context) {
		if err := blp.Instance().Control(ptz.Manual, ptz.None, 0, 0, 0); err != nil {
			c.JSON(http.StatusBadRequest, err)
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

		if err := blp.Instance().Position(&pos); err != nil {
			c.JSON(http.StatusBadRequest, err)
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

	engine.Run(":8089")

}

func presetRouter(engine *gin.Engine) {
	presetGroup := engine.Group("/v1/ptz/preset")

	// 获取所有预置点
	presetGroup.GET("/getpresets", func(c *gin.Context) {
		list := blp.Instance().ListPreset()

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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().Control(ptz.Manual, ptz.Preset, idNum, 0, 0); err != nil {
			c.JSON(http.StatusBadRequest, err)
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		name := c.Query("name")

		if err := blp.Instance().UpdatePreset(dsd.PresetID(idNum), dsd.PresetName(name)); err != nil {
			c.JSON(http.StatusBadRequest, err)
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().DeletePreset(dsd.PresetID(idNum)); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, "delete preset ok")
		return
	})

	// 删除全部预置点
	presetGroup.DELETE("/removepresets", func(c *gin.Context) {
		if err := blp.Instance().DeleteAllPreset(); err != nil {
			c.JSON(http.StatusBadRequest, err)
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		name := c.Query("name")

		if err := blp.Instance().SetPreset(dsd.PresetID(idNum), dsd.PresetName(name)); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, "set preset ok")
		return
	})
}

func lineRouter(engine *gin.Engine) {
	lineGroup := engine.Group("/v1/ptz/linearscan")

	// 获取线性扫描配置
	lineGroup.GET("", func(c *gin.Context) {

		list := blp.Instance().ListLine()

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

		if err := blp.Instance().DefaultLine(); err != nil {
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
			c.JSON(400, gin.H{"status": "bad param"})
			return
		}

		if err := blp.Instance().SetLine(&line); err != nil {
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().Control(ptz.Manual, ptz.LineScan, idNum, 0, 0); err != nil {
			log.Error(err.Error())
			c.JSON(500, "line start failed")
			return
		}

		c.JSON(200, "success")
		return
	})

	// 停止线扫
	lineGroup.POST("/stop", func(c *gin.Context) {
		if err := blp.Instance().Control(ptz.Manual, ptz.None, 0, 0, 0); err != nil {
			log.Error(err.Error())
			c.JSON(500, "line start failed")
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		limit := c.Query("limit")
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, err)
			return
		}

		clear := c.Query("clear")
		clearBool, err := strconv.ParseBool(clear)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().SetLineMargin(dsd.LineScanID(idNum), limitNum, clearBool); err != nil {
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
	cruiseGroup := engine.Group("/v1/ptz/tour")

	// 获取巡航组配置
	cruiseGroup.GET("/gettours", func(c *gin.Context) {

		cruises := blp.Instance().ListCruise()

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

		if err := blp.Instance().DefaultCruise(); err != nil {
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		name := c.Query("name")

		if err := blp.Instance().UpdateCruise(dsd.CruiseID(idNum), dsd.CruiseName(name)); err != nil {
			c.JSON(http.StatusBadRequest, err)
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().DeleteCruise(dsd.CruiseID(idNum)); err != nil {
			c.JSON(http.StatusBadRequest, err)
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

		//if err := cruise.Validate(); err != nil {
		//	log.Error(err.Error())
		//	c.JSON(400, gin.H{"status": "bad param"})
		//	return
		//}

		if err := blp.Instance().SetCruise(&cruise); err != nil {
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
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if err := blp.Instance().Control(ptz.Manual, ptz.Cruise, idNum, 0, 0); err != nil {
			log.Error(err.Error())
			c.JSON(500, "cruise start failed")
			return
		}

		c.JSON(200, "success")
		return
	})

	// 停止
	cruiseGroup.PUT("/stoptour", func(c *gin.Context) {
		if err := blp.Instance().Control(ptz.Manual, ptz.None, 0, 0, 0); err != nil {
			log.Error(err.Error())
			c.JSON(500, "cruise start failed")
			return
		}

		c.JSON(200, "success")
		return
	})
}

func powerRouter(engine *gin.Engine) {
	powerGroup := engine.Group("/v1/ptz/powerUpAction")

	// 获取开机动作
	powerGroup.GET("", func(c *gin.Context) {

		up, err := blp.Instance().GetPowerUp()
		if err != nil {
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

		if err := blp.Instance().SetPowerUp(&ups); err != nil {
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

		if err := blp.Instance().DefaultPowerUp(); err != nil {
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
	idleGroup := engine.Group("/v1/ptz/idlemotion")

	// 获取空闲动作配置
	idleGroup.GET("", func(c *gin.Context) {
		idle, err := blp.Instance().GetIdle()
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
			Data:      idle,
			Detail:    "",
			Message:   "",
			Translate: "",
		})
		return
	})

	// 空闲动作恢复默认配置
	idleGroup.PUT("", func(c *gin.Context) {
		err := blp.Instance().DefaultIdle()
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

		err = blp.Instance().SetIdle(&motion)
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
