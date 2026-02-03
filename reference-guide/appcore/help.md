# help.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/help.go`
**Lines:** ~305
**Purpose:** Generate help text from keybindings with OS-specific modifier display

## Overview

Generates formatted help text including:
- Global keybindings
- Mode-specific keybindings
- Command list
- Special sequences (gg, dd, etc.)

**OS-Specific Display:** On macOS, shows "Cmd" and "Option" instead of "Ctrl" and "Alt".

## Code Blocks

### Lines 1-24: Package, Imports, and OS Detection

```go
package appcore
import (
    "fmt"
    "runtime"
    "strings"

    "gioui.org/io/key"
)

// OS-specific modifier key display names
var (
    modCtrlDisplay string
    modAltDisplay  string
)

func init() {
    if runtime.GOOS == "darwin" {
        modCtrlDisplay = "Cmd"
        modAltDisplay = "Option"
    } else {
        modCtrlDisplay = "Ctrl"
        modAltDisplay = "Alt"
    }
}
```

### Lines 26-78: generateHelpText

Generates complete help text:
1. Header with title
2. Global keybindings section
3. Mode sections (NORMAL, INSERT, VISUAL, EXPLORER, TERMINAL)
4. Commands section
5. Special sequences section

### Lines 80-107: appendGlobalKeybindings

Lists global keybindings using `modCtrlDisplay` and `modAltDisplay` for OS-appropriate display:

| Keys (Linux/Windows) | Keys (macOS) | Description |
|----------------------|--------------|-------------|
| Ctrl+T | Cmd+T | Toggle file explorer |
| Ctrl+H | Cmd+H | Focus file explorer |
| Ctrl+L | Cmd+L | Focus editor |
| Ctrl+F | Cmd+F | Open fuzzy finder |
| Ctrl+U | Cmd+U | Undo last edit |
| Ctrl+C | Cmd+C | Copy current line |
| Ctrl+P | Cmd+P | Paste from clipboard |
| Ctrl+Shift+R | Cmd+Shift+R | Resize panes |
| Ctrl+X | Cmd+X | Close pane/buffer |
| Ctrl+` | Cmd+` | Open/toggle terminal |
| Alt+h/j/k/l | Option+h/j/k/l | Focus pane left/down/up/right |
| Shift+Tab | Shift+Tab | Cycle to next pane |
| Shift+Enter | Shift+Enter | Toggle fullscreen |

### Lines 109-122: appendModeKeybindings

Iterates mode keybindings from `modeKeybindings` map and formats each.

### Lines 124-152: appendCommands

Lists command-mode commands:

| Command | Description |
|---------|-------------|
| :q | Close current pane/buffer |
| :q! | Force close |
| :qa | Quit entire application |
| :qa! | Force quit all |
| :w | Save current buffer |
| :w <file> | Save as |
| :wq | Save and close |
| :e <file> | Open file |
| :bn | Next buffer |
| :bp | Previous buffer |
| :bd | Delete buffer |
| :ls | List buffers |
| :ex | Toggle explorer |
| :cd <path> | Change directory |
| :pwd | Print working directory |
| :term | Open terminal |
| :help | Show help |

### Lines 154-178: appendSpecialSequences

Lists multi-key sequences (uses `modCtrlDisplay` for split commands):

| Sequence (Linux/Windows) | Sequence (macOS) | Description |
|--------------------------|------------------|-------------|
| gg | gg | Jump to first line |
| G | G | Jump to last line |
| <count>G | <count>G | Jump to line |
| <count>j/k | <count>j/k | Move lines |
| dd | dd | Delete line |
| zz | zz | Center cursor |
| zt | zt | Cursor to top |
| zb | zb | Cursor to bottom |
| Ctrl+S v | Cmd+S v | Split vertically |
| Ctrl+S h | Cmd+S h | Split horizontally |
| Ctrl+S = | Cmd+S = | Equalize panes |
| Ctrl+S o | Cmd+S o | Zoom pane |

### Lines 180-202: formatKeybinding

Formats a KeyBinding for display using OS-specific modifier names:

```go
if binding.Modifiers.Contain(key.ModCtrl) {
    parts = append(parts, modCtrlDisplay)  // "Cmd" on macOS, "Ctrl" elsewhere
}
if binding.Modifiers.Contain(key.ModAlt) {
    parts = append(parts, modAltDisplay)   // "Option" on macOS, "Alt" elsewhere
}
```

### Lines 204-230: formatKeyName

Converts key.Name to display string:
- `key.NameEscape` → "Esc"
- `key.NameReturn` → "Enter"
- Arrow keys → "←", "→", "↑", "↓"
- etc.

### Lines 232-303: actionDescription

Returns human-readable descriptions for all actions:
- Navigation actions (move left/right/up/down)
- Mode transitions (enter insert/visual/delete)
- Editing actions (undo, copy, paste, delete)
- Explorer actions (open, collapse, expand)
- Pane actions (split, focus, close)
- And more...

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Called when `:help` command is executed
- Help text displayed in new buffer
- Uses keybindings.go definitions
- Uses `runtime.GOOS` for OS detection

---
*Last Updated: Added OS-specific modifier display (Cmd/Option for macOS)*
