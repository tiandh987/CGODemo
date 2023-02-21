package controller

import (
	"encoding/json"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"github.com/tiandh987/CGODemo/example/rolex/utils/apis"
	"github.com/tiandh987/CGODemo/example/rolex/utils/errors"
)

func (c *Controller) GetIdle() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	cfg := c.blp.Idle().Get()

	resp.Success(cfg)
}

func (c *Controller) SetIdle() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	var motion dsd.IdleMotion
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &motion); err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, "unmarshal body failed")
		return
	}

	if err := motion.Validate(); err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	if err := c.blp.Idle().Set(&motion); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) DefaultIdle() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	if err := c.blp.Idle().Default(); err != nil {
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	resp.Success()
}
