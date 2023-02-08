package config

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/file"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"sync"
)

// 为兼容老版本，为每个配置文件封装 GetXXX、SetXXX 方法，后续再有新的配置文件建议直接使用 Viper

type Config struct {
	mu    sync.RWMutex
	viper *viper.Viper
}

func New(opt *Option) *Config {
	if opt == nil {
		opt = NewOption()
	}

	log.Debugf("new Config, option: %s", opt.String())

	v := viper.New()
	v.AddConfigPath(opt.Path)
	v.SetConfigName(opt.Name)
	v.SetConfigType(opt.Type)

	exists, err := file.PathExists(opt.String())
	if err != nil {
		log.Panicf("check file failed, file: %s, err: %s", opt.String(), err.Error())
	}

	if !exists {
		if err := v.SafeWriteConfigAs(opt.String()); err != nil {
			log.Panicf("write config failed, file: %s, err: %s", opt.String(), err.Error())
		}
	}

	if err := v.ReadInConfig(); err != nil {
		log.Panicf("readInConfig failed, file: %s, err: %s", v.ConfigFileUsed(), err.Error())
	}

	return &Config{
		viper: v,
	}
}

func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if isSet := c.viper.IsSet(key); !isSet {
		log.Errorf("key(%s) is not set", key)
		return errors.New("key is not set")
	}

	if err := c.viper.UnmarshalKey(key, rawVal); err != nil {
		log.Errorf("UnmarshalKey(%s) failed, err: %s", key, err.Error())
		return err
	}

	//log.Debugf("UnmarshalKey file: %s, key: %s, value: %+v\n", c.viper.ConfigFileUsed(), key, rawVal)

	return nil
}

func (c *Config) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	//log.Debugf("set file: %s, key: %s, value: %+v\n", c.viper.ConfigFileUsed(), key, value)

	c.viper.Set(key, value)

	if err := c.viper.WriteConfig(); err != nil {
		log.Errorf("WriteConfig failed, err: %s", err.Error())
		return errors.New("write config failed")
	}

	return nil
}

func (c *Config) IsSet(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.viper.IsSet(key)
}
