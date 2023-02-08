package controller

import (
	"github.com/astaxie/beego"
	"rolex/ptz/blp"
)

type Controller struct {
	beego.Controller
	blp *blp.Blp
}

func (c *Controller) Prepare() {
	if c.blp != nil {
		return
	}

	//c.blp = service.Service()
}
