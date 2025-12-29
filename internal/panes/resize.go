package panes

type splitPath struct {
	node   *PaneNode
	inLeft bool // true if target pane is in left/top subtree
}

const (
	minSplitRatio = 0.1
)

// ResizeActivePane adjusts the nearest split ratio for the active pane.
// Returns true if a resize was applied.
func (pm *PaneManager) ResizeActivePane(dir Direction, delta float32) bool {
	if pm == nil || pm.root == nil || pm.activePane == nil {
		return false
	}
	if delta <= 0 {
		return false
	}

	path, ok := pm.pathToPane(pm.root, pm.activePane)
	if !ok {
		return false
	}

	target, found := findResizeTarget(path, dir)
	if !found {
		if dir == DirLeft || dir == DirRight {
			fallbackDir := DirUp
			if dir == DirRight {
				fallbackDir = DirDown
			}
			target, found = findResizeTarget(path, fallbackDir)
			if found {
				dir = fallbackDir
			}
		}
	}
	if !found {
		return false
	}

	switch dir {
	case DirRight, DirDown:
		target.node.Ratio += delta
	case DirLeft, DirUp:
		target.node.Ratio -= delta
	default:
		return false
	}

	if target.node.Ratio < minSplitRatio {
		target.node.Ratio = minSplitRatio
	}
	if target.node.Ratio > 1-minSplitRatio {
		target.node.Ratio = 1 - minSplitRatio
	}

	return true
}

func (pm *PaneManager) pathToPane(node *PaneNode, target *Pane) ([]splitPath, bool) {
	if node == nil {
		return nil, false
	}
	if node.IsLeaf() {
		return nil, node.Pane == target
	}

	if path, ok := pm.pathToPane(node.Left, target); ok {
		return append([]splitPath{{node: node, inLeft: true}}, path...), true
	}
	if path, ok := pm.pathToPane(node.Right, target); ok {
		return append([]splitPath{{node: node, inLeft: false}}, path...), true
	}

	return nil, false
}

func findResizeTarget(path []splitPath, dir Direction) (splitPath, bool) {
	for i := len(path) - 1; i >= 0; i-- {
		entry := path[i]
		switch dir {
		case DirLeft:
			if entry.node.Split == SplitHorizontal {
				return entry, true
			}
		case DirRight:
			if entry.node.Split == SplitHorizontal {
				return entry, true
			}
		case DirUp:
			if entry.node.Split == SplitVertical {
				return entry, true
			}
		case DirDown:
			if entry.node.Split == SplitVertical {
				return entry, true
			}
		}
	}
	return splitPath{}, false
}
