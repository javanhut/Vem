package appcore

// handleBufferCompletionTrigger triggers buffer word completion.
// This looks for words in the current buffer that match what the user is typing.
func (s *appState) handleBufferCompletionTrigger() {
	buf := s.activeBuffer()
	if buf == nil {
		return
	}

	// Get the word prefix at cursor
	prefix := buf.GetCurrentWordPrefix()
	if prefix == "" || len(prefix) < 1 {
		s.bufferCompletionActive = false
		return
	}

	// Get matching words from buffer
	matches := buf.GetWordsMatching(prefix)
	if len(matches) == 0 {
		s.status = "No completions found"
		s.bufferCompletionActive = false
		return
	}

	// Activate buffer completion
	s.bufferCompletionActive = true
	s.bufferCompletionItems = matches
	s.bufferCompletionIndex = 0
	s.bufferCompletionPrefix = prefix

	// Cancel LSP completion if active
	s.completionActive = false
	s.completionItems = nil
}

// handleBufferCompletionAccept accepts the selected buffer completion.
func (s *appState) handleBufferCompletionAccept() {
	if !s.bufferCompletionActive || len(s.bufferCompletionItems) == 0 {
		return
	}

	buf := s.activeBuffer()
	if buf == nil {
		return
	}

	// Get selected completion
	completion := s.bufferCompletionItems[s.bufferCompletionIndex]

	// Delete the prefix and insert the completion
	cursor := buf.Cursor()
	line := buf.Line(cursor.Line)
	runes := []rune(line)

	// Find start of word
	wordStart := cursor.Col
	for wordStart > 0 && isBufferWordChar(runes[wordStart-1]) {
		wordStart--
	}

	// Delete from word start to cursor
	if wordStart < cursor.Col {
		buf.DeleteCharRange(cursor.Line, wordStart, cursor.Line, cursor.Col)
	}

	// Insert completion text
	buf.InsertText(completion)

	// Clear completion state
	s.bufferCompletionActive = false
	s.bufferCompletionItems = nil
	s.bufferCompletionPrefix = ""

	// Update syntax highlighting
	s.invalidateSyntaxCache()
}

// handleBufferCompletionCancel cancels buffer completion.
func (s *appState) handleBufferCompletionCancel() {
	s.bufferCompletionActive = false
	s.bufferCompletionItems = nil
	s.bufferCompletionPrefix = ""
}

// handleBufferCompletionNext moves to the next buffer completion item.
func (s *appState) handleBufferCompletionNext() {
	if !s.bufferCompletionActive || len(s.bufferCompletionItems) == 0 {
		return
	}
	s.bufferCompletionIndex = (s.bufferCompletionIndex + 1) % len(s.bufferCompletionItems)
}

// handleBufferCompletionPrev moves to the previous buffer completion item.
func (s *appState) handleBufferCompletionPrev() {
	if !s.bufferCompletionActive || len(s.bufferCompletionItems) == 0 {
		return
	}
	s.bufferCompletionIndex--
	if s.bufferCompletionIndex < 0 {
		s.bufferCompletionIndex = len(s.bufferCompletionItems) - 1
	}
}

// isBufferWordChar checks if a rune is a word character for buffer completion.
func isBufferWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}
