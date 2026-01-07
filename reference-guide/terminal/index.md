# terminal Package Index

**Path:** `/home/javanhut/Development/Vem/internal/terminal/`
**Purpose:** Integrated terminal emulator with PTY support

## Files Overview

| File | Lines | Purpose |
|------|-------|---------|
| [terminal.go](terminal.md) | 400+ | VT100 terminal emulator |
| [buffer.go](buffer.md) | 237 | Screen buffer with cells |
| [colors.go](colors.md) | 93 | ANSI color palette |
| [input.go](input.md) | 138 | Key to escape sequence conversion |
| [pty_unix.go](pty_unix.md) | 107 | Unix PTY implementation |
| [pty_windows.go](pty_windows.md) | 95 | Windows ConPTY implementation |
| [conpty_wrapper.go](conpty_wrapper.md) | 15 | Windows ConPTY wrapper |

## Key Types

### Terminal (terminal.go)
Main terminal controller:
- PTY/ConPTY management
- vt10x emulator integration
- Screen buffer
- Input/output goroutines

### ScreenBuffer (buffer.go)
Terminal display state:
- Cell grid (character + attributes)
- Cursor position and style
- Dirty line tracking

### Cell (buffer.go)
```go
type Cell struct {
    Rune      rune
    FG        color.NRGBA
    BG        color.NRGBA
    Bold      bool
    Dim       bool
    Italic    bool
    Underline bool
    Blink     bool
    Reverse   bool
}
```

## Platform Support

| Platform | PTY Library | Default Shell |
|----------|-------------|---------------|
| Linux/macOS | creack/pty | $SHELL or /bin/sh |
| Windows | UserExistsError/conpty | pwsh > powershell > cmd |

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Vem Application                    │
│                                                     │
│  ┌─────────────────────────────────────────────┐   │
│  │              Terminal struct                 │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────────┐  │   │
│  │  │ vt10x   │  │ Screen  │  │ Input Chan  │  │   │
│  │  │ parser  │  │ Buffer  │  │             │  │   │
│  │  └────┬────┘  └────┬────┘  └──────┬──────┘  │   │
│  └───────│───────────│──────────────│──────────┘   │
│          │           │              │              │
└──────────│───────────│──────────────│──────────────┘
           │           │              │
           ▼           ▼              ▼
    ┌──────────────────────────────────────────┐
    │              PTY / ConPTY                │
    │         (kernel pseudo-terminal)         │
    └────────────────────┬─────────────────────┘
                         │
                         ▼
    ┌──────────────────────────────────────────┐
    │            Shell Process                 │
    │         (bash, zsh, pwsh, etc.)          │
    └──────────────────────────────────────────┘
```

## Data Flow

```
User Input → KeyToTerminalSequence() → inputChan → writeLoop() → PTY
                                                                  ↓
                                                           Shell Process
                                                                  ↓
PTY → readLoop() → vt10x parser → ScreenBuffer → Rendering
```

## Key Sequences

| Key | Sequence |
|-----|----------|
| Enter | `\r` |
| Tab | `\t` |
| Backspace | `\x7f` |
| Up Arrow | `\x1b[A` |
| Ctrl+C | `\x03` |
| Ctrl+D | `\x04` |

See [input.md](input.md) for complete mapping.

## Known Issues Summary

1. **terminal.go:177** - 100ms read deadline may cause performance issues
2. **terminal.go:298** - 2 second wait timeout may leave zombie processes

---
*Last Updated: Reference guide creation*
