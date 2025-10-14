package main

import (
	"fmt"

	"github.com/owen-6936/llm-cortex/examples/models/vision"
	"github.com/owen-6936/llm-cortex/scheduler"
)

func main() {
	// Example of running all vision models in parallel using the new scheduler.
	fmt.Println("--- Running All Vision Models in Parallel ---")
	taskRunner := scheduler.NewTaskRunner(
		vision.BlipExample,
		vision.ClipExample,
		vision.CliptionExample,
	)
	taskRunner.Run()
	fmt.Println("--- All Vision Models Finished ---")

	// The web server logic can be added back here if needed.
}
