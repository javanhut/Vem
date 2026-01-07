# features.go

**Path:** `/home/javanhut/Development/Vem/internal/lsp/features.go`
**Lines:** 451
**Purpose:** LSP feature implementations (completion, hover, definition, etc.)

## Overview

Implements all LSP feature requests with:
- Capability checking before requests
- 10 second timeout for all requests
- Flexible response parsing (handles multiple formats)
- Support checks for UI integration

## Code Blocks

### Lines 1-11: Package and Constants

```go
package lsp
import (...)

const defaultTimeout = 10 * time.Second
```

### Lines 13-34: GotoDefinition

#### GotoDefinition (Lines 13-34)
Requests definition location for symbol:
1. Checks `DefinitionProvider` capability
2. Sends `textDocument/definition` request
3. Parses response via `parseLocationResult()`

### Lines 36-72: GetHover

#### GetHover (Lines 36-72)
Requests hover information:
1. Checks `HoverProvider` capability
2. Sends `textDocument/hover` request
3. Handles null response
4. Parses as `Hover` or plain `MarkupContent`

### Lines 74-100: FindReferences

#### FindReferences (Lines 74-100)
Finds all references to symbol:
1. Checks `ReferencesProvider` capability
2. Sends `textDocument/references` request
3. Includes `includeDeclaration` context option

### Lines 102-152: GetCompletion

#### GetCompletion (Lines 102-152)
Requests completions at position:
1. Checks `CompletionProvider` capability
2. Creates completion context (trigger kind + character)
3. Sends `textDocument/completion` request
4. Parses as `CompletionList` or `[]CompletionItem`

**Trigger Kinds:**
- `Invoked` (1): Manual trigger (Ctrl+Space)
- `TriggerCharacter` (2): Triggered by character (e.g., `.`)

### Lines 154-169: ResolveCompletion

#### ResolveCompletion (Lines 154-169)
Resolves additional completion item details:
1. Checks `ResolveProvider` capability
2. Sends `completionItem/resolve` request
3. Returns original item if not supported

### Lines 171-210: PrepareRename

#### PrepareRename (Lines 171-210)
Checks if rename is valid at position:
1. Checks `RenameProvider` capability
2. Sends `textDocument/prepareRename` request
3. Parses response as:
   - Simple `Range`
   - `{ range: Range, placeholder: string }`

### Lines 212-236: Rename

#### Rename (Lines 212-236)
Performs symbol rename:
1. Checks `RenameProvider` capability
2. Sends `textDocument/rename` request
3. Returns `WorkspaceEdit` with changes

### Lines 238-275: GetCodeActions

#### GetCodeActions (Lines 238-275)
Requests code actions for range:
1. Checks `CodeActionProvider` capability
2. Sends range and associated diagnostics
3. Returns list of `CodeAction` items

### Lines 277-304: FormatDocument

#### FormatDocument (Lines 277-304)
Formats entire document:
1. Checks `DocumentFormattingProvider` capability
2. Sends formatting options:
   - TabSize: 4
   - InsertSpaces: false (use tabs)
   - TrimTrailingWhitespace: true
   - InsertFinalNewline: true

### Lines 306-327: GetSignatureHelp

#### GetSignatureHelp (Lines 306-327)
Requests signature help:
1. Checks `SignatureHelpProvider` capability
2. Sends `textDocument/signatureHelp` request
3. Returns function signature information

### Lines 329-373: Type/Implementation Navigation

| Function | Lines | Purpose |
|----------|-------|---------|
| `GotoTypeDefinition()` | 329-350 | Navigate to type definition |
| `GotoImplementation()` | 352-373 | Navigate to implementation |

### Lines 375-407: parseLocationResult

#### parseLocationResult (Lines 375-407)
Parses various location response formats:
1. Tries single `Location`
2. Tries `[]Location`
3. Tries `[]LocationLink` and converts to `Location`

### Lines 409-450: Support Checks

| Function | Lines | Checks |
|----------|-------|--------|
| `SupportsDefinition()` | 409-412 | DefinitionProvider |
| `SupportsHover()` | 414-417 | HoverProvider |
| `SupportsReferences()` | 419-422 | ReferencesProvider |
| `SupportsCompletion()` | 424-427 | CompletionProvider |
| `SupportsRename()` | 429-432 | RenameProvider |
| `SupportsCodeAction()` | 434-437 | CodeActionProvider |
| `SupportsFormatting()` | 439-442 | DocumentFormattingProvider |
| `GetTriggerCharacters()` | 444-450 | Completion triggers |

## Known Issues / Potential Bugs

1. **Line 11: 10 second timeout may be too short**
   - Complex operations (references, rename) may need more time
   - Consider per-feature timeouts

2. **Line 289-294: Hardcoded formatting options**
   - TabSize: 4 and InsertSpaces: false are hardcoded
   - Should be configurable per-project

## Dead/Unused Code

None identified.

## Integration Points

- Called from `appcore/lsp_actions.go`
- Results displayed by `appcore/lsp_rendering.go`
- Capability checks used by UI to show/hide features

## Feature Matrix

| Feature | Method | Capability |
|---------|--------|------------|
| Go to Definition | textDocument/definition | DefinitionProvider |
| Hover | textDocument/hover | HoverProvider |
| Find References | textDocument/references | ReferencesProvider |
| Completion | textDocument/completion | CompletionProvider |
| Rename | textDocument/rename | RenameProvider |
| Code Actions | textDocument/codeAction | CodeActionProvider |
| Format | textDocument/formatting | DocumentFormattingProvider |
| Signature Help | textDocument/signatureHelp | SignatureHelpProvider |

---
*Last Updated: Reference guide creation*
