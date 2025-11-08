# Navigation in ProjectVem

This document describes the navigation features in ProjectVem, including pane navigation between the file tree explorer and the text editor.

## Pane Navigation

ProjectVem supports Vim-like split pane navigation with keyboard shortcuts to quickly jump between the file tree explorer and the text editor.

### Key Bindings

#### Ctrl+H - Jump to File Tree
When in **NORMAL mode** with the editor focused:
- Pressing `Ctrl+H` will switch focus to the file tree explorer (left pane)
- If the tree view is hidden, it will be opened and focused automatically
- The focused pane is indicated by a blue border on its edge
- Status bar will display: "Focus: Tree View (use Ctrl+L to return to editor)"

#### Ctrl+L - Jump to Editor
When in **EXPLORER mode** (tree view focused):
- Pressing `Ctrl+L` will switch focus back to the text editor (right pane)
- The tree view remains visible but unfocused
- The focused pane is indicated by a blue border on its edge
- Status bar will display: "Focus: Editor (use Ctrl+H to return to tree view)"

#### Ctrl+T - Toggle Tree View
In any mode:
- Pressing `Ctrl+T` will toggle the visibility of the file tree explorer
- If the tree view is hidden, it will be shown and focused (EXPLORER mode)
- If the tree view is visible, it will be hidden and focus returns to the editor (NORMAL mode)

### Visual Indicators

The currently focused pane is indicated by:
1. A blue border (3px wide) on the edge of the focused pane
   - File tree: border on the right edge
   - Editor: border on the left edge
2. Mode indicator in the status bar
   - `EXPLORER` when tree view is focused
   - `NORMAL` when editor is focused
3. Helpful status messages showing current focus and how to switch

### Navigation Flow

```
NORMAL mode (Editor focused)
    |
    | Ctrl+H
    v
EXPLORER mode (Tree view focused)
    |
    | Ctrl+L
    v
NORMAL mode (Editor focused)
```

### File Tree Navigation

When the file tree is focused (EXPLORER mode), you can:
- `j` or Down Arrow: Move selection down
- `k` or Up Arrow: Move selection up
- `h` or Left Arrow: Collapse directory or move to parent
- `l` or Right Arrow: Expand directory or move to first child
- `Enter`: Open file or toggle directory expansion
- `r`: Refresh the tree
- `u`: Navigate to parent directory
- `q` or `Esc`: Return to NORMAL mode (exit explorer)

### Command Mode Integration

You can also control the explorer from command mode:
- `:ex` or `:explore` - Toggle the file tree explorer
- `:cd <path>` - Change the working directory and show the tree
- `:pwd` - Print the current working directory

## Examples

### Opening and Navigating Files

1. Press `Ctrl+T` to open the file tree
2. Use `j`/`k` to navigate to the desired file
3. Press `Enter` to open the file
4. The editor will be focused automatically with the opened file
5. Press `Ctrl+H` to jump back to the tree view
6. Navigate to another file
7. Press `Ctrl+L` to jump back to the editor

### Quick Pane Switching

1. While editing, press `Ctrl+H` to quickly check the file tree
2. Navigate with `j`/`k` to find a file
3. Press `Ctrl+L` to return to editing without opening anything

## Implementation Details

The pane navigation feature is implemented in `internal/appcore/app.go`:

- Focus state is tracked using the `explorerFocused` boolean flag
- Mode transitions between `modeNormal` and `modeExplorer`
- Visual borders are rendered using Gio's `clip.Rect` and `paint.Fill` operations
- Keyboard event handlers in `handleNormalMode` and `handleExplorerMode` manage the transitions

### Code References

- Ctrl+H handler: `internal/appcore/app.go:505-519`
- Ctrl+L handler: `internal/appcore/app.go:768-772`
- Focus border rendering (tree): `internal/appcore/app.go:448-454`
- Focus border rendering (editor): `internal/appcore/app.go:260-268`
