# filesystem Package Index

**Path:** `/home/javanhut/Development/Vem/internal/filesystem/`
**Purpose:** File system navigation and file tree management

## Files Overview

| File | Lines | Purpose |
|------|-------|---------|
| [tree.go](tree.md) | 503 | File tree structure |

**Additional files** (not yet documented):
- `finder.go` - File finding utilities
- `icons.go` - File type icons
- `loader.go` - Directory loading

## Key Types

### FileTree (tree.go)
Top-level tree container:
- Root node
- Flat list for rendering
- Selected index
- Ignore patterns

### TreeNode (tree.go)
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

## Known Issues Summary

1. **tree.go:380-381** - Returns nil error on nil input

## Ignore Patterns

Default patterns hidden from tree:
- `.git`
- `.gocache`
- `node_modules`
- `.DS_Store`
- `*.swp`, `*.swo`, `*~`

## File Operations

| Operation | Method | Description |
|-----------|--------|-------------|
| Expand | `Expand()` | Show directory children |
| Collapse | `Collapse()` | Hide directory children |
| Toggle | `Toggle()` | Toggle expand state |
| Navigate Up | `MoveUp()` | Move selection up |
| Navigate Down | `MoveDown()` | Move selection down |
| Change Root | `ChangeRoot()` | Set new root directory |
| Parent Dir | `NavigateToParent()` | Go to parent directory |

## File Management

| Operation | Method | Description |
|-----------|--------|-------------|
| Create File | `CreateFile()` | Create new file/nested path |
| Rename | `RenameNode()` | Rename file/directory |
| Delete | `DeleteNode()` | Delete file/directory |

## Flat List Rendering

```
Tree Structure:              Flat List (for rendering):
project/                     [0] project/
├── src/                     [1]   src/
│   ├── main.go              [2]     main.go
│   └── util.go              [3]     util.go
├── README.md                [4]   README.md
└── go.mod                   [5]   go.mod

Navigation: MoveUp/MoveDown moves through flat list
            Expand/Collapse rebuilds flat list
```

## Nested Path Creation

`CreateFile("dir/subdir/file.txt")`:
1. Creates `dir/` if not exists
2. Creates `dir/subdir/` if not exists
3. Creates `dir/subdir/file.txt`
4. Adds TreeNode hierarchy

---
*Last Updated: Reference guide creation*
