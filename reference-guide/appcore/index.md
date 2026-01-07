# appcore Package Index

**Path:** `/home/javanhut/Development/Vem/internal/appcore/`
**Purpose:** Core application logic, UI rendering, and event handling

## Files Overview

| File | Lines | Purpose |
|------|-------|---------|
| [app.go](app.md) | 4000+ | Main application state and event loop |
| [keybindings.go](keybindings.md) | 870+ | Keybinding definitions and action execution |
| [buffer_completion.go](buffer_completion.md) | 112 | Buffer word completion (Ctrl+N style) |
| [command_completion.go](command_completion.md) | 175 | Command mode path completion (Tab) |
| [fuzzy.go](fuzzy.md) | 179 | Fuzzy file finder algorithm |
| [help.go](help.md) | 289 | Help text generation |
| [input_unix.go](input.md) | 43 | Unix modifier key handling |
| [input_windows.go](input.md) | 79 | Windows modifier key handling |
| [lsp_actions.go](lsp_actions.md) | 1249 | LSP feature handlers |
| [lsp_rendering.go](lsp_rendering.md) | 905 | LSP UI components + buffer completion + diagnostics list |
| [pane_actions.go](pane_actions.md) | 375 | Pane management handlers |
| [pane_rendering.go](pane_rendering.md) | 319 | Pane tree rendering |

## Key Types

### appState (app.go)
Central application state containing:
- Mode management
- Buffer manager
- Pane manager
- LSP manager
- Terminal instances
- UI state (viewport, selection, etc.)

### Action Constants (keybindings.go)
All editor actions like `ActionMoveLeft`, `ActionEnterInsert`, etc.

### KeyBinding (keybindings.go)
```go
type KeyBinding struct {
    Modifiers key.Modifiers
    Key       key.Name
    Modes     []Mode
    Action    Action
}
```

## Known Issues Summary

1. ~~**app.go:650** - `os.ReadFile()` loads entire file into memory~~ MITIGATED - Large files >5MB show warning, >50MB rejected
2. ~~**keybindings.go:257** - `u` key not bound to undo in Normal mode~~ FIXED
3. **pane_actions.go:16-48** - Debug print statements should be removed
4. **pane_actions.go:349-350** - j/k inverted for vertical resize
5. **lsp_rendering.go:533** - Hardcoded 8px character width

## Recent Fixes

- **Syntax highlighter cache** - Fixed highlighting disappearing when opening files from explorer (buffer index mismatch after closing empty buffer)
- **Cursor positioning** - Fixed visual offset caused by token-based rendering vs single-string measurement
- **Tab key stuck** - Fixed `skipNextEdit` timing race condition
- **Format on save** - Added `:formatonsave` and `:format` commands
- **Soft wrap** - Added `:wrap` command for dynamic text wrapping (enabled by default)

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                    main.go                          │
│              (creates window, calls Run)            │
└────────────────────┬────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────┐
│                   appcore.Run()                     │
│              (initializes appState)                 │
└────────────────────┬────────────────────────────────┘
                     │
     ┌───────────────┴───────────────┐
     ▼                               ▼
┌──────────────┐              ┌──────────────┐
│  Event Loop  │              │   Layout     │
│   (run())    │◄────────────►│  Rendering   │
└──────────────┘              └──────────────┘
     │                               │
     │ key.Event                     │ layout.Context
     ▼                               ▼
┌──────────────┐              ┌──────────────┐
│ keybindings  │              │ drawBuffer() │
│ handleKey()  │              │ drawPanes()  │
│ executeAction│              │ drawStatus() │
└──────────────┘              └──────────────┘
```

## Event Flow

1. `run()` receives Gio events
2. `key.Event` dispatched to `handleKeyEvent()`
3. Modifiers synced via `syncModifierState()`
4. Keybinding matched in `modeKeybindings` or `globalKeybindings`
5. Action executed via `executeAction()`
6. Buffer/UI state updated
7. `window.Invalidate()` triggers re-render
8. Layout functions called to draw UI

---
*Last Updated: Reference guide creation*
