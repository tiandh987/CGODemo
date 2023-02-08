package controller

//func (c *Controller) ListPreset() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	cfg, err := c.srv.ListPreset(false)
//	if err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success(cfg)
//}
//
//func (c *PtzController) GotoPreset() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	id, err := c.GetInt("id", 0)
//	if err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	if err := c.srv.GotoPreset(id); err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success()
//}
//
//func (c *PtzController) SetPresetName() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	id, err := c.GetInt("id", 0)
//	if err != nil {
//		infra.Error(err)
//		resp.Error(errors.ErrGetCfgFailed, "get id failed")
//		return
//	}
//
//	name := c.GetString("name", "")
//	if len(strings.TrimSpace(name)) == 0 {
//		resp.Error(errors.ErrGetCfgFailed, "name is empty")
//		return
//	}
//
//	if err := c.srv.SetPresetName(id, name); err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success()
//}
//
//func (c *PtzController) GotoPresetOk() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	id, err := c.GetInt("id", 0)
//	if err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	ok, err := c.srv.GotoPresetOk(id)
//	if err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success(ok)
//}
//
//func (c *PtzController) RemovePreset() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	id, err := c.GetInt("id", 0)
//	if err != nil {
//		infra.Error(err)
//		resp.Error(errors.ErrGetCfgFailed, "get id failed")
//		return
//	}
//
//	if err := c.srv.RemovePreset(id); err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success()
//}
//
//func (c *PtzController) RemoveAllPreset() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	if err := c.srv.RemoveAllPreset(); err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success()
//}
//
//func (c *PtzController) SetPreset() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	id, err := c.GetInt("id")
//	if err != nil {
//		infra.Error(err)
//		resp.Error(errors.ErrOprFailed, "get param id failed")
//		return
//	}
//
//	name := c.GetString("name")
//	if len(name) == 0 {
//		resp.Error(errors.ErrOprFailed, "param name invalid")
//		return
//	}
//
//	if err := c.srv.SetPresetPoint(id, name); err != nil {
//		resp.Error(errors.ErrOprFailed, err.Error())
//		return
//	}
//
//	resp.Success()
//}
//
//func (c *PtzController) SetPresetSpeed() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	speed, err := c.GetInt("speed", 0)
//	if err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	if err := c.srv.SetPresetSpeed(speed); err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success()
//}
//
//func (c *PtzController) GetPresetId() {
//	infra.Info("enter GetPresetId\n")
//
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	id, err := c.srv.GetPresetId()
//	if err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success(id)
//}
