package vision

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

// ClipResponse represents the output of the CLIP model.
type ClipResponse struct {
	Results map[string]float32 `json:"results"`
	Latency float32            `json:"latency"`
	Image   string             `json:"image"`
}

// InvokeClip runs the CLIP model on the given image against a list of text labels.
func InvokeClip(modelPath string, imagePath string, texts []string, useFast bool, device string) ClipResponse {
	args := []string{
		"python/models/vision/clip.py",
		"--model-path", modelPath,
		"--image-path", imagePath,
		"--device", device,
	}

	if useFast {
		args = append(args, "--use-fast")
	}

	// Append the text labels for classification
	args = append(args, "--texts")
	args = append(args, texts...)

	cmd := exec.Command("python3", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Panicf("Failed to run clip.py: %v\nStderr: %s", err, stderr.String())
	}

	var response ClipResponse
	// Check for a JSON error object from the script
	var errorResponse struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(out.Bytes(), &errorResponse); err == nil && errorResponse.Error != "" {
		log.Panicf("clip.py returned an error: %s", errorResponse.Error)
	}

	// Parse the successful response
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		log.Panicf("Failed to parse clip.py output: %v\nStdout: %s", err, out.String())
	}

	return response
}