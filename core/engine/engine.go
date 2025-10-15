package engine

import (
	"github.com/owen-6936/llm-cortex/core/config"
)

// Engine is the central orchestrator for the application.
// It manages the lifecycle of models and other core services.
type Engine struct {
	config *config.AppConfig
	// You could add model managers or other services here.
}

// New creates a new application engine.
func New(cfg *config.AppConfig) (*Engine, error) {
	return &Engine{
		config: cfg,
	}, nil
}

// Start initializes and starts all the core services.
func (e *Engine) Start() error {
	// Example: Initialize model managers here.
	// vision.Initialize(e.config.PythonVenvPath)
	return nil
}