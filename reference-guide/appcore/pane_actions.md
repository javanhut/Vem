# pane_actions.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/pane_actions.go`
**Lines:** 375
**Purpose:** Pane management action handlers

## Overview

Implements all pane-related actions:
- Split creation (vertical/horizontal)
- Focus navigation (left/right/up/down)
- Pane cycling and closing
- Resize mode handling
- Zoom toggle

## Code Blocks

### Lines 1-12: Package and Imports

```go
package appcore
import (...)

const paneResizeStep = 0.05  // 5% resize per step
```

### Lines 14-49: handleSplitVertical

Creates vertical split (left | right):
1. Creates new empty buffer
2. Calls `paneManager.SplitHorizontal()` (creates vertical divider)
3. Updates status with pane count

**Note:** Contains debug print statements (potential cleanup).

### Lines 51-86: handleSplitHorizontal

Creates horizontal split (top / bottom):
1. Creates new empty buffer
2. Calls `paneManager.SplitVertical()` (creates horizontal divider)
3. Updates status with pane count

**Note:** Contains debug print statements.

### Lines 88-138: Focus Navigation

| Function | Lines | Purpose |
|----------|-------|---------|
| `handlePaneFocusLeft()` | 88-99 | Focus pane to left |
| `handlePaneFocusRight()` | 101-112 | Focus pane to right |
| `handlePaneFocusUp()` | 114-125 | Focus pane above |
| `handlePaneFocusDown()` | 127-138 | Focus pane below |

Each calls corresponding `paneManager.Navigate*()` method.

### Lines 140-177: handlePaneCycleNext

Cycles to next pane:
1. Checks if more than one pane exists
2. Calls `paneManager.CycleNextPane()`
3. Updates status with pane index

**Note:** Contains debug print statements.

### Lines 179-242: handlePaneClose

Closes active pane:
1. Checks for unsaved changes
2. Closes terminal if present
3. If multiple panes: close pane and buffer
4. If last pane: close buffer, switch to buffer 0

### Lines 244-267: Equalize and Zoom

| Function | Lines | Purpose |
|----------|-------|---------|
| `handlePaneEqualize()` | 244-252 | Make all panes equal size |
| `handlePaneZoomToggle()` | 254-267 | Toggle pane zoom |

### Lines 269-291: handlePaneCommand

Handles Ctrl+S prefix commands:

| Key | Action |
|-----|--------|
| v | Split vertical |
| h | Split horizontal |
| = | Equalize panes |
| o | Toggle zoom |

### Lines 293-312: togglePaneResizeMode

Enters/exits pane resize mode:
1. Checks if multiple panes exist
2. Checks if not zoomed
3. Sets `paneResizeMode = true`
4. Shows resize instructions

### Lines 314-357: handlePaneResizeKey

Handles keys in resize mode:

| Key | Action |
|-----|--------|
| Escape | Exit resize mode |
| ←/h | Resize left |
| →/l | Resize right |
| ↑/j | Resize up (note: j is up) |
| ↓/k | Resize down (note: k is down) |

### Lines 359-374: resizeActivePane

Applies resize step:
1. Checks not zoomed
2. Calls `paneManager.ResizeActivePane(dir, paneResizeStep)`

## Known Issues / Potential Bugs

1. **Lines 16-48, 53-84, 142-176: Debug print statements**
   - `fmt.Printf("[PANE_SPLIT]...")` statements
   - Should be removed or use proper logging

2. **Lines 349-350: Inverted j/k for vertical resize**
   - j moves up, k moves down (opposite of Vim convention)
   - May confuse users expecting Vim-like behavior

## Dead/Unused Code

None identified.

## Integration Points

- Called from keybindings.go action handlers
- Uses panes/manager.go for operations
- Interacts with bufferMgr for buffer creation

## Pane Operations Summary

```
Ctrl+S v    → Split vertical (left|right)
Ctrl+S h    → Split horizontal (top/bottom)
Ctrl+S =    → Equalize all panes
Ctrl+S o    → Toggle zoom

Alt+h/j/k/l → Focus left/down/up/right
Shift+Tab   → Cycle to next pane
Ctrl+X      → Close pane

Ctrl+Shift+R → Enter resize mode
  h/←       → Shrink left
  l/→       → Grow right
  j/↑       → Shrink up
  k/↓       → Grow down
  Esc       → Exit resize mode
```

---
*Last Updated: Reference guide creation*
