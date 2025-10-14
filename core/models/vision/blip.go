package vision

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/owen-6936/llm-cortex/spawn"
	"github.com/owen-6936/llm-cortex/utils"
)

// BlipResponse represents the JSON output from the blip.py script.
type BlipResponse struct {
	Caption string  `json:"caption"`
	Latency float32 `json:"latency"`
	Prompt  string  `json:"prompt"`
	Image   string  `json:"image"`
}

var (
	blipManager = NewModelManager()
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
	sessionID, err := blipManager.Load(
		modelPath,
		device,
		"python/models/vision/blip.py",
		"[BLIP] Ready.",
		120*time.Second,
	)
	if err != nil {
		return nil, err
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
	utils.HandleError(err, "failed to marshal blip request")

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
	err := blipManager.Unload(b.ModelPath)
	if err != nil {
		return fmt.Errorf("failed to close blip session for %s: %w", b.ModelPath, err)
	}
	return nil
}