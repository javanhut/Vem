# help.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/help.go`
**Lines:** 289
**Purpose:** Generate help text from keybindings

## Overview

Generates formatted help text including:
- Global keybindings
- Mode-specific keybindings
- Command list
- Special sequences (gg, dd, etc.)

## Code Blocks

### Lines 1-8: Package and Imports

```go
package appcore
import (
    "fmt"
    "strings"

    "gioui.org/io/key"
)
```

### Lines 10-63: generateHelpText

Generates complete help text:
1. Header with title
2. Global keybindings section
3. Mode sections (NORMAL, INSERT, VISUAL, EXPLORER, TERMINAL)
4. Commands section
5. Special sequences section

### Lines 65-92: appendGlobalKeybindings

Lists global keybindings:

| Keys | Description |
|------|-------------|
| Ctrl+T | Toggle file explorer |
| Ctrl+H | Focus file explorer |
| Ctrl+L | Focus editor |
| Ctrl+F | Open fuzzy finder |
| Ctrl+U | Undo last edit |
| Ctrl+C | Copy current line |
| Ctrl+P | Paste from clipboard |
| Ctrl+Shift+R | Resize panes |
| Ctrl+X | Close pane/buffer |
| Ctrl+` | Open/toggle terminal |
| Alt+h/j/k/l | Focus pane left/down/up/right |
| Shift+Tab | Cycle to next pane |
| Shift+Enter | Toggle fullscreen |

### Lines 94-107: appendModeKeybindings

Iterates mode keybindings from `modeKeybindings` map and formats each.

### Lines 109-137: appendCommands

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

### Lines 139-163: appendSpecialSequences

Lists multi-key sequences:

| Sequence | Description |
|----------|-------------|
| gg | Jump to first line |
| G | Jump to last line |
| <count>G | Jump to line |
| <count>j/k | Move lines |
| dd | Delete line |
| zz | Center cursor |
| zt | Cursor to top |
| zb | Cursor to bottom |
| Ctrl+S v | Split vertically |
| Ctrl+S h | Split horizontally |
| Ctrl+S = | Equalize panes |
| Ctrl+S o | Zoom pane |

### Lines 165-215: formatKeybinding / formatKeyName

Formatting utilities:
- `formatKeybinding()` - Formats full keybinding with modifiers
- `formatKeyName()` - Converts key.Name to display string

### Lines 217-288: actionDescription

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

---
*Last Updated: Reference guide creation*
