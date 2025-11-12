//go:build !windows

package appcore

import (
	"log"

	"gioui.org/io/key"
)

// handleModifierEvent handles modifier key events on Unix/Linux/macOS.
// On these platforms, Gio correctly sends both Press and Release events
// for modifier keys, and the timing is correct.
func (s *appState) handleModifierEvent(e key.Event) bool {
	if e.Name == key.NameCtrl {
		s.ctrlPressed = (e.State == key.Press)
		if e.State == key.Press {
			log.Printf("⌨ [CTRL] Pressed")
		} else {
			log.Printf("⌨ [CTRL] Released")
		}
		return true
	}

	if e.Name == key.NameShift {
		s.shiftPressed = (e.State == key.Press)
		if e.State == key.Press {
			log.Printf("⌨ [SHIFT] Pressed")
		} else {
			log.Printf("⌨ [SHIFT] Released")
		}
		return true
	}

	if e.Name == key.NameAlt {
		log.Printf("⌨ [ALT] %v", e.State)
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
		log.Printf("🔍 [UNIX-SYNC] Ctrl detected via ev.Modifiers (missed Press event)")
	}
	if e.Modifiers.Contain(key.ModShift) && !s.shiftPressed {
		s.shiftPressed = true
		log.Printf("🔍 [UNIX-SYNC] Shift detected via ev.Modifiers (missed Press event)")
	}
}
