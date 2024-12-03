package configs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

func GetConfigFromEnv() (map[string]string, error) {
	var configMap map[string]string
	configStr := os.Getenv("config")
	decoded, err := base64.StdEncoding.DecodeString(configStr)
	if err != nil {
		fmt.Println("[Parse config] Convert B64 config string error: " + err.Error())
		return nil, err
	}
	err = json.Unmarshal(decoded, &configMap)
	if err != nil {
		fmt.Println("[Parse config] Parse JSON with config string error: " + err.Error())
		return nil, err
	}
	return configMap, err

}
