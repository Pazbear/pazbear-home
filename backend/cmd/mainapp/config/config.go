package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Target struct {
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
}

type MsgQueue struct {
	URL    string `json:"url"`
	Target Target `json:"target"`
}

type Config struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	MQ      MsgQueue `json:"mq"`
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
