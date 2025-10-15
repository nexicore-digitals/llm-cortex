package vision

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/owen-6936/llm-cortex/spawn"
	"github.com/owen-6936/llm-cortex/utils"
)

// ClipResponse represents the JSON output from the clip.py script.
type ClipResponse struct {
	Results map[string]float32 `json:"results"`
	Latency float32            `json:"latency"`
	Image   string             `json:"image"`
}

var (
	clipManager = NewModelManager()
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
	sessionID, err := clipManager.Load(
		modelPath,
		device,
		"python/models/vision/clip.py",
		"[CLIP] Ready.",
		90*time.Second,
	)
	if err != nil {
		return nil, err
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
	utils.HandleError(err, "failed to marshal clip request")

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
	err := clipManager.Unload(c.ModelPath)
	if err != nil {
		return fmt.Errorf("failed to close clip session for %s: %w", c.ModelPath, err)
	}
	return nil
}