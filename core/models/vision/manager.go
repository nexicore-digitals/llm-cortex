package vision

import (
	"fmt"
	"sync"
	"time"

	"github.com/owen-6936/llm-cortex/spawn"
)

// ModelManager handles the lifecycle of persistent Python model processes.
type ModelManager struct {
	sessions map[string]string
	mutex    *sync.Mutex
}

// NewModelManager creates a new manager for model sessions.
func NewModelManager() *ModelManager {
	return &ModelManager{
		sessions: make(map[string]string),
		mutex:    &sync.Mutex{},
	}
}

// Load ensures a model is loaded, starting a new session if one doesn't exist for the given model path.
func (m *ModelManager) Load(modelPath, device, pythonScript, readyString string, timeout time.Duration) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	sessionID, ok := m.sessions[modelPath]
	if !ok {
		var err error
		cmd := []string{
			PythonVenvPath,
			pythonScript,
			"--model-path", modelPath,
			"--device", device,
			"--interactive",
		}
		sessionID, err = spawn.NewShellWithCommand(cmd...)
		if err != nil {
			return "", fmt.Errorf("failed to start session for %s: %w", pythonScript, err)
		}
		m.sessions[modelPath] = sessionID
		spawn.StartReading(sessionID, spawn.OutputHandler, spawn.ErrorOutputHandler)

		err = spawn.WaitForString(sessionID, readyString, timeout)
		if err != nil {
			spawn.CloseSession(sessionID)
			delete(m.sessions, modelPath)
			return "", fmt.Errorf("error waiting for model '%s' to load: %w", modelPath, err)
		}
	}
	return sessionID, nil
}

// Unload closes the session for a given model path.
func (m *ModelManager) Unload(modelPath string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if sessionID, ok := m.sessions[modelPath]; ok {
		delete(m.sessions, modelPath)
		return spawn.CloseSession(sessionID)
	}
	return nil
}