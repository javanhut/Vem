//go:build darwin

package appcore

import (
	"fmt"
	"os"

	"gioui.org/io/key"
)

// darwinDebug controls debug logging for macOS input
var darwinDebug = os.Getenv("VEM_DEBUG_DARWIN") != ""

func darwinInputLog(format string, args ...interface{}) {
	if darwinDebug {
		fmt.Fprintf(os.Stderr, "[darwin] "+format+"\n", args...)
	}
}

// handleModifierEvent handles modifier key events on macOS.
// Treats Command key as equivalent to Ctrl for keybinding purposes,
// allowing users to use Cmd+key shortcuts (macOS convention) while
// keeping all existing key.ModCtrl bindings working.
func (s *appState) handleModifierEvent(e key.Event) bool {
	if darwinDebug {
		darwinInputLog("handleModifierEvent: Name=%q State=%v Modifiers=%v", e.Name, e.State, e.Modifiers)
	}

	// Track Command key as ctrlPressed (macOS convention)
	if e.Name == key.NameCommand {
		s.ctrlPressed = (e.State == key.Press)
		if darwinDebug {
			darwinInputLog("  -> Command key: ctrlPressed=%v", s.ctrlPressed)
		}
		return true
	}

	// Also track actual Ctrl (for users who prefer Ctrl+key)
	if e.Name == key.NameCtrl {
		s.ctrlPressed = (e.State == key.Press)
		if darwinDebug {
			darwinInputLog("  -> Ctrl key: ctrlPressed=%v", s.ctrlPressed)
		}
		return true
	}

	if e.Name == key.NameShift {
		s.shiftPressed = (e.State == key.Press)
		if darwinDebug {
			darwinInputLog("  -> Shift key: shiftPressed=%v", s.shiftPressed)
		}
		return true
	}

	if e.Name == key.NameAlt {
		if darwinDebug {
			darwinInputLog("  -> Alt key (ignored)")
		}
		return true
	}

	return false
}

// syncModifierState syncs the tracked modifier state before handling character keys.
// On macOS, check both ModCommand and ModCtrl since either can trigger "Ctrl" shortcuts.
func (s *appState) syncModifierState(e key.Event) {
	if darwinDebug {
		darwinInputLog("syncModifierState: Name=%q Modifiers=%v (before: ctrl=%v shift=%v)",
			e.Name, e.Modifiers, s.ctrlPressed, s.shiftPressed)
	}

	// On macOS, Command key should trigger Ctrl shortcuts
	if e.Modifiers.Contain(key.ModCommand) && !s.ctrlPressed {
		s.ctrlPressed = true
		if darwinDebug {
			darwinInputLog("  -> ModCommand detected, set ctrlPressed=true")
		}
	}
	// Also support actual Ctrl key
	if e.Modifiers.Contain(key.ModCtrl) && !s.ctrlPressed {
		s.ctrlPressed = true
		if darwinDebug {
			darwinInputLog("  -> ModCtrl detected, set ctrlPressed=true")
		}
	}
	if e.Modifiers.Contain(key.ModShift) && !s.shiftPressed {
		s.shiftPressed = true
		if darwinDebug {
			darwinInputLog("  -> ModShift detected, set shiftPressed=true")
		}
	}

	if darwinDebug {
		darwinInputLog("  -> after sync: ctrl=%v shift=%v", s.ctrlPressed, s.shiftPressed)
	}
}
