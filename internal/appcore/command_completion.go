package appcore

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// handleCommandTabComplete handles Tab completion in command mode.
// Supports path completion for commands like :cd, :e, :edit
func (s *appState) handleCommandTabComplete() {
	cmd := s.cmdText

	// Check if we're already in completion mode - cycle to next
	if s.cmdCompletionActive && len(s.cmdCompletionItems) > 0 {
		if s.shiftPressed {
			// Shift+Tab: go backward
			s.cmdCompletionIndex--
			if s.cmdCompletionIndex < 0 {
				s.cmdCompletionIndex = len(s.cmdCompletionItems) - 1
			}
		} else {
			// Tab: go forward
			s.cmdCompletionIndex = (s.cmdCompletionIndex + 1) % len(s.cmdCompletionItems)
		}
		// Update command text with selected completion
		s.cmdText = s.cmdCompletionPrefix + s.cmdCompletionItems[s.cmdCompletionIndex]
		return
	}

	// Parse the command to find what we're completing
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return
	}

	cmdName := strings.ToLower(fields[0])

	// Commands that support path completion
	pathCommands := map[string]bool{
		"cd":      true,
		"e":       true,
		"edit":    true,
		"w":       true,
		"write":   true,
		"saveas":  true,
		"install": true,
	}

	if !pathCommands[cmdName] {
		return
	}

	// Get the path portion (everything after the command)
	var pathArg string
	var prefix string
	if len(fields) > 1 {
		// Find where the path starts in the original string
		cmdLen := len(fields[0])
		rest := strings.TrimLeft(cmd[cmdLen:], " ")
		pathArg = rest
		prefix = cmd[:len(cmd)-len(rest)]
	} else if strings.HasSuffix(cmd, " ") {
		// Command with trailing space, complete from current directory
		pathArg = ""
		prefix = cmd
	} else {
		// No path yet
		return
	}

	// Expand tilde
	expandedPath := expandTilde(pathArg)

	// Get completions
	completions := getPathCompletions(expandedPath)
	if len(completions) == 0 {
		s.status = "No completions"
		return
	}

	// If only one completion, use it directly
	if len(completions) == 1 {
		s.cmdText = prefix + completions[0]
		s.cmdCompletionActive = false
		return
	}

	// Multiple completions - enter completion mode
	s.cmdCompletionActive = true
	s.cmdCompletionItems = completions
	s.cmdCompletionIndex = 0
	s.cmdCompletionPrefix = prefix
	s.cmdText = prefix + completions[0]
	s.status = "Tab to cycle completions"
}

// expandTilde expands ~ to the user's home directory
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return home + path[1:]
		}
	}
	return path
}

// getPathCompletions returns possible path completions for the given partial path
func getPathCompletions(partial string) []string {
	var dir, prefix string

	if partial == "" {
		// Complete from current directory
		dir = "."
		prefix = ""
	} else if strings.HasSuffix(partial, "/") {
		// Completing inside a directory
		dir = partial
		prefix = ""
	} else {
		// Completing a partial name
		dir = filepath.Dir(partial)
		prefix = filepath.Base(partial)
	}

	// Read directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var completions []string
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless prefix starts with .
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(prefix, ".") {
			continue
		}

		// Check if name matches prefix
		if prefix != "" && !strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			continue
		}

		// Build the full path for completion
		var completion string
		if dir == "." {
			completion = name
		} else if strings.HasSuffix(partial, "/") {
			completion = partial + name
		} else {
			completion = filepath.Join(filepath.Dir(partial), name)
		}

		// Add trailing slash for directories
		if entry.IsDir() {
			completion += "/"
		}

		// Preserve tilde if original had it
		if strings.HasPrefix(partial, "~") {
			home, _ := os.UserHomeDir()
			if strings.HasPrefix(completion, home) {
				completion = "~" + completion[len(home):]
			}
		}

		completions = append(completions, completion)
	}

	// Sort completions
	sort.Strings(completions)

	return completions
}

// cancelCommandCompletion cancels command completion state
func (s *appState) cancelCommandCompletion() {
	s.cmdCompletionActive = false
	s.cmdCompletionItems = nil
	s.cmdCompletionIndex = 0
	s.cmdCompletionPrefix = ""
}
