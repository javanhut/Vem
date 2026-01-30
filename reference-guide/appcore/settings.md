# settings.go Reference

**Path:** `internal/appcore/settings.go`
**Lines:** ~190
**Purpose:** Settings infrastructure - struct definitions, TOML persistence, cross-platform paths, runtime integration

## Overview

Provides the Settings system for Vem, allowing users to configure editor behavior via a TOML file. Uses `os.UserConfigDir()` for cross-platform path support. Settings are loaded on startup and synced to appState runtime fields.

## Types

### Settings (lines 13-17)
```go
type Settings struct {
    Editor EditorSettings `toml:"editor"`
    UI     UISettings     `toml:"ui"`
    LSP    LSPSettings    `toml:"lsp"`
}
```
Top-level configuration container with three sections.

### EditorSettings (lines 20-28)
```go
type EditorSettings struct {
    TabWidth     int  `toml:"tab_width"`      // Default: 4
    UseSpaces    bool `toml:"use_spaces"`     // Default: false
    AutoIndent   bool `toml:"auto_indent"`    // Default: true
    AutoPairs    bool `toml:"auto_pairs"`     // Default: true
    ScrollOffset int  `toml:"scroll_offset"`  // Default: 5
    SoftWrap     bool `toml:"soft_wrap"`      // Default: true
    FormatOnSave bool `toml:"format_on_save"` // Default: false
}
```

### UISettings (lines 31-37)
```go
type UISettings struct {
    FontSize        int    `toml:"font_size"`           // Default: 14
    Theme           string `toml:"theme"`               // Default: "dark"
    ShowLineNumbers bool   `toml:"show_line_numbers"`   // Default: true
    RelativeNumbers bool   `toml:"relative_line_numbers"` // Default: false
    ExplorerWidth   int    `toml:"explorer_width"`      // Default: 250
}
```

### LSPSettings (lines 40-43)
```go
type LSPSettings struct {
    Enabled    bool `toml:"enabled"`     // Default: true
    AutoDetect bool `toml:"auto_detect"` // Default: true
}
```

## Functions

### NewSettings() (lines 46-68)
```go
func NewSettings() *Settings
```
Factory function returning Settings with sensible defaults. Always call this before Load() to ensure defaults are present.

### ConfigDir() (lines 76-81)
```go
func ConfigDir() (string, error)
```
Returns platform-appropriate config directory:
- Linux: `~/.config/vem` (or `$XDG_CONFIG_HOME/vem`)
- macOS: `~/Library/Application Support/vem`
- Windows: `%APPDATA%\vem`

### ConfigPath() (lines 84-90)
```go
func ConfigPath() (string, error)
```
Returns full path to settings.toml (e.g., `~/.config/vem/settings.toml`).

### Load() (lines 96-110)
```go
func (s *Settings) Load() error
```
Loads settings from config file, merging over existing defaults in the receiver. Returns nil if file doesn't exist (uses defaults silently). Returns error for permission issues or malformed TOML.

**Pattern:** Call `NewSettings()` first, then `Load()` on the result.

### Save() (lines 114-136)
```go
func (s *Settings) Save() error
```
Writes settings to config file. Creates config directory if it doesn't exist (0755 permissions).

### applySettings() (lines 139-172)
```go
func (s *appState) applySettings()
```
Syncs Settings values to appState runtime fields. Called after Load() and after settings modal saves.

**Fields synced:**
- Editor: autoIndentEnabled, autoPairsEnabled, scrollOffsetLines, softWrapEnabled, formatOnSave, indentString
- UI: explorerWidth (FontSize and Theme in Phase 2)
- LSP: lspEnabled, lspAutoEnabled

### openSettingsModal() (lines 175-179)
```go
func (s *appState) openSettingsModal()
```
Opens the settings modal. Initializes widget values from current settings and activates the modal.

See also: [settings_modal.md](settings_modal.md) for modal rendering and tab navigation.

### settingsPath() (lines 184-189)
```go
func (s *appState) settingsPath() string
```
Returns the config file path for display purposes.

## Config File Format

```toml
[editor]
tab_width = 4
use_spaces = false
auto_indent = true
auto_pairs = true
scroll_offset = 5
soft_wrap = true
format_on_save = false

[ui]
font_size = 14
theme = "dark"
show_line_numbers = true
relative_line_numbers = false
explorer_width = 250

[lsp]
enabled = true
auto_detect = true
```

## Integration Points

- **app.go:99** - appState has `settings *Settings` field
- **app.go:447-452** - Settings loaded in newAppState()
- **app.go:3319** - `:settings` command
- **keybindings.go:146** - ActionOpenSettings constant
- **keybindings.go:210** - Ctrl+M keybinding in Normal mode
- **keybindings.go:448-449** - executeAction case for ActionOpenSettings

## Triggers

- **:settings** - Command mode
- **Ctrl+M** - Normal mode keybinding

## Usage Example

```go
// In newAppState()
state.settings = NewSettings()
if err := state.settings.Load(); err != nil {
    state.status = fmt.Sprintf("Settings load error: %v (using defaults)", err)
}
state.applySettings()
```

---
*Updated: 2026-01-30 - Phase 2 Plan 02-01 Complete*
