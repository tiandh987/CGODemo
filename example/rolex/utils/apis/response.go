package apis

import (
	"rolex/conf"
	"rolex/config"
	"rolex/infra"
	"rolex/utils/errors"
)

type Response struct {
	Code      int         `description: "返回码"`
	Message   string      `description: "返回码描述"`
	Translate string      `description: "返回码提示翻译"`
	Detail    string      `description: "详细错误信息，不展示，仅供调试时使用"`
	Data      interface{} `description: "返回数据"`
}

func init() {
	//infra.Info("init response enter")
	var (
		devLanguage conf.DevLanguage
	)
	err := config.GetPd("dev_language", &devLanguage)
	if err != nil {
		infra.Error("get dev_language value from pd config faild. err:", err)
	}

	if devLanguage.Language == "chinese" {
		errors.CodeMap = errors.CodeMapCN
	}

	if devLanguage.Language == "english" {
		errors.CodeMap = errors.CodeMapEN
	}

	for code, val := range errors.CodeMap {
		errors.CodeMsgMap[val.Message] = errors.CodeKey{code, val.Translate}
	}
	infra.Info("init response after")
}

// 成功返回
func (this *Response) Success(data ...interface{}) {
	this.Code = int(errors.Success)
	this.Message = errors.Success.GetMessage()
	this.Translate = errors.Success.GetTranslate()
	this.Detail = ""
	if len(data) > 0 {
		this.Data = data[0]
	}
}

// 成功返回但需要重启设备
func (this *Response) SuccessNeedReboot(data ...interface{}) {
	this.Code = int(errors.SuccessNeedReboot)
	this.Message = errors.SuccessNeedReboot.GetMessage()
	this.Translate = errors.SuccessNeedReboot.GetTranslate()
	this.Detail = ""
	if len(data) > 0 {
		this.Data = data[0]
	}
}

// 错误返回
func (this *Response) Error(err error, detail ...string) {
	str := err.Error()
	value, ok := errors.CodeMsgMap[str]
	if ok {
		this.Code = int(value.Code)
		this.Message = str
		this.Translate = value.Translate
		if len(detail) > 0 {
			this.Detail = detail[0]
		}
		this.Data = nil
	} else {
		this.Code = int(errors.ErrSystem)
		this.Message = errors.ErrSystem.GetMessage()
		this.Translate = errors.ErrSystem.GetTranslate()
		this.Detail = str
		this.Data = nil
	}
}
