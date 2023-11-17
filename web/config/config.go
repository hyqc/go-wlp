package config

import (
	"github.com/go-micro/plugins/v4/auth/jwt/token"
	"web/config/server"
	"web/pkg/core"
)

type Config struct {
	Server core.Configure `yaml:"server"`
}

var (
	App = NewConfig()
	JWT token.Provider
)

func NewConfig() *Config {
	App = &Config{
		Server: server.Config,
	}
	return App
}
