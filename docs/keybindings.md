# Keybindings Reference

Complete reference for all keybindings in ProjectVem.

## Global Keybindings

These keybindings work in all modes.

| Keybinding | Action | Description |
|------------|--------|-------------|
| `Ctrl+T` | Toggle Explorer | Show/hide the file tree explorer |
| `Ctrl+H` | Focus Explorer | Switch focus to the file tree (if visible) |
| `Ctrl+L` | Focus Editor | Switch focus to the text editor |
| `Shift+Enter` | Toggle Fullscreen | Enter or exit fullscreen mode |

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
| `$` | Line End | Jump to last character of line |

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
| `$` | Line End | Move to end of line |
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
| `:q` | Close current buffer (fails if modified) |
| `:q!` | Force close current buffer (discard changes) |

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

### Directory Operations

| Key | Action | Description |
|-----|--------|-------------|
| `r` | Refresh | Reload file tree from disk |
| `u` | Navigate Up | Change to parent directory |

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

ProjectVem uses robust modifier key tracking to handle platform differences:

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

## Keybinding Conflicts

ProjectVem uses a priority-based keybinding system to prevent conflicts:

1. **Global keybindings** (highest priority): Ctrl+T, Ctrl+H, Ctrl+L, Shift+Enter
2. **Mode-specific keybindings**: Commands that only work in specific modes
3. **Special handlers**: Complex sequences like `gg`, counts, etc.

Exception: COMMAND mode keybindings take priority over global keybindings to ensure commands execute correctly.

## Architecture

For details on how the keybinding system works internally, see [Architecture.md](Architecture.md).

## See Also

- [Navigation Guide](navigation.md) - Pane navigation and fullscreen features
- [Tutorial](tutorial.md) - Step-by-step guide for new users
- [Architecture](Architecture.md) - Technical implementation details
