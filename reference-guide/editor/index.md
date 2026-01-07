# editor Package Index

**Path:** `/home/javanhut/Development/Vem/internal/editor/`
**Purpose:** Text buffer and editing operations

## Files Overview

| File | Lines | Purpose |
|------|-------|---------|
| [buffer.go](buffer.md) | 888 | Core text buffer with editing operations |
| [buffer_manager.go](buffer_manager.md) | 309 | Multi-buffer management |

## Key Types

### Buffer (buffer.go)
Core text editing type containing:
- Line storage
- Cursor position
- Selection state
- Undo/redo history
- Modified flag
- LSP callback

### BufferManager (buffer_manager.go)
Multi-buffer container with:
- Buffer list
- Active buffer tracking
- Path-to-index mapping
- Buffer creation/deletion

### Cursor (buffer.go)
```go
type Cursor struct {
    Line int  // 0-based line number
    Col  int  // 0-based column (rune index)
}
```

## Known Issues Summary

1. **buffer.go:650** - `os.ReadFile()` loads entire file into memory
2. **buffer.go:830** - Undo limited to 100 entries

## Key Operations

### Navigation
- `MoveLeft()`, `MoveRight()`, `MoveUp()`, `MoveDown()`
- `MoveWordForward()`, `MoveWordBackward()`, `MoveWordEnd()`
- `MoveToLine()`, `MoveToLineStart()`, `MoveToLineEnd()`

### Editing
- `InsertRune()`, `InsertText()`, `InsertNewline()`
- `DeleteBackward()`, `DeleteForward()`, `DeleteLine()`
- `DeleteCharRange()` - Range deletion

### Undo/Redo
- `Undo()`, `Redo()`
- `pushUndo()` - Internal undo stack management

### File Operations
- `LoadFile()`, `Save()`, `SaveAs()`
- `GetContent()`, `Line()`, `LineCount()`

## Buffer Lifecycle

```
┌─────────────────────┐
│  BufferManager      │
│  CreateEmptyBuffer()│
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Buffer Created    │
│   (empty, unnamed)  │
└──────────┬──────────┘
           │ LoadFile() or manual editing
           ▼
┌─────────────────────┐
│   Buffer Modified   │
│   modified = true   │
└──────────┬──────────┘
           │ Save()
           ▼
┌─────────────────────┐
│   Buffer Saved      │
│   modified = false  │
└──────────┬──────────┘
           │ CloseBuffer()
           ▼
┌─────────────────────┐
│   Buffer Closed     │
│   (removed from mgr)│
└─────────────────────┘
```

---
*Last Updated: Reference guide creation*
