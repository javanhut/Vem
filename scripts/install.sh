#!/usr/bin/env bash
# Install binary, config, desktop entry (Linux)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

GOOS=$(go env GOOS)

# Determine binary name
BINARY_NAME="vem"
if [ "$GOOS" = "windows" ]; then
    BINARY_NAME="vem.exe"
fi

# Install locations
INSTALL_DIR="/usr/local/bin"
INSTALL_PATH="$INSTALL_DIR/vem"
DESKTOP_DIR="/usr/local/share/applications"
DESKTOP_FILE="vem.desktop"
DESKTOP_PATH="$DESKTOP_DIR/$DESKTOP_FILE"
DESKTOP_SOURCE="assets/$DESKTOP_FILE"
ICON_DIR="/usr/local/share/pixmaps"
ICON_FILE="vem.png"
ICON_PATH="$ICON_DIR/$ICON_FILE"
ICON_SOURCE="assets/Vem.png"

if [ "$GOOS" = "windows" ]; then
    echo "Windows detected: Installation to system PATH not automated."
    echo "Binary built as $BINARY_NAME"
    echo "To use Vem, either:"
    echo "  1. Add the current directory to your PATH"
    echo "  2. Move $BINARY_NAME to a directory in your PATH"
    exit 0
fi

# Check if binary exists
if [ ! -f "$BINARY_NAME" ]; then
    echo "Binary '$BINARY_NAME' not found. Run 'imlazy build' first."
    exit 1
fi

echo "Installing $BINARY_NAME to $INSTALL_PATH..."

if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating $INSTALL_DIR..."
    sudo mkdir -p "$INSTALL_DIR"
fi

sudo install -m 755 "$BINARY_NAME" "$INSTALL_PATH"

echo "Installing default config..."
mkdir -p "$HOME/.config/vem"

if [ ! -f "$HOME/.config/vem/keybindings.toml" ]; then
    if [ -f "configs/keybindings.toml" ]; then
        cp configs/keybindings.toml "$HOME/.config/vem/keybindings.toml"
        echo "Installed default keybindings to $HOME/.config/vem/keybindings.toml"
    else
        echo "Warning: configs/keybindings.toml not found, skipping..."
    fi
else
    echo "Keybindings config already exists, skipping..."
fi

# Linux-specific: desktop entry and icon
if [ "$GOOS" = "linux" ]; then
    echo "Installing desktop entry and icon..."

    if [ ! -d "$DESKTOP_DIR" ]; then
        echo "Creating $DESKTOP_DIR..."
        sudo mkdir -p "$DESKTOP_DIR"
    fi

    if [ ! -d "$ICON_DIR" ]; then
        echo "Creating $ICON_DIR..."
        sudo mkdir -p "$ICON_DIR"
    fi

    if [ -f "$DESKTOP_SOURCE" ]; then
        sudo install -m 644 "$DESKTOP_SOURCE" "$DESKTOP_PATH"
        echo "Desktop entry installed to $DESKTOP_PATH"
    else
        echo "Warning: $DESKTOP_SOURCE not found, skipping..."
    fi

    if [ -f "$ICON_SOURCE" ]; then
        sudo install -m 644 "$ICON_SOURCE" "$ICON_PATH"
        echo "Icon installed to $ICON_PATH"
    else
        echo "Warning: $ICON_SOURCE not found, skipping..."
    fi
fi

echo "Cleaning up build artifact from source directory..."
rm -f "$BINARY_NAME"

echo "✓ Installation complete. Run 'vem' to start the editor."
