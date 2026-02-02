#!/usr/bin/env bash
# Remove installed files

set -e

GOOS=$(go env GOOS)

INSTALL_PATH="/usr/local/bin/vem"
DESKTOP_PATH="/usr/local/share/applications/vem.desktop"
ICON_PATH="/usr/local/share/pixmaps/vem.png"

if [ "$GOOS" = "windows" ]; then
    echo "Windows detected: No system installation to remove."
    echo "If you manually installed vem.exe, please remove it manually."
    exit 0
fi

REMOVED=0

if [ -f "$INSTALL_PATH" ]; then
    echo "Removing $INSTALL_PATH..."
    sudo rm -f "$INSTALL_PATH"
    REMOVED=1
fi

if [ "$GOOS" = "linux" ]; then
    if [ -f "$DESKTOP_PATH" ]; then
        echo "Removing $DESKTOP_PATH..."
        sudo rm -f "$DESKTOP_PATH"
        REMOVED=1
    fi

    if [ -f "$ICON_PATH" ]; then
        echo "Removing $ICON_PATH..."
        sudo rm -f "$ICON_PATH"
        REMOVED=1
    fi
fi

if [ $REMOVED -eq 1 ]; then
    echo "✓ Uninstall complete."
else
    echo "Vem is not installed at $INSTALL_PATH"
fi
