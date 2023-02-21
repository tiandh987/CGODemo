package controller

import (
	"encoding/json"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"github.com/tiandh987/CGODemo/example/rolex/utils/apis"
	"github.com/tiandh987/CGODemo/example/rolex/utils/errors"
)

func (c *Controller) GetPower() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	cfg := c.blp.Power().Get()

	resp.Success(cfg)
}

func (c *Controller) SetPower() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	var ups dsd.PowerUps
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ups); err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, "unmarshal body failed")
		return
	}

	if err := ups.Validate(); err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	if err := c.blp.Power().Set(&ups); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) DefaultPower() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	if err := c.blp.Power().Default(); err != nil {
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	resp.Success()
}
