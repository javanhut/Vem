# Makefile for ProjectVem
# Cross-platform build, install, and dependency management

.PHONY: help build install uninstall clean test check-deps check-deps-linux check-deps-windows check-deps-darwin check-go check-vulkan check-xkbcommon check-xkbcommon-x11 check-wayland check-wayland-cursor check-x11 check-egl check-xcursor check-xfixes install-linux-deps build-windows build-linux build-darwin

# Detect OS and Architecture
GOOS := $(shell go env GOOS 2>/dev/null)
GOARCH := $(shell go env GOARCH 2>/dev/null)

# Minimum required Go version
MIN_GO_VERSION := 1.21

# Binary name
BINARY_NAME := vem
ifeq ($(GOOS),windows)
	BINARY_NAME := vem.exe
endif

# Install location (Linux/macOS only)
INSTALL_DIR := /usr/local/bin
INSTALL_PATH := $(INSTALL_DIR)/vem

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "ProjectVem Makefile"
	@echo ""
	@echo "Detected OS: $(GOOS)"
	@echo "Detected Architecture: $(GOARCH)"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

check-go: ## Check if Go is installed and meets version requirements
	@echo "Checking Go installation..."
	@if ! command -v go >/dev/null 2>&1; then \
		echo "✗ Go is not installed."; \
		echo "Please install Go from https://golang.org/dl/"; \
		exit 1; \
	fi
	@echo "✓ Go is installed: $$(go version)"
	@GO_VERSION=$$(go version | sed -n 's/.*go\([0-9]*\.[0-9]*\).*/\1/p'); \
	MAJOR=$$(echo $$GO_VERSION | cut -d. -f1); \
	MINOR=$$(echo $$GO_VERSION | cut -d. -f2); \
	MIN_MAJOR=$$(echo $(MIN_GO_VERSION) | cut -d. -f1); \
	MIN_MINOR=$$(echo $(MIN_GO_VERSION) | cut -d. -f2); \
	if [ $$MAJOR -lt $$MIN_MAJOR ] || ([ $$MAJOR -eq $$MIN_MAJOR ] && [ $$MINOR -lt $$MIN_MINOR ]); then \
		echo "✗ Go version $$GO_VERSION is too old. Minimum required: $(MIN_GO_VERSION)"; \
		exit 1; \
	fi
	@echo "✓ Go version meets requirements (>= $(MIN_GO_VERSION))"
	@if [ -z "$(GOOS)" ] || [ -z "$(GOARCH)" ]; then \
		echo "✗ Cannot detect GOOS or GOARCH"; \
		exit 1; \
	fi
	@echo "✓ GOOS=$(GOOS), GOARCH=$(GOARCH)"

check-vulkan: ## Check if Vulkan headers are installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for Vulkan headers..."
	@if [ -f /usr/include/vulkan/vulkan.h ]; then \
		echo "✓ Vulkan headers found."; \
	else \
		echo "✗ Vulkan headers not found."; \
		exit 1; \
	fi
else
	@echo "Vulkan check not required on $(GOOS)"
endif

check-xkbcommon: ## Check if xkbcommon is installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for xkbcommon..."
	@if pkg-config --exists xkbcommon 2>/dev/null || [ -f /usr/include/xkbcommon/xkbcommon.h ]; then \
		echo "✓ xkbcommon found."; \
	else \
		echo "✗ xkbcommon not found."; \
		exit 1; \
	fi
else
	@echo "xkbcommon check not required on $(GOOS)"
endif

check-wayland: ## Check if wayland-client is installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for wayland-client..."
	@if pkg-config --exists wayland-client 2>/dev/null || [ -f /usr/include/wayland-client.h ]; then \
		echo "✓ wayland-client found."; \
	else \
		echo "✗ wayland-client not found."; \
		exit 1; \
	fi
else
	@echo "wayland-client check not required on $(GOOS)"
endif

check-x11: ## Check if X11 headers are installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for X11 headers..."
	@if pkg-config --exists x11 2>/dev/null || [ -f /usr/include/X11/Xlib.h ]; then \
		echo "✓ X11 headers found."; \
	else \
		echo "✗ X11 headers not found."; \
		exit 1; \
	fi
else
	@echo "X11 check not required on $(GOOS)"
endif

check-egl: ## Check if EGL (mesa) is installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for EGL..."
	@if pkg-config --exists egl 2>/dev/null || [ -f /usr/include/EGL/egl.h ]; then \
		echo "✓ EGL found."; \
	else \
		echo "✗ EGL not found."; \
		exit 1; \
	fi
else
	@echo "EGL check not required on $(GOOS)"
endif

check-xkbcommon-x11: ## Check if xkbcommon-x11 is installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for xkbcommon-x11..."
	@if pkg-config --exists xkbcommon-x11 2>/dev/null; then \
		echo "✓ xkbcommon-x11 found."; \
	else \
		echo "✗ xkbcommon-x11 not found."; \
		exit 1; \
	fi
else
	@echo "xkbcommon-x11 check not required on $(GOOS)"
endif

check-xcursor: ## Check if libxcursor is installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for libxcursor..."
	@if pkg-config --exists xcursor 2>/dev/null || [ -f /usr/include/X11/Xcursor/Xcursor.h ]; then \
		echo "✓ libxcursor found."; \
	else \
		echo "✗ libxcursor not found."; \
		exit 1; \
	fi
else
	@echo "libxcursor check not required on $(GOOS)"
endif

check-xfixes: ## Check if libxfixes is installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for libxfixes..."
	@if pkg-config --exists xfixes 2>/dev/null || [ -f /usr/include/X11/extensions/Xfixes.h ]; then \
		echo "✓ libxfixes found."; \
	else \
		echo "✗ libxfixes not found."; \
		exit 1; \
	fi
else
	@echo "libxfixes check not required on $(GOOS)"
endif

check-wayland-cursor: ## Check if wayland-cursor is installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for wayland-cursor..."
	@if pkg-config --exists wayland-cursor 2>/dev/null; then \
		echo "✓ wayland-cursor found."; \
	else \
		echo "✗ wayland-cursor not found."; \
		exit 1; \
	fi
else
	@echo "wayland-cursor check not required on $(GOOS)"
endif

check-windows-conpty: ## Check if Windows supports ConPTY (Windows 10 1809+)
ifeq ($(GOOS),windows)
	@echo "Checking Windows version for ConPTY support..."
	@powershell -Command "if ([System.Environment]::OSVersion.Version.Build -lt 17763) { Write-Host '✗ Windows version too old. ConPTY requires Windows 10 1809 (build 17763) or later'; exit 1 } else { Write-Host '✓ Windows version supports ConPTY' }"
else
	@echo "Windows ConPTY check not required on $(GOOS)"
endif

check-windows-shell: ## Check if a suitable shell is available on Windows
ifeq ($(GOOS),windows)
	@echo "Checking for available shells on Windows..."
	@SHELL_FOUND=0; \
	if command -v pwsh.exe >/dev/null 2>&1; then \
		echo "✓ PowerShell Core (pwsh.exe) found"; \
		SHELL_FOUND=1; \
	fi; \
	if command -v powershell.exe >/dev/null 2>&1; then \
		echo "✓ Windows PowerShell (powershell.exe) found"; \
		SHELL_FOUND=1; \
	fi; \
	if [ ! -z "$$COMSPEC" ]; then \
		echo "✓ Command Prompt ($$COMSPEC) found"; \
		SHELL_FOUND=1; \
	fi; \
	if [ $$SHELL_FOUND -eq 0 ]; then \
		echo "✗ No suitable shell found"; \
		exit 1; \
	fi
else
	@echo "Windows shell check not required on $(GOOS)"
endif

install-linux-deps: ## Install all required Linux dependencies (Vulkan, xkbcommon, wayland, X11, EGL, etc.)
ifeq ($(GOOS),linux)
	@echo "Detecting package manager and installing dependencies..."
	@if command -v apt-get >/dev/null 2>&1; then \
		echo "Using apt-get (Debian/Ubuntu)..."; \
		sudo apt-get update && sudo apt-get install -y libvulkan-dev libxkbcommon-dev libxkbcommon-x11-dev libwayland-dev libx11-dev libegl1-mesa-dev libxcursor-dev libxfixes-dev wayland-protocols; \
	elif command -v dnf >/dev/null 2>&1; then \
		echo "Using dnf (Fedora/RHEL/CentOS)..."; \
		sudo dnf install -y vulkan-devel libxkbcommon-devel libxkbcommon-x11-devel wayland-devel libX11-devel mesa-libEGL-devel libXcursor-devel libXfixes-devel wayland-protocols-devel; \
	elif command -v pacman >/dev/null 2>&1; then \
		echo "Using pacman (Arch/Manjaro)..."; \
		sudo pacman -S --noconfirm vulkan-headers vulkan-icd-loader libxkbcommon libxkbcommon-x11 wayland wayland-protocols libx11 mesa libxcursor libxfixes; \
	elif command -v zypper >/dev/null 2>&1; then \
		echo "Using zypper (openSUSE)..."; \
		sudo zypper install -y vulkan-devel libxkbcommon-devel libxkbcommon-x11-devel wayland-devel libX11-devel Mesa-libEGL-devel libXcursor-devel libXfixes-devel wayland-protocols-devel; \
	elif command -v apk >/dev/null 2>&1; then \
		echo "Using apk (Alpine Linux)..."; \
		sudo apk add vulkan-headers vulkan-loader-dev libxkbcommon-dev libxkbcommon-x11 wayland-dev wayland-protocols libx11-dev mesa-dev libxcursor-dev libxfixes-dev; \
	else \
		echo "Error: No supported package manager found."; \
		echo "Please install the following dependencies manually:"; \
		echo ""; \
		echo "Debian/Ubuntu:"; \
		echo "  sudo apt-get install libvulkan-dev libxkbcommon-dev libxkbcommon-x11-dev libwayland-dev libx11-dev libegl1-mesa-dev libxcursor-dev libxfixes-dev wayland-protocols"; \
		echo ""; \
		echo "Fedora/RHEL/CentOS:"; \
		echo "  sudo dnf install vulkan-devel libxkbcommon-devel libxkbcommon-x11-devel wayland-devel libX11-devel mesa-libEGL-devel libXcursor-devel libXfixes-devel wayland-protocols-devel"; \
		echo ""; \
		echo "Arch/Manjaro:"; \
		echo "  sudo pacman -S vulkan-headers vulkan-icd-loader libxkbcommon libxkbcommon-x11 wayland wayland-protocols libx11 mesa libxcursor libxfixes"; \
		echo ""; \
		echo "openSUSE:"; \
		echo "  sudo zypper install vulkan-devel libxkbcommon-devel libxkbcommon-x11-devel wayland-devel libX11-devel Mesa-libEGL-devel libXcursor-devel libXfixes-devel wayland-protocols-devel"; \
		echo ""; \
		echo "Alpine Linux:"; \
		echo "  sudo apk add vulkan-headers vulkan-loader-dev libxkbcommon-dev libxkbcommon-x11 wayland-dev wayland-protocols libx11-dev mesa-dev libxcursor-dev libxfixes-dev"; \
		exit 1; \
	fi
	@echo "All dependencies installed successfully!"
else
	@echo "Linux dependency installation not required on $(GOOS)"
endif

check-deps-linux: ## Check all Linux-specific dependencies
ifeq ($(GOOS),linux)
	@echo "Checking all Linux dependencies..."
	@MISSING=0; \
	if ! $(MAKE) check-vulkan 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-xkbcommon 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-xkbcommon-x11 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-wayland 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-wayland-cursor 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-x11 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-egl 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-xcursor 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-xfixes 2>/dev/null; then MISSING=1; fi; \
	if [ $$MISSING -eq 1 ]; then \
		echo ""; \
		echo "Some dependencies are missing. Installing..."; \
		$(MAKE) install-linux-deps; \
	else \
		echo ""; \
		echo "✓ All Linux dependencies satisfied!"; \
	fi
else
	@echo "Linux dependency checks not required on $(GOOS)"
endif

check-deps-windows: ## Check all Windows-specific dependencies
ifeq ($(GOOS),windows)
	@echo "Checking all Windows dependencies..."
	@MISSING=0; \
	if ! $(MAKE) check-windows-conpty 2>/dev/null; then MISSING=1; fi; \
	if ! $(MAKE) check-windows-shell 2>/dev/null; then MISSING=1; fi; \
	if [ $$MISSING -eq 1 ]; then \
		echo ""; \
		echo "✗ Some Windows dependencies are missing."; \
		echo "Please ensure you have:"; \
		echo "  - Windows 10 1809 (build 17763) or later for ConPTY support"; \
		echo "  - PowerShell or Command Prompt available"; \
		exit 1; \
	else \
		echo ""; \
		echo "✓ All Windows dependencies satisfied!"; \
	fi
else
	@echo "Windows dependency checks not required on $(GOOS)"
endif

check-deps-darwin: ## Check all macOS-specific dependencies
ifeq ($(GOOS),darwin)
	@echo "Checking macOS dependencies..."
	@echo "✓ macOS has all required dependencies built-in (Metal, Cocoa)"
else
	@echo "macOS dependency checks not required on $(GOOS)"
endif

check-deps: check-go ## Check and install all platform-specific dependencies
	@echo ""
	@echo "Checking platform-specific dependencies for $(GOOS)..."
ifeq ($(GOOS),linux)
	@$(MAKE) check-deps-linux
else ifeq ($(GOOS),windows)
	@$(MAKE) check-deps-windows
else ifeq ($(GOOS),darwin)
	@$(MAKE) check-deps-darwin
else
	@echo "✓ No platform-specific dependencies required for $(GOOS)"
endif
	@echo ""
	@echo "Checking Go module dependencies..."
	@go mod download
	@go mod verify
	@echo "✓ All dependencies satisfied!"

build: check-deps ## Build Vem binary for current OS/architecture
	@echo ""
	@echo "==========================================="
	@echo "Building Vem for $(GOOS)/$(GOARCH)..."
	@echo "==========================================="
	@echo ""
	@echo "Target binary: $(BINARY_NAME)"
	@if go build -v -o $(BINARY_NAME) . 2>&1; then \
		echo ""; \
		echo "✓ Build successful: $(BINARY_NAME)"; \
		echo ""; \
		if [ "$(GOOS)" = "windows" ]; then \
			echo "Windows build complete."; \
			echo "You can run the editor with: ./$(BINARY_NAME)"; \
		else \
			echo "Run 'make install' to install to $(INSTALL_DIR)"; \
			echo "Or run directly with: ./$(BINARY_NAME)"; \
		fi; \
	else \
		echo ""; \
		echo "✗ Build failed!"; \
		echo ""; \
		if [ "$(GOOS)" = "windows" ]; then \
			echo "Windows build troubleshooting:"; \
			echo "  - Ensure you have Windows 10 1809 or later for ConPTY support"; \
			echo "  - Check that Go CGO is properly configured"; \
			echo "  - Verify all Go module dependencies are available"; \
		elif [ "$(GOOS)" = "linux" ]; then \
			echo "Linux build troubleshooting:"; \
			echo "  - Run 'make check-deps' to verify all dependencies"; \
			echo "  - Ensure development headers are installed"; \
		else \
			echo "Build troubleshooting:"; \
			echo "  - Run 'make check-go' to verify Go installation"; \
			echo "  - Run 'go mod tidy' to fix module issues"; \
		fi; \
		exit 1; \
	fi

build-windows: check-go ## Cross-compile for Windows (amd64)
	@echo "Cross-compiling for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build -v -o vem.exe .
	@echo "✓ Windows build complete: vem.exe"

build-linux: check-go ## Cross-compile for Linux (amd64)
	@echo "Cross-compiling for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build -v -o vem .
	@echo "✓ Linux build complete: vem"
	@echo "Note: This binary may require Linux system libraries (Vulkan, X11, Wayland) on the target system."

build-darwin: check-go ## Cross-compile for macOS (amd64 and arm64)
	@echo "Cross-compiling for macOS (amd64)..."
	@if GOOS=darwin GOARCH=amd64 go build -v -o vem-darwin-amd64 . 2>&1; then \
		echo "✓ macOS amd64 build complete: vem-darwin-amd64"; \
		echo ""; \
		echo "Cross-compiling for macOS (arm64/Apple Silicon)..."; \
		if GOOS=darwin GOARCH=arm64 go build -v -o vem-darwin-arm64 . 2>&1; then \
			echo "✓ macOS arm64 build complete: vem-darwin-arm64"; \
		else \
			echo "✗ macOS arm64 build failed"; \
			echo "Note: Cross-compiling to macOS may require macOS SDK and CGO"; \
			exit 1; \
		fi; \
	else \
		echo "✗ macOS amd64 build failed"; \
		echo "Note: Cross-compiling to macOS requires:"; \
		echo "  - macOS SDK (for CGO dependencies)"; \
		echo "  - Proper CGO cross-compilation setup"; \
		echo "  - Consider building natively on macOS instead"; \
		exit 1; \
	fi

install: build ## Install Vem to system
ifeq ($(GOOS),windows)
	@echo "Windows detected: Installation to system PATH not automated."
	@echo "Binary built as $(BINARY_NAME)"
	@echo "To use Vem, either:"
	@echo "  1. Add the current directory to your PATH"
	@echo "  2. Move $(BINARY_NAME) to a directory in your PATH"
else
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@if [ ! -d $(INSTALL_DIR) ]; then \
		echo "Creating $(INSTALL_DIR)..."; \
		sudo mkdir -p $(INSTALL_DIR); \
	fi
	@sudo install -m 755 $(BINARY_NAME) $(INSTALL_PATH)
	@echo "Installation complete. Run 'vem' to start the editor."
endif

uninstall: ## Uninstall Vem from system
ifeq ($(GOOS),windows)
	@echo "Windows detected: No system installation to remove."
	@echo "If you manually installed $(BINARY_NAME), please remove it manually."
else
	@if [ -f $(INSTALL_PATH) ]; then \
		echo "Removing $(INSTALL_PATH)..."; \
		sudo rm -f $(INSTALL_PATH); \
		echo "Uninstall complete."; \
	else \
		echo "Vem is not installed at $(INSTALL_PATH)"; \
	fi
endif

clean: ## Remove built binaries
	@echo "Cleaning build artifacts..."
	@rm -f vem vem.exe
	@echo "Clean complete."

test: ## Run all tests
	@echo "Running tests..."
	@go test ./...
