package vision

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/owen-6936/llm-cortex/spawn"
)

// ClipResponse represents the JSON output from the clip.py script.
type ClipResponse struct {
	Results map[string]float32 `json:"results"`
	Latency float32            `json:"latency"`
	Image   string             `json:"image"`
}

var (
	clipSessions = make(map[string]string)
	clipMutex    = &sync.Mutex{}
)

const CLIP_JSON_DELIMITER = "END_OF_JSON"

// InvokeClip runs the CLIP model on the given image against a list of text labels.
// This function is a convenience wrapper that loads the model, sends a single prompt,
// and then unloads the model. It is less efficient for multiple sequential prompts.
func InvokeClip(modelPath string, imagePath string, texts []string, useFast bool, device string) (ClipResponse, error) {
	clipInstance, err := NewClip(modelPath, device)
	if err != nil {
		return ClipResponse{}, err
	}
	defer clipInstance.UnloadClipModel()

	return clipInstance.SendPrompt(imagePath, texts, useFast)
}

// Clip represents a loaded CLIP model instance, managed as a persistent
// interactive Python process.
type Clip struct {
	ModelPath string // Path to the model files.
	SessionID string // The unique ID for the underlying shell session.
	Device    string // The device the model is running on ('cpu' or 'cuda').
}

// NewClip loads a CLIP model into memory by starting a persistent Python process
// in interactive mode. It returns a Clip struct instance which can be used to
// send multiple prompts efficiently.
func NewClip(modelPath string, device string) (*Clip, error) {
	clipMutex.Lock()
	defer clipMutex.Unlock()

	sessionID, ok := clipSessions[modelPath]
	if !ok {
		var err error
		cmd := []string{
			PythonVenvPath,
			"python/models/vision/clip.py",
			"--model-path", modelPath,
			"--interactive",
			"--device", device,
		}
		sessionID, err = spawn.NewShellWithCommand(cmd...)
		if err != nil {
			return nil, fmt.Errorf("failed to start clip session: %w", err)
		}
		clipSessions[modelPath] = sessionID
		spawn.StartReading(sessionID, spawn.OutputHandler, spawn.ErrorOutputHandler)

		// Wait for the Python script to signal that the model is ready.
		err = spawn.WaitForString(sessionID, "[CLIP] Ready.", 90*time.Second) // 90-second timeout
		if err != nil {
			// Clean up the failed session
			spawn.CloseSession(sessionID)
			delete(clipSessions, modelPath)
			return nil, fmt.Errorf("error waiting for CLIP model to load: %w", err)
		}
	}

	return &Clip{
		ModelPath: modelPath,
		SessionID: sessionID,
		Device:    device,
	}, nil
}

// SendPrompt sends a request to the loaded CLIP model.
// It marshals the request, sends it to the Python process, and parses the JSON response.
func (c *Clip) SendPrompt(imagePath string, texts []string, useFast bool) (ClipResponse, error) {
	request := map[string]interface{}{
		"image_path": imagePath,
		"texts":      texts,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return ClipResponse{}, fmt.Errorf("failed to marshal clip request: %w", err)
	}

	output, err := spawn.SendCommandAndWait(c.SessionID, string(jsonRequest), CLIP_JSON_DELIMITER)
	if err != nil {
		return ClipResponse{}, fmt.Errorf("failed to execute clip command: %w", err)
	}

	var response ClipResponse
	// Check for a JSON error object from the script
	var errorResponse struct{ Error string `json:"error"` }
	if err := json.Unmarshal([]byte(output), &errorResponse); err == nil && errorResponse.Error != "" {
		return ClipResponse{}, fmt.Errorf("clip.py returned an error: %s", errorResponse.Error)
	}

	// Parse the successful response
	if err := json.Unmarshal([]byte(output), &response); err != nil {
		return ClipResponse{}, fmt.Errorf("failed to parse clip.py output: %w\nOutput: %s", err, output)
	}

	return response, nil
}

// UnloadModel terminates the persistent Python process and cleans up resources.
func (c *Clip) UnloadClipModel() error {
	clipMutex.Lock()
	defer clipMutex.Unlock()

	if sessionID, ok := clipSessions[c.ModelPath]; ok {
		err := spawn.CloseSession(sessionID)
		if err != nil {
			return fmt.Errorf("failed to close clip session %s: %w", sessionID, err)
		}
		delete(clipSessions, c.ModelPath)
	}
	return nil
}