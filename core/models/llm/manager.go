package llm

import (
	"fmt"
	"sync"
	"time"

	"github.com/owen-6936/llm-cortex/spawn"
)

// LLMManager handles the lifecycle of persistent `llama-cli` processes.
type LLMManager struct {
	sessions map[string]string
	mutex    *sync.Mutex
}

// NewLLMManager creates a new manager for GGUF model sessions.
func NewLLMManager() *LLMManager {
	return &LLMManager{
		sessions: make(map[string]string),
		mutex:    &sync.Mutex{},
	}
}

// Load ensures a GGUF model is loaded, starting a new `llama-cli` session if one doesn't exist.
func (m *LLMManager) Load(config Settings) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	sessionID, ok := m.sessions[config.ModelPath]
	if !ok {
		var err error
		args := config.ToArgs(true) // Start in interactive mode
		cmd := append([]string{"bin/llama-cli"}, args...)

		sessionID, err = spawn.NewShellWithCommand(cmd...)
		if err != nil {
			return "", fmt.Errorf("failed to start llama-cli session: %w", err)
		}
		m.sessions[config.ModelPath] = sessionID
		spawn.StartReading(sessionID, spawn.OutputHandler, spawn.InfoOutputHandler)

		// Wait for llama.cpp to be ready for input.
		err = spawn.WaitForString(sessionID, "\n> ", 180*time.Second)
		if err != nil {
			spawn.CloseSession(sessionID)
			delete(m.sessions, config.ModelPath)
			return "", fmt.Errorf("error waiting for GGUF model '%s' to load: %w", config.ModelPath, err)
		}
	}
	return sessionID, nil
}

// Unload closes the session for a given model path.
func (m *LLMManager) Unload(modelPath string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if sessionID, ok := m.sessions[modelPath]; ok {
		delete(m.sessions, modelPath)
		// In interactive mode, llama-cli exits on EOF, which `CloseSession` handles.
		return spawn.CloseSession(sessionID)
	}
	return nil
}