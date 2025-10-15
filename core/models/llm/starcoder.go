package llm

// InvokeStarCoder executes the StarCoder GGUF model using llama-cli.
// This is a convenience wrapper that loads the model, sends a single prompt,
// and then unloads the model.
func InvokeStarCoder(config Settings) (string, error) {
	model, err := NewGGUFModel(config)
	if err != nil {
		return "", err
	}
	defer model.Unload()
	return model.SendPrompt(config.Prompt)
}