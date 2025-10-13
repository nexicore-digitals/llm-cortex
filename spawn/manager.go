package spawn

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	uuid "github.com/google/uuid"
)

// ShellSession represents an active, interactive shell process.
type ShellSession struct {
	ID        string
	Cmd       *exec.Cmd
	Stdin     io.WriteCloser
	Stdout    io.ReadCloser
	CreatedAt time.Time
	OutputBuf bytes.Buffer
}

var (
	mu       sync.Mutex
	sessions = make(map[string]*ShellSession)
)

// NewShell creates, starts, and registers a new interactive bash session.
// It returns the unique session ID for future interactions.
func NewShell() (string, error) {
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
		OutputBuf: *bytes.NewBuffer([]byte{}),
	}

	mu.Lock()
	sessions[session.ID] = session
	mu.Unlock()

	fmt.Printf("🧠 New shell started: %s\n", id)
	return id, nil
}

// SendCommand writes a command string to the Stdin of a specific shell session.
// The command should not include a newline character, as it is appended automatically.
func SendCommand(sessionID string, command string) error {
	mu.Lock()
	session, ok := sessions[sessionID] // Use package-level sessions
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

// OutputHandler is a default handler that processes raw byte output from the shell.
// It appends the output to the session's buffer and prints it to the console.
func OutputHandler(output []byte, sessionID string, session *ShellSession) {
	// Convert the raw bytes to a string for display
	mu.Lock()
	session.OutputBuf.Write(output)
	mu.Unlock()

	// Print the output clearly labeled with the session ID
	fmt.Printf("[Output from Session %s]: %s", sessionID, string(output))
}

// StartReading launches a goroutine to continuously read from the shell's Stdout pipe.
// It calls the provided outputHandler for each chunk of data read until the pipe is closed.
func StartReading(sessionID string, outputHandler func(output []byte, sessionID string, session *ShellSession)) error {
	// Use a buffer for efficient reading
	mu.Lock()
	session, ok := sessions[sessionID]
	mu.Unlock()

	if !ok {
		return fmt.Errorf("session %s not found for starting reader", sessionID)
	}

	buf := make([]byte, 1024)
	// Launch the dedicated reader goroutine
	go func() {
		defer session.Stdout.Close()
		for {
			// Read blocks until data is available or the pipe closes
			n, err := session.Stdout.Read(buf)
			if n > 0 {
				// Pass the read bytes to the handler for processing/logging/sending
				outputHandler(buf[:n], session.ID, session)
			}

			if err != nil {
				// io.EOF is expected when the shell exits normally
				if err != io.EOF {
					fmt.Printf("Error reading from session %s: %v\n", session.ID, err)
				} else {
					fmt.Printf("Shell session %s finished (EOF).\n", session.ID)
					return
				}
				// Once the pipe closes (EOF or error), the goroutine exits
				return
			}
		}
	}()
	return nil
}

// CloseSession sends an 'exit' command to the shell, closes its pipes,
// and waits for the process to terminate, releasing all resources.
func CloseSession(sessionID string) error {
	mu.Lock()
	session, ok := sessions[sessionID] // Use package-level sessions
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

// GetSession safely retrieves a session by its ID.
func GetSession(sessionID string) (*ShellSession, bool) {
	mu.Lock()
	session, ok := sessions[sessionID] // Use package-level sessions
	mu.Unlock()
	return session, ok
}

// IsRunning checks if the underlying process for a session is still active.
func IsRunning(sessionID string) bool {
	session, ok := GetSession(sessionID)
	if !ok {
		return false
	}

	// Check if the underlying process is still running.
	// If Cmd.ProcessState is nil, it usually means the command is still running.
	return session.Cmd.ProcessState == nil || !session.Cmd.ProcessState.Exited()
}
