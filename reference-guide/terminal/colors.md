# colors.go

**Path:** `/home/javanhut/Development/Vem/internal/terminal/colors.go`
**Lines:** 93
**Purpose:** ANSI color palette and vt10x color conversion

## Overview

Provides:
- Standard 16-color ANSI palette
- 256-color extended palette support
- True color (24-bit RGB) support
- vt10x color format conversion

## Code Blocks

### Lines 1-3: Package and Imports

```go
package terminal
import "image/color"
```

### Lines 5-26: ANSI 16-Color Palette

```go
var ansiColors = [16]color.NRGBA{
    // Normal colors (0-7)
    {R: 0x00, G: 0x00, B: 0x00, A: 0xff}, // 0: Black
    {R: 0xcc, G: 0x00, B: 0x00, A: 0xff}, // 1: Red
    {R: 0x4e, G: 0x9a, B: 0x06, A: 0xff}, // 2: Green
    {R: 0xc4, G: 0xa0, B: 0x00, A: 0xff}, // 3: Yellow
    {R: 0x34, G: 0x65, B: 0xa4, A: 0xff}, // 4: Blue
    {R: 0x75, G: 0x50, B: 0x7b, A: 0xff}, // 5: Magenta
    {R: 0x06, G: 0x98, B: 0x9a, A: 0xff}, // 6: Cyan
    {R: 0xd3, G: 0xd7, B: 0xcf, A: 0xff}, // 7: White

    // Bright colors (8-15)
    {R: 0x55, G: 0x57, B: 0x53, A: 0xff}, // 8: Bright Black
    {R: 0xef, G: 0x29, B: 0x29, A: 0xff}, // 9: Bright Red
    {R: 0x8a, G: 0xe2, B: 0x34, A: 0xff}, // 10: Bright Green
    {R: 0xfc, G: 0xe9, B: 0x4f, A: 0xff}, // 11: Bright Yellow
    {R: 0x72, G: 0x9f, B: 0xcf, A: 0xff}, // 12: Bright Blue
    {R: 0xad, G: 0x7f, B: 0xa8, A: 0xff}, // 13: Bright Magenta
    {R: 0x34, G: 0xe2, B: 0xe2, A: 0xff}, // 14: Bright Cyan
    {R: 0xee, G: 0xee, B: 0xec, A: 0xff}, // 15: Bright White
}
```

### Lines 28-52: GetANSIColor

#### GetANSIColor (Lines 28-52)
Returns color for ANSI color code:
1. **Codes 0-15:** Returns from ansiColors array
2. **Codes 16-231:** 6x6x6 color cube (216 colors)
   ```go
   code -= 16
   r := (code / 36) * 51
   g := ((code % 36) / 6) * 51
   b := (code % 6) * 51
   ```
3. **Codes 232-255:** Grayscale ramp (24 shades)
   ```go
   gray := (code - 232) * 10 + 8
   ```

### Lines 54-58: Default Colors

```go
var (
    DefaultFG = ansiColors[7]  // White
    DefaultBG = color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}  // Match Vem bg
)
```

### Lines 60-92: vt10x Color Conversion

#### vt10xColorToNRGBA (Lines 60-92)
Converts vt10x color format to NRGBA:
1. **DefaultFG value:** `1 << 24`
2. **DefaultBG value:** `(1 << 24) + 1`
3. **ANSI codes:** `< 256`
4. **True color:** `(1 << 25) | (r << 16) | (g << 8) | b`

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Used by buffer.go for default cell colors
- Called during vt10x state sync in terminal.go
- Colors match Vem's default theme (#1a1f2e background)

## Color Encoding

```
256-Color Palette Layout:
┌─────────────────────────────────────────────┐
│ 0-15:   Standard ANSI colors               │
│ 16-231: 6x6x6 color cube (216 colors)      │
│ 232-255: Grayscale ramp (24 shades)        │
└─────────────────────────────────────────────┘

vt10x True Color Format:
┌─────┬────────┬────────┬────────┐
│ 1   │ R (8)  │ G (8)  │ B (8)  │
│ bit │ bits   │ bits   │ bits   │
└─────┴────────┴────────┴────────┘
  ^25   ^16-23   ^8-15    ^0-7
```

---
*Last Updated: Reference guide creation*
