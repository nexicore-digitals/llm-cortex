package vision

import (
	"fmt"
	"log"

	"github.com/owen-6936/llm-cortex/core/models/vision"
)

func CliptionExample() {
	modelPath := "models/CLIPtion"
	imagePath := "samples/images/istockphoto-1155240408-612x612.jpg"
	device := "cpu"

	fmt.Println("--- Testing CLIPtion Model ---")

	// Load the model
	cliptionInstance, err := vision.NewCLIPtion(modelPath, device)
	if err != nil {
		log.Fatalf("Failed to load CLIPtion model: %v", err)
	}

	defer cliptionInstance.UnloadCLIPtionModel()

	fmt.Printf("Model loaded. Sending prompt for image: %s\n", imagePath)

	// Send a prompt to the loaded model using beam search for higher quality
	response, err := cliptionInstance.SendPrompt(imagePath, true, true, 5, 5, 1.0)
	if err != nil {
		log.Fatalf("Failed to send prompt to CLIPtion model: %v", err)
	}

	fmt.Println("--- CLIPtion Response ---")
	fmt.Printf("Caption: %s\n", response.Caption)
	fmt.Printf("Latency: %.2f seconds\n", response.Latency)
	fmt.Println("")
	fmt.Println("-------------------------")
	fmt.Println("")
}