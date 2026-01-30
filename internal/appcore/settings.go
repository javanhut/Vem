package appcore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Settings represents all user-configurable options for Vem.
// Settings are loaded from a TOML file and applied to appState on startup.
type Settings struct {
	Editor EditorSettings `toml:"editor"`
	UI     UISettings     `toml:"ui"`
	LSP    LSPSettings    `toml:"lsp"`
}

// EditorSettings controls text editing behavior.
type EditorSettings struct {
	TabWidth     int  `toml:"tab_width"`
	UseSpaces    bool `toml:"use_spaces"`
	AutoIndent   bool `toml:"auto_indent"`
	AutoPairs    bool `toml:"auto_pairs"`
	ScrollOffset int  `toml:"scroll_offset"`
	SoftWrap     bool `toml:"soft_wrap"`
	FormatOnSave bool `toml:"format_on_save"`
}

// UISettings controls visual appearance.
type UISettings struct {
	FontSize        int    `toml:"font_size"`
	Theme           string `toml:"theme"`
	ShowLineNumbers bool   `toml:"show_line_numbers"`
	RelativeNumbers bool   `toml:"relative_line_numbers"`
	ExplorerWidth   int    `toml:"explorer_width"`
}

// LSPSettings controls Language Server Protocol integration.
type LSPSettings struct {
	Enabled    bool `toml:"enabled"`
	AutoDetect bool `toml:"auto_detect"`
}

// NewSettings returns a Settings instance with sensible defaults.
func NewSettings() *Settings {
	return &Settings{
		Editor: EditorSettings{
			TabWidth:     4,
			UseSpaces:    false,
			AutoIndent:   true,
			AutoPairs:    true,
			ScrollOffset: 5,
			SoftWrap:     true,
			FormatOnSave: false,
		},
		UI: UISettings{
			FontSize:        14,
			Theme:           "dark",
			ShowLineNumbers: true,
			RelativeNumbers: false,
			ExplorerWidth:   250,
		},
		LSP: LSPSettings{
			Enabled:    true,
			AutoDetect: true,
		},
	}
}

// ConfigDir returns the platform-appropriate config directory for Vem.
// Uses os.UserConfigDir() for cross-platform support:
//   - Linux: ~/.config/vem (or $XDG_CONFIG_HOME/vem)
//   - macOS: ~/Library/Application Support/vem
//   - Windows: %APPDATA%\vem
func ConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting user config dir: %w", err)
	}
	return filepath.Join(base, "vem"), nil
}

// ConfigPath returns the full path to the settings.toml file.
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.toml"), nil
}

// Load reads settings from the config file, merging over defaults.
// If the config file doesn't exist, Load returns nil (using defaults).
// Other errors (permission denied, malformed TOML) are returned.
func (s *Settings) Load() error {
	path, err := ConfigPath()
	if err != nil {
		return fmt.Errorf("getting config path: %w", err)
	}

	_, err = toml.DecodeFile(path, s)
	if os.IsNotExist(err) {
		return nil // Use defaults silently
	}
	if err != nil {
		return fmt.Errorf("decoding settings file: %w", err)
	}
	return nil
}

// Save writes the current settings to the config file.
// Creates the config directory if it doesn't exist.
func (s *Settings) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return fmt.Errorf("getting config path: %w", err)
	}

	// Ensure config directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating settings file: %w", err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(s); err != nil {
		return fmt.Errorf("encoding settings: %w", err)
	}
	return nil
}

// applySettings syncs Settings values to appState runtime fields.
// Call after Load() and after settings modal saves.
func (s *appState) applySettings() {
	if s.settings == nil {
		return
	}

	// Editor settings
	// Note: tabWidth is a package constant - will need refactoring in Phase 2
	s.autoIndentEnabled = s.settings.Editor.AutoIndent
	s.autoPairsEnabled = s.settings.Editor.AutoPairs
	s.scrollOffsetLines = s.settings.Editor.ScrollOffset
	s.softWrapEnabled = s.settings.Editor.SoftWrap
	s.formatOnSave = s.settings.Editor.FormatOnSave

	// Set indent string based on UseSpaces setting
	if s.settings.Editor.UseSpaces {
		// Create spaces string based on TabWidth
		spaces := ""
		for i := 0; i < s.settings.Editor.TabWidth; i++ {
			spaces += " "
		}
		s.indentString = spaces
	} else {
		s.indentString = "\t"
	}

	// UI settings
	s.explorerWidth = s.settings.UI.ExplorerWidth

	// Theme settings
	s.loadTheme(s.settings.UI.Theme)

	// LSP settings
	s.lspEnabled = s.settings.LSP.Enabled
	s.lspAutoEnabled = s.settings.LSP.AutoDetect
}

// openSettingsModal opens the settings modal.
func (s *appState) openSettingsModal() {
	s.initSettingsModal()
	s.settingsModal.active = true
	s.status = "Settings"
}

// settingsPath returns the config file path for display purposes.
func (s *appState) settingsPath() string {
	path, err := ConfigPath()
	if err != nil {
		return "unknown"
	}
	return path
}
