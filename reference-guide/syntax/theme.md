# theme.go

**Path:** `/home/javanhut/Development/Vem/internal/syntax/theme.go`
**Lines:** 113
**Purpose:** Color theme handling and token color resolution

## Overview

Provides theme color resolution for syntax highlighting with:
- Token type to color mapping with hierarchy fallback
- Background color extraction
- Dark/light theme detection
- Preset theme list with descriptions

## Code Blocks

### Lines 1-8: Package and Imports

```go
package syntax
import (
    "image/color"
    "github.com/alecthomas/chroma/v2"
)
```

### Lines 9-33: Token Color Resolution

#### GetTokenColor (Lines 9-33)
Resolves token type to RGBA color:
1. Check if style is nil → return default (#dfe7ff)
2. Get style entry for token type
3. If color is set, convert and return
4. Otherwise, try parent token type (hierarchy)
5. Fall back to default color

**Token Hierarchy Example:**
```
Keyword.Declaration → Keyword → Token
```

### Lines 35-61: Color Utilities

#### chromaColorToNRGBA (Lines 35-45)
Converts Chroma color (0xRRGGBB integer) to `color.NRGBA`.

#### GetBackgroundColor (Lines 47-61)
Extracts background color from theme style.
Default: `#1a1f2e`

### Lines 63-96: Preset Themes

#### PresetThemes (Lines 63-75)
```go
var PresetThemes = []string{
    "monokai",          // Default, vibrant
    "dracula",          // Purple accents
    "github-dark",      // GitHub dark
    "nord",             // Arctic colors
    "one-dark",         // Atom One Dark
    "solarized-dark",   // Precision colors
    "solarized-light",  // Light variant
    "vim",              // Classic Vim
    "catppuccin-mocha", // Warm pastels
    "gruvbox",          // Retro groove
}
```

#### GetThemeDescription (Lines 77-96)
Returns human-readable description for theme.

### Lines 98-112: Theme Detection

#### IsDarkTheme (Lines 98-112)
Determines if theme is dark using perceived brightness:
```go
brightness = 0.299*R + 0.587*G + 0.114*B
return brightness < 128  // Dark if < 128
```

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Called from `drawBufferLine()` in app.go
- Used by highlighter.go for token coloring
- Theme switching via `:theme <name>` command

## Adding New Themes

1. Add to `PresetThemes` slice
2. Add description to `GetThemeDescription()` map
3. Theme must exist in Chroma library (`styles.Get()`)

---
*Last Updated: Reference guide creation*
