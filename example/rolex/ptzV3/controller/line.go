package controller

import (
	"encoding/json"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"github.com/tiandh987/CGODemo/example/rolex/utils/apis"
	"github.com/tiandh987/CGODemo/example/rolex/utils/errors"
)

func (c *Controller) ListLine() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	cfg := c.blp.Line().List()

	resp.Success(cfg)
}

func (c *Controller) DefaultLine() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	if err := c.blp.Line().Default(); err != nil {
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) SetLine() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	var line dsd.LineScan
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &line); err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, "unmarshal body failed")
		return
	}

	if err := line.Validate(); err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	if err := c.blp.Line().Set(&line); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) SetLineMargin() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, "get param id failed")
		return
	}

	limit, err := c.GetInt("limit")
	if err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, "get param limit failed")
		return
	}

	clear, err := c.GetBool("clear")
	if err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, "get param clear failed")
		return
	}

	var op dsd.LineMarginOp
	if clear && limit == 1 {
		op = dsd.ClearLeftMargin
	} else if clear && limit == 2 {
		op = dsd.ClearRightMargin
	} else if !clear && limit == 1 {
		op = dsd.SetLeftMargin
	} else if !clear && limit == 2 {
		op = dsd.SetRightMargin
	} else {
		resp.Error(errors.ErrOprFailed, "param (limit/clear) is invalid")
		return
	}

	if err := c.blp.Line().SetMargin(dsd.LineScanID(id), op); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}
