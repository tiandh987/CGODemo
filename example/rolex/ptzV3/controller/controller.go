package controller

import (
	"github.com/astaxie/beego"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/blp"
)

type Controller struct {
	beego.Controller
	blp *blp.Blp
}

func (c *Controller) Prepare() {
	//if c.blp != nil {
	//	return
	//}
	//
	//c.blp = blp.Instance()
}
