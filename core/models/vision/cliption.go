package vision

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/owen-6936/llm-cortex/spawn"
)

// CLIPtionResponse represents the JSON output from the cliption.py script.
type CLIPtionResponse struct {
	Caption string  `json:"caption"`
	Latency float32 `json:"latency"`
	Image   string  `json:"image"`
}

var (
	cliptionManager = NewModelManager()
)

const CLIPTION_JSON_DELIMITER = "END_OF_JSON"

// InvokeCLIPtion is a convenience wrapper that loads the model, sends a single prompt,
// and then unloads the model. It is less efficient for multiple sequential prompts.
func InvokeCLIPtion(modelPath string, imagePath string, useFast bool, beamSearch bool, beamWidth int, bestOf int, temperature float32, device string) (CLIPtionResponse, error) {
	cliptionInstance, err := NewCLIPtion(modelPath, device)
	if err != nil {
		return CLIPtionResponse{}, err
	}
	defer cliptionInstance.UnloadCLIPtionModel()

	return cliptionInstance.SendPrompt(imagePath, useFast, beamSearch, beamWidth, bestOf, temperature)
}

// CLIPtion represents a loaded CLIPtion model instance, managed as a persistent
// interactive Python process.
type CLIPtion struct {
	ModelPath string // Path to the model files.
	SessionID string // The unique ID for the underlying shell session.
	Device    string // The device the model is running on ('cpu' or 'cuda').
}

// NewCLIPtion loads a CLIPtion model into memory by starting a persistent Python process
// in interactive mode. It returns a CLIPtion struct instance which can be used to
// send multiple prompts efficiently.
func NewCLIPtion(modelPath string, device string) (*CLIPtion, error) {
	sessionID, err := cliptionManager.Load(
		modelPath,
		device,
		"python/models/vision/cliption/cliption.py",
		"[CLIPtion] Ready.",
		90*time.Second,
	)
	if err != nil {
		return nil, err
	}

	return &CLIPtion{
		ModelPath: modelPath,
		SessionID: sessionID,
		Device:    device,
	}, nil
}

// SendPrompt sends a request to the loaded CLIPtion model.
// It marshals the request, sends it to the Python process, and parses the JSON response.
func (c *CLIPtion) SendPrompt(imagePath string, useFast bool, beamSearch bool, beamWidth int, bestOf int, temperature float32) (CLIPtionResponse, error) {
	request := map[string]interface{}{
		"image_path":  imagePath,
		"use_fast":    useFast,
		"beam_search": beamSearch,
		"beam_width":  beamWidth,
		"best_of":     bestOf,
		"temperature": temperature,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return CLIPtionResponse{}, fmt.Errorf("failed to marshal cliption request: %w", err)
	}

	output, err := spawn.SendCommandAndWait(c.SessionID, string(jsonRequest), CLIPTION_JSON_DELIMITER)
	if err != nil {
		return CLIPtionResponse{}, fmt.Errorf("failed to execute cliption command: %w", err)
	}

	var response CLIPtionResponse
	// Check for a JSON error object from the script
	var errorResponse struct{ Error string `json:"error"` }
	if err := json.Unmarshal([]byte(output), &errorResponse); err == nil && errorResponse.Error != "" {
		return CLIPtionResponse{}, fmt.Errorf("cliption.py returned an error: %s", errorResponse.Error)
	}

	// Parse the successful response
	if err := json.Unmarshal([]byte(output), &response); err != nil {
		return CLIPtionResponse{}, fmt.Errorf("failed to parse cliption.py output: %w\nOutput: %s", err, output)
	}

	return response, nil
}

// UnloadModel terminates the persistent Python process and cleans up resources.
func (c *CLIPtion) UnloadCLIPtionModel() error {
	err := cliptionManager.Unload(c.ModelPath)
	if err != nil {
		return fmt.Errorf("failed to close cliption session for %s: %w", c.ModelPath, err)
	}
	return nil
}