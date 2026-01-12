# leader.go & leader_rendering.go

**Paths:**
- `/home/javanhut/Development/Vem/internal/appcore/leader.go`
- `/home/javanhut/Development/Vem/internal/appcore/leader_rendering.go`

**Purpose:** Leader bar feature - which-key style popup for custom keybindings

## Overview

The leader bar provides a Neovim which-key style popup that appears when pressing Space twice in normal mode. It displays available custom keybindings loaded from `~/.config/vem/keybindings.toml` and allows multi-key sequences (e.g., `gd` for goto definition).

## Configuration File

**Location:** `~/.config/vem/keybindings.toml`

**Format:**
```toml
[[keybind]]
name = "Go to Definition"
keys = "gd"
action = "lsp_goto_definition"

[[keybind]]
name = "Run Tests"
keys = "rt"
command = "go test ./..."
```

**Fields:**
- `name` - Display name shown in popup
- `keys` - Key sequence after leader (e.g., "gd" means press g then d)
- `action` - Vem action name (see actionMap)
- `command` - Shell command to execute (optional, alternative to action)

## Structs

### LeaderBinding
```go
type LeaderBinding struct {
    Name    string `toml:"name"`
    Keys    string `toml:"keys"`
    Action  string `toml:"action"`
    Command string `toml:"command"`
}
```

### LeaderConfig
```go
type LeaderConfig struct {
    Keybinds []LeaderBinding `toml:"keybind"`
}
```

## State Fields (in appState)

| Field | Type | Purpose |
|-------|------|---------|
| `leaderBarActive` | `bool` | Whether leader bar is visible |
| `leaderBarSequence` | `string` | Accumulated key sequence |
| `leaderBarBindings` | `[]LeaderBinding` | All loaded bindings |
| `leaderBarMatches` | `[]LeaderBinding` | Bindings matching current sequence |
| `leaderBarIndex` | `int` | Selected item in popup |
| `lastSpaceTime` | `time.Time` | For double-space detection |

## Functions

### leader.go

| Function | Purpose |
|----------|---------|
| `loadLeaderConfig()` | Load bindings from TOML file |
| `getDefaultBindings()` | Return minimal default bindings |
| `handleSpaceKey()` | Detect double-space for leader activation |
| `enterLeaderBar()` | Show popup, reset state |
| `exitLeaderBar()` | Hide popup, clear state |
| `handleLeaderKey()` | Process key input, filter matches |
| `getMatchingBindings()` | Filter bindings by sequence prefix |
| `executeLeaderBinding()` | Run action or shell command |
| `executeShellCommand()` | Run shell command via `:!` |
| `reloadLeaderConfig()` | Reload bindings (`:leaderreload`) |

### leader_rendering.go

| Function | Purpose |
|----------|---------|
| `drawLeaderBar()` | Render centered popup overlay |

## Action Mappings

The `actionMap` translates config action names to Action constants:

| Config Name | Action Constant |
|-------------|-----------------|
| `lsp_goto_definition` | `ActionLSPGotoDefinition` |
| `lsp_find_references` | `ActionLSPReferences` |
| `lsp_hover` | `ActionLSPHover` |
| `lsp_rename` | `ActionLSPRename` |
| `lsp_code_action` | `ActionLSPCodeAction` |
| `lsp_format` | `ActionLSPFormat` |
| `fuzzy_finder` | `ActionOpenFuzzyFinder` |
| `toggle_explorer` | `ActionToggleExplorer` |
| `split_vertical` | `ActionSplitVertical` |
| `split_horizontal` | `ActionSplitHorizontal` |
| `pane_close` | `ActionPaneClose` |
| `next_buffer` | `ActionNextBuffer` |
| `prev_buffer` | `ActionPrevBuffer` |

The `commandMap` handles actions that execute as commands:

| Config Name | Command |
|-------------|---------|
| `save_file` | `:w` |
| `close_buffer` | `:bd` |

## Usage

1. Press Space twice quickly (within 300ms) in normal mode
2. Leader bar popup appears with all available bindings
3. Type key sequence to filter (e.g., `g` shows all bindings starting with `g`)
4. When sequence matches exactly one binding, it executes automatically
5. Use Up/Down arrows to select, Enter to confirm, Escape to cancel
6. Backspace removes last key from sequence

## Commands

| Command | Description |
|---------|-------------|
| `:leaderreload` | Reload keybindings from config file |

## Default Bindings

See `configs/keybindings.toml` for the default bindings installed via `make install`.

---
*Last Updated: After leader bar feature implementation*
