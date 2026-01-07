# lsp_actions.go

**Path:** `/home/javanhut/Development/Vem/internal/appcore/lsp_actions.go`
**Lines:** 1249
**Purpose:** LSP feature handlers and document synchronization

## Overview

Implements all LSP action handlers:
- Go to definition, hover, references
- Completion with debouncing
- Rename, format, code actions
- Diagnostics navigation
- Document synchronization

## Code Blocks

### Lines 1-18: Constants and Types

```go
const (
    lspChangeDebounce     = 50 * time.Millisecond
    lspCompletionDebounce = 80 * time.Millisecond
)

type lspCompletionRequest struct {
    filePath string
    line     int
    col      int
    trigger  string
}
```

### Lines 27-75: setLSPHint

Sets helpful hints when LSP fails:
- Rust: "Install rust-analyzer via: rustup component add rust-analyzer"
- General: "Install X and ensure it is on PATH"
- Project root: "Project root not found (need go.mod, etc.)"

### Lines 77-114: handleLSPGotoDefinition

Navigates to symbol definition:
1. Syncs document
2. Calls `server.GotoDefinition()`
3. Navigates to first location

### Lines 116-147: handleLSPHover

Shows hover information:
1. Syncs document
2. Calls `server.GetHover()`
3. Sets `hoverActive` and `hoverContent`

### Lines 149-186: handleLSPReferences

Finds all references:
1. Calls `server.FindReferences()`
2. Sets `referencesActive` and `referencesItems`
3. Shows count in status

### Lines 188-219: handleLSPRename

Renames symbol:
1. Calls `server.Rename()`
2. Applies workspace edit
3. Reports files changed

### Lines 221-256: handleLSPFormat

Formats document:
1. Calls `server.FormatDocument()`
2. Applies text edits
3. Updates status

### Lines 258-305: handleLSPCodeAction

Shows available code actions:
1. Gets diagnostics at cursor
2. Calls `server.GetCodeActions()`
3. Shows code action menu

### Lines 307-340: handleLSPCompletion

Triggers completion manually:
1. Syncs document
2. Calls `server.GetCompletion()`
3. Sets completion state
4. Resolves first item

### Lines 342-426: Completion Navigation

| Function | Lines | Purpose |
|----------|-------|---------|
| `handleLSPCompletionAccept()` | 342-399 | Accept selected completion |
| `handleLSPCompletionCancel()` | 401-405 | Cancel completion |
| `handleLSPCompletionNext()` | 407-414 | Move to next item |
| `handleLSPCompletionPrev()` | 416-426 | Move to previous item |

### Lines 428-475: resolveCompletionItem

Resolves additional completion details in background goroutine:
1. Checks version to avoid stale updates
2. Calls `server.ResolveCompletion()`
3. Updates item in place
4. Invalidates window

### Lines 477-546: Diagnostics Navigation

| Function | Lines | Purpose |
|----------|-------|---------|
| `handleLSPNextDiagnostic()` | 477-510 | Move to next diagnostic |
| `handleLSPPrevDiagnostic()` | 512-546 | Move to previous diagnostic |

### Lines 548-614: Dismiss Handlers

| Function | Lines | Purpose |
|----------|-------|---------|
| `handleLSPDismissHover()` | 548-553 | Hide hover tooltip |
| `handleLSPDismissReferences()` | 555-559 | Hide references list |
| `handleLSPDismissCodeActions()` | 561-565 | Hide code actions menu |
| `handleLSPReferencesNext/Prev()` | 567-584 | Navigate references |
| `handleLSPReferencesOpen()` | 586-595 | Open selected reference |
| `handleLSPCodeActionSelect()` | 597-614 | Apply selected action |

### Lines 616-674: Status and Control

| Function | Lines | Purpose |
|----------|-------|---------|
| `handleLSPStatus()` | 617-652 | Show LSP server status |
| `handleLSPRestart()` | 654-674 | Restart language server |

### Lines 676-726: Navigation Helper

#### navigateToLocation (Lines 692-726)
Navigates to LSP location:
1. Opens file if different
2. Sets up LSP for buffer
3. Jumps to position

### Lines 728-821: Workspace Edits

| Function | Lines | Purpose |
|----------|-------|---------|
| `applyWorkspaceEdit()` | 728-761 | Apply multi-file edit |
| `getOrOpenBuffer()` | 763-782 | Get or open buffer |
| `applyTextEdits()` | 784-821 | Apply text edits to buffer |

**Note:** Edits are sorted in reverse order to preserve positions.

### Lines 823-903: Debounced Change Notification

| Function | Lines | Purpose |
|----------|-------|---------|
| `scheduleLSPChange()` | 831-849 | Schedule debounced change |
| `flushLSPChange()` | 851-880 | Send change to server |
| `flushPendingLSPChange()` | 882-903 | Flush immediately |
| `cancelLSPChange()` | 905-918 | Cancel pending change |

### Lines 920-1063: Debounced Completion

| Function | Lines | Purpose |
|----------|-------|---------|
| `scheduleLSPCompletion()` | 920-938 | Schedule completion request |
| `cancelLSPCompletion()` | 940-953 | Cancel pending request |
| `runLSPCompletion()` | 955-1001 | Execute completion |
| `maybeTriggerLSPCompletion()` | 1003-1063 | Auto-trigger on typing |

### Lines 1065-1148: Buffer Setup

| Function | Lines | Purpose |
|----------|-------|---------|
| `startLSPForFile()` | 1065-1089 | Start LSP for file |
| `maybeAutoInstallLSP()` | 1091-1122 | Auto-install missing servers |
| `setupLSPForBuffer()` | 1124-1147 | Set up LSP for buffer |
| `cleanupLSPForBuffer()` | 1149-1169 | Clean up on buffer close |

### Lines 1171-1248: Utilities

| Function | Lines | Purpose |
|----------|-------|---------|
| `getLSPDiagnosticsForLine()` | 1171-1185 | Get diagnostics for line |
| `getDiagnosticCountString()` | 1187-1212 | Format "E:N W:N" string |
| `processLSPCommand()` | 1214-1248 | Process :command LSP commands |

## Known Issues / Potential Bugs

1. **Line 719-721: Column positioning loop**
   - Uses loop to move cursor to column
   - Could be replaced with direct column set

## Dead/Unused Code

None identified.

## Integration Points

- Called from keybindings.go action handlers
- Uses lsp/features.go for server calls
- Rendering in lsp_rendering.go

---
*Last Updated: Reference guide creation*
