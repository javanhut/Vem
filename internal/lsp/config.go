package lsp

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ServerConfig holds configuration for a language server.
type ServerConfig struct {
	Name           string                 // Display name
	Command        string                 // Executable name or path
	Args           []string               // Command arguments
	FileExtensions []string               // Associated file extensions
	TriggerChars   []string               // Completion trigger characters
	RootPatterns   []string               // Patterns to identify project root
	InitOptions    map[string]interface{} // Server-specific init options
	InstallCommand []string               // Install command and args
}

// DefaultServers returns hardcoded configurations for common language servers.
func DefaultServers() map[string]ServerConfig {
	return map[string]ServerConfig{
		"go": {
			Name:           "gopls",
			Command:        "gopls",
			Args:           []string{"serve"},
			FileExtensions: []string{".go"},
			TriggerChars:   []string{".", "/"},
			RootPatterns:   []string{"go.mod", "go.work", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"go", "install", "golang.org/x/tools/gopls@latest"},
		},
		"python": {
			Name:           "pyright",
			Command:        "pyright-langserver",
			Args:           []string{"--stdio"},
			FileExtensions: []string{".py", ".pyi"},
			TriggerChars:   []string{".", "/"},
			RootPatterns:   []string{"pyproject.toml", "setup.py", "requirements.txt", "pyrightconfig.json", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "pyright"},
		},
		"python-pylsp": {
			Name:           "python-lsp-server",
			Command:        "pylsp",
			Args:           []string{},
			FileExtensions: []string{".py", ".pyi"},
			TriggerChars:   []string{".", "/"},
			RootPatterns:   []string{"pyproject.toml", "setup.py", "requirements.txt", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"pip", "install", "python-lsp-server"},
		},
		"typescript": {
			Name:           "typescript-language-server",
			Command:        "typescript-language-server",
			Args:           []string{"--stdio"},
			FileExtensions: []string{".ts", ".tsx"},
			TriggerChars:   []string{".", "/", "<"},
			RootPatterns:   []string{"tsconfig.json", "package.json", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "typescript-language-server", "typescript"},
		},
		"javascript": {
			Name:           "typescript-language-server",
			Command:        "typescript-language-server",
			Args:           []string{"--stdio"},
			FileExtensions: []string{".js", ".jsx", ".mjs", ".cjs"},
			TriggerChars:   []string{".", "/", "<"},
			RootPatterns:   []string{"jsconfig.json", "package.json", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "typescript-language-server", "typescript"},
		},
		"rust": {
			Name:           "rust-analyzer",
			Command:        "rust-analyzer",
			Args:           []string{},
			FileExtensions: []string{".rs"},
			TriggerChars:   []string{".", ":", "<"},
			RootPatterns:   []string{"Cargo.toml", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"rustup", "component", "add", "rust-analyzer"},
		},
		"c": {
			Name:           "clangd",
			Command:        "clangd",
			Args:           []string{"--background-index"},
			FileExtensions: []string{".c", ".h"},
			TriggerChars:   []string{".", ":", "<", "/"},
			RootPatterns:   []string{"compile_commands.json", "CMakeLists.txt", "Makefile", ".git"},
			InitOptions:    map[string]interface{}{},
		},
		"cpp": {
			Name:           "clangd",
			Command:        "clangd",
			Args:           []string{"--background-index"},
			FileExtensions: []string{".cpp", ".hpp", ".cc", ".hh", ".cxx", ".hxx", ".c++", ".h++"},
			TriggerChars:   []string{".", ":", "<", "/"},
			RootPatterns:   []string{"compile_commands.json", "CMakeLists.txt", "Makefile", ".git"},
			InitOptions:    map[string]interface{}{},
		},
		"java": {
			Name:           "jdtls",
			Command:        "jdtls",
			Args:           []string{},
			FileExtensions: []string{".java"},
			TriggerChars:   []string{".", "@"},
			RootPatterns:   []string{"pom.xml", "build.gradle", "build.gradle.kts", ".git"},
			InitOptions:    map[string]interface{}{},
		},
		"lua": {
			Name:           "lua-language-server",
			Command:        "lua-language-server",
			Args:           []string{},
			FileExtensions: []string{".lua"},
			TriggerChars:   []string{".", ":"},
			RootPatterns:   []string{".luarc.json", ".luarc.jsonc", ".git"},
			InitOptions:    map[string]interface{}{},
		},
		"bash": {
			Name:           "bash-language-server",
			Command:        "bash-language-server",
			Args:           []string{"start"},
			FileExtensions: []string{".sh", ".bash", ".zsh"},
			TriggerChars:   []string{"$", "/"},
			RootPatterns:   []string{".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "bash-language-server"},
		},
		"html": {
			Name:           "vscode-html-language-server",
			Command:        "vscode-html-language-server",
			Args:           []string{"--stdio"},
			FileExtensions: []string{".html", ".htm"},
			TriggerChars:   []string{"<", "/", "\"", "="},
			RootPatterns:   []string{"package.json", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "vscode-langservers-extracted"},
		},
		"css": {
			Name:           "vscode-css-language-server",
			Command:        "vscode-css-language-server",
			Args:           []string{"--stdio"},
			FileExtensions: []string{".css", ".scss", ".less"},
			TriggerChars:   []string{":", " "},
			RootPatterns:   []string{"package.json", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "vscode-langservers-extracted"},
		},
		"json": {
			Name:           "vscode-json-language-server",
			Command:        "vscode-json-language-server",
			Args:           []string{"--stdio"},
			FileExtensions: []string{".json", ".jsonc"},
			TriggerChars:   []string{},
			RootPatterns:   []string{".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "vscode-langservers-extracted"},
		},
		"yaml": {
			Name:           "yaml-language-server",
			Command:        "yaml-language-server",
			Args:           []string{"--stdio"},
			FileExtensions: []string{".yaml", ".yml"},
			TriggerChars:   []string{":"},
			RootPatterns:   []string{".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "yaml-language-server"},
		},
		"zig": {
			Name:           "zls",
			Command:        "zls",
			Args:           []string{},
			FileExtensions: []string{".zig"},
			TriggerChars:   []string{".", "@"},
			RootPatterns:   []string{"build.zig", ".git"},
			InitOptions:    map[string]interface{}{},
		},
		"ruby": {
			Name:           "solargraph",
			Command:        "solargraph",
			Args:           []string{"stdio"},
			FileExtensions: []string{".rb", ".rake"},
			TriggerChars:   []string{".", ":"},
			RootPatterns:   []string{"Gemfile", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"gem", "install", "solargraph"},
		},
		"php": {
			Name:           "intelephense",
			Command:        "intelephense",
			Args:           []string{"--stdio"},
			FileExtensions: []string{".php"},
			TriggerChars:   []string{"$", ">", ":"},
			RootPatterns:   []string{"composer.json", ".git"},
			InitOptions:    map[string]interface{}{},
			InstallCommand: []string{"npm", "install", "-g", "intelephense"},
		},
		"elixir": {
			Name:           "elixir-ls",
			Command:        "elixir-ls",
			Args:           []string{},
			FileExtensions: []string{".ex", ".exs"},
			TriggerChars:   []string{".", "@"},
			RootPatterns:   []string{"mix.exs", ".git"},
			InitOptions:    map[string]interface{}{},
		},
		"kotlin": {
			Name:           "kotlin-language-server",
			Command:        "kotlin-language-server",
			Args:           []string{},
			FileExtensions: []string{".kt", ".kts"},
			TriggerChars:   []string{".", ":"},
			RootPatterns:   []string{"build.gradle", "build.gradle.kts", ".git"},
			InitOptions:    map[string]interface{}{},
		},
		"swift": {
			Name:           "sourcekit-lsp",
			Command:        "sourcekit-lsp",
			Args:           []string{},
			FileExtensions: []string{".swift"},
			TriggerChars:   []string{"."},
			RootPatterns:   []string{"Package.swift", ".git"},
			InitOptions:    map[string]interface{}{},
		},
		"csharp": {
			Name:           "omnisharp",
			Command:        "OmniSharp",
			Args:           []string{"-lsp"},
			FileExtensions: []string{".cs"},
			TriggerChars:   []string{".", " "},
			RootPatterns:   []string{"*.csproj", "*.sln", ".git"},
			InitOptions:    map[string]interface{}{},
		},
	}
}

// GetConfigForFile returns the server config for a file path based on extension.
func GetConfigForFile(filePath string) *ServerConfig {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		// Check for extensionless shell scripts
		base := filepath.Base(filePath)
		if isShellScript(filePath) {
			cfg := DefaultServers()["bash"]
			return &cfg
		}
		// No extension and not a shell script
		if base == "Makefile" || base == "makefile" || base == "GNUmakefile" {
			return nil // No LSP for makefiles
		}
		return nil
	}

	servers := DefaultServers()

	for _, cfg := range servers {
		for _, e := range cfg.FileExtensions {
			if e == ext {
				cfgCopy := cfg
				return &cfgCopy
			}
		}
	}

	return nil
}

// GetConfigByIdentifier returns a server config by language key or server name.
func GetConfigByIdentifier(identifier string) *ServerConfig {
	id := strings.ToLower(strings.TrimSpace(identifier))
	if id == "" {
		return nil
	}

	servers := DefaultServers()
	if cfg, ok := servers[id]; ok {
		cfgCopy := cfg
		return &cfgCopy
	}

	for _, cfg := range servers {
		if strings.ToLower(cfg.Name) == id {
			cfgCopy := cfg
			return &cfgCopy
		}
	}

	return nil
}

// isShellScript checks if a file is a shell script by reading the shebang.
func isShellScript(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first line
	buf := make([]byte, 256)
	n, err := file.Read(buf)
	if err != nil || n < 2 {
		return false
	}

	line := string(buf[:n])
	if strings.HasPrefix(line, "#!") {
		firstLine := strings.Split(line, "\n")[0]
		return strings.Contains(firstLine, "sh") ||
			strings.Contains(firstLine, "bash") ||
			strings.Contains(firstLine, "zsh")
	}

	return false
}

// IsServerAvailable checks if a language server is installed and executable.
func IsServerAvailable(cfg *ServerConfig) bool {
	if cfg == nil {
		return false
	}

	if resolved, ok := resolveServerPath(cfg); ok {
		cfg.Command = resolved
		return true
	}

	return false
}

func resolveServerPath(cfg *ServerConfig) (string, bool) {
	if cfg == nil || cfg.Command == "" {
		return "", false
	}

	if path, err := exec.LookPath(cfg.Command); err == nil {
		return path, true
	}

	switch cfg.Command {
	case "gopls":
		if path := findGoBinary("gopls"); path != "" {
			return path, true
		}
	case "rust-analyzer":
		if path := findCargoBinary("rust-analyzer"); path != "" {
			return path, true
		}
	}

	return "", false
}

func findGoBinary(name string) string {
	if !commandExists("go") {
		return ""
	}
	binName := withExecutableSuffix(name)
	if gobin := strings.TrimSpace(runCommandOutput("go", "env", "GOBIN")); gobin != "" {
		if path := filepath.Join(gobin, binName); fileExists(path) {
			return path
		}
	}
	if gopath := strings.TrimSpace(runCommandOutput("go", "env", "GOPATH")); gopath != "" {
		if path := filepath.Join(gopath, "bin", binName); fileExists(path) {
			return path
		}
	}
	return ""
}

func findCargoBinary(name string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	path := filepath.Join(home, ".cargo", "bin", withExecutableSuffix(name))
	if fileExists(path) {
		return path
	}
	return ""
}

func withExecutableSuffix(name string) string {
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(name), ".exe") {
		return name + ".exe"
	}
	return name
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func runCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// FindProjectRoot finds the project root directory based on marker files.
func FindProjectRoot(startPath string, patterns []string) string {
	if startPath == "" {
		return ""
	}

	// Get absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return filepath.Dir(startPath)
	}

	// If startPath is a file, start from its directory
	info, err := os.Stat(absPath)
	if err != nil {
		return filepath.Dir(absPath)
	}

	var dir string
	if info.IsDir() {
		dir = absPath
	} else {
		dir = filepath.Dir(absPath)
	}

	// Walk up the directory tree looking for root markers
	for {
		for _, pattern := range patterns {
			// Handle glob patterns
			if strings.Contains(pattern, "*") {
				matches, err := filepath.Glob(filepath.Join(dir, pattern))
				if err == nil && len(matches) > 0 {
					return dir
				}
			} else {
				// Exact match
				markerPath := filepath.Join(dir, pattern)
				if _, err := os.Stat(markerPath); err == nil {
					return dir
				}
			}
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	// No root markers found, use the file's directory
	if info.IsDir() {
		return absPath
	}
	return filepath.Dir(absPath)
}

// LanguageID returns the LSP language identifier for a file.
func LanguageID(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".go":
		return "go"
	case ".py", ".pyi":
		return "python"
	case ".ts":
		return "typescript"
	case ".tsx":
		return "typescriptreact"
	case ".js", ".mjs", ".cjs":
		return "javascript"
	case ".jsx":
		return "javascriptreact"
	case ".rs":
		return "rust"
	case ".c":
		return "c"
	case ".h":
		return "c" // C headers treated as C
	case ".cpp", ".cc", ".cxx", ".c++":
		return "cpp"
	case ".hpp", ".hh", ".hxx", ".h++":
		return "cpp"
	case ".java":
		return "java"
	case ".lua":
		return "lua"
	case ".sh", ".bash":
		return "shellscript"
	case ".zsh":
		return "shellscript"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".scss":
		return "scss"
	case ".less":
		return "less"
	case ".json":
		return "json"
	case ".jsonc":
		return "jsonc"
	case ".yaml", ".yml":
		return "yaml"
	case ".xml":
		return "xml"
	case ".md", ".markdown":
		return "markdown"
	case ".zig":
		return "zig"
	case ".rb", ".rake":
		return "ruby"
	case ".php":
		return "php"
	case ".ex", ".exs":
		return "elixir"
	case ".kt", ".kts":
		return "kotlin"
	case ".swift":
		return "swift"
	case ".cs":
		return "csharp"
	case ".sql":
		return "sql"
	case ".toml":
		return "toml"
	case ".ini", ".cfg":
		return "ini"
	case ".vim":
		return "vim"
	default:
		// Check for shell scripts without extension
		if isShellScript(filePath) {
			return "shellscript"
		}
		return "plaintext"
	}
}

// FilePathToURI converts a file path to a document URI.
func FilePathToURI(path string) DocumentURI {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	// Clean path and convert to forward slashes
	absPath = filepath.ToSlash(absPath)
	return DocumentURI("file://" + absPath)
}

// URIToFilePath converts a document URI to a file path.
func URIToFilePath(uri DocumentURI) string {
	s := string(uri)
	if strings.HasPrefix(s, "file://") {
		path := s[7:]
		// Convert forward slashes to OS-specific separator
		return filepath.FromSlash(path)
	}
	return s
}
