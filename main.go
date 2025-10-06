package main

import (
	"fmt"
	"llm-cortex/router"
	"llm-cortex/spawn"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	// fmt.Println("Starting server at port 8080")
	// err := http.ListenAndServe(":8080", mux)
	// if err != nil {
	// 	panic(err)
	// }
	var sessions = make(map[string]*spawn.ShellSession)
	sessionId, err := spawn.NewShell(sessions)
	if err != nil {
		panic(err)
	}

	spawn.StartReading(sessions[sessionId], spawn.OutputHandler)

	fmt.Println("\n--- Sending Command 2: 'pwd' ---")
	sendCmdErr := spawn.SendCommand(sessions, sessionId, "pwd")
	if sendCmdErr != nil {
		fmt.Println("Error sending command:", err)
	}

	fmt.Println("\n--- Sending Command 2:  ---")
	sendCmdErr2 := spawn.SendCommand(sessions, sessionId, "echo hello from Owen")
	if sendCmdErr2 != nil {
		fmt.Println("Error sending command:", err)
	}

	fmt.Println("\n--- Sending Command 3: ---")
	sendCmdErr3 := spawn.SendCommand(sessions, sessionId, "echo Hmm imagine if this is from an LLM")
	if sendCmdErr3 != nil {
		fmt.Println("Error sending command:", err)
	}

	fmt.Println("\n--- Closing Session ---")
	if err := spawn.CloseSession(sessions, sessionId); err != nil {
		fmt.Println("Error closing session:", err)
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
