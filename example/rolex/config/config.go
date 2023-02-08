package config

import "github.com/tiandh987/CGODemo/example/rolex/pkg/log"

func init() {
	log.Info("start init config")

	initDefaultConfig()
	initCurrentConfig()

	log.Info("init config success")
}
