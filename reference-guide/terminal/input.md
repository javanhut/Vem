# input.go

**Path:** `/home/javanhut/Development/Vem/internal/terminal/input.go`
**Lines:** 138
**Purpose:** Convert Gio key events to VT100 terminal escape sequences

## Overview

Translates Gio keyboard events to:
- VT100/ANSI escape sequences for special keys
- Control character sequences for Ctrl combinations
- Alt prefix sequences (ESC + character)

## Code Blocks

### Lines 1-6: Package and Imports

```go
package terminal
import (
    "gioui.org/io/key"
)
```

### Lines 7-37: Special Keys

#### KeyToTerminalSequence (Lines 7-37)
Converts special keys to escape sequences:

| Key | Sequence |
|-----|----------|
| Return/Enter | `\r` |
| Tab | `\t` |
| Backspace | `\x7f` (DEL, 127) |
| Delete | `\x1b[3~` |
| Up Arrow | `\x1b[A` |
| Down Arrow | `\x1b[B` |
| Right Arrow | `\x1b[C` |
| Left Arrow | `\x1b[D` |
| Home | `\x1b[H` |
| End | `\x1b[F` |
| PageUp | `\x1b[5~` |
| PageDown | `\x1b[6~` |
| Escape | `\x1b` |

### Lines 39-65: Function Keys

| Key | Sequence |
|-----|----------|
| F1 | `\x1bOP` |
| F2 | `\x1bOQ` |
| F3 | `\x1bOR` |
| F4 | `\x1bOS` |
| F5 | `\x1b[15~` |
| F6 | `\x1b[17~` |
| F7 | `\x1b[18~` |
| F8 | `\x1b[19~` |
| F9 | `\x1b[20~` |
| F10 | `\x1b[21~` |
| F11 | `\x1b[23~` |
| F12 | `\x1b[24~` |

### Lines 67-79: Shift+Arrow Keys

Uses CSI modifier format `\x1b[1;2X`:

| Key | Sequence |
|-----|----------|
| Shift+Up | `\x1b[1;2A` |
| Shift+Down | `\x1b[1;2B` |
| Shift+Right | `\x1b[1;2C` |
| Shift+Left | `\x1b[1;2D` |

### Lines 81-125: Ctrl Combinations

#### Ctrl+Arrow (Lines 81-93)
Uses modifier 5: `\x1b[1;5X`

#### Ctrl+Letter (Lines 94-104)
Converts to control character:
```go
// Ctrl+A = 0x01, Ctrl+B = 0x02, ..., Ctrl+Z = 0x1A
if r >= 'a' && r <= 'z' {
    return string(rune(r - 'a' + 1))
}
```

#### Special Ctrl Combinations (Lines 106-125)

| Combination | Sequence | Name |
|-------------|----------|------|
| Ctrl+Space | `\x00` | NUL |
| Ctrl+2 / Ctrl+@ | `\x00` | NUL |
| Ctrl+3 / Ctrl+[ | `\x1b` | ESC |
| Ctrl+4 / Ctrl+\ | `\x1c` | FS |
| Ctrl+5 / Ctrl+] | `\x1d` | GS |
| Ctrl+6 / Ctrl+^ | `\x1e` | RS |
| Ctrl+7 / Ctrl+_ | `\x1f` | US |
| Ctrl+8 | `\x7f` | DEL |

### Lines 127-137: Alt Combinations

#### Alt+Key (Lines 127-133)
Sends ESC prefix:
```go
if ev.Modifiers.Contain(key.ModAlt) {
    if len(ev.Name) == 1 {
        return "\x1b" + string(ev.Name)
    }
}
```

**Note:** Regular character input comes through EditEvent, not KeyEvent.

## Known Issues / Potential Bugs

1. **Line 95-104: Only handles single-character key names**
   - Keys with multi-character names won't generate Ctrl sequences
   - This is intentional - special keys are handled separately

## Dead/Unused Code

None identified.

## Integration Points

- Called from `handleTerminalKey()` in app.go
- Output sent to terminal via `Terminal.Write()`
- Works with EditEvent for regular text input

## Escape Sequence Format

```
Standard VT100:
  \x1b[   = CSI (Control Sequence Introducer)
  \x1bO   = SS3 (Single Shift 3)

CSI Format:
  \x1b[<params><final>

  Final bytes:
    A = Up, B = Down, C = Right, D = Left
    H = Home, F = End

  Modifier params:
    1;2 = Shift
    1;5 = Ctrl
    1;3 = Alt
```

---
*Last Updated: Reference guide creation*
