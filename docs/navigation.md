# Navigation in Vem

This document describes the navigation features in Vem, including cursor movement, viewport scrolling, pane navigation between the file tree explorer and the text editor, and fullscreen mode management.

## Cursor Movement

Vem provides Vim-compatible cursor movement commands for efficient navigation.

### Character Movement

- **`h`** or **Left Arrow** - Move cursor left one character
- **`l`** or **Right Arrow** - Move cursor right one character
- **`j`** or **Down Arrow** - Move cursor down one line
- **`k`** or **Up Arrow** - Move cursor up one line

### Line Movement

- **`0`** - Jump to start of line (column 0)
- **`$`** - Jump to end of line
- **`gg`** - Jump to first line of buffer
- **`G`** - Jump to last line of buffer
- **`[count]G`** - Jump to line number [count] (e.g., `42G` jumps to line 42)

### Word Movement

Word navigation allows efficient movement through code and text:

- **`w`** - Move forward to start of next **word**
  - A word is a sequence of letters, digits, and underscores OR other non-blank characters
  - Examples: `foo_bar` (1 word), `foo-bar` (3 words: `foo`, `-`, `bar`)
  
- **`b`** - Move backward to start of previous **word**
  - Moves to the beginning of the current or previous word
  
- **`e`** - Move forward to end of current or next **word**
  - Positions cursor on the last character of the word

**Word Definition:**
- **Word characters**: Letters (a-z, A-Z), digits (0-9), underscore (_)
- **Punctuation**: All other non-whitespace characters are treated as separate words
- **Whitespace**: Spaces, tabs, newlines separate words

**Examples:**
```
Text: "foo_bar = get_value(123);"

Starting at 'f' in foo_bar:
w → moves to '=' (start of next word)
w → moves to 'g' in get_value
w → moves to '(' (punctuation is separate word)
w → moves to '1' in 123
e → moves to '3' (end of 123)
b → moves to '1' (start of previous word)
```

### Search Movement

- **`/`** - Enter search mode (type pattern and press Enter)
- **`n`** - Jump to next search match
- **`N`** - Jump to previous search match

## Visual Mode

Vem implements Vim-style visual selection with two distinct modes for selecting text: character-wise and line-wise.

### Visual Mode Types

#### Character-Wise Visual Mode (`v`)

- **`v`** - Enter character-wise visual mode
  - Allows precise, per-character selection
  - Selection extends from the anchor point to the cursor position
  - Works across multiple lines
  - Visual indication: selected characters are highlighted with a colored background
  - Status bar shows: `-- VISUAL --`

**Usage:**
```
v           # Start character selection at cursor
w           # Extend selection forward by word
e           # Extend to end of word
b           # Contract selection backward by word
h/j/k/l     # Extend/contract by character or line
c           # Copy selected text
d           # Delete selected text
p           # Paste (replace selection with clipboard)
Escape      # Exit visual mode
```

#### Line-Wise Visual Mode (`Shift+V`)

- **`Shift+V`** - Enter line-wise visual mode
  - Selects entire lines at a time
  - Selection always includes complete lines
  - Visual indication: selected lines are highlighted with a colored background
  - Status bar shows: `-- VISUAL LINE --`

**Usage:**
```
Shift+V     # Start line selection at cursor line
j/k         # Extend selection down/up by line
c           # Copy selected lines
d           # Delete selected lines
p           # Paste (insert lines at selection)
Escape      # Exit visual mode
```

### Visual Mode Operations

Once text is selected in either visual mode, you can perform these operations:

- **`c`** - Copy selection to clipboard
  - Character mode: copies exact character range
  - Line mode: copies complete lines
  - Status shows: "Copied N character(s)" or "Copied N line(s)"

- **`d`** - Delete selection
  - Character mode: deletes character range, cursor moves to start
  - Line mode: deletes complete lines, cursor moves to start line
  - Status shows: "Deleted selection"

- **`p`** - Paste clipboard, replacing selection
  - Character mode: replaces selected characters with clipboard text
  - Line mode: replaces selected lines with clipboard lines
  - Status shows: "Pasted N character(s)" or "Inserted N line(s)"

- **`Escape`** - Exit visual mode without making changes
  - Returns to NORMAL mode
  - Selection is cleared
  - Cursor position is preserved

### Visual Mode Highlighting

Vem uses different highlighting strategies for character-wise and line-wise visual modes to maximize clarity:

#### Character-Wise Mode Highlighting

In character-wise visual mode (`v`), only the selected characters are highlighted:
- **Selection highlight**: Purple background on selected characters only
- **Cursor line highlight**: Disabled to show precise character selection
- **Cursor indicator**: Block cursor remains visible to show current position

This ensures you can see exactly which characters are selected without visual confusion.

**Example:**
```
v + 2w (select two words)
Line 10: int foo = get_value();
         ^^^^^^^^                ← Only "int foo " is highlighted (purple)
```

#### Line-Wise Mode Highlighting

In line-wise visual mode (`Shift+V`), entire lines are highlighted:
- **Selection highlight**: Purple background on entire selected lines
- **Cursor line highlight**: Blue background on cursor line (overlays with selection)
- **Visual indication**: Full line width is highlighted

**Example:**
```
Shift+V + 2j (select three lines)
Line 10: int foo = get_value();
████████████████████████████████ ← Entire line highlighted (purple + blue)
Line 11: int bar = 42;
████████████████████████████████ ← Entire line highlighted (purple)
Line 12: return foo + bar;
████████████████████████████████ ← Entire line highlighted (purple)
```

#### Normal Mode Highlighting

In NORMAL and INSERT modes, the cursor line has a subtle blue highlight:
- **Cursor line highlight**: Blue background on the line containing the cursor
- **Purpose**: Helps identify current cursor position during navigation

### Visual Mode Navigation

All normal mode navigation commands work in visual mode to extend or contract the selection:

**Character Movement:**
- `h`, `l`, `←`, `→` - Extend/contract by character
- `j`, `k`, `↓`, `↑` - Extend/contract by line

**Word Movement:**
- `w` - Extend forward to next word start
- `b` - Contract backward to previous word start
- `e` - Extend forward to word end

**Line Movement:**
- `0` - Extend to line start
- `$` - Extend to line end
- `gg` - Extend to first line
- `G` - Extend to last line
- `[count]G` - Extend to line number

### Visual Mode Examples

**Select and copy a word:**
```
v           # Start character selection
w           # Extend to next word
c           # Copy to clipboard
```

**Select and delete multiple lines:**
```
Shift+V     # Start line selection
3j          # Extend down 3 lines (total 4 lines selected)
d           # Delete all selected lines
```

**Replace a function name across lines:**
```
v           # Start character selection
e           # Extend to end of current word
c           # Copy old name
/new_func   # Search for location to replace
v           # Start new selection
w           # Select word
p           # Paste, replacing selection
```

**Select code block and copy:**
```
Shift+V     # Start line selection at function start
}           # Jump to end of block (when implemented)
c           # Copy entire block
```

### Visual Mode vs Vem Keybindings

Vem uses **Vem-style copy/delete/paste** keybindings instead of Vim's traditional `y` (yank):

| Operation | Vem Key | Vim Key | Notes |
|-----------|---------|---------|-------|
| Copy      | `c`     | `y`     | Vem: "copy" is more intuitive |
| Delete    | `d`     | `d`     | Same in both |
| Paste     | `p`     | `p`     | Same in both |

This makes the keybindings more intuitive for new users while maintaining the Vim-style modal editing workflow.

## Viewport Scrolling

Vem implements Vim-style viewport scrolling to ensure the cursor is always visible and to provide fine-grained control over the viewport position.

### Automatic Scrolling

The viewport automatically scrolls to keep the cursor visible with a configurable scroll offset (default: 3 lines of context above/below the cursor). This happens automatically when you:

- Jump to a line with `gg`, `G`, or `[count]G`
- Navigate with `h`, `j`, `k`, `l` or arrow keys
- Move by word with `w`, `b`, `e`
- Search and jump to matches with `/` and `n`/`N`
- Use any other navigation command

### Manual Scroll Commands

#### Cursor Positioning Commands

- **`zz`** - Center cursor line in viewport (Vim's `zz`)
  - Scrolls the viewport so the cursor line is in the middle of the screen
  - Cursor position doesn't change, only the viewport
  - Status: "Centered cursor"

- **`zt`** - Position cursor line at top of viewport (Vim's `zt`)
  - Scrolls the viewport so the cursor line is at the top
  - Cursor position doesn't change, only the viewport
  - Status: "Cursor at top"

- **`zb`** - Position cursor line at bottom of viewport (Vim's `zb`)
  - Scrolls the viewport so the cursor line is at the bottom
  - Cursor position doesn't change, only the viewport
  - Status: "Cursor at bottom"

#### Line-by-Line Scrolling

- **`Ctrl+E`** - Scroll viewport down one line (Vim's `Ctrl+E`)
  - Moves the viewport down by one line
  - Cursor position doesn't change (unless it goes off-screen)
  - Status: "Scrolled down (top line: N)"

- **`Ctrl+Y`** - Scroll viewport up one line (Vim's `Ctrl+Y`)
  - Moves the viewport up by one line
  - Cursor position doesn't change (unless it goes off-screen)
  - Status: "Scrolled up (top line: N)"

### Scroll Offset

By default, Vem maintains a scroll offset of 3 lines around the cursor. This means:

- When scrolling up, the cursor stays at least 3 lines from the top edge
- When scrolling down, the cursor stays at least 3 lines from the bottom edge
- This provides visual context and prevents the cursor from being at the screen edge

### Viewport Scrolling in Different Modes

All scroll commands work in:
- **NORMAL mode** - Full scroll functionality available
- **VISUAL mode** - Scroll while maintaining selection (same commands)
- **INSERT mode** - Automatic scrolling only (manual scroll commands not available)

### Examples

**Jump to top of file and center:**
```
gg      # Jump to line 1
zz      # Center line 1 in viewport
```

**Navigate to a specific line and position at top:**
```
42G     # Jump to line 42
zt      # Position line 42 at top of viewport
```

**Fine-tune viewport position:**
```
/search_term    # Search for something
Ctrl+E          # Scroll down 1 line to see more context below
Ctrl+E          # Scroll down another line
```

**Quick viewport adjustments:**
```
zz      # Center current line
zt      # Move to top
zb      # Move to bottom
```

## Fullscreen Mode

### Shift+Enter - Toggle Fullscreen

Vem supports fullscreen mode for a distraction-free editing experience.

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

Vem supports Vim-like split pane navigation with keyboard shortcuts to quickly jump between the file tree explorer and the text editor.

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
- If the tree view is hidden, it will be shown (stays in current mode)
- If the tree view is visible, it will be hidden and focus returns to the editor
- This is the ONLY way to show/hide the file tree
- Note: Plain 't' key does NOT toggle the tree (removed to avoid confusion)

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
- Viewport scrolling: `internal/appcore/app.go:ensureCursorVisible()`, `internal/appcore/app.go:drawBuffer()`
- Key input handling: `internal/appcore/app.go:printableKey()` (shift-aware letter case conversion)
- Focus border rendering (tree): `internal/appcore/app.go:drawFileExplorer()`
- Focus border rendering (editor): `internal/appcore/app.go:drawBuffer()`

For detailed architecture documentation, see `docs/keybindings.md`.

## Technical Notes

### Keyboard Input Handling

Vem handles keyboard input with special attention to platform quirks:

- **Letter Case Handling**: Gio UI reports all letter keys as uppercase in the key name (e.g., `ev.Name = "G"` for both `g` and `Shift+G`). The `printableKey()` method checks the tracked shift state to correctly convert letters to lowercase when shift is not pressed, ensuring commands like `gg` (jump to top) and `G` (jump to bottom) work correctly.

- **Modifier Tracking**: Due to platform limitations where `ev.Modifiers` is not reliably reported, Vem tracks modifier key state (Ctrl, Shift) through explicit press/release events in `s.ctrlPressed` and `s.shiftPressed`.

- **Smart Reset**: After executing commands, modifiers are automatically reset to prevent them from "sticking" between keypress events.
