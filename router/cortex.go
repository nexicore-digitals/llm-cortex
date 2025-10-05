package cortex

import (
    "fmt"
    llama "github.com/go-skynet/go-llama.cpp"
)

type Cortex struct {
    Model *llama.LLama
}

func NewCortex(modelPath string, threads, gpulayers int) (*Cortex, error) {
    model, err := llama.New(modelPath,
        llama.EnableF16Memory,
        llama.SetContext(2048),
        llama.SetThreads(threads),
        llama.SetGPULayers(gpulayers),
    )
    if err != nil {
        return nil, err
    }
    return &Cortex{Model: model}, nil
}

func (c *Cortex) Run(prompt string, tokens int) error {
    _, err := c.Model.Predict(prompt,
        llama.SetTokens(tokens),
        llama.SetTopK(40),
        llama.SetTopP(0.9),
        llama.SetTokenCallback(func(token string) bool {
            fmt.Print(token)
            return true
        }),
    )
    return err
}
