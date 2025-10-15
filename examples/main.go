package main

import (
	"github.com/owen-6936/llm-cortex/examples/models/llm"
	"github.com/owen-6936/llm-cortex/examples/models/vision"
)

func main() {
	// Running all llms in parallel
	llm.QwenExample()

	// Running all vision models in parallel
	vision.ParallelLoadExample()

	// // Running Blip Model...
	// vision.BlipExample()

	// // Running CLIP Model...
	// vision.ClipExample()

	// // Running CLIPtion Model...
	// vision.CliptionExample()
}