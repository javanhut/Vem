# input_unix.go / input_darwin.go / input_windows.go / debug_darwin.go / debug_stub.go

**Paths:**
- `/home/javanhut/Development/Vem/internal/appcore/input_unix.go` (43 lines)
- `/home/javanhut/Development/Vem/internal/appcore/input_darwin.go` (98 lines)
- `/home/javanhut/Development/Vem/internal/appcore/input_windows.go` (79 lines)
- `/home/javanhut/Development/Vem/internal/appcore/debug_darwin.go` (59 lines)
- `/home/javanhut/Development/Vem/internal/appcore/debug_stub.go` (14 lines)

**Purpose:** Platform-specific modifier key handling and debug logging

## Overview

Handles platform-specific differences in modifier key events:
- **Windows:** Works around a Gio bug where Ctrl/Shift Press events never arrive
- **macOS:** Treats Command key as equivalent to Ctrl for keybinding purposes
- **Linux/BSD:** Standard modifier tracking

## Platform Build Tags

| File | Build Tag | Platforms |
|------|-----------|-----------|
| `input_unix.go` | `!windows && !darwin` | Linux, BSD, other Unix |
| `input_darwin.go` | `darwin` | macOS |
| `input_windows.go` | `windows` | Windows |
| `debug_darwin.go` | `darwin` | macOS debug logging |
| `debug_stub.go` | `!darwin` | Empty stubs for other platforms |

## input_unix.go

Build tag: `//go:build !windows && !darwin`

### handleModifierEvent (Lines 9-28)
Handles modifier key events on Linux/BSD:
- Tracks Ctrl press/release state
- Tracks Shift press/release state
- Returns true if event was a modifier

```go
if e.Name == key.NameCtrl {
    s.ctrlPressed = (e.State == key.Press)
    return true
}
```

### syncModifierState (Lines 30-42)
Syncs modifier state before handling character keys:
- Uses `ev.Modifiers` as fallback
- On Unix, ev.Modifiers usually works correctly

## input_darwin.go

Build tag: `//go:build darwin`

### macOS Command Key Support

On macOS, users expect Cmd+key shortcuts (not Ctrl+key). This file maps the Command key to the internal `ctrlPressed` state, allowing all existing `key.ModCtrl` keybindings to work with Cmd.

**Both Cmd+key AND Ctrl+key work on macOS** - users can use either based on preference.

### Debug Logging

Debug logging is controlled by the `VEM_DEBUG_DARWIN` environment variable. When set, comprehensive logging is output to stderr showing:
- All modifier key events (Press/Release)
- Modifier state before and after sync
- Which modifiers are detected from `ev.Modifiers`

### handleModifierEvent (Lines 25-64)
Handles modifier key events on macOS with optional debug logging:
- **Command key**: Sets `ctrlPressed` (macOS convention)
- **Ctrl key**: Also sets `ctrlPressed` (for user preference)
- **Shift key**: Sets `shiftPressed`
- **Alt key**: Logged but ignored

```go
// Track Command key as ctrlPressed (macOS convention)
if e.Name == key.NameCommand {
    s.ctrlPressed = (e.State == key.Press)
    if darwinDebug {
        darwinInputLog("  -> Command key: ctrlPressed=%v", s.ctrlPressed)
    }
    return true
}
```

### syncModifierState (Lines 68-98)
Syncs modifier state before handling character keys with debug logging:
- Checks **both** `ModCommand` and `ModCtrl` from `ev.Modifiers`
- Either modifier will set `ctrlPressed = true`
- Logs state before and after synchronization

```go
if e.Modifiers.Contain(key.ModCommand) && !s.ctrlPressed {
    s.ctrlPressed = true
    if darwinDebug {
        darwinInputLog("  -> ModCommand detected, set ctrlPressed=true")
    }
}
```

## debug_darwin.go

Build tag: `//go:build darwin`

Contains debug logging functions for macOS input troubleshooting. All functions are no-ops unless `VEM_DEBUG_DARWIN=1` is set.

### Functions

| Function | Purpose | Called From |
|----------|---------|-------------|
| `debugKeyEvent(e key.Event)` | Logs key events with name, state, modifiers | `app.go` key event handler |
| `debugEditEvent(e key.EditEvent)` | Logs edit events and skipNextEdit state | `app.go` edit event handler |
| `debugPrintableKey(ev, r, ok)` | Logs printableKey conversion result | `app.go:printableKey()` |
| `debugModifierEvent(e, keyType, newValue)` | Logs modifier state changes | `input_darwin.go` |
| `debugSyncModifierBefore(e)` | Logs state before modifier sync | `input_darwin.go` |
| `debugSyncModifierAfter()` | Logs state after modifier sync | `input_darwin.go` |

### Usage

Enable debugging on macOS:
```bash
VEM_DEBUG_DARWIN=1 ./vem
```

### Expected Debug Output

For Shift+G (working correctly):
```
[darwin] handleModifierEvent: Name="Shift" State=Press Modifiers=0
[darwin]   -> Shift key: shiftPressed=true
[darwin-key] Event: Name="G" State=Press Modifiers=Shift
[darwin-key]   tracked: ctrl=false shift=true
[darwin-key] printableKey: Name="G" -> rune='G' ok=true (shiftPressed=true)
[darwin] handleModifierEvent: Name="Shift" State=Release Modifiers=0
[darwin]   -> Shift key: shiftPressed=false
```

For Shift+G (broken - shift not tracked):
```
[darwin-key] Event: Name="G" State=Press Modifiers=Shift
[darwin-key]   tracked: ctrl=false shift=false  <- Problem: shift not tracked!
[darwin-key] printableKey: Name="G" -> rune='g' ok=true (shiftPressed=false)  <- Wrong case!
```

## debug_stub.go

Build tag: `//go:build !darwin`

Contains empty stub implementations of all debug functions for non-darwin platforms. These are compiled into Linux, BSD, and Windows builds but do nothing.

## input_windows.go

Build tag: `//go:build windows`

### Windows Bug Description

```
Example timeline when user presses Ctrl+T:
 1. User presses Ctrl     -> NO EVENT (bug!)
 2. User presses T        -> NO EVENT YET
 3. User releases Ctrl    -> Ctrl Release event fires
 4. Character "T" arrives -> key.Event with ev.Modifiers == empty
```

### handleModifierEvent (Lines 11-51)
Uses **temporal logic** to work around Windows bug:
1. Records timestamp when modifier is released
2. A character key arriving within 200ms means modifier was held

```go
if e.Name == key.NameCtrl {
    if e.State == key.Release {
        s.ctrlReleaseTime = time.Now()
    } else {
        s.ctrlPressed = true  // Press events rarely arrive on Windows
    }
    return true
}
```

### syncModifierState (Lines 53-78)
Syncs modifier state using release timestamps:

```go
ctrlWindow := now.Sub(s.ctrlReleaseTime)
if ctrlWindow < 200*time.Millisecond && ctrlWindow >= 0 {
    s.ctrlPressed = true
}
```

## appState Fields Used

```go
type appState struct {
    // Modifier tracking
    ctrlPressed      bool
    shiftPressed     bool
    ctrlReleaseTime  time.Time  // Windows only
    shiftReleaseTime time.Time  // Windows only
}
```

## Known Issues / Potential Bugs

1. **Windows: 200ms window may be too long**
   - Could cause false positives with fast typing
   - Could be too short for slow typists

2. **macOS: Shift+G and modifier issues being investigated**
   - Use `VEM_DEBUG_DARWIN=1` to debug
   - Hypothesis: Gio may not send Shift Press/Release events on macOS

## Dead/Unused Code

None identified.

## Integration Points

- Called from main key event handler in app.go (lines ~727-775)
- Debug calls added at:
  - `app.go:730` - `debugKeyEvent()` after `case key.Event:`
  - `app.go:779` - `debugEditEvent()` after `case key.EditEvent:`
  - `app.go:4218` - `debugPrintableKey()` before return in `printableKey()`
- Used before processing keybindings
- Affects all keyboard input handling

## Platform Differences

| Behavior | Linux/BSD | macOS | Windows |
|----------|-----------|-------|---------|
| Modifier Press events | Work | Work | Never arrive |
| ev.Modifiers field | Accurate | Accurate | Often empty |
| Command key handling | N/A | Maps to Ctrl | N/A |
| Solution | Direct tracking | Cmd->Ctrl mapping | Temporal logic |
| Timing dependency | No | No | 200ms window |
| Debug logging | No | VEM_DEBUG_DARWIN | No |

## Keybinding Behavior by Platform

| Shortcut | Linux | macOS | Windows |
|----------|-------|-------|---------|
| Toggle Explorer | Ctrl+T | Cmd+T or Ctrl+T | Ctrl+T |
| Fuzzy Finder | Ctrl+F | Cmd+F or Ctrl+F | Ctrl+F |
| Undo | Ctrl+U | Cmd+U or Ctrl+U | Ctrl+U |
| Copy | Ctrl+C | Cmd+C or Ctrl+C | Ctrl+C |
| Paste | Ctrl+P | Cmd+P or Ctrl+P | Ctrl+P |

---
*Last Updated: Added debug logging for macOS input troubleshooting (debug_darwin.go, debug_stub.go)*
