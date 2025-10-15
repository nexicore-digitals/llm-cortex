package main

import (
	"fmt"
	"net/http"
	"strings"

	llm_example "github.com/owen-6936/llm-cortex/examples/models/llm"

	"github.com/owen-6936/llm-cortex/examples/models/vision"
	"github.com/owen-6936/llm-cortex/handlers"
	"github.com/owen-6936/llm-cortex/router"
	"github.com/owen-6936/llm-cortex/scheduler"
	"github.com/owen-6936/llm-cortex/utils"
)

func main() {
	// Example of running all vision models in parallel using the new scheduler.
	fmt.Println("--- Running All Vision Models in Parallel ---")
	taskRunner := scheduler.NewTaskRunner(
		vision.BlipExample,
		vision.ClipExample,
		vision.CliptionExample,
		llm_example.QwenExample,
		llm_example.StarcoderExample,
	)
	taskRunner.Run()
	fmt.Println("--- All Vision Models Finished ---")

	mux := http.NewServeMux()
	mux.HandleFunc("/serve", rootHandler)

	os := http.FileServer(http.Dir("ui"))
	mux.Handle("/", os)

	mux.HandleFunc("/shell/start", handlers.StartShellHandler)
	mux.HandleFunc("/shell/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/send"):
			handlers.SendCommandHandler(w, r)
		case strings.HasSuffix(r.URL.Path, "/stream"):
			handlers.StreamOutputHandler(w, r)
		case strings.HasSuffix(r.URL.Path, "/close"):
			handlers.CloseShellHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	fmt.Println("Starting server at port 8080")
	err := http.ListenAndServe(":8080", mux)
	utils.HandleError(err, "Failed to start server", true)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-type", "plain/text")
	fmt.Fprint(w, "Hello from the go root router")
	params := r.URL.Query()
	if run, ok := params["run"]; ok && run[0] == "true" {
		llmReply := router.InvokeLLM()
		fmt.Printf("build: %s\nmodel: %s\nresponse: %s\ntimeElasped: %s\n", llmReply.Build, llmReply.Model, llmReply.Response, llmReply.TimeElasped)
	}
}
