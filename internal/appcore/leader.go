package appcore

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gioui.org/io/key"
)

// LeaderBinding represents a single leader key binding
type LeaderBinding struct {
	Name    string `toml:"name"`
	Keys    string `toml:"keys"`
	Action  string `toml:"action"`
	Command string `toml:"command"`
}

// LeaderConfig represents the leader keybindings configuration file
type LeaderConfig struct {
	Keybinds []LeaderBinding `toml:"keybind"`
}

// actionMap maps action names from config to Action constants
var actionMap = map[string]Action{
	"lsp_goto_definition": ActionLSPGotoDefinition,
	"lsp_find_references": ActionLSPReferences,
	"lsp_hover":           ActionLSPHover,
	"lsp_rename":          ActionLSPRename,
	"lsp_code_action":     ActionLSPCodeAction,
	"lsp_format":          ActionLSPFormat,
	"fuzzy_finder":        ActionOpenFuzzyFinder,
	"toggle_explorer":     ActionToggleExplorer,
	"split_vertical":      ActionSplitVertical,
	"split_horizontal":    ActionSplitHorizontal,
	"pane_close":          ActionPaneClose,
	"next_buffer":         ActionNextBuffer,
	"prev_buffer":         ActionPrevBuffer,
}

// commandMap maps action names that need to be executed as commands
var commandMap = map[string]string{
	"save_file":    "w",
	"close_buffer": "bd",
}

// loadLeaderConfig loads the leader keybindings from the config file
func (s *appState) loadLeaderConfig() {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "vem", "keybindings.toml")

	// Try user config first
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No user config, use defaults
		s.leaderBarBindings = getDefaultBindings()
		return
	}

	var config LeaderConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		s.status = "Error loading keybindings.toml: " + err.Error()
		s.leaderBarBindings = getDefaultBindings()
		return
	}

	s.leaderBarBindings = config.Keybinds
}

// getDefaultBindings returns a minimal set of default leader bindings
func getDefaultBindings() []LeaderBinding {
	return []LeaderBinding{
		{Name: "Go to Definition", Keys: "gd", Action: "lsp_goto_definition"},
		{Name: "Find References", Keys: "gr", Action: "lsp_find_references"},
		{Name: "Hover Info", Keys: "K", Action: "lsp_hover"},
		{Name: "Save File", Keys: "w", Action: "save_file"},
		{Name: "Fuzzy Find", Keys: "ff", Action: "fuzzy_finder"},
		{Name: "Toggle Explorer", Keys: "e", Action: "toggle_explorer"},
		{Name: "Split Vertical", Keys: "sv", Action: "split_vertical"},
		{Name: "Split Horizontal", Keys: "sh", Action: "split_horizontal"},
		{Name: "Close Pane", Keys: "x", Action: "pane_close"},
	}
}

// handleSpaceKey handles space key press in normal mode for leader bar detection
func (s *appState) handleSpaceKey() bool {
	now := time.Now()
	if now.Sub(s.lastSpaceTime) < 300*time.Millisecond {
		// Double-space detected, enter leader bar
		s.enterLeaderBar()
		s.lastSpaceTime = time.Time{} // Reset to prevent triple-space issues
		return true
	}
	s.lastSpaceTime = now
	return false
}

// enterLeaderBar activates the leader bar popup
func (s *appState) enterLeaderBar() {
	s.leaderBarActive = true
	s.leaderBarSequence = ""
	s.leaderBarMatches = s.leaderBarBindings // Show all bindings initially
	s.leaderBarIndex = 0
	s.status = "LEADER: Press keys for command, Esc to cancel"
}

// exitLeaderBar closes the leader bar popup
func (s *appState) exitLeaderBar() {
	s.leaderBarActive = false
	s.leaderBarSequence = ""
	s.leaderBarMatches = nil
	s.leaderBarIndex = 0
	s.status = ""
}

// handleLeaderKey processes key input while the leader bar is active
func (s *appState) handleLeaderKey(keyName string) {
	if keyName == "Escape" {
		s.exitLeaderBar()
		return
	}

	// Handle special keys
	switch keyName {
	case "↑", "Up":
		if s.leaderBarIndex > 0 {
			s.leaderBarIndex--
		}
		return
	case "↓", "Down":
		if s.leaderBarIndex < len(s.leaderBarMatches)-1 {
			s.leaderBarIndex++
		}
		return
	case "Return", "Enter":
		// Execute selected binding
		if len(s.leaderBarMatches) > 0 && s.leaderBarIndex < len(s.leaderBarMatches) {
			s.executeLeaderBinding(s.leaderBarMatches[s.leaderBarIndex])
			s.exitLeaderBar()
		}
		return
	case "DeleteBackward", "Backspace":
		// Remove last character from sequence
		if len(s.leaderBarSequence) > 0 {
			s.leaderBarSequence = s.leaderBarSequence[:len(s.leaderBarSequence)-1]
			s.leaderBarMatches = s.getMatchingBindings(s.leaderBarSequence)
			s.leaderBarIndex = 0
		}
		return
	}

	// Ignore modifier keys and special keys
	if len(keyName) > 1 && keyName != "Space" {
		return
	}

	// Convert Space to actual space character for matching
	if keyName == "Space" {
		keyName = " "
	}

	// Append to sequence
	s.leaderBarSequence += keyName

	// Filter matches
	s.leaderBarMatches = s.getMatchingBindings(s.leaderBarSequence)
	s.leaderBarIndex = 0

	// If exactly one match and sequence is complete, execute immediately
	if len(s.leaderBarMatches) == 1 && s.leaderBarMatches[0].Keys == s.leaderBarSequence {
		s.executeLeaderBinding(s.leaderBarMatches[0])
		s.exitLeaderBar()
		return
	}

	// If no matches, show error and close
	if len(s.leaderBarMatches) == 0 {
		s.status = "No matching command for: " + s.leaderBarSequence
		s.exitLeaderBar()
	}
}

// getMatchingBindings returns bindings that match the given sequence prefix
func (s *appState) getMatchingBindings(sequence string) []LeaderBinding {
	if sequence == "" {
		return s.leaderBarBindings
	}

	var matches []LeaderBinding
	for _, b := range s.leaderBarBindings {
		if strings.HasPrefix(b.Keys, sequence) {
			matches = append(matches, b)
		}
	}
	return matches
}

// executeLeaderBinding executes the action or command from a leader binding
func (s *appState) executeLeaderBinding(binding LeaderBinding) {
	// If it's a shell command, execute it
	if binding.Command != "" {
		s.executeShellCommand(binding.Command)
		return
	}

	// If it's a Vem action, look it up and execute
	if binding.Action != "" {
		// First check if it maps to an Action constant
		if action, ok := actionMap[binding.Action]; ok {
			s.executeAction(action, key.Event{})
			return
		}
		// Then check if it maps to a command
		if cmd, ok := commandMap[binding.Action]; ok {
			s.cmdText = cmd
			s.executeCommandLine()
			return
		}
		s.status = "Unknown action: " + binding.Action
	}
}

// executeShellCommand runs a shell command (uses the :! command internally)
func (s *appState) executeShellCommand(cmd string) {
	// Use the existing command line execution with ! prefix
	s.cmdText = "!" + cmd
	s.executeCommandLine()
}

// reloadLeaderConfig reloads the leader keybindings from the config file
func (s *appState) reloadLeaderConfig() {
	s.loadLeaderConfig()
	s.status = "Leader keybindings reloaded"
}
