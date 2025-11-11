package panes

import "math"

// Direction represents a navigation direction.
type Direction int

const (
	DirLeft Direction = iota
	DirRight
	DirDown
	DirUp
)

// NavigateDirection navigates to the pane in the specified direction.
// Returns true if navigation was successful.
func (pm *PaneManager) NavigateDirection(dir Direction) bool {
	if pm.activePane == nil {
		return false
	}

	// For now, implement simple cycling behavior
	// In a full implementation, this would do geometric pane selection
	// based on actual pane positions on screen
	allPanes := pm.AllPanes()
	if len(allPanes) <= 1 {
		return false
	}

	// Find current index
	currentIdx := -1
	for i, p := range allPanes {
		if p == pm.activePane {
			currentIdx = i
			break
		}
	}

	if currentIdx == -1 {
		return false
	}

	var nextIdx int
	switch dir {
	case DirLeft:
		// Move to previous pane (wraps around)
		nextIdx = (currentIdx - 1 + len(allPanes)) % len(allPanes)
	case DirRight:
		// Move to next pane (wraps around)
		nextIdx = (currentIdx + 1) % len(allPanes)
	case DirUp:
		// Move to previous pane (wraps around)
		nextIdx = (currentIdx - 1 + len(allPanes)) % len(allPanes)
	case DirDown:
		// Move to next pane (wraps around)
		nextIdx = (currentIdx + 1) % len(allPanes)
	default:
		return false
	}

	pm.SetActivePane(allPanes[nextIdx])
	return true
}

// NavigateLeft navigates to the pane on the left.
func (pm *PaneManager) NavigateLeft() bool {
	return pm.NavigateDirection(DirLeft)
}

// NavigateRight navigates to the pane on the right.
func (pm *PaneManager) NavigateRight() bool {
	return pm.NavigateDirection(DirRight)
}

// NavigateUp navigates to the pane above.
func (pm *PaneManager) NavigateUp() bool {
	return pm.NavigateDirection(DirUp)
}

// NavigateDown navigates to the pane below.
func (pm *PaneManager) NavigateDown() bool {
	return pm.NavigateDirection(DirDown)
}

// Equalize makes all panes equal size by resetting all ratios to 0.5.
func (pm *PaneManager) Equalize() {
	pm.equalizeNode(pm.root)
}

// equalizeNode recursively sets all split ratios to 0.5.
func (pm *PaneManager) equalizeNode(node *PaneNode) {
	if node == nil || node.IsLeaf() {
		return
	}

	node.Ratio = 0.5
	pm.equalizeNode(node.Left)
	pm.equalizeNode(node.Right)
}

// PaneGeometry represents the on-screen position and size of a pane.
type PaneGeometry struct {
	Pane   *Pane
	X      int // Left edge
	Y      int // Top edge
	Width  int
	Height int
}

// CalculateGeometry calculates the on-screen geometry of all panes.
// This is used for geometric pane navigation.
func (pm *PaneManager) CalculateGeometry(width, height int) []PaneGeometry {
	var geometries []PaneGeometry
	pm.calculateNodeGeometry(pm.root, 0, 0, width, height, &geometries)
	return geometries
}

// calculateNodeGeometry recursively calculates pane geometries.
func (pm *PaneManager) calculateNodeGeometry(node *PaneNode, x, y, width, height int, geometries *[]PaneGeometry) {
	if node == nil {
		return
	}

	if node.IsLeaf() {
		*geometries = append(*geometries, PaneGeometry{
			Pane:   node.Pane,
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		})
		return
	}

	// Calculate split position
	if node.Split == SplitHorizontal {
		// Left | Right split
		leftWidth := int(float32(width) * node.Ratio)
		rightWidth := width - leftWidth - 1 // -1 for separator

		pm.calculateNodeGeometry(node.Left, x, y, leftWidth, height, geometries)
		pm.calculateNodeGeometry(node.Right, x+leftWidth+1, y, rightWidth, height, geometries)
	} else {
		// Top / Bottom split
		topHeight := int(float32(height) * node.Ratio)
		bottomHeight := height - topHeight - 1 // -1 for separator

		pm.calculateNodeGeometry(node.Left, x, y, width, topHeight, geometries)
		pm.calculateNodeGeometry(node.Right, x, y+topHeight+1, width, bottomHeight, geometries)
	}
}

// FindPaneInDirection finds the pane in the given direction from the active pane.
// Uses geometric calculation for more intuitive navigation.
func (pm *PaneManager) FindPaneInDirection(dir Direction, width, height int) *Pane {
	if pm.activePane == nil {
		return nil
	}

	geometries := pm.CalculateGeometry(width, height)
	if len(geometries) <= 1 {
		return nil
	}

	// Find current pane geometry
	var currentGeom *PaneGeometry
	for i := range geometries {
		if geometries[i].Pane == pm.activePane {
			currentGeom = &geometries[i]
			break
		}
	}

	if currentGeom == nil {
		return nil
	}

	// Find closest pane in the specified direction
	var bestPane *Pane
	bestDistance := math.MaxFloat64

	for i := range geometries {
		if geometries[i].Pane == pm.activePane {
			continue
		}

		candidateGeom := &geometries[i]

		// Check if pane is in the right direction
		var isInDirection bool
		var distance float64

		switch dir {
		case DirLeft:
			isInDirection = candidateGeom.X < currentGeom.X
			if isInDirection {
				dx := float64(currentGeom.X - (candidateGeom.X + candidateGeom.Width))
				dy := float64(currentGeom.Y - candidateGeom.Y)
				distance = math.Sqrt(dx*dx + dy*dy)
			}
		case DirRight:
			isInDirection = candidateGeom.X > currentGeom.X
			if isInDirection {
				dx := float64(candidateGeom.X - (currentGeom.X + currentGeom.Width))
				dy := float64(currentGeom.Y - candidateGeom.Y)
				distance = math.Sqrt(dx*dx + dy*dy)
			}
		case DirUp:
			isInDirection = candidateGeom.Y < currentGeom.Y
			if isInDirection {
				dx := float64(currentGeom.X - candidateGeom.X)
				dy := float64(currentGeom.Y - (candidateGeom.Y + candidateGeom.Height))
				distance = math.Sqrt(dx*dx + dy*dy)
			}
		case DirDown:
			isInDirection = candidateGeom.Y > currentGeom.Y
			if isInDirection {
				dx := float64(currentGeom.X - candidateGeom.X)
				dy := float64(candidateGeom.Y - (currentGeom.Y + currentGeom.Height))
				distance = math.Sqrt(dx*dx + dy*dy)
			}
		}

		if isInDirection && distance < bestDistance {
			bestDistance = distance
			bestPane = candidateGeom.Pane
		}
	}

	return bestPane
}
