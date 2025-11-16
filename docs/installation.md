# Installation Guide

This guide covers all methods for installing Vem on Linux, macOS, and Windows.

## Table of Contents

- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
  - [Using Make (Recommended)](#using-make-recommended)
  - [Manual Build](#manual-build)
- [Platform-Specific Notes](#platform-specific-notes)
- [Troubleshooting](#troubleshooting)
- [Uninstallation](#uninstallation)

## Quick Start

The fastest way to install Vem on Linux or macOS:

```bash
git clone https://github.com/javanhut/Vem.git
cd Vem
make install
```

For Windows, see [Windows Installation](#windows).

## Prerequisites

### All Platforms

- **Go 1.25.3 or later**: [Download Go](https://golang.org/dl/)
- **Git**: For cloning the repository

### Linux-Specific

Vem requires several system libraries for GUI rendering and input handling on Linux. The Makefile will automatically detect and install these for you, but you can also install them manually:

#### Debian/Ubuntu
```bash
sudo apt-get install libvulkan-dev libxkbcommon-dev libwayland-dev
```

#### Fedora/RHEL/CentOS
```bash
sudo dnf install vulkan-devel libxkbcommon-devel wayland-devel
```

#### Arch Linux/Manjaro
```bash
sudo pacman -S vulkan-headers vulkan-icd-loader libxkbcommon wayland
```

#### openSUSE
```bash
sudo zypper install vulkan-devel libxkbcommon-devel wayland-devel
```

#### Alpine Linux
```bash
sudo apk add vulkan-headers vulkan-loader-dev libxkbcommon-dev wayland-dev
```

**What these libraries do:**
- **Vulkan**: GPU-accelerated rendering backend
- **xkbcommon**: Keyboard input handling
- **Wayland**: Display server support (also works on X11)

### macOS

No additional dependencies required. Vem uses Metal for GPU acceleration, which is built into macOS.

### Windows

No additional dependencies required. Vem uses Direct3D 11 for GPU acceleration, which is built into Windows.

## Installation Methods

### Using Make (Recommended)

The Makefile automates the build and installation process, including dependency checking.

#### 1. Clone the Repository

```bash
git clone https://github.com/javanhut/Vem.git
cd Vem
```

#### 2. View Available Commands

```bash
make help
```

This displays all available Makefile targets:

```
Available targets:
  help            Show this help message
  check-vulkan    Check if Vulkan headers are installed (Linux only)
  install-vulkan  Install Vulkan headers (Linux only)
  check-deps      Check and install all dependencies
  build           Build Vem binary for current OS/architecture
  install         Install Vem to system
  uninstall       Uninstall Vem from system
  clean           Remove built binaries
  test            Run all tests
```

#### 3. Install Vem

**Linux and macOS:**
```bash
make install
```

This will:
1. Automatically detect your OS and architecture
2. Check for Vulkan headers (Linux only) and install if missing
3. Build the Vem binary
4. Install to `/usr/local/bin/vem` (requires sudo)

After installation, run Vem from anywhere:
```bash
vem
```

**Windows:**
```bash
make build
```

On Windows, the Makefile builds `vem.exe` in the current directory. You can then:
- Add the Vem directory to your PATH, or
- Move `vem.exe` to a directory already in your PATH

### Manual Build

If you prefer not to use Make, you can build manually:

#### 1. Install Dependencies (Linux Only)

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

#### 2. Build the Binary

```bash
git clone https://github.com/javanhut/Vem.git
cd Vem

# Optional: Use local build cache to avoid permission issues
export GOCACHE="$(pwd)/.gocache"

# Build
go build -o vem
```

On Windows, this creates `vem.exe` automatically.

#### 3. Install (Optional)

**Linux/macOS:**
```bash
sudo install -m 755 vem /usr/local/bin/vem
```

**Windows:**
Move `vem.exe` to a directory in your PATH or add the current directory to your PATH.

## Platform-Specific Notes

### Linux

- **Installation Location**: `/usr/local/bin/vem`
- **Permissions**: Installation requires `sudo` for system-wide access
- **Graphics Backend**: Vulkan
- **Display Servers**: Works on both X11 and Wayland
- **Tested Distributions**: Ubuntu 22.04, Debian 12, Fedora 40, Arch Linux

### macOS

- **Installation Location**: `/usr/local/bin/vem`
- **Permissions**: Installation requires `sudo` for system-wide access
- **Graphics Backend**: Metal (built-in)
- **Architectures**: Both Intel (x86_64) and Apple Silicon (arm64)
- **Tested Versions**: macOS 13 (Ventura) and later

### Windows

- **Installation**: Not automated. Binary is built as `vem.exe`
- **Usage**: Run from current directory or add to PATH
- **Graphics Backend**: Direct3D 11 (built-in)
- **Tested Versions**: Windows 10, Windows 11

To add Vem to your PATH on Windows:
1. Build with `make build`
2. Move `vem.exe` to `C:\Program Files\Vem\` (or any preferred location)
3. Add that directory to your System PATH environment variable

## Troubleshooting

### Linux: Vulkan Headers Not Found

If the automatic installation fails:

```bash
# Try installing manually based on your distribution
# Debian/Ubuntu
sudo apt-get update && sudo apt-get install libvulkan-dev libxkbcommon-dev libwayland-dev

# Fedora/RHEL/CentOS
sudo dnf install vulkan-devel libxkbcommon-devel wayland-devel

# Arch/Manjaro
sudo pacman -S vulkan-headers vulkan-icd-loader libxkbcommon wayland

# openSUSE
sudo zypper install vulkan-devel libxkbcommon-devel wayland-devel

# Alpine Linux
sudo apk add vulkan-headers vulkan-loader-dev libxkbcommon-dev wayland-dev
```

Then retry:
```bash
make build
```

### Build Cache Permission Issues

If you encounter permission errors with the Go build cache:

```bash
export GOCACHE="$(pwd)/.gocache"
make clean
make build
```

### "Command Not Found" After Installation

**Linux/macOS:**

Ensure `/usr/local/bin` is in your PATH:
```bash
echo $PATH | grep /usr/local/bin
```

If not present, add to your shell configuration (`~/.bashrc`, `~/.zshrc`, etc.):
```bash
export PATH="/usr/local/bin:$PATH"
```

Then reload:
```bash
source ~/.bashrc  # or ~/.zshrc
```

**Windows:**

If `vem.exe` is not found:
1. Verify the directory containing `vem.exe` is in your PATH
2. Open a new terminal window to refresh environment variables
3. Try running with full path: `C:\path\to\vem.exe`

### Package Manager Not Detected (Linux)

If you're using a distribution with a different package manager, install the required libraries manually. Look for packages named:

- **Vulkan development files**: `libvulkan-dev`, `vulkan-devel`, or `vulkan-headers`
- **xkbcommon development files**: `libxkbcommon-dev`, `libxkbcommon-devel`, or `libxkbcommon`
- **Wayland client development files**: `libwayland-dev`, `wayland-devel`, or `wayland`

Then build:
```bash
make build
```

### Go Version Too Old

Check your Go version:
```bash
go version
```

If older than 1.25.3, update Go from [golang.org/dl/](https://golang.org/dl/)

## Uninstallation

### Linux and macOS

```bash
cd Vem
make uninstall
```

Or manually:
```bash
sudo rm /usr/local/bin/vem
```

### Windows

Simply delete `vem.exe` from wherever you placed it. If you added a directory to your PATH, remove that PATH entry from your environment variables.

## Running Tests

To verify your installation is working correctly:

```bash
cd Vem
make test
```

Or:
```bash
go test ./...
```

## Next Steps

After installation:

1. Read the [Tutorial](tutorial.md) for a guided introduction
2. Check the [Keybindings Reference](keybindings.md) for all commands
3. Explore the [Architecture Guide](Architecture.md) if you want to contribute

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/javanhut/Vem/issues)
- **Documentation**: Check the `docs/` directory
- **Repository**: [github.com/javanhut/Vem](https://github.com/javanhut/Vem)
