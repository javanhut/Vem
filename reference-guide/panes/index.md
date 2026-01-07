# panes Package Index

**Path:** `/home/javanhut/Development/Vem/internal/panes/`
**Purpose:** Window pane management with binary tree structure

## Files Overview

| File | Lines | Purpose |
|------|-------|---------|
| [manager.go](manager.md) | 307 | Pane tree management |

**Additional files** (not yet documented):
- `layout.go` - Layout calculations
- `navigation.go` - Focus navigation
- `pane.go` - Pane struct definition
- `resize.go` - Pane resize operations

## Key Types

### PaneManager (manager.go)
Top-level pane controller:
- Binary tree root
- Active pane tracking
- Zoom state
- Layout dimensions

### PaneNode (manager.go)
Binary tree node:
```go
type PaneNode struct {
    Pane  *Pane     // nil if internal node
    Split SplitType // SplitHorizontal or SplitVertical
    Left  *PaneNode // First child
    Right *PaneNode // Second child
    Ratio float32   // Split ratio (0.0-1.0)
}
```

### Pane
Individual pane instance:
- Unique ID
- Buffer index
- Active flag
- Viewport position

## Known Issues Summary

1. **manager.go:79** - Debug print statement should be removed

## Pane Operations

| Operation | Method | Description |
|-----------|--------|-------------|
| Split Horizontal | `SplitHorizontal()` | Creates left\|right split |
| Split Vertical | `SplitVertical()` | Creates top/bottom split |
| Close | `ClosePane()` | Removes active pane |
| Navigate | `NavigateLeft/Right/Up/Down()` | Focus adjacent pane |
| Cycle | `CycleNextPane()` | Focus next pane |
| Equalize | `Equalize()` | Set all ratios to 0.5 |
| Zoom | `ToggleZoom()` | Fullscreen active pane |
| Resize | `ResizeActivePane()` | Adjust split ratio |

## Tree Structure

```
Single Pane:
┌─────────────────────────────────────────────────────┐
│                     Root Node                       │
│                     (Leaf)                          │
│                     Pane A                          │
└─────────────────────────────────────────────────────┘

After Vertical Split (SplitHorizontal):
┌─────────────────────────────────────────────────────┐
│                Root Node (Internal)                 │
│              Split: Horizontal (│)                  │
│              Ratio: 0.5                             │
├─────────────────────┬───────────────────────────────┤
│  Left Node (Leaf)   │   Right Node (Leaf)           │
│     Pane A          │      Pane B                   │
└─────────────────────┴───────────────────────────────┘

After Additional Horizontal Split (SplitVertical):
┌─────────────────────────────────────────────────────┐
│                Root Node (Internal)                 │
│              Split: Horizontal (│)                  │
├─────────────────────┬───────────────────────────────┤
│  Left Node (Leaf)   │   Right Node (Internal)       │
│     Pane A          │   Split: Vertical (─)         │
│                     ├───────────────────────────────┤
│                     │  Top Node    │  Bottom Node   │
│                     │   Pane B     │    Pane C      │
└─────────────────────┴──────────────┴────────────────┘
```

## Navigation Algorithm

```
NavigateLeft:
1. Find current pane in tree
2. Walk up to find parent with horizontal split
3. If current is in right subtree, find rightmost leaf in left subtree
4. If no left sibling, navigation fails

NavigateDown:
1. Find current pane in tree
2. Walk up to find parent with vertical split
3. If current is in top subtree, find topmost leaf in bottom subtree
4. If no bottom sibling, navigation fails
```

---
*Last Updated: Reference guide creation*
