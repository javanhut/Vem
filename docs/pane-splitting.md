# Pane Splitting in Vem

Pane splitting allows you to view and edit multiple files side-by-side in Vem, similar to Tmux or Vim splits.

## Overview

Vem supports **mutually exclusive panes** where each pane displays exactly one buffer. You can split panes vertically (top/bottom) or horizontally (left/right), navigate between them, and manage their layout.

## Key Features

- Split panes vertically (horizontal divider: `─`) or horizontally (vertical divider: `│`)
- Each pane shows a unique buffer (no duplicate files across panes)
- Active pane is brighter, inactive panes are dimmer
- Subtle separator lines between panes
- Vim-style navigation with `Alt+hjkl`
- Zoom individual panes to full screen temporarily

## Keybindings

### Creating Splits

| Keybinding | Action | Visual Result |
|------------|--------|---------------|
| `Ctrl+S v` | Split vertical | Creates horizontal divider (top/bottom) |
| `Ctrl+S h` | Split horizontal | Creates vertical divider (left/right) |

**Note**: After splitting, the new pane starts with an empty buffer. Use `:e filename` or `Ctrl+P` to open a file.

### Navigating Panes

| Keybinding | Action |
|------------|--------|
| `Alt+h` | Focus pane to the left |
| `Alt+j` | Focus pane below |
| `Alt+k` | Focus pane above |
| `Alt+l` | Focus pane to the right |
| `Shift+Tab` | Cycle to next pane |

### Managing Panes

| Keybinding | Action |
|------------|--------|
| `Ctrl+X` | Close active pane (prompts if unsaved changes) |
| `Ctrl+S =` | Equalize all panes (make them 50/50) |
| `Ctrl+S o` | Toggle zoom (maximize/restore active pane) |

## Workflow Examples

### Example 1: Side-by-Side File Editing

```
1. Open main.go
   ┌──────────────────────────────────┐
   │ main.go                          │
   │                                  │
   └──────────────────────────────────┘

2. Press Ctrl+S h (horizontal split)
   ┌──────────────┬───────────────────┐
   │ main.go      │ [No Name]         │
   │              │                   │
   └──────────────┴───────────────────┘

3. Press :e config.go
   ┌──────────────┬───────────────────┐
   │ main.go      │ config.go         │
   │              │ [ACTIVE]          │
   └──────────────┴───────────────────┘

4. Press Alt+h to focus left pane
   ┌──────────────┬───────────────────┐
   │ main.go      │ config.go         │
   │ [ACTIVE]     │                   │
   └──────────────┴───────────────────┘
```

### Example 2: Complex Layout

```
1. Start with main.go

2. Ctrl+S h → Split horizontal
   ┌──────────┬────────────┐
   │ main.go  │ [New]      │
   └──────────┴────────────┘

3. Alt+l, :e app.go
   ┌──────────┬────────────┐
   │ main.go  │ app.go     │
   └──────────┴────────────┘

4. Ctrl+S v → Split vertical on right pane
   ┌──────────┬────────────┐
   │ main.go  │ app.go     │
   │          ├────────────┤
   │          │ [New]      │
   └──────────┴────────────┘

5. :e test.go
   ┌──────────┬────────────┐
   │ main.go  │ app.go     │
   │          ├────────────┤
   │          │ test.go    │
   └──────────┴────────────┘
```

## Visual Styling

### Active vs Inactive Panes

- **Active pane**: Normal background color (brighter)
- **Inactive panes**: Slightly darker background (15% darker)
- **Separator**: Subtle 1px gray line

### Colors

```go
activePaneBg   = #1a1f2e  // Brighter (active)
inactivePaneBg = #141824  // Dimmer (inactive)
paneSeparator  = #303544  // Subtle gray
```

## Pane Lifecycle

### Opening Files in Panes

After creating a new pane, it starts with an empty buffer:

```
Option 1: Use :e command
  :e path/to/file.go

Option 2: Use Fuzzy Finder
  Ctrl+P, then type filename
```

### Closing Panes

When you close a pane with `Ctrl+X`:
1. Checks if buffer has unsaved changes
2. Closes the pane
3. Closes the buffer (since it's unique to that pane)
4. Focuses the next available pane

**Note**: You cannot close the last pane.

## Advanced Features

### Zoom Mode

Press `Ctrl+S o` to temporarily maximize the active pane:

```
Before zoom (3 panes):
┌──────┬─────┐
│  A   │  B  │
├──────┴─────┤
│     C      │
└────────────┘

After Ctrl+S o on pane A:
┌────────────┐
│            │
│     A      │
│  (zoomed)  │
│            │
└────────────┘

Press Ctrl+S o again to restore:
┌──────┬─────┐
│  A   │  B  │
├──────┴─────┤
│     C      │
└────────────┘
```

### Equalize Panes

Press `Ctrl+S =` to make all panes equal size (50/50 splits):

```
Before (uneven):
┌─────┬──────────┐
│  A  │    B     │
└─────┴──────────┘

After Ctrl+S =:
┌──────┬──────┐
│  A   │  B   │
└──────┴──────┘
```

## Status Bar Information

The status bar shows pane information:

```
MODE NORMAL | FILE main.go | CURSOR 42:15 | PANE 2/3 | Ready
                                           ^^^^^^^^
                                           Active pane 2 of 3 total
```

When zoomed:
```
MODE NORMAL | FILE main.go | CURSOR 42:15 | PANE 1/3 | ZOOMED | Ready
```

## Tips and Tricks

1. **Quick file comparison**: Use `Ctrl+S h` to split, then open related files side-by-side

2. **Reference while editing**: Keep documentation in one pane, code in another

3. **Test and implementation**: Put test file in left pane, implementation in right pane

4. **Cycle through panes**: Use `Shift+Tab` when you have many panes open

5. **Temporary full screen**: Use `Ctrl+S o` to focus on one file, then `Ctrl+S o` again to restore

## Limitations

- Each pane must show a different buffer (no duplicate files)
- All splits are 50/50 (no custom ratios yet)
- Maximum practical panes: ~4-6 (more becomes hard to navigate)
- Pane layout is not persistent (resets on restart)

## Architecture

Vem uses a binary tree structure for pane layout:

- **Leaf nodes**: Contain actual panes showing buffers
- **Internal nodes**: Represent splits (horizontal or vertical)
- **Navigation**: Uses geometric calculation to find panes in each direction

For technical details, see [Architecture.md](Architecture.md#pane-system).

## See Also

- [Keybindings Reference](keybindings.md) - Complete keybinding documentation
- [Tutorial](tutorial.md) - Getting started with Vem
- [Architecture](Architecture.md) - Technical implementation details
