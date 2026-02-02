#!/usr/bin/env bash
# Check Windows dependencies (ConPTY, shell availability)

set -e

GOOS=$(go env GOOS 2>/dev/null || echo "")

if [ "$GOOS" != "windows" ]; then
    echo "Windows dependency checks not required on $GOOS"
    exit 0
fi

echo "Checking all Windows dependencies..."

MISSING=0

# Check Windows version for ConPTY support
echo "Checking Windows version for ConPTY support..."
if command -v powershell >/dev/null 2>&1; then
    BUILD=$(powershell -Command "[System.Environment]::OSVersion.Version.Build" 2>/dev/null || echo "0")
    if [ "$BUILD" -lt 17763 ]; then
        echo "✗ Windows version too old. ConPTY requires Windows 10 1809 (build 17763) or later"
        MISSING=1
    else
        echo "✓ Windows version supports ConPTY (build $BUILD)"
    fi
else
    echo "⚠ Cannot verify Windows version (PowerShell not available in this shell)"
fi

# Check for available shells
echo "Checking for available shells on Windows..."
SHELL_FOUND=0

if command -v pwsh.exe >/dev/null 2>&1; then
    echo "✓ PowerShell Core (pwsh.exe) found"
    SHELL_FOUND=1
fi

if command -v powershell.exe >/dev/null 2>&1; then
    echo "✓ Windows PowerShell (powershell.exe) found"
    SHELL_FOUND=1
fi

if [ -n "$COMSPEC" ]; then
    echo "✓ Command Prompt ($COMSPEC) found"
    SHELL_FOUND=1
fi

if [ $SHELL_FOUND -eq 0 ]; then
    echo "✗ No suitable shell found"
    MISSING=1
fi

if [ $MISSING -eq 1 ]; then
    echo ""
    echo "✗ Some Windows dependencies are missing."
    echo "Please ensure you have:"
    echo "  - Windows 10 1809 (build 17763) or later for ConPTY support"
    echo "  - PowerShell or Command Prompt available"
    exit 1
else
    echo ""
    echo "✓ All Windows dependencies satisfied!"
fi
