package filesystem

import (
	"os"
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
	Root           *TreeNode
	flatList       []*TreeNode
	selectedIndex  int
	needsRebuild   bool
	ignorePatterns []string
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

// GetIcon returns the appropriate Nerd Font icon for this node
func (node *TreeNode) GetIcon() string {
	return GetFileIcon(node.Name, node.IsDir)
}

// GetExpandIcon returns the expand/collapse icon if this is a directory
func (node *TreeNode) GetExpandIcon() string {
	if !node.IsDir {
		return ""
	}
	return GetExpandIcon(node.Expanded)
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

// RenameNode renames a file or directory.
func (ft *FileTree) RenameNode(node *TreeNode, newName string) error {
	if node == nil {
		return nil
	}

	oldPath := node.Path
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newName)

	// Rename on disk
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	// Update node
	node.Name = newName
	node.Path = newPath

	// If it's a directory, update all children's paths recursively
	if node.IsDir {
		ft.updateChildPaths(node)
	}

	ft.needsRebuild = true
	return nil
}

// updateChildPaths recursively updates paths for all children of a node.
func (ft *FileTree) updateChildPaths(node *TreeNode) {
	for _, child := range node.Children {
		child.Path = filepath.Join(node.Path, child.Name)
		if child.IsDir {
			ft.updateChildPaths(child)
		}
	}
}

// DeleteNode removes a file or directory from disk.
func (ft *FileTree) DeleteNode(node *TreeNode) error {
	if node == nil {
		return nil
	}

	// Remove from disk
	if node.IsDir {
		if err := os.RemoveAll(node.Path); err != nil {
			return err
		}
	} else {
		if err := os.Remove(node.Path); err != nil {
			return err
		}
	}

	// Remove from parent's children
	if node.Parent != nil {
		parent := node.Parent
		for i, child := range parent.Children {
			if child == node {
				parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
				break
			}
		}
	}

	ft.needsRebuild = true
	return nil
}

// CreateFile creates a new file in the specified directory.
// Supports creating nested directories with path separators (e.g., "dir/subdir/file.txt").
func (ft *FileTree) CreateFile(parentNode *TreeNode, fileName string) error {
	if parentNode == nil {
		return nil
	}

	// Determine the base directory
	var baseDir string
	var targetNode *TreeNode

	if parentNode.IsDir {
		baseDir = parentNode.Path
		targetNode = parentNode
	} else {
		baseDir = filepath.Dir(parentNode.Path)
		targetNode = parentNode.Parent
	}

	if targetNode == nil {
		return nil
	}

	// Check if fileName contains path separators
	if strings.Contains(fileName, "/") || strings.Contains(fileName, string(filepath.Separator)) {
		// Split into directory path and final filename
		dir := filepath.Dir(fileName)
		finalName := filepath.Base(fileName)

		// Create full path for directories
		fullDirPath := filepath.Join(baseDir, dir)

		// Create all intermediate directories
		if err := os.MkdirAll(fullDirPath, 0755); err != nil {
			return err
		}

		// Create the file in the final directory
		fullFilePath := filepath.Join(fullDirPath, finalName)
		file, err := os.Create(fullFilePath)
		if err != nil {
			return err
		}
		file.Close()

		// Add directory and file nodes to tree
		if err := ft.addNestedPath(targetNode, fileName); err != nil {
			return err
		}

	} else {
		// Simple case: just a filename, no directories
		filePath := filepath.Join(baseDir, fileName)

		// Create the file on disk
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		file.Close()

		// Create new TreeNode
		newNode := &TreeNode{
			Path:     filePath,
			Name:     fileName,
			IsDir:    false,
			Expanded: false,
		}

		// Add to parent's children
		targetNode.AddChild(newNode)
	}

	ft.needsRebuild = true
	return nil
}

// addNestedPath adds directory and file nodes for a nested path like "dir1/dir2/file.txt"
func (ft *FileTree) addNestedPath(parentNode *TreeNode, path string) error {
	if parentNode == nil {
		return nil
	}

	// Split path into components (use filepath.Separator for cross-platform support)
	// But normalize to forward slash first for consistent splitting
	normalizedPath := filepath.ToSlash(path)
	parts := strings.Split(normalizedPath, "/")

	currentNode := parentNode
	currentPath := parentNode.Path

	// Process each component
	for i, part := range parts {
		if part == "" {
			continue
		}

		isLastPart := (i == len(parts)-1)
		currentPath = filepath.Join(currentPath, part)

		// Check if this node already exists in children
		found := false
		for _, child := range currentNode.Children {
			if child.Name == part {
				currentNode = child
				found = true
				break
			}
		}

		if !found {
			// Create new node
			newNode := &TreeNode{
				Path:     currentPath,
				Name:     part,
				IsDir:    !isLastPart,
				Expanded: false,
			}

			currentNode.AddChild(newNode)
			currentNode = newNode
		}
	}

	return nil
}
