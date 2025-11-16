package editor

import (
	"fmt"
	"os"
	"path/filepath"
)

// BufferManager manages multiple buffers and tracks the active buffer.
type BufferManager struct {
	buffers     []*Buffer
	activeIndex int
	pathToIndex map[string]int
}

// NewBufferManager creates a new buffer manager with a default empty buffer.
func NewBufferManager() *BufferManager {
	defaultBuffer := NewBuffer("")
	return &BufferManager{
		buffers:     []*Buffer{defaultBuffer},
		activeIndex: 0,
		pathToIndex: make(map[string]int),
	}
}

// NewBufferManagerWithBuffer creates a buffer manager with an initial buffer.
func NewBufferManagerWithBuffer(buf *Buffer) *BufferManager {
	bm := &BufferManager{
		buffers:     []*Buffer{buf},
		activeIndex: 0,
		pathToIndex: make(map[string]int),
	}

	if buf.FilePath() != "" {
		bm.pathToIndex[buf.FilePath()] = 0
	}

	return bm
}

// ActiveBuffer returns the currently active buffer.
func (bm *BufferManager) ActiveBuffer() *Buffer {
	if bm.activeIndex >= 0 && bm.activeIndex < len(bm.buffers) {
		return bm.buffers[bm.activeIndex]
	}
	return nil
}

// BufferCount returns the total number of buffers.
func (bm *BufferManager) BufferCount() int {
	return len(bm.buffers)
}

// ActiveIndex returns the index of the active buffer.
func (bm *BufferManager) ActiveIndex() int {
	return bm.activeIndex
}

// GetBuffer returns the buffer at the specified index.
func (bm *BufferManager) GetBuffer(index int) *Buffer {
	if index >= 0 && index < len(bm.buffers) {
		return bm.buffers[index]
	}
	return nil
}

// GetBufferByPath returns the buffer for the given file path, or nil if not found.
func (bm *BufferManager) GetBufferByPath(path string) *Buffer {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil
	}

	if index, exists := bm.pathToIndex[absPath]; exists {
		return bm.buffers[index]
	}
	return nil
}

// OpenFile opens a file into a new or existing buffer and makes it active.
// If the file is already open, it switches to that buffer instead.
func (bm *BufferManager) OpenFile(path string) (*Buffer, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Check if already open
	if existing := bm.GetBufferByPath(absPath); existing != nil {
		// Switch to existing buffer
		for i, buf := range bm.buffers {
			if buf == existing {
				bm.activeIndex = i
				return buf, nil
			}
		}
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		// Create new empty buffer for new file
		buf := NewBuffer("")
		buf.SetFilePath(absPath)
		buf.SetModified(false)
		return bm.addBuffer(buf), nil
	}

	// Load existing file
	buf, err := NewBufferFromFile(absPath)
	if err != nil {
		return nil, err
	}

	return bm.addBuffer(buf), nil
}

// addBuffer adds a buffer to the manager and makes it active.
func (bm *BufferManager) addBuffer(buf *Buffer) *Buffer {
	bm.buffers = append(bm.buffers, buf)
	bm.activeIndex = len(bm.buffers) - 1

	if buf.FilePath() != "" {
		bm.pathToIndex[buf.FilePath()] = bm.activeIndex
	}

	return buf
}

// CreateEmptyBuffer creates a new empty buffer and returns its index.
func (bm *BufferManager) CreateEmptyBuffer() int {
	buf := NewBuffer("")
	bm.buffers = append(bm.buffers, buf)
	return len(bm.buffers) - 1
}

// CreateTerminalBuffer creates a new buffer for a terminal and returns its index.
func (bm *BufferManager) CreateTerminalBuffer() int {
	buf := &Buffer{
		lines:      []string{""},
		cursor:     Cursor{},
		bufferType: BufferTypeTerminal,
		undoStack:  make([]UndoEntry, 0),
		maxUndos:   100,
	}
	bm.buffers = append(bm.buffers, buf)
	return len(bm.buffers) - 1
}

// CreateBufferWithContent creates a new buffer with the given content and returns its index.
func (bm *BufferManager) CreateBufferWithContent(content string) int {
	buf := NewBuffer(content)
	bm.buffers = append(bm.buffers, buf)
	return len(bm.buffers) - 1
}

// SaveActiveBuffer saves the currently active buffer.
func (bm *BufferManager) SaveActiveBuffer() error {
	buf := bm.ActiveBuffer()
	if buf == nil {
		return fmt.Errorf("no active buffer")
	}

	if buf.FilePath() == "" {
		return fmt.Errorf("no file name")
	}

	return buf.Save()
}

// SaveAs saves the active buffer to a new file path.
func (bm *BufferManager) SaveAs(path string) error {
	buf := bm.ActiveBuffer()
	if buf == nil {
		return fmt.Errorf("no active buffer")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Remove old path mapping
	if buf.FilePath() != "" {
		delete(bm.pathToIndex, buf.FilePath())
	}

	// Save to new path
	if err := buf.SaveToFile(absPath); err != nil {
		return err
	}

	// Update path mapping
	bm.pathToIndex[absPath] = bm.activeIndex

	return nil
}

// CloseBuffer closes the buffer at the specified index.
// If it's modified, returns an error unless force is true.
func (bm *BufferManager) CloseBuffer(index int, force bool) error {
	if index < 0 || index >= len(bm.buffers) {
		return fmt.Errorf("invalid buffer index")
	}

	buf := bm.buffers[index]

	// Terminal buffers don't have unsaved changes
	if !buf.IsTerminal() {
		// Check for unsaved changes in text buffers
		if !force && buf.Modified() {
			return fmt.Errorf("buffer has unsaved changes (use :q! to force)")
		}
	}

	// Remove from path mapping
	if buf.FilePath() != "" {
		delete(bm.pathToIndex, buf.FilePath())
	}

	// Remove buffer
	bm.buffers = append(bm.buffers[:index], bm.buffers[index+1:]...)

	// Update path mappings for shifted buffers
	for i := index; i < len(bm.buffers); i++ {
		if bm.buffers[i].FilePath() != "" {
			bm.pathToIndex[bm.buffers[i].FilePath()] = i
		}
	}

	// Adjust active index
	if len(bm.buffers) == 0 {
		// Create default empty buffer
		defaultBuf := NewBuffer("")
		bm.buffers = []*Buffer{defaultBuf}
		bm.activeIndex = 0
	} else if bm.activeIndex >= len(bm.buffers) {
		bm.activeIndex = len(bm.buffers) - 1
	} else if bm.activeIndex > index {
		bm.activeIndex--
	}

	return nil
}

// CloseActiveBuffer closes the currently active buffer.
func (bm *BufferManager) CloseActiveBuffer(force bool) error {
	return bm.CloseBuffer(bm.activeIndex, force)
}

// NextBuffer switches to the next buffer (wraps around).
func (bm *BufferManager) NextBuffer() bool {
	if len(bm.buffers) <= 1 {
		return false
	}

	bm.activeIndex = (bm.activeIndex + 1) % len(bm.buffers)
	return true
}

// PrevBuffer switches to the previous buffer (wraps around).
func (bm *BufferManager) PrevBuffer() bool {
	if len(bm.buffers) <= 1 {
		return false
	}

	bm.activeIndex--
	if bm.activeIndex < 0 {
		bm.activeIndex = len(bm.buffers) - 1
	}
	return true
}

// SwitchToBuffer switches to the buffer at the specified index.
func (bm *BufferManager) SwitchToBuffer(index int) bool {
	if index >= 0 && index < len(bm.buffers) {
		bm.activeIndex = index
		return true
	}
	return false
}

// ListBuffers returns a slice of buffer info for display.
func (bm *BufferManager) ListBuffers() []string {
	result := make([]string, len(bm.buffers))
	for i, buf := range bm.buffers {
		prefix := " "
		if i == bm.activeIndex {
			prefix = "*"
		}

		modFlag := " "
		if buf.Modified() {
			modFlag = "+"
		}

		name := buf.FilePath()
		if name == "" {
			if buf.IsTerminal() {
				name = "[Terminal]"
			} else {
				name = "[No Name]"
			}
		}

		result[i] = fmt.Sprintf("%s %d %s %s", prefix, i+1, modFlag, name)
	}
	return result
}
