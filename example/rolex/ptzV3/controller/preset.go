package controller

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"github.com/tiandh987/CGODemo/example/rolex/utils/apis"
	"github.com/tiandh987/CGODemo/example/rolex/utils/errors"
)

func (c *Controller) ListPreset() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	cfg := c.blp.Preset().List()

	resp.Success(cfg)
}

//	func (c *Controller) GotoPreset() {
//		resp := apis.Response{}
//		c.Data["json"] = &resp
//		defer c.ServeJSON()
//
//		id, err := c.GetInt("id")
//		if err != nil {
//			log.Error(err.Error())
//			resp.Error(errors.ErrGetCfgFailed, "get param id failed")
//			return
//		}
//
//		if err := dsd.PresetID(id).Validate(); err != nil {
//			resp.Error(errors.ErrGetCfgFailed, err.Error())
//			return
//		}
//
//		if err := c.blp.Preset.Goto(dsd.PresetID(id)); err != nil {
//			resp.Error(errors.ErrGetCfgFailed, err.Error())
//			return
//		}
//
//		resp.Success()
//	}
func (c *Controller) UpdatePresetName() {
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

	if err := c.blp.Preset().Update(dsd.PresetID(id), name); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}

//	func (c *Controller) GotoPresetOk() {
//		resp := apis.Response{}
//		c.Data["json"] = &resp
//		defer c.ServeJSON()
//
//		id, err := c.GetInt("id")
//		if err != nil {
//			log.Error(err.Error())
//			resp.Error(errors.ErrGetCfgFailed, "get param id failed")
//			return
//		}
//
//		if err := dsd.PresetID(id).Validate(); err != nil {
//			resp.Error(errors.ErrGetCfgFailed, err.Error())
//			return
//		}
//
//		ok, err := c.blp.Preset.GotoOk(dsd.PresetID(id))
//		if err != nil {
//			resp.Error(errors.ErrGetCfgFailed, err.Error())
//			return
//		}
//
//		resp.Success(ok)
//	}
func (c *Controller) DeletePreset() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		resp.Error(errors.ErrGetCfgFailed, "get param id failed")
		return
	}

	if err := dsd.PresetID(id).Validate(); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	if err := c.blp.Preset().Delete(dsd.PresetID(id)); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) DeleteAllPreset() {
	resp := apis.Response{}
	c.Data["json"] = &resp
	defer c.ServeJSON()

	if err := c.blp.Preset().DeleteAll(); err != nil {
		resp.Error(errors.ErrGetCfgFailed, err.Error())
		return
	}

	resp.Success()
}

func (c *Controller) SetPreset() {
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

	if err := c.blp.Preset().Set(dsd.PresetID(id), name); err != nil {
		resp.Error(errors.ErrOprFailed, err.Error())
		return
	}

	resp.Success()
}

//
//func (c *Controller) SetPresetSpeed() {
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
//	if err := state.Speed(speed).Validate(); err != nil {
//		resp.Error(errors.ErrOprFailed, err.Error())
//		return
//	}
//
//	if err := c.blp.Preset.Speed(state.Speed(speed)); err != nil {
//		resp.Error(errors.ErrGetCfgFailed, err.Error())
//		return
//	}
//
//	resp.Success()
//}
//
//func (c *Controller) GetPresetId() {
//	resp := apis.Response{}
//	c.Data["json"] = &resp
//	defer c.ServeJSON()
//
//	id := c.blp.Preset.PresetID()
//
//	resp.Success(id)
//}
