package main

import (
	"encoding/json"
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
	Code      int64  `json:"Code"`      // 200
	Data      string `json:"Data"`      // 返回数据
	Detail    string `json:"Detail"`    // 详细错误信息
	Message   string `json:"Message"`   // msg: success
	Translate string `json:"Translate"` // 操作成功
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
		// TODO

		c.JSON(http.StatusOK, "status")
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

	engine.Run(":8089")

}

func presetRouter(engine *gin.Engine) {
	presetGroup := engine.Group("/v1/ptz/preset")

	// 获取所有预置点
	presetGroup.GET("/getpresets", func(c *gin.Context) {
		list, err := blp.Instance().Preset.List()
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		listStr, _ := json.Marshal(list)
		c.JSON(http.StatusOK, Response{
			Code:      200,
			Data:      string(listStr),
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

		if err := blp.Instance().Preset.Update(dsd.PresetID(idNum), dsd.PresetName(name)); err != nil {
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

		if err := blp.Instance().Preset.Delete(dsd.PresetID(idNum)); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, "delete preset ok")
		return
	})

	// 删除全部预置点
	presetGroup.DELETE("/removepresets", func(c *gin.Context) {
		if err := blp.Instance().Preset.DeleteAll(); err != nil {
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

		if err := blp.Instance().Preset.Set(blp.Instance().GetControl(), dsd.PresetID(idNum), dsd.PresetName(name)); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, "set preset ok")
		return
	})
}
