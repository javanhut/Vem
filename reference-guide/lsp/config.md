# config.go

**Path:** `/home/javanhut/Development/Vem/internal/lsp/config.go`
**Lines:** ~545
**Purpose:** Language server configuration, command discovery, and project root detection

## Overview

Provides:
- Pre-configured language server definitions for 21 languages
- File extension to server mapping
- Shebang detection for extensionless scripts
- Project root detection via marker files
- Language ID resolution for LSP protocol
- URI/path conversion utilities

## Code Blocks

### Lines 1-9: Package and Imports

```go
package lsp
import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)
```

### Lines 10-20: ServerConfig Struct

```go
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
```

### Lines 22-237: Default Server Configurations

#### DefaultServers (Lines 22-237)
Returns map of 21 language server configs:

| Language | Server | Command | Extensions |
|----------|--------|---------|------------|
| go | gopls | `gopls serve` | .go |
| python | pyright | `pyright-langserver --stdio` | .py, .pyi |
| python-pylsp | python-lsp-server | `pylsp` | .py, .pyi |
| typescript | typescript-language-server | `--stdio` | .ts, .tsx |
| javascript | typescript-language-server | `--stdio` | .js, .jsx, .mjs, .cjs |
| rust | rust-analyzer | (no args) | .rs |
| c | clangd | `--background-index` | .c, .h |
| cpp | clangd | `--background-index` | .cpp, .hpp, .cc, etc. |
| java | jdtls | (no args) | .java |
| lua | lua-language-server | (no args) | .lua |
| bash | bash-language-server | `start` | .sh, .bash, .zsh |
| html | vscode-html-language-server | `--stdio` | .html, .htm |
| css | vscode-css-language-server | `--stdio` | .css, .scss, .less |
| json | vscode-json-language-server | `--stdio` | .json, .jsonc |
| yaml | yaml-language-server | `--stdio` | .yaml, .yml |
| zig | zls | (no args) | .zig |
| ruby | solargraph | `stdio` | .rb, .rake |
| php | intelephense | `--stdio` | .php |
| elixir | elixir-ls | (no args) | .ex, .exs |
| kotlin | kotlin-language-server | (no args) | .kt, .kts |
| swift | sourcekit-lsp | (no args) | .swift |
| csharp | OmniSharp | `-lsp` | .cs |

### Lines 239-268: Config Lookup

#### GetConfigForFile (Lines 239-268)
Gets server config for file path:
1. Gets lowercase file extension
2. If no extension, checks for shell script via shebang
3. Handles Makefile (returns nil - no LSP)
4. Searches all configs for matching extension
5. Returns copy of matching config

### Lines 270-291: Config by Name

#### GetConfigByIdentifier (Lines 270-291)
Gets config by language key or server name:
1. Tries exact language key match (e.g., "go", "rust")
2. Tries server name match (e.g., "gopls", "rust-analyzer")

### Lines 293-317: Shell Script Detection

#### isShellScript (Lines 293-317)
Checks if file is shell script by shebang:
1. Opens file and reads first 256 bytes
2. Checks for `#!` prefix
3. Looks for "sh", "bash", or "zsh" in first line

### Lines 320-381: Server Command Discovery

#### findGoExecutable (Lines 320-354)
Locates Go-installed executables by checking:
1. `GOBIN` environment variable
2. `GOPATH/bin` directory
3. Default `~/go/bin` location
4. `/usr/local/go/bin` fallback

This allows gopls and other Go tools to be found even when not in system PATH.

#### FindServerCommand (Lines 356-376)
Returns full path to server executable:
1. First tries `exec.LookPath()` (system PATH)
2. For gopls, falls back to Go-specific locations via `findGoExecutable()`
3. Returns empty string if not found

```go
func FindServerCommand(cfg *ServerConfig) string {
    // Try PATH first
    if path, err := exec.LookPath(cfg.Command); err == nil {
        return path
    }
    // For Go tools, check Go-specific locations
    if cfg.Command == "gopls" {
        if path, err := findGoExecutable("gopls"); err == nil {
            return path
        }
    }
    return ""
}
```

#### IsServerAvailable (Lines 378-381)
Checks if server is installed using `FindServerCommand()`:
```go
func IsServerAvailable(cfg *ServerConfig) bool {
    return FindServerCommand(cfg) != ""
}
```

### Lines 383-440: Project Root Detection

#### FindProjectRoot (Lines 383-440)
Finds project root by marker files:
1. Gets absolute path
2. Determines starting directory
3. Walks up directory tree looking for patterns
4. Handles glob patterns (e.g., `*.csproj`)
5. Falls back to file's directory if no markers found

**Example Root Patterns:**
- Go: `go.mod`, `go.work`, `.git`
- Node: `package.json`, `.git`
- Rust: `Cargo.toml`, `.git`

### Lines 387-469: Language ID Mapping

#### LanguageID (Lines 387-469)
Returns LSP language identifier for file:

| Extension | Language ID |
|-----------|-------------|
| .go | go |
| .py, .pyi | python |
| .ts | typescript |
| .tsx | typescriptreact |
| .js, .mjs, .cjs | javascript |
| .jsx | javascriptreact |
| .rs | rust |
| .c, .h | c |
| .cpp, .cc, .cxx | cpp |
| .java | java |
| .sh, .bash, .zsh | shellscript |
| .html, .htm | html |
| .css | css |
| .json | json |
| .yaml, .yml | yaml |
| .md, .markdown | markdown |
| (unknown) | plaintext |

### Lines 471-491: URI Utilities

| Function | Lines | Purpose |
|----------|-------|---------|
| `FilePathToURI()` | 471-480 | Convert path to `file://` URI |
| `URIToFilePath()` | 482-491 | Convert URI back to path |

## Known Issues / Potential Bugs

1. **Line 293-317: isShellScript opens files**
   - Opens every extensionless file to check shebang
   - Could be slow for many files

2. **Line 250-251: Makefile handling**
   - Returns nil for Makefiles
   - Could support make-language-server

## Dead/Unused Code

None identified.

## Integration Points

- Called by `Manager.GetServerForFile()` in manager.go
- Used by `document.go` for language IDs
- Used by `features.go` for URI conversion

## Adding New Language Servers

1. Add entry to `DefaultServers()` map
2. Include all fields:
   - `Name`: Display name
   - `Command`: Executable name
   - `Args`: Command line arguments
   - `FileExtensions`: Supported extensions
   - `TriggerChars`: Completion triggers
   - `RootPatterns`: Project markers
   - `InstallCommand`: (optional) Installation command

---
*Last Updated: After gopls PATH detection fix*
