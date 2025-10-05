package main

import (
    "flag"
    "fmt"
    "os"
    "runtime"

    "github.com/nexicore-digitals/llm-cortex/router/cortex"
)

func main() {
    var modelPath string
    var prompt string
    var tokens int
    var threads int

    flag.StringVar(&modelPath, "model", "./models/deepseek-7b/deepseek-llm-7b-chat.Q4_K_M.gguf", "Path to GGUF model")
    flag.StringVar(&prompt, "prompt", "", "Prompt to send to the model")
    flag.IntVar(&tokens, "tokens", 128, "Number of tokens to predict")
    flag.IntVar(&threads, "threads", runtime.NumCPU(), "Number of threads to use")
    flag.Parse()

    if prompt == "" {
        fmt.Println("Error: --prompt is required")
        os.Exit(1)
    }

    fmt.Printf("🧠 Loading model: %s\n", modelPath)
    c, err := cortex.NewCortex(modelPath, threads, 0)
    if err != nil {
        fmt.Println("Error loading model:", err)
        os.Exit(1)
    }

    fmt.Printf("🚀 Running prompt: %s\n\n", prompt)
    err = c.Run(prompt, tokens)
    if err != nil {
        fmt.Println("Error during inference:", err)
        os.Exit(1)
    }
}
