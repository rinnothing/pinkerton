package config

import (
	"fmt"
	"time"

	config_internal "github.com/rinnothing/pinkerton/config/internal"
	"github.com/rinnothing/pinkerton/internal/model"
)

type (
	Config struct {
		HTTP
		HealthCheck
		InitModels
		ShutdownTimeout time.Duration
	}

	HTTP struct {
		Host string
		Port string
	}

	HealthCheck struct {
		Threads int
		Timeout time.Duration
	}

	TargetRequest = model.TargetRequest

	InitModels []TargetRequest
)

func parseDuration(s string) time.Duration {
	tm, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("duration string %s is in wrong format: %s", s, err))
	}

	return tm
}

func parseModel(req config_internal.TargetRequest) TargetRequest {
	return TargetRequest{
		URL:    req.URL,
		Period: parseDuration(req.Period),
	}
}

func ReadConfig(path string) *Config {
	cfg_int := config_internal.ReadConfig(path)

	models := make([]TargetRequest, len(cfg_int.InitModels))
	for i, mdl := range cfg_int.InitModels {
		models[i] = parseModel(mdl)
	}

	return &Config{
		HTTP: HTTP{
			Host: cfg_int.Host,
			Port: cfg_int.Port,
		},
		HealthCheck: HealthCheck{
			Threads: cfg_int.Threads,
			Timeout: parseDuration(cfg_int.Timeout),
		},
		InitModels:      models,
		ShutdownTimeout: parseDuration(cfg_int.ShutdownTimeout),
	}
}
