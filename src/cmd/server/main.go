package main

import (
	"log"

	"ne-project/src/internal/app"
	"ne-project/src/internal/config/appconfig"
)

// @title           GreenMart API
// @version         1.0.1.02072026
// @description     API Server for GreenMart Marketplace.
// @BasePath        /api/v1
func main() {
	// Load Configuration
	cfg, err := appconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run Application
	app.Run(cfg)
}
