//go:build !windows

package appcore

import (
	"gioui.org/io/key"
)

// handleModifierEvent handles modifier key events on Unix/Linux/macOS.
// On these platforms, Gio correctly sends both Press and Release events
// for modifier keys, and the timing is correct.
func (s *appState) handleModifierEvent(e key.Event) bool {
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
// On Unix/Linux/macOS, ev.Modifiers should contain accurate information, so we use
// it as a fallback in case Press events were missed.
func (s *appState) syncModifierState(e key.Event) {
	// On Unix, ev.Modifiers usually works correctly
	// Use it as a fallback to catch any missed Press events
	if e.Modifiers.Contain(key.ModCtrl) && !s.ctrlPressed {
		s.ctrlPressed = true
	}
	if e.Modifiers.Contain(key.ModShift) && !s.shiftPressed {
		s.shiftPressed = true
	}
}
