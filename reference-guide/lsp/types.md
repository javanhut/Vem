# types.go

**Path:** `/home/javanhut/Development/Vem/internal/lsp/types.go`
**Lines:** 733
**Purpose:** LSP protocol type definitions (LSP 3.17 specification)

## Overview

Defines all LSP protocol types including:
- Position and range types
- Document identifiers
- Diagnostics and completion items
- Code actions and workspace edits
- Client and server capabilities
- Request/response parameter types

## Code Blocks

### Lines 1-6: Package Declaration

```go
// Package lsp provides Language Server Protocol client functionality.
// It implements LSP 3.17 specification for communication with language servers.
package lsp
import "encoding/json"
```

### Lines 7-34: Basic Position Types

| Type | Lines | Fields |
|------|-------|--------|
| `Position` | 7-11 | Line, Character (0-based) |
| `Range` | 13-17 | Start, End positions |
| `Location` | 19-23 | URI, Range |
| `LocationLink` | 25-31 | OriginSelectionRange, TargetURI, TargetRange, TargetSelectionRange |
| `DocumentURI` | 33-34 | string alias for `file://` URIs |

### Lines 36-60: Document Identifiers

| Type | Lines | Purpose |
|------|-------|---------|
| `TextDocumentIdentifier` | 36-39 | Identifies document by URI |
| `VersionedTextDocumentIdentifier` | 41-45 | Adds version to identifier |
| `TextDocumentItem` | 47-53 | Full document for transfer |
| `TextDocumentPositionParams` | 55-60 | Document + position for requests |

### Lines 62-101: Diagnostics

#### DiagnosticSeverity (Lines 62-74)
```go
const (
    DiagnosticSeverityError       DiagnosticSeverity = 1
    DiagnosticSeverityWarning     DiagnosticSeverity = 2
    DiagnosticSeverityInformation DiagnosticSeverity = 3
    DiagnosticSeverityHint        DiagnosticSeverity = 4
)
```

#### DiagnosticTag (Lines 76-84)
```go
const (
    DiagnosticTagUnnecessary DiagnosticTag = 1  // Unused code
    DiagnosticTagDeprecated  DiagnosticTag = 2  // Obsolete code
)
```

#### Diagnostic (Lines 86-94)
```go
type Diagnostic struct {
    Range    Range
    Severity DiagnosticSeverity
    Code     interface{}
    Source   string
    Message  string
    Tags     []DiagnosticTag
}
```

### Lines 103-203: Completion Types

#### CompletionItemKind (Lines 103-132)
25 kinds including: Text, Method, Function, Constructor, Field, Variable, Class, Interface, Module, Property, Enum, Keyword, Snippet, etc.

#### CompletionItem (Lines 134-153)
```go
type CompletionItem struct {
    Label               string
    Kind                CompletionItemKind
    Detail              string
    Documentation       interface{}
    InsertText          string
    TextEdit            *TextEdit
    AdditionalTextEdits []TextEdit
    // ... more fields
}
```

#### CompletionTriggerKind (Lines 196-203)
```go
const (
    CompletionTriggerKindInvoked          = 1  // Manual (Ctrl+Space)
    CompletionTriggerKindTriggerCharacter = 2  // Typed character
    CompletionTriggerKindTriggerForIncompleteCompletions = 3
)
```

### Lines 205-223: Hover Types

```go
type Hover struct {
    Contents MarkupContent
    Range    *Range
}

type MarkupContent struct {
    Kind  MarkupKind  // "plaintext" or "markdown"
    Value string
}
```

### Lines 225-287: Edit Types

| Type | Lines | Purpose |
|------|-------|---------|
| `TextEdit` | 225-229 | Single text replacement |
| `WorkspaceEdit` | 231-235 | Multi-file edits |
| `TextDocumentEdit` | 237-241 | Edits for versioned document |
| `CodeAction` | 243-252 | Quickfix or refactoring |
| `CodeActionKind` | 254-267 | quickfix, refactor, source |
| `Command` | 282-287 | Server command reference |

### Lines 289-319: Request Parameters

| Type | Lines | Purpose |
|------|-------|---------|
| `ReferenceParams` | 289-298 | Find references request |
| `RenameParams` | 300-304 | Rename request |
| `DocumentFormattingParams` | 306-310 | Format request |
| `FormattingOptions` | 312-319 | Tab size, spaces, etc. |

### Lines 321-357: Document Notifications

| Type | Lines | Purpose |
|------|-------|---------|
| `DidOpenTextDocumentParams` | 321-324 | Open notification |
| `DidChangeTextDocumentParams` | 326-330 | Change notification |
| `TextDocumentContentChangeEvent` | 332-337 | Change event |
| `DidSaveTextDocumentParams` | 339-343 | Save notification |
| `DidCloseTextDocumentParams` | 345-348 | Close notification |

### Lines 350-364: Initialize Types

```go
type InitializeParams struct {
    ProcessID             int
    RootURI               DocumentURI
    Capabilities          ClientCapabilities
    InitializationOptions interface{}
    Trace                 string
    WorkspaceFolders      []WorkspaceFolder
}
```

### Lines 366-617: Client Capabilities

Extensive capability definitions for:
- Workspace capabilities (applyEdit, workspaceFolders)
- Text document capabilities:
  - Synchronization
  - Completion (snippet support, documentation formats)
  - Hover (content formats)
  - Signature help
  - Definition, TypeDefinition, Implementation
  - References, DocumentHighlight
  - DocumentSymbol (hierarchical support)
  - CodeAction, CodeLens
  - Formatting, Rename
  - PublishDiagnostics

### Lines 619-687: Server Capabilities

#### ServerCapabilities (Lines 631-652)
```go
type ServerCapabilities struct {
    TextDocumentSync           interface{}
    CompletionProvider         *CompletionOptions
    HoverProvider              interface{}
    SignatureHelpProvider      *SignatureHelpOptions
    DefinitionProvider         interface{}
    ReferencesProvider         interface{}
    CodeActionProvider         interface{}
    DocumentFormattingProvider interface{}
    RenameProvider             interface{}
    // ... more providers
}
```

### Lines 689-733: Additional Types

| Type | Lines | Purpose |
|------|-------|---------|
| `TextDocumentSyncKind` | 689-699 | None, Full, Incremental |
| `TextDocumentSyncOptions` | 701-707 | Sync configuration |
| `SignatureHelp` | 714-719 | Signature information |
| `SignatureInformation` | 721-726 | Single signature |
| `ParameterInformation` | 728-732 | Parameter details |

## Known Issues / Potential Bugs

None identified - this file contains only type definitions.

## Dead/Unused Code

Some capability types may not be fully utilized:
- `CodeLensClientCapabilities` - CodeLens not implemented
- `DocumentHighlightClientCaps` - Highlight not implemented

## Integration Points

- Used throughout `lsp/` package
- Serialized to JSON for LSP protocol
- Parsed from JSON responses

---
*Last Updated: Reference guide creation*
