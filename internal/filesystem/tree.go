package filesystem

import (
	"path/filepath"
	"sort"
	"strings"
)

// TreeNode represents a file or directory in the file tree.
type TreeNode struct {
	Path     string
	Name     string
	IsDir    bool
	Expanded bool
	Children []*TreeNode
	Parent   *TreeNode
	Depth    int
}

// FileTree manages the file system tree structure and navigation.
type FileTree struct {
	Root            *TreeNode
	flatList        []*TreeNode
	selectedIndex   int
	needsRebuild    bool
	ignorePatterns  []string
}

// NewFileTree creates a new file tree rooted at the given path.
func NewFileTree(rootPath string) (*FileTree, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	root := &TreeNode{
		Path:     absPath,
		Name:     filepath.Base(absPath),
		IsDir:    true,
		Expanded: true,
		Depth:    0,
	}

	tree := &FileTree{
		Root:           root,
		selectedIndex:  0,
		needsRebuild:   true,
		ignorePatterns: defaultIgnorePatterns(),
	}

	return tree, nil
}

// defaultIgnorePatterns returns common patterns to ignore in file trees.
func defaultIgnorePatterns() []string {
	return []string{
		".git",
		".gocache",
		"node_modules",
		".DS_Store",
		"*.swp",
		"*.swo",
		"*~",
	}
}

// GetFlatList returns a flattened list of visible nodes for rendering.
func (ft *FileTree) GetFlatList() []*TreeNode {
	if ft.needsRebuild {
		ft.rebuildFlatList()
	}
	return ft.flatList
}

// rebuildFlatList creates a flat representation of the tree for rendering.
func (ft *FileTree) rebuildFlatList() {
	ft.flatList = make([]*TreeNode, 0, 100)
	ft.flattenNode(ft.Root)
	ft.needsRebuild = false
	
	// Clamp selected index
	if ft.selectedIndex >= len(ft.flatList) {
		ft.selectedIndex = len(ft.flatList) - 1
	}
	if ft.selectedIndex < 0 && len(ft.flatList) > 0 {
		ft.selectedIndex = 0
	}
}

// flattenNode recursively adds nodes to the flat list if they're visible.
func (ft *FileTree) flattenNode(node *TreeNode) {
	if node == nil {
		return
	}
	
	ft.flatList = append(ft.flatList, node)
	
	if node.IsDir && node.Expanded {
		for _, child := range node.Children {
			ft.flattenNode(child)
		}
	}
}

// SelectedNode returns the currently selected node.
func (ft *FileTree) SelectedNode() *TreeNode {
	list := ft.GetFlatList()
	if ft.selectedIndex >= 0 && ft.selectedIndex < len(list) {
		return list[ft.selectedIndex]
	}
	return nil
}

// SelectedIndex returns the current selection index.
func (ft *FileTree) SelectedIndex() int {
	return ft.selectedIndex
}

// MoveUp moves the selection up one item.
func (ft *FileTree) MoveUp() bool {
	if ft.selectedIndex > 0 {
		ft.selectedIndex--
		return true
	}
	return false
}

// MoveDown moves the selection down one item.
func (ft *FileTree) MoveDown() bool {
	list := ft.GetFlatList()
	if ft.selectedIndex < len(list)-1 {
		ft.selectedIndex++
		return true
	}
	return false
}

// Toggle expands or collapses the selected directory.
func (ft *FileTree) Toggle() bool {
	node := ft.SelectedNode()
	if node == nil || !node.IsDir {
		return false
	}
	
	node.Expanded = !node.Expanded
	ft.needsRebuild = true
	return true
}

// Expand expands the selected directory.
func (ft *FileTree) Expand() bool {
	node := ft.SelectedNode()
	if node == nil || !node.IsDir {
		return false
	}
	
	if !node.Expanded {
		node.Expanded = true
		ft.needsRebuild = true
		return true
	}
	
	// If already expanded, move to first child
	if len(node.Children) > 0 {
		return ft.MoveDown()
	}
	
	return false
}

// Collapse collapses the selected directory or moves to parent.
func (ft *FileTree) Collapse() bool {
	node := ft.SelectedNode()
	if node == nil {
		return false
	}
	
	if node.IsDir && node.Expanded {
		node.Expanded = false
		ft.needsRebuild = true
		return true
	}
	
	// Move to parent
	if node.Parent != nil {
		list := ft.GetFlatList()
		for i, n := range list {
			if n == node.Parent {
				ft.selectedIndex = i
				return true
			}
		}
	}
	
	return false
}

// shouldIgnore checks if a path matches any ignore pattern.
func (ft *FileTree) shouldIgnore(name string) bool {
	for _, pattern := range ft.ignorePatterns {
		if strings.HasPrefix(pattern, "*.") {
			// Simple suffix match for *.ext patterns
			ext := pattern[1:]
			if strings.HasSuffix(name, ext) {
				return true
			}
		} else if name == pattern {
			return true
		}
	}
	return false
}

// AddChild adds a child node to a directory, maintaining sorted order.
func (node *TreeNode) AddChild(child *TreeNode) {
	child.Parent = node
	child.Depth = node.Depth + 1
	node.Children = append(node.Children, child)
	
	// Sort: directories first, then alphabetically
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return strings.ToLower(node.Children[i].Name) < strings.ToLower(node.Children[j].Name)
	})
}

// ClearChildren removes all children from a directory node.
func (node *TreeNode) ClearChildren() {
	for _, child := range node.Children {
		child.Parent = nil
	}
	node.Children = nil
}

// IsRoot returns true if this is the root node.
func (node *TreeNode) IsRoot() bool {
	return node.Parent == nil
}

// ChangeRoot changes the root directory of the tree to a new path.
func (ft *FileTree) ChangeRoot(newPath string) error {
	absPath, err := filepath.Abs(newPath)
	if err != nil {
		return err
	}

	root := &TreeNode{
		Path:     absPath,
		Name:     filepath.Base(absPath),
		IsDir:    true,
		Expanded: true,
		Depth:    0,
	}

	ft.Root = root
	ft.selectedIndex = 0
	ft.needsRebuild = true

	return nil
}

// NavigateToParent changes the root to the parent directory.
func (ft *FileTree) NavigateToParent() error {
	parentPath := filepath.Dir(ft.Root.Path)
	
	// Check if we're already at root (e.g., "/" or "C:\")
	if parentPath == ft.Root.Path {
		return nil // Already at filesystem root
	}

	return ft.ChangeRoot(parentPath)
}

// CurrentPath returns the current root path.
func (ft *FileTree) CurrentPath() string {
	if ft.Root == nil {
		return ""
	}
	return ft.Root.Path
}

// IsAtFilesystemRoot returns true if we're at the filesystem root.
func (ft *FileTree) IsAtFilesystemRoot() bool {
	if ft.Root == nil {
		return false
	}
	parentPath := filepath.Dir(ft.Root.Path)
	return parentPath == ft.Root.Path
}
