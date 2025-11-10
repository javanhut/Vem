package panes

// Pane represents a single editor pane with its own view state.
// Each pane displays exactly one buffer and maintains independent scroll position.
type Pane struct {
	ID          string // Unique identifier for this pane
	BufferIndex int    // Index into BufferManager.buffers
	Active      bool   // Is this pane currently focused?
	ViewportTop int    // First visible line (0-based) for independent scrolling
}

// NewPane creates a new pane with the given buffer index.
func NewPane(id string, bufferIndex int) *Pane {
	return &Pane{
		ID:          id,
		BufferIndex: bufferIndex,
		Active:      false,
		ViewportTop: 0,
	}
}

// SetActive marks this pane as active (focused).
func (p *Pane) SetActive(active bool) {
	p.Active = active
}

// SetBufferIndex changes which buffer this pane displays.
func (p *Pane) SetBufferIndex(index int) {
	p.BufferIndex = index
	// Reset viewport when switching buffers
	p.ViewportTop = 0
}

// SetViewportTop sets the first visible line for this pane.
func (p *Pane) SetViewportTop(line int) {
	if line < 0 {
		line = 0
	}
	p.ViewportTop = line
}
