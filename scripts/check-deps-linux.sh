#!/usr/bin/env bash
# Check all Linux dependencies (vulkan, xkbcommon, wayland, x11, egl, etc.)

set -e

GOOS=$(go env GOOS 2>/dev/null || echo "")

if [ "$GOOS" != "linux" ]; then
    echo "Linux dependency checks not required on $GOOS"
    exit 0
fi

echo "Checking all Linux dependencies..."

MISSING=0

# Check Vulkan headers
echo "Checking for Vulkan headers..."
if [ -f /usr/include/vulkan/vulkan.h ]; then
    echo "✓ Vulkan headers found."
else
    echo "✗ Vulkan headers not found."
    MISSING=1
fi

# Check xkbcommon
echo "Checking for xkbcommon..."
if pkg-config --exists xkbcommon 2>/dev/null || [ -f /usr/include/xkbcommon/xkbcommon.h ]; then
    echo "✓ xkbcommon found."
else
    echo "✗ xkbcommon not found."
    MISSING=1
fi

# Check xkbcommon-x11
echo "Checking for xkbcommon-x11..."
if pkg-config --exists xkbcommon-x11 2>/dev/null; then
    echo "✓ xkbcommon-x11 found."
else
    echo "✗ xkbcommon-x11 not found."
    MISSING=1
fi

# Check wayland-client
echo "Checking for wayland-client..."
if pkg-config --exists wayland-client 2>/dev/null || [ -f /usr/include/wayland-client.h ]; then
    echo "✓ wayland-client found."
else
    echo "✗ wayland-client not found."
    MISSING=1
fi

# Check wayland-cursor
echo "Checking for wayland-cursor..."
if pkg-config --exists wayland-cursor 2>/dev/null; then
    echo "✓ wayland-cursor found."
else
    echo "✗ wayland-cursor not found."
    MISSING=1
fi

# Check X11 headers
echo "Checking for X11 headers..."
if pkg-config --exists x11 2>/dev/null || [ -f /usr/include/X11/Xlib.h ]; then
    echo "✓ X11 headers found."
else
    echo "✗ X11 headers not found."
    MISSING=1
fi

# Check x11-xcb
echo "Checking for x11-xcb..."
if pkg-config --exists x11-xcb 2>/dev/null; then
    echo "✓ x11-xcb found."
else
    echo "✗ x11-xcb not found. Install libx11-xcb-dev (Debian/Ubuntu) or libX11-xcb (Fedora/Arch)"
    MISSING=1
fi

# Check EGL
echo "Checking for EGL..."
if pkg-config --exists egl 2>/dev/null || [ -f /usr/include/EGL/egl.h ]; then
    echo "✓ EGL found."
else
    echo "✗ EGL not found."
    MISSING=1
fi

# Check libxcursor
echo "Checking for libxcursor..."
if pkg-config --exists xcursor 2>/dev/null || [ -f /usr/include/X11/Xcursor/Xcursor.h ]; then
    echo "✓ libxcursor found."
else
    echo "✗ libxcursor not found."
    MISSING=1
fi

# Check libxfixes
echo "Checking for libxfixes..."
if pkg-config --exists xfixes 2>/dev/null || [ -f /usr/include/X11/extensions/Xfixes.h ]; then
    echo "✓ libxfixes found."
else
    echo "✗ libxfixes not found."
    MISSING=1
fi

if [ $MISSING -eq 1 ]; then
    echo ""
    echo "Some dependencies are missing."
    echo "Run 'imlazy deps:install' to install them automatically."
    exit 1
else
    echo ""
    echo "✓ All Linux dependencies satisfied!"
fi
