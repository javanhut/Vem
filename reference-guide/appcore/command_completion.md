# command_completion.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/command_completion.go`
**Lines:** ~175
**Purpose:** Path completion for command mode (Tab completion for :cd, :e, etc.)

## Overview

Implements Tab completion for file paths in command mode. Supports commands like `:cd`, `:e`, `:edit`, `:w`, `:write`, `:saveas`, `:install`.

## Features

- Triggered with Tab in COMMAND mode
- Supports tilde expansion (`~` → home directory)
- Directories shown with trailing `/`
- Tab cycles forward, Shift+Tab cycles backward
- Cancels on typing or backspace

## Code Blocks

### Lines 1-30: handleCommandTabComplete

Main handler for Tab key in command mode:
1. If completion active, cycles to next/prev item
2. Parses command to find path argument
3. Checks if command supports path completion
4. Gets completions and updates `cmdText`

### Lines 32-50: Supported Commands

```go
pathCommands := map[string]bool{
    "cd":      true,
    "e":       true,
    "edit":    true,
    "w":       true,
    "write":   true,
    "saveas":  true,
    "install": true,
}
```

### Lines 95-105: expandTilde

Expands `~` to user's home directory:
```go
func expandTilde(path string) string {
    if strings.HasPrefix(path, "~") {
        home, err := os.UserHomeDir()
        if err == nil {
            return home + path[1:]
        }
    }
    return path
}
```

### Lines 107-165: getPathCompletions

Returns matching paths for partial input:
1. Determines directory and prefix from partial path
2. Reads directory entries
3. Filters by prefix (case-insensitive)
4. Skips hidden files unless prefix starts with `.`
5. Adds trailing `/` for directories
6. Preserves tilde in output if input had tilde

### Lines 167-175: cancelCommandCompletion

Clears completion state:
```go
func (s *appState) cancelCommandCompletion() {
    s.cmdCompletionActive = false
    s.cmdCompletionItems = nil
    s.cmdCompletionIndex = 0
    s.cmdCompletionPrefix = ""
}
```

## Integration Points

### State Fields (in app.go)
```go
cmdCompletionActive  bool
cmdCompletionItems   []string
cmdCompletionIndex   int
cmdCompletionPrefix  string
cmdCompletionReverse bool
```

### Keybindings (in keybindings.go)
- Tab in COMMAND mode: `ActionCommandTabComplete`
- Shift+Tab in COMMAND mode: `ActionCommandTabComplete` (cycles backward)

### Auto-Cancel (in app.go)
- `appendCommandText()` calls `cancelCommandCompletion()` when user types
- `deleteCommandChar()` calls `cancelCommandCompletion()` when user deletes

## Usage Example

```
:cd ~/Dev     [Tab]  →  :cd ~/Development/
:cd ~/Dev     [Tab]  →  :cd ~/DevOps/        (cycles to next)
:e src/       [Tab]  →  :e src/main.go
:e src/       [Tab]  →  :e src/utils/        (cycles)
```

---
*Last Updated: After command path completion implementation*
