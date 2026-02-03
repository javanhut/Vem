# input_unix.go / input_darwin.go / input_windows.go

**Paths:**
- `/home/javanhut/Development/Vem/internal/appcore/input_unix.go` (43 lines)
- `/home/javanhut/Development/Vem/internal/appcore/input_darwin.go` (51 lines)
- `/home/javanhut/Development/Vem/internal/appcore/input_windows.go` (79 lines)

**Purpose:** Platform-specific modifier key handling

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

### handleModifierEvent (Lines 9-33)
Handles modifier key events on macOS:
- **Command key**: Sets `ctrlPressed` (macOS convention)
- **Ctrl key**: Also sets `ctrlPressed` (for user preference)
- **Shift key**: Sets `shiftPressed`

```go
// Track Command key as ctrlPressed (macOS convention)
if e.Name == key.NameCommand {
    s.ctrlPressed = (e.State == key.Press)
    return true
}

// Also track actual Ctrl (for users who prefer Ctrl+key)
if e.Name == key.NameCtrl {
    s.ctrlPressed = (e.State == key.Press)
    return true
}
```

### syncModifierState (Lines 35-48)
Syncs modifier state before handling character keys:
- Checks **both** `ModCommand` and `ModCtrl` from `ev.Modifiers`
- Either modifier will set `ctrlPressed = true`

```go
if e.Modifiers.Contain(key.ModCommand) && !s.ctrlPressed {
    s.ctrlPressed = true
}
if e.Modifiers.Contain(key.ModCtrl) && !s.ctrlPressed {
    s.ctrlPressed = true
}
```

## input_windows.go

Build tag: `//go:build windows`

### Windows Bug Description

```
Example timeline when user presses Ctrl+T:
 1. User presses Ctrl     → NO EVENT (bug!)
 2. User presses T        → NO EVENT YET
 3. User releases Ctrl    → Ctrl Release event fires
 4. Character "T" arrives → key.Event with ev.Modifiers == empty
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

## Dead/Unused Code

None identified.

## Integration Points

- Called from main key event handler in app.go
- Used before processing keybindings
- Affects all keyboard input handling

## Platform Differences

| Behavior | Linux/BSD | macOS | Windows |
|----------|-----------|-------|---------|
| Modifier Press events | Work | Work | Never arrive |
| ev.Modifiers field | Accurate | Accurate | Often empty |
| Command key handling | N/A | Maps to Ctrl | N/A |
| Solution | Direct tracking | Cmd→Ctrl mapping | Temporal logic |
| Timing dependency | No | No | 200ms window |

## Keybinding Behavior by Platform

| Shortcut | Linux | macOS | Windows |
|----------|-------|-------|---------|
| Toggle Explorer | Ctrl+T | Cmd+T or Ctrl+T | Ctrl+T |
| Fuzzy Finder | Ctrl+F | Cmd+F or Ctrl+F | Ctrl+F |
| Undo | Ctrl+U | Cmd+U or Ctrl+U | Ctrl+U |
| Copy | Ctrl+C | Cmd+C or Ctrl+C | Ctrl+C |
| Paste | Ctrl+P | Cmd+P or Ctrl+P | Ctrl+P |

---
*Last Updated: macOS Command key support added*
