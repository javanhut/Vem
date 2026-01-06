package appcore

import (
	"image"
	"image/color"
	"strings"

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
		return s.drawBuffer(gtx, true)
	}

	s.paneManager.SetLayoutSize(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

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
		return s.drawBuffer(gtx, true)
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
		dims := s.drawBuffer(gtx, false)

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

	dims := s.drawBuffer(gtx, true)

	// Save any viewport changes back to the pane
	pane.SetViewportTop(s.viewportTopLine) // Sync back FROM global TO pane

	s.updateViewStateForPane(pane)
	s.maybeSaveViewState()

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
		for y := viewportTop; y < viewportEnd; y++ {
			line := screen.GetLine(y)
			cellY := (y - viewportTop) * charHeight

			// Batch contiguous cells with same colors into runs.
			runStart := 0
			var runFG, runBG color.NRGBA
			runInitialized := false
			flushRun := func(runEnd int) {
				if !runInitialized || runEnd <= runStart {
					return
				}
				cellX := runStart * charWidth
				width := (runEnd - runStart) * charWidth
				// Background
				bgRect := clip.Rect{
					Min: image.Pt(cellX, cellY),
					Max: image.Pt(cellX+width, cellY+charHeight),
				}.Push(gtx.Ops)
				paint.Fill(gtx.Ops, runBG)
				bgRect.Pop()

				// Text
				var b strings.Builder
				for x := runStart; x < runEnd && x < len(line.Cells) && x < cols; x++ {
					ch := line.Cells[x].Rune
					if ch == 0 {
						ch = ' '
					}
					b.WriteRune(ch)
				}
				label := material.Body1(s.theme, b.String())
				label.Font.Typeface = "JetBrainsMono"
				label.Color = runFG
				offset := op.Offset(image.Pt(cellX, cellY)).Push(gtx.Ops)
				label.Layout(gtx)
				offset.Pop()
			}

			maxCells := cols
			if len(line.Cells) < maxCells {
				maxCells = len(line.Cells)
			}

			for x := 0; x < maxCells; x++ {
				cell := line.Cells[x]
				cellFG := cell.FG
				cellBG := cell.BG

				// Cursor cell gets handled separately for correct colors.
				if x == cursorX && y == cursorY && cursorStyle == terminal.CursorBlock {
					flushRun(x)
					runInitialized = false

					// Calculate actual cursor X position by measuring text width
					// This ensures cursor aligns with rendered text (avoids font kerning issues)
					var prefixText strings.Builder
					for px := 0; px < x && px < len(line.Cells); px++ {
						ch := line.Cells[px].Rune
						if ch == 0 {
							ch = ' '
						}
						prefixText.WriteRune(ch)
					}
					actualCellX := 0
					if prefixText.Len() > 0 {
						prefixLabel := material.Body1(s.theme, prefixText.String())
						prefixLabel.Font.Typeface = "JetBrainsMono"
						prefixMeasure := gtx
						prefixMeasure.Constraints = layout.Constraints{Max: image.Point{X: 10000, Y: 10000}}
						prefixDims := prefixLabel.Layout(prefixMeasure)
						actualCellX = prefixDims.Size.X
					}

					cursorRect := clip.Rect{
						Min: image.Pt(actualCellX, cellY),
						Max: image.Pt(actualCellX+charWidth, cellY+charHeight),
					}.Push(gtx.Ops)
					paint.Fill(gtx.Ops, cursorColor)
					cursorRect.Pop()

					ch := cell.Rune
					if ch == 0 {
						ch = ' '
					}
					label := material.Body1(s.theme, string(ch))
					label.Font.Typeface = "JetBrainsMono"
					label.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
					offset := op.Offset(image.Pt(actualCellX, cellY)).Push(gtx.Ops)
					label.Layout(gtx)
					offset.Pop()
					runStart = x + 1
					continue
				}

				if !runInitialized {
					runStart = x
					runFG = cellFG
					runBG = cellBG
					runInitialized = true
					continue
				}

				if cellFG != runFG || cellBG != runBG {
					flushRun(x)
					runStart = x
					runFG = cellFG
					runBG = cellBG
					runInitialized = true
				}
			}

			flushRun(maxCells)
		}

		return layout.Dimensions{
			Size: image.Pt(cols*charWidth, linesPerPage*charHeight),
		}
	})
}
