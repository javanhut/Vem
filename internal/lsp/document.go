package lsp

// OpenDocument notifies the server that a document was opened.
func (si *ServerInstance) OpenDocument(filePath, content string) error {
	uri := FilePathToURI(filePath)
	langID := LanguageID(filePath)

	si.DocumentsMu.Lock()
	si.Documents[uri] = &DocumentState{
		URI:        uri,
		Version:    1,
		Content:    content,
		LanguageID: langID,
	}
	si.DocumentsMu.Unlock()

	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: langID,
			Version:    1,
			Text:       content,
		},
	}

	return si.Client.Notify("textDocument/didOpen", params)
}

// ChangeDocument notifies the server of document changes.
// Uses full document sync for maximum compatibility.
func (si *ServerInstance) ChangeDocument(filePath, content string) error {
	uri := FilePathToURI(filePath)

	si.DocumentsMu.Lock()
	state, ok := si.Documents[uri]
	if !ok {
		si.DocumentsMu.Unlock()
		// Document not open, open it first
		return si.OpenDocument(filePath, content)
	}

	// Increment version
	state.Version++
	state.Content = content
	version := state.Version
	si.DocumentsMu.Unlock()

	// Use full document sync (works with all servers)
	params := DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{URI: uri},
			Version:                version,
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{Text: content}, // Full content sync
		},
	}

	return si.Client.Notify("textDocument/didChange", params)
}

// SaveDocument notifies the server that a document was saved.
func (si *ServerInstance) SaveDocument(filePath string) error {
	uri := FilePathToURI(filePath)

	si.DocumentsMu.RLock()
	state, ok := si.Documents[uri]
	var content string
	if ok {
		content = state.Content
	}
	si.DocumentsMu.RUnlock()

	if !ok {
		return nil // Document not tracked
	}

	params := DidSaveTextDocumentParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
		Text:         content, // Include text if server supports it
	}

	return si.Client.Notify("textDocument/didSave", params)
}

// CloseDocument notifies the server that a document was closed.
func (si *ServerInstance) CloseDocument(filePath string) error {
	uri := FilePathToURI(filePath)

	si.DocumentsMu.Lock()
	delete(si.Documents, uri)
	si.DocumentsMu.Unlock()

	// Also clear diagnostics for this file
	si.DiagnosticsMu.Lock()
	delete(si.Diagnostics, uri)
	si.DiagnosticsMu.Unlock()

	params := DidCloseTextDocumentParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
	}

	return si.Client.Notify("textDocument/didClose", params)
}

// IsDocumentOpen returns true if a document is tracked by this server.
func (si *ServerInstance) IsDocumentOpen(filePath string) bool {
	uri := FilePathToURI(filePath)

	si.DocumentsMu.RLock()
	_, ok := si.Documents[uri]
	si.DocumentsMu.RUnlock()

	return ok
}

// GetDocumentVersion returns the current version of a document.
func (si *ServerInstance) GetDocumentVersion(filePath string) int {
	uri := FilePathToURI(filePath)

	si.DocumentsMu.RLock()
	defer si.DocumentsMu.RUnlock()

	if state, ok := si.Documents[uri]; ok {
		return state.Version
	}
	return 0
}

// SyncDocument ensures a document is open and synced with the given content.
// This is useful when first needing to use LSP features on a file.
func (si *ServerInstance) SyncDocument(filePath, content string) error {
	uri := FilePathToURI(filePath)

	si.DocumentsMu.RLock()
	state, ok := si.Documents[uri]
	si.DocumentsMu.RUnlock()

	if !ok {
		// Document not open, open it
		return si.OpenDocument(filePath, content)
	}

	// Check if content differs
	if state.Content != content {
		return si.ChangeDocument(filePath, content)
	}

	return nil
}
