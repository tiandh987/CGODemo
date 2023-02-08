package controller

import (
	"encoding/json"
	"rolex/pkg/log"
	"rolex/ptz/dsd"
	"rolex/utils/apis"
	"rolex/utils/errors"
)

// GetSerialConfig 获取云台协议参数（串口设置）
func (c *Controller) GetSerialConfig() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	def, err := c.GetBool("default", false)
	if err != nil {
		log.Errorf("get param default failed, err: %s", err.Error())
		resp.Error(errors.ErrGetCfgFailed, "get default failed")
		return
	}

	param, err := c.blp.Serial.Get(def)
	if err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success(param)
}

// SetSerialConfig 设置云台协议参数（串口设置）
func (c *Controller) SetSerialConfig() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	var cfg dsd.PTZ
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &cfg); err != nil {
		log.Errorf("unmarshal serial config failed, err: %s", err.Error())
		resp.Error(errors.ErrSetCfgFailed, "unmarshal serial config failed")
		return
	}

	if err := c.blp.Serial.Set(&cfg); err != nil {
		resp.Error(errors.ErrSetCfgFailed, err.Error())
		return
	}

	resp.Success()
}
