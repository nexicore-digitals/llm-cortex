package router

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func InvokeLLM() LLMReply {
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
		panic("Error: --prompt is required")
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
		panic(err)
	}
	fullResponse := string(output)
	elapsed := time.Since(start)

	var llmReply LLMReply
	const BUILDSTARTINDEXSUBSTRCOUNT = 7
	const RESPONSESTARTINDEXSUBSTRCOUNT = 10
	buildIndexStart := strings.Index(fullResponse, "build: ") + BUILDSTARTINDEXSUBSTRCOUNT
	buildIndexEnd := strings.Index(fullResponse, "\nmain: llama backend init")
	responseIndexStart := strings.LastIndex(fullResponse, "Assistant:") + RESPONSESTARTINDEXSUBSTRCOUNT
	responseIndexEnd := strings.Index(fullResponse, "\n> EOF by user")
	build := fullResponse[buildIndexStart:buildIndexEnd]
	response := fullResponse[responseIndexStart:responseIndexEnd]
	llmReply.Build = build
	llmReply.Model = strings.Split(modelPath, "/")[2]
	llmReply.TimeElasped = elapsed
	llmReply.Response = response

	return llmReply
}

type LLMReply struct {
	Build       string
	Model       string
	Response    string
	TimeElasped time.Duration
}
