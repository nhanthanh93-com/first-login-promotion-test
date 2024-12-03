package app

import (
	"time"
	"trinity/configs"
)

type Config struct {
	Config  configs.Config
	Timeout time.Duration
}

func NewAppConfig(c configs.Config) Config {
	return Config{
		Config:  c,
		Timeout: 10 * time.Second,
	}
}
