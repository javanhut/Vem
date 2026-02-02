#!/usr/bin/env bash
# Auto-install Linux dependencies via apt/dnf/pacman/zypper/apk

set -e

GOOS=$(go env GOOS 2>/dev/null || echo "")

if [ "$GOOS" != "linux" ]; then
    echo "Linux dependency installation not required on $GOOS"
    exit 0
fi

echo "Detecting package manager and installing dependencies..."

if command -v apt-get >/dev/null 2>&1; then
    echo "Using apt-get (Debian/Ubuntu)..."
    sudo apt-get update && sudo apt-get install -y \
        libvulkan-dev \
        libxkbcommon-dev \
        libxkbcommon-x11-dev \
        libwayland-dev \
        libx11-dev \
        libx11-xcb-dev \
        libegl1-mesa-dev \
        libxcursor-dev \
        libxfixes-dev \
        wayland-protocols

elif command -v dnf >/dev/null 2>&1; then
    echo "Using dnf (Fedora/RHEL/CentOS)..."
    sudo dnf install -y \
        vulkan-devel \
        libxkbcommon-devel \
        libxkbcommon-x11-devel \
        wayland-devel \
        libX11-devel \
        libX11-xcb \
        mesa-libEGL-devel \
        libXcursor-devel \
        libXfixes-devel \
        wayland-protocols-devel

elif command -v pacman >/dev/null 2>&1; then
    echo "Using pacman (Arch/Manjaro)..."
    sudo pacman -S --noconfirm \
        vulkan-headers \
        vulkan-icd-loader \
        libxkbcommon \
        libxkbcommon-x11 \
        wayland \
        wayland-protocols \
        libx11 \
        mesa \
        libxcursor \
        libxfixes

elif command -v zypper >/dev/null 2>&1; then
    echo "Using zypper (openSUSE)..."
    sudo zypper install -y \
        vulkan-devel \
        libxkbcommon-devel \
        libxkbcommon-x11-devel \
        wayland-devel \
        libX11-devel \
        libX11-xcb1 \
        Mesa-libEGL-devel \
        libXcursor-devel \
        libXfixes-devel \
        wayland-protocols-devel

elif command -v apk >/dev/null 2>&1; then
    echo "Using apk (Alpine Linux)..."
    sudo apk add \
        vulkan-headers \
        vulkan-loader-dev \
        libxkbcommon-dev \
        libxkbcommon-x11 \
        wayland-dev \
        wayland-protocols \
        libx11-dev \
        mesa-dev \
        libxcursor-dev \
        libxfixes-dev

else
    echo "Error: No supported package manager found."
    echo "Please install the following dependencies manually:"
    echo ""
    echo "Debian/Ubuntu:"
    echo "  sudo apt-get install libvulkan-dev libxkbcommon-dev libxkbcommon-x11-dev libwayland-dev libx11-dev libx11-xcb-dev libegl1-mesa-dev libxcursor-dev libxfixes-dev wayland-protocols"
    echo ""
    echo "Fedora/RHEL/CentOS:"
    echo "  sudo dnf install vulkan-devel libxkbcommon-devel libxkbcommon-x11-devel wayland-devel libX11-devel libX11-xcb mesa-libEGL-devel libXcursor-devel libXfixes-devel wayland-protocols-devel"
    echo ""
    echo "Arch/Manjaro:"
    echo "  sudo pacman -S vulkan-headers vulkan-icd-loader libxkbcommon libxkbcommon-x11 wayland wayland-protocols libx11 mesa libxcursor libxfixes"
    echo ""
    echo "openSUSE:"
    echo "  sudo zypper install vulkan-devel libxkbcommon-devel libxkbcommon-x11-devel wayland-devel libX11-devel Mesa-libEGL-devel libXcursor-devel libXfixes-devel wayland-protocols-devel"
    echo ""
    echo "Alpine Linux:"
    echo "  sudo apk add vulkan-headers vulkan-loader-dev libxkbcommon-dev libxkbcommon-x11 wayland-dev wayland-protocols libx11-dev mesa-dev libxcursor-dev libxfixes-dev"
    exit 1
fi

echo "✓ All dependencies installed successfully!"
