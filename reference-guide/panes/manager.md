# manager.go

**Path:** `/home/javanhut/Development/Vem/internal/panes/manager.go`
**Lines:** 307
**Purpose:** Pane tree management with splitting, navigation, and zoom

## Overview

Implements window pane management similar to tmux/vim:
- Binary tree structure for pane layout
- Horizontal and vertical splitting
- Clockwise pane cycling
- Zoom (maximize) mode

## Code Blocks

### Lines 1-17: Package and Struct

```go
package panes

import (
    "fmt"
    "math"
    "sort"
)

type PaneManager struct {
    root       *PaneNode    // Root of pane tree
    activePane *Pane        // Currently focused pane
    nextPaneID int          // Auto-incrementing pane IDs
    zoomed     *Pane        // Maximized pane (nil if not zoomed)
    lastWidth  int          // Last layout width
    lastHeight int          // Last layout height
}
```

### Lines 19-30: Constructor

```go
func NewPaneManager(initialBufferIndex int) *PaneManager
```
Creates manager with single pane displaying given buffer.

### Lines 32-62: Basic Accessors

| Function | Lines | Returns |
|----------|-------|---------|
| `Root()` | 32-35 | Pane tree root |
| `SetLayoutSize()` | 37-41 | Store layout dimensions |
| `ActivePane()` | 43-46 | Currently active pane |
| `PaneCount()` | 48-54 | Total pane count |
| `AllPanes()` | 56-62 | All panes in tree |

### Lines 64-86: Active Pane Management

#### SetActivePane (Lines 64-80)
Sets pane as active:
1. Deactivates all panes
2. Activates target pane
3. Updates `activePane` reference

**NOTE:** Line 79 contains debug print statement:
```go
fmt.Printf("[PANE_MANAGER] SetActivePane: ID=%s, BufferIndex=%d\n", ...)
```

#### SetActivePaneQuiet (Lines 82-86)
Low-level setter without deactivation (for rendering).

### Lines 88-148: Pane Splitting

#### SplitVertical (Lines 88-106)
Creates horizontal divider, new pane below:
```
┌───────┐      ┌───────┐
│   A   │  =>  │   A   │
│       │      ├───────┤
│       │      │  NEW  │
└───────┘      └───────┘
```

#### SplitHorizontal (Lines 108-126)
Creates vertical divider, new pane right:
```
┌───────┐      ┌───┬───┐
│   A   │  =>  │ A │NEW│
│       │      │   │   │
└───────┘      └───┴───┘
```

#### splitNodeContainingPane (Lines 128-148)
Recursive helper to find and split target pane node.

### Lines 150-202: Pane Closing

#### ClosePane (Lines 150-174)
Closes active pane:
1. Prevents closing last pane
2. Removes from tree
3. Sets new active pane

#### removeNodeContainingPane (Lines 176-202)
Recursive helper to remove pane and collapse parent split.

### Lines 204-275: Pane Navigation

#### CycleNextPane (Lines 204-223)
Cycles to next pane in clockwise order.

#### clockwiseOrder (Lines 225-275)
Calculates clockwise ordering of panes:
1. Gets all pane geometries
2. Computes center point
3. Calculates angle from center for each pane
4. Sorts by angle (0 at top, clockwise)

### Lines 277-306: Zoom and Lookup

#### ToggleZoom (Lines 277-286)
Toggles zoom (maximize) for active pane.

#### IsZoomed / ZoomedPane (Lines 288-296)
Check and get zoomed pane.

#### FindPaneByBufferIndex (Lines 298-306)
Find pane displaying given buffer.

## Known Issues / Potential Bugs

1. **Line 79: Debug print statement in production**
   ```go
   fmt.Printf("[PANE_MANAGER] SetActivePane: ID=%s, BufferIndex=%d\n", ...)
   ```
   - Should be removed or use proper logging
   - Pollutes stdout in production

2. **Line 227: Empty pane list handling**
   - Returns allPanes directly if <= 1
   - Could return empty slice instead

## Dead/Unused Code

1. **Lines 38-41: lastWidth/lastHeight**
   - Only used for clockwise ordering
   - Could be passed as parameters instead

## Integration Points

- Used by `appState.paneManager` in app.go
- Rendering in pane_rendering.go
- Keybindings in keybindings.go (Ctrl+S prefix commands)

## Pane Tree Structure

```
PaneNode (Split Horizontal)
├── Left: PaneNode (Leaf) -> Pane A
└── Right: PaneNode (Split Vertical)
    ├── Left: PaneNode (Leaf) -> Pane B
    └── Right: PaneNode (Leaf) -> Pane C

Layout:
┌─────┬─────┐
│     │  B  │
│  A  ├─────┤
│     │  C  │
└─────┴─────┘
```

---
*Last Updated: Reference guide creation*
