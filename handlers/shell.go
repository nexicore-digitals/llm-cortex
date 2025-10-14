package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/owen-6936/llm-cortex/spawn"
)

var mu sync.Mutex

// StartShellHandler spawns a new shell and returns its session ID
func StartShellHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := spawn.NewShell()
	if err != nil {
		http.Error(w, "Failed to start shell", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": sessionID})

	// The output handler now needs the session object, which we get from the spawn manager
	spawn.StartReading(sessionID, spawn.OutputHandler, spawn.ErrorOutputHandler)

}

// SendCommandHandler sends a command to a shell session
func SendCommandHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		http.Error(w, "Missing session ID", http.StatusBadRequest)
		return
	}
	session, ok := spawn.GetSession(id)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	// reset the output string after each command execution
	mu.Lock()
	session.OutputBuf.Reset()
	mu.Unlock()

	var payload struct {
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid command payload", http.StatusBadRequest)
		return
	}

	if err := spawn.SendCommand(id, payload.Command); err != nil {
		http.Error(w, "Failed to send command", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// StreamOutputHandler returns the latest output from a shell session
func StreamOutputHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	session, ok := spawn.GetSession(id)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-type", "plain/text")
	// Inside StreamOutputHandler:
	// CRITICAL FIX: Lock before reading shared state
	mu.Lock()
	output := session.OutputBuf.Bytes()
	mu.Unlock()
	w.Write(output)
}

// CloseShellHandler gracefully shuts down a shell session
func CloseShellHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if err := spawn.CloseSession(id); err != nil {
		http.Error(w, "Failed to close session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// extractID parses the session ID from the URL path
func extractID(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}
