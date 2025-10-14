package vision

import (
	"fmt"
	"log"

	"github.com/owen-6936/llm-cortex/core/models/vision"
)

func ClipExample() {
	modelPath := "models/clip-vit-b-32"
	imagePath := "samples/images/istockphoto-1155240408-612x612.jpg"
	device := "cpu"
	texts := []string{"a photo of food", "a drawing of a car", "a picture of people"}

	fmt.Println("--- Testing CLIP Model ---")

	// Load the model
	clipInstance, err := vision.NewClip(modelPath, device)
	if err != nil {
		log.Fatalf("Failed to load CLIP model: %v", err)
	}
	defer clipInstance.UnloadClipModel()

	fmt.Printf("Model loaded. Sending prompt for image: %s\n", imagePath)

	// Send a prompt to the loaded model
	response, err := clipInstance.SendPrompt(imagePath, texts, true)
	if err != nil {
		log.Fatalf("Failed to send prompt to CLIP model: %v", err)
	}

	fmt.Println("--- CLIP Response ---")
	fmt.Println("Classification Probabilities:")
	for text, prob := range response.Results {
		fmt.Printf("  - \"%s\": %.4f\n", text, prob)
	}
	fmt.Printf("Latency: %.2f seconds\n", response.Latency)
	fmt.Println("")
	fmt.Println("-------------------------")
	fmt.Println("")
}