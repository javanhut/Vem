#!/usr/bin/env bash
# Check macOS dependencies (minimal - built-in)

set -e

GOOS=$(go env GOOS 2>/dev/null || echo "")

if [ "$GOOS" != "darwin" ]; then
    echo "macOS dependency checks not required on $GOOS"
    exit 0
fi

echo "Checking macOS dependencies..."
echo "✓ macOS has all required dependencies built-in (Metal, Cocoa)"
