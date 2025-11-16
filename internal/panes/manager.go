package panes

import (
	"fmt"
)

// PaneManager manages the pane tree and active pane state.
type PaneManager struct {
	root       *PaneNode
	activePane *Pane
	nextPaneID int
	zoomed     *Pane // If set, this pane is temporarily maximized
}

// NewPaneManager creates a new pane manager with a single initial pane.
func NewPaneManager(initialBufferIndex int) *PaneManager {
	pane := NewPane("pane-0", initialBufferIndex)
	pane.SetActive(true)

	return &PaneManager{
		root:       NewPaneNode(pane),
		activePane: pane,
		nextPaneID: 1,
		zoomed:     nil,
	}
}

// Root returns the root of the pane tree.
func (pm *PaneManager) Root() *PaneNode {
	return pm.root
}

// ActivePane returns the currently active pane.
func (pm *PaneManager) ActivePane() *Pane {
	return pm.activePane
}

// PaneCount returns the total number of panes.
func (pm *PaneManager) PaneCount() int {
	if pm.root == nil {
		return 0
	}
	return pm.root.CountPanes()
}

// AllPanes returns all panes in the tree.
func (pm *PaneManager) AllPanes() []*Pane {
	if pm.root == nil {
		return nil
	}
	return pm.root.CollectPanes()
}

// SetActivePane sets the given pane as active and deactivates others.
func (pm *PaneManager) SetActivePane(pane *Pane) {
	if pane == nil {
		return
	}

	// Deactivate all panes
	for _, p := range pm.AllPanes() {
		p.SetActive(false)
	}

	// Activate the target pane
	pane.SetActive(true)
	pm.activePane = pane

	fmt.Printf("[PANE_MANAGER] SetActivePane: ID=%s, BufferIndex=%d\n", pane.ID, pane.BufferIndex)
}

// SetActivePaneQuiet is a low-level setter for the active pane that doesn't
// trigger side effects like deactivating other panes. Used only for rendering.
func (pm *PaneManager) SetActivePaneQuiet(pane *Pane) {
	pm.activePane = pane
}

// SplitVertical splits the active pane vertically (creates horizontal divider).
// Creates a new pane below the active pane.
func (pm *PaneManager) SplitVertical(newBufferIndex int) error {
	if pm.activePane == nil {
		return fmt.Errorf("no active pane to split")
	}

	// Create new pane
	newPane := NewPane(fmt.Sprintf("pane-%d", pm.nextPaneID), newBufferIndex)
	pm.nextPaneID++

	// Find the node containing the active pane and replace it with a split
	pm.root = pm.splitNodeContainingPane(pm.root, pm.activePane, SplitVertical, newPane)

	// Activate the new pane
	pm.SetActivePane(newPane)

	return nil
}

// SplitHorizontal splits the active pane horizontally (creates vertical divider).
// Creates a new pane to the right of the active pane.
func (pm *PaneManager) SplitHorizontal(newBufferIndex int) error {
	if pm.activePane == nil {
		return fmt.Errorf("no active pane to split")
	}

	// Create new pane
	newPane := NewPane(fmt.Sprintf("pane-%d", pm.nextPaneID), newBufferIndex)
	pm.nextPaneID++

	// Find the node containing the active pane and replace it with a split
	pm.root = pm.splitNodeContainingPane(pm.root, pm.activePane, SplitHorizontal, newPane)

	// Activate the new pane
	pm.SetActivePane(newPane)

	return nil
}

// splitNodeContainingPane recursively finds and splits the node containing the target pane.
func (pm *PaneManager) splitNodeContainingPane(node *PaneNode, targetPane *Pane, direction SplitDirection, newPane *Pane) *PaneNode {
	if node == nil {
		return nil
	}

	// If this is the leaf containing the target pane, replace it with a split
	if node.IsLeaf() && node.Pane == targetPane {
		oldPaneNode := NewPaneNode(node.Pane)
		newPaneNode := NewPaneNode(newPane)
		return NewSplitNode(direction, oldPaneNode, newPaneNode)
	}

	// Recurse into children if this is an internal node
	if !node.IsLeaf() {
		node.Left = pm.splitNodeContainingPane(node.Left, targetPane, direction, newPane)
		node.Right = pm.splitNodeContainingPane(node.Right, targetPane, direction, newPane)
	}

	return node
}

// ClosePane closes the active pane and removes it from the tree.
func (pm *PaneManager) ClosePane() error {
	if pm.activePane == nil {
		return fmt.Errorf("no active pane to close")
	}

	// Don't allow closing the last pane
	if pm.PaneCount() <= 1 {
		return fmt.Errorf("cannot close the last pane")
	}

	// Find the parent of the node containing the active pane and collapse it
	paneToClose := pm.activePane
	pm.root = pm.removeNodeContainingPane(pm.root, paneToClose)

	// Set a new active pane (first available)
	allPanes := pm.AllPanes()
	if len(allPanes) > 0 {
		pm.SetActivePane(allPanes[0])
	} else {
		pm.activePane = nil
	}

	return nil
}

// removeNodeContainingPane recursively finds and removes the node containing the target pane.
func (pm *PaneManager) removeNodeContainingPane(node *PaneNode, targetPane *Pane) *PaneNode {
	if node == nil {
		return nil
	}

	// If this is a split node, check if either child contains the target
	if !node.IsLeaf() {
		// Check if left child is the target leaf
		if node.Left.IsLeaf() && node.Left.Pane == targetPane {
			// Replace this split with the right child
			return node.Right
		}

		// Check if right child is the target leaf
		if node.Right.IsLeaf() && node.Right.Pane == targetPane {
			// Replace this split with the left child
			return node.Left
		}

		// Recurse into children
		node.Left = pm.removeNodeContainingPane(node.Left, targetPane)
		node.Right = pm.removeNodeContainingPane(node.Right, targetPane)
	}

	return node
}

// CycleNextPane cycles to the next pane in the list.
func (pm *PaneManager) CycleNextPane() {
	allPanes := pm.AllPanes()
	if len(allPanes) <= 1 {
		return
	}

	// Find current active pane index
	currentIdx := -1
	for i, p := range allPanes {
		if p == pm.activePane {
			currentIdx = i
			break
		}
	}

	// Cycle to next pane
	nextIdx := (currentIdx + 1) % len(allPanes)
	pm.SetActivePane(allPanes[nextIdx])
}

// ToggleZoom toggles zoom (maximize/restore) for the active pane.
func (pm *PaneManager) ToggleZoom() {
	if pm.zoomed != nil {
		// Unzoom
		pm.zoomed = nil
	} else {
		// Zoom active pane
		pm.zoomed = pm.activePane
	}
}

// IsZoomed returns true if a pane is currently zoomed.
func (pm *PaneManager) IsZoomed() bool {
	return pm.zoomed != nil
}

// ZoomedPane returns the currently zoomed pane, or nil if none.
func (pm *PaneManager) ZoomedPane() *Pane {
	return pm.zoomed
}

// FindPaneByBufferIndex finds a pane displaying the given buffer index.
func (pm *PaneManager) FindPaneByBufferIndex(bufferIndex int) *Pane {
	for _, pane := range pm.AllPanes() {
		if pane.BufferIndex == bufferIndex {
			return pane
		}
	}
	return nil
}
