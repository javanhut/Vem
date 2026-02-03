//go:build !darwin

package appcore

import "gioui.org/io/key"

// Stub functions for non-darwin platforms - these are no-ops

func (s *appState) debugKeyEvent(e key.Event)                       {}
func (s *appState) debugEditEvent(e key.EditEvent)                  {}
func (s *appState) debugPrintableKey(ev key.Event, r rune, ok bool) {}
func (s *appState) debugModifierEvent(e key.Event, keyType string, newValue bool) {}
func (s *appState) debugSyncModifierBefore(e key.Event)             {}
func (s *appState) debugSyncModifierAfter()                         {}
