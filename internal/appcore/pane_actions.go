package appcore

import (
	"fmt"
	"strings"

	"gioui.org/io/key"
)

// handleSplitVertical creates a vertical split (vertical divider - left|right).
func (s *appState) handleSplitVertical() {
	fmt.Printf("[PANE_SPLIT] Starting vertical split (left|right)\n")

	if s.paneManager == nil {
		s.status = "Pane manager not initialized"
		fmt.Printf("[PANE_SPLIT] ERROR: Pane manager is nil\n")
		return
	}

	fmt.Printf("[PANE_SPLIT] Current pane count: %d\n", s.paneManager.PaneCount())
	fmt.Printf("[PANE_SPLIT] Current buffer count: %d\n", s.bufferMgr.BufferCount())

	// Create a new empty buffer for the new pane
	newBufferIndex := s.bufferMgr.CreateEmptyBuffer()

	fmt.Printf("[PANE_SPLIT] Created new buffer, total buffers: %d\n", s.bufferMgr.BufferCount())
	fmt.Printf("[PANE_SPLIT] New buffer index: %d\n", newBufferIndex)

	// Split the active pane horizontally (creates vertical divider)
	fmt.Printf("[PANE_SPLIT] Calling SplitHorizontal with buffer index %d\n", newBufferIndex)
	if err := s.paneManager.SplitHorizontal(newBufferIndex); err != nil {
		s.status = fmt.Sprintf("Split failed: %v", err)
		fmt.Printf("[PANE_SPLIT] ERROR: Split failed: %v\n", err)
	} else {
		paneCount := s.paneManager.PaneCount()
		s.status = fmt.Sprintf("Split vertical (│) - %d panes total | Use :e or Ctrl+P to open file", paneCount)
		fmt.Printf("[PANE_SPLIT] SUCCESS: Vertical split created, now have %d panes\n", paneCount)

		// Debug: Print all panes
		allPanes := s.paneManager.AllPanes()
		for i, p := range allPanes {
			fmt.Printf("[PANE_SPLIT]   Pane %d: ID=%s BufferIndex=%d Active=%v\n", i, p.ID, p.BufferIndex, p.Active)
		}
	}
}

// handleSplitHorizontal creates a horizontal split (horizontal divider - top/bottom).
func (s *appState) handleSplitHorizontal() {
	fmt.Printf("[PANE_SPLIT] Starting horizontal split (top/bottom)\n")

	if s.paneManager == nil {
		s.status = "Pane manager not initialized"
		fmt.Printf("[PANE_SPLIT] ERROR: Pane manager is nil\n")
		return
	}

	fmt.Printf("[PANE_SPLIT] Current pane count: %d\n", s.paneManager.PaneCount())
	fmt.Printf("[PANE_SPLIT] Current buffer count: %d\n", s.bufferMgr.BufferCount())

	// Create a new empty buffer for the new pane
	newBufferIndex := s.bufferMgr.CreateEmptyBuffer()

	fmt.Printf("[PANE_SPLIT] Created new buffer, total buffers: %d\n", s.bufferMgr.BufferCount())
	fmt.Printf("[PANE_SPLIT] New buffer index: %d\n", newBufferIndex)

	// Split the active pane vertically (creates horizontal divider)
	fmt.Printf("[PANE_SPLIT] Calling SplitVertical with buffer index %d\n", newBufferIndex)
	if err := s.paneManager.SplitVertical(newBufferIndex); err != nil {
		s.status = fmt.Sprintf("Split failed: %v", err)
		fmt.Printf("[PANE_SPLIT] ERROR: Split failed: %v\n", err)
	} else {
		paneCount := s.paneManager.PaneCount()
		s.status = fmt.Sprintf("Split horizontal (─) - %d panes total | Use :e or Ctrl+P to open file", paneCount)
		fmt.Printf("[PANE_SPLIT] SUCCESS: Horizontal split created, now have %d panes\n", paneCount)

		// Debug: Print all panes
		allPanes := s.paneManager.AllPanes()
		for i, p := range allPanes {
			fmt.Printf("[PANE_SPLIT]   Pane %d: ID=%s BufferIndex=%d Active=%v\n", i, p.ID, p.BufferIndex, p.Active)
		}
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
	fmt.Printf("[PANE_CYCLE] Starting pane cycle\n")

	if s.paneManager == nil {
		fmt.Printf("[PANE_CYCLE] ERROR: Pane manager is nil\n")
		return
	}

	paneCount := s.paneManager.PaneCount()
	fmt.Printf("[PANE_CYCLE] Current pane count: %d\n", paneCount)

	if paneCount <= 1 {
		s.status = "Only one pane open"
		fmt.Printf("[PANE_CYCLE] Only one pane, nothing to cycle\n")
		return
	}

	// Get current active pane before cycling
	oldActivePane := s.paneManager.ActivePane()
	fmt.Printf("[PANE_CYCLE] Before cycle - Active pane: %s\n", oldActivePane.ID)

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

	fmt.Printf("[PANE_CYCLE] After cycle - Active pane: %s (index %d/%d)\n", newActivePane.ID, activeIdx, paneCount)
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
		s.bufferMgr.CloseBuffer(bufferIndex, false)
		s.status = fmt.Sprintf("Pane closed - %d panes remaining", s.paneManager.PaneCount())
		return
	}

	// Last pane - close buffer but keep editor open
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
