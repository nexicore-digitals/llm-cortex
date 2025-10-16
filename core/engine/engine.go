package engine

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/owen-6936/llm-cortex/core/config"
	"github.com/owen-6936/llm-cortex/core/models/vision"
	"github.com/owen-6936/llm-cortex/handlers"
)

// Engine is the central orchestrator for the application.
// It manages the lifecycle of models and other core services.
type Engine struct {
	config       *config.AppConfig
	modelPlugins map[string]interface{} // A simple map to hold loaded model instances
}

// New creates a new application engine.
func New(cfg *config.AppConfig) (*Engine, error) {
	return &Engine{
		config:       cfg,
		modelPlugins: make(map[string]interface{}),
	}, nil
}

// Start initializes and starts all the core services.
func (e *Engine) Start() error {
	// 1. Initialize models based on configuration
	if err := e.initializeModels(); err != nil {
		return fmt.Errorf("failed to initialize models: %w", err)
	}

	// 2. Set up HTTP server and handlers
	mux := http.NewServeMux()

	// Serve UI
	os := http.FileServer(http.Dir("ui"))
	mux.Handle("/", os)

	// Shell handlers
	mux.HandleFunc("/shell/start", handlers.StartShellHandler)
	mux.HandleFunc("/shell/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/send"):
			handlers.SendCommandHandler(w, r)
		case strings.HasSuffix(r.URL.Path, "/stream"):
			handlers.StreamOutputHandler(w, r)
		case strings.HasSuffix(r.URL.Path, "/close"):
			handlers.CloseShellHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	// TODO: Add API handlers that use the loaded modelPlugins

	// 3. Start the server
	log.Printf("Starting server at port %s", e.config.ServerPort)
	return http.ListenAndServe(":"+e.config.ServerPort, mux)
}

// initializeModels loads all models specified in the config file.
func (e *Engine) initializeModels() error {
	log.Println("--- Initializing Models from Config ---")
	vision.PythonVenvPath = e.config.PythonVenvPath // Set the python path for vision models

	for _, modelCfg := range e.config.Models {
		log.Printf("Loading model: %s (type: %s, path: %s)", modelCfg.Name, modelCfg.Type, modelCfg.Path)
		switch modelCfg.Type {
		case "blip", "clip", "cliption":
			// For now, we are just logging. The next step is to create a plugin interface.
			// For example: plugin, err := vision.New(modelCfg)
		case "gguf":
			// For example: plugin, err := llm.NewGGUFModel(...)
		default:
			log.Printf("Warning: Unknown model type '%s' for model '%s'", modelCfg.Type, modelCfg.Name)
		}
	}
	log.Println("--- Model Initialization Complete ---")
	return nil
}