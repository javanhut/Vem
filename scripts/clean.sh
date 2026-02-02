#!/usr/bin/env bash
# Remove build artifacts

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "Cleaning build artifacts..."
rm -f vem vem.exe vem-darwin-amd64 vem-darwin-arm64
echo "✓ Clean complete."
