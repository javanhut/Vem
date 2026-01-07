# Vem Reference Guide

A comprehensive code documentation for the Vem text editor project.

## Project Overview

**Vem** is a modern Vim emulator built in Go with GPU-accelerated rendering using Gio UI.

**Total Files:** 42+ Go source files
**Lines of Code:** ~15,000+ lines
**License:** GPLv2

## Quick Start

This reference guide documents every significant file in the Vem codebase. Use this when:
- Making code changes (understand context before editing)
- Fixing bugs (find related code quickly)
- Adding features (understand existing patterns)
- Refactoring (identify dead code and issues)

## Directory Structure

```
Vem/
├── main.go                    Entry point
├── reference-guide/           This documentation
└── internal/
    ├── appcore/               Application core (10 files, ~8000 lines)
    │   ├── app.go             Main state and event loop
    │   ├── keybindings.go     Key definitions and actions
    │   ├── fuzzy.go           Fuzzy finder algorithm
    │   ├── help.go            Help text generation
    │   ├── input_*.go         Platform-specific key handling
    │   ├── lsp_actions.go     LSP feature handlers
    │   ├── lsp_rendering.go   LSP UI components
    │   ├── pane_actions.go    Pane management
    │   └── pane_rendering.go  Pane tree rendering
    ├── editor/                Text editing (2 files, ~1200 lines)
    │   ├── buffer.go          Core buffer operations
    │   └── buffer_manager.go  Multi-buffer management
    ├── lsp/                   Language Server Protocol (6 files, ~2800 lines)
    │   ├── manager.go         Server lifecycle
    │   ├── client.go          JSON-RPC client
    │   ├── config.go          Language configurations
    │   ├── document.go        Document sync
    │   ├── features.go        LSP features
    │   └── types.go           Protocol types
    ├── panes/                 Window management (5 files, ~600 lines)
    │   └── manager.go         Pane tree management
    ├── syntax/                Highlighting (3 files, ~400 lines)
    │   ├── highlighter.go     Chroma integration
    │   ├── theme.go           Color themes
    │   └── lexers/carrion.go  Custom lexer
    ├── terminal/              Terminal emulator (7 files, ~700 lines)
    │   ├── terminal.go        VT100 emulator
    │   ├── buffer.go          Screen buffer
    │   ├── colors.go          ANSI colors
    │   ├── input.go           Key sequences
    │   └── pty_*.go           Platform PTY
    ├── filesystem/            File explorer (4 files, ~700 lines)
    │   └── tree.go            File tree structure
    └── fonts/                 Font embedding (1 file, ~50 lines)
        └── fonts.go           JetBrains Mono loading
```

## Package Index

| Package | Purpose | Index |
|---------|---------|-------|
| [appcore](./appcore/) | Main application logic | [index.md](./appcore/index.md) |
| [editor](./editor/) | Text buffer management | [index.md](./editor/index.md) |
| [lsp](./lsp/) | Language server support | [index.md](./lsp/index.md) |
| [panes](./panes/) | Window splitting | [index.md](./panes/index.md) |
| [syntax](./syntax/) | Syntax highlighting | [index.md](./syntax/index.md) |
| [terminal](./terminal/) | Terminal emulator | [index.md](./terminal/index.md) |
| [filesystem](./filesystem/) | File explorer | [index.md](./filesystem/index.md) |
| [fonts](./fonts/) | Font handling | [index.md](./fonts/index.md) |

## Key Entry Points

- **main.go** - Application entry point, creates window
- **appcore/app.go:223** - `Run()` function starts the app
- **appcore/app.go:228** - `run()` main event loop

## Editing Modes

| Mode | Constant | Description |
|------|----------|-------------|
| NORMAL | `modeNormal` | Default navigation mode |
| INSERT | `modeInsert` | Text insertion mode |
| VISUAL | `modeVisual` | Selection mode |
| DELETE | `modeDelete` | Deletion operations |
| COMMAND | `modeCommand` | Command-line (:commands) |
| EXPLORER | `modeExplorer` | File tree navigation |
| SEARCH | `modeSearch` | Pattern search |
| FUZZY_FINDER | `modeFuzzyFinder` | Quick file finding |
| TERMINAL | `modeTerminal` | Integrated terminal |

## Known Issues Summary

### High Priority

| Location | Issue | Impact |
|----------|-------|--------|
| `editor/buffer.go:650` | `os.ReadFile()` loads entire file | Large file handling |
| `appcore/keybindings.go:257` | `u` key not bound to undo | Missing Vim feature |
| `panes/manager.go:79` | Debug print statement | Console noise |

### Medium Priority

| Location | Issue | Impact |
|----------|-------|--------|
| `appcore/pane_actions.go:16-84` | Debug print statements | Console noise |
| `appcore/pane_actions.go:349-350` | j/k inverted for resize | UX confusion |
| `syntax/highlighter.go:33` | Unused `formatter` field | Dead code |
| `lsp/features.go:289` | Hardcoded formatting options | Inflexibility |

### Low Priority

| Location | Issue | Impact |
|----------|-------|--------|
| `appcore/lsp_rendering.go:533` | Hardcoded 8px char width | Potential misalignment |
| `appcore/fuzzy.go:175-178` | Unused `isWordBoundary()` | Dead code |
| `syntax/highlighter.go:57` | Hardcoded "monokai" theme | Configuration |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      main.go                                │
│                 (creates Gio window)                        │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    appcore.Run()                            │
│              (initializes all subsystems)                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ BufferMgr   │  │ PaneMgr     │  │ LSP Manager         │ │
│  │ (editor)    │  │ (panes)     │  │ (lsp)               │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ FileTree    │  │ Highlighter │  │ Terminal Manager    │ │
│  │ (filesystem)│  │ (syntax)    │  │ (terminal)          │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                     Event Loop                              │
│  Key Events → handleKeyEvent() → executeAction() → Render  │
└─────────────────────────────────────────────────────────────┘
```

## File Documentation Format

Each file is documented with:
- **Path** and **Lines** count
- **Purpose** summary
- **Code Blocks** with line ranges
- **Known Issues** specific to that file
- **Dead/Unused Code** identification
- **Integration Points** with other files

---

*Last Updated: Reference guide creation*
