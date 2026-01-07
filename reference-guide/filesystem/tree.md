# tree.go

**Path:** `/home/javanhut/Development/Vem/internal/filesystem/tree.go`
**Lines:** 503
**Purpose:** File system tree structure for explorer navigation

## Overview

Implements file tree management with:
- Tree structure with expand/collapse
- Flat list rendering for scrolling
- Ignore patterns (.git, node_modules, etc.)
- File operations (rename, delete, create)
- Nested path creation

## Code Blocks

### Lines 1-28: Package and Structs

#### TreeNode (Lines 10-19)
```go
type TreeNode struct {
    Path     string       // Full path
    Name     string       // Filename only
    IsDir    bool         // Directory flag
    Expanded bool         // Expand state
    Children []*TreeNode  // Child nodes
    Parent   *TreeNode    // Parent reference
    Depth    int          // Nesting level
}
```

#### FileTree (Lines 21-28)
```go
type FileTree struct {
    Root           *TreeNode    // Root node
    flatList       []*TreeNode  // Flattened visible nodes
    selectedIndex  int          // Current selection
    needsRebuild   bool         // Rebuild flag
    ignorePatterns []string     // Patterns to hide
}
```

### Lines 30-66: Constructor and Ignore Patterns

#### NewFileTree (Lines 30-53)
Creates tree rooted at path with default ignore patterns.

#### defaultIgnorePatterns (Lines 55-66)
```go
return []string{
    ".git",
    ".gocache",
    "node_modules",
    ".DS_Store",
    "*.swp", "*.swo", "*~",
}
```

### Lines 68-104: Flat List Generation

#### GetFlatList (Lines 68-74)
Returns flattened visible nodes, rebuilding if needed.

#### rebuildFlatList (Lines 76-89)
Creates flat list by traversing tree, clamping selected index.

#### flattenNode (Lines 91-104)
Recursive traversal adding visible nodes to list.

### Lines 106-137: Navigation

| Function | Lines | Purpose |
|----------|-------|---------|
| `SelectedNode()` | 106-113 | Get selected node |
| `SelectedIndex()` | 115-118 | Get selection index |
| `MoveUp()` | 120-127 | Move selection up |
| `MoveDown()` | 129-137 | Move selection down |

### Lines 139-197: Expand/Collapse

| Function | Lines | Behavior |
|----------|-------|----------|
| `Toggle()` | 139-149 | Toggle expand state |
| `Expand()` | 151-170 | Expand or move to first child |
| `Collapse()` | 172-197 | Collapse or move to parent |

### Lines 199-254: Utilities

| Function | Lines | Purpose |
|----------|-------|---------|
| `shouldIgnore()` | 199-213 | Check ignore patterns |
| `AddChild()` | 215-228 | Add sorted child node |
| `ClearChildren()` | 230-236 | Remove all children |
| `IsRoot()` | 238-241 | Check if root node |
| `GetIcon()` | 243-246 | Get file icon |
| `GetExpandIcon()` | 248-254 | Get expand/collapse icon |

### Lines 256-305: Directory Navigation

| Function | Lines | Purpose |
|----------|-------|---------|
| `ChangeRoot()` | 256-276 | Change root directory |
| `NavigateToParent()` | 278-288 | Go to parent dir |
| `CurrentPath()` | 290-296 | Get current root path |
| `IsAtFilesystemRoot()` | 298-305 | Check if at "/" |

### Lines 307-375: File Operations

#### RenameNode (Lines 307-333)
Renames file/directory:
1. Compute new path
2. `os.Rename()` on disk
3. Update node properties
4. Recursively update child paths if directory

#### DeleteNode (Lines 345-375)
Deletes file/directory:
1. `os.RemoveAll()` for directories
2. `os.Remove()` for files
3. Remove from parent's children

### Lines 377-502: File Creation

#### CreateFile (Lines 377-452)
Creates new file, supports nested paths:
- Simple: `"file.txt"` → create file
- Nested: `"dir/subdir/file.txt"` → create directories and file

Uses `os.MkdirAll()` for intermediate directories.

#### addNestedPath (Lines 454-502)
Adds TreeNode hierarchy for nested path.

## Known Issues / Potential Bugs

1. **Line 380-381: nil return on nil parentNode**
   - Returns nil error even on nil input
   - Should return error

2. **Line 78: Magic number 100**
   - Initial capacity of flatList is 100
   - Could be a constant

## Dead/Unused Code

None identified.

## Integration Points

- Used by `appState.fileTree` in app.go
- Rendering in `drawFileExplorer()` in app.go
- Icons from `filesystem/icons.go`
- Loading from `filesystem/loader.go`

---
*Last Updated: Reference guide creation*
