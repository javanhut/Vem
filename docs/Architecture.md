# ProjectVem Architecture

Technical architecture and design decisions for ProjectVem.

## Table of Contents

- [Overview](#overview)
- [Core Components](#core-components)
- [Rendering Pipeline](#rendering-pipeline)
- [Keybinding System](#keybinding-system)
- [Buffer Management](#buffer-management)
- [File System Integration](#file-system-integration)
- [Platform Abstraction](#platform-abstraction)
- [Design Decisions](#design-decisions)

## Overview

ProjectVem is built as a cross-platform text editor using Go and Gio UI. The architecture prioritizes:

1. **Cross-platform compatibility** - Single codebase for Linux, macOS, Windows, WebAssembly
2. **Zero external dependencies** - Bundles fonts and assets, no system packages required
3. **GPU acceleration** - Smooth rendering via Vulkan/Metal/Direct3D
4. **Modal editing** - Vim-like editing paradigm with robust state management
5. **Extensibility** - Plugin system planned for Lua/Python/Carrion scripts

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
    fileTree     *filesystem.FileTree
    mode         mode
    status       string
    
    // Explorer state
    explorerVisible bool
    explorerFocused bool
    
    // Modifier tracking
    ctrlPressed  bool
    shiftPressed bool
}
```

**Responsibilities:**
- Window creation and lifecycle management
- Event loop (keyboard, mouse, window events)
- Modal state machine (NORMAL, INSERT, VISUAL, DELETE, COMMAND, EXPLORER)
- Rendering pipeline orchestration
- Status bar and UI chrome
- Caret blinking animation

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

**Key Methods:**
- `matchGlobalKeybinding()` - Check global shortcuts
- `matchModeKeybinding()` - Check mode-specific bindings
- `modifiersMatch()` - Handle platform-specific modifier detection
- `executeAction()` - Dispatch actions to implementation

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
}

type Cursor struct {
    Line int
    Col  int
}
```

**Responsibilities:**
- In-memory text storage as slice of lines
- UTF-8 aware cursor positioning
- Text mutation operations (insert, delete, modify)
- Line-based navigation
- Undo/redo support (planned)

**Key Methods:**
- `InsertText(text string)` - Insert text at cursor
- `DeleteBackward()` - Backspace operation
- `DeleteForward()` - Delete key operation
- `MoveLeft/Right/Up/Down()` - Cursor movement
- `JumpLineStart/End()` - Line boundary jumps
- `DeleteLines(start, end)` - Multi-line deletion

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

**Key Methods:**
- `OpenFile(path)` - Load file into new buffer
- `SaveActiveBuffer()` - Save current buffer
- `NextBuffer() / PrevBuffer()` - Switch buffers
- `CloseActiveBuffer(force)` - Close with modified check

### 3. File System Integration (`internal/filesystem/`)

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

**Key Methods:**
- `LoadInitial()` - Load initial directory contents
- `Expand() / Collapse()` - Toggle directory
- `MoveUp() / MoveDown()` - Navigate selection
- `SelectedNode()` - Get current selection
- `GetFlatList()` - Flatten tree for rendering

#### `loader.go`

Asynchronous directory loading:

**Responsibilities:**
- Read directory contents from filesystem
- Sort files and directories
- Handle filesystem errors gracefully
- Add parent directory (..) entries

## Rendering Pipeline

ProjectVem uses Gio UI for GPU-accelerated rendering.

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
  - drawFileExplorer() (if visible)
  - drawBuffer() (text content)
  - drawStatusBar() / drawCommandBar()
    ↓
GPU command buffer
    ↓
Frame submitted
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

### Cursor/Caret Rendering

**Normal Mode**: Block cursor
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
  2. Global keybindings (Ctrl+T, Shift+Enter, etc.)
  3. Mode-specific keybindings (i, j, k, l, etc.)
  4. Special handlers (counts, goto sequences)
    ↓
Action enum
    ↓
executeAction()
    ↓
State change + UI update
```

### Priority System

1. **COMMAND Mode Priority**: When in COMMAND mode, mode-specific bindings are checked first to ensure Enter executes the command instead of triggering global Shift+Enter

2. **Global Keybindings**: Checked first for all other modes, ensuring Ctrl+T, Ctrl+H, Ctrl+L, Shift+Enter always work

3. **Mode-Specific Keybindings**: Checked after globals, only active in specific modes

4. **Special Handlers**: Complex logic for counts, goto sequences, colon commands

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

// Smart reset: Prevent modifiers from sticking
if s.mode == modeNormal || s.mode == modeInsert || s.mode == modeCommand {
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
    case ActionToggleExplorer:
        s.toggleExplorer()
    case ActionMoveUp:
        if s.mode == modeExplorer {
            s.fileTree.MoveUp()
        } else {
            s.activeBuffer().MoveUp()
        }
    // ... more actions
    }
}
```

This allows the same action to behave differently based on mode context.

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

**Trade-offs:**
- Not ideal for very large files (millions of lines)
- No rope or piece table (yet)
- Future: Consider rope data structure for large files

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

**Future improvements:**
- Lazy loading for large files
- Memory-mapped I/O
- Incremental saving

## File System Integration

### Directory Structure

File tree is loaded on-demand:

1. Load root directory contents
2. Display directories and files
3. On expand: Load subdirectory contents
4. Cache loaded directories

### File Operations

- **Open File**: Load into new buffer, switch to editor
- **Navigate Directory**: Change root, reload tree
- **Refresh**: Re-read current directory from disk

## Platform Abstraction

Gio UI provides platform abstraction:

- **Window Management**: Create/destroy, resize, fullscreen
- **Input Events**: Keyboard, mouse, touch
- **Graphics**: Vulkan/Metal/Direct3D/WebGL backend selection
- **Text Input**: Native IME integration

ProjectVem adds additional abstractions:

- **Modifier Key Tracking**: Compensate for platform reporting inconsistencies
- **File Paths**: Use `filepath` package for cross-platform path handling
- **Line Endings**: Normalize to `\n` internally, convert on save

## Design Decisions

### Why Go?

- **Cross-compilation**: Single codebase for all platforms
- **Performance**: Near-native performance, efficient memory use
- **Simplicity**: Easy to understand, minimal cognitive overhead
- **Concurrency**: Goroutines for async operations (future)
- **Tooling**: Excellent tooling (go fmt, go test, go build)

### Why Gio UI?

- **GPU-Accelerated**: Smooth 60fps rendering
- **Cross-Platform**: Linux, macOS, Windows, WebAssembly
- **No C Dependencies**: Pure Go (except platform graphics APIs)
- **Immediate Mode**: Simple to reason about, no retained state
- **Small Binary**: ~10MB statically linked binary

### Why Modal Editing?

- **Vim Familiarity**: Large user base familiar with Vim
- **Efficiency**: Separates navigation from editing
- **Keyboard Focus**: Minimal mouse usage required
- **Extensibility**: Clear separation of concerns

### Why Bundled Fonts?

- **Zero Setup**: No user configuration required
- **Consistency**: Same appearance across platforms
- **Offline**: Works without internet or package managers
- **Licensing**: GoFont family (Apache 2.0 compatible with GPLv2)

## Future Architecture

### Plugin System (Phase 3)

Planned plugin architecture:

```
Plugin Host (Go)
    ↓
Language Runtimes:
  - Lua VM
  - Python interpreter
  - Carrion VM
    ↓
Plugin API:
  - Keybinding registration
  - Command registration
  - Event hooks
  - Buffer manipulation
```

**Sandboxing:**
- Restrict filesystem access
- Limit network access
- CPU/memory quotas
- Permission system

### LSP Integration (Phase 3)

Language Server Protocol client:

```
Editor ←→ LSP Client ←→ Language Server
```

**Features:**
- Multi-server support (one per language)
- Completion, diagnostics, hover
- Go-to-definition, find-references
- Code actions, formatting

### Syntax Highlighting (Phase 3)

Treesitter-style syntax pipeline:

```
Source Code
    ↓
Parser (per language)
    ↓
Syntax Tree
    ↓
Query (highlighting patterns)
    ↓
Highlighted Ranges
    ↓
Rendering (with colors)
```

## Performance Considerations

### Current Performance

- **Rendering**: 60fps on most hardware (GPU-accelerated)
- **Text Editing**: Sub-millisecond for typical operations
- **File Loading**: Limited by disk I/O
- **File Tree**: Limited by directory size

### Planned Optimizations

- **Large Files**: Rope data structure, lazy loading
- **Syntax Highlighting**: Incremental parsing
- **File Tree**: Virtual scrolling for large directories
- **Undo/Redo**: Efficient diff-based storage

## Testing Strategy

### Current Testing

- **Unit Tests**: Buffer operations (`buffer_test.go`)
- **Manual Testing**: Interactive testing of UI features

### Planned Testing

- **Integration Tests**: Full editor workflows
- **Fuzzing**: Random input testing for buffer operations
- **Snapshot Tests**: UI rendering regression tests
- **Performance Tests**: Benchmark large files and operations

## Code Organization

```
internal/
├── appcore/          # UI and event handling
│   ├── app.go       # Main application state
│   └── keybindings.go  # Keybinding system
├── editor/           # Text editing logic
│   ├── buffer.go    # Buffer abstraction
│   ├── buffer_test.go  # Buffer tests
│   └── buffer_manager.go  # Multi-buffer management
└── filesystem/       # File tree
    ├── tree.go      # Tree data structure
    └── loader.go    # Directory loading
```

**Design Principles:**
- **Separation of Concerns**: UI, editing, filesystem are independent
- **Testability**: Core logic has no UI dependencies
- **Extensibility**: Clear interfaces for future plugins

## See Also

- [Keybindings Reference](keybindings.md) - Complete keybinding documentation
- [Tutorial](tutorial.md) - User-facing feature documentation
- [ROADMAP.md](../ROADMAP.md) - Development phases and milestones
- [PROJECT_DESCRIPTION.md](../PROJECT_DESCRIPTION.md) - Project vision and goals
