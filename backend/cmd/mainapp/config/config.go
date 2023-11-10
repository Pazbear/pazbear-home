package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (c *Config) Init() {
}

const (
	cnfPath = "./config/config.conf"
)

func AppConfig() (Config, error) {
	var config Config
	var f *os.File
	var err error
	f, err = os.Open(cnfPath)
	if err != nil {
		return Config{}, err
	}
	if f == nil {
		return Config{}, fmt.Errorf("No config file")
	}
	cnfbytes, err := io.ReadAll(f)
	if err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(cnfbytes, &config); err != nil {
		return Config{}, err
	}
	config.Init()
	return config, nil
}
