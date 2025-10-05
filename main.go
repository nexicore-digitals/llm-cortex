package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"
)

func main() {
	var modelPath string
	var prompt string
	var tokens int
	var threads int

	flag.StringVar(&modelPath, "model", "./models/deepseek-7b/deepseek-llm-7b-chat.Q4_K_M.gguf", "Path to GGUF model")
	flag.StringVar(&prompt, "prompt", "Hello what is your name", "Prompt to send to the model")
	flag.IntVar(&tokens, "n_predict", 128, "Number of tokens to predict")
	flag.IntVar(&threads, "threads", 8, "Number of threads to use")
	flag.Parse()

	if prompt == "" {
		fmt.Println("Error: --prompt is required")
		return
	}

	start := time.Now()

	cmd := exec.Command("./bin/llama-cli",
		"--model", modelPath,
		"--prompt", prompt,
		"--n_predict", fmt.Sprintf("%d", tokens),
		"--threads", fmt.Sprintf("%d", threads),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running llama-cli:", err)
		return
	}
	fmt.Println(string(output))

	elapsed := time.Since(start)
	fmt.Printf("\n⏱️ Eval time: %s\n", elapsed)
}
