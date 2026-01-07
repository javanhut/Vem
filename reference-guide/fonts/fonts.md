# fonts.go

**Path:** `/home/javanhut/Development/Vem/internal/fonts/fonts.go`
**Lines:** 46
**Purpose:** Embedded JetBrains Mono Nerd Font loading

## Overview

Provides embedded font loading for:
- JetBrains Mono Nerd Font Regular
- JetBrains Mono Nerd Font Bold

Uses Go's `//go:embed` directive to bundle fonts in binary.

## Code Blocks

### Lines 1-9: Package and Imports

```go
package fonts
import (
    _ "embed"
    "fmt"

    "gioui.org/font"
    "gioui.org/font/opentype"
)
```

### Lines 11-15: Embedded Font Data

```go
//go:embed JetBrainsMonoNerdFont-Regular.ttf
var jetbrainsMonoRegular []byte

//go:embed JetBrainsMonoNerdFont-Bold.ttf
var jetbrainsMonoBold []byte
```

**Note:** The TTF files must exist in the same directory as fonts.go.

### Lines 17-45: Collection

#### Collection (Lines 17-45)
Returns collection of parsed font faces:
1. Parses regular font with `opentype.Parse()`
2. Creates font face with typeface "JetBrainsMono"
3. Parses bold font
4. Creates bold font face with `Weight: font.Bold`

```go
fonts = append(fonts, font.FontFace{
    Font: font.Font{Typeface: "JetBrainsMono"},
    Face: regularFace,
})

fonts = append(fonts, font.FontFace{
    Font: font.Font{
        Typeface: "JetBrainsMono",
        Weight:   font.Bold,
    },
    Face: boldFace,
})
```

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Called from `appcore.Run()` during startup
- Registered with Gio's font shaper
- Used throughout app for text rendering

## Font Files Required

| File | Purpose |
|------|---------|
| `JetBrainsMonoNerdFont-Regular.ttf` | Normal text |
| `JetBrainsMonoNerdFont-Bold.ttf` | Bold text, headers |

## Why Nerd Font?

JetBrains Mono Nerd Font includes:
- Monospace characters for code
- Programming ligatures
- Powerline symbols
- File icons (used in file explorer)
- Git symbols
- Developer glyphs

---
*Last Updated: Reference guide creation*
