package appcore

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/javanhut/ProjectVem/internal/editor"
	"github.com/javanhut/ProjectVem/internal/filesystem"
	"github.com/javanhut/ProjectVem/internal/fonts"
	"github.com/javanhut/ProjectVem/internal/panes"
)

type mode string

type visualModeType int

const (
	visualModeNone visualModeType = iota
	visualModeChar                // Character-wise selection (v)
	visualModeLine                // Line-wise selection (Shift+V)
)

type SearchMatch struct {
	Line int
	Col  int
	Len  int
}

type FuzzyMatch struct {
	FilePath string
	Score    int
	Indices  []int
}

const (
	modeNormal      mode = "NORMAL"
	modeInsert      mode = "INSERT"
	modeVisual      mode = "VISUAL"
	modeDelete      mode = "DELETE"
	modeCommand     mode = "COMMAND"
	modeExplorer    mode = "EXPLORER"
	modeSearch      mode = "SEARCH"
	modeFuzzyFinder mode = "FUZZY_FINDER"
)

const caretBlinkInterval = 600 * time.Millisecond

var (
	highlightColor    = color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x55}
	selectionColor    = color.NRGBA{R: 0x1c, G: 0x39, B: 0x60, A: 0x99}
	background        = color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}
	statusBg          = color.NRGBA{R: 0x12, G: 0x17, B: 0x22, A: 0xff}
	headerColor       = color.NRGBA{R: 0xa1, G: 0xc6, B: 0xff, A: 0xff}
	cursorColor       = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	focusBorder       = color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
	searchMatchColor  = color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0x77}
	currentMatchColor = color.NRGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xaa}

	// Pane colors
	activePaneBg   = color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff} // Same as background (active is brighter)
	inactivePaneBg = color.NRGBA{R: 0x14, G: 0x18, B: 0x24, A: 0xff} // 15% darker
	paneSeparator  = color.NRGBA{R: 0x30, G: 0x35, B: 0x44, A: 0xff} // Subtle gray line
)

type appState struct {
	theme              *material.Theme
	bufferMgr          *editor.BufferManager
	paneManager        *panes.PaneManager
	fileTree           *filesystem.FileTree
	mode               mode
	status             string
	lastKey            string
	focusTag           *int
	pendingCount       int
	pendingGoto        bool
	pendingScroll      bool
	pendingPaneCmd     bool
	visualMode         visualModeType
	visualStartLine    int
	visualStartCol     int
	skipNextEdit       bool
	skipNextFileOpEdit bool
	skipNextSearchEdit bool
	skipNextFuzzyEdit  bool
	caretVisible       bool
	nextBlink          time.Time
	caretReset         bool
	clipLines          []string
	cmdText            string
	window             *app.Window

	// Explorer state
	explorerVisible      bool
	explorerWidth        int
	explorerFocused      bool
	explorerListPosition layout.List

	// File operation state
	fileOpMode         string
	fileOpInput        string
	fileOpOriginalName string
	fileOpTarget       *filesystem.TreeNode

	// Search state
	searchPattern   string
	searchMatches   []SearchMatch
	currentMatchIdx int
	searchActive    bool

	// Fuzzy finder state
	fuzzyFinderActive      bool
	fuzzyFinderInput       string
	fuzzyFinderFiles       []string
	fuzzyFinderMatches     []FuzzyMatch
	fuzzyFinderSelectedIdx int

	// Modifier tracking (some platforms don't report modifiers correctly)
	ctrlPressed  bool
	shiftPressed bool

	// Fullscreen state tracking
	currentWindowMode app.WindowMode
	wasFullscreen     bool

	// Viewport scrolling state
	viewportTopLine   int // First visible line in viewport (0-based)
	scrollOffsetLines int // Context lines around cursor (Vim's scrolloff)
	listPosition      layout.List
}

func Run(w *app.Window) error {
	state := newAppState()
	return state.run(w)
}

func (s *appState) run(w *app.Window) error {
	s.window = w
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.ConfigEvent:
			// Track window mode changes (fullscreen, maximized, etc.)
			s.currentWindowMode = e.Config.Mode
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			s.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func newAppState() *appState {
	theme := material.NewTheme()

	// Try to load JetBrains Mono Nerd Font, fall back to gofont if it fails
	customFonts, err := fonts.Collection()
	if err != nil {
		log.Printf("Warning: Failed to load JetBrains Mono Nerd Font, using default: %v", err)
		theme.Shaper = text.NewShaper(
			text.NoSystemFonts(),
			text.WithCollection(gofont.Collection()),
		)
	} else {
		log.Printf("Loaded JetBrains Mono Nerd Font successfully")
		theme.Shaper = text.NewShaper(
			text.NoSystemFonts(),
			text.WithCollection(customFonts),
		)
	}

	buf := editor.NewBuffer(strings.TrimSpace(sampleBuffer))
	bufferMgr := editor.NewBufferManagerWithBuffer(buf)

	// Initialize pane manager with the initial buffer (index 0)
	paneManager := panes.NewPaneManager(0)

	// Initialize file tree from current directory
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "."
	}
	fileTree, err := filesystem.NewFileTree(workDir)
	if err != nil {
		fileTree = nil
	} else {
		fileTree.LoadInitial()
	}

	return &appState{
		theme:                theme,
		bufferMgr:            bufferMgr,
		paneManager:          paneManager,
		fileTree:             fileTree,
		mode:                 modeNormal,
		status:               "Ready",
		focusTag:             new(int),
		visualMode:           visualModeNone,
		visualStartLine:      0,
		visualStartCol:       0,
		caretVisible:         true,
		explorerVisible:      false,
		explorerWidth:        275,
		explorerFocused:      false,
		explorerListPosition: layout.List{Axis: layout.Vertical},
		currentWindowMode:    app.Windowed,
		wasFullscreen:        false,
		viewportTopLine:      0,
		scrollOffsetLines:    3,
		listPosition:         layout.List{Axis: layout.Vertical},
	}
}

// activeBuffer returns the buffer for the active pane.
func (s *appState) activeBuffer() *editor.Buffer {
	if s.paneManager == nil || s.bufferMgr == nil {
		return nil
	}

	activePane := s.paneManager.ActivePane()
	if activePane == nil {
		return nil
	}

	return s.bufferMgr.GetBuffer(activePane.BufferIndex)
}

// activePaneViewportTop returns the viewport top line for the active pane.
func (s *appState) activePaneViewportTop() int {
	if s.paneManager == nil {
		return 0
	}

	activePane := s.paneManager.ActivePane()
	if activePane == nil {
		return 0
	}

	return activePane.ViewportTop
}

// setActivePaneViewportTop sets the viewport top line for the active pane.
func (s *appState) setActivePaneViewportTop(line int) {
	if s.paneManager == nil {
		return
	}

	activePane := s.paneManager.ActivePane()
	if activePane != nil {
		activePane.SetViewportTop(line)
	}
}

func (s *appState) layout(gtx layout.Context) layout.Dimensions {
	s.handleEvents(gtx)
	s.updateCaretBlink(gtx)

	canvas := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, background)
	canvas.Pop()

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.drawHeader(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if s.explorerVisible && s.fileTree != nil {
				// Horizontal split: explorer | panes
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return s.drawFileExplorer(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return s.drawPanes(gtx)
					}),
				)
			}
			return s.drawPanes(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if s.mode == modeCommand {
				return s.drawCommandBar(gtx)
			}
			return s.drawStatusBar(gtx)
		}),
	)

	// Draw fuzzy finder overlay on top if active
	if s.fuzzyFinderActive {
		s.drawFuzzyFinder(gtx)
	}

	return dims
}

func (s *appState) handleEvents(gtx layout.Context) {
	event.Op(gtx.Ops, s.focusTag)
	if s.mode == modeInsert || s.mode == modeCommand || s.mode == modeSearch || s.mode == modeFuzzyFinder {
		key.InputHintOp{Tag: s.focusTag, Hint: key.HintText}.Add(gtx.Ops)
		gtx.Execute(key.SoftKeyboardCmd{Show: true})
	} else {
		// In normal mode, we still want to receive Tab events
		gtx.Execute(key.SoftKeyboardCmd{Show: false})
	}
	gtx.Execute(key.FocusCmd{Tag: s.focusTag})
	for {
		ev, ok := gtx.Event(
			key.FocusFilter{Target: s.focusTag},
			key.Filter{Focus: s.focusTag},
		)
		if !ok {
			break
		}
		switch e := ev.(type) {
		case key.FocusEvent:
			if e.Focus {
				s.status = "Ready"
			}
		case key.Event:
			// Track Ctrl key press/release
			if e.Name == key.NameCtrl {
				s.ctrlPressed = (e.State == key.Press)
				if e.State == key.Press {
					log.Printf("⌨ [CTRL] Pressed")
				} else {
					log.Printf("⌨ [CTRL] Released")
				}
				continue
			}
			// Track Shift key press/release
			if e.Name == key.NameShift {
				s.shiftPressed = (e.State == key.Press)
				if e.State == key.Press {
					log.Printf("⌨ [SHIFT] Pressed (waiting for character key...)")
				} else {
					log.Printf("⌨ [SHIFT] Released")
				}
				continue
			}
			// Track Alt key press/release
			if e.Name == key.NameAlt {
				log.Printf("⌨ [ALT] %v", e.State)
				continue
			}

			// Save modifier state before handling key
			hadCtrl := s.ctrlPressed
			hadShift := s.shiftPressed

			s.handleKey(e)

			// Smart reset: If modifiers are still set after handleKey and we're in certain modes,
			// it likely means the user released the modifier but platform didn't send Release event.
			// Reset modifiers after successful command execution to prevent them from sticking.
			// Exception: Don't reset if we just entered a mode or are waiting for a pane command
			// Exception: In INSERT mode, don't reset modifiers for printable keys that generate EditEvents
			//            (the EditEvent needs the modifier state to determine capitalization)
			shouldResetModifiers := false

			if s.mode == modeNormal || s.mode == modeCommand || s.mode == modeExplorer || s.mode == modeSearch || s.mode == modeFuzzyFinder {
				// In non-insert modes, reset modifiers after command keys unless waiting for pane command
				shouldResetModifiers = !s.pendingPaneCmd
			} else if s.mode == modeInsert {
				// In INSERT mode, only reset for special keys (Escape, arrow keys, function keys)
				// Don't reset for: printable keys, Tab (needs Shift state), or when pending pane command
				isSpecialKey := (e.Name == key.NameEscape ||
					e.Name == key.NameLeftArrow || e.Name == key.NameRightArrow ||
					e.Name == key.NameUpArrow || e.Name == key.NameDownArrow ||
					e.Name == key.NameDeleteBackward || e.Name == key.NameDeleteForward)
				shouldResetModifiers = isSpecialKey && !s.pendingPaneCmd
			}

			if shouldResetModifiers {
				if hadCtrl && s.ctrlPressed {
					s.ctrlPressed = false
				}
				if hadShift && s.shiftPressed {
					s.shiftPressed = false
				}
			}
		case key.EditEvent:
			if e.Text == "" {
				continue
			}

			// Handle file operation input if active
			if s.fileOpMode == "rename" || s.fileOpMode == "create" {
				if s.skipNextFileOpEdit {
					s.skipNextFileOpEdit = false
					continue
				}
				s.appendFileOpInput(e.Text)
				continue
			}

			// Check for colon to enter command mode (except in INSERT and COMMAND modes)
			if e.Text == ":" && s.mode != modeInsert && s.mode != modeCommand {
				s.enterCommandMode()
				continue
			}

			switch s.mode {
			case modeInsert:
				if s.skipNextEdit {
					s.skipNextEdit = false
					log.Printf("✓ [FIX_ACTIVE] Skipped EditEvent %q (already handled by KeyEvent)", e.Text)
					continue
				}
				// Platform didn't send KeyEvent, only EditEvent - use it
				log.Printf("⚠ [PLATFORM_QUIRK] EditEvent %q arrived without KeyEvent (platform limitation)", e.Text)
				s.insertText(e.Text)
				// Reset modifiers after EditEvent insertion
				if s.shiftPressed {
					log.Printf("⚠ [PLATFORM_QUIRK] Resetting Shift after EditEvent")
					s.shiftPressed = false
				}
				if s.ctrlPressed {
					log.Printf("⚠ [PLATFORM_QUIRK] Resetting Ctrl after EditEvent")
					s.ctrlPressed = false
				}
			case modeCommand:
				s.appendCommandText(e.Text)
				// Reset modifiers after text insertion to prevent sticking
				if s.shiftPressed {
					log.Printf("[EDIT_RESET] Resetting Shift after command text=%q", e.Text)
					s.shiftPressed = false
				}
				if s.ctrlPressed {
					log.Printf("[EDIT_RESET] Resetting Ctrl after command text=%q", e.Text)
					s.ctrlPressed = false
				}
			case modeSearch:
				if s.skipNextSearchEdit {
					s.skipNextSearchEdit = false
					continue
				}
				s.appendSearchText(e.Text)
				// Reset modifiers after text insertion to prevent sticking
				if s.shiftPressed {
					log.Printf("[EDIT_RESET] Resetting Shift after search text=%q", e.Text)
					s.shiftPressed = false
				}
				if s.ctrlPressed {
					log.Printf("[EDIT_RESET] Resetting Ctrl after search text=%q", e.Text)
					s.ctrlPressed = false
				}
			case modeFuzzyFinder:
				if s.skipNextFuzzyEdit {
					s.skipNextFuzzyEdit = false
					continue
				}
				s.appendFuzzyInput(e.Text)
				// Reset modifiers after text insertion to prevent sticking
				if s.shiftPressed {
					log.Printf("[EDIT_RESET] Resetting Shift after fuzzy text=%q", e.Text)
					s.shiftPressed = false
				}
				if s.ctrlPressed {
					log.Printf("[EDIT_RESET] Resetting Ctrl after fuzzy text=%q", e.Text)
					s.ctrlPressed = false
				}
			}
		}
	}
}

func (s *appState) drawHeader(gtx layout.Context) layout.Dimensions {
	label := material.H5(s.theme, "Vem")
	label.Color = headerColor
	label.Font.Weight = font.Bold
	inset := layout.Inset{
		Top:    unit.Dp(12),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(4),
		Left:   unit.Dp(16),
	}
	return inset.Layout(gtx, label.Layout)
}

func (s *appState) drawBuffer(gtx layout.Context) layout.Dimensions {
	lines := s.activeBuffer().LineCount()
	cursorLine := s.activeBuffer().Cursor().Line
	selStart, selEnd, hasSel := s.visualSelectionRange()
	cursorCol := s.activeBuffer().Cursor().Col

	// Calculate approximate lines per page for viewport scrolling
	// Use a rough estimate: line height ~20dp, inset ~16dp top+bottom
	lineHeightDp := 20
	insetDp := 32
	availableHeight := gtx.Constraints.Max.Y - gtx.Dp(unit.Dp(insetDp))
	linesPerPage := availableHeight / gtx.Dp(unit.Dp(lineHeightDp))
	if linesPerPage < 1 {
		linesPerPage = 1
	}

	// Ensure cursor is visible in viewport
	s.ensureCursorVisible(linesPerPage)

	// Set scroll position to viewport top line
	s.listPosition.Position.First = s.viewportTopLine
	s.listPosition.Position.Offset = 0

	// Draw focus border on the left edge if editor is focused (not in explorer mode)
	editorFocused := !s.explorerFocused && s.explorerVisible
	if editorFocused {
		borderWidth := 3
		borderRect := clip.Rect{
			Min: image.Pt(0, 0),
			Max: image.Pt(borderWidth, gtx.Constraints.Max.Y),
		}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, focusBorder)
		borderRect.Pop()
	}

	inset := layout.Inset{
		Top:    unit.Dp(8),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(8),
		Left:   unit.Dp(16),
	}
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return s.listPosition.Layout(gtx, lines, func(gtx layout.Context, index int) layout.Dimensions {
			// Important: Add gutter BEFORE expanding tabs so tab stops align with cursor positioning
			lineText := fmt.Sprintf("%4d  %s", index+1, s.activeBuffer().Line(index))
			lineText = expandTabs(lineText, 4)
			label := material.Body1(s.theme, lineText)
			label.Font.Typeface = "JetBrainsMono"
			label.Color = color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
			macro := op.Record(gtx.Ops)
			dims := label.Layout(gtx)
			call := macro.Stop()

			// Draw selection highlighting
			if s.visualMode == visualModeChar {
				s.drawCharSelection(gtx, index, dims.Size.Y)
			} else if hasSel && index >= selStart && index <= selEnd {
				// Line selection
				rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}.Push(gtx.Ops)
				paint.Fill(gtx.Ops, selectionColor)
				rect.Pop()
			}

			// Draw cursor line highlight (skip in character-wise visual mode to show precise selection)
			if index == cursorLine && s.visualMode != visualModeChar {
				rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}.Push(gtx.Ops)
				paint.Fill(gtx.Ops, highlightColor)
				rect.Pop()
			}

			// Draw search highlights
			if s.searchActive && len(s.searchMatches) > 0 {
				s.drawSearchHighlights(gtx, index, dims.Size.Y)
			}

			call.Add(gtx.Ops)

			if index == cursorLine {
				gutter := fmt.Sprintf("%4d  ", index+1)
				prefix := s.activeBuffer().LinePrefix(index, cursorCol)
				charUnder := s.getCharAtCursor(index, cursorCol)
				s.drawCursor(gtx, gutter, prefix, charUnder, dims.Size.Y)
			}
			return dims
		})
	})
}

func (s *appState) drawSearchHighlights(gtx layout.Context, lineIdx int, lineHeight int) {
	gutter := fmt.Sprintf("%4d  ", lineIdx+1)
	gutterWidth := s.measureTextWidth(gtx, gutter)

	for i, match := range s.searchMatches {
		if match.Line != lineIdx {
			continue
		}

		// Calculate position of match
		lineContent := s.activeBuffer().Line(lineIdx)
		prefix := string([]rune(lineContent)[:match.Col])
		matchText := string([]rune(lineContent)[match.Col : match.Col+match.Len])

		prefixWidth := s.measureTextWidth(gtx, prefix)
		matchWidth := s.measureTextWidth(gtx, matchText)

		// Determine highlight color (current match vs other matches)
		highlightCol := searchMatchColor
		if i == s.currentMatchIdx {
			highlightCol = currentMatchColor
		}

		// Draw highlight rectangle
		x := gutterWidth + prefixWidth
		rect := clip.Rect{
			Min: image.Pt(x, 0),
			Max: image.Pt(x+matchWidth, lineHeight),
		}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, highlightCol)
		rect.Pop()
	}
}

func (s *appState) drawCharSelection(gtx layout.Context, lineIdx int, lineHeight int) {
	startLine, startCol, endLine, endCol, ok := s.visualSelectionRangeChar()
	if !ok {
		return
	}

	// Check if this line is in the selection range
	if lineIdx < startLine || lineIdx > endLine {
		return
	}

	lineContent := s.activeBuffer().Line(lineIdx)
	runes := []rune(lineContent)

	// Calculate selection range for this line
	selStart := 0
	selEnd := len(runes)

	if lineIdx == startLine {
		selStart = startCol
		if selStart > len(runes) {
			selStart = len(runes)
		}
	}
	if lineIdx == endLine {
		selEnd = endCol
		if selEnd > len(runes) {
			selEnd = len(runes)
		}
	}

	// Nothing to select on this line
	if selStart >= selEnd {
		return
	}

	// Measure text widths
	gutter := fmt.Sprintf("%4d  ", lineIdx+1)
	gutterWidth := s.measureTextWidth(gtx, gutter)

	prefix := string(runes[:selStart])
	prefixWidth := s.measureTextWidth(gtx, prefix)

	selected := string(runes[selStart:selEnd])
	selectedWidth := s.measureTextWidth(gtx, selected)

	// Draw highlight rectangle
	x := gutterWidth + prefixWidth
	rect := clip.Rect{
		Min: image.Pt(x, 0),
		Max: image.Pt(x+selectedWidth, lineHeight),
	}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, selectionColor)
	rect.Pop()
}

func (s *appState) drawStatusBar(gtx layout.Context) layout.Dimensions {
	var status string

	// If search mode is active, show search prompt
	if s.mode == modeSearch {
		status = "/" + s.searchPattern
	} else if s.fileOpMode != "" {
		// If file operation is active, show ONLY the file operation prompt for clarity
		status = s.getFileOpPrompt()
	} else {
		// Check if we have an active buffer
		buf := s.activeBuffer()
		if buf == nil {
			// No active buffer - show minimal status
			paneInfo := ""
			if s.paneManager != nil {
				paneCount := s.paneManager.PaneCount()
				if paneCount > 1 {
					allPanes := s.paneManager.AllPanes()
					activeIdx := 1
					for i, p := range allPanes {
						if p == s.paneManager.ActivePane() {
							activeIdx = i + 1
							break
						}
					}
					paneInfo = fmt.Sprintf(" | PANE %d/%d", activeIdx, paneCount)
				}
			}
			status = fmt.Sprintf("MODE %s | No active buffer%s | %s", s.mode, paneInfo, s.status)
		} else {
			cur := buf.Cursor()

			// Build status line with file info
			fileName := buf.FilePath()
			if fileName == "" {
				fileName = "[No Name]"
			}

			modFlag := ""
			if buf.Modified() {
				modFlag = " [+]"
			}

			// Add pane information
			paneInfo := ""
			if s.paneManager != nil && s.paneManager.PaneCount() > 1 {
				// Find active pane index
				allPanes := s.paneManager.AllPanes()
				activeIdx := 1
				for i, p := range allPanes {
					if p == s.paneManager.ActivePane() {
						activeIdx = i + 1
						break
					}
				}
				paneInfo = fmt.Sprintf(" | PANE %d/%d", activeIdx, s.paneManager.PaneCount())
			}

			// Add fullscreen indicator
			fullscreenInfo := ""
			if s.currentWindowMode == app.Fullscreen {
				fullscreenInfo = " | FULLSCREEN"
			}

			// Add zoom indicator
			zoomInfo := ""
			if s.paneManager != nil && s.paneManager.IsZoomed() {
				zoomInfo = " | ZOOMED"
			}

			status = fmt.Sprintf("MODE %s | FILE %s%s | CURSOR %d:%d%s%s%s | %s",
				s.mode, fileName, modFlag, cur.Line+1, cur.Col+1, paneInfo, fullscreenInfo, zoomInfo, s.status,
			)
		}
	}

	label := material.Body2(s.theme, status)
	label.Font.Typeface = "JetBrainsMono"
	label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	macro := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(8)).Layout(gtx, label.Layout)
	call := macro.Stop()

	rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, statusBg)
	rect.Pop()

	call.Add(gtx.Ops)
	return layout.Dimensions{
		Size: image.Pt(gtx.Constraints.Max.X, dims.Size.Y),
	}
}

func (s *appState) drawFuzzyFinder(gtx layout.Context) layout.Dimensions {
	// Overlay background (semi-transparent)
	overlayBg := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xcc}
	overlayRect := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, overlayBg)
	overlayRect.Pop()

	// Calculate centered fuzzy finder dimensions
	finderWidth := gtx.Constraints.Max.X * 3 / 4
	if finderWidth > 800 {
		finderWidth = 800
	}
	finderHeight := gtx.Constraints.Max.Y * 2 / 3
	if finderHeight > 600 {
		finderHeight = 600
	}

	offsetX := (gtx.Constraints.Max.X - finderWidth) / 2
	offsetY := (gtx.Constraints.Max.Y - finderHeight) / 4

	// Draw fuzzy finder box
	boxBg := color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}
	boxBorder := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}

	// Position the fuzzy finder
	offset := op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops)
	defer offset.Pop()

	// Draw border
	borderRect := clip.Rect{Max: image.Pt(finderWidth, finderHeight)}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, boxBorder)
	borderRect.Pop()

	// Draw background (slightly inset for border effect)
	bgRect := clip.Rect{
		Min: image.Pt(2, 2),
		Max: image.Pt(finderWidth-2, finderHeight-2),
	}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, boxBg)
	bgRect.Pop()

	// Constrain drawing to fuzzy finder area
	gtx.Constraints.Max.X = finderWidth - 4
	gtx.Constraints.Max.Y = finderHeight - 4

	inset := layout.Inset{
		Top:    unit.Dp(8),
		Right:  unit.Dp(8),
		Bottom: unit.Dp(8),
		Left:   unit.Dp(8),
	}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Input field
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				prompt := "Fuzzy Finder: " + s.fuzzyFinderInput
				label := material.Body1(s.theme, prompt)
				label.Font.Typeface = "JetBrainsMono"
				label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
				return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, label.Layout)
			}),
			// Match count
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				matchInfo := fmt.Sprintf("%d matches", len(s.fuzzyFinderMatches))
				label := material.Body2(s.theme, matchInfo)
				label.Font.Typeface = "JetBrainsMono"
				label.Color = color.NRGBA{R: 0xa1, G: 0xc6, B: 0xff, A: 0xff}
				return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, label.Layout)
			}),
			// Results list
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				list := layout.List{Axis: layout.Vertical}
				return list.Layout(gtx, len(s.fuzzyFinderMatches), func(gtx layout.Context, index int) layout.Dimensions {
					match := s.fuzzyFinderMatches[index]

					// Highlight selected item
					if index == s.fuzzyFinderSelectedIdx {
						selectedBg := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x88}
						rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(24)))}.Push(gtx.Ops)
						paint.Fill(gtx.Ops, selectedBg)
						rect.Pop()
					}

					// Draw file path with highlighted matched characters
					label := material.Body2(s.theme, match.FilePath)
					label.Font.Typeface = "JetBrainsMono"
					if index == s.fuzzyFinderSelectedIdx {
						label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
					} else {
						label.Color = color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
					}

					return layout.Inset{
						Top:    unit.Dp(2),
						Bottom: unit.Dp(2),
						Left:   unit.Dp(4),
					}.Layout(gtx, label.Layout)
				})
			}),
		)
	})
}

func (s *appState) drawCommandBar(gtx layout.Context) layout.Dimensions {
	prompt := ":" + s.cmdText
	label := material.Body2(s.theme, prompt)
	label.Font.Typeface = "JetBrainsMono"
	label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	macro := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(8)).Layout(gtx, label.Layout)
	call := macro.Stop()

	rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, statusBg)
	rect.Pop()

	call.Add(gtx.Ops)
	return layout.Dimensions{
		Size: image.Pt(gtx.Constraints.Max.X, dims.Size.Y),
	}
}

func (s *appState) drawFileExplorer(gtx layout.Context) layout.Dimensions {
	if s.fileTree == nil {
		return layout.Dimensions{}
	}

	explorerBg := color.NRGBA{R: 0x15, G: 0x1a, B: 0x28, A: 0xff}
	selectedBg := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x88}
	dirColor := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
	fileColor := color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}

	width := gtx.Dp(unit.Dp(s.explorerWidth))
	gtx.Constraints.Max.X = width
	gtx.Constraints.Min.X = width

	// Draw background with focus border if explorer is focused
	macro := op.Record(gtx.Ops)
	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Header showing current path
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			currentPath := s.fileTree.CurrentPath()
			pathLabel := material.Body2(s.theme, currentPath)
			pathLabel.Font.Typeface = "JetBrainsMono"
			pathLabel.Color = color.NRGBA{R: 0xa1, G: 0xc6, B: 0xff, A: 0xff}

			return layout.Inset{
				Top:    unit.Dp(4),
				Right:  unit.Dp(8),
				Bottom: unit.Dp(4),
				Left:   unit.Dp(8),
			}.Layout(gtx, pathLabel.Layout)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			nodes := s.fileTree.GetFlatList()
			selectedIndex := s.fileTree.SelectedIndex()

			inset := layout.Inset{
				Top:    unit.Dp(8),
				Right:  unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(8),
			}

			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return s.explorerListPosition.Layout(gtx, len(nodes), func(gtx layout.Context, index int) layout.Dimensions {
					node := nodes[index]

					// Draw selection highlight
					if index == selectedIndex {
						rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(20)))}.Push(gtx.Ops)
						paint.Fill(gtx.Ops, selectedBg)
						rect.Pop()
					}

					// Build line text with indentation and icon
					indent := strings.Repeat("  ", node.Depth)

					var icon string
					if node.IsDir {
						// Directory: expand/collapse icon + folder icon
						expandIcon := node.GetExpandIcon()
						folderIcon := node.GetIcon()
						icon = expandIcon + " " + folderIcon + " "
					} else {
						// File: file type icon with spacing
						fileIcon := node.GetIcon()
						icon = "  " + fileIcon + " "
					}

					lineText := indent + icon + node.Name

					// Render text
					label := material.Body2(s.theme, lineText)
					label.Font.Typeface = "JetBrainsMono"
					if node.IsDir {
						label.Color = dirColor
					} else {
						label.Color = fileColor
					}

					return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, label.Layout)
				})
			})
		}),
	)
	call := macro.Stop()

	// Fill background
	rect := clip.Rect{Max: image.Pt(width, gtx.Constraints.Max.Y)}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, explorerBg)
	rect.Pop()

	// Draw focus border on the right edge if explorer is focused
	if s.explorerFocused {
		borderWidth := 3
		borderRect := clip.Rect{
			Min: image.Pt(width-borderWidth, 0),
			Max: image.Pt(width, gtx.Constraints.Max.Y),
		}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, focusBorder)
		borderRect.Pop()
	}

	call.Add(gtx.Ops)

	return layout.Dimensions{
		Size: image.Pt(width, dims.Size.Y),
	}
}

func (s *appState) handleKey(ev key.Event) {
	// Special case: Tab events only come through on Release (consumed by focus system on Press)
	// Process Tab on Release, all other keys on Press
	isTabRelease := (ev.Name == key.NameTab && ev.State == key.Release)

	if ev.State != key.Press && !isTabRelease {
		return
	}
	s.lastKey = describeKey(ev)

	// Clear skipNextEdit at the start of each KeyEvent to prevent stale state
	// This ensures we only skip EditEvents that correspond to THIS KeyEvent
	if s.mode == modeInsert {
		s.skipNextEdit = false
	}

	// Handle file operation input if active
	if s.fileOpMode != "" {
		if s.handleFileOpKey(ev) {
			return
		}
	}

	// Handle Ctrl+S prefix for pane commands
	if s.ctrlPressed && strings.ToLower(string(ev.Name)) == "s" && !s.pendingPaneCmd {
		s.pendingPaneCmd = true
		s.status = "Pane: v=vsplit h=hsplit Alt+hjkl=nav Tab=cycle Ctrl+X=close ==equalize o=zoom"
		return
	}

	if s.pendingPaneCmd {
		s.handlePaneCommand(ev)
		s.pendingPaneCmd = false
		return
	}

	// Debug logging
	modStr := s.formatModifiers(ev.Modifiers)
	log.Printf("[KEY] Key=%q Modifiers=%s Mode=%s ExplorerVisible=%v ExplorerFocused=%v",
		ev.Name, modStr, s.mode, s.explorerVisible, s.explorerFocused)

	// Phase 1: Try mode-specific keybindings first for COMMAND mode
	// (COMMAND mode keys should take priority over global shortcuts)
	if s.mode == modeCommand {
		if action := s.matchModeKeybinding(s.mode, ev); action != ActionNone {
			log.Printf("[MATCH] Mode-specific keybinding matched: Mode=%s Action=%v", s.mode, action)
			s.executeAction(action, ev)
			return
		}
	}

	// Phase 2: Try global keybindings (highest priority for other modes)
	if action := s.matchGlobalKeybinding(ev); action != ActionNone {
		log.Printf("[MATCH] Global keybinding matched: Action=%v", action)
		s.executeAction(action, ev)
		return
	}

	// Phase 3: Try mode-specific keybindings
	if action := s.matchModeKeybinding(s.mode, ev); action != ActionNone {
		log.Printf("[MATCH] Mode-specific keybinding matched: Mode=%s Action=%v", s.mode, action)
		s.executeAction(action, ev)
		return
	}

	// Phase 4: Handle special cases that need custom logic
	switch s.mode {
	case modeNormal:
		if s.handleNormalModeSpecial(ev) {
			log.Printf("[SPECIAL] Normal mode special handler matched")
			return
		}
	case modeInsert:
		if s.handleInsertModeSpecial(ev) {
			log.Printf("[SPECIAL] Insert mode special handler matched")
			return
		}
	case modeDelete:
		if s.handleDeleteModeSpecial(ev) {
			log.Printf("[SPECIAL] Delete mode special handler matched")
			return
		}
	case modeVisual:
		if s.handleVisualModeSpecial(ev) {
			log.Printf("[SPECIAL] Visual mode special handler matched")
			return
		}
	case modeCommand:
		return
	case modeExplorer:
		return
	}
	log.Printf("[NO_MATCH] No keybinding matched for key=%q modifiers=%s", ev.Name, modStr)
	s.status = "Waiting for motion"
}

func (s *appState) formatModifiers(mods key.Modifiers) string {
	var parts []string
	var tracked []string

	if mods.Contain(key.ModCtrl) {
		parts = append(parts, "Ctrl")
	}
	if s.ctrlPressed && !mods.Contain(key.ModCtrl) {
		tracked = append(tracked, "Ctrl")
	}

	if mods.Contain(key.ModShift) {
		parts = append(parts, "Shift")
	}
	if s.shiftPressed && !mods.Contain(key.ModShift) {
		tracked = append(tracked, "Shift")
	}

	if mods.Contain(key.ModAlt) {
		parts = append(parts, "Alt")
	}

	result := strings.Join(parts, "+")
	if len(tracked) > 0 {
		if result != "" {
			result += "+"
		}
		result += strings.Join(tracked, "+") + "(tracked)"
	}

	if result == "" {
		return "none"
	}
	return result
}

func (s *appState) handleNormalModeSpecial(ev key.Event) bool {
	if s.isColonKey(ev) {
		s.enterCommandMode()
		return true
	}

	if r, ok := s.printableKey(ev); ok {
		if unicode.IsDigit(r) {
			if s.handleCountDigit(int(r - '0')) {
				return true
			}
		}
		if s.pendingGoto {
			if s.handleGotoSequence(r) {
				return true
			}
			s.pendingGoto = false
		}
		if s.pendingScroll {
			if s.handleScrollSequence(r) {
				return true
			}
			s.pendingScroll = false
		}
		switch r {
		case 'G':
			s.gotoLineWithCount()
			return true
		case 'g':
			s.startGotoSequence()
			return true
		case 'z':
			s.startScrollSequence()
			return true
		}
	}
	return false
}

func (s *appState) handleInsertModeSpecial(ev key.Event) bool {
	if r, ok := s.printableKey(ev); ok {
		// FIX WORKING: Insert character immediately from KeyEvent (bypassing delayed EditEvent)
		log.Printf("✓ [FIX_ACTIVE] KeyEvent insert: %q | Shift=%v Ctrl=%v", string(r), s.shiftPressed, s.ctrlPressed)
		s.insertText(string(r))

		// Reset modifiers immediately after insertion
		modifiersReset := false
		if s.shiftPressed {
			s.shiftPressed = false
			modifiersReset = true
		}
		if s.ctrlPressed {
			s.ctrlPressed = false
			modifiersReset = true
		}
		if modifiersReset {
			log.Printf("✓ [FIX_ACTIVE] Modifiers reset (no sticking)")
		}

		// Mark that we should skip the corresponding EditEvent
		// Note: Sometimes platform doesn't send KeyEvent, only EditEvent
		s.skipNextEdit = true
		return true
	}
	return false
}

func (s *appState) handleDeleteModeSpecial(ev key.Event) bool {
	if r, ok := s.printableKey(ev); ok {
		if unicode.IsDigit(r) {
			s.handleCountDigit(int(r - '0'))
			if s.pendingCount > 0 {
				s.status = "DELETE line " + string(rune('0'+s.pendingCount)) + " (pending)"
			} else {
				s.status = "DELETE mode"
			}
			return true
		}
		if unicode.ToLower(r) == 'd' {
			s.executeDeleteCommand()
			return true
		}
	}
	return false
}

func (s *appState) handleVisualModeSpecial(ev key.Event) bool {
	if s.isColonKey(ev) {
		s.enterCommandMode()
		return true
	}

	if r, ok := s.printableKey(ev); ok {
		if unicode.IsDigit(r) && s.handleCountDigit(int(r-'0')) {
			return true
		}
		if s.pendingGoto {
			if s.handleGotoSequence(r) {
				return true
			}
			s.pendingGoto = false
		}
		if s.pendingScroll {
			if s.handleScrollSequence(r) {
				return true
			}
			s.pendingScroll = false
		}
		switch r {
		case 'G':
			s.gotoLineWithCount()
			return true
		case 'g':
			s.startGotoSequence()
			return true
		case 'z':
			s.startScrollSequence()
			return true
		}
	}
	return false
}

func (s *appState) moveCursor(direction string) {
	var moved bool
	switch direction {
	case "left":
		moved = s.activeBuffer().MoveLeft()
	case "right":
		moved = s.activeBuffer().MoveRight()
	case "up":
		moved = s.activeBuffer().MoveUp()
	case "down":
		moved = s.activeBuffer().MoveDown()
	default:
		return
	}
	if moved {
		s.setCursorStatus(fmt.Sprintf("Moved %s", direction))
	} else {
		s.status = fmt.Sprintf("Hit %s boundary", direction)
	}
}

func (s *appState) setCursorStatus(action string) {
	cur := s.activeBuffer().Cursor()
	s.status = fmt.Sprintf("%s → %d:%d", action, cur.Line+1, cur.Col+1)
	s.caretReset = true
}

func (s *appState) enterInsertMode() {
	if s.mode == modeInsert {
		return
	}
	s.mode = modeInsert
	s.skipNextEdit = true
	s.resetCount()
	s.status = "Switched to INSERT"
	s.caretReset = true
}

func (s *appState) getCharAtCursor(lineIdx, col int) string {
	line := s.activeBuffer().Line(lineIdx)
	if col >= len([]rune(line)) {
		return " "
	}
	runes := []rune(line)
	return string(runes[col])
}

func (s *appState) drawCursor(gtx layout.Context, gutter, prefix, charUnder string, height int) {
	if !s.caretVisible {
		return
	}

	full := gutter + prefix
	x := s.measureTextWidth(gtx, full)

	// Cursor drawing (debug logs removed for clarity)

	if s.mode == modeInsert {
		// Thin line cursor for INSERT mode
		width := gtx.Dp(unit.Dp(2))
		if width < 2 {
			width = 2
		}
		rect := image.Rect(x, 0, x+width, height)
		stack := clip.Rect(rect).Push(gtx.Ops)
		paint.Fill(gtx.Ops, cursorColor)
		stack.Pop()
	} else {
		// Block cursor for NORMAL/VISUAL/DELETE modes
		// Expand tab character for display (tabs should render as spaces)
		displayChar := charUnder
		if charUnder == "\t" {
			displayChar = expandTabs("\t", 4)
		}

		charWidth := s.measureTextWidth(gtx, charUnder)
		if charWidth < 8 {
			charWidth = 8
		}
		rect := image.Rect(x, 0, x+charWidth, height)
		stack := clip.Rect(rect).Push(gtx.Ops)
		paint.Fill(gtx.Ops, cursorColor)
		stack.Pop()

		// Draw the character on top of the cursor in contrasting color
		label := material.Body1(s.theme, displayChar)
		label.Font.Typeface = "JetBrainsMono"
		label.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		offset := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)
		label.Layout(gtx)
		offset.Pop()
	}
}

func (s *appState) measureTextWidth(gtx layout.Context, txt string) int {
	// Expand tabs to spaces before measuring so measurements match visual rendering
	expandedTxt := expandTabs(txt, 4)

	label := material.Body1(s.theme, expandedTxt)
	label.Font.Typeface = "JetBrainsMono"
	measureGtx := gtx
	measureGtx.Constraints = layout.Constraints{
		Min: image.Point{},
		Max: image.Point{X: math.MaxInt32, Y: math.MaxInt32},
	}
	macro := op.Record(measureGtx.Ops)
	dims := label.Layout(measureGtx)
	macro.Stop()

	return dims.Size.X
}

func (s *appState) updateCaretBlink(gtx layout.Context) {
	if s.caretReset {
		s.caretVisible = true
		s.nextBlink = gtx.Now.Add(caretBlinkInterval)
		s.caretReset = false
	}
	if s.mode != modeInsert {
		s.caretVisible = true
		return
	}
	if s.nextBlink.IsZero() {
		s.nextBlink = gtx.Now.Add(caretBlinkInterval)
	}
	if !s.nextBlink.After(gtx.Now) {
		s.caretVisible = !s.caretVisible
		s.nextBlink = gtx.Now.Add(caretBlinkInterval)
	}
	gtx.Execute(op.InvalidateCmd{At: s.nextBlink})
}

// ensureCursorVisible adjusts the viewport to ensure the cursor is visible with context lines.
// linesPerPage is the number of lines that fit in the current viewport.
func (s *appState) ensureCursorVisible(linesPerPage int) {
	if linesPerPage <= 0 {
		linesPerPage = 20 // Fallback default
	}

	cursorLine := s.activeBuffer().Cursor().Line
	totalLines := s.activeBuffer().LineCount()

	// Calculate viewport bounds
	viewportBottom := s.viewportTopLine + linesPerPage - 1

	// Check if cursor is above viewport (need to scroll up)
	if cursorLine < s.viewportTopLine+s.scrollOffsetLines {
		s.viewportTopLine = cursorLine - s.scrollOffsetLines
		if s.viewportTopLine < 0 {
			s.viewportTopLine = 0
		}
	}

	// Check if cursor is below viewport (need to scroll down)
	if cursorLine > viewportBottom-s.scrollOffsetLines {
		s.viewportTopLine = cursorLine - linesPerPage + s.scrollOffsetLines + 1
		if s.viewportTopLine < 0 {
			s.viewportTopLine = 0
		}
	}

	// Clamp viewport to valid range
	maxTopLine := totalLines - linesPerPage
	if maxTopLine < 0 {
		maxTopLine = 0
	}
	if s.viewportTopLine > maxTopLine {
		s.viewportTopLine = maxTopLine
	}
}

// scrollToCenter centers the cursor line in the viewport (Vim's zz command).
func (s *appState) scrollToCenter(linesPerPage int) {
	if linesPerPage <= 0 {
		linesPerPage = 20
	}
	cursorLine := s.activeBuffer().Cursor().Line
	s.viewportTopLine = cursorLine - (linesPerPage / 2)
	if s.viewportTopLine < 0 {
		s.viewportTopLine = 0
	}
	maxTopLine := s.activeBuffer().LineCount() - linesPerPage
	if maxTopLine < 0 {
		maxTopLine = 0
	}
	if s.viewportTopLine > maxTopLine {
		s.viewportTopLine = maxTopLine
	}
	s.status = "Centered cursor"
}

// scrollToTop positions cursor line at top of viewport (Vim's zt command).
func (s *appState) scrollToTop() {
	cursorLine := s.activeBuffer().Cursor().Line
	s.viewportTopLine = cursorLine
	maxTopLine := s.activeBuffer().LineCount() - 1
	if s.viewportTopLine > maxTopLine {
		s.viewportTopLine = maxTopLine
	}
	s.status = "Cursor at top"
}

// scrollToBottom positions cursor line at bottom of viewport (Vim's zb command).
func (s *appState) scrollToBottom(linesPerPage int) {
	if linesPerPage <= 0 {
		linesPerPage = 20
	}
	cursorLine := s.activeBuffer().Cursor().Line
	s.viewportTopLine = cursorLine - linesPerPage + 1
	if s.viewportTopLine < 0 {
		s.viewportTopLine = 0
	}
	s.status = "Cursor at bottom"
}

// scrollLineUp scrolls viewport up by one line (Vim's Ctrl+Y).
func (s *appState) scrollLineUp() {
	if s.viewportTopLine > 0 {
		s.viewportTopLine--
		s.status = fmt.Sprintf("Scrolled up (top line: %d)", s.viewportTopLine+1)
	}
}

// scrollLineDown scrolls viewport down by one line (Vim's Ctrl+E).
func (s *appState) scrollLineDown() {
	maxTopLine := s.activeBuffer().LineCount() - 1
	if s.viewportTopLine < maxTopLine {
		s.viewportTopLine++
		s.status = fmt.Sprintf("Scrolled down (top line: %d)", s.viewportTopLine+1)
	}
}

func (s *appState) handleCountDigit(d int) bool {
	if s.mode == modeInsert {
		return false
	}
	if d == 0 && s.pendingCount == 0 {
		return false
	}
	s.pendingCount = s.pendingCount*10 + d
	s.status = fmt.Sprintf("Count %d", s.pendingCount)
	return true
}

func (s *appState) consumeCount(defaultVal int) int {
	if s.pendingCount > 0 {
		val := s.pendingCount
		s.pendingCount = 0
		return val
	}
	return defaultVal
}

func (s *appState) resetCount() {
	s.pendingCount = 0
	s.pendingGoto = false
	s.pendingScroll = false
}

func (s *appState) gotoLine(target int) {
	total := s.activeBuffer().LineCount()
	if total == 0 {
		return
	}
	if target <= 0 {
		target = total
	}
	if target > total {
		target = total
	}
	s.activeBuffer().MoveToLine(target - 1)
	s.setCursorStatus(fmt.Sprintf("Goto line %d", target))
}

func (s *appState) gotoLineWithCount() {
	target := s.consumeCount(-1)
	if target <= 0 {
		target = s.activeBuffer().LineCount()
	}
	s.gotoLine(target)
}

func (s *appState) enterDeleteMode() {
	s.mode = modeDelete
	s.pendingCount = 0
	s.pendingGoto = false
	s.status = "DELETE mode: type line # and press d"
}

func (s *appState) exitDeleteMode() {
	if s.mode == modeDelete {
		s.mode = modeNormal
	}
	s.resetCount()
	s.status = "Back to NORMAL"
}

func (s *appState) executeDeleteCommand() {
	target := s.pendingCount
	if target <= 0 {
		target = s.activeBuffer().Cursor().Line + 1
	}
	total := s.activeBuffer().LineCount()
	if target < 1 || target > total {
		s.status = fmt.Sprintf("Line %d out of range", target)
		s.exitDeleteMode()
		return
	}
	s.activeBuffer().DeleteLines(target-1, target-1)
	s.setCursorStatus(fmt.Sprintf("Deleted line %d", target))
	s.exitDeleteMode()
}

func (s *appState) startGotoSequence() {
	s.pendingGoto = true
	s.status = "goto line: awaiting g/G"
}

func (s *appState) handleGotoSequence(r rune) bool {
	if !s.pendingGoto {
		return false
	}
	s.pendingGoto = false
	switch r {
	case 'g':
		target := s.consumeCount(1)
		s.gotoLine(target)
		return true
	case 'G':
		target := s.consumeCount(s.activeBuffer().LineCount())
		if target <= 0 {
			target = s.activeBuffer().LineCount()
		}
		s.gotoLine(target)
		return true
	default:
		return false
	}
}

func (s *appState) startScrollSequence() {
	s.pendingScroll = true
	s.status = "scroll: awaiting z/t/b"
}

func (s *appState) handleScrollSequence(r rune) bool {
	if !s.pendingScroll {
		return false
	}
	s.pendingScroll = false

	// Calculate approximate lines per page
	linesPerPage := 20

	switch r {
	case 'z':
		s.scrollToCenter(linesPerPage)
		return true
	case 't':
		s.scrollToTop()
		return true
	case 'b':
		s.scrollToBottom(linesPerPage)
		return true
	default:
		s.status = "Unknown scroll command"
		return false
	}
}

func (s *appState) enterVisualChar() {
	s.mode = modeVisual
	s.visualMode = visualModeChar
	s.visualStartLine = s.activeBuffer().Cursor().Line
	s.visualStartCol = s.activeBuffer().Cursor().Col
	s.resetCount()
	s.status = "VISUAL (char)"
}

func (s *appState) enterVisualLine() {
	s.mode = modeVisual
	s.visualMode = visualModeLine
	s.visualStartLine = s.activeBuffer().Cursor().Line
	s.visualStartCol = 0
	s.resetCount()
	s.status = "VISUAL (line)"
}

func (s *appState) enterCommandMode() {
	if s.mode == modeCommand {
		return
	}
	s.mode = modeCommand
	s.cmdText = ""
	s.status = "COMMAND (:...)"
}

func (s *appState) exitVisualMode() {
	if s.mode == modeVisual {
		s.mode = modeNormal
	}
	s.visualMode = visualModeNone
	s.visualStartLine = 0
	s.visualStartCol = 0
}

func (s *appState) exitCommandMode() {
	if s.mode == modeCommand {
		s.mode = modeNormal
	}
	s.cmdText = ""
}

func (s *appState) enterExplorerMode() {
	if s.mode == modeExplorer {
		return
	}
	s.mode = modeExplorer
	s.explorerVisible = true
	s.explorerFocused = true
	s.status = "EXPLORER (j/k nav, Enter open, u up-dir, r refresh, Esc exit, Ctrl+T toggle)"
}

func (s *appState) exitExplorerMode() {
	s.mode = modeNormal
	s.explorerFocused = false
	s.status = "Back to NORMAL"
}

func (s *appState) toggleExplorer() {
	if s.explorerVisible {
		// Hide explorer completely
		s.explorerVisible = false
		s.explorerFocused = false
		if s.mode == modeExplorer {
			s.mode = modeNormal
		}
		s.status = "Explorer closed"
	} else {
		// Show explorer AND focus it immediately
		s.enterExplorerMode()
	}
}

// ensureExplorerItemVisible scrolls the explorer to keep selected item visible
func (s *appState) ensureExplorerItemVisible() {
	if s.fileTree == nil {
		return
	}

	selectedIndex := s.fileTree.SelectedIndex()
	nodes := s.fileTree.GetFlatList()

	if selectedIndex < 0 || selectedIndex >= len(nodes) {
		return
	}

	// Scroll to show selected item
	s.explorerListPosition.Position.First = selectedIndex
	s.explorerListPosition.Position.Offset = 0
}

func (s *appState) toggleFullscreen() {
	if s.window == nil {
		return
	}

	if s.currentWindowMode == app.Fullscreen {
		// Exit fullscreen - return to windowed mode
		s.window.Option(app.Windowed.Option())
		s.status = "Exited fullscreen"
		s.wasFullscreen = false
	} else {
		// Enter fullscreen
		s.window.Perform(system.ActionFullscreen)
		s.status = "Entered fullscreen (Shift+Enter to exit)"
		s.wasFullscreen = true
	}
}

func (s *appState) openSelectedNode() {
	if s.fileTree == nil {
		return
	}

	node := s.fileTree.SelectedNode()
	if node == nil {
		return
	}

	if node.IsDir {
		// Special handling for ".." parent directory
		if node.Name == ".." {
			if err := s.fileTree.ChangeRoot(node.Path); err != nil {
				s.status = fmt.Sprintf("Error navigating to parent: %v", err)
			} else {
				s.fileTree.LoadInitial()
				s.status = fmt.Sprintf("Changed to %s", s.fileTree.CurrentPath())
			}
			return
		}

		// Toggle directory expansion
		if node.Expanded {
			s.fileTree.Collapse()
		} else {
			s.fileTree.Expand()
			s.fileTree.ExpandAndLoad(node)
		}
		return
	}

	// Open file
	_, err := s.bufferMgr.OpenFile(node.Path)
	if err != nil {
		s.status = fmt.Sprintf("Error opening %s: %v", node.Name, err)
		return
	}

	// Update the active pane to display the newly opened buffer
	if s.paneManager != nil {
		activePane := s.paneManager.ActivePane()
		if activePane != nil {
			activePane.SetBufferIndex(s.bufferMgr.ActiveIndex())
		}
	}

	s.exitExplorerMode()
	s.status = fmt.Sprintf("Opened %s", node.Name)
}

func (s *appState) visualSelectionRange() (int, int, bool) {
	if s.visualMode == visualModeNone {
		return 0, 0, false
	}
	cur := s.activeBuffer().Cursor().Line
	start := s.visualStartLine
	if cur < start {
		return cur, start, true
	}
	return start, cur, true
}

func (s *appState) visualSelectionRangeChar() (startLine, startCol, endLine, endCol int, ok bool) {
	if s.visualMode != visualModeChar {
		return 0, 0, 0, 0, false
	}

	curLine := s.activeBuffer().Cursor().Line
	curCol := s.activeBuffer().Cursor().Col

	// Determine start and end based on cursor position relative to anchor
	if curLine < s.visualStartLine || (curLine == s.visualStartLine && curCol < s.visualStartCol) {
		// Selection goes backward
		return curLine, curCol, s.visualStartLine, s.visualStartCol, true
	}
	// Selection goes forward
	return s.visualStartLine, s.visualStartCol, curLine, curCol, true
}

func (s *appState) deleteVisualSelection() {
	if s.visualMode == visualModeChar {
		// Character-wise deletion
		startLine, startCol, endLine, endCol, ok := s.visualSelectionRangeChar()
		if !ok {
			s.status = "No selection"
			return
		}
		s.activeBuffer().DeleteCharRange(startLine, startCol, endLine, endCol)
		s.exitVisualMode()
		s.setCursorStatus("Deleted selection")
	} else if s.visualMode == visualModeLine {
		// Line-wise deletion
		start, end, ok := s.visualSelectionRange()
		if !ok {
			s.status = "No selection"
			return
		}
		s.activeBuffer().DeleteLines(start, end)
		s.exitVisualMode()
		s.setCursorStatus("Deleted selection")
	} else {
		s.status = "No selection"
	}
}

func (s *appState) copyVisualSelection() {
	if s.visualMode == visualModeChar {
		// Character-wise copy
		startLine, startCol, endLine, endCol, ok := s.visualSelectionRangeChar()
		if !ok {
			s.status = "No selection to copy"
			return
		}
		text := s.activeBuffer().GetCharRange(startLine, startCol, endLine, endCol)
		if len(text) == 0 {
			s.status = "No selection to copy"
			return
		}
		// Store as a single line in clipboard
		s.clipLines = []string{text}
		s.status = fmt.Sprintf("Copied %d character(s)", len(text))
	} else if s.visualMode == visualModeLine {
		// Line-wise copy
		start, end, ok := s.visualSelectionRange()
		if !ok {
			s.status = "No selection to copy"
			return
		}
		lines := s.activeBuffer().LinesRange(start, end)
		if len(lines) == 0 {
			s.status = "No selection to copy"
			return
		}
		s.clipLines = append([]string(nil), lines...)
		s.status = fmt.Sprintf("Copied %d line(s)", len(lines))
	} else {
		s.status = "No selection to copy"
	}
}

func (s *appState) pasteClipboard() {
	if len(s.clipLines) == 0 {
		s.status = "Clipboard empty"
		return
	}

	if s.visualMode == visualModeChar {
		// Character-wise paste: replace selection with clipboard text
		startLine, startCol, endLine, endCol, ok := s.visualSelectionRangeChar()
		if !ok {
			s.status = "Select destination in VISUAL mode"
			return
		}
		buf := s.activeBuffer()
		// Delete the selected range (this positions cursor at startLine, startCol)
		buf.DeleteCharRange(startLine, startCol, endLine, endCol)
		// Insert clipboard text at cursor position
		text := s.clipLines[0] // Character copy stores as single line
		buf.InsertText(text)
		s.exitVisualMode()
		s.setCursorStatus(fmt.Sprintf("Pasted %d character(s)", len(text)))
	} else if s.visualMode == visualModeLine {
		// Line-wise paste
		start, _, ok := s.visualSelectionRange()
		if !ok {
			s.status = "Select destination in VISUAL mode"
			return
		}
		lines := append([]string(nil), s.clipLines...)
		s.activeBuffer().InsertLines(start, lines)
		s.exitVisualMode()
		s.setCursorStatus(fmt.Sprintf("Inserted %d line(s)", len(lines)))
	} else {
		s.status = "Select destination in VISUAL mode"
	}
}

func (s *appState) isColonKey(ev key.Event) bool {
	if string(ev.Name) == ":" {
		return true
	}
	if ev.Modifiers.Contain(key.ModShift) && string(ev.Name) == ";" {
		return true
	}
	return false
}

func (s *appState) appendCommandText(text string) {
	if text == "" {
		return
	}
	for _, r := range text {
		if r == '\n' || r == '\r' {
			continue
		}
		s.cmdText += string(r)
	}
}

func (s *appState) deleteCommandChar() {
	if s.cmdText == "" {
		return
	}
	runes := []rune(s.cmdText)
	if len(runes) == 0 {
		return
	}
	s.cmdText = string(runes[:len(runes)-1])
}

func (s *appState) executeCommandLine() {
	cmd := strings.TrimSpace(s.cmdText)
	s.exitCommandMode()
	if cmd == "" {
		s.status = "No command"
		return
	}
	if strings.HasPrefix(cmd, ":") {
		cmd = strings.TrimSpace(cmd[1:])
	}
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		s.status = "No command"
		return
	}
	name := strings.ToLower(fields[0])
	args := ""
	if len(fields) > 1 {
		args = strings.Join(fields[1:], " ")
	}
	switch name {
	case "q", "quit":
		s.handleQuitCommand(false)
	case "q!":
		s.handleQuitCommand(true)
	case "w", "write":
		s.handleWriteCommand(strings.TrimSpace(args), false)
	case "wq":
		s.handleWriteCommand(strings.TrimSpace(args), true)
	case "e", "edit":
		s.handleEditCommand(strings.TrimSpace(args))
	case "bn", "bnext":
		if s.bufferMgr.NextBuffer() {
			s.status = "Switched to next buffer"
		} else {
			s.status = "Already at last buffer"
		}
	case "bp", "bprev":
		if s.bufferMgr.PrevBuffer() {
			s.status = "Switched to previous buffer"
		} else {
			s.status = "Already at first buffer"
		}
	case "bd", "bdelete":
		s.handleBufferDeleteCommand(false)
	case "bd!":
		s.handleBufferDeleteCommand(true)
	case "ls", "buffers":
		s.handleListBuffersCommand()
	case "ex", "explore":
		s.toggleExplorer()
	case "cd":
		s.handleChangeDirectoryCommand(strings.TrimSpace(args))
	case "pwd":
		s.handlePrintWorkingDirectoryCommand()
	default:
		s.status = fmt.Sprintf("Unknown command: %s", name)
	}
}

func (s *appState) handleQuitCommand(force bool) {
	// With pane system, quit closes the active pane and its buffer
	if s.paneManager == nil {
		s.requestClose()
		return
	}

	activePane := s.paneManager.ActivePane()
	if activePane == nil {
		s.requestClose()
		return
	}

	// Get the buffer for this pane
	buf := s.bufferMgr.GetBuffer(activePane.BufferIndex)
	if buf == nil {
		// No buffer, just close the pane
		if s.paneManager.PaneCount() > 1 {
			s.paneManager.ClosePane()
			s.status = fmt.Sprintf("Pane closed - %d panes remaining", s.paneManager.PaneCount())
		} else {
			s.requestClose()
		}
		return
	}

	// Check if buffer is modified (unless force)
	if !force && buf.Modified() {
		s.status = "Buffer has unsaved changes (use :q! to force)"
		return
	}

	// Close the buffer
	if err := s.bufferMgr.CloseBuffer(activePane.BufferIndex, force); err != nil {
		s.status = fmt.Sprintf("Error closing buffer: %v", err)
		return
	}

	// Close the pane
	if s.paneManager.PaneCount() > 1 {
		if err := s.paneManager.ClosePane(); err != nil {
			s.status = fmt.Sprintf("Error closing pane: %v", err)
		} else {
			s.status = fmt.Sprintf("Pane closed - %d panes remaining", s.paneManager.PaneCount())
		}
	} else {
		// Last pane - close the application
		s.requestClose()
	}
}

func (s *appState) handleWriteCommand(arg string, andQuit bool) {
	// Check if we have an active buffer
	buf := s.activeBuffer()
	if buf == nil {
		s.status = "No active buffer to save"
		return
	}

	var err error
	if arg == "" {
		// Save to current file
		err = s.bufferMgr.SaveActiveBuffer()
	} else {
		// Save as
		err = s.bufferMgr.SaveAs(arg)
	}

	if err != nil {
		s.status = fmt.Sprintf("Write failed: %v", err)
		return
	}

	filename := buf.FilePath()
	if andQuit {
		s.handleQuitCommand(false)
		s.status = fmt.Sprintf("Wrote + closed %s", filename)
		return
	}
	s.status = fmt.Sprintf("Wrote %d line(s) → %s", buf.LineCount(), filename)
}

func (s *appState) handleEditCommand(path string) {
	if path == "" {
		s.status = "E471: Argument required"
		return
	}

	_, err := s.bufferMgr.OpenFile(path)
	if err != nil {
		s.status = fmt.Sprintf("Error opening %s: %v", path, err)
		return
	}

	// Update the active pane to display the newly opened buffer
	if s.paneManager != nil {
		activePane := s.paneManager.ActivePane()
		if activePane != nil {
			activePane.SetBufferIndex(s.bufferMgr.ActiveIndex())
		}
	}

	s.status = fmt.Sprintf("Opened %s", path)
}

func (s *appState) handleBufferDeleteCommand(force bool) {
	if err := s.bufferMgr.CloseActiveBuffer(force); err != nil {
		s.status = fmt.Sprintf("Error: %v", err)
	} else {
		s.status = "Buffer deleted"
	}
}

func (s *appState) handleListBuffersCommand() {
	buffers := s.bufferMgr.ListBuffers()
	s.status = fmt.Sprintf("Buffers: %s", strings.Join(buffers, " | "))
}

func (s *appState) handleChangeDirectoryCommand(path string) {
	if path == "" {
		// No argument - go to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			s.status = fmt.Sprintf("Error getting home directory: %v", err)
			return
		}
		path = homeDir
	} else if path == "~" {
		// Expand ~ to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			s.status = fmt.Sprintf("Error getting home directory: %v", err)
			return
		}
		path = homeDir
	} else if strings.HasPrefix(path, "~/") {
		// Expand ~/path to home/path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			s.status = fmt.Sprintf("Error getting home directory: %v", err)
			return
		}
		path = filepath.Join(homeDir, path[2:])
	}

	if s.fileTree == nil {
		s.status = "File tree not initialized"
		return
	}

	if err := s.fileTree.ChangeRoot(path); err != nil {
		s.status = fmt.Sprintf("Error changing directory: %v", err)
		return
	}

	if err := s.fileTree.LoadInitial(); err != nil {
		s.status = fmt.Sprintf("Error loading directory: %v", err)
		return
	}

	s.status = fmt.Sprintf("Changed directory to %s", s.fileTree.CurrentPath())

	// Show explorer if not already visible
	if !s.explorerVisible {
		s.explorerVisible = true
	}
}

func (s *appState) handlePrintWorkingDirectoryCommand() {
	if s.fileTree == nil {
		s.status = "File tree not initialized"
		return
	}
	s.status = fmt.Sprintf("Current directory: %s", s.fileTree.CurrentPath())
}

func (s *appState) enterRenameMode() {
	if s.fileTree == nil {
		return
	}
	node := s.fileTree.SelectedNode()
	if node == nil {
		return
	}
	s.fileOpMode = "rename"
	s.fileOpOriginalName = node.Name
	s.fileOpInput = ""
	s.fileOpTarget = node
	s.skipNextFileOpEdit = true
	s.status = fmt.Sprintf("Rename File: %s -> ", s.fileOpOriginalName)
}

func (s *appState) enterCreateMode() {
	if s.fileTree == nil {
		return
	}
	node := s.fileTree.SelectedNode()
	if node == nil {
		return
	}
	s.fileOpMode = "create"
	s.fileOpInput = ""
	s.fileOpTarget = node
	s.skipNextFileOpEdit = true
	s.status = "New file: "
}

func (s *appState) enterFileDeleteMode() {
	if s.fileTree == nil {
		return
	}
	node := s.fileTree.SelectedNode()
	if node == nil {
		return
	}
	s.fileOpMode = "delete"
	s.fileOpTarget = node
	s.status = fmt.Sprintf("Delete '%s'? (y/n)", node.Name)
}

func (s *appState) completeFileOp() {
	if s.fileTree == nil || s.fileOpTarget == nil {
		s.cancelFileOp()
		return
	}

	var err error
	switch s.fileOpMode {
	case "rename":
		if s.fileOpInput == "" {
			s.status = "Error: filename cannot be empty"
			s.cancelFileOp()
			return
		}
		err = s.fileTree.RenameNode(s.fileOpTarget, s.fileOpInput)
		if err != nil {
			s.status = fmt.Sprintf("Rename failed: %v", err)
		} else {
			s.status = fmt.Sprintf("Renamed to '%s'", s.fileOpInput)
			s.fileTree.Refresh()
		}

	case "create":
		if s.fileOpInput == "" {
			s.status = "Error: filename cannot be empty"
			s.cancelFileOp()
			return
		}
		err = s.fileTree.CreateFile(s.fileOpTarget, s.fileOpInput)
		if err != nil {
			s.status = fmt.Sprintf("Create failed: %v", err)
		} else {
			// Check if directories were created
			if strings.Contains(s.fileOpInput, "/") || strings.Contains(s.fileOpInput, string(filepath.Separator)) {
				s.status = fmt.Sprintf("Created path '%s'", s.fileOpInput)
			} else {
				s.status = fmt.Sprintf("Created '%s'", s.fileOpInput)
			}
			s.fileTree.Refresh()
		}

	case "delete":
		err = s.fileTree.DeleteNode(s.fileOpTarget)
		if err != nil {
			s.status = fmt.Sprintf("Delete failed: %v", err)
		} else {
			s.status = fmt.Sprintf("Deleted '%s'", s.fileOpTarget.Name)
			s.fileTree.Refresh()
		}
	}

	s.cancelFileOp()
}

func (s *appState) cancelFileOp() {
	s.fileOpMode = ""
	s.fileOpInput = ""
	s.fileOpOriginalName = ""
	s.fileOpTarget = nil
	s.skipNextFileOpEdit = false
}

func (s *appState) appendFileOpInput(text string) {
	if text == "" {
		return
	}
	for _, r := range text {
		if r == '\n' || r == '\r' {
			continue
		}
		s.fileOpInput += string(r)
	}
	s.status = s.getFileOpPrompt()
}

func (s *appState) deleteFileOpChar() {
	if s.fileOpInput == "" {
		return
	}
	runes := []rune(s.fileOpInput)
	if len(runes) == 0 {
		return
	}
	s.fileOpInput = string(runes[:len(runes)-1])
	s.status = s.getFileOpPrompt()
}

func (s *appState) getFileOpPrompt() string {
	switch s.fileOpMode {
	case "rename":
		return fmt.Sprintf("Rename File: %s -> %s", s.fileOpOriginalName, s.fileOpInput)
	case "create":
		return fmt.Sprintf("New file: %s", s.fileOpInput)
	case "delete":
		if s.fileOpTarget != nil {
			return fmt.Sprintf("Delete '%s'? (y/n)", s.fileOpTarget.Name)
		}
		return "Delete? (y/n)"
	default:
		return s.status
	}
}

func (s *appState) handleFileOpKey(ev key.Event) bool {
	// Handle Escape to cancel
	if ev.Name == key.NameEscape {
		s.status = "Cancelled"
		s.cancelFileOp()
		return true
	}

	// Handle delete confirmation (y/n)
	if s.fileOpMode == "delete" {
		if r, ok := s.printableKey(ev); ok {
			switch unicode.ToLower(r) {
			case 'y':
				s.completeFileOp()
				return true
			case 'n':
				s.cancelFileOp()
				return true
			}
		}
		return true
	}

	// Handle Enter to complete rename/create
	if ev.Name == key.NameReturn || ev.Name == key.NameEnter {
		s.completeFileOp()
		return true
	}

	// Handle Backspace for rename/create
	if ev.Name == key.NameDeleteBackward {
		s.deleteFileOpChar()
		return true
	}

	// Consume ALL other keys when in rename/create mode
	// This prevents 'd', 'r', 'n' from triggering explorer actions while typing
	// Text input is handled separately via EditEvent
	if s.fileOpMode == "rename" || s.fileOpMode == "create" {
		return true
	}

	return false
}

func (s *appState) insertText(text string) {
	if text == "" {
		return
	}
	buf := s.activeBuffer()
	buf.InsertText(text)
	s.setCursorStatus(fmt.Sprintf("Insert %q", text))
}

func (s *appState) saveBufferToFile(path string) error {
	if path == "" {
		return fmt.Errorf("no file name")
	}
	lines := make([]string, s.activeBuffer().LineCount())
	for i := 0; i < len(lines); i++ {
		lines[i] = s.activeBuffer().Line(i)
	}
	content := strings.Join(lines, "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func (s *appState) requestClose() {
	if s.window == nil {
		return
	}
	s.window.Perform(system.ActionClose)
}

func (s *appState) printableKey(ev key.Event) (rune, bool) {
	str := string(ev.Name)
	if utf8.RuneCountInString(str) != 1 {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(str)
	if r == utf8.RuneError || r == 0 {
		return 0, false
	}
	if unicode.IsControl(r) {
		return 0, false
	}

	// Handle letter case based on shift state
	// Gio reports all letters as uppercase in ev.Name, so we need to check shift state
	if unicode.IsLetter(r) {
		// If shift is NOT pressed, convert to lowercase
		if !s.shiftPressed {
			r = unicode.ToLower(r)
		}
		// If shift IS pressed, keep uppercase (already uppercase from Gio)
	}

	// Gio reports punctuation keys like Shift+; as ';' with ModShift; normalize to ':'.
	if r == ';' && ev.Modifiers.Contain(key.ModShift) {
		r = ':'
	}
	return r, true
}

func describeKey(ev key.Event) string {
	if ev.Name != "" {
		return string(ev.Name)
	}
	return "key"
}

// isPrintableKey returns true if the key generates a printable character via EditEvent.
// These keys should not have their modifiers reset immediately in INSERT mode,
// as the EditEvent needs the modifier state to determine capitalization.
func isPrintableKey(keyName key.Name) bool {
	nameStr := string(keyName)

	// Single character keys (letters and digits)
	if len(nameStr) == 1 {
		r := rune(nameStr[0])
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
		// Also include common printable symbols
		if strings.ContainsRune("!@#$%^&*()_+-=[]{}\\|;:'\",.<>/?`~", r) {
			return true
		}
	}

	// Special printable keys
	return keyName == key.NameSpace
}

// expandTabs converts tab characters to spaces.
// tabWidth specifies how many spaces each tab should expand to.
func expandTabs(s string, tabWidth int) string {
	if !strings.Contains(s, "\t") {
		return s
	}
	var result strings.Builder
	col := 0
	for _, r := range s {
		if r == '\t' {
			spaces := tabWidth - (col % tabWidth)
			result.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			result.WriteRune(r)
			col++
		}
	}
	return result.String()
}

// Search mode methods

func (s *appState) enterSearchMode() {
	s.mode = modeSearch
	s.searchPattern = ""
	s.searchMatches = nil
	s.currentMatchIdx = -1
	s.skipNextSearchEdit = true
	s.status = "/"
}

func (s *appState) exitSearchMode() {
	s.mode = modeNormal
	s.status = "Search cancelled"
}

func (s *appState) executeSearch() {
	if s.searchPattern == "" {
		s.searchMatches = nil
		s.currentMatchIdx = -1
		s.searchActive = false
		s.status = "No search pattern"
		s.mode = modeNormal
		return
	}

	s.searchMatches = s.findAllMatches(s.searchPattern)

	if len(s.searchMatches) == 0 {
		s.status = fmt.Sprintf("Pattern not found: %s", s.searchPattern)
		s.searchActive = false
		s.currentMatchIdx = -1
	} else {
		s.currentMatchIdx = s.findNextMatchFromCursor()
		if s.currentMatchIdx >= 0 {
			match := s.searchMatches[s.currentMatchIdx]
			s.activeBuffer().MoveToLine(match.Line)
			// Move to start of line then right to column
			s.activeBuffer().JumpLineStart()
			for i := 0; i < match.Col; i++ {
				s.activeBuffer().MoveRight()
			}
		}
		s.searchActive = true
		s.status = fmt.Sprintf("/%s [%d/%d]", s.searchPattern, s.currentMatchIdx+1, len(s.searchMatches))
	}

	s.mode = modeNormal
}

func (s *appState) findAllMatches(pattern string) []SearchMatch {
	var matches []SearchMatch

	for lineIdx := 0; lineIdx < s.activeBuffer().LineCount(); lineIdx++ {
		line := s.activeBuffer().Line(lineIdx)
		lowerLine := strings.ToLower(line)
		lowerPattern := strings.ToLower(pattern)

		startPos := 0
		for {
			idx := strings.Index(lowerLine[startPos:], lowerPattern)
			if idx == -1 {
				break
			}

			actualPos := startPos + idx
			matches = append(matches, SearchMatch{
				Line: lineIdx,
				Col:  len([]rune(line[:actualPos])),
				Len:  len([]rune(pattern)),
			})

			startPos = actualPos + 1
		}
	}

	return matches
}

func (s *appState) findNextMatchFromCursor() int {
	if len(s.searchMatches) == 0 {
		return -1
	}

	curLine := s.activeBuffer().Cursor().Line
	curCol := s.activeBuffer().Cursor().Col

	for i, match := range s.searchMatches {
		if match.Line > curLine || (match.Line == curLine && match.Col >= curCol) {
			return i
		}
	}

	return 0
}

func (s *appState) jumpToNextMatch() {
	if !s.searchActive || len(s.searchMatches) == 0 {
		s.status = "No active search"
		return
	}

	s.currentMatchIdx = (s.currentMatchIdx + 1) % len(s.searchMatches)
	match := s.searchMatches[s.currentMatchIdx]

	s.activeBuffer().MoveToLine(match.Line)
	s.activeBuffer().JumpLineStart()
	for i := 0; i < match.Col; i++ {
		s.activeBuffer().MoveRight()
	}
	s.status = fmt.Sprintf("/%s [%d/%d]", s.searchPattern, s.currentMatchIdx+1, len(s.searchMatches))
}

func (s *appState) jumpToPrevMatch() {
	if !s.searchActive || len(s.searchMatches) == 0 {
		s.status = "No active search"
		return
	}

	s.currentMatchIdx--
	if s.currentMatchIdx < 0 {
		s.currentMatchIdx = len(s.searchMatches) - 1
	}
	match := s.searchMatches[s.currentMatchIdx]

	s.activeBuffer().MoveToLine(match.Line)
	s.activeBuffer().JumpLineStart()
	for i := 0; i < match.Col; i++ {
		s.activeBuffer().MoveRight()
	}
	s.status = fmt.Sprintf("/%s [%d/%d]", s.searchPattern, s.currentMatchIdx+1, len(s.searchMatches))
}

func (s *appState) clearSearch() {
	s.searchActive = false
	s.searchMatches = nil
	s.currentMatchIdx = -1
	s.searchPattern = ""
	s.status = "Search cleared"
}

func (s *appState) appendSearchText(text string) {
	if text == "" {
		return
	}
	for _, r := range text {
		if r == '\n' || r == '\r' {
			continue
		}
		s.searchPattern += string(r)
	}
	s.status = "/" + s.searchPattern
}

func (s *appState) deleteSearchChar() {
	if s.searchPattern == "" {
		return
	}
	runes := []rune(s.searchPattern)
	if len(runes) == 0 {
		return
	}
	s.searchPattern = string(runes[:len(runes)-1])
	s.status = "/" + s.searchPattern
}

// Fuzzy finder methods

func (s *appState) enterFuzzyFinder() {
	if s.fileTree == nil {
		s.status = "File tree not available"
		return
	}

	// Discover all files in the workspace
	workDir := s.fileTree.CurrentPath()
	files, err := filesystem.FindAllFiles(workDir)
	if err != nil {
		s.status = fmt.Sprintf("Error discovering files: %v", err)
		return
	}

	s.mode = modeFuzzyFinder
	s.fuzzyFinderActive = true
	s.fuzzyFinderInput = ""
	s.fuzzyFinderFiles = files
	s.fuzzyFinderMatches = PerformFuzzyMatch("", files, 50)
	s.fuzzyFinderSelectedIdx = 0
	s.skipNextFuzzyEdit = true
	s.status = fmt.Sprintf("Fuzzy Finder: %d files", len(files))
}

func (s *appState) exitFuzzyFinder() {
	s.mode = modeNormal
	s.fuzzyFinderActive = false
	s.fuzzyFinderInput = ""
	s.fuzzyFinderFiles = nil
	s.fuzzyFinderMatches = nil
	s.fuzzyFinderSelectedIdx = 0
	s.status = "Fuzzy finder cancelled"
}

func (s *appState) updateFuzzyMatches() {
	s.fuzzyFinderMatches = PerformFuzzyMatch(s.fuzzyFinderInput, s.fuzzyFinderFiles, 50)
	s.fuzzyFinderSelectedIdx = 0
}

func (s *appState) appendFuzzyInput(text string) {
	if text == "" {
		return
	}
	for _, r := range text {
		if r == '\n' || r == '\r' {
			continue
		}
		s.fuzzyFinderInput += string(r)
	}
	s.updateFuzzyMatches()
}

func (s *appState) deleteFuzzyChar() {
	if s.fuzzyFinderInput == "" {
		return
	}
	runes := []rune(s.fuzzyFinderInput)
	if len(runes) == 0 {
		return
	}
	s.fuzzyFinderInput = string(runes[:len(runes)-1])
	s.updateFuzzyMatches()
}

func (s *appState) fuzzyFinderMoveUp() {
	if s.fuzzyFinderSelectedIdx > 0 {
		s.fuzzyFinderSelectedIdx--
	}
}

func (s *appState) fuzzyFinderMoveDown() {
	if s.fuzzyFinderSelectedIdx < len(s.fuzzyFinderMatches)-1 {
		s.fuzzyFinderSelectedIdx++
	}
}

func (s *appState) fuzzyFinderConfirm() {
	if s.fuzzyFinderSelectedIdx < 0 || s.fuzzyFinderSelectedIdx >= len(s.fuzzyFinderMatches) {
		s.exitFuzzyFinder()
		return
	}

	match := s.fuzzyFinderMatches[s.fuzzyFinderSelectedIdx]
	fullPath := filepath.Join(s.fileTree.CurrentPath(), match.FilePath)

	_, err := s.bufferMgr.OpenFile(fullPath)
	if err != nil {
		s.status = fmt.Sprintf("Error opening %s: %v", match.FilePath, err)
		s.exitFuzzyFinder()
		return
	}

	// Update the active pane to display the newly opened buffer
	if s.paneManager != nil {
		activePane := s.paneManager.ActivePane()
		if activePane != nil {
			activePane.SetBufferIndex(s.bufferMgr.ActiveIndex())
		}
	}

	s.exitFuzzyFinder()
	s.status = fmt.Sprintf("Opened %s", match.FilePath)
}

const sampleBuffer = ``
