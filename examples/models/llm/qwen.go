package llm

import (
	"fmt"
	"log"

	"github.com/owen-6936/llm-cortex/core/models/llm"
)

func QwenExample() {
	fmt.Println("--- Testing Qwen LLM (Session-based) ---")

	modelPath := "models/qwen/Qwen2.5-Coder-7B-Instruct-Q6_K_L.gguf"
	nPredict := 128

	// Get performance settings. The prompt is no longer needed here for initialization.
	config := llm.Performance(modelPath, "", nPredict)

	// Create a new session-based model instance
	model, err := llm.NewGGUFModel(config)
	if err != nil {
		log.Printf("Failed to load GGUF model: %v", err)
		return
	}
	defer model.Unload()

	// --- Send first prompt ---
	prompt1 := "def fibonacci(n):"
	_,err1:= model.SendPrompt(prompt1)
	if err1 != nil {
		
	} else {
		
	}
	// --- Send second prompt ---
	prompt2 := "write a go function to check if a number is prime"
	_, err2 := model.SendPrompt(prompt2)
	if err2 != nil {
		
	} else {
		
	}
}