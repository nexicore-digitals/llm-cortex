package spawn

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	uuid "github.com/google/uuid"
)

type ShellSession struct {
	ID        string
	Cmd       *exec.Cmd
	Stdin     io.WriteCloser
	Stdout    io.ReadCloser
	CreatedAt time.Time
}

var mu sync.Mutex

func NewShell(sessions map[string]*ShellSession) (string, error) {
	cmd := exec.Command("bash", "-i")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	id := uuid.New().String()

	session := &ShellSession{
		ID:        id,
		Cmd:       cmd,
		Stdin:     stdin,
		Stdout:    stdout,
		CreatedAt: time.Now(),
	}

	mu.Lock()
	sessions[session.ID] = session
	mu.Unlock()

	fmt.Printf("🧠 New shell started: %s\n", id)
	return id, nil
}

// Helper function to send a command to a specific session ID
func SendCommand(sessions map[string]*ShellSession, sessionID string, command string) error {
	mu.Lock()
	session, ok := sessions[sessionID]
	mu.Unlock()

	if !ok {
		return fmt.Errorf("shell session %s not found", sessionID)
	}

	// Append a newline to ensure the shell executes the command
	commandWithNewline := command + "\n"

	// Write the command bytes to the shell's input pipe
	_, err := session.Stdin.Write([]byte(commandWithNewline))

	return err
}

// This is the function that processes the raw byte output from the shell.
// In a real application, this is where you'd integrate with a WebSocket
// or a logging system.
func OutputHandler(output []byte, sessionID string) {
	// Convert the raw bytes to a string for display
	outputString := string(output)

	// Print the output clearly labeled with the session ID
	fmt.Printf("[Output from Session %s]: %s", sessionID, outputString)
}

// StartReading launches a goroutine to continuously read output
// from the shell's Stdout pipe until the pipe is closed (shell exits).
// outputHandler is a function that processes the raw byte output.
func StartReading(session *ShellSession, outputHandler func(output []byte, sessionID string)) {
	// Launch the dedicated reader goroutine
	go func() {
		defer session.Stdout.Close()

		// Use a buffer for efficient reading
		buf := make([]byte, 1024)

		for {
			// Read blocks until data is available or the pipe closes
			n, err := session.Stdout.Read(buf)

			if n > 0 {
				// Pass the read bytes to the handler for processing/logging/sending
				outputHandler(buf[:n], session.ID)
			}

			if err != nil {
				// io.EOF is expected when the shell exits normally
				if err != io.EOF {
					fmt.Printf("Error reading from session %s: %v\n", session.ID, err)
				} else {
					fmt.Printf("Shell session %s finished (EOF).\n", session.ID)
				}
				// Once the pipe closes (EOF or error), the goroutine exits
				return
			}
		}
	}()
}

// Helper function to close a session
func CloseSession(sessions map[string]*ShellSession, sessionID string) error {
	mu.Lock()
	session, ok := sessions[sessionID]
	delete(sessions, sessionID) // Remove from the map
	mu.Unlock()

	if !ok {
		return fmt.Errorf("shell session %s not found", sessionID)
	}

	// 1. Send the 'exit' command to the shell
	session.Stdin.Write([]byte("exit\n"))
	session.Stdin.Close() // Close the input pipe

	// 2. Wait for the command to finish and release resources
	// This blocks until the shell process terminates
	return session.Cmd.Wait()
}

func GetSession(sessions map[string]*ShellSession, sessionID string) (*ShellSession, bool) {
	mu.Lock()
	session, ok := sessions[sessionID]
	mu.Unlock()
	return session, ok
}

func IsRunning(session *ShellSession) bool {
	// Check if the underlying process is still running.
	// If Cmd.ProcessState is nil, it usually means the command is still running.
	return session.Cmd.ProcessState == nil || !session.Cmd.ProcessState.Exited()
}
