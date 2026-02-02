#!/usr/bin/env bash
# Cross-compile for specified platform (linux/windows/darwin)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

TARGET="${1:-}"

if [ -z "$TARGET" ]; then
    echo "Usage: $0 <platform>"
    echo "  Platforms: linux, windows, darwin"
    exit 1
fi

case "$TARGET" in
    linux)
        echo "Cross-compiling for Linux (amd64)..."
        GOOS=linux GOARCH=amd64 go build -v -o vem .
        echo "✓ Linux build complete: vem"
        echo "Note: This binary may require Linux system libraries (Vulkan, X11, Wayland) on the target system."
        ;;
    windows)
        echo "Cross-compiling for Windows (amd64)..."
        GOOS=windows GOARCH=amd64 go build -v -o vem.exe .
        echo "✓ Windows build complete: vem.exe"
        ;;
    darwin)
        echo "Cross-compiling for macOS (amd64)..."
        if GOOS=darwin GOARCH=amd64 go build -v -o vem-darwin-amd64 . 2>&1; then
            echo "✓ macOS amd64 build complete: vem-darwin-amd64"
            echo ""
            echo "Cross-compiling for macOS (arm64/Apple Silicon)..."
            if GOOS=darwin GOARCH=arm64 go build -v -o vem-darwin-arm64 . 2>&1; then
                echo "✓ macOS arm64 build complete: vem-darwin-arm64"
            else
                echo "✗ macOS arm64 build failed"
                echo "Note: Cross-compiling to macOS may require macOS SDK and CGO"
                exit 1
            fi
        else
            echo "✗ macOS amd64 build failed"
            echo "Note: Cross-compiling to macOS requires:"
            echo "  - macOS SDK (for CGO dependencies)"
            echo "  - Proper CGO cross-compilation setup"
            echo "  - Consider building natively on macOS instead"
            exit 1
        fi
        ;;
    *)
        echo "Unknown platform: $TARGET"
        echo "  Platforms: linux, windows, darwin"
        exit 1
        ;;
esac
