package llm

// InvokeQwen executes the Qwen GGUF model using llama-cli.
// This is a convenience wrapper that loads the model, sends a single prompt,
// and then unloads the model.
func InvokeQwen(config Settings) (string, error) {
	model, err := NewGGUFModel(config)
	if err != nil {
		return "", err
	}
	defer model.Unload()

	// The prompt is part of the config struct for one-shot invocation.
	return model.SendPrompt(config.Prompt)
}
