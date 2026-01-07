# buffer.go

**Path:** `/home/javanhut/Development/Vem/internal/editor/buffer.go`
**Lines:** ~900+
**Purpose:** Text buffer implementation with cursor, undo, file operations, and word completion

## Overview

Implements the core text buffer abstraction with:
- Line-based text storage
- Cursor position tracking
- Undo/redo stack
- File I/O operations
- Word navigation (Vim-style w/b/e)
- LSP change notification callbacks

## Code Blocks

### Lines 1-7: Package and Imports
```go
package editor
import (
    "os"
    "strings"
    "unicode/utf8"
)
```

### Lines 9-38: Type Definitions

#### BufferType (Lines 9-15)
```go
type BufferType int
const (
    BufferTypeText BufferType = iota
    BufferTypeTerminal
)
```

#### UndoEntry (Lines 17-22)
```go
type UndoEntry struct {
    lines       []string    // Snapshot of buffer lines
    cursor      Cursor      // Cursor position at snapshot
    description string      // Operation description
}
```

#### Buffer (Lines 24-39)
```go
type Buffer struct {
    lines       []string        // Text content
    cursor      Cursor          // Current position
    filePath    string          // Associated file
    modified    bool            // Unsaved changes flag
    undoStack   []UndoEntry     // Undo history (max 100)
    maxUndos    int             // Default: 100
    bufferType  BufferType      // Text or Terminal
    terminal    interface{}     // Terminal reference
    readOnly    bool            // Prevent edits
    isLargeFile bool            // File > 5MB (NEW)
    lspOnChange func(string)    // LSP callback
}
```

#### Cursor (Lines 40-44)
```go
type Cursor struct {
    Line int  // 0-based line index
    Col  int  // 0-based column (rune index, not byte)
}
```

### Lines 46-58: Constructor

`NewBuffer(text string)`:
- Splits text by newlines
- Ensures at least one empty line
- Initializes empty undo stack (max 100 entries)

### Lines 60-106: Basic Accessors

| Function | Lines | Returns |
|----------|-------|---------|
| `LineCount()` | 60-63 | Number of lines |
| `Line(i)` | 65-71 | Line at index |
| `LinesRange(start, end)` | 73-93 | Copy of line range |
| `LinePrefix(lineIdx, cols)` | 95-106 | First N runes of line |
| `Cursor()` | 108-111 | Current cursor position |

### Lines 113-161: Line Operations

#### MoveToLine (Lines 115-126)
Moves cursor to line, clamping to valid range.

#### SetCursor (Lines 128-141)
Moves cursor to specific line and column:
1. Clamps line to valid range
2. Sets cursor line and column
3. Clamps column to line length

#### ReplaceLine (Lines 143-154)
Replaces content of a specific line:
1. Checks bounds and read-only flag
2. Saves undo state
3. Replaces line content
4. Marks modified

#### DeleteLines (Lines 156-189)
Deletes inclusive line range:
1. Checks read-only flag
2. Saves state for undo
3. Removes lines
4. Ensures buffer has at least one line
5. Repositions cursor
6. Marks modified

### Lines 191-253: Text Insertion

#### InsertLines (Lines 163-190)
Inserts lines at index, saves undo state.

#### InsertText (Lines 192-225)
Inserts text at cursor:
1. Checks read-only flag
2. Saves undo state
3. Splits current line at cursor
4. Handles multi-line insertion
5. Updates cursor position
6. Triggers LSP callback

### Lines 227-286: Deletion Operations

#### DeleteBackward (Lines 227-260)
Backspace semantics:
- At start of line: merge with previous
- Otherwise: delete character before cursor

#### DeleteForward (Lines 262-286)
Delete key semantics:
- At end of line: merge with next
- Otherwise: delete character at cursor

### Lines 288-354: Cursor Movement

| Function | Lines | Vim Key | Behavior |
|----------|-------|---------|----------|
| `MoveLeft()` | 288-300 | `h` | Left, wrap to prev line |
| `MoveRight()` | 302-315 | `l` | Right, wrap to next line |
| `MoveUp()` | 317-325 | `k` | Up, clamp column |
| `MoveDown()` | 327-335 | `j` | Down, clamp column |
| `JumpLineStart()` | 337-344 | `0` | Start of line |
| `JumpLineEnd()` | 346-354 | `$` | End of line |

### Lines 356-508: Word Navigation

#### MoveWordForward (Lines 356-400)
Vim's `w` command:
1. Skip current word characters
2. Skip whitespace (across lines)
3. Stop at start of next word

#### MoveWordBackward (Lines 402-456)
Vim's `b` command:
1. Move back one position
2. Skip whitespace backward
3. Find start of word

#### MoveWordEnd (Lines 458-508)
Vim's `e` command:
1. Move forward one position
2. Skip whitespace
3. Find end of current/next word

### Lines 510-596: Helper Functions

#### Character Type Classification (Lines 510-542)
```go
type charType int  // Space, Word, Punct

func isSpace(r rune) bool      // Whitespace check
func isWordChar(r rune) bool   // [a-zA-Z0-9_]
func getCharType(r rune) charType
```

#### String Utilities (Lines 544-596)
| Function | Purpose |
|----------|---------|
| `clampColumn()` | Ensure cursor within line |
| `lineLength()` | Rune count of line |
| `splitAtRune()` | Split string at rune index |
| `runeCount()` | UTF-8 safe length |
| `byteIndexForRune()` | Byte offset for rune index |
| `removeLine()` | Remove line from slice |

### Lines 598-646: File Path and Modified State

| Function | Lines | Purpose |
|----------|-------|---------|
| `FilePath()` | 598-601 | Get file path |
| `SetFilePath()` | 603-606 | Set file path |
| `Modified()` | 608-611 | Check modified flag |
| `SetModified()` | 613-616 | Set modified flag |
| `markModified()` | 618-625 | Internal: set + notify LSP |
| `SetLSPChangeCallback()` | 627-631 | Set LSP callback |
| `ClearLSPChangeCallback()` | 633-636 | Remove LSP callback |
| `SetReadOnly()` | 638-641 | Set read-only flag |
| `IsReadOnly()` | 643-646 | Check read-only flag |

### Lines 648-718: File I/O

#### LoadFromFile (Lines 648-673)
```go
func (b *Buffer) LoadFromFile(path string) error
```
**POTENTIAL ISSUE:** Uses `os.ReadFile()` (line 650) which loads entire file into memory. No chunked loading for large files.

#### SaveToFile (Lines 675-691)
Saves buffer with trailing newline, sets permissions 0644.

#### NewBufferFromFile (Lines 706-718)
Factory function that creates buffer and loads file.

### Lines 720-827: Character Range Operations

#### GetCharRange (Lines 720-767)
Returns text in character range (for visual selection).

#### DeleteCharRange (Lines 769-827)
Deletes text in character range:
1. Saves undo state
2. Handles single-line deletion
3. Handles multi-line deletion (merges lines)
4. Repositions cursor

### Lines 830-866: Undo System

#### saveState (Lines 830-848)
```go
func (b *Buffer) saveState(description string)
```
- Creates deep copy of lines
- Saves cursor position
- Limits stack to `maxUndos` (100) entries

#### Undo (Lines 850-866)
```go
func (b *Buffer) Undo() bool
```
- Pops last entry from undo stack
- Restores lines and cursor
- Marks buffer as modified

### Lines 868-888: Terminal Support

| Function | Lines | Purpose |
|----------|-------|---------|
| `BufferType()` | 868-871 | Get buffer type |
| `Terminal()` | 873-876 | Get terminal interface |
| `SetTerminal()` | 878-882 | Associate terminal |
| `IsTerminal()` | 884-887 | Check if terminal buffer |

### Lines 649-657: Large File Methods (NEW)

| Function | Lines | Purpose |
|----------|-------|---------|
| `SetLargeFile(bool)` | 649-652 | Mark buffer as containing large file |
| `IsLargeFile()` | 654-657 | Check if buffer contains large file (>5MB) |

### Lines 760-820: Word Completion Methods (NEW)

#### GetCurrentWordPrefix (Lines 760-785)
Returns the word prefix at the cursor position:
- Scans backward from cursor to find word start
- Returns characters from word start to cursor
- Used by buffer word completion (Ctrl+N)

#### GetWordsMatching (Lines 787-820)
Returns unique words matching a prefix:
- Scans all lines for words starting with prefix
- Deduplicates results
- Sorts by frequency (most common first)
- Used by buffer word completion menu

#### Helper Functions (Lines 822-840)
| Function | Purpose |
|----------|---------|
| `isWordRune(r)` | Check if rune is [a-zA-Z0-9_] |
| `extractWords(line)` | Extract all words from a line |

## Known Issues / Potential Bugs

1. ~~**Line 650: Large file loading**~~ FIXED
   - Now handled in buffer_manager.go
   - Files > 5MB show warning, > 50MB rejected

2. **Line 856-858: Undo marks buffer as modified**
   - `markModified()` is called after undo
   - This triggers LSP callback even though content reverted
   - Could cause unnecessary LSP traffic

3. **No redo functionality**
   - Only undo stack, no redo
   - Could add redo stack for full undo/redo support

## New Features Added

1. **Word Completion Methods** - `GetCurrentWordPrefix()` and `GetWordsMatching(prefix)`
2. **Large File Flag** - `isLargeFile` field with getter/setter

---
*Last Updated: After enhancement plan implementation*
