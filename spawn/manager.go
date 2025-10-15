package spawn

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	uuid "github.com/google/uuid"
)

// ShellSession represents an active, interactive shell process.
type ShellSession struct {
	ID        string
	Cmd       *exec.Cmd
	Stdin     io.WriteCloser
	Stdout    io.ReadCloser // Pipe for standard output
	Stderr    io.ReadCloser // Pipe for standard error
	CreatedAt time.Time
	OutputBuf bytes.Buffer // Buffer for stdout
	StderrBuf bytes.Buffer // Buffer for stderr
	mu        sync.Mutex   // Mutex to protect this session's buffers
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
		OutputBuf: *bytes.NewBuffer(nil),
		// Stderr is not captured for basic shells, only for command shells.
		StderrBuf: *bytes.NewBuffer(nil),
	}

	mu.Lock()
	sessions[session.ID] = session
	mu.Unlock()

	fmt.Printf("🧠 New shell started: %s\n", id)
	return id, nil
}

// NewShellWithCommand creates, starts, and registers a new interactive session with a custom command.
// It is used for launching persistent Python model scripts.
// It returns the unique session ID for future interactions.
func NewShellWithCommand(command ...string) (string, error) {
	if len(command) == 0 {
		return "", fmt.Errorf("NewShellWithCommand requires a command to execute")
	}
	cmd := exec.Command(command[0], command[1:]...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	// Create a separate pipe for stderr to distinguish errors from normal output.
	stderr, err := cmd.StderrPipe()
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
		Stderr:    stderr,
		CreatedAt: time.Now(),
		OutputBuf: *bytes.NewBuffer(nil),
		StderrBuf: *bytes.NewBuffer(nil),
	}

	mu.Lock()
	sessions[id] = session
	mu.Unlock()

	fmt.Printf("🧠 New command shell started: %s with command %v\n", id, command)
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

// SendCommandAndWait sends a command and waits for a specific delimiter in the response.
// It polls the session's output buffer until the delimiter is found or a timeout occurs.
func SendCommandAndWait(sessionID string, command string, delimiter string) (string, error) {
	mu.Lock()
	session, ok := sessions[sessionID]
	mu.Unlock()

	if !ok {
		return "", fmt.Errorf("shell session %s not found", sessionID)
	}

	// Clear the buffer before sending a new command to ensure we only capture the new output.
	mu.Lock()
	session.OutputBuf.Reset()
	mu.Unlock()

	// Append a newline to ensure the shell executes the command
	commandWithNewline := command + "\n"

	// Write the command bytes to the shell's input pipe
	_, err := session.Stdin.Write([]byte(commandWithNewline))
	if err != nil {
		return "", err
	}

	// Wait for the response by polling the buffer until the delimiter is found.
	// This is a simple polling mechanism. For high-performance scenarios,
	// a condition variable or channel-based approach might be more efficient.
	timeout := time.After(30 * time.Second) // 30-second timeout for the command to respond
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-timeout:
			return "", fmt.Errorf("timed out waiting for response delimiter: %s", delimiter)
		case <-tick.C:
			if output, ok := checkBufferForDelimiter(session, delimiter); ok {
				return output, nil
			}
		}
	}
}

// OutputHandler is a default handler that processes raw byte output from the shell.
// It appends the output to the session's buffer and prints it to the console.
func OutputHandler(output []byte, sessionID string, session *ShellSession) {
	// Convert the raw bytes to a string for display
	session.mu.Lock()
	session.OutputBuf.Write(output)
	session.mu.Unlock()

	// Print the output clearly labeled with the session ID
	fmt.Printf("[Output from Session %s]: %s", sessionID, string(output))
}

// ErrorOutputHandler is a handler that processes raw byte output from the shell's stderr.
// It appends the output to the session's StderrBuf and prints it to the console as an error.
func ErrorOutputHandler(output []byte, sessionID string, session *ShellSession) {
	session.mu.Lock()
	session.StderrBuf.Write(output)
	session.mu.Unlock()

	// Print the error output clearly labeled with the session ID
	fmt.Printf("[Error from Session %s]: %s", sessionID, string(output))
}

// StartReading launches a goroutine to continuously read from the shell's Stdout pipe.
// It calls the provided outputHandler for each chunk of data read until the pipe is closed.
func StartReading(sessionID string, stdoutHandler func(output []byte, sessionID string, session *ShellSession), stderrHandler func(output []byte, sessionID string, session *ShellSession)) error {
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
				stdoutHandler(buf[:n], session.ID, session)
			}

			if err != nil {
				// io.EOF is expected when the shell exits normally
				if err != io.EOF {
					fmt.Printf("Error reading from session %s: %v\n", session.ID, err)
				} else {
					fmt.Printf("Shell session %s finished (EOF).\n", session.ID)
					return
					// EOF is expected when the shell exits normally.
					// The process finishing will be logged by the stderr handler if it exits with an error.
				}
				// Once the pipe closes (EOF or error), the goroutine exits
				return
			}
		}
	}()

	// Launch a dedicated reader goroutine for stderr if it exists
	if session.Stderr != nil {
		go func() {
			defer session.Stderr.Close()
			errBuf := make([]byte, 1024)
			for {
				n, err := session.Stderr.Read(errBuf)
				if n > 0 {
					stderrHandler(errBuf[:n], session.ID, session)
				}

				if err != nil {
					if err != io.EOF {
						fmt.Printf("Error reading from session stderr %s: %v\n", session.ID, err)
					} else {
						// EOF is expected when the process closes its stderr.
					}
					return
				}
			}
		}()
	}
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

	// 1. Send a JSON exit command to gracefully shut down the Python script.
	// The script is designed to exit when it receives this command.
	// This avoids pipe closing issues that can affect subsequent process spawns.
	exitCommand := `{"command": "exit"}` + "\n"
	session.Stdin.Write([]byte(exitCommand))

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

// WaitForString polls the session's output buffer until a specific string is found
// or a timeout is reached. This is useful for waiting for a "Ready" signal from a script.
func WaitForString(sessionID string, target string, timeout time.Duration) error {
	startTime := time.Now()
	for {
		if time.Since(startTime) > timeout {
			mu.Lock()
			session, ok := sessions[sessionID]
			var output, stderrOutput string
			if ok {
				output = session.OutputBuf.String()
				stderrOutput = session.StderrBuf.String()
			}
			mu.Unlock()
			return fmt.Errorf("timed out waiting for string '%s'.\nLast stdout: %s\nLast stderr: %s", target, output, stderrOutput)
		}

		time.Sleep(200 * time.Millisecond) // Poll every 200ms

		mu.Lock()
		session, ok := sessions[sessionID]
		if !ok {
			mu.Unlock()
			return fmt.Errorf("session %s not found while waiting for string", sessionID)
		}
		output := session.OutputBuf.String()
		mu.Unlock()

		if strings.Contains(output, target) {
			return nil
		}
	}
}

func checkBufferForDelimiter(session *ShellSession, delimiter string) (string, bool) {
	session.mu.Lock()
	defer session.mu.Unlock()
	output := session.OutputBuf.String()
	if strings.Contains(output, delimiter) {
		return strings.Split(output, delimiter)[0], true
	}
	return "", false
}
