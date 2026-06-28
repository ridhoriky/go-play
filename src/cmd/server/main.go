package main

import (
	"log"

	"ne-project/src/internal/app"
	"ne-project/src/internal/config/appconfig"
)

func main() {
	// Load Configuration
	cfg, err := appconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run Application
	app.Run(cfg)
}
