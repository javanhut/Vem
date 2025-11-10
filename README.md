# ProjectVem

A NeoVim-inspired text editor written in Go with cross-platform support and zero external dependencies.

## Overview

ProjectVem is a modern text editor that brings the power and efficiency of NeoVim's modal editing to a cross-platform desktop application. Built with Go and Gio UI, it runs natively on Linux, macOS, Windows, and WebAssembly without requiring users to install fonts, packages, or system dependencies.

## Features

### Current Implementation (Phase 1)

- **Modal Editing**: Vim-like modes (NORMAL, INSERT, VISUAL, DELETE, EXPLORER, COMMAND)
- **File Management**: Built-in file tree explorer with directory navigation
- **Multi-Buffer Support**: Open and edit multiple files simultaneously
- **Vim Commands**: Standard command-line interface (`:e`, `:w`, `:q`, `:wq`, etc.)
- **Pane Navigation**: Keyboard-driven navigation between file tree and editor
- **Fullscreen Mode**: Distraction-free editing with `Shift+Enter`
- **GPU-Accelerated Rendering**: Smooth, responsive interface using Gio UI
- **Cross-Platform**: Runs identically on Linux, macOS, Windows, and WebAssembly

### Planned Features (See ROADMAP.md)

- Vim macro recording and playback
- LSP integration for language intelligence
- Syntax highlighting via Treesitter-style pipeline
- Plugin system with Lua/Python/Carrion scripting support
- Command palette and fuzzy file finder
- Session persistence and workspace management
- Customizable themes and keybindings

## Installation

### Quick Install (Recommended)

**Linux and macOS:**
```bash
git clone https://github.com/javanhut/ProjectVem.git
cd ProjectVem
make install
```

The Makefile automatically detects your OS/architecture, checks for dependencies (including Vulkan headers on Linux), and installs Vem to `/usr/local/bin`.

**Windows:**
```bash
git clone https://github.com/javanhut/ProjectVem.git
cd ProjectVem
make build
```

This creates `vem.exe` in the current directory. Add it to your PATH or run directly.

For detailed installation instructions, troubleshooting, and manual build options, see [docs/installation.md](docs/installation.md).

### Prerequisites

#### All Platforms
- Go 1.25.3 or later
- Git (for cloning the repository)
- Make (for automated installation)

#### Linux-Specific
- Vulkan headers (automatically installed by `make install`)
  - **Debian/Ubuntu**: `libvulkan-dev libxkbcommon-dev libwayland-dev`
  - **Fedora/RHEL/CentOS**: `vulkan-devel libxkbcommon-devel wayland-devel`
  - **Arch/Manjaro**: `vulkan-headers vulkan-icd-loader libxkbcommon wayland`
  - **openSUSE**: `vulkan-devel libxkbcommon-devel wayland-devel`
  - **Alpine**: `vulkan-headers vulkan-loader-dev libxkbcommon-dev wayland-dev`

These libraries provide GPU-accelerated rendering (Vulkan), keyboard input handling (xkbcommon), and display server support (Wayland/X11) on Linux. The Makefile detects your package manager and installs them automatically. On macOS, Metal is used; on Windows, Direct3D is used (no extra dependencies needed).

### Manual Build

If you prefer not to use Make:

```bash
# Clone the repository
git clone https://github.com/javanhut/ProjectVem.git
cd ProjectVem

# Install Vulkan headers (Linux only)
# Debian/Ubuntu: sudo apt-get install libvulkan-dev libxkbcommon-dev libwayland-dev
# Fedora/RHEL: sudo dnf install vulkan-devel libxkbcommon-devel wayland-devel
# Arch: sudo pacman -S vulkan-headers vulkan-icd-loader libxkbcommon wayland
# openSUSE: sudo zypper install vulkan-devel libxkbcommon-devel wayland-devel
# Alpine: sudo apk add vulkan-headers vulkan-loader-dev libxkbcommon-dev wayland-dev

# Set local build cache (recommended to avoid permission issues)
export GOCACHE="$(pwd)/.gocache"

# Build the binary
go build -o vem

# Run the editor
./vem

# Or install manually to /usr/local/bin (Linux/macOS)
sudo install -m 755 vem /usr/local/bin/vem
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/editor

# Run with verbose output
go test -v ./...
```

## Quick Start

### Basic Navigation

1. Launch the editor: `go run .`
2. Open the file explorer: Press `Ctrl+T`
3. Navigate files: Use `j`/`k` or arrow keys
4. Open a file: Press `Enter` on a file
5. Edit text: Press `i` to enter INSERT mode
6. Save file: Press `Esc`, then type `:w` and press `Enter`
7. Quit: Type `:q` and press `Enter`

### Essential Keybindings

#### Window Management
- `Ctrl+T` - Toggle file explorer visibility
- `Ctrl+H` - Focus file explorer (when visible)
- `Ctrl+L` - Focus editor
- `Shift+Enter` - Toggle fullscreen mode

#### Mode Switching
- `i` - Enter INSERT mode
- `v` - Enter VISUAL line mode
- `d` - Enter DELETE mode
- `:` - Enter COMMAND mode
- `Esc` - Return to NORMAL mode

#### Text Navigation (NORMAL mode)
- `h/j/k/l` - Move left/down/up/right
- `0` - Jump to start of line
- `$` - Jump to end of line
- `gg` - Jump to first line
- `G` - Jump to last line
- `<count>G` - Jump to specific line (e.g., `42G`)

#### Editing (INSERT mode)
- Type normally to insert text
- `Enter` - Insert newline
- `Backspace` - Delete character before cursor
- `Delete` - Delete character after cursor
- `Esc` - Return to NORMAL mode

For complete documentation, see `docs/keybindings.md`.

## Documentation

- **[Installation Guide](docs/installation.md)** - Detailed installation instructions for all platforms
- **[Keybindings Reference](docs/keybindings.md)** - Complete list of all keybindings
- **[Architecture Guide](docs/Architecture.md)** - Technical architecture and design decisions
- **[Tutorial](docs/tutorial.md)** - Step-by-step guide for new users
- **[Navigation Guide](docs/navigation.md)** - Pane navigation and fullscreen features
- **[Project Description](PROJECT_DESCRIPTION.md)** - Project goals and vision
- **[Roadmap](ROADMAP.md)** - Development phases and milestones

## Project Structure

```
ProjectVem/
├── main.go                     # Application entry point
├── internal/
│   ├── appcore/               # Main application loop and rendering
│   │   ├── app.go            # Core event handling and UI layout
│   │   └── keybindings.go    # Keybinding system and actions
│   ├── editor/                # Text editing logic
│   │   ├── buffer.go         # Buffer abstraction (lines, cursor)
│   │   ├── buffer_test.go    # Buffer unit tests
│   │   └── buffer_manager.go # Multi-buffer management
│   └── filesystem/            # File tree navigation
│       ├── tree.go           # File tree data structure
│       └── loader.go         # Directory loading and caching
├── docs/                      # Documentation
│   ├── keybindings.md        # Keybinding reference
│   ├── Architecture.md       # Architecture documentation
│   ├── tutorial.md           # User tutorial
│   └── navigation.md         # Navigation features guide
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── PROJECT_DESCRIPTION.md     # Project vision
├── ROADMAP.md                 # Development roadmap
└── README.md                  # This file
```

## Dependencies

ProjectVem uses minimal, carefully selected dependencies:

### Core Dependencies
- **[Gio UI](https://gioui.org)** v0.9.0 - GPU-accelerated cross-platform UI framework
  - Vulkan backend on Linux
  - Metal backend on macOS
  - Direct3D backend on Windows
  - WebGL backend for WebAssembly

### Transitive Dependencies (Automatic)
- `gioui.org/shader` v1.0.8 - Shader compilation
- `github.com/go-text/typesetting` v0.3.0 - Text layout and shaping
- `golang.org/x/exp/shiny` - Gio platform abstraction
- `golang.org/x/image` v0.26.0 - Image handling
- `golang.org/x/sys` v0.33.0 - System calls
- `golang.org/x/text` v0.24.0 - Text processing

All dependencies are managed via Go modules and installed automatically with `go build` or `go run`.

## Platform Support

### Linux
- **Tested on**: Ubuntu 22.04, Debian 12, Fedora 40, Arch Linux
- **Display servers**: X11 and Wayland
- **Graphics**: Vulkan (requires `libvulkan-dev` or equivalent)

### macOS
- **Tested on**: macOS 13 (Ventura) and later
- **Graphics**: Metal (built-in, no extra dependencies)
- **Architecture**: Both Intel (x86_64) and Apple Silicon (arm64)

### Windows
- **Tested on**: Windows 10, Windows 11
- **Graphics**: Direct3D 11 (built-in, no extra dependencies)

### WebAssembly
- **Supported**: Experimental support via Gio's WASM backend
- **Graphics**: WebGL

## Development

### Setting Up Development Environment

```bash
# Clone repository
git clone https://github.com/javanhut/ProjectVem.git
cd ProjectVem

# Install dependencies (automatic with Go modules)
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
```

### Code Style

- Follow standard Go formatting: `gofmt` and `go vet`
- Use descriptive variable names (no single-letter names except idioms)
- Document exported functions and types
- Keep functions focused and testable
- See `CLAUDE.md` for detailed coding guidelines

### Contributing

ProjectVem is currently in active development (Phase 1). Contributions are welcome once the architecture stabilizes after Milestone 1 completion.

1. Check the [ROADMAP.md](ROADMAP.md) for current phase and priorities
2. Open an issue to discuss significant changes
3. Follow the code style guidelines in `CLAUDE.md`
4. Write tests for new features
5. Update documentation to reflect changes

## Troubleshooting

### Linux: Vulkan Headers Not Found

If using `make install`, Vulkan headers are installed automatically. If installing manually or if automatic installation fails:

```bash
# Debian/Ubuntu
sudo apt-get install libvulkan-dev libxkbcommon-dev libwayland-dev

# Fedora/RHEL/CentOS
sudo dnf install vulkan-devel libxkbcommon-devel wayland-devel

# Arch/Manjaro
sudo pacman -S vulkan-headers vulkan-icd-loader libxkbcommon wayland

# openSUSE
sudo zypper install vulkan-devel libxkbcommon-devel wayland-devel

# Alpine Linux
sudo apk add vulkan-headers vulkan-loader-dev libxkbcommon-dev wayland-dev
```

### Build Cache Permission Issues

```bash
# Use local build cache
export GOCACHE="$(pwd)/.gocache"
make clean
make build
```

For more troubleshooting help, see [docs/installation.md](docs/installation.md#troubleshooting).

### Platform-Specific Keybinding Issues

Some platforms may not report modifier keys correctly. ProjectVem includes workarounds for these platform quirks. See `DEBUG_FINDINGS.md` for details on modifier key tracking.

## License

ProjectVem is licensed under the GNU General Public License v2.0 (GPLv2).

See [LICENSE](LICENSE) for the full license text.

## Current Status

ProjectVem is currently in **Phase 1: Foundations & Architecture** (Weeks 1-4 of the roadmap).

### Completed
- Modal editing system (NORMAL, INSERT, VISUAL, DELETE, EXPLORER, COMMAND)
- File tree explorer with navigation
- Multi-buffer support
- Basic Vim commands (`:e`, `:w`, `:q`, `:wq`, `:bd`, etc.)
- Pane navigation (Ctrl+H, Ctrl+L)
- Fullscreen mode (Shift+Enter)
- GPU-accelerated rendering with Gio UI
- Cross-platform build system

### In Progress
- Architecture documentation finalization
- Buffer representation improvements
- Window split prototyping

### Next Milestones
- **M2 (Weeks 5-10)**: Core editing experience - Vim-parity motions, macros, multi-buffer improvements
- **M3 (Weeks 11-18)**: Language intelligence - LSP, syntax highlighting, plugin system
- **M4 (Weeks 19-24)**: User fluency - Command palette, fuzzy finder, themes
- **M5 (Weeks 25-28)**: Packaging and community launch

See [ROADMAP.md](ROADMAP.md) for detailed milestone planning.

## Contact

- **Repository**: https://github.com/javanhut/ProjectVem
- **Issues**: https://github.com/javanhut/ProjectVem/issues

## Acknowledgments

ProjectVem is inspired by:
- **NeoVim** - Modal editing paradigm and command interface
- **Vim** - Classic text editing motions and philosophy
- **Gio UI** - Cross-platform GPU-accelerated UI framework

Special thanks to the Gio community for their excellent documentation and support.
