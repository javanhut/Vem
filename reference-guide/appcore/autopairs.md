# autopairs.go - Auto-close Brackets/Quotes and Auto-indent

**Path:** `/home/javanhut/Development/Vem/internal/appcore/autopairs.go`
**Lines:** ~280
**Purpose:** Automatic bracket/quote pairing and smart indentation

## Overview

This file implements two related features for insert mode:
1. **Auto-pairs** - Automatically insert matching closing brackets/quotes
2. **Auto-indent** - Smart indentation on newline with bracket expansion

## Configuration

Fields in `appState` (app.go):
```go
autoPairsEnabled  bool   // Toggle auto-close brackets/quotes
autoIndentEnabled bool   // Toggle smart indentation
indentString      string // "\t" or "    " (spaces)
```

## Commands

| Command | Effect |
|---------|--------|
| `:autopairs` | Enable auto-close brackets/quotes |
| `:noautopairs` | Disable auto-close brackets/quotes |
| `:autoindent` | Enable smart indentation |
| `:noautoindent` | Disable smart indentation |
| `:expandtab` | Use 4 spaces for indent |
| `:noexpandtab` | Use tabs for indent |

## Auto-pairs Behavior

### Supported Pairs
```go
var autoPairs = map[rune]rune{
    '(':  ')',
    '[':  ']',
    '{':  '}',
    '"':  '"',
    '\'': '\'',
}
```

### Rules
1. **Opening char** - Inserts both opening and closing, cursor between them
2. **Closing char** - If cursor is before same closing char, skips over it (no duplicate)
3. **Quote exceptions**:
   - Skip after `\` (escape sequence like `\'`)
   - Skip after alphanumeric (contractions like `it's`)

### Integration Points
Called from two places in app.go:

1. **KeyEvent path** - `handleInsertModeSpecial()`:
```go
if s.handleAutoPairInsertion(r) {
    // Handled by auto-pairs
    return true
}
// Fall through to normal insertion
```

2. **EditEvent path** - `handleEvents()` (for platforms that use EditEvent for character input):
```go
if len(e.Text) == 1 {
    r := rune(e.Text[0])
    if s.handleAutoPairInsertion(r) {
        // Handled by auto-pairs
        continue
    }
}
// Fall through to normal insertion
```

## Auto-indent Behavior

### Rules
1. **Preserve indent** - New line inherits leading whitespace from current line
2. **Extra indent** - Added after lines ending with `{`, `(`, `[`, or `:` (Python files)
3. **Bracket expansion** - Cursor between `{}`, `[]`, `()` + Enter expands to 3 lines:
   ```
   {
       |cursor
   }
   ```
4. **Auto-dedent** - Typing `}`, `]`, `)` at line start (only whitespace before) removes one indent level
5. **Comment detection** - Skipped when cursor is inside a comment (via Chroma tokenization)

### Integration Point
Called from `ActionInsertNewline` in keybindings.go:
```go
if s.autoIndentEnabled && !s.isCursorInComment() {
    s.insertNewlineWithIndent()
} else {
    s.insertText("\n")
}
```

## Key Functions

### Auto-pairs Functions

| Function | Line | Purpose |
|----------|------|---------|
| `isOpeningChar(r)` | ~30 | Returns closing char if r is opening bracket/quote |
| `isClosingChar(r)` | ~36 | Returns true if r is closing bracket/quote |
| `charBeforeCursor()` | ~40 | Gets rune before cursor position |
| `charAtCursor()` | ~55 | Gets rune at cursor position |
| `shouldSkipClosingChar(r)` | ~70 | Checks if char at cursor matches typed closing char |
| `insertAutoPair(open, close)` | ~80 | Inserts both chars, positions cursor between |
| `handleAutoPairInsertion(r)` | ~90 | Main entry point for auto-pair logic |

### Auto-indent Functions

| Function | Line | Purpose |
|----------|------|---------|
| `getLeadingWhitespace(line)` | ~140 | Extracts leading spaces/tabs from line |
| `shouldAddExtraIndent(text, path)` | ~150 | Checks if line ends with block-opener |
| `isPythonLike(filePath)` | ~175 | Checks if file is Python/YAML |
| `isCursorBetweenBrackets()` | ~180 | Checks if cursor is between matching brackets |
| `isCursorInComment()` | ~195 | Uses Chroma tokens to detect comments |
| `isCommentTokenType(t)` | ~225 | Checks if token type is comment variant |
| `insertNewlineWithIndent()` | ~235 | Main smart newline function |
| `handleAutoDedent(r)` | ~285 | Handles dedent when typing closing bracket |
| `dedentOnce(indent, unit)` | ~320 | Removes one indent level |

## Comment Detection

Uses Chroma syntax highlighting tokens:
```go
func isCommentTokenType(t chroma.TokenType) bool {
    return t == chroma.Comment ||
           t == chroma.CommentSingle ||
           t == chroma.CommentMultiline ||
           t == chroma.CommentSpecial ||
           t == chroma.CommentPreproc ||
           t == chroma.CommentPreprocFile
}
```

## Edge Cases Handled

1. **Empty buffer** - Nil checks throughout
2. **Terminal buffers** - Auto-features skipped
3. **End of line** - `charAtCursor()` returns 0
4. **Start of line** - `charBeforeCursor()` returns 0
5. **Read-only buffers** - Checked in Buffer methods
6. **Escape sequences** - `\'` and `\"` not auto-closed

---
*Last Updated: Auto-pairs and auto-indent implementation*
