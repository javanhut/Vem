//go:build windows

package appcore

import (
	"syscall"
	"time"

	"gioui.org/io/key"
)

const (
	vkShift    = 0x10
	vkControl  = 0x11
	vkMenu     = 0x12
	vkLShift   = 0xA0
	vkRShift   = 0xA1
	vkLControl = 0xA2
	vkRControl = 0xA3
	vkLMenu    = 0xA4
	vkRMenu    = 0xA5
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

func keyDown(vk int) bool {
	state, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return state&0x8000 != 0
}

func keyDownAny(vks ...int) bool {
	for _, vk := range vks {
		if keyDown(vk) {
			return true
		}
	}
	return false
}

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
			debugf("[MOD_EVENT] Ctrl release")
		} else {
			// Press events don't arrive on Windows, but handle it just in case
			s.ctrlPressed = true
			debugf("[MOD_EVENT] Ctrl press")
		}
		return true
	}

	if e.Name == key.NameShift {
		if e.State == key.Release {
			s.shiftReleaseTime = time.Now()
			debugf("[MOD_EVENT] Shift release")
		} else {
			s.shiftPressed = true
			debugf("[MOD_EVENT] Shift press")
		}
		return true
	}

	if e.Name == key.NameAlt {
		if e.State == key.Release {
			s.altReleaseTime = time.Now()
			debugf("[MOD_EVENT] Alt release")
		} else {
			s.altPressed = true
			debugf("[MOD_EVENT] Alt press")
		}
		return true
	}

	return false
}

// syncModifierState syncs the tracked modifier state before handling character keys.
// On Windows, use actual key state when available, then fall back to temporal logic
// for the Gio modifier ordering bug.
func (s *appState) syncModifierState(e key.Event) {
	now := time.Now()

	ctrlDown := e.Modifiers.Contain(key.ModCtrl) || keyDownAny(vkControl, vkLControl, vkRControl)
	shiftDown := e.Modifiers.Contain(key.ModShift) || keyDownAny(vkShift, vkLShift, vkRShift)
	altDown := e.Modifiers.Contain(key.ModAlt) || keyDownAny(vkMenu, vkLMenu, vkRMenu)

	// Temporal fallback: if a modifier was released recently, treat it as held.
	ctrlWindow := now.Sub(s.ctrlReleaseTime)
	if !ctrlDown && ctrlWindow < 200*time.Millisecond && ctrlWindow >= 0 {
		ctrlDown = true
	}

	shiftWindow := now.Sub(s.shiftReleaseTime)
	if !shiftDown && shiftWindow < 200*time.Millisecond && shiftWindow >= 0 {
		shiftDown = true
	}

	altWindow := now.Sub(s.altReleaseTime)
	if !altDown && altWindow < 200*time.Millisecond && altWindow >= 0 {
		altDown = true
	}

	s.ctrlPressed = ctrlDown
	s.shiftPressed = shiftDown
	s.altPressed = altDown
}
