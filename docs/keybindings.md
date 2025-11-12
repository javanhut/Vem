# Keybindings Reference

Complete reference for all keybindings in Vem.

## Global Keybindings

These keybindings work in all modes.

| Keybinding | Action | Description |
|------------|--------|-------------|
| `Ctrl+T` | Toggle Explorer | Show/hide the file tree explorer |
| `Ctrl+H` | Focus Explorer | Switch focus to the file tree (if visible) |
| `Ctrl+L` | Focus Editor | Switch focus to the text editor |
| `Ctrl+F` | Fuzzy Finder | Open fuzzy file finder |
| `Ctrl+U` | Undo | Undo last edit operation |
| `Ctrl+C` | Copy Line | Copy current line to clipboard (NORMAL mode only) |
| `Ctrl+P` | Paste | Paste clipboard content at cursor |
| `Shift+Enter` | Toggle Fullscreen | Enter or exit fullscreen mode (NORMAL mode only) |
| `Shift+Tab` | Cycle Panes | Cycle to next pane (when multiple panes open) |
| `Ctrl+X` | Close Pane | Close active pane/buffer (like `:q`, keeps editor open) |
| `Alt+h` | Focus Left Pane | Move focus to pane on the left |
| `Alt+j` | Focus Down Pane | Move focus to pane below |
| `Alt+k` | Focus Up Pane | Move focus to pane above |
| `Alt+l` | Focus Right Pane | Move focus to pane on the right |
| ``Ctrl+` `` | Open/Toggle Terminal | Open new terminal or switch to TERMINAL INPUT mode |

## NORMAL Mode

NORMAL mode is the default mode for navigation and executing commands.

### Mode Transitions

| Key | Action | Description |
|-----|--------|-------------|
| `i` | Enter INSERT | Switch to INSERT mode at cursor |
| `v` | Enter VISUAL | Switch to VISUAL line mode |
| `d` | Enter DELETE | Switch to DELETE mode |
| `:` | Enter COMMAND | Open command-line interface |
| `Esc` | Exit Mode | Return to NORMAL (if in other mode) |

### Navigation

#### Basic Movement

| Key | Action | Description |
|-----|--------|-------------|
| `h` | Move Left | Move cursor one character left |
| `j` | Move Down | Move cursor one line down |
| `k` | Move Up | Move cursor one line up |
| `l` | Move Right | Move cursor one character right |
| `←` | Move Left | Same as `h` |
| `↓` | Move Down | Same as `j` |
| `↑` | Move Up | Same as `k` |
| `→` | Move Right | Same as `l` |

#### Line Movement

| Key | Action | Description |
|-----|--------|-------------|
| `0` | Line Start | Jump to first character of line |
| `$` or `Shift+4` | Line End | Jump to last character of line |

#### Search

| Key | Action | Description |
|-----|--------|-------------|
| `/` | Enter SEARCH | Open search prompt |
| `n` | Next Match | Jump to next search match |
| `Shift+N` | Previous Match | Jump to previous search match |

#### Document Movement

| Key | Action | Description |
|-----|--------|-------------|
| `gg` | First Line | Jump to first line of buffer |
| `G` | Last Line | Jump to last line of buffer |
| `<count>G` | Goto Line | Jump to line `<count>` (e.g., `42G`) |

### Counts

Many navigation commands accept a count prefix:

| Command | Description |
|---------|-------------|
| `5j` | Move down 5 lines |
| `10k` | Move up 10 lines |
| `3h` | Move left 3 characters |
| `42G` | Jump to line 42 |

## INSERT Mode

INSERT mode is for typing and editing text.

### Mode Control

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Exit INSERT | Return to NORMAL mode |

### Text Input

| Key | Action | Description |
|-----|--------|-------------|
| `Enter` | Newline | Insert a new line |
| `Space` | Space | Insert a space character |
| `Tab` | Insert Tab | Insert a tab character |
| `Backspace` | Delete Backward | Delete character before cursor |
| `Delete` | Delete Forward | Delete character after cursor |

### Navigation (in INSERT mode)

| Key | Action | Description |
|-----|--------|-------------|
| `←` | Move Left | Move cursor left |
| `→` | Move Right | Move cursor right |
| `↑` | Move Up | Move cursor up |
| `↓` | Move Down | Move cursor down |

## VISUAL Mode

VISUAL line mode is for selecting and manipulating multiple lines.

### Mode Control

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Exit VISUAL | Return to NORMAL mode |
| `v` | Toggle VISUAL | Exit VISUAL mode back to NORMAL |

### Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `h` | Move Left | Extend/reduce selection left |
| `j` | Move Down | Extend/reduce selection down |
| `k` | Move Up | Extend/reduce selection up |
| `l` | Move Right | Extend/reduce selection right |
| `←` | Move Left | Same as `h` |
| `↓` | Move Down | Same as `j` |
| `↑` | Move Up | Same as `k` |
| `→` | Move Right | Same as `l` |
| `0` | Line Start | Move to start of line |
| `$` or `Shift+4` | Line End | Move to end of line |
| `gg` | First Line | Extend selection to first line |
| `G` | Last Line | Extend selection to last line |

### Selection Operations

| Key | Action | Description |
|-----|--------|-------------|
| `c` | Copy | Copy selected lines to clipboard |
| `d` | Delete | Delete selected lines |
| `p` | Paste | Paste clipboard at selection |

## DELETE Mode

DELETE mode is for deleting specific lines.

### Mode Control

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Exit DELETE | Return to NORMAL mode |

### Deletion

| Key | Action | Description |
|-----|--------|-------------|
| `<count>d` | Delete Line | Delete line `<count>` (e.g., `5d` deletes line 5) |
| `d` | Delete Current | Delete current line (if no count) |

### Counts

Type digits to specify which line to delete:

| Command | Description |
|---------|-------------|
| `5d` | Delete line 5 |
| `42d` | Delete line 42 |
| `d` | Delete current line |

## COMMAND Mode

COMMAND mode provides a Vim-style command-line interface.

### Mode Control

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Cancel | Exit COMMAND mode without executing |
| `Enter` | Execute | Execute the typed command |
| `Backspace` | Delete Char | Delete character from command line |

### Commands

#### File Operations

| Command | Description |
|---------|-------------|
| `:e <file>` | Open `<file>` for editing |
| `:w` | Save current buffer to its file |
| `:w <file>` | Save current buffer as `<file>` |
| `:wq` | Save and close current buffer |
| `:q` | Close current pane/buffer (keeps editor open, switches to buffer 0 if last pane) |
| `:q!` | Force close current pane/buffer (discard changes) |
| `:qa` or `:qall` | Quit entire application (fails if buffers have unsaved changes) |
| `:qa!` | Force quit entire application (discard all unsaved changes) |

#### Buffer Management

| Command | Description |
|---------|-------------|
| `:bn` or `:bnext` | Switch to next buffer |
| `:bp` or `:bprev` | Switch to previous buffer |
| `:bd` or `:bdelete` | Close current buffer (fails if modified) |
| `:bd!` | Force close current buffer (discard changes) |
| `:ls` or `:buffers` | List all open buffers |

#### File Explorer

| Command | Description |
|---------|-------------|
| `:ex` or `:explore` | Toggle file tree explorer |
| `:cd <path>` | Change working directory to `<path>` |
| `:cd` | Change to home directory |
| `:pwd` | Print current working directory |
| `:term` or `:terminal` | Open embedded terminal in current pane |

## SEARCH Mode

SEARCH mode is for finding text within the current buffer.

### Mode Control

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Cancel Search | Exit SEARCH mode without searching |
| `Enter` | Execute Search | Execute search and jump to first match |
| `Backspace` | Delete Char | Delete character from search pattern |

### Search Pattern

Type any text to build your search pattern. The search is:
- **Case-insensitive**: `hello` matches `Hello`, `HELLO`, etc.
- **Substring matching**: `the` matches `the`, `there`, `weather`, etc.
- **Highlights all matches**: All occurrences are highlighted in the buffer

### Search Navigation

After executing a search (pressing `Enter`):

| Key | Action | Description |
|-----|--------|-------------|
| `n` | Next Match | Jump to next occurrence (wraps around) |
| `Shift+N` | Previous Match | Jump to previous occurrence (wraps around) |

### Visual Feedback

- **Yellow highlight**: All search matches
- **Orange highlight**: Current match (where cursor is)
- **Status bar**: Shows search pattern and match count (e.g., `/hello [2/5]`)

### Example Workflow

1. Press `/` in NORMAL mode to enter SEARCH mode
2. Type your search pattern (e.g., `function`)
3. Press `Enter` to execute search
4. Use `n` and `Shift+N` to navigate between matches
5. Press `Esc` in NORMAL mode to clear search highlights

## EXPLORER Mode

EXPLORER mode is for navigating the file tree.

### Mode Control

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Exit EXPLORER | Return to NORMAL mode |
| `q` | Quit | Same as `Esc` |
| `Ctrl+L` | Focus Editor | Switch focus to editor |

### Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `j` | Move Down | Move selection down |
| `k` | Move Up | Move selection up |
| `↑` | Move Up | Same as `k` |
| `↓` | Move Down | Same as `j` |

### File Operations

| Key | Action | Description |
|-----|--------|-------------|
| `Enter` | Open/Toggle | Open file or toggle directory expansion |
| `h` | Collapse | Collapse current directory |
| `l` | Expand | Expand current directory |
| `←` | Collapse | Same as `h` |
| `→` | Expand | Same as `l` |

### File Operations

| Key | Action | Description |
|-----|--------|-------------|
| `r` | Rename | Rename selected file or directory |
| `d` | Delete | Delete selected file or directory (with confirmation) |
| `n` | New File | Create a new file |

### Directory Operations

| Key | Action | Description |
|-----|--------|-------------|
| `u` | Navigate Up | Change to parent directory |

## FUZZY_FINDER Mode

FUZZY_FINDER mode provides quick file navigation using fuzzy matching.

### Mode Control

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Cancel | Exit fuzzy finder without opening a file |
| `Enter` | Open File | Open the selected file in the editor |
| `Backspace` | Delete Char | Delete character from search pattern |

### Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `↑` | Move Up | Move selection up in the results list |
| `↓` | Move Down | Move selection down in the results list |

## TERMINAL Mode

TERMINAL mode provides an embedded terminal emulator within Vem.

### Opening a Terminal

| Method | Description |
|--------|-------------|
| ``Ctrl+` `` | Open new terminal and enter TERMINAL INPUT mode immediately |
| `:term` or `:terminal` | Open new terminal and enter TERMINAL INPUT mode immediately |

### Mode Control

| Key | Action | Description |
|-----|--------|-------------|
| ``Ctrl+` `` | Open/Toggle Terminal | Open new terminal or switch to TERMINAL INPUT mode |
| `i` | Enter TERMINAL INPUT | Enter TERMINAL INPUT mode (when in terminal buffer in NORMAL mode) |
| `Esc` | Exit to NORMAL | Return to NORMAL mode (can navigate terminal output) |
| `Shift+Tab` | Exit to NORMAL | Alternative way to return to NORMAL mode (also cycles panes) |
| `Ctrl+X` | Close Terminal | Close terminal pane/buffer (works in TERMINAL INPUT mode) |

### Terminal Features

- Full VT100/xterm-256color terminal emulation with ANSI color support
- Runs your default shell (bash, zsh, fish, etc.)
- Properly renders colors, bold, italic, underline attributes
- Integrates with Vem's pane system (can split terminals)
- Automatically starts in current working directory
- Easy access with ``Ctrl+` `` keybinding

### Usage

**Quick Start:**
1. Press ``Ctrl+` `` → Terminal opens and enters TERMINAL INPUT mode immediately
2. Type commands → Works like a normal terminal
3. Press `Esc` → Return to NORMAL mode
4. Press ``Ctrl+` `` or `i` → Re-enter TERMINAL INPUT mode
5. Press `Shift+Tab` → Switch to other panes/buffers

**Alternative (Command Mode):**
1. Type `:term` and press Enter → Opens terminal in TERMINAL INPUT mode
2. Type commands as you would in a normal terminal
3. Press `Esc` → Return to NORMAL mode
4. Press `i` → Resume typing commands

### Terminal Input (When in TERMINAL INPUT Mode)

All keyboard input is sent directly to the terminal:

- Arrow keys, function keys, and special keys work as expected
- Ctrl+key combinations are sent to the terminal (except ``Ctrl+` `` which toggles mode)
- Alt+key combinations are sent to the terminal
- Tab completion works normally
- Colors and text attributes are properly displayed

### Example Workflow

1. Press ``Ctrl+` `` → Opens terminal in TERMINAL INPUT mode
2. Type `ls -la` → Run command with colored output
3. Type `git status` → Check git status
4. Press `Esc` → Return to NORMAL mode (output visible, can't type)
5. Use `:e file.go` → Open a file
6. Press `Shift+Tab` → Cycle back to terminal
7. Press `i` or ``Ctrl+` `` → Resume typing in terminal
8. Use `Ctrl+S v` → Split vertically for side-by-side terminal and editor
9. Press `Alt+h/j/k/l` → Navigate between panes
10. Press ``Ctrl+` `` while in another pane → Quickly jump back to terminal

### Search Pattern

Type any text to filter files. The fuzzy matcher:
- Matches characters in sequence (not necessarily consecutive)
- Prioritizes matches at word boundaries
- Ranks shorter paths higher
- Shows up to 50 best matches

### Visual Feedback

- **Blue overlay**: Semi-transparent background showing fuzzy finder is active
- **Blue border**: Fuzzy finder box
- **Highlighted row**: Currently selected file (blue background)
- **Match count**: Shows total number of matching files

### Example Workflow

1. Press `Ctrl+F` from any mode to open fuzzy finder
2. Type partial file name (e.g., `bufgo` to find `internal/editor/buffer.go`)
3. Use `↑`/`↓` to navigate through matches
4. Press `Enter` to open the selected file
5. Press `Esc` to cancel

### Excluded Directories

The fuzzy finder automatically excludes:
- Hidden directories (starting with `.`)
- `node_modules`
- `vendor`
- `.git`
- `.gocache`
- `dist`
- `build`
- `target`

## Special Sequences

### Goto Sequences

The `g` key initiates a goto sequence:

| Sequence | Description |
|----------|-------------|
| `gg` | Jump to first line |
| `<count>gg` | Jump to line `<count>` |
| `gG` | Jump to last line (same as `G`) |

### Count Accumulation

Type digits to build a count before a motion or command:

| Example | Description |
|---------|-------------|
| `5j` | Move down 5 lines |
| `10k` | Move up 10 lines |
| `42G` | Jump to line 42 |
| `5dd` | Delete 5 lines starting at line 5 |

## Platform-Specific Notes

### Modifier Keys

Vem uses robust modifier key tracking to handle platform differences:

- **Linux (X11/Wayland)**: Ctrl and Shift are tracked explicitly due to platform quirks
- **macOS**: Command key is mapped to Ctrl for consistency
- **Windows**: Ctrl key works as expected

### Fullscreen Behavior

The fullscreen toggle (`Shift+Enter`) may behave differently depending on the platform:

- **Linux**: Uses window manager's fullscreen mode
- **macOS**: Uses native fullscreen (separate space)
- **Windows**: Uses borderless maximized window

## Customization (Future)

In future releases, keybindings will be customizable via:

- Configuration file (`~/.vemrc`)
- Lua/Python/Carrion scripts
- GUI keybinding editor

See the [ROADMAP.md](../ROADMAP.md) for planned customization features.

## Pane Management (Ctrl+S prefix)

Vem provides powerful pane splitting and management capabilities. Press `Ctrl+S` to enter pane command mode, then press one of the following keys:

| Keybinding | Action | Description |
|------------|--------|-------------|
| `Ctrl+S v` | Split Vertical | Create vertical split (top/bottom with horizontal divider) |
| `Ctrl+S h` | Split Horizontal | Create horizontal split (left/right with vertical divider) |
| `Ctrl+S =` | Equalize Panes | Make all panes equal size (50/50) |
| `Ctrl+S o` | Zoom Toggle | Maximize/restore active pane |

**Note**: After splitting, use `:e filename` or `Ctrl+P` to open a file in the new pane.

For detailed pane usage, see [Pane Splitting Guide](pane-splitting.md).

## Clipboard Operations

Vem provides both system clipboard integration and mode-specific clipboard operations:

### System Clipboard (Ctrl+C / Ctrl+P)

| Keybinding | Mode | Action | Description |
|------------|------|--------|-------------|
| `Ctrl+C` | NORMAL | Copy Line | Copy current line to system clipboard |
| `Ctrl+P` | NORMAL, INSERT | Paste | Paste clipboard content at cursor position |

The system clipboard integration:
- Works with your operating system's clipboard
- Allows copying/pasting between Vem and other applications
- In NORMAL mode, `Ctrl+C` copies the entire current line
- `Ctrl+P` pastes at the cursor position in both NORMAL and INSERT modes
- Automatically syncs with Visual mode copy operations

### Visual Mode Clipboard (c / p keys)

| Keybinding | Mode | Action | Description |
|------------|------|--------|-------------|
| `c` | VISUAL | Copy Selection | Copy selected text/lines to clipboard |
| `p` | VISUAL | Paste | Replace selection with clipboard content |

Visual mode clipboard operations:
- `c` in Visual mode copies selection to both system and internal clipboard
- `p` in Visual mode replaces selection with clipboard content
- Works with both character-wise (`v`) and line-wise (`Shift+V`) selections

Example usage:
1. **Copy current line**: In NORMAL mode, press `Ctrl+C`
2. **Paste at cursor**: Press `Ctrl+P` to paste
3. **Copy selection**: Press `v` to enter VISUAL mode, select text, press `c`
4. **Paste in another app**: After copying in Vem, paste in any other application using system shortcuts

## Undo Functionality

Vem provides undo functionality for text editing operations:

| Keybinding | Action | Description |
|------------|--------|-------------|
| `Ctrl+U` | Undo | Undo the last edit operation (insert, delete, etc.) |

The undo system:
- Tracks up to 100 edit operations
- Works for insertions, deletions, and line operations
- Available in all modes (most useful in NORMAL and INSERT modes)
- Shows "Nothing to undo" when the undo stack is empty

Example usage:
1. Make some edits in INSERT mode
2. Press `Esc` to return to NORMAL mode
3. Press `Ctrl+U` to undo the last change
4. Press `Ctrl+U` repeatedly to undo multiple changes

## Keybinding Conflicts

Vem uses a priority-based keybinding system to prevent conflicts:

1. **Global keybindings** (highest priority): Ctrl+T, Ctrl+H, Ctrl+L, Ctrl+F, Ctrl+U, Alt+hjkl, Shift+Tab, Ctrl+X
2. **Mode-specific keybindings**: Commands that only work in specific modes (e.g., Shift+Enter for fullscreen only in NORMAL mode)
3. **Special handlers**: Complex sequences like `gg`, counts, Ctrl+S prefix, etc.

Exception: COMMAND mode keybindings take priority over global keybindings to ensure commands execute correctly.

## Architecture

For details on how the keybinding system works internally, see [Architecture.md](Architecture.md).

## See Also

- [Pane Splitting Guide](pane-splitting.md) - Complete guide to pane management
- [Navigation Guide](navigation.md) - Pane navigation and fullscreen features
- [Tutorial](tutorial.md) - Step-by-step guide for new users
- [Architecture](Architecture.md) - Technical implementation details
