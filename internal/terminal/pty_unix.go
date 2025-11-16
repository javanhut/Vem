//go:build unix

package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
)

// startPTY creates PTY and starts shell (Unix implementation)
func (t *Terminal) startPTY() error {
	// Create PTY
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
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true, // Create new session
		Setctty: true, // Set controlling terminal
	}

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
	}()

	return nil
}

// Resize updates PTY window size
func (t *Terminal) Resize(width, height int) error {
	t.mu.Lock()
	t.width = width
	t.height = height
	ptyFile := t.pty
	t.mu.Unlock()

	if ptyFile == nil {
		return fmt.Errorf("PTY not initialized")
	}

	// Resize screen buffer
	if t.screen != nil {
		t.screen.Resize(width, height)
	}

	// Resize vt10x emulator
	if t.vt != nil {
		t.vt.Resize(width, height)
	}

	return pty.Setsize(ptyFile, &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
	})
}

// DefaultShell returns default shell for Unix
func DefaultShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}
	return "/bin/sh"
}

// DefaultArgs returns default shell args for Unix
func DefaultArgs() []string {
	return []string{"-i"} // Interactive shell
}
