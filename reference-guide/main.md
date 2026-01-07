# main.go

**Path:** `/home/javanhut/Development/Vem/main.go`
**Lines:** 27
**Purpose:** Application entry point

## Overview

The main package creates the Gio window and launches the application. It's a minimal entry point that delegates all logic to the `appcore` package.

## Code Blocks

### Lines 1-10: Package and Imports
```go
package main

import (
    "os"
    gioapp "gioui.org/app"
    "gioui.org/unit"
    "github.com/javanhut/vem/internal/appcore"
)
```
- Imports Gio UI framework for cross-platform GPU-accelerated rendering
- Imports internal `appcore` package for application logic

### Lines 12-26: Main Function
```go
func main() {
    go func() {
        w := new(gioapp.Window)
        w.Option(
            gioapp.Title("Vem - Vim Emulator"),
            gioapp.Size(unit.Dp(960), unit.Dp(640)),
        )
        filePaths := os.Args[1:]
        if err := appcore.Run(w, filePaths); err != nil {
            // Silently handle app exit errors
        }
        os.Exit(0)
    }()
    gioapp.Main()
}
```

| Lines | Description |
|-------|-------------|
| 13 | Starts application in a goroutine (Gio requirement) |
| 14 | Creates new Gio window instance |
| 15-18 | Sets window title and default size (960x640 dp) |
| 19 | Captures command-line arguments as file paths to open |
| 20-22 | Calls `appcore.Run()` to start the application |
| 23 | Exits cleanly when app closes |
| 25 | `gioapp.Main()` runs the platform event loop |

## Dependencies

- **gioui.org/app** - Window management
- **gioui.org/unit** - Device-independent units (dp)
- **internal/appcore** - Application logic

## Flow

```
main()
  └── goroutine
       ├── Create gioapp.Window
       ├── Set window options (title, size)
       ├── Parse CLI args for file paths
       └── Call appcore.Run(window, filePaths)
  └── gioapp.Main() - Event loop (blocks)
```

## Potential Issues

None identified.

## Dead Code

None identified.

---
*Last Updated: Reference guide creation*
