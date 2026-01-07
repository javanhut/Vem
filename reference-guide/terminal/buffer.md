# buffer.go

**Path:** `/home/javanhut/Development/Vem/internal/terminal/buffer.go`
**Lines:** 237
**Purpose:** Terminal screen buffer with cell-based storage

## Overview

Implements terminal screen buffer with:
- Cell-based storage (character + attributes)
- Thread-safe access via RWMutex
- Dirty line tracking for efficient rendering
- Resize support with content preservation

## Code Blocks

### Lines 1-6: Package and Imports

```go
package terminal
import (
    "image/color"
    "sync"
)
```

### Lines 8-19: Cell Struct

```go
type Cell struct {
    Rune      rune        // Unicode character
    FG        color.NRGBA // Foreground color
    BG        color.NRGBA // Background color
    Bold      bool        // Bold attribute
    Dim       bool        // Dim attribute
    Italic    bool        // Italic attribute
    Underline bool        // Underline attribute
    Blink     bool        // Blink attribute
    Reverse   bool        // Reverse video
}
```

### Lines 21-45: Line and ScreenBuffer

#### Line (Lines 21-25)
```go
type Line struct {
    Cells []Cell
    Dirty bool  // Whether line needs redraw
}
```

#### ScreenBuffer (Lines 27-36)
```go
type ScreenBuffer struct {
    lines       []Line
    width       int
    height      int
    cursorX     int
    cursorY     int
    cursorStyle CursorStyle
    mu          sync.RWMutex
}
```

#### CursorStyle (Lines 38-45)
```go
const (
    CursorBlock CursorStyle = iota
    CursorUnderline
    CursorBar
)
```

### Lines 47-70: Constructor

#### NewScreenBuffer (Lines 47-70)
Creates new screen buffer:
1. Allocates lines array
2. Initializes each line with cells
3. Sets all lines dirty
4. Fills cells with space and default colors

### Lines 72-100: Getters

| Function | Lines | Purpose |
|----------|-------|---------|
| `Dimensions()` | 72-77 | Get width and height |
| `GetLine()` | 79-93 | Get copy of line (thread-safe) |
| `GetCursor()` | 95-100 | Get cursor position and style |

**Note:** `GetLine()` returns a copy to prevent data races.

### Lines 102-135: Setters

| Function | Lines | Purpose |
|----------|-------|---------|
| `SetCursor()` | 102-122 | Set cursor with bounds clamping |
| `SetCell()` | 124-135 | Set cell value and mark dirty |

### Lines 137-181: Clear Operations

| Function | Lines | Purpose |
|----------|-------|---------|
| `ClearLine()` | 137-154 | Clear single line to spaces |
| `Clear()` | 156-171 | Clear entire buffer |
| `MarkClean()` | 173-181 | Mark all lines as clean |

### Lines 183-229: Resize

#### Resize (Lines 183-229)
Resizes buffer preserving content:
1. Early return if dimensions unchanged
2. Creates new lines array
3. Copies old content where possible
4. Fills new cells with defaults
5. Clamps cursor to new bounds

### Lines 231-236: Writer Interface

#### Write (Lines 231-236)
Implements `io.Writer` for vt10x:
```go
func (sb *ScreenBuffer) Write(p []byte) (n int, err error) {
    return len(p), nil  // vt10x handles parsing
}
```

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Created by `Terminal.NewTerminal()` in terminal.go
- Updated by vt10x emulator during readLoop
- Read by rendering code in app.go

## Thread Safety

All public methods use `sync.RWMutex`:
- Read methods use `RLock()`
- Write methods use `Lock()`
- `GetLine()` returns copy to prevent races

---
*Last Updated: Reference guide creation*
