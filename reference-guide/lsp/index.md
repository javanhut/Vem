# lsp Package Index

**Path:** `/home/javanhut/Development/Vem/internal/lsp/`
**Purpose:** Language Server Protocol client implementation

## Files Overview

| File | Lines | Purpose |
|------|-------|---------|
| [manager.go](manager.md) | 418 | Server lifecycle management |
| [client.go](client.md) | 532 | JSON-RPC 2.0 client |
| [config.go](config.md) | 492 | Server configurations |
| [document.go](document.md) | 151 | Document synchronization |
| [features.go](features.md) | 451 | LSP feature implementations |
| [types.go](types.md) | 733 | LSP protocol types |

## Key Types

### Manager (manager.go)
Top-level LSP coordinator:
- Maps workspace roots to server instances
- Handles server startup/shutdown
- Routes diagnostic notifications

### Client (client.go)
JSON-RPC transport:
- Process management
- Request/response tracking
- Notification handling

### ServerInstance (manager.go)
Per-workspace server:
- Client connection
- Document state
- Diagnostics storage
- Server capabilities

## Supported Languages

| Language | Server | Command |
|----------|--------|---------|
| Go | gopls | `gopls serve` |
| Python | pyright | `pyright-langserver --stdio` |
| TypeScript | typescript-language-server | `--stdio` |
| Rust | rust-analyzer | (no args) |
| C/C++ | clangd | `--background-index` |
| ... | ... | ... |

See [config.md](config.md) for full list.

## LSP Features

| Feature | Method | Status |
|---------|--------|--------|
| Go to Definition | textDocument/definition | ✅ |
| Hover | textDocument/hover | ✅ |
| References | textDocument/references | ✅ |
| Completion | textDocument/completion | ✅ |
| Rename | textDocument/rename | ✅ |
| Format | textDocument/formatting | ✅ |
| Code Actions | textDocument/codeAction | ✅ |
| Diagnostics | textDocument/publishDiagnostics | ✅ |

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Manager                          │
│  servers: map[workspace]*ServerInstance             │
└─────────────────────┬───────────────────────────────┘
                      │
        ┌─────────────┴─────────────┐
        ▼                           ▼
┌───────────────────┐      ┌───────────────────┐
│  ServerInstance   │      │  ServerInstance   │
│  workspace: /proj1│      │  workspace: /proj2│
│  client: *Client  │      │  client: *Client  │
│  documents: map[] │      │  documents: map[] │
└─────────┬─────────┘      └───────────────────┘
          │
          ▼
┌───────────────────┐
│     Client        │
│  cmd: gopls       │
│  stdin/stdout     │
│  pendingReqs      │
└─────────┬─────────┘
          │ JSON-RPC 2.0
          ▼
┌───────────────────┐
│  Language Server  │
│  (gopls, etc.)    │
└───────────────────┘
```

## Known Issues Summary

1. **manager.go:244** - 30 second timeout for initialize may be too long/short
2. **features.go:11** - 10 second timeout may be too short for complex operations
3. **features.go:289** - Hardcoded formatting options

---
*Last Updated: Reference guide creation*
