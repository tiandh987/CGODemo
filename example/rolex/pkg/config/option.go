package config

import "fmt"

type Option struct {
	Path string
	Name string
	Type string
}

func NewOption() *Option {
	return &Option{
		Path: "/mnt/config",
		Name: "rolex",
		Type: "yaml",
	}
}

func (o *Option) String() string {
	return fmt.Sprintf("%s/%s.%s", o.Path, o.Name, o.Type)
}
