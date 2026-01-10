# app.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/app.go`
**Lines:** ~4000+
**Purpose:** Main application state, event loop, rendering, and command handling

## Overview

This is the largest file in the codebase, containing the core application logic including:
- Application state management (`appState` struct)
- Main event loop and window handling
- Buffer and pane rendering
- Input event handling
- Command-line execution (`:commands`)
- Mode management

## Code Blocks

### Lines 1-40: Package and Imports
```go
package appcore
import (
    // Gio UI framework imports
    // Internal package imports
    // Clipboard support
)
```

### Lines 42-74: Mode and Type Definitions

| Lines | Definition | Purpose |
|-------|------------|---------|
| 42 | `type mode string` | Editing mode type |
| 44-50 | `visualModeType` | None, Char, Line |
| 52-56 | `SearchMatch` | Search result location |
| 58-62 | `FuzzyMatch` | Fuzzy finder result |
| 64-74 | Mode constants | All 9 editing modes |

### Lines 76-96: Constants and Colors

| Lines | Constant | Value |
|-------|----------|-------|
| 77 | `caretBlinkInterval` | 600ms |
| 78 | `tabWidth` | 4 |
| 82 | `highlightColor` | Current line highlight |
| 83 | `selectionColor` | Visual selection |
| 84 | `background` | Editor background (#1a1f2e) |
| 85 | `statusBg` | Status bar background |
| 94-95 | `activePaneBg`, `inactivePaneBg` | Pane colors |

### Lines 98-221: appState Struct

The central state struct containing all application state:

| Lines | Category | Fields |
|-------|----------|--------|
| 99-101 | Core | `theme`, `bufferMgr`, `paneManager` |
| 102-105 | State | `fileTree`, `mode`, `status`, `lastKey` |
| 106-111 | Pending ops | `pendingCount`, `pendingGoto`, etc. |
| 112-124 | Visual/Edit | `visualMode`, `skipNextEdit`, etc. |
| 128-139 | Explorer | `explorerVisible`, `explorerWidth`, etc. |
| 141-144 | Search | `searchPattern`, `searchMatches`, etc. |
| 146-151 | Fuzzy finder | `fuzzyFinderInput`, `fuzzyFinderMatches`, etc. |
| 153-161 | Modifiers | `ctrlPressed`, `shiftPressed`, time tracking |
| 163-169 | Viewport | `viewportTopLine`, `scrollOffsetLines` |
| 171-175 | Syntax/Terminal | `syntaxHighlighters`, `syntaxEnabled`, `currentTheme`, `terminals` |
| 181-221 | LSP | All LSP-related state fields |
| 218-222 | Diagnostics List | `diagnosticsListActive`, `diagnosticsListItems`, `diagnosticsListIndex` |
| 223-227 | Buffer Completion | `bufferCompletionActive`, `bufferCompletionItems`, `bufferCompletionIndex` |
| 229-231 | Git Branch Cache | `gitBranch`, `gitBranchLastCheck` |

### Lines 223-268: Entry Point and Event Loop

#### Run() (Line 223)
```go
func Run(w *app.Window, filePaths []string) error {
    state := newAppState(filePaths)
    return state.run(w)
}
```

#### run() (Lines 228-251)
Main event loop handling:
- `app.DestroyEvent` - Exit
- `app.ConfigEvent` - Window mode/resize
- `app.FrameEvent` - Render frame

#### cleanup() (Lines 254-268)
Shutdown logic:
- Closes all terminals
- Stops all LSP servers

### Lines 297-408: State Initialization

`newAppState()` creates the initial application state:
- Loads fonts (JetBrains Mono Nerd Font)
- Creates buffer manager from file paths
- Initializes pane manager
- Initializes file tree from working directory
- Sets up LSP manager and callbacks

### Lines 456-517: Buffer and Highlighter Helpers

| Function | Lines | Purpose |
|----------|-------|---------|
| `activeBuffer()` | 456-468 | Get buffer for active pane |
| `getOrCreateHighlighter()` | 470-507 | Get/create syntax highlighter |
| `invalidateSyntaxCache()` | 509-517 | Clear syntax cache on edit |

### Lines 545-595: Layout Function

Main rendering layout:
```
Vertical Flex:
├── Rigid: drawHeader()
├── Rigid: drawTabBar()          ← NEW: Tab bar for open buffers
├── Flexed(1):
│   └── Horizontal (if explorer visible):
│       ├── Rigid: drawFileExplorer()
│       └── Flexed(1): drawPanes()
│   └── drawPanes() (if explorer hidden)
└── Rigid: drawCommandBar() or drawStatusBar()

Overlay: drawFuzzyFinder() (if active)
```

### Lines 761-854: Tab Bar (NEW)

#### drawTabBar()
Draws a tab bar showing all open buffers:
- Only renders when multiple buffers open
- Shows filename for each buffer
- Modified buffers show `[+]` suffix
- Active buffer highlighted with different colors
- Colors: dark background, lighter active tab

### Lines 587-740: Event Handling

`handleEvents()` processes:
- `key.FocusEvent` - Focus changes
- `key.Event` - Keyboard input
- `key.EditEvent` - Text input

**Key processing phases:**
1. Handle modifier events
2. Sync modifier state
3. Call `handleKey()`
4. Smart modifier reset based on mode

### Lines 755-908: Buffer Rendering

#### drawBuffer() (755-826)
- Calculates viewport and lines per page
- Ensures cursor visibility
- Renders line list with syntax highlighting
- Draws overlays (completion menu, hover, references)

#### drawBufferLine() (828-908)
- Renders gutter (line numbers: 4-digit format)
- Tokenizes line via syntax highlighter
- Draws backgrounds (selection, cursor line)
- Draws search highlights
- Draws cursor

### Lines 1111-1240: Status Bar (Enhanced)

`drawStatusBar()` renders:
```
MODE {mode} | FILE {name}[+][RO][Large] | CURSOR {line}:{col} [| PANE x/y] [| FULLSCREEN] [| ZOOMED] [| GIT branch] [| LSP] [| E:n W:n] | {status} [| HINT: ...]
```

New status components:
- `[Large]` - Indicates file > 5MB
- `GIT branch` - Current git branch (cached 5 seconds)
- `LSP` - Shows when LSP server active for file
- `E:n W:n` - Error/warning counts from diagnostics

#### Helper Functions (Lines 1225-1290)
| Function | Purpose |
|----------|---------|
| `getGitBranch()` | Returns cached git branch name |
| `getDiagnosticCounts()` | Returns error/warning counts for current file |
| `getLSPStatus()` | Returns "LSP" if server active for file |

### Lines 1110-1232: Overlays

| Function | Lines | Purpose |
|----------|-------|---------|
| `drawFuzzyFinder()` | 1110-1212 | Fuzzy finder modal |
| `drawCommandBar()` | 1214-1232 | Command input bar |
| `drawFileExplorer()` | 1234-1341 | File tree sidebar |

### Lines 1343-1441: Key Handling

`handleKey()` processing phases:
1. Check pane resize mode (Ctrl+Shift+R)
2. Handle terminal mode
3. Handle file operations
4. Check Ctrl+S prefix for pane commands
5. Match global keybindings
6. Match mode-specific keybindings
7. Handle mode special cases

### Lines 2521-2594: Command Execution

`executeCommandLine()` parses and executes commands:

| Command | Lines | Action |
|---------|-------|--------|
| `:q` `:quit` | 2542-2545 | Quit pane/buffer |
| `:w` `:write` | 2550-2553 | Save file |
| `:wq` | 2552-2553 | Save and quit |
| `:e` `:edit` | 2554-2555 | Open file |
| `:bn` `:bp` | 2556-2567 | Buffer navigation |
| `:bd` | 2568-2571 | Delete buffer |
| `:ls` `:buffers` | 2572-2573 | List buffers |
| `:ex` `:explore` | 2574-2575 | Toggle explorer |
| `:cd` `:pwd` | 2576-2579 | Directory commands |
| `:term` | 2580-2581 | Open terminal |
| `:help` | 2582-2583 | Show help |
| `:install` | 2584-2585 | Install LSP server |
| `:lspauto` | 2586-2587 | Toggle LSP auto-install |
| `:theme <name>` | 2698-2699 | Switch color theme (NEW) |
| `:themes` | 2700-2701 | List available themes (NEW) |
| `:diagnostics` `:diag` | 2804-2805 | Show diagnostics list (NEW) |

### Lines 2597-2660: Quit Command Logic

`handleQuitCommand()`:
- Checks for unsaved changes
- Closes terminal if buffer has one
- Handles multiple panes vs single pane
- Cleans up LSP

## Known Issues / Potential Bugs

1. ~~**No `u` for undo in Normal mode**~~ FIXED - `u` now bound in Normal mode

2. ~~**Large file memory usage**~~ MITIGATED - Large files (>5MB) show warning, >50MB rejected

3. **Platform modifier quirks** (Lines 332-336)
   - `ev.Modifiers` always empty on some platforms
   - Workaround uses tracked state

## New Features Added

1. **Tab Bar** (Lines 761-854) - Shows open buffers when multiple present
2. **Theme Switching** - `:theme <name>` and `:themes` commands
3. **Enhanced Status Line** - Git branch, LSP status, diagnostics count
4. **Diagnostics List** - `:diagnostics` command for navigating errors/warnings
5. **Buffer Word Completion** - Ctrl+N triggers word completion from buffer
6. **Large File Handling** - Warning for >5MB, rejection for >50MB
7. **Automatic Completion** - `maybeTriggerAutoCompletion()` falls back to buffer completion when LSP unavailable
8. **Improved Viewport Scrolling** - Uses `cachedLineHeight` for accurate scroll calculations
9. **Command Path Completion** - Tab completion for `:cd`, `:e`, `:w` paths (see `command_completion.go`)
10. **Format on Save** - `:formatonsave on/off` toggle and `:format` manual command for LSP formatting
11. **Soft Wrap** - `:wrap` toggle for dynamic text wrapping to viewport width
12. **Tab Key Fix** - Fixed `skipNextEdit` timing race to prevent double key presses

### New Commands
| Command | Description |
|---------|-------------|
| `:formatonsave on/off` | Toggle LSP format on save |
| `:format` or `:fmt` | Manual LSP format current buffer |
| `:wrap` or `:wrap on/off` | Toggle soft wrap (no arg toggles) |

### New State Fields
```go
cachedLineHeight      int      // Cached line height for accurate scrolling
cmdCompletionActive   bool     // Command mode path completion state
cmdCompletionItems    []string
cmdCompletionIndex    int
cmdCompletionPrefix   string
commandHistory        []string // List of executed commands (in-memory)
commandHistoryIndex   int      // Current position in history (-1 = new input)
commandHistorySaved   string   // Saves current input when navigating history
formatOnSave          bool     // Auto-format on save when LSP available
softWrapEnabled       bool     // Soft wrap long lines to viewport width
```

### Command History Functions
- `commandHistoryUp()` - Navigate to previous command (Up arrow in command mode)
- `commandHistoryDown()` - Navigate to next command (Down arrow in command mode)
- History is saved when commands are executed via `executeCommandLine()`
- Duplicate consecutive commands are not saved

### Soft Wrap Implementation
- `wrapLine()` function splits lines into segments based on visual width
- `drawBufferLineWrapped()` renders wrapped segments with continuation indicators (`↪`)
- Line numbers shown only on first segment
- Cursor properly positioned within wrapped segments

### Cursor Positioning Fix
- **Problem**: Cursor was visually offset from actual character position due to measuring prefix as one string while rendering tokens separately (font kerning differences)
- **Solution**: Token-based cursor positioning that accumulates actual rendered token widths
- `drawCursorAtX()` - New function that draws cursor at pre-calculated X position
- `drawBufferLine()` - Tracks cursor X by summing token widths during rendering
- `drawBufferLineWrapped()` - Same token-based tracking for wrapped segments

### Syntax Highlighter Cache Fix
- **Problem**: Syntax highlighting disappeared when opening files from explorer
- **Root Cause**: When closing buffer 0 (empty `[No Name]` buffer), the highlighter cache became stale. The new file's highlighter was stored at index 1, but after closing buffer 0, the file shifted to index 0 while the cache still pointed to index 1.
- **Solution**: Added `updateHighlighterCacheOnBufferClose()` function (lines 556-572) that:
  1. Deletes the highlighter for the closed buffer
  2. Shifts all higher-indexed highlighters down by one
- This function is now called before every `CloseBuffer()` call throughout the codebase

### Helper Functions (Lines 546-572)
| Function | Purpose |
|----------|---------|
| `invalidateSyntaxCache()` | Clears syntax cache for active buffer on content change |
| `updateHighlighterCacheOnBufferClose()` | Updates highlighter indices when a buffer is closed |

---
*Last Updated: After syntax highlighter cache fix*
