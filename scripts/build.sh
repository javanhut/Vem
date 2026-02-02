#!/usr/bin/env bash
# Build vem binary for current platform

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

# Determine binary name
BINARY_NAME="vem"
if [ "$GOOS" = "windows" ]; then
    BINARY_NAME="vem.exe"
fi

echo ""
echo "==========================================="
echo "Building Vem for $GOOS/$GOARCH..."
echo "==========================================="
echo ""
echo "Target binary: $BINARY_NAME"

if go build -v -o "$BINARY_NAME" .; then
    echo ""
    echo "✓ Build successful: $BINARY_NAME"
    echo ""
    if [ "$GOOS" = "windows" ]; then
        echo "Windows build complete."
        echo "You can run the editor with: ./$BINARY_NAME"
    else
        echo "Run 'imlazy install' to install to /usr/local/bin"
        echo "Or run directly with: ./$BINARY_NAME"
    fi
else
    echo ""
    echo "✗ Build failed!"
    echo ""
    if [ "$GOOS" = "windows" ]; then
        echo "Windows build troubleshooting:"
        echo "  - Ensure you have Windows 10 1809 or later for ConPTY support"
        echo "  - Check that Go CGO is properly configured"
        echo "  - Verify all Go module dependencies are available"
    elif [ "$GOOS" = "linux" ]; then
        echo "Linux build troubleshooting:"
        echo "  - Run 'imlazy check:deps' to verify all dependencies"
        echo "  - Ensure development headers are installed"
    else
        echo "Build troubleshooting:"
        echo "  - Run 'imlazy check:go' to verify Go installation"
        echo "  - Run 'go mod tidy' to fix module issues"
    fi
    exit 1
fi
