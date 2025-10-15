package main

import (
	"github.com/owen-6936/llm-cortex/examples/models/vision"
)

func main() {
	// Running all vision models in parallel
	vision.ParallelLoadExample()

	// // Running Blip Model...
	// vision.BlipExample()

	// // Running CLIP Model...
	// vision.ClipExample()

	// // Running CLIPtion Model...
	// vision.CliptionExample()
}