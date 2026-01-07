# conpty_wrapper.go

**Path:** `/home/javanhut/Development/Vem/internal/terminal/conpty_wrapper.go`
**Lines:** 15
**Purpose:** Wrapper to adapt ConPTY to ConPtyIO interface

## Overview

Simple wrapper that embeds `conpty.ConPty` to implement the `ConPtyIO` interface defined in terminal.go.

Build tag: `//go:build windows`

## Code Blocks

### Lines 1-5: Package and Imports

```go
//go:build windows

package terminal
import "github.com/UserExistsError/conpty"
```

### Lines 7-14: ConPtyWrapper

```go
// ConPtyWrapper wraps the ConPTY to implement ConPtyIO interface
type ConPtyWrapper struct {
    *conpty.ConPty
}

func (c *ConPtyWrapper) Resize(width, height int) error {
    return c.ConPty.Resize(width, height)
}
```

## Interface Implementation

The wrapper implements `ConPtyIO` (defined in terminal.go):
```go
type ConPtyIO interface {
    Read(p []byte) (int, error)
    Write(p []byte) (int, error)
    Close() error
    Resize(width, height int) error
}
```

- `Read`, `Write`, `Close` - inherited from embedded `*conpty.ConPty`
- `Resize` - explicitly implemented to adapt interface

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Created in `pty_windows.go:startPTY()`
- Stored in `Terminal.conpty` field
- Used for all PTY I/O on Windows

---
*Last Updated: Reference guide creation*
