package syntax

import (
	"image/color"

	"github.com/alecthomas/chroma/v2"
)

// GetTokenColor returns the color for a given token type from the style.
func GetTokenColor(tokenType chroma.TokenType, style *chroma.Style) color.NRGBA {
	if style == nil {
		// Default text color
		return color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
	}

	// Get the style entry for this token type
	entry := style.Get(tokenType)

	// Check if a color is set
	if entry.Colour.IsSet() {
		return chromaColorToNRGBA(entry.Colour)
	}

	// Try parent token types if no color is set
	// Chroma uses a hierarchy: Keyword.Declaration -> Keyword -> Token
	parentType := tokenType.Parent()
	if parentType != tokenType && parentType != chroma.None {
		return GetTokenColor(parentType, style)
	}

	// Fallback to default text color
	return color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
}

// chromaColorToNRGBA converts a Chroma color to color.NRGBA.
func chromaColorToNRGBA(c chroma.Colour) color.NRGBA {
	// Chroma colors are stored as integers (0xRRGGBB format)
	rgb := uint32(c)
	return color.NRGBA{
		R: uint8((rgb >> 16) & 0xFF),
		G: uint8((rgb >> 8) & 0xFF),
		B: uint8(rgb & 0xFF),
		A: 0xff,
	}
}

// GetBackgroundColor returns the background color from the style.
func GetBackgroundColor(style *chroma.Style) color.NRGBA {
	if style == nil || !style.Has(chroma.Background) {
		// Default background
		return color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}
	}

	bg := style.Get(chroma.Background)
	if bg.Background.IsSet() {
		return chromaColorToNRGBA(bg.Background)
	}

	// Fallback
	return color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}
}

// PresetThemes defines recommended themes that work well with Vem.
var PresetThemes = []string{
	"monokai",
	"dracula",
	"github-dark",
	"nord",
	"one-dark",
	"solarized-dark",
	"solarized-light",
	"vim",
	"catppuccin-mocha",
	"gruvbox",
}

// GetThemeDescription returns a human-readable description for known themes.
func GetThemeDescription(themeName string) string {
	descriptions := map[string]string{
		"monokai":          "Dark theme with vibrant colors (default)",
		"dracula":          "Dark theme with purple accents",
		"github-dark":      "GitHub's dark theme",
		"nord":             "Arctic-inspired dark theme",
		"one-dark":         "Atom One Dark theme",
		"solarized-dark":   "Precision colors for machines and people",
		"solarized-light":  "Light variant of Solarized",
		"vim":              "Classic Vim color scheme",
		"catppuccin-mocha": "Warm, pastel dark theme",
		"gruvbox":          "Retro groove dark theme",
	}

	if desc, ok := descriptions[themeName]; ok {
		return desc
	}
	return "Color theme"
}

// IsDarkTheme returns true if the theme is dark (for contrast adjustments).
func IsDarkTheme(style *chroma.Style) bool {
	if style == nil {
		return true // Default to dark
	}

	bg := GetBackgroundColor(style)

	// Calculate perceived brightness using the formula:
	// brightness = (0.299*R + 0.587*G + 0.114*B)
	brightness := (0.299*float64(bg.R) + 0.587*float64(bg.G) + 0.114*float64(bg.B))

	// If brightness is less than 128 (middle of 0-255), it's dark
	return brightness < 128
}
