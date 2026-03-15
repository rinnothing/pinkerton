package config_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/rinnothing/pinkerton/config"
)

var desiredConfig = &config.Config{
	HTTP: config.HTTP{
		Host: "somehost",
		Port: "1234",
	},
	HealthCheck: config.HealthCheck{
		Threads: 4,
		Timeout: time.Second * 10,
	},
	InitModels: config.InitModels{
		config.TargetRequest{
			URL:    "url1",
			Period: time.Second * 30,
		},
		config.TargetRequest{
			URL:    "url2",
			Period: time.Second * 5,
		},
	},
	ShutdownTimeout: time.Minute,
}

func TestReadConfig(t *testing.T) {
	t.Parallel()

	defer func() {
		if panicErr := recover(); panicErr != nil {
			t.Fatalf("failed because of panic: %v", panicErr)
		}
	}()

	cfg := config.ReadConfig("config_test.json")
	if !reflect.DeepEqual(cfg, desiredConfig) {
		t.Fatalf("desired and test configs should be equal: %v != %v", desiredConfig, cfg)
	}
}

func TestReadConfigEnv(t *testing.T) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			t.Fatalf("failed because of panic: %v", panicErr)
		}
	}()

	t.Setenv("PORT", "1234")

	cfg := config.ReadConfig("config_env_test.json")
	if !reflect.DeepEqual(cfg, desiredConfig) {
		t.Fatalf("desired and test configs should be equal: %v != %v", desiredConfig, cfg)
	}
}

func TestReadConfigEmptyPort(t *testing.T) {
	t.Parallel()

	defer func() {
		if panicErr := recover(); panicErr == nil {
			t.Fatalf("should fail because of empty port")
		}
	}()

	config.ReadConfig("config_env_test.json")
}

func TestServerConfigCorrect(t *testing.T) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			t.Fatalf("failed because of panic: %v", panicErr)
		}
	}()

	cfg := config.ReadConfig("server.json")
	if cfg == nil {
		t.Fatal("server config shouldn't be nil")
	}
}
