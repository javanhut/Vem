package appcore

import (
	"fmt"
	"image"
	"image/color"
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

const (
	modeNormal   mode = "NORMAL"
	modeInsert   mode = "INSERT"
	modeVisual   mode = "VISUAL"
	modeDelete   mode = "DELETE"
	modeCommand  mode = "COMMAND"
	modeExplorer mode = "EXPLORER"
)

const caretBlinkInterval = 600 * time.Millisecond

var (
	highlightColor = color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x55}
	selectionColor = color.NRGBA{R: 0x1c, G: 0x39, B: 0x60, A: 0x99}
	background     = color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}
	statusBg       = color.NRGBA{R: 0x12, G: 0x17, B: 0x22, A: 0xff}
	headerColor    = color.NRGBA{R: 0xa1, G: 0xc6, B: 0xff, A: 0xff}
	cursorColor    = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	focusBorder    = color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
)

type appState struct {
	theme        *material.Theme
	bufferMgr    *editor.BufferManager
	fileTree     *filesystem.FileTree
	mode         mode
	status       string
	lastKey      string
	focusTag     *int
	pendingCount int
	pendingGoto  bool
	visualStart  int
	visualActive bool
	skipNextEdit bool
	caretVisible bool
	nextBlink    time.Time
	caretReset   bool
	clipLines    []string
	cmdText      string
	window       *app.Window

	// Explorer state
	explorerVisible bool
	explorerWidth   int
	explorerFocused bool

	// Modifier tracking (some platforms don't report Ctrl in key modifiers)
	ctrlPressed bool

	// Fullscreen state tracking
	currentWindowMode app.WindowMode
	wasFullscreen     bool
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
		visualStart:       -1,
		visualActive:      false,
		caretVisible:      true,
		explorerVisible:   false,
		explorerWidth:     250,
		explorerFocused:   false,
		currentWindowMode: app.Windowed,
		wasFullscreen:     false,
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

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
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
}

func (s *appState) handleEvents(gtx layout.Context) {
	event.Op(gtx.Ops, s.focusTag)
	if s.mode == modeInsert || s.mode == modeCommand {
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
			if e.Name == key.NameCtrl {
				s.ctrlPressed = (e.State == key.Press)
				continue
			}
			s.handleKey(e)
		case key.EditEvent:
			if e.Text == "" {
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
	list := layout.List{Axis: layout.Vertical}
	cursorLine := s.activeBuffer().Cursor().Line
	selStart, selEnd, hasSel := s.visualSelectionRange()
	cursorCol := s.activeBuffer().Cursor().Col

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
		return list.Layout(gtx, lines, func(gtx layout.Context, index int) layout.Dimensions {
			lineContent := expandTabs(s.activeBuffer().Line(index), 4)
			lineText := fmt.Sprintf("%4d  %s", index+1, lineContent)
			label := material.Body1(s.theme, lineText)
			label.Font.Typeface = "GoMono"
			label.Color = color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
			macro := op.Record(gtx.Ops)
			dims := label.Layout(gtx)
			call := macro.Stop()

			if hasSel && index >= selStart && index <= selEnd {
				rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}.Push(gtx.Ops)
				paint.Fill(gtx.Ops, selectionColor)
				rect.Pop()
			}
			if index == cursorLine {
				rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}.Push(gtx.Ops)
				paint.Fill(gtx.Ops, highlightColor)
				rect.Pop()
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

func (s *appState) drawStatusBar(gtx layout.Context) layout.Dimensions {
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

	status := fmt.Sprintf("MODE %s | FILE %s%s | CURSOR %d:%d%s%s | %s",
		s.mode, fileName, modFlag, cur.Line+1, cur.Col+1, bufferInfo, fullscreenInfo, s.status,
	)
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

	if s.handleGlobalShortcuts(ev) {
		return
	}

	switch s.mode {
	case modeNormal:
		if s.handleNormalMode(ev) {
			return
		}
	case modeInsert:
		if s.handleInsertMode(ev) {
			return
		}
	case modeDelete:
		if s.handleDeleteMode(ev) {
			return
		}
	case modeVisual:
		if s.handleVisualMode(ev) {
			return
		}
	case modeCommand:
		if s.handleCommandMode(ev) {
			return
		}
	case modeExplorer:
		if s.handleExplorerMode(ev) {
			return
		}
	}
	s.status = "Waiting for motion"
}

func (s *appState) handleGlobalShortcuts(ev key.Event) bool {
	ctrlActive := ev.Modifiers.Contain(key.ModCtrl) || s.ctrlPressed
	if !ctrlActive {
		return false
	}

	name := strings.ToLower(string(ev.Name))

	switch name {
	case "t":
		return s.handleCtrlToggleExplorerShortcut()
	case "h":
		return s.handleCtrlFocusExplorerShortcut()
	case "l":
		return s.handleCtrlFocusEditorShortcut()
	}

	// Some platforms send DeleteBackward for Ctrl+H
	if ev.Name == key.NameDeleteBackward {
		return s.handleCtrlFocusExplorerShortcut()
	}

	return false
}

func (s *appState) handleCtrlToggleExplorerShortcut() bool {
	if s.fileTree == nil {
		s.status = "File tree not available"
		return true
	}
	s.toggleExplorer()
	return true
}

func (s *appState) handleCtrlFocusExplorerShortcut() bool {
	if s.fileTree == nil {
		s.status = "File tree not available"
		return true
	}
	if s.mode == modeExplorer {
		s.status = "Tree view already focused (Ctrl+L to return to editor)"
		return true
	}
	if s.mode != modeNormal {
		s.status = "Ctrl+H available from NORMAL mode"
		return true
	}
	s.enterExplorerMode()
	s.status = "Focus: Tree View (use Ctrl+L to return to editor)"
	return true
}

func (s *appState) handleCtrlFocusEditorShortcut() bool {
	if !s.explorerVisible {
		s.status = "Tree view hidden (Ctrl+T to open)"
		return true
	}
	if s.mode != modeExplorer {
		s.status = "Editor already focused"
		return true
	}
	s.exitExplorerMode()
	s.status = "Focus: Editor (use Ctrl+H to return to tree view)"
	return true
}

func (s *appState) handleNormalMode(ev key.Event) bool {
	// Check for Shift+Enter to toggle fullscreen
	if (ev.Name == key.NameReturn || ev.Name == key.NameEnter) && ev.Modifiers.Contain(key.ModShift) {
		s.toggleFullscreen()
		return true
	}

	if s.isColonKey(ev) {
		s.enterCommandMode()
		return true
	}

	switch ev.Name {
	case key.NameEscape:
		s.exitVisualMode()
		s.resetCount()
		s.status = "Staying in NORMAL"
		return true
	case key.NameLeftArrow:
		s.moveCursor("left")
		return true
	case key.NameRightArrow:
		s.moveCursor("right")
		return true
	case key.NameDownArrow:
		s.moveCursor("down")
		return true
	case key.NameUpArrow:
		s.moveCursor("up")
		return true
	}
	if r, ok := printableKey(ev); ok {
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
		switch r {
		case 'G':
			s.gotoLineWithCount()
			return true
		case 'g':
			s.startGotoSequence()
			return true
		}
		switch unicode.ToLower(r) {
		case 'i':
			s.enterInsertMode()
			return true
		case 't':
			s.enterExplorerMode()
			return true
		case 'v':
			s.enterVisualLine()
			return true
		case 'd':
			s.enterDeleteMode()
			return true
		case 'h':
			s.moveCursor("left")
			return true
		case 'j':
			s.moveCursor("down")
			return true
		case 'k':
			s.moveCursor("up")
			return true
		case 'l':
			s.moveCursor("right")
			return true
		case '0':
			if s.activeBuffer().JumpLineStart() {
				s.setCursorStatus("Line start")
			} else {
				s.status = "Already at line start"
			}
			return true
		case '$':
			if s.activeBuffer().JumpLineEnd() {
				s.setCursorStatus("Line end")
			} else {
				s.status = "Already at line end"
			}
			return true
		}
	}
	return false
}

func (s *appState) handleInsertMode(ev key.Event) bool {
	// Check for Shift+Enter to toggle fullscreen (before regular Enter handling)
	if (ev.Name == key.NameReturn || ev.Name == key.NameEnter) && ev.Modifiers.Contain(key.ModShift) {
		s.toggleFullscreen()
		return true
	}

	switch ev.Name {
	case key.NameEscape:
		s.mode = modeNormal
		s.skipNextEdit = false
		s.resetCount()
		s.status = "Back to NORMAL"
		return true
	case key.NameReturn, key.NameEnter:
		s.insertText("\n")
		return true
	case key.NameSpace:
		s.insertText(" ")
		return true
	case key.NameDeleteBackward:
		if s.activeBuffer().DeleteBackward() {
			s.setCursorStatus("Backspace")
		} else {
			s.status = "Start of buffer"
		}
		return true
	case key.NameDeleteForward:
		if s.activeBuffer().DeleteForward() {
			s.setCursorStatus("Delete")
		} else {
			s.status = "End of buffer"
		}
		return true
	case key.NameLeftArrow:
		s.moveCursor("left")
		return true
	case key.NameRightArrow:
		s.moveCursor("right")
		return true
	case key.NameUpArrow:
		s.moveCursor("up")
		return true
	case key.NameDownArrow:
		s.moveCursor("down")
		return true
	}
	if _, ok := printableKey(ev); ok {
		// Text insertion is driven by key.EditEvent to avoid double characters.
		return true
	}
	return false
}

func (s *appState) handleDeleteMode(ev key.Event) bool {
	// Check for Shift+Enter to toggle fullscreen
	if (ev.Name == key.NameReturn || ev.Name == key.NameEnter) && ev.Modifiers.Contain(key.ModShift) {
		s.toggleFullscreen()
		return true
	}

	switch ev.Name {
	case key.NameEscape:
		s.exitDeleteMode()
		return true
	}
	if r, ok := printableKey(ev); ok {
		if unicode.IsDigit(r) {
			s.handleCountDigit(int(r - '0'))
			if s.pendingCount > 0 {
				s.status = fmt.Sprintf("DELETE line %d (pending)", s.pendingCount)
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

func (s *appState) handleVisualMode(ev key.Event) bool {
	// Check for Shift+Enter to toggle fullscreen
	if (ev.Name == key.NameReturn || ev.Name == key.NameEnter) && ev.Modifiers.Contain(key.ModShift) {
		s.toggleFullscreen()
		return true
	}

	if s.isColonKey(ev) {
		s.enterCommandMode()
		return true
	}
	switch ev.Name {
	case key.NameEscape:
		s.exitVisualMode()
		s.resetCount()
		s.status = "Exited VISUAL"
		return true
	case key.NameLeftArrow:
		s.moveCursor("left")
		return true
	case key.NameRightArrow:
		s.moveCursor("right")
		return true
	case key.NameDownArrow:
		s.moveCursor("down")
		return true
	case key.NameUpArrow:
		s.moveCursor("up")
		return true
	}
	if r, ok := printableKey(ev); ok {
		if unicode.IsDigit(r) && s.handleCountDigit(int(r-'0')) {
			return true
		}
		if s.pendingGoto {
			if s.handleGotoSequence(r) {
				return true
			}
			s.pendingGoto = false
		}
		switch r {
		case 'G':
			s.gotoLineWithCount()
			return true
		case 'g':
			s.startGotoSequence()
			return true
		}
		switch unicode.ToLower(r) {
		case 'c':
			s.copyVisualSelection()
			return true
		case 'd':
			s.deleteVisualSelection()
			return true
		case 'p':
			s.pasteClipboard()
			return true
		case 'v':
			s.exitVisualMode()
			return true
		case 'h':
			s.moveCursor("left")
			return true
		case 'j':
			s.moveCursor("down")
			return true
		case 'k':
			s.moveCursor("up")
			return true
		case 'l':
			s.moveCursor("right")
			return true
		case '0':
			if s.activeBuffer().JumpLineStart() {
				s.setCursorStatus("Line start")
			}
			return true
		case '$':
			if s.activeBuffer().JumpLineEnd() {
				s.setCursorStatus("Line end")
			}
			return true
		}
	}
	return false
}

func (s *appState) handleCommandMode(ev key.Event) bool {
	// Check for Shift+Enter to toggle fullscreen (before regular Enter handling)
	if (ev.Name == key.NameReturn || ev.Name == key.NameEnter) && ev.Modifiers.Contain(key.ModShift) {
		s.toggleFullscreen()
		return true
	}

	switch ev.Name {
	case key.NameEscape:
		s.exitCommandMode()
		s.status = "Command cancelled"
		return true
	case key.NameReturn, key.NameEnter:
		s.executeCommandLine()
		return true
	case key.NameDeleteBackward:
		s.deleteCommandChar()
		return true
	}
	return false
}

func (s *appState) handleExplorerMode(ev key.Event) bool {
	if s.fileTree == nil {
		return false
	}

	// Check for Shift+Enter to toggle fullscreen (before regular Enter handling)
	if (ev.Name == key.NameReturn || ev.Name == key.NameEnter) && ev.Modifiers.Contain(key.ModShift) {
		s.toggleFullscreen()
		return true
	}

	switch ev.Name {
	case key.NameEscape:
		s.exitExplorerMode()
		return true
	case key.NameReturn, key.NameEnter:
		s.openSelectedNode()
		return true
	case key.NameUpArrow:
		if s.fileTree.MoveUp() {
			s.status = "Explorer: moved up"
		}
		return true
	case key.NameDownArrow:
		if s.fileTree.MoveDown() {
			s.status = "Explorer: moved down"
		}
		return true
	case key.NameLeftArrow:
		if s.fileTree.Collapse() {
			s.status = "Explorer: collapsed"
		}
		return true
	case key.NameRightArrow:
		if s.fileTree.Expand() {
			if node := s.fileTree.SelectedNode(); node != nil && node.IsDir {
				s.fileTree.ExpandAndLoad(node)
			}
			s.status = "Explorer: expanded"
		}
		return true
	}

	if r, ok := printableKey(ev); ok {
		switch unicode.ToLower(r) {
		case 'j':
			if s.fileTree.MoveDown() {
				s.status = "Explorer: moved down"
			}
			return true
		case 'k':
			if s.fileTree.MoveUp() {
				s.status = "Explorer: moved up"
			}
			return true
		case 'h':
			if s.fileTree.Collapse() {
				s.status = "Explorer: collapsed"
			}
			return true
		case 'l':
			if s.fileTree.Expand() {
				if node := s.fileTree.SelectedNode(); node != nil && node.IsDir {
					s.fileTree.ExpandAndLoad(node)
				}
				s.status = "Explorer: expanded"
			}
			return true
		case 'r':
			if err := s.fileTree.Refresh(); err != nil {
				s.status = fmt.Sprintf("Refresh error: %v", err)
			} else {
				s.status = "Tree refreshed"
			}
			return true
		case 'u':
			// Navigate to parent directory
			if err := s.fileTree.NavigateToParent(); err != nil {
				s.status = fmt.Sprintf("Error navigating up: %v", err)
			} else {
				s.fileTree.LoadInitial()
				s.status = fmt.Sprintf("Up to %s", s.fileTree.CurrentPath())
			}
			return true
		case 'q':
			s.exitExplorerMode()
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

func (s *appState) enterVisualLine() {
	s.mode = modeVisual
	s.visualActive = true
	s.visualStart = s.activeBuffer().Cursor().Line
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
	s.visualActive = false
	s.visualStart = -1
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
		// Show and enter explorer
		s.explorerVisible = true
		s.enterExplorerMode()
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
	if !s.visualActive || s.visualStart < 0 {
		return 0, 0, false
	}
	cur := s.activeBuffer().Cursor().Line
	start := s.visualStart
	if cur < start {
		return cur, start, true
	}
	return start, cur, true
}

func (s *appState) deleteVisualSelection() {
	start, end, ok := s.visualSelectionRange()
	if !ok {
		s.status = "No selection"
		return
	}
	s.activeBuffer().DeleteLines(start, end)
	s.exitVisualMode()
	s.setCursorStatus("Deleted selection")
}

func (s *appState) copyVisualSelection() {
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
}

func (s *appState) pasteClipboard() {
	if len(s.clipLines) == 0 {
		s.status = "Clipboard empty"
		return
	}
	start, _, ok := s.visualSelectionRange()
	if !ok {
		s.status = "Select destination in VISUAL mode"
		return
	}
	lines := append([]string(nil), s.clipLines...)
	s.activeBuffer().InsertLines(start, lines)
	s.exitVisualMode()
	s.setCursorStatus(fmt.Sprintf("Inserted %d line(s)", len(lines)))
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

func printableKey(ev key.Event) (rune, bool) {
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

const sampleBuffer = `
Goal: Build a NeoVim inspired text editor that runs anywhere.

Constraints:
  - Written in Go with a modern GPU UI (Gio spike).
  - Ships fonts + dependencies, zero extra setup.
  - Extensible with Lua, Python, Carrion scripting.
  - Needs Vim motions, macros, and LSP plumbing.

Next steps after spike:
  1. Solidify architecture doc.
  2. Validate rendering + text metrics decisions.
  3. Expand buffer representation for edits.
  4. Prove plugin host boundary + config flow.
`
