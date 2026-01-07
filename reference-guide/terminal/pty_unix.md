# pty_unix.go

**Path:** `/home/javanhut/Development/Vem/internal/terminal/pty_unix.go`
**Lines:** 107
**Purpose:** Unix PTY implementation using creack/pty

## Overview

Provides Unix-specific terminal functionality:
- PTY creation and management
- Shell process startup
- Window resize handling
- Default shell detection

Build tag: `//go:build unix`

## Code Blocks

### Lines 1-12: Package and Imports

```go
//go:build unix

package terminal
import (
    "fmt"
    "os"
    "os/exec"
    "syscall"

    "github.com/creack/pty"
)
```

### Lines 14-65: startPTY

#### startPTY (Lines 14-65)
Creates PTY and starts shell:
1. Opens PTY pair (master/slave)
2. Sets initial size
3. Creates shell command
4. Connects command to TTY:
   - stdin = ttyFile
   - stdout = ttyFile
   - stderr = ttyFile
5. Sets process attributes:
   - `Setsid: true` - Creates new session
   - `Setctty: true` - Sets controlling terminal
6. Starts process
7. Closes TTY in parent
8. Starts wait goroutine

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Setsid:  true,  // Create new session
    Setctty: true,  // Set controlling terminal
}
```

### Lines 67-93: Resize

#### Resize (Lines 67-93)
Updates PTY window size:
1. Updates Terminal dimensions
2. Resizes screen buffer
3. Resizes vt10x emulator
4. Calls `pty.Setsize()` with new dimensions

```go
return pty.Setsize(ptyFile, &pty.Winsize{
    Rows: uint16(height),
    Cols: uint16(width),
})
```

### Lines 95-106: Default Shell

#### DefaultShell (Lines 95-101)
Returns default shell:
```go
if shell := os.Getenv("SHELL"); shell != "" {
    return shell
}
return "/bin/sh"
```

#### DefaultArgs (Lines 103-106)
Returns default shell arguments:
```go
return []string{"-i"}  // Interactive shell
```

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Called by `Terminal.Start()` in terminal.go
- PTY file used for read/write in terminal.go
- Resize called when terminal pane dimensions change

## PTY Architecture

```
┌─────────────────────────────────────────────┐
│                   Parent Process            │
│  ┌──────────────┐                           │
│  │   Terminal   │                           │
│  │   struct     │───reads/writes────┐       │
│  └──────────────┘                   │       │
│                                     ▼       │
│                              ┌──────────┐   │
│                              │  PTY     │   │
│                              │  Master  │   │
│                              └────┬─────┘   │
└───────────────────────────────────│─────────┘
                                    │
                     (kernel provides pipe)
                                    │
┌───────────────────────────────────│─────────┐
│                              ┌────▼─────┐   │
│                              │  TTY     │   │
│                              │  Slave   │   │
│                              └────┬─────┘   │
│                                   │         │
│  ┌────────┐  ┌────────┐  ┌───────▼───────┐ │
│  │ stdin  │──│ stdout │──│ Shell Process │ │
│  └────────┘  │ stderr │  │   (/bin/sh)   │ │
│              └────────┘  └───────────────┘ │
│                   Child Process            │
└─────────────────────────────────────────────┘
```

---
*Last Updated: Reference guide creation*
