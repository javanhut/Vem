package appcore

import (
	"fmt"
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/javanhut/ProjectVem/internal/panes"
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

	// Determine background color based on active state
	bgColor := inactivePaneBg
	if pane.Active {
		bgColor = activePaneBg
	}

	// Draw background first
	bgRect := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, bgColor)
	bgRect.Pop()

	// Temporarily switch active pane for rendering
	oldActivePane := s.paneManager.ActivePane()
	wasActive := pane.Active

	// If this is not the active pane, temporarily make it active for cursor rendering
	// but restore it after
	if !wasActive {
		pane.SetActive(true)
		s.paneManager.SetActivePane(pane)
	}

	// Draw buffer content using existing drawBuffer logic
	dims := s.drawBuffer(gtx)

	// Restore active pane
	if !wasActive {
		pane.SetActive(false)
		if oldActivePane != nil {
			s.paneManager.SetActivePane(oldActivePane)
		}
	}

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
