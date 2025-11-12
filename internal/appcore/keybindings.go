package appcore

import (
	"log"
	"strings"
	"unicode"

	"gioui.org/io/key"
)

type Action int

const (
	ActionNone Action = iota

	// Global actions (work in any mode)
	ActionToggleExplorer
	ActionFocusExplorer
	ActionFocusEditor
	ActionToggleFullscreen

	// Mode transitions
	ActionEnterInsert
	ActionEnterVisualChar
	ActionEnterVisualLine
	ActionEnterDelete
	ActionEnterCommand
	ActionEnterExplorer
	ActionExitMode

	// Navigation
	ActionMoveLeft
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
	ActionJumpLineStart
	ActionJumpLineEnd
	ActionGotoLine
	ActionStartGotoSequence
	ActionWordForward
	ActionWordBackward
	ActionWordEnd

	// Editing
	ActionInsertNewline
	ActionInsertSpace
	ActionInsertTab
	ActionDeleteBackward
	ActionDeleteForward
	ActionDeleteLine
	ActionUndo

	// Visual mode
	ActionCopySelection
	ActionDeleteSelection
	ActionPasteClipboard

	// Clipboard (Normal mode)
	ActionCopyLine
	ActionPaste

	// Explorer
	ActionOpenNode
	ActionCollapseNode
	ActionExpandNode
	ActionRefreshTree
	ActionNavigateUp
	ActionRenameFile
	ActionDeleteFile
	ActionCreateFile

	// Search
	ActionEnterSearch
	ActionNextMatch
	ActionPrevMatch
	ActionClearSearch

	// Fuzzy Finder
	ActionOpenFuzzyFinder
	ActionFuzzyFinderConfirm

	// Buffer management
	ActionNextBuffer
	ActionPrevBuffer

	// Viewport scrolling
	ActionScrollToCenter
	ActionScrollToTop
	ActionScrollToBottom
	ActionScrollLineUp
	ActionScrollLineDown

	// Pane management
	ActionSplitVertical
	ActionSplitHorizontal
	ActionPaneFocusLeft
	ActionPaneFocusRight
	ActionPaneFocusUp
	ActionPaneFocusDown
	ActionPaneCycleNext
	ActionPaneClose
	ActionPaneEqualize
	ActionPaneZoomToggle

	// Terminal
	ActionOpenTerminal
	ActionTerminalExit
)

type KeyBinding struct {
	Modifiers key.Modifiers
	Key       key.Name
	Modes     []mode
	Action    Action
}

var globalKeybindings = []KeyBinding{
	{Modifiers: key.ModCtrl, Key: "t", Modes: nil, Action: ActionToggleExplorer},
	{Modifiers: key.ModCtrl, Key: "h", Modes: nil, Action: ActionFocusExplorer},
	{Modifiers: key.ModCtrl, Key: "l", Modes: nil, Action: ActionFocusEditor},
	{Modifiers: key.ModCtrl, Key: "f", Modes: nil, Action: ActionOpenFuzzyFinder},
	{Modifiers: key.ModCtrl, Key: "u", Modes: nil, Action: ActionUndo},
	{Modifiers: key.ModShift, Key: key.NameReturn, Modes: []mode{modeNormal}, Action: ActionToggleFullscreen},
	{Modifiers: key.ModShift, Key: key.NameEnter, Modes: []mode{modeNormal}, Action: ActionToggleFullscreen},

	// Clipboard operations
	{Modifiers: key.ModCtrl, Key: "c", Modes: []mode{modeNormal}, Action: ActionCopyLine},
	{Modifiers: key.ModCtrl, Key: "p", Modes: nil, Action: ActionPaste},

	// Pane navigation (Alt+hjkl)
	{Modifiers: key.ModAlt, Key: "h", Modes: nil, Action: ActionPaneFocusLeft},
	{Modifiers: key.ModAlt, Key: "j", Modes: nil, Action: ActionPaneFocusDown},
	{Modifiers: key.ModAlt, Key: "k", Modes: nil, Action: ActionPaneFocusUp},
	{Modifiers: key.ModAlt, Key: "l", Modes: nil, Action: ActionPaneFocusRight},

	// Pane closing
	{Modifiers: key.ModCtrl, Key: "x", Modes: nil, Action: ActionPaneClose},

	// Terminal - Ctrl+` opens/toggles terminal
	{Modifiers: key.ModCtrl, Key: "`", Modes: nil, Action: ActionOpenTerminal},
}

var modeKeybindings = map[mode][]KeyBinding{
	modeNormal: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameLeftArrow, Modes: nil, Action: ActionMoveLeft},
		{Modifiers: 0, Key: key.NameRightArrow, Modes: nil, Action: ActionMoveRight},
		{Modifiers: 0, Key: key.NameUpArrow, Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: key.NameDownArrow, Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: "i", Modes: nil, Action: ActionEnterInsert},
		{Modifiers: 0, Key: "v", Modes: nil, Action: ActionEnterVisualChar},
		{Modifiers: key.ModShift, Key: "v", Modes: nil, Action: ActionEnterVisualLine},
		{Modifiers: 0, Key: "d", Modes: nil, Action: ActionEnterDelete},
		{Modifiers: 0, Key: "h", Modes: nil, Action: ActionMoveLeft},
		{Modifiers: 0, Key: "j", Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: "k", Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: "l", Modes: nil, Action: ActionMoveRight},
		{Modifiers: 0, Key: "w", Modes: nil, Action: ActionWordForward},
		{Modifiers: 0, Key: "b", Modes: nil, Action: ActionWordBackward},
		{Modifiers: 0, Key: "e", Modes: nil, Action: ActionWordEnd},
		{Modifiers: 0, Key: "0", Modes: nil, Action: ActionJumpLineStart},
		{Modifiers: 0, Key: "$", Modes: nil, Action: ActionJumpLineEnd},
		{Modifiers: key.ModShift, Key: "4", Modes: nil, Action: ActionJumpLineEnd},
		{Modifiers: 0, Key: "/", Modes: nil, Action: ActionEnterSearch},
		{Modifiers: 0, Key: "n", Modes: nil, Action: ActionNextMatch},
		{Modifiers: key.ModShift, Key: "n", Modes: nil, Action: ActionPrevMatch},
		{Modifiers: key.ModCtrl, Key: "e", Modes: nil, Action: ActionScrollLineDown},
		{Modifiers: key.ModCtrl, Key: "y", Modes: nil, Action: ActionScrollLineUp},
		{Modifiers: key.ModShift, Key: key.NameTab, Modes: nil, Action: ActionPaneCycleNext},
	},
	modeInsert: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameReturn, Modes: nil, Action: ActionInsertNewline},
		{Modifiers: 0, Key: key.NameEnter, Modes: nil, Action: ActionInsertNewline},
		{Modifiers: 0, Key: key.NameSpace, Modes: nil, Action: ActionInsertSpace},
		{Modifiers: 0, Key: key.NameTab, Modes: nil, Action: ActionInsertTab},
		{Modifiers: 0, Key: key.NameDeleteBackward, Modes: nil, Action: ActionDeleteBackward},
		{Modifiers: 0, Key: key.NameDeleteForward, Modes: nil, Action: ActionDeleteForward},
		{Modifiers: 0, Key: key.NameLeftArrow, Modes: nil, Action: ActionMoveLeft},
		{Modifiers: 0, Key: key.NameRightArrow, Modes: nil, Action: ActionMoveRight},
		{Modifiers: 0, Key: key.NameUpArrow, Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: key.NameDownArrow, Modes: nil, Action: ActionMoveDown},
	},
	modeVisual: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameLeftArrow, Modes: nil, Action: ActionMoveLeft},
		{Modifiers: 0, Key: key.NameRightArrow, Modes: nil, Action: ActionMoveRight},
		{Modifiers: 0, Key: key.NameUpArrow, Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: key.NameDownArrow, Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: "h", Modes: nil, Action: ActionMoveLeft},
		{Modifiers: 0, Key: "j", Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: "k", Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: "l", Modes: nil, Action: ActionMoveRight},
		{Modifiers: 0, Key: "w", Modes: nil, Action: ActionWordForward},
		{Modifiers: 0, Key: "b", Modes: nil, Action: ActionWordBackward},
		{Modifiers: 0, Key: "e", Modes: nil, Action: ActionWordEnd},
		{Modifiers: 0, Key: "0", Modes: nil, Action: ActionJumpLineStart},
		{Modifiers: 0, Key: "$", Modes: nil, Action: ActionJumpLineEnd},
		{Modifiers: key.ModShift, Key: "4", Modes: nil, Action: ActionJumpLineEnd},
		{Modifiers: 0, Key: "c", Modes: nil, Action: ActionCopySelection},
		{Modifiers: 0, Key: "d", Modes: nil, Action: ActionDeleteSelection},
		{Modifiers: 0, Key: "p", Modes: nil, Action: ActionPasteClipboard},
		{Modifiers: 0, Key: "v", Modes: nil, Action: ActionExitMode},
		{Modifiers: key.ModShift, Key: key.NameTab, Modes: nil, Action: ActionPaneCycleNext},
	},
	modeDelete: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: key.ModShift, Key: key.NameTab, Modes: nil, Action: ActionPaneCycleNext},
	},
	modeCommand: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameReturn, Modes: nil, Action: ActionInsertNewline},
		{Modifiers: 0, Key: key.NameEnter, Modes: nil, Action: ActionInsertNewline},
		{Modifiers: 0, Key: key.NameDeleteBackward, Modes: nil, Action: ActionDeleteBackward},
	},
	modeExplorer: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameReturn, Modes: nil, Action: ActionOpenNode},
		{Modifiers: 0, Key: key.NameEnter, Modes: nil, Action: ActionOpenNode},
		{Modifiers: 0, Key: key.NameUpArrow, Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: key.NameDownArrow, Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: key.NameLeftArrow, Modes: nil, Action: ActionCollapseNode},
		{Modifiers: 0, Key: key.NameRightArrow, Modes: nil, Action: ActionExpandNode},
		{Modifiers: 0, Key: "j", Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: "k", Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: "h", Modes: nil, Action: ActionCollapseNode},
		{Modifiers: 0, Key: "l", Modes: nil, Action: ActionExpandNode},
		{Modifiers: 0, Key: "r", Modes: nil, Action: ActionRenameFile},
		{Modifiers: 0, Key: "d", Modes: nil, Action: ActionDeleteFile},
		{Modifiers: 0, Key: "n", Modes: nil, Action: ActionCreateFile},
		{Modifiers: 0, Key: "u", Modes: nil, Action: ActionNavigateUp},
		{Modifiers: 0, Key: "q", Modes: nil, Action: ActionExitMode},
		{Modifiers: key.ModShift, Key: key.NameTab, Modes: nil, Action: ActionPaneCycleNext},
	},
	modeSearch: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameReturn, Modes: nil, Action: ActionNextMatch},
		{Modifiers: 0, Key: key.NameEnter, Modes: nil, Action: ActionNextMatch},
		{Modifiers: 0, Key: key.NameDeleteBackward, Modes: nil, Action: ActionDeleteBackward},
	},
	modeFuzzyFinder: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameReturn, Modes: nil, Action: ActionFuzzyFinderConfirm},
		{Modifiers: 0, Key: key.NameEnter, Modes: nil, Action: ActionFuzzyFinderConfirm},
		{Modifiers: 0, Key: key.NameUpArrow, Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: key.NameDownArrow, Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: key.NameDeleteBackward, Modes: nil, Action: ActionDeleteBackward},
	},
	modeTerminal: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionTerminalExit},
		{Modifiers: key.ModShift, Key: key.NameTab, Modes: nil, Action: ActionTerminalExit},
	},
}

func (s *appState) matchGlobalKeybinding(ev key.Event) Action {
	for _, binding := range globalKeybindings {
		if !s.modifiersMatch(ev, binding.Modifiers) {
			continue
		}

		if s.keysMatch(ev.Name, binding.Key) {
			if len(binding.Modes) == 0 {
				return binding.Action
			}
			for _, m := range binding.Modes {
				if m == s.mode {
					return binding.Action
				}
			}
		}
	}

	return ActionNone
}

func (s *appState) matchModeKeybinding(m mode, ev key.Event) Action {
	bindings, exists := modeKeybindings[m]
	if !exists {
		return ActionNone
	}

	for _, binding := range bindings {
		if !s.modifiersMatch(ev, binding.Modifiers) {
			continue
		}

		if s.keysMatch(ev.Name, binding.Key) {
			return binding.Action
		}
	}

	return ActionNone
}

func (s *appState) keysMatch(actual, expected key.Name) bool {
	return strings.EqualFold(string(actual), string(expected))
}

func (s *appState) modifiersMatch(ev key.Event, required key.Modifiers) bool {
	// If no modifiers are required, ensure no modifiers are pressed
	if required == 0 {
		return ev.Modifiers == 0
	}

	// Build the actual modifiers state
	// PLATFORM QUIRK: ev.Modifiers is ALWAYS empty on some platforms!
	// We MUST rely on tracked state from explicit Press/Release events
	ctrlHeld := s.ctrlPressed                   // Trust tracked state, not ev.Modifiers
	shiftHeld := s.shiftPressed                 // Trust tracked state, not ev.Modifiers
	altHeld := ev.Modifiers.Contain(key.ModAlt) // Alt not tracked yet

	// Check if required modifiers are present
	ctrlRequired := required.Contain(key.ModCtrl)
	shiftRequired := required.Contain(key.ModShift)
	altRequired := required.Contain(key.ModAlt)

	// All required modifiers must be present
	if ctrlRequired && !ctrlHeld {
		return false
	}
	if shiftRequired && !shiftHeld {
		return false
	}
	if altRequired && !altHeld {
		return false
	}

	// No extra modifiers should be present (exact match)
	if !ctrlRequired && ctrlHeld {
		return false
	}
	if !shiftRequired && shiftHeld {
		return false
	}
	if !altRequired && altHeld {
		return false
	}

	return true
}

func (s *appState) matchPrintableKey(ev key.Event, target rune) bool {
	r, ok := s.printableKey(ev)
	if !ok {
		return false
	}
	return unicode.ToLower(r) == unicode.ToLower(target)
}

func (s *appState) executeAction(action Action, ev key.Event) {
	log.Printf("[ACTION] Executing action=%v mode=%s", action, s.mode)

	switch action {
	case ActionToggleExplorer:
		log.Printf("[TOGGLE_EXPLORER] Before: visible=%v focused=%v mode=%s",
			s.explorerVisible, s.explorerFocused, s.mode)
		if s.fileTree == nil {
			s.status = "File tree not available"
			log.Printf("[TOGGLE_EXPLORER] Aborted: file tree not available")
			return
		}
		s.toggleExplorer()
		log.Printf("[TOGGLE_EXPLORER] After: visible=%v focused=%v mode=%s",
			s.explorerVisible, s.explorerFocused, s.mode)

	case ActionFocusExplorer:
		if s.fileTree == nil {
			s.status = "File tree not available"
			return
		}
		if !s.explorerVisible {
			s.status = "Tree view hidden (Ctrl+T to open)"
			return
		}
		if s.mode == modeExplorer {
			s.status = "Tree view already focused (Ctrl+L to return to editor)"
			return
		}
		if s.mode != modeNormal {
			s.status = "Ctrl+H available from NORMAL mode"
			return
		}
		s.enterExplorerMode()
		s.status = "Focus: Tree View (use Ctrl+L to return to editor)"

	case ActionFocusEditor:
		if !s.explorerVisible {
			s.status = "Tree view hidden (Ctrl+T to open)"
			return
		}
		if s.mode != modeExplorer {
			s.status = "Editor already focused"
			return
		}
		s.exitExplorerMode()
		s.status = "Focus: Editor (use Ctrl+H to return to tree view)"

	case ActionToggleFullscreen:
		s.toggleFullscreen()

	case ActionEnterInsert:
		s.enterInsertMode()

	case ActionEnterVisualChar:
		s.enterVisualChar()

	case ActionEnterVisualLine:
		s.enterVisualLine()

	case ActionEnterDelete:
		s.enterDeleteMode()

	case ActionEnterCommand:
		s.enterCommandMode()

	case ActionEnterExplorer:
		log.Printf("[MODE_CHANGE] Entering EXPLORER mode from %s", s.mode)
		s.enterExplorerMode()
		log.Printf("[MODE_CHANGE] Now in mode=%s explorerFocused=%v", s.mode, s.explorerFocused)

	case ActionExitMode:
		log.Printf("[MODE_CHANGE] Exiting mode=%s", s.mode)
		oldMode := s.mode
		switch s.mode {
		case modeInsert:
			s.mode = modeNormal
			s.skipNextEdit = false
			s.resetCount()
			s.status = "Back to NORMAL"
		case modeVisual:
			s.exitVisualMode()
			s.resetCount()
			s.status = "Exited VISUAL"
		case modeDelete:
			s.exitDeleteMode()
		case modeCommand:
			s.exitCommandMode()
			s.status = "Command cancelled"
		case modeExplorer:
			s.exitExplorerMode()
		case modeSearch:
			s.exitSearchMode()
		case modeFuzzyFinder:
			s.exitFuzzyFinder()
		case modeNormal:
			s.exitVisualMode()
			s.resetCount()
			s.status = "Staying in NORMAL"
		}
		log.Printf("[MODE_CHANGE] Exited %s -> now in %s", oldMode, s.mode)

	case ActionMoveLeft:
		s.moveCursor("left")

	case ActionMoveRight:
		s.moveCursor("right")

	case ActionMoveUp:
		if s.mode == modeExplorer && s.fileTree != nil {
			if s.fileTree.MoveUp() {
				s.ensureExplorerItemVisible()
				s.status = "Explorer: moved up"
			}
		} else if s.mode == modeFuzzyFinder {
			s.fuzzyFinderMoveUp()
		} else {
			s.moveCursor("up")
		}

	case ActionMoveDown:
		if s.mode == modeExplorer && s.fileTree != nil {
			if s.fileTree.MoveDown() {
				s.ensureExplorerItemVisible()
				s.status = "Explorer: moved down"
			}
		} else if s.mode == modeFuzzyFinder {
			s.fuzzyFinderMoveDown()
		} else {
			s.moveCursor("down")
		}

	case ActionJumpLineStart:
		if s.activeBuffer().JumpLineStart() {
			s.setCursorStatus("Line start")
		} else {
			s.status = "Already at line start"
		}

	case ActionJumpLineEnd:
		if s.activeBuffer().JumpLineEnd() {
			s.setCursorStatus("Line end")
		} else {
			s.status = "Already at line end"
		}

	case ActionWordForward:
		if s.activeBuffer().MoveWordForward() {
			s.setCursorStatus("Word forward")
		} else {
			s.status = "End of buffer"
		}

	case ActionWordBackward:
		if s.activeBuffer().MoveWordBackward() {
			s.setCursorStatus("Word backward")
		} else {
			s.status = "Start of buffer"
		}

	case ActionWordEnd:
		if s.activeBuffer().MoveWordEnd() {
			s.setCursorStatus("Word end")
		} else {
			s.status = "End of buffer"
		}

	case ActionInsertNewline:
		if s.mode == modeInsert {
			buf := s.activeBuffer()
			cursorBefore := buf.Cursor()
			lineCountBefore := buf.LineCount()
			log.Printf("[NEWLINE_DEBUG] BEFORE: Line=%d Col=%d LineCount=%d",
				cursorBefore.Line, cursorBefore.Col, lineCountBefore)

			s.insertText("\n")
			s.skipNextEdit = true // Prevent EditEvent from inserting again

			cursorAfter := buf.Cursor()
			lineCountAfter := buf.LineCount()
			log.Printf("[NEWLINE_DEBUG] AFTER: Line=%d Col=%d LineCount=%d (added %d lines)",
				cursorAfter.Line, cursorAfter.Col, lineCountAfter, lineCountAfter-lineCountBefore)
		} else if s.mode == modeCommand {
			s.executeCommandLine()
		}

	case ActionInsertSpace:
		if s.mode == modeInsert {
			buf := s.activeBuffer()
			cursorBefore := buf.Cursor()
			lineBefore := buf.Line(cursorBefore.Line)
			log.Printf("[SPACE_DEBUG] BEFORE: Line=%d Col=%d LineContent=%q",
				cursorBefore.Line, cursorBefore.Col, lineBefore)

			s.insertText(" ")
			s.skipNextEdit = true // Prevent EditEvent from inserting again

			cursorAfter := buf.Cursor()
			lineAfter := buf.Line(cursorAfter.Line)
			log.Printf("[SPACE_DEBUG] AFTER: Line=%d Col=%d LineContent=%q",
				cursorAfter.Line, cursorAfter.Col, lineAfter)
			log.Printf("[SPACE_DEBUG] Cursor moved from Col %d to Col %d (delta: %d)",
				cursorBefore.Col, cursorAfter.Col, cursorAfter.Col-cursorBefore.Col)
		}

	case ActionInsertTab:
		if s.mode == modeInsert {
			buf := s.activeBuffer()
			cursorBefore := buf.Cursor()
			lineBefore := buf.Line(cursorBefore.Line)
			log.Printf("[TAB_DEBUG] BEFORE: Line=%d Col=%d LineContent=%q",
				cursorBefore.Line, cursorBefore.Col, lineBefore)

			s.insertText("\t")
			s.skipNextEdit = true // Prevent EditEvent from inserting again

			cursorAfter := buf.Cursor()
			lineAfter := buf.Line(cursorAfter.Line)
			log.Printf("[TAB_DEBUG] AFTER: Line=%d Col=%d LineContent=%q",
				cursorAfter.Line, cursorAfter.Col, lineAfter)
			log.Printf("[TAB_DEBUG] Tab character inserted, line now has tab character at correct position")
		}

	case ActionDeleteBackward:
		if s.mode == modeInsert {
			if s.activeBuffer().DeleteBackward() {
				s.setCursorStatus("Backspace")
			} else {
				s.status = "Start of buffer"
			}
		} else if s.mode == modeCommand {
			s.deleteCommandChar()
		} else if s.mode == modeSearch {
			s.deleteSearchChar()
		} else if s.mode == modeFuzzyFinder {
			s.deleteFuzzyChar()
		}

	case ActionDeleteForward:
		if s.mode == modeInsert {
			if s.activeBuffer().DeleteForward() {
				s.setCursorStatus("Delete")
			} else {
				s.status = "End of buffer"
			}
		}

	case ActionUndo:
		if s.activeBuffer().Undo() {
			s.status = "Undo successful"
		} else {
			s.status = "Nothing to undo"
		}

	case ActionCopySelection:
		s.copyVisualSelection()

	case ActionDeleteSelection:
		s.deleteVisualSelection()

	case ActionPasteClipboard:
		s.pasteClipboard()

	case ActionCopyLine:
		s.copyCurrentLine()

	case ActionPaste:
		s.pasteAtCursor()

	case ActionOpenNode:
		s.openSelectedNode()

	case ActionCollapseNode:
		if s.fileTree != nil {
			if s.fileTree.Collapse() {
				s.status = "Explorer: collapsed"
			}
		}

	case ActionExpandNode:
		if s.fileTree != nil {
			if s.fileTree.Expand() {
				if node := s.fileTree.SelectedNode(); node != nil && node.IsDir {
					s.fileTree.ExpandAndLoad(node)
				}
				s.status = "Explorer: expanded"
			}
		}

	case ActionRefreshTree:
		if s.fileTree != nil {
			if err := s.fileTree.Refresh(); err != nil {
				s.status = "Refresh error: " + err.Error()
			} else {
				s.status = "Tree refreshed"
			}
		}

	case ActionNavigateUp:
		if s.fileTree != nil {
			if err := s.fileTree.NavigateToParent(); err != nil {
				s.status = "Error navigating up: " + err.Error()
			} else {
				s.fileTree.LoadInitial()
				s.status = "Up to " + s.fileTree.CurrentPath()
			}
		}

	case ActionRenameFile:
		if s.mode == modeExplorer && s.fileTree != nil {
			s.enterRenameMode()
		}

	case ActionDeleteFile:
		if s.mode == modeExplorer && s.fileTree != nil {
			s.enterFileDeleteMode()
		}

	case ActionCreateFile:
		if s.mode == modeExplorer && s.fileTree != nil {
			s.enterCreateMode()
		}

	case ActionEnterSearch:
		s.enterSearchMode()

	case ActionNextMatch:
		if s.mode == modeSearch {
			s.executeSearch()
		} else {
			s.jumpToNextMatch()
		}

	case ActionPrevMatch:
		s.jumpToPrevMatch()

	case ActionClearSearch:
		s.clearSearch()

	case ActionOpenFuzzyFinder:
		s.enterFuzzyFinder()

	case ActionFuzzyFinderConfirm:
		s.fuzzyFinderConfirm()

	case ActionScrollToCenter:
		linesPerPage := 20
		s.scrollToCenter(linesPerPage)

	case ActionScrollToTop:
		s.scrollToTop()

	case ActionScrollToBottom:
		linesPerPage := 20
		s.scrollToBottom(linesPerPage)

	case ActionScrollLineUp:
		s.scrollLineUp()

	case ActionScrollLineDown:
		s.scrollLineDown()

	case ActionSplitVertical:
		s.handleSplitVertical()

	case ActionSplitHorizontal:
		s.handleSplitHorizontal()

	case ActionPaneFocusLeft:
		s.handlePaneFocusLeft()

	case ActionPaneFocusRight:
		s.handlePaneFocusRight()

	case ActionPaneFocusUp:
		s.handlePaneFocusUp()

	case ActionPaneFocusDown:
		s.handlePaneFocusDown()

	case ActionPaneCycleNext:
		s.handlePaneCycleNext()

	case ActionPaneClose:
		s.handlePaneClose()

	case ActionPaneEqualize:
		s.handlePaneEqualize()

	case ActionPaneZoomToggle:
		s.handlePaneZoomToggle()

	case ActionOpenTerminal:
		s.handleOpenTerminal()

	case ActionTerminalExit:
		s.handleTerminalExit()
	}
}
