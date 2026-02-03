//go:build darwin

package appcore

import (
	"fmt"
	"os"
	"time"

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
//
// macOS/Gio has a bug similar to Windows: modifier Release events arrive
// BEFORE the character key event. Additionally, ev.Modifiers is often empty.
//
// Example timeline when user presses Shift+G:
//  1. User presses Shift   → Shift Press event fires, shiftPressed=true
//  2. User presses G       → NO EVENT YET
//  3. User releases Shift  → Shift Release event fires, shiftPressed=false
//  4. Character "G" arrives → key.Event with ev.Modifiers == empty, shift=false!
//
// Solution: Track the timestamp of modifier Release events. If a character
// key arrives within 200ms of a modifier Release, we know that modifier was
// held during the key press.
//
// Additionally, on macOS we treat Command key as equivalent to Ctrl for
// keybinding purposes, allowing users to use Cmd+key shortcuts (macOS convention).
func (s *appState) handleModifierEvent(e key.Event) bool {
	if darwinDebug {
		darwinInputLog("handleModifierEvent: Name=%q State=%v Modifiers=%v", e.Name, e.State, e.Modifiers)
	}

	// Track Command key as ctrlPressed (macOS convention)
	if e.Name == key.NameCommand {
		if e.State == key.Release {
			// Mark when Command was released - a character key may be coming soon
			s.ctrlReleaseTime = time.Now()
			if darwinDebug {
				darwinInputLog("  -> Command key Release: recorded ctrlReleaseTime")
			}
		} else {
			s.ctrlPressed = true
			if darwinDebug {
				darwinInputLog("  -> Command key Press: ctrlPressed=true")
			}
		}
		return true
	}

	// Also track actual Ctrl (for users who prefer Ctrl+key)
	if e.Name == key.NameCtrl {
		if e.State == key.Release {
			s.ctrlReleaseTime = time.Now()
			if darwinDebug {
				darwinInputLog("  -> Ctrl key Release: recorded ctrlReleaseTime")
			}
		} else {
			s.ctrlPressed = true
			if darwinDebug {
				darwinInputLog("  -> Ctrl key Press: ctrlPressed=true")
			}
		}
		return true
	}

	if e.Name == key.NameShift {
		if e.State == key.Release {
			s.shiftReleaseTime = time.Now()
			if darwinDebug {
				darwinInputLog("  -> Shift key Release: recorded shiftReleaseTime")
			}
		} else {
			s.shiftPressed = true
			if darwinDebug {
				darwinInputLog("  -> Shift key Press: shiftPressed=true")
			}
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
// On macOS, we use temporal logic: if a modifier was released within 200ms,
// it was held during this key press. This handles the out-of-order event delivery.
func (s *appState) syncModifierState(e key.Event) {
	if darwinDebug {
		darwinInputLog("syncModifierState: Name=%q Modifiers=%v (before: ctrl=%v shift=%v)",
			e.Name, e.Modifiers, s.ctrlPressed, s.shiftPressed)
	}

	now := time.Now()

	// Check if Command/Ctrl was released within last 200ms
	ctrlWindow := now.Sub(s.ctrlReleaseTime)
	if ctrlWindow < 200*time.Millisecond && ctrlWindow >= 0 {
		s.ctrlPressed = true
		if darwinDebug {
			darwinInputLog("  -> Ctrl/Cmd released %v ago, set ctrlPressed=true", ctrlWindow)
		}
	}

	// Check if Shift was released within last 200ms
	shiftWindow := now.Sub(s.shiftReleaseTime)
	if shiftWindow < 200*time.Millisecond && shiftWindow >= 0 {
		s.shiftPressed = true
		if darwinDebug {
			darwinInputLog("  -> Shift released %v ago, set shiftPressed=true", shiftWindow)
		}
	}

	// Also check ev.Modifiers as a fallback (often empty on macOS, but try anyway)
	if e.Modifiers.Contain(key.ModCommand) && !s.ctrlPressed {
		s.ctrlPressed = true
		if darwinDebug {
			darwinInputLog("  -> ModCommand in ev.Modifiers, set ctrlPressed=true")
		}
	}
	if e.Modifiers.Contain(key.ModCtrl) && !s.ctrlPressed {
		s.ctrlPressed = true
		if darwinDebug {
			darwinInputLog("  -> ModCtrl in ev.Modifiers, set ctrlPressed=true")
		}
	}
	if e.Modifiers.Contain(key.ModShift) && !s.shiftPressed {
		s.shiftPressed = true
		if darwinDebug {
			darwinInputLog("  -> ModShift in ev.Modifiers, set shiftPressed=true")
		}
	}

	if darwinDebug {
		darwinInputLog("  -> after sync: ctrl=%v shift=%v", s.ctrlPressed, s.shiftPressed)
	}
}
