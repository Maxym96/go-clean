package main

import (
	"go-clean/config"
	"go-clean/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		return
	}

	// Run
	app.Run(cfg)
}
