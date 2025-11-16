//go:build windows

package appcore

import (
	"time"

	"gioui.org/io/key"
)

// handleModifierEvent handles modifier key events on Windows.
// Windows/Gio has a critical bug: Ctrl/Shift Press events NEVER arrive.
// Only Release events are sent, and they arrive BEFORE character keys.
//
// Example timeline when user presses Ctrl+T:
//  1. User presses Ctrl     → NO EVENT (bug!)
//  2. User presses T        → NO EVENT YET
//  3. User releases Ctrl    → Ctrl Release event fires
//  4. Character "T" arrives → key.Event with ev.Modifiers == empty
//
// Solution: Track the timestamp of modifier Release events. If a character
// key arrives within 200ms of a modifier Release, we know that modifier was
// held during the key press. The 200ms window accounts for Windows event
// buffering and user typing speed variability.
func (s *appState) handleModifierEvent(e key.Event) bool {
	if e.Name == key.NameCtrl {
		if e.State == key.Release {
			// Mark when Ctrl was released - a character key is coming soon!
			s.ctrlReleaseTime = time.Now()
		} else {
			// Press events don't arrive on Windows, but handle it just in case
			s.ctrlPressed = true
		}
		return true
	}

	if e.Name == key.NameShift {
		if e.State == key.Release {
			s.shiftReleaseTime = time.Now()
		} else {
			s.shiftPressed = true
		}
		return true
	}

	if e.Name == key.NameAlt {
		return true
	}

	return false
}

// syncModifierState syncs the tracked modifier state before handling character keys.
// On Windows, we use temporal logic: if a modifier was released within 200ms,
// it was held during this key press.
func (s *appState) syncModifierState(e key.Event) {
	now := time.Now()

	// Check if Ctrl was released within last 200ms
	ctrlWindow := now.Sub(s.ctrlReleaseTime)
	if ctrlWindow < 200*time.Millisecond && ctrlWindow >= 0 {
		s.ctrlPressed = true
	}

	// Check if Shift was released within last 200ms
	shiftWindow := now.Sub(s.shiftReleaseTime)
	if shiftWindow < 200*time.Millisecond && shiftWindow >= 0 {
		s.shiftPressed = true
	}

	// Also check ev.Modifiers as a fallback (usually empty on Windows, but try anyway)
	if e.Modifiers.Contain(key.ModCtrl) {
		s.ctrlPressed = true
	}
	if e.Modifiers.Contain(key.ModShift) {
		s.shiftPressed = true
	}
}
