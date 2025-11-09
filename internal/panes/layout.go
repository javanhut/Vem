package panes

// SplitDirection defines how a pane is split.
type SplitDirection int

const (
	// SplitHorizontal creates a left | right split (vertical divider).
	SplitHorizontal SplitDirection = iota
	// SplitVertical creates a top / bottom split (horizontal divider).
	SplitVertical
)

// PaneNode represents a node in the pane tree.
// The tree is a binary split tree where:
// - Leaf nodes contain a Pane (actual editor pane)
// - Internal nodes contain a Split direction and two children
type PaneNode struct {
	// Leaf node fields (if this is a pane)
	Pane *Pane

	// Internal node fields (if this is a split container)
	Split SplitDirection
	Ratio float32   // Split ratio (always 0.5 for 50/50 splits)
	Left  *PaneNode // Left or top child
	Right *PaneNode // Right or bottom child
}

// NewPaneNode creates a leaf node containing a pane.
func NewPaneNode(pane *Pane) *PaneNode {
	return &PaneNode{
		Pane: pane,
	}
}

// NewSplitNode creates an internal node representing a split.
func NewSplitNode(direction SplitDirection, left, right *PaneNode) *PaneNode {
	return &PaneNode{
		Split: direction,
		Ratio: 0.5, // Always 50/50 splits
		Left:  left,
		Right: right,
	}
}

// IsLeaf returns true if this node is a leaf (contains a pane).
func (n *PaneNode) IsLeaf() bool {
	return n.Pane != nil
}

// FindPane recursively searches for a pane by ID.
func (n *PaneNode) FindPane(id string) *Pane {
	if n == nil {
		return nil
	}

	if n.IsLeaf() {
		if n.Pane.ID == id {
			return n.Pane
		}
		return nil
	}

	// Search left subtree
	if pane := n.Left.FindPane(id); pane != nil {
		return pane
	}

	// Search right subtree
	return n.Right.FindPane(id)
}

// CollectPanes returns all panes in the tree (in-order traversal).
func (n *PaneNode) CollectPanes() []*Pane {
	if n == nil {
		return nil
	}

	if n.IsLeaf() {
		return []*Pane{n.Pane}
	}

	var panes []*Pane
	panes = append(panes, n.Left.CollectPanes()...)
	panes = append(panes, n.Right.CollectPanes()...)
	return panes
}

// CountPanes returns the total number of panes in the tree.
func (n *PaneNode) CountPanes() int {
	if n == nil {
		return 0
	}

	if n.IsLeaf() {
		return 1
	}

	return n.Left.CountPanes() + n.Right.CountPanes()
}
