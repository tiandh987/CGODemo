package config

import (
	"errors"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/config"

	"sync"
)

var (
	callbackMu  sync.RWMutex
	callbackMap = make(map[string]map[string]Callback)
)

var (
	curConfig *config.Config
	defConfig *config.Config
)

var (
	DftConfigPath = "/mnt/config/default.json"
	CurConfigPath = "/mnt/config/current.json"
)

func initDefaultConfig() {
	opt := &config.Option{
		Path: "/mnt/custom/tian/rolex_nb/config",
		Name: "default",
		Type: "yaml",
	}

	defConfig = config.New(opt)
}

func SetDefault(key string, value interface{}) (err error) {
	if err := defConfig.Set(key, value); err != nil {
		return err
	}

	// 判断当前配置（current.json）是否为空，若为空则设置为默认值
	if !curConfig.IsSet(key) {
		if err := curConfig.Set(key, value); err != nil {
			return err
		}
	}

	return nil
}

func GetDefault(key string, value interface{}) (err error) {
	return defConfig.UnmarshalKey(key, value)
}

func initCurrentConfig() {
	opt := &config.Option{
		Path: "/mnt/custom/tian/rolex_nb/config",
		Name: "current",
		Type: "yaml",
	}

	curConfig = config.New(opt)
}

func SetConfig(key string, value interface{}) (err error) {
	// 执行配置修改回调函数
	callbackMu.RLock()

	cbMap, ok := callbackMap[key]
	if ok {
		for _, callback := range cbMap {
			if !callback.OnApplyConfig(value) {
				return errors.New("exec config callback failed")
			}
		}
	}

	callbackMu.RUnlock()

	return curConfig.Set(key, value)
}

func GetConfig(key string, value interface{}) (err error) {
	return curConfig.UnmarshalKey(key, value)
}

// Callback 配置回调接口 (current.json)
type Callback interface {
	OnApplyConfig(val interface{}) bool
}

// Attach 注册配置回调函数
func Attach(cfgName string, cbName string, callback Callback) bool {
	callbackMu.Lock()
	defer callbackMu.Unlock()

	cbMap, ok := callbackMap[cfgName]
	if !ok {
		cbMap = make(map[string]Callback)
		callbackMap[cfgName] = cbMap
	}
	cbMap[cbName] = callback
	return true
}

// Detach 注销配置回调函数
func Detach(cfgName string, cbName string) bool {
	callbackMu.Lock()
	defer callbackMu.Unlock()

	if _, ok := callbackMap[cfgName]; !ok {
		return false
	}

	cbMap, ok := callbackMap[cfgName]
	if !ok {
		return false
	}

	delete(cbMap, cbName)

	return true
}
