package llm

import (
	"fmt"

	"github.com/owen-6936/llm-cortex/spawn"
)

var (
	llmManager = NewLLMManager()
)

// GGUFModel represents a loaded GGUF model instance, managed as a persistent
// interactive `llama-cli` process.
type GGUFModel struct {
	Settings  Settings
	SessionID string
}

// NewGGUFModel loads a GGUF model into memory by starting a persistent `llama-cli` process
// in interactive mode.
func NewGGUFModel(config Settings) (*GGUFModel, error) {
	sessionID, err := llmManager.Load(config)
	if err != nil {
		return nil, err
	}

	return &GGUFModel{
		Settings:  config,
		SessionID: sessionID,
	}, nil
}

// SendPrompt sends a prompt to the loaded GGUF model.
// It sends the prompt to the `llama-cli` process's stdin and waits for the response.
func (m *GGUFModel) SendPrompt(prompt string) (string, error) {
	// The first prompt needs to be handled differently as the buffer is not reset.
	// Subsequent prompts will use the standard SendCommandAndWait.
	// A simple way to check is to see if the buffer contains just the initial "> ".
	session, ok := spawn.GetSession(m.SessionID)
	if !ok {
		return "", fmt.Errorf("session not found for GGUF model")
	}

	// The delimiter `\n>` indicates it's ready for the next prompt.
	output, err := spawn.SendCommandAndWait(m.SessionID, prompt, "\n> ")
	if err != nil {
		return "", fmt.Errorf("failed to execute GGUF prompt: %w", err)
	}
	session.OutputBuf.Reset() // Clean buffer after successful command
	return output, nil
}

// Unload terminates the persistent `llama-cli` process.
func (m *GGUFModel) Unload() error {
	if err := llmManager.Unload(m.Settings.ModelPath); err != nil {
		return fmt.Errorf("failed to close GGUF session for %s: %w", m.Settings.ModelPath, err)
	}
	return nil
}