# syntax Package Index

**Path:** `/home/javanhut/Development/Vem/internal/syntax/`
**Purpose:** Syntax highlighting with Chroma library integration

## Files Overview

| File | Lines | Purpose |
|------|-------|---------|
| [highlighter.go](highlighter.md) | 223 | Main highlighter with caching |
| [theme.go](theme.md) | 113 | Theme handling and colors |

**Additional files:**
- `lexers/carrion.go` - Custom Carrion language lexer

## Key Types

### Highlighter (highlighter.go)
Per-buffer highlighter:
- Chroma lexer
- Style/theme
- Line cache (hash-based)
- Enable/disable state

### Token (highlighter.go)
```go
type Token struct {
    Text  string
    Type  chroma.TokenType
    Style *chroma.Style
}
```

### HighlightedLine (highlighter.go)
```go
type HighlightedLine struct {
    Tokens []Token
    Hash   uint64  // FNV-64a hash for change detection
}
```

## Available Themes

| Theme | Description |
|-------|-------------|
| monokai | Default, vibrant colors |
| dracula | Purple accents |
| github-dark | GitHub dark mode |
| nord | Arctic colors |
| one-dark | Atom One Dark |
| solarized-dark | Precision colors |
| solarized-light | Light variant |
| vim | Classic Vim |
| catppuccin-mocha | Warm pastels |
| gruvbox | Retro groove |

## Known Issues Summary

1. **highlighter.go:33** - Unused `formatter` field
2. **highlighter.go:57** - Hardcoded default theme "monokai"

## Token Color Resolution

```
Token Type Hierarchy:
  Keyword.Declaration
        ↓ (not found)
      Keyword
        ↓ (not found)
       Token
        ↓ (not found)
    Default Color (#dfe7ff)
```

## Caching Algorithm

```
HighlightLine(lineNumber, lineContent):
  hash = FNV64(lineContent)

  if cache[lineNumber].Hash == hash:
    return cache[lineNumber].Tokens  // Cache hit

  tokens = lexer.Tokenise(lineContent)
  cache[lineNumber] = HighlightedLine{tokens, hash}
  return tokens
```

## Binary File Detection

Files with these extensions skip highlighting:
- Binaries: `.bin`, `.exe`, `.so`, `.dylib`, `.dll`
- Images: `.jpg`, `.png`, `.gif`, `.bmp`, `.ico`
- Media: `.mp3`, `.mp4`, `.avi`, `.wav`
- Archives: `.zip`, `.tar`, `.gz`, `.rar`
- Documents: `.pdf`, `.doc`, `.xls`, `.ppt`

---
*Last Updated: Reference guide creation*
