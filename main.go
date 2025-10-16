package main

import (
	"log"

	"github.com/owen-6936/llm-cortex/core/config"
	"github.com/owen-6936/llm-cortex/core/engine"
	"github.com/owen-6936/llm-cortex/utils"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load("config.yaml")
	utils.HandleError(err, "Failed to load configuration", true)
	// 2. Create a new engine instance
	appEngine, err := engine.New(cfg)
	utils.HandleError(err, "Failed to create engine", true)
	// 3. Start the engine (which starts the server)
	if err := appEngine.Start(); err != nil {
		log.Fatalf("Failed to start engine: %v", err)
	}
}
