package config

import (
	"os"
)

// AppConfig holds all configuration for the application.
type AppConfig struct {
	PythonVenvPath string
	ServerPort     string
}

// Load returns a new configuration for the application, loading values
// from environment variables with sensible defaults.
func Load() *AppConfig {
	pythonPath := os.Getenv("PYTHON_VENV_PATH")
	if pythonPath == "" {
		pythonPath = "/home/owen/repos/llm-cortex/python_venv/bin/python3" // Default value
	}

	return &AppConfig{
		PythonVenvPath: pythonPath,
		ServerPort:     "8080", // Default port
	}
}