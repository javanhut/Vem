//go:build darwin

package appcore

import (
	"fmt"
	"os"

	"gioui.org/io/key"
)

// darwinDebugEnabled controls debug logging for macOS input
var darwinDebugEnabled = os.Getenv("VEM_DEBUG_DARWIN") != ""

func darwinLog(format string, args ...interface{}) {
	if darwinDebugEnabled {
		fmt.Fprintf(os.Stderr, "[darwin] "+format+"\n", args...)
	}
}

func (s *appState) debugKeyEvent(e key.Event) {
	if !darwinDebugEnabled {
		return
	}
	fmt.Fprintf(os.Stderr, "[darwin-key] Event: Name=%q State=%v Modifiers=%v\n", e.Name, e.State, e.Modifiers)
	fmt.Fprintf(os.Stderr, "[darwin-key]   tracked: ctrl=%v shift=%v\n", s.ctrlPressed, s.shiftPressed)
}

func (s *appState) debugEditEvent(e key.EditEvent) {
	if !darwinDebugEnabled {
		return
	}
	fmt.Fprintf(os.Stderr, "[darwin-edit] EditEvent: Text=%q\n", e.Text)
	fmt.Fprintf(os.Stderr, "[darwin-edit]   skipNextEdit=%v mode=%v\n", s.skipNextEdit, s.mode)
}

func (s *appState) debugPrintableKey(ev key.Event, r rune, ok bool) {
	if !darwinDebugEnabled {
		return
	}
	fmt.Fprintf(os.Stderr, "[darwin-key] printableKey: Name=%q -> rune=%q ok=%v (shiftPressed=%v)\n",
		ev.Name, r, ok, s.shiftPressed)
}

func (s *appState) debugModifierEvent(e key.Event, keyType string, newValue bool) {
	if !darwinDebugEnabled {
		return
	}
	darwinLog("  -> %s key: %s=%v", keyType, keyType+"Pressed", newValue)
}

func (s *appState) debugSyncModifierBefore(e key.Event) {
	if !darwinDebugEnabled {
		return
	}
	darwinLog("syncModifierState: Name=%q Modifiers=%v (before: ctrl=%v shift=%v)",
		e.Name, e.Modifiers, s.ctrlPressed, s.shiftPressed)
}

func (s *appState) debugSyncModifierAfter() {
	if !darwinDebugEnabled {
		return
	}
	darwinLog("  -> after sync: ctrl=%v shift=%v", s.ctrlPressed, s.shiftPressed)
}
