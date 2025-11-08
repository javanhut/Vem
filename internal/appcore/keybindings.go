package appcore

import (
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
	ActionEnterVisual
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

	// Editing
	ActionInsertNewline
	ActionInsertSpace
	ActionDeleteBackward
	ActionDeleteForward
	ActionDeleteLine

	// Visual mode
	ActionCopySelection
	ActionDeleteSelection
	ActionPasteClipboard

	// Explorer
	ActionOpenNode
	ActionCollapseNode
	ActionExpandNode
	ActionRefreshTree
	ActionNavigateUp

	// Buffer management
	ActionNextBuffer
	ActionPrevBuffer
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
	{Modifiers: key.ModShift, Key: key.NameReturn, Modes: nil, Action: ActionToggleFullscreen},
	{Modifiers: key.ModShift, Key: key.NameEnter, Modes: nil, Action: ActionToggleFullscreen},
}

var modeKeybindings = map[mode][]KeyBinding{
	modeNormal: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameLeftArrow, Modes: nil, Action: ActionMoveLeft},
		{Modifiers: 0, Key: key.NameRightArrow, Modes: nil, Action: ActionMoveRight},
		{Modifiers: 0, Key: key.NameUpArrow, Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: key.NameDownArrow, Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: "i", Modes: nil, Action: ActionEnterInsert},
		{Modifiers: 0, Key: "t", Modes: nil, Action: ActionEnterExplorer},
		{Modifiers: 0, Key: "v", Modes: nil, Action: ActionEnterVisual},
		{Modifiers: 0, Key: "d", Modes: nil, Action: ActionEnterDelete},
		{Modifiers: 0, Key: "h", Modes: nil, Action: ActionMoveLeft},
		{Modifiers: 0, Key: "j", Modes: nil, Action: ActionMoveDown},
		{Modifiers: 0, Key: "k", Modes: nil, Action: ActionMoveUp},
		{Modifiers: 0, Key: "l", Modes: nil, Action: ActionMoveRight},
		{Modifiers: 0, Key: "0", Modes: nil, Action: ActionJumpLineStart},
		{Modifiers: 0, Key: "$", Modes: nil, Action: ActionJumpLineEnd},
	},
	modeInsert: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
		{Modifiers: 0, Key: key.NameReturn, Modes: nil, Action: ActionInsertNewline},
		{Modifiers: 0, Key: key.NameEnter, Modes: nil, Action: ActionInsertNewline},
		{Modifiers: 0, Key: key.NameSpace, Modes: nil, Action: ActionInsertSpace},
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
		{Modifiers: 0, Key: "0", Modes: nil, Action: ActionJumpLineStart},
		{Modifiers: 0, Key: "$", Modes: nil, Action: ActionJumpLineEnd},
		{Modifiers: 0, Key: "c", Modes: nil, Action: ActionCopySelection},
		{Modifiers: 0, Key: "d", Modes: nil, Action: ActionDeleteSelection},
		{Modifiers: 0, Key: "p", Modes: nil, Action: ActionPasteClipboard},
		{Modifiers: 0, Key: "v", Modes: nil, Action: ActionExitMode},
	},
	modeDelete: {
		{Modifiers: 0, Key: key.NameEscape, Modes: nil, Action: ActionExitMode},
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
		{Modifiers: 0, Key: "r", Modes: nil, Action: ActionRefreshTree},
		{Modifiers: 0, Key: "u", Modes: nil, Action: ActionNavigateUp},
		{Modifiers: 0, Key: "q", Modes: nil, Action: ActionExitMode},
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
	if required == 0 {
		return ev.Modifiers == 0
	}

	ctrlHeld := ev.Modifiers.Contain(key.ModCtrl) || s.ctrlPressed

	if required.Contain(key.ModCtrl) && !ctrlHeld {
		return false
	}

	if required.Contain(key.ModShift) && !ev.Modifiers.Contain(key.ModShift) {
		return false
	}

	if required.Contain(key.ModAlt) && !ev.Modifiers.Contain(key.ModAlt) {
		return false
	}

	return true
}

func (s *appState) matchPrintableKey(ev key.Event, target rune) bool {
	r, ok := printableKey(ev)
	if !ok {
		return false
	}
	return unicode.ToLower(r) == unicode.ToLower(target)
}

func (s *appState) executeAction(action Action, ev key.Event) {
	switch action {
	case ActionToggleExplorer:
		if s.fileTree == nil {
			s.status = "File tree not available"
			return
		}
		s.toggleExplorer()

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

	case ActionEnterVisual:
		s.enterVisualLine()

	case ActionEnterDelete:
		s.enterDeleteMode()

	case ActionEnterCommand:
		s.enterCommandMode()

	case ActionEnterExplorer:
		s.enterExplorerMode()

	case ActionExitMode:
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
		case modeNormal:
			s.exitVisualMode()
			s.resetCount()
			s.status = "Staying in NORMAL"
		}

	case ActionMoveLeft:
		s.moveCursor("left")

	case ActionMoveRight:
		s.moveCursor("right")

	case ActionMoveUp:
		if s.mode == modeExplorer && s.fileTree != nil {
			if s.fileTree.MoveUp() {
				s.status = "Explorer: moved up"
			}
		} else {
			s.moveCursor("up")
		}

	case ActionMoveDown:
		if s.mode == modeExplorer && s.fileTree != nil {
			if s.fileTree.MoveDown() {
				s.status = "Explorer: moved down"
			}
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

	case ActionInsertNewline:
		if s.mode == modeInsert {
			s.insertText("\n")
		} else if s.mode == modeCommand {
			s.executeCommandLine()
		}

	case ActionInsertSpace:
		if s.mode == modeInsert {
			s.insertText(" ")
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
		}

	case ActionDeleteForward:
		if s.mode == modeInsert {
			if s.activeBuffer().DeleteForward() {
				s.setCursorStatus("Delete")
			} else {
				s.status = "End of buffer"
			}
		}

	case ActionCopySelection:
		s.copyVisualSelection()

	case ActionDeleteSelection:
		s.deleteVisualSelection()

	case ActionPasteClipboard:
		s.pasteClipboard()

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
	}
}
