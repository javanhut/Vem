package terminal

import (
	"image/color"
	"sync"
)

// Cell represents a single terminal cell
type Cell struct {
	Rune      rune        // Unicode character
	FG        color.NRGBA // Foreground color
	BG        color.NRGBA // Background color
	Bold      bool        // Bold attribute
	Dim       bool        // Dim attribute
	Italic    bool        // Italic attribute
	Underline bool        // Underline attribute
	Blink     bool        // Blink attribute
	Reverse   bool        // Reverse video
}

// Line represents a row of cells
type Line struct {
	Cells []Cell
	Dirty bool // Whether line needs redraw
}

// ScreenBuffer represents the terminal screen
type ScreenBuffer struct {
	lines       []Line
	width       int
	height      int
	cursorX     int
	cursorY     int
	cursorStyle CursorStyle
	mu          sync.RWMutex // Protects buffer
}

// CursorStyle represents cursor appearance
type CursorStyle int

const (
	CursorBlock CursorStyle = iota
	CursorUnderline
	CursorBar
)

// NewScreenBuffer creates a new screen buffer
func NewScreenBuffer(width, height int) *ScreenBuffer {
	sb := &ScreenBuffer{
		width:  width,
		height: height,
		lines:  make([]Line, height),
	}

	// Initialize all cells
	for i := range sb.lines {
		sb.lines[i].Cells = make([]Cell, width)
		sb.lines[i].Dirty = true
		// Initialize cells with default colors
		for j := range sb.lines[i].Cells {
			sb.lines[i].Cells[j] = Cell{
				Rune: ' ',
				FG:   DefaultFG,
				BG:   DefaultBG,
			}
		}
	}

	return sb
}

// Dimensions returns the buffer dimensions
func (sb *ScreenBuffer) Dimensions() (width, height int) {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return sb.width, sb.height
}

// GetLine returns a copy of line at index (thread-safe)
func (sb *ScreenBuffer) GetLine(index int) Line {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	if index < 0 || index >= len(sb.lines) {
		return Line{}
	}

	// Return copy
	line := sb.lines[index]
	cells := make([]Cell, len(line.Cells))
	copy(cells, line.Cells)
	return Line{Cells: cells, Dirty: line.Dirty}
}

// GetCursor returns cursor position
func (sb *ScreenBuffer) GetCursor() (x, y int, style CursorStyle) {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return sb.cursorX, sb.cursorY, sb.cursorStyle
}

// SetCursor sets cursor position
func (sb *ScreenBuffer) SetCursor(x, y int) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if x < 0 {
		x = 0
	}
	if x >= sb.width {
		x = sb.width - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= sb.height {
		y = sb.height - 1
	}

	sb.cursorX = x
	sb.cursorY = y
}

// SetCell sets a cell value
func (sb *ScreenBuffer) SetCell(x, y int, cell Cell) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if y < 0 || y >= sb.height || x < 0 || x >= sb.width {
		return
	}

	sb.lines[y].Cells[x] = cell
	sb.lines[y].Dirty = true
}

// ClearLine clears a line
func (sb *ScreenBuffer) ClearLine(y int) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if y < 0 || y >= sb.height {
		return
	}

	for x := range sb.lines[y].Cells {
		sb.lines[y].Cells[x] = Cell{
			Rune: ' ',
			FG:   DefaultFG,
			BG:   DefaultBG,
		}
	}
	sb.lines[y].Dirty = true
}

// Clear clears the entire buffer
func (sb *ScreenBuffer) Clear() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	for y := range sb.lines {
		for x := range sb.lines[y].Cells {
			sb.lines[y].Cells[x] = Cell{
				Rune: ' ',
				FG:   DefaultFG,
				BG:   DefaultBG,
			}
		}
		sb.lines[y].Dirty = true
	}
}

// MarkClean marks all lines as clean (after render)
func (sb *ScreenBuffer) MarkClean() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	for i := range sb.lines {
		sb.lines[i].Dirty = false
	}
}

// Resize resizes the buffer
func (sb *ScreenBuffer) Resize(width, height int) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if width == sb.width && height == sb.height {
		return
	}

	newLines := make([]Line, height)
	for i := range newLines {
		newLines[i].Cells = make([]Cell, width)
		newLines[i].Dirty = true

		// Copy old content if available
		if i < len(sb.lines) {
			copyWidth := width
			if copyWidth > len(sb.lines[i].Cells) {
				copyWidth = len(sb.lines[i].Cells)
			}
			copy(newLines[i].Cells, sb.lines[i].Cells[:copyWidth])
		}

		// Fill remaining cells with defaults
		for j := range newLines[i].Cells {
			if i >= len(sb.lines) || j >= len(sb.lines[i].Cells) {
				newLines[i].Cells[j] = Cell{
					Rune: ' ',
					FG:   DefaultFG,
					BG:   DefaultBG,
				}
			}
		}
	}

	sb.lines = newLines
	sb.width = width
	sb.height = height

	// Clamp cursor
	if sb.cursorX >= width {
		sb.cursorX = width - 1
	}
	if sb.cursorY >= height {
		sb.cursorY = height - 1
	}
}

// Implement io.Writer for vt10x emulator
func (sb *ScreenBuffer) Write(p []byte) (n int, err error) {
	// This will be called by vt10x to update the buffer
	// vt10x handles parsing and calls our methods
	return len(p), nil
}
