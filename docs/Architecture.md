# Vem Architecture

Technical architecture and design decisions for Vem - a modern Vim emulator built from scratch in Go.

## Table of Contents

- [Overview](#overview)
- [Core Components](#core-components)
- [Rendering Pipeline](#rendering-pipeline)
- [Keybinding System](#keybinding-system)
- [Buffer Management](#buffer-management)
- [Pane System](#pane-system)
- [File System Integration](#file-system-integration)
- [Fuzzy Finder](#fuzzy-finder)
- [Search System](#search-system)
- [Undo System](#undo-system)
- [Platform Abstraction](#platform-abstraction)
- [Design Decisions](#design-decisions)

## Overview

Vem is a cross-platform Vim emulator built with Go and Gio UI. The architecture prioritizes:

1. **Cross-platform compatibility** - Single codebase for Linux, macOS, Windows, WebAssembly
2. **Zero external dependencies** - Bundles fonts and assets, no system packages required
3. **GPU acceleration** - Smooth rendering via Vulkan/Metal/Direct3D/WebGL
4. **Modal editing** - Complete Vim-like editing paradigm with robust state management
5. **Modern features** - Pane splitting, fuzzy finding, search highlighting, undo support

### Technology Stack

- **Language**: Go 1.25.3+
- **UI Framework**: Gio UI v0.9.0
- **Graphics**: Vulkan (Linux), Metal (macOS), Direct3D (Windows), WebGL (WebAssembly)
- **License**: GPLv2

## Core Components

### 1. Application Core (`internal/appcore/`)

The application core manages the main event loop, rendering, and user interaction.

#### `app.go`

Main application state and event handling:

```go
type appState struct {
    theme        *material.Theme
    bufferMgr    *editor.BufferManager
    paneManager  *panes.PaneManager
    fileTree     *filesystem.FileTree
    mode         mode
    status       string

    // Explorer state
    explorerVisible bool
    explorerFocused bool

    // Modifier tracking
    ctrlPressed  bool
    shiftPressed bool

    // Fuzzy finder state
    fuzzyFinderActive bool
    fuzzyPattern      string
    fuzzyResults      []FuzzyMatch
    fuzzySelected     int

    // Search state
    searchPattern     string
    searchMatches     []SearchMatch
    searchCurrentIdx  int
}
```

**Responsibilities:**
- Window creation and lifecycle management
- Event loop (keyboard, mouse, window events)
- Modal state machine (NORMAL, INSERT, VISUAL, DELETE, COMMAND, EXPLORER, SEARCH, FUZZY_FINDER)
- Rendering pipeline orchestration
- Status bar and UI chrome
- Caret blinking animation
- Pane layout coordination

**Key Methods:**
- `run()` - Main event loop
- `handleEvents()` - Process input events
- `handleKey()` - Dispatch key events to appropriate handlers
- `layout()` - Render UI components
- `executeAction()` - Execute keybinding actions

#### `keybindings.go`

Keybinding system and action dispatch:

```go
type Action int

const (
    ActionToggleExplorer
    ActionMoveUp
    ActionEnterInsert
    ActionSplitVertical
    ActionSplitHorizontal
    ActionPaneFocusLeft
    ActionOpenFuzzyFinder
    ActionEnterSearch
    ActionUndo
    // ... more actions
)

type KeyBinding struct {
    Modifiers key.Modifiers
    Key       key.Name
    Modes     []mode
    Action    Action
}
```

**Responsibilities:**
- Define all editor actions as enum constants
- Map key combinations to actions
- Handle global vs mode-specific keybindings
- Modifier key matching with platform quirk handling
- Action execution routing
- Support for Ctrl+S prefix sequences

**Key Methods:**
- `matchGlobalKeybinding()` - Check global shortcuts
- `matchModeKeybinding()` - Check mode-specific bindings
- `modifiersMatch()` - Handle platform-specific modifier detection
- `executeAction()` - Dispatch actions to implementation

#### `pane_actions.go`

Pane management actions:

**Responsibilities:**
- Handle pane splitting (vertical/horizontal)
- Navigate between panes (Alt+hjkl)
- Close panes with buffer cleanup
- Zoom/unzoom panes
- Equalize pane sizes
- Cycle through panes

#### `pane_rendering.go`

Pane rendering logic:

**Responsibilities:**
- Calculate pane layout bounds
- Render active/inactive pane backgrounds
- Draw pane separators
- Handle zoomed pane rendering
- Coordinate buffer rendering within panes

#### `fuzzy.go`

Fuzzy file finder implementation:

**Responsibilities:**
- Fuzzy matching algorithm with scoring
- File path filtering and ranking
- Match highlighting
- Result caching and sorting

### 2. Buffer Management (`internal/editor/`)

Text buffer abstraction and multi-buffer management.

#### `buffer.go`

Core text buffer implementation:

```go
type Buffer struct {
    lines    []string
    cursor   Cursor
    filePath string
    modified bool
    undoStack []UndoState

    // Visual mode state
    visualStart *Cursor
    visualMode  VisualMode

    // Search state
    searchMatches []SearchMatch
}

type Cursor struct {
    Line int
    Col  int
}

type UndoState struct {
    Lines  []string
    Cursor Cursor
}
```

**Responsibilities:**
- In-memory text storage as slice of lines
- UTF-8 aware cursor positioning
- Text mutation operations (insert, delete, modify)
- Line-based and word-based navigation
- Undo/redo support (up to 100 operations)
- Visual mode selection tracking
- Search match storage and navigation

**Key Methods:**
- `InsertText(text string)` - Insert text at cursor
- `DeleteBackward()` - Backspace operation
- `DeleteForward()` - Delete key operation
- `MoveLeft/Right/Up/Down()` - Cursor movement
- `MoveWordForward/Backward/End()` - Word motions
- `JumpLineStart/End()` - Line boundary jumps
- `DeleteLines(start, end)` - Multi-line deletion
- `Undo()` - Undo last operation
- `SaveUndoState()` - Capture state for undo
- `SetVisualStart()` - Begin visual selection
- `GetVisualSelection()` - Get selected range

#### `buffer_manager.go`

Multi-buffer management and file I/O:

```go
type BufferManager struct {
    buffers     []*Buffer
    activeIndex int
}
```

**Responsibilities:**
- Maintain list of open buffers
- Track active buffer
- File loading and saving
- Buffer switching (next/prev)
- Buffer lifecycle (open, close, save)
- Prevent duplicate file loading

**Key Methods:**
- `OpenFile(path)` - Load file into new buffer
- `SaveActiveBuffer()` - Save current buffer
- `SaveBufferAs(path)` - Save as new file
- `NextBuffer() / PrevBuffer()` - Switch buffers
- `CloseActiveBuffer(force)` - Close with modified check
- `ListBuffers()` - Get all open buffers
- `FindBufferByPath(path)` - Check if file already open

### 3. Pane Management (`internal/panes/`)

Window pane splitting and layout management.

#### `manager.go`

Pane tree manager:

```go
type PaneManager struct {
    root       *PaneNode
    activePane *Pane
    nextPaneID int
    zoomed     *Pane
}
```

**Responsibilities:**
- Manage binary tree of panes
- Track active pane
- Handle pane creation and destruction
- Zoom/unzoom functionality
- Pane cycling

**Key Methods:**
- `SplitVertical(bufferIndex)` - Create top/bottom split
- `SplitHorizontal(bufferIndex)` - Create left/right split
- `ClosePane()` - Remove pane from tree
- `SetActivePane(pane)` - Change active pane
- `CycleNextPane()` - Move to next pane
- `ToggleZoom()` - Maximize/restore pane

#### `layout.go`

Pane layout calculation:

**Responsibilities:**
- Convert pane tree to screen coordinates
- Calculate split positions (50/50)
- Handle recursive layout for nested splits
- Account for separator widths

#### `navigation.go`

Geometric pane navigation:

**Responsibilities:**
- Find pane in direction (left/right/up/down)
- Calculate pane centers for direction matching
- Handle complex layouts with multiple splits

**Key Methods:**
- `FindPaneInDirection(current, direction)` - Navigate by direction
- Uses geometric overlap detection

#### `pane.go`

Individual pane abstraction:

```go
type Pane struct {
    ID          string
    BufferIndex int
    Active      bool
    Bounds      image.Rectangle
}
```

**Responsibilities:**
- Track buffer index for this pane
- Store layout bounds
- Maintain active/inactive state

### 4. File System Integration (`internal/filesystem/`)

File tree navigation and directory browsing.

#### `tree.go`

File tree data structure and navigation:

```go
type FileTree struct {
    root         *Node
    selected     *Node
    currentPath  string
}

type Node struct {
    Name     string
    Path     string
    IsDir    bool
    Expanded bool
    Children []*Node
    Depth    int
}
```

**Responsibilities:**
- Hierarchical file tree representation
- Directory expansion/collapse
- Selection tracking
- Parent directory navigation
- Tree flattening for display
- Auto-scroll to keep selection visible

**Key Methods:**
- `LoadInitial()` - Load initial directory contents
- `Expand() / Collapse()` - Toggle directory
- `MoveUp() / MoveDown()` - Navigate selection
- `SelectedNode()` - Get current selection
- `GetFlatList()` - Flatten tree for rendering
- `NavigateToParent()` - Go up one directory
- `Refresh()` - Reload from filesystem

#### `loader.go`

Asynchronous directory loading:

**Responsibilities:**
- Read directory contents from filesystem
- Sort files and directories (dirs first, then alphabetically)
- Handle filesystem errors gracefully
- Add parent directory (..) entries

#### `finder.go`

File finding for fuzzy search:

**Responsibilities:**
- Recursively scan directories for files
- Exclude common directories (.git, node_modules, vendor, etc.)
- Return file paths relative to root
- Handle permission errors gracefully

#### `icons.go`

File type icon mapping:

**Responsibilities:**
- Map file extensions to Nerd Font icons
- Provide default icon for unknown types
- Support common file types (go, js, py, md, etc.)

### 5. Font Management (`internal/fonts/`)

Font loading and text measurement.

#### `fonts.go`

**Responsibilities:**
- Load bundled fonts (Go Mono, Go Regular)
- Provide font collections for Gio
- Text width measurement for cursor positioning
- Consistent monospace rendering

## Rendering Pipeline

Vem uses Gio UI for GPU-accelerated rendering with immediate-mode UI.

### Rendering Flow

```
User Action
    ↓
Event (app.FrameEvent)
    ↓
layout() orchestration
    ↓
Component rendering:
  - drawHeader()
  - drawPanes() (tree-based layout)
    - drawFileExplorer() (if visible)
    - drawBuffer() (text content per pane)
  - drawFuzzyFinder() (if active)
  - drawStatusBar() / drawCommandBar()
    ↓
GPU command buffer
    ↓
Frame submitted (60fps)
```

### Gio UI Integration

Gio uses immediate-mode rendering:

1. **Layout Pass**: Calculate dimensions and positions
2. **Paint Pass**: Generate GPU commands for each frame
3. **Submit**: Send commands to graphics API

**Graphics API by Platform:**
- Linux: Vulkan
- macOS: Metal
- Windows: Direct3D 11
- WebAssembly: WebGL

### Text Rendering

Text is rendered using Gio's material design components:

```go
label := material.Body1(s.theme, text)
label.Font.Typeface = "GoMono"  // Bundled monospace font
label.Color = textColor
dims := label.Layout(gtx)
```

**Font System:**
- Bundled fonts via `gofont.Collection()`
- No system font dependencies
- Custom glyph measurement for cursor positioning
- Text shaping via `text.Shaper`
- Consistent monospace rendering for code editing

### Cursor/Caret Rendering

**Normal Mode**: Block cursor (full character width)
```go
// Measure text width up to cursor position
x := measureTextWidth(gtx, prefix)
charWidth := measureTextWidth(gtx, charUnderCursor)

// Draw block
rect := image.Rect(x, 0, x+charWidth, height)
paint.Fill(gtx.Ops, cursorColor)

// Draw character on top in contrasting color
```

**Insert Mode**: Thin line cursor (2px wide)
```go
rect := image.Rect(x, 0, x+2, height)
paint.Fill(gtx.Ops, cursorColor)
```

**Blinking**: Toggled every 600ms using `time.Time` invalidation

### Pane Rendering

Panes are rendered recursively:

```go
func renderPaneNode(node *PaneNode, bounds image.Rectangle) {
    if node.IsLeaf() {
        // Render pane background (dimmed if inactive)
        // Render buffer content within pane bounds
        // Render cursor if this is the active pane
    } else {
        // Split bounds based on direction
        leftBounds, rightBounds := calculateSplitBounds(bounds, node.Direction)

        // Render separator
        drawSeparator(separatorBounds)

        // Recursively render children
        renderPaneNode(node.Left, leftBounds)
        renderPaneNode(node.Right, rightBounds)
    }
}
```

**Active Pane Dimming:**
- Active pane: Normal background (#1a1f2e)
- Inactive panes: 15% darker (#141824)
- Separator: Subtle gray (#303544)

### Search Highlighting

Search matches are highlighted in real-time:

```go
// For each line
for each match in line {
    // Draw background highlight
    if match is current match {
        highlightColor = orange  // Current match
    } else {
        highlightColor = yellow  // Other matches
    }
    drawRect(matchBounds, highlightColor)
}

// Draw text on top
drawText(lineText, normalColor)
```

## Keybinding System

The keybinding system uses a **Command/Action pattern** for robustness and extensibility.

### Architecture

```
Key Event
    ↓
handleKey()
    ↓
Priority Matching:
  1. COMMAND mode bindings (if in COMMAND mode)
  2. Global keybindings (Ctrl+T, Ctrl+F, Ctrl+U, Alt+hjkl, etc.)
  3. Ctrl+S prefix sequences (for pane commands)
  4. Mode-specific keybindings (i, j, k, l, etc.)
  5. Special handlers (counts, goto sequences, g-commands)
    ↓
Action enum
    ↓
executeAction()
    ↓
State change + UI update
```

### Priority System

1. **COMMAND Mode Priority**: When in COMMAND mode, mode-specific bindings are checked first to ensure Enter executes the command

2. **Global Keybindings**: Checked next for all modes:
   - `Ctrl+T` - Toggle explorer
   - `Ctrl+H` - Focus explorer
   - `Ctrl+L` - Focus editor
   - `Ctrl+F` - Fuzzy finder
   - `Ctrl+U` - Undo
   - `Ctrl+X` - Close pane
   - `Alt+h/j/k/l` - Pane navigation
   - `Shift+Tab` - Cycle panes
   - `Shift+Enter` - Fullscreen (NORMAL mode only)

3. **Ctrl+S Prefix Sequences**: Two-key sequences for pane management:
   - `Ctrl+S v` - Split vertical
   - `Ctrl+S h` - Split horizontal
   - `Ctrl+S =` - Equalize panes
   - `Ctrl+S o` - Zoom toggle

4. **Mode-Specific Keybindings**: Checked after globals, only active in specific modes

5. **Special Handlers**: Complex logic for counts, goto sequences (`gg`, `G`), colon commands

### Modifier Key Handling

Platform-specific quirks are handled transparently:

```go
// Track Ctrl/Shift explicitly via Press/Release events
if e.Name == key.NameCtrl {
    s.ctrlPressed = (e.State == key.Press)
    return
}

// Use tracked state, not ev.Modifiers (unreliable on some platforms)
ctrlHeld := s.ctrlPressed
shiftHeld := s.shiftPressed
altHeld := ev.Modifiers.Contain(key.ModAlt)

// Smart reset: Prevent modifiers from sticking
if s.mode == modeNormal || s.mode == modeInsert {
    if hadCtrl && s.ctrlPressed {
        s.ctrlPressed = false
    }
}
```

### Action Execution

Actions are executed in `executeAction()`:

```go
func (s *appState) executeAction(action Action, ev key.Event) {
    switch action {
    case ActionSplitVertical:
        s.handleSplitVertical()
    case ActionPaneFocusLeft:
        s.handlePaneFocusLeft()
    case ActionOpenFuzzyFinder:
        s.enterFuzzyFinder()
    case ActionUndo:
        if s.activeBuffer().Undo() {
            s.status = "Undo successful"
        }
    // ... more actions
    }
}
```

This allows the same action to behave differently based on context.

## Buffer Management

### Text Representation

Text is stored as a slice of strings (one per line):

```go
type Buffer struct {
    lines []string  // Each string is one line (no \n)
}
```

**Advantages:**
- Simple and efficient for line-based operations
- Natural mapping to display coordinates
- Easy to implement multi-line selection
- Fast line access O(1)

**Trade-offs:**
- Not ideal for very large files (millions of lines)
- No rope or piece table (yet)
- Memory usage proportional to file size

### Cursor Movement

Cursor movement handles edge cases:

```go
func (b *Buffer) MoveRight() bool {
    line := b.lines[b.cursor.Line]
    maxCol := len([]rune(line))

    if b.cursor.Col < maxCol {
        b.cursor.Col++
        return true
    }

    // Move to next line if at end
    if b.cursor.Line < len(b.lines)-1 {
        b.cursor.Line++
        b.cursor.Col = 0
        return true
    }

    return false  // At end of buffer
}
```

**Column clamping**: When moving up/down, cursor column is clamped to line length if the target line is shorter.

**Word motions**: Respect word boundaries (alphanumeric, punctuation, whitespace)

### File I/O

Files are loaded entirely into memory:

```go
func LoadFile(path string) (*Buffer, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(content), "\n")
    return NewBuffer(strings.Join(lines, "\n")), nil
}
```

**Saving:**
```go
func (bm *BufferManager) SaveActiveBuffer() error {
    buf := bm.ActiveBuffer()
    content := strings.Join(buf.Lines(), "\n")
    return os.WriteFile(buf.FilePath(), []byte(content), 0644)
}
```

## Pane System

### Binary Tree Structure

Panes are organized in a binary tree:

```
        [Split H]
       /          \
   [Pane 1]    [Split V]
               /        \
           [Pane 2]   [Pane 3]
```

**Node Types:**
- **Leaf nodes**: Contain actual panes with buffers
- **Internal nodes**: Represent splits (H = horizontal divider, V = vertical divider)

### Split Terminology

- **Vertical Split**: Creates horizontal divider (top/bottom panes)
- **Horizontal Split**: Creates vertical divider (left/right panes)

This matches Vim's terminology.

### Layout Algorithm

```go
func calculateLayout(node *PaneNode, bounds image.Rectangle) {
    if node.IsLeaf() {
        node.Pane.Bounds = bounds
        return
    }

    if node.Direction == SplitVertical {
        // Split horizontally (top/bottom)
        midY := bounds.Min.Y + bounds.Dy()/2
        topBounds := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, midY)
        bottomBounds := image.Rect(bounds.Min.X, midY, bounds.Max.X, bounds.Max.Y)
        calculateLayout(node.Left, topBounds)
        calculateLayout(node.Right, bottomBounds)
    } else {
        // Split vertically (left/right)
        midX := bounds.Min.X + bounds.Dx()/2
        leftBounds := image.Rect(bounds.Min.X, bounds.Min.Y, midX, bounds.Max.Y)
        rightBounds := image.Rect(midX, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
        calculateLayout(node.Left, leftBounds)
        calculateLayout(node.Right, rightBounds)
    }
}
```

All splits are 50/50.

### Pane Navigation

Navigation uses geometric calculation:

```go
func FindPaneInDirection(current *Pane, direction Direction, allPanes []*Pane) *Pane {
    currentCenter := center(current.Bounds)

    var candidates []*Pane
    for _, pane := range allPanes {
        if pane == current {
            continue
        }

        candidateCenter := center(pane.Bounds)

        // Check if pane is in the right direction
        if direction == Left && candidateCenter.X < currentCenter.X {
            candidates = append(candidates, pane)
        }
        // ... similar for Right, Up, Down
    }

    // Return closest candidate
    return findClosest(candidates, currentCenter, direction)
}
```

### Zoom Mode

Zoom temporarily hides all panes except the active one:

```go
func (pm *PaneManager) ToggleZoom() {
    if pm.zoomed != nil {
        pm.zoomed = nil  // Restore all panes
    } else {
        pm.zoomed = pm.activePane  // Show only active pane
    }
}
```

During rendering, if zoomed is set, only that pane gets the full window bounds.

## Fuzzy Finder

### Fuzzy Matching Algorithm

The fuzzy finder uses a scoring algorithm to rank file paths:

```go
func FuzzyScore(pattern, target string) (int, []int) {
    // Find all sequential character matches
    // Score based on:
    // - Sequential matches: +10 each
    // - Consecutive matches: +15 bonus
    // - Word boundary matches: +5 bonus
    // - Start of string: +10 bonus
    // - Case match: +2 bonus
    // - Shorter paths: bonus
    // - Gaps between matches: penalty
}
```

**Example:**
- Pattern: `bufgo`
- Target: `internal/editor/buffer.go`
- Matches: **b**uffer.**g****o**
- Score: High (word boundary + consecutive + short path)

### File Discovery

Files are discovered recursively:

```go
func FindFiles(root string, exclude []string) []string {
    var files []string
    filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        // Skip excluded directories
        if info.IsDir() && isExcluded(path, exclude) {
            return filepath.SkipDir
        }

        // Add files
        if !info.IsDir() {
            relPath := relativePath(path, root)
            files = append(files, relPath)
        }
        return nil
    })
    return files
}
```

**Excluded by default:**
- `.git`, `.gocache`, `node_modules`, `vendor`
- `dist`, `build`, `target`
- Hidden directories (starting with `.`)

### UI Implementation

Fuzzy finder renders as an overlay:

```go
// Semi-transparent background overlay
drawRect(fullScreen, rgba(0, 0, 0, 0.5))

// Centered finder box with border
drawRect(finderBox, backgroundColor)
drawBorder(finderBox, borderColor)

// Input field showing pattern
drawText(fuzzyPattern, cursorPos)

// Results list with highlighted matches
for i, result := range fuzzyResults {
    if i == fuzzySelected {
        drawRect(resultBounds, selectionColor)  // Highlight selected
    }
    drawTextWithHighlights(result.Path, result.MatchIndices)
}
```

## Search System

### Search Implementation

Search is case-insensitive substring matching:

```go
func (b *Buffer) Search(pattern string) []SearchMatch {
    var matches []SearchMatch
    lowerPattern := strings.ToLower(pattern)

    for lineIdx, line := range b.lines {
        lowerLine := strings.ToLower(line)
        startPos := 0

        for {
            pos := strings.Index(lowerLine[startPos:], lowerPattern)
            if pos == -1 {
                break
            }

            matches = append(matches, SearchMatch{
                Line:   lineIdx,
                StartCol: startPos + pos,
                EndCol:   startPos + pos + len(pattern),
            })

            startPos += pos + 1
        }
    }

    return matches
}
```

### Match Navigation

Navigation wraps around:

```go
func (s *appState) jumpToNextMatch() {
    if len(s.searchMatches) == 0 {
        return
    }

    s.searchCurrentIdx = (s.searchCurrentIdx + 1) % len(s.searchMatches)
    match := s.searchMatches[s.searchCurrentIdx]

    // Move cursor to match
    buf := s.activeBuffer()
    buf.SetCursor(match.Line, match.StartCol)

    s.status = fmt.Sprintf("/%s [%d/%d]", s.searchPattern,
                          s.searchCurrentIdx+1, len(s.searchMatches))
}
```

### Highlighting

All matches are highlighted:
- **Yellow background**: All search matches
- **Orange background**: Current match (where cursor is)

## Undo System

### Implementation

Undo uses a simple stack of buffer states:

```go
type Buffer struct {
    lines     []string
    cursor    Cursor
    undoStack []UndoState
}

type UndoState struct {
    Lines  []string      // Snapshot of all lines
    Cursor Cursor        // Cursor position at save time
}

func (b *Buffer) SaveUndoState() {
    // Limit stack size to 100
    if len(b.undoStack) >= 100 {
        b.undoStack = b.undoStack[1:]
    }

    // Deep copy lines
    linesCopy := make([]string, len(b.lines))
    copy(linesCopy, b.lines)

    b.undoStack = append(b.undoStack, UndoState{
        Lines:  linesCopy,
        Cursor: b.cursor,
    })
}

func (b *Buffer) Undo() bool {
    if len(b.undoStack) == 0 {
        return false
    }

    // Pop last state
    state := b.undoStack[len(b.undoStack)-1]
    b.undoStack = b.undoStack[:len(b.undoStack)-1]

    // Restore state
    b.lines = state.Lines
    b.cursor = state.Cursor
    b.modified = true

    return true
}
```

**Undo triggers:**
- Text insertion
- Text deletion
- Line deletion
- Paste operations

## Platform Abstraction

Gio UI provides platform abstraction:

- **Window Management**: Create/destroy, resize, fullscreen
- **Input Events**: Keyboard, mouse, touch
- **Graphics**: Vulkan/Metal/Direct3D/WebGL backend selection
- **Text Input**: Native IME integration

Vem adds additional abstractions:

- **Modifier Key Tracking**: Compensate for platform reporting inconsistencies
- **File Paths**: Use `filepath` package for cross-platform path handling
- **Line Endings**: Normalize to `\n` internally, convert on save (future)

## Design Decisions

### Why Go?

- **Cross-compilation**: Single codebase for all platforms
- **Performance**: Near-native performance, efficient memory use
- **Simplicity**: Easy to understand, minimal cognitive overhead
- **Concurrency**: Goroutines for async operations (future)
- **Tooling**: Excellent tooling (go fmt, go test, go build)
- **Static binary**: Single executable, no runtime dependencies

### Why Gio UI?

- **GPU-Accelerated**: Smooth 60fps rendering
- **Cross-Platform**: Linux, macOS, Windows, WebAssembly
- **No C Dependencies**: Pure Go (except platform graphics APIs)
- **Immediate Mode**: Simple to reason about, no retained state
- **Small Binary**: ~12MB statically linked binary
- **Active Development**: Regular updates and responsive community

### Why Modal Editing?

- **Vim Familiarity**: Large user base familiar with Vim
- **Efficiency**: Separates navigation from editing
- **Keyboard Focus**: Minimal mouse usage required
- **Extensibility**: Clear separation of concerns
- **Power**: Complex operations with few keystrokes

### Why Bundled Fonts?

- **Zero Setup**: No user configuration required
- **Consistency**: Same appearance across platforms
- **Offline**: Works without internet or package managers
- **Licensing**: GoFont family (BSD license, compatible with GPLv2)

### Why Binary Tree for Panes?

- **Simplicity**: Easy to understand and implement
- **Recursive Layout**: Natural fit for split layout
- **Flexibility**: Supports arbitrary nesting
- **Vim-like**: Matches Vim's split model

## Performance Considerations

### Current Performance

- **Rendering**: 60fps on most hardware (GPU-accelerated)
- **Text Editing**: Sub-millisecond for typical operations
- **File Loading**: Limited by disk I/O (~1GB/s for typical SSDs)
- **File Tree**: Handles directories with thousands of files
- **Fuzzy Finder**: Searches tens of thousands of files in <100ms
- **Undo**: O(n) space where n is number of undo states (max 100)

### Known Limitations

- **Large Files**: Loading multi-megabyte files into memory
- **Deep Undo**: Each undo state duplicates entire buffer
- **File Tree**: Entire tree loaded into memory
- **No Syntax Highlighting**: Pure text rendering only

### Future Optimizations

- **Large Files**: Rope data structure, lazy loading, memory mapping
- **Undo**: Diff-based storage instead of full snapshots
- **Syntax Highlighting**: Incremental tree-sitter parsing
- **File Tree**: Virtual scrolling for large directories
- **Fuzzy Finder**: Incremental search with result caching

## Testing Strategy

### Current Testing

- **Unit Tests**: Buffer operations (`buffer_test.go`)
- **Manual Testing**: Interactive testing of UI features

### Planned Testing

- **Integration Tests**: Full editor workflows
- **Fuzzing**: Random input testing for buffer operations
- **Snapshot Tests**: UI rendering regression tests
- **Performance Tests**: Benchmark large files and operations
- **Cross-platform**: Automated testing on Linux/macOS/Windows

## Code Organization

```
internal/
├── appcore/              # UI and event handling
│   ├── app.go           # Main application state
│   ├── keybindings.go   # Keybinding system
│   ├── pane_actions.go  # Pane management actions
│   ├── pane_rendering.go # Pane rendering
│   └── fuzzy.go         # Fuzzy finder
├── editor/               # Text editing logic
│   ├── buffer.go        # Buffer abstraction
│   ├── buffer_test.go   # Buffer tests
│   └── buffer_manager.go # Multi-buffer management
├── filesystem/           # File tree and operations
│   ├── tree.go          # Tree data structure
│   ├── loader.go        # Directory loading
│   ├── finder.go        # File finding for fuzzy search
│   └── icons.go         # File type icons
├── panes/                # Pane management
│   ├── manager.go       # Pane tree manager
│   ├── layout.go        # Layout calculation
│   ├── navigation.go    # Geometric navigation
│   └── pane.go          # Pane abstraction
└── fonts/                # Font management
    └── fonts.go         # Font loading and rendering
```

**Design Principles:**
- **Separation of Concerns**: UI, editing, filesystem, panes are independent
- **Testability**: Core logic has no UI dependencies
- **Extensibility**: Clear interfaces for future plugins
- **Immutability**: Minimize mutable state where possible

## See Also

- [Keybindings Reference](keybindings.md) - Complete keybinding documentation
- [Pane Splitting Guide](pane-splitting.md) - Detailed pane management
- [Tutorial](tutorial.md) - User-facing getting started guide
- [Installation Guide](installation.md) - Platform-specific installation
