# fonts Package Index

**Path:** `/home/javanhut/Development/Vem/internal/fonts/`
**Purpose:** Embedded font loading for the application

## Files Overview

| File | Lines | Purpose |
|------|-------|---------|
| [fonts.go](fonts.md) | 46 | JetBrains Mono Nerd Font loading |

## Embedded Files

| File | Purpose |
|------|---------|
| `JetBrainsMonoNerdFont-Regular.ttf` | Normal text |
| `JetBrainsMonoNerdFont-Bold.ttf` | Bold text, headers |

## Key Function

### Collection()
Returns `[]font.FontFace` for Gio font shaper:
- Regular weight face
- Bold weight face

```go
fonts, err := fonts.Collection()
if err != nil {
    // handle error
}
// Register with Gio shaper
```

## Why JetBrains Mono Nerd Font?

1. **Monospace** - Essential for code editing
2. **Programming Ligatures** - Optional ligature support
3. **Nerd Font Glyphs**:
   - File type icons (used in file explorer)
   - Git symbols
   - Powerline symbols
   - Developer glyphs

## Font Registration

In app.go initialization:
```go
fonts, err := fonts.Collection()
// ...
shaper := text.NewShaper(text.WithCollection(fonts))
```

## Glyph Examples

| Glyph | Use |
|-------|-----|
|  | Go files |
|  | Folder |
|  | Git |
|  | Config |
|  | Text file |

---
*Last Updated: Reference guide creation*
