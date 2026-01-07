# input_unix.go / input_windows.go

**Paths:**
- `/home/javanhut/Development/Vem/internal/appcore/input_unix.go` (43 lines)
- `/home/javanhut/Development/Vem/internal/appcore/input_windows.go` (79 lines)

**Purpose:** Platform-specific modifier key handling

## Overview

Handles a **critical Windows bug** where Ctrl/Shift Press events never arrive in Gio. Implements platform-specific workarounds for modifier key tracking.

## Windows Bug Description

```
Example timeline when user presses Ctrl+T:
 1. User presses Ctrl     → NO EVENT (bug!)
 2. User presses T        → NO EVENT YET
 3. User releases Ctrl    → Ctrl Release event fires
 4. Character "T" arrives → key.Event with ev.Modifiers == empty
```

## input_unix.go

Build tag: `//go:build !windows`

### handleModifierEvent (Lines 9-28)
Handles modifier key events on Unix/Linux/macOS:
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

## input_windows.go

Build tag: `//go:build windows`

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

1. **200ms window may be too long**
   - Could cause false positives with fast typing
   - Could be too short for slow typists

## Dead/Unused Code

None identified.

## Integration Points

- Called from main key event handler in app.go
- Used before processing keybindings
- Affects all keyboard input handling

## Platform Differences

| Behavior | Unix | Windows |
|----------|------|---------|
| Modifier Press events | Work | Never arrive |
| ev.Modifiers field | Accurate | Often empty |
| Solution | Direct tracking | Temporal logic |
| Timing dependency | No | 200ms window |

---
*Last Updated: Reference guide creation*
