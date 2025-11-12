package terminal

import "image/color"

// ANSI 16-color palette
var ansiColors = [16]color.NRGBA{
	// Normal colors (0-7)
	{R: 0x00, G: 0x00, B: 0x00, A: 0xff}, // Black
	{R: 0xcc, G: 0x00, B: 0x00, A: 0xff}, // Red
	{R: 0x4e, G: 0x9a, B: 0x06, A: 0xff}, // Green
	{R: 0xc4, G: 0xa0, B: 0x00, A: 0xff}, // Yellow
	{R: 0x34, G: 0x65, B: 0xa4, A: 0xff}, // Blue
	{R: 0x75, G: 0x50, B: 0x7b, A: 0xff}, // Magenta
	{R: 0x06, G: 0x98, B: 0x9a, A: 0xff}, // Cyan
	{R: 0xd3, G: 0xd7, B: 0xcf, A: 0xff}, // White

	// Bright colors (8-15)
	{R: 0x55, G: 0x57, B: 0x53, A: 0xff}, // Bright Black
	{R: 0xef, G: 0x29, B: 0x29, A: 0xff}, // Bright Red
	{R: 0x8a, G: 0xe2, B: 0x34, A: 0xff}, // Bright Green
	{R: 0xfc, G: 0xe9, B: 0x4f, A: 0xff}, // Bright Yellow
	{R: 0x72, G: 0x9f, B: 0xcf, A: 0xff}, // Bright Blue
	{R: 0xad, G: 0x7f, B: 0xa8, A: 0xff}, // Bright Magenta
	{R: 0x34, G: 0xe2, B: 0xe2, A: 0xff}, // Bright Cyan
	{R: 0xee, G: 0xee, B: 0xec, A: 0xff}, // Bright White
}

// GetANSIColor returns color for ANSI color code
func GetANSIColor(code int) color.NRGBA {
	if code >= 0 && code < 16 {
		return ansiColors[code]
	}

	// 256-color mode (codes 16-255)
	if code >= 16 && code < 232 {
		// 216-color cube (6x6x6)
		code -= 16
		r := (code / 36) * 51
		g := ((code % 36) / 6) * 51
		b := (code % 6) * 51
		return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}
	}

	// Grayscale (codes 232-255)
	if code >= 232 {
		gray := uint8((code-232)*10 + 8)
		return color.NRGBA{R: gray, G: gray, B: gray, A: 0xff}
	}

	// Default to white
	return ansiColors[7]
}

// Default foreground/background
var (
	DefaultFG = ansiColors[7]                                   // White
	DefaultBG = color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff} // Match Vem bg
)

// vt10xColorToNRGBA converts vt10x.Color to color.NRGBA
func vt10xColorToNRGBA(vtColor uint32) color.NRGBA {
	// vt10x uses special values for default colors
	const (
		DefaultFGValue = 1 << 24       // DefaultFG in vt10x
		DefaultBGValue = (1 << 24) + 1 // DefaultBG in vt10x
	)

	// Check for default colors
	if vtColor == DefaultFGValue {
		return DefaultFG
	}
	if vtColor == DefaultBGValue {
		return DefaultBG
	}

	// Check if it's a standard ANSI color (0-255)
	if vtColor < 256 {
		return GetANSIColor(int(vtColor))
	}

	// Check if it's a true color (RGB)
	// vt10x stores RGB as: (1<<25) | (r<<16) | (g<<8) | b
	if vtColor&(1<<25) != 0 {
		r := uint8((vtColor >> 16) & 0xFF)
		g := uint8((vtColor >> 8) & 0xFF)
		b := uint8(vtColor & 0xFF)
		return color.NRGBA{R: r, G: g, B: b, A: 0xff}
	}

	// Fallback to default
	return DefaultFG
}
