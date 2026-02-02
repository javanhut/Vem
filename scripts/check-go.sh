#!/usr/bin/env bash
# Check if Go is installed and meets version requirements

set -e

MIN_GO_VERSION="1.21"

echo "Checking Go installation..."

if ! command -v go >/dev/null 2>&1; then
    echo "✗ Go is not installed."
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "✓ Go is installed: $(go version)"

GO_VERSION=$(go version | sed -n 's/.*go\([0-9]*\.[0-9]*\).*/\1/p')
MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
MINOR=$(echo "$GO_VERSION" | cut -d. -f2)
MIN_MAJOR=$(echo "$MIN_GO_VERSION" | cut -d. -f1)
MIN_MINOR=$(echo "$MIN_GO_VERSION" | cut -d. -f2)

if [ "$MAJOR" -lt "$MIN_MAJOR" ] || ([ "$MAJOR" -eq "$MIN_MAJOR" ] && [ "$MINOR" -lt "$MIN_MINOR" ]); then
    echo "✗ Go version $GO_VERSION is too old. Minimum required: $MIN_GO_VERSION"
    exit 1
fi

echo "✓ Go version meets requirements (>= $MIN_GO_VERSION)"

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

if [ -z "$GOOS" ] || [ -z "$GOARCH" ]; then
    echo "✗ Cannot detect GOOS or GOARCH"
    exit 1
fi

echo "✓ GOOS=$GOOS, GOARCH=$GOARCH"
