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
)

type appState struct {
	theme              *material.Theme
	bufferMgr          *editor.BufferManager
	fileTree           *filesystem.FileTree
	mode               mode
	status             string
	lastKey            string
	focusTag           *int
	pendingCount       int
	pendingGoto        bool
	pendingScroll      bool
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
	explorerVisible bool
	explorerWidth   int
	explorerFocused bool

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
	theme.Shaper = text.NewShaper(
		text.NoSystemFonts(),
		text.WithCollection(gofont.Collection()),
	)

	buf := editor.NewBuffer(strings.TrimSpace(sampleBuffer))
	bufferMgr := editor.NewBufferManagerWithBuffer(buf)

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
		theme:             theme,
		bufferMgr:         bufferMgr,
		fileTree:          fileTree,
		mode:              modeNormal,
		status:            "Ready",
		focusTag:          new(int),
		visualMode:        visualModeNone,
		visualStartLine:   0,
		visualStartCol:    0,
		caretVisible:      true,
		explorerVisible:   false,
		explorerWidth:     250,
		explorerFocused:   false,
		currentWindowMode: app.Windowed,
		wasFullscreen:     false,
		viewportTopLine:   0,
		scrollOffsetLines: 3,
		listPosition:      layout.List{Axis: layout.Vertical},
	}
}

// activeBuffer returns the active buffer (helper method).
func (s *appState) activeBuffer() *editor.Buffer {
	if s.bufferMgr == nil {
		return nil
	}
	return s.bufferMgr.ActiveBuffer()
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
				// Horizontal split: explorer | editor
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return s.drawFileExplorer(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return s.drawBuffer(gtx)
					}),
				)
			}
			return s.drawBuffer(gtx)
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
				wasPressed := s.ctrlPressed
				s.ctrlPressed = (e.State == key.Press)
				log.Printf("[CTRL_TRACK] Ctrl %v -> %v (state=%v)", wasPressed, s.ctrlPressed, e.State)
				continue
			}
			// Track Shift key press/release
			if e.Name == key.NameShift {
				wasPressed := s.shiftPressed
				s.shiftPressed = (e.State == key.Press)
				log.Printf("[SHIFT_TRACK] Shift %v -> %v (state=%v)", wasPressed, s.shiftPressed, e.State)
				continue
			}
			// Track Alt key press/release
			if e.Name == key.NameAlt {
				log.Printf("[ALT_TRACK] Alt key event received, state=%v", e.State)
				continue
			}

			// Platform quirk: ev.Modifiers is ALWAYS empty on this platform!
			// We must rely ONLY on our tracked modifier state (ctrlPressed, shiftPressed)
			log.Printf("[PLATFORM_QUIRK] ev.Modifiers=%v (platform never sets this)", e.Modifiers)

			// Save modifier state before handling key
			hadCtrl := s.ctrlPressed
			hadShift := s.shiftPressed

			s.handleKey(e)

			// Smart reset: If modifiers are still set after handleKey and we're in certain modes,
			// it likely means the user released the modifier but platform didn't send Release event.
			// Reset modifiers after successful command execution to prevent them from sticking.
			// Exception: Don't reset if we just entered a mode (mode transitions are intentional)
			if s.mode == modeNormal || s.mode == modeInsert || s.mode == modeCommand {
				if hadCtrl && s.ctrlPressed {
					// Ctrl was held before and is still held - user might have released it
					// Reset it so next key doesn't have Ctrl stuck
					log.Printf("[SMART_RESET] Resetting Ctrl after key=%q to prevent sticking", e.Name)
					s.ctrlPressed = false
				}
				if hadShift && s.shiftPressed {
					log.Printf("[SMART_RESET] Resetting Shift after key=%q to prevent sticking", e.Name)
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
					continue
				}
				s.insertText(e.Text)
			case modeCommand:
				s.appendCommandText(e.Text)
			case modeSearch:
				if s.skipNextSearchEdit {
					s.skipNextSearchEdit = false
					continue
				}
				s.appendSearchText(e.Text)
			case modeFuzzyFinder:
				if s.skipNextFuzzyEdit {
					s.skipNextFuzzyEdit = false
					continue
				}
				s.appendFuzzyInput(e.Text)
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
			lineContent := expandTabs(s.activeBuffer().Line(index), 4)
			lineText := fmt.Sprintf("%4d  %s", index+1, lineContent)
			label := material.Body1(s.theme, lineText)
			label.Font.Typeface = "GoMono"
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
		cur := s.activeBuffer().Cursor()

		// Build status line with file info
		fileName := s.activeBuffer().FilePath()
		if fileName == "" {
			fileName = "[No Name]"
		}

		modFlag := ""
		if s.activeBuffer().Modified() {
			modFlag = " [+]"
		}

		bufferInfo := ""
		if s.bufferMgr.BufferCount() > 1 {
			bufferInfo = fmt.Sprintf(" | BUFFER %d/%d", s.bufferMgr.ActiveIndex()+1, s.bufferMgr.BufferCount())
		}

		// Add fullscreen indicator
		fullscreenInfo := ""
		if s.currentWindowMode == app.Fullscreen {
			fullscreenInfo = " | FULLSCREEN"
		}

		status = fmt.Sprintf("MODE %s | FILE %s%s | CURSOR %d:%d%s%s | %s",
			s.mode, fileName, modFlag, cur.Line+1, cur.Col+1, bufferInfo, fullscreenInfo, s.status,
		)
	}

	label := material.Body2(s.theme, status)
	label.Font.Typeface = "GoMono"
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
				label.Font.Typeface = "GoMono"
				label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
				return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, label.Layout)
			}),
			// Match count
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				matchInfo := fmt.Sprintf("%d matches", len(s.fuzzyFinderMatches))
				label := material.Body2(s.theme, matchInfo)
				label.Font.Typeface = "GoMono"
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
					label.Font.Typeface = "GoMono"
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
	label.Font.Typeface = "GoMono"
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
			pathLabel.Font.Typeface = "GoMono"
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

			list := layout.List{Axis: layout.Vertical}
			inset := layout.Inset{
				Top:    unit.Dp(8),
				Right:  unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(8),
			}

			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return list.Layout(gtx, len(nodes), func(gtx layout.Context, index int) layout.Dimensions {
					node := nodes[index]

					// Draw selection highlight
					if index == selectedIndex {
						rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(20)))}.Push(gtx.Ops)
						paint.Fill(gtx.Ops, selectedBg)
						rect.Pop()
					}

					// Build line text with indentation
					indent := strings.Repeat("  ", node.Depth)
					icon := "  "
					if node.IsDir {
						if node.Expanded {
							icon = "▼ "
						} else {
							icon = "▶ "
						}
					}
					lineText := indent + icon + node.Name

					// Render text
					label := material.Body2(s.theme, lineText)
					label.Font.Typeface = "GoMono"
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
	if ev.State != key.Press {
		return
	}
	s.lastKey = describeKey(ev)

	// Handle file operation input if active
	if s.fileOpMode != "" {
		if s.handleFileOpKey(ev) {
			return
		}
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
	if _, ok := s.printableKey(ev); ok {
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
		charWidth := s.measureTextWidth(gtx, charUnder)
		if charWidth < 8 {
			charWidth = 8
		}
		rect := image.Rect(x, 0, x+charWidth, height)
		stack := clip.Rect(rect).Push(gtx.Ops)
		paint.Fill(gtx.Ops, cursorColor)
		stack.Pop()

		// Draw the character on top of the cursor in contrasting color
		label := material.Body1(s.theme, charUnder)
		label.Font.Typeface = "GoMono"
		label.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		offset := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)
		label.Layout(gtx)
		offset.Pop()
	}
}

func (s *appState) measureTextWidth(gtx layout.Context, txt string) int {
	label := material.Body1(s.theme, txt)
	label.Font.Typeface = "GoMono"
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
		// Show explorer without focusing
		s.explorerVisible = true
		s.status = "Explorer opened"
	}
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
	if _, err := s.bufferMgr.OpenFile(node.Path); err != nil {
		s.status = fmt.Sprintf("Error opening %s: %v", node.Name, err)
	} else {
		s.exitExplorerMode()
		s.status = fmt.Sprintf("Opened %s", node.Name)
	}
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
	if err := s.bufferMgr.CloseActiveBuffer(force); err != nil {
		s.status = fmt.Sprintf("Error: %v", err)
		return
	}

	// If no buffers left, close the window
	if s.bufferMgr.BufferCount() == 0 || (s.bufferMgr.BufferCount() == 1 && s.activeBuffer().FilePath() == "") {
		s.requestClose()
	} else {
		s.status = "Buffer closed"
	}
}

func (s *appState) handleWriteCommand(arg string, andQuit bool) {
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

	filename := s.activeBuffer().FilePath()
	if andQuit {
		s.handleQuitCommand(false)
		s.status = fmt.Sprintf("Wrote + closed %s", filename)
		return
	}
	s.status = fmt.Sprintf("Wrote %d line(s) → %s", s.activeBuffer().LineCount(), filename)
}

func (s *appState) handleEditCommand(path string) {
	if path == "" {
		s.status = "E471: Argument required"
		return
	}

	if _, err := s.bufferMgr.OpenFile(path); err != nil {
		s.status = fmt.Sprintf("Error opening %s: %v", path, err)
	} else {
		s.status = fmt.Sprintf("Opened %s", path)
	}
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
	s.activeBuffer().InsertText(text)
	s.setCursorStatus(fmt.Sprintf("Insert %q", text))
	s.skipNextEdit = false
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

	if _, err := s.bufferMgr.OpenFile(fullPath); err != nil {
		s.status = fmt.Sprintf("Error opening %s: %v", match.FilePath, err)
		s.exitFuzzyFinder()
	} else {
		s.exitFuzzyFinder()
		s.status = fmt.Sprintf("Opened %s", match.FilePath)
	}
}

const sampleBuffer = ``
