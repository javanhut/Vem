#!/usr/bin/env bash
# Main entry point for dependency checking - detects OS and calls appropriate checker

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Check Go first
"$SCRIPT_DIR/check-go.sh"

echo ""
echo "Checking platform-specific dependencies..."

GOOS=$(go env GOOS)

case "$GOOS" in
    linux)
        "$SCRIPT_DIR/check-deps-linux.sh"
        ;;
    windows)
        "$SCRIPT_DIR/check-deps-windows.sh"
        ;;
    darwin)
        "$SCRIPT_DIR/check-deps-darwin.sh"
        ;;
    *)
        echo "✓ No platform-specific dependencies required for $GOOS"
        ;;
esac

echo ""
echo "Checking Go module dependencies..."
go mod download
go mod verify
echo "✓ All dependencies satisfied!"
