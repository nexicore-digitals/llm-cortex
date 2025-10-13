package vision

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// CLIPtionResponse represents the output of the CLIPtion model.
type CLIPtionResponse struct {
	Caption string  `json:"caption"`
	Latency float32 `json:"latency"`
	Image   string  `json:"image"`
}

// InvokeCLIPtion runs the CLIPtion model on the given image.
func InvokeCLIPtion(modelPath string, imagePath string, useFast bool, beamSearch bool, beamWidth int, bestOf int, temperature float32, device string) CLIPtionResponse {
	args := []string{
		"python/models/vision/cliption/cliption.py",
		"--model-path", modelPath,
		"--image-path", imagePath,
		"--device", device,
	}

	if useFast {
		args = append(args, "--use-fast")
	}

	if beamSearch {
		args = append(args, "--beam-search")
		args = append(args, "--beam-width", fmt.Sprintf("%d", beamWidth))
	}
	args = append(args, "--best-of", fmt.Sprintf("%d", bestOf))
	args = append(args, "--temperature", fmt.Sprintf("%f", temperature))

	cmd := exec.Command("python3", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Panicf("Failed to run cliption.py: %v\nStderr: %s", err, stderr.String())
	}

	var response CLIPtionResponse
	// Check for a JSON error object from the script
	var errorResponse struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(out.Bytes(), &errorResponse); err == nil && errorResponse.Error != "" {
		log.Panicf("cliption.py returned an error: %s", errorResponse.Error)
	}

	// Parse the successful response
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		log.Panicf("Failed to parse cliption.py output: %v\nStdout: %s", err, out.String())
	}

	return response
}