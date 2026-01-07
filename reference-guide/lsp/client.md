# client.go

**Path:** `/home/javanhut/Development/Vem/internal/lsp/client.go`
**Lines:** 532
**Purpose:** JSON-RPC 2.0 client for LSP communication

## Overview

Implements LSP client with:
- Process management for language server
- JSON-RPC 2.0 request/response protocol
- Content-Length header framing
- Notification handlers for server events
- Graceful shutdown with timeout

## Code Blocks

### Lines 1-16: Package and Imports

```go
package lsp
import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)
```

### Lines 18-53: Client Struct

```go
type Client struct {
    cmd    *exec.Cmd
    stdin  io.WriteCloser
    stdout *bufio.Reader
    stderr io.ReadCloser

    // Request/response tracking
    nextID      int64
    pendingReqs map[int64]chan *Response
    pendingMu   sync.Mutex

    // Notification handlers
    handlers   map[string]NotificationHandler
    handlersMu sync.RWMutex

    // Lifecycle
    ctx       context.Context
    cancel    context.CancelFunc
    running   bool
    runningMu sync.RWMutex
    lastErr   error
    lastErrMu sync.RWMutex

    // Server capabilities
    capabilities   ServerCapabilities
    capabilitiesMu sync.RWMutex

    // Error callback
    onError func(error)

    writeMu sync.Mutex

    stderrMu       sync.RWMutex
    lastStderrLine string
}
```

### Lines 55-83: JSON-RPC Types

| Type | Lines | Purpose |
|------|-------|---------|
| `Request` | 55-61 | JSON-RPC request with optional ID |
| `Response` | 63-69 | JSON-RPC response with result/error |
| `ResponseError` | 71-80 | Error in JSON-RPC response |
| `NotificationHandler` | 82-83 | Function type for handling notifications |

### Lines 85-121: Constructor

#### NewClient (Lines 85-121)
Creates new LSP client:
1. Creates context with cancel
2. Creates exec.Cmd with working directory
3. Sets up stdin, stdout, stderr pipes
4. Initializes maps for pending requests and handlers

### Lines 123-148: Startup

#### Start (Lines 123-148)
Starts language server:
1. Checks if already running
2. Starts the command process
3. Launches 3 goroutines:
   - `readLoop()` for reading messages
   - `readStderr()` for stderr logging
   - `monitorProcess()` for exit detection

### Lines 150-187: Shutdown

#### Stop (Lines 150-187)
Graceful shutdown:
1. Sends `shutdown` request (2 second timeout)
2. Sends `exit` notification
3. Closes stdin
4. Cancels context
5. Waits for process (5 second timeout)
6. Force kills if necessary

### Lines 189-218: State Queries

| Function | Lines | Purpose |
|----------|-------|---------|
| `IsRunning()` | 189-194 | Check if client is running |
| `GetLastError()` | 196-201 | Get last error that occurred |
| `setLastStderrLine()` | 203-207 | Store last stderr line |
| `getLastStderrLine()` | 209-213 | Get last stderr line |
| `SetErrorCallback()` | 215-218 | Set error callback |

### Lines 220-275: Request/Response

#### Call (Lines 220-275)
Sends request and waits for response:
1. Checks if running
2. Generates unique ID atomically
3. Creates response channel
4. Sends request via `writeMessage()`
5. Waits on context, client context, or response
6. Unmarshals result if provided

### Lines 277-314: Notifications and Capabilities

| Function | Lines | Purpose |
|----------|-------|---------|
| `Notify()` | 277-293 | Send notification (no response) |
| `OnNotification()` | 295-300 | Register notification handler |
| `Capabilities()` | 302-307 | Get server capabilities |
| `SetCapabilities()` | 309-314 | Set server capabilities |

### Lines 316-340: Message Writing

#### writeMessage (Lines 316-340)
Writes LSP message with Content-Length header:
1. Marshals message to JSON
2. Writes header: `Content-Length: N\r\n\r\n`
3. Writes body
4. Uses mutex for atomic writes

### Lines 342-379: Read Loop

#### readLoop (Lines 342-379)
Continuously reads server messages:
1. Checks context for cancellation
2. Reads message via `readMessage()`
3. Stores errors and calls error callback
4. Handles message via `handleMessage()`
5. Exits on EOF or pipe closed

### Lines 381-421: Message Reading

#### readMessage (Lines 381-421)
Reads single LSP message:
1. Reads headers line by line until empty line
2. Parses Content-Length header
3. Reads exact body bytes
4. Returns raw JSON

### Lines 423-481: Message Handling

#### handleMessage (Lines 423-481)
Routes incoming messages:
1. Parses message to determine type
2. **Response** (has ID, no method): Routes to pending request channel
3. **Notification** (has method, no ID): Calls registered handler in goroutine
4. **Server Request** (has method and ID): Sends empty acknowledgment

### Lines 483-531: Background Goroutines

| Function | Lines | Purpose |
|----------|-------|---------|
| `readStderr()` | 483-490 | Read stderr to prevent blocking |
| `monitorProcess()` | 492-531 | Monitor for unexpected exits |

#### monitorProcess Details (Lines 492-531)
1. Waits for process exit
2. Sets running = false
3. Stores error with optional stderr context
4. Calls error callback
5. Closes all pending request channels

## Known Issues / Potential Bugs

None identified.

## Dead/Unused Code

None identified.

## Integration Points

- Created by `Manager.startServer()` in manager.go
- Used by `ServerInstance` for all LSP requests
- Notification handlers set for `textDocument/publishDiagnostics`

## Message Protocol

```
Client                              Server
   |                                   |
   |-- Content-Length: N\r\n\r\n ----->|
   |-- {jsonrpc, id, method, params} ->|
   |                                   |
   |<-- Content-Length: N\r\n\r\n -----|
   |<-- {jsonrpc, id, result/error} ---|
```

---
*Last Updated: Reference guide creation*
