package config

import (
	"web/config/server"
	"web/pkg/core"
)

type Config struct {
	Server core.Configure `yaml:"server"`
}

var (
	App = NewConfig()
)

func NewConfig() *Config {
	App = &Config{
		Server: server.Config,
	}
	return App
}
