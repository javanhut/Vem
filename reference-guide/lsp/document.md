# document.go

**Path:** `/home/javanhut/Development/Vem/internal/lsp/document.go`
**Lines:** 151
**Purpose:** Document synchronization with language servers

## Overview

Implements LSP document lifecycle notifications:
- `textDocument/didOpen`
- `textDocument/didChange`
- `textDocument/didSave`
- `textDocument/didClose`

Uses full document sync (not incremental) for maximum compatibility.

## Code Blocks

### Lines 1-27: OpenDocument

#### OpenDocument (Lines 1-27)
Notifies server that document was opened:
1. Converts file path to URI
2. Gets language ID from extension
3. Creates `DocumentState` with version 1
4. Stores in `Documents` map
5. Sends `textDocument/didOpen` notification

```go
params := DidOpenTextDocumentParams{
    TextDocument: TextDocumentItem{
        URI:        uri,
        LanguageID: langID,
        Version:    1,
        Text:       content,
    },
}
return si.Client.Notify("textDocument/didOpen", params)
```

### Lines 29-60: ChangeDocument

#### ChangeDocument (Lines 29-60)
Notifies server of document changes:
1. Converts path to URI
2. Checks if document is open
3. If not open, calls `OpenDocument()` instead
4. Increments version
5. Updates stored content
6. Sends full content (not incremental changes)

**Note:** Uses full document sync for compatibility with all servers.

### Lines 62-84: SaveDocument

#### SaveDocument (Lines 62-84)
Notifies server that document was saved:
1. Gets document state
2. If not tracked, returns nil (no-op)
3. Sends `textDocument/didSave` with content

```go
params := DidSaveTextDocumentParams{
    TextDocument: TextDocumentIdentifier{URI: uri},
    Text:         content,  // Include text if server supports it
}
```

### Lines 86-104: CloseDocument

#### CloseDocument (Lines 86-104)
Notifies server that document was closed:
1. Removes from `Documents` map
2. Clears diagnostics for file
3. Sends `textDocument/didClose` notification

### Lines 106-128: Query Functions

| Function | Lines | Purpose |
|----------|-------|---------|
| `IsDocumentOpen()` | 106-115 | Check if document is tracked |
| `GetDocumentVersion()` | 117-128 | Get current version number |

### Lines 130-150: Sync Utility

#### SyncDocument (Lines 130-150)
Ensures document is open and synced:
1. Checks if document is already open
2. If not open, opens it
3. If open but content differs, sends change notification
4. Useful before making LSP requests

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Called from `appcore/lsp_actions.go` when buffers change
- Used before LSP feature requests (completion, hover, etc.)
- State stored in `ServerInstance.Documents` map

## Document Lifecycle

```
Buffer Opened → OpenDocument()
     ↓
Buffer Modified → ChangeDocument() (increments version)
     ↓
Buffer Saved → SaveDocument()
     ↓
Buffer Closed → CloseDocument() (clears diagnostics)
```

## Version Tracking

- Version starts at 1 on open
- Increments on each change
- Server uses version to track document state
- Prevents applying stale edits

---
*Last Updated: Reference guide creation*
