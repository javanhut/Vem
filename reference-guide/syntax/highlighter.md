# highlighter.go

**Path:** `/home/javanhut/Development/Vem/internal/syntax/highlighter.go`
**Lines:** 223
**Purpose:** Syntax highlighting with Chroma integration and caching

## Overview

Provides syntax highlighting using the Chroma library with:
- Auto-detection of language from file extension
- Per-line caching with hash-based change detection
- Theme support (10 preset themes)
- Binary file detection

## Code Blocks

### Lines 1-13: Package and Imports

```go
package syntax
import (
    "hash/fnv"
    "path/filepath"
    "strings"
    "github.com/alecthomas/chroma/v2"
    "github.com/alecthomas/chroma/v2/lexers"
    "github.com/alecthomas/chroma/v2/styles"
    _ "github.com/javanhut/vem/internal/syntax/lexers"  // Custom lexers
)
```

### Lines 15-35: Type Definitions

#### Token (Lines 15-20)
```go
type Token struct {
    Text  string            // Token text
    Type  chroma.TokenType  // Keyword, String, etc.
    Style *chroma.Style     // Color theme
}
```

#### HighlightedLine (Lines 22-26)
```go
type HighlightedLine struct {
    Tokens []Token  // Tokenized line
    Hash   uint64   // Hash for change detection
}
```

#### Highlighter (Lines 28-35)
```go
type Highlighter struct {
    lexer     chroma.Lexer           // Language tokenizer
    style     *chroma.Style          // Color theme
    cache     map[int]*HighlightedLine  // Line number -> cached result
    formatter *chroma.Formatter      // (unused)
    enabled   bool                   // Highlighting on/off
}
```

### Lines 37-78: Constructors

#### NewHighlighter (Lines 37-68)
Creates highlighter for file:
1. Match lexer by filename
2. Fall back to extension matching
3. Fall back to plain text lexer
4. Default theme: `monokai`

#### NewPlainHighlighter (Lines 70-78)
Creates disabled highlighter for plain text.

### Lines 80-129: Core Highlighting

#### HighlightLine (Lines 80-129)
Highlights single line with caching:
1. If disabled, return plain text token
2. Check cache by line number + content hash
3. If cache miss, tokenize with Chroma lexer
4. Convert Chroma tokens to internal Token type
5. Cache and return result

**Cache Key:** Line number + FNV-64a hash of line content

### Lines 131-166: Cache and Theme Management

| Function | Lines | Purpose |
|----------|-------|---------|
| `InvalidateLine()` | 131-134 | Remove line from cache |
| `InvalidateAll()` | 136-139 | Clear entire cache |
| `SetTheme()` | 141-148 | Change theme, clear cache |
| `GetThemeName()` | 150-153 | Get current theme name |
| `SetEnabled()` | 155-161 | Enable/disable highlighting |
| `IsEnabled()` | 163-166 | Check if enabled |

### Lines 168-195: Utilities

| Function | Lines | Purpose |
|----------|-------|---------|
| `GetLanguage()` | 168-177 | Get detected language name |
| `ListAvailableThemes()` | 179-188 | List all Chroma themes |
| `hashString()` | 190-195 | FNV-64a hash for caching |

### Lines 197-222: File Detection

#### ShouldHighlight (Lines 197-222)
Determines if file should be highlighted:
- Returns false for empty path
- Returns false for binary extensions:
  - Binaries: `.bin`, `.exe`, `.so`, `.dylib`, `.dll`
  - Images: `.jpg`, `.png`, `.gif`, etc.
  - Media: `.mp3`, `.mp4`, `.avi`, etc.
  - Archives: `.zip`, `.tar`, `.gz`, etc.
  - Documents: `.pdf`, `.doc`, `.xls`, etc.

## Known Issues / Potential Bugs

1. **Line 33: `formatter` field unused**
   - Field declared but never used
   - Could be removed

2. **Line 57: Hardcoded default theme**
   - `"monokai"` is hardcoded
   - Should use a constant or config

3. **No file size check in ShouldHighlight**
   - Large text files will still be highlighted
   - Could add size threshold

## Dead/Unused Code

1. **Line 33: `formatter *chroma.Formatter`**
   - Never assigned or used
   - Safe to remove

## Integration Points

- Used by `appState.getOrCreateHighlighter()` in app.go
- Token colors resolved via `theme.go:GetTokenColor()`
- Custom Carrion lexer in `lexers/carrion.go`

---
*Last Updated: Reference guide creation*
