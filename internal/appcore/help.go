package appcore

import (
	"fmt"
	"strings"

	"gioui.org/io/key"
)

// generateHelpText creates formatted help text from keybindings
func generateHelpText() string {
	var sb strings.Builder

	sb.WriteString("═══════════════════════════════════════════════════════════\n")
	sb.WriteString("                   VEM HELP - KEYBINDINGS                  \n")
	sb.WriteString("═══════════════════════════════════════════════════════════\n")
	sb.WriteString("\n")
	sb.WriteString("Press / to search, :q to close\n")
	sb.WriteString("\n")

	// Global Keybindings
	sb.WriteString("GLOBAL KEYBINDINGS (work in all modes)\n")
	sb.WriteString("───────────────────────────────────────────────────────────\n")
	appendGlobalKeybindings(&sb)
	sb.WriteString("\n")

	// Mode-specific keybindings
	sb.WriteString("NORMAL MODE\n")
	sb.WriteString("───────────────────────────────────────────────────────────\n")
	appendModeKeybindings(&sb, modeNormal)
	sb.WriteString("\n")

	sb.WriteString("INSERT MODE\n")
	sb.WriteString("───────────────────────────────────────────────────────────\n")
	appendModeKeybindings(&sb, modeInsert)
	sb.WriteString("\n")

	sb.WriteString("VISUAL MODE\n")
	sb.WriteString("───────────────────────────────────────────────────────────\n")
	appendModeKeybindings(&sb, modeVisual)
	sb.WriteString("\n")

	sb.WriteString("EXPLORER MODE\n")
	sb.WriteString("───────────────────────────────────────────────────────────\n")
	appendModeKeybindings(&sb, modeExplorer)
	sb.WriteString("\n")

	sb.WriteString("TERMINAL MODE\n")
	sb.WriteString("───────────────────────────────────────────────────────────\n")
	appendModeKeybindings(&sb, modeTerminal)
	sb.WriteString("\n")

	sb.WriteString("COMMANDS\n")
	sb.WriteString("───────────────────────────────────────────────────────────\n")
	appendCommands(&sb)
	sb.WriteString("\n")

	sb.WriteString("SPECIAL SEQUENCES\n")
	sb.WriteString("───────────────────────────────────────────────────────────\n")
	appendSpecialSequences(&sb)

	return sb.String()
}

// appendGlobalKeybindings adds global keybinding help
func appendGlobalKeybindings(sb *strings.Builder) {
	bindings := []struct {
		keys string
		desc string
	}{
		{"Ctrl+T", "Toggle file explorer"},
		{"Ctrl+H", "Focus file explorer"},
		{"Ctrl+L", "Focus editor"},
		{"Ctrl+F", "Open fuzzy finder"},
		{"Ctrl+U", "Undo last edit"},
		{"Ctrl+C", "Copy current line (NORMAL mode)"},
		{"Ctrl+P", "Paste from clipboard"},
		{"Ctrl+X", "Close pane/buffer"},
		{"Ctrl+`", "Open/toggle terminal"},
		{"Alt+h", "Focus pane left"},
		{"Alt+j", "Focus pane down"},
		{"Alt+k", "Focus pane up"},
		{"Alt+l", "Focus pane right"},
		{"Shift+Tab", "Cycle to next pane"},
		{"Shift+Enter", "Toggle fullscreen (NORMAL mode)"},
	}

	for _, b := range bindings {
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", b.keys, b.desc))
	}
}

// appendModeKeybindings adds mode-specific keybinding help
func appendModeKeybindings(sb *strings.Builder, mode mode) {
	bindings, exists := modeKeybindings[mode]
	if !exists {
		sb.WriteString("  No keybindings defined\n")
		return
	}

	for _, binding := range bindings {
		keys := formatKeybinding(binding)
		desc := actionDescription(binding.Action)
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", keys, desc))
	}
}

// appendCommands adds command help
func appendCommands(sb *strings.Builder) {
	commands := []struct {
		cmd  string
		desc string
	}{
		{":q", "Close current pane/buffer"},
		{":q!", "Force close (discard changes)"},
		{":qa", "Quit entire application"},
		{":qa!", "Force quit (discard all changes)"},
		{":w", "Save current buffer"},
		{":w <file>", "Save as <file>"},
		{":wq", "Save and close"},
		{":e <file>", "Open file for editing"},
		{":bn", "Next buffer"},
		{":bp", "Previous buffer"},
		{":bd", "Delete buffer"},
		{":ls", "List all buffers"},
		{":ex", "Toggle file explorer"},
		{":cd <path>", "Change working directory"},
		{":pwd", "Print working directory"},
		{":term", "Open embedded terminal"},
		{":help", "Show this help"},
	}

	for _, c := range commands {
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", c.cmd, c.desc))
	}
}

// appendSpecialSequences adds special sequence help
func appendSpecialSequences(sb *strings.Builder) {
	sequences := []struct {
		seq  string
		desc string
	}{
		{"gg", "Jump to first line"},
		{"G", "Jump to last line"},
		{"<count>G", "Jump to line <count> (e.g., 42G)"},
		{"<count>j/k", "Move <count> lines (e.g., 5j)"},
		{"dd", "Delete current line"},
		{"<count>dd", "Delete line <count>"},
		{"zz", "Center cursor in viewport"},
		{"zt", "Cursor to top of viewport"},
		{"zb", "Cursor to bottom of viewport"},
		{"Ctrl+S v", "Split vertically"},
		{"Ctrl+S h", "Split horizontally"},
		{"Ctrl+S =", "Equalize panes"},
		{"Ctrl+S o", "Zoom/unzoom pane"},
	}

	for _, s := range sequences {
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", s.seq, s.desc))
	}
}

// formatKeybinding formats a keybinding for display
func formatKeybinding(binding KeyBinding) string {
	var parts []string

	// Format modifiers
	if binding.Modifiers.Contain(key.ModCtrl) {
		parts = append(parts, "Ctrl")
	}
	if binding.Modifiers.Contain(key.ModShift) {
		parts = append(parts, "Shift")
	}
	if binding.Modifiers.Contain(key.ModAlt) {
		parts = append(parts, "Alt")
	}

	// Format key name
	keyName := formatKeyName(binding.Key)

	if len(parts) > 0 {
		return strings.Join(parts, "+") + "+" + keyName
	}
	return keyName
}

// formatKeyName formats a key name for display
func formatKeyName(k key.Name) string {
	switch k {
	case key.NameEscape:
		return "Esc"
	case key.NameReturn, key.NameEnter:
		return "Enter"
	case key.NameLeftArrow:
		return "←"
	case key.NameRightArrow:
		return "→"
	case key.NameUpArrow:
		return "↑"
	case key.NameDownArrow:
		return "↓"
	case key.NameDeleteBackward:
		return "Backspace"
	case key.NameDeleteForward:
		return "Delete"
	case key.NameSpace:
		return "Space"
	case key.NameTab:
		return "Tab"
	default:
		return string(k)
	}
}

// actionDescription returns a human-readable description for an action
func actionDescription(action Action) string {
	descriptions := map[Action]string{
		ActionNone:               "No action",
		ActionToggleExplorer:     "Toggle file explorer",
		ActionFocusExplorer:      "Focus explorer",
		ActionFocusEditor:        "Focus editor",
		ActionToggleFullscreen:   "Toggle fullscreen",
		ActionEnterInsert:        "Enter INSERT mode",
		ActionEnterVisualChar:    "Enter VISUAL (char) mode",
		ActionEnterVisualLine:    "Enter VISUAL (line) mode",
		ActionEnterDelete:        "Enter DELETE mode",
		ActionEnterCommand:       "Enter COMMAND mode",
		ActionEnterExplorer:      "Enter EXPLORER mode",
		ActionExitMode:           "Exit current mode",
		ActionMoveLeft:           "Move cursor left",
		ActionMoveRight:          "Move cursor right",
		ActionMoveUp:             "Move cursor up",
		ActionMoveDown:           "Move cursor down",
		ActionJumpLineStart:      "Jump to line start",
		ActionJumpLineEnd:        "Jump to line end",
		ActionWordForward:        "Move to next word",
		ActionWordBackward:       "Move to previous word",
		ActionWordEnd:            "Move to end of word",
		ActionInsertNewline:      "Insert newline",
		ActionInsertSpace:        "Insert space",
		ActionInsertTab:          "Insert tab",
		ActionDeleteBackward:     "Delete backward",
		ActionDeleteForward:      "Delete forward",
		ActionUndo:               "Undo last edit",
		ActionCopySelection:      "Copy selection",
		ActionDeleteSelection:    "Delete selection",
		ActionPasteClipboard:     "Paste clipboard",
		ActionCopyLine:           "Copy current line",
		ActionPaste:              "Paste at cursor",
		ActionOpenNode:           "Open file/folder",
		ActionCollapseNode:       "Collapse folder",
		ActionExpandNode:         "Expand folder",
		ActionRenameFile:         "Rename file",
		ActionDeleteFile:         "Delete file",
		ActionCreateFile:         "Create new file",
		ActionNavigateUp:         "Navigate to parent dir",
		ActionEnterSearch:        "Enter search mode",
		ActionNextMatch:          "Next search match",
		ActionPrevMatch:          "Previous search match",
		ActionClearSearch:        "Clear search",
		ActionOpenFuzzyFinder:    "Open fuzzy finder",
		ActionFuzzyFinderConfirm: "Confirm selection",
		ActionScrollToCenter:     "Center viewport",
		ActionScrollToTop:        "Scroll to top",
		ActionScrollToBottom:     "Scroll to bottom",
		ActionScrollLineUp:       "Scroll up one line",
		ActionScrollLineDown:     "Scroll down one line",
		ActionSplitVertical:      "Split vertically",
		ActionSplitHorizontal:    "Split horizontally",
		ActionPaneFocusLeft:      "Focus pane left",
		ActionPaneFocusRight:     "Focus pane right",
		ActionPaneFocusUp:        "Focus pane up",
		ActionPaneFocusDown:      "Focus pane down",
		ActionPaneCycleNext:      "Cycle to next pane",
		ActionPaneClose:          "Close pane",
		ActionPaneEqualize:       "Equalize panes",
		ActionPaneZoomToggle:     "Toggle pane zoom",
		ActionOpenTerminal:       "Open terminal",
		ActionTerminalExit:       "Exit terminal mode",
	}

	if desc, exists := descriptions[action]; exists {
		return desc
	}
	return "Unknown action"
}
