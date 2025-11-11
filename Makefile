# Makefile for ProjectVem
# Cross-platform build, install, and dependency management

.PHONY: help build install uninstall clean test check-deps check-vulkan check-xkbcommon check-xkbcommon-x11 check-wayland check-wayland-cursor check-x11 check-egl check-xcursor check-xfixes install-linux-deps

# Detect OS and Architecture
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

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
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

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

check-deps: ## Check and install all dependencies
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
		echo "✓ All dependencies satisfied!"; \
	fi
else
	@echo "Dependency checks not required on $(GOOS)"
endif

build: check-deps ## Build Vem binary for current OS/architecture
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@go build -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

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
