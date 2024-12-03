package configs

import (
	"os"
	"strconv"
)

type Config struct {
	ConfigMap map[string]string
	Env       string
	Port      int
	Version   string
}

func NewConfig(configMap map[string]string) Config {
	var port int
	portStr := os.Getenv("PORT")
	if os.Getenv("PORT") != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			panic(err)
		}
		port = p
	} else {
		port = 9888
	}
	c := Config{
		Env:     os.Getenv("env"),
		Version: os.Getenv("version"),
		Port:    port,
	}

	switch c.Env {
	case "stg":
		{
			c.ConfigMap = configMap
		}
	default:
		{
			c.ConfigMap = configMap
		}
	}

	return c
}
