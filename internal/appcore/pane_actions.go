package appcore

import (
	"fmt"
	"strings"

	"gioui.org/io/key"

	"github.com/javanhut/vem/internal/panes"
)

const paneResizeStep = 0.05

// handleSplitVertical creates a vertical split (vertical divider - left|right).
func (s *appState) handleSplitVertical() {
	if s.paneManager == nil {
		s.status = "Pane manager not initialized"
		return
	}

	// Create a new empty buffer for the new pane
	newBufferIndex := s.bufferMgr.CreateEmptyBuffer()

	// Split the active pane horizontally (creates vertical divider)
	if err := s.paneManager.SplitHorizontal(newBufferIndex); err != nil {
		s.status = fmt.Sprintf("Split failed: %v", err)
	} else {
		paneCount := s.paneManager.PaneCount()
		s.status = fmt.Sprintf("Split vertical (│) - %d panes total | Use :e or Ctrl+P to open file", paneCount)
	}
}

// handleSplitHorizontal creates a horizontal split (horizontal divider - top/bottom).
func (s *appState) handleSplitHorizontal() {
	if s.paneManager == nil {
		s.status = "Pane manager not initialized"
		return
	}

	// Create a new empty buffer for the new pane
	newBufferIndex := s.bufferMgr.CreateEmptyBuffer()

	// Split the active pane vertically (creates horizontal divider)
	if err := s.paneManager.SplitVertical(newBufferIndex); err != nil {
		s.status = fmt.Sprintf("Split failed: %v", err)
	} else {
		paneCount := s.paneManager.PaneCount()
		s.status = fmt.Sprintf("Split horizontal (─) - %d panes total | Use :e or Ctrl+P to open file", paneCount)
	}
}

// handlePaneFocusLeft focuses the pane to the left.
func (s *appState) handlePaneFocusLeft() {
	if s.paneManager == nil {
		return
	}

	if s.paneManager.NavigateLeft() {
		s.status = "Focused left pane ←"
	} else {
		s.status = "No pane to the left"
	}
}

// handlePaneFocusRight focuses the pane to the right.
func (s *appState) handlePaneFocusRight() {
	if s.paneManager == nil {
		return
	}

	if s.paneManager.NavigateRight() {
		s.status = "Focused right pane →"
	} else {
		s.status = "No pane to the right"
	}
}

// handlePaneFocusUp focuses the pane above.
func (s *appState) handlePaneFocusUp() {
	if s.paneManager == nil {
		return
	}

	if s.paneManager.NavigateUp() {
		s.status = "Focused pane above ↑"
	} else {
		s.status = "No pane above"
	}
}

// handlePaneFocusDown focuses the pane below.
func (s *appState) handlePaneFocusDown() {
	if s.paneManager == nil {
		return
	}

	if s.paneManager.NavigateDown() {
		s.status = "Focused pane below ↓"
	} else {
		s.status = "No pane below"
	}
}

// handlePaneCycleNext cycles to the next pane.
func (s *appState) handlePaneCycleNext() {
	if s.paneManager == nil {
		return
	}

	paneCount := s.paneManager.PaneCount()
	if paneCount <= 1 {
		s.status = "Only one pane open"
		return
	}

	s.paneManager.CycleNextPane()

	// Get active pane index after cycling
	allPanes := s.paneManager.AllPanes()
	activeIdx := -1
	newActivePane := s.paneManager.ActivePane()
	for i, p := range allPanes {
		if p == newActivePane {
			activeIdx = i + 1 // 1-based for display
			break
		}
	}

	s.status = fmt.Sprintf("Cycled to pane %d/%d", activeIdx, paneCount)
}

// handlePaneClose closes the active pane.
func (s *appState) handlePaneClose() {
	if s.paneManager == nil {
		return
	}

	activePane := s.paneManager.ActivePane()
	if activePane == nil {
		return
	}

	// Get the buffer for this pane
	buf := s.bufferMgr.GetBuffer(activePane.BufferIndex)
	if buf == nil {
		// No buffer, just close the pane if multiple exist
		if s.paneManager.PaneCount() > 1 {
			s.paneManager.ClosePane()
			s.status = fmt.Sprintf("Pane closed - %d panes remaining", s.paneManager.PaneCount())
		} else {
			// Last pane with no buffer - switch to buffer 0 (sample buffer)
			activePane.SetBufferIndex(0)
			s.status = "No buffer to close"
		}
		return
	}

	// Check if buffer is modified (unless terminal)
	if buf.Modified() && !buf.IsTerminal() {
		s.status = "Buffer has unsaved changes (use :q! to force close)"
		return
	}

	bufferIndex := activePane.BufferIndex

	// Close terminal if this buffer has one
	s.closeTerminal(bufferIndex)

	// Multiple panes - close this pane and buffer
	if s.paneManager.PaneCount() > 1 {
		if err := s.paneManager.ClosePane(); err != nil {
			s.status = fmt.Sprintf("Error closing pane: %v", err)
			return
		}
		s.cleanupLSPForBuffer(buf)
		s.updateHighlighterCacheOnBufferClose(bufferIndex)
		s.bufferMgr.CloseBuffer(bufferIndex, false)
		s.status = fmt.Sprintf("Pane closed - %d panes remaining", s.paneManager.PaneCount())
		return
	}

	// Last pane - close buffer but keep editor open
	s.cleanupLSPForBuffer(buf)
	s.updateHighlighterCacheOnBufferClose(bufferIndex)
	s.bufferMgr.CloseBuffer(bufferIndex, false)

	// Ensure we have at least one buffer (switch to buffer 0 - sample buffer)
	if s.bufferMgr.BufferCount() == 0 || s.bufferMgr.GetBuffer(0) == nil {
		// This shouldn't happen, but handle gracefully
		activePane.SetBufferIndex(0)
		s.status = "Buffer closed"
	} else {
		// Switch to buffer 0 (sample buffer)
		activePane.SetBufferIndex(0)
		s.status = "Buffer closed"
	}
}

// handlePaneEqualize makes all panes equal size.
func (s *appState) handlePaneEqualize() {
	if s.paneManager == nil {
		return
	}

	s.paneManager.Equalize()
	s.status = "All panes equalized (50/50)"
}

// handlePaneZoomToggle toggles zoom for the active pane.
func (s *appState) handlePaneZoomToggle() {
	if s.paneManager == nil {
		return
	}

	s.paneManager.ToggleZoom()

	if s.paneManager.IsZoomed() {
		s.status = "Pane zoomed (Ctrl+S o to restore)"
	} else {
		s.status = "Pane restored to normal view"
	}
}

// handlePaneCommand handles Ctrl+S prefix pane commands.
func (s *appState) handlePaneCommand(ev key.Event) {
	// Convert to lowercase for case-insensitive matching
	keyName := strings.ToLower(string(ev.Name))

	// Check for split commands
	switch keyName {
	case "v":
		s.executeAction(ActionSplitVertical, ev)
		return
	case "h":
		s.executeAction(ActionSplitHorizontal, ev)
		return
	case "=":
		s.executeAction(ActionPaneEqualize, ev)
		return
	case "o":
		s.executeAction(ActionPaneZoomToggle, ev)
		return
	default:
		s.status = "Unknown pane command (v=vsplit h=hsplit ==equalize o=zoom)"
	}
}

func (s *appState) togglePaneResizeMode() {
	if s.paneResizeMode {
		s.paneResizeMode = false
		s.status = "Pane resize off"
		return
	}

	if s.paneManager == nil || s.paneManager.PaneCount() <= 1 {
		s.status = "Only one pane open"
		return
	}

	if s.paneManager.IsZoomed() {
		s.status = "Exit zoom before resizing panes"
		return
	}

	s.paneResizeMode = true
	s.status = "Resize panes: arrows or h/j/k/l, Esc to exit"
}

func (s *appState) handlePaneResizeKey(ev key.Event) bool {
	if !s.paneResizeMode {
		return false
	}

	switch ev.Name {
	case key.NameEscape:
		s.paneResizeMode = false
		s.status = "Pane resize off"
		return true
	case key.NameLeftArrow:
		s.resizeActivePane(panes.DirLeft)
		return true
	case key.NameRightArrow:
		s.resizeActivePane(panes.DirRight)
		return true
	case key.NameUpArrow:
		s.resizeActivePane(panes.DirUp)
		return true
	case key.NameDownArrow:
		s.resizeActivePane(panes.DirDown)
		return true
	}

	if r, ok := s.printableKey(ev); ok {
		switch r {
		case 'h':
			s.resizeActivePane(panes.DirLeft)
			return true
		case 'l':
			s.resizeActivePane(panes.DirRight)
			return true
		case 'j':
			s.resizeActivePane(panes.DirDown)
			return true
		case 'k':
			s.resizeActivePane(panes.DirUp)
			return true
		}
	}

	s.status = "Resize panes: arrows or h/j/k/l, Esc to exit"
	return true
}

func (s *appState) resizeActivePane(dir panes.Direction) {
	if s.paneManager == nil {
		return
	}

	if s.paneManager.IsZoomed() {
		s.status = "Pane resize disabled while zoomed"
		return
	}

	if s.paneManager.ResizeActivePane(dir, paneResizeStep) {
		s.status = "Pane resized"
	} else {
		s.status = "No pane edge to resize"
	}
}
