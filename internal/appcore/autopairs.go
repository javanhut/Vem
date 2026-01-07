package appcore

import (
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2"
)

// Auto-pair character mappings
var autoPairs = map[rune]rune{
	'(':  ')',
	'[':  ']',
	'{':  '}',
	'"':  '"',
	'\'': '\'',
}

// Closing characters for skip-over detection
var closingChars = map[rune]bool{
	')':  true,
	']':  true,
	'}':  true,
	'"':  true,
	'\'': true,
}

// isOpeningChar returns the closing char if r is an opening bracket/quote
func isOpeningChar(r rune) (rune, bool) {
	close, ok := autoPairs[r]
	return close, ok
}

// isClosingChar returns true if r is a closing bracket/quote
func isClosingChar(r rune) bool {
	return closingChars[r]
}

// charBeforeCursor returns the rune immediately before cursor (or 0 if none)
func (s *appState) charBeforeCursor() rune {
	buf := s.activeBuffer()
	if buf == nil {
		return 0
	}

	cursor := buf.Cursor()
	line := buf.Line(cursor.Line)
	runes := []rune(line)

	if cursor.Col > 0 && cursor.Col <= len(runes) {
		return runes[cursor.Col-1]
	}
	return 0
}

// charAtCursor returns the rune at cursor position (or 0 if none)
func (s *appState) charAtCursor() rune {
	buf := s.activeBuffer()
	if buf == nil {
		return 0
	}

	cursor := buf.Cursor()
	line := buf.Line(cursor.Line)
	runes := []rune(line)

	if cursor.Col < len(runes) {
		return runes[cursor.Col]
	}
	return 0
}

// shouldSkipClosingChar checks if cursor is before the same closing char
func (s *appState) shouldSkipClosingChar(r rune) bool {
	return s.charAtCursor() == r
}

// isAlphaNumeric checks if rune is alphanumeric or underscore
func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

// insertAutoPair inserts opening char, closing char, and positions cursor between
func (s *appState) insertAutoPair(open, close rune) {
	buf := s.activeBuffer()
	if buf == nil {
		return
	}

	// Insert both characters
	buf.InsertText(string(open) + string(close))

	// Move cursor back one position (between the pair)
	buf.MoveLeft()
}

// handleAutoPairInsertion processes a character for auto-pairing
// Returns true if the character was handled (inserted or skipped)
func (s *appState) handleAutoPairInsertion(r rune) bool {
	if !s.autoPairsEnabled {
		return false
	}

	buf := s.activeBuffer()
	if buf == nil {
		return false
	}

	// Skip for terminal buffers
	if buf.BufferType() != 0 { // BufferTypeText is 0
		return false
	}

	// Case 1: Typing a closing char that matches char at cursor - skip over it
	if isClosingChar(r) && s.shouldSkipClosingChar(r) {
		buf.MoveRight()
		return true
	}

	// Case 2: Typing an opening char - insert pair
	if close, ok := isOpeningChar(r); ok {
		// For quotes, apply special rules
		if r == '"' || r == '\'' {
			// Skip if previous char is backslash (escaping)
			if s.charBeforeCursor() == '\\' {
				return false
			}
			// Skip if previous char is alphanumeric (contractions like it's)
			prev := s.charBeforeCursor()
			if prev != 0 && isAlphaNumeric(prev) {
				return false
			}
		}

		s.insertAutoPair(r, close)
		return true
	}

	// Case 3: Auto-dedent when typing } at start of line
	if s.autoIndentEnabled && (r == '}' || r == ']' || r == ')') {
		if s.handleAutoDedent(r) {
			return true
		}
	}

	return false
}

// getLeadingWhitespace returns the leading whitespace of a line
func getLeadingWhitespace(line string) string {
	var ws strings.Builder
	for _, r := range line {
		if r == ' ' || r == '\t' {
			ws.WriteRune(r)
		} else {
			break
		}
	}
	return ws.String()
}

// shouldAddExtraIndent checks if line ends with a block-opening character
func shouldAddExtraIndent(textBeforeCursor string, filePath string) bool {
	trimmed := strings.TrimRight(textBeforeCursor, " \t")
	if trimmed == "" {
		return false
	}

	lastChar := rune(trimmed[len(trimmed)-1])

	// C-like languages: indent after {
	if lastChar == '{' {
		return true
	}

	// Indent after ( and [ for multi-line arguments
	if lastChar == '(' || lastChar == '[' {
		return true
	}

	// Python-like: indent after :
	if lastChar == ':' && isPythonLike(filePath) {
		return true
	}

	return false
}

// isPythonLike checks if file is Python or similar language
func isPythonLike(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".py" || ext == ".pyw" || ext == ".yaml" || ext == ".yml"
}

// isCursorBetweenBrackets checks if cursor is between matching brackets {}[]()
func (s *appState) isCursorBetweenBrackets() (bool, rune, rune) {
	before := s.charBeforeCursor()
	after := s.charAtCursor()

	pairs := map[rune]rune{
		'{': '}',
		'[': ']',
		'(': ')',
	}

	if close, ok := pairs[before]; ok && after == close {
		return true, before, close
	}

	return false, 0, 0
}

// isCursorInComment uses Chroma tokenization to detect if cursor is in a comment
func (s *appState) isCursorInComment() bool {
	buf := s.activeBuffer()
	if buf == nil {
		return false
	}

	highlighter := s.getOrCreateHighlighter()
	if highlighter == nil {
		return false
	}

	cursor := buf.Cursor()
	line := buf.Line(cursor.Line)
	tokens := highlighter.HighlightLine(cursor.Line, line)

	// Find which token contains the cursor
	col := cursor.Col
	currentCol := 0

	for _, token := range tokens {
		tokenLen := utf8.RuneCountInString(token.Text)
		if col >= currentCol && col < currentCol+tokenLen {
			// Cursor is in this token - check if it's a comment
			return isCommentTokenType(token.Type)
		}
		currentCol += tokenLen
	}

	// If cursor is at end of line, check last token
	if len(tokens) > 0 && col >= currentCol {
		lastToken := tokens[len(tokens)-1]
		return isCommentTokenType(lastToken.Type)
	}

	return false
}

// isCommentTokenType checks if a Chroma token type is a comment
func isCommentTokenType(t chroma.TokenType) bool {
	return t == chroma.Comment ||
		t == chroma.CommentSingle ||
		t == chroma.CommentMultiline ||
		t == chroma.CommentSpecial ||
		t == chroma.CommentPreproc ||
		t == chroma.CommentPreprocFile
}

// insertNewlineWithIndent handles smart newline insertion with indentation
func (s *appState) insertNewlineWithIndent() {
	buf := s.activeBuffer()
	if buf == nil {
		s.insertText("\n")
		return
	}

	cursor := buf.Cursor()
	currentLine := buf.Line(cursor.Line)
	baseIndent := getLeadingWhitespace(currentLine)

	// Get text before cursor on current line
	runes := []rune(currentLine)
	textBeforeCursor := ""
	if cursor.Col > 0 && cursor.Col <= len(runes) {
		textBeforeCursor = string(runes[:cursor.Col])
	}

	filePath := buf.FilePath()

	// Check for bracket expansion: cursor between {}[]()
	if between, _, _ := s.isCursorBetweenBrackets(); between {
		// Insert: newline + indent + extra + newline + indent
		extraIndent := s.indentString

		// Insert first newline with extra indent
		buf.InsertText("\n" + baseIndent + extraIndent)

		// Remember cursor position (we want cursor on middle line)
		middleLine := buf.Cursor().Line
		middleCol := buf.Cursor().Col

		// Insert closing line
		buf.InsertText("\n" + baseIndent)

		// Move cursor back to middle line
		buf.SetCursor(middleLine, middleCol)

		return
	}

	// Check if we should add extra indent
	needsExtraIndent := shouldAddExtraIndent(textBeforeCursor, filePath)

	if needsExtraIndent {
		buf.InsertText("\n" + baseIndent + s.indentString)
	} else {
		buf.InsertText("\n" + baseIndent)
	}
}

// handleAutoDedent checks if typing } at line start should dedent
// Returns true if handled
func (s *appState) handleAutoDedent(r rune) bool {
	buf := s.activeBuffer()
	if buf == nil {
		return false
	}

	cursor := buf.Cursor()
	line := buf.Line(cursor.Line)

	// Check if everything before cursor is whitespace
	runes := []rune(line)
	beforeCursor := ""
	if cursor.Col > 0 {
		beforeCursor = string(runes[:cursor.Col])
	}

	// Only dedent if line only has whitespace before cursor
	if strings.TrimSpace(beforeCursor) != "" {
		return false
	}

	// Need at least some indent to dedent
	currentIndent := getLeadingWhitespace(line)
	if currentIndent == "" {
		return false
	}

	// Remove one level of indent
	newIndent := dedentOnce(currentIndent, s.indentString)

	// Get text after cursor
	afterCursor := ""
	if cursor.Col < len(runes) {
		afterCursor = string(runes[cursor.Col:])
	}

	// Replace the line content
	newLine := newIndent + string(r) + afterCursor

	// Use buffer methods to replace line
	buf.ReplaceLine(cursor.Line, newLine)

	// Position cursor after the bracket
	newCol := utf8.RuneCountInString(newIndent) + 1
	buf.SetCursor(cursor.Line, newCol)

	return true
}

// dedentOnce removes one level of indentation
func dedentOnce(indent, indentUnit string) string {
	if strings.HasSuffix(indent, indentUnit) {
		return indent[:len(indent)-len(indentUnit)]
	}
	// Handle mixed tabs/spaces - try removing a tab first, then spaces
	if strings.HasSuffix(indent, "\t") {
		return indent[:len(indent)-1]
	}
	// Remove up to 4 trailing spaces
	spaces := 0
	runes := []rune(indent)
	for i := len(runes) - 1; i >= 0 && spaces < 4; i-- {
		if runes[i] == ' ' {
			spaces++
		} else {
			break
		}
	}
	if spaces > 0 {
		return string(runes[:len(runes)-spaces])
	}
	return indent
}
