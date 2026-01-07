# keybindings.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/keybindings.go`
**Lines:** ~850+
**Purpose:** Defines all keyboard bindings and action dispatch

## Overview

This file defines the keybinding system that maps keyboard inputs to actions across all editing modes. It contains action constants, binding definitions, and the action execution logic.

## Code Blocks

### Lines 1-8: Package and Imports
```go
package appcore
import (
    "strings"
    "unicode"
    "gioui.org/io/key"
)
```

### Lines 10-128: Action Constants
Defines all possible actions as an enum:

| Lines | Action Group | Examples |
|-------|--------------|----------|
| 12-20 | Global actions | `ActionToggleExplorer`, `ActionFocusEditor` |
| 22-28 | Mode transitions | `ActionEnterInsert`, `ActionExitMode` |
| 30-42 | Navigation | `ActionMoveLeft`, `ActionWordForward` |
| 44-51 | Editing | `ActionInsertNewline`, `ActionDeleteBackward`, `ActionUndo` |
| 53-60 | Clipboard | `ActionCopySelection`, `ActionPaste` |
| 62-70 | Explorer | `ActionOpenNode`, `ActionRenameFile` |
| 72-76 | Search | `ActionEnterSearch`, `ActionNextMatch` |
| 78-80 | Fuzzy finder | `ActionOpenFuzzyFinder` |
| 82-84 | Buffer mgmt | `ActionNextBuffer`, `ActionPrevBuffer` |
| 86-91 | Scrolling | `ActionScrollToCenter` |
| 93-103 | Pane mgmt | `ActionSplitVertical`, `ActionPaneClose` |
| 105-107 | Terminal | `ActionOpenTerminal` |
| 109-128 | LSP | `ActionLSPGotoDefinition`, `ActionLSPCompletion` |
| 130-134 | Buffer Completion | `ActionBufferCompletionTrigger`, `ActionBufferCompletionAccept`, etc. |
| 136-137 | Command Completion | `ActionCommandTabComplete` |

### Lines 140-145: KeyBinding Struct
```go
type KeyBinding struct {
    Modifiers key.Modifiers
    Key       key.Name
    Modes     []mode      // nil means all modes
    Action    Action
}
```

### Lines 137-161: Global Keybindings
Bindings that work in any mode:

| Line | Keys | Action |
|------|------|--------|
| 138 | Ctrl+T | Toggle explorer |
| 139 | Ctrl+H | Focus explorer |
| 140 | Ctrl+L | Focus editor |
| 141 | Ctrl+F | Open fuzzy finder |
| 142 | Ctrl+U | **Undo** (global) |
| 143-144 | Shift+Enter | Toggle fullscreen |
| 147 | Ctrl+C | Copy line (Normal) |
| 148 | Ctrl+P | Paste |
| 151-154 | Alt+H/J/K/L | Pane navigation |
| 157 | Ctrl+X | Close pane |
| 160 | Ctrl+` | Open terminal |

### Lines 163-279: Mode-Specific Keybindings

#### modeNormal (Lines 164-192)
| Line | Keys | Action |
|------|------|--------|
| 170 | i | Enter INSERT mode |
| 171 | v | Enter VISUAL (char) |
| 172 | Shift+V | Enter VISUAL (line) |
| 173 | d | Enter DELETE mode |
| 174-177 | h/j/k/l | Navigation |
| 178-180 | w/b/e | Word movement |
| 181-183 | 0/$ | Line start/end |
| 184 | / | Enter search |
| 185-186 | n/N | Next/prev match |
| 191 | Shift+K | LSP Hover |
| 200 | u | Undo (NEW) |

**NOTE:** `u` is now bound to undo in Normal mode (in addition to Ctrl+U).

#### modeInsert (Lines 202-219)
| Line | Keys | Action |
|------|------|--------|
| 203 | Escape | Exit to Normal (or cancel buffer completion) |
| 204-205 | Enter | Insert newline (or accept completion) |
| 207 | Tab | Insert tab (or cycle completion) |
| 208 | Shift+Tab | Prev completion (NEW) |
| 209 | Backspace | Delete backward |
| 216 | Ctrl+Space | LSP completion |
| 217 | Ctrl+N | Next completion / trigger buffer completion |
| 218 | Ctrl+P | Prev completion |

**Smart Completion Behavior:**
- Tab/Enter check for active LSP completion first, then buffer completion
- Escape cancels buffer completion before exiting INSERT mode
- Ctrl+N triggers buffer completion if LSP not active

#### modeVisual (Lines 210-231)
Same navigation as Normal, plus:
| Line | Keys | Action |
|------|------|--------|
| 226 | c | Copy selection |
| 227 | d | Delete selection |
| 228 | p | Paste |

#### modeDelete (Lines 232-235)
- Escape exits
- Shift+Tab cycles panes

#### modeCommand (Lines 249-256)
- Escape exits
- Enter executes command
- Backspace deletes char
- Tab/Shift+Tab: Path completion for `:cd`, `:e`, `:w` commands

#### modeExplorer (Lines 242-260)
| Line | Keys | Action |
|------|------|--------|
| 244-245 | Enter | Open node |
| 254 | r | Rename file |
| 255 | d | Delete file |
| 256 | n | Create file |
| 257 | u | Navigate up directory |

#### modeSearch (Lines 261-266)
- Escape exits
- Enter finds next match
- Backspace deletes char

#### modeFuzzyFinder (Lines 267-274)
- Escape exits
- Enter confirms selection
- Up/Down navigate results
- Backspace deletes char

#### modeTerminal (Lines 275-278)
- Escape or Shift+Tab exits terminal mode

### Lines 281-319: Keybinding Matching Functions

#### matchGlobalKeybinding (Lines 281-300)
Checks if key event matches any global binding.

#### matchModeKeybinding (Lines 302-319)
Checks if key event matches any binding for the current mode.

### Lines 321-374: Helper Functions

#### keysMatch (Lines 321-323)
Case-insensitive key name comparison.

#### modifiersMatch (Lines 325-366)
Modifier key matching with platform quirk handling.
- **IMPORTANT:** Lines 332-336 document that `ev.Modifiers` is always empty on some platforms!
- Uses tracked `s.ctrlPressed` and `s.shiftPressed` state instead

### Lines 376-791: Action Execution

`executeAction()` is a massive switch statement handling all actions:

| Lines | Action Category |
|-------|-----------------|
| 378-415 | Toggle/Focus explorer |
| 417-419 | Fullscreen |
| 420-464 | Mode transitions |
| 466-508 | Cursor movement |
| 510-543 | Insert operations |
| 561-590 | Delete operations & Undo |
| 592-605 | Clipboard operations |
| 607-659 | Explorer operations |
| 661-676 | Search operations |
| 677-681 | Fuzzy finder |
| 683-698 | Scrolling |
| 700-728 | Pane management |
| 730-734 | Terminal |
| 736-790 | LSP actions |

## Known Issues / Potential Bugs

1. **Line 257: `u` in Explorer = Navigate Up, not Undo**
   - Users may expect `u` to undo in Explorer mode
   - Undo only works via Ctrl+U globally

2. ~~**No `u` binding for undo in Normal mode**~~ FIXED
   - `u` now bound to ActionUndo in Normal mode (line 200)

## New Features Added

1. **Buffer Word Completion Actions** (Lines 130-134)
   - `ActionBufferCompletionTrigger`, `Accept`, `Cancel`, `Next`, `Prev`

2. **Smart Action Handlers** (Lines 813-838)
   - `ActionLSPReferencesNext/Prev/Open` also handle diagnostics list navigation
   - `ActionLSPDismissReferences` also closes diagnostics list
   - `ActionInsertTab` cycles buffer completion when active
   - `ActionInsertNewline` accepts buffer completion when active
   - `ActionExitMode` cancels buffer completion before exiting mode

3. **Tab Key Timing Fix**
   - Fixed `skipNextEdit` flag race condition for Tab key in INSERT mode
   - Flag is now set BEFORE action processing to prevent duplicate tab insertion

4. **Cursor Positioning Fix**
   - Fixed visual cursor offset caused by measuring prefix as single string vs rendering tokens separately
   - New `drawCursorAtX()` function uses pre-calculated X position from accumulated token widths

---
*Last Updated: After cursor positioning fix*
