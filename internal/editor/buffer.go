package editor

import (
	"os"
	"strings"
	"unicode/utf8"
)

// Buffer represents an in-memory text buffer with a Vim-style cursor.
type Buffer struct {
	lines    []string
	cursor   Cursor
	filePath string
	modified bool
}

// Cursor stores the current line/column position (1 rune == 1 column).
type Cursor struct {
	Line int
	Col  int
}

// NewBuffer builds a Buffer from a block of text.
func NewBuffer(text string) *Buffer {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	return &Buffer{
		lines:  lines,
		cursor: Cursor{},
	}
}

// LineCount returns the number of lines in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// Line returns the line at the supplied index or an empty string if out of bounds.
func (b *Buffer) Line(i int) string {
	if i < 0 || i >= len(b.lines) {
		return ""
	}
	return b.lines[i]
}

// LinesRange returns a copy of lines between start and end (inclusive), clamped to buffer bounds.
func (b *Buffer) LinesRange(start, end int) []string {
	if len(b.lines) == 0 {
		return []string{}
	}
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if end >= len(b.lines) {
		end = len(b.lines) - 1
	}
	if start >= len(b.lines) {
		return []string{}
	}
	lines := make([]string, end-start+1)
	copy(lines, b.lines[start:end+1])
	return lines
}

// LinePrefix returns the first prefixCols runes of the line at index.
func (b *Buffer) LinePrefix(lineIdx, prefixCols int) string {
	line := b.Line(lineIdx)
	if prefixCols <= 0 || line == "" {
		return ""
	}
	if prefixCols >= runeCount(line) {
		return line
	}
	byteIdx := byteIndexForRune(line, prefixCols)
	return line[:byteIdx]
}

// Cursor returns the current cursor position.
func (b *Buffer) Cursor() Cursor {
	return b.cursor
}

// MoveToLine moves the cursor to the provided zero-based line index.
func (b *Buffer) MoveToLine(line int) {
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}
	if line < 0 {
		line = 0
	} else if line >= len(b.lines) {
		line = len(b.lines) - 1
	}
	b.cursor.Line = line
	b.clampColumn()
}

// DeleteLines removes the inclusive line range and repositions the cursor.
func (b *Buffer) DeleteLines(start, end int) {
	if len(b.lines) == 0 {
		return
	}
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if end >= len(b.lines) {
		end = len(b.lines) - 1
	}
	if start >= len(b.lines) {
		return
	}
	b.lines = append(b.lines[:start], b.lines[end+1:]...)
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}
	if start >= len(b.lines) {
		start = len(b.lines) - 1
	}
	b.cursor.Line = start
	b.clampColumn()
	b.markModified()
}

// InsertLines inserts the provided lines at the given index, adjusting the cursor to the end of the block.
func (b *Buffer) InsertLines(at int, lines []string) {
	if len(lines) == 0 {
		return
	}
	if at < 0 {
		at = 0
	}
	if at > len(b.lines) {
		at = len(b.lines)
	}
	linesCopy := append([]string(nil), lines...)
	newLines := make([]string, 0, len(b.lines)+len(linesCopy))
	newLines = append(newLines, b.lines[:at]...)
	newLines = append(newLines, linesCopy...)
	newLines = append(newLines, b.lines[at:]...)
	b.lines = newLines
	b.cursor.Line = at + len(linesCopy) - 1
	b.clampColumn()
	b.markModified()
}

// InsertText inserts the provided text at the cursor position and moves the cursor
// to the end of the inserted text.
func (b *Buffer) InsertText(text string) {
	if text == "" {
		return
	}
	left, right := splitAtRune(b.lines[b.cursor.Line], b.cursor.Col)
	segments := strings.Split(text, "\n")
	lastIdx := len(segments) - 1
	lastSegmentLen := runeCount(segments[lastIdx])

	segments[0] = left + segments[0]
	segments[lastIdx] = segments[lastIdx] + right

	prefix := append([]string{}, b.lines[:b.cursor.Line]...)
	suffix := append([]string{}, b.lines[b.cursor.Line+1:]...)

	b.lines = append(append(prefix, segments...), suffix...)

	b.cursor.Line += lastIdx
	if lastIdx == 0 {
		b.cursor.Col = runeCount(left) + lastSegmentLen
	} else {
		b.cursor.Col = lastSegmentLen
	}
	b.markModified()
}

// DeleteBackward deletes the rune before the cursor (backspace semantics).
// When invoked at the start of a line, it merges with the previous line.
func (b *Buffer) DeleteBackward() bool {
	if b.cursor.Col == 0 {
		if b.cursor.Line == 0 {
			return false
		}
		prev := b.cursor.Line - 1
		prevLen := runeCount(b.lines[prev])
		b.lines[prev] = b.lines[prev] + b.lines[b.cursor.Line]
		b.lines = removeLine(b.lines, b.cursor.Line)
		b.cursor.Line = prev
		b.cursor.Col = prevLen
		b.markModified()
		return true
	}

	line := []rune(b.lines[b.cursor.Line])
	if b.cursor.Col > len(line) {
		b.cursor.Col = len(line)
	}
	line = append(line[:b.cursor.Col-1], line[b.cursor.Col:]...)
	b.lines[b.cursor.Line] = string(line)
	b.cursor.Col--
	b.markModified()
	return true
}

// DeleteForward deletes the rune at the cursor (delete semantics).
// When at the end of a line, it merges with the following line.
func (b *Buffer) DeleteForward() bool {
	lineRunes := []rune(b.lines[b.cursor.Line])
	if b.cursor.Col < len(lineRunes) {
		lineRunes = append(lineRunes[:b.cursor.Col], lineRunes[b.cursor.Col+1:]...)
		b.lines[b.cursor.Line] = string(lineRunes)
		b.markModified()
		return true
	}
	if b.cursor.Line >= len(b.lines)-1 {
		return false
	}
	b.lines[b.cursor.Line] = b.lines[b.cursor.Line] + b.lines[b.cursor.Line+1]
	b.lines = removeLine(b.lines, b.cursor.Line+1)
	b.markModified()
	return true
}

// MoveLeft moves the cursor left, spilling to the previous line when needed.
func (b *Buffer) MoveLeft() bool {
	if b.cursor.Col > 0 {
		b.cursor.Col--
		return true
	}
	if b.cursor.Line == 0 {
		return false
	}
	b.cursor.Line--
	b.cursor.Col = b.lineLength(b.cursor.Line)
	return true
}

// MoveRight moves the cursor right, spilling to the next line when needed.
func (b *Buffer) MoveRight() bool {
	lineLen := b.lineLength(b.cursor.Line)
	if b.cursor.Col < lineLen {
		b.cursor.Col++
		return true
	}
	if b.cursor.Line >= len(b.lines)-1 {
		return false
	}
	b.cursor.Line++
	b.cursor.Col = 0
	return true
}

// MoveUp moves the cursor to the previous line, clamped by line length.
func (b *Buffer) MoveUp() bool {
	if b.cursor.Line == 0 {
		return false
	}
	b.cursor.Line--
	b.clampColumn()
	return true
}

// MoveDown moves the cursor to the next line, clamped by line length.
func (b *Buffer) MoveDown() bool {
	if b.cursor.Line >= len(b.lines)-1 {
		return false
	}
	b.cursor.Line++
	b.clampColumn()
	return true
}

// JumpLineStart places the cursor at the start of the current line.
func (b *Buffer) JumpLineStart() bool {
	if b.cursor.Col == 0 {
		return false
	}
	b.cursor.Col = 0
	return true
}

// JumpLineEnd places the cursor at the end of the current line.
func (b *Buffer) JumpLineEnd() bool {
	lineLen := b.lineLength(b.cursor.Line)
	if b.cursor.Col == lineLen {
		return false
	}
	b.cursor.Col = lineLen
	return true
}

func (b *Buffer) clampColumn() {
	lineLen := b.lineLength(b.cursor.Line)
	if b.cursor.Col > lineLen {
		b.cursor.Col = lineLen
	}
}

func (b *Buffer) lineLength(line int) int {
	if line < 0 || line >= len(b.lines) {
		return 0
	}
	return utf8.RuneCountInString(b.lines[line])
}

func splitAtRune(text string, index int) (string, string) {
	if index <= 0 {
		return "", text
	}
	length := runeCount(text)
	if index >= length {
		return text, ""
	}
	byteIdx := byteIndexForRune(text, index)
	return text[:byteIdx], text[byteIdx:]
}

func runeCount(s string) int {
	return utf8.RuneCountInString(s)
}

func byteIndexForRune(s string, idx int) int {
	if idx <= 0 {
		return 0
	}
	if idx >= runeCount(s) {
		return len(s)
	}
	count := 0
	byteIdx := 0
	for byteIdx < len(s) && count < idx {
		_, size := utf8.DecodeRuneInString(s[byteIdx:])
		byteIdx += size
		count++
	}
	return byteIdx
}

func removeLine(lines []string, index int) []string {
	if index < 0 || index >= len(lines) {
		return lines
	}
	return append(lines[:index], lines[index+1:]...)
}

// FilePath returns the file path associated with this buffer.
func (b *Buffer) FilePath() string {
	return b.filePath
}

// SetFilePath sets the file path for this buffer.
func (b *Buffer) SetFilePath(path string) {
	b.filePath = path
}

// Modified returns true if the buffer has unsaved changes.
func (b *Buffer) Modified() bool {
	return b.modified
}

// SetModified sets the modified flag.
func (b *Buffer) SetModified(modified bool) {
	b.modified = modified
}

// MarkModified marks the buffer as modified (used internally after edits).
func (b *Buffer) markModified() {
	b.modified = true
}

// LoadFromFile loads the buffer content from a file.
func (b *Buffer) LoadFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	text := string(content)
	lines := strings.Split(text, "\n")
	
	// Remove trailing empty line if file ends with newline
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	
	if len(lines) == 0 {
		lines = []string{""}
	}

	b.lines = lines
	b.cursor = Cursor{Line: 0, Col: 0}
	b.filePath = path
	b.modified = false

	return nil
}

// SaveToFile saves the buffer content to a file.
func (b *Buffer) SaveToFile(path string) error {
	content := b.GetContent()
	
	// Ensure file ends with newline
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	b.filePath = path
	b.modified = false
	return nil
}

// Save saves the buffer to its associated file path.
func (b *Buffer) Save() error {
	if b.filePath == "" {
		return os.ErrInvalid
	}
	return b.SaveToFile(b.filePath)
}

// GetContent returns the entire buffer content as a string.
func (b *Buffer) GetContent() string {
	return strings.Join(b.lines, "\n")
}

// NewBufferFromFile creates a new buffer and loads content from a file.
func NewBufferFromFile(path string) (*Buffer, error) {
	buf := &Buffer{
		lines:  []string{""},
		cursor: Cursor{},
	}
	
	if err := buf.LoadFromFile(path); err != nil {
		return nil, err
	}
	
	return buf, nil
}
