package appcore

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/javanhut/vem/internal/editor"
	"github.com/javanhut/vem/internal/lsp"
)

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

func (s *appState) setLSPHint(filePath string, err error) {
	if err == nil {
		s.lspHint = ""
		s.lspHintFile = ""
		return
	}

	cfg := lsp.GetConfigForFile(filePath)
	msg := err.Error()
	stderrHint := ""
	if idx := strings.LastIndex(msg, "stderr: "); idx != -1 {
		stderrHint = strings.TrimSpace(msg[idx+len("stderr: "):])
	}

	switch {
	case strings.Contains(stderrHint, "Unknown binary 'rust-analyzer'"):
		s.lspHint = "Install rust-analyzer via: rustup component add rust-analyzer"
	case strings.Contains(msg, "not installed"):
		if cfg != nil {
			s.lspHint = fmt.Sprintf("Install %s and ensure it is on PATH", cfg.Name)
		} else {
			s.lspHint = "Language server not installed"
		}
	case strings.Contains(msg, "project root"):
		if cfg != nil && len(cfg.RootPatterns) > 0 {
			s.lspHint = fmt.Sprintf("Project root not found (need %s)", strings.Join(cfg.RootPatterns, ", "))
		} else {
			s.lspHint = "Project root not found"
		}
	case strings.Contains(msg, "no language server configured"):
		s.lspHint = "No LSP configured for this file type"
	default:
		if cfg != nil {
			if stderrHint != "" {
				s.lspHint = fmt.Sprintf("Failed to start %s (%s)", cfg.Name, stderrHint)
			} else {
				s.lspHint = fmt.Sprintf("Failed to start %s (%s)", cfg.Name, msg)
			}
		} else {
			if stderrHint != "" {
				s.lspHint = fmt.Sprintf("LSP start failed (%s)", stderrHint)
			} else {
				s.lspHint = fmt.Sprintf("LSP start failed (%s)", msg)
			}
		}
	}

	s.lspHintFile = filePath
}

// handleLSPGotoDefinition navigates to symbol definition.
func (s *appState) handleLSPGotoDefinition() {
	if !s.lspEnabled || s.lspManager == nil {
		s.status = "LSP not enabled"
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		s.status = "No file open"
		return
	}

	// Ensure document is synced
	s.syncDocumentWithLSP(buf)

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		s.status = fmt.Sprintf("LSP: %v", err)
		return
	}

	cursor := buf.Cursor()
	locations, err := server.GotoDefinition(buf.FilePath(), cursor.Line, cursor.Col)
	if err != nil {
		s.status = fmt.Sprintf("Definition: %v", err)
		return
	}

	if len(locations) == 0 {
		s.status = "No definition found"
		return
	}

	// Navigate to first location
	loc := locations[0]
	s.navigateToLocation(loc)
}

// handleLSPHover shows hover information.
func (s *appState) handleLSPHover() {
	if !s.lspEnabled || s.lspManager == nil {
		s.status = "LSP not enabled"
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	// Ensure document is synced
	s.syncDocumentWithLSP(buf)

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		return
	}

	cursor := buf.Cursor()
	hover, err := server.GetHover(buf.FilePath(), cursor.Line, cursor.Col)
	if err != nil || hover == nil || hover.Contents.Value == "" {
		s.status = "No hover information"
		return
	}

	s.hoverActive = true
	s.hoverContent = hover.Contents.Value
	s.hoverRange = hover.Range
	s.status = ""
}

// handleLSPReferences finds all references.
func (s *appState) handleLSPReferences() {
	if !s.lspEnabled || s.lspManager == nil {
		s.status = "LSP not enabled"
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	// Ensure document is synced
	s.syncDocumentWithLSP(buf)

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		s.status = fmt.Sprintf("LSP: %v", err)
		return
	}

	cursor := buf.Cursor()
	refs, err := server.FindReferences(buf.FilePath(), cursor.Line, cursor.Col, true)
	if err != nil {
		s.status = fmt.Sprintf("References: %v", err)
		return
	}

	if len(refs) == 0 {
		s.status = "No references found"
		return
	}

	s.referencesActive = true
	s.referencesItems = refs
	s.referencesIndex = 0
	s.status = fmt.Sprintf("Found %d references", len(refs))
}

// handleLSPRename renames a symbol.
func (s *appState) handleLSPRename(newName string) {
	if !s.lspEnabled || s.lspManager == nil {
		s.status = "LSP not enabled"
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	// Ensure document is synced
	s.syncDocumentWithLSP(buf)

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		s.status = fmt.Sprintf("LSP: %v", err)
		return
	}

	cursor := buf.Cursor()
	edit, err := server.Rename(buf.FilePath(), cursor.Line, cursor.Col, newName)
	if err != nil {
		s.status = fmt.Sprintf("Rename: %v", err)
		return
	}

	// Apply workspace edit
	filesChanged := s.applyWorkspaceEdit(edit)
	s.status = fmt.Sprintf("Renamed in %d file(s)", filesChanged)
}

// handleLSPFormat formats the document.
func (s *appState) handleLSPFormat() {
	if !s.lspEnabled || s.lspManager == nil {
		s.status = "LSP not enabled"
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	// Ensure document is synced
	s.syncDocumentWithLSP(buf)

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		s.status = fmt.Sprintf("LSP: %v", err)
		return
	}

	edits, err := server.FormatDocument(buf.FilePath())
	if err != nil {
		s.status = fmt.Sprintf("Format: %v", err)
		return
	}

	if len(edits) == 0 {
		s.status = "Document already formatted"
		return
	}

	// Apply edits
	s.applyTextEdits(buf, edits)
	s.status = "Document formatted"
}

// handleLSPCodeAction shows available code actions.
func (s *appState) handleLSPCodeAction() {
	if !s.lspEnabled || s.lspManager == nil {
		s.status = "LSP not enabled"
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	// Ensure document is synced
	s.syncDocumentWithLSP(buf)

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		s.status = fmt.Sprintf("LSP: %v", err)
		return
	}

	cursor := buf.Cursor()

	// Get diagnostics at cursor
	var diags []lsp.Diagnostic
	if fileDiags, ok := s.lspDiagnostics[buf.FilePath()]; ok {
		for _, d := range fileDiags {
			if d.Range.Start.Line == cursor.Line {
				diags = append(diags, d)
			}
		}
	}

	actions, err := server.GetCodeActions(buf.FilePath(), cursor.Line, cursor.Col, cursor.Line, cursor.Col, diags)
	if err != nil {
		s.status = fmt.Sprintf("Code actions: %v", err)
		return
	}

	if len(actions) == 0 {
		s.status = "No code actions available"
		return
	}

	s.codeActionsActive = true
	s.codeActionItems = actions
	s.codeActionIndex = 0
}

// handleLSPCompletion triggers completion.
func (s *appState) handleLSPCompletion() {
	if !s.lspEnabled || s.lspManager == nil {
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	// Ensure document is synced
	s.syncDocumentWithLSP(buf)

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		return
	}

	cursor := buf.Cursor()
	completions, err := server.GetCompletion(buf.FilePath(), cursor.Line, cursor.Col, "")
	if err != nil || completions == nil || len(completions.Items) == 0 {
		s.completionActive = false
		return
	}

	s.completionActive = true
	s.completionItems = completions.Items
	s.completionIndex = 0
	s.completionTrigger = cursor
	s.completionVersion++
	s.completionResolved = make(map[int]bool)
	s.resolveCompletionItem(0)
}

// handleLSPCompletionAccept accepts the selected completion.
func (s *appState) handleLSPCompletionAccept() {
	if !s.completionActive || len(s.completionItems) == 0 {
		return
	}

	buf := s.activeBuffer()
	if buf == nil {
		return
	}

	item := s.completionItems[s.completionIndex]

	// Determine what text to insert
	insertText := item.InsertText
	if insertText == "" {
		insertText = item.Label
	}

	// If there's a text edit, use that instead
	if item.TextEdit != nil {
		// Delete the range and insert new text
		buf.DeleteCharRange(
			item.TextEdit.Range.Start.Line,
			item.TextEdit.Range.Start.Character,
			item.TextEdit.Range.End.Line,
			item.TextEdit.Range.End.Character,
		)
		buf.InsertText(item.TextEdit.NewText)
	} else {
		// Simple insert: delete word prefix and insert completion
		// First, find the start of the word being completed
		cursor := buf.Cursor()
		line := buf.Line(cursor.Line)
		runes := []rune(line)

		// Find start of word
		wordStart := cursor.Col
		for wordStart > 0 && isWordChar(runes[wordStart-1]) {
			wordStart--
		}

		// Delete from word start to cursor
		if wordStart < cursor.Col {
			buf.DeleteCharRange(cursor.Line, wordStart, cursor.Line, cursor.Col)
		}

		// Insert completion text
		buf.InsertText(insertText)
	}

	// Clear completion state
	s.completionActive = false
	s.completionItems = nil

	// Update syntax highlighting
	s.invalidateSyntaxCache()
}

// handleLSPCompletionCancel cancels completion.
func (s *appState) handleLSPCompletionCancel() {
	s.completionActive = false
	s.completionItems = nil
}

// handleLSPCompletionNext moves to next completion item.
func (s *appState) handleLSPCompletionNext() {
	if !s.completionActive || len(s.completionItems) == 0 {
		return
	}
	s.completionIndex = (s.completionIndex + 1) % len(s.completionItems)
	s.resolveCompletionItem(s.completionIndex)
}

// handleLSPCompletionPrev moves to previous completion item.
func (s *appState) handleLSPCompletionPrev() {
	if !s.completionActive || len(s.completionItems) == 0 {
		return
	}
	s.completionIndex--
	if s.completionIndex < 0 {
		s.completionIndex = len(s.completionItems) - 1
	}
	s.resolveCompletionItem(s.completionIndex)
}

func (s *appState) resolveCompletionItem(index int) {
	if !s.lspEnabled || s.lspManager == nil || !s.completionActive {
		return
	}
	if index < 0 || index >= len(s.completionItems) {
		return
	}

	if s.completionResolved == nil {
		s.completionResolved = make(map[int]bool)
	}
	if s.completionResolved[index] {
		return
	}
	s.completionResolved[index] = true
	version := s.completionVersion
	item := s.completionItems[index]

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}
	filePath := buf.FilePath()

	go func() {
		server, err := s.lspManager.GetServerForFile(filePath)
		if err != nil || server == nil {
			return
		}

		resolved, err := server.ResolveCompletion(item)
		if err != nil || resolved == nil {
			return
		}

		if s.completionVersion != version || index >= len(s.completionItems) {
			return
		}
		if s.completionItems[index].Label != item.Label {
			return
		}

		s.completionItems[index] = *resolved
		if s.window != nil {
			s.window.Invalidate()
		}
	}()
}

// handleLSPNextDiagnostic moves to the next diagnostic.
func (s *appState) handleLSPNextDiagnostic() {
	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	diags, ok := s.lspDiagnostics[buf.FilePath()]
	if !ok || len(diags) == 0 {
		s.status = "No diagnostics"
		return
	}

	cursor := buf.Cursor()

	// Find next diagnostic after cursor
	for _, d := range diags {
		if d.Range.Start.Line > cursor.Line ||
			(d.Range.Start.Line == cursor.Line && d.Range.Start.Character > cursor.Col) {
			buf.MoveToLine(d.Range.Start.Line)
			s.status = d.Message
			s.ensureCursorVisible(0)
			return
		}
	}

	// Wrap around to first diagnostic
	if len(diags) > 0 {
		d := diags[0]
		buf.MoveToLine(d.Range.Start.Line)
		s.status = d.Message
		s.ensureCursorVisible(0)
	}
}

// handleLSPPrevDiagnostic moves to the previous diagnostic.
func (s *appState) handleLSPPrevDiagnostic() {
	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	diags, ok := s.lspDiagnostics[buf.FilePath()]
	if !ok || len(diags) == 0 {
		s.status = "No diagnostics"
		return
	}

	cursor := buf.Cursor()

	// Find previous diagnostic before cursor
	for i := len(diags) - 1; i >= 0; i-- {
		d := diags[i]
		if d.Range.Start.Line < cursor.Line ||
			(d.Range.Start.Line == cursor.Line && d.Range.Start.Character < cursor.Col) {
			buf.MoveToLine(d.Range.Start.Line)
			s.status = d.Message
			s.ensureCursorVisible(0)
			return
		}
	}

	// Wrap around to last diagnostic
	if len(diags) > 0 {
		d := diags[len(diags)-1]
		buf.MoveToLine(d.Range.Start.Line)
		s.status = d.Message
		s.ensureCursorVisible(0)
	}
}

// handleLSPDismissHover hides the hover tooltip.
func (s *appState) handleLSPDismissHover() {
	s.hoverActive = false
	s.hoverContent = ""
	s.hoverRange = nil
}

// handleLSPDismissReferences hides the references list.
func (s *appState) handleLSPDismissReferences() {
	s.referencesActive = false
	s.referencesItems = nil
}

// handleLSPDismissCodeActions hides the code actions menu.
func (s *appState) handleLSPDismissCodeActions() {
	s.codeActionsActive = false
	s.codeActionItems = nil
}

// handleLSPReferencesNext moves to the next reference.
func (s *appState) handleLSPReferencesNext() {
	if !s.referencesActive || len(s.referencesItems) == 0 {
		return
	}
	s.referencesIndex = (s.referencesIndex + 1) % len(s.referencesItems)
}

// handleLSPReferencesPrev moves to the previous reference.
func (s *appState) handleLSPReferencesPrev() {
	if !s.referencesActive || len(s.referencesItems) == 0 {
		return
	}
	s.referencesIndex--
	if s.referencesIndex < 0 {
		s.referencesIndex = len(s.referencesItems) - 1
	}
}

// handleLSPReferencesOpen opens the selected reference.
func (s *appState) handleLSPReferencesOpen() {
	if !s.referencesActive || len(s.referencesItems) == 0 {
		return
	}

	loc := s.referencesItems[s.referencesIndex]
	s.navigateToLocation(loc)
	s.referencesActive = false
}

// handleLSPCodeActionSelect applies the selected code action.
func (s *appState) handleLSPCodeActionSelect() {
	if !s.codeActionsActive || len(s.codeActionItems) == 0 {
		return
	}

	action := s.codeActionItems[s.codeActionIndex]

	// Apply workspace edit if present
	if action.Edit != nil {
		filesChanged := s.applyWorkspaceEdit(action.Edit)
		s.status = fmt.Sprintf("Applied: %s (%d files)", action.Title, filesChanged)
	} else {
		s.status = fmt.Sprintf("Applied: %s", action.Title)
	}

	s.codeActionsActive = false
	s.codeActionItems = nil
}

// handleLSPStatus shows LSP server status.
func (s *appState) handleLSPStatus() {
	if !s.lspEnabled || s.lspManager == nil {
		s.status = "LSP disabled"
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		s.status = "No file open"
		return
	}

	cfg := lsp.GetConfigForFile(buf.FilePath())
	if cfg == nil {
		s.status = "No language server for this file type"
		s.setLSPHint(buf.FilePath(), fmt.Errorf("no language server configured for this file type"))
		return
	}

	if !lsp.IsServerAvailable(cfg) {
		s.status = fmt.Sprintf("LSP: %s not installed", cfg.Name)
		s.setLSPHint(buf.FilePath(), fmt.Errorf("language server %s not installed", cfg.Name))
		return
	}

	if s.lspManager.IsServerRunningForFile(buf.FilePath()) {
		s.status = fmt.Sprintf("LSP: %s running", cfg.Name)
		s.setLSPHint(buf.FilePath(), nil)
	} else {
		s.status = fmt.Sprintf("LSP: %s available (not started)", cfg.Name)
		if s.lspHint == "" || s.lspHintFile != buf.FilePath() {
			s.setLSPHint(buf.FilePath(), fmt.Errorf("language server %s not started", cfg.Name))
		}
	}
}

// handleLSPRestart restarts the language server for the current file.
func (s *appState) handleLSPRestart() {
	if !s.lspEnabled || s.lspManager == nil {
		s.status = "LSP disabled"
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		s.status = "No file open"
		return
	}

	_, err := s.lspManager.RestartServerForFile(buf.FilePath())
	if err != nil {
		s.status = fmt.Sprintf("LSP restart failed: %v", err)
		return
	}

	s.status = "LSP server restarted"
}

// syncDocumentWithLSP ensures the buffer content is synced with the LSP server.
func (s *appState) syncDocumentWithLSP(buf *editor.Buffer) {
	if !s.lspEnabled || s.lspManager == nil || buf == nil || buf.FilePath() == "" {
		return
	}

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		return
	}

	// Sync the document
	content := buf.GetContent()
	_ = server.SyncDocument(buf.FilePath(), content)
}

// navigateToLocation navigates to a location.
func (s *appState) navigateToLocation(loc lsp.Location) {
	targetPath := lsp.URIToFilePath(loc.URI)
	buf := s.activeBuffer()

	// Open file if different
	if buf == nil || targetPath != buf.FilePath() {
		buf, err := s.bufferMgr.OpenFile(targetPath)
		if err != nil {
			s.status = fmt.Sprintf("Cannot open %s: %v", targetPath, err)
			return
		}
		s.setupLSPForBuffer(buf)
		// Update pane
		if s.paneManager != nil {
			if pane := s.paneManager.ActivePane(); pane != nil {
				pane.SetBufferIndex(s.bufferMgr.ActiveIndex())
			}
		}
		// Get the new buffer
		buf = s.activeBuffer()
	}

	// Jump to position
	if buf != nil {
		buf.MoveToLine(loc.Range.Start.Line)
		// Move to column
		for i := 0; i < loc.Range.Start.Character && buf.Cursor().Col < loc.Range.Start.Character; i++ {
			buf.MoveRight()
		}
		s.ensureCursorVisible(0)
	}

	s.status = fmt.Sprintf("Jumped to %s:%d", filepath.Base(targetPath), loc.Range.Start.Line+1)
}

// applyWorkspaceEdit applies a workspace edit to buffers.
func (s *appState) applyWorkspaceEdit(edit *lsp.WorkspaceEdit) int {
	if edit == nil {
		return 0
	}

	filesChanged := 0

	// Handle Changes (map of URI -> []TextEdit)
	if edit.Changes != nil {
		for uri, edits := range edit.Changes {
			filePath := lsp.URIToFilePath(uri)
			buf := s.getOrOpenBuffer(filePath)
			if buf != nil {
				s.applyTextEdits(buf, edits)
				filesChanged++
			}
		}
	}

	// Handle DocumentChanges (array of TextDocumentEdit)
	if edit.DocumentChanges != nil {
		for _, docEdit := range edit.DocumentChanges {
			filePath := lsp.URIToFilePath(docEdit.TextDocument.URI)
			buf := s.getOrOpenBuffer(filePath)
			if buf != nil {
				s.applyTextEdits(buf, docEdit.Edits)
				filesChanged++
			}
		}
	}

	return filesChanged
}

// getOrOpenBuffer gets an existing buffer or opens a file.
func (s *appState) getOrOpenBuffer(filePath string) *editor.Buffer {
	// Check if buffer is already open
	for i := 0; i < s.bufferMgr.BufferCount(); i++ {
		if buf := s.bufferMgr.GetBuffer(i); buf != nil && buf.FilePath() == filePath {
			s.setupLSPForBuffer(buf)
			return buf
		}
	}

	// Open the file
	_, err := s.bufferMgr.OpenFile(filePath)
	if err != nil {
		return nil
	}

	buf := s.bufferMgr.ActiveBuffer()
	s.setupLSPForBuffer(buf)
	return buf
}

// applyTextEdits applies text edits to a buffer.
func (s *appState) applyTextEdits(buf *editor.Buffer, edits []lsp.TextEdit) {
	if buf == nil || len(edits) == 0 {
		return
	}

	// Sort edits in reverse order to apply from end to start
	// This preserves positions for earlier edits
	sortedEdits := make([]lsp.TextEdit, len(edits))
	copy(sortedEdits, edits)
	sort.Slice(sortedEdits, func(i, j int) bool {
		if sortedEdits[i].Range.Start.Line != sortedEdits[j].Range.Start.Line {
			return sortedEdits[i].Range.Start.Line > sortedEdits[j].Range.Start.Line
		}
		return sortedEdits[i].Range.Start.Character > sortedEdits[j].Range.Start.Character
	})

	for _, edit := range sortedEdits {
		// Delete range
		buf.DeleteCharRange(
			edit.Range.Start.Line, edit.Range.Start.Character,
			edit.Range.End.Line, edit.Range.End.Character,
		)

		// Move cursor to start position
		buf.MoveToLine(edit.Range.Start.Line)
		for buf.Cursor().Col < edit.Range.Start.Character {
			buf.MoveRight()
		}

		// Insert new text
		if edit.NewText != "" {
			buf.InsertText(edit.NewText)
		}
	}

	s.invalidateSyntaxCache()
}

// isWordChar checks if a rune is a word character for completion.
func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}

func (s *appState) scheduleLSPChange(filePath, content string) {
	if !s.lspEnabled || s.lspManager == nil || filePath == "" {
		return
	}

	s.lspChangeMu.Lock()
	s.lspChangeContents[filePath] = content
	s.lspChangeSeq[filePath]++
	seq := s.lspChangeSeq[filePath]

	if timer, ok := s.lspChangeTimers[filePath]; ok {
		timer.Stop()
	}

	s.lspChangeTimers[filePath] = time.AfterFunc(lspChangeDebounce, func() {
		s.flushLSPChange(filePath, seq)
	})
	s.lspChangeMu.Unlock()
}

func (s *appState) flushLSPChange(filePath string, seq int) {
	s.lspChangeMu.Lock()
	if s.lspChangeSeq[filePath] != seq {
		s.lspChangeMu.Unlock()
		return
	}
	content, ok := s.lspChangeContents[filePath]
	s.lspChangeMu.Unlock()

	if !ok || !s.lspEnabled || s.lspManager == nil {
		return
	}

	if !s.lspManager.IsServerRunningForFile(filePath) {
		return
	}

	s.lspChangeMu.Lock()
	delete(s.lspChangeContents, filePath)
	delete(s.lspChangeTimers, filePath)
	delete(s.lspChangeSeq, filePath)
	s.lspChangeMu.Unlock()

	server, err := s.lspManager.GetServerForFile(filePath)
	if err != nil {
		return
	}

	_ = server.ChangeDocument(filePath, content)
}

func (s *appState) flushPendingLSPChange(filePath string) {
	s.lspChangeMu.Lock()
	content, ok := s.lspChangeContents[filePath]
	delete(s.lspChangeContents, filePath)
	if timer, okTimer := s.lspChangeTimers[filePath]; okTimer {
		timer.Stop()
		delete(s.lspChangeTimers, filePath)
	}
	delete(s.lspChangeSeq, filePath)
	s.lspChangeMu.Unlock()

	if !ok || !s.lspEnabled || s.lspManager == nil {
		return
	}

	server, err := s.lspManager.GetServerForFile(filePath)
	if err != nil {
		return
	}

	_ = server.ChangeDocument(filePath, content)
}

func (s *appState) cancelLSPChange(filePath string) {
	if filePath == "" {
		return
	}

	s.lspChangeMu.Lock()
	if timer, ok := s.lspChangeTimers[filePath]; ok {
		timer.Stop()
		delete(s.lspChangeTimers, filePath)
	}
	delete(s.lspChangeContents, filePath)
	delete(s.lspChangeSeq, filePath)
	s.lspChangeMu.Unlock()
}

func (s *appState) scheduleLSPCompletion(req lspCompletionRequest) {
	if !s.lspEnabled || s.lspManager == nil || req.filePath == "" {
		return
	}

	s.lspCompletionMu.Lock()
	s.lspCompletionReq[req.filePath] = req
	s.lspCompletionSeq[req.filePath]++
	seq := s.lspCompletionSeq[req.filePath]

	if timer, ok := s.lspCompletionTimers[req.filePath]; ok {
		timer.Stop()
	}

	s.lspCompletionTimers[req.filePath] = time.AfterFunc(lspCompletionDebounce, func() {
		s.runLSPCompletion(req.filePath, seq)
	})
	s.lspCompletionMu.Unlock()
}

func (s *appState) cancelLSPCompletion(filePath string) {
	if filePath == "" {
		return
	}

	s.lspCompletionMu.Lock()
	if timer, ok := s.lspCompletionTimers[filePath]; ok {
		timer.Stop()
		delete(s.lspCompletionTimers, filePath)
	}
	delete(s.lspCompletionReq, filePath)
	delete(s.lspCompletionSeq, filePath)
	s.lspCompletionMu.Unlock()
}

func (s *appState) runLSPCompletion(filePath string, seq int) {
	s.lspCompletionMu.Lock()
	if s.lspCompletionSeq[filePath] != seq {
		s.lspCompletionMu.Unlock()
		return
	}
	req, ok := s.lspCompletionReq[filePath]
	delete(s.lspCompletionReq, filePath)
	delete(s.lspCompletionTimers, filePath)
	delete(s.lspCompletionSeq, filePath)
	s.lspCompletionMu.Unlock()

	if !ok || !s.lspEnabled || s.lspManager == nil {
		return
	}

	if !s.lspManager.IsServerRunningForFile(filePath) {
		return
	}

	server, err := s.lspManager.GetServerForFile(filePath)
	if err != nil {
		return
	}

	completions, err := server.GetCompletion(filePath, req.line, req.col, req.trigger)
	if err != nil || completions == nil || len(completions.Items) == 0 {
		s.completionActive = false
		s.completionItems = nil
		if s.window != nil {
			s.window.Invalidate()
		}
		return
	}

	s.completionActive = true
	s.completionItems = completions.Items
	s.completionIndex = 0
	s.completionTrigger = editor.Cursor{Line: req.line, Col: req.col}
	s.completionVersion++
	s.completionResolved = make(map[int]bool)
	s.resolveCompletionItem(0)

	if s.window != nil {
		s.window.Invalidate()
	}
}

func (s *appState) maybeTriggerLSPCompletion(inserted string) {
	if !s.lspEnabled || s.lspManager == nil || inserted == "" {
		return
	}

	buf := s.activeBuffer()
	if buf == nil || buf.FilePath() == "" {
		return
	}

	if !s.lspManager.IsServerRunningForFile(buf.FilePath()) {
		return
	}

	runes := []rune(inserted)
	last := runes[len(runes)-1]
	if last == '\n' || last == '\r' || unicode.IsSpace(last) {
		s.completionActive = false
		s.completionItems = nil
		return
	}

	trigger := ""
	cfg := lsp.GetConfigForFile(buf.FilePath())
	if cfg != nil {
		for _, ch := range cfg.TriggerChars {
			if ch == string(last) {
				trigger = ch
				break
			}
		}
	}

	if trigger == "" && !isWordChar(last) {
		s.completionActive = false
		s.completionItems = nil
		return
	}

	cursor := buf.Cursor()
	if trigger == "" && isWordChar(last) {
		line := []rune(buf.Line(cursor.Line))
		if cursor.Col > len(line) {
			return
		}
		start := cursor.Col
		for start > 0 && isWordChar(line[start-1]) {
			start--
		}
		if cursor.Col-start < 1 {
			return
		}
	}

	s.scheduleLSPCompletion(lspCompletionRequest{
		filePath: buf.FilePath(),
		line:     cursor.Line,
		col:      cursor.Col,
		trigger:  trigger,
	})
}

func (s *appState) startLSPForFile(filePath, content string) {
	if !s.lspEnabled || s.lspManager == nil || filePath == "" {
		return
	}

	go func() {
		server, err := s.lspManager.GetServerForFile(filePath)
		if err != nil {
			s.setLSPHint(filePath, err)
			s.maybeAutoInstallLSP(filePath)
			return
		}

		s.lspHint = ""
		s.lspHintFile = filePath

		if server.IsDocumentOpen(filePath) {
			_ = server.SyncDocument(filePath, content)
		} else {
			_ = server.OpenDocument(filePath, content)
		}

		s.flushPendingLSPChange(filePath)
	}()
}

func (s *appState) maybeAutoInstallLSP(filePath string) {
	if !s.lspAutoEnabled || s.lspManager == nil || filePath == "" {
		return
	}

	cfg := lsp.GetConfigForFile(filePath)
	if cfg == nil {
		return
	}

	if lsp.IsServerAvailable(cfg) {
		return
	}

	if s.lspAutoAttempted == nil {
		s.lspAutoAttempted = make(map[string]bool)
	}
	if s.lspAutoAttempted[cfg.Name] {
		return
	}
	s.lspAutoAttempted[cfg.Name] = true

	if len(cfg.InstallCommand) == 0 {
		if s.lspHint == "" || s.lspHintFile != filePath {
			s.lspHint = fmt.Sprintf("Auto-install unavailable for %s", cfg.Name)
			s.lspHintFile = filePath
		}
		return
	}

	s.installLSPServer(cfg, filePath, true)
}

// setupLSPForBuffer sets up LSP document synchronization for a buffer.
func (s *appState) setupLSPForBuffer(buf *editor.Buffer) {
	if !s.lspEnabled || s.lspManager == nil || buf == nil || buf.FilePath() == "" {
		return
	}

	// Check if there's a server for this file type
	if !s.lspManager.HasServerForFile(buf.FilePath()) {
		return
	}

	// Set up change callback
	buf.SetLSPChangeCallback(func(content string) {
		filePath := buf.FilePath()
		if filePath == "" {
			return
		}
		s.scheduleLSPChange(filePath, content)
	})

	filePath := buf.FilePath()
	content := buf.GetContent()
	s.startLSPForFile(filePath, content)
}

// cleanupLSPForBuffer cleans up LSP for a buffer.
func (s *appState) cleanupLSPForBuffer(buf *editor.Buffer) {
	if buf == nil {
		return
	}

	buf.ClearLSPChangeCallback()
	s.cancelLSPChange(buf.FilePath())
	s.cancelLSPCompletion(buf.FilePath())

	if s.lspManager == nil || buf.FilePath() == "" || !s.lspManager.IsServerRunningForFile(buf.FilePath()) {
		return
	}

	server, err := s.lspManager.GetServerForFile(buf.FilePath())
	if err != nil {
		return
	}

	_ = server.CloseDocument(buf.FilePath())
}

// getLSPDiagnosticsForLine returns diagnostics for a specific line.
func (s *appState) getLSPDiagnosticsForLine(filePath string, line int) []lsp.Diagnostic {
	diags, ok := s.lspDiagnostics[filePath]
	if !ok {
		return nil
	}

	var result []lsp.Diagnostic
	for _, d := range diags {
		if d.Range.Start.Line <= line && d.Range.End.Line >= line {
			result = append(result, d)
		}
	}
	return result
}

// getDiagnosticCountString returns a string with diagnostic counts for status bar.
func (s *appState) getDiagnosticCountString(filePath string) string {
	diags, ok := s.lspDiagnostics[filePath]
	if !ok || len(diags) == 0 {
		return ""
	}

	var errors, warnings int
	for _, d := range diags {
		switch d.Severity {
		case lsp.DiagnosticSeverityError:
			errors++
		case lsp.DiagnosticSeverityWarning:
			warnings++
		}
	}

	if errors > 0 && warnings > 0 {
		return fmt.Sprintf("E:%d W:%d", errors, warnings)
	} else if errors > 0 {
		return fmt.Sprintf("E:%d", errors)
	} else if warnings > 0 {
		return fmt.Sprintf("W:%d", warnings)
	}
	return ""
}

// processLSPCommand handles LSP-related command line commands.
func (s *appState) processLSPCommand(cmd string, args string) bool {
	switch cmd {
	case "gd", "godef", "definition":
		s.handleLSPGotoDefinition()
		return true
	case "hover":
		s.handleLSPHover()
		return true
	case "refs", "references":
		s.handleLSPReferences()
		return true
	case "rename":
		if args == "" {
			s.status = "Usage: :rename <newname>"
		} else {
			s.handleLSPRename(strings.TrimSpace(args))
		}
		return true
	case "format":
		s.handleLSPFormat()
		return true
	case "codeaction", "ca":
		s.handleLSPCodeAction()
		return true
	case "lspstatus":
		s.handleLSPStatus()
		return true
	case "lsprestart":
		s.handleLSPRestart()
		return true
	default:
		return false
	}
}
