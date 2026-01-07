# pane_rendering.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/pane_rendering.go`
**Lines:** 319
**Purpose:** Pane tree rendering and terminal display

## Overview

Implements pane layout rendering:
- Recursive pane tree traversal
- Per-pane viewport management
- Terminal content rendering
- Pane separator drawing

## Code Blocks

### Lines 1-17: Package and Imports

```go
package appcore
import (
    "image"
    "image/color"

    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/op/clip"
    "gioui.org/op/paint"
    "gioui.org/unit"
    "gioui.org/widget/material"

    "github.com/javanhut/vem/internal/editor"
    "github.com/javanhut/vem/internal/panes"
    "github.com/javanhut/vem/internal/terminal"
)
```

### Lines 19-43: drawPanes

Entry point for pane rendering:
1. Falls back to single buffer if no pane manager
2. Sets layout size on pane manager
3. Handles zoomed pane case
4. Calls `renderPaneNode()` for pane tree

### Lines 45-84: renderPaneNode

Recursively renders pane tree:

**Leaf Node (Lines 52-54):**
```go
if node.IsLeaf() {
    return s.drawSinglePane(gtx, node.Pane)
}
```

**Horizontal Split - Left|Right (Lines 57-69):**
```go
layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
    layout.Flexed(node.Ratio, ...),     // Left pane
    layout.Rigid(s.drawPaneSeparator),  // Vertical divider
    layout.Flexed(1-node.Ratio, ...),   // Right pane
)
```

**Vertical Split - Top/Bottom (Lines 70-83):**
```go
layout.Flex{Axis: layout.Vertical}.Layout(gtx,
    layout.Flexed(node.Ratio, ...),     // Top pane
    layout.Rigid(s.drawPaneSeparator),  // Horizontal divider
    layout.Flexed(1-node.Ratio, ...),   // Bottom pane
)
```

### Lines 86-158: drawSinglePane

Renders a single pane with its buffer:

**Inactive Pane Handling (Lines 119-141):**
1. Saves current viewport state
2. Temporarily switches active pane (quietly)
3. Draws buffer content
4. Restores original state

**Active Pane Handling (Lines 144-157):**
1. Syncs viewport from pane to global
2. Draws buffer
3. Syncs viewport back to pane

### Lines 160-176: drawPaneSeparator

Draws 1px separator line:
```go
if vertical {
    width = 1
    height = gtx.Constraints.Max.Y
} else {
    width = gtx.Constraints.Max.X
    height = 1
}
```

### Lines 178-204: drawTerminalPane

Renders terminal pane:
1. Gets terminal instance
2. Checks if running
3. Shows error if not initialized
4. Calls `drawTerminalContent()`

### Lines 206-318: drawTerminalContent

Renders terminal screen buffer:

**Character Measurement (Lines 211-224):**
- Uses 'M' character for width measurement
- Defaults to 8x16 if measurement fails

**Viewport Calculation (Lines 226-251):**
- Calculates lines per page
- Auto-scrolls to keep cursor visible
- Gets viewport top line

**Cell Rendering (Lines 260-311):**
For each visible cell:
1. Draws cell background
2. Draws cursor if at position
3. Draws character with color
4. Inverts colors for cursor

## Known Issues / Potential Bugs

1. **Lines 219-224: Hardcoded fallback dimensions**
   - Falls back to 8x16 if measurement fails
   - May not match actual font metrics

## Dead/Unused Code

None identified.

## Integration Points

- Called from app.go main layout
- Uses panes/manager.go for tree structure
- Interacts with terminal package for terminal panes

## Pane Tree Structure

```
                    Root (SplitHorizontal)
                   /                      \
          Left (ratio=0.5)         Right (SplitVertical)
              │                         /              \
         [Pane A]               Top (ratio=0.6)    Bottom
                                     │                │
                                [Pane B]          [Pane C]

Visual Layout:
┌──────────────────┬─────────────────────┐
│                  │                     │
│     Pane A       │      Pane B         │
│                  │                     │
│                  ├─────────────────────┤
│                  │      Pane C         │
│                  │                     │
└──────────────────┴─────────────────────┘
```

---
*Last Updated: Reference guide creation*
