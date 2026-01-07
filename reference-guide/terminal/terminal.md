# terminal.go

**Path:** `/home/javanhut/Development/Vem/internal/terminal/terminal.go`
**Lines:** ~400+
**Purpose:** VT100 terminal emulator with PTY support

## Overview

Implements integrated terminal with:
- VT100/ANSI escape sequence support via vt10x library
- Cross-platform PTY (Unix) / ConPTY (Windows)
- Screen buffer management
- Input/output goroutines

## Code Blocks

### Lines 1-15: Package and Imports

```go
package terminal
import (
    "context"
    "fmt"
    "io"
    "os"
    "os/exec"
    "runtime"
    "sync"
    "time"
    "gioui.org/app"
    "github.com/hinshun/vt10x"
)
```

### Lines 17-70: Type Definitions

#### ConPtyIO Interface (Lines 17-24)
```go
type ConPtyIO interface {
    Read(p []byte) (int, error)
    Write(p []byte) (int, error)
    Close() error
    Resize(width, height int) error
}
```

#### Terminal Struct (Lines 25-70)
```go
type Terminal struct {
    // PTY and process
    pty    *os.File   // Unix PTY master
    conpty ConPtyIO   // Windows ConPTY
    cmd    *exec.Cmd  // Shell process

    // VT100 emulator
    vt vt10x.Terminal

    // Screen buffer
    screen *ScreenBuffer

    // Dimensions
    width, height int

    // Shell info
    shell, workingDir string
    args, env []string

    // Lifecycle
    ctx, cancel context.Context/CancelFunc
    running bool
    mu sync.RWMutex

    // Channels
    inputChan  chan []byte    // UI -> PTY
    updateChan chan struct{}  // Screen updates

    // Window invalidation
    window *app.Window

    // Exit callback
    onExit func()
}
```

#### Config Struct (Lines 72-82)
```go
type Config struct {
    Width, Height int
    Shell         string
    Args          []string
    WorkingDir    string
    Env           []string
    Window        *app.Window
    OnExit        func()
}
```

### Lines 84-126: Constructor

#### NewTerminal (Lines 84-126)
Creates terminal with config:
1. Validates dimensions
2. Sets defaults for shell, args, workingDir
3. Creates context for lifecycle
4. Creates screen buffer
5. Creates VT10x emulator with size

### Lines 128-153: Startup

#### Start (Lines 128-153)
Starts terminal:
1. Checks if already running
2. Calls `startPTY()` (platform-specific)
3. Starts read and write goroutines

### Lines 155-214: Read Loop

#### readLoop (Lines 155-214)
Reads PTY output:
1. Reads from PTY/ConPTY in 4KB chunks
2. Writes to vt10x parser
3. Updates screen buffer from vt10x state
4. Signals update via channel
5. Invalidates window for redraw
6. Handles EOF and timeout errors

### Lines 216-240: Write Loop

#### writeLoop (Lines 216-240)
Sends input to PTY:
1. Waits on context or input channel
2. Writes to PTY/ConPTY
3. Handles errors

### Lines 242-300: Public API

| Function | Lines | Purpose |
|----------|-------|---------|
| `Write()` | 242-254 | Send input to PTY |
| `GetScreen()` | 256-259 | Get screen buffer |
| `Close()` | 261-300 | Stop terminal |

#### Close Details (Lines 261-300)
1. Sets running = false
2. Cancels context
3. Closes ConPTY (Windows)
4. Closes PTY (Unix)
5. Kills process if still running
6. Waits for goroutines (2 second timeout)

## Known Issues / Potential Bugs

1. **Line 177: 100ms read deadline**
   - Fixed timeout may cause performance issues
   - Consider adaptive timeout

2. **Line 298: 2 second wait timeout**
   - May leave zombie processes
   - Consider more aggressive cleanup

## Dead/Unused Code

None identified.

## Platform-Specific Files

- `pty_unix.go` - Unix PTY implementation
- `pty_windows.go` - Windows ConPTY implementation
- `conpty_wrapper.go` - Windows ConPTY wrapper

## Integration Points

- Created in `handleOpenTerminal()` in app.go
- Buffer stored in terminal map keyed by buffer index
- Input handled in `handleTerminalKey()` and `handleTerminalEdit()`
- Rendering in terminal section of app.go

## Data Flow

```
User Input → inputChan → writeLoop → PTY
                                      ↓
                              Shell Process
                                      ↓
PTY → readLoop → vt10x → screen → Window.Invalidate()
```

---
*Last Updated: Reference guide creation*
