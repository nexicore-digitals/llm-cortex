package main

import (
	"fmt"
	"llm-cortex/router"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

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
