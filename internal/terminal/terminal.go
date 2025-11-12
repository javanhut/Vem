package terminal

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"gioui.org/app"
	"github.com/hinshun/vt10x"
)

// Terminal represents a running terminal emulator session
type Terminal struct {
	// PTY and process
	pty *os.File  // PTY master file descriptor
	cmd *exec.Cmd // Shell process

	// VT100 emulator
	vt vt10x.Terminal // Terminal interface (VT100)

	// Screen buffer
	screen *ScreenBuffer // Current screen content

	// Terminal size
	width  int // Columns (e.g., 80)
	height int // Rows (e.g., 24)

	// Shell info
	shell      string   // Shell path (/bin/bash, etc.)
	args       []string // Shell arguments
	workingDir string   // Working directory
	env        []string // Environment variables

	// Lifecycle
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	mu      sync.RWMutex // Protects running flag

	// Communication channels
	inputChan  chan []byte   // Input from UI to PTY
	updateChan chan struct{} // Signals screen update

	// Window invalidation
	window *app.Window // For triggering redraws

	// Goroutine management
	wg sync.WaitGroup // Waits for goroutines to finish

	// Error tracking
	lastError error
	errorMu   sync.Mutex

	// Exit callback
	onExit func() // Called when terminal process exits
}

// Config holds terminal configuration
type Config struct {
	Width      int
	Height     int
	Shell      string
	Args       []string
	WorkingDir string
	Env        []string
	Window     *app.Window // For invalidation
	OnExit     func()      // Called when terminal process exits
}

// NewTerminal creates a new terminal with given config
func NewTerminal(cfg Config) (*Terminal, error) {
	// Validation
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: %dx%d", cfg.Width, cfg.Height)
	}
	if cfg.Shell == "" {
		cfg.Shell = DefaultShell()
	}
	if cfg.Args == nil {
		cfg.Args = DefaultArgs()
	}
	if cfg.WorkingDir == "" {
		cfg.WorkingDir, _ = os.Getwd()
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())

	// Create terminal
	t := &Terminal{
		width:      cfg.Width,
		height:     cfg.Height,
		shell:      cfg.Shell,
		args:       cfg.Args,
		workingDir: cfg.WorkingDir,
		env:        cfg.Env,
		ctx:        ctx,
		cancel:     cancel,
		inputChan:  make(chan []byte, 256), // Buffered for responsiveness
		updateChan: make(chan struct{}, 1), // Buffered, drop duplicates
		window:     cfg.Window,
		onExit:     cfg.OnExit,
	}

	// Create screen buffer
	t.screen = NewScreenBuffer(cfg.Width, cfg.Height)

	// Create VT100 emulator with size
	t.vt = vt10x.New(vt10x.WithSize(cfg.Width, cfg.Height))

	return t, nil
}

// Start starts the shell process and begins I/O loops
func (t *Terminal) Start() error {
	t.mu.Lock()
	if t.running {
		t.mu.Unlock()
		return fmt.Errorf("terminal already running")
	}
	t.running = true
	t.mu.Unlock()

	// Create PTY
	if err := t.startPTY(); err != nil {
		t.setError(err)
		t.mu.Lock()
		t.running = false
		t.mu.Unlock()
		return fmt.Errorf("failed to start PTY: %w", err)
	}

	// Start goroutines
	t.wg.Add(2)
	go t.readLoop()
	go t.writeLoop()

	log.Printf("[TERMINAL] Started: %dx%d shell=%s", t.width, t.height, t.shell)
	return nil
}

// readLoop reads from PTY and feeds to emulator
func (t *Terminal) readLoop() {
	defer t.wg.Done()
	defer log.Println("[TERMINAL] Read loop exited")

	buf := make([]byte, 4096)

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

		// Set read timeout to allow checking context
		if t.pty != nil {
			t.pty.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		}

		// Read from PTY
		n, err := t.pty.Read(buf)

		if n > 0 {
			// Write to vt10x parser - it will parse ANSI sequences and update its internal state
			if _, writeErr := t.vt.Write(buf[:n]); writeErr != nil {
				log.Printf("[TERMINAL] vt10x write error: %v", writeErr)
			}

			// Update our screen buffer from vt10x state
			t.updateScreenFromVT10x()

			// Signal update (non-blocking)
			select {
			case t.updateChan <- struct{}{}:
			default:
			}

			// Invalidate window to trigger redraw
			if t.window != nil {
				t.window.Invalidate()
			}
		}

		if err != nil {
			if err == io.EOF {
				log.Println("[TERMINAL] PTY closed (EOF)")
				return
			}
			// Ignore timeout errors
			if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() {
				continue
			}
			log.Printf("[TERMINAL] Read error: %v", err)
			t.setError(err)
			return
		}
	}
}

// writeLoop writes input to PTY
func (t *Terminal) writeLoop() {
	defer t.wg.Done()
	defer log.Println("[TERMINAL] Write loop exited")

	for {
		select {
		case <-t.ctx.Done():
			return
		case data := <-t.inputChan:
			if t.pty != nil {
				if _, err := t.pty.Write(data); err != nil {
					log.Printf("[TERMINAL] Write error: %v", err)
					t.setError(err)
					return
				}
			}
		}
	}
}

// Write sends input to PTY
func (t *Terminal) Write(data []byte) error {
	if !t.IsRunning() {
		return fmt.Errorf("terminal not running")
	}

	select {
	case t.inputChan <- data:
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("write timeout")
	}
}

// GetScreen returns current screen buffer
func (t *Terminal) GetScreen() *ScreenBuffer {
	return t.screen
}

// Close stops the terminal
func (t *Terminal) Close() error {
	log.Println("[TERMINAL] Close() called")

	t.mu.Lock()
	if !t.running {
		t.mu.Unlock()
		return nil
	}
	t.running = false
	t.mu.Unlock()

	// Cancel context to stop goroutines
	t.cancel()

	// Close PTY (will cause shell to exit)
	if t.pty != nil {
		t.pty.Close()
	}

	// Kill process if still running
	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}

	// Wait for goroutines with timeout
	done := make(chan struct{})
	go func() {
		t.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[TERMINAL] Cleanup complete")
	case <-time.After(2 * time.Second):
		log.Println("[TERMINAL] Cleanup timeout")
	}

	return nil
}

// IsRunning returns whether terminal is active
func (t *Terminal) IsRunning() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.running
}

// GetLastError returns last error encountered
func (t *Terminal) GetLastError() error {
	t.errorMu.Lock()
	defer t.errorMu.Unlock()
	return t.lastError
}

func (t *Terminal) setError(err error) {
	t.errorMu.Lock()
	t.lastError = err
	t.errorMu.Unlock()
}

// getEnvironment returns environment variables for shell
func (t *Terminal) getEnvironment() []string {
	env := os.Environ()

	// Set TERM
	env = append(env, "TERM=xterm-256color")

	// Set terminal size
	env = append(env, fmt.Sprintf("COLUMNS=%d", t.width))
	env = append(env, fmt.Sprintf("LINES=%d", t.height))

	// Platform-specific
	if runtime.GOOS != "windows" {
		// Unix: Set COLORTERM for better color support
		env = append(env, "COLORTERM=truecolor")
	}

	// Merge with user-provided env
	env = append(env, t.env...)

	return env
}

// updateScreenFromVT10x updates our screen buffer from vt10x's parsed state
func (t *Terminal) updateScreenFromVT10x() {
	// Lock vt10x state while reading
	t.vt.Lock()
	defer t.vt.Unlock()

	// Get terminal dimensions
	cols, rows := t.vt.Size()

	// Attribute bit masks (from vt10x source)
	const (
		attrBold      = 1 << 0
		attrDim       = 1 << 1
		attrItalic    = 1 << 2
		attrUnderline = 1 << 3
		attrBlink     = 1 << 4
		attrReverse   = 1 << 5
	)

	// Update each cell from vt10x
	for y := 0; y < rows && y < t.height; y++ {
		for x := 0; x < cols && x < t.width; x++ {
			// Get glyph from vt10x
			glyph := t.vt.Cell(x, y)

			// Convert vt10x.Color to our color format
			fg := vt10xColorToNRGBA(uint32(glyph.FG))
			bg := vt10xColorToNRGBA(uint32(glyph.BG))

			// Handle reverse video
			if glyph.Mode&attrReverse != 0 {
				fg, bg = bg, fg
			}

			// Set cell in our buffer
			t.screen.SetCell(x, y, Cell{
				Rune:      glyph.Char,
				FG:        fg,
				BG:        bg,
				Bold:      glyph.Mode&attrBold != 0,
				Dim:       glyph.Mode&attrDim != 0,
				Italic:    glyph.Mode&attrItalic != 0,
				Underline: glyph.Mode&attrUnderline != 0,
				Blink:     glyph.Mode&attrBlink != 0,
				Reverse:   glyph.Mode&attrReverse != 0,
			})
		}
	}

	// Update cursor position
	cursor := t.vt.Cursor()
	t.screen.SetCursor(cursor.X, cursor.Y)
}
