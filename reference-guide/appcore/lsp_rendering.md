# lsp_rendering.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/lsp_rendering.go`
**Lines:** ~905
**Purpose:** LSP UI rendering (completion menu, hover, diagnostics, buffer completion)

## Overview

Renders all LSP-related UI components:
- Diagnostic underlines (wavy squiggles)
- Completion popup menu
- Hover tooltips
- References list
- Code actions menu
- Buffer word completion popup (NEW)
- Diagnostics list view (NEW)

## Code Blocks

### Lines 1-38: Package and Colors

#### LSP UI Colors
```go
var (
    diagnosticErrorColor   = color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
    diagnosticWarningColor = color.NRGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff}
    diagnosticInfoColor    = color.NRGBA{R: 0x00, G: 0xbf, B: 0xff, A: 0xff}
    diagnosticHintColor    = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}

    completionBgColor       = color.NRGBA{R: 0x2b, G: 0x2b, B: 0x2b, A: 0xf0}
    completionSelectedColor = color.NRGBA{R: 0x00, G: 0x7a, B: 0xcc, A: 0xff}
    completionBorderColor   = color.NRGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xff}

    hoverBgColor     = color.NRGBA{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xf5}
    hoverBorderColor = color.NRGBA{R: 0x45, G: 0x45, B: 0x45, A: 0xff}

    referencesBgColor = color.NRGBA{R: 0x1a, G: 0x1a, B: 0x1a, A: 0xf0}
)
```

### Lines 40-97: Diagnostic Underlines

#### drawDiagnosticUnderlines (Lines 40-97)
Draws wavy underlines under diagnostics:
1. Gets diagnostics for line
2. Calculates start/end columns
3. Chooses color by severity
4. Draws wavy line

### Lines 99-127: drawWavyUnderline

Draws simulated wavy line using rectangles:
```go
step := 4
amplitude := 2

up := true
for px := x; px < x+width; px += step {
    dy := 0
    if up {
        dy = -amplitude
    }
    // Draw segment...
    up = !up
}
```

### Lines 129-276: Completion Menu

#### drawCompletionMenu (Lines 129-276)
Draws the completion popup:
1. Calculates dimensions (max 400px wide, 10 items visible)
2. Positions below cursor
3. Draws background and border
4. Renders visible items with:
   - Selection highlight
   - Icon for completion kind
   - Label text
   - Detail (type info)
5. Draws documentation panel to the right

### Lines 278-345: Hover Tooltip

#### drawHoverTooltip (Lines 278-345)
Draws hover information:
1. Wraps text to 80 characters
2. Calculates popup dimensions
3. Positions above cursor (or below if no room)
4. Draws background, border, content

### Lines 347-440: References List

#### drawReferencesList (Lines 347-440)
Draws references at bottom of screen:
1. Shows up to 8 visible items
2. Highlights selected item
3. Shows file:line:col for each
4. Displays navigation hint

### Lines 442-503: Code Actions Menu

#### drawCodeActionsMenu (Lines 442-503)
Draws code actions popup:
1. Positions near cursor
2. Shows all available actions
3. Highlights selected action

### Lines 505-530: drawBorder

Draws 1px border using 4 rectangles (top, bottom, left, right).

### Lines 532-536: measureCharWidth

Returns approximate character width (hardcoded 8px).

### Lines 538-594: getCompletionIcon

Returns icon character for completion kind:

| Kind | Icon |
|------|------|
| Method | m |
| Function | f |
| Variable | v |
| Class | c |
| Interface | i |
| Keyword | k |
| Snippet | s |
| ... | ... |

### Lines 596-616: getCompletionIconColor

Returns color for completion kind:
- Method/Function: Yellow
- Variable/Field: Light blue
- Class/Struct/Interface: Teal
- Module/Folder: Orange
- Keyword: Purple

### Lines 618-695: Text Utilities

| Function | Lines | Purpose |
|----------|-------|---------|
| `truncateString()` | 618-624 | Truncate with "..." |
| `completionDocText()` | 626-639 | Extract doc text |
| `completionDocString()` | 641-656 | Parse doc from interface{} |
| `wrapText()` | 658-687 | Word wrap text |
| `formatInt()` | 689-695 | Format int to string |

## Known Issues / Potential Bugs

1. **Line 533-536: Hardcoded character width**
   - Returns 8 regardless of actual font metrics
   - Could cause alignment issues

## Dead/Unused Code

None identified.

## Integration Points

- Called from buffer rendering in app.go
- Uses lsp_actions.go state variables
- Triggered by completion, hover, references state

## Visual Layout

```
┌─────────────────────────────────────────────────────────┐
│ Editor Content                                          │
│                                                         │
│  function example() {                                   │
│      let value = something.meth|                        │
│                           ┌──────────────┬──────────┐   │
│                           │ m method     │ Docs...  │   │
│                           │ f methodTwo  │          │   │
│                           │ v methodProp │          │   │
│                           └──────────────┴──────────┘   │
│                              ▲ Completion menu          │
│                                                         │
├─────────────────────────────────────────────────────────┤
│ 1. file.go:42:10                   [j/k: nav] [Enter]  │
│ 2. other.go:15:5                                        │
│ 3. test.go:100:20                                       │
└─────────────────────────────────────────────────────────┘
  ▲ References list (bottom of screen)
```

### Lines 697-793: Buffer Word Completion Menu (NEW)

#### drawBufferCompletionMenu (Lines 697-793)
Draws completion popup for buffer word completion (Ctrl+N style):
1. Only renders when `bufferCompletionActive` is true
2. Positions popup below cursor
3. Shows matching words from current buffer
4. Highlights selected item
5. Shows navigation hints (Tab/Esc/Ctrl+N)

### Lines 795-904: Diagnostics List View (NEW)

#### drawDiagnosticsList (Lines 795-904)
Draws diagnostics list at bottom of screen (`:diagnostics` command):
1. Shows up to 10 visible items with scrolling
2. Color-coded by severity:
   - Red for errors
   - Yellow for warnings
   - Blue for information
   - Gray for hints
3. Format: `[E] L12:C4 message text`
4. Header shows total count and navigation hints
5. j/k to navigate, Enter to jump, Esc to close

---
*Last Updated: After enhancement plan implementation*
