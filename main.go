package main

import (
	"fmt"
	"llm-cortex/handlers"
	"llm-cortex/router"
	"net/http"
	"strings"
)

func main() {
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
	if err != nil {
		panic(err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-type", "plain/text")
	fmt.Fprint(w, "Hello from the go root router")
	params := r.URL.Query()
	for key, value := range params {
		fmt.Printf("%s: %s\n", key, value)
		if key == "run" && value[0] == "true" {
			llmReply := router.InvokeLLM()
			fmt.Printf("build: %s\nmodel: %s\nresponse: %s\ntimeElasped: %s\n", llmReply.Build, llmReply.Model, llmReply.Response, llmReply.TimeElasped)
		}
	}
}
