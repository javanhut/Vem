package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Default timeout for LSP requests
const defaultTimeout = 10 * time.Second

// GotoDefinition requests the definition location for a symbol.
func (si *ServerInstance) GotoDefinition(filePath string, line, col int) ([]Location, error) {
	// Check if server supports definition
	if si.Capabilities.DefinitionProvider == nil {
		return nil, fmt.Errorf("server does not support go-to-definition")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := TextDocumentPositionParams{
		TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
		Position:     Position{Line: line, Character: col},
	}

	var result json.RawMessage
	if err := si.Client.Call(ctx, "textDocument/definition", params, &result); err != nil {
		return nil, err
	}

	return parseLocationResult(result)
}

// GetHover requests hover information for a position.
func (si *ServerInstance) GetHover(filePath string, line, col int) (*Hover, error) {
	// Check if server supports hover
	if si.Capabilities.HoverProvider == nil {
		return nil, fmt.Errorf("server does not support hover")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := TextDocumentPositionParams{
		TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
		Position:     Position{Line: line, Character: col},
	}

	var result json.RawMessage
	if err := si.Client.Call(ctx, "textDocument/hover", params, &result); err != nil {
		return nil, err
	}

	if result == nil || string(result) == "null" {
		return nil, nil
	}

	// Parse hover result - can be Hover or just MarkupContent
	var hover Hover
	if err := json.Unmarshal(result, &hover); err != nil {
		// Try parsing as simple content
		var content MarkupContent
		if err := json.Unmarshal(result, &content); err != nil {
			return nil, fmt.Errorf("parse hover: %w", err)
		}
		hover.Contents = content
	}

	return &hover, nil
}

// FindReferences finds all references to a symbol.
func (si *ServerInstance) FindReferences(filePath string, line, col int, includeDecl bool) ([]Location, error) {
	// Check if server supports references
	if si.Capabilities.ReferencesProvider == nil {
		return nil, fmt.Errorf("server does not support find references")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := ReferenceParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
			Position:     Position{Line: line, Character: col},
		},
		Context: ReferenceContext{
			IncludeDeclaration: includeDecl,
		},
	}

	var result []Location
	if err := si.Client.Call(ctx, "textDocument/references", params, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCompletion requests completions at a position.
func (si *ServerInstance) GetCompletion(filePath string, line, col int, triggerChar string) (*CompletionList, error) {
	// Check if server supports completion
	if si.Capabilities.CompletionProvider == nil {
		return nil, fmt.Errorf("server does not support completion")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := CompletionParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
			Position:     Position{Line: line, Character: col},
		},
	}

	// Add context if triggered by a character
	if triggerChar != "" {
		params.Context = &CompletionContext{
			TriggerKind:      CompletionTriggerKindTriggerCharacter,
			TriggerCharacter: triggerChar,
		}
	} else {
		params.Context = &CompletionContext{
			TriggerKind: CompletionTriggerKindInvoked,
		}
	}

	var result json.RawMessage
	if err := si.Client.Call(ctx, "textDocument/completion", params, &result); err != nil {
		return nil, err
	}

	if result == nil || string(result) == "null" {
		return &CompletionList{}, nil
	}

	// Result can be CompletionList or []CompletionItem
	var list CompletionList
	if err := json.Unmarshal(result, &list); err != nil {
		// Try parsing as array
		var items []CompletionItem
		if err := json.Unmarshal(result, &items); err != nil {
			return nil, fmt.Errorf("parse completion: %w", err)
		}
		list.Items = items
	}

	return &list, nil
}

// ResolveCompletion resolves additional information for a completion item.
func (si *ServerInstance) ResolveCompletion(item CompletionItem) (*CompletionItem, error) {
	if si.Capabilities.CompletionProvider == nil || !si.Capabilities.CompletionProvider.ResolveProvider {
		return &item, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var result CompletionItem
	if err := si.Client.Call(ctx, "completionItem/resolve", item, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// PrepareRename checks if rename is valid at position and returns the range to rename.
func (si *ServerInstance) PrepareRename(filePath string, line, col int) (*Range, string, error) {
	// Check if server supports rename
	if si.Capabilities.RenameProvider == nil {
		return nil, "", fmt.Errorf("server does not support rename")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := TextDocumentPositionParams{
		TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
		Position:     Position{Line: line, Character: col},
	}

	var result json.RawMessage
	if err := si.Client.Call(ctx, "textDocument/prepareRename", params, &result); err != nil {
		return nil, "", err
	}

	if result == nil || string(result) == "null" {
		return nil, "", fmt.Errorf("cannot rename at this position")
	}

	// Result can be Range or { range: Range, placeholder: string }
	var rangeResult Range
	if err := json.Unmarshal(result, &rangeResult); err == nil {
		return &rangeResult, "", nil
	}

	var prepareResult struct {
		Range       Range  `json:"range"`
		Placeholder string `json:"placeholder"`
	}
	if err := json.Unmarshal(result, &prepareResult); err != nil {
		return nil, "", fmt.Errorf("parse prepare rename: %w", err)
	}

	return &prepareResult.Range, prepareResult.Placeholder, nil
}

// Rename performs a symbol rename.
func (si *ServerInstance) Rename(filePath string, line, col int, newName string) (*WorkspaceEdit, error) {
	// Check if server supports rename
	if si.Capabilities.RenameProvider == nil {
		return nil, fmt.Errorf("server does not support rename")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := RenameParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
			Position:     Position{Line: line, Character: col},
		},
		NewName: newName,
	}

	var result WorkspaceEdit
	if err := si.Client.Call(ctx, "textDocument/rename", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCodeActions requests code actions for a range.
func (si *ServerInstance) GetCodeActions(filePath string, startLine, startCol, endLine, endCol int, diagnostics []Diagnostic) ([]CodeAction, error) {
	// Check if server supports code actions
	if si.Capabilities.CodeActionProvider == nil {
		return nil, fmt.Errorf("server does not support code actions")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := CodeActionParams{
		TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
		Range: Range{
			Start: Position{Line: startLine, Character: startCol},
			End:   Position{Line: endLine, Character: endCol},
		},
		Context: CodeActionContext{
			Diagnostics: diagnostics,
		},
	}

	var result json.RawMessage
	if err := si.Client.Call(ctx, "textDocument/codeAction", params, &result); err != nil {
		return nil, err
	}

	if result == nil || string(result) == "null" {
		return nil, nil
	}

	// Result can be []CodeAction or [](CodeAction | Command)
	var actions []CodeAction
	if err := json.Unmarshal(result, &actions); err != nil {
		return nil, fmt.Errorf("parse code actions: %w", err)
	}

	return actions, nil
}

// FormatDocument formats the entire document.
func (si *ServerInstance) FormatDocument(filePath string) ([]TextEdit, error) {
	// Check if server supports formatting
	if si.Capabilities.DocumentFormattingProvider == nil {
		return nil, fmt.Errorf("server does not support document formatting")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := DocumentFormattingParams{
		TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
		Options: FormattingOptions{
			TabSize:                4,
			InsertSpaces:           false, // Use tabs
			TrimTrailingWhitespace: true,
			InsertFinalNewline:     true,
			TrimFinalNewlines:      true,
		},
	}

	var result []TextEdit
	if err := si.Client.Call(ctx, "textDocument/formatting", params, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetSignatureHelp requests signature help at a position.
func (si *ServerInstance) GetSignatureHelp(filePath string, line, col int) (*SignatureHelp, error) {
	// Check if server supports signature help
	if si.Capabilities.SignatureHelpProvider == nil {
		return nil, fmt.Errorf("server does not support signature help")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := TextDocumentPositionParams{
		TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
		Position:     Position{Line: line, Character: col},
	}

	var result SignatureHelp
	if err := si.Client.Call(ctx, "textDocument/signatureHelp", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GotoTypeDefinition requests the type definition location for a symbol.
func (si *ServerInstance) GotoTypeDefinition(filePath string, line, col int) ([]Location, error) {
	// Check if server supports type definition
	if si.Capabilities.TypeDefinitionProvider == nil {
		return nil, fmt.Errorf("server does not support go-to-type-definition")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := TextDocumentPositionParams{
		TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
		Position:     Position{Line: line, Character: col},
	}

	var result json.RawMessage
	if err := si.Client.Call(ctx, "textDocument/typeDefinition", params, &result); err != nil {
		return nil, err
	}

	return parseLocationResult(result)
}

// GotoImplementation requests the implementation location for a symbol.
func (si *ServerInstance) GotoImplementation(filePath string, line, col int) ([]Location, error) {
	// Check if server supports implementation
	if si.Capabilities.ImplementationProvider == nil {
		return nil, fmt.Errorf("server does not support go-to-implementation")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	params := TextDocumentPositionParams{
		TextDocument: TextDocumentIdentifier{URI: FilePathToURI(filePath)},
		Position:     Position{Line: line, Character: col},
	}

	var result json.RawMessage
	if err := si.Client.Call(ctx, "textDocument/implementation", params, &result); err != nil {
		return nil, err
	}

	return parseLocationResult(result)
}

// parseLocationResult parses various location response formats.
func parseLocationResult(result json.RawMessage) ([]Location, error) {
	if result == nil || string(result) == "null" {
		return nil, nil
	}

	// Try parsing as single Location
	var loc Location
	if err := json.Unmarshal(result, &loc); err == nil && loc.URI != "" {
		return []Location{loc}, nil
	}

	// Try parsing as []Location
	var locs []Location
	if err := json.Unmarshal(result, &locs); err == nil {
		return locs, nil
	}

	// Try parsing as []LocationLink
	var links []LocationLink
	if err := json.Unmarshal(result, &links); err == nil {
		locations := make([]Location, len(links))
		for i, link := range links {
			locations[i] = Location{
				URI:   link.TargetURI,
				Range: link.TargetSelectionRange,
			}
		}
		return locations, nil
	}

	return nil, fmt.Errorf("unexpected location response format")
}

// SupportsDefinition returns true if the server supports go-to-definition.
func (si *ServerInstance) SupportsDefinition() bool {
	return si.Capabilities.DefinitionProvider != nil
}

// SupportsHover returns true if the server supports hover.
func (si *ServerInstance) SupportsHover() bool {
	return si.Capabilities.HoverProvider != nil
}

// SupportsReferences returns true if the server supports find references.
func (si *ServerInstance) SupportsReferences() bool {
	return si.Capabilities.ReferencesProvider != nil
}

// SupportsCompletion returns true if the server supports completion.
func (si *ServerInstance) SupportsCompletion() bool {
	return si.Capabilities.CompletionProvider != nil
}

// SupportsRename returns true if the server supports rename.
func (si *ServerInstance) SupportsRename() bool {
	return si.Capabilities.RenameProvider != nil
}

// SupportsCodeAction returns true if the server supports code actions.
func (si *ServerInstance) SupportsCodeAction() bool {
	return si.Capabilities.CodeActionProvider != nil
}

// SupportsFormatting returns true if the server supports document formatting.
func (si *ServerInstance) SupportsFormatting() bool {
	return si.Capabilities.DocumentFormattingProvider != nil
}

// GetTriggerCharacters returns the completion trigger characters for this server.
func (si *ServerInstance) GetTriggerCharacters() []string {
	if si.Capabilities.CompletionProvider == nil {
		return nil
	}
	return si.Capabilities.CompletionProvider.TriggerCharacters
}
