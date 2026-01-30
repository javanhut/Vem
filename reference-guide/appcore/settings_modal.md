# settings_modal.go Reference

**Path:** `internal/appcore/settings_modal.go`
**Lines:** ~400
**Purpose:** Settings modal UI - tabbed interface, widget state, keyboard navigation, save/cancel flow

## Overview

Provides a modal overlay for editing Vem settings. Uses the same overlay pattern as leader_rendering.go (semi-transparent background + centered popup). Features tab navigation (General, Editor, LSP) with keyboard and mouse support.

## Types

### settingsModalState (lines 26-59)
```go
type settingsModalState struct {
    active      bool
    selectedTab int

    // Tab buttons
    tabGeneral widget.Clickable
    tabEditor  widget.Clickable
    tabLSP     widget.Clickable

    // General settings widgets
    fontSizeMinus, fontSizePlus widget.Clickable
    fontSizeValue int

    // Editor settings widgets
    autoIndent, autoPairs, softWrap, formatOnSave, useSpaces widget.Bool
    tabWidthMinus, tabWidthPlus widget.Clickable
    tabWidthValue int
    scrollOffsetMinus, scrollOffsetPlus widget.Clickable
    scrollOffsetValue int

    // LSP settings widgets
    lspEnabled, lspAutoDetect widget.Bool

    // Action buttons
    saveButton, cancelButton widget.Clickable
}
```
Persistent widget state for the settings modal. Added as `settingsModal` field on appState in app.go:274.

### Tab Constants (lines 20-24)
```go
const (
    tabGeneral = iota
    tabEditor
    tabLSP
)
```

## Functions

### initSettingsModal() (lines 62-78)
```go
func (s *appState) initSettingsModal()
```
Initializes widget values from current `s.settings`. Called when modal opens.

**Copies from Settings:**
- UI: FontSize
- Editor: AutoIndent, AutoPairs, SoftWrap, FormatOnSave, UseSpaces, TabWidth, ScrollOffset
- LSP: Enabled, AutoDetect

### drawSettingsModal() (lines 81-155)
```go
func (s *appState) drawSettingsModal(gtx layout.Context) layout.Dimensions
```
Main modal rendering function. Draws:
- Semi-transparent overlay
- Centered popup (max 600x450, responsive)
- Title, tab bar, content area, button row

### drawSettingsTabs() (lines 158-222)
```go
func (s *appState) drawSettingsTabs(gtx layout.Context) layout.Dimensions
```
Renders horizontal tab bar. Handles click detection and visual state (selected/unselected).

### drawTabContent() (lines 225-236)
```go
func (s *appState) drawTabContent(gtx layout.Context) layout.Dimensions
```
Switches on `selectedTab` to render appropriate content.

### drawGeneralTab() (lines 239-282)
```go
func (s *appState) drawGeneralTab(gtx layout.Context) layout.Dimensions
```
Renders General settings:
- Font Size stepper (+/-, range 8-32)
- Theme display (read-only placeholder for Phase 3)

### drawEditorTab() (lines 285-290)
```go
func (s *appState) drawEditorTab(gtx layout.Context) layout.Dimensions
```
Stub for Plan 02-02. Will render editor behavior settings.

### drawLSPTab() (lines 293-298)
```go
func (s *appState) drawLSPTab(gtx layout.Context) layout.Dimensions
```
Stub for Plan 02-02. Will render LSP settings.

### drawSectionHeader() (lines 301-308)
```go
func (s *appState) drawSectionHeader(gtx layout.Context, title string) layout.Dimensions
```
Renders section headers with blue bold text.

### drawNumberStepper() (lines 311-357)
```go
func (s *appState) drawNumberStepper(gtx layout.Context, label string, value *int, minus, plus *widget.Clickable, minVal, maxVal int) layout.Dimensions
```
Reusable stepper control: [Label] [-] [value] [+]
Handles click events and value clamping.

### drawSettingsButtons() (lines 360-413)
```go
func (s *appState) drawSettingsButtons(gtx layout.Context) layout.Dimensions
```
Renders Cancel and Save buttons. Handles click events.

### handleSettingsModalKey() (lines 416-434)
```go
func (s *appState) handleSettingsModalKey(ev key.Event)
```
Handles keyboard input when modal is active:
- Escape: Close without saving
- Enter: Save and close
- Tab: Cycle tabs
- 1/2/3: Jump to specific tab

### closeSettingsModal() (lines 437-440)
```go
func (s *appState) closeSettingsModal()
```
Closes modal without saving changes. Clears status.

### saveAndCloseSettings() (lines 443-455)
```go
func (s *appState) saveAndCloseSettings()
```
Copies widget values to Settings, calls Save() and applySettings(), closes modal.
Currently only saves FontSize (Editor/LSP added in Plan 02-02).

## Integration Points

- **app.go:274** - `settingsModal settingsModalState` field
- **app.go:677-679** - Modal rendering in layout() (before fuzzy finder)
- **app.go:1923-1926** - Key interception in handleKey()
- **settings.go:175-179** - openSettingsModal() activates modal

## Triggers

- **Ctrl+M** - Normal mode keybinding (opens modal)
- **:settings** - Command mode

## Keyboard Shortcuts (when modal active)

| Key | Action |
|-----|--------|
| Escape | Close without saving |
| Enter | Save and close |
| Tab | Cycle to next tab |
| 1 | Jump to General tab |
| 2 | Jump to Editor tab |
| 3 | Jump to LSP tab |

## Modal Layout

```
┌─────────────────────────────────────────┐
│  Settings                                │ <- Title
├──────────┬──────────┬──────────┬────────┤
│ [General]│ [Editor] │  [LSP]   │        │ <- Tab bar
├──────────┴──────────┴──────────┴────────┤
│                                          │
│  Appearance                              │ <- Section header
│  Font Size          [-]  14  [+]         │ <- Number stepper
│  Theme              dark                 │ <- Display only
│  Theme selection coming in Phase 3       │ <- Note
│                                          │
├──────────────────────────────────────────┤
│        [Cancel]  [Save]                  │ <- Buttons
└──────────────────────────────────────────┘
```

---
*Created: 2026-01-30 - Phase 2 Plan 02-01*
