# Makefile for ProjectVem
# Cross-platform build, install, and dependency management

.PHONY: help build install uninstall clean test check-deps check-vulkan install-vulkan

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
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

check-vulkan: ## Check if Vulkan headers are installed (Linux only)
ifeq ($(GOOS),linux)
	@echo "Checking for Vulkan headers..."
	@if [ -f /usr/include/vulkan/vulkan.h ]; then \
		echo "Vulkan headers found."; \
	else \
		echo "Vulkan headers not found. Attempting to install..."; \
		$(MAKE) install-vulkan; \
	fi
else
	@echo "Vulkan check not required on $(GOOS)"
endif

install-vulkan: ## Install Vulkan headers (Linux only)
ifeq ($(GOOS),linux)
	@echo "Detecting package manager..."
	@if command -v apt-get >/dev/null 2>&1; then \
		echo "Using apt-get..."; \
		sudo apt-get update && sudo apt-get install -y libvulkan-dev; \
	elif command -v dnf >/dev/null 2>&1; then \
		echo "Using dnf..."; \
		sudo dnf install -y vulkan-devel; \
	elif command -v pacman >/dev/null 2>&1; then \
		echo "Using pacman..."; \
		sudo pacman -S --noconfirm vulkan-headers vulkan-icd-loader; \
	else \
		echo "Error: No supported package manager found (apt-get, dnf, or pacman)"; \
		echo "Please install Vulkan headers manually:"; \
		echo "  Debian/Ubuntu: sudo apt-get install libvulkan-dev"; \
		echo "  Fedora/RHEL: sudo dnf install vulkan-devel"; \
		echo "  Arch: sudo pacman -S vulkan-headers vulkan-icd-loader"; \
		exit 1; \
	fi
else
	@echo "Vulkan installation not required on $(GOOS)"
endif

check-deps: check-vulkan ## Check and install all dependencies

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
