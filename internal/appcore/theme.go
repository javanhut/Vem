package appcore

import (
	"image/color"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/javanhut/vem/internal/syntax"
)

// UIColors contains all UI colors derived from the active theme.
// These colors are computed from the Chroma syntax theme to provide
// a consistent visual appearance across the entire application.
type UIColors struct {
	// Core colors
	Background color.NRGBA
	Foreground color.NRGBA
	Accent     color.NRGBA

	// Selection and focus
	Selection color.NRGBA
	Focus     color.NRGBA
	FocusBg   color.NRGBA

	// Status bar
	StatusBar color.NRGBA
	StatusFg  color.NRGBA

	// Tabs
	TabActive   color.NRGBA
	TabInactive color.NRGBA
	TabText     color.NRGBA

	// File explorer
	Explorer          color.NRGBA
	ExplorerText      color.NRGBA
	ExplorerHighlight color.NRGBA

	// Modal dialogs
	Modal       color.NRGBA
	ModalBorder color.NRGBA
	ModalTitle  color.NRGBA

	// Diagnostics
	DiagnosticError color.NRGBA
	DiagnosticWarn  color.NRGBA
	DiagnosticInfo  color.NRGBA
	DiagnosticHint  color.NRGBA

	// Search
	SearchMatch color.NRGBA

	// Scrollbar
	Scrollbar color.NRGBA
}

// ColorsFromChromaStyle derives UI colors from a Chroma syntax highlighting style.
// This creates a unified look where syntax theme controls the entire UI appearance.
func ColorsFromChromaStyle(style *chroma.Style) UIColors {
	if style == nil {
		return GetDefaultColors()
	}

	// Extract base colors from Chroma style
	bg := syntax.GetBackgroundColor(style)
	fg := getForegroundColor(style)
	isDark := syntax.IsDarkTheme(style)

	// Compute accent color - use a contrasting blue/cyan
	accent := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}

	// Try to get keyword color as accent
	keywordEntry := style.Get(chroma.Keyword)
	if keywordEntry.Colour.IsSet() {
		accent = chromaColorToNRGBA(keywordEntry.Colour)
	}

	// Derive other colors based on dark/light theme
	var statusBar, statusFg color.NRGBA
	var tabActive, tabInactive, tabText color.NRGBA
	var selection, focus, focusBg color.NRGBA
	var explorer, explorerText, explorerHighlight color.NRGBA
	var modal, modalBorder, modalTitle color.NRGBA
	var scrollbar color.NRGBA

	if isDark {
		// Dark theme derivations
		statusBar = adjustBrightness(bg, 1.2)
		statusFg = fg
		tabActive = adjustBrightness(bg, 1.3)
		tabInactive = adjustBrightness(bg, 0.9)
		tabText = fg
		selection = withAlpha(accent, 0x44)
		focus = accent
		focusBg = withAlpha(accent, 0x22)
		explorer = adjustBrightness(bg, 0.9)
		explorerText = fg
		explorerHighlight = withAlpha(accent, 0x33)
		modal = adjustBrightness(bg, 1.1)
		modalBorder = adjustBrightness(bg, 1.5)
		modalTitle = fg
		scrollbar = withAlpha(fg, 0x44)
	} else {
		// Light theme derivations
		statusBar = adjustBrightness(bg, 0.9)
		statusFg = fg
		tabActive = adjustBrightness(bg, 0.95)
		tabInactive = adjustBrightness(bg, 1.05)
		tabText = fg
		selection = withAlpha(accent, 0x33)
		focus = accent
		focusBg = withAlpha(accent, 0x18)
		explorer = adjustBrightness(bg, 1.02)
		explorerText = fg
		explorerHighlight = withAlpha(accent, 0x22)
		modal = adjustBrightness(bg, 0.98)
		modalBorder = adjustBrightness(bg, 0.8)
		modalTitle = fg
		scrollbar = withAlpha(fg, 0x33)
	}

	return UIColors{
		Background: bg,
		Foreground: fg,
		Accent:     accent,

		Selection: selection,
		Focus:     focus,
		FocusBg:   focusBg,

		StatusBar: statusBar,
		StatusFg:  statusFg,

		TabActive:   tabActive,
		TabInactive: tabInactive,
		TabText:     tabText,

		Explorer:          explorer,
		ExplorerText:      explorerText,
		ExplorerHighlight: explorerHighlight,

		Modal:       modal,
		ModalBorder: modalBorder,
		ModalTitle:  modalTitle,

		// Fixed diagnostic colors for visibility
		DiagnosticError: color.NRGBA{R: 0xff, G: 0x55, B: 0x55, A: 0xff},
		DiagnosticWarn:  color.NRGBA{R: 0xff, G: 0xaa, B: 0x00, A: 0xff},
		DiagnosticInfo:  color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff},
		DiagnosticHint:  color.NRGBA{R: 0x8b, G: 0xe9, B: 0xfd, A: 0xff},

		SearchMatch: color.NRGBA{R: 0xff, G: 0xcc, B: 0x00, A: 0x88},

		Scrollbar: scrollbar,
	}
}

// GetDefaultColors returns the default dark theme colors.
// Used as fallback when no Chroma style is available.
func GetDefaultColors() UIColors {
	bg := color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}
	fg := color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
	accent := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}

	return UIColors{
		Background: bg,
		Foreground: fg,
		Accent:     accent,

		Selection: color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x88},
		Focus:     accent,
		FocusBg:   color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0x22},

		StatusBar: color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff},
		StatusFg:  fg,

		TabActive:   color.NRGBA{R: 0x28, G: 0x2e, B: 0x40, A: 0xff},
		TabInactive: color.NRGBA{R: 0x18, G: 0x1c, B: 0x28, A: 0xff},
		TabText:     fg,

		Explorer:          color.NRGBA{R: 0x16, G: 0x1a, B: 0x26, A: 0xff},
		ExplorerText:      fg,
		ExplorerHighlight: color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0x33},

		Modal:       color.NRGBA{R: 0x1e, G: 0x24, B: 0x34, A: 0xff},
		ModalBorder: color.NRGBA{R: 0x3a, G: 0x42, B: 0x58, A: 0xff},
		ModalTitle:  fg,

		DiagnosticError: color.NRGBA{R: 0xff, G: 0x55, B: 0x55, A: 0xff},
		DiagnosticWarn:  color.NRGBA{R: 0xff, G: 0xaa, B: 0x00, A: 0xff},
		DiagnosticInfo:  color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff},
		DiagnosticHint:  color.NRGBA{R: 0x8b, G: 0xe9, B: 0xfd, A: 0xff},

		SearchMatch: color.NRGBA{R: 0xff, G: 0xcc, B: 0x00, A: 0x88},

		Scrollbar: color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0x44},
	}
}

// getForegroundColor extracts the foreground/text color from a Chroma style.
func getForegroundColor(style *chroma.Style) color.NRGBA {
	if style == nil {
		return color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
	}

	// Try to get the Text token color
	textEntry := style.Get(chroma.Text)
	if textEntry.Colour.IsSet() {
		return chromaColorToNRGBA(textEntry.Colour)
	}

	// Try Name token
	nameEntry := style.Get(chroma.Name)
	if nameEntry.Colour.IsSet() {
		return chromaColorToNRGBA(nameEntry.Colour)
	}

	// Fallback to default text color
	return color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
}

// chromaColorToNRGBA converts a Chroma color to color.NRGBA.
func chromaColorToNRGBA(c chroma.Colour) color.NRGBA {
	rgb := uint32(c)
	return color.NRGBA{
		R: uint8((rgb >> 16) & 0xFF),
		G: uint8((rgb >> 8) & 0xFF),
		B: uint8(rgb & 0xFF),
		A: 0xff,
	}
}

// adjustBrightness adjusts the brightness of a color by a factor.
// factor > 1.0 lightens, factor < 1.0 darkens.
func adjustBrightness(c color.NRGBA, factor float64) color.NRGBA {
	r := float64(c.R) * factor
	g := float64(c.G) * factor
	b := float64(c.B) * factor

	// Clamp to 0-255
	if r > 255 {
		r = 255
	}
	if r < 0 {
		r = 0
	}
	if g > 255 {
		g = 255
	}
	if g < 0 {
		g = 0
	}
	if b > 255 {
		b = 255
	}
	if b < 0 {
		b = 0
	}

	return color.NRGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: c.A,
	}
}

// withAlpha returns the color with a new alpha value.
func withAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	return color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: alpha,
	}
}

// mapThemeName maps user-friendly theme names to Chroma style names.
// Returns the input unchanged if it's already a valid Chroma theme name.
func mapThemeName(themeName string) string {
	switch themeName {
	case "dark":
		return "monokai"
	case "light":
		return "solarized-light"
	default:
		return themeName
	}
}

// LoadChromaStyle loads a Chroma style by name, with fallback.
func LoadChromaStyle(themeName string) *chroma.Style {
	mappedName := mapThemeName(themeName)
	style := styles.Get(mappedName)
	if style == nil {
		// Fallback to monokai
		style = styles.Get("monokai")
	}
	return style
}
