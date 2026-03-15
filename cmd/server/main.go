package main

import (
	"github.com/rinnothing/pinkerton/config"
	"github.com/rinnothing/pinkerton/internal/controller"
)

func main() {
	cfg := config.ReadConfig("config/server.json")

	controller.Run(cfg)
}
