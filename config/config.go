package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rinnothing/pinkerton/internal/model"
)

type (
	Config struct {
		HTTP            `json:"http"`
		HealthCheck     `json:"healthcheck"`
		InitModels      `json:"init_models"`
		ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	}

	HTTP struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}

	HealthCheck struct {
		Threads int           `json:"threads"`
		Timeout time.Duration `json:"timeout"`
	}

	InitModels []model.TargetRequest
)

func ReadConfig(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("can't open configuration file %s: %s", path, err))
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var cfg *Config
	err = dec.Decode(cfg)
	if err != nil {
		panic(fmt.Sprintf("can't decode configration file %s: %s", path, err))
	}

	return cfg
}
