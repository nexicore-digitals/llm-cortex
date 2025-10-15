package llm

import (
	"fmt"

	"github.com/owen-6936/llm-cortex/core/models/llm"
)

func StarcoderExample() {
	fmt.Println("--- Testing StarCoder LLM (Session-based) ---")

	// Make sure this path points to your downloaded StarCoder GGUF model
	modelPath := "models/starcoder/starcoder2-15b-instruct-v0.1-Q4_K_M.gguf"
	nPredict := 256

	// Use balanced settings for this example
	config := llm.Balanced(modelPath, "", nPredict)

	// Create a new session-based model instance
	model, err := llm.NewGGUFModel(config)
	if err != nil {
		return
	}
	defer model.Unload()

	// --- Send first prompt ---
	prompt1 := "create a rust function that returns the factorial of a number"

	_,promptErr := model.SendPrompt(prompt1)
	if promptErr != nil {
		return
	} else {
	}

	fmt.Println("-------------------------")
}