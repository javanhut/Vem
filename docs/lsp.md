# Language Server Protocol (LSP) Support

Vem includes built-in Language Server Protocol support, providing IDE-like features for many programming languages.

## Supported Languages

Vem automatically detects and connects to language servers for the following languages:

| Language | Server | Installation |
|----------|--------|--------------|
| Go | gopls | `go install golang.org/x/tools/gopls@latest` |
| Python | pyright | `npm install -g pyright` |
| TypeScript/JavaScript | typescript-language-server | `npm install -g typescript-language-server typescript` |
| Rust | rust-analyzer | `rustup component add rust-analyzer` |
| C/C++ | clangd | Install via package manager |
| Java | jdtls | Install Eclipse JDT Language Server |
| Lua | lua-language-server | Install via package manager or LuaRocks |
| Bash | bash-language-server | `npm install -g bash-language-server` |
| HTML | vscode-html-language-server | `npm install -g vscode-langservers-extracted` |
| CSS | vscode-css-language-server | `npm install -g vscode-langservers-extracted` |
| JSON | vscode-json-language-server | `npm install -g vscode-langservers-extracted` |
| YAML | yaml-language-server | `npm install -g yaml-language-server` |
| Zig | zls | Install from https://github.com/zigtools/zls |
| Ruby | solargraph | `gem install solargraph` |
| PHP | intelephense | `npm install -g intelephense` |
| Elixir | elixir-ls | Install from https://github.com/elixir-lsp/elixir-ls |
| Kotlin | kotlin-language-server | Install from https://github.com/fwcd/kotlin-language-server |
| Swift | sourcekit-lsp | Included with Swift toolchain |
| C# | OmniSharp | Install from https://github.com/OmniSharp/omnisharp-roslyn |

## Features

### Go to Definition

Navigate to the definition of a symbol under the cursor.

- **Command**: `:goDef` or `:gd` or `:definition`
- **Keybinding**: `gd` (in normal mode, press `g` then `d`)

### Hover Information

Display type information and documentation for the symbol under the cursor.

- **Command**: `:hover`
- **Keybinding**: `K` (Shift+k in normal mode)

Press Escape to dismiss the hover tooltip.

### Find References

Find all references to the symbol under the cursor.

- **Command**: `:refs` or `:references`

Navigation in the references list:
- `j` / Down Arrow: Next reference
- `k` / Up Arrow: Previous reference
- `Enter`: Open selected reference
- `Escape`: Close references list

### Rename Symbol

Rename a symbol across all files in the project.

- **Command**: `:rename <newname>`

Example: `:rename newFunctionName`

### Document Formatting

Format the current document using the language server's formatter.

- **Command**: `:format`

### Code Actions

Show available quick fixes and refactorings at the cursor position.

- **Command**: `:codeaction` or `:ca`

Navigation in the code actions menu:
- `j` / Down Arrow: Next action
- `k` / Up Arrow: Previous action
- `Enter`: Apply selected action
- `Escape`: Close menu

### Auto-completion

Trigger completion suggestions while typing.

- **Keybinding**: `Ctrl+Space` (in insert mode)
- **Auto-trigger**: Completion also triggers automatically after `.` and other trigger characters

Navigation in the completion menu:
- `Ctrl+n` / Down Arrow: Next item
- `Ctrl+p` / Up Arrow: Previous item
- `Tab` or `Enter`: Accept completion
- `Escape`: Cancel completion

### Diagnostics

Errors and warnings from the language server are displayed as underlines in the editor:
- Red wavy underline: Error
- Orange wavy underline: Warning
- Blue wavy underline: Information
- Gray wavy underline: Hint

Navigate between diagnostics:
- `]d`: Jump to next diagnostic
- `[d`: Jump to previous diagnostic

The status bar shows a count of errors and warnings for the current file.

## LSP Commands

| Command | Description |
|---------|-------------|
| `:goDef` | Go to definition |
| `:hover` | Show hover information |
| `:refs` | Find all references |
| `:rename <name>` | Rename symbol |
| `:format` | Format document |
| `:codeaction` | Show code actions |
| `:lspstatus` | Show LSP server status |
| `:lsprestart` | Restart language server |
| `:install <lang>` | Install language server |
| `:uninstall <lang>` | Uninstall language server |
| `:lspauto [on|off|toggle|status]` | Auto-install missing language servers |
| `:formatonsave [on|off|toggle|status]` | Toggle format on save |
| `:lintonsave [on|off|toggle|status]` | Toggle lint on save |

## Keybindings Summary

### Normal Mode

| Key | Action |
|-----|--------|
| `gd` | Go to definition |
| `K` | Show hover information |
| `]d` | Next diagnostic |
| `[d` | Previous diagnostic |

### Insert Mode

| Key | Action |
|-----|--------|
| `Ctrl+Space` | Trigger completion |
| `Ctrl+n` | Next completion item |
| `Ctrl+p` | Previous completion item |
| `Tab` / `Enter` | Accept completion |
| `Escape` | Cancel completion |

## How It Works

1. **Automatic Detection**: When you open a file, Vem checks if there's a language server configured for that file type.

2. **Project Root Detection**: Vem finds the project root by looking for marker files (like `go.mod`, `package.json`, `Cargo.toml`, etc.).

3. **Server Lifecycle**: Language servers are started on-demand when you first use an LSP feature. They remain running until you close Vem.

4. **Document Synchronization**: Changes to your files are automatically synchronized with the language server, enabling real-time diagnostics and accurate completions.

## Troubleshooting

### Language Server Not Found

If you see "language server not installed", ensure the server is in your PATH:

```bash
# Check if gopls is available
which gopls

# Check if typescript-language-server is available
which typescript-language-server
```

### No Diagnostics Appearing

1. Ensure the language server supports diagnostics
2. Save the file to trigger a full diagnostic check
3. Check `:lspstatus` to verify the server is running

### Slow Completions

Some language servers (especially for large projects) may take time to index. Wait for initial indexing to complete.

### Server Crashes

Use `:lsprestart` to restart the language server. If crashes persist, check the language server's logs or update to the latest version.

## Configuration

LSP support is enabled by default. Language server configurations are built-in and do not require manual configuration.

Supported servers are automatically detected based on file extension. If multiple servers support the same language, Vem uses the first available server.

## Architecture

The LSP implementation consists of:

- **JSON-RPC Client**: Handles communication with language servers via stdin/stdout
- **Server Manager**: Manages lifecycle of language server processes
- **Document Sync**: Keeps servers updated with file changes
- **Feature Handlers**: Implement go-to-definition, hover, completion, etc.
- **UI Rendering**: Draws diagnostics, completion popups, hover tooltips
