# input_unix.go / input_darwin.go / input_windows.go / debug_darwin.go / debug_stub.go

**Paths:**
- `/home/javanhut/Development/Vem/internal/appcore/input_unix.go` (43 lines)
- `/home/javanhut/Development/Vem/internal/appcore/input_darwin.go` (155 lines)
- `/home/javanhut/Development/Vem/internal/appcore/input_windows.go` (79 lines)
- `/home/javanhut/Development/Vem/internal/appcore/debug_darwin.go` (59 lines)
- `/home/javanhut/Development/Vem/internal/appcore/debug_stub.go` (14 lines)

**Purpose:** Platform-specific modifier key handling and debug logging

## Overview

Handles platform-specific differences in modifier key events:
- **Windows:** Works around a Gio bug where modifier Release events arrive before character keys
- **macOS:** Same bug as Windows + treats Command key as equivalent to Ctrl
- **Linux/BSD:** Standard modifier tracking (events arrive in correct order)

## The macOS/Windows Bug

Both macOS and Windows have the same Gio bug: **modifier Release events arrive BEFORE the character key event**.

```
Example timeline when user presses Shift+G:
 1. User presses Shift   → Shift Press event fires, shiftPressed=true
 2. User presses G       → NO EVENT YET
 3. User releases Shift  → Shift Release event fires, shiftPressed=false  ← PROBLEM!
 4. Character "G" arrives → key.Event with ev.Modifiers == empty, shift=false!
```

**Solution:** Track the timestamp of modifier Release events. If a character key arrives within 200ms of a modifier Release, we know that modifier was held during the key press.

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

### Temporal Logic Fix

macOS uses the same temporal tracking as Windows because modifier Release events arrive before character key events:

1. When Shift/Ctrl/Cmd is released, record `time.Now()` in release time fields
2. When a character key arrives, check if release time is within 200ms
3. If within window, set modifier as pressed

### Debug Logging

Debug logging is controlled by the `VEM_DEBUG_DARWIN` environment variable.

### handleModifierEvent (Lines 39-100)
Handles modifier key events on macOS with temporal tracking:
- **Command key Release**: Records `ctrlReleaseTime`
- **Command key Press**: Sets `ctrlPressed = true`
- **Ctrl key**: Same behavior (for user preference)
- **Shift key**: Records `shiftReleaseTime` on Release, sets `shiftPressed` on Press

```go
if e.Name == key.NameCommand {
    if e.State == key.Release {
        s.ctrlReleaseTime = time.Now()  // Track release time
    } else {
        s.ctrlPressed = true
    }
    return true
}
```

### syncModifierState (Lines 105-154)
Syncs modifier state using temporal logic:

```go
now := time.Now()

// Check if Command/Ctrl was released within last 200ms
ctrlWindow := now.Sub(s.ctrlReleaseTime)
if ctrlWindow < 200*time.Millisecond && ctrlWindow >= 0 {
    s.ctrlPressed = true
}

// Check if Shift was released within last 200ms
shiftWindow := now.Sub(s.shiftReleaseTime)
if shiftWindow < 200*time.Millisecond && shiftWindow >= 0 {
    s.shiftPressed = true
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

### Usage

Enable debugging on macOS:
```bash
VEM_DEBUG_DARWIN=1 ./vem
```

### Expected Debug Output (After Fix)

For Shift+G (now working correctly):
```
[darwin] handleModifierEvent: Name="Shift" State=Press Modifiers=
[darwin]   -> Shift key Press: shiftPressed=true
[darwin] handleModifierEvent: Name="Shift" State=Release Modifiers=
[darwin]   -> Shift key Release: recorded shiftReleaseTime
[darwin-key] Event: Name="G" State=Press Modifiers=
[darwin-key]   tracked: ctrl=false shift=false
[darwin] syncModifierState: Name="G" Modifiers= (before: ctrl=false shift=false)
[darwin]   -> Shift released 5.123ms ago, set shiftPressed=true    ← TEMPORAL FIX!
[darwin]   -> after sync: ctrl=false shift=true
[darwin-key] printableKey: Name="G" -> rune='G' ok=true (shiftPressed=true)
```

## debug_stub.go

Build tag: `//go:build !darwin`

Contains empty stub implementations of all debug functions for non-darwin platforms.

## input_windows.go

Build tag: `//go:build windows`

### Windows Bug Description

Same bug as macOS - Release events arrive before character keys.

### handleModifierEvent (Lines 11-51)
Uses temporal logic:

```go
if e.Name == key.NameCtrl {
    if e.State == key.Release {
        s.ctrlReleaseTime = time.Now()
    } else {
        s.ctrlPressed = true
    }
    return true
}
```

### syncModifierState (Lines 53-78)
Syncs modifier state using release timestamps (same as macOS).

## appState Fields Used

```go
type appState struct {
    // Modifier tracking
    ctrlPressed      bool
    shiftPressed     bool
    ctrlReleaseTime  time.Time  // Used by Windows and macOS
    shiftReleaseTime time.Time  // Used by Windows and macOS
}
```

## Known Issues / Potential Bugs

1. **200ms window may be too long**
   - Could cause false positives with fast typing
   - Could be too short for slow typists

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
| Modifier Release timing | Correct | Out-of-order | Out-of-order |
| ev.Modifiers field | Accurate | Often empty | Often empty |
| Command key handling | N/A | Maps to Ctrl | N/A |
| Solution | Direct tracking | Temporal logic | Temporal logic |
| Timing dependency | No | 200ms window | 200ms window |
| Debug logging | No | VEM_DEBUG_DARWIN | No |

## Keybinding Behavior by Platform

| Shortcut | Linux | macOS | Windows |
|----------|-------|-------|---------|
| Toggle Explorer | Ctrl+T | Cmd+T or Ctrl+T | Ctrl+T |
| Fuzzy Finder | Ctrl+F | Cmd+F or Ctrl+F | Ctrl+F |
| Undo | Ctrl+U | Cmd+U or Ctrl+U | Ctrl+U |
| Copy | Ctrl+C | Cmd+C or Ctrl+C | Ctrl+C |
| Paste | Ctrl+P | Cmd+P or Ctrl+P | Ctrl+P |
| Jump to last line | Shift+G | Shift+G | Shift+G |

---
*Last Updated: Fixed macOS modifier timing bug using temporal logic (same approach as Windows)*
