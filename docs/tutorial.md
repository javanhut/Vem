# Vem Tutorial

A step-by-step guide to getting started with Vem.

## Table of Contents

- [Installation](#installation)
- [First Launch](#first-launch)
- [Understanding Modes](#understanding-modes)
- [Basic Editing](#basic-editing)
- [File Navigation](#file-navigation)
- [Working with Multiple Files](#working-with-multiple-files)
- [Advanced Features](#advanced-features)
- [Next Steps](#next-steps)

## Installation

### Prerequisites

Before installing Vem, ensure you have:

1. **Go 1.25.3 or later**: Download from [golang.org](https://golang.org/dl/)
2. **Git**: For cloning the repository
3. **Platform-specific requirements**:
   - **Linux**: Vulkan development headers
     ```bash
     # Debian/Ubuntu
     sudo apt-get install libvulkan-dev
     
     # Fedora/RHEL
     sudo dnf install vulkan-devel
     
     # Arch
     sudo pacman -S vulkan-headers vulkan-icd-loader
     ```
   - **macOS**: No additional requirements
   - **Windows**: No additional requirements

### Building from Source

```bash
# Clone the repository
git clone https://github.com/javanhut/Vem.git
cd Vem

# Set local build cache (recommended)
export GOCACHE="$(pwd)/.gocache"

# Run the editor
go run .
```

Alternatively, build a binary:

```bash
go build -o vem
./vem
```

## First Launch

### Opening Files

You can launch Vem in several ways:

```bash
# Launch with sample text
vem

# Open a specific file
vem main.go

# Open multiple files (first file is active, switch with :bn/:bp)
vem file1.txt file2.go file3.md

# Create a new file
vem newfile.txt
```

When you first launch Vem without arguments, you'll see:

```
┌─────────────────────────────────────────────────┐
│ Vem                                             │
├─────────────────────────────────────────────────┤
│                                                 │
│   1  Goal: Build a NeoVim inspired text editor │
│   2                                             │
│   3  Constraints:                               │
│   4    - Written in Go with a modern GPU UI    │
│   5    - Ships fonts + dependencies...         │
│                                                 │
├─────────────────────────────────────────────────┤
│ MODE NORMAL | FILE [No Name] | CURSOR 1:1 |... │
└─────────────────────────────────────────────────┘
```

The interface consists of:
- **Header**: Application title
- **Editor Area**: Text content with line numbers
- **Status Bar**: Mode, file name, cursor position, status messages

The editor starts in **NORMAL mode** with sample text loaded (or your specified file).

## Understanding Modes

Vem, like Vim, has different modes for different tasks:

### NORMAL Mode

The default mode for navigation and issuing commands.

**Purpose**: Navigate the file, switch to other modes, execute commands

**Indicator**: `MODE NORMAL` in status bar

**Common Keys**:
- `h/j/k/l` - Move cursor left/down/up/right
- `i` - Enter INSERT mode
- `v` - Enter VISUAL mode
- `:` - Enter COMMAND mode

### INSERT Mode

For typing and editing text.

**Purpose**: Insert and modify text

**Indicator**: `MODE INSERT` in status bar

**How to enter**: Press `i` from NORMAL mode

**How to exit**: Press `Esc` to return to NORMAL mode

### VISUAL Mode

For selecting multiple lines.

**Purpose**: Select and manipulate ranges of text

**Indicator**: `MODE VISUAL` in status bar

**How to enter**: Press `v` from NORMAL mode

**How to exit**: Press `Esc` or `v` to return to NORMAL mode

### COMMAND Mode

For executing commands (like `:w` to save).

**Purpose**: Execute editor commands

**Indicator**: `:` prompt at bottom of screen

**How to enter**: Press `:` from NORMAL or VISUAL mode

**How to exit**: Press `Esc` to cancel, or `Enter` to execute

### EXPLORER Mode

For navigating the file tree.

**Purpose**: Browse and open files

**Indicator**: `MODE EXPLORER` in status bar

**How to enter**: Press `Ctrl+T` to show tree, then `Ctrl+H` to focus it

**How to exit**: Press `Esc` or `Ctrl+L` to return to NORMAL mode

## Basic Editing

### Opening a File

There are two ways to open a file:

**Method 1: File Explorer**

1. Press `Ctrl+T` to open the file tree
2. Use `j`/`k` to navigate to your file
3. Press `Enter` to open it

**Method 2: Command Line**

1. Press `:` to enter COMMAND mode
2. Type `e /path/to/file.txt`
3. Press `Enter`

Example:
```
:e example.txt
```

### Navigating Text

In NORMAL mode, use these keys to move around:

**Character Movement**:
- `h` - Move left
- `l` - Move right
- `j` - Move down
- `k` - Move up

Or use arrow keys: `←` `↓` `↑` `→`

**Line Movement**:
- `0` - Jump to start of line
- `$` - Jump to end of line

**Document Movement**:
- `gg` - Jump to first line
- `G` - Jump to last line
- `42G` - Jump to line 42 (replace 42 with any line number)

**Pro Tip**: You can use counts with movements. For example, `5j` moves down 5 lines.

### Inserting Text

1. Move to where you want to insert text (use `h/j/k/l`)
2. Press `i` to enter INSERT mode
3. Type your text
4. Press `Esc` when done to return to NORMAL mode

Example workflow:
```
1. Position cursor on line 5
2. Press i
3. Type "Hello, world!"
4. Press Esc
```

### Deleting Text

**Single Character**:
- In INSERT mode: Use `Backspace` or `Delete`

**Entire Line**:
1. Press `d` to enter DELETE mode
2. Type the line number (e.g., `5`)
3. Press `d` again to confirm
4. Or just press `d` immediately to delete current line

### Saving a File

**If file already has a name**:
1. Press `:` to enter COMMAND mode
2. Type `w`
3. Press `Enter`

**If file needs a name (Save As)**:
1. Press `:`
2. Type `w filename.txt`
3. Press `Enter`

Example:
```
:w myfile.txt
```

### Closing a File

**Close without saving**:
```
:q!
```

**Close with save**:
```
:wq
```

Or save first, then close:
```
:w
:q
```

## File Navigation

### Using the File Explorer

Vem includes a built-in file tree navigator.

**Opening the Explorer**:
```
Press Ctrl+T
```

The file tree appears on the left side showing the current directory.

**Navigating**:
- `j` or `↓` - Move down
- `k` or `↑` - Move up
- `Enter` - Open file or toggle directory
- `h` or `←` - Collapse directory
- `l` or `→` - Expand directory

**Additional Operations**:
- `r` - Refresh tree
- `u` - Navigate to parent directory
- `q` or `Esc` - Exit back to editor

**Switching Between Editor and Explorer**:
- `Ctrl+H` - Focus file explorer
- `Ctrl+L` - Focus editor
- `Ctrl+T` - Toggle explorer visibility

### Changing Directories

You can change the working directory using commands:

```
:cd /path/to/directory
:cd ~/projects
:cd ..
```

To see current directory:
```
:pwd
```

### Fullscreen Mode

For distraction-free editing:

```
Press Shift+Enter
```

This toggles fullscreen mode. Press `Shift+Enter` again to exit.

## Working with Multiple Files

### Opening Multiple Files

Open additional files using `:e`:

```
:e file1.txt
:e file2.txt
:e file3.txt
```

Each file opens in its own buffer.

### Switching Between Buffers

**Next Buffer**:
```
:bn
```
or
```
:bnext
```

**Previous Buffer**:
```
:bp
```
or
```
:bprev
```

### Listing Buffers

To see all open buffers:
```
:ls
```
or
```
:buffers
```

### Closing Buffers

**Close current buffer**:
```
:bd
```

**Force close (discard changes)**:
```
:bd!
```

## Advanced Features

### Visual Line Selection

Visual mode lets you select and manipulate multiple lines:

1. Press `v` to enter VISUAL mode
2. Use `j`/`k` to extend selection
3. Press `c` to copy selection
4. Press `d` to delete selection
5. Press `p` to paste copied lines

Example: Copying lines 5-10
```
1. Move to line 5 (5G)
2. Press v to enter VISUAL mode
3. Press 5j to select 5 lines down
4. Press c to copy
5. Move to destination
6. Press v to enter VISUAL mode
7. Press p to paste
```

### Using Counts

Many commands accept a count prefix:

**Movement**:
- `5j` - Move down 5 lines
- `10k` - Move up 10 lines
- `3h` - Move left 3 characters

**Goto**:
- `42G` - Jump to line 42
- `1G` - Jump to line 1
- `99999G` - Jump to last line (if less than 99999 lines)

### Goto Sequences

The `g` key starts a goto sequence:

- `gg` - Jump to first line
- `10gg` - Jump to line 10
- `gG` - Jump to last line (same as `G`)

## Practical Examples

### Example 1: Creating a New File

```
1. Launch Vem
2. Press :
3. Type: e mynotes.txt
4. Press Enter
5. Press i to enter INSERT mode
6. Type your notes
7. Press Esc when done
8. Type :w to save
9. Type :q to quit
```

### Example 2: Editing an Existing File

```
1. Launch Vem
2. Press Ctrl+T to open file tree
3. Navigate to your file with j/k
4. Press Enter to open
5. Navigate to the line you want to edit
6. Press i to enter INSERT mode
7. Make your changes
8. Press Esc when done
9. Type :w to save
```

### Example 3: Working with Multiple Files

```
1. Open first file: :e file1.txt
2. Make some edits, save with :w
3. Open second file: :e file2.txt
4. Make some edits, save with :w
5. Switch back to first file: :bp
6. Continue editing
7. Close current file: :bd
8. Close remaining files: :q
```

### Example 4: Copying Lines Between Files

```
1. Open source file: :e source.txt
2. Move to lines you want to copy
3. Press v to enter VISUAL mode
4. Select lines with j/k
5. Press c to copy
6. Open destination file: :e dest.txt
7. Move to where you want to paste
8. Press v to enter VISUAL mode
9. Press p to paste
10. Save both: :w then :bp then :w
```

## Common Tasks Cheat Sheet

| Task | Command |
|------|---------|
| Open file | `:e filename` |
| Save file | `:w` |
| Save as | `:w newname` |
| Close file | `:q` |
| Save and close | `:wq` |
| Force close (no save) | `:q!` |
| Next buffer | `:bn` |
| Previous buffer | `:bp` |
| List buffers | `:ls` |
| Toggle file tree | `Ctrl+T` |
| Toggle fullscreen | `Shift+Enter` |
| Enter INSERT mode | `i` |
| Enter VISUAL mode | `v` |
| Return to NORMAL | `Esc` |
| Jump to line | `<line>G` |
| Jump to first line | `gg` |
| Jump to last line | `G` |
| Delete line | `d` then `d` |

## Tips for Vim Users

If you're coming from Vim or NeoVim:

**What Works the Same**:
- Modal editing (NORMAL, INSERT, VISUAL)
- Basic motions (h/j/k/l, 0/$, gg/G)
- Commands (:w, :q, :e, etc.)
- Counts (5j, 10k, etc.)

**What's Different**:
- Keybindings for file tree navigation (Ctrl+T, Ctrl+H, Ctrl+L)
- Fullscreen toggle (Shift+Enter instead of window manager)
- DELETE mode for deleting specific lines
- Limited text objects (currently only line-based)

**What's Not Yet Implemented** (see ROADMAP.md):
- Macros (recording/playback)
- Registers (numbered/named)
- Advanced text objects (word, paragraph, etc.)
- Marks
- Folds
- Splits/windows

## Troubleshooting

### "Vulkan headers not found" (Linux)

Install Vulkan development headers:
```bash
# Debian/Ubuntu
sudo apt-get install libvulkan-dev

# Fedora
sudo dnf install vulkan-devel

# Arch
sudo pacman -S vulkan-headers
```

### Keybindings not working

Some platforms have quirks with modifier key reporting. Vem includes workarounds, but if you experience issues:

1. Check the console output for debug messages
2. See `DEBUG_FINDINGS.md` for known platform issues
3. Report the issue with your platform details

### File tree not showing

Make sure the working directory has readable files:
```
:pwd  # Check current directory
:cd /path/to/files  # Change directory
```

## Next Steps

Now that you know the basics:

1. **Practice basic editing**: Open files, make changes, save
2. **Explore the file tree**: Navigate your project structure
3. **Try multiple buffers**: Work with several files at once
4. **Read the keybindings reference**: Learn all available commands
5. **Check the roadmap**: See what features are coming next

### Further Reading

- [Keybindings Reference](keybindings.md) - Complete list of all keybindings
- [Navigation Guide](navigation.md) - Detailed pane navigation documentation
- [Architecture](Architecture.md) - Technical implementation details
- [ROADMAP.md](../ROADMAP.md) - Development phases and milestones

### Getting Help

- Check the documentation in `docs/`
- Read `DEBUG_FINDINGS.md` for known issues and workarounds
- Open an issue on GitHub if you find bugs

## Welcome to Vem

You're now ready to start using Vem productively. Remember:

- Start in NORMAL mode
- Press `i` to insert text
- Press `Esc` to return to NORMAL
- Use `:w` to save and `:q` to quit
- Use `Ctrl+T` for the file tree

Happy editing!
