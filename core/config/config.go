package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ModelConfig defines the configuration for a single model.
type ModelConfig struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"` // e.g., "blip", "clip", "gguf"
	Path   string `yaml:"path"`
	Device string `yaml:"device"` // e.g., "cpu", "cuda"
}

// AppConfig holds all configuration for the application.
type AppConfig struct {
	PythonVenvPath string        `yaml:"python_venv_path"`
	ServerPort     string        `yaml:"server_port"`
	Models         []ModelConfig `yaml:"models"`
}

// Load returns a new configuration for the application, loading values
// from a YAML file and allowing overrides from environment variables.
func Load(configPath string) (*AppConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Override with environment variables if they exist
	if pythonPath := os.Getenv("PYTHON_VENV_PATH"); pythonPath != "" {
		cfg.PythonVenvPath = pythonPath
	}
	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		cfg.ServerPort = serverPort
	}

	return &cfg, nil
}