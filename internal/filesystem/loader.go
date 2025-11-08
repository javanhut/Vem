package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
)

// LoadDirectory loads the immediate children of a directory node.
func (ft *FileTree) LoadDirectory(node *TreeNode) error {
	if node == nil || !node.IsDir {
		return nil
	}

	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return err
	}

	// Clear existing children
	node.ClearChildren()

	// Add ".." parent directory entry if not at root AND this is the tree root
	if node == ft.Root && !ft.IsAtFilesystemRoot() {
		parentPath := filepath.Dir(node.Path)
		parentNode := &TreeNode{
			Path:     parentPath,
			Name:     "..",
			IsDir:    true,
			Expanded: false,
		}
		node.AddChild(parentNode)
	}

	for _, entry := range entries {
		name := entry.Name()
		
		// Skip ignored files
		if ft.shouldIgnore(name) {
			continue
		}

		childPath := filepath.Join(node.Path, name)
		child := &TreeNode{
			Path:     childPath,
			Name:     name,
			IsDir:    entry.IsDir(),
			Expanded: false,
		}

		node.AddChild(child)
	}

	ft.needsRebuild = true
	return nil
}

// Refresh reloads the tree from the filesystem.
func (ft *FileTree) Refresh() error {
	// Save expanded state
	expandedPaths := make(map[string]bool)
	ft.collectExpandedPaths(ft.Root, expandedPaths)

	// Reload root
	if err := ft.LoadDirectory(ft.Root); err != nil {
		return err
	}

	// Recursively reload expanded directories
	if err := ft.reloadExpanded(ft.Root, expandedPaths); err != nil {
		return err
	}

	ft.needsRebuild = true
	return nil
}

// collectExpandedPaths stores which directories are expanded.
func (ft *FileTree) collectExpandedPaths(node *TreeNode, paths map[string]bool) {
	if node == nil {
		return
	}

	if node.IsDir && node.Expanded {
		paths[node.Path] = true
		for _, child := range node.Children {
			ft.collectExpandedPaths(child, paths)
		}
	}
}

// reloadExpanded reloads directories that were previously expanded.
func (ft *FileTree) reloadExpanded(node *TreeNode, expandedPaths map[string]bool) error {
	if node == nil || !node.IsDir {
		return nil
	}

	if expandedPaths[node.Path] {
		node.Expanded = true
		if err := ft.LoadDirectory(node); err != nil {
			return err
		}

		for _, child := range node.Children {
			if err := ft.reloadExpanded(child, expandedPaths); err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadInitial loads the initial tree structure (root + first level).
func (ft *FileTree) LoadInitial() error {
	if err := ft.LoadDirectory(ft.Root); err != nil {
		return err
	}
	ft.needsRebuild = true
	return nil
}

// ExpandAndLoad expands a directory and loads its children if not already loaded.
func (ft *FileTree) ExpandAndLoad(node *TreeNode) error {
	if node == nil || !node.IsDir {
		return nil
	}

	// Load children if not already loaded
	if len(node.Children) == 0 {
		if err := ft.LoadDirectory(node); err != nil {
			return err
		}
	}

	node.Expanded = true
	ft.needsRebuild = true
	return nil
}

// WalkTree walks the file tree up to maxDepth and calls fn for each node.
func WalkTree(rootPath string, maxDepth int, fn func(path string, info fs.FileInfo, depth int) error) error {
	return walkTreeRecursive(rootPath, 0, maxDepth, fn)
}

// walkTreeRecursive is the recursive helper for WalkTree.
func walkTreeRecursive(path string, currentDepth, maxDepth int, fn func(string, fs.FileInfo, int) error) error {
	if currentDepth > maxDepth {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if err := fn(path, info, currentDepth); err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		if err := walkTreeRecursive(childPath, currentDepth+1, maxDepth, fn); err != nil {
			return err
		}
	}

	return nil
}

// IsTextFile checks if a file is likely a text file based on extension.
func IsTextFile(path string) bool {
	textExtensions := []string{
		".txt", ".md", ".go", ".py", ".js", ".ts", ".jsx", ".tsx",
		".c", ".h", ".cpp", ".hpp", ".rs", ".java", ".rb", ".php",
		".html", ".css", ".scss", ".json", ".yaml", ".yml", ".toml",
		".xml", ".sh", ".bash", ".zsh", ".vim", ".lua", ".sql",
		".conf", ".config", ".ini", ".env", ".gitignore", ".log",
	}

	ext := filepath.Ext(path)
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}

	// Files without extension might be text (Makefile, Dockerfile, etc.)
	if ext == "" {
		name := filepath.Base(path)
		textNames := []string{
			"Makefile", "Dockerfile", "README", "LICENSE", "CHANGELOG",
			"TODO", "NOTES", "AUTHORS", "CONTRIBUTING",
		}
		for _, textName := range textNames {
			if name == textName {
				return true
			}
		}
	}

	return false
}

// GetFileSize returns a human-readable file size.
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
