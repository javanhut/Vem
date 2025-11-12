//go:build windows

package terminal

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// startPTY creates ConPTY and starts shell (Windows implementation)
func (t *Terminal) startPTY() error {
	// creack/pty handles Windows ConPTY automatically
	ptyFile, ttyFile, err := pty.Open()
	if err != nil {
		return err
	}

	t.pty = ptyFile

	// Set PTY size
	if err := pty.Setsize(ptyFile, &pty.Winsize{
		Rows: uint16(t.height),
		Cols: uint16(t.width),
	}); err != nil {
		ptyFile.Close()
		return err
	}

	// Create command
	cmd := exec.Command(t.shell, t.args...)
	cmd.Dir = t.workingDir
	cmd.Env = t.getEnvironment()

	// Connect to TTY
	cmd.Stdin = ttyFile
	cmd.Stdout = ttyFile
	cmd.Stderr = ttyFile

	// Start process
	if err := cmd.Start(); err != nil {
		ptyFile.Close()
		ttyFile.Close()
		return err
	}

	// Close TTY in parent process
	ttyFile.Close()

	t.cmd = cmd

	// Wait for process in goroutine
	go func() {
		cmd.Wait()
		log.Println("[TERMINAL] Shell process exited")
	}()

	return nil
}

// Resize updates ConPTY window size
func (t *Terminal) Resize(width, height int) error {
	t.mu.Lock()
	t.width = width
	t.height = height
	ptyFile := t.pty
	t.mu.Unlock()

	if ptyFile == nil {
		return fmt.Errorf("ConPTY not initialized")
	}

	// Resize screen buffer
	if t.screen != nil {
		t.screen.Resize(width, height)
	}

	return pty.Setsize(ptyFile, &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
	})
}

// DefaultShell returns default shell for Windows
func DefaultShell() string {
	// Try PowerShell Core first
	if _, err := exec.LookPath("pwsh.exe"); err == nil {
		return "pwsh.exe"
	}

	// Try Windows PowerShell
	if _, err := exec.LookPath("powershell.exe"); err == nil {
		return "powershell.exe"
	}

	// Fallback to cmd.exe
	if comspec := os.Getenv("COMSPEC"); comspec != "" {
		return comspec
	}

	return "cmd.exe"
}

// DefaultArgs returns default shell args for Windows
func DefaultArgs() []string {
	return []string{} // No args needed
}
