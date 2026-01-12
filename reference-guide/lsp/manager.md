# manager.go

**Path:** `/home/javanhut/Development/Vem/internal/lsp/manager.go`
**Lines:** 418
**Purpose:** Language Server Protocol server lifecycle management

## Overview

Manages multiple LSP server instances with:
- Per-workspace server instances
- Automatic server startup
- Diagnostic notification handling
- Server health monitoring
- Document state tracking

## Code Blocks

### Lines 1-30: Package and Struct

#### Manager (Lines 12-30)
```go
type Manager struct {
    servers   map[string]*ServerInstance  // workspace root -> server
    serversMu sync.RWMutex
    configs   map[string]ServerConfig     // language -> config

    // Callbacks
    onDiagnostics func(uri DocumentURI, diagnostics []Diagnostic)
    onStatus      func(message string)
    onError       func(err error)
    onInvalidate  func()  // UI refresh
}
```

#### ServerInstance (Lines 32-42)
```go
type ServerInstance struct {
    Config        ServerConfig
    Client        *Client
    WorkspaceRoot string
    Documents     map[DocumentURI]*DocumentState
    DocumentsMu   sync.RWMutex
    Capabilities  ServerCapabilities
    Diagnostics   map[DocumentURI][]Diagnostic
    DiagnosticsMu sync.RWMutex
}
```

#### DocumentState (Lines 44-50)
```go
type DocumentState struct {
    URI        DocumentURI
    Version    int
    Content    string
    LanguageID string
}
```

### Lines 52-78: Constructor and Callbacks

| Function | Lines | Purpose |
|----------|-------|---------|
| `NewManager()` | 52-58 | Create manager with default configs |
| `OnDiagnostics()` | 60-63 | Set diagnostic callback |
| `OnStatus()` | 65-68 | Set status message callback |
| `OnError()` | 70-73 | Set error callback |
| `OnInvalidate()` | 75-78 | Set UI refresh callback |

### Lines 80-107: Server Lookup

#### GetServerForFile (Lines 80-107)
Gets or creates server for file:
1. Get config for file type
2. Check if server is available
3. Find project root using root patterns
4. Check for existing server at root
5. Start new server if needed

### Lines 109-280: Server Startup

#### startServer (Lines 109-280)
Starts new language server:
1. Resolves full command path via `FindServerCommand()` (handles gopls in ~/go/bin)
2. Creates LSP client with resolved path
3. Starts client process
4. Sets up error callback
5. Registers diagnostic notification handler
6. Sends `initialize` request with full capabilities
7. Sends `initialized` notification
8. Starts server monitor goroutine

**Client Capabilities Registered:**
- Workspace: applyEdit, documentChanges, workspaceFolders
- TextDocument: sync, completion, hover, definition, references
- CodeAction: quickfix, refactor, source organize
- Formatting, Rename, PublishDiagnostics

### Lines 277-299: Server Monitoring

#### monitorServer (Lines 277-299)
Background goroutine checking server health:
- Polls every second
- Removes dead servers from map
- Reports status changes
- Reports last error if any

### Lines 301-337: Server Control

| Function | Lines | Purpose |
|----------|-------|---------|
| `StopAll()` | 301-310 | Stop all servers |
| `StopServer()` | 312-321 | Stop specific server |
| `RestartServerForFile()` | 323-337 | Restart server for file |

### Lines 339-375: Diagnostics Access

| Function | Lines | Purpose |
|----------|-------|---------|
| `GetDiagnostics()` | 339-356 | Get diagnostics for file |
| `GetAllDiagnostics()` | 358-375 | Get all diagnostics |

### Lines 377-417: Status and Queries

| Function | Lines | Purpose |
|----------|-------|---------|
| `GetServerStatus()` | 377-392 | Get status of all servers |
| `HasServerForFile()` | 394-398 | Check if config exists |
| `IsServerRunningForFile()` | 400-417 | Check if server running |

## Known Issues / Potential Bugs

1. **Line 282: Ticker never stops on crash**
   - Monitor goroutine returns on crash
   - Ticker is properly deferred, but pattern could be cleaner

2. **Line 244: 30 second timeout for initialize**
   - May be too long for user experience
   - May be too short for slow servers

## Dead/Unused Code

None identified.

## Integration Points

- Used by `appState.lspManager` in app.go
- Config from `lsp/config.go`
- Client from `lsp/client.go`
- Features accessed via `lsp/features.go`

## Server Lifecycle

```
GetServerForFile()
  ├── Check config exists
  ├── Check server available
  ├── Find project root
  ├── Check existing server
  └── startServer()
       ├── NewClient()
       ├── Client.Start()
       ├── Register handlers
       ├── initialize request
       ├── initialized notification
       └── monitorServer()
```

---
*Last Updated: After gopls PATH detection fix*
