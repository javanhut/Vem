# buffer_manager.go

**Path:** `/home/javanhut/Development/Vem/internal/editor/buffer_manager.go`
**Lines:** ~320+
**Purpose:** Multi-buffer management with file path tracking and large file handling

## Overview

Manages multiple open buffers with:
- Buffer collection and active buffer tracking
- File path to buffer index mapping
- Buffer creation, opening, closing
- Buffer navigation (next/prev)
- Large file handling (NEW)

## Code Blocks

### Lines 1-14: Package and Struct

```go
package editor

import (
    "fmt"
    "os"
    "path/filepath"
)

type BufferManager struct {
    buffers     []*Buffer          // All open buffers
    activeIndex int                // Currently active buffer
    pathToIndex map[string]int     // File path -> buffer index lookup
}
```

### Lines 16-39: Constructors

#### NewBufferManager (Lines 16-24)
Creates manager with default empty buffer.

#### NewBufferManagerWithBuffer (Lines 26-39)
Creates manager with initial buffer, registers file path.

### Lines 82-86: Large File Constants (NEW)

```go
const (
    LargeFileWarningBytes = 5 * 1024 * 1024  // 5MB - show warning
    LargeFileRejectBytes  = 50 * 1024 * 1024 // 50MB - reject file
)
```

### Lines 41-78: Basic Accessors

| Function | Lines | Returns |
|----------|-------|---------|
| `ActiveBuffer()` | 41-47 | Current active buffer |
| `BufferCount()` | 49-52 | Total buffer count |
| `ActiveIndex()` | 54-57 | Active buffer index |
| `GetBuffer(index)` | 59-65 | Buffer at index |
| `GetBufferByPath(path)` | 67-78 | Buffer by file path |

### Lines 88-137: File Opening (Updated)

#### OpenFile (Lines 88-137)
Opens file into new or existing buffer:
1. Converts to absolute path
2. **Checks if already open** - switches to existing buffer
3. If file doesn't exist - creates empty buffer
4. **NEW: Check file size**
   - > 50MB: Returns error "file too large"
   - > 5MB: Marks buffer with `SetLargeFile(true)`
5. Loads file content into new buffer
6. Registers path mapping

#### addBuffer (Lines 139-149)
Internal helper to add buffer and update mappings.

### Lines 129-154: Buffer Creation

| Function | Lines | Purpose |
|----------|-------|---------|
| `CreateEmptyBuffer()` | 129-134 | New empty text buffer |
| `CreateTerminalBuffer()` | 136-147 | New terminal buffer |
| `CreateBufferWithContent()` | 149-154 | Buffer with initial content |

### Lines 156-196: Save Operations

#### SaveActiveBuffer (Lines 156-168)
Saves current buffer to its file path.

#### SaveAs (Lines 170-196)
Saves buffer to new path:
1. Removes old path mapping
2. Saves to new path
3. Updates path mapping

### Lines 198-243: Buffer Closing

#### CloseBuffer (Lines 198-243)
Closes buffer at index:
1. Validates index
2. **Skips unsaved check for terminals**
3. Checks for unsaved changes (unless force)
4. Removes from path mapping
5. Updates mappings for shifted buffers
6. Creates default buffer if all closed
7. Adjusts active index

#### CloseActiveBuffer (Lines 245-248)
Wrapper for `CloseBuffer(activeIndex, force)`.

### Lines 250-280: Buffer Navigation

#### NextBuffer (Lines 250-258)
Cycles to next buffer (wraps around):
```go
bm.activeIndex = (bm.activeIndex + 1) % len(bm.buffers)
```

#### PrevBuffer (Lines 260-271)
Cycles to previous buffer (wraps around).

#### SwitchToBuffer (Lines 273-280)
Switches to buffer by index.

### Lines 282-308: Buffer Listing

#### ListBuffers (Lines 282-308)
Returns formatted buffer list for display:
```
* 1 + /path/to/file.go    (* = active, + = modified)
  2   [No Name]
  3   [Terminal]
```

## Known Issues / Potential Bugs

1. **Line 88-97: Duplicate buffer detection**
   - Loops through buffers to find match after `GetBufferByPath`
   - Could simplify by storing index in pathToIndex

2. **Line 221-228: Path mapping update on close**
   - Loops through all remaining buffers
   - O(n) operation on each close
   - Could optimize with different data structure

## Dead/Unused Code

None identified.

## Integration Points

- Used by `appState.bufferMgr` in app.go
- Pane system uses `BufferIndex` to reference buffers
- LSP integration uses file paths for document tracking

## Tab Bar Integration

For adding a tab bar, use these methods:
- `BufferCount()` - Get number of tabs
- `ActiveIndex()` - Highlight active tab
- `ListBuffers()` - Get tab labels
- `SwitchToBuffer(index)` - Handle tab clicks
- `CloseBuffer(index, force)` - Handle tab close

## New Features Added

1. **Large File Handling** (Lines 82-86, 119-133)
   - Constants for warning (5MB) and reject (50MB) thresholds
   - OpenFile checks file size before loading
   - Marks large buffers with `SetLargeFile(true)` for status bar display

---
*Last Updated: After enhancement plan implementation*
