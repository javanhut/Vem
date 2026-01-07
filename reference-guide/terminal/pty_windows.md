# pty_windows.go

**Path:** `/home/javanhut/Development/Vem/internal/terminal/pty_windows.go`
**Lines:** 95
**Purpose:** Windows ConPTY implementation

## Overview

Provides Windows-specific terminal functionality:
- ConPTY creation and management
- PowerShell/CMD process startup
- Window resize handling
- Default shell detection (pwsh > powershell > cmd)

Build tag: `//go:build windows`

## Code Blocks

### Lines 1-13: Package and Imports

```go
//go:build windows

package terminal
import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/UserExistsError/conpty"
)
```

### Lines 15-49: startPTY

#### startPTY (Lines 15-49)
Creates ConPTY and starts shell:
1. Builds command line string (ConPTY requires single string)
2. Creates ConPTY with options:
   - `ConPtyDimensions(width, height)`
   - `ConPtyWorkDir(workingDir)`
   - `ConPtyEnv(environment)`
3. Wraps in `ConPtyWrapper`
4. Starts wait goroutine with onExit callback

```go
cpty, err := conpty.Start(
    commandLine,
    conpty.ConPtyDimensions(t.width, t.height),
    conpty.ConPtyWorkDir(t.workingDir),
    conpty.ConPtyEnv(t.getEnvironment()),
)
```

### Lines 51-69: Resize

#### Resize (Lines 51-69)
Updates ConPTY window size:
1. Updates Terminal dimensions
2. Resizes screen buffer
3. Calls `cpty.Resize(width, height)`

### Lines 71-94: Default Shell

#### DefaultShell (Lines 71-89)
Returns default shell with priority:
1. **PowerShell Core** (`pwsh.exe`) - if available
2. **Windows PowerShell** (`powershell.exe`) - if available
3. **COMSPEC** environment variable
4. **Fallback** `cmd.exe`

```go
if _, err := exec.LookPath("pwsh.exe"); err == nil {
    return "pwsh.exe"
}
if _, err := exec.LookPath("powershell.exe"); err == nil {
    return "powershell.exe"
}
if comspec := os.Getenv("COMSPEC"); comspec != "" {
    return comspec
}
return "cmd.exe"
```

#### DefaultArgs (Lines 91-94)
Returns empty args for Windows:
```go
return []string{}  // No args needed
```

## Known Issues / Potential Bugs

1. **Line 21-22: Command line string building**
   - Simple space joining may break with quoted arguments
   - Consider proper escaping for complex commands

## Dead/Unused Code

None identified.

## Integration Points

- Called by `Terminal.Start()` in terminal.go
- ConPTY wrapper used for read/write via ConPtyIO interface
- Resize called when terminal pane dimensions change

## ConPTY Architecture

```
┌─────────────────────────────────────────────┐
│                   Vem Process               │
│  ┌──────────────┐                           │
│  │   Terminal   │                           │
│  │   struct     │───reads/writes────┐       │
│  └──────────────┘                   │       │
│                                     ▼       │
│                           ┌─────────────┐   │
│                           │ ConPtyWrapper│   │
│                           │  (ConPtyIO)  │   │
│                           └──────┬──────┘   │
└──────────────────────────────────│──────────┘
                                   │
                     (Windows ConPTY API)
                                   │
┌──────────────────────────────────│──────────┐
│                           ┌──────▼──────┐   │
│                           │   ConPTY    │   │
│                           │ Pseudoconsole│   │
│                           └──────┬──────┘   │
│                                  │          │
│  ┌────────┐  ┌────────┐  ┌──────▼────────┐ │
│  │ stdin  │──│ stdout │──│ pwsh.exe /    │ │
│  └────────┘  │ stderr │  │ powershell /  │ │
│              └────────┘  │ cmd.exe       │ │
│                   Child Process           │
└─────────────────────────────────────────────┘
```

## Differences from Unix

| Aspect | Unix | Windows |
|--------|------|---------|
| PTY Library | creack/pty | UserExistsError/conpty |
| Default Shell | $SHELL or /bin/sh | pwsh > powershell > cmd |
| Shell Args | `-i` (interactive) | (none) |
| Command Line | argv array | single string |
| Session Control | Setsid/Setctty | Handled by ConPTY |

---
*Last Updated: Reference guide creation*
