//go:build darwin

package appcore

import (
	"gioui.org/io/key"
)

// handleModifierEvent handles modifier key events on macOS.
// Treats Command key as equivalent to Ctrl for keybinding purposes,
// allowing users to use Cmd+key shortcuts (macOS convention) while
// keeping all existing key.ModCtrl bindings working.
func (s *appState) handleModifierEvent(e key.Event) bool {
	// Track Command key as ctrlPressed (macOS convention)
	if e.Name == key.NameCommand {
		s.ctrlPressed = (e.State == key.Press)
		return true
	}

	// Also track actual Ctrl (for users who prefer Ctrl+key)
	if e.Name == key.NameCtrl {
		s.ctrlPressed = (e.State == key.Press)
		return true
	}

	if e.Name == key.NameShift {
		s.shiftPressed = (e.State == key.Press)
		return true
	}

	if e.Name == key.NameAlt {
		return true
	}

	return false
}

// syncModifierState syncs the tracked modifier state before handling character keys.
// On macOS, check both ModCommand and ModCtrl since either can trigger "Ctrl" shortcuts.
func (s *appState) syncModifierState(e key.Event) {
	// On macOS, Command key should trigger Ctrl shortcuts
	if e.Modifiers.Contain(key.ModCommand) && !s.ctrlPressed {
		s.ctrlPressed = true
	}
	// Also support actual Ctrl key
	if e.Modifiers.Contain(key.ModCtrl) && !s.ctrlPressed {
		s.ctrlPressed = true
	}
	if e.Modifiers.Contain(key.ModShift) && !s.shiftPressed {
		s.shiftPressed = true
	}
}
