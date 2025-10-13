package vision

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// BlipResponse represents the output of the BLIP model.
type BlipResponse struct {
	Caption string  // generated caption
	Latency float32 // time taken for inference
	Prompt  string  // input prompt
	Image   string  // image path
}


// InvokeBlip runs the BLIP model on the given image with optional prompt.
// model_path: path to the BLIP model
// image_path: path to the input image
// prompt: optional text prompt
// use_fast: whether to use the fast processor
// max_length: maximum caption length
func InvokeBlip(modelPath string, imagePath string, prompt string, useFast bool, legacy bool, maxLength int16) BlipResponse {
	args := []string{
		"python/models/vision/blip.py",
		"--model-path", modelPath,
		"--image-path", imagePath,
		"--prompt", prompt,
		"--max-length", fmt.Sprintf("%d", maxLength),
	}

	if useFast {
		args = append(args, "--use-fast")
	}
	if !legacy {
		args = append(args, "--no-legacy")
	}

	cmd := exec.Command("python3", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Panicf("Failed to run blip.py: %v\nStderr: %s", err, stderr.String())
	}

	var response BlipResponse
	// The python script might return an error object.
	// We'll try to unmarshal into that first.
	var errorResponse struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(out.Bytes(), &errorResponse); err == nil && errorResponse.Error != "" {
		log.Panicf("blip.py returned an error: %s", errorResponse.Error)
	}

	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		log.Panicf("Failed to parse blip.py output: %v\nStdout: %s", err, out.String())
	}

	return response
}