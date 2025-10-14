package vision

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/owen-6936/llm-cortex/spawn"
)

// BlipResponse represents the JSON output from the blip.py script.
type BlipResponse struct {
	Caption string  `json:"caption"`
	Latency float32 `json:"latency"`
	Prompt  string  `json:"prompt"`
	Image   string  `json:"image"`
}

var (
	blipSessions = make(map[string]string)
	blipMutex    = &sync.Mutex{}
)

const BLIP_JSON_DELIMITER = "END_OF_JSON"

// InvokeBlip runs the BLIP model on the given image with optional prompt.
// This function is a convenience wrapper that loads the model, sends a single prompt,
// and then unloads the model. It is less efficient for multiple sequential prompts.
func InvokeBlip(modelPath string, imagePath string, prompt string, useFast bool, legacy bool, maxLength int16, device string) (BlipResponse, error) {
	blipInstance, err := NewBlip(modelPath, device)
	if err != nil {
		return BlipResponse{}, err
	}
	defer blipInstance.UnloadBlipModel()
	return blipInstance.SendPrompt(imagePath, prompt, useFast, legacy, maxLength)
}

// Blip represents a loaded BLIP model instance, managed as a persistent
// interactive Python process.
type Blip struct {
	ModelPath string // Path to the model files.
	SessionID string // The unique ID for the underlying shell session.
	Device    string // The device the model is running on ('cpu' or 'cuda').
}

// NewBlip loads a BLIP model into memory by starting a persistent Python process
// in interactive mode. It returns a Blip struct instance which can be used to
// send multiple prompts efficiently.
func NewBlip(modelPath string, device string) (*Blip, error) {
	blipMutex.Lock()
	defer blipMutex.Unlock()

	sessionID, ok := blipSessions[modelPath]
	if !ok {
		var err error
		cmd := []string{
			PythonVenvPath,
			"python/models/vision/blip.py",
			"--model-path", modelPath,
			"--device", device,
			"--interactive",
		}
		sessionID, err = spawn.NewShellWithCommand(cmd...)
		if err != nil {
			return nil, fmt.Errorf("failed to start blip session: %w", err)
		}
		blipSessions[modelPath] = sessionID
		spawn.StartReading(sessionID, spawn.OutputHandler, spawn.ErrorOutputHandler)

		// Wait for the Python script to signal that the model is ready.
		// This is more robust than a fixed time.Sleep().
		err = spawn.WaitForString(sessionID, "[BLIP] Ready.", 120*time.Second) // 60-second timeout for model loading
		if err != nil {
			// Clean up the failed session
			spawn.CloseSession(sessionID)
			delete(blipSessions, modelPath)
			return nil, fmt.Errorf("error waiting for BLIP model to load: %w", err)
		}
	}

	return &Blip{
		ModelPath: modelPath,
		SessionID: sessionID,
		Device:    device,
	}, nil
}

// SendPrompt sends a request to the loaded BLIP model.
// It marshals the request, sends it to the Python process, and parses the JSON response.
func (b *Blip) SendPrompt(imagePath string, prompt string, useFast bool, legacy bool, maxLength int16) (BlipResponse, error) {
	// Create JSON request for the interactive Python script
	request := map[string]interface{}{
		"image_path": imagePath,
		"prompt":     prompt,
		"max_length": maxLength,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return BlipResponse{}, fmt.Errorf("failed to marshal blip request: %w", err)
	}

	// Send command and wait for response
	output, err := spawn.SendCommandAndWait(b.SessionID, string(jsonRequest), BLIP_JSON_DELIMITER)
	if err != nil {
		return BlipResponse{}, fmt.Errorf("failed to execute blip command: %w", err)
	}

	var response BlipResponse
	var errorResponse struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal([]byte(output), &errorResponse); err == nil && errorResponse.Error != "" {
		return BlipResponse{}, fmt.Errorf("blip.py script error: %s", errorResponse.Error)
	}

	if err := json.Unmarshal([]byte(output), &response); err != nil {
		return BlipResponse{}, fmt.Errorf("failed to parse blip.py output: %w\nOutput: %s", err, output)
	}

	return response, nil
}

// UnloadModel terminates the persistent Python process and cleans up resources.
func (b *Blip) UnloadBlipModel() error {
	blipMutex.Lock()
	defer blipMutex.Unlock()

	if sessionID, ok := blipSessions[b.ModelPath]; ok {
		err := spawn.CloseSession(sessionID)
		if err != nil {
			return fmt.Errorf("failed to close blip session %s: %w", sessionID, err)
		}
		delete(blipSessions, b.ModelPath)
	}
	return nil
}