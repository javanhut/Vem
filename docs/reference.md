# Vem Reference Guide

Complete reference documentation for Vem text editor.

## Table of Contents

- [Command Line Usage](#command-line-usage)
- [Modes](#modes)
- [Keybindings](#keybindings)
- [Commands](#commands)
- [Buffer Management](#buffer-management)
- [File Operations](#file-operations)
- [Search and Navigation](#search-and-navigation)
- [Terminal](#terminal)
- [Pane Management](#pane-management)
- [Configuration](#configuration)

## Command Line Usage

### Basic Usage

```bash
# Launch with sample buffer
vem

# Open a single file
vem filename.txt

# Open multiple files (first file becomes active)
vem file1.go file2.md file3.txt

# Create a new file (opens empty buffer)
vem newfile.txt

# Open files with paths
vem /path/to/file.go
vem ../relative/path.txt
```

### Behavior

- **No arguments**: Opens with sample buffer showing welcome text
- **Single file**: Opens file in buffer 0, or creates empty buffer if file doesn't exist
- **Multiple files**: Opens all files, first file is active, switch with `:bn`/`:bp`
- **Invalid paths**: Logs warning and skips the file
- **All files fail**: Falls back to sample buffer

## Modes

Vem uses modal editing similar to Vim. Each mode serves a specific purpose.

### NORMAL Mode

**Purpose**: Navigation and command execution

**Indicator**: `MODE NORMAL` in status bar

**Key Characteristics**:
- Default mode on startup
- All navigation commands available
- Can switch to other modes
- Cannot insert text directly

**Switching to NORMAL**:
- Press `Esc` from any mode
- Automatic after command execution

### INSERT Mode

**Purpose**: Text insertion and editing

**Indicator**: `MODE INSERT` in status bar

**Key Characteristics**:
- All typed characters insert into buffer
- Arrow keys for navigation
- Backspace/Delete for deletion
- Tab inserts tab character
- Enter creates new line

**Entering INSERT Mode**:
- Press `i` in NORMAL mode
- Cursor position remains unchanged

**Exiting INSERT Mode**:
- Press `Esc` to return to NORMAL mode

### VISUAL Mode

**Purpose**: Text selection for copy/delete operations

**Indicator**: `MODE VISUAL` or `MODE VISUAL LINE` in status bar

**Types**:
- **Character-wise** (`v`): Select individual characters
- **Line-wise** (`Shift+V`): Select entire lines

**Key Characteristics**:
- Highlighted selection region
- Navigation extends selection
- Operations apply to selection

**Entering VISUAL Mode**:
- Press `v` for character-wise
- Press `Shift+V` for line-wise

**Operations**:
- `c` - Copy selection to clipboard
- `d` - Delete selection
- `p` - Paste from clipboard

**Exiting VISUAL Mode**:
- Press `Esc` or `v` to return to NORMAL mode

### COMMAND Mode

**Purpose**: Execute editor commands

**Indicator**: `:` prompt in status bar

**Key Characteristics**:
- Vim-style colon commands
- File operations
- Buffer management
- Settings changes

**Entering COMMAND Mode**:
- Press `:` in NORMAL mode

**Command Syntax**:
```
:command [arguments]
```

**Exiting COMMAND Mode**:
- Press `Enter` to execute
- Press `Esc` to cancel

### EXPLORER Mode

**Purpose**: File tree navigation and file operations

**Indicator**: `MODE EXPLORER` in status bar

**Key Characteristics**:
- Tree-based directory view
- File/directory operations
- Visual navigation with j/k
- Nerd Font icons for file types

**Entering EXPLORER Mode**:
- Press `Ctrl+T` to toggle
- Press `Ctrl+H` to focus
- Use `:ex` or `:explore` command

**Operations**:
- Navigate with `j`/`k`
- Open file/toggle directory with `Enter`
- Create file with `n`
- Rename with `r`
- Delete with `d`
- Navigate up with `u`

**Exiting EXPLORER Mode**:
- Press `Esc` to return to NORMAL mode
- Press `q` to close explorer
- Press `Ctrl+L` to focus editor

### SEARCH Mode

**Purpose**: Find text in current buffer

**Indicator**: `/` prompt in status bar with search pattern

**Key Characteristics**:
- Case-insensitive search
- Real-time match highlighting
- Pattern building
- Navigate matches with n/N

**Entering SEARCH Mode**:
- Press `/` in NORMAL mode

**Building Pattern**:
- Type characters to add to pattern
- Press `Backspace` to remove last character

**Executing Search**:
- Press `Enter` to find first match
- Use `n` for next match (in NORMAL mode)
- Use `Shift+N` for previous match

**Exiting SEARCH Mode**:
- Press `Enter` to execute and return to NORMAL
- Press `Esc` to cancel

### FUZZY_FINDER Mode

**Purpose**: Fast file finding with fuzzy matching

**Indicator**: `FUZZY FINDER` in status bar with query

**Key Characteristics**:
- Searches entire project directory
- Fuzzy matching algorithm
- Real-time results filtering
- Shows matching score

**Entering FUZZY_FINDER Mode**:
- Press `Ctrl+F` in any mode

**Using Fuzzy Finder**:
- Type characters to filter files
- Use `↑`/`↓` to navigate results
- Press `Enter` to open selected file
- Press `Backspace` to remove characters

**Exiting FUZZY_FINDER Mode**:
- Press `Enter` to open file
- Press `Esc` to cancel

### TERMINAL Mode

**Purpose**: Integrated shell access

**Indicator**: `MODE TERMINAL` in status bar

**Key Characteristics**:
- Full VT100/ANSI terminal emulator
- 256-color and true color support
- Bold, italic, underline attributes
- Runs in buffer system
- Auto-closes on shell exit

**Entering TERMINAL Mode**:
- Press `Ctrl+` ` (backtick) to toggle

**Terminal Features**:
- All input passes to shell
- Full escape sequence support
- Color output rendering
- PTY integration (Unix/Windows)

**Exiting TERMINAL Mode**:
- Press `Esc` to return to NORMAL mode
- Press `Ctrl+` ` to close terminal
- Terminal auto-closes when shell exits

## Keybindings

### Global Keybindings

Work in all modes:

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+T` | Toggle Explorer | Show/hide file tree |
| `Ctrl+H` | Focus Explorer | Switch to file tree |
| `Ctrl+L` | Focus Editor | Switch to editor pane |
| `Ctrl+F` | Fuzzy Finder | Quick file search |
| `Ctrl+U` | Undo | Undo last operation |
| `Ctrl+X` | Close Pane | Close active pane/buffer |
| `Ctrl+` ` | Toggle Terminal | Open/close terminal |
| `Shift+Enter` | Fullscreen | Toggle fullscreen mode |
| `Shift+Tab` | Cycle Panes | Move to next pane |

### Pane Management

Press `Ctrl+S` followed by:

| Key | Action | Description |
|-----|--------|-------------|
| `v` | Split Vertical | Create top/bottom split |
| `h` | Split Horizontal | Create left/right split |
| `=` | Equalize | Make all panes equal size |
| `o` | Zoom Toggle | Maximize/restore active pane |

### Pane Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `Alt+h` | Focus Left | Move to left pane |
| `Alt+j` | Focus Down | Move to pane below |
| `Alt+k` | Focus Up | Move to pane above |
| `Alt+l` | Focus Right | Move to right pane |

### NORMAL Mode Keybindings

#### Mode Switching

| Key | Mode | Description |
|-----|------|-------------|
| `i` | INSERT | Enter insert mode at cursor |
| `v` | VISUAL | Character-wise visual mode |
| `Shift+V` | VISUAL LINE | Line-wise visual mode |
| `d` | DELETE | Enter delete mode |
| `:` | COMMAND | Open command prompt |
| `/` | SEARCH | Start search |

#### Cursor Movement

| Key | Movement | Description |
|-----|----------|-------------|
| `h` | Left | Move one character left |
| `j` | Down | Move one line down |
| `k` | Up | Move one line up |
| `l` | Right | Move one character right |
| `w` | Word Forward | Jump to next word start |
| `b` | Word Backward | Jump to previous word start |
| `e` | Word End | Jump to end of current word |
| `0` | Line Start | Jump to beginning of line |
| `$` | Line End | Jump to end of line |
| `gg` | First Line | Jump to top of file |
| `G` | Last Line | Jump to bottom of file |
| `<n>G` | Goto Line | Jump to line number n |

#### Viewport Scrolling

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+E` | Scroll Down | Scroll viewport down one line |
| `Ctrl+Y` | Scroll Up | Scroll viewport up one line |

#### Search Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `n` | Next Match | Jump to next search result |
| `Shift+N` | Previous Match | Jump to previous search result |
| `Esc` | Clear Search | Clear search highlights |

### INSERT Mode Keybindings

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Exit | Return to NORMAL mode |
| `Enter` | New Line | Insert newline and move cursor |
| `Tab` | Insert Tab | Insert tab character |
| `Backspace` | Delete Back | Delete character before cursor |
| `Delete` | Delete Forward | Delete character at cursor |
| Arrow keys | Navigate | Move cursor while in INSERT mode |
| Any character | Insert | Insert character at cursor |

### VISUAL Mode Keybindings

#### Navigation (Extends Selection)

| Key | Movement | Description |
|-----|----------|-------------|
| `h` | Left | Extend selection left |
| `j` | Down | Extend selection down |
| `k` | Up | Extend selection up |
| `l` | Right | Extend selection right |
| `w` | Word Forward | Extend to next word |
| `b` | Word Backward | Extend to previous word |
| `e` | Word End | Extend to word end |
| `0` | Line Start | Extend to line start |
| `$` | Line End | Extend to line end |
| `gg` | First Line | Extend to file start |
| `G` | Last Line | Extend to file end |

#### Operations

| Key | Action | Description |
|-----|--------|-------------|
| `c` | Copy | Copy selection to clipboard |
| `d` | Delete | Delete selected text |
| `p` | Paste | Paste from clipboard |
| `v` | Exit | Return to NORMAL mode |
| `Esc` | Exit | Return to NORMAL mode |

### EXPLORER Mode Keybindings

| Key | Action | Description |
|-----|--------|-------------|
| `j` | Down | Move selection down |
| `k` | Up | Move selection up |
| `Enter` | Open/Toggle | Open file or expand/collapse directory |
| `h` | Collapse | Collapse current directory |
| `l` | Expand | Expand current directory |
| `r` | Rename | Rename file/directory |
| `d` | Delete | Delete file/directory |
| `n` | New File | Create new file |
| `u` | Navigate Up | Move to parent directory |
| `q` | Quit | Close explorer |
| `Esc` | Exit | Return to NORMAL mode |

### SEARCH Mode Keybindings

| Key | Action | Description |
|-----|--------|-------------|
| Any character | Add to Pattern | Build search pattern |
| `Backspace` | Delete Char | Remove last character |
| `Enter` | Execute | Find first match |
| `Esc` | Cancel | Exit without searching |

### FUZZY_FINDER Mode Keybindings

| Key | Action | Description |
|-----|--------|-------------|
| Any character | Filter | Add character to filter |
| `↑` or `k` | Previous | Select previous file |
| `↓` or `j` | Next | Select next file |
| `Enter` | Open | Open selected file |
| `Backspace` | Delete Char | Remove last character |
| `Esc` | Cancel | Close fuzzy finder |

### TERMINAL Mode Keybindings

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Exit to NORMAL | Return to NORMAL mode |
| `Ctrl+` ` | Close Terminal | Close terminal and buffer |
| All others | Shell Input | Pass directly to shell |

## Commands

All commands start with `:` in NORMAL mode.

### File Operations

| Command | Arguments | Description |
|---------|-----------|-------------|
| `:e` | `<file>` | Open file for editing |
| `:w` | None | Save current buffer |
| `:w` | `<file>` | Save buffer as new file |
| `:wq` | None | Save and close buffer |
| `:q` | None | Close buffer (fails if unsaved) |
| `:q!` | None | Force close buffer (discard changes) |
| `:qa` | None | Quit all (close editor) |
| `:qa!` | None | Force quit all (discard all changes) |
| `:qall` | None | Quit all (alias for :qa) |

### Buffer Management

| Command | Arguments | Description |
|---------|-----------|-------------|
| `:bn` | None | Switch to next buffer |
| `:bnext` | None | Switch to next buffer (alias) |
| `:bp` | None | Switch to previous buffer |
| `:bprev` | None | Switch to previous buffer (alias) |
| `:bd` | None | Close current buffer |
| `:bdelete` | None | Close current buffer (alias) |
| `:bd!` | None | Force close buffer |
| `:ls` | None | List all open buffers |
| `:buffers` | None | List all open buffers (alias) |

### File Explorer

| Command | Arguments | Description |
|---------|-----------|-------------|
| `:ex` | None | Toggle file explorer |
| `:explore` | None | Toggle file explorer (alias) |
| `:cd` | `<path>` | Change working directory |
| `:pwd` | None | Print working directory |

### Help System

| Command | Arguments | Description |
|---------|-----------|-------------|
| `:help` | None | Open comprehensive help buffer |
| `:h` | None | Open help (alias) |

### Pane Management

Pane operations are primarily handled through keybindings (`Ctrl+S` prefix), but some operations can be done via commands.

**Note**: Most pane operations use `Ctrl+S` + key combinations. See Pane Management Keybindings section.

## Buffer Management

### Buffer Lifecycle

1. **Creation**:
   - Command line: `vem file.txt`
   - `:e <file>` command
   - Explorer: Press Enter on file
   - Fuzzy finder: Select and open file
   - Terminal: `Ctrl+` ` creates terminal buffer

2. **Active Buffer**:
   - Only one buffer active at a time per pane
   - Shown in status bar: `FILE filename.txt`
   - Changes apply to active buffer only

3. **Switching**:
   - `:bn` / `:bp` - Navigate buffers
   - `:ls` - List all buffers
   - Numbers in `:ls` can be used with `:b<n>`

4. **Closing**:
   - `:q` - Close if no unsaved changes
   - `:q!` - Force close (discard changes)
   - `:bd` - Close buffer, keep editor open
   - Terminal buffers close automatically on shell exit

### Buffer Types

1. **Text Buffers**:
   - Standard file editing
   - Can be modified
   - Require save before close
   - Shows `+` if modified

2. **Terminal Buffers**:
   - Created with `Ctrl+` `
   - No unsaved changes warning
   - Auto-close on shell exit
   - Shows `[Terminal]` in buffer list

3. **Read-Only Buffers**:
   - Help pages (`:help`)
   - Shows `[RO]` in status bar
   - Cannot enter INSERT mode
   - Cannot be modified

4. **Special Buffers**:
   - Sample buffer (no file path)
   - Shows `[No Name]`
   - Can be saved with `:w <filename>`

### Buffer Indicators

In status bar:
- `FILE filename.txt` - Normal buffer
- `FILE filename.txt +` - Modified buffer
- `FILE [No Name]` - Unnamed buffer
- `FILE [Terminal]` - Terminal buffer
- `FILE helpfile.txt [RO]` - Read-only buffer

In `:ls` output:
```
* 1 + main.go
  2   README.md
  3   [Terminal]
  4   [No Name]
```
- `*` - Active buffer
- Number - Buffer index
- `+` - Modified
- Name - File path or special name

## File Operations

### Opening Files

**From Command Line**:
```bash
vem file1.txt file2.go
```

**From Editor**:
- `:e filename.txt` - Open specific file
- Explorer (`Ctrl+T`) - Navigate and press Enter
- Fuzzy Finder (`Ctrl+F`) - Type to filter, Enter to open

**Behavior**:
- If file exists: Loads content into buffer
- If file doesn't exist: Creates empty buffer with path
- If already open: Switches to existing buffer

### Saving Files

**Save Current Buffer**:
```
:w
```

**Save As**:
```
:w newname.txt
```

**Save and Close**:
```
:wq
```

**Behavior**:
- Creates file if it doesn't exist
- Overwrites existing file
- Updates modification time
- Clears modified flag
- Terminal buffers cannot be saved

### Creating Files

**Method 1: Command Line**
```bash
vem newfile.txt
```

**Method 2: Explorer**
- Open Explorer (`Ctrl+T`)
- Press `n` for new file
- Enter filename
- Press Enter to create

**Method 3: Save Unnamed Buffer**
```
:w filename.txt
```

### Renaming Files

**Using Explorer**:
- Open Explorer (`Ctrl+T`)
- Navigate to file
- Press `r`
- Enter new name
- Press Enter

**Using Save As**:
```
:w newname.txt
```
Creates new file, keeps old file unchanged.

### Deleting Files

**Using Explorer**:
- Open Explorer (`Ctrl+T`)
- Navigate to file
- Press `d`
- Confirm deletion

**Warning**: Deletion is permanent. There is no undo for file system operations.

## Search and Navigation

### Text Search

**Starting Search**:
```
/searchterm
```

**Features**:
- Case-insensitive by default
- Real-time match highlighting
- Pattern displayed in status bar
- All matches highlighted

**Navigation**:
- `n` - Next match
- `Shift+N` - Previous match
- `Esc` - Clear highlights

**Search Highlights**:
- Yellow (`#ffff00`) with transparency - All matches
- Orange (`#ffa500`) with stronger highlight - Current match

### Fuzzy File Finding

**Activation**:
```
Ctrl+F
```

**Algorithm**:
- Fuzzy matching (non-consecutive characters)
- Scores based on match quality
- Consecutive matches score higher
- Searches from current directory

**Usage**:
1. Press `Ctrl+F`
2. Type characters (order matters, but gaps allowed)
3. Use `↑`/`↓` to select
4. Press `Enter` to open

**Example**:
- Query: `mgo` 
- Matches: `main.go`, `models/user.go`, `migrations/001.go`
- Higher score: `main.go` (consecutive `mgo`)

### Cursor Navigation

**Line-Based**:
- `j`/`k` - Up/down one line
- `gg` - First line of file
- `G` - Last line of file
- `42G` - Jump to line 42

**Character-Based**:
- `h`/`l` - Left/right one character
- `0` - Start of line
- `$` - End of line

**Word-Based**:
- `w` - Next word start
- `b` - Previous word start  
- `e` - End of current word

**Viewport**:
- `Ctrl+E` - Scroll down
- `Ctrl+Y` - Scroll up
- Automatic scroll offset: 3 lines from top/bottom

## Terminal

### Overview

Vem includes a full-featured terminal emulator integrated into the buffer system.

### Activation

**Toggle Terminal**:
```
Ctrl+` (backtick)
```

**Behavior**:
- First press: Creates terminal buffer, starts shell
- Second press: Closes terminal buffer
- Terminal appears in current pane
- Replaces active buffer (can switch back with `:bp`)

### Features

**VT100/ANSI Support**:
- Full escape sequence parsing
- Cursor movement commands
- Clear screen, erase line
- Insert/delete line
- Character attributes

**Color Support**:
- 256-color mode
- True color (24-bit RGB)
- Standard ANSI colors (0-15)
- Extended colors (16-255)

**Text Attributes**:
- Bold
- Italic
- Underline
- Reverse video
- Dim
- Combinations supported

**PTY Integration**:
- Unix: Uses `/dev/ptmx` with `creack/pty`
- Windows: Uses ConPTY
- Full shell interaction
- Signal handling
- Window resize support

### Shell Integration

**Default Shell**:
- Unix/Linux: `$SHELL` environment variable (usually `/bin/bash` or `/bin/zsh`)
- macOS: `$SHELL` (usually `/bin/zsh`)
- Windows: `cmd.exe` or PowerShell

**Environment**:
- Inherits parent environment
- `TERM=xterm-256color`
- Working directory: Same as editor

**Lifecycle**:
1. Terminal created with `Ctrl+` `
2. Shell starts in PTY
3. User interacts with shell
4. Shell exits (e.g., `exit` command)
5. Terminal buffer auto-closes
6. Editor returns to previous buffer

### Terminal Buffer Behavior

**Buffer Properties**:
- Shows as `[Terminal]` in `:ls`
- Cannot be saved (`:w` ignored)
- No unsaved changes warning
- `:q` and `:bd` work normally
- Terminal buffer index tracked separately

**Integration**:
- Works with pane system
- Can split with terminal active
- Can have multiple terminals in different panes
- Switching buffers preserves terminal state

**Limitations**:
- Cannot enter INSERT mode in terminal buffer
- Cannot edit terminal output
- Search not available in terminal
- Undo not available

### Terminal Input

**In TERMINAL Mode**:
- All keys pass to shell (except `Esc` and `Ctrl+` `)
- Function keys supported
- Arrow keys for shell history
- Tab completion works
- Ctrl sequences pass through

**Special Keys**:
- `Esc` - Return to NORMAL mode (doesn't close terminal)
- `Ctrl+` ` - Close terminal buffer
- `Ctrl+C` - Sends SIGINT to shell process
- `Ctrl+D` - EOF (may exit shell)

### Window Resize

Terminal automatically resizes when:
- Editor window resized
- Pane layout changes
- Fullscreen toggled

**Terminal Dimensions**:
- Calculated based on pane size
- Minimum: 20 cols × 6 rows
- Updates sent to PTY
- Shell receives `SIGWINCH`

## Pane Management

### Pane Basics

**What is a Pane?**
- Subdivision of editor window
- Each pane shows one buffer
- One pane is active at a time
- Panes can be split horizontally or vertically

**Pane Indicators**:
- Active pane: Brighter background
- Inactive pane: Dimmed background (~15% darker)
- Status bar shows: `PANE 1/3` (current/total)

### Creating Panes

**Split Vertical** (top/bottom):
```
Ctrl+S v
```

**Split Horizontal** (left/right):
```
Ctrl+S h
```

**Behavior**:
- Current pane splits in two
- New pane shows same buffer
- Current pane remains active
- Both panes independently scrollable

### Navigating Panes

**Directional Navigation**:
- `Alt+h` - Focus left pane
- `Alt+j` - Focus down pane
- `Alt+k` - Focus up pane
- `Alt+l` - Focus right pane

**Cycle Through Panes**:
```
Shift+Tab
```

**Focus Specific Pane**:
- `Ctrl+H` - Focus explorer (if open)
- `Ctrl+L` - Focus editor

### Pane Layout

**Equalize Panes**:
```
Ctrl+S =
```
Makes all panes equal size.

**Zoom Pane**:
```
Ctrl+S o
```
- First press: Maximize active pane (hides others)
- Second press: Restore original layout
- Other panes remain in memory

### Closing Panes

**Close Active Pane**:
```
Ctrl+X
```
or
```
:q
```

**Behavior**:
- Closes active pane
- Closes buffer in that pane
- Prompts if unsaved changes
- If last pane: Switches to buffer 0 instead of exiting
- Focus moves to remaining pane

**Force Close**:
```
:q!
```
Discards unsaved changes.

### Pane and Buffer Relationship

**Key Concepts**:
- Pane = Visual window
- Buffer = File content
- One buffer can appear in multiple panes
- Editing in one pane affects all panes showing that buffer
- Cursor position independent per pane

**Example**:
```
1. Open file.txt in pane 1
2. Split vertical (Ctrl+S v) → pane 2 shows file.txt
3. Edit in pane 1 → changes appear in pane 2
4. Scroll in pane 1 → pane 2 scroll position unchanged
5. Move cursor in pane 2 → pane 1 cursor position unchanged
```

### Layout Types

**Horizontal Split**:
```
┌─────────────┐
│   Pane 1    │
├─────────────┤
│   Pane 2    │
└─────────────┘
```

**Vertical Split**:
```
┌──────┬──────┐
│      │      │
│ Pane │ Pane │
│  1   │  2   │
│      │      │
└──────┴──────┘
```

**Complex Layout**:
```
┌──────┬──────┐
│      │ P2   │
│ Pane ├──────┤
│  1   │ P3   │
│      │      │
└──────┴──────┘
```

## Configuration

### Current Status

Vem does not currently support runtime configuration files (e.g., `.vemrc`). All settings are hardcoded defaults.

### Defaults

**Editor**:
- Tab size: Uses system default
- Scroll offset: 3 lines
- Line numbers: Always shown
- Syntax highlighting: Enabled by default

**Colors**:
- Background: Dark blue (`#1a1f2e`)
- Text: Based on Chroma theme
- Selection: Translucent blue
- Search: Yellow highlights

**Fonts**:
- Primary: JetBrains Mono Nerd Font
- Fallback: Go default fonts
- Size: Based on window size

**Terminal**:
- Shell: `$SHELL` environment variable
- TERM: `xterm-256color`
- Columns: Calculated from pane width
- Rows: Calculated from pane height

### Future Configuration

Planned features for future releases:
- `.vemrc` configuration file
- Color scheme selection
- Custom keybindings
- Font size adjustment
- Tab width settings
- Line wrap options
- Syntax theme selection

## Status Bar Reference

### Status Bar Format

```
MODE <mode> | FILE <filename> | CURSOR <line>:<col> | PANE <current>/<total> | <message>
```

### Status Bar Components

**Mode Indicator**:
- `MODE NORMAL` - Normal mode
- `MODE INSERT` - Insert mode
- `MODE VISUAL` - Visual character mode
- `MODE VISUAL LINE` - Visual line mode
- `MODE DELETE` - Delete mode
- `MODE COMMAND` - Command mode (shows `:` prompt)
- `MODE EXPLORER` - Explorer mode
- `MODE SEARCH` - Search mode (shows `/` prompt)
- `FUZZY FINDER` - Fuzzy finder mode
- `MODE TERMINAL` - Terminal mode

**File Indicator**:
- `FILE filename.txt` - Normal file
- `FILE filename.txt +` - Modified file
- `FILE [No Name]` - Unnamed buffer
- `FILE [Terminal]` - Terminal buffer
- `FILE filename.txt [RO]` - Read-only buffer

**Cursor Position**:
- `CURSOR 42:15` - Line 42, Column 15
- 1-indexed (first line = 1, first column = 1)

**Pane Indicator**:
- `PANE 2/3` - Pane 2 of 3 total panes
- Only shown when multiple panes exist

**Message Area**:
- Shows command output
- Displays errors
- Confirmation messages
- Search results count

### Example Status Bars

```
MODE NORMAL | FILE main.go | CURSOR 42:15 | PANE 1/2 | Ready
```

```
MODE INSERT | FILE main.go + | CURSOR 42:15 | Editing...
```

```
MODE VISUAL | FILE README.md | CURSOR 10:5 | 3 lines selected
```

```
:w main.go
```

```
/search term
```

```
MODE TERMINAL | FILE [Terminal] | CURSOR 1:1 | Running shell
```

## Tips and Tricks

### Workflow Tips

1. **Quick File Opening**:
   - Use `Ctrl+F` fuzzy finder for fast navigation
   - Type partial name, no need for full path

2. **Multi-File Editing**:
   - Open multiple files from command line: `vem *.go`
   - Switch quickly with `:bn` / `:bp`

3. **Terminal Workflow**:
   - Keep terminal open in split pane
   - Run commands without leaving editor
   - Use `Alt+h/j/k/l` to switch between editor and terminal

4. **Search Efficiency**:
   - Use `/` for quick search
   - Navigate with `n` for next occurrence
   - Clear with `Esc` when done

5. **Buffer Management**:
   - Use `:ls` to see all open files
   - Close unused buffers with `:bd`
   - Force close all with `:qa!` if needed

### Keyboard Efficiency

1. **Stay in Normal Mode**:
   - Use navigation commands instead of arrow keys
   - `w`/`b` faster than repeated `h`/`l`

2. **Use Visual Mode**:
   - Select text with `v`
   - Copy with `c`, paste with `p`
   - Line-wise with `Shift+V` for full lines

3. **Pane Shortcuts**:
   - `Ctrl+S v` / `Ctrl+S h` for splits
   - `Alt+hjkl` for fast pane switching
   - `Ctrl+S o` to zoom current pane

4. **Command Mode**:
   - `:wq` faster than `:w` then `:q`
   - `:e!` to reload file (discard changes)

### Common Patterns

**Edit-Compile-Test Cycle**:
```
1. Edit file in pane 1
2. Ctrl+` to open terminal
3. Run: go build
4. Fix errors
5. Ctrl+` to close terminal
6. Repeat
```

**Side-by-Side Comparison**:
```
1. Open file1.go
2. Ctrl+S h (horizontal split)
3. :e file2.go in second pane
4. Alt+h / Alt+l to switch panes
```

**Project Navigation**:
```
1. Ctrl+T for explorer
2. Browse project structure
3. Enter to open file
4. Ctrl+L to focus editor
5. Repeat
```

## Troubleshooting

### Common Issues

**File Won't Save**:
- Check file permissions
- Ensure directory exists
- Use `:w!` to force (if appropriate)

**Can't Enter Insert Mode**:
- Check if buffer is read-only (`[RO]` in status)
- Help buffers are read-only by design

**Search Not Working**:
- Press `Esc` to clear previous search
- Check pattern in status bar
- Use `/` to start fresh search

**Terminal Not Opening**:
- Check if `$SHELL` is set (Unix/Linux/macOS)
- Verify terminal emulator support on platform

**Pane Navigation Confusing**:
- Use `:ls` to see all buffers
- Use `Shift+Tab` to cycle if unsure
- Close extra panes with `Ctrl+X`

### Platform-Specific Notes

**Linux**:
- Requires Vulkan drivers
- Wayland and X11 supported
- Shell typically `bash` or `zsh`

**macOS**:
- Uses Metal graphics
- Shell typically `zsh`
- Intel and Apple Silicon supported

**Windows**:
- Uses Direct3D 11
- Terminal uses ConPTY
- Shell is `cmd.exe` or PowerShell

## Version Information

This reference guide corresponds to Vem Phase 1.

**Features Included**:
- Modal editing (9 modes)
- Syntax highlighting (200+ languages)
- Integrated terminal
- Pane management
- Fuzzy file finding
- File explorer
- Multi-buffer support
- Command line file opening
- Built-in help system
- Read-only buffers

**Not Yet Implemented**:
- Configuration files
- Custom keybindings
- Plugin system
- Macros
- Registers (advanced)
- Split window resizing
- Tab pages

For the latest updates, see [GitHub Repository](https://github.com/javanhut/Vem).
