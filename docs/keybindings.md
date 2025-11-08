# Keybinding System Architecture

## Overview

The keybinding system in ProjectVem uses a **Command/Action pattern** inspired by Neovim's robust modal keybinding architecture. This design ensures that keybindings work consistently across all modes and platforms, avoiding fragile logic that can break due to platform-specific modifier key reporting.

## Key Concepts

### Actions

Actions are enum constants that represent discrete operations the editor can perform. Each action is mode-agnostic and can be triggered from different contexts.

Examples:
- `ActionToggleExplorer` - Toggle file explorer visibility
- `ActionMoveUp` - Move cursor/selection up
- `ActionEnterInsert` - Enter INSERT mode

### KeyBindings

A `KeyBinding` maps a key combination to an action. It consists of:
- **Modifiers**: Ctrl, Shift, Alt (optional)
- **Key**: The key name (e.g., "t", key.NameEscape)
- **Modes**: Which modes this binding applies to (empty = all modes)
- **Action**: The action to execute

### Two-Phase Matching

When a key is pressed, the system checks keybindings in priority order:

1. **Global keybindings** - Checked first, work in ANY mode
2. **Mode-specific keybindings** - Checked second, only active in specific modes
3. **Special handlers** - Custom logic for complex cases (counts, goto sequences)

This ensures that global shortcuts like `Ctrl+T` always work, regardless of which mode you're in.

## Architecture Benefits

### 1. Platform Independence

The system handles platform-specific quirks transparently:
- Some platforms don't report Ctrl modifier correctly
- Some platforms send different key codes for the same physical key
- The `modifiersMatch()` function handles these cases with fallback logic

### 2. Testability

Each component can be tested independently:
- `matchGlobalKeybinding()` - Test global shortcut matching
- `matchModeKeybinding()` - Test mode-specific bindings
- `executeAction()` - Test action execution logic

### 3. Debuggability

Clear separation makes debugging easier:
- Log which binding matched
- Log which action was executed
- Trace the full path from keypress → binding → action → result

### 4. Extensibility

Adding new keybindings is straightforward:

```go
// Add a new action
const (
    ActionSplitWindow Action = iota + 100
)

// Add a global keybinding
var globalKeybindings = []KeyBinding{
    {Modifiers: key.ModCtrl, Key: "w", Modes: nil, Action: ActionSplitWindow},
}

// Implement the action
func (s *appState) executeAction(action Action, ev key.Event) {
    switch action {
    case ActionSplitWindow:
        s.splitWindow()
    }
}
```

### 5. User Customization (Future)

This architecture makes it easy to support user-customizable keybindings:
- Load keybindings from config file
- Allow per-mode keybinding overrides
- Support keybinding "layers" (base + user + plugin)

## How It Works

### Event Flow

```
User presses key
    ↓
handleEvents() receives key.Event
    ↓
handleKey() (if State == Press)
    ↓
Phase 1: matchGlobalKeybinding()
    ↓ (if no match)
Phase 2: matchModeKeybinding()
    ↓ (if no match)
Phase 3: handleModeSpecial() (counts, etc.)
    ↓
executeAction()
    ↓
State update + re-render
```

### Example: Ctrl+T

1. User presses `Ctrl+T` while in EXPLORER mode
2. `handleKey()` receives the event
3. `matchGlobalKeybinding()` checks all global bindings
4. Finds match: `{Modifiers: key.ModCtrl, Key: "t", Action: ActionToggleExplorer}`
5. Calls `executeAction(ActionToggleExplorer, ev)`
6. `executeAction()` toggles explorer visibility
7. Event handling returns (mode-specific handlers never run)

### Example: Normal Mode 'j'

1. User presses `j` in NORMAL mode
2. `handleKey()` receives the event
3. `matchGlobalKeybinding()` finds no match (no global binding for 'j')
4. `matchModeKeybinding(modeNormal, ev)` checks NORMAL mode bindings
5. Finds match: `{Key: "j", Action: ActionMoveDown}`
6. Calls `executeAction(ActionMoveDown, ev)`
7. Cursor moves down one line

## Modifier Matching

The `modifiersMatch()` function handles platform quirks:

```go
func (s *appState) modifiersMatch(ev key.Event, required key.Modifiers) bool {
    if required == 0 {
        return ev.Modifiers == 0  // No modifiers required
    }
    
    // Use both ev.Modifiers AND s.ctrlPressed for robustness
    ctrlHeld := ev.Modifiers.Contain(key.ModCtrl) || s.ctrlPressed
    
    if required.Contain(key.ModCtrl) && !ctrlHeld {
        return false
    }
    // ... check other modifiers
}
```

This ensures Ctrl bindings work even if `ev.Modifiers` doesn't report Ctrl correctly.

## Special Handlers

Some features require more complex logic than simple key → action mapping:

### Counts (e.g., "5j" to move down 5 lines)

Handled in `handleNormalModeSpecial()`:
- Accumulates digits into `pendingCount`
- Applies count when motion is executed

### Goto Sequences (e.g., "gg" to go to top)

Handled via state machine:
- First 'g' sets `pendingGoto = true`
- Second 'g' or 'G' executes goto with accumulated count

### Colon Commands (e.g., ":w" to save)

Detected via `isColonKey()` which handles both `:` and `Shift+;`

## File Organization

- `internal/appcore/keybindings.go` - Action enum, bindings registry, matching logic
- `internal/appcore/app.go` - Event handling, special handlers, action execution

## Future Enhancements

1. **Config File Support**: Load user keybindings from `~/.vemrc`
2. **Plugin Keybindings**: Allow plugins to register their own actions and bindings
3. **Keybinding Conflicts**: Detect and warn about conflicting bindings
4. **Visual Keybinding Editor**: GUI for customizing keybindings
5. **Macro Recording**: Record and replay action sequences
