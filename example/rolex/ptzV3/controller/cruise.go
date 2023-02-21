package controller

import (
	"encoding/json"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"github.com/tiandh987/CGODemo/example/rolex/utils/apis"
	"github.com/tiandh987/CGODemo/example/rolex/utils/errors"
)

func (c *Controller) ListCruise() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	cfg := c.blp.Cruise().List()

	resp.Success(cfg)
}

func (c *Controller) DefaultCruise() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	if err := c.blp.Cruise().Default(); err != nil {
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) UpdateCruiseName() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrGetCfgFailed, "get param id failed")
		return
	}

	name := c.GetString("name")
	if err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrGetCfgFailed, "get param name failed")
		return
	}

	if err := c.blp.Cruise().Update(dsd.CruiseID(id), name); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) DeleteCruise() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrGetCfgFailed, "get param id failed")
		return
	}

	if err := dsd.CruiseID(id).Validate(); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	if err := c.blp.Cruise().Delete(dsd.CruiseID(id)); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) SetCruise() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	var cruise dsd.TourPreset
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &cruise); err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, "unmarshal body failed")
		return
	}

	if err := cruise.Validate(); err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	if err := c.blp.Cruise().Set(&cruise); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}
