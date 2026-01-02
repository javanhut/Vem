package appcore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/javanhut/vem/internal/panes"
)

type viewStateFile struct {
	Files map[string]bufferViewState `json:"files"`
}

func (s *appState) initViewState() {
	s.viewStatePath = defaultViewStatePath()
	s.loadViewState()
	if s.paneManager != nil {
		if pane := s.paneManager.ActivePane(); pane != nil {
			s.applyViewStateForPane(pane)
		}
	}
}

func defaultViewStatePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" {
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil || homeDir == "" {
			return ""
		}
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, "vem", "state.json")
}

func (s *appState) loadViewState() {
	if s.viewStatePath == "" {
		return
	}
	data, err := os.ReadFile(s.viewStatePath)
	if err != nil {
		return
	}
	var state viewStateFile
	if err := json.Unmarshal(data, &state); err != nil {
		return
	}
	if state.Files != nil {
		s.viewState = state.Files
	}
}

func (s *appState) saveViewStateNow() {
	if s.viewStatePath == "" {
		return
	}
	state := viewStateFile{Files: s.viewState}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(s.viewStatePath), 0755); err != nil {
		return
	}
	if err := os.WriteFile(s.viewStatePath, data, 0644); err != nil {
		return
	}
	s.viewStateDirty = false
	s.viewStateLastSaved = time.Now()
}

func (s *appState) maybeSaveViewState() {
	if !s.viewStateDirty {
		return
	}
	if time.Since(s.viewStateLastSaved) < 8*time.Second {
		return
	}
	s.saveViewStateNow()
}

func (s *appState) updateViewStateForPane(pane *panes.Pane) {
	if pane == nil || s.bufferMgr == nil {
		return
	}
	buf := s.bufferMgr.GetBuffer(pane.BufferIndex)
	if buf == nil {
		return
	}
	path := buf.FilePath()
	if path == "" {
		return
	}
	cur := buf.Cursor()
	state := bufferViewState{
		Line:        cur.Line,
		Col:         cur.Col,
		ViewportTop: pane.ViewportTop,
	}
	if prev, ok := s.viewState[path]; ok && prev == state {
		return
	}
	s.viewState[path] = state
	s.viewStateDirty = true
}

func (s *appState) applyViewStateForPane(pane *panes.Pane) {
	if pane == nil || s.bufferMgr == nil {
		return
	}
	buf := s.bufferMgr.GetBuffer(pane.BufferIndex)
	if buf == nil {
		return
	}
	path := buf.FilePath()
	if path == "" {
		return
	}
	state, ok := s.viewState[path]
	if !ok {
		return
	}
	buf.SetCursor(state.Line, state.Col)
	pane.SetViewportTop(state.ViewportTop)
}

func (s *appState) recordBufferAccess(index int) {
	if index < 0 {
		return
	}
	if s.bufferAccessTimes == nil {
		s.bufferAccessTimes = make(map[int]time.Time)
	}
	s.bufferAccessTimes[index] = time.Now()
}

func (s *appState) switchPaneBuffer(pane *panes.Pane, index int) {
	if pane == nil || s.bufferMgr == nil {
		return
	}
	s.updateViewStateForPane(pane)
	s.bufferMgr.SwitchToBuffer(index)
	pane.SetBufferIndex(index)
	s.applyViewStateForPane(pane)
	s.recordBufferAccess(index)
}
