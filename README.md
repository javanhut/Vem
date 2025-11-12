# Vem - Vim Emulator

A modern, opinionated Vim emulator built from scratch in Go with GPU-accelerated rendering. Vem brings the power of modal editing with a clean, NeoVim-inspired interface and cross-platform support.

## Overview

Vem is a lightweight yet powerful text editor that combines Vim's modal editing paradigm with modern features like fuzzy file finding, pane splitting, and GPU-accelerated rendering. Built with Go and Gio UI, it runs natively on Linux, macOS, Windows, and WebAssembly without external dependencies.

## Features

### Core Editing

- **Modal Editing**: Full Vim-like modes (NORMAL, INSERT, VISUAL, DELETE, COMMAND, EXPLORER, SEARCH, FUZZY_FINDER, TERMINAL)
- **Vim Motions**: Complete navigation with hjkl, word motions (w/b/e), line jumps (0/$), document jumps (gg/G)
- **Visual Mode**: Line and character selection with copy/delete/paste operations
- **Undo System**: Full undo support for all edit operations
- **Multi-Buffer Support**: Open and edit multiple files simultaneously, including from command line
- **Search & Highlight**: Case-insensitive search with match highlighting and navigation
- **Syntax Highlighting**: Powered by Chroma with support for 200+ languages and multiple color themes
- **Integrated Terminal**: Full-featured terminal emulator with VT100/ANSI support and true color

### Window Management

- **Pane Splitting**: Split windows horizontally or vertically to view multiple files side-by-side
- **Pane Navigation**: Navigate between panes with Alt+hjkl or Shift+Tab
- **Zoom Mode**: Temporarily maximize any pane for focused editing
- **Equalize Layout**: Balance all pane sizes with a single command
- **Active Pane Dimming**: Clear visual indication of which pane is active

### File Management

- **File Explorer**: Built-in tree view with directory navigation
- **Fuzzy Finder**: Fast file search with fuzzy matching (Ctrl+F)
- **File Operations**: Create, rename, and delete files directly from the explorer
- **Directory Navigation**: Move up/down directory hierarchy easily
- **Auto-Scroll**: Explorer automatically scrolls to keep selection visible

### User Experience

- **Fullscreen Mode**: Distraction-free editing with Shift+Enter
- **Status Bar**: Shows mode, file name, cursor position, pane info, and messages
- **File Type Icons**: Visual file type indicators in the explorer using Nerd Font icons
- **Command-Line Interface**: Vim-style commands (:e, :w, :q, :wq, :bd, :help, etc.)
- **Built-in Help**: Comprehensive help system with :help command
- **Read-Only Buffers**: Support for read-only buffers (help pages, system info)
- **GPU-Accelerated**: Smooth, responsive interface using Gio UI
- **Cross-Platform**: Identical experience on Linux, macOS, Windows, and WebAssembly

## Installation

### Quick Install

**Linux and macOS:**
```bash
git clone https://github.com/javanhut/Vem.git
cd Vem
make install
```

The Makefile automatically detects your OS/architecture, checks for dependencies, and installs Vem to `/usr/local/bin`.

**Windows:**
```bash
git clone https://github.com/javanhut/Vem.git
cd Vem
make build
```

This creates `vem.exe` in the current directory. Add it to your PATH or run directly.

### Prerequisites

#### All Platforms
- Go 1.25.3 or later
- Git
- Make

#### Linux-Specific
- Vulkan headers (automatically installed by `make install`)
  - **Debian/Ubuntu**: `libvulkan-dev libxkbcommon-dev libwayland-dev`
  - **Fedora/RHEL**: `vulkan-devel libxkbcommon-devel wayland-devel`
  - **Arch/Manjaro**: `vulkan-headers vulkan-icd-loader libxkbcommon wayland`
  - **openSUSE**: `vulkan-devel libxkbcommon-devel wayland-devel`
  - **Alpine**: `vulkan-headers vulkan-loader-dev libxkbcommon-dev wayland-dev`

### Manual Build

```bash
git clone https://github.com/javanhut/Vem.git
cd Vem

# Set local build cache (recommended)
export GOCACHE="$(pwd)/.gocache"

# Build
go build -o vem

# Run
./vem

# Install to /usr/local/bin (optional)
sudo install -m 755 vem /usr/local/bin/vem
```

## Quick Start

```bash
# Launch Vem
vem

# Open a specific file
vem main.go

# Open multiple files
vem file1.txt file2.go file3.md

# Create a new file
vem newfile.txt
```

### Basic Workflow

1. **Open explorer**: Press `Ctrl+T`
2. **Navigate files**: Use `j`/`k` to move up/down
3. **Open file**: Press `Enter` on a file
4. **Edit text**: Press `i` to enter INSERT mode
5. **Save**: Press `Esc`, then `:w` and `Enter`
6. **Quit**: Type `:q` and press `Enter`

## Keybindings

### Global (Work in All Modes)

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+T` | Toggle Explorer | Show/hide file tree |
| `Ctrl+H` | Focus Explorer | Switch to file tree |
| `Ctrl+L` | Focus Editor | Switch to editor |
| `Ctrl+F` | Fuzzy Finder | Quick file search |
| `Ctrl+U` | Undo | Undo last operation |
| `Ctrl+X` | Close Pane | Close active pane |
| `Ctrl+` ` | Toggle Terminal | Open/close integrated terminal |
| `Shift+Enter` | Fullscreen | Toggle fullscreen mode |
| `Shift+Tab` | Cycle Panes | Move to next pane |

### Pane Management (Ctrl+S Prefix)

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
| `Alt+j` | Focus Down | Move to below pane |
| `Alt+k` | Focus Up | Move to above pane |
| `Alt+l` | Focus Right | Move to right pane |

### NORMAL Mode

#### Mode Switching

| Key | Action | Description |
|-----|--------|-------------|
| `i` | INSERT | Enter insert mode |
| `v` | VISUAL | Enter visual character mode |
| `Shift+V` | VISUAL LINE | Enter visual line mode |
| `d` | DELETE | Enter delete mode |
| `:` | COMMAND | Open command line |
| `/` | SEARCH | Start search |

#### Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `h`/`j`/`k`/`l` | Move | Left/Down/Up/Right |
| `w` | Word Forward | Next word start |
| `b` | Word Backward | Previous word start |
| `e` | Word End | End of current word |
| `0` | Line Start | Jump to line beginning |
| `$` | Line End | Jump to line end |
| `gg` | First Line | Jump to top of file |
| `G` | Last Line | Jump to bottom of file |
| `<n>G` | Goto Line | Jump to line n (e.g., `42G`) |

#### Scrolling

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+E` | Scroll Down | Scroll one line down |
| `Ctrl+Y` | Scroll Up | Scroll one line up |

#### Search

| Key | Action | Description |
|-----|--------|-------------|
| `/` | Start Search | Enter search mode |
| `n` | Next Match | Jump to next result |
| `Shift+N` | Previous Match | Jump to previous result |
| `Esc` | Clear Search | Clear search highlights |

### INSERT Mode

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Exit | Return to NORMAL mode |
| `Enter` | New Line | Insert newline |
| `Tab` | Insert Tab | Insert tab character |
| `Backspace` | Delete Back | Delete previous character |
| `Delete` | Delete Forward | Delete next character |
| Arrow keys | Navigate | Move cursor while typing |

### VISUAL Mode

#### Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `h`/`j`/`k`/`l` | Extend Selection | Move selection boundaries |
| `w`/`b`/`e` | Word Motion | Move by words |
| `0`/`$` | Line Bounds | Move to line start/end |
| `gg`/`G` | Document Bounds | Extend to top/bottom |

#### Operations

| Key | Action | Description |
|-----|--------|-------------|
| `c` | Copy | Copy selection to clipboard |
| `d` | Delete | Delete selected text |
| `p` | Paste | Paste from clipboard |
| `v` | Exit | Return to NORMAL mode |
| `Esc` | Exit | Return to NORMAL mode |

### COMMAND Mode

#### File Operations

| Command | Description |
|---------|-------------|
| `:e <file>` | Open file for editing |
| `:w` | Save current file |
| `:w <file>` | Save as new file |
| `:wq` | Save and close |
| `:q` | Close (fails if unsaved) |
| `:q!` | Force close (discard changes) |

#### Buffer Management

| Command | Description |
|---------|-------------|
| `:bn` or `:bnext` | Next buffer |
| `:bp` or `:bprev` | Previous buffer |
| `:bd` or `:bdelete` | Close buffer |
| `:bd!` | Force close buffer |
| `:ls` or `:buffers` | List all buffers |

#### File Explorer

| Command | Description |
|---------|-------------|
| `:ex` or `:explore` | Toggle file explorer |
| `:cd <path>` | Change directory |
| `:pwd` | Print working directory |

#### Help System

| Command | Description |
|---------|-------------|
| `:help` or `:h` | Open comprehensive keybindings help |

### EXPLORER Mode

| Key | Action | Description |
|-----|--------|-------------|
| `j`/`k` | Navigate | Move up/down in tree |
| `Enter` | Open/Toggle | Open file or toggle directory |
| `h` | Collapse | Collapse directory |
| `l` | Expand | Expand directory |
| `r` | Rename | Rename file/directory |
| `d` | Delete | Delete file/directory |
| `n` | New File | Create new file |
| `u` | Navigate Up | Go to parent directory |
| `q` | Quit | Return to editor |
| `Esc` | Exit | Return to NORMAL mode |

### SEARCH Mode

| Key | Action | Description |
|-----|--------|-------------|
| Type text | Build Pattern | Add to search pattern |
| `Enter` | Execute | Find first match |
| `Backspace` | Delete Char | Remove last character |
| `Esc` | Cancel | Exit search mode |

After search, use `n` and `Shift+N` in NORMAL mode to navigate matches.

### FUZZY_FINDER Mode

| Key | Action | Description |
|-----|--------|-------------|
| Type text | Filter Files | Show matching files |
| `↑`/`↓` | Navigate | Select different file |
| `Enter` | Open | Open selected file |
| `Backspace` | Delete Char | Remove last character |
| `Esc` | Cancel | Close fuzzy finder |

### TERMINAL Mode

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+` ` | Toggle Terminal | Open/close terminal |
| `Esc` | Exit to NORMAL | Return to NORMAL mode |
| All other keys | Pass to Shell | Direct terminal input |

**Features**:
- Full VT100/ANSI escape sequence support
- 256-color and true color (24-bit) support
- Bold, italic, underline, and other text attributes
- Auto-closes when shell exits
- Integrates with buffer system (switch with `:bn`/`:bp`)

## Documentation

- **[Reference Guide](docs/reference.md)** - Complete command and feature reference
- **[Keybindings Reference](docs/keybindings.md)** - Complete keybinding documentation
- **[Architecture Guide](docs/Architecture.md)** - Technical implementation details
- **[Tutorial](docs/tutorial.md)** - Step-by-step getting started guide
- **[Installation Guide](docs/installation.md)** - Platform-specific installation instructions
- **[Pane Splitting Guide](docs/pane-splitting.md)** - Detailed pane management guide
- **[Navigation Guide](docs/navigation.md)** - Pane navigation and fullscreen features
- **[Search Guide](docs/search.md)** - Search functionality documentation
- **[Syntax Highlighting](docs/syntax-highlighting.md)** - Color themes and language support

## Platform Support

### Linux
- **Display Servers**: X11 and Wayland
- **Graphics**: Vulkan
- **Tested On**: Ubuntu 22.04, Debian 12, Fedora 40, Arch Linux

### macOS
- **Graphics**: Metal (built-in)
- **Architecture**: Intel (x86_64) and Apple Silicon (arm64)
- **Tested On**: macOS 13 (Ventura) and later

### Windows
- **Graphics**: Direct3D 11 (built-in)
- **Tested On**: Windows 10, Windows 11

### WebAssembly
- **Graphics**: WebGL
- **Support**: Experimental via Gio's WASM backend

## Project Structure

```
Vem/
├── main.go                     # Application entry point
├── internal/
│   ├── appcore/               # Core application and rendering
│   │   ├── app.go            # Event handling and UI layout
│   │   ├── help.go           # Built-in help system
│   │   ├── keybindings.go    # Keybinding system
│   │   ├── pane_actions.go   # Pane management actions
│   │   ├── pane_rendering.go # Pane rendering logic
│   │   └── fuzzy.go          # Fuzzy finder implementation
│   ├── editor/                # Text editing engine
│   │   ├── buffer.go         # Buffer abstraction with read-only support
│   │   ├── buffer_manager.go # Multi-buffer management
│   │   └── buffer_test.go    # Unit tests
│   ├── filesystem/            # File system operations
│   │   ├── tree.go           # File tree data structure
│   │   ├── loader.go         # Directory loading
│   │   ├── finder.go         # File finding logic
│   │   └── icons.go          # File type icons
│   ├── panes/                 # Pane management
│   │   ├── manager.go        # Pane layout manager
│   │   ├── layout.go         # Layout calculations
│   │   ├── navigation.go     # Pane navigation
│   │   └── pane.go           # Pane abstraction
│   ├── syntax/                # Syntax highlighting
│   │   ├── highlighter.go    # Language detection and tokenization
│   │   └── theme.go          # Color theme management
│   ├── terminal/              # Integrated terminal
│   │   ├── terminal.go       # Terminal emulator core
│   │   ├── pty_unix.go       # Unix PTY implementation
│   │   ├── pty_windows.go    # Windows PTY implementation
│   │   ├── buffer.go         # Terminal buffer management
│   │   ├── colors.go         # ANSI color handling
│   │   └── input.go          # Terminal input processing
│   └── fonts/                 # Font management
│       ├── fonts.go          # Font loading and rendering
│       ├── JetBrainsMonoNerdFont-Regular.ttf
│       └── JetBrainsMonoNerdFont-Bold.ttf
├── docs/                      # Documentation
├── go.mod                     # Go module definition
├── Makefile                   # Build automation
└── README.md                  # This file
```

## Dependencies

Vem uses minimal dependencies:

### Core
- **[Gio UI](https://gioui.org)** v0.9.0 - GPU-accelerated UI framework
  - Vulkan (Linux)
  - Metal (macOS)
  - Direct3D (Windows)
  - WebGL (WebAssembly)
- **[Chroma](https://github.com/alecthomas/chroma)** v2.20.0 - Syntax highlighting engine
  - 200+ language lexers
  - Multiple color themes
  - Fast tokenization
- **[vt10x](https://github.com/hinshun/vt10x)** - VT100/ANSI terminal emulator
  - Full escape sequence support
  - 256-color and true color
- **[pty](https://github.com/creack/pty)** v1.1.21 - Cross-platform PTY support
- **[clipboard](https://golang.design/x/clipboard)** v0.7.1 - System clipboard integration

### Transitive (Automatic)
- `gioui.org/shader` v1.0.8 - Shader compilation
- `github.com/go-text/typesetting` v0.3.0 - Text layout
- `github.com/dlclark/regexp2` v1.11.5 - Advanced regex (for Chroma)
- `golang.org/x/exp/shiny` - Platform abstraction
- `golang.org/x/image` v0.28.0 - Image handling
- `golang.org/x/sys` v0.33.0 - System calls
- `golang.org/x/text` v0.26.0 - Text processing

All dependencies are managed via Go modules and installed automatically.

## Development

### Setup Development Environment

```bash
git clone https://github.com/javanhut/Vem.git
cd Vem

# Install dependencies (automatic)
go mod download

# Set local build cache
export GOCACHE="$(pwd)/.gocache"

# Run in development mode
go run .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v -run TestBufferInsert ./internal/editor

# Verbose output
go test -v ./...
```

### Code Style

- Follow standard Go formatting (`gofmt`)
- Run `go vet` to catch common issues
- Use descriptive variable names
- Document exported functions and types
- Keep functions focused and testable
- Write unit tests for new features

## Troubleshooting

### Linux: Vulkan Headers Not Found

If using `make install`, Vulkan headers are installed automatically. For manual installation:

```bash
# Debian/Ubuntu
sudo apt-get install libvulkan-dev libxkbcommon-dev libwayland-dev

# Fedora/RHEL
sudo dnf install vulkan-devel libxkbcommon-devel wayland-devel

# Arch/Manjaro
sudo pacman -S vulkan-headers vulkan-icd-loader libxkbcommon wayland

# openSUSE
sudo zypper install vulkan-devel libxkbcommon-devel wayland-devel

# Alpine
sudo apk add vulkan-headers vulkan-loader-dev libxkbcommon-dev wayland-dev
```

### Build Cache Permission Issues

```bash
# Use local build cache
export GOCACHE="$(pwd)/.gocache"
make clean
make build
```

### Platform-Specific Keybinding Issues

Some platforms may not report modifier keys correctly. Vem includes workarounds for these platform quirks. See `docs/debugging.md` for details.

## Contributing

Vem is in active development. Contributions are welcome once the architecture stabilizes.

1. Check existing issues and documentation
2. Open an issue to discuss significant changes
3. Follow the code style guidelines
4. Write tests for new features
5. Update documentation to reflect changes

## License

Vem is licensed under the GNU General Public License v2.0 (GPLv2).

See [LICENSE](LICENSE) for the full license text.

## Current Status

Vem is feature-complete for Phase 1 and includes:

- Full modal editing system (NORMAL, INSERT, VISUAL, DELETE, COMMAND, EXPLORER, SEARCH, FUZZY_FINDER, TERMINAL)
- Syntax highlighting with 200+ languages
- Pane splitting and management
- Fuzzy file finder
- File explorer with operations
- Search with highlighting
- Multi-buffer support with command line file opening
- Integrated terminal emulator with full VT100/ANSI support
- Built-in help system (:help command)
- Read-only buffer support
- Undo functionality
- Cross-platform support
- GPU-accelerated rendering

## Acknowledgments

Vem is inspired by:
- **Vim** - Modal editing philosophy
- **NeoVim** - Modern text editing paradigm
- **Gio UI** - Cross-platform GPU-accelerated framework

Special thanks to the Go and Gio communities for excellent tools and documentation.

## Contact

- **Repository**: https://github.com/javanhut/Vem
- **Issues**: https://github.com/javanhut/Vem/issues
