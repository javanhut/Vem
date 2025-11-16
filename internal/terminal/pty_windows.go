//go:build windows

package terminal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/UserExistsError/conpty"
)

// startPTY creates ConPTY and starts shell (Windows implementation)
func (t *Terminal) startPTY() error {
	// Build command line - Windows ConPTY API requires a single command line string
	commandLine := t.shell
	if len(t.args) > 0 {
		// For PowerShell, join args with spaces
		commandLine = t.shell + " " + strings.Join(t.args, " ")
	}

	// Create ConPTY with proper dimensions
	cpty, err := conpty.Start(
		commandLine,
		conpty.ConPtyDimensions(t.width, t.height),
		conpty.ConPtyWorkDir(t.workingDir),
		conpty.ConPtyEnv(t.getEnvironment()),
	)
	if err != nil {
		return fmt.Errorf("failed to create ConPTY: %w", err)
	}

	// Store ConPTY instance
	t.conpty = &ConPtyWrapper{cpty}

	// Start wait goroutine
	go func() {
		_, _ = cpty.Wait(context.Background())

		// Call onExit callback if set
		if t.onExit != nil {
			t.onExit()
		}
	}()

	return nil
}

// Resize updates ConPTY window size
func (t *Terminal) Resize(width, height int) error {
	t.mu.Lock()
	t.width = width
	t.height = height
	cpty := t.conpty
	t.mu.Unlock()

	if cpty == nil {
		return fmt.Errorf("ConPTY not initialized")
	}

	// Resize screen buffer
	if t.screen != nil {
		t.screen.Resize(width, height)
	}

	return cpty.Resize(width, height)
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
