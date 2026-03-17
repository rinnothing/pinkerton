package main

import (
	"context"
	"log/slog"

	"github.com/rinnothing/pinkerton/config"
	"github.com/rinnothing/pinkerton/internal/controller"
)

func main() {
	cfg := config.ReadConfig("config/server.json")
	if cfg.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	controller.Run(context.Background(), cfg)
}
