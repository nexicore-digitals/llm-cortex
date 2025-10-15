package llm

import (
	"runtime"
	"strconv"
)

// Settings holds the parameters for running a GGUF model with llama-cli.
type Settings struct {
	ModelPath     string
	Prompt        string
	Threads       int
	NPredict      int
	NoMMap        bool
	BatchSize     int
	JinjaTemplate string
	// Add other general parameters here as needed.
}

// ToArgs converts the settings to a slice of command-line arguments for llama-cli.
func (s *Settings) ToArgs(interactive bool) []string {
	args := []string{"-m", s.ModelPath}

	if interactive {
		args = append(args, "-i")
		// In interactive mode, prompt is sent to stdin, not as a startup argument.
	} else if s.Prompt != "" {
		args = append(args, "-p", s.Prompt)
	}
	if s.NPredict > 0 {
		args = append(args, "--n-predict", strconv.Itoa(s.NPredict))
	}
	if s.Threads > 0 {
		args = append(args, "--threads", strconv.Itoa(s.Threads))
	}
	if s.BatchSize > 0 {
		args = append(args, "--batch-size", strconv.Itoa(s.BatchSize))
	}
	if s.NoMMap {
		args = append(args, "--no-mmap")
	}
	if s.JinjaTemplate != "" {
		// Assuming llama-cli will have a --jinja flag
		// args = append(args, "--jinja", s.JinjaTemplate)
	}

	return args
}

// Performance returns a config optimized for performance.
// It uses all available CPU threads and a large batch size.
func Performance(modelPath, prompt string, nPredict int) Settings {
	return Settings{
		ModelPath: modelPath,
		Prompt:    prompt,
		Threads:   runtime.NumCPU(),
		NPredict:  nPredict,
		NoMMap:    true,
		BatchSize: 4096,
	}
}

// Balanced returns a config with sensible, general-purpose defaults.
func Balanced(modelPath, prompt string, nPredict int) Settings {
	return Settings{
		ModelPath: modelPath,
		Prompt:    prompt,
		Threads:   runtime.NumCPU() / 2, // Use half of the available cores
		NPredict:  nPredict,
		NoMMap:    false,
		BatchSize: 512, // Default llama.cpp batch size
	}
}