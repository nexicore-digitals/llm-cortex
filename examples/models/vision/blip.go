package vision

import (
	"fmt"
	"log"

	"github.com/owen-6936/llm-cortex/core/models/vision"
)

func BlipExample() {
	modelPath := "models/blip2-flan-t5-xl"
	imagePath := "samples/images/istockphoto-1155240408-612x612.jpg"
	prompt := "Question: describe this image in detail. Answer:"
	device := "cpu"

	fmt.Println("--- Testing BLIP Model ---")

	// Load the model (starts the persistent Python process)
	blipResponse, err := vision.InvokeBlip(modelPath, imagePath, prompt, true, false, 75, device)
	if err != nil {
		log.Fatalf("Failed to load BLIP model: %v", err)
	}

	fmt.Printf("Model loaded. Sending prompt for image: %s\n", imagePath)

	if err != nil {
		log.Fatalf("Failed to send prompt to BLIP model: %v", err)
	}

	fmt.Println("--- BLIP Response ---")
	fmt.Printf("Caption: %s\n", blipResponse.Caption)
	fmt.Printf("Latency: %.2f seconds\n", blipResponse.Latency)
	fmt.Println("")
	fmt.Println("-------------------------")
	fmt.Println("")
}