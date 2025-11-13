package appcore

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/javanhut/vem/internal/editor"
	"github.com/javanhut/vem/internal/panes"
	"github.com/javanhut/vem/internal/terminal"
)

// drawPanes is the entry point for rendering all panes.
func (s *appState) drawPanes(gtx layout.Context) layout.Dimensions {
	if s.paneManager == nil {
		// Fallback to single buffer view
		return s.drawBuffer(gtx)
	}

	paneCount := s.paneManager.PaneCount()

	// If zoomed, just draw the zoomed pane
	if s.paneManager.IsZoomed() {
		zoomedPane := s.paneManager.ZoomedPane()
		if zoomedPane != nil {
			return s.drawSinglePane(gtx, zoomedPane)
		}
	}

	// Render the pane tree
	root := s.paneManager.Root()
	if root == nil {
		fmt.Printf("[PANE_RENDER] WARNING: Root is nil, paneCount=%d\n", paneCount)
		return s.drawBuffer(gtx)
	}

	if paneCount > 1 {
		fmt.Printf("[PANE_RENDER] Rendering %d panes\n", paneCount)
	}

	return s.renderPaneNode(gtx, root)
}

// renderPaneNode recursively renders a pane node (either a split or a leaf pane).
func (s *appState) renderPaneNode(gtx layout.Context, node *panes.PaneNode) layout.Dimensions {
	if node == nil {
		return layout.Dimensions{}
	}

	// Leaf node: render the actual pane
	if node.IsLeaf() {
		return s.drawSinglePane(gtx, node.Pane)
	}

	// Internal node: render split with separator
	if node.Split == panes.SplitHorizontal {
		// Left | Right split (vertical divider)
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(node.Ratio, func(gtx layout.Context) layout.Dimensions {
				return s.renderPaneNode(gtx, node.Left)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.drawPaneSeparator(gtx, true)
			}),
			layout.Flexed(1-node.Ratio, func(gtx layout.Context) layout.Dimensions {
				return s.renderPaneNode(gtx, node.Right)
			}),
		)
	} else {
		// Top / Bottom split (horizontal divider)
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(node.Ratio, func(gtx layout.Context) layout.Dimensions {
				return s.renderPaneNode(gtx, node.Left)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.drawPaneSeparator(gtx, false)
			}),
			layout.Flexed(1-node.Ratio, func(gtx layout.Context) layout.Dimensions {
				return s.renderPaneNode(gtx, node.Right)
			}),
		)
	}
}

// drawSinglePane renders a single pane with its buffer content.
func (s *appState) drawSinglePane(gtx layout.Context, pane *panes.Pane) layout.Dimensions {
	if pane == nil {
		return layout.Dimensions{}
	}

	// Get buffer for this pane
	buf := s.bufferMgr.GetBuffer(pane.BufferIndex)
	if buf == nil {
		return layout.Dimensions{}
	}

	// Debug: Log pane geometry
	if s.mode == modeInsert && pane.Active {
		fmt.Printf("[PANE_GEOMETRY] Pane=%s Active=%v Constraints: Min=%v Max=%v IsTerminal=%v\n",
			pane.ID, pane.Active, gtx.Constraints.Min, gtx.Constraints.Max, buf.IsTerminal())
	}

	// Determine background color based on active state
	bgColor := inactivePaneBg
	if pane.Active {
		bgColor = activePaneBg
	}

	// Draw background first
	bgRect := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, bgColor)
	bgRect.Pop()

	// Check if this is a terminal buffer
	if buf.IsTerminal() {
		return s.drawTerminalPane(gtx, pane, buf)
	}

	// For non-active panes, we need to render with the correct buffer context
	// but without actually changing the global active pane state (which would
	// interfere with input handling).
	//
	// We temporarily swap the viewport state so drawBuffer renders the right view.
	wasActive := pane.Active
	if !wasActive {
		// Save current viewport state
		oldViewportTop := s.viewportTopLine

		// Use this pane's viewport
		s.viewportTopLine = pane.ViewportTop

		// Temporarily override activeBuffer to return this pane's buffer
		// by quietly swapping the pane manager's active pane ONLY for rendering
		oldActivePane := s.paneManager.ActivePane()
		s.paneManager.SetActivePaneQuiet(pane)

		// Draw buffer content
		dims := s.drawBuffer(gtx)

		// Restore original active pane (quietly, without triggering side effects)
		s.paneManager.SetActivePaneQuiet(oldActivePane)

		// Restore viewport state
		s.viewportTopLine = oldViewportTop

		return dims
	}

	// For active pane, sync viewport state between pane and global state
	// This ensures that each pane maintains its own scroll position
	oldViewportTop := s.viewportTopLine
	s.viewportTopLine = pane.ViewportTop // Sync FROM pane TO global for rendering

	dims := s.drawBuffer(gtx)

	// Save any viewport changes back to the pane
	pane.SetViewportTop(s.viewportTopLine) // Sync back FROM global TO pane

	// Restore previous global viewport state
	s.viewportTopLine = oldViewportTop

	return dims
}

// drawPaneSeparator draws a 1px separator line between panes.
func (s *appState) drawPaneSeparator(gtx layout.Context, vertical bool) layout.Dimensions {
	var width, height int
	if vertical {
		width = 1
		height = gtx.Constraints.Max.Y
	} else {
		width = gtx.Constraints.Max.X
		height = 1
	}

	rect := clip.Rect{Max: image.Pt(width, height)}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, paneSeparator)
	rect.Pop()

	return layout.Dimensions{Size: image.Pt(width, height)}
}

// drawTerminalPane renders a terminal pane
func (s *appState) drawTerminalPane(gtx layout.Context, pane *panes.Pane, buf *editor.Buffer) layout.Dimensions {
	// Get terminal instance
	term, exists := s.terminals[pane.BufferIndex]
	if !exists || term == nil {
		// Terminal not found - show error message
		label := material.Body1(s.theme, "Terminal not initialized")
		label.Color = color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, label.Layout)
	}

	// Check if terminal is running
	if !term.IsRunning() {
		label := material.Body1(s.theme, "Terminal exited (press Ctrl+X to close)")
		label.Color = color.NRGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff}
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, label.Layout)
	}

	// Get screen buffer
	screen := term.GetScreen()
	if screen == nil {
		return layout.Dimensions{}
	}

	// Draw terminal content
	return s.drawTerminalContent(gtx, screen, pane.BufferIndex)
}

// drawTerminalContent renders the terminal screen buffer with viewport scrolling
func (s *appState) drawTerminalContent(gtx layout.Context, screen *terminal.ScreenBuffer, bufferIndex int) layout.Dimensions {
	cols, rows := screen.Dimensions()
	cursorX, cursorY, cursorStyle := screen.GetCursor()

	// Calculate character dimensions using actual text measurement
	testLabel := material.Body1(s.theme, "M") // Use 'M' as widest character
	testLabel.Font.Typeface = "JetBrainsMono"
	testGtx := gtx
	testGtx.Constraints = layout.Constraints{Max: image.Point{X: 1000, Y: 1000}}
	testDims := testLabel.Layout(testGtx)
	charWidth := testDims.Size.X
	charHeight := testDims.Size.Y
	if charWidth == 0 {
		charWidth = 8
	}
	if charHeight == 0 {
		charHeight = 16
	}

	// Calculate lines per page for viewport
	insetDp := 16 // Top + bottom inset
	availableHeight := gtx.Constraints.Max.Y - gtx.Dp(unit.Dp(insetDp))
	linesPerPage := availableHeight / charHeight
	if linesPerPage < 1 {
		linesPerPage = 1
	}
	if linesPerPage > rows {
		linesPerPage = rows
	}

	// Ensure cursor is visible (auto-scroll)
	s.ensureTerminalCursorVisible(bufferIndex, linesPerPage, screen)

	// Get viewport top line
	viewportTop, exists := s.terminalViewports[bufferIndex]
	if !exists {
		viewportTop = 0
		s.terminalViewports[bufferIndex] = 0
	}

	// Calculate viewport end
	viewportEnd := viewportTop + linesPerPage
	if viewportEnd > rows {
		viewportEnd = rows
	}

	inset := layout.Inset{
		Top:    unit.Dp(8),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(8),
		Left:   unit.Dp(16),
	}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Draw only visible lines in viewport
		for y := viewportTop; y < viewportEnd; y++ {
			line := screen.GetLine(y)
			for x := 0; x < len(line.Cells) && x < cols; x++ {
				cell := line.Cells[x]

				// Calculate cell position (adjusted for viewport)
				cellX := x * charWidth
				cellY := (y - viewportTop) * charHeight // Subtract viewportTop to adjust Y position

				// Draw cell background
				bgRect := clip.Rect{
					Min: image.Pt(cellX, cellY),
					Max: image.Pt(cellX+charWidth, cellY+charHeight),
				}.Push(gtx.Ops)
				paint.Fill(gtx.Ops, cell.BG)
				bgRect.Pop()

				// Draw cursor if at this position
				if x == cursorX && y == cursorY && cursorStyle == terminal.CursorBlock {
					cursorRect := clip.Rect{
						Min: image.Pt(cellX, cellY),
						Max: image.Pt(cellX+charWidth, cellY+charHeight),
					}.Push(gtx.Ops)
					paint.Fill(gtx.Ops, cursorColor)
					cursorRect.Pop()
				}

				// Draw character
				char := cell.Rune
				if char == 0 {
					char = ' '
				}

				label := material.Body1(s.theme, string(char))
				label.Font.Typeface = "JetBrainsMono"

				// Use cell foreground color (or cursor color if cursor is here)
				if x == cursorX && y == cursorY && cursorStyle == terminal.CursorBlock {
					// Invert color for cursor
					label.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
				} else {
					label.Color = cell.FG
				}

				// Position and draw the character
				offset := op.Offset(image.Pt(cellX, cellY)).Push(gtx.Ops)
				label.Layout(gtx)
				offset.Pop()
			}
		}

		// Return dimensions based on visible area
		return layout.Dimensions{
			Size: image.Pt(cols*charWidth, linesPerPage*charHeight),
		}
	})
}
