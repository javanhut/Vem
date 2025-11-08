# Navigation in ProjectVem

This document describes the navigation features in ProjectVem, including pane navigation between the file tree explorer and the text editor, and fullscreen mode management.

## Fullscreen Mode

### Shift+Enter - Toggle Fullscreen

ProjectVem supports fullscreen mode for a distraction-free editing experience.

#### Key Binding
- **Shift+Enter**: Toggle fullscreen mode (works in all modes)

#### Behavior
- Press `Shift+Enter` to enter fullscreen mode
- Press `Shift+Enter` again to return to windowed mode
- Status bar displays "FULLSCREEN" indicator when in fullscreen mode
- Works across all editor modes (NORMAL, INSERT, VISUAL, DELETE, EXPLORER, COMMAND)
- Supported on all platforms (Linux, macOS, Windows, WebAssembly)

#### Mode-Specific Behavior
- **INSERT mode**: Regular `Enter` still inserts a newline; only `Shift+Enter` toggles fullscreen
- **COMMAND mode**: `Shift+Enter` toggles fullscreen instead of executing the command
- **EXPLORER mode**: `Shift+Enter` toggles fullscreen; regular `Enter` still opens files/directories
- **NORMAL, VISUAL, DELETE modes**: `Shift+Enter` toggles fullscreen without side effects

#### Visual Feedback
- Status bar shows "FULLSCREEN" indicator when in fullscreen mode
- Status message confirms "Entered fullscreen (Shift+Enter to exit)" or "Exited fullscreen"
- Window state is automatically tracked even if changed via OS window controls

#### Platform Notes
- On some platforms, fullscreen behavior may vary based on OS window management
- Gio UI handles platform-specific fullscreen implementation automatically
- Wayland, X11, macOS, and Windows are all fully supported

## Pane Navigation

ProjectVem supports Vim-like split pane navigation with keyboard shortcuts to quickly jump between the file tree explorer and the text editor.

### Key Bindings

#### Ctrl+H - Jump to File Tree
When in **NORMAL mode** with the editor focused and tree view visible:
- Pressing `Ctrl+H` will switch focus to the file tree explorer (left pane)
- **Note**: The tree view must already be open (use `Ctrl+T` to open it first)
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
1. Open tree view with Ctrl+T

NORMAL mode (Editor focused, tree visible)
    |
    | Ctrl+H (jump to tree)
    v
EXPLORER mode (Tree view focused)
    |
    | Ctrl+L (jump to editor)
    v
NORMAL mode (Editor focused, tree visible)
```

**Important**: `Ctrl+H` and `Ctrl+L` only work when the tree view is already open. Use `Ctrl+T` to toggle the tree view visibility first.

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

1. Press `Ctrl+T` to open the file tree (if not already visible)
2. While editing, press `Ctrl+H` to quickly jump to the file tree
3. Navigate with `j`/`k` to find a file
4. Press `Ctrl+L` to return to editing without opening anything

### Workflow Tips

- Use `Ctrl+T` to show/hide the tree view when you need more screen space
- Once the tree is visible, use `Ctrl+H` and `Ctrl+L` to quickly switch focus
- Normal Vim navigation (`h`, `j`, `k`, `l`) works as expected in the editor
- In EXPLORER mode, `h`/`l` collapse/expand directories instead of moving the cursor

## Implementation Details

The pane navigation feature uses a robust **Command/Action pattern** for keybinding handling:

- Focus state is tracked using the `explorerFocused` boolean flag
- Mode transitions between `modeNormal` and `modeExplorer`
- Visual borders are rendered using Gio's `clip.Rect` and `paint.Fill` operations
- Global keybindings (Ctrl+T, Ctrl+H, Ctrl+L) work in ANY mode
- Mode-specific keybindings only apply to their respective modes

### Architecture

The keybinding system uses **two-phase matching**:

1. **Phase 1**: Check global keybindings (highest priority)
   - Ctrl+T (toggle explorer)
   - Ctrl+H (focus explorer)
   - Ctrl+L (focus editor)
   - Shift+Enter (toggle fullscreen)

2. **Phase 2**: Check mode-specific keybindings
   - NORMAL mode: i, v, d, h/j/k/l, etc.
   - EXPLORER mode: j/k (navigate), Enter (open), etc.
   - INSERT mode: Escape, arrow keys, etc.

3. **Phase 3**: Special handlers for complex cases
   - Count accumulation (e.g., "5j")
   - Goto sequences (e.g., "gg")
   - Colon commands (e.g., ":w")

This ensures global shortcuts always work, regardless of mode.

### Code References

- Keybinding system: `internal/appcore/keybindings.go`
- Action execution: `internal/appcore/keybindings.go:executeAction()`
- Event handling: `internal/appcore/app.go:handleKey()`
- Focus border rendering (tree): `internal/appcore/app.go:drawFileExplorer()`
- Focus border rendering (editor): `internal/appcore/app.go:drawBuffer()`

For detailed architecture documentation, see `docs/keybindings.md`.
