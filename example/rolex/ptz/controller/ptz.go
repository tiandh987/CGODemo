package controller

import (
	"rolex/utils/apis"
	"rolex/utils/errors"
)

func (c *Controller) Version() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	version, err := c.blp.Ptz.Version()
	if err != nil {
		resp.Error(errors.ErrOprFailed, "get ptz version failed")
		return
	}

	resp.Success(version)
}

func (c *Controller) Model() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	model, err := c.blp.Ptz.Model()
	if err != nil {
		resp.Error(errors.ErrOprFailed, "get ptz model failed")
		return
	}

	resp.Success(model)
}
