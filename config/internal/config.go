package config_internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Config struct {
		HTTP            `json:"http"`
		HealthCheck     `json:"healthcheck"`
		InitModels      `json:"init_models"`
		ShutdownTimeout string `json:"shutdown_timeout"`
		Debug           bool   `json:"debug"`
	}

	HTTP struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}

	HealthCheck struct {
		Threads int    `json:"threads"`
		Timeout string `json:"timeout"`
	}

	TargetRequest struct {
		URL    string `json:"url"`
		Period string `json:"period"`
	}

	InitModels []TargetRequest
)

func ReadConfig(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("can't open configuration file %s: %s", path, err))
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var cfg Config
	err = dec.Decode(&cfg)
	if err != nil {
		panic(fmt.Sprintf("can't decode configration file %s: %s", path, err))
	}

	return &cfg
}
